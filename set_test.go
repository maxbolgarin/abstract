package abstract_test

import (
	"sync"
	"testing"

	"github.com/maxbolgarin/abstract"
)

// TestNewSet tests creating a new Set and adding elements.
func TestNewSet(t *testing.T) {
	s := abstract.NewSet[int]()
	if !s.IsEmpty() {
		t.Error("New set should be empty")
	}

	// Test addition of elements
	s.Add(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("Expected set length to be 3, got %d", s.Len())
	}

	// Test Has method
	if !s.Has(1) || !s.Has(2) || !s.Has(3) {
		t.Error("Set should contain elements 1, 2, and 3")
	}
	if s.Has(4) {
		t.Error("Set should not contain element 4")
	}

	// Test removal of elements
	s.Delete(2)
	if s.Len() != 2 {
		t.Errorf("Expected set length to be 2, got %d", s.Len())
	}

	a := s.Values()
	if len(a) != 2 {
		t.Errorf("Expected set length to be 2, got %d", len(a))
	}
}

// TestSetClear tests clearing the Set.
func TestSetClear(t *testing.T) {
	s := abstract.NewSetFromItems(1, 2, 3)
	s.Clear()

	if !s.IsEmpty() {
		t.Error("Set should be empty after clear")
	}
	if s.Len() != 0 {
		t.Errorf("Set length should be 0 after clear, got %d", s.Len())
	}
}

// TestSetTransform tests transforming the Set.
func TestSetTransform(t *testing.T) {
	s := abstract.NewSet([]int{1, 2, 3})
	s.Transform(func(k int) int {
		return k * 2
	})

	if !s.Has(2) || !s.Has(4) || !s.Has(6) {
		t.Error("Set should transform its elements correctly")
	}
	if s.Has(1) || s.Has(3) {
		t.Error("Set should not have old values after transform")
	}
}

// TestSetRange tests iterating over the Set.
func TestSetRange(t *testing.T) {
	s := abstract.NewSetWithSize[int](3)
	s.Add(1, 2, 3)
	if s.Range(func(k int) bool {
		if k == 3 {
			return false
		}
		if k != 1 && k != 2 && k != 3 {
			t.Errorf("Set should iterate over all elements, got %d", k)
		}
		return true
	}) {
		t.Error("Expected Range to return false, but got true")
	}

	if !s.Range(func(k int) bool {
		return true
	}) {
		t.Error("Expected Range to return true, but got false")
	}
}

// TestNewSafeSet tests creating a SafeSet and concurrent use.
func TestNewSafeSet(t *testing.T) {
	s := abstract.NewSafeSet[int]()
	if !s.IsEmpty() {
		t.Error("New SafeSet should be empty")
	}

	// Test addition of elements
	s.Add(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("SafeSet length should be 3, got %d", s.Len())
	}

	// Test Has method
	if !s.Has(1) || !s.Has(2) || !s.Has(3) {
		t.Error("SafeSet should contain elements 1, 2, and 3")
	}
	if s.Has(4) {
		t.Error("SafeSet should not contain element 4")
	}

	s.Delete(2)
	if s.Len() != 2 {
		t.Errorf("SafeSet length should be 2, got %d", s.Len())
	}

	a := s.Values()
	if len(a) != 2 {
		t.Errorf("SafeSet length should be 2, got %d", len(a))
	}
}

// TestSafeSetConcurrency tests thread safety of SafeSet.
func TestSafeSetConcurrency(t *testing.T) {
	s := abstract.NewSafeSet[int]()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Add(i)
		}(i)
	}

	wg.Wait()

	if s.Len() != 100 {
		t.Errorf("SafeSet length should be 100, got %d", s.Len())
	}
}

// TestSafeSetRange tests iterating over the SafeSet.
func TestSafeSetRange(t *testing.T) {
	s := abstract.NewSafeSetWithSize[int](3)
	s.Add(1, 2, 3)
	if s.Range(func(k int) bool {
		if k == 3 {
			return false
		}
		if k != 1 && k != 2 && k != 3 {
			t.Errorf("SafeSet should iterate over all elements, got %d", k)
		}
		return true
	}) {
		t.Error("Expected Range to return false, but got true")
	}

	if !s.Range(func(k int) bool {
		return true
	}) {
		t.Error("Expected Range to return true, but got false")
	}
}

// TestSafeSetClear tests clearing the SafeSet.
func TestSafeSetClear(t *testing.T) {
	s := abstract.NewSafeSet[int]([]int{1, 2, 3})
	s.Clear()

	if !s.IsEmpty() {
		t.Error("SafeSet should be empty after clear")
	}
	if s.Len() != 0 {
		t.Errorf("SafeSet length should be 0 after clear, got %d", s.Len())
	}
}

// TestSafeSetTransform tests transforming the SafeSet.
func TestSafeSetTransform(t *testing.T) {
	s := abstract.NewSafeSet([]int{1, 2, 3})
	s.Transform(func(k int) int {
		return k * 2
	})

	if !s.Has(2) || !s.Has(4) || !s.Has(6) {
		t.Error("SafeSet should transform its elements correctly")
	}
	if s.Has(1) || s.Has(3) {
		t.Error("SafeSet should not have old values after transform")
	}

	s = abstract.NewSafeSetFromItems(1, 2, 3)
	s.Transform(func(k int) int {
		return k * 2
	})

	if !s.Has(2) || !s.Has(4) || !s.Has(6) {
		t.Error("SafeSet should transform its elements correctly")
	}
	if s.Has(1) || s.Has(3) {
		t.Error("SafeSet should not have old values after transform")
	}
}
