package abstract

import (
	"crypto/rand"
	"math"
	mr "math/rand"
	"strconv"
	"time"
)

var (
	defaultAlphabet = []byte("0123456789abcdef")
	alphabetLen     = uint8(math.Min(float64(len(defaultAlphabet)), float64(math.MaxUint8)))
	// Predefined character sets
	lowerAlpha    = []byte("abcdefghijklmnopqrstuvwxyz")
	upperAlpha    = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	alphaNumeric  = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	strictNumeric = []byte("0123456789")
)

// GetRandomString returns a random string of length n contains letters from hex.
func GetRandomString(n int) string {
	return string(GetRandomBytes(n))
}

// GetRandomBytes returns a random bytes of length n contains letters from hex.
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

// GetRandListenAddress generates a random port number between 10000 and 63000.
func GetRandListenAddress() (port string) {
	r := mr.New(mr.NewSource(time.Now().UnixNano()))
	return ":" + strconv.Itoa(10000+r.Intn(53000))
}

// GetRandomStringWithAlphabet returns a random string of length n using characters
// from the provided alphabet.
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

// GetRandomLowerAlpha returns a random lowercase alphabetic string of length n.
func GetRandomLowerAlpha(n int) string {
	return GetRandomStringWithAlphabet(n, lowerAlpha)
}

// GetRandomUpperAlpha returns a random uppercase alphabetic string of length n.
func GetRandomUpperAlpha(n int) string {
	return GetRandomStringWithAlphabet(n, upperAlpha)
}

// GetRandomAlphaNumeric returns a random alphanumeric string of length n.
func GetRandomAlphaNumeric(n int) string {
	return GetRandomStringWithAlphabet(n, alphaNumeric)
}

// GetRandomNumeric returns a random numeric string of length n.
func GetRandomNumeric(n int) string {
	return GetRandomStringWithAlphabet(n, strictNumeric)
}

// GetRandomInt returns a random integer in the range [min, max].
// If min > max, it will swap them automatically.
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

// GetRandomBool returns a random boolean value.
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
// If the slice is empty, returns the zero value for the type.
func GetRandomChoice[T any](slice []T) (T, bool) {
	var zero T
	if len(slice) == 0 {
		return zero, false
	}

	r := mr.New(mr.NewSource(time.Now().UnixNano()))
	return slice[r.Intn(len(slice))], true
}

// ShuffleSlice randomly shuffles the elements in the provided slice.
func ShuffleSlice[T any](slice []T) {
	r := mr.New(mr.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
