# eio - Extended IO
Common additions i find myself rewriting from time to time. Highly related to functionality found in the [io package](https://golang.org/pkg/io). No external deps except for std lib.

[![Go Report Card](https://goreportcard.com/badge/github.com/zatte/eio?style=flat-square)](https://goreportcard.com/report/github.com/zatte/eio)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/zatte/eio)

## License
Apache 2.0

## Forked
Forked and rewritten & cleaned version of [github.com/kvanticoss/goutils/eioutil](https://github.com/kvanticoss/goutils/tree/master/writerfactory)
(remove most external deps and some cleanup; since ioutil is becoming deprecated no point in having extended ioutil)

##
Utilities

```golang

// Custom Close methods
eio.NewReadCustomCloser(r io.Reader, closer func() error) io.ReadCloser
eio.NewWriteCustomCloser(w io.Writer, closer func() error) io.WriteCloser

// If a Close methods exists on the reader use it; other wise tag add a nop Close(); either way, you now have a WriterCloser
eio.NewReadMaybeCloser(r io.Reader, closer func() error) io.ReadCloser
eio.NewWriteMaybeCloser(w io.Writer, closer func() error) io.WriteCloser

// Writer hooks
eio.NewPreWriteCallbacks(w io.WriteCloser, callbacks ...func([]byte) error) io.WriteCloser // errors abort a write
eio.NewPostWriteCallbacks(w io.WriteCloser, callbacks ...func([]byte, int, error)) io.WriteCloser
eio.NewPreCloseCallbacks(w io.WriteCloser, callbacks ...func() error) io.WriteCloser // errors abort a close
eio.NewPostCloseCallbacks(w io.WriteCloser, callbacks ...func(error)) io.WriteCloser

// Limited Writer
eio.NewLimitedWriter(w io.WriteCloser, options ...LimitedWriterOption) io.WriteCloser

// SpanWriter
eio.NewSpanWriter(format string, maxBytes int, writerFactory func(blobId string) io.WriteCloser)


// Limited Writer example
buffer := bytes.NewBuffer(nil)
eio.NewLimitedWriter(eio.NewWriteMaybeCloser(buffer), eio.WithMaxBytes(32)) // Will return eio.ErrTooLargeWrite if more than 32 would be written.


// SpanWriter Example
db := map[string]*bytes.Buffer{}
writerFactory := func(id string) (io.WriteCloser, error) {
  db[id] = bytes.NewBuffer(nil)
  return eio.NewWriteMaybeCloser(db[id]), nil
}
spanWriter := eio.NewSpanWriter("%08d.data", 1024*1024, writerFactory)

// Writes to span writer will start filling up 00000000.data up to 1MB
// once maxBytes is reach the previous writer will be closed and a new on
// created with a new sequence number (00000001.data).
spanWriter.Write(....)
```