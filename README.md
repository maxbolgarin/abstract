# abstract

[![Go Version][version-img]][doc] [![GoDoc][doc-img]][doc] [![Build][ci-img]][ci] [![Coverage][coverage-img]][coverage] [![GoReport][report-img]][report]

**A comprehensive Go library providing generic data structures, cryptographic utilities, concurrency helpers, and powerful abstractions to accelerate development and reduce boilerplate.**

```bash
go get -u github.com/maxbolgarin/abstract
```

## 🚀 Overview

The `abstract` package provides a comprehensive collection of utilities and data structures designed to simplify common programming tasks in Go. It leverages Go's generics to provide type-safe, efficient implementations that work across different types while maintaining excellent performance.

### 🎯 Key Features

- **🗂️ Generic Data Structures**: Maps, Sets, Stacks, Slices, LinkedLists with thread-safe variants
- **🔐 Cryptographic Utilities**: AES encryption, HMAC generation, ECDSA signing/verification
- **⚡ Concurrency Tools**: Futures/Promises, Worker pools, Rate limiting, Concurrent helpers
- **🎲 Random Generation**: Secure random strings, numbers, choices with customizable alphabets
- **📊 CSV Processing**: Advanced table manipulation and data transformation
- **📏 Mathematical Utilities**: Generic math functions with type-safe constraints
- **🆔 ID Generation**: Structured entity ID creation with type safety
- **⏱️ Timing Utilities**: Precise timing, lap timing, deadline management
- **🔧 Type Constraints**: Comprehensive type constraints for generic programming

### ✅ Advantages

- **Clean Code**: Focus on business logic, not boilerplate
- **Type Safety**: Leverage Go's generics for compile-time type checking
- **Performance**: Optimized implementations with minimal overhead
- **Concurrency**: Built-in thread-safe alternatives for concurrent environments
- **Flexibility**: Extensive customization options and extensible design

### ⚠️ Considerations

- **Learning Curve**: Requires familiarity with Go generics and constraints
- **Complexity**: Abstraction layers might add complexity for simple use cases
- **Memory Overhead**: Thread-safe variants use mutexes which add memory overhead

## 📚 Table of Contents

- [Data Structures](#-data-structures)
  - [Maps](#maps)
  - [Sets](#sets)
  - [Stacks](#stacks)
  - [Slices](#slices)
  - [LinkedLists](#linkedlists)
- [Cryptographic Utilities](#-cryptographic-utilities)
- [Concurrency Tools](#-concurrency-tools)
- [Random Generation](#-random-generation)
- [CSV Processing](#-csv-processing)
- [Mathematical Utilities](#-mathematical-utilities)
- [ID Generation](#-id-generation)
- [Timing Utilities](#-timing-utilities)
- [Type Constraints](#-type-constraints)
- [Installation](#-installation)
- [API Reference](#-api-reference)

## 🗂️ Data Structures

All data structures provide both regular and thread-safe variants (prefixed with "Safe"). Thread-safe variants use RWMutex for concurrent access while maintaining performance.

### Maps

Generic maps with enhanced functionality for key-value storage.

```go
// Basic Map operations
m := abstract.NewMap[string, int]()
m.Set("apple", 5)
m.Set("banana", 3)

fmt.Println(m.Get("apple"))   // 5
fmt.Println(m.Has("orange"))  // false
fmt.Println(m.Keys())         // [apple, banana]
fmt.Println(m.Values())       // [5, 3]

// Thread-safe variant
safeMap := abstract.NewSafeMap[string, int]()
safeMap.Set("concurrent", 42)
value := safeMap.Get("concurrent") // Safe for concurrent access

// Advanced: EntityMap for objects with ordering
type User struct {
    id    string
    name  string
    order int
}

func (u User) GetID() string    { return u.id }
func (u User) GetName() string  { return u.name }
func (u User) GetOrder() int    { return u.order }
func (u User) SetOrder(o int) abstract.Entity[string] {
    return User{id: u.id, name: u.name, order: o}
}

entityMap := abstract.NewEntityMap[string, User]()
entityMap.Set(User{id: "1", name: "Alice", order: 1})
entityMap.Set(User{id: "2", name: "Bob", order: 2})

users := entityMap.AllOrdered() // Returns users in order
user, found := entityMap.LookupByName("Alice")
```

### Sets

Efficient set implementation for unique value storage.

```go
// String set
stringSet := abstract.NewSet[string]()
stringSet.Add("apple", "banana", "apple") // Duplicates ignored
fmt.Println(stringSet.Has("apple"))       // true
fmt.Println(stringSet.Len())              // 2

// Thread-safe set
safeSet := abstract.NewSafeSet[int]()
safeSet.Add(1, 2, 3)
safeSet.Remove(2)
values := safeSet.ToSlice() // [1, 3]

// Set operations
set1 := abstract.NewSet[int]()
set1.Add(1, 2, 3)
set2 := abstract.NewSet[int]()
set2.Add(3, 4, 5)

intersection := set1.Intersection(set2) // [3]
union := set1.Union(set2)               // [1, 2, 3, 4, 5]
```

### Stacks

LIFO (Last In, First Out) data structure implementations.

```go
// Basic stack
stack := abstract.NewStack[string]()
stack.Push("first")
stack.Push("second")
stack.Push("third")

fmt.Println(stack.Pop())  // "third"
fmt.Println(stack.Last()) // "second" (peek without removing)
fmt.Println(stack.Len())  // 2

// Thread-safe stack
safeStack := abstract.NewSafeStack[int]()
safeStack.Push(10, 20, 30)
value := safeStack.Pop() // 30

// UniqueStack - no duplicates
uniqueStack := abstract.NewUniqueStack[string]()
uniqueStack.Push("a", "b", "a") // Only adds "a" once
fmt.Println(uniqueStack.Len())  // 2
```

### Slices

Enhanced slice operations with generic support.

```go
// Basic slice operations
slice := abstract.NewSlice[int]()
slice.Append(1, 2, 3, 4, 5)
slice.Insert(2, 99)     // Insert 99 at index 2
fmt.Println(slice.Get(2)) // 99

value := slice.Pop()      // Remove and return last element
slice.Delete(0)          // Remove first element
slice.Reverse()          // Reverse in-place

// Thread-safe slice
safeSlice := abstract.NewSafeSlice[string]()
safeSlice.Append("hello", "world")
safeSlice.Sort() // Sorts in-place for comparable types

// Slice utilities
numbers := []int{1, 2, 3, 4, 5}
doubled := abstract.Map(numbers, func(x int) int { return x * 2 })
evens := abstract.Filter(numbers, func(x int) bool { return x%2 == 0 })
```

### LinkedLists

Doubly linked list implementation with efficient insertion/deletion.

```go
// Basic linked list
list := abstract.NewLinkedList[string]()
list.PushFront("first")
list.PushBack("last")
list.PushFront("new first")

fmt.Println(list.Front()) // "new first"
fmt.Println(list.Back())  // "last"
fmt.Println(list.Len())   // 3

// Iterate through list
for elem := list.Front(); elem != nil; elem = elem.Next() {
    fmt.Println(elem.Value)
}

// Thread-safe linked list
safeList := abstract.NewSafeLinkedList[int]()
safeList.PushBack(1, 2, 3)
value := safeList.PopFront() // 1
```

## 🔐 Cryptographic Utilities

Secure cryptographic operations with best-practice implementations.

```go
// AES Encryption
key := abstract.NewEncryptionKey()
plaintext := []byte("sensitive data")

// Encrypt
ciphertext, err := abstract.EncryptAES(plaintext, key)
if err != nil {
    log.Fatal(err)
}

// Decrypt
decrypted, err := abstract.DecryptAES(ciphertext, key)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Decrypted: %s\n", decrypted)

// HMAC Generation
hmacKey := abstract.NewHMACKey()
message := []byte("important message")
mac := abstract.GenerateHMAC(message, hmacKey)

// Verify HMAC
isValid := abstract.CheckHMAC(message, mac, hmacKey)
fmt.Printf("HMAC valid: %v\n", isValid)

// ECDSA Signing
privateKey, err := abstract.NewSigningKey()
if err != nil {
    log.Fatal(err)
}

signature, err := abstract.SignData(message, privateKey)
if err != nil {
    log.Fatal(err)
}

// Verify signature
publicKey := &privateKey.PublicKey
isValid = abstract.VerifySign(message, signature, publicKey)
fmt.Printf("Signature valid: %v\n", isValid)

// Key encoding/decoding
pemPrivateKey, _ := abstract.EncodePrivateKey(privateKey)
pemPublicKey, _ := abstract.EncodePublicKey(publicKey)

decodedPrivate, _ := abstract.DecodePrivateKey(pemPrivateKey)
decodedPublic, _ := abstract.DecodePublicKey(pemPublicKey)
```

## ⚡ Concurrency Tools

Powerful concurrency utilities for async operations and parallel processing.

### Futures and Promises

```go
// Basic Future
ctx := context.Background()
future := abstract.NewFuture(ctx, logger, func(ctx context.Context) (string, error) {
    time.Sleep(100 * time.Millisecond)
    return "async result", nil
})

result, err := future.Get(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %s\n", result)

// Future with timeout
result, err = future.GetWithTimeout(ctx, 50*time.Millisecond)
if err == abstract.ErrTimeout {
    fmt.Println("Operation timed out")
}

// Waiter for void operations
waiter := abstract.NewWaiter(ctx, logger, func(ctx context.Context) error {
    // Perform some work
    return nil
})
err = waiter.Await(ctx)

// WaiterSet for multiple operations
waiterSet := abstract.NewWaiterSet(logger)
waiterSet.Add(ctx, func(ctx context.Context) error {
    // Task 1
    return nil
})
waiterSet.Add(ctx, func(ctx context.Context) error {
    // Task 2
    return nil
})
err = waiterSet.Await(ctx) // Wait for all tasks
```

### Worker Pools

```go
// Create generic worker pool
pool := abstract.NewWorkerPoolV2[string](5, 100) // 5 workers, queue capacity 100
pool.Start()
defer pool.Stop()

// Submit tasks
for i := 0; i < 10; i++ {
    i := i // Capture loop variable
    task := func() (string, error) {
        return fmt.Sprintf("Task %d result", i), nil
    }
    
    if !pool.Submit(task) {
        fmt.Println("Failed to submit task")
    }
}

// Fetch results (waits for all submitted tasks at call time)
results, errors := pool.FetchResults(5 * time.Second)
for i, result := range results {
    if errors[i] != nil {
        fmt.Printf("Task error: %v\n", errors[i])
    } else {
        fmt.Printf("Task result: %v\n", result)
    }
}

// Submit with timeout
if !pool.Submit(task, 100*time.Millisecond) {
    fmt.Println("Task submission timed out")
}

// Monitor pool status
fmt.Printf("Submitted: %d, Running: %d, Finished: %d\n", 
    pool.Submitted(), pool.Running(), pool.Finished())

// Fetch all results (including tasks submitted after call)
allResults, allErrors := pool.FetchAllResults(10 * time.Second)
```

### Rate Limiting

```go
// Rate-limited processing
ctx := context.Background()
processor := abstract.NewRateProcessor(ctx, 10) // Max 10 operations per second

// Add rate-limited tasks
for i := 0; i < 100; i++ {
    i := i
    processor.AddTask(func(ctx context.Context) error {
        // This will be rate-limited
        fmt.Printf("Processing item %d\n", i)
        return nil
    })
}

// Wait for completion
errors := processor.Wait()
if len(errors) > 0 {
    fmt.Printf("Encountered %d errors\n", len(errors))
}
```

### Concurrent Helpers

```go
// Periodic updater
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

abstract.StartUpdater(ctx, 30*time.Second, logger, func() {
    fmt.Println("Periodic health check")
})

// Immediate execution, then periodic
abstract.StartUpdaterNow(ctx, time.Minute, logger, func() {
    fmt.Println("Sync operation")
})

// With shutdown handler
abstract.StartUpdaterWithShutdown(ctx, 10*time.Second, logger, 
    func() {
        // Periodic work
    },
    func() {
        fmt.Println("Shutting down...")
    },
)
```

## 🎲 Random Generation

Cryptographically secure random generation with multiple character sets.

```go
// Basic random strings
token := abstract.GetRandomString(32)        // Hexadecimal
sessionID := abstract.GetRandomAlphaNumeric(16) // Letters + numbers
pin := abstract.GetRandomNumeric(6)          // Numbers only
code := abstract.GetRandomLowerAlpha(8)      // Lowercase letters
id := abstract.GetRandomUpperAlpha(4)        // Uppercase letters

// Custom alphabet
customAlphabet := []byte("!@#$%^&*")
symbols := abstract.GetRandomStringWithAlphabet(8, customAlphabet)

// Random numbers and choices
dice := abstract.GetRandomInt(1, 6)
coinFlip := abstract.GetRandomBool()

colors := []string{"red", "green", "blue", "yellow"}
color, ok := abstract.GetRandomChoice(colors)
if ok {
    fmt.Printf("Random color: %s\n", color)
}

// Shuffle operations
cards := []string{"A", "K", "Q", "J", "10", "9", "8", "7"}
abstract.ShuffleSlice(cards)
fmt.Printf("Shuffled: %v\n", cards)

// Random network addresses (for testing)
addr := abstract.GetRandListenAddress() // ":12345" (random port)
```

## 📊 CSV Processing

Advanced CSV table manipulation with indexing and querying capabilities.

```go
// Load CSV from file
table, err := abstract.NewCSVTableFromFilePath("data.csv")
if err != nil {
    log.Fatal(err)
}

// Create from data
data := map[string]map[string]string{
    "user1": {"name": "Alice", "age": "30", "city": "NYC"},
    "user2": {"name": "Bob", "age": "25", "city": "LA"},
}
table = abstract.NewCSVTableFromMap(data, "userID")

// Query operations
row := table.Row("user1")
value := table.Value("user1", "name") // "Alice"
exists := table.Has("user3")          // false

// Find operations
criteria := map[string]string{"city": "NYC"}
foundID, foundRow := table.FindRow(criteria)
allNYCUsers := table.Find(criteria)

// Modifications
table.AddRow("user3", map[string]string{
    "name": "Charlie",
    "age": "35",
    "city": "Chicago",
})

table.UpdateRow("user1", map[string]string{
    "age": "31",
})

table.AppendColumn("country", []string{"USA", "USA", "USA"})
table.DeleteColumn("age")
table.DeleteRow("user2")

// Export
csvBytes := table.Bytes()
fmt.Printf("CSV Data:\n%s", csvBytes)

// Thread-safe operations
safeTable := abstract.NewCSVTableSafe(records)
safeTable.AddRow("concurrent", map[string]string{
    "name": "Thread Safe",
    "value": "42",
})
```

## 📏 Mathematical Utilities

Generic mathematical functions with type-safe constraints.

```go
// Type-safe arithmetic
result := abstract.Min(1, 2, 3, 4, 5)     // 1
maximum := abstract.Max(3.14, 2.71, 1.41) // 3.14
absolute := abstract.Abs(-42)             // 42
power := abstract.Pow(2, 3)               // 8
rounded := abstract.Round(3.7)            // 4

// String conversions
str := abstract.Itoa(42)        // "42"
num, err := abstract.Atoi[int]("123") // 123

// Type constraints in action
func Average[T abstract.Number](values []T) T {
    if len(values) == 0 {
        return 0
    }
    
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum / T(len(values))
}

// Usage with different numeric types
intAvg := Average([]int{1, 2, 3, 4, 5})           // 3
floatAvg := Average([]float64{1.5, 2.5, 3.5})     // 2.5
```

## 🆔 ID Generation

Structured, type-safe ID generation with entity types.

```go
// Register entity types
const (
    UserEntity  = abstract.RegisterEntityType("USER")
    PostEntity  = abstract.RegisterEntityType("POST")
    AdminEntity = abstract.RegisterEntityType("ADMN")
)

// Generate IDs
userID := abstract.NewID(UserEntity)   // "USERa1b2c3d4e5f6"
postID := abstract.NewID(PostEntity)   // "POST9z8y7x6w5v4u"
adminID := abstract.NewID(AdminEntity) // "ADMNx1y2z3w4v5u6"

// Test IDs
testID := abstract.NewTestID() // "00x0" + random

// ID operations
entityType := abstract.FetchEntityType(userID) // "USER"
convertedID := abstract.FromID(userID, AdminEntity) // Convert user ID to admin ID

// Builder pattern for repeated ID generation
userBuilder := abstract.WithEntityType(UserEntity)
postBuilder := abstract.WithEntityType(PostEntity)

// Generate multiple IDs of the same type
users := make([]string, 10)
for i := range users {
    users[i] = userBuilder.NewID()
}

// Configure entity type size
abstract.SetEntitySize(5) // Use 5-character entity types
longEntity := abstract.RegisterEntityType("USERS")
```

## ⏱️ Timing Utilities

Precise timing measurements with advanced features.

```go
// Basic timing
timer := abstract.StartTimer()
time.Sleep(100 * time.Millisecond)
elapsed := timer.ElapsedTime()
fmt.Printf("Elapsed: %v\n", elapsed)

// Different time units
seconds := timer.ElapsedSeconds()
minutes := timer.ElapsedMinutes()
milliseconds := timer.ElapsedMilliseconds()
microseconds := timer.ElapsedMicroseconds()

// Lap timing
timer.Reset()
doWork1()
lap1 := timer.Lap()
doWork2()
lap2 := timer.Lap()
doWork3()
lap3 := timer.Lap()

fmt.Printf("Lap 1: %v, Lap 2: %v, Lap 3: %v\n", lap1, lap2, lap3)

// Pause and resume
timer.Pause()
time.Sleep(time.Second) // This won't count
timer.Resume()

// Deadline management
deadlineTimer := abstract.Deadline(5 * time.Minute)
for !deadlineTimer.IsExpired() {
    // Do work with deadline
    remaining := deadlineTimer.TimeRemaining()
    fmt.Printf("Time remaining: %v\n", remaining)
    time.Sleep(time.Second)
}

// Formatting
formatted := timer.Format("%02d:%02d:%02d.%03d") // "00:01:23.456"
short := timer.FormatShort()                     // "1m23s" or "456ms"

// Conditional timing
if timer.HasElapsed(30 * time.Second) {
    fmt.Println("30 seconds have passed")
}
```

## 🔧 Type Constraints

Comprehensive type constraints for generic programming.

```go
// Numeric constraints
func AddNumbers[T abstract.Number](a, b T) T {
    return a + b
}

// Works with integers and floats
intResult := AddNumbers(5, 3)        // 8
floatResult := AddNumbers(3.14, 2.86) // 6.0

// Ordered types (support comparison)
func Clamp[T abstract.Ordered](value, min, max T) T {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

clamped := Clamp(15, 10, 20)    // 15
clamped = Clamp(5, 10, 20)      // 10
clamped = Clamp("m", "a", "z")  // "m"

// Integer-specific operations
func IsEven[T abstract.Integer](n T) bool {
    return n%2 == 0
}

even := IsEven(4)    // true
odd := IsEven(7)     // false

// Signed integer operations
func AbsInt[T abstract.Signed](x T) T {
    if x < 0 {
        return -x
    }
    return x
}

positive := AbsInt(-42) // 42

// Available constraints:
// - abstract.Signed: ~int | ~int8 | ~int16 | ~int32 | ~int64
// - abstract.Unsigned: ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
// - abstract.Integer: Signed | Unsigned
// - abstract.Float: ~float32 | ~float64
// - abstract.Complex: ~complex64 | ~complex128
// - abstract.Number: Integer | Float
// - abstract.Ordered: Integer | Float | ~string
```

## 🛠️ Installation

```bash
go get -u github.com/maxbolgarin/abstract
```

**Requirements:**
- Go 1.23 or higher
- Dependencies: `github.com/maxbolgarin/lang v1.5.0`

## 📖 API Reference

### Core Data Structures

| Type | Description | Safe Alternative | Key Methods |
|------|-------------|------------------|-------------|
| `Map[K, V]` | Generic key-value map | `SafeMap[K, V]` | `Get`, `Set`, `Delete`, `Keys`, `Values`, `Has` |
| `Set[T]` | Unique value collection | `SafeSet[T]` | `Add`, `Remove`, `Has`, `Union`, `Intersection` |
| `Stack[T]` | LIFO data structure | `SafeStack[T]` | `Push`, `Pop`, `Last`, `Len` |
| `UniqueStack[T]` | LIFO without duplicates | `SafeUniqueStack[T]` | `Push`, `Pop`, `Last`, `Has` |
| `Slice[T]` | Enhanced slice operations | `SafeSlice[T]` | `Append`, `Get`, `Insert`, `Delete`, `Sort` |
| `LinkedList[T]` | Doubly linked list | `SafeLinkedList[T]` | `PushFront`, `PushBack`, `PopFront`, `PopBack` |
| `EntityMap[K, T]` | Map with entity ordering | `SafeEntityMap[K, T]` | `Set`, `AllOrdered`, `LookupByName` |
| `OrderedPairs[K, V]` | Ordered key-value pairs | `SafeOrderedPairs[K, V]` | `Add`, `Get`, `Keys`, `Rand` |

### Utility Types

| Type | Description | Key Methods |
|------|-------------|-------------|
| `Orderer[T]` | Order management | `Add`, `Apply`, `Get`, `Clear` |
| `Memorizer[T]` | Thread-safe single value store | `Set`, `Get`, `Pop` |
| `CSVTable` | CSV data manipulation | `Row`, `Find`, `AddRow`, `UpdateRow` |
| `Timer` | Precise timing measurements | `ElapsedTime`, `Lap`, `Pause`, `Resume` |
| `WorkerPoolV2[T]` | Generic concurrent task processing | `Submit`, `FetchResults`, `FetchAllResults`, `Submitted`, `Running`, `Finished`, `Stop` |
| `RateProcessor` | Rate-limited processing | `AddTask`, `Wait` |
| `Future[T]` | Async operation result | `Get`, `GetWithTimeout` |

### Function Categories

| Category | Functions |
|----------|-----------|
| **Math** | `Min`, `Max`, `Abs`, `Pow`, `Round`, `Itoa`, `Atoi` |
| **Random** | `GetRandomString`, `GetRandomInt`, `GetRandomBool`, `GetRandomChoice`, `ShuffleSlice` |
| **Crypto** | `EncryptAES`, `DecryptAES`, `GenerateHMAC`, `SignData`, `VerifySign` |
| **Timing** | `StartTimer`, `Deadline`, `StartUpdater`, `StartCycle` |
| **ID Gen** | `NewID`, `RegisterEntityType`, `FetchEntityType`, `WithEntityType` |
| **CSV** | `NewCSVTableFromFilePath`, `NewCSVTableFromMap`, `NewCSVTableFromReader` |

## 🤝 Contributing

We welcome contributions! Please feel free to:

1. Report bugs and issues
2. Suggest new features
3. Submit pull requests
4. Improve documentation

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**Happy coding with abstract! 🎉**

[MIT License]: LICENSE
[version-img]: https://img.shields.io/badge/Go-%3E%3D%201.23-%23007d9c
[doc-img]: https://pkg.go.dev/badge/github.com/maxbolgarin/abstract
[doc]: https://pkg.go.dev/github.com/maxbolgarin/abstract
[ci-img]: https://github.com/maxbolgarin/abstract/actions/workflows/go.yaml/badge.svg
[ci]: https://github.com/maxbolgarin/abstract/actions
[report-img]: https://goreportcard.com/badge/github.com/maxbolgarin/abstract
[report]: https://goreportcard.com/report/github.com/maxbolgarin/abstract
[coverage-img]: https://codecov.io/gh/maxbolgarin/abstract/branch/main/graph/badge.svg
[coverage]: https://codecov.io/gh/maxbolgarin/abstract