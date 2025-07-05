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
	m := abstract.NewSafeMap[string, int](map[string]int{
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
