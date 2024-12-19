// Package abstract provides abstract objects to make work easier and simplify code.
package abstract

import (
	"math"
	"strconv"
	"sync"
)

// Signed is a constraint that permits any signed integer type.
// If future releases of Go add new predeclared signed integer types,
// this constraint will be modified to include them.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
// If future releases of Go add new predeclared unsigned integer types,
// this constraint will be modified to include them.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type.
// If future releases of Go add new predeclared integer types,
// this constraint will be modified to include them.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
// If future releases of Go add new predeclared floating-point types,
// this constraint will be modified to include them.
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
// If future releases of Go add new predeclared complex numeric types,
// this constraint will be modified to include them.
type Complex interface {
	~complex64 | ~complex128
}

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
// If future releases of Go add new ordered types,
// this constraint will be modified to include them.
type Ordered interface {
	Integer | Float | ~string
}

// Number is a constraint that permits any numeric type.
type Number interface {
	Integer | Float
}

// Orderer is a struct that holds an order of comparable items.
type Orderer[T comparable] struct {
	order         map[T]int
	applyCallback func(order map[T]int)
	mu            sync.Mutex
}

// NewOrderer returns a new orderer.
func NewOrderer[T comparable](f func(order map[T]int)) *Orderer[T] {
	return &Orderer[T]{
		order:         make(map[T]int),
		applyCallback: f,
	}
}

// Add adds an item to the orderer.
func (m *Orderer[T]) Add(id T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.order[id] = len(m.order)
}

// Get returns the current order.
func (m *Orderer[T]) Get() map[T]int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.order
}

// Apply applies the order using the callback.
func (m *Orderer[T]) Apply() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.order) > 0 {
		m.applyCallback(m.order)
	}

	m.order = make(map[T]int)
}

// Clear clears the orderer.
func (m *Orderer[T]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.order = make(map[T]int)
}

// Memorizer is a struct that holds a single item.
type Memorizer[T any] struct {
	item  T
	isSet bool
	mu    sync.Mutex
}

// NewMemorizer returns a new Memorizer.
func NewMemorizer[T any]() *Memorizer[T] {
	return &Memorizer[T]{}
}

// Set sets the value to the Memorizer.
func (m *Memorizer[T]) Set(c T) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.isSet = true
	m.item = c
}

// Get returns the value from the Memorizer.
func (m *Memorizer[T]) Get() (T, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.item, m.isSet
}

// Pop returns the value for the provided key and deletes it from map or default type value if key is not present.
func (m *Memorizer[T]) Pop() (T, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isSet {
		var zero T
		return zero, false
	}

	m.isSet = false
	return m.item, true
}

// Itoa converts an integer to a string.
func Itoa[T Number](i T) string {
	return strconv.Itoa(int(i))
}

// Atoi converts a string to an integer.
func Atoi[T Number](s string) (T, error) {
	i, err := strconv.Atoi(s)
	return T(i), err
}

// Round returns the nearest integer, rounding half away from zero.
func Round[T Number](f T) T {
	return T(math.Round(float64(f)))
}

// Min returns the minimum value from the provided values.
func Min[T Number](xs ...T) T {
	var min T
	if len(xs) == 0 {
		return min
	}
	min = xs[0]
	for _, x := range xs {
		if x < min {
			min = x
		}
	}
	return min
}

// Max returns the maximum value from the provided values.
func Max[T Number](xs ...T) T {
	var max T
	if len(xs) == 0 {
		return max
	}
	max = xs[0]
	for _, x := range xs {
		if x > max {
			max = x
		}
	}
	return max
}

// Abs returns the absolute value of the provided value.
func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Pow returns the value raised to the provided power.
func Pow[T1, T2 Number](x T1, y T2) T1 {
	return T1(math.Pow(float64(x), float64(y)))
}
