# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

`abstract` is a comprehensive Go library providing generic data structures, cryptographic utilities, concurrency helpers, and powerful abstractions. The library leverages Go 1.23+ generics to provide type-safe, efficient implementations.

## Development Commands

### Testing
```bash
# Run all tests with race detection and coverage
go test -v -race -coverprofile=coverage.txt ./...

# Run tests for a specific file (e.g., workerpool)
go test -v -race -run TestWorkerPool ./workerpool_test.go

# Run a specific test function
go test -v -race -run TestNewOrderer

# Run tests without verbose output
go test -race ./...
```

### Building
```bash
# Build the package
go build -v ./...

# Check if code compiles
go build
```

### Linting and Code Quality
```bash
# Format code
go fmt ./...

# Run vet
go vet ./...

# Run staticcheck (if installed)
staticcheck ./...
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

## Code Architecture

### Design Principles

1. **Generic-First**: The library extensively uses Go generics (1.23+) to provide type-safe abstractions
2. **Thread-Safe Variants**: Most data structures have both regular and thread-safe (`Safe*`) variants
3. **Zero External Dependencies**: Only depends on `github.com/maxbolgarin/lang` for logging utilities
4. **Minimal Allocations**: Focus on performance with careful memory management

### Core Components

#### Type Constraints (`abstract.go`)
The foundation of the library's generic system. Defines constraints like:
- `Number` - all numeric types (Integer | Float)
- `Ordered` - types supporting comparison (Integer | Float | ~string)
- `Integer`, `Signed`, `Unsigned` - various integer categories
- `Float`, `Complex` - floating point and complex number types

All generic math functions (`Min`, `Max`, `Abs`, `Pow`, etc.) use these constraints.

#### Data Structures Pattern
Each data structure follows a consistent pattern:
- Base implementation (`Map`, `Set`, `Stack`, `Slice`, `LinkedList`)
- Thread-safe wrapper (`SafeMap`, `SafeSet`, etc.) using `sync.RWMutex`
- Specialized variants (e.g., `UniqueStack`, `EntityMap`, `OrderedPairs`)

Thread-safe variants delegate to the base implementation while holding appropriate locks.

#### Worker Pools and Job Queues
The library provides multiple concurrent task execution implementations:

1. **WorkerPool[T]** (`workerpool.go`): Context-aware worker pool
   - Generic over result type `T`
   - Task signature: `func(ctx context.Context) (T, error)` (**context-aware**)
   - External context passed to `Start(ctx)`, `Submit(ctx, task)`, `Shutdown(ctx)`
   - Granular metrics: `TasksInQueue()`, `OnFlyRunningTasks()`, `FinishedTasks()`, `TotalTasks()`
   - Logger support
   - Graceful shutdown with `Shutdown(ctx)` that waits for completion
   - Use case: Context-aware async work with results, cancellation support
   - **This is the recommended worker pool for new code**

2. **JobQueue** (`jobqueue.go`): Fire-and-forget job execution
   - Task signature: `func(ctx context.Context)` (no return values)
   - Context-aware like WorkerPool
   - No result channels (lighter weight than WorkerPool)
   - `Wait(ctx)` method to block until all tasks complete
   - Same metrics as WorkerPool: `TasksInQueue()`, `OnFlyRunningTasks()`, etc.
   - Use case: Side-effect tasks without return values (logging, notifications, cleanup, etc.)

3. **WorkerPoolV2[T]** (`workerpool_v2.go`): Generic worker pool (simple)
   - Generic over result type `T`
   - Task signature: `func() (T, error)` (no context parameter)
   - Tracks submitted/running/finished task counts
   - Methods: `Submit`, `FetchResults`, `FetchAllResults`, `Submitted()`, `Running()`, `Finished()`
   - `FetchResults()` returns results for tasks submitted at call time
   - `FetchAllResults()` waits for all submitted tasks
   - Use case: Simple async work with results, no context propagation needed

4. **LegacyWorkerPool** (`legacy_workerpool.go`): Original implementation
   - Non-generic, returns `Result` structs with `any` type
   - Task signature: `func() (any, error)`
   - Internal context management
   - Methods: `Submit`, `SubmitWait`, `Results()`, `Stop`, `StopAndWait`
   - **Deprecated: Use WorkerPool for new code. This implementation is kept for backward compatibility only**

**Key Differences:**
- **Context propagation**: WorkerPoolV2 tasks don't receive context; WorkerPool and JobQueue tasks do
- **Results**: LegacyWorkerPool/V2/WorkerPool collect results; JobQueue doesn't (fire-and-forget)
- **Lifecycle**: WorkerPool and JobQueue require external context for Start/Shutdown

### Worker Pool Selection
- Use **WorkerPool[T]** for new context-aware work with results (recommended)
- Use **JobQueue** for fire-and-forget tasks (no results needed)
- Use **WorkerPoolV2[T]** only if context propagation is not needed
- **LegacyWorkerPool** is deprecated, kept for backward compatibility only

### Cryptographic Code
All crypto utilities in `crypto.go` use standard library crypto packages with secure defaults:
- AES-256-GCM for encryption
- HMAC-SHA256 for message authentication
- ECDSA P-256 for digital signatures
- Never modify crypto code without security review

### Dependencies
- Primary dependency: `github.com/maxbolgarin/lang` (provides `Logger` interface and utilities)
- Keep external dependencies minimal
- Any new dependency requires justification

### Release Process
- Automated via semantic-release (`.releaserc.json`)
- GitHub Actions workflow (`release.yaml`) handles releases
- Follow conventional commits for automatic versioning
