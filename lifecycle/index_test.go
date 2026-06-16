package lifecycle

import (
	"sync"
	"testing"
)

func TestLifecycleManager_OnExit_LIFO(t *testing.T) {
	lm := New()
	var order []int

	lm.OnExit(func() {
		order = append(order, 1)
	})
	lm.OnExit(func() {
		order = append(order, 2)
	})
	lm.OnExit(func() {
		order = append(order, 3)
	})

	lm.Exit()
	lm.Wait()

	if len(order) != 3 {
		t.Fatalf("expected 3 callbacks, got %d", len(order))
	}

	// Should be 3, 2, 1 (LIFO)
	if order[0] != 3 || order[1] != 2 || order[2] != 1 {
		t.Errorf("expected LIFO order [3, 2, 1], got %v", order)
	}
}

func TestLifecycleManager_ThreadSafety(t *testing.T) {
	lm := New()
	const count = 100
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			lm.OnExit(func() {})
		}()
	}

	wg.Wait()

	lm.Exit()
	lm.Wait()

	if len(lm.exitCallbacks) != count {
		t.Errorf("expected %d callbacks, got %d", count, len(lm.exitCallbacks))
	}
}

func TestLifecycleManager_Options(t *testing.T) {
	var printfCalled bool
	var printlnCalled bool

	lm := New(
		WithPrintf(func(format string, v ...any) {
			printfCalled = true
		}),
		WithPrintln(func(v ...any) {
			printlnCalled = true
		}),
	)

	lm.Exitf("test")
	lm.Exitln("test")
	lm.Wait()

	if !printfCalled {
		t.Error("expected custom printf to be called")
	}
	if !printlnCalled {
		t.Error("expected custom println to be called")
	}
}
