package uniq

import (
	"strings"
	"testing"
)

func TestGenerateDefaults(t *testing.T) {
	got, err := Generate(20)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 20 {
		t.Fatalf("length = %d, want 20", len(got))
	}
	for _, c := range got {
		if c < 'a' || c > 'z' {
			t.Errorf("default alphabet must be lowercase letters, got %q", c)
		}
	}
}

func TestGenerateWithUppercase(t *testing.T) {
	tests := []struct {
		name      string
		length    int
		upper     int
		wantUpper int
	}{
		{"some", 20, 5, 5},
		{"none", 20, 0, 0},
		{"all", 10, 10, 10},
		{"exceeds length", 5, 50, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MustGenerate(tt.length, WithUppercase(tt.upper))
			if len(got) != tt.length {
				t.Errorf("length = %d, want %d (%q)", len(got), tt.length, got)
			}
			if n := countUpper(got); n != tt.wantUpper {
				t.Errorf("uppercase = %d, want %d (%q)", n, tt.wantUpper, got)
			}
		})
	}
}

func TestGenerateWithDigits(t *testing.T) {
	got := MustGenerate(2000, WithDigits())
	hasDigit := false
	for _, c := range got {
		isLower := c >= 'a' && c <= 'z'
		isDigit := c >= '0' && c <= '9'
		if isDigit {
			hasDigit = true
		}
		if !isLower && !isDigit {
			t.Fatalf("unexpected character %q", c)
		}
	}
	if !hasDigit {
		t.Error("expected at least one digit in a 2000-char string")
	}
}

func TestGenerateWithAlphabet(t *testing.T) {
	const alphabet = "abc123"
	got := MustGenerate(500, WithAlphabet(alphabet))
	for _, c := range got {
		if !strings.ContainsRune(alphabet, c) {
			t.Fatalf("character %q not in custom alphabet", c)
		}
	}
}

func TestGenerateCombinedOptions(t *testing.T) {
	got := MustGenerate(30, WithDigits(), WithUppercase(10))
	if n := countUpper(got); n != 10 {
		t.Errorf("uppercase = %d, want 10 (%q)", n, got)
	}
}

func TestGenerateErrors(t *testing.T) {
	if _, err := Generate(-1); err == nil {
		t.Error("negative length: expected error")
	}
	if _, err := Generate(10, WithAlphabet("")); err == nil {
		t.Error("empty alphabet: expected error")
	}
	if got, err := Generate(0); err != nil || got != "" {
		t.Errorf("zero length: got %q, %v; want empty, nil", got, err)
	}
}

func TestGenerateNoCollisions(t *testing.T) {
	const iterations = 10000
	seen := make(map[string]struct{}, iterations)
	for i := 0; i < iterations; i++ {
		id := MustGenerate(20, WithDigits(), WithUppercase(5))
		if _, dup := seen[id]; dup {
			t.Fatalf("collision after %d iterations: %q", i, id)
		}
		seen[id] = struct{}{}
	}
}

// TestGenerateUppercasePositionDistribution guards against bias in the
// partial shuffle: with one uppercase letter in a length-10 string, each
// position should be hit close to 1/10 of the time.
func TestGenerateUppercasePositionDistribution(t *testing.T) {
	const samples = 20000
	const length = 10
	counts := make([]int, length)
	for i := 0; i < samples; i++ {
		id := MustGenerate(length, WithUppercase(1))
		for p := 0; p < length; p++ {
			if id[p] >= 'A' && id[p] <= 'Z' {
				counts[p]++
			}
		}
	}
	expected := float64(samples) / length
	for p, got := range counts {
		if float64(got) < expected*0.85 || float64(got) > expected*1.15 {
			t.Errorf("position %d hit %d times, expected ~%.0f (±15%%)", p, got, expected)
		}
	}
}

func BenchmarkGenerate10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustGenerate(10)
	}
}

func BenchmarkGenerate20Digits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustGenerate(20, WithDigits())
	}
}

func BenchmarkGenerate20Upper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustGenerate(20, WithUppercase(6))
	}
}

func BenchmarkGenerate32All(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustGenerate(32, WithDigits(), WithUppercase(8))
	}
}

func BenchmarkGenerateParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			MustGenerate(20, WithDigits())
		}
	})
}
