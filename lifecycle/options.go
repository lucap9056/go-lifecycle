package lifecycle

import "log"

type Options struct {
	printf  func(format string, v ...any)
	println func(v ...any)
}

type Option = func(*Options)

func defaultOptions() *Options {
	return &Options{
		printf: func(format string, v ...any) {
			log.Printf(format, v...)
		},
		println: func(v ...any) {
			log.Println(v...)
		},
	}
}

// WithPrintf sets the printf function used for logging formatted messages.
func WithPrintf(f func(format string, v ...any)) Option {
	return func(o *Options) {
		o.printf = f
	}
}

// WithPrintln sets the println function used for logging messages with a newline.
func WithPrintln(f func(v ...any)) Option {
	return func(o *Options) {
		o.println = f
	}
}
