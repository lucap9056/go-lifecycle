package lifecycle

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Package lifecycle provides a simple way to manage the lifecycle of your
// application, including graceful shutdown on receiving common operating system
// signals like SIGINT (Ctrl+C) and SIGTERM (termination request).

// LifecycleManager helps coordinate the startup and shutdown of your application
// components. Create a new instance using New(), and then have your components
// listen to the context provided by the LifecycleManager to know when to shut down.
type LifecycleManager struct {
	ctx           context.Context
	cancel        context.CancelFunc
	exitCallbacks []func()
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
// Call this function in your main goroutine to keep the application running
// until a shutdown is initiated.
//
// Example:
//
//	lifecycle := lifecycle.New()
//	// ... start other services ...
//	lifecycle.Wait() // Wait for shutdown
//	log.Println("Application finished shutting down.")
func (lifecycle *LifecycleManager) Wait() {
	<-lifecycle.ctx.Done()
	for _, callback := range lifecycle.exitCallbacks {
		callback()
	}
}

// Done returns a channel that is closed when the application is shutting down.
//
// Services and background goroutines can select on this channel to listen for
// the shutdown signal and perform a graceful exit. This is the preferred way to
// integrate with the lifecycle manager.
//
// Example:
//
//	go func() {
//	    select {
//	    case <-lifecycle.Done():
//	        log.Println("Shutting down worker...")
//	        return
//	    default:
//	        // ... do work ...
//	    }
//	}()
func (lifecycle *LifecycleManager) Done() <-chan struct{} {
	return lifecycle.ctx.Done()
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

// OnExit registers a function to be executed when the application is shutting down.
//
// You can call this function multiple times to register multiple cleanup tasks.
// The registered functions will be executed in the order they were registered.
// This is useful for tasks like closing database connections, flushing logs, or
// waiting for background goroutines to finish.
//
// Example:
//
//	lifecycle.OnExit(func() {
//	    log.Println("Closing database connection...")
//	    db.Close()
//	})
//
//	lifecycle.OnExit(func() {
//	    log.Println("Shutting down HTTP server...")
//	    server.Shutdown(context.Background())
//	})
func (lifecycle *LifecycleManager) OnExit(callback func()) {
	lifecycle.exitCallbacks = append(lifecycle.exitCallbacks, callback)
}

func (lifecycle *LifecycleManager) handleSignals() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	lifecycle.cancel()
}
