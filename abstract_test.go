package abstract_test

import (
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
