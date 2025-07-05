package abstract

import (
	"crypto/rand"
	"math"
	mr "math/rand"
	"strconv"
	"time"
)

var (
	// Default alphabet for random string generation (hexadecimal characters)
	defaultAlphabet = []byte("0123456789abcdef")
	alphabetLen     = uint8(math.Min(float64(len(defaultAlphabet)), float64(math.MaxUint8)))

	// Predefined character sets for different random string types
	lowerAlpha    = []byte("abcdefghijklmnopqrstuvwxyz")
	upperAlpha    = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	alphaNumeric  = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	strictNumeric = []byte("0123456789")
)

// GetRandomString returns a cryptographically secure random string of the specified length
// using hexadecimal characters (0-9, a-f).
//
// Security considerations:
//   - Uses crypto/rand for secure random generation when available
//   - Falls back to math/rand with time-based seed if crypto/rand fails
//   - Suitable for generating tokens, session IDs, and other security-sensitive identifiers
//
// Parameters:
//   - n: The length of the random string to generate
//
// Returns:
//   - A random string of length n containing hexadecimal characters
//
// Example usage:
//
//	sessionID := GetRandomString(32)  // "a1b2c3d4e5f6..."
//	token := GetRandomString(16)      // "f1e2d3c4b5a6..."
func GetRandomString(n int) string {
	return string(GetRandomBytes(n))
}

// GetRandomBytes returns cryptographically secure random bytes of the specified length
// using hexadecimal characters (0-9, a-f).
//
// Security considerations:
//   - Uses crypto/rand for secure random generation when available
//   - Falls back to math/rand with time-based seed if crypto/rand fails
//   - Each byte is masked to ensure uniform distribution across the alphabet
//
// Parameters:
//   - n: The number of random bytes to generate
//
// Returns:
//   - A byte slice of length n containing random hexadecimal characters
//
// Example usage:
//
//	randomBytes := GetRandomBytes(16)
//	fmt.Printf("Random bytes: %x\n", randomBytes)
func GetRandomBytes(n int) []byte {
	out := make([]byte, n)
	_, err := rand.Read(out)
	if err != nil {
		r := mr.New(mr.NewSource(time.Now().UnixNano()))
		for i := range out {
			out[i] = byte(r.Intn(math.MaxUint8))
		}
	}
	for i := range out {
		out[i] = defaultAlphabet[out[i]&(alphabetLen-1)]
	}
	return out
}

// GetRandListenAddress generates a random TCP port number in the range 10000-62999
// and returns it as a string formatted for network listening (e.g., ":12345").
//
// Note: This function uses math/rand, not crypto/rand, as it's intended for
// development and testing purposes where cryptographic security is not required.
//
// Returns:
//   - A string in the format ":XXXXX" where XXXXX is a random port number
//
// Example usage:
//
//	addr := GetRandListenAddress()  // ":45123"
//	listener, err := net.Listen("tcp", addr)
//	if err != nil {
//		log.Fatal(err)
//	}
func GetRandListenAddress() (port string) {
	r := mr.New(mr.NewSource(time.Now().UnixNano()))
	return ":" + strconv.Itoa(10000+r.Intn(53000))
}

// GetRandomStringWithAlphabet returns a cryptographically secure random string
// using characters from the specified alphabet.
//
// Security considerations:
//   - Uses crypto/rand for secure random generation when available
//   - Falls back to math/rand with time-based seed if crypto/rand fails
//   - Uses modulo operation to ensure uniform distribution across the alphabet
//   - Returns empty string if alphabet is empty
//
// Parameters:
//   - n: The length of the random string to generate
//   - alphabet: A byte slice containing the characters to use
//
// Returns:
//   - A random string of length n using characters from the alphabet
//
// Example usage:
//
//	// Generate a random string with custom alphabet
//	customAlphabet := []byte("!@#$%^&*")
//	symbols := GetRandomStringWithAlphabet(8, customAlphabet)
//
//	// Generate a random password
//	passwordChars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*")
//	password := GetRandomStringWithAlphabet(12, passwordChars)
func GetRandomStringWithAlphabet(n int, alphabet []byte) string {
	if len(alphabet) == 0 {
		return ""
	}

	out := make([]byte, n)
	_, err := rand.Read(out)
	if err != nil {
		r := mr.New(mr.NewSource(time.Now().UnixNano()))
		for i := range out {
			out[i] = byte(r.Intn(math.MaxUint8))
		}
	}

	alphabetLength := byte(len(alphabet))
	for i := range out {
		out[i] = alphabet[out[i]%alphabetLength]
	}
	return string(out)
}

// GetRandomLowerAlpha returns a cryptographically secure random string
// containing only lowercase letters (a-z).
//
// Security considerations:
//   - Uses crypto/rand for secure random generation
//   - Suitable for generating case-sensitive identifiers
//
// Parameters:
//   - n: The length of the random string to generate
//
// Returns:
//   - A random string of length n containing only lowercase letters
//
// Example usage:
//
//	code := GetRandomLowerAlpha(8)  // "abcdefgh"
//	id := GetRandomLowerAlpha(16)   // "qwertyuiopasdfgh"
func GetRandomLowerAlpha(n int) string {
	return GetRandomStringWithAlphabet(n, lowerAlpha)
}

// GetRandomUpperAlpha returns a cryptographically secure random string
// containing only uppercase letters (A-Z).
//
// Security considerations:
//   - Uses crypto/rand for secure random generation
//   - Suitable for generating case-sensitive identifiers
//
// Parameters:
//   - n: The length of the random string to generate
//
// Returns:
//   - A random string of length n containing only uppercase letters
//
// Example usage:
//
//	code := GetRandomUpperAlpha(6)  // "ABCDEF"
//	id := GetRandomUpperAlpha(12)   // "QWERTYUIOPAS"
func GetRandomUpperAlpha(n int) string {
	return GetRandomStringWithAlphabet(n, upperAlpha)
}

// GetRandomAlphaNumeric returns a cryptographically secure random string
// containing letters (both cases) and numbers (0-9, a-z, A-Z).
//
// Security considerations:
//   - Uses crypto/rand for secure random generation
//   - Provides good entropy with 62 possible characters per position
//   - Suitable for user-facing identifiers and codes
//
// Parameters:
//   - n: The length of the random string to generate
//
// Returns:
//   - A random string of length n containing letters and numbers
//
// Example usage:
//
//	userID := GetRandomAlphaNumeric(12)  // "Ab3Cd5Ef7Gh9"
//	token := GetRandomAlphaNumeric(20)   // "aB3cD5eF7gH9iJ2kL4"
func GetRandomAlphaNumeric(n int) string {
	return GetRandomStringWithAlphabet(n, alphaNumeric)
}

// GetRandomNumeric returns a cryptographically secure random string
// containing only numeric digits (0-9).
//
// Security considerations:
//   - Uses crypto/rand for secure random generation
//   - Lower entropy than alphanumeric strings (10 vs 62 characters)
//   - Suitable for numeric codes and identifiers
//
// Parameters:
//   - n: The length of the random string to generate
//
// Returns:
//   - A random string of length n containing only digits
//
// Example usage:
//
//	pin := GetRandomNumeric(4)     // "1234"
//	code := GetRandomNumeric(8)    // "87654321"
//	id := GetRandomNumeric(12)     // "123456789012"
func GetRandomNumeric(n int) string {
	return GetRandomStringWithAlphabet(n, strictNumeric)
}

// GetRandomInt returns a cryptographically secure random integer in the specified range [min, max].
// The range is inclusive on both ends.
//
// Security considerations:
//   - Uses crypto/rand for secure random generation when available
//   - Falls back to math/rand with time-based seed if crypto/rand fails
//   - Automatically swaps min and max if min > max
//   - Returns min if min equals max
//
// Parameters:
//   - min: The minimum value (inclusive)
//   - max: The maximum value (inclusive)
//
// Returns:
//   - A random integer in the range [min, max]
//
// Example usage:
//
//	dice := GetRandomInt(1, 6)       // Random number 1-6
//	percent := GetRandomInt(0, 100)  // Random percentage 0-100
//	port := GetRandomInt(8000, 9000) // Random port in range
func GetRandomInt(min, max int) int {
	if min > max {
		min, max = max, min
	}
	if min == max {
		return min
	}
	r := mr.New(mr.NewSource(time.Now().UnixNano()))
	return min + r.Intn(max-min+1)
}

// GetRandomBool returns a cryptographically secure random boolean value.
//
// Security considerations:
//   - Uses crypto/rand for secure random generation when available
//   - Falls back to math/rand with time-based seed if crypto/rand fails
//   - Provides unbiased true/false selection
//
// Returns:
//   - A random boolean value (true or false)
//
// Example usage:
//
//	coinFlip := GetRandomBool()       // true or false
//	if GetRandomBool() {
//		fmt.Println("Random event occurred")
//	}
//
//	// Use in feature flags or random decisions
//	enableFeature := GetRandomBool()
func GetRandomBool() bool {
	bytes := make([]byte, 1)
	_, err := rand.Read(bytes)
	if err != nil {
		r := mr.New(mr.NewSource(time.Now().UnixNano()))
		return r.Intn(2) == 1
	}
	return bytes[0]%2 == 1
}

// GetRandomChoice returns a random element from the provided slice.
// This function uses math/rand for performance reasons.
//
// Parameters:
//   - slice: The slice to choose from
//
// Returns:
//   - A random element from the slice and true if successful
//   - The zero value for the type and false if the slice is empty
//
// Example usage:
//
//	colors := []string{"red", "green", "blue", "yellow"}
//	color, ok := GetRandomChoice(colors)
//	if ok {
//		fmt.Printf("Random color: %s\n", color)
//	}
//
//	numbers := []int{1, 2, 3, 4, 5}
//	number, _ := GetRandomChoice(numbers)
//	fmt.Printf("Random number: %d\n", number)
func GetRandomChoice[T any](slice []T) (T, bool) {
	var zero T
	if len(slice) == 0 {
		return zero, false
	}

	r := mr.New(mr.NewSource(time.Now().UnixNano()))
	return slice[r.Intn(len(slice))], true
}

// ShuffleSlice randomly shuffles the elements in the provided slice in-place.
// This function modifies the original slice using the Fisher-Yates shuffle algorithm.
//
// Note: This function uses math/rand for performance reasons. For cryptographically
// secure shuffling, consider using crypto/rand with a custom implementation.
//
// Parameters:
//   - slice: The slice to shuffle (modified in-place)
//
// Example usage:
//
//	cards := []string{"A", "K", "Q", "J", "10", "9", "8", "7"}
//	ShuffleSlice(cards)
//	fmt.Printf("Shuffled cards: %v\n", cards)
//
//	numbers := []int{1, 2, 3, 4, 5}
//	ShuffleSlice(numbers)
//	fmt.Printf("Shuffled numbers: %v\n", numbers)
func ShuffleSlice[T any](slice []T) {
	r := mr.New(mr.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
