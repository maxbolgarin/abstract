package abstract

import (
	"iter"
	"maps"
	"sync"

	"github.com/maxbolgarin/lang"
)

// Set represents a data structure that behaves like a common map but is more lightweight.
// It is used to store unique keys without associated values.
type Set[K comparable] struct {
	items map[K]struct{}
}

// NewSet returns a [Set] inited using the provided data.
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

// NewSetFromItems returns a [Set] inited using the provided data.
func NewSetFromItems[K comparable](data ...K) *Set[K] {
	out := &Set[K]{
		items: make(map[K]struct{}, len(data)),
	}
	for _, v := range data {
		out.items[v] = struct{}{}
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

// Delete removes the keys from the set, does nothing if the key is not present in the set.
func (m *Set[K]) Delete(keys ...K) (deleted bool) {
	for _, key := range keys {
		if _, ok := m.items[key]; ok {
			delete(m.items, key)
			deleted = true
		}
	}
	return deleted
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
func (m *Set[K]) Range(f func(K) bool) bool {
	for k := range m.items {
		if !f(k) {
			return false
		}
	}
	return true
}

// Raw returns the underlying map.
func (m *Set[K]) Raw() map[K]struct{} {
	return m.items
}

// Iter returns a channel that yields each key in the set.
func (m *Set[K]) Iter() iter.Seq[K] {
	return maps.Keys(m.items)
}

// Copy returns a copy of the set.
func (m *Set[K]) Copy() map[K]struct{} {
	out := make(map[K]struct{}, len(m.items))
	maps.Copy(out, m.items)
	return out
}

// Union returns a new set with the union of the current set and the provided set.
func (m *Set[K]) Union(set map[K]struct{}) *Set[K] {
	out := NewSet[K]()
	for k := range m.items {
		out.items[k] = struct{}{}
	}
	for k := range set {
		out.items[k] = struct{}{}
	}
	return out
}

// Intersection returns a new set with the intersection of the current set and the provided set.
func (m *Set[K]) Intersection(set map[K]struct{}) *Set[K] {
	out := NewSet[K]()
	for k := range m.items {
		if _, ok := set[k]; ok {
			out.items[k] = struct{}{}
		}
	}
	return out
}

// Difference returns a new set with the difference of the current set and the provided set.
func (m *Set[K]) Difference(set map[K]struct{}) *Set[K] {
	out := NewSet[K]()
	for k := range m.items {
		if _, ok := set[k]; !ok {
			out.items[k] = struct{}{}
		}
	}
	return out
}

// SymmetricDifference returns a new set with the symmetric difference of the current set and the provided set.
func (m *Set[K]) SymmetricDifference(set map[K]struct{}) *Set[K] {
	out := NewSet[K]()
	for k := range m.items {
		if _, ok := set[k]; !ok {
			out.items[k] = struct{}{}
		}
	}
	for k := range set {
		if _, ok := m.items[k]; !ok {
			out.items[k] = struct{}{}
		}
	}
	return out
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

// NewSafeSetFromItems returns a new [SafeSet] with empty set.
func NewSafeSetFromItems[K comparable](data ...K) *SafeSet[K] {
	out := &SafeSet[K]{
		set: make(map[K]struct{}, len(data)),
	}
	for _, v := range data {
		out.set[v] = struct{}{}
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

// Delete removes keys from the set. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Delete(keys ...K) (deleted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range keys {
		if _, ok := m.set[key]; ok {
			delete(m.set, key)
			deleted = true
		}
	}
	return deleted
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
func (m *SafeSet[K]) Range(f func(K) bool) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k := range m.set {
		if !f(k) {
			return false
		}
	}
	return true
}

// Raw returns the underlying map.
func (m *SafeSet[K]) Raw() map[K]struct{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.set
}

// Iter returns a channel that yields each key in the set.
// It is safe for concurrent/parallel use.
// DON'T USE SAFE SET METHOD INSIDE LOOP TO PREVENT FROM DEADLOCK!
func (m *SafeSet[K]) Iter() iter.Seq[K] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return maps.Keys(m.set)
}

// Copy returns a copy of the set. It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Copy() map[K]struct{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(map[K]struct{}, len(m.set))
	maps.Copy(out, m.set)

	return out
}

// Union returns a new set with the union of the current set and the provided set.
// It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Union(set map[K]struct{}) *Set[K] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := NewSet[K]()
	for k := range m.set {
		out.items[k] = struct{}{}
	}
	for k := range set {
		out.items[k] = struct{}{}
	}
	return out
}

// Intersection returns a new set with the intersection of the current set and the provided set.
// It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Intersection(set map[K]struct{}) *Set[K] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := NewSet[K]()
	for k := range m.set {
		if _, ok := set[k]; ok {
			out.items[k] = struct{}{}
		}
	}
	return out
}

// Difference returns a new set with the difference of the current set and the provided set.
// It is safe for concurrent/parallel use.
func (m *SafeSet[K]) Difference(set map[K]struct{}) *Set[K] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := NewSet[K]()
	for k := range m.set {
		if _, ok := set[k]; !ok {
			out.items[k] = struct{}{}
		}
	}
	return out
}

// SymmetricDifference returns a new set with the symmetric difference of the current set and the provided set.
// It is safe for concurrent/parallel use.
func (m *SafeSet[K]) SymmetricDifference(set map[K]struct{}) *Set[K] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := NewSet[K]()
	for k := range m.set {
		if _, ok := set[k]; !ok {
			out.items[k] = struct{}{}
		}
	}
	for k := range set {
		if !m.Has(k) {
			out.items[k] = struct{}{}
		}
	}
	return out
}
