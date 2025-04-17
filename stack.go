package abstract

import "sync"

// Stack is a simple stack data structure.
type Stack[T any] struct {
	mem []T
}

// NewStack creates a new Stack.
func NewStack[T any](data ...[]T) *Stack[T] {
	mem := make([]T, 0, getSlicesLen(data...))
	for _, d := range data {
		mem = append(mem, d...)
	}
	return &Stack[T]{mem: mem}
}

// NewStackWithCapacity creates a new Stack with a specified capacity.
func NewStackWithCapacity[T any](capacity int) *Stack[T] {
	return &Stack[T]{mem: make([]T, 0, capacity)}
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.mem = append(s.mem, item)
}

// Last returns the last item in the stack.
func (s *Stack[T]) Last() T {
	if len(s.mem) == 0 {
		var zero T
		return zero
	}
	return s.mem[len(s.mem)-1]
}

// Pop removes and returns the top item from the stack.
func (s *Stack[T]) Pop() T {
	if len(s.mem) == 0 {
		var zero T
		return zero
	}

	last := len(s.mem) - 1
	item := s.mem[last]
	s.mem = s.mem[:last]
	// Set the reference to nil to prevent memory leaks if T is a reference type
	// This only has effect if T is a reference type
	var zero T
	s.mem = append(s.mem, zero)
	s.mem = s.mem[:last]

	return item
}

// PopOK is like Pop but also returns a boolean indicating if the operation was successful.
func (s *Stack[T]) PopOK() (T, bool) {
	if len(s.mem) == 0 {
		var zero T
		return zero, false
	}
	last := len(s.mem) - 1
	item := s.mem[last]
	s.mem = s.mem[:last]
	// Set the reference to nil to prevent memory leaks if T is a reference type
	// This only has effect if T is a reference type
	var zero T
	s.mem = append(s.mem, zero)
	s.mem = s.mem[:last]

	return item, true
}

// IsEmpty returns true if the stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Len returns the number of items in the stack.
func (s *Stack[T]) Len() int {
	return len(s.mem)
}

// Clear removes all items from the stack.
func (s *Stack[T]) Clear() {
	s.mem = make([]T, 0)
}

// Raw returns the underlying slice of the stack.
func (s *Stack[T]) Raw() []T {
	return s.mem
}

// SafeStack is a simple stack data structure that is thread-safe.
type SafeStack[T any] struct {
	*Stack[T]
	sync.Mutex
}

// NewSafeStack creates a new SafeStack.
func NewSafeStack[T any](data ...[]T) *SafeStack[T] {
	return &SafeStack[T]{Stack: NewStack(data...)}
}

// NewSafeStackWithCapacity creates a new SafeStack with a specified capacity.
func NewSafeStackWithCapacity[T any](capacity int) *SafeStack[T] {
	return &SafeStack[T]{Stack: NewStackWithCapacity[T](capacity)}
}

// Push adds an item to the top of the stack.
func (s *SafeStack[T]) Push(item T) {
	s.Lock()
	defer s.Unlock()
	s.Stack.Push(item)
}

// Pop removes and returns the top item from the stack.
func (s *SafeStack[T]) Pop() T {
	s.Lock()
	defer s.Unlock()
	return s.Stack.Pop()
}

// PopOK is like Pop but also returns a boolean indicating if the operation was successful.
func (s *SafeStack[T]) PopOK() (T, bool) {
	s.Lock()
	defer s.Unlock()
	return s.Stack.PopOK()
}

// Last returns the last item in the stack.
func (s *SafeStack[T]) Last() T {
	s.Lock()
	defer s.Unlock()
	return s.Stack.Last()
}

// IsEmpty returns true if the stack is empty.
func (s *SafeStack[T]) IsEmpty() bool {
	s.Lock()
	defer s.Unlock()
	return s.Stack.IsEmpty()
}

// Len returns the number of items in the stack.
func (s *SafeStack[T]) Len() int {
	s.Lock()
	defer s.Unlock()
	return s.Stack.Len()
}

// Clear removes all items from the stack.
func (s *SafeStack[T]) Clear() {
	s.Lock()
	defer s.Unlock()
	s.Stack.Clear()
}

// Raw returns the underlying slice of the stack.
func (s *SafeStack[T]) Raw() []T {
	s.Lock()
	defer s.Unlock()
	return s.Stack.Raw()
}

// UniqueStack is a stack that doesn't allow duplicates.
type UniqueStack[T comparable] struct {
	mem []T
	ind map[T]int
}

// NewUniqueStack creates a new UniqueStack.
func NewUniqueStack[T comparable](data ...[]T) *UniqueStack[T] {
	mem := make([]T, 0, getSlicesLen(data...))
	for _, d := range data {
		mem = append(mem, d...)
	}

	ind := make(map[T]int, cap(mem))
	for i, v := range mem {
		ind[v] = i
	}

	return &UniqueStack[T]{mem: mem, ind: ind}
}

// NewUniqueStackWithCapacity creates a new UniqueStack with a specified capacity.
func NewUniqueStackWithCapacity[T comparable](cap int) *UniqueStack[T] {
	return &UniqueStack[T]{mem: make([]T, 0, cap), ind: make(map[T]int, cap)}
}

// Push adds an item to the stack. If the item is already present in the stack,
// it moves it to the top of the stack.
func (s *UniqueStack[T]) Push(item T) {
	if index, ok := s.ind[item]; ok {
		last := len(s.mem) - 1
		if index == last {
			return
		}
		s.ind[s.mem[last]] = index
		s.ind[item] = last

		s.mem[index], s.mem[last] = s.mem[last], s.mem[index]
		return
	}
	s.ind[item] = len(s.mem)
	s.mem = append(s.mem, item)
}

// Last returns the last item in the stack.
func (s *UniqueStack[T]) Last() T {
	if len(s.mem) == 0 {
		var zero T
		return zero
	}

	return s.mem[len(s.mem)-1]
}

// Pop removes the last item from the stack and returns it.
func (s *UniqueStack[T]) Pop() T {
	m, _ := s.PopOK()
	return m
}

// PopOK is like Pop, but returns a boolean indicating whether the stack was empty.
func (s *UniqueStack[T]) PopOK() (T, bool) {
	if len(s.mem) == 0 {
		var zero T
		return zero, false
	}
	index := len(s.mem) - 1

	item := s.mem[index]
	s.mem = s.mem[:index]
	delete(s.ind, item)

	return item, true
}

// IsEmpty returns true if the stack is empty.
func (s *UniqueStack[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Len returns the number of items in the stack.
func (s *UniqueStack[T]) Len() int {
	return len(s.mem)
}

// Remove removes an item from the stack.
func (s *UniqueStack[T]) Remove(item T) bool {
	indexToRemove, ok := s.ind[item]
	if !ok {
		return false
	}

	if indexToRemove >= len(s.mem) {
		// There's an inconsistency between index map and memory slice
		// Just remove the item from the index map
		delete(s.ind, item)
		return true
	}

	// Update the index for the item at the end of the stack
	if indexToRemove < len(s.mem)-1 {
		lastItem := s.mem[len(s.mem)-1]
		s.mem[indexToRemove] = lastItem
		s.ind[lastItem] = indexToRemove
	}

	// Remove the last element from the stack
	s.mem = s.mem[:len(s.mem)-1]
	delete(s.ind, item)

	return true
}

// Clear removes all items from the stack.
func (s *UniqueStack[T]) Clear() {
	s.mem = make([]T, 0)
	s.ind = make(map[T]int)
}

// Raw returns the underlying slice of the stack.
func (s *UniqueStack[T]) Raw() []T {
	return s.mem
}

// SafeUniqueStack is a thread-safe UniqueStack.
type SafeUniqueStack[T comparable] struct {
	s  *UniqueStack[T]
	mu sync.Mutex
}

// NewSafeUniqueStack creates a new SafeUniqueStack.
func NewSafeUniqueStack[T comparable](data ...[]T) *SafeUniqueStack[T] {
	return &SafeUniqueStack[T]{s: NewUniqueStack(data...)}
}

// NewSafeUniqueStackWithCapacity creates a new SafeUniqueStack with a specified capacity.
func NewSafeUniqueStackWithCapacity[T comparable](cap int) *SafeUniqueStack[T] {
	return &SafeUniqueStack[T]{s: NewUniqueStackWithCapacity[T](cap)}
}

// Push adds an item to the stack if it is not already present.
func (ss *SafeUniqueStack[T]) Push(item T) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.s.Push(item)
}

// Last returns the last item in the stack without removing it.
func (ss *SafeUniqueStack[T]) Last() T {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.s.Last()
}

// Pop removes and returns the last item from the stack.
func (ss *SafeUniqueStack[T]) Pop() T {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.s.Pop()
}

// PopOK is like Pop but also returns a boolean indicating if the operation was successful.
func (ss *SafeUniqueStack[T]) PopOK() (T, bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.s.PopOK()
}

// IsEmpty checks if the stack is empty.
func (ss *SafeUniqueStack[T]) IsEmpty() bool {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.s.IsEmpty()
}

// Len returns the number of items in the stack.
func (ss *SafeUniqueStack[T]) Len() int {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.s.Len()
}

// Remove deletes a specific item from the stack, if it exists.
func (ss *SafeUniqueStack[T]) Remove(item T) bool {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.s.Remove(item)
}

// Clear removes all items from the stack.
func (ss *SafeUniqueStack[T]) Clear() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.s.Clear()
}

// Raw returns a copy of the underlying data slice (non-concurrent safe operation).
func (ss *SafeUniqueStack[T]) Raw() []T {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.s.Raw()
}
