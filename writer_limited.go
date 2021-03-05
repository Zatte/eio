package eio

import (
	"fmt"
	"io"
)

// NewLimitedWriter will automatically Close the writer on certain conditions defined by options.
// Since some options include async events the writes will be mutex protected meaning only
// 1 thread will be able to call Write/Close at the time
func NewLimitedWriter(w io.WriteCloser, options ...LimitedWriterOption) io.WriteCloser {
	for _, o := range options {
		w = o(w)
	}

	return NewSyncedWriteCloser(w)
}

// LimitedWriterOption modifies writers with self destruct behaviour; calling close on specific conditions.
type LimitedWriterOption func(w io.WriteCloser) io.WriteCloser

// WithMaxBytes will block writes which would make the total stream larger than maxBytes.
func WithMaxBytes(maxBytes int) LimitedWriterOption {
	bytesWritten := 0
	return func(w io.WriteCloser) io.WriteCloser {
		preCheck := NewPreWriteCallbacks(w, func(p []byte) error {
			if bytesWritten+len(p) > maxBytes {
				if err := w.Close(); err != nil {
					return fmt.Errorf("failed to close WriteCloser writing maxBytes; Close error was: %w", err)
				}
				return ErrTooLargeWrite
			}
			return nil
		})

		return NewPostWriteCallbacks(preCheck, func(p []byte, n int, err error) {
			bytesWritten += n
		})
	}
}
