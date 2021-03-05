package eio

import (
	"io"
)

var _ io.WriteCloser = &customWriteCloser{}

type customWriteCloser struct {
	closer func() error
	writer func(p []byte) (n int, err error)
}

// Write implements the io.WriteCloser interface but defferes the write to the callback.
func (cwc *customWriteCloser) Write(p []byte) (n int, err error) {
	return cwc.writer(p)
}

func (cwc *customWriteCloser) Close() (err error) {
	return cwc.closer()
}

// NewPreWriteCallbacks allows hooks to be called prior to writes. Errors in any callback will
// block/stop writes to the underlying writer. Callbacks are called in blocking sequence.
// Can be use to block/delay writes as well as conditionally abort writes.
func NewPreWriteCallbacks(w io.WriteCloser, callbacks ...func([]byte) error) io.WriteCloser {
	return &customWriteCloser{
		closer: w.Close,
		writer: func(p []byte) (n int, err error) {
			for _, cb := range callbacks {
				if err := cb(p); err != nil {
					return 0, err
				}
			}
			return w.Write(p)
		},
	}
}

// NewPostWriteCallbacks will trigger callbacks after a write has taken place. As such they can not block
// nor cancel a write but only be informed about the success of it. The callback will be passed with the
// original data to be written and the bytes written as well as any error from the Write-operation.
func NewPostWriteCallbacks(w io.WriteCloser, callbacks ...func([]byte, int, error)) io.WriteCloser {
	return &customWriteCloser{
		closer: w.Close,
		writer: func(p []byte) (n int, err error) {
			n, err = w.Write(p)

			for _, cb := range callbacks {
				cb(p, n, err)
			}

			return n, err
		},
	}
}

// NewPreCloseCallbacks will be called before the Close() call on the WriteCloser.
// Eny error by a callback weill abort the call to the w.Close().
func NewPreCloseCallbacks(w io.WriteCloser, callbacks ...func() error) io.WriteCloser {
	return &customWriteCloser{
		closer: func() error {
			for _, cb := range callbacks {
				if err := cb(); err != nil {
					return err
				}
			}
			return w.Close()
		},
		writer: w.Write,
	}
}

// NewPostCloseCallbacks will be called directly after a Close() call on the WriteCloser have finished.
// any error will be forwards to the callbacks which are invoked in blocking sequence
func NewPostCloseCallbacks(w io.WriteCloser, callbacks ...func(error)) io.WriteCloser {
	return &customWriteCloser{
		closer: func() error {
			err := w.Close()
			for _, cb := range callbacks {
				cb(err)
			}
			return err
		},
		writer: w.Write,
	}
}
