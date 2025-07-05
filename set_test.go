package abstract_test

import (
	"sync"
	"testing"

	"github.com/maxbolgarin/abstract"
)

// TestNewSet tests creating a new Set and adding elements.
func TestNewSet(t *testing.T) {
	var s = &abstract.Set[int]{}
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

// TestSetCopy tests copying the Set.
func TestSetCopy(t *testing.T) {
	s := abstract.NewSet([]int{1, 2, 3})
	copy := s.Copy()
	if len(copy) != 3 {
		t.Errorf("Expected set length to be 3, got %d", len(copy))
	}
}

func TestSetIter(t *testing.T) {
	s := abstract.NewSet([]int{1, 2, 3})
	iter := s.Iter()
	var counter int
	for range iter {
		counter++
	}
	if counter != 3 {
		t.Errorf("Expected 3, got %d", counter)
	}
}

// TestSetUnion tests the Union method of Set.
func TestSetUnion(t *testing.T) {
	set1 := abstract.NewSet([]int{1, 2, 3})
	set2 := abstract.NewSet([]int{4, 5, 6})

	union := set1.Union(set2.Copy())
	if union.Len() != 6 {
		t.Errorf("Expected union length to be 6, got %d", union.Len())
	}
}

// TestSetIntersection tests the Intersection method of Set.
func TestSetIntersection(t *testing.T) {
	set1 := abstract.NewSet([]int{1, 2, 3})
	set2 := abstract.NewSet([]int{4, 5, 6})

	intersection := set1.Intersection(set2.Copy())
	if intersection.Len() != 0 {
		t.Errorf("Expected intersection length to be 0, got %d", intersection.Len())
	}
}

// TestSetDifference tests the Difference method of Set.
func TestSetDifference(t *testing.T) {
	set1 := abstract.NewSet([]int{1, 2, 3})
	set2 := abstract.NewSet([]int{4, 5, 6})

	difference := set1.Difference(set2.Copy())
	if difference.Len() != 3 {
		t.Errorf("Expected difference length to be 3, got %d", difference.Len())
	}
}

// TestSetSymmetricDifference tests the SymmetricDifference method of Set.
func TestSetSymmetricDifference(t *testing.T) {
	set1 := abstract.NewSet([]int{1, 2, 3})
	set2 := abstract.NewSet([]int{4, 5, 6})

	symmetricDifference := set1.SymmetricDifference(set2.Copy())
	if symmetricDifference.Len() != 6 {
		t.Errorf("Expected symmetric difference length to be 6, got %d", symmetricDifference.Len())
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
	s := abstract.NewSafeSet([]int{1, 2, 3})
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

// TestSafeSetCopy tests copying the SafeSet.
func TestSafeSetCopy(t *testing.T) {
	s := abstract.NewSafeSet([]int{1, 2, 3})
	copy := s.Copy()
	if len(copy) != 3 {
		t.Errorf("Expected set length to be 3, got %d", len(copy))
	}
}

func TestSafeSetIter(t *testing.T) {
	s := abstract.NewSafeSet([]int{1, 2, 3})
	iter := s.Iter()
	var counter int
	for range iter {
		counter++
	}
	if counter != 3 {
		t.Errorf("Expected 3, got %d", counter)
	}
}

// TestSafeSetUnion tests the Union method of SafeSet.
func TestSafeSetUnion(t *testing.T) {
	set1 := abstract.NewSafeSet([]int{1, 2, 3})
	set2 := abstract.NewSafeSet([]int{4, 5, 6})

	union := set1.Union(set2.Copy())
	if union.Len() != 6 {
		t.Errorf("Expected union length to be 6, got %d", union.Len())
	}
}

// TestSafeSetIntersection tests the Intersection method of SafeSet.
func TestSafeSetIntersection(t *testing.T) {
	set1 := abstract.NewSafeSet([]int{1, 2, 3})
	set2 := abstract.NewSafeSet([]int{4, 5, 6})

	intersection := set1.Intersection(set2.Copy())
	if intersection.Len() != 0 {
		t.Errorf("Expected intersection length to be 0, got %d", intersection.Len())
	}
}

// TestSafeSetDifference tests the Difference method of SafeSet.
func TestSafeSetDifference(t *testing.T) {
	set1 := abstract.NewSafeSet([]int{1, 2, 3})
	set2 := abstract.NewSafeSet([]int{4, 5, 6})

	difference := set1.Difference(set2.Copy())
	if difference.Len() != 3 {
		t.Errorf("Expected difference length to be 3, got %d", difference.Len())
	}
}

// TestSafeSetSymmetricDifference tests the SymmetricDifference method of SafeSet.
func TestSafeSetSymmetricDifference(t *testing.T) {
	set1 := abstract.NewSafeSet([]int{1, 2, 3})
	set2 := abstract.NewSafeSet([]int{4, 5, 6})

	symmetricDifference := set1.SymmetricDifference(set2.Copy())
	if symmetricDifference.Len() != 6 {
		t.Errorf("Expected symmetric difference length to be 6, got %d", symmetricDifference.Len())
	}
}

// ===== UNINITIALIZED SET TESTS =====

func TestSet_UninitializedMethods(t *testing.T) {
	// Test Add with uninitialized set
	var s1 abstract.Set[int]
	s1.Add(1, 2, 3)
	if s1.Len() != 3 {
		t.Errorf("Expected length 3 after Add on uninitialized set, got %d", s1.Len())
	}

	// Test Has with uninitialized set
	var s2 abstract.Set[int]
	if s2.Has(1) {
		t.Error("Expected false for uninitialized set")
	}

	// Test Delete with uninitialized set
	var s3 abstract.Set[int]
	deleted := s3.Delete(1)
	if deleted {
		t.Error("Expected false from Delete on uninitialized set")
	}

	// Test Len with uninitialized set
	var s4 abstract.Set[int]
	if s4.Len() != 0 {
		t.Errorf("Expected length 0 for uninitialized set, got %d", s4.Len())
	}

	// Test IsEmpty with uninitialized set
	var s5 abstract.Set[int]
	if !s5.IsEmpty() {
		t.Error("Expected true from IsEmpty on uninitialized set")
	}

	// Test Values with uninitialized set
	var s6 abstract.Set[int]
	values := s6.Values()
	if len(values) != 0 {
		t.Errorf("Expected empty values slice, got length %d", len(values))
	}

	// Test Clear with uninitialized set
	var s7 abstract.Set[int]
	s7.Clear()
	if s7.Len() != 0 {
		t.Errorf("Expected length 0 after Clear on uninitialized set, got %d", s7.Len())
	}

	// Test Transform with uninitialized set
	var s8 abstract.Set[int]
	s8.Transform(func(k int) int { return k + 1 })
	if s8.Len() != 0 {
		t.Errorf("Expected no items after Transform on uninitialized set, got %d", s8.Len())
	}

	// Test Range with uninitialized set
	var s9 abstract.Set[int]
	called := false
	result := s9.Range(func(k int) bool {
		called = true
		return true
	})
	if !result || called {
		t.Error("Expected Range to return true without calling function on uninitialized set")
	}

	// Test Raw with uninitialized set
	var s10 abstract.Set[int]
	raw := s10.Raw()
	if len(raw) != 0 {
		t.Errorf("Expected empty raw map, got length %d", len(raw))
	}

	// Test Iter with uninitialized set
	var s11 abstract.Set[int]
	count := 0
	for range s11.Iter() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 iterations from Iter on uninitialized set, got %d", count)
	}

	// Test Copy with uninitialized set
	var s12 abstract.Set[int]
	copied := s12.Copy()
	if len(copied) != 0 {
		t.Errorf("Expected empty copied map, got length %d", len(copied))
	}

	// Test Union with uninitialized set
	var s13 abstract.Set[int]
	other := map[int]struct{}{1: {}, 2: {}}
	union := s13.Union(other)
	if union.Len() != 2 {
		t.Errorf("Expected union length 2, got %d", union.Len())
	}

	// Test Intersection with uninitialized set
	var s14 abstract.Set[int]
	intersection := s14.Intersection(other)
	if intersection.Len() != 0 {
		t.Errorf("Expected intersection length 0, got %d", intersection.Len())
	}

	// Test Difference with uninitialized set
	var s15 abstract.Set[int]
	difference := s15.Difference(other)
	if difference.Len() != 0 {
		t.Errorf("Expected difference length 0, got %d", difference.Len())
	}

	// Test SymmetricDifference with uninitialized set
	var s16 abstract.Set[int]
	symmetricDiff := s16.SymmetricDifference(other)
	if symmetricDiff.Len() != 2 {
		t.Errorf("Expected symmetric difference length 2, got %d", symmetricDiff.Len())
	}
}

func TestSafeSet_UninitializedMethods(t *testing.T) {
	// Test Add with uninitialized safe set
	var s1 abstract.SafeSet[int]
	s1.Add(1, 2, 3)
	if s1.Len() != 3 {
		t.Errorf("Expected length 3 after Add on uninitialized safe set, got %d", s1.Len())
	}

	// Test Has with uninitialized safe set
	var s2 abstract.SafeSet[int]
	if s2.Has(1) {
		t.Error("Expected false for uninitialized safe set")
	}

	// Test Delete with uninitialized safe set
	var s3 abstract.SafeSet[int]
	deleted := s3.Delete(1)
	if deleted {
		t.Error("Expected false from Delete on uninitialized safe set")
	}

	// Test Len with uninitialized safe set
	var s4 abstract.SafeSet[int]
	if s4.Len() != 0 {
		t.Errorf("Expected length 0 for uninitialized safe set, got %d", s4.Len())
	}

	// Test IsEmpty with uninitialized safe set
	var s5 abstract.SafeSet[int]
	if !s5.IsEmpty() {
		t.Error("Expected true from IsEmpty on uninitialized safe set")
	}

	// Test Values with uninitialized safe set
	var s6 abstract.SafeSet[int]
	values := s6.Values()
	if len(values) != 0 {
		t.Errorf("Expected empty values slice, got length %d", len(values))
	}

	// Test Clear with uninitialized safe set
	var s7 abstract.SafeSet[int]
	s7.Clear()
	if s7.Len() != 0 {
		t.Errorf("Expected length 0 after Clear on uninitialized safe set, got %d", s7.Len())
	}

	// Test Transform with uninitialized safe set
	var s8 abstract.SafeSet[int]
	s8.Transform(func(k int) int { return k + 1 })
	if s8.Len() != 0 {
		t.Errorf("Expected no items after Transform on uninitialized safe set, got %d", s8.Len())
	}

	// Test Range with uninitialized safe set
	var s9 abstract.SafeSet[int]
	called := false
	result := s9.Range(func(k int) bool {
		called = true
		return true
	})
	if !result || called {
		t.Error("Expected Range to return true without calling function on uninitialized safe set")
	}

	// Test Raw with uninitialized safe set
	var s10 abstract.SafeSet[int]
	raw := s10.Raw()
	if len(raw) != 0 {
		t.Errorf("Expected empty raw map, got length %d", len(raw))
	}

	// Test Iter with uninitialized safe set
	var s11 abstract.SafeSet[int]
	count := 0
	for range s11.Iter() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 iterations from Iter on uninitialized safe set, got %d", count)
	}

	// Test Copy with uninitialized safe set
	var s12 abstract.SafeSet[int]
	copied := s12.Copy()
	if len(copied) != 0 {
		t.Errorf("Expected empty copied map, got length %d", len(copied))
	}

	// Test Union with uninitialized safe set
	var s13 abstract.SafeSet[int]
	other := map[int]struct{}{1: {}, 2: {}}
	union := s13.Union(other)
	if union.Len() != 2 {
		t.Errorf("Expected union length 2, got %d", union.Len())
	}

	// Test Intersection with uninitialized safe set
	var s14 abstract.SafeSet[int]
	intersection := s14.Intersection(other)
	if intersection.Len() != 0 {
		t.Errorf("Expected intersection length 0, got %d", intersection.Len())
	}

	// Test Difference with uninitialized safe set
	var s15 abstract.SafeSet[int]
	difference := s15.Difference(other)
	if difference.Len() != 0 {
		t.Errorf("Expected difference length 0, got %d", difference.Len())
	}

	// Test SymmetricDifference with uninitialized safe set
	var s16 abstract.SafeSet[int]
	symmetricDiff := s16.SymmetricDifference(other)
	if symmetricDiff.Len() != 2 {
		t.Errorf("Expected symmetric difference length 2, got %d", symmetricDiff.Len())
	}
}

func TestSet_NilInitializationSequence(t *testing.T) {
	// Test that multiple operations work correctly on an uninitialized set
	var s abstract.Set[string]

	// First operation should initialize the map
	s.Add("first")
	if !s.Has("first") {
		t.Error("Expected 'first' to be present after Add")
	}

	// Subsequent operations should work normally
	s.Add("second", "third")
	if s.Len() != 3 {
		t.Errorf("Expected length 3, got %d", s.Len())
	}

	// Test operations that return data
	values := s.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Test modification operations
	deleted := s.Delete("second")
	if !deleted {
		t.Error("Expected Delete to return true")
	}

	if s.Len() != 2 {
		t.Errorf("Expected length 2 after delete, got %d", s.Len())
	}

	// Test set operations
	other := map[string]struct{}{"fourth": {}, "fifth": {}}
	union := s.Union(other)
	if union.Len() != 4 {
		t.Errorf("Expected union length 4, got %d", union.Len())
	}
}

func TestSafeSet_NilInitializationSequence(t *testing.T) {
	// Test that multiple operations work correctly on an uninitialized safe set
	var s abstract.SafeSet[string]

	// First operation should initialize the map
	s.Add("first")
	if !s.Has("first") {
		t.Error("Expected 'first' to be present after Add")
	}

	// Subsequent operations should work normally
	s.Add("second", "third")
	if s.Len() != 3 {
		t.Errorf("Expected length 3, got %d", s.Len())
	}

	// Test operations that return data
	values := s.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Test modification operations
	deleted := s.Delete("second")
	if !deleted {
		t.Error("Expected Delete to return true")
	}

	if s.Len() != 2 {
		t.Errorf("Expected length 2 after delete, got %d", s.Len())
	}

	// Test set operations
	other := map[string]struct{}{"fourth": {}, "fifth": {}}
	union := s.Union(other)
	if union.Len() != 4 {
		t.Errorf("Expected union length 4, got %d", union.Len())
	}
}

func TestSet_ConcurrentNilInitialization(t *testing.T) {
	// Test that concurrent operations on uninitialized SafeSet work correctly
	var s abstract.SafeSet[int]
	var wg sync.WaitGroup

	const numGoroutines = 50

	// Test concurrent Add operations on uninitialized set
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Add(i)
		}(i)
	}

	wg.Wait()

	if s.Len() != numGoroutines {
		t.Errorf("Expected length %d after concurrent operations, got %d", numGoroutines, s.Len())
	}

	// Test concurrent read operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if !s.Has(i) {
				t.Errorf("Expected item %d to be present", i)
			}
		}(i)
	}

	wg.Wait()
}

func TestSet_NilMapBehavior(t *testing.T) {
	// Test that all methods properly handle nil items map
	var s abstract.Set[int]

	// Test that methods don't panic and properly initialize
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic, but got: %v", r)
		}
	}()

	// Each method should handle nil items and initialize properly
	if s.Has(1) {
		t.Error("Expected Has to return false for uninitialized set")
	}

	if s.Len() != 0 {
		t.Error("Expected Len to return 0 for uninitialized set")
	}

	if !s.IsEmpty() {
		t.Error("Expected IsEmpty to return true for uninitialized set")
	}

	values := s.Values()
	if values == nil || len(values) != 0 {
		t.Error("Expected Values to return empty slice for uninitialized set")
	}

	raw := s.Raw()
	if raw == nil || len(raw) != 0 {
		t.Error("Expected Raw to return empty map for uninitialized set")
	}

	copy := s.Copy()
	if copy == nil || len(copy) != 0 {
		t.Error("Expected Copy to return empty map for uninitialized set")
	}

	// Test that iterator works with empty set
	count := 0
	for range s.Iter() {
		count++
	}
	if count != 0 {
		t.Error("Expected Iter to yield no items for uninitialized set")
	}
}

func TestSafeSet_NilMapBehavior(t *testing.T) {
	// Test that all SafeSet methods properly handle nil items map
	var s abstract.SafeSet[int]

	// Test that methods don't panic and properly initialize
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic, but got: %v", r)
		}
	}()

	// Each method should handle nil items and initialize properly
	if s.Has(1) {
		t.Error("Expected Has to return false for uninitialized safe set")
	}

	if s.Len() != 0 {
		t.Error("Expected Len to return 0 for uninitialized safe set")
	}

	if !s.IsEmpty() {
		t.Error("Expected IsEmpty to return true for uninitialized safe set")
	}

	values := s.Values()
	if values == nil || len(values) != 0 {
		t.Error("Expected Values to return empty slice for uninitialized safe set")
	}

	raw := s.Raw()
	if raw == nil || len(raw) != 0 {
		t.Error("Expected Raw to return empty map for uninitialized safe set")
	}

	copy := s.Copy()
	if copy == nil || len(copy) != 0 {
		t.Error("Expected Copy to return empty map for uninitialized safe set")
	}

	// Test that iterator works with empty set
	count := 0
	for range s.Iter() {
		count++
	}
	if count != 0 {
		t.Error("Expected Iter to yield no items for uninitialized safe set")
	}
}
