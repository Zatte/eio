package eio

import "errors"

// ErrAlreadyClosed is returned when any operation id done on a already closed reader/writer
var ErrAlreadyClosed = errors.New("stream is closed")

// ErrTooLargeWrite is returend by a Writer that has been protected with NewLimitedWriter(w, WithMaxBytes())
var ErrTooLargeWrite = errors.New("the write would violate maxBytes contraint")
