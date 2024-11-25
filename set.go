package abstract

import (
	"sync"

	"github.com/maxbolgarin/lang"
)

// Set represents a data structure that behaves like a common map but is more lightweight.
// It is used to store unique keys without associated values.
type Set[K comparable] struct {
	items map[K]struct{}
}

// NewSet returns a [Set] with an empty map.
func NewSet[K comparable](data ...[]K) *Set[K] {
	out := &Set[K]{
		items: make(map[K]struct{}, getSlicesLen(data...)),
	}
	for _, v := range data {
		for _, v := range v {
			out.items[v] = struct{}{}
		}
	}
	return out
}

// NewSetWithSize returns a [Set] with a map inited using the provided size.
func NewSetWithSize[K comparable](size int) *Set[K] {
	return &Set[K]{
		items: make(map[K]struct{}, size),
	}
}

// Add adds keys to the set.
func (m *Set[K]) Add(key ...K) {
	for _, v := range key {
		m.items[v] = struct{}{}
	}
}

// Has returns true if the key is present in the set, false otherwise.
func (m *Set[K]) Has(key K) bool {
	_, ok := m.items[key]
	return ok
}

// Delete removes the key from the set, does nothing if the key is not present in the set.
func (m *Set[K]) Delete(key K) {
	delete(m.items, key)
}

// Len returns the length of the set.
func (m *Set[K]) Len() int {
	return len(m.items)
}

// IsEmpty returns true if the set is empty. It is safe for concurrent/parallel use.
func (m *Set[K]) IsEmpty() bool {
	return len(m.items) == 0
}

// Keys returns a slice of keys of the set.
func (m *Set[K]) Values() []K {
	return lang.Keys(m.items)
}

// Clear creates a new map using make without size.
func (m *Set[K]) Clear() {
	m.items = make(map[K]struct{})
}

// Transform transforms all values of the set using provided function.
func (m *Set[K]) Transform(f func(K) K) {
	raw := make(map[K]struct{}, len(m.items))
	for k := range m.items {
		raw[f(k)] = struct{}{}
	}
	m.items = raw
}

// Range calls the provided function for each key in the set.
func (m *Set[K]) Range(f func(K) bool) {
	for k := range m.items {
		if !f(k) {
			return
		}
	}
}

// SafeSet is used like a set, but it is protected with RW mutex, so it can be used in many goroutines.
type SafeSet[K comparable] struct {
	set map[K]struct{}
	mu  sync.RWMutex
}

// NewSafeSet returns a new [SafeSet] with empty set.
func NewSafeSet[K comparable](data ...[]K) *SafeSet[K] {
	out := &SafeSet[K]{
		set: make(map[K]struct{}, getSlicesLen(data...)),
	}
	for _, v := range data {
		for _, v := range v {
			out.set[v] = struct{}{}
		}
	}
	return out
}

// NewSafeSetWithSize returns a new [SafeSet] with empty set.
func NewSafeSetWithSize[K comparable](size int) *SafeSet[K] {
	return &SafeSet[K]{
		set: make(map[K]struct{}, size),
	}
}

// Add adds keys to the set. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Add(key ...K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, v := range key {
		m.set[v] = struct{}{}
	}
}

// Has returns true if key is present in the set, false otherwise. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Has(key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.set[key]
	return ok
}

// Delete removes a key from the set. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.set, key)
}

// Len returns the number of keys in set. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.set)
}

// Empty returns true if the set is empty. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.set) == 0
}

// Values returns a slice of values of the set. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Values() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return lang.Keys(m.set)
}

// Clear removes all keys from the set. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.set = make(map[K]struct{})
}

// Transform transforms all values of the set using provided function. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Transform(f func(K) K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	raw := make(map[K]struct{}, len(m.set))
	for k := range m.set {
		raw[f(k)] = struct{}{}
	}
	m.set = raw
}

// Range calls the provided function for each key in the set. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Range(f func(K) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k := range m.set {
		if !f(k) {
			return
		}
	}
}
