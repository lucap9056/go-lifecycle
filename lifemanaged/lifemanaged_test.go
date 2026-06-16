package lifemanaged

import (
	"errors"
	"testing"
	"time"

	"github.com/lucap9056/go-lifecycle/lifecycle"
)

func TestRun_Success(t *testing.T) {
	err := Run(func(lm *lifecycle.LifecycleManager) error {
		return nil
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestRun_Error(t *testing.T) {
	expectedErr := errors.New("business logic error")
	err := Run(func(lm *lifecycle.LifecycleManager) error {
		return expectedErr
	})
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestRun_Cleanup(t *testing.T) {
	cleanupCalled := false
	err := Run(func(lm *lifecycle.LifecycleManager) error {
		lm.OnExit(func() {
			cleanupCalled = true
		})
		return nil
	})

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if !cleanupCalled {
		t.Error("expected cleanup to be called")
	}
}

func TestRun_ExternalExit(t *testing.T) {
	err := Run(func(lm *lifecycle.LifecycleManager) error {
		go func() {
			time.Sleep(100 * time.Millisecond)
			lm.Exit()
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
}
