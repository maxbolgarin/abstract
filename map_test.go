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
	m := abstract.NewMap(map[string]int{
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
	m.Set("key3", 300)
	m.Set("key4", 400)

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

	if !m.Has("key3") {
		t.Errorf("Expected 'key3' to be present")
	}

	if !m.Has("key4") {
		t.Errorf("Expected 'key4' to be present")
	}

	deleted = m.Delete("key3", "key4")
	if !deleted {
		t.Errorf("Expected 'key3' and 'key4' to be deleted")
	}

	if m.Has("key3") {
		t.Errorf("Expected 'key3' to not be present after deletion")
	}

	if m.Has("key4") {
		t.Errorf("Expected 'key4' to not be present after deletion")
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

func TestChange(t *testing.T) {
	m := abstract.NewMap[string, int]()
	m.Set("key1", 1)
	m.Change("key1", func(k string, v int) int {
		return v * 2
	})

	if v := m.Get("key1"); v != 2 {
		t.Errorf("Expected value for 'key1' to be transformed to 2, got %d", v)
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

func TestClear(t *testing.T) {
	m := abstract.NewMap[string, int]()
	m.Set("key1", 1)

	m.Clear()
	if m.Len() != 0 {
		t.Errorf("Expected map to be clear, but got length %d", m.Len())
	}
}

func TestMapIter(t *testing.T) {
	m := abstract.NewMap[string, int]()
	m.Set("key1", 1)
	m.Set("key2", 2)
	iter := m.Iter()
	for k, v := range iter {
		if k != "key1" && k != "key2" {
			t.Errorf("Expected to visit key 'key1' and 'key2', got %s", k)
		}
		if v != 1 && v != 2 {
			t.Errorf("Expected to visit value 1 and 2, got %d", v)
		}
	}

	iter2 := m.IterKeys()
	for k := range iter2 {
		if k != "key1" && k != "key2" {
			t.Errorf("Expected to visit key 'key1' and 'key2', got %s", k)
		}
	}

	iter3 := m.IterValues()
	for v := range iter3 {
		if v != 1 && v != 2 {
			t.Errorf("Expected to visit value 1 and 2, got %d", v)
		}
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
	m := abstract.NewSafeMap(map[string]int{
		"key1": 10,
	})

	if value := m.Get("key1"); value != 10 {
		t.Errorf("Expected value 10, got %d", value)
	}
}

func TestSafeMap_Lookup(t *testing.T) {
	m := &abstract.SafeMap[string, int]{}

	if value, ok := m.Lookup("key1"); ok || value != 0 {
		t.Errorf("Expected value 0, got %d, ok %v", value, ok)
	}

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
	m.Set("key1", 100)
	m.Set("key1", 200) // overwrite
	m.Set("key3", 300)
	m.Set("key4", 400)

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

	if !m.Has("key3") {
		t.Errorf("Expected 'key3' to be present")
	}

	if !m.Has("key4") {
		t.Errorf("Expected 'key4' to be present")
	}

	deleted = m.Delete("key3", "key4")
	if !deleted {
		t.Errorf("Expected 'key3' and 'key4' to be deleted")
	}

	if m.Has("key3") {
		t.Errorf("Expected 'key3' to not be present after deletion")
	}

	if m.Has("key4") {
		t.Errorf("Expected 'key4' to not be present after deletion")
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

func TestSafeMap_Change(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 1)
	m.Change("key1", func(k string, v int) int {
		return v * 2
	})

	if v := m.Get("key1"); v != 2 {
		t.Errorf("Expected value for 'key1' to be transformed to 2, got %d", v)
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

func TestSafeMap_Iter(t *testing.T) {
	m := abstract.NewSafeMap[string, int]()
	m.Set("key1", 1)
	m.Set("key2", 2)
	iter := m.Iter()
	for k, v := range iter {
		if k != "key1" && k != "key2" {
			t.Errorf("Expected to visit key 'key1' and 'key2', got %s", k)
		}
		if v != 1 && v != 2 {
			t.Errorf("Expected to visit value 1 and 2, got %d", v)
		}
	}
	iter2 := m.IterKeys()
	for k := range iter2 {
		if k != "key1" && k != "key2" {
			t.Errorf("Expected to visit key 'key1' and 'key2', got %s", k)
		}
	}
	iter3 := m.IterValues()
	for v := range iter3 {
		if v != 1 && v != 2 {
			t.Errorf("Expected to visit value 1 and 2, got %d", v)
		}
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
	entity := &testEntity{id: 1, name: "Entity1"}

	m.Set(entity)
	if got := m.Get(1); got != entity {
		t.Errorf("Expected %v, got %v", entity, got)
	}
	order := m.Set(&testEntity{id: 2, name: "Entity2"})
	if order != 1 {
		t.Errorf("Expected order to be 1, got %d", order)
	}
	if got := m.Get(2); got.order != 1 {
		t.Errorf("Expected order to be 1, got %d", got.order)
	}
	m.Set(entity)
	if got := m.Get(1); got.order != 0 {
		t.Errorf("Expected order to be 0, got %d", got.order)
	}
}

func TestEntityMap_SetManualOrderAndGet(t *testing.T) {
	m := abstract.NewEntityMapWithSize[int, *testEntity](10)
	Entity1 := &testEntity{id: 1, name: "Entity1"}
	Entity2 := &testEntity{id: 2, name: "Entity2"}
	Entity3 := &testEntity{id: 3, name: "Entity3"}

	order := m.SetManualOrder(Entity1)
	if order != 0 {
		t.Errorf("Expected order to be 0, got %d", order)
	}
	if got := m.Get(1); got != Entity1 {
		t.Errorf("Expected %v, got %v", Entity1, got)
	}
	m.SetManualOrder(Entity2)
	if got := m.Get(2); got.order != 0 {
		t.Errorf("Expected order to be 0, got %d", got.order)
	}
	m.SetManualOrder(Entity3)
	if got := m.Get(2); got.order != 0 {
		t.Errorf("Expected order to be 0, got %d", got.order)
	}
	ordered := m.AllOrdered()
	if len(ordered) != 3 {
		t.Errorf("Expected 3 entities, got %d", len(ordered))
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

	expectedOrder := []*testEntity{entities[0], entities[1], entities[2]}
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

	if m.AllOrdered()[1].GetName() != "Entity3" {
		t.Errorf("Expected Entity3 at position 1, got %s", m.AllOrdered()[1].GetName())
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
	entity = &testEntity{id: 1, name: "Entity1", order: -1}

	order := m.Set(entity)
	if order != 0 {
		t.Error("Expected order to be 0")
	}
	if got := m.Get(1); got.order != 0 {
		t.Errorf("Expected order to be 0, got %d", got.order)
	}
}

func TestSafeEntityMap_SetManualOrderAndGet(t *testing.T) {
	m := abstract.NewSafeEntityMapWithSize[int, *testEntity](10)
	Entity1 := &testEntity{id: 1, name: "Entity1"}
	Entity2 := &testEntity{id: 2, name: "Entity2"}
	Entity3 := &testEntity{id: 3, name: "Entity3"}

	order := m.SetManualOrder(Entity1)
	if order != 0 {
		t.Error("Expected order to be 0")
	}
	if got := m.Get(1); got != Entity1 {
		t.Errorf("Expected %v, got %v", Entity1, got)
	}
	m.SetManualOrder(Entity2)
	if got := m.Get(2); got.order != 0 {
		t.Errorf("Expected order to be 0, got %d", got.order)
	}
	m.SetManualOrder(Entity3)
	if got := m.Get(2); got.order != 0 {
		t.Errorf("Expected order to be 0, got %d", got.order)
	}
	ordered := m.AllOrdered()
	if len(ordered) != 3 {
		t.Errorf("Expected 3 entities, got %d", len(ordered))
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

	expectedOrder := []*testEntity{entities[0], entities[1], entities[2]}
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

	if m.AllOrdered()[1].GetName() != "Entity3" {
		t.Errorf("Expected Entity3 at position 1, got %s", m.AllOrdered()[1].GetName())
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

// Tests for MapOfMaps[K1, K2, V]

func TestMapOfMaps_NewMapOfMaps(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	if m.Len() != 0 {
		t.Errorf("Expected map length to be 0, got %d", m.Len())
	}
	if m.OuterLen() != 0 {
		t.Errorf("Expected outer map length to be 0, got %d", m.OuterLen())
	}
	if !m.IsEmpty() {
		t.Error("Expected map to be empty")
	}
}

func TestMapOfMaps_NewMapOfMapsWithSize(t *testing.T) {
	m := abstract.NewMapOfMapsWithSize[string, int, float64](10)
	if m.Len() != 0 {
		t.Errorf("Expected map length to be 0, got %d", m.Len())
	}
}

func TestMapOfMaps_NewMapOfMapsFromExisting(t *testing.T) {
	existing := map[string]map[int]float64{
		"group1": {1: 1.1, 2: 2.2},
		"group2": {3: 3.3, 4: 4.4},
	}
	m := abstract.NewMapOfMaps(existing)

	if m.Len() != 4 {
		t.Errorf("Expected total length to be 4, got %d", m.Len())
	}
	if m.OuterLen() != 2 {
		t.Errorf("Expected outer length to be 2, got %d", m.OuterLen())
	}
}

func TestMapOfMaps_SetAndGet(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()

	// Test setting and getting values
	m.Set("users", 1, 10.5)
	m.Set("users", 2, 20.7)
	m.Set("products", 100, 99.99)

	if val := m.Get("users", 1); val != 10.5 {
		t.Errorf("Expected value 10.5, got %f", val)
	}

	if val := m.Get("users", 2); val != 20.7 {
		t.Errorf("Expected value 20.7, got %f", val)
	}

	if val := m.Get("products", 100); val != 99.99 {
		t.Errorf("Expected value 99.99, got %f", val)
	}

	// Test getting non-existent values
	if val := m.Get("nonexistent", 1); val != 0.0 {
		t.Errorf("Expected default value 0.0, got %f", val)
	}
}

func TestMapOfMaps_Lookup(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)

	// Test existing value
	val, ok := m.Lookup("users", 1)
	if !ok || val != 10.5 {
		t.Errorf("Expected value 10.5 and ok=true, got %f and ok=%v", val, ok)
	}

	// Test non-existent outer key
	val, ok = m.Lookup("nonexistent", 1)
	if ok || val != 0.0 {
		t.Errorf("Expected default value 0.0 and ok=false, got %f and ok=%v", val, ok)
	}

	// Test non-existent inner key
	val, ok = m.Lookup("users", 999)
	if ok || val != 0.0 {
		t.Errorf("Expected default value 0.0 and ok=false, got %f and ok=%v", val, ok)
	}
}

func TestMapOfMaps_GetMapAndSetMap(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()

	// Test getting non-existent map
	innerMap := m.GetMap("nonexistent")
	if innerMap != nil {
		t.Error("Expected nil for non-existent map")
	}

	// Test setting and getting map
	testMap := map[int]float64{1: 1.1, 2: 2.2}
	m.SetMap("test", testMap)

	retrieved := m.GetMap("test")
	if len(retrieved) != 2 {
		t.Errorf("Expected map length 2, got %d", len(retrieved))
	}

	if retrieved[1] != 1.1 || retrieved[2] != 2.2 {
		t.Error("Retrieved map values don't match")
	}

	// Verify it's a copy (modifying original shouldn't affect stored)
	testMap[3] = 3.3
	retrieved2 := m.GetMap("test")
	if len(retrieved2) != 2 {
		t.Error("Expected stored map to be unaffected by original modification")
	}
}

func TestMapOfMaps_LookupMap(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	testMap := map[int]float64{1: 1.1, 2: 2.2}
	m.SetMap("test", testMap)

	// Test existing map
	retrieved, ok := m.LookupMap("test")
	if !ok || len(retrieved) != 2 {
		t.Errorf("Expected map with length 2 and ok=true, got length %d and ok=%v", len(retrieved), ok)
	}

	// Test non-existent map
	retrieved, ok = m.LookupMap("nonexistent")
	if ok || retrieved != nil {
		t.Errorf("Expected nil and ok=false, got %v and ok=%v", retrieved, ok)
	}
}

func TestMapOfMaps_HasAndHasMap(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)

	// Test Has
	if !m.Has("users", 1) {
		t.Error("Expected Has to return true for existing nested key")
	}

	if m.Has("users", 999) {
		t.Error("Expected Has to return false for non-existent inner key")
	}

	if m.Has("nonexistent", 1) {
		t.Error("Expected Has to return false for non-existent outer key")
	}

	// Test HasMap
	if !m.HasMap("users") {
		t.Error("Expected HasMap to return true for existing outer key")
	}

	if m.HasMap("nonexistent") {
		t.Error("Expected HasMap to return false for non-existent outer key")
	}
}

func TestMapOfMaps_PopAndPopMap(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)
	m.Set("users", 2, 20.7)
	m.Set("products", 100, 99.99)

	// Test Pop
	val := m.Pop("users", 1)
	if val != 10.5 {
		t.Errorf("Expected popped value 10.5, got %f", val)
	}

	if m.Has("users", 1) {
		t.Error("Expected key to be removed after pop")
	}

	if !m.Has("users", 2) {
		t.Error("Expected other keys in same map to remain")
	}

	// Test popping non-existent value
	val = m.Pop("nonexistent", 1)
	if val != 0.0 {
		t.Errorf("Expected default value 0.0, got %f", val)
	}

	// Test PopMap
	poppedMap := m.PopMap("users")
	if len(poppedMap) != 1 || poppedMap[2] != 20.7 {
		t.Error("PopMap didn't return correct map")
	}

	if m.HasMap("users") {
		t.Error("Expected outer key to be removed after PopMap")
	}

	if !m.HasMap("products") {
		t.Error("Expected other outer keys to remain")
	}
}

func TestMapOfMaps_SetIfNotPresent(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)

	// Test setting when key exists
	val := m.SetIfNotPresent("users", 1, 99.9)
	if val != 10.5 {
		t.Errorf("Expected existing value 10.5, got %f", val)
	}

	if m.Get("users", 1) != 10.5 {
		t.Error("Expected existing value to be unchanged")
	}

	// Test setting when key doesn't exist
	val = m.SetIfNotPresent("users", 2, 20.7)
	if val != 20.7 {
		t.Errorf("Expected new value 20.7, got %f", val)
	}

	if m.Get("users", 2) != 20.7 {
		t.Error("Expected new value to be set")
	}
}

func TestMapOfMaps_Swap(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)

	// Test swapping existing value
	old := m.Swap("users", 1, 99.9)
	if old != 10.5 {
		t.Errorf("Expected old value 10.5, got %f", old)
	}

	if m.Get("users", 1) != 99.9 {
		t.Errorf("Expected new value 99.9, got %f", m.Get("users", 1))
	}

	// Test swapping non-existent value
	old = m.Swap("users", 2, 20.7)
	if old != 0.0 {
		t.Errorf("Expected default old value 0.0, got %f", old)
	}

	if m.Get("users", 2) != 20.7 {
		t.Errorf("Expected new value 20.7, got %f", m.Get("users", 2))
	}
}

func TestMapOfMaps_DeleteAndDeleteMap(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)
	m.Set("users", 2, 20.7)
	m.Set("users", 3, 30.9)
	m.Set("products", 100, 99.99)

	// Test Delete single key
	deleted := m.Delete("users", 1)
	if !deleted {
		t.Error("Expected Delete to return true")
	}

	if m.Has("users", 1) {
		t.Error("Expected key to be deleted")
	}

	// Test Delete multiple keys
	deleted = m.Delete("users", 2, 3)
	if !deleted {
		t.Error("Expected Delete to return true for multiple keys")
	}

	if m.HasMap("users") {
		t.Error("Expected outer key to be removed when inner map becomes empty")
	}

	// Test Delete non-existent key
	deleted = m.Delete("nonexistent", 1)
	if deleted {
		t.Error("Expected Delete to return false for non-existent key")
	}

	// Test DeleteMap
	deleted = m.DeleteMap("products")
	if !deleted {
		t.Error("Expected DeleteMap to return true")
	}

	if m.HasMap("products") {
		t.Error("Expected outer key to be deleted")
	}
}

func TestMapOfMaps_LenAndOuterLen(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()

	if m.Len() != 0 || m.OuterLen() != 0 {
		t.Error("Expected empty map to have zero lengths")
	}

	m.Set("users", 1, 10.5)
	m.Set("users", 2, 20.7)
	m.Set("products", 100, 99.99)

	if m.Len() != 3 {
		t.Errorf("Expected total length 3, got %d", m.Len())
	}

	if m.OuterLen() != 2 {
		t.Errorf("Expected outer length 2, got %d", m.OuterLen())
	}
}

func TestMapOfMaps_KeysAndValues(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)
	m.Set("users", 2, 20.7)
	m.Set("products", 100, 99.99)

	// Test OuterKeys
	outerKeys := m.OuterKeys()
	if len(outerKeys) != 2 {
		t.Errorf("Expected 2 outer keys, got %d", len(outerKeys))
	}

	expectedOuter := map[string]bool{"users": true, "products": true}
	for _, key := range outerKeys {
		if !expectedOuter[key] {
			t.Errorf("Unexpected outer key: %s", key)
		}
	}

	// Test AllKeys
	allKeys := m.AllKeys()
	if len(allKeys) != 3 {
		t.Errorf("Expected 3 inner keys, got %d", len(allKeys))
	}

	expectedInner := map[int]bool{1: true, 2: true, 100: true}
	for _, key := range allKeys {
		if !expectedInner[key] {
			t.Errorf("Unexpected inner key: %d", key)
		}
	}

	// Test AllValues
	allValues := m.AllValues()
	if len(allValues) != 3 {
		t.Errorf("Expected 3 values, got %d", len(allValues))
	}

	expectedValues := map[float64]bool{10.5: true, 20.7: true, 99.99: true}
	for _, val := range allValues {
		if !expectedValues[val] {
			t.Errorf("Unexpected value: %f", val)
		}
	}
}

func TestMapOfMaps_Change(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)

	m.Change("users", 1, func(outerKey string, innerKey int, value float64) float64 {
		return value * 2
	})

	if val := m.Get("users", 1); val != 21.0 {
		t.Errorf("Expected changed value 21.0, got %f", val)
	}

	// Test changing non-existent key
	m.Change("users", 2, func(outerKey string, innerKey int, value float64) float64 {
		return value + 100
	})

	if val := m.Get("users", 2); val != 100.0 {
		t.Errorf("Expected new value 100.0, got %f", val)
	}
}

func TestMapOfMaps_Transform(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)
	m.Set("users", 2, 20.7)
	m.Set("products", 100, 99.99)

	m.Transform(func(outerKey string, innerKey int, value float64) float64 {
		if outerKey == "users" {
			return value * 2
		}
		return value
	})

	if val := m.Get("users", 1); val != 21.0 {
		t.Errorf("Expected transformed value 21.0, got %f", val)
	}

	if val := m.Get("users", 2); val != 41.4 {
		t.Errorf("Expected transformed value 41.4, got %f", val)
	}

	if val := m.Get("products", 100); val != 99.99 {
		t.Errorf("Expected unchanged value 99.99, got %f", val)
	}
}

func TestMapOfMaps_Range(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)
	m.Set("users", 2, 20.7)
	m.Set("products", 100, 99.99)

	visited := make(map[string]map[int]float64)
	result := m.Range(func(outerKey string, innerKey int, value float64) bool {
		if visited[outerKey] == nil {
			visited[outerKey] = make(map[int]float64)
		}
		visited[outerKey][innerKey] = value
		return value < 50.0 // Stop when we hit a value >= 50
	})

	if result {
		t.Error("Expected Range to return false when stopped early")
	}

	if len(visited) == 0 {
		t.Error("Expected some values to be visited")
	}
}

func TestMapOfMaps_CopyAndRaw(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)
	m.Set("products", 100, 99.99)

	// Test Copy
	copied := m.Copy()
	if len(copied) != 2 {
		t.Errorf("Expected copied map to have 2 outer keys, got %d", len(copied))
	}

	// Modify copy and ensure original is unchanged
	copied["users"][1] = 999.9
	if m.Get("users", 1) != 10.5 {
		t.Error("Expected original to be unchanged after modifying copy")
	}

	// Test Raw
	raw := m.Raw()
	if len(raw) != 2 {
		t.Errorf("Expected raw map to have 2 outer keys, got %d", len(raw))
	}

	// Modifying raw affects original
	raw["users"][1] = 888.8
	if m.Get("users", 1) != 888.8 {
		t.Error("Expected original to be affected when modifying raw")
	}
}

func TestMapOfMaps_ClearAndRefill(t *testing.T) {
	m := abstract.NewMapOfMaps[string, int, float64]()
	m.Set("users", 1, 10.5)
	m.Set("products", 100, 99.99)

	// Test Clear
	m.Clear()
	if !m.IsEmpty() {
		t.Error("Expected map to be empty after Clear")
	}

	// Test Refill
	newData := map[string]map[int]float64{
		"categories": {1: 1.1, 2: 2.2},
		"items":      {10: 10.1, 20: 20.2},
	}

	m.Refill(newData)
	if m.OuterLen() != 2 {
		t.Errorf("Expected 2 outer keys after refill, got %d", m.OuterLen())
	}

	if m.Get("categories", 1) != 1.1 {
		t.Error("Expected refilled data to be accessible")
	}
}

// Tests for SafeMapOfMaps[K1, K2, V]

func TestSafeMapOfMaps_BasicOperations(t *testing.T) {
	m := abstract.NewSafeMapOfMaps[string, int, float64]()

	// Test concurrent Set operations
	var wg sync.WaitGroup
	numGoroutines := 50

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			outerKey := "group" + strconv.Itoa(i%5)
			innerKey := i
			value := float64(i) * 1.1
			m.Set(outerKey, innerKey, value)
		}(i)
	}

	wg.Wait()

	if m.Len() != numGoroutines {
		t.Errorf("Expected %d total items, got %d", numGoroutines, m.Len())
	}

	if m.OuterLen() != 5 {
		t.Errorf("Expected 5 outer keys, got %d", m.OuterLen())
	}
}

func TestSafeMapOfMaps_ConcurrentReadWrite(t *testing.T) {
	m := abstract.NewSafeMapOfMaps[string, int, float64]()

	// Pre-populate
	for i := 0; i < 10; i++ {
		m.Set("test", i, float64(i)*1.1)
	}

	var wg sync.WaitGroup
	numReaders := 10
	numWriters := 5

	// Start readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				val := m.Get("test", j%10)
				if val < 0 {
					t.Errorf("Unexpected negative value: %f", val)
				}

				if ok := m.Has("test", j%10); !ok {
					t.Error("Expected key to exist")
				}
			}
		}(i)
	}

	// Start writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				key := (i*50 + j) % 10
				m.Set("test", key, float64(i*50+j)*1.1)
			}
		}(i)
	}

	wg.Wait()
}

func TestSafeMapOfMaps_AllMethods(t *testing.T) {
	m := abstract.NewSafeMapOfMapsWithSize[string, int, float64](10)

	// Test all methods work the same as non-safe version
	m.Set("users", 1, 10.5)
	m.Set("users", 2, 20.7)

	if val := m.Get("users", 1); val != 10.5 {
		t.Errorf("Expected 10.5, got %f", val)
	}

	if val, ok := m.Lookup("users", 1); !ok || val != 10.5 {
		t.Errorf("Expected 10.5 and true, got %f and %v", val, ok)
	}

	if !m.Has("users", 1) {
		t.Error("Expected Has to return true")
	}

	if !m.HasMap("users") {
		t.Error("Expected HasMap to return true")
	}

	userMap := m.GetMap("users")
	if len(userMap) != 2 {
		t.Errorf("Expected user map length 2, got %d", len(userMap))
	}

	if userMap, ok := m.LookupMap("users"); !ok || len(userMap) != 2 {
		t.Errorf("Expected user map with length 2, got length %d and ok=%v", len(userMap), ok)
	}

	val := m.Pop("users", 1)
	if val != 10.5 {
		t.Errorf("Expected popped value 10.5, got %f", val)
	}

	testMap := map[int]float64{10: 10.1, 20: 20.2}
	m.SetMap("products", testMap)

	poppedMap := m.PopMap("products")
	if len(poppedMap) != 2 {
		t.Errorf("Expected popped map length 2, got %d", len(poppedMap))
	}

	old := m.SetIfNotPresent("users", 2, 99.9)
	if old != 20.7 {
		t.Errorf("Expected existing value 20.7, got %f", old)
	}

	old = m.Swap("users", 2, 30.9)
	if old != 20.7 {
		t.Errorf("Expected old value 20.7, got %f", old)
	}

	if deleted := m.Delete("users", 2); !deleted {
		t.Error("Expected Delete to return true")
	}

	m.Set("test", 1, 1.1)
	m.Set("test", 2, 2.2)

	if deleted := m.DeleteMap("test"); !deleted {
		t.Error("Expected DeleteMap to return true")
	}

	// Test utility methods
	m.Set("a", 1, 1.1)
	m.Set("a", 2, 2.2)
	m.Set("b", 3, 3.3)

	if m.Len() != 3 {
		t.Errorf("Expected length 3, got %d", m.Len())
	}

	if m.OuterLen() != 2 {
		t.Errorf("Expected outer length 2, got %d", m.OuterLen())
	}

	if m.IsEmpty() {
		t.Error("Expected map not to be empty")
	}

	outerKeys := m.OuterKeys()
	if len(outerKeys) != 2 {
		t.Errorf("Expected 2 outer keys, got %d", len(outerKeys))
	}

	allKeys := m.AllKeys()
	if len(allKeys) != 3 {
		t.Errorf("Expected 3 inner keys, got %d", len(allKeys))
	}

	allValues := m.AllValues()
	if len(allValues) != 3 {
		t.Errorf("Expected 3 values, got %d", len(allValues))
	}

	// Test Change
	m.Change("a", 1, func(outer string, inner int, val float64) float64 {
		return val * 2
	})

	if val := m.Get("a", 1); val != 2.2 {
		t.Errorf("Expected changed value 2.2, got %f", val)
	}

	// Test Transform
	m.Transform(func(outer string, inner int, val float64) float64 {
		return val + 1.0
	})

	if val := m.Get("a", 1); val != 3.2 {
		t.Errorf("Expected transformed value 3.2, got %f", val)
	}

	// Test Range
	count := 0
	result := m.Range(func(outer string, inner int, val float64) bool {
		count++
		return count < 2 // Stop after 2 iterations
	})

	if result || count != 2 {
		t.Errorf("Expected Range to stop after 2 iterations, got %d and result %v", count, result)
	}

	// Test Copy
	copied := m.Copy()
	if len(copied) != 2 {
		t.Errorf("Expected copied map to have 2 outer keys, got %d", len(copied))
	}

	// Test Raw
	raw := m.Raw()
	if len(raw) != 2 {
		t.Errorf("Expected raw map to have 2 outer keys, got %d", len(raw))
	}

	// Test Clear
	m.Clear()
	if !m.IsEmpty() {
		t.Error("Expected map to be empty after Clear")
	}

	// Test Refill
	refillData := map[string]map[int]float64{
		"new1": {1: 1.1},
		"new2": {2: 2.2},
	}
	m.Refill(refillData)

	if m.OuterLen() != 2 {
		t.Errorf("Expected 2 outer keys after refill, got %d", m.OuterLen())
	}
}

func TestMapOfMaps_DifferentTypes(t *testing.T) {
	// Test with different type combinations to verify the three type parameters work correctly

	// String -> Int -> String
	m1 := abstract.NewMapOfMaps[string, int, string]()
	m1.Set("group1", 1, "value1")
	m1.Set("group1", 2, "value2")

	if val := m1.Get("group1", 1); val != "value1" {
		t.Errorf("Expected 'value1', got '%s'", val)
	}

	// Int -> String -> Bool
	m2 := abstract.NewMapOfMaps[int, string, bool]()
	m2.Set(100, "key1", true)
	m2.Set(100, "key2", false)
	m2.Set(200, "key1", true)

	if val := m2.Get(100, "key2"); val != false {
		t.Errorf("Expected false, got %v", val)
	}

	if m2.OuterLen() != 2 {
		t.Errorf("Expected 2 outer keys, got %d", m2.OuterLen())
	}

	// Test that outer and inner key types are truly independent
	outerKeys := m2.OuterKeys() // Should be []int
	innerKeys := m2.AllKeys()   // Should be []string

	if len(outerKeys) != 2 {
		t.Errorf("Expected 2 outer keys (int), got %d", len(outerKeys))
	}

	if len(innerKeys) != 3 {
		t.Errorf("Expected 3 inner keys (string), got %d", len(innerKeys))
	}
}

// Tests for nil values in all map types

func TestMap_NilValues(t *testing.T) {
	// Test with *int as value type to allow nil
	var m abstract.Map[string, *int]

	// Test setting nil value
	m.Set("nilkey", nil)

	// Test Get with nil value
	val := m.Get("nilkey")
	if val != nil {
		t.Errorf("Expected nil value, got %v", val)
	}

	// Test Lookup with nil value
	val, ok := m.Lookup("nilkey")
	if !ok || val != nil {
		t.Errorf("Expected nil value and true, got %v and %v", val, ok)
	}

	// Test Has with nil value
	if !m.Has("nilkey") {
		t.Error("Expected key with nil value to exist")
	}

	// Test Pop with nil value
	val = m.Pop("nilkey")
	if val != nil {
		t.Errorf("Expected popped nil value, got %v", val)
	}

	// Test that key is removed after pop
	if m.Has("nilkey") {
		t.Error("Expected key to be removed after pop")
	}

	// Test SetIfNotPresent with nil
	m.Set("key1", nil)
	returnedVal := m.SetIfNotPresent("key1", &[]int{42}[0])
	if returnedVal != nil {
		t.Errorf("Expected existing nil value, got %v", returnedVal)
	}

	// Test Swap with nil
	newVal := &[]int{100}[0]
	oldVal := m.Swap("key1", newVal)
	if oldVal != nil {
		t.Errorf("Expected old nil value, got %v", oldVal)
	}

	// Test Change with nil
	m.Set("key2", nil)
	m.Change("key2", func(k string, v *int) *int {
		if v == nil {
			return &[]int{42}[0]
		}
		return v
	})
	if val := m.Get("key2"); val == nil || *val != 42 {
		t.Errorf("Expected changed value 42, got %v", val)
	}

	// Test Transform with nil values
	m.Set("key3", nil)
	m.Set("key4", &[]int{10}[0])
	m.Transform(func(k string, v *int) *int {
		if v == nil {
			return &[]int{0}[0]
		}
		return v
	})

	if val := m.Get("key3"); val == nil || *val != 0 {
		t.Errorf("Expected transformed nil to 0, got %v", val)
	}

	// Test Range with nil values
	m.Set("key5", nil)
	m.Range(func(k string, v *int) bool {
		if k == "key5" && v != nil {
			t.Errorf("Expected nil value for key5, got %v", v)
		}
		return true
	})

	// Test Values with nil
	values := m.Values()
	hasNil := false
	for _, v := range values {
		if v == nil {
			hasNil = true
			break
		}
	}
	if !hasNil {
		t.Error("Expected at least one nil value in Values slice")
	}

	// Test Copy with nil values
	copied := m.Copy()
	if val, ok := copied["key5"]; !ok || val != nil {
		t.Errorf("Expected copied map to have nil value for key5, got %v", val)
	}

	// Test iterators with nil values
	for k, v := range m.Iter() {
		if k == "key5" && v != nil {
			t.Errorf("Expected nil value for key5 in iterator, got %v", v)
		}
	}

	// Test IterKeys with nil values
	for k := range m.IterKeys() {
		if k == "key5" {
			// Key should exist even if value is nil
			if !m.Has(k) {
				t.Error("Expected key5 to exist in IterKeys")
			}
		}
	}

	// Test IterValues with nil values
	hasNilInValues := false
	for v := range m.IterValues() {
		if v == nil {
			hasNilInValues = true
			break
		}
	}
	if !hasNilInValues {
		t.Error("Expected at least one nil value in IterValues")
	}
}

func TestSafeMap_NilValues(t *testing.T) {
	// Test with *int as value type to allow nil
	var m abstract.SafeMap[string, *int]

	// Test setting nil value
	m.Set("nilkey", nil)

	// Test Get with nil value
	val := m.Get("nilkey")
	if val != nil {
		t.Errorf("Expected nil value, got %v", val)
	}

	// Test Lookup with nil value
	val, ok := m.Lookup("nilkey")
	if !ok || val != nil {
		t.Errorf("Expected nil value and true, got %v and %v", val, ok)
	}

	// Test Has with nil value
	if !m.Has("nilkey") {
		t.Error("Expected key with nil value to exist")
	}

	// Test Pop with nil value
	val = m.Pop("nilkey")
	if val != nil {
		t.Errorf("Expected popped nil value, got %v", val)
	}

	// Test SetIfNotPresent with nil
	m.Set("key1", nil)
	returnedVal := m.SetIfNotPresent("key1", &[]int{42}[0])
	if returnedVal != nil {
		t.Errorf("Expected existing nil value, got %v", returnedVal)
	}

	// Test Swap with nil
	newVal := &[]int{100}[0]
	oldVal := m.Swap("key1", newVal)
	if oldVal != nil {
		t.Errorf("Expected old nil value, got %v", oldVal)
	}

	// Test Change with nil
	m.Set("key2", nil)
	m.Change("key2", func(k string, v *int) *int {
		if v == nil {
			return &[]int{42}[0]
		}
		return v
	})
	if val := m.Get("key2"); val == nil || *val != 42 {
		t.Errorf("Expected changed value 42, got %v", val)
	}

	// Test Transform with nil values
	m.Set("key3", nil)
	m.Set("key4", &[]int{10}[0])
	m.Transform(func(k string, v *int) *int {
		if v == nil {
			return &[]int{0}[0]
		}
		return v
	})

	if val := m.Get("key3"); val == nil || *val != 0 {
		t.Errorf("Expected transformed nil to 0, got %v", val)
	}

	// Test Range with nil values
	m.Set("key5", nil)
	m.Range(func(k string, v *int) bool {
		if k == "key5" && v != nil {
			t.Errorf("Expected nil value for key5, got %v", v)
		}
		return true
	})

	// Test concurrent access with nil values
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "concurrent" + strconv.Itoa(i)
			if i%2 == 0 {
				m.Set(key, nil)
			} else {
				m.Set(key, &[]int{i}[0])
			}

			val := m.Get(key)
			if i%2 == 0 && val != nil {
				t.Errorf("Expected nil value for even index %d, got %v", i, val)
			} else if i%2 == 1 && (val == nil || *val != i) {
				t.Errorf("Expected value %d for odd index %d, got %v", i, i, val)
			}
		}(i)
	}
	wg.Wait()

	// Test Copy with nil values
	copied := m.Copy()
	if val, ok := copied["key5"]; !ok || val != nil {
		t.Errorf("Expected copied map to have nil value for key5, got %v", val)
	}

	// Test iterators with nil values
	for k, v := range m.Iter() {
		if k == "key5" && v != nil {
			t.Errorf("Expected nil value for key5 in iterator, got %v", v)
		}
	}

	// Test IterKeys with nil values
	for k := range m.IterKeys() {
		if k == "key5" {
			// Key should exist even if value is nil
			if !m.Has(k) {
				t.Error("Expected key5 to exist in IterKeys")
			}
		}
	}

	// Test IterValues with nil values
	hasNilInValues := false
	for v := range m.IterValues() {
		if v == nil {
			hasNilInValues = true
			break
		}
	}
	if !hasNilInValues {
		t.Error("Expected at least one nil value in IterValues")
	}
}

func TestEntityMap_NilValues(t *testing.T) {
	// Test with *testEntity as value type to allow nil
	m := abstract.NewEntityMap[int, *testEntity]()

	// Note: EntityMap methods like Set, SetManualOrder, AllOrdered, LookupByName, etc.
	// are not designed to handle nil entities because they call methods on the entities.
	// This is expected behavior. We test basic map operations that can handle nil values.

	// Manually set nil entity to underlying map for testing
	rawMap := m.Raw()
	rawMap[0] = nil

	// Test Get with nil entity
	val := m.Get(0)
	if val != nil {
		t.Errorf("Expected nil entity, got %v", val)
	}

	// Test Has with nil entities
	if !m.Has(0) {
		t.Error("Expected key with nil entity to exist")
	}

	// Test Len with nil entities
	if m.Len() != 1 {
		t.Errorf("Expected length 1, got %d", m.Len())
	}

	// Test basic operations that don't call entity methods
	// Range can handle nil entities as long as we don't use them
	m.Range(func(k int, v *testEntity) bool {
		if k == 0 && v != nil {
			t.Errorf("Expected nil entity for key 0, got %v", v)
		}
		return true
	})
}

func TestSafeEntityMap_NilValues(t *testing.T) {
	// Test with *testEntity as value type to allow nil
	m := abstract.NewSafeEntityMap[int, *testEntity]()

	// Note: EntityMap methods like Set, SetManualOrder, AllOrdered, LookupByName, etc.
	// are not designed to handle nil entities because they call methods on the entities.
	// This is expected behavior. We test basic map operations that can handle nil values.

	// Manually set nil entity to underlying map for testing
	rawMap := m.Raw()
	rawMap[0] = nil

	// Test Get with nil entity
	val := m.Get(0)
	if val != nil {
		t.Errorf("Expected nil entity, got %v", val)
	}

	// Test Has with nil entities
	if !m.Has(0) {
		t.Error("Expected key with nil entity to exist")
	}

	// Test Len with nil entities
	if m.Len() != 1 {
		t.Errorf("Expected length 1, got %d", m.Len())
	}

	// Test concurrent access with valid entities (not nil)
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			entity := &testEntity{id: i, name: "test"}
			m.Set(entity)

			val := m.Get(i)
			if val == nil || val.id != i {
				t.Errorf("Expected entity with id %d, got %v", i, val)
			}
		}(i)
	}
	wg.Wait()
}

func TestOrderedPairs_NilValues(t *testing.T) {
	// Test with *int as value type to allow nil
	var pairs abstract.OrderedPairs[int, *int]

	// Test adding nil values
	pairs.Add(1, nil)
	pairs.Add(2, &[]int{42}[0])
	pairs.Add(3, nil)

	// Test Get with nil values
	val := pairs.Get(1)
	if val != nil {
		t.Errorf("Expected nil value, got %v", val)
	}

	val = pairs.Get(2)
	if val == nil || *val != 42 {
		t.Errorf("Expected value 42, got %v", val)
	}

	// Test Keys with nil values
	keys := pairs.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Test Rand with nil values (should handle gracefully)
	for i := 0; i < 10; i++ {
		randVal := pairs.Rand()
		// Should be either nil or pointer to 42
		if randVal != nil && *randVal != 42 {
			t.Errorf("Expected either nil or 42, got %v", randVal)
		}
	}
}

func TestSafeOrderedPairs_NilValues(t *testing.T) {
	// Test with *int as value type to allow nil
	pairs := abstract.NewSafeOrderedPairs[int, *int]()

	// Test adding nil values
	pairs.Add(1, nil)
	pairs.Add(2, &[]int{42}[0])
	pairs.Add(3, nil)

	// Test Get with nil values
	val := pairs.Get(1)
	if val != nil {
		t.Errorf("Expected nil value, got %v", val)
	}

	val = pairs.Get(2)
	if val == nil || *val != 42 {
		t.Errorf("Expected value 42, got %v", val)
	}

	// Test concurrent access with nil values
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				pairs.Add(i+10, nil)
			} else {
				pairs.Add(i+10, &[]int{i}[0])
			}

			val := pairs.Get(i + 10)
			if i%2 == 0 && val != nil {
				t.Errorf("Expected nil value for even index %d, got %v", i, val)
			} else if i%2 == 1 && (val == nil || *val != i) {
				t.Errorf("Expected value %d for odd index %d, got %v", i, i, val)
			}
		}(i)
	}
	wg.Wait()

	// Test Rand with nil values (should handle gracefully)
	for i := 0; i < 10; i++ {
		randVal := pairs.Rand()
		// Should be either nil or some integer pointer
		if randVal != nil && *randVal < 0 {
			t.Errorf("Expected non-negative value or nil, got %v", randVal)
		}
	}
}

func TestMapOfMaps_NilValues(t *testing.T) {
	// Test with *int as value type to allow nil
	var m abstract.MapOfMaps[string, int, *int]

	// Test setting nil values
	m.Set("group1", 1, nil)
	m.Set("group1", 2, &[]int{42}[0])
	m.Set("group2", 1, nil)

	// Test Get with nil values
	val := m.Get("group1", 1)
	if val != nil {
		t.Errorf("Expected nil value, got %v", val)
	}

	val = m.Get("group1", 2)
	if val == nil || *val != 42 {
		t.Errorf("Expected value 42, got %v", val)
	}

	// Test Lookup with nil values
	val, ok := m.Lookup("group1", 1)
	if !ok || val != nil {
		t.Errorf("Expected nil value and true, got %v and %v", val, ok)
	}

	// Test Has with nil values
	if !m.Has("group1", 1) {
		t.Error("Expected key with nil value to exist")
	}

	// Test Pop with nil values
	val = m.Pop("group1", 1)
	if val != nil {
		t.Errorf("Expected popped nil value, got %v", val)
	}

	// Test SetIfNotPresent with nil
	m.Set("group3", 1, nil)
	returnedVal := m.SetIfNotPresent("group3", 1, &[]int{100}[0])
	if returnedVal != nil {
		t.Errorf("Expected existing nil value, got %v", returnedVal)
	}

	// Test Swap with nil
	newVal := &[]int{200}[0]
	oldVal := m.Swap("group3", 1, newVal)
	if oldVal != nil {
		t.Errorf("Expected old nil value, got %v", oldVal)
	}

	// Test Change with nil
	m.Set("group4", 1, nil)
	m.Change("group4", 1, func(outer string, inner int, v *int) *int {
		if v == nil {
			return &[]int{42}[0]
		}
		return v
	})
	if val := m.Get("group4", 1); val == nil || *val != 42 {
		t.Errorf("Expected changed value 42, got %v", val)
	}

	// Test Transform with nil values
	m.Set("group5", 1, nil)
	m.Set("group5", 2, &[]int{10}[0])
	m.Transform(func(outer string, inner int, v *int) *int {
		if v == nil {
			return &[]int{0}[0]
		}
		return v
	})

	if val := m.Get("group5", 1); val == nil || *val != 0 {
		t.Errorf("Expected transformed nil to 0, got %v", val)
	}

	// Test Range with nil values
	m.Set("group6", 1, nil)
	m.Range(func(outer string, inner int, v *int) bool {
		if outer == "group6" && inner == 1 && v != nil {
			t.Errorf("Expected nil value for group6[1], got %v", v)
		}
		return true
	})

	// Test AllValues with nil
	allValues := m.AllValues()
	hasNil := false
	for _, v := range allValues {
		if v == nil {
			hasNil = true
			break
		}
	}
	if !hasNil {
		t.Error("Expected at least one nil value in AllValues")
	}

	// Test Copy with nil values
	copied := m.Copy()
	if val := copied["group6"][1]; val != nil {
		t.Errorf("Expected copied map to have nil value for group6[1], got %v", val)
	}

	// Test SetMap with nil values
	nilMap := map[int]*int{
		10: nil,
		20: &[]int{30}[0],
	}
	m.SetMap("nilGroup", nilMap)

	if val := m.Get("nilGroup", 10); val != nil {
		t.Errorf("Expected nil value for nilGroup[10], got %v", val)
	}

	// Test GetMap with nil values
	retrievedMap := m.GetMap("nilGroup")
	if retrievedMap[10] != nil {
		t.Errorf("Expected nil value in retrieved map, got %v", retrievedMap[10])
	}
}

func TestSafeMapOfMaps_NilValues(t *testing.T) {
	// Test with *int as value type to allow nil
	var m abstract.SafeMapOfMaps[string, int, *int]

	// Test setting nil values
	m.Set("group1", 1, nil)
	m.Set("group1", 2, &[]int{42}[0])
	m.Set("group2", 1, nil)

	// Test Get with nil values
	val := m.Get("group1", 1)
	if val != nil {
		t.Errorf("Expected nil value, got %v", val)
	}

	val = m.Get("group1", 2)
	if val == nil || *val != 42 {
		t.Errorf("Expected value 42, got %v", val)
	}

	// Test concurrent access with nil values
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			outerKey := "concurrent" + strconv.Itoa(i%5)
			innerKey := i

			if i%2 == 0 {
				m.Set(outerKey, innerKey, nil)
			} else {
				m.Set(outerKey, innerKey, &[]int{i}[0])
			}

			val := m.Get(outerKey, innerKey)
			if i%2 == 0 && val != nil {
				t.Errorf("Expected nil value for even index %d, got %v", i, val)
			} else if i%2 == 1 && (val == nil || *val != i) {
				t.Errorf("Expected value %d for odd index %d, got %v", i, i, val)
			}
		}(i)
	}
	wg.Wait()

	// Test Lookup with nil values
	val, ok := m.Lookup("group1", 1)
	if !ok || val != nil {
		t.Errorf("Expected nil value and true, got %v and %v", val, ok)
	}

	// Test Transform with nil values
	m.Set("group3", 1, nil)
	m.Set("group3", 2, &[]int{10}[0])
	m.Transform(func(outer string, inner int, v *int) *int {
		if v == nil {
			return &[]int{0}[0]
		}
		return v
	})

	if val := m.Get("group3", 1); val == nil || *val != 0 {
		t.Errorf("Expected transformed nil to 0, got %v", val)
	}

	// Test Range with nil values
	m.Set("group4", 1, nil)
	m.Range(func(outer string, inner int, v *int) bool {
		if outer == "group4" && inner == 1 && v != nil {
			t.Errorf("Expected nil value for group4[1], got %v", v)
		}
		return true
	})

	// Test SetMap with nil values
	nilMap := map[int]*int{
		10: nil,
		20: &[]int{30}[0],
	}
	m.SetMap("nilGroup", nilMap)

	if val := m.Get("nilGroup", 10); val != nil {
		t.Errorf("Expected nil value for nilGroup[10], got %v", val)
	}

	// Test GetMap with nil values
	retrievedMap := m.GetMap("nilGroup")
	if retrievedMap[10] != nil {
		t.Errorf("Expected nil value in retrieved map, got %v", retrievedMap[10])
	}
}

func TestMap_NilMapInitialization(t *testing.T) {
	// Test that methods handle nil map initialization properly
	var m *abstract.Map[string, int]

	// These should not panic due to nil map initialization in methods
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic, but got: %v", r)
		}
	}()

	m = &abstract.Map[string, int]{}

	// Test that methods initialize the map properly
	m.Set("key1", 42)
	if val := m.Get("key1"); val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}

	if !m.Has("key1") {
		t.Error("Expected key to exist after Set")
	}

	if m.Len() != 1 {
		t.Errorf("Expected length 1, got %d", m.Len())
	}
}

func TestSafeMap_NilMapInitialization(t *testing.T) {
	// Test that methods handle nil map initialization properly
	var m *abstract.SafeMap[string, int]

	// These should not panic due to nil map initialization in methods
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic, but got: %v", r)
		}
	}()

	m = &abstract.SafeMap[string, int]{}

	// Test that methods initialize the map properly
	m.Set("key1", 42)
	if val := m.Get("key1"); val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}

	if !m.Has("key1") {
		t.Error("Expected key to exist after Set")
	}

	if m.Len() != 1 {
		t.Errorf("Expected length 1, got %d", m.Len())
	}
}

func TestMapOfMaps_NilMapInitialization(t *testing.T) {
	// Test that methods handle nil map initialization properly
	var m *abstract.MapOfMaps[string, int, float64]

	// These should not panic due to nil map initialization in methods
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic, but got: %v", r)
		}
	}()

	m = &abstract.MapOfMaps[string, int, float64]{}

	// Test that methods initialize the map properly
	m.Set("outer1", 1, 1.1)
	if val := m.Get("outer1", 1); val != 1.1 {
		t.Errorf("Expected 1.1, got %f", val)
	}

	if !m.Has("outer1", 1) {
		t.Error("Expected key to exist after Set")
	}

	if m.Len() != 1 {
		t.Errorf("Expected length 1, got %d", m.Len())
	}
}

func TestSafeMapOfMaps_NilMapInitialization(t *testing.T) {
	// Test that methods handle nil map initialization properly
	var m *abstract.SafeMapOfMaps[string, int, float64]

	// These should not panic due to nil map initialization in methods
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic, but got: %v", r)
		}
	}()

	m = &abstract.SafeMapOfMaps[string, int, float64]{}

	// Test that methods initialize the map properly
	m.Set("outer1", 1, 1.1)
	if val := m.Get("outer1", 1); val != 1.1 {
		t.Errorf("Expected 1.1, got %f", val)
	}

	if !m.Has("outer1", 1) {
		t.Error("Expected key to exist after Set")
	}

	if m.Len() != 1 {
		t.Errorf("Expected length 1, got %d", m.Len())
	}
}

// Tests for uninitialized maps (zero-value with nil items)

func TestMap_UninitializedMethods(t *testing.T) {
	// Test Get with uninitialized map
	var m1 abstract.Map[string, int]
	val := m1.Get("key")
	if val != 0 {
		t.Errorf("Expected default value 0, got %d", val)
	}

	// Test Lookup with uninitialized map
	var m2 abstract.Map[string, int]
	val, ok := m2.Lookup("key")
	if ok || val != 0 {
		t.Errorf("Expected default value 0 and false, got %d and %v", val, ok)
	}

	// Test Has with uninitialized map
	var m3 abstract.Map[string, int]
	if m3.Has("key") {
		t.Error("Expected false for uninitialized map")
	}

	// Test Pop with uninitialized map
	var m4 abstract.Map[string, int]
	val = m4.Pop("key")
	if val != 0 {
		t.Errorf("Expected default value 0, got %d", val)
	}

	// Test Set with uninitialized map
	var m5 abstract.Map[string, int]
	m5.Set("key", 42)
	if m5.Get("key") != 42 {
		t.Errorf("Expected 42 after Set on uninitialized map, got %d", m5.Get("key"))
	}

	// Test SetIfNotPresent with uninitialized map
	var m6 abstract.Map[string, int]
	val = m6.SetIfNotPresent("key", 42)
	if val != 42 {
		t.Errorf("Expected 42 from SetIfNotPresent on uninitialized map, got %d", val)
	}

	// Test Swap with uninitialized map
	var m7 abstract.Map[string, int]
	old := m7.Swap("key", 42)
	if old != 0 {
		t.Errorf("Expected default value 0 from Swap on uninitialized map, got %d", old)
	}

	// Test Delete with uninitialized map
	var m8 abstract.Map[string, int]
	deleted := m8.Delete("key")
	if deleted {
		t.Error("Expected false from Delete on uninitialized map")
	}

	// Test Len with uninitialized map
	var m9 abstract.Map[string, int]
	if m9.Len() != 0 {
		t.Errorf("Expected length 0 for uninitialized map, got %d", m9.Len())
	}

	// Test IsEmpty with uninitialized map
	var m10 abstract.Map[string, int]
	if !m10.IsEmpty() {
		t.Error("Expected true from IsEmpty on uninitialized map")
	}

	// Test Keys with uninitialized map
	var m11 abstract.Map[string, int]
	keys := m11.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected empty keys slice, got length %d", len(keys))
	}

	// Test Values with uninitialized map
	var m12 abstract.Map[string, int]
	values := m12.Values()
	if len(values) != 0 {
		t.Errorf("Expected empty values slice, got length %d", len(values))
	}

	// Test Change with uninitialized map
	var m13 abstract.Map[string, int]
	m13.Change("key", func(k string, v int) int { return v + 1 })
	if m13.Get("key") != 1 {
		t.Errorf("Expected 1 from Change on uninitialized map, got %d", m13.Get("key"))
	}

	// Test Transform with uninitialized map
	var m14 abstract.Map[string, int]
	m14.Transform(func(k string, v int) int { return v + 1 })
	if m14.Len() != 0 {
		t.Errorf("Expected no items after Transform on uninitialized map, got %d", m14.Len())
	}

	// Test Range with uninitialized map
	var m15 abstract.Map[string, int]
	called := false
	result := m15.Range(func(k string, v int) bool {
		called = true
		return true
	})
	if !result || called {
		t.Error("Expected Range to return true without calling function on uninitialized map")
	}

	// Test Copy with uninitialized map
	var m16 abstract.Map[string, int]
	copied := m16.Copy()
	if len(copied) != 0 {
		t.Errorf("Expected empty copied map, got length %d", len(copied))
	}

	// Test Raw with uninitialized map
	var m17 abstract.Map[string, int]
	raw := m17.Raw()
	if len(raw) != 0 {
		t.Errorf("Expected empty raw map, got length %d", len(raw))
	}

	// Test Clear with uninitialized map
	var m18 abstract.Map[string, int]
	m18.Clear()
	if m18.Len() != 0 {
		t.Errorf("Expected length 0 after Clear on uninitialized map, got %d", m18.Len())
	}

	// Test IterKeys with uninitialized map
	var m19 abstract.Map[string, int]
	count := 0
	for range m19.IterKeys() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 iterations from IterKeys on uninitialized map, got %d", count)
	}

	// Test IterValues with uninitialized map
	var m20 abstract.Map[string, int]
	count = 0
	for range m20.IterValues() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 iterations from IterValues on uninitialized map, got %d", count)
	}

	// Test Iter with uninitialized map
	var m21 abstract.Map[string, int]
	count = 0
	for range m21.Iter() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 iterations from Iter on uninitialized map, got %d", count)
	}
}

func TestSafeMap_UninitializedMethods(t *testing.T) {
	// Test Get with uninitialized map
	var m1 abstract.SafeMap[string, int]
	val := m1.Get("key")
	if val != 0 {
		t.Errorf("Expected default value 0, got %d", val)
	}

	// Test Lookup with uninitialized map
	var m2 abstract.SafeMap[string, int]
	val, ok := m2.Lookup("key")
	if ok || val != 0 {
		t.Errorf("Expected default value 0 and false, got %d and %v", val, ok)
	}

	// Test Has with uninitialized map
	var m3 abstract.SafeMap[string, int]
	if m3.Has("key") {
		t.Error("Expected false for uninitialized map")
	}

	// Test Pop with uninitialized map
	var m4 abstract.SafeMap[string, int]
	val = m4.Pop("key")
	if val != 0 {
		t.Errorf("Expected default value 0, got %d", val)
	}

	// Test Set with uninitialized map
	var m5 abstract.SafeMap[string, int]
	m5.Set("key", 42)
	if m5.Get("key") != 42 {
		t.Errorf("Expected 42 after Set on uninitialized map, got %d", m5.Get("key"))
	}

	// Test SetIfNotPresent with uninitialized map
	var m6 abstract.SafeMap[string, int]
	val = m6.SetIfNotPresent("key", 42)
	if val != 42 {
		t.Errorf("Expected 42 from SetIfNotPresent on uninitialized map, got %d", val)
	}

	// Test Swap with uninitialized map
	var m7 abstract.SafeMap[string, int]
	old := m7.Swap("key", 42)
	if old != 0 {
		t.Errorf("Expected default value 0 from Swap on uninitialized map, got %d", old)
	}

	// Test Delete with uninitialized map
	var m8 abstract.SafeMap[string, int]
	deleted := m8.Delete("key")
	if deleted {
		t.Error("Expected false from Delete on uninitialized map")
	}

	// Test Len with uninitialized map
	var m9 abstract.SafeMap[string, int]
	if m9.Len() != 0 {
		t.Errorf("Expected length 0 for uninitialized map, got %d", m9.Len())
	}

	// Test IsEmpty with uninitialized map
	var m10 abstract.SafeMap[string, int]
	if !m10.IsEmpty() {
		t.Error("Expected true from IsEmpty on uninitialized map")
	}

	// Test Keys with uninitialized map
	var m11 abstract.SafeMap[string, int]
	keys := m11.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected empty keys slice, got length %d", len(keys))
	}

	// Test Values with uninitialized map
	var m12 abstract.SafeMap[string, int]
	values := m12.Values()
	if len(values) != 0 {
		t.Errorf("Expected empty values slice, got length %d", len(values))
	}

	// Test Change with uninitialized map
	var m13 abstract.SafeMap[string, int]
	m13.Change("key", func(k string, v int) int { return v + 1 })
	if m13.Get("key") != 1 {
		t.Errorf("Expected 1 from Change on uninitialized map, got %d", m13.Get("key"))
	}

	// Test Transform with uninitialized map
	var m14 abstract.SafeMap[string, int]
	m14.Transform(func(k string, v int) int { return v + 1 })
	if m14.Len() != 0 {
		t.Errorf("Expected no items after Transform on uninitialized map, got %d", m14.Len())
	}

	// Test Range with uninitialized map
	var m15 abstract.SafeMap[string, int]
	called := false
	result := m15.Range(func(k string, v int) bool {
		called = true
		return true
	})
	if !result || called {
		t.Error("Expected Range to return true without calling function on uninitialized map")
	}

	// Test Copy with uninitialized map
	var m16 abstract.SafeMap[string, int]
	copied := m16.Copy()
	if len(copied) != 0 {
		t.Errorf("Expected empty copied map, got length %d", len(copied))
	}

	// Test Clear with uninitialized map
	var m17 abstract.SafeMap[string, int]
	m17.Clear()
	if m17.Len() != 0 {
		t.Errorf("Expected length 0 after Clear on uninitialized map, got %d", m17.Len())
	}

	// Test Refill with uninitialized map
	var m18 abstract.SafeMap[string, int]
	m18.Refill(map[string]int{"key": 42})
	if m18.Get("key") != 42 {
		t.Errorf("Expected 42 after Refill on uninitialized map, got %d", m18.Get("key"))
	}

	// Test Raw with uninitialized map
	var m19 abstract.SafeMap[string, int]
	raw := m19.Raw()
	if len(raw) != 0 {
		t.Errorf("Expected empty raw map, got length %d", len(raw))
	}

	// Test IterKeys with uninitialized map
	var m20 abstract.SafeMap[string, int]
	count := 0
	for range m20.IterKeys() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 iterations from IterKeys on uninitialized map, got %d", count)
	}

	// Test IterValues with uninitialized map
	var m21 abstract.SafeMap[string, int]
	count = 0
	for range m21.IterValues() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 iterations from IterValues on uninitialized map, got %d", count)
	}

	// Test Iter with uninitialized map
	var m22 abstract.SafeMap[string, int]
	count = 0
	for range m22.Iter() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 iterations from Iter on uninitialized map, got %d", count)
	}
}

func TestEntityMap_UninitializedMethods(t *testing.T) {
	// Note: Most EntityMap methods cannot be tested with zero-value because
	// the embedded Map would be nil, causing panics. We test the basic embedded methods
	// that would work if the Map were initialized.

	// These tests will panic as expected because the embedded Map is nil,
	// which demonstrates that EntityMap requires proper initialization
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling methods on zero-value EntityMap")
		}
	}()

	var m abstract.EntityMap[int, *testEntity]
	// This should panic because the embedded Map is nil
	_ = m.Len()
}

func TestSafeEntityMap_UninitializedMethods(t *testing.T) {
	// Note: Most SafeEntityMap methods cannot be tested with zero-value because
	// the embedded SafeMap would be nil, causing panics. We test the basic embedded methods
	// that would work if the SafeMap were initialized.

	// These tests will panic as expected because the embedded SafeMap is nil,
	// which demonstrates that SafeEntityMap requires proper initialization
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling methods on zero-value SafeEntityMap")
		}
	}()

	var m abstract.SafeEntityMap[int, *testEntity]
	// This should panic because the embedded SafeMap is nil
	_ = m.Len()
}

func TestOrderedPairs_UninitializedMethods(t *testing.T) {
	// Test Add with uninitialized OrderedPairs
	var m1 abstract.OrderedPairs[int, string]
	m1.Add(1, "one")
	val := m1.Get(1)
	if val != "one" {
		t.Errorf("Expected 'one', got %s", val)
	}

	// Test Get with uninitialized OrderedPairs
	var m2 abstract.OrderedPairs[int, string]
	val = m2.Get(1)
	if val != "" {
		t.Errorf("Expected empty string, got %s", val)
	}

	// Test Keys with uninitialized OrderedPairs
	var m3 abstract.OrderedPairs[int, string]
	keys := m3.Keys()
	if keys != nil {
		t.Errorf("Expected nil keys slice, got %v", keys)
	}

	// Test Rand with uninitialized OrderedPairs
	var m4 abstract.OrderedPairs[int, string]
	val = m4.Rand()
	if val != "" {
		t.Errorf("Expected empty string from Rand on empty OrderedPairs, got %s", val)
	}

	// Test RandKey with uninitialized OrderedPairs
	var m5 abstract.OrderedPairs[int, string]
	key := m5.RandKey()
	if key != 0 {
		t.Errorf("Expected 0 from RandKey on empty OrderedPairs, got %d", key)
	}
}

func TestSafeOrderedPairs_UninitializedMethods(t *testing.T) {
	// Note: SafeOrderedPairs methods cannot be tested with zero-value because
	// the embedded OrderedPairs would be nil, causing panics.

	// These tests will panic as expected because the embedded OrderedPairs is nil,
	// which demonstrates that SafeOrderedPairs requires proper initialization
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling methods on zero-value SafeOrderedPairs")
		}
	}()

	var m abstract.SafeOrderedPairs[int, string]
	// This should panic because the embedded OrderedPairs is nil
	m.Add(1, "one")
}

func TestMapOfMaps_UninitializedMethods(t *testing.T) {
	// Test Get with uninitialized map
	var m1 abstract.MapOfMaps[string, int, float64]
	val := m1.Get("group", 1)
	if val != 0.0 {
		t.Errorf("Expected default value 0.0, got %f", val)
	}

	// Test GetMap with uninitialized map
	var m2 abstract.MapOfMaps[string, int, float64]
	innerMap := m2.GetMap("group")
	if innerMap != nil {
		t.Errorf("Expected nil inner map, got %v", innerMap)
	}

	// Test Lookup with uninitialized map
	var m3 abstract.MapOfMaps[string, int, float64]
	val, ok := m3.Lookup("group", 1)
	if ok || val != 0.0 {
		t.Errorf("Expected default value 0.0 and false, got %f and %v", val, ok)
	}

	// Test LookupMap with uninitialized map
	var m4 abstract.MapOfMaps[string, int, float64]
	innerMap, ok = m4.LookupMap("group")
	if ok || innerMap != nil {
		t.Errorf("Expected nil inner map and false, got %v and %v", innerMap, ok)
	}

	// Test Has with uninitialized map
	var m5 abstract.MapOfMaps[string, int, float64]
	if m5.Has("group", 1) {
		t.Error("Expected false for uninitialized map")
	}

	// Test HasMap with uninitialized map
	var m6 abstract.MapOfMaps[string, int, float64]
	if m6.HasMap("group") {
		t.Error("Expected false for uninitialized map")
	}

	// Test Pop with uninitialized map
	var m7 abstract.MapOfMaps[string, int, float64]
	val = m7.Pop("group", 1)
	if val != 0.0 {
		t.Errorf("Expected default value 0.0, got %f", val)
	}

	// Test PopMap with uninitialized map
	var m8 abstract.MapOfMaps[string, int, float64]
	innerMap = m8.PopMap("group")
	if innerMap != nil {
		t.Errorf("Expected nil inner map, got %v", innerMap)
	}

	// Test Set with uninitialized map
	var m9 abstract.MapOfMaps[string, int, float64]
	m9.Set("group", 1, 1.1)
	if m9.Get("group", 1) != 1.1 {
		t.Errorf("Expected 1.1 after Set on uninitialized map, got %f", m9.Get("group", 1))
	}

	// Test SetMap with uninitialized map
	var m10 abstract.MapOfMaps[string, int, float64]
	testInnerMap := map[int]float64{1: 1.1, 2: 2.2}
	m10.SetMap("group", testInnerMap)
	if m10.Get("group", 1) != 1.1 {
		t.Errorf("Expected 1.1 after SetMap on uninitialized map, got %f", m10.Get("group", 1))
	}

	// Test SetIfNotPresent with uninitialized map
	var m11 abstract.MapOfMaps[string, int, float64]
	val = m11.SetIfNotPresent("group", 1, 1.1)
	if val != 1.1 {
		t.Errorf("Expected 1.1 from SetIfNotPresent on uninitialized map, got %f", val)
	}

	// Test Swap with uninitialized map
	var m12 abstract.MapOfMaps[string, int, float64]
	old := m12.Swap("group", 1, 1.1)
	if old != 0.0 {
		t.Errorf("Expected default value 0.0 from Swap on uninitialized map, got %f", old)
	}

	// Test Delete with uninitialized map
	var m13 abstract.MapOfMaps[string, int, float64]
	deleted := m13.Delete("group", 1)
	if deleted {
		t.Error("Expected false from Delete on uninitialized map")
	}

	// Test DeleteMap with uninitialized map
	var m14 abstract.MapOfMaps[string, int, float64]
	deleted = m14.DeleteMap("group")
	if deleted {
		t.Error("Expected false from DeleteMap on uninitialized map")
	}

	// Test Len with uninitialized map
	var m15 abstract.MapOfMaps[string, int, float64]
	if m15.Len() != 0 {
		t.Errorf("Expected length 0 for uninitialized map, got %d", m15.Len())
	}

	// Test OuterLen with uninitialized map
	var m16 abstract.MapOfMaps[string, int, float64]
	if m16.OuterLen() != 0 {
		t.Errorf("Expected outer length 0 for uninitialized map, got %d", m16.OuterLen())
	}

	// Test IsEmpty with uninitialized map
	var m17 abstract.MapOfMaps[string, int, float64]
	if !m17.IsEmpty() {
		t.Error("Expected true from IsEmpty on uninitialized map")
	}

	// Test OuterKeys with uninitialized map
	var m18 abstract.MapOfMaps[string, int, float64]
	outerKeys := m18.OuterKeys()
	if len(outerKeys) != 0 {
		t.Errorf("Expected empty outer keys slice, got length %d", len(outerKeys))
	}

	// Test AllKeys with uninitialized map
	var m19 abstract.MapOfMaps[string, int, float64]
	allKeys := m19.AllKeys()
	if len(allKeys) != 0 {
		t.Errorf("Expected empty all keys slice, got length %d", len(allKeys))
	}

	// Test AllValues with uninitialized map
	var m20 abstract.MapOfMaps[string, int, float64]
	allValues := m20.AllValues()
	if len(allValues) != 0 {
		t.Errorf("Expected empty all values slice, got length %d", len(allValues))
	}

	// Test Change with uninitialized map
	var m21 abstract.MapOfMaps[string, int, float64]
	m21.Change("group", 1, func(outer string, inner int, v float64) float64 { return v + 1.0 })
	if m21.Get("group", 1) != 1.0 {
		t.Errorf("Expected 1.0 from Change on uninitialized map, got %f", m21.Get("group", 1))
	}

	// Test Transform with uninitialized map
	var m22 abstract.MapOfMaps[string, int, float64]
	m22.Transform(func(outer string, inner int, v float64) float64 { return v + 1.0 })
	if m22.Len() != 0 {
		t.Errorf("Expected no items after Transform on uninitialized map, got %d", m22.Len())
	}

	// Test Range with uninitialized map
	var m23 abstract.MapOfMaps[string, int, float64]
	called := false
	result := m23.Range(func(outer string, inner int, v float64) bool {
		called = true
		return true
	})
	if !result || called {
		t.Error("Expected Range to return true without calling function on uninitialized map")
	}

	// Test Copy with uninitialized map
	var m24 abstract.MapOfMaps[string, int, float64]
	copied := m24.Copy()
	if len(copied) != 0 {
		t.Errorf("Expected empty copied map, got length %d", len(copied))
	}

	// Test Raw with uninitialized map
	var m25 abstract.MapOfMaps[string, int, float64]
	raw := m25.Raw()
	if len(raw) != 0 {
		t.Errorf("Expected empty raw map, got length %d", len(raw))
	}

	// Test Clear with uninitialized map
	var m26 abstract.MapOfMaps[string, int, float64]
	m26.Clear()
	if m26.Len() != 0 {
		t.Errorf("Expected length 0 after Clear on uninitialized map, got %d", m26.Len())
	}

	// Test Refill with uninitialized map
	var m27 abstract.MapOfMaps[string, int, float64]
	refillData := map[string]map[int]float64{"group": {1: 1.1}}
	m27.Refill(refillData)
	if m27.Get("group", 1) != 1.1 {
		t.Errorf("Expected 1.1 after Refill on uninitialized map, got %f", m27.Get("group", 1))
	}
}

func TestSafeMapOfMaps_UninitializedMethods(t *testing.T) {
	// Test Get with uninitialized map
	var m1 abstract.SafeMapOfMaps[string, int, float64]
	val := m1.Get("group", 1)
	if val != 0.0 {
		t.Errorf("Expected default value 0.0, got %f", val)
	}

	// Test GetMap with uninitialized map
	var m2 abstract.SafeMapOfMaps[string, int, float64]
	innerMap := m2.GetMap("group")
	if innerMap != nil {
		t.Errorf("Expected nil inner map, got %v", innerMap)
	}

	// Test Lookup with uninitialized map
	var m3 abstract.SafeMapOfMaps[string, int, float64]
	val, ok := m3.Lookup("group", 1)
	if ok || val != 0.0 {
		t.Errorf("Expected default value 0.0 and false, got %f and %v", val, ok)
	}

	// Test LookupMap with uninitialized map
	var m4 abstract.SafeMapOfMaps[string, int, float64]
	innerMap, ok = m4.LookupMap("group")
	if ok || innerMap != nil {
		t.Errorf("Expected nil inner map and false, got %v and %v", innerMap, ok)
	}

	// Test Has with uninitialized map
	var m5 abstract.SafeMapOfMaps[string, int, float64]
	if m5.Has("group", 1) {
		t.Error("Expected false for uninitialized map")
	}

	// Test HasMap with uninitialized map
	var m6 abstract.SafeMapOfMaps[string, int, float64]
	if m6.HasMap("group") {
		t.Error("Expected false for uninitialized map")
	}

	// Test Pop with uninitialized map
	var m7 abstract.SafeMapOfMaps[string, int, float64]
	val = m7.Pop("group", 1)
	if val != 0.0 {
		t.Errorf("Expected default value 0.0, got %f", val)
	}

	// Test PopMap with uninitialized map
	var m8 abstract.SafeMapOfMaps[string, int, float64]
	innerMap = m8.PopMap("group")
	if innerMap != nil {
		t.Errorf("Expected nil inner map, got %v", innerMap)
	}

	// Test Set with uninitialized map
	var m9 abstract.SafeMapOfMaps[string, int, float64]
	m9.Set("group", 1, 1.1)
	if m9.Get("group", 1) != 1.1 {
		t.Errorf("Expected 1.1 after Set on uninitialized map, got %f", m9.Get("group", 1))
	}

	// Test SetMap with uninitialized map
	var m10 abstract.SafeMapOfMaps[string, int, float64]
	testInnerMap := map[int]float64{1: 1.1, 2: 2.2}
	m10.SetMap("group", testInnerMap)
	if m10.Get("group", 1) != 1.1 {
		t.Errorf("Expected 1.1 after SetMap on uninitialized map, got %f", m10.Get("group", 1))
	}

	// Test SetIfNotPresent with uninitialized map
	var m11 abstract.SafeMapOfMaps[string, int, float64]
	val = m11.SetIfNotPresent("group", 1, 1.1)
	if val != 1.1 {
		t.Errorf("Expected 1.1 from SetIfNotPresent on uninitialized map, got %f", val)
	}

	// Test Swap with uninitialized map
	var m12 abstract.SafeMapOfMaps[string, int, float64]
	old := m12.Swap("group", 1, 1.1)
	if old != 0.0 {
		t.Errorf("Expected default value 0.0 from Swap on uninitialized map, got %f", old)
	}

	// Test Delete with uninitialized map
	var m13 abstract.SafeMapOfMaps[string, int, float64]
	deleted := m13.Delete("group", 1)
	if deleted {
		t.Error("Expected false from Delete on uninitialized map")
	}

	// Test DeleteMap with uninitialized map
	var m14 abstract.SafeMapOfMaps[string, int, float64]
	deleted = m14.DeleteMap("group")
	if deleted {
		t.Error("Expected false from DeleteMap on uninitialized map")
	}

	// Test Len with uninitialized map
	var m15 abstract.SafeMapOfMaps[string, int, float64]
	if m15.Len() != 0 {
		t.Errorf("Expected length 0 for uninitialized map, got %d", m15.Len())
	}

	// Test OuterLen with uninitialized map
	var m16 abstract.SafeMapOfMaps[string, int, float64]
	if m16.OuterLen() != 0 {
		t.Errorf("Expected outer length 0 for uninitialized map, got %d", m16.OuterLen())
	}

	// Test IsEmpty with uninitialized map
	var m17 abstract.SafeMapOfMaps[string, int, float64]
	if !m17.IsEmpty() {
		t.Error("Expected true from IsEmpty on uninitialized map")
	}

	// Test OuterKeys with uninitialized map
	var m18 abstract.SafeMapOfMaps[string, int, float64]
	outerKeys := m18.OuterKeys()
	if len(outerKeys) != 0 {
		t.Errorf("Expected empty outer keys slice, got length %d", len(outerKeys))
	}

	// Test AllKeys with uninitialized map
	var m19 abstract.SafeMapOfMaps[string, int, float64]
	allKeys := m19.AllKeys()
	if len(allKeys) != 0 {
		t.Errorf("Expected empty all keys slice, got length %d", len(allKeys))
	}

	// Test AllValues with uninitialized map
	var m20 abstract.SafeMapOfMaps[string, int, float64]
	allValues := m20.AllValues()
	if len(allValues) != 0 {
		t.Errorf("Expected empty all values slice, got length %d", len(allValues))
	}

	// Test Change with uninitialized map
	var m21 abstract.SafeMapOfMaps[string, int, float64]
	m21.Change("group", 1, func(outer string, inner int, v float64) float64 { return v + 1.0 })
	if m21.Get("group", 1) != 1.0 {
		t.Errorf("Expected 1.0 from Change on uninitialized map, got %f", m21.Get("group", 1))
	}

	// Test Transform with uninitialized map
	var m22 abstract.SafeMapOfMaps[string, int, float64]
	m22.Transform(func(outer string, inner int, v float64) float64 { return v + 1.0 })
	if m22.Len() != 0 {
		t.Errorf("Expected no items after Transform on uninitialized map, got %d", m22.Len())
	}

	// Test Range with uninitialized map
	var m23 abstract.SafeMapOfMaps[string, int, float64]
	called := false
	result := m23.Range(func(outer string, inner int, v float64) bool {
		called = true
		return true
	})
	if !result || called {
		t.Error("Expected Range to return true without calling function on uninitialized map")
	}

	// Test Copy with uninitialized map
	var m24 abstract.SafeMapOfMaps[string, int, float64]
	copied := m24.Copy()
	if len(copied) != 0 {
		t.Errorf("Expected empty copied map, got length %d", len(copied))
	}

	// Test Raw with uninitialized map
	var m25 abstract.SafeMapOfMaps[string, int, float64]
	raw := m25.Raw()
	if len(raw) != 0 {
		t.Errorf("Expected empty raw map, got length %d", len(raw))
	}

	// Test Clear with uninitialized map
	var m26 abstract.SafeMapOfMaps[string, int, float64]
	m26.Clear()
	if m26.Len() != 0 {
		t.Errorf("Expected length 0 after Clear on uninitialized map, got %d", m26.Len())
	}

	// Test Refill with uninitialized map
	var m27 abstract.SafeMapOfMaps[string, int, float64]
	refillData := map[string]map[int]float64{"group": {1: 1.1}}
	m27.Refill(refillData)
	if m27.Get("group", 1) != 1.1 {
		t.Errorf("Expected 1.1 after Refill on uninitialized map, got %f", m27.Get("group", 1))
	}
}
