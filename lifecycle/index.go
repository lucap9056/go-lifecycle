package lifecycle

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Package lifecycle provides a simple way to manage the lifecycle of your
// application, including graceful shutdown on receiving common operating system
// signals like SIGINT (Ctrl+C) and SIGTERM (termination request).

// LifecycleManager helps coordinate the startup and shutdown of your application
// components. Create a new instance using New(), and then have your components
// listen to the context provided by the LifecycleManager to know when to shut down.
type LifecycleManager struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates and returns a new LifecycleManager.
//
// This function initializes the lifecycle management system by setting up a
// context that will be cancelled when a termination signal is received. It also
// starts a background goroutine to monitor for these signals.
//
// Example usage:
//
//	lifecycle := lifecycle.New()
//	lifecycle.Wait(5 * time.Second) // Wait for shutdown, with a potential 5-second delay
//
//	// The application will now shut down gracefully when SIGINT or SIGTERM is received.
func New() *LifecycleManager {
	ctx, cancel := context.WithCancel(context.Background())

	lifecycle := &LifecycleManager{
		ctx:    ctx,
		cancel: cancel,
	}

	go lifecycle.handleSignals()

	return lifecycle
}

// Wait blocks the current goroutine until the application receives a termination
// signal (SIGINT or SIGTERM) or one of the Exit functions is called.
//
// The optional duration 'd' allows you to specify a delay after the shutdown
// signal is received before this function returns. This can be useful to give
// your application components some time to gracefully shut down. If 'd' is zero,
// it returns immediately after the context is cancelled.
//
// Call this function in your main goroutine to keep the application running
// until a shutdown is initiated.
//
// Example:
//
//	lifecycle := lifecycle.New()
//	// ... start other services ...
//	lifecycle.Wait(10 * time.Second) // Wait for shutdown, allowing up to 10 seconds for cleanup
//	log.Println("Application finished shutting down.")
func (lifecycle *LifecycleManager) Wait(d time.Duration) {
	<-lifecycle.ctx.Done()
	if d != 0 {
		time.Sleep(d)
	}
}

// Exit initiates a graceful shutdown of the application.
//
// This function cancels the internal context, signaling to all components
// that are listening to this context to begin their shutdown procedures.
// It does not log any message before initiating the shutdown.
//
// **This function is particularly useful for triggering shutdown from within
// asynchronous functions or goroutines.** When a critical error or a specific
// condition is met in a background process, calling `Exit()` will signal the
// main application to start its shutdown sequence.
//
// Example:
//
//	go func() {
//	    // ... some background work ...
//	    if criticalErrorOccurred {
//	        lifecycleManager.Exit() // Trigger shutdown from within the goroutine
//	        return
//	    }
//	    // ... more work ...
//	}()
func (lifecycle *LifecycleManager) Exit() {
	lifecycle.cancel()
}

// Exitf logs a formatted message and then initiates a graceful shutdown of the
// application.
//
// This function is useful when you want to log a specific reason or error
// message before the application begins its shutdown process. Like `Exit()`,
// this can be called from asynchronous functions.
//
// Example:
//
//	go func() {
//	    // ... some background task ...
//	    if dataValidationFailed {
//	        lifecycleManager.Exitf("Data validation failed for ID: %d", id)
//	        return
//	    }
//	    // ... continue processing ...
//	}()
func (lifecycle *LifecycleManager) Exitf(format string, v ...any) {
	log.Printf(format, v...)
	lifecycle.cancel()
}

// Exitln logs a message with a newline and then initiates a graceful shutdown
// of the application.
//
// Similar to Exitf, this function allows you to log a message before shutting
// down, but it formats the output using log.Println. This can also be called
// from within goroutines.
//
// Example:
//
//	go func() {
//	    // ... monitoring task ...
//	    if resourceUsageExceeded {
//	        lifecycleManager.Exitln("Resource usage limit exceeded, initiating shutdown.")
//	        return
//	    }
//	    // ... continue monitoring ...
//	}()
func (lifecycle *LifecycleManager) Exitln(v ...any) {
	log.Println(v...)
	lifecycle.cancel()
}

func (lifecycle *LifecycleManager) handleSignals() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	lifecycle.cancel()
}
