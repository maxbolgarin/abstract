package abstract

import (
	"sync"

	"github.com/maxbolgarin/lang"
)

// Map is used like a common map.
type Map[K comparable, V any] map[K]V

// NewMap returns a [Map] with an empty map.
func NewMap[K comparable, V any](raw ...map[K]V) Map[K, V] {
	out := make(map[K]V)
	for _, v := range raw {
		for k, v := range v {
			out[k] = v
		}
	}
	return out
}

// NewMapFromPairs returns a [Map] with a map inited using the provided pairs.
func NewMapFromPairs[K comparable, V any](pairs ...any) Map[K, V] {
	out := make(map[K]V, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out[pairs[i].(K)] = pairs[i+1].(V)
	}
	return out
}

// NewMapWithSize returns a [Map] with a map inited using the provided size.
func NewMapWithSize[K comparable, V any](size int) Map[K, V] {
	return make(map[K]V, size)
}

// Get returns the value for the provided key or the default type value if the key is not present in the map.
func (m Map[K, V]) Get(key K) V {
	return m[key]
}

// Lookup returns the value for the provided key and true if the key is present in the map, the default value and false otherwise.
func (m Map[K, V]) Lookup(key K) (V, bool) {
	v, ok := m[key]
	return v, ok
}

// Has returns true if the key is present in the map, false otherwise.
func (m Map[K, V]) Has(key K) bool {
	_, ok := m[key]
	return ok
}

// Pop returns the value for the provided key and deletes it from map or default type value if key is not present.
func (m Map[K, V]) Pop(key K) V {
	val, ok := m[key]
	if ok {
		delete(m, key)
	}
	return val
}

// Set sets the value to the map.
func (m Map[K, V]) Set(key K, value V) {
	m[key] = value
}

// SetIfNotPresent sets the value to the map if the key is not present,
// returns the old value if the key was set, new value otherwise.
func (m Map[K, V]) SetIfNotPresent(key K, value V) V {
	if _, ok := m[key]; !ok {
		m[key] = value
		return value
	}
	return m[key]
}

// Swap swaps the values for the provided keys and returns the old value.
func (m Map[K, V]) Swap(key K, value V) V {
	old := m[key]
	m[key] = value
	return old
}

// Delete removes the key and associated value from the map, does nothing if the key is not present in the map,
// returns true if the key was deleted
func (m Map[K, V]) Delete(key K) bool {
	_, ok := m[key]
	if !ok {
		return false
	}
	delete(m, key)

	return true
}

// Len returns the length of the map.
func (m Map[K, V]) Len() int {
	return len(m)
}

// IsEmpty returns true if the map is empty. It is safe for concurrent/parallel use.
func (m Map[K, V]) IsEmpty() bool {
	return len(m) == 0
}

// Keys returns a slice of keys of the map.
func (m Map[K, V]) Keys() []K {
	return lang.Keys(m)
}

// Values returns a slice of values of the map.
func (m Map[K, V]) Values() []V {
	return lang.Values(m)
}

// Transform transforms all values of the map using provided function.
func (m Map[K, V]) Transform(f func(K, V) V) {
	for k, v := range m {
		m[k] = f(k, v)
	}
}

// Range calls the provided function for each key-value pair in the map.
func (m Map[K, V]) Range(f func(K, V) bool) {
	for k, v := range m {
		if !f(k, v) {
			return
		}
	}
}

// Copy returns another map that is a copy of the underlying map.
func (m Map[K, V]) Copy() map[K]V {
	return lang.CopyMap(m)
}

// SafeMap is used like a common map, but it is protected with RW mutex, so it can be used in many goroutines.
type SafeMap[K comparable, V any] struct {
	items map[K]V
	mu    sync.RWMutex
}

// NewSafeMap returns a new [SafeMap] with empty map.
func NewSafeMap[K comparable, V any](raw ...map[K]V) *SafeMap[K, V] {
	out := &SafeMap[K, V]{
		items: make(map[K]V),
	}
	for _, v := range raw {
		for k, v := range v {
			out.items[k] = v
		}
	}
	return out
}

// NewSafeMapFromPairs returns a [SafeMap] with a map inited using the provided pairs.
func NewSafeMapFromPairs[K comparable, V any](pairs ...any) *SafeMap[K, V] {
	out := &SafeMap[K, V]{
		items: make(map[K]V, len(pairs)/2),
	}
	for i := 0; i < len(pairs); i += 2 {
		out.items[pairs[i].(K)] = pairs[i+1].(V)
	}
	return out
}

// NewSafeMapWithSize returns a new [SafeMap] with map inited using provided size.
func NewSafeMapWithSize[K comparable, V any](size int) *SafeMap[K, V] {
	return &SafeMap[K, V]{
		items: make(map[K]V, size),
	}
}

// Get returns the value for the provided key or default type value if key is not present in the map.
// It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Get(key K) V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.items[key]
}

// Lookup returns the value for the provided key and true if key is present in the map, default value and false otherwise.
// It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Lookup(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.items[key]
	return v, ok
}

// Has returns true if key is present in the map, false otherwise. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Has(key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.items[key]
	return ok
}

// Pop returns the value for the provided key and deletes it from map or default type value if key is not present.
// It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Pop(key K) V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.items[key]
	if ok {
		delete(m.items, key)
	}
	return val
}

// Set sets the value to the map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = value
}

// SetIfNotPresent sets the value to the map if the key is not present,
// returns the old value if the key was set, new value otherwise. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) SetIfNotPresent(key K, value V) V {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.items[key]; !ok {
		m.items[key] = value
		return value
	}
	return m.items[key]
}

// Swap swaps the values for the provided keys and returns the old value. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Swap(key K, value V) V {
	m.mu.Lock()
	defer m.mu.Unlock()

	old := m.items[key]
	m.items[key] = value
	return old
}

// Delete removes key and associated value from map, does nothing if key is not present in map,
// returns true if key was deleted. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Delete(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.items[key]
	if !ok {
		return false
	}
	delete(m.items, key)

	return true
}

// Len returns the length of the map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.items)
}

// IsEmpty returns true if the map is empty. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.items) == 0
}

// Keys returns a slice of keys of the map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return lang.Keys(m.items)
}

// Values returns a slice of values of the map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return lang.Values(m.items)
}

// Update updates the map using provided function. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Transform(upd func(K, V) V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range m.items {
		m.items[k] = upd(k, v)
	}
}

// Range calls the provided function for each key-value pair in the map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Range(f func(K, V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.items {
		if !f(k, v) {
			return
		}
	}
}

// Copy returns a new map that is a copy of the underlying map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Copy() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return lang.CopyMap(m.items)
}

// Clear creates a new map using make without size.
func (m *SafeMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[K]V)
}

// Refill creates a new map with values from the provided one.
func (m *SafeMap[K, V]) Refill(raw map[K]V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = lang.CopyMap(raw)
}
