package lifemanaged

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/lucap9056/go-lifecycle/lifecycle"
)

func TestRun_Success(t *testing.T) {
	done := make(chan struct{})
	go func() {
		err := Run(func(lm *lifecycle.LifecycleManager) error {
			go func() {
				time.Sleep(100 * time.Millisecond)
				syscall.Kill(os.Getpid(), syscall.SIGTERM)
			}()
			return nil
		})
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("test timed out")
	}
}

func TestRun_Error(t *testing.T) {
	expectedErr := errors.New("business logic error")
	done := make(chan struct{})
	go func() {
		err := Run(func(lm *lifecycle.LifecycleManager) error {
			return expectedErr
		})
		if err != expectedErr {
			t.Errorf("expected %v, got %v", expectedErr, err)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("test timed out")
	}
}

func TestRun_Cleanup(t *testing.T) {
	cleanupCalled := false
	done := make(chan struct{})
	go func() {
		err := Run(func(lm *lifecycle.LifecycleManager) error {
			lm.OnExit(func() {
				cleanupCalled = true
			})
			go func() {
				time.Sleep(100 * time.Millisecond)
				syscall.Kill(os.Getpid(), syscall.SIGTERM)
			}()
			return nil
		})

		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if !cleanupCalled {
			t.Error("expected cleanup to be called")
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("test timed out")
	}
}

func TestRun_ExternalExit(t *testing.T) {
	done := make(chan struct{})
	go func() {
		err := Run(func(lm *lifecycle.LifecycleManager) error {
			go func() {
				time.Sleep(100 * time.Millisecond)
				syscall.Kill(os.Getpid(), syscall.SIGTERM)
			}()

			select {
			case <-lm.Done():
				return nil
			case <-time.After(1 * time.Second):
				return errors.New("timeout")
			}
		})

		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("test timed out")
	}
}
