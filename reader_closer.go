package eio

import (
	"io"
)

// CustomReadCloser implements custom methods for Read and Close calls
type CustomReadCloser struct {
	Closer func() error
	Reader func(p []byte) (n int, err error)
}

// Close is a NOP if not specified by a user.
func (nc CustomReadCloser) Close() error {
	if nc.Closer != nil {
		return nc.Closer()
	}
	return nil
}

func (nc CustomReadCloser) Read(p []byte) (n int, err error) {
	return nc.Reader(p)
}

// NewReadCustomCloser returns a ReadCloser; will decorate the Reader with the a user provider closer
func NewReadCustomCloser(r io.Reader, closer func() error) io.ReadCloser {
	return &CustomReadCloser{
		Reader: r.Read,
		Closer: closer,
	}
}

// NewReadMaybeCloser returns a ReadCloser where IF the reader implements the Close() method
// it will be invoked, otherwise the Close() will be a NOP
func NewReadMaybeCloser(r io.Reader) io.ReadCloser {
	return NewReadCustomCloser(r, func() error {
		if closer, ok := r.(io.Closer); ok {
			closer.Close()
		}
		return nil
	})
}
