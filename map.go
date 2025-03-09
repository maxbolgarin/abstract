package abstract

import (
	"crypto/rand"
	"iter"
	"maps"
	"math/big"
	"sort"
	"strings"
	"sync"

	"github.com/maxbolgarin/lang"
)

// Map is used like a common map.
type Map[K comparable, V any] struct {
	items map[K]V
}

// NewMap returns a [Map] with an empty map.
func NewMap[K comparable, V any](raw ...map[K]V) *Map[K, V] {
	out := make(map[K]V, getMapsLength(raw...))
	for _, v := range raw {
		for k, v := range v {
			out[k] = v
		}
	}
	return &Map[K, V]{
		items: out,
	}
}

// NewMapFromPairs returns a [Map] with a map inited using the provided pairs.
func NewMapFromPairs[K comparable, V any](pairs ...any) *Map[K, V] {
	out := make(map[K]V, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out[pairs[i].(K)] = pairs[i+1].(V)
	}
	return &Map[K, V]{
		items: out,
	}
}

// NewMapWithSize returns a [Map] with a map inited using the provided size.
func NewMapWithSize[K comparable, V any](size int) *Map[K, V] {
	return &Map[K, V]{
		items: make(map[K]V, size),
	}
}

// Get returns the value for the provided key or the default type value if the key is not present in the map.
func (m *Map[K, V]) Get(key K) V {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	return m.items[key]
}

// Lookup returns the value for the provided key and true if the key is present in the map, the default value and false otherwise.
func (m *Map[K, V]) Lookup(key K) (V, bool) {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	v, ok := m.items[key]
	return v, ok
}

// Has returns true if the key is present in the map, false otherwise.
func (m *Map[K, V]) Has(key K) bool {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	_, ok := m.items[key]
	return ok
}

// Pop returns the value for the provided key and deletes it from map or default type value if key is not present.
func (m *Map[K, V]) Pop(key K) V {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	val, ok := m.items[key]
	if ok {
		delete(m.items, key)
	}
	return val
}

// Set sets the value to the map.
func (m *Map[K, V]) Set(key K, value V) {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	m.items[key] = value
}

// SetIfNotPresent sets the value to the map if the key is not present,
// returns the old value if the key was set, new value otherwise.
func (m *Map[K, V]) SetIfNotPresent(key K, value V) V {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	if _, ok := m.items[key]; !ok {
		m.items[key] = value
		return value
	}
	return m.items[key]
}

// Swap swaps the values for the provided keys and returns the old value.
func (m *Map[K, V]) Swap(key K, value V) V {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	old := m.items[key]
	m.items[key] = value
	return old
}

// Delete removes keys and associated values from the map, does nothing if the key is not present in the map,
// returns true if the key was deleted
func (m *Map[K, V]) Delete(keys ...K) (deleted bool) {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	for _, key := range keys {
		if _, ok := m.items[key]; ok {
			deleted = true
			delete(m.items, key)
		}
	}
	return deleted
}

// Len returns the length of the map.
func (m *Map[K, V]) Len() int {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	return len(m.items)
}

// IsEmpty returns true if the map is empty. It is safe for concurrent/parallel use.
func (m *Map[K, V]) IsEmpty() bool {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	return len(m.items) == 0
}

// Keys returns a slice of keys of the map.
func (m *Map[K, V]) Keys() []K {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	return lang.Keys(m.items)
}

// Values returns a slice of values of the map.
func (m *Map[K, V]) Values() []V {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	return lang.Values(m.items)
}

// Change changes the value for the provided key using provided function.
func (m *Map[K, V]) Change(key K, f func(K, V) V) {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	m.items[key] = f(key, m.items[key])
}

// Transform transforms all values of the map using provided function.
func (m *Map[K, V]) Transform(f func(K, V) V) {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	for k, v := range m.items {
		m.items[k] = f(k, v)
	}
}

// Range calls the provided function for each key-value pair in the map.
func (m *Map[K, V]) Range(f func(K, V) bool) bool {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	for k, v := range m.items {
		if !f(k, v) {
			return false
		}
	}
	return true
}

// Copy returns another map that is a copy of the underlying map.
func (m *Map[K, V]) Copy() map[K]V {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	return lang.CopyMap(m.items)
}

// Raw returns the underlying map.
func (m *Map[K, V]) Raw() map[K]V {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	return m.items
}

// Clear creates a new map using make without size.
func (m *Map[K, V]) Clear() {
	m.items = make(map[K]V)
}

// IterKeys returns an iterator over the map keys.
func (m *Map[K, V]) IterKeys() iter.Seq[K] {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	return maps.Keys(m.items)
}

// IterValues returns an iterator over the map values.
func (m *Map[K, V]) IterValues() iter.Seq[V] {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	return maps.Values(m.items)
}

// Iter returns an iterator over the map.
func (m *Map[K, V]) Iter() iter.Seq2[K, V] {
	if m.items == nil {
		m.items = make(map[K]V)
	}
	return maps.All(m.items)
}

// SafeMap is used like a common map, but it is protected with RW mutex, so it can be used in many goroutines.
type SafeMap[K comparable, V any] struct {
	items map[K]V
	mu    sync.RWMutex
}

// NewSafeMap returns a new [SafeMap] with empty map.
func NewSafeMap[K comparable, V any](raw ...map[K]V) *SafeMap[K, V] {
	out := &SafeMap[K, V]{
		items: make(map[K]V, getMapsLength(raw...)),
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

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	return m.items[key]
}

// Lookup returns the value for the provided key and true if key is present in the map, default value and false otherwise.
// It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Lookup(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	v, ok := m.items[key]
	return v, ok
}

// Has returns true if key is present in the map, false otherwise. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Has(key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	_, ok := m.items[key]
	return ok
}

// Pop returns the value for the provided key and deletes it from map or default type value if key is not present.
// It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Pop(key K) V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

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

	if m.items == nil {
		m.items = make(map[K]V)
	}

	m.items[key] = value
}

// SetIfNotPresent sets the value to the map if the key is not present,
// returns the old value if the key was set, new value otherwise. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) SetIfNotPresent(key K, value V) V {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.items == nil {
		m.items = make(map[K]V)
	}

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

	if m.items == nil {
		m.items = make(map[K]V)
	}

	old := m.items[key]
	m.items[key] = value
	return old
}

// Delete removes keys and associated values from map, does nothing if key is not present in map,
// returns true if key was deleted. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Delete(keys ...K) (deleted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.items == nil {
		m.items = make(map[K]V)
	}

	for _, key := range keys {
		if _, ok := m.items[key]; ok {
			deleted = true
			delete(m.items, key)
		}
	}

	return deleted
}

// Len returns the length of the map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	return len(m.items)
}

// IsEmpty returns true if the map is empty. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	return len(m.items) == 0
}

// Keys returns a slice of keys of the map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	return lang.Keys(m.items)
}

// Values returns a slice of values of the map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	return lang.Values(m.items)
}

// Change changes the value for the provided key using provided function. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Change(key K, f func(K, V) V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.items == nil {
		m.items = make(map[K]V)
	}

	m.items[key] = f(key, m.items[key])
}

// Update updates the map using provided function. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Transform(upd func(K, V) V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.items == nil {
		m.items = make(map[K]V)
	}

	for k, v := range m.items {
		m.items[k] = upd(k, v)
	}
}

// Range calls the provided function for each key-value pair in the map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Range(f func(K, V) bool) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	for k, v := range m.items {
		if !f(k, v) {
			return false
		}
	}
	return true
}

// Copy returns a new map that is a copy of the underlying map. It is safe for concurrent/parallel use.
func (m *SafeMap[K, V]) Copy() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

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

	if m.items == nil {
		m.items = make(map[K]V)
	}

	m.items = lang.CopyMap(raw)
}

// Raw returns the underlying map.
func (m *SafeMap[K, V]) Raw() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	return m.items
}

// IterValues returns an iterator over the map values.
// It is safe for concurrent/parallel use.
// DON'T USE SAFE MAP METHOD INSIDE LOOP TO PREVENT FROM DEADLOCK!
func (m *SafeMap[K, V]) IterValues() iter.Seq[V] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	return maps.Values(m.items)
}

// IterKeys returns an iterator over the map keys.
// It is safe for concurrent/parallel use.
// DON'T USE SAFE MAP METHOD INSIDE LOOP TO PREVENT FROM DEADLOCK!
func (m *SafeMap[K, V]) IterKeys() iter.Seq[K] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	return maps.Keys(m.items)
}

// Iter returns an iterator over the map.
// It is safe for concurrent/parallel use.
// DON'T USE SAFE MAP METHOD INSIDE LOOP TO PREVENT FROM DEADLOCK!
func (m *SafeMap[K, V]) Iter() iter.Seq2[K, V] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.items == nil {
		m.mu.RUnlock()
		m.mu.Lock()
		m.items = make(map[K]V)
		m.mu.Unlock()
		m.mu.RLock()
	}

	return maps.All(m.items)
}

func getMapsLength[K comparable, V any](maps ...map[K]V) int {
	length := 0
	for _, m := range maps {
		length += len(m)
	}
	return length
}

// Entity is an interface for an object that has an ID, a name, and an order.
type Entity[K comparable] interface {
	GetID() K
	GetName() string
	GetOrder() int
	SetOrder(int) Entity[K]
}

// EntityMap is a map of entities. It has all methods of Map with some new ones.
// It is not safe for concurrent/parallel, use [SafeEntityMap] if you need it.
type EntityMap[K comparable, T Entity[K]] struct {
	*Map[K, T]
}

// NewEntityMap returns a new EntityMap from the provided map.
func NewEntityMap[K comparable, T Entity[K]](raw ...map[K]T) *EntityMap[K, T] {
	return &EntityMap[K, T]{
		Map: NewMap(raw...),
	}
}

// NewEntityMapWithSize returns a new EntityMap with the provided size.
func NewEntityMapWithSize[K comparable, T Entity[K]](size int) *EntityMap[K, T] {
	return &EntityMap[K, T]{
		Map: NewMapWithSize[K, T](size),
	}
}

// LookupByName returns the value for the provided name.
// It is not case-sensetive according to name.
func (s *EntityMap[K, T]) LookupByName(name string) (T, bool) {
	if s.Map.items == nil {
		s.Map.items = make(map[K]T)
	}

	name = strings.ToLower(name)

	for _, h := range s.Map.items {
		if strings.ToLower(h.GetName()) == name {
			return h, true
		}
	}

	var zero T
	return zero, false
}

// Set sets the value for the provided key.
// It sets last order to the entity's order, so it adds to the end of the list.
// It sets the same order of existing entity in case of conflict.
// It returns the order of the entity.
func (s *EntityMap[K, T]) Set(info T) int {
	if s.Map.items == nil {
		s.Map.items = make(map[K]T)
	}

	id := info.GetID()
	old, ok := s.Map.items[id]
	if ok {
		info = info.SetOrder(old.GetOrder()).(T)
	} else {
		info = info.SetOrder(len(s.Map.items)).(T)
	}
	s.Map.items[id] = info

	return info.GetOrder()
}

// SetManualOrder sets the value for the provided key.
// Better to use [EntityMap.Set] to prevent from order errors.
// It returns the order of the entity.
func (s *EntityMap[K, T]) SetManualOrder(info T) int {
	s.Map.items[info.GetID()] = info
	return info.GetOrder()
}

// AllOrdered returns all values in order.
func (s *EntityMap[K, T]) AllOrdered() []T {
	if s.Map.items == nil {
		s.Map.items = make(map[K]T)
	}

	var (
		nOfItems   = len(s.Map.items)
		out        = make([]T, nOfItems)
		seen       = make([]bool, nOfItems)
		broken     []T
		seenBroken bool
	)

	for _, h := range s.Map.items {
		order := h.GetOrder()
		if order < 0 || order >= nOfItems || seen[order] {
			seenBroken = true
			broken = append(broken, h)
			continue
		}
		out[order] = h
		seen[order] = true
	}
	if seenBroken {
		sort.Slice(broken, func(i, j int) bool {
			orderI := broken[i].GetOrder()
			orderJ := broken[j].GetOrder()
			if orderI < 0 || orderJ < 0 {
				return orderI < orderJ
			}
			return orderI < orderJ
		})
		var i int
		for ind, isFound := range seen {
			if isFound {
				continue
			}
			out[ind] = broken[i]
			i++
		}
	}

	return out
}

// NextOrder returns the next order.
func (s *EntityMap[K, T]) NextOrder() int {
	if s.Map.items == nil {
		s.Map.items = make(map[K]T)
	}

	return len(s.Map.items)
}

// ChangeOrder changes the order of the values.
func (s *EntityMap[K, T]) ChangeOrder(draft map[K]int) {
	if s.Map.items == nil {
		s.Map.items = make(map[K]T)
	}

	ordered := s.AllOrdered()

	maxOrder := len(draft)
	for _, item := range ordered {
		ord, ok := draft[item.GetID()]
		if !ok {
			ord = maxOrder
			maxOrder++
		}
		s.Map.items[item.GetID()] = item.SetOrder(ord).(T)
	}
}

// Delete deletes values for the provided keys.
// It reorders all remaining values.
func (s *EntityMap[K, T]) Delete(keys ...K) (deleted bool) {
	if s.Map.items == nil {
		s.Map.items = make(map[K]T)
	}

	for _, key := range keys {
		toDelete, ok := s.Map.items[key]
		if !ok {
			continue
		}

		deleteOrder := toDelete.GetOrder()
		ordered := s.AllOrdered()

		for i, h := range ordered {
			if i == deleteOrder {
				delete(s.Map.items, key)
				deleted = true
				continue
			}
			if i > deleteOrder {
				h = h.SetOrder(h.GetOrder() - 1).(T)
			}
			s.Map.items[h.GetID()] = h
		}
	}
	return deleted
}

// SafeEntityMap is a thread-safe map of entities.
// It is safe for concurrent/parallel use.
type SafeEntityMap[K comparable, T Entity[K]] struct {
	*SafeMap[K, T]
}

// NewSafeEntityMap returns a new SafeEntityMap from the provided map.
func NewSafeEntityMap[K comparable, T Entity[K]](raw ...map[K]T) *SafeEntityMap[K, T] {
	return &SafeEntityMap[K, T]{
		SafeMap: NewSafeMap(raw...),
	}
}

// NewSafeEntityMapWithSize returns a new SafeEntityMap with the provided size.
func NewSafeEntityMapWithSize[K comparable, T Entity[K]](size int) *SafeEntityMap[K, T] {
	return &SafeEntityMap[K, T]{
		SafeMap: NewSafeMapWithSize[K, T](size),
	}
}

// LookupByName returns the value for the provided name.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) LookupByName(name string) (T, bool) {
	name = strings.ToLower(name)

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.items == nil {
		s.mu.RUnlock()
		s.mu.Lock()
		s.items = make(map[K]T)
		s.mu.Unlock()
		s.mu.RLock()
	}

	for _, h := range s.items {
		if strings.ToLower(h.GetName()) == name {
			return h, true
		}
	}

	var zero T
	return zero, false
}

// Set sets the value for the provided key.
// If the key is not present in the map, it will be added.
// It sets last order to the entity's order.
// It sets the same order of existing entity in case of conflict.
// It returns the order of the entity.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) Set(info T) int {
	id := info.GetID()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.items == nil {
		s.items = make(map[K]T)
	}

	old, ok := s.items[id]
	if ok {
		info = info.SetOrder(old.GetOrder()).(T)
	} else {
		info = info.SetOrder(len(s.items)).(T)
	}
	s.items[id] = info

	return info.GetOrder()
}

// SetManualOrder sets the value for the provided key.
// Better to use [SafeEntityMap.Set] to prevent from order errors.
// It returns the order of the entity.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) SetManualOrder(info T) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.items == nil {
		s.items = make(map[K]T)
	}

	s.items[info.GetID()] = info

	return info.GetOrder()
}

// AllOrdered returns all values in the map sorted by their order.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) AllOrdered() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.items == nil {
		s.mu.RUnlock()
		s.mu.Lock()
		s.items = make(map[K]T)
		s.mu.Unlock()
		s.mu.RLock()
	}

	return s.allOrderedNoLock()
}

func (s *SafeEntityMap[K, T]) allOrderedNoLock() []T {
	var (
		nOfItems   = len(s.items)
		out        = make([]T, nOfItems)
		seen       = make([]bool, nOfItems)
		broken     []T
		seenBroken bool
	)

	for _, h := range s.items {
		order := h.GetOrder()
		if order < 0 || order >= nOfItems {
			seenBroken = true
			broken = append(broken, h)
			continue
		}
		out[order] = h
		seen[order] = true
	}
	if seenBroken {
		sort.Slice(broken, func(i, j int) bool {
			orderI := broken[i].GetOrder()
			orderJ := broken[j].GetOrder()
			if orderI < 0 || orderJ < 0 {
				return orderI < orderJ
			}
			return orderI < orderJ
		})
		var i int
		for ind, isFound := range seen {
			if isFound {
				continue
			}
			out[ind] = broken[i]
			i++
		}
	}

	return out
}

// NextOrder returns the next order number.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) NextOrder() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.items == nil {
		s.mu.RUnlock()
		s.mu.Lock()
		s.items = make(map[K]T)
		s.mu.Unlock()
		s.mu.RLock()
	}

	return len(s.items)
}

// ChangeOrder changes the order of the values in the map based on the provided map.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) ChangeOrder(draft map[K]int) {
	ordered := s.AllOrdered()

	maxOrder := len(draft)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.items == nil {
		s.items = make(map[K]T)
	}

	for _, item := range ordered {
		ord, ok := draft[item.GetID()]
		if !ok {
			ord = maxOrder
			maxOrder++
		}
		s.items[item.GetID()] = item.SetOrder(ord).(T)
	}
}

// Delete deletes values for the provided keys.
// It reorders all remaining values.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) Delete(keys ...K) (deleted bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.items == nil {
		s.items = make(map[K]T)
	}

	for _, key := range keys {
		toDelete, ok := s.items[key]
		if !ok {
			continue
		}

		deleteOrder := toDelete.GetOrder()
		ordered := s.allOrderedNoLock()

		for i, h := range ordered {
			if i == deleteOrder {
				delete(s.items, key)
				deleted = true
				continue
			}
			if i > deleteOrder {
				h = h.SetOrder(h.GetOrder() - 1).(T)
			}
			s.items[h.GetID()] = h
		}
	}

	return deleted
}

// OrderedPairs is a data structure that behaves like a map but remembers
// the order in which the items were added. It is also possible to get a random
// value or key from the structure. It allows duplicate keys.
// It is NOT safe for concurrent/parallel use.
//
// The type parameter K must implement the Ordered interface.
type OrderedPairs[K Ordered, V any] struct {
	elems   []V
	keys    []K
	indexes map[K]int
}

// NewOrderedPairs creates a new OrderedPairs from the provided pairs. It allows duplicate keys.
func NewOrderedPairs[K Ordered, V any](pairs ...any) *OrderedPairs[K, V] {
	if len(pairs)%2 == 1 {
		pairs = pairs[:len(pairs)-1]
	}
	m := &OrderedPairs[K, V]{
		elems:   make([]V, 0, len(pairs)/2),
		keys:    make([]K, 0, len(pairs)/2),
		indexes: make(map[K]int, len(pairs)/2),
	}
	for i := 0; i < len(pairs)-1; i += 2 {
		key := pairs[i].(K)
		value := pairs[i+1].(V)
		m.Add(key, value)
	}
	return m
}

// Add adds a key-value pair to the structure. It allows duplicate keys.
func (m *OrderedPairs[K, V]) Add(key K, value V) {
	if m.indexes == nil {
		m.indexes = make(map[K]int)
	}
	if index, ok := m.indexes[key]; ok {
		m.elems[index] = value
	}
	m.indexes[key] = len(m.elems)
	m.elems = append(m.elems, value)
	m.keys = append(m.keys, key)
}

// Get returns the value associated with the key.
func (m *OrderedPairs[K, V]) Get(key K) (res V) {
	if m.indexes == nil {
		m.indexes = make(map[K]int)
	}
	if index, ok := m.indexes[key]; ok {
		return m.elems[index]
	}
	return res
}

// Keys returns a slice of all keys in the structure.
func (m *OrderedPairs[K, V]) Keys() []K {
	return m.keys
}

// Rand returns a random value from the structure.
func (m *OrderedPairs[K, V]) Rand() V {
	if len(m.elems) == 0 {
		return *new(V)
	}
	return m.elems[getRand(len(m.elems))]
}

// RandKey returns a random key from the structure.
func (m *OrderedPairs[K, V]) RandKey() K {
	if len(m.keys) == 0 {
		return *new(K)
	}
	return m.keys[getRand(len(m.keys))]
}

func getRand(max int) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return nBig.Int64()
}

// SafeOrderedPairs is a thread-safe variant of the OrderedPairs type.
// It uses a RW mutex to protect the underlying structure.
//
// The type parameter K must implement the Ordered interface.
type SafeOrderedPairs[K Ordered, V any] struct {
	*OrderedPairs[K, V]
	mu sync.RWMutex
}

// NewSafeOrderedPairs returns a new SafeOrderedPairs from the provided pairs.
// It is a thread-safe variant of the NewOrderedPairs function.
func NewSafeOrderedPairs[K Ordered, V any](pairs ...any) *SafeOrderedPairs[K, V] {
	return &SafeOrderedPairs[K, V]{
		OrderedPairs: NewOrderedPairs[K, V](pairs...),
	}
}

// Add adds a key-value pair to the structure. It allows duplicate keys.
// It is a thread-safe variant of the Add method.
func (s *SafeOrderedPairs[K, V]) Add(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.OrderedPairs.Add(key, value)
}

// Get returns the value associated with the key.
// It is a thread-safe variant of the Get method.
func (s *SafeOrderedPairs[K, V]) Get(key K) (res V) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.OrderedPairs.Get(key)
}

// Rand returns a random value from the structure.
// It is a thread-safe variant of the Rand method.
func (s *SafeOrderedPairs[K, V]) Rand() V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.OrderedPairs.Rand()
}

// RandKey returns a random key from the structure.
// It is a thread-safe variant of the RandKey method.
func (s *SafeOrderedPairs[K, V]) RandKey() K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.OrderedPairs.RandKey()
}
