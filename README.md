# abstract

[![Go Version][version-img]][doc] [![GoDoc][doc-img]][doc] [![Build][ci-img]][ci] [![GoReport][report-img]][report]


**`abstract` provides a suite of generic data structures to focus on the core business logic and avoid unnecessary boilerplate**

```
go get -u github.com/maxbolgarin/abstract
```

## Overview

`abstract` package provides a suite of generic data structures to enhance code organization, concurrency safety, and utility in Go. The package aims to simplify and abstract the working processes with generic and concurrent data structures such as maps, stacks, sets, and pairs. 

### Key Features
- **Generic Data Structures**: Includes Map, Orderer, Stack, Set, and their thread-safe counterparts.
- **Clean Code**: Using of this package allowes you to focus on the core business logic and avoid unnecessary complexity.
- **Concurrency Safety**: Thread-safe structures provided for concurrent access across multiple goroutines.
- **Ease of Use**: Designed to integrate seamlessly and abstract away complex operations into simple APIs.

### Cons
- **Complexity**: Abstraction layers might add complexity for simpler use cases.
- **Performance Overhead**: Concurrency safety through mutexes can add overhead.
- **Learning Curve**: Understanding generics and constraints may require additional learning.


### Data Structures and Safe Alternatives

| Data Structure | Description | Safe Alternative | Key Methods |
|----------------|-------------|------------------|-------------|
| `Map`          | Generic map for key-value pairs. | `SafeMap` | Get, Set, Delete, Keys, Values |
| `EntityMap`    | Map of entities, providing ordering features. | `SafeEntityMap` | Set, Delete, AllOrdered |
| `Set`          | Stores keys uniquely without associated values. | `SafeSet` | Add, Remove, Has |
| `Slice`        | Generic slice implementation. | `SafeSlice` | Append, Get, Pop, Delete |
| `Stack`        | Generic stack implementation. | `SafeStack` | Push, Pop, Last |
| `UniqueStack`  | Stack without duplicates. | `SafeUniqueStack` | Push, Pop, Last |
| `LinkedList`   | Generic doubly linked list implementation. | `SafeLinkedList` | Front, Back |
| `OrderedPairs` | Maintains order of key-value pairs. | `SafeOrderedPairs` | Add, Get, Keys, Rand |


## Usage Examples

### Using Map and SafeMap

```go
package main

import (
	"fmt"
	"github.com/maxbolgarin/abstract"
)

func main() {
	// Example using Map
	m := abstract.NewMap[string, int]()
	m.Set("apple", 5)
	fmt.Println(m.Get("apple")) // Output: 5

	// Example using SafeMap
	sm := abstract.NewSafeMap[string, int]()
	sm.Set("banana", 10)
	fmt.Println(sm.Get("banana")) // Output: 10
}
```

### Using Stack and SafeStack

```go
package main

import (
	"fmt"
	"abstract"
)

func main() {
	// Example using Stack
	s := abstract.NewStack[int]()
	s.Push(10)
	fmt.Println(s.Pop()) // Output: 10

	// Example using SafeStack
	ss := abstract.NewSafeStack[int]()
	ss.Push(20)
	fmt.Println(ss.Pop()) // Output: 20
}
```

### Using EntityMap and SafeEntityMap

```go
package main

import (
	"fmt"
	"abstract"
)

// Define a simple Entity type
type MyEntity struct {
	id    string
	name  string
	order int
}

func (e MyEntity) ID() string    { return e.id }
func (e MyEntity) Name() string  { return e.name }
func (e MyEntity) Order() int    { return e.order }
func (e MyEntity) SetOrder(o int) abstract.Entity[string] {
	return MyEntity{id: e.id, name: e.name, order: o}
}

func main() {
	eMap := abstract.NewEntityMap[string, MyEntity]()
	eMap.Set(MyEntity{id: "1", name: "entity1", order: 0})
	foundEntity, ok := eMap.LookupByName("entity1")
	if ok {
		fmt.Println(foundEntity.Name()) // Output: entity1
	}
}
```


## Contributions and Issues

Feel free to contribute to this package or report issues you encounter during usage

## License

This project is licensed under the terms of the [MIT License](LICENSE).

[MIT License]: LICENSE.txt
[version-img]: https://img.shields.io/badge/Go-%3E%3D%201.19-%23007d9c
[doc-img]: https://pkg.go.dev/badge/github.com/maxbolgarin/abstract
[doc]: https://pkg.go.dev/github.com/maxbolgarin/abstract
[ci-img]: https://github.com/maxbolgarin/abstract/actions/workflows/go.yaml/badge.svg
[ci]: https://github.com/maxbolgarin/abstract/actions
[report-img]: https://goreportcard.com/badge/github.com/maxbolgarin/abstract
[report]: https://goreportcard.com/report/github.com/maxbolgarin/abstract
