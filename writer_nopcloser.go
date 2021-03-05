package eio

import "io"

// NewWriteCustomCloser taggs on a callback-close method on writers which lacks them (such as buffer)
func NewWriteCustomCloser(w io.Writer, f func() error) io.WriteCloser {
	if f == nil {
		f = func() error { return nil }
	}

	return &customWriteCloser{
		writer: w.Write,
		closer: f,
	}
}

// NewWriteMaybeCloser returns a WriteCloser; IF the writer implements the Close() method
// it will be invoked, otherwise the Close() will be a NOP
func NewWriteMaybeCloser(w io.Writer) io.WriteCloser {
	return NewWriteCustomCloser(w, func() error {
		if closer, ok := w.(io.Closer); ok {
			closer.Close()
		}
		return nil
	})
}
