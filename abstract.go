// Package abstract provides a comprehensive collection of abstract data structures,
// utilities, and helper functions designed to simplify common programming tasks in Go.
//
// The package includes:
//
// - Generic data structures: Map, Set, Stack, Slice, LinkedList with both regular and thread-safe variants
// - Cryptographic utilities: AES encryption/decryption, HMAC generation, ECDSA signing/verification
// - Concurrency helpers: Future/Promise patterns, Worker pools, Rate limiting
// - Random generation: Secure random strings, numbers, and choices with various character sets
// - CSV processing: Table manipulation and data transformation
// - Mathematical utilities: Type-safe numeric operations with generic constraints
// - ID generation: Structured entity ID creation with type safety
// - Timing utilities: Precise timing measurements and deadline management
//
// All data structures provide both regular and thread-safe variants (prefixed with "Safe").
// The thread-safe variants use RWMutex for concurrent access while maintaining performance.
//
// Example usage:
//
//	// Create a thread-safe map
//	m := abstract.NewSafeMap[string, int]()
//	m.Set("key", 42)
//
//	// Generate secure random strings
//	token := abstract.GetRandomString(32)
//
//	// Use futures for concurrent operations
//	future := abstract.NewFuture(ctx, logger, func(ctx context.Context) (string, error) {
//		return "result", nil
//	})
//	result, err := future.Get(ctx)
package abstract

import (
	"math"
	"strconv"
	"sync"
)

// Signed is a constraint that permits any signed integer type.
// This constraint is useful for generic functions that need to work with
// signed integers while maintaining type safety.
// If future releases of Go add new predeclared signed integer types,
// this constraint will be modified to include them.
//
// Example usage:
//
//	func AbsInt[T Signed](x T) T {
//		if x < 0 {
//			return -x
//		}
//		return x
//	}
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
// This constraint is useful for generic functions that need to work with
// unsigned integers while maintaining type safety.
// If future releases of Go add new predeclared unsigned integer types,
// this constraint will be modified to include them.
//
// Example usage:
//
//	func NextPowerOfTwo[T Unsigned](x T) T {
//		if x == 0 {
//			return 1
//		}
//		return T(1) << (64 - bits.LeadingZeros64(uint64(x-1)))
//	}
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type, both signed and unsigned.
// This constraint is useful for generic functions that need to work with
// any integer type while maintaining type safety.
// If future releases of Go add new predeclared integer types,
// this constraint will be modified to include them.
//
// Example usage:
//
//	func IsEven[T Integer](x T) bool {
//		return x%2 == 0
//	}
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
// This constraint is useful for generic functions that need to work with
// floating-point numbers while maintaining type safety.
// If future releases of Go add new predeclared floating-point types,
// this constraint will be modified to include them.
//
// Example usage:
//
//	func IsNaN[T Float](x T) bool {
//		return math.IsNaN(float64(x))
//	}
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
// This constraint is useful for generic functions that need to work with
// complex numbers while maintaining type safety.
// If future releases of Go add new predeclared complex numeric types,
// this constraint will be modified to include them.
//
// Example usage:
//
//	func ComplexAbs[T Complex](x T) float64 {
//		return cmplx.Abs(complex128(x))
//	}
type Complex interface {
	~complex64 | ~complex128
}

// Ordered is a constraint that permits any ordered type: any type
// that supports the comparison operators < <= >= >.
// This constraint is essential for sorting operations and comparisons.
// If future releases of Go add new ordered types,
// this constraint will be modified to include them.
//
// Example usage:
//
//	func Clamp[T Ordered](x, min, max T) T {
//		if x < min {
//			return min
//		}
//		if x > max {
//			return max
//		}
//		return x
//	}
type Ordered interface {
	Integer | Float | ~string
}

// Number is a constraint that permits any numeric type (integers and floats).
// This constraint is useful for generic mathematical operations that work
// with both integer and floating-point types.
//
// Example usage:
//
//	func Square[T Number](x T) T {
//		return x * x
//	}
type Number interface {
	Integer | Float
}

// Orderer is a struct that holds an order of comparable items and provides
// methods to manage ordering operations in a thread-safe manner.
// It's useful for scenarios where you need to track the order of items
// and apply ordering operations atomically.
//
// Example usage:
//
//	orderer := NewOrderer[string](func(order map[string]int) {
//		// Apply the order to your data structure
//		fmt.Println("Applying order:", order)
//	})
//	orderer.Add("item1")
//	orderer.Add("item2")
//	orderer.Apply() // Calls the callback with the current order
type Orderer[T comparable] struct {
	order         map[T]int
	applyCallback func(order map[T]int)
	mu            sync.Mutex
}

// NewOrderer creates a new Orderer with the specified callback function.
// The callback function is called when Apply() is invoked, receiving
// the current order mapping as a parameter.
//
// Parameters:
//   - f: A callback function that receives the order mapping when Apply() is called.
//
// Returns:
//   - A new Orderer instance ready for use.
func NewOrderer[T comparable](f func(order map[T]int)) *Orderer[T] {
	return &Orderer[T]{
		order:         make(map[T]int),
		applyCallback: f,
	}
}

// Add adds an item to the orderer with the next available order index.
// The order index is determined by the current number of items in the orderer.
// This method is thread-safe.
//
// Parameters:
//   - id: The item to add to the orderer.
func (m *Orderer[T]) Add(id T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.order[id] = len(m.order)
}

// Get returns a copy of the current order mapping.
// This method is thread-safe and returns a snapshot of the current state.
//
// Returns:
//   - A map containing the current order mapping.
func (m *Orderer[T]) Get() map[T]int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.order
}

// Apply applies the current order using the callback function and then
// clears the order mapping. This method is thread-safe.
// If the order mapping is empty, the callback is not called.
func (m *Orderer[T]) Apply() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.order) > 0 {
		m.applyCallback(m.order)
	}

	m.order = make(map[T]int)
}

// Clear removes all items from the orderer without calling the callback.
// This method is thread-safe.
func (m *Orderer[T]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.order = make(map[T]int)
}

// Memorizer is a thread-safe container that holds a single item of any type.
// It's useful for scenarios where you need to store and retrieve a single
// value safely across multiple goroutines, with the ability to check if
// the value has been set.
//
// Example usage:
//
//	memo := NewMemorizer[string]()
//	memo.Set("important value")
//	if value, ok := memo.Get(); ok {
//		fmt.Println("Value:", value)
//	}
type Memorizer[T any] struct {
	item  T
	isSet bool
	mu    sync.Mutex
}

// NewMemorizer creates a new Memorizer instance for the specified type.
// The memorizer starts empty (isSet = false).
//
// Returns:
//   - A new Memorizer instance ready for use.
func NewMemorizer[T any]() *Memorizer[T] {
	return &Memorizer[T]{}
}

// Set stores a value in the Memorizer and marks it as set.
// This method is thread-safe.
//
// Parameters:
//   - c: The value to store in the memorizer.
func (m *Memorizer[T]) Set(c T) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.isSet = true
	m.item = c
}

// Get retrieves the value from the Memorizer along with a boolean indicating
// whether the value has been set. This method is thread-safe.
//
// Returns:
//   - The stored value (or zero value if not set).
//   - A boolean indicating whether the value has been set.
func (m *Memorizer[T]) Get() (T, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.item, m.isSet
}

// Pop retrieves the value from the Memorizer and marks it as unset.
// This method is thread-safe and atomically retrieves and clears the value.
//
// Returns:
//   - The stored value (or zero value if not set).
//   - A boolean indicating whether the value was set before popping.
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

// Itoa converts a numeric value to its string representation.
// This is a generic version of strconv.Itoa that works with any numeric type.
//
// Parameters:
//   - i: The numeric value to convert.
//
// Returns:
//   - The string representation of the numeric value.
//
// Example usage:
//
//	str := Itoa(42)        // "42"
//	str := Itoa(3.14)      // "3"
//	str := Itoa(int64(99)) // "99"
func Itoa[T Number](i T) string {
	return strconv.Itoa(int(i))
}

// Atoi converts a string to a numeric value of the specified type.
// This is a generic version of strconv.Atoi that works with any numeric type.
//
// Parameters:
//   - s: The string to convert.
//
// Returns:
//   - The numeric value as the specified type.
//   - An error if the string cannot be converted.
//
// Example usage:
//
//	val, err := Atoi[int]("42")      // 42, nil
//	val, err := Atoi[float64]("99")  // 99.0, nil
//	val, err := Atoi[int8]("300")    // 44, nil (overflow)
func Atoi[T Number](s string) (T, error) {
	i, err := strconv.Atoi(s)
	return T(i), err
}

// Round returns the nearest integer to the input value, rounding half away from zero.
// This function works with any numeric type and returns the same type.
//
// Parameters:
//   - f: The numeric value to round.
//
// Returns:
//   - The rounded value as the same type as the input.
//
// Example usage:
//
//	rounded := Round(3.7)    // 4.0
//	rounded := Round(-2.3)   // -2.0
//	rounded := Round(42)     // 42
func Round[T Number](f T) T {
	return T(math.Round(float64(f)))
}

// Min returns the minimum value from the provided values.
// If no values are provided, returns the zero value for the type.
//
// Parameters:
//   - xs: Variable number of values to compare.
//
// Returns:
//   - The minimum value among the provided values.
//
// Example usage:
//
//	min := Min(1, 2, 3)           // 1
//	min := Min(3.14, 2.71, 1.41)  // 1.41
//	min := Min[int]()             // 0 (zero value)
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
// If no values are provided, returns the zero value for the type.
//
// Parameters:
//   - xs: Variable number of values to compare.
//
// Returns:
//   - The maximum value among the provided values.
//
// Example usage:
//
//	max := Max(1, 2, 3)           // 3
//	max := Max(3.14, 2.71, 1.41)  // 3.14
//	max := Max[int]()             // 0 (zero value)
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

// Abs returns the absolute value of the provided numeric value.
// This function works with any numeric type and returns the same type.
//
// Parameters:
//   - x: The numeric value to get the absolute value of.
//
// Returns:
//   - The absolute value of the input.
//
// Example usage:
//
//	abs := Abs(-5)    // 5
//	abs := Abs(3.14)  // 3.14
//	abs := Abs(-2.71) // 2.71
func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Pow returns the value x raised to the power of y.
// This function works with any numeric types and returns the same type as the base.
//
// Parameters:
//   - x: The base value.
//   - y: The exponent value.
//
// Returns:
//   - The result of x raised to the power of y.
//
// Example usage:
//
//	result := Pow(2, 3)      // 8
//	result := Pow(2.0, 0.5)  // 1.414... (square root of 2)
//	result := Pow(10, -1)    // 0.1
func Pow[T1, T2 Number](x T1, y T2) T1 {
	return T1(math.Pow(float64(x), float64(y)))
}
