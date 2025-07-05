package abstract_test

import (
	"sync"
	"testing"

	"github.com/maxbolgarin/abstract"
)

// TestSlice tests all methods for the Slice type.
func TestSlice(t *testing.T) {
	// Test NewSlice and Append
	slice := abstract.NewSlice[int]()
	if !slice.IsEmpty() {
		t.Error("expected slice to be empty")
	}
	slice.Append(1)
	if slice.Len() != 1 {
		t.Errorf("expected length 1, got %d", slice.Len())
	}

	// Test Get
	if slice.Get(0) != 1 {
		t.Errorf("expected element 1, got %d", slice.Get(0))
	}
	if slice.Get(1) != 0 {
		t.Errorf("expected default zero value, got %d", slice.Get(1))
	}

	// Test AddFront
	slice.AddFront(2)
	if slice.Get(0) != 2 {
		t.Errorf("expected element 2, got %d", slice.Get(0))
	}
	if slice.Get(1) != 1 {
		t.Errorf("expected element 1, got %d", slice.Get(1))
	}

	// Test Pop
	popped := slice.Pop()
	if popped != 1 {
		t.Errorf("expected popped element 1, got %d", popped)
	}
	slice.Pop()
	if !slice.IsEmpty() {
		t.Error("expected slice to be empty after pop")
	}

	// Test Delete
	slice.Append(1)
	slice.Append(2)
	deleted := slice.Delete(0)
	if !deleted {
		t.Errorf("expected delete to be successful")
	}
	if slice.Len() != 1 {
		t.Errorf("expected length 1, got %d", slice.Len())
	}

	// Test IsEmpty and Clear
	slice.Clear()
	if !slice.IsEmpty() {
		t.Error("expected slice to be empty after clear")
	}

	// Test Copy
	slice.Append(3)
	copy := slice.Copy()
	if len(copy) != 1 || copy[0] != 3 {
		t.Error("copy did not match original slice")
	}

	// Test Transform
	slice.Transform(func(x int) int { return x * 2 })
	if slice.Get(0) != 6 {
		t.Errorf("expected transformed element 6, got %d", slice.Get(0))
	}

	// Test Range
	slice.Append(8) // Now [6, 8]
	sum := 0
	slice.Range(func(x int) bool {
		sum += x
		return true
	})
	if sum != 14 {
		t.Errorf("expected sum 14, got %d", sum)
	}

	slice = abstract.NewSlice([]int{1, 2, 3})
	slice.Transform(func(x int) int { return x * 2 })
	if slice.Get(0) != 2 || slice.Get(1) != 4 || slice.Get(2) != 6 {
		t.Error("expected transformed elements to match original slice")
	}

	slice = abstract.NewSliceWithSize[int](10)
	if slice.Len() != 0 {
		t.Errorf("expected length 0, got %d", slice.Len())
	}
	if slice.Pop() != 0 {
		t.Errorf("expected default zero value, got %d", slice.Pop())
	}
	if slice.Delete(0) {
		t.Errorf("expected delete to fail")
	}

	slice.Append(1, 2, 3, 4, 5)
	if slice.Len() != 5 {
		t.Errorf("expected length 5, got %d", slice.Len())
	}
	slice.Truncate(2)
	if slice.Len() != 2 {
		t.Errorf("expected length 2, got %d", slice.Len())
	}
	var i int
	if !slice.Range(func(x int) bool {
		i++
		return true
	}) {
		t.Error("expected Range to return true")
	}
	if i != 2 {
		t.Errorf("expected Range to iterate 2 times, got %d", i)
	}
	i = 0
	if slice.Range(func(x int) bool {
		i++
		return false
	}) {
		t.Error("expected Range to return false")
	}
	if i != 1 {
		t.Errorf("expected Range to iterate 1 time, got %d", i)
	}
	slice = abstract.NewSliceFromItems(1, 2, 3)
	if slice.Len() != 3 {
		t.Errorf("expected length 3, got %d", slice.Len())
	}
	slice.Transform(func(x int) int { return x * 2 })
	if slice.Get(0) != 2 || slice.Get(1) != 4 || slice.Get(2) != 6 {
		t.Error("expected transformed elements to match original slice")
	}
	var counter = 1
	iter := slice.Iter()
	for x := range iter {
		if x != counter*2 {
			t.Errorf("Expected %d, got %d", counter, x)
		}
		counter++
	}
	if counter != 4 {
		t.Errorf("Expected 3, got %d", counter)
	}
	counter = 1
	iter2 := slice.Iter2()
	for x, y := range iter2 {
		if x != counter-1 {
			t.Errorf("Expected %d, got %d", counter, x)
		}
		if y != counter*2 {
			t.Errorf("Expected %d, got %d", counter*2, y)
		}
		counter++
	}
	if counter != 4 {
		t.Errorf("Expected 3, got %d", counter)
	}
}

func TestChangeSlice(t *testing.T) {
	slice := abstract.NewSliceFromItems(1, 2, 3)
	slice.Change(0, func(x int) int { return x * 2 })
	if slice.Get(0) != 2 {
		t.Errorf("expected transformed element 2, got %d", slice.Get(0))
	}
	slice.Change(1, func(x int) int { return x * 2 })
	if slice.Get(1) != 4 {
		t.Errorf("expected transformed element 4, got %d", slice.Get(1))
	}
}

func TestSlice_Raw(t *testing.T) {
	slice := abstract.NewSliceFromItems(1, 2, 3)
	raw := slice.Raw()
	if len(raw) != 3 {
		t.Errorf("expected length 3, got %d", len(raw))
	}
}

// TestSafeSlice tests all methods for the SafeSlice type with concurrency.
func TestSafeSlice(t *testing.T) {
	var wg sync.WaitGroup
	safeSlice := abstract.NewSafeSlice[int]()

	for i := 0; i < 1000; i++ {
		safeSlice.Append(i)
	}

	for i := 1000; i < 2000; i++ {
		safeSlice.Append(i)
	}

	if safeSlice.Len() != 2000 {
		t.Errorf("expected length 2000, got %d", safeSlice.Len())
	}

	// Test concurrent Get
	wg.Add(1)
	go func() {
		defer wg.Done()
		if safeSlice.Get(500) == 0 {
			t.Error("expected non-zero value")
		}
	}()
	wg.Wait()

	// Test Pop
	element := safeSlice.Pop()
	if element != 1999 {
		t.Errorf("expected 1999, got %d", element)
	}

	// Test Delete
	wg.Add(1)
	go func() {
		defer wg.Done()
		if !safeSlice.Delete(0) {
			t.Error("expected successful deletion")
		}
	}()
	wg.Wait()

	// Test Transform
	safeSlice.Transform(func(x int) int { return x + 1 })
	if safeSlice.Get(0) != 2 { // 1 + 1 after transformation
		t.Errorf("expected transformed element 2, got %d", safeSlice.Get(0))
	}

	// Test Range
	wg.Add(1)
	go func() {
		defer wg.Done()
		safeSlice.Range(func(x int) bool {
			return x > 0
		})
	}()
	wg.Wait()

	// Test Clear
	safeSlice.Clear()
	if !safeSlice.IsEmpty() {
		t.Error("expected safe slice to be empty after clear")
	}

	slice := abstract.NewSafeSlice([]int{1, 2, 3})
	slice.Transform(func(x int) int { return x * 2 })
	if slice.Get(0) != 2 || slice.Get(1) != 4 || slice.Get(2) != 6 {
		t.Error("expected transformed elements to match original slice")
	}

	slice = abstract.NewSafeSliceWithSize[int](10)
	if slice.Len() != 0 {
		t.Errorf("expected length 0, got %d", slice.Len())
	}
	if slice.Pop() != 0 {
		t.Errorf("expected default zero value, got %d", slice.Pop())
	}
	if slice.Delete(0) {
		t.Errorf("expected delete to fail")
	}

	slice.Append(1, 2, 3, 4, 5)
	if slice.Len() != 5 {
		t.Errorf("expected length 5, got %d", slice.Len())
	}
	slice.AddFront(6, 7, 8)
	if slice.Get(0) != 6 {
		t.Errorf("expected element 6, got %d", slice.Get(0))
	}
	if slice.Get(1) != 7 {
		t.Errorf("expected element 7, got %d", slice.Get(1))
	}
	iter := slice.Iter()
	for x := range iter {
		t.Log(x)
	}
	iter2 := slice.Iter2()
	for x, y := range iter2 {
		t.Log(x, y)
	}
	slice.Truncate(2)
	if slice.Len() != 2 {
		t.Errorf("expected length 2, got %d", slice.Len())
	}
	var i int
	if !slice.Range(func(x int) bool {
		i++
		return true
	}) {
		t.Error("expected Range to return true")
	}
	if i != 2 {
		t.Errorf("expected Range to iterate 2 times, got %d", i)
	}
	i = 0
	if slice.Range(func(x int) bool {
		i++
		return false
	}) {
		t.Error("expected Range to return false")
	}
	if i != 1 {
		t.Errorf("expected Range to iterate 1 time, got %d", i)
	}
	if slice.Get(590) != 0 {
		t.Errorf("expected default zero value, got %d", slice.Get(590))
	}

	slice2 := slice.Copy()
	if len(slice2) != 2 {
		t.Errorf("expected length 2, got %d", len(slice2))
	}

	slice3 := abstract.NewSafeSliceFromItems(1, 2, 3)
	if slice3.Len() != 3 {
		t.Errorf("expected length 3, got %d", slice3.Len())
	}
	slice3.Transform(func(x int) int { return x * 2 })
	if slice3.Get(0) != 2 || slice3.Get(1) != 4 || slice3.Get(2) != 6 {
		t.Error("expected transformed elements to match original slice")
	}
}

func TestChangeSafeSlice(t *testing.T) {
	slice := abstract.NewSafeSliceFromItems(1, 2, 3)
	slice.Change(0, func(x int) int { return x * 2 })
	if slice.Get(0) != 2 {
		t.Errorf("expected transformed element 2, got %d", slice.Get(0))
	}
	slice.Change(1, func(x int) int { return x * 2 })
	if slice.Get(1) != 4 {
		t.Errorf("expected transformed element 4, got %d", slice.Get(1))
	}
}

func TestSafeSlice_Raw(t *testing.T) {
	slice := abstract.NewSafeSliceFromItems(1, 2, 3)
	raw := slice.Raw()
	if len(raw) != 3 {
		t.Errorf("expected length 3, got %d", len(raw))
	}
}
