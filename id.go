package abstract

import "strconv"

// EntityType represents a type identifier for entities in the system.
// It's used as a prefix for generated IDs to ensure type safety and
// easy identification of entity types from their IDs.
//
// Example usage:
//
//	userType := RegisterEntityType("USER")
//	adminType := RegisterEntityType("ADMN")
//
//	userID := NewID(userType)   // "USER" + random string
//	adminID := NewID(adminType) // "ADMN" + random string
type EntityType string

// String returns the string representation of the EntityType.
// This method allows EntityType to be used as a string in contexts
// where string representation is needed.
//
// Returns:
//   - The string value of the EntityType
func (e EntityType) String() string {
	return string(e)
}

const (
	// TestIDEntity is a predefined EntityType for testing purposes.
	// It uses "00x0" as the type identifier.
	TestIDEntity EntityType = "00x0"

	// defaultIDSize is the default length of the random portion of generated IDs.
	// The total ID length will be entityTypeSize + defaultIDSize.
	defaultIDSize = 12
)

// entityTypeSize specifies the required length for entity type identifiers.
// This can be modified using SetEntitySize() but should be consistent
// across your application.
var entityTypeSize = 4

// RegisterEntityType creates a new EntityType with the specified identifier.
// The identifier must be exactly entityTypeSize characters long.
//
// This function is typically called during application initialization
// to define the entity types used in your system.
//
// Parameters:
//   - entityType: A string identifier for the entity type (must be exactly entityTypeSize characters)
//
// Returns:
//   - An EntityType that can be used with NewID() and other ID functions
//
// Panics:
//   - If the entityType length doesn't match entityTypeSize
//
// Example usage:
//
//	const (
//		UserEntity = RegisterEntityType("USER")
//		PostEntity = RegisterEntityType("POST")
//		AdminEntity = RegisterEntityType("ADMN")
//	)
//
//	// Or in init function:
//	func init() {
//		UserEntity = RegisterEntityType("USER")
//		PostEntity = RegisterEntityType("POST")
//	}
func RegisterEntityType(entityType string) EntityType {
	if len(entityType) != entityTypeSize {
		panic("entity type must be " + strconv.Itoa(entityTypeSize) + " characters long")
	}
	return EntityType(entityType)
}

// SetEntitySize configures the required length for entity type identifiers.
// This should be called before registering any entity types and should be
// consistent across your entire application.
//
// Parameters:
//   - size: The required length for entity type identifiers
//
// Example usage:
//
//	func init() {
//		SetEntitySize(5)  // Use 5-character entity types
//		UserEntity := RegisterEntityType("USER_")
//		AdminEntity := RegisterEntityType("ADMIN")
//	}
func SetEntitySize(size int) {
	entityTypeSize = size
}

// init registers the TestIDEntity during package initialization.
// This ensures the test entity type is always available.
func init() {
	RegisterEntityType(TestIDEntity.String())
}

// NewID generates a new unique identifier with the specified entity type prefix.
// The ID consists of the entity type followed by a random string of defaultIDSize length.
//
// Parameters:
//   - entityType: The EntityType to use as a prefix for the ID
//
// Returns:
//   - A unique identifier string in the format: entityType + randomString
//
// Example usage:
//
//	userType := RegisterEntityType("USER")
//	postType := RegisterEntityType("POST")
//
//	userID := NewID(userType)  // "USERa1b2c3d4e5f6"
//	postID := NewID(postType)  // "POST9z8y7x6w5v4u"
//
//	// Use in your application
//	type User struct {
//		ID   string
//		Name string
//	}
//
//	user := User{
//		ID:   NewID(userType),
//		Name: "John Doe",
//	}
func NewID(entityType EntityType) string {
	return entityType.String() + GetRandomString(defaultIDSize)
}

// NewTestID generates a new test identifier using the TestIDEntity type.
// This is a convenience function for testing and development purposes.
//
// Returns:
//   - A test identifier string in the format: "00x0" + randomString
//
// Example usage:
//
//	testID := NewTestID()  // "00x0a1b2c3d4e5f6"
//
//	// Use in tests
//	func TestUserCreation(t *testing.T) {
//		user := User{
//			ID:   NewTestID(),
//			Name: "Test User",
//		}
//		// ... test logic
//	}
func NewTestID() string {
	return NewID(TestIDEntity)
}

// FromID converts an existing ID to use a different entity type.
// This function replaces the entity type prefix of an ID while
// preserving the random portion (or as much as possible).
//
// Parameters:
//   - id: The existing ID to convert
//   - t: The new EntityType to use
//
// Returns:
//   - A new ID with the specified entity type prefix
//
// Example usage:
//
//	userType := RegisterEntityType("USER")
//	adminType := RegisterEntityType("ADMN")
//
//	userID := NewID(userType)                    // "USERa1b2c3d4e5f6"
//	adminID := FromID(userID, adminType)         // "ADMNa1b2c3d4e5f6"
//
//	// Useful for type conversions or migrations
//	func promoteUserToAdmin(userID string) string {
//		return FromID(userID, AdminEntity)
//	}
func FromID(id string, t EntityType) string {
	if len(id) <= len(t) {
		return t.String() + id
	}
	return t.String() + id[len(t):]
}

// FetchEntityType extracts the entity type from an ID.
// This function returns the entity type prefix from an ID string.
//
// Parameters:
//   - id: The ID string to extract the entity type from
//
// Returns:
//   - The EntityType extracted from the ID prefix
//
// Example usage:
//
//	userType := RegisterEntityType("USER")
//	userID := NewID(userType)                    // "USERa1b2c3d4e5f6"
//	extractedType := FetchEntityType(userID)     // "USER"
//
//	// Use for type checking or routing
//	func handleEntity(id string) {
//		entityType := FetchEntityType(id)
//		switch entityType {
//		case "USER":
//			handleUser(id)
//		case "POST":
//			handlePost(id)
//		case "ADMN":
//			handleAdmin(id)
//		}
//	}
func FetchEntityType(id string) EntityType {
	if len(id) < entityTypeSize {
		return EntityType(id)
	}
	return EntityType(id[:entityTypeSize])
}

// Builder provides a convenient way to generate IDs of a specific entity type.
// It encapsulates an EntityType and provides methods to generate IDs without
// needing to pass the entity type each time.
//
// Example usage:
//
//	userType := RegisterEntityType("USER")
//	userBuilder := WithEntityType(userType)
//
//	// Generate multiple user IDs
//	user1ID := userBuilder.NewID()  // "USERa1b2c3d4e5f6"
//	user2ID := userBuilder.NewID()  // "USERx9y8z7w6v5u4"
//	user3ID := userBuilder.NewID()  // "USERm3n4o5p6q7r8"
type Builder struct {
	t EntityType
}

// WithEntityType creates a new Builder for the specified entity type.
// This is a factory function that returns a Builder configured to generate
// IDs of the specified type.
//
// Parameters:
//   - t: The EntityType to use for ID generation
//
// Returns:
//   - A Builder instance configured for the specified entity type
//
// Example usage:
//
//	userType := RegisterEntityType("USER")
//	postType := RegisterEntityType("POST")
//
//	userBuilder := WithEntityType(userType)
//	postBuilder := WithEntityType(postType)
//
//	// Use builders to generate IDs
//	users := make([]User, 5)
//	for i := range users {
//		users[i] = User{
//			ID:   userBuilder.NewID(),
//			Name: fmt.Sprintf("User %d", i+1),
//		}
//	}
func WithEntityType(t EntityType) Builder {
	return Builder{
		t: t,
	}
}

// NewID generates a new ID using the Builder's configured entity type.
// This method provides a convenient way to generate multiple IDs of the
// same type without repeatedly specifying the entity type.
//
// Returns:
//   - A new ID string with the Builder's entity type prefix
//
// Example usage:
//
//	userBuilder := WithEntityType(RegisterEntityType("USER"))
//
//	// Generate multiple user IDs
//	for i := 0; i < 10; i++ {
//		userID := userBuilder.NewID()
//		fmt.Printf("Generated user ID: %s\n", userID)
//	}
//
//	// Use in struct initialization
//	type UserService struct {
//		idBuilder Builder
//	}
//
//	func NewUserService() *UserService {
//		return &UserService{
//			idBuilder: WithEntityType(RegisterEntityType("USER")),
//		}
//	}
//
//	func (s *UserService) CreateUser(name string) User {
//		return User{
//			ID:   s.idBuilder.NewID(),
//			Name: name,
//		}
//	}
func (b Builder) NewID() string {
	return NewID(b.t)
}
