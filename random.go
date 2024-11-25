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
