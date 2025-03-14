package abstract_test

import (
	"testing"

	"github.com/maxbolgarin/abstract"
)

// TestEntityTypeString tests the String method of EntityType
func TestEntityTypeString(t *testing.T) {
	entityType := abstract.EntityType("TEST")
	if entityType.String() != "TEST" {
		t.Errorf("Expected 'TEST', got '%s'", entityType.String())
	}
}

// TestRegisterEntityType tests the RegisterEntityType function
func TestRegisterEntityType(t *testing.T) {
	// Set entity size to 4 for this test
	abstract.SetEntitySize(4)

	entityType := abstract.RegisterEntityType("ABCD")
	if entityType.String() != "ABCD" {
		t.Errorf("Expected 'ABCD', got '%s'", entityType.String())
	}
}

// TestRegisterEntityTypePanic tests that RegisterEntityType panics with incorrect size
func TestRegisterEntityTypePanic(t *testing.T) {
	abstract.SetEntitySize(4)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when it should have")
		}
	}()

	abstract.RegisterEntityType("ABC") // Should panic as it's not 4 chars
}

// TestSetEntitySize tests the SetEntitySize function
func TestSetEntitySize(t *testing.T) {
	// First set it to a different value
	abstract.SetEntitySize(5)

	// Test that registering with new size works
	entityType := abstract.RegisterEntityType("ABCDE")
	if entityType.String() != "ABCDE" {
		t.Errorf("Expected 'ABCDE', got '%s'", entityType.String())
	}

	// Reset to 4 for other tests
	abstract.SetEntitySize(4)
}

// TestNew tests the New function
func TestNew(t *testing.T) {
	abstract.SetEntitySize(4)
	entityType := abstract.EntityType("TEST")
	id := abstract.New(entityType)

	// Check the prefix is our entity type
	if id[:4] != "TEST" {
		t.Errorf("Expected ID to start with 'TEST', got '%s'", id[:4])
	}

	// Check the total length (entity type + default ID size)
	expectedLength := 4 + 12 // entityTypeSize + defaultIDSize
	if len(id) != expectedLength {
		t.Errorf("Expected ID length to be %d, got %d", expectedLength, len(id))
	}
}

// TestNewTest tests the NewTest function
func TestNewTest(t *testing.T) {
	id := abstract.NewTest()

	// Check the prefix is the Test entity type
	if id[:4] != abstract.Test.String() {
		t.Errorf("Expected ID to start with '%s', got '%s'", abstract.Test.String(), id[:4])
	}

	// Check the total length
	expectedLength := 4 + 12 // entityTypeSize + defaultIDSize
	if len(id) != expectedLength {
		t.Errorf("Expected ID length to be %d, got %d", expectedLength, len(id))
	}
}

// TestFrom tests the From function
func TestFrom(t *testing.T) {
	abstract.SetEntitySize(4)
	originalType := abstract.EntityType("ORIG")
	newType := abstract.EntityType("NEWS")

	originalID := abstract.New(originalType)
	newID := abstract.From(originalID, newType)

	// Check the prefix changed
	if newID[:4] != "NEWS" {
		t.Errorf("Expected ID to start with 'NEWS', got '%s'", newID[:4])
	}

	// Check the random part remains the same
	if originalID[4:] != newID[4:] {
		t.Errorf("Expected random part to be the same, got '%s' vs '%s'", originalID[4:], newID[4:])
	}
}

// TestFromShort tests the From function with a short ID
func TestFromShort(t *testing.T) {
	abstract.SetEntitySize(4)
	newType := abstract.EntityType("NEWS")

	shortID := "123" // Shorter than entity type
	newID := abstract.From(shortID, newType)

	expected := "NEWS123"
	if newID != expected {
		t.Errorf("Expected '%s', got '%s'", expected, newID)
	}
}

// TestFetchEntityType tests the FetchEntityType function
func TestFetchEntityType(t *testing.T) {
	abstract.SetEntitySize(4)
	entityType := abstract.EntityType("TEST")
	id := abstract.New(entityType)

	fetchedType := abstract.FetchEntityType(id)
	if fetchedType.String() != "TEST" {
		t.Errorf("Expected 'TEST', got '%s'", fetchedType.String())
	}
}

// TestFetchEntityTypeShort tests the FetchEntityType function with a short ID
func TestFetchEntityTypeShort(t *testing.T) {
	abstract.SetEntitySize(4)
	shortID := "AB"

	fetchedType := abstract.FetchEntityType(shortID)
	if fetchedType.String() != "AB" {
		t.Errorf("Expected 'AB', got '%s'", fetchedType.String())
	}
}

// TestWith tests the With function and Builder.New
func TestWith(t *testing.T) {
	abstract.SetEntitySize(4)
	entityType := abstract.EntityType("CUST")

	builder := abstract.With(entityType)
	id := builder.New()

	// Check the prefix is our entity type
	if id[:4] != "CUST" {
		t.Errorf("Expected ID to start with 'CUST', got '%s'", id[:4])
	}

	// Check the total length
	expectedLength := 4 + 12 // entityTypeSize + defaultIDSize
	if len(id) != expectedLength {
		t.Errorf("Expected ID length to be %d, got %d", expectedLength, len(id))
	}
}
