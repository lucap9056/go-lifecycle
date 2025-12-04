# go-lifecycle

[![Go Reference](https://pkg.go.dev/badge/github.com/lucap9056/go-lifecycle.svg)](https://pkg.go.dev/github.com/lucap9056/go-lifecycle/lifecycle)

`go-lifecycle` is a lightweight Go module that provides a simple and effective way to manage your application's lifecycle, ensuring a graceful shutdown when termination signals are received.

## Features

-   **Graceful Shutdown:** Automatically handles `SIGINT` (Ctrl+C) and `SIGTERM` signals to initiate a clean shutdown process.
-   **Cleanup Hooks:** Register multiple `OnExit` callback functions for cleanup tasks like closing database connections, flushing logs, or stopping services.
-   **Manual Shutdown:** Trigger a shutdown programmatically from anywhere in your application using `Exit()`, `Exitf()`, or `Exitln()`.
-   **Context-Based:** Uses Go's `context.Context` to propagate cancellation signals throughout your application's components.
-   **Easy to Use:** Simple and intuitive API for straightforward integration into any Go application.

## Installation

```bash
go get github.com/lucap9056/go-lifecycle/lifecycle
```

## Usage

Here is a basic example of how to use the `LifecycleManager` to manage an application that runs a background task.

```go
package main

import (
	"fmt"
	"time"

	"github.com/lucap9056/go-lifecycle/lifecycle"
)

func main() {
	// 1. Initialize the lifecycle manager
	lm := lifecycle.New()
	fmt.Println("Application starting... Press Ctrl+C to exit.")

	// 2. Register cleanup functions
	lm.OnExit(func() {
		fmt.Println("Cleanup task 1: Closing database connections...")
		time.Sleep(500 * time.Millisecond) // Simulate cleanup work
	})

	lm.OnExit(func() {
		fmt.Println("Cleanup task 2: Shutting down services...")
		time.Sleep(300 * time.Millisecond) // Simulate cleanup work
	})

	// 3. Start your application's services or background tasks
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-lm.Done(): // Listen for shutdown signal
				fmt.Println("Background task received shutdown signal.")
				return
			case t := <-ticker.C:
				fmt.Printf("Background task is running at %v\n", t.Format(time.RFC3339))
			}
		}
	}()

	// 4. Wait for a shutdown signal
	// The Wait() function blocks until a signal is received or Exit() is called.
	// It then executes the OnExit callbacks before unblocking.
	lm.Wait()

	fmt.Println("Application has shut down gracefully.")
}
```

## API Reference

### `func New() *LifecycleManager`

Creates and returns a new `LifecycleManager`. It sets up a context that is cancelled upon receiving `SIGINT` or `SIGTERM`.

### `func (lm *LifecycleManager) Wait()`

Blocks the current goroutine until a shutdown signal is received. Once the signal is caught, it executes all registered `OnExit` callbacks in the order they were added, and then unblocks.

### `func (lm *LifecycleManager) Done() <-chan struct{}`

Returns a channel that is closed when the application is shutting down. Components in your application can select on this channel to know when to begin their own shutdown procedure.

### `func (lm *LifecycleManager) OnExit(callback func())`

Registers a callback function to be executed during the shutdown sequence. Callbacks are executed sequentially in the order they are registered.

### `func (lm *LifecycleManager) Exit()`

Manually initiates a graceful shutdown. This is useful for triggering a shutdown from within your application's logic (e.g., in a goroutine after a critical error).

### `func (lm *LifecycleManager) Exitf(format string, v ...any)`

Logs a formatted message and then initiates a graceful shutdown.

### `func (lm *LifecycleManager) Exitln(v ...any)`

Logs a message using `log.Println` and then initiates a graceful shutdown.

