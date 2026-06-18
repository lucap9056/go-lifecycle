package lifemanaged

import (
	"github.com/lucap9056/go-lifecycle/lifecycle"
)

func Run(callback func(*lifecycle.LifecycleManager) error, opts ...lifecycle.Option) error {
	lm := lifecycle.New(opts...)
	err := callback(lm)

	if err != nil {
		lm.Exit()
	}

	lm.Wait()
	return err
}
