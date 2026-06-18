package lifecycle

import (
	"testing"
)

func TestOnExitDeadlock(t *testing.T) {
	lm := New()
	lm.OnExit(func() {
		// This should not deadlock anymore
		lm.OnExit(func() {})
	})
	lm.Exit()
	lm.Wait()
}

func TestWaitIdempotency(t *testing.T) {
	count := 0
	lm := New()
	lm.OnExit(func() {
		count++
	})
	lm.Exit()
	lm.Wait()
	lm.Wait()
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
}
