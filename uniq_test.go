package uniq

import (
	"testing"
)

func countUpper(s string) int {
	n := 0
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			n++
		}
	}
	return n
}

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		length        int
		numUpperCase  int
		wantLen       int
		wantUpperCase int
	}{
		{"basic", 10, 3, 10, 3},
		{"medium", 15, 5, 15, 5},
		{"long", 50, 20, 50, 20},
		{"no uppercase", 10, 0, 10, 0},
		{"all uppercase", 10, 10, 10, 10},
		{"uppercase exceeds length", 5, 100, 5, 5},
		{"negative uppercase", 10, -3, 10, 0},
		{"zero length", 0, 5, 0, 0},
		{"negative length", -1, 5, 0, 0},
		{"length one", 1, 1, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.length, tt.numUpperCase)

			if len(result) != tt.wantLen {
				t.Errorf("New(%d, %d) length = %d, want %d (%q)",
					tt.length, tt.numUpperCase, len(result), tt.wantLen, result)
			}
			if got := countUpper(result); got != tt.wantUpperCase {
				t.Errorf("New(%d, %d) uppercase count = %d, want %d (%q)",
					tt.length, tt.numUpperCase, got, tt.wantUpperCase, result)
			}
		})
	}
}

func TestNewOnlyLetters(t *testing.T) {
	result := New(1000, 300)
	for i, c := range result {
		isLower := c >= 'a' && c <= 'z'
		isUpper := c >= 'A' && c <= 'Z'
		if !isLower && !isUpper {
			t.Fatalf("unexpected character %q at position %d", c, i)
		}
	}
}

// TestNewNoCollisions verifies the fix for per-call reseeding: rapid
// successive calls must not produce identical IDs.
func TestNewNoCollisions(t *testing.T) {
	const iterations = 10000
	seen := make(map[string]struct{}, iterations)
	for i := 0; i < iterations; i++ {
		id := New(20, 5)
		if _, dup := seen[id]; dup {
			t.Fatalf("collision after %d iterations: %q", i, id)
		}
		seen[id] = struct{}{}
	}
}

func TestNewConcurrent(t *testing.T) {
	const goroutines = 50
	done := make(chan string, goroutines)
	for i := 0; i < goroutines; i++ {
		go func() { done <- New(20, 5) }()
	}
	seen := make(map[string]struct{}, goroutines)
	for i := 0; i < goroutines; i++ {
		id := <-done
		if len(id) != 20 {
			t.Errorf("length = %d, want 20", len(id))
		}
		if _, dup := seen[id]; dup {
			t.Errorf("concurrent collision: %q", id)
		}
		seen[id] = struct{}{}
	}
}

func BenchmarkNew10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(10, 3)
	}
}

func BenchmarkNew20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(20, 6)
	}
}

func BenchmarkNew50(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(50, 15)
	}
}

func BenchmarkNew100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(100, 30)
	}
}

func BenchmarkNewParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			New(20, 6)
		}
	})
}
