package abstract_test

import (
	"strconv"
	"testing"

	"github.com/maxbolgarin/abstract"
)

var defaultAlphabet = []byte("0123456789abcdef")

// TestGetRandomString ensures that GetRandomString returns a string of the requested length.
func TestGetRandomString(t *testing.T) {
	lengths := []int{0, 1, 8, 16, 32}

	for _, length := range lengths {
		result := abstract.GetRandomString(length)
		if len(result) != length {
			t.Errorf("expected length %d, got %d", length, len(result))
		}
		// Check if each character in the result is a valid hex character
		for i := range result {
			if !isHexChar(result[i]) {
				t.Errorf("expected hex character, got %c", result[i])
			}
		}
	}
}

// TestGetRandomBytes ensures that GetRandomBytes returns a byte slice of the requested length.
func TestGetRandomBytes(t *testing.T) {
	lengths := []int{0, 1, 8, 16, 32}

	for _, length := range lengths {
		result := abstract.GetRandomBytes(length)
		if len(result) != length {
			t.Errorf("expected length %d, got %d", length, len(result))
		}
		// Check if each byte in the result is a valid index of defaultAlphabet
		for i := range result {
			if !isHexChar(result[i]) {
				t.Errorf("expected valid index in defaultAlphabet, got %c", result[i])
			}
		}
	}
}

// TestGetRandListenAddress ensures that GetRandListenAddress returns a valid port number.
func TestGetRandListenAddress(t *testing.T) {
	for i := 0; i < 100; i++ { // run multiple tests to ensure randomness
		result := abstract.GetRandListenAddress()
		// Check if result starts with ':'
		if result[0] != ':' {
			t.Errorf("expected ':', got %s", result)
		}

		portString := result[1:]
		port, err := strconv.Atoi(portString)
		if err != nil {
			t.Errorf("expected a valid integer port, got %s", portString)
		}

		if port < 10000 || port > 63000 {
			t.Errorf("expected port between 10000 and 63000, got %d", port)
		}
	}
}

// Helper function to check if a character is a valid hex character
func isHexChar(c byte) bool {
	for _, hc := range defaultAlphabet {
		if c == hc {
			return true
		}
	}
	return false
}
