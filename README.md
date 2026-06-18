# go-lifecycle

[![Go Reference](https://pkg.go.dev/badge/github.com/lucap9056/go-lifecycle.svg)](https://pkg.go.dev/github.com/lucap9056/go-lifecycle/lifecycle)

`go-lifecycle` is a lightweight Go module that provides a simple and effective way to manage your application's lifecycle, ensuring a graceful shutdown when termination signals are received.

## Features
-   **Graceful Shutdown:** Automatically handles `SIGINT` (Ctrl+C) and `SIGTERM` signals.
-   **Cleanup Hooks (LIFO):** Register `OnExit` callbacks that execute in reverse order of registration, ensuring dependencies are shut down correctly.

## Installation

```bash
go get github.com/lucap9056/go-lifecycle
```

## Usage

### Tier 1: Core Lifecycle (`lifecycle`)

Use this when you want full control over when to initialize and wait for shutdown.

```go
package main

import (
	"fmt"
	"github.com/lucap9056/go-lifecycle/lifecycle"
)

func main() {
	lm := lifecycle.New()

	lm.OnExit(func() {
		fmt.Println("Cleaning up resources...")
	})

	// Main application logic...
	
	lm.Wait() // Blocks until SIGINT/SIGTERM or Exit()
}
```

### Tier 2: Managed Runner (`lifemanaged`)

Use this for a more declarative approach with centralized error handling. `Run` will block until the provided function returns an error, or a termination signal is received.

```go
package main

import (
	"fmt"
	"github.com/lucap9056/go-lifecycle/lifecycle"
	"github.com/lucap9056/go-lifecycle/lifemanaged"
)

func main() {
	err := lifemanaged.Run(func(lm *lifecycle.LifecycleManager) error {
		lm.OnExit(func() {
			fmt.Println("Cleaning up...")
		})

		fmt.Println("Application is running. Press Ctrl+C to exit.")
		
		// Run your server or business logic here.
		// If you return nil, lifemanaged.Run will wait for a signal.
		// If you return an error, it triggers immediate graceful shutdown.
		return nil 
	})

	if err != nil {
		fmt.Printf("Application exited with error: %v\n", err)
	}
}
```

## Configuration (Options)

You can customize the logging behavior using the Options pattern:

```go
lm := lifecycle.New(
    lifecycle.WithPrintf(myCustomLogger.Printf),
    lifecycle.WithPrintln(myCustomLogger.Println),
)
```

## API Reference

### `lifecycle` Package

- **`New(opts ...Option) *LifecycleManager`**: Creates a new manager.
- **`Wait()`**: Blocks until shutdown, then executes callbacks in **LIFO** order.
- **`OnExit(callback func())`**: Registers a cleanup function (Thread-safe).
- **`Done() <-chan struct{}`**: Returns the shutdown notification channel.
- **`Exit()` / `Exitf()` / `Exitln()`**: Manually triggers graceful shutdown.

### `lifemanaged` Package

- **`Run(fn func(*lifecycle.LifecycleManager) error, opts ...lifecycle.Option) error`**: Encapsulates the entire lifecycle. Runs the provided function. If the function returns an error, it triggers shutdown immediately. If it returns `nil`, it waits for a termination signal. Cleanup hooks are guaranteed to run before `Run` returns.

