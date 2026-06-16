package lifemanaged

import (
	"github.com/lucap9056/go-lifecycle/lifecycle"
)

func Run(callback func(*lifecycle.LifecycleManager) error, opts ...lifecycle.Option) error {
	lm := lifecycle.New(opts...)
	errChan := make(chan error, 1)
	go func() {
		errChan <- callback(lm)
	}()

	select {
	case err := <-errChan:
		lm.Exit()
		lm.Wait()
		return err
	case <-lm.Done():
	}

	lm.Wait()
	return nil
}
