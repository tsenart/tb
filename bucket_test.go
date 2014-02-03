package tb

import (
	"testing"
	"time"
)

func TestNewBucket(t *testing.T) {
	t.Parallel()

	if b := NewBucket(10); b.capacity != 10 || b.tokens != 10 {
		t.Errorf("Wrong capacity or tokens. Want 10, Got [%d, %d]",
			b.capacity, b.tokens)
	}
}

func TestTake(t *testing.T) {
	t.Parallel()

	b := NewBucket(10)
	defer b.Stop()

	ex := [...]int64{5, 5, 1, 1, 5, 4, 1, 0}
	for i := 0; i < len(ex)-1; i += 2 {
		if got, want := b.Take(ex[i]), ex[i+1]; got != want {
			t.Errorf("Want: %d, Got: %d", want, got)
		}
	}

	time.Sleep(1 * time.Second) // Wait for bucket to fill up

	exs := [2][]int64{{4, 4, 2, 2, 1, 1}, {2, 2, 1, 1, 1, 0}}
	for i := 0; i < 2; i++ {
		go func(i int) {
			for j := 0; j < len(exs[i])-1; j += 2 {
				if got, want := b.Take(exs[i][j]), exs[i][j+1]; got != want {
					t.Errorf("Want: %d, Got: %d", want, got)
				}
			}
		}(i)
	}
}
