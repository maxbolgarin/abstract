package abstract_test

import (
	"sync"
	"testing"

	"github.com/maxbolgarin/abstract"
)

func TestStack(t *testing.T) {
	stack := abstract.NewStack([]int{1, 2})

	// Test Push
	stack.Push(1)
	stack.Push(2)
	if stack.Len() != 4 {
		t.Errorf("Expected length 2, got %d", stack.Len())
	}

	// Test Last
	if last := stack.Last(); last != 2 {
		t.Errorf("Expected last item 2, got %d", last)
	}

	// Test Pop
	if popped := stack.Pop(); popped != 2 {
		t.Errorf("Expected popped item 2, got %d", popped)
	}
	if stack.Len() != 3 {
		t.Errorf("Expected length 1, got %d", stack.Len())
	}

	// Test PopOK
	item, ok := stack.PopOK()
	if !ok || item != 1 {
		t.Errorf("Expected true and item 1, got %v and %d", ok, item)
	}

	stack.Pop()
	stack.Pop()
	stack.Pop()
	stack.Pop()
	stack.PopOK()

	// Test IsEmpty
	if !stack.IsEmpty() {
		t.Errorf("Expected stack to be empty")
	}

	stack = abstract.NewStackWithCapacity[int](10)

	// Test Clear
	stack.Push(1)
	stack.Clear()
	if stack.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", stack.Len())
	}

	// Test Raw
	if raw := stack.Raw(); len(raw) != 0 {
		t.Errorf("Expected raw length 0, got %d", len(raw))
	}

	if stack.Last() != 0 {
		t.Errorf("Expected last item 0, got %d", stack.Last())
	}
}

func TestSafeStack(t *testing.T) {
	safeStack := abstract.NewSafeStack([]int{1, 2})

	// Test Push
	safeStack.Push(1)
	safeStack.Push(2)
	if safeStack.Len() != 4 {
		t.Errorf("Expected length 4, got %d", safeStack.Len())
	}

	// Test Last
	if last := safeStack.Last(); last != 2 {
		t.Errorf("Expected last item 2, got %d", last)
	}

	// Test Pop
	if popped := safeStack.Pop(); popped != 2 {
		t.Errorf("Expected popped item 2, got %d", popped)
	}
	if safeStack.Len() != 3 {
		t.Errorf("Expected length 3, got %d", safeStack.Len())
	}

	// Test PopOK
	item, ok := safeStack.PopOK()
	if !ok || item != 1 {
		t.Errorf("Expected true and item 1, got %v and %d", ok, item)
	}

	safeStack.Pop()
	safeStack.Pop()
	safeStack.Pop()
	safeStack.Pop()
	safeStack.PopOK()

	// Test IsEmpty
	if !safeStack.IsEmpty() {
		t.Errorf("Expected stack to be empty")
	}

	safeStack = abstract.NewSafeStackWithCapacity[int](10)

	// Test Clear
	safeStack.Push(1)
	safeStack.Clear()
	if safeStack.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", safeStack.Len())
	}

	// Test Raw
	if raw := safeStack.Raw(); len(raw) != 0 {
		t.Errorf("Expected raw length 0, got %d", len(raw))
	}
}

func TestNewUniqueStack(t *testing.T) {
	stack := abstract.NewUniqueStack[int]()
	if stack == nil {
		t.Fatalf("Expected non-nil Stack object")
	}
	if stack.Len() != 0 {
		t.Errorf("Expected the stack to be empty, got length %d", stack.Len())
	}

	// Test NewStack with initial data
	initialData := []int{1, 2, 3}
	stackWithData := abstract.NewUniqueStack(initialData)
	if stackWithData.Len() != len(initialData) {
		t.Errorf("Expected stack length %d, got %d", len(initialData), stackWithData.Len())
	}
}

func TestNewUniqueStackWithCapacity(t *testing.T) {
	capacity := 10
	stack := abstract.NewUniqueStackWithCapacity[int](capacity)
	if stack.Len() != 0 {
		t.Errorf("Expected the stack to be empty, got length %d", stack.Len())
	}
}

func TestUniqueStack_Push(t *testing.T) {
	stack := abstract.NewUniqueStack[int]()
	stack.Push(1)
	for i := 0; i <= 10; i++ {
		stack.Push(2)
	}
	for i := 0; i <= 10; i++ {
		stack.Push(1)
	}

	if stack.Len() != 2 {
		t.Errorf("Expected stack length 2, got %d", stack.Len())
	}

	last := stack.Last()
	if last != 1 {
		t.Errorf("Expected last element to be 2, got %d", last)
	}

	for i := 0; i <= 10; i++ {
		stack.Push(i)
	}

	last = stack.Last()
	if last != 10 {
		t.Errorf("Expected last element to be 2, got %d", last)
	}

	stack.Push(5)

	last = stack.Last()
	if last != 5 {
		t.Errorf("Expected last element to be 2, got %d", last)
	}
}

func TestUniqueStack_Last(t *testing.T) {
	stack := abstract.NewUniqueStack[int]()
	if last := stack.Last(); last != 0 {
		t.Errorf("Expected zero value for last element, got %v", last)
	}

	stack.Push(1)
	if last := stack.Last(); last != 1 {
		t.Errorf("Expected last element to be 1, got %d", last)
	}
}

func TestStack_Pop(t *testing.T) {
	stack := abstract.NewUniqueStack[int]()
	item, ok := stack.PopOK()
	if ok || item != 0 {
		t.Errorf("Expected false and zero value when popping from empty stack, got %v, %v", item, ok)
	}

	stack.Push(1)
	stack.Push(2)

	item = stack.Pop()
	if item != 2 {
		t.Errorf("Expected 2 when popping, got %d", item)
	}

	if stack.Len() != 1 {
		t.Errorf("Expected stack length 1 after popping, got %d", stack.Len())
	}

	if stack.Last() != 1 {
		t.Errorf("Expected last element to be 1, got %d", stack.Last())
	}
}

func TestUniqueStack_IsEmpty(t *testing.T) {
	stack := abstract.NewUniqueStack[int]()
	if !stack.IsEmpty() {
		t.Errorf("Expected stack to be empty")
	}

	stack.Push(1)
	if stack.IsEmpty() {
		t.Errorf("Expected stack not to be empty")
	}
}

func TestUniqueStack_Remove(t *testing.T) {
	stack := abstract.NewUniqueStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	ok := stack.Remove(2)
	if !ok {
		t.Errorf("Expected to successfully remove 2")
	}

	if stack.Len() != 2 {
		t.Errorf("Expected stack length 2 after removal, got %d", stack.Len())
	}

	// Test removing a non-existent item
	ok = stack.Remove(4)
	if ok {
		t.Errorf("Expected to fail to remove non-existent item")
	}

	// Test removing the last item
	ok = stack.Remove(1)
	if !ok {
		t.Errorf("Expected to successfully remove 1")
	}

	if stack.Len() != 1 {
		t.Errorf("Expected stack length 1 after removal, got %d", stack.Len())
	}

	// Test removing the only item
	ok = stack.Remove(3)
	if !ok {
		t.Errorf("Expected to successfully remove 3")
	}
}

func TestUniqueStack_Clear(t *testing.T) {
	stack := abstract.NewUniqueStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Clear()

	if stack.Len() != 0 {
		t.Errorf("Expected stack to be empty after Clear, got length %d", stack.Len())
	}
}

func TestUniqueStack_Raw(t *testing.T) {
	stack := abstract.NewUniqueStack[int]()
	stack.Push(1)
	stack.Push(2)

	raw := stack.Raw()
	if len(raw) != 2 || raw[0] != 1 || raw[1] != 2 {
		t.Errorf("Raw method returned unexpected slice: %v", raw)
	}
}

func TestSafeUniqueStack_Push(t *testing.T) {
	stack := abstract.NewSafeUniqueStack[int]()
	stack.Push(1)

	if stack.Len() != 1 {
		t.Errorf("Expected length 1, got %d", stack.Len())
	}

	stack.Push(1) // Push duplicate
	if stack.Len() != 1 {
		t.Errorf("Duplicate item added, expected length 1, got %d", stack.Len())
	}

	stack.Push(2)
	if stack.Len() != 2 {
		t.Errorf("Expected length 2, got %d", stack.Len())
	}

	stack.Clear()

	if stack.Len() != 0 {
		t.Errorf("Expected length 0, got %d", stack.Len())
	}
}

func TestSafeUniqueStack_Last(t *testing.T) {
	stack := abstract.NewSafeUniqueStack[int]()
	stack.Push(1)
	stack.Push(2)

	last := stack.Last()
	if last != 2 {
		t.Errorf("Expected last value 2, got %d", last)
	}

	// Ensure the stack size is unchanged
	if stack.Len() != 2 {
		t.Errorf("Expected length 2, got %d", stack.Len())
	}
}

func TestSafeUniqueStack_Pop(t *testing.T) {
	stack := abstract.NewSafeUniqueStack[int]()
	stack.Push(1)
	stack.Push(2)

	item := stack.Pop()
	if item != 2 {
		t.Errorf("Expected popped value 2, got %d", item)
	}

	if stack.Len() != 1 {
		t.Errorf("Expected length 1, got %d", stack.Len())
	}

	if !stack.Remove(1) {
		t.Errorf("Expected to successfully remove 1")
	}

	if len(stack.Raw()) != 0 {
		t.Errorf("Expected length 0, got %d", len(stack.Raw()))
	}

	if !stack.IsEmpty() {
		t.Errorf("Expected stack to be empty")
	}
}

func TestSafeUniqueStack_PopEmpty(t *testing.T) {
	stack := abstract.NewSafeUniqueStack[int]()

	item := stack.Pop()
	if item != *new(int) {
		t.Errorf("Expected default zero value, got %d", item)
	}

	if stack.Len() != 0 {
		t.Errorf("Expected length 0, got %d", stack.Len())
	}

}

func TestSafeUniqueStack_PopOK(t *testing.T) {
	stack := abstract.NewSafeUniqueStackWithCapacity[int](10)
	stack.Push(1)

	item, ok := stack.PopOK()
	if !ok || item != 1 {
		t.Errorf("Expected pop (1, true), got (%d, %v)", item, ok)
	}

	if stack.Len() != 0 {
		t.Errorf("Expected length 0, got %d", stack.Len())
	}

	item, ok = stack.PopOK()
	if ok || item != *new(int) {
		t.Errorf("Expected pop (0, false), got (%d, %v)", item, ok)
	}
}

func TestSafeUniqueStack_Concurrency(t *testing.T) {
	stack := abstract.NewSafeUniqueStack[int]()

	var wg sync.WaitGroup
	const numGoroutines = 100
	const numInserts = 1000

	// Concurrently push the same set of elements
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numInserts; j++ {
				stack.Push(j)
			}
		}()
	}

	wg.Wait()

	// The length should equal numInserts because each number should only appear once
	if stack.Len() != numInserts {
		t.Errorf("Expected length %d, got %d", numInserts, stack.Len())
	}
}
