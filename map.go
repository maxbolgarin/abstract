package abstract

import (
	"crypto/rand"
	"math/big"
	"sort"
	"strings"
	"sync"

	"github.com/maxbolgarin/lang"
)

// Map is used like a common map.
type Map[K comparable, V any] map[K]V

// NewMap returns a [Map] with an empty map.
func NewMap[K comparable, V any](raw ...map[K]V) Map[K, V] {
	out := make(map[K]V, getMapsLength(raw...))
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
func (m Map[K, V]) Range(f func(K, V) bool) bool {
	for k, v := range m {
		if !f(k, v) {
			return false
		}
	}
	return true
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
func (m *SafeMap[K, V]) Range(f func(K, V) bool) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

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
	Map[K, T]
}

// NewEntityMap returns a new EntityMap from the provided map.
func NewEntityMap[K comparable, T Entity[K]](raw ...map[K]T) EntityMap[K, T] {
	return EntityMap[K, T]{
		Map: NewMap(raw...),
	}
}

// NewEntityMapWithSize returns a new EntityMap with the provided size.
func NewEntityMapWithSize[K comparable, T Entity[K]](size int) EntityMap[K, T] {
	return EntityMap[K, T]{
		Map: NewMapWithSize[K, T](size),
	}
}

// LookupByName returns the value for the provided name.
func (s EntityMap[K, T]) LookupByName(name string) (T, bool) {
	name = strings.ToLower(name)

	for _, h := range s.Map {
		if strings.ToLower(h.GetName()) == name {
			return h, true
		}
	}

	var zero T
	return zero, false
}

// Set sets the value for the provided key. It delete the value if the order is -1.
func (s *EntityMap[K, T]) Set(info T) {
	if info.GetOrder() == -1 {
		s.Delete(info.GetID())
		return
	}
	s.Map[info.GetID()] = info
}

// AllOrdered returns all values in order.
func (s *EntityMap[K, T]) AllOrdered() []T {
	var (
		nOfItems   = len(s.Map)
		out        = make([]T, nOfItems)
		seen       = make([]bool, nOfItems)
		broken     []T
		seenBroken bool
	)

	for _, h := range s.Map {
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

// NextOrder returns the next order.
func (s EntityMap[K, T]) NextOrder() int {
	return len(s.Map)
}

// ChangeOrder changes the order of the values.
func (s *EntityMap[K, T]) ChangeOrder(draft map[K]int) {
	ordered := s.AllOrdered()

	maxOrder := len(draft)
	for _, item := range ordered {
		ord, ok := draft[item.GetID()]
		if !ok {
			ord = maxOrder
			maxOrder++
		}

		item, ok := item.SetOrder(ord).(T)
		if !ok {
			panic("you should use the same Entity type as return value from SetOrder")
		}

		s.Map[item.GetID()] = item
	}
}

// Delete deletes the value for the provided key.
func (s EntityMap[K, T]) Delete(key K) bool {
	toDelete, ok := s.Map[key]
	if !ok {
		return false
	}

	// if Deleted has -1 order
	deleteOrder := toDelete.GetOrder()
	if deleteOrder == -1 {
		delete(s.Map, key)
		return true
	}

	ordered := s.AllOrdered()

	var flag bool
	for i, h := range ordered {
		if i == deleteOrder {
			delete(s.Map, key)
			flag = true
			continue
		}
		if i > deleteOrder {
			h, ok = h.SetOrder(h.GetOrder() - 1).(T)
			if !ok {
				panic("you should use the same Entity type as return value from SetOrder")
			}
		}
		s.Map[h.GetID()] = h
	}

	return flag
}

// SafeEntityMap is a thread-safe map of entities.
// It is safe for concurrent/parallel use.
type SafeEntityMap[K comparable, T Entity[K]] struct {
	*SafeMap[K, T]
}

// NewSafeEntityMap returns a new SafeEntityMap from the provided map.
func NewSafeEntityMap[K comparable, T Entity[K]](raw ...map[K]T) SafeEntityMap[K, T] {
	return SafeEntityMap[K, T]{
		SafeMap: NewSafeMap(raw...),
	}
}

// NewSafeEntityMapWithSize returns a new SafeEntityMap with the provided size.
func NewSafeEntityMapWithSize[K comparable, T Entity[K]](size int) SafeEntityMap[K, T] {
	return SafeEntityMap[K, T]{
		SafeMap: NewSafeMapWithSize[K, T](size),
	}
}

// LookupByName returns the value for the provided name.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) LookupByName(name string) (T, bool) {
	name = strings.ToLower(name)

	s.mu.RLock()
	defer s.mu.RUnlock()

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
// If the key is present and the value has order -1, the value will be deleted.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) Set(info T) {
	if info.GetOrder() == -1 {
		s.Delete(info.GetID())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[info.GetID()] = info
}

// AllOrdered returns all values in the map sorted by their order.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) AllOrdered() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

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
func (s SafeEntityMap[K, T]) NextOrder() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

// ChangeOrder changes the order of the values in the map based on the provided map.
// It is safe for concurrent/parallel use.
func (s *SafeEntityMap[K, T]) ChangeOrder(draft map[K]int) {
	ordered := s.AllOrdered()

	maxOrder := len(draft)

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, item := range ordered {
		ord, ok := draft[item.GetID()]
		if !ok {
			ord = maxOrder
			maxOrder++
		}

		item, ok = item.SetOrder(ord).(T)
		if !ok {
			panic("you should use the same Entity type as return value from SetOrder")
		}

		s.items[item.GetID()] = item
	}
}

// Delete deletes the value for the provided key.
// If the value has order -1, the value will be deleted.
// It is safe for concurrent/parallel use.
func (s SafeEntityMap[K, T]) Delete(key K) bool {
	s.mu.RLock()
	toDelete, ok := s.items[key]
	s.mu.RUnlock()

	if !ok {
		return false
	}

	// if Deleted has -1 order
	deleteOrder := toDelete.GetOrder()
	if deleteOrder == -1 {
		s.mu.Lock()
		defer s.mu.Unlock()

		delete(s.items, key)
		return true
	}

	ordered := s.AllOrdered()

	s.mu.Lock()
	defer s.mu.Unlock()

	var flag bool
	for i, h := range ordered {
		if i == deleteOrder {
			delete(s.items, key)
			flag = true
			continue
		}
		if i > deleteOrder {
			h, ok = h.SetOrder(h.GetOrder() - 1).(T)
			if !ok {
				panic("you should use the same Entity type as return value from SetOrder")
			}
		}
		s.items[h.GetID()] = h
	}

	return flag
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
	if index, ok := m.indexes[key]; ok {
		m.elems[index] = value
	}
	m.indexes[key] = len(m.elems)
	m.elems = append(m.elems, value)
	m.keys = append(m.keys, key)
}

// Get returns the value associated with the key.
func (m *OrderedPairs[K, V]) Get(key K) (res V) {
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
