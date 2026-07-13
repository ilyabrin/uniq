package uniq

import (
	"strings"
	"testing"
)

func TestToken(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantLen int
	}{
		{"short", 8, 8},
		{"typical", 32, 32},
		{"long", 128, 128},
		{"length one", 1, 1},
		{"zero length", 0, 0},
		{"negative length", -5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Token(tt.length)
			if err != nil {
				t.Fatalf("Token(%d) error: %v", tt.length, err)
			}
			if len(got) != tt.wantLen {
				t.Errorf("Token(%d) length = %d, want %d (%q)", tt.length, len(got), tt.wantLen, got)
			}
			for i, c := range got {
				if !strings.ContainsRune(tokenAlphabet, c) {
					t.Errorf("Token(%d) contains %q at %d, not in alphabet", tt.length, c, i)
				}
			}
		})
	}
}

func TestTokenNoCollisions(t *testing.T) {
	const iterations = 10000
	seen := make(map[string]struct{}, iterations)
	for i := 0; i < iterations; i++ {
		id := MustToken(32)
		if _, dup := seen[id]; dup {
			t.Fatalf("collision after %d iterations: %q", i, id)
		}
		seen[id] = struct{}{}
	}
}

// TestTokenDistribution sanity-checks uniformity: over a large sample,
// each alphabet character should appear within a loose tolerance of the
// expected frequency (guards against modulo bias).
func TestTokenDistribution(t *testing.T) {
	const samples = 200
	const tokenLen = 500
	counts := make(map[rune]int)
	for i := 0; i < samples; i++ {
		for _, c := range MustToken(tokenLen) {
			counts[c]++
		}
	}
	total := samples * tokenLen
	expected := float64(total) / float64(len(tokenAlphabet))
	for _, c := range tokenAlphabet {
		got := float64(counts[c])
		if got < expected*0.8 || got > expected*1.2 {
			t.Errorf("char %q count %v deviates >20%% from expected %v", c, got, expected)
		}
	}
}

func BenchmarkToken16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustToken(16)
	}
}

func BenchmarkToken32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustToken(32)
	}
}

func BenchmarkToken64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustToken(64)
	}
}

func BenchmarkTokenParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			MustToken(32)
		}
	})
}
