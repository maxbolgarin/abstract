package abstract

import "sync"

// Slice is used like a common slice.
type Slice[T any] []T

// NewSlice returns a new [Slice] with empty slice.
func NewSlice[T any](data ...[]T) Slice[T] {
	out := make([]T, 0, getSlicesLen(data...))
	for _, v := range data {
		out = append(out, v...)
	}
	return out
}

// NewSliceWithSize returns a new [Slice] with slice inited using the provided size.
func NewSliceWithSize[T any](size int) Slice[T] {
	return make([]T, 0, size)
}

// Get returns the value for the provided key or the default type value if the key is not present in the slice.
func (s *Slice[T]) Get(index int) T {
	if index >= len(*s) {
		var empty T
		return empty
	}
	return (*s)[index]
}

// Append adds a new element to the end of the slice.
func (s *Slice[T]) Append(v ...T) {
	*s = append(*s, v...)
}

// Pop removes the last element of the slice and returns it.
func (s *Slice[T]) Pop() T {
	if len(*s) == 0 {
		var empty T
		return empty
	}

	elem := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return elem
}

// Delete removes the key and associated value from the slice, does nothing if the key is not present in the slice,
// returns true if the key was deleted.
func (s *Slice[T]) Delete(index int) bool {
	if index >= len(*s) {
		return false
	}
	*s = append((*s)[:index], (*s)[index+1:]...)
	return true
}

// Len returns the length of the slice.
func (s *Slice[T]) Len() int {
	return len(*s)
}

// IsEmpty returns true if the slice is empty. It is safe for concurrent/parallel use.
func (s *Slice[T]) IsEmpty() bool {
	return len(*s) == 0
}

// Truncate truncates the slice to the provided size.
func (s *Slice[T]) Truncate(size int) {
	*s = (*s)[:size]
}

// Clear creates a new slice using make without size.
func (s *Slice[T]) Clear() {
	*s = make([]T, 0)
}

// Copy returns a copy of the slice.
func (s *Slice[T]) Copy() Slice[T] {
	return append(make([]T, 0, len(*s)), *s...)
}

// Transform transforms all values of the slice using provided function.
func (s *Slice[T]) Transform(f func(T) T) {
	for i, v := range *s {
		(*s)[i] = f(v)
	}
}

// Range calls the provided function for each element in the slice.
func (s *Slice[T]) Range(f func(T) bool) {
	for _, v := range *s {
		if !f(v) {
			return
		}
	}
}

// SafeSlice is used like a common slice, but it is protected with RW mutex, so it can be used in many goroutines.
type SafeSlice[T any] struct {
	items []T
	mu    sync.RWMutex
}

// NewSafeSlice returns a new [SafeSlice] with empty slice.
func NewSafeSlice[T any](data ...[]T) *SafeSlice[T] {
	out := &SafeSlice[T]{
		items: make([]T, 0, getSlicesLen(data...)),
	}
	for _, v := range data {
		out.items = append(out.items, v...)
	}
	return out
}

// NewSafeSliceWithSize returns a new [SafeSlice] with slice inited using the provided size.
func NewSafeSliceWithSize[T any](size int) *SafeSlice[T] {
	return &SafeSlice[T]{
		items: make([]T, 0, size),
	}
}

// Get returns the value for the provided key or the default type value if the key is not present in the slice.
func (s *SafeSlice[T]) Get(index int) T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if index >= len(s.items) {
		var empty T
		return empty
	}
	return s.items[index]
}

// Append adds a new element to the end of the slice.
func (s *SafeSlice[T]) Append(v ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = append(s.items, v...)
}

// Pop removes the last element of the slice and returns it.
func (s *SafeSlice[T]) Pop() T {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.items) == 0 {
		var empty T
		return empty
	}

	elem := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return elem
}

// Delete removes the key and associated value from the slice, does nothing if the key is not present in the slice,
// returns true if the key was deleted.
func (s *SafeSlice[T]) Delete(index int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index >= len(s.items) {
		return false
	}
	s.items = append(s.items[:index], s.items[index+1:]...)
	return true
}

// Len returns the length of the slice. It is safe for concurrent/parallel use.
func (s *SafeSlice[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

// IsEmpty returns true if the slice is empty. It is safe for concurrent/parallel use.
func (s *SafeSlice[T]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items) == 0
}

// Truncate truncates the slice to the provided size.
func (s *SafeSlice[T]) Truncate(size int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = s.items[:size]
}

// Clear creates a new slice using make without size.
func (s *SafeSlice[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = make([]T, 0)
}

// Copy returns a copy of the slice.
func (s *SafeSlice[T]) Copy() *SafeSlice[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &SafeSlice[T]{
		items: append(make([]T, 0, len(s.items)), s.items...),
	}
}

// Transform transforms all values of the slice using provided function.
func (s *SafeSlice[T]) Transform(f func(T) T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, v := range s.items {
		s.items[i] = f(v)
	}
}

// Range calls the provided function for each element in the slice.
func (s *SafeSlice[T]) Range(f func(T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, v := range s.items {
		if !f(v) {
			return
		}
	}
}

func getSlicesLen[T any](slices ...[]T) int {
	var length int
	for _, slice := range slices {
		length += len(slice)
	}
	return length
}
