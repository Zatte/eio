package eio

import (
	"io"
)

// readCloserImpl ads a NOP Close method to any Reader.
type readCloserImpl struct {
	io.Reader
	closer func() error
}

// Close is a NOP if not specified by a user.
func (nc readCloserImpl) Close() error {
	if nc.closer != nil {
		return nc.closer()
	}
	return nil
}

// NewReadCustomCloser returns a ReadCloser; will decorate the Reader with the a user provider closer
func NewReadCustomCloser(r io.Reader, closer func() error) io.ReadCloser {
	return &readCloserImpl{
		Reader: r,
		closer: closer,
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
