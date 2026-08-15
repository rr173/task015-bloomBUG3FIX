package bloom

import (
	"strconv"
	"sync"
	"testing"
)

// TestProbeConcurrentAddCount verifies that concurrent Add calls produce
// an exact final count equal to the number of calls, ensuring no lost updates.
func TestProbeConcurrentAddCount(t *testing.T) {
	f, err := New(10000, 0.01)
	if err != nil {
		t.Fatal(err)
	}
	const goroutines = 8
	const perG = 500
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < perG; i++ {
				f.Add("g" + strconv.Itoa(id) + "-" + strconv.Itoa(i))
			}
		}(g)
	}
	wg.Wait()
	s := f.Stats()
	want := uint64(goroutines * perG)
	if s.Count != want {
		t.Errorf("concurrent add count = %d, want %d (lost updates)", s.Count, want)
	}
}

// TestProbeOptimalKRounding verifies that OptimalK uses rounding (not
// truncation) so that k=round(6.64)=7 for n=1000,p=0.01 (m=9586).
func TestProbeOptimalKRounding(t *testing.T) {
	// m=9586, n=1000: k_exact = (9586/1000)*ln2 ≈ 6.643
	// Correct: round(6.643) = 7
	k := OptimalK(9586, 1000)
	if k != 7 {
		t.Errorf("OptimalK(9586,1000) = %d, want 7 (should round, not floor)", k)
	}
}
