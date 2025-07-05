package abstract_test

import (
	"math"
	"sync"
	"testing"

	"github.com/maxbolgarin/abstract"
)

// Helper function for testing Apply method
func callback(order map[string]int) {
	// Dummy callback; implement if needed for complex tests
}

func TestNewOrderer(t *testing.T) {
	orderer := abstract.NewOrderer(callback)
	if orderer == nil {
		t.Fatalf("Expected non-nil Orderer object")
	}
}

func TestOrderer_AddAndGet(t *testing.T) {
	orderer := abstract.NewOrderer(callback)

	testData := []string{"a", "b", "c"}
	for _, item := range testData {
		orderer.Add(item)
	}

	order := orderer.Get()
	expectedLength := len(testData)

	if len(order) != expectedLength {
		t.Errorf("Expected order length %d, got %d", expectedLength, len(order))
	}

	for i, item := range testData {
		if order[item] != i {
			t.Errorf("Expected order[%s] to be %d, got %d", item, i, order[item])
		}
	}
}

func TestOrderer_Apply(t *testing.T) {
	appliedOrder := make(map[string]int)
	var mu sync.Mutex

	orderer := abstract.NewOrderer(func(order map[string]int) {
		mu.Lock()
		defer mu.Unlock()
		for k, v := range order {
			appliedOrder[k] = v
		}
	})

	// Add some data
	testData := []string{"a", "b", "c"}
	for _, item := range testData {
		orderer.Add(item)
	}

	// Apply the order
	orderer.Apply()

	// Check if the callback has the correct data
	mu.Lock()
	defer mu.Unlock()
	expectedLength := len(testData)
	if len(appliedOrder) != expectedLength {
		t.Errorf("Expected applied order length %d, got %d", expectedLength, len(appliedOrder))
	}

	for i, item := range testData {
		if appliedOrder[item] != i {
			t.Errorf("Expected applied order[%s] to be %d, got %d", item, i, appliedOrder[item])
		}
	}

	// Check if the order is cleared after Apply
	order := orderer.Get()
	if len(order) != 0 {
		t.Errorf("Expected order to be cleared, but got %d items", len(order))
	}
}

func TestOrderer_Clear(t *testing.T) {
	orderer := abstract.NewOrderer(callback)

	// Add some data
	testData := []string{"a", "b", "c"}
	for _, item := range testData {
		orderer.Add(item)
	}

	// Clear the order
	orderer.Clear()

	order := orderer.Get()
	if len(order) != 0 {
		t.Errorf("Expected order to be cleared, but got %d items", len(order))
	}
}

func TestNewMemorizer(t *testing.T) {
	memorizer := abstract.NewMemorizer[int]()
	if memorizer == nil {
		t.Fatalf("Expected non-nil Memorizer object")
	}

	// Test that a newly created Memorizer has no set value
	value, isSet := memorizer.Get()
	if isSet {
		t.Errorf("Expected isSet to be false, got true")
	}

	var zeroValue int
	if value != zeroValue {
		t.Errorf("Expected zero value, got %v", value)
	}
}

func TestMemorizer_SetAndGet(t *testing.T) {
	memorizer := abstract.NewMemorizer[string]()

	valueToSet := "test value"
	memorizer.Set(valueToSet)

	value, isSet := memorizer.Get()
	if !isSet {
		t.Errorf("Expected isSet to be true, got false")
	}

	if value != valueToSet {
		t.Errorf("Expected value %v, got %v", valueToSet, value)
	}
}

func TestMemorizer_Pop(t *testing.T) {
	memorizer := abstract.NewMemorizer[float64]()

	valueToSet := 3.14
	memorizer.Set(valueToSet)

	value, isSet := memorizer.Pop()
	if !isSet {
		t.Errorf("Expected isSet after Pop to be true, got false")
	}

	if value != valueToSet {
		t.Errorf("Expected value %v after Pop, got %v", valueToSet, value)
	}

	// Ensure the item is removed after Pop
	_, isSet = memorizer.Get()
	if isSet {
		t.Errorf("Expected isSet to be false after Pop, got true")
	}

	// Pop again should return zero value and false
	value, isSet = memorizer.Pop()
	var zeroValue float64
	if value != zeroValue || isSet {
		t.Errorf("Expected zero value and false, got %v and %v", value, isSet)
	}
}

func TestMemorizer_SetAndPop(t *testing.T) {
	memorizer := abstract.NewMemorizer[int]()

	valueToSet := 42
	memorizer.Set(valueToSet)

	// Get the value before pop to ensure it's set
	value, isSet := memorizer.Get()
	if !isSet || value != valueToSet {
		t.Errorf("Expected value %v and isSet true, got %v and isSet %v", valueToSet, value, isSet)
	}

	// Pop should retrieve the same set value
	value, isSet = memorizer.Pop()
	if !isSet || value != valueToSet {
		t.Errorf("Expected value %v and isSet true after Pop, got %v and isSet %v", valueToSet, value, isSet)
	}

	// Check if popping again returns default zero value and false
	value, isSet = memorizer.Pop()
	var zeroValue int
	if value != zeroValue || isSet {
		t.Errorf("Expected zero value %v and isSet false after additional Pop, got %v and isSet %v", zeroValue, value, isSet)
	}
}

func TestItoa(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"Integer", 42, "42"},
		{"Negative Integer", -123, "-123"},
		{"Zero", 0, "0"},
		{"Float converted to int", 3.14, "3"},
		{"Large number", 9999999, "9999999"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result string
			switch v := test.input.(type) {
			case int:
				result = abstract.Itoa(v)
			case float64:
				result = abstract.Itoa(v)
			}
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestAtoi(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{"Valid Integer", "42", 42, false},
		{"Negative Integer", "-123", -123, false},
		{"Zero", "0", 0, false},
		{"Invalid String", "abc", 0, true},
		{"Empty String", "", 0, true},
		{"Mixed String", "123abc", 0, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := abstract.Atoi[int](test.input)
			if test.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !test.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if !test.expectError && result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"Integer", 42.0, 42.0},
		{"Float Up", 3.7, 4.0},
		{"Float Down", 3.2, 3.0},
		{"Float Midpoint", 3.5, 4.0},
		{"Negative Float Up", -3.7, -4.0},
		{"Negative Float Down", -3.2, -3.0},
		{"Negative Float Midpoint", -3.5, -4.0},
		{"Zero", 0.0, 0.0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := abstract.Round(test.input)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestMin(t *testing.T) {
	t.Run("Integer values", func(t *testing.T) {
		result := abstract.Min(5, 3, 8, 1, 9)
		expected := 1
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Float values", func(t *testing.T) {
		result := abstract.Min(3.5, 2.1, 4.9)
		expected := 2.1
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Negative values", func(t *testing.T) {
		result := abstract.Min(-5, -3, -8, -1, -9)
		expected := -9
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Mixed values", func(t *testing.T) {
		result := abstract.Min(-5, 3, -8, 1, 9)
		expected := -8
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Single value", func(t *testing.T) {
		result := abstract.Min(5)
		expected := 5
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Empty slice", func(t *testing.T) {
		result := abstract.Min[int]()
		expected := 0
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestMax(t *testing.T) {
	t.Run("Integer values", func(t *testing.T) {
		result := abstract.Max(5, 3, 8, 1, 9)
		expected := 9
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Float values", func(t *testing.T) {
		result := abstract.Max(3.5, 2.1, 4.9)
		expected := 4.9
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Negative values", func(t *testing.T) {
		result := abstract.Max(-5, -3, -8, -1, -9)
		expected := -1
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Mixed values", func(t *testing.T) {
		result := abstract.Max(-5, 3, -8, 1, 9)
		expected := 9
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Single value", func(t *testing.T) {
		result := abstract.Max(5)
		expected := 5
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Empty slice", func(t *testing.T) {
		result := abstract.Max[int]()
		expected := 0
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{"Positive Integer", 42, 42},
		{"Negative Integer", -123, 123},
		{"Zero", 0, 0},
		{"Positive Float", 3.14, 3.14},
		{"Negative Float", -3.14, 3.14},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch input := test.input.(type) {
			case int:
				result := abstract.Abs(input)
				if result != test.expected.(int) {
					t.Errorf("Expected %v, got %v", test.expected, result)
				}
			case float64:
				result := abstract.Abs(input)
				if result != test.expected.(float64) {
					t.Errorf("Expected %v, got %v", test.expected, result)
				}
			}
		})
	}
}

func TestPow(t *testing.T) {
	t.Run("Integer base and exponent", func(t *testing.T) {
		result := abstract.Pow(5, 2)
		expected := 25
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Cube", func(t *testing.T) {
		result := abstract.Pow(3, 3)
		expected := 27
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Zero Power", func(t *testing.T) {
		result := abstract.Pow(5, 0)
		expected := 1
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("One Power", func(t *testing.T) {
		result := abstract.Pow(5, 1)
		expected := 5
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Negative Exponent", func(t *testing.T) {
		result := abstract.Pow(4.0, -1.0)
		expected := 0.25
		if math.Abs(float64(result)-expected) > 0.000001 {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Negative Base", func(t *testing.T) {
		result := abstract.Pow(-2, 2)
		expected := 4
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Negative Base Odd Exponent", func(t *testing.T) {
		result := abstract.Pow(-2, 3)
		expected := -8
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Fractional Exponent", func(t *testing.T) {
		result := abstract.Pow(4.0, 0.5)
		expected := 2.0
		if math.Abs(float64(result)-expected) > 0.000001 {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Zero Base", func(t *testing.T) {
		result := abstract.Pow(0, 5)
		expected := 0
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

// Test mixed type operations
func TestMixedTypes(t *testing.T) {
	t.Run("Pow with mixed types", func(t *testing.T) {
		baseInt := 2
		expFloat := 3.0
		expected := 8.0

		result := abstract.Pow(baseInt, expFloat)
		if math.Abs(float64(result)-expected) > 0.000001 {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Atoi with float conversion", func(t *testing.T) {
		input := "42"
		expected := 42.0

		result, err := abstract.Atoi[float64](input)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
