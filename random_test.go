package abstract_test

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maxbolgarin/abstract"
)

var defaultAlphabet = []byte("0123456789abcdef")

// TestGetRandomString ensures that GetRandomString returns a string of the requested length.
func TestGetRandomString(t *testing.T) {
	const length = 10
	result := abstract.GetRandomString(length)

	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}

	// Verify only hex chars are used
	if !regexp.MustCompile(`^[0-9a-f]+$`).MatchString(result) {
		t.Errorf("Result contains non-hex characters: %s", result)
	}
}

// TestGetRandomBytes ensures that GetRandomBytes returns a byte slice of the requested length.
func TestGetRandomBytes(t *testing.T) {
	const length = 10
	result := abstract.GetRandomBytes(length)

	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}

	// Verify only hex chars are used
	for _, b := range result {
		if !((b >= '0' && b <= '9') || (b >= 'a' && b <= 'f')) {
			t.Errorf("Result contains non-hex character: %c", b)
		}
	}
}

// TestGetRandListenAddress ensures that GetRandListenAddress returns a valid port number.
func TestGetRandListenAddress(t *testing.T) {
	result := abstract.GetRandListenAddress()

	// Should start with a colon
	if !strings.HasPrefix(result, ":") {
		t.Errorf("Expected to start with ':', got %s", result)
	}

	// Extract port number
	port, err := strconv.Atoi(result[1:])
	if err != nil {
		t.Errorf("Failed to parse port number: %v", err)
	}

	// Validate port range
	if port < 10000 || port > 63000 {
		t.Errorf("Port %d outside of expected range [10000, 63000]", port)
	}
}

// TestGetRandomStringWithAlphabet ensures that GetRandomStringWithAlphabet returns a string of the requested length using the specified alphabet.
func TestGetRandomStringWithAlphabet(t *testing.T) {
	alphabet := []byte("ABC123")
	const length = 15

	result := abstract.GetRandomStringWithAlphabet(length, alphabet)

	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}

	// Verify only chars from our alphabet are used
	for _, char := range result {
		found := false
		for _, validChar := range alphabet {
			if char == rune(validChar) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Result contains character not in alphabet: %c", char)
		}
	}

	// Test with empty alphabet
	emptyResult := abstract.GetRandomStringWithAlphabet(10, []byte{})
	if emptyResult != "" {
		t.Errorf("Expected empty string for empty alphabet, got %s", emptyResult)
	}
}

// TestGetRandomLowerAlpha ensures that GetRandomLowerAlpha returns a string of the requested length containing only lowercase alphabetic characters.
func TestGetRandomLowerAlpha(t *testing.T) {
	const length = 10
	result := abstract.GetRandomLowerAlpha(length)

	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}

	// Verify only lowercase alpha chars are used
	if !regexp.MustCompile(`^[a-z]+$`).MatchString(result) {
		t.Errorf("Result contains non-lowercase-alpha characters: %s", result)
	}
}

// TestGetRandomUpperAlpha ensures that GetRandomUpperAlpha returns a string of the requested length containing only uppercase alphabetic characters.
func TestGetRandomUpperAlpha(t *testing.T) {
	const length = 10
	result := abstract.GetRandomUpperAlpha(length)

	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}

	// Verify only uppercase alpha chars are used
	if !regexp.MustCompile(`^[A-Z]+$`).MatchString(result) {
		t.Errorf("Result contains non-uppercase-alpha characters: %s", result)
	}
}

// TestGetRandomAlphaNumeric ensures that GetRandomAlphaNumeric returns a string of the requested length containing only alphanumeric characters.
func TestGetRandomAlphaNumeric(t *testing.T) {
	const length = 10
	result := abstract.GetRandomAlphaNumeric(length)

	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}

	// Verify only alphanumeric chars are used
	if !regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString(result) {
		t.Errorf("Result contains non-alphanumeric characters: %s", result)
	}
}

// TestGetRandomNumeric ensures that GetRandomNumeric returns a string of the requested length containing only numeric characters.
func TestGetRandomNumeric(t *testing.T) {
	const length = 10
	result := abstract.GetRandomNumeric(length)

	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}

	// Verify only numeric chars are used
	if !regexp.MustCompile(`^[0-9]+$`).MatchString(result) {
		t.Errorf("Result contains non-numeric characters: %s", result)
	}
}

// TestGetRandomInt ensures that GetRandomInt returns a random integer within the specified range.
func TestGetRandomInt(t *testing.T) {
	// Test normal range
	min, max := 10, 20
	for i := 0; i < 100; i++ {
		result := abstract.GetRandomInt(min, max)
		if result < min || result > max {
			t.Errorf("Random int %d outside of range [%d, %d]", result, min, max)
		}
	}

	// Test with min > max (should swap)
	for i := 0; i < 100; i++ {
		result := abstract.GetRandomInt(20, 10)
		if result < 10 || result > 20 {
			t.Errorf("Random int %d outside of range [%d, %d] after swapping", result, 10, 20)
		}
	}

	// Test with min = max
	result := abstract.GetRandomInt(15, 15)
	if result != 15 {
		t.Errorf("Expected %d for equal min/max, got %d", 15, result)
	}
}

// TestGetRandomBool ensures that GetRandomBool returns a random boolean value.
func TestGetRandomBool(t *testing.T) {
	// Run multiple times to ensure both values occur
	trueCount, falseCount := 0, 0
	iterations := 1000

	for i := 0; i < iterations; i++ {
		if abstract.GetRandomBool() {
			trueCount++
		} else {
			falseCount++
		}
	}

	// With enough iterations, we should get a reasonable distribution
	if trueCount == 0 || falseCount == 0 {
		t.Errorf("Expected both true and false values, got %d true and %d false", trueCount, falseCount)
	}
}

// TestGetRandomChoice ensures that GetRandomChoice returns a random element from the slice and indicates whether the element was found.
func TestGetRandomChoice(t *testing.T) {
	// Test with non-empty slice
	slice := []string{"apple", "banana", "cherry", "date"}

	// Run multiple times to ensure different values are selected
	seen := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		result, ok := abstract.GetRandomChoice(slice)
		if !ok {
			t.Errorf("Expected ok=true for non-empty slice")
		}

		// Check that result is from the slice
		found := false
		for _, item := range slice {
			if result == item {
				found = true
				seen[result] = true
				break
			}
		}

		if !found {
			t.Errorf("Result %v not found in slice", result)
		}
	}

	// With enough iterations, we should see all items
	if len(seen) < len(slice) {
		t.Errorf("Expected to see all items from slice, only saw %d out of %d", len(seen), len(slice))
	}

	// Test with empty slice
	emptySlice := []string{}
	_, ok := abstract.GetRandomChoice(emptySlice)
	if ok {
		t.Errorf("Expected ok=false for empty slice")
	}
}

// TestShuffleSlice ensures that ShuffleSlice shuffles the slice and maintains the original elements.
func TestShuffleSlice(t *testing.T) {
	// Create a slice of unique items
	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	shuffled := make([]int, len(original))
	copy(shuffled, original)

	// Shuffle multiple times to ensure different orders
	matchCount := 0
	iterations := 10

	for i := 0; i < iterations; i++ {
		abstract.ShuffleSlice(shuffled)

		// Count matches with original order
		matches := 0
		for j := range original {
			if original[j] == shuffled[j] {
				matches++
			}
		}

		// If matches == len(original), the shuffle didn't change anything
		if matches == len(original) {
			matchCount++
		}
	}

	// It's extremely unlikely all iterations match the original order
	if matchCount == iterations {
		t.Errorf("After %d shuffles, slice remained in original order", iterations)
	}

	// Verify all original elements are still in the shuffled slice
	if len(shuffled) != len(original) {
		t.Errorf("Shuffled slice length changed: expected %d, got %d", len(original), len(shuffled))
	}

	counts := make(map[int]int)
	for _, v := range shuffled {
		counts[v]++
	}

	for _, v := range original {
		if counts[v] != 1 {
			t.Errorf("Element %d appears %d times in shuffled slice, expected once", v, counts[v])
		}
	}
}

// TestGetRandomStringFast ensures that GetRandomStringFast returns a string of the requested length.
func TestGetRandomStringFast(t *testing.T) {
	const length = 10
	result := abstract.GetRandomStringFast(length)

	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}

	// Verify only hex chars are used (same as regular GetRandomString)
	if !regexp.MustCompile(`^[0-9a-f]+$`).MatchString(result) {
		t.Errorf("Result contains non-hex characters: %s", result)
	}

	// Test with zero length
	zeroResult := abstract.GetRandomStringFast(0)
	if zeroResult != "" {
		t.Errorf("Expected empty string for length 0, got %s", zeroResult)
	}

	// Test with negative length
	negativeResult := abstract.GetRandomStringFast(-1)
	if negativeResult != "" {
		t.Errorf("Expected empty string for negative length, got %s", negativeResult)
	}
}

// TestGetRandomBytesFast ensures that GetRandomBytesFast returns a byte slice of the requested length.
func TestGetRandomBytesFast(t *testing.T) {
	const length = 10
	result := abstract.GetRandomBytesFast(length)

	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}

	// Verify only hex chars are used
	for _, b := range result {
		if !((b >= '0' && b <= '9') || (b >= 'a' && b <= 'f')) {
			t.Errorf("Result contains non-hex character: %c", b)
		}
	}

	// Test with zero length
	zeroResult := abstract.GetRandomBytesFast(0)
	if zeroResult != nil {
		t.Errorf("Expected nil for length 0, got %v", zeroResult)
	}

	// Test with negative length
	negativeResult := abstract.GetRandomBytesFast(-1)
	if negativeResult != nil {
		t.Errorf("Expected nil for negative length, got %v", negativeResult)
	}
}

// TestGetRandomStringFastConsistency ensures that GetRandomStringFast produces consistent results
// for the same input length (though not necessarily the same string due to time-based seeding).
func TestGetRandomStringFastConsistency(t *testing.T) {
	const length = 8
	results := make(map[string]bool)
	iterations := 100

	// Generate multiple strings and check they're all valid
	for i := 0; i < iterations; i++ {
		result := abstract.GetRandomStringFast(length)

		if len(result) != length {
			t.Errorf("Expected length %d, got %d", length, len(result))
		}

		// Verify hex characters
		if !regexp.MustCompile(`^[0-9a-f]+$`).MatchString(result) {
			t.Errorf("Result contains non-hex characters: %s", result)
		}

		results[result] = true
	}

	// With enough iterations, we should see some variety (though not guaranteed due to LCG)
	if len(results) < iterations/10 {
		t.Logf("Warning: Only %d unique strings generated out of %d iterations", len(results), iterations)
	}
}

// TestGetRandomBytesFastConsistency ensures that GetRandomBytesFast produces consistent results
// for the same input length.
func TestGetRandomBytesFastConsistency(t *testing.T) {
	const length = 8
	results := make(map[string]bool)
	iterations := 100

	// Generate multiple byte slices and check they're all valid
	for i := 0; i < iterations; i++ {
		result := abstract.GetRandomBytesFast(length)

		if len(result) != length {
			t.Errorf("Expected length %d, got %d", length, len(result))
		}

		// Verify hex characters
		for _, b := range result {
			if !((b >= '0' && b <= '9') || (b >= 'a' && b <= 'f')) {
				t.Errorf("Result contains non-hex character: %c", b)
			}
		}

		results[string(result)] = true
	}

	// With enough iterations, we should see some variety (though not guaranteed due to LCG)
	if len(results) < iterations/10 {
		t.Logf("Warning: Only %d unique byte slices generated out of %d iterations", len(results), iterations)
	}
}

// TestGetRandomStringFastPerformance compares the performance characteristics
// of the fast vs regular random string generation.
func TestGetRandomStringFastPerformance(t *testing.T) {
	const length = 100
	const iterations = 1000

	// Test fast version
	fastStart := time.Now()
	for i := 0; i < iterations; i++ {
		abstract.GetRandomStringFast(length)
	}
	fastDuration := time.Since(fastStart)

	// Test regular version
	regularStart := time.Now()
	for i := 0; i < iterations; i++ {
		abstract.GetRandomString(length)
	}
	regularDuration := time.Since(regularStart)

	// Log performance comparison
	t.Logf("Fast version took %v for %d iterations", fastDuration, iterations)
	t.Logf("Regular version took %v for %d iterations", regularDuration, iterations)
	t.Logf("Fast version is %.2fx faster", float64(regularDuration)/float64(fastDuration))
}

// TestGetRandomBytesFastEdgeCases tests edge cases for GetRandomBytesFast.
func TestGetRandomBytesFastEdgeCases(t *testing.T) {
	// Test very large length
	largeResult := abstract.GetRandomBytesFast(1000)
	if len(largeResult) != 1000 {
		t.Errorf("Expected length 1000, got %d", len(largeResult))
	}

	// Verify all bytes are valid hex characters
	for i, b := range largeResult {
		if !((b >= '0' && b <= '9') || (b >= 'a' && b <= 'f')) {
			t.Errorf("Byte at position %d is not a hex character: %c", i, b)
		}
	}

	// Test length 1
	singleResult := abstract.GetRandomBytesFast(1)
	if len(singleResult) != 1 {
		t.Errorf("Expected length 1, got %d", len(singleResult))
	}

	// Test that the single byte is a valid hex character
	b := singleResult[0]
	if !((b >= '0' && b <= '9') || (b >= 'a' && b <= 'f')) {
		t.Errorf("Single byte is not a hex character: %c", b)
	}
}

// TestGetRandomStringFastEdgeCases tests edge cases for GetRandomStringFast.
func TestGetRandomStringFastEdgeCases(t *testing.T) {
	// Test very large length
	largeResult := abstract.GetRandomStringFast(1000)
	if len(largeResult) != 1000 {
		t.Errorf("Expected length 1000, got %d", len(largeResult))
	}

	// Verify all characters are valid hex characters
	if !regexp.MustCompile(`^[0-9a-f]+$`).MatchString(largeResult) {
		t.Errorf("Large result contains non-hex characters: %s", largeResult)
	}

	// Test length 1
	singleResult := abstract.GetRandomStringFast(1)
	if len(singleResult) != 1 {
		t.Errorf("Expected length 1, got %d", len(singleResult))
	}

	// Test that the single character is a valid hex character
	if !regexp.MustCompile(`^[0-9a-f]$`).MatchString(singleResult) {
		t.Errorf("Single character is not a hex character: %s", singleResult)
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
