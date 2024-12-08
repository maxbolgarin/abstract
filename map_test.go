package abstract_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/maxbolgarin/abstract"
)

func TestNewMap(t *testing.T) {
	m := abstract.NewMap[string, int]()

	if m.Len() != 0 {
		t.Errorf("Expected map length to be 0, got %d", m.Len())
	}

	m.Set("key1", 1)
	if !m.Has("key1") {
		t.Errorf("Expected key 'key1' to be present")
	}
}

func TestNewMapFromPairs(t *testing.T) {
	m := abstract.NewMapFromPairs[string, int]("key1", 1, "key2", 2)

	if len := m.Len(); len != 2 {
		t.Errorf("Expected map length to be 2, got %d", len)
	}

	if val := m.Get("key1"); val != 1 {
		t.Errorf("Expected value for 'key1' to be 1, got %d", val)
	}
}

func TestGetAndLookup(t *testing.T) {
	m := abstract.NewMap[string, int](map[string]int{
		"key1": 100,
	})

	if val := m.Get("key1"); val != 100 {
		t.Errorf("Expected 'key1' to have value of 100, got %d", val)
	}

	val, ok := m.Lookup("key2")
	if ok || val != 0 {
		t.Errorf("Expected 'key2' to not exist, but got %v with value %d", ok, val)
	}
}

func TestSetAndDelete(t *testing.T) {
	m := abstract.NewMapWithSize[string, int](10)

	m.Set("key1", 100)
	m.Set("key1", 200) // overwrite

	if val := m.Get("key1"); val != 200 {
		t.Errorf("Expected 'key1' to have new value of 200, got %d", val)
	}

	deleted := m.Delete("key1")
	if !deleted {
		t.Errorf("Expected 'key1' to be deleted")
	}

	if m.Has("key1") {
		t.Errorf("Expected 'key1' to not be present after deletion")
	}

	deleted = m.Delete("key2")
	if deleted {
		t.Errorf("Expected 'key2' to not be deleted")
	}
}

func TestPop(t *testing.T) {
	m := abstract.NewMap[string, int]()
	m.Set("key1", 100)

	val := m.Pop("key1")
	if val != 100 {
		t.Errorf("Expected to pop value 100, got %d", val)
	}

	if m.Has("key1") {
		t.Errorf("Expected 'key1' to be removed after pop")
	}
}

func TestKeysAndValues(t *testing.T) {
	m := abstract.NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	keys := m.Keys()
	values := m.Values()

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}
}

func TestIsEmpty(t *testing.T) {
	m := abstract.NewMap[string, int]()

	if !m.IsEmpty() {
		t.Errorf("Expected map to be empty initially")
	}

	m.Set("key", 10)
	if m.IsEmpty() {
		t.Errorf("Expected map not to be empty after adding an item")
	}

	m.Delete("key")
	if !m.IsEmpty() {
		t.Errorf("Expected map to be empty after deleting the item")
	}
}

func TestSwap(t *testing.T) {
	m := abstract.NewMap[string, int]()
	m.Set("key1", 100)

	oldValue := m.Swap("key1", 200)
	if oldValue != 100 {
		t.Errorf("Expected old value to be 100, got %d", oldValue)
	}

	if newVal := m.Get("key1"); newVal != 200 {
		t.Errorf("Expected new value to be 200, got %d", newVal)
	}
}

func TestSetIfNotPresent(t *testing.T) {
	m := abstract.NewMap[string, int]()
	m.Set("key1", 100)

	existedValue := m.SetIfNotPresent("key1", 200)
	if existedValue != 100 {
		t.Errorf("Expected existing value to be 100, got %d", existedValue)
	}

	newValue := m.SetIfNotPresent("key2", 300)
	if newValue != 300 {
		t.Errorf("Expected new value to be set to 300, got %d", newValue)
	}
}

func TestTransform(t *testing.T) {
	m := abstract.NewMap[string, int]()
	m.Set("key1", 1)
	m.Set("key2", 2)

	m.Transform(func(k string, v int) int {
		return v * 2
	})

	if v := m.Get("key1"); v != 2 {
		t.Errorf("Expected value for 'key1' to be transformed to 2, got %d", v)
	}
	if v := m.Get("key2"); v != 4 {
		t.Errorf("Expected value for 'key2' to be transformed to 4, got %d", v)
	}
}

func TestRange(t *testing.T) {
	m := abstract.NewMap[string, int]()
	m.Set("key1", 1)
	m.Set("key2", 2)

	if m.Range(func(k string, v int) bool {
		if k != "key1" && k != "key2" {
			t.Errorf("Expected to visit key 'key1' and 'key2', got %s", k)
		}
		if v == 2 {
			return false
		}
		return true
	}) {
		t.Error("Expected Range to return false, but got true")
	}

	if !m.Range(func(k string, v int) bool {
		return true
	}) {
		t.Error("Expected Range to return true, but got false")
	}
}

func TestCopy(t *testing.T) {
	m := abstract.NewMap[string, int]()
	m.Set("key1", 1)

	copyMap := m.Copy()
	copyMap["key1"] = 10 // Modify the copy

	// Check original is unchanged
	if original := m.Get("key1"); original != 1 {
		t.Errorf("Expected original map value for 'key1' to be 1, got %d", original)
	}
}

func TestSafeMap_NewSafeMap(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	if m.Len() != 0 {
		t.Errorf("Expected map length to be 0, got %d", m.Len())
	}
}

func TestNewSafeMapFromPairs(t *testing.T) {
	m := abstract.NewSafeMapFromPairs[string, int]("key1", 1, "key2", 2)

	if len := m.Len(); len != 2 {
		t.Errorf("Expected map length to be 2, got %d", len)
	}

	if val := m.Get("key1"); val != 1 {
		t.Errorf("Expected value for 'key1' to be 1, got %d", val)
	}
}

func TestSafeMap_SetAndGet(t *testing.T) {
	m := abstract.NewSafeMap[string, int](map[string]int{
		"key1": 10,
	})

	if value := m.Get("key1"); value != 10 {
		t.Errorf("Expected value 10, got %d", value)
	}
}

func TestSafeMap_Lookup(t *testing.T) {
	m := abstract.NewSafeMapWithSize[string, int](2)

	m.Set("key1", 10)

	if value, ok := m.Lookup("key1"); !ok || value != 10 {
		t.Errorf("Expected value 10, got %d, ok %v", value, ok)
	}

	if _, ok := m.Lookup("key2"); ok {
		t.Errorf("Expected key2 to be absent")
	}
}

func TestSafeMap_Has(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 10)

	if !m.Has("key1") {
		t.Errorf("Expected key1 to be present")
	}

	if m.Has("key2") {
		t.Errorf("Expected key2 to be absent")
	}
}

func TestSafeMap_Delete(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 10)

	if !m.Delete("key1") {
		t.Errorf("Expected successful deletion of key1")
	}

	if m.Has("key1") {
		t.Errorf("Expected key1 to be deleted")
	}

	if m.Delete("key2") {
		t.Errorf("Expected failed deletion of key2")
	}
}

func TestSafeMap_Empty(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()

	if !m.IsEmpty() {
		t.Errorf("Expected map to be empty")
	}

	m.Set("key1", 10)
	if m.IsEmpty() {
		t.Errorf("Expected map to not be empty")
	}
}

func TestSafeMap_Len(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 10)
	m.Set("key2", 20)

	if m.Len() != 2 {
		t.Errorf("Expected map length to be 2, got %d", m.Len())
	}
}

func TestSafeMap_Pop(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 100)

	val := m.Pop("key1")
	if val != 100 {
		t.Errorf("Expected to pop value 100, got %d", val)
	}

	if m.Has("key1") {
		t.Errorf("Expected 'key1' to be removed after pop")
	}
}

func TestSafeMap_SetIfNotPresent(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 100)

	existedValue := m.SetIfNotPresent("key1", 200)
	if existedValue != 100 {
		t.Errorf("Expected existing value to be 100, got %d", existedValue)
	}

	newValue := m.SetIfNotPresent("key2", 300)
	if newValue != 300 {
		t.Errorf("Expected new value to be set to 300, got %d", newValue)
	}
}

func TestSafeMap_Swap(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 100)

	oldValue := m.Swap("key1", 200)
	if oldValue != 100 {
		t.Errorf("Expected old value to be 100, got %d", oldValue)
	}

	if newVal := m.Get("key1"); newVal != 200 {
		t.Errorf("Expected new value to be 200, got %d", newVal)
	}
}

func TestSafeMap_Keys(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 10)
	m.Set("key2", 20)

	keys := m.Keys()
	expectedKeys := []string{"key1", "key2"}

	for _, expectedKey := range expectedKeys {
		found := false
		for _, key := range keys {
			if key == expectedKey {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected key %s in keys slice", expectedKey)
		}
	}
}

func TestSafeMap_Values(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 10)
	m.Set("key2", 20)

	values := m.Values()
	expectedValues := []int{10, 20}

	for _, expectedValue := range expectedValues {
		found := false
		for _, value := range values {
			if value == expectedValue {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected value %d in values slice", expectedValue)
		}
	}
}

func TestSafeMap_ConcurrentAccess(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	var wg sync.WaitGroup

	const numGoroutines = 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			m.Set(key, i)
			if val := m.Get(key); val != i {
				t.Errorf("Expected value %d for key %s, got %d", i, key, val)
			}
		}(i)
	}

	wg.Wait()

	if m.Len() != numGoroutines {
		t.Errorf("Expected map length to be %d, got %d", numGoroutines, m.Len())
	}
}
func TestSafeMap_Transform(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 1)
	m.Set("key2", 2)

	m.Transform(func(k string, v int) int {
		return v * 2
	})

	if v := m.Get("key1"); v != 2 {
		t.Errorf("Expected value for 'key1' to be transformed to 2, got %d", v)
	}
	if v := m.Get("key2"); v != 4 {
		t.Errorf("Expected value for 'key2' to be transformed to 4, got %d", v)
	}
}

func TestSafeMap_Range(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 1)
	m.Set("key2", 2)

	if m.Range(func(k string, v int) bool {
		if k != "key1" && k != "key2" {
			t.Errorf("Expected to visit key 'key1' and 'key2', got %s", k)
		}
		if v == 2 {
			return false
		}
		return true
	}) {
		t.Error("Expected Range to return false, but got true")
	}

	if !m.Range(func(k string, v int) bool {
		return true
	}) {
		t.Error("Expected Range to return true, but got false")
	}
}

func TestSafeMap_Copy(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 1)

	copyMap := m.Copy()
	copyMap["key1"] = 10 // Modify the copy

	// Check original is unchanged
	if original := m.Get("key1"); original != 1 {
		t.Errorf("Expected original map value for 'key1' to be 1, got %d", original)
	}
}

func TestSafeMap_Clear(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 10)
	m.Set("key2", 20)

	m.Clear()
	if m.Len() != 0 {
		t.Errorf("Expected map to be clear, but got length %d", m.Len())
	}
}

func TestSafeMap_Refill(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	newData := map[string]int{"key3": 30, "key4": 40}

	m.Refill(newData)
	if m.Len() != 2 {
		t.Errorf("Expected map length to be 2 after refill, got %d", m.Len())
	}

	if val, ok := m.Lookup("key3"); !ok || val != 30 {
		t.Errorf("Expected key3 to have value 30, got %d", val)
	}

	if val, ok := m.Lookup("key4"); !ok || val != 40 {
		t.Errorf("Expected key4 to have value 40, got %d", val)
	}
}

// Define a simple Entity implementation for testing
type testEntity struct {
	id    int
	name  string
	order int
}

func (e *testEntity) GetID() int {
	return e.id
}

func (e *testEntity) GetName() string {
	return e.name
}

func (e *testEntity) GetOrder() int {
	return e.order
}

func (e *testEntity) SetOrder(order int) abstract.Entity[int] {
	e.order = order
	return e
}

func TestEntityMap_NewEntityMap(t *testing.T) {
	m := abstract.NewEntityMap[int, *testEntity]()
	if m.Len() != 0 {
		t.Errorf("Expected map length to be 0, got %d", m.Len())
	}
}

func TestEntityMap_SetAndGet(t *testing.T) {
	m := abstract.NewEntityMapWithSize[int, *testEntity](10)
	entity := &testEntity{id: 1, name: "Entity1", order: 0}

	m.Set(entity)
	if got := m.Get(1); got != entity {
		t.Errorf("Expected %v, got %v", entity, got)
	}
	entity.order = -1

	m.Set(entity)
	if got := m.Get(1); got == entity {
		t.Errorf("Expected empty entity, got %v", got)
	}
}

func TestEntityMap_LookupByName(t *testing.T) {
	m := abstract.NewEntityMap[int, *testEntity]()
	entity := &testEntity{id: 1, name: "Entity1", order: 0}

	m.Set(entity)

	if got, ok := m.LookupByName("Entity1"); !ok || got != entity {
		t.Errorf("Expected %v, got %v, ok %v", entity, got, ok)
	}

	if _, ok := m.LookupByName("Nonexistent"); ok {
		t.Error("Expected name to be absent")
	}
}

func TestEntityMap_AllOrdered(t *testing.T) {
	m := abstract.NewEntityMap[int, *testEntity]()
	entities := []*testEntity{
		{id: 1, name: "Entity1", order: 2},
		{id: 2, name: "Entity2", order: 0},
		{id: 3, name: "Entity3", order: 1},
	}

	for _, e := range entities {
		m.Set(e)
	}

	expectedOrder := []*testEntity{entities[1], entities[2], entities[0]}
	ordered := m.AllOrdered()

	for i, e := range expectedOrder {
		if ordered[i] != e {
			t.Errorf("Expected %v at position %d, got %v", e, i, ordered[i])
		}
	}
}

func TestEntityMap_NextOrder(t *testing.T) {
	m := abstract.NewEntityMap[int, *testEntity]()
	if order := m.NextOrder(); order != 0 {
		t.Errorf("Expected next order to be 0, got %d", order)
	}

	m.Set(&testEntity{id: 1, order: 0})
	if order := m.NextOrder(); order != 1 {
		t.Errorf("Expected next order to be 1, got %d", order)
	}
}

func TestEntityMap_ChangeOrder(t *testing.T) {
	m := abstract.NewEntityMap[int, *testEntity]()
	entities := []*testEntity{
		{id: 1, name: "Entity1", order: 2},
		{id: 2, name: "Entity2", order: 0},
		{id: 3, name: "Entity3", order: 1},
	}

	for _, e := range entities {
		m.Set(e)
	}

	newOrders := map[int]int{
		1: 0,
		2: 1,
		3: 2,
	}

	m.ChangeOrder(newOrders)
	expectedOrder := []*testEntity{entities[0], entities[1], entities[2]} // new orders applied
	ordered := m.AllOrdered()

	for i := range expectedOrder {
		if ordered[i].GetOrder() != newOrders[ordered[i].GetID()] {
			t.Errorf("Expected order for %v to be %d, got %d", ordered[i].GetName(), newOrders[ordered[i].GetID()], ordered[i].GetOrder())
		}
	}
}

func TestEntityMap_Delete(t *testing.T) {
	m := abstract.NewEntityMap[int, *testEntity]()
	entity := &testEntity{id: 1, name: "Entity1", order: 0}

	m.Set(entity)

	if !m.Delete(1) {
		t.Error("Expected deletion to be successful")
	}

	if m.Has(1) {
		t.Error("Expected the entity to be deleted")
	}

	entities := []*testEntity{
		{id: 1, name: "Entity1", order: 2},
		{id: 2, name: "Entity2", order: 0},
		{id: 3, name: "Entity3", order: 1},
		{id: 4, name: "Entity4", order: -10},
		{id: 5, name: "Entity5", order: -11},
	}

	for _, e := range entities {
		m.Set(e)
	}

	if !m.Delete(2) {
		t.Error("Expected deletion to be successful")
	}

	if m.Has(2) {
		t.Error("Expected the entity to be deleted")
	}

	if m.AllOrdered()[1].GetName() != "Entity1" {
		t.Errorf("Expected Entity1 at position 1, got %s", m.AllOrdered()[1].GetName())
	}
}

func TestSafeEntityMap_NewEntityMap(t *testing.T) {
	m := abstract.NewSafeEntityMap[int, *testEntity]()
	if m.Len() != 0 {
		t.Errorf("Expected map length to be 0, got %d", m.Len())
	}
}

func TestSafeEntityMap_SetAndGet(t *testing.T) {
	m := abstract.NewSafeEntityMapWithSize[int, *testEntity](10)
	entity := &testEntity{id: 1, name: "Entity1", order: 0}

	m.Set(entity)
	if got := m.Get(1); got != entity {
		t.Errorf("Expected %v, got %v", entity, got)
	}
	entity.order = -1

	m.Set(entity)
	if got := m.Get(1); got == entity {
		t.Errorf("Expected empty entity, got %v", got)
	}
}

func TestSafeEntityMap_LookupByName(t *testing.T) {
	m := abstract.NewSafeEntityMap[int, *testEntity]()
	entity := &testEntity{id: 1, name: "Entity1", order: 0}

	m.Set(entity)

	if got, ok := m.LookupByName("Entity1"); !ok || got != entity {
		t.Errorf("Expected %v, got %v, ok %v", entity, got, ok)
	}

	if _, ok := m.LookupByName("Nonexistent"); ok {
		t.Error("Expected name to be absent")
	}
}

func TestSafeEntityMap_AllOrdered(t *testing.T) {
	m := abstract.NewSafeEntityMap[int, *testEntity]()
	entities := []*testEntity{
		{id: 1, name: "Entity1", order: 2},
		{id: 2, name: "Entity2", order: 0},
		{id: 3, name: "Entity3", order: 1},
	}

	for _, e := range entities {
		m.Set(e)
	}

	expectedOrder := []*testEntity{entities[1], entities[2], entities[0]}
	ordered := m.AllOrdered()

	for i, e := range expectedOrder {
		if ordered[i] != e {
			t.Errorf("Expected %v at position %d, got %v", e, i, ordered[i])
		}
	}
}

func TestSafeEntityMap_NextOrder(t *testing.T) {
	m := abstract.NewSafeEntityMap[int, *testEntity]()
	if order := m.NextOrder(); order != 0 {
		t.Errorf("Expected next order to be 0, got %d", order)
	}

	m.Set(&testEntity{id: 1, order: 0})
	if order := m.NextOrder(); order != 1 {
		t.Errorf("Expected next order to be 1, got %d", order)
	}
}

func TestSafeEntityMap_ChangeOrder(t *testing.T) {
	m := abstract.NewSafeEntityMap[int, *testEntity]()
	entities := []*testEntity{
		{id: 1, name: "Entity1", order: 2},
		{id: 2, name: "Entity2", order: 0},
		{id: 3, name: "Entity3", order: 1},
	}

	for _, e := range entities {
		m.Set(e)
	}

	newOrders := map[int]int{
		1: 0,
		2: 1,
		3: 2,
	}

	m.ChangeOrder(newOrders)
	expectedOrder := []*testEntity{entities[0], entities[1], entities[2]} // new orders applied
	ordered := m.AllOrdered()

	for i := range expectedOrder {
		if ordered[i].GetOrder() != newOrders[ordered[i].GetID()] {
			t.Errorf("Expected order for %v to be %d, got %d", ordered[i].GetName(), newOrders[ordered[i].GetID()], ordered[i].GetOrder())
		}
	}
}

func TestSafeEntityMap_Delete(t *testing.T) {
	m := abstract.NewSafeEntityMap[int, *testEntity]()
	entity := &testEntity{id: 1, name: "Entity1", order: 0}

	m.Set(entity)

	if !m.Delete(1) {
		t.Error("Expected deletion to be successful")
	}

	if m.Has(1) {
		t.Error("Expected the entity to be deleted")
	}

	entities := []*testEntity{
		{id: 1, name: "Entity1", order: 2},
		{id: 2, name: "Entity2", order: 0},
		{id: 3, name: "Entity3", order: 1},
		{id: 4, name: "Entity4", order: -10},
		{id: 5, name: "Entity5", order: -11},
	}

	for _, e := range entities {
		m.Set(e)
	}

	if !m.Delete(2) {
		t.Error("Expected deletion to be successful")
	}

	if m.Has(2) {
		t.Error("Expected the entity to be deleted")
	}

	if m.AllOrdered()[1].GetName() != "Entity1" {
		t.Errorf("Expected Entity1 at position 1, got %s", m.AllOrdered()[1].GetName())
	}
}

func TestOrderedPairs_AddAndGet(t *testing.T) {
	pairs := abstract.NewOrderedPairs[int, string]()

	// Test adding elements
	pairs.Add(1, "one")
	pairs.Add(2, "two")
	pairs.Add(1, "uno") // Duplicate key with new value

	val := pairs.Get(1)
	if val != "uno" {
		t.Errorf("Expected value 'uno', but got %v", val)
	}

	val = pairs.Get(2)
	if val != "two" {
		t.Errorf("Expected value 'two', but got %v", val)
	}

	val = pairs.Get(3)
	if val != "" {
		t.Errorf("Expected empty string for non-existent key, but got %v", val)
	}
}

func TestOrderedPairs_Keys(t *testing.T) {
	pairs := abstract.NewOrderedPairs[int, string]()
	pairs.Add(1, "one")
	pairs.Add(2, "two")
	pairs.Add(1, "uno")

	keys := pairs.Keys()
	expectedKeys := []int{1, 2, 1}

	if len(keys) != len(expectedKeys) {
		t.Fatalf("Expected keys length %v, but got %v", len(expectedKeys), len(keys))
	}

	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Errorf("Expected key %v at index %v, but got %v", expectedKeys[i], i, key)
		}
	}
}

func TestOrderedPairs_Rand(t *testing.T) {
	pairs := abstract.NewOrderedPairs[int, string]()
	pairs.Add(1, "one")
	pairs.Add(2, "two")
	pairs.Add(3, "three")

	randomValue := pairs.Rand()
	if randomValue == "" {
		t.Error("Expected a random value from the set, but got an empty result")
	}

	// Should handle single element scenario
	singlePair := abstract.NewOrderedPairs[int, string]()
	singlePair.Add(1, "only")
	if randomValue := singlePair.Rand(); randomValue != "only" {
		t.Errorf("Expected 'only' for singleton pair, got %v", randomValue)
	}

	// Should handle empty scenario gracefully
	emptyPair := abstract.NewOrderedPairs[int, string]()
	if randomValue := emptyPair.Rand(); randomValue != "" {
		t.Errorf("Expected empty value for empty pair map, got %v", randomValue)
	}
}

func TestOrderedPairs_RandKey(t *testing.T) {
	pairs := abstract.NewOrderedPairs[int, string](1, "one", 2, "two", 3, "three")

	randomKey := pairs.RandKey()
	if (randomKey > 3) || (randomKey < 1) {
		t.Errorf("Expected random key from 1 to 3, but got %v", randomKey)
	}

	// Test with a single key
	singleKeyPair := abstract.NewOrderedPairs[int, string]()
	singleKeyPair.Add(1, "only")
	if randomKey := singleKeyPair.RandKey(); randomKey != 1 {
		t.Errorf("Expected key '1' for single key pair, got %v", randomKey)
	}

	// Test with an empty OrderedPairs
	emptyPair := abstract.NewOrderedPairs[int, string]()
	if randomKey := emptyPair.RandKey(); randomKey != 0 {
		t.Errorf("Expected zero value for empty pair map, got %v", randomKey)
	}
}

func TestSafeOrderedPairs_AddAndGet(t *testing.T) {
	pairs := abstract.NewSafeOrderedPairs[int, string]()

	// Test adding elements
	pairs.Add(1, "one")
	pairs.Add(2, "two")
	pairs.Add(1, "uno") // Duplicate key with new value

	val := pairs.Get(1)
	if val != "uno" {
		t.Errorf("Expected value 'uno', but got %v", val)
	}

	val = pairs.Get(2)
	if val != "two" {
		t.Errorf("Expected value 'two', but got %v", val)
	}

	val = pairs.Get(3)
	if val != "" {
		t.Errorf("Expected empty string for non-existent key, but got %v", val)
	}
}

func TesSafeOrderedPairs_Keys(t *testing.T) {
	pairs := abstract.NewSafeOrderedPairs[int, string]()
	pairs.Add(1, "one")
	pairs.Add(2, "two")
	pairs.Add(1, "uno")

	keys := pairs.Keys()
	expectedKeys := []int{1, 2, 1}

	if len(keys) != len(expectedKeys) {
		t.Fatalf("Expected keys length %v, but got %v", len(expectedKeys), len(keys))
	}

	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Errorf("Expected key %v at index %v, but got %v", expectedKeys[i], i, key)
		}
	}
}

func TestSafeOrderedPairs_Rand(t *testing.T) {
	pairs := abstract.NewSafeOrderedPairs[int, string]()
	pairs.Add(1, "one")
	pairs.Add(2, "two")
	pairs.Add(3, "three")

	randomValue := pairs.Rand()
	if randomValue == "" {
		t.Error("Expected a random value from the set, but got an empty result")
	}

	// Should handle single element scenario
	singlePair := abstract.NewOrderedPairs[int, string]()
	singlePair.Add(1, "only")
	if randomValue := singlePair.Rand(); randomValue != "only" {
		t.Errorf("Expected 'only' for singleton pair, got %v", randomValue)
	}

	// Should handle empty scenario gracefully
	emptyPair := abstract.NewOrderedPairs[int, string]()
	if randomValue := emptyPair.Rand(); randomValue != "" {
		t.Errorf("Expected empty value for empty pair map, got %v", randomValue)
	}
}

func TestSafeOrderedPairs_RandKey(t *testing.T) {
	pairs := abstract.NewSafeOrderedPairs[int, string](1, "one", 2, "two", 3, "three")

	randomKey := pairs.RandKey()
	if (randomKey > 3) || (randomKey < 1) {
		t.Errorf("Expected random key from 1 to 3, but got %v", randomKey)
	}

	// Test with a single key
	singleKeyPair := abstract.NewOrderedPairs[int, string]()
	singleKeyPair.Add(1, "only")
	if randomKey := singleKeyPair.RandKey(); randomKey != 1 {
		t.Errorf("Expected key '1' for single key pair, got %v", randomKey)
	}

	// Test with an empty OrderedPairs
	emptyPair := abstract.NewOrderedPairs[int, string]()
	if randomKey := emptyPair.RandKey(); randomKey != 0 {
		t.Errorf("Expected zero value for empty pair map, got %v", randomKey)
	}
}
