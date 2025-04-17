package abstract

import (
	"iter"
	"slices"
	"sync"
)

// Slice is used like a common slice.
type Slice[T any] struct {
	items []T
}

// NewSlice returns a new [Slice] with empty slice.
func NewSlice[T any](data ...[]T) *Slice[T] {
	out := make([]T, 0, getSlicesLen(data...))
	for _, v := range data {
		out = append(out, v...)
	}
	return &Slice[T]{
		items: out,
	}
}

// NewSliceFromItems returns a new [Slice] with the provided items.
func NewSliceFromItems[T any](data ...T) *Slice[T] {
	return &Slice[T]{
		items: data,
	}
}

// NewSliceWithSize returns a new [Slice] with slice inited using the provided size.
func NewSliceWithSize[T any](size int) *Slice[T] {
	return &Slice[T]{
		items: make([]T, 0, size),
	}
}

// Get returns the value for the provided key or the default type value if the key is not present in the slice.
func (s *Slice[T]) Get(index int) T {
	if index >= len(s.items) {
		var empty T
		return empty
	}
	return s.items[index]
}

// Pop removes the last element of the slice and returns it.
func (s *Slice[T]) Pop() T {
	if len(s.items) == 0 {
		var empty T
		return empty
	}

	elem := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return elem
}

// Append adds a new element to the end of the slice.
func (s *Slice[T]) Append(v ...T) {
	s.items = append(s.items, v...)
}

// AddFront adds a new elements to the front of the slice.
func (s *Slice[T]) AddFront(v ...T) {
	s.items = append(v, s.items...)
}

// Delete removes the key and associated value from the slice, does nothing if the key is not present in the slice,
// returns true if the key was deleted.
func (s *Slice[T]) Delete(index int) bool {
	if index < 0 || index >= len(s.items) {
		return false
	}
	s.items = append(s.items[:index], s.items[index+1:]...)
	return true
}

// Truncate truncates the slice to the provided size.
func (s *Slice[T]) Truncate(size int) {
	s.items = s.items[:size]
}

// Clear creates a new slice using make without size.
func (s *Slice[T]) Clear() {
	s.items = make([]T, 0)
}

// Len returns the length of the slice.
func (s *Slice[T]) Len() int {
	return len(s.items)
}

// IsEmpty returns true if the slice is empty. It is safe for concurrent/parallel use.
func (s *Slice[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Copy returns a copy of the slice.
func (s *Slice[T]) Copy() []T {
	return append(make([]T, 0, len(s.items)), s.items...)
}

// Change changes the value for the provided key using provided function.
func (s *Slice[T]) Change(index int, f func(T) T) {
	s.items[index] = f(s.items[index])
}

// Transform transforms all values of the slice using provided function.
func (s *Slice[T]) Transform(f func(T) T) {
	for i, v := range s.items {
		s.items[i] = f(v)
	}
}

// Range calls the provided function for each element in the slice.
func (s *Slice[T]) Range(f func(T) bool) bool {
	for _, v := range s.items {
		if !f(v) {
			return false
		}
	}
	return true
}

// Raw returns the underlying slice.
func (s *Slice[T]) Raw() []T {
	return s.items
}

// Iter returns an iterator over the slice values.
func (s *Slice[T]) Iter() iter.Seq[T] {
	return slices.Values(s.items)
}

// Iter2 returns an iterator over the slice values and their indexes.
func (s *Slice[T]) Iter2() iter.Seq2[int, T] {
	return slices.All(s.items)
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

// NewSafeSliceFromItems returns a new [SafeSlice] with the provided items.
func NewSafeSliceFromItems[T any](data ...T) *SafeSlice[T] {
	return &SafeSlice[T]{
		items: data,
	}
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

// AddFront adds a new elements to the front of the slice.
func (s *SafeSlice[T]) AddFront(v ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = append(v, s.items...)
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
func (s *SafeSlice[T]) Copy() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append(make([]T, 0, len(s.items)), s.items...)
}

// Change changes the value for the provided key using provided function.
func (s *SafeSlice[T]) Change(index int, f func(T) T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[index] = f(s.items[index])
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
func (s *SafeSlice[T]) Range(f func(T) bool) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, v := range s.items {
		if !f(v) {
			return false
		}
	}
	return true
}

// Raw returns the underlying slice.
func (s *SafeSlice[T]) Raw() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.items
}

// Iter returns an iterator over the slice values.
// It is safe for concurrent/parallel use.
// DON'T USE SAFE SLICE METHOD INSIDE LOOP TO PREVENT FROM DEADLOCK!
func (s *SafeSlice[T]) Iter() iter.Seq[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return slices.Values(s.items)
}

// Iter2 returns an iterator over the slice values and their indexes.
// It is safe for concurrent/parallel use.
// DON'T USE SAFE SLICE METHOD INSIDE LOOP TO PREVENT FROM DEADLOCK!
func (s *SafeSlice[T]) Iter2() iter.Seq2[int, T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return slices.All(s.items)
}

func getSlicesLen[T any](slices ...[]T) int {
	var length int
	for _, slice := range slices {
		length += len(slice)
	}
	return length
}
