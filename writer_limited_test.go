package eio_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zatte/eio"
)

func TestTimedSelfDestructAfterMaxBytes(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	writeCloser := eio.NewLimitedWriter(eio.NewWriteMaybeCloser(buffer), eio.WithMaxBytes(32))

	// Valid withing max limits
	testString1 := []byte("This write should not be blocked")
	n, err := writeCloser.Write(testString1)
	require.NoError(t, err)

	assert.Equal(t, len(testString1), n, "string bytes should have been reported written")
	assert.Equal(t, len(testString1), len(buffer.Bytes()), "string bytes should have been actually written")

	// Blocked
	testString2 := []byte("This write should be blocked")
	n, err = writeCloser.Write(testString2)
	assert.EqualError(t, eio.ErrTooLargeWrite, err.Error())

	assert.Equal(t, 0, n, "no bytes should have been reported written")
	assert.Equal(t, len(testString1), len(buffer.Bytes()), "only existing writes bytes should have been written")

}
