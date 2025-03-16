package abstract

import "strconv"

type EntityType string

func (e EntityType) String() string {
	return string(e)
}

const (
	TestIDEntity EntityType = "00x0"

	defaultIDSize = 12
)

var entityTypeSize = 4

func RegisterEntityType(entityType string) EntityType {
	if len(entityType) != entityTypeSize {
		panic("entity type must be " + strconv.Itoa(entityTypeSize) + " characters long")
	}
	return EntityType(entityType)
}

func SetEntitySize(size int) {
	entityTypeSize = size
}

func init() {
	RegisterEntityType(TestIDEntity.String())
}

// NewID is used to generate a new ID based on provided entity type.
func NewID(entityType EntityType) string {
	return entityType.String() + GetRandomString(defaultIDSize)
}

// NewTestID is used to generate a new ID based on Test entity type.
func NewTestID() string {
	return NewID(TestIDEntity)
}

// FromID changes entity type for the provided ID.
func FromID(id string, t EntityType) string {
	if len(id) <= len(t) {
		return t.String() + id
	}
	return t.String() + id[len(t):]
}

// FetchEntityType is used to extract entity type from provided ID.
func FetchEntityType(id string) EntityType {
	if len(id) < entityTypeSize {
		return EntityType(id)
	}
	return EntityType(id[:entityTypeSize])
}

type Builder struct {
	t EntityType
}

func WithEntityType(t EntityType) Builder {
	return Builder{
		t: t,
	}
}

func (b Builder) NewID() string {
	return NewID(b.t)
}
