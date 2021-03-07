package eio

import (
	"fmt"
	"io"
)

// NewSpanWriter returns a writer which splits writes into multiple blobs.
// fileIDPattern can be any Sprintf pattern that accepts an integer argument and yeild a valid filename.
// maxBytes is the largest a file can be before blocking the write and creating a new file.
// writerfactory is the method used to intialize a new writer
// Write a segment that is larger than maxBytes will always fail.
func NewSpanWriter(
	fileIDPattern string,
	maxBytes int,
	wf func(fileId string) (io.WriteCloser, error),
) io.WriteCloser {
	fileSequence := 0

	var currentWriter io.WriteCloser
	var lastErr error

	reInitWriter := func() {
		currentWriter, lastErr = wf(fmt.Sprintf(fileIDPattern, fileSequence))
		if lastErr == nil {
			currentWriter = NewLimitedWriter(currentWriter, WithMaxBytes(maxBytes))
		}
		fileSequence++
	}

	reInitWriter()

	return &CustomWriteCloser{
		Closer: func() error {
			if err := currentWriter.Close(); err != nil {
				return err
			}
			// future invokations will get this error
			lastErr = ErrAlreadyClosed
			return nil
		},
		Writer: func(p []byte) (n int, err error) {
			if currentWriter == nil {
				reInitWriter()
			}
			if lastErr != nil {
				return 0, lastErr
			}

			n, err = currentWriter.Write(p)
			if err == ErrTooLargeWrite {
				// If the write would never fit; return ErrTooLargeWrite
				if len(p) > maxBytes {
					return n, err
				}
				reInitWriter()
				if lastErr != nil {
					return 0, lastErr
				}
				n, lastErr = currentWriter.Write(p)
			}

			return n, lastErr
		},
	}
}
