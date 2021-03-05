package eio_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zatte/eio"
)

func TestWriterSpan(t *testing.T) {
	db := map[string]*bytes.Buffer{}

	writerFactory := func(id string) (io.WriteCloser, error) {
		db[id] = bytes.NewBuffer(nil)
		return eio.NewWriteMaybeCloser(db[id]), nil
	}

	spanWriter := eio.NewSpanWriter("%08d.data", 5, writerFactory)

	_, err := spanWriter.Write([]byte("1234")) // goes into 1st
	require.NoError(t, err)
	_, err = spanWriter.Write([]byte("56")) // goes into 2nd
	require.NoError(t, err)
	_, err = spanWriter.Write([]byte("56")) // goes into 2nd
	require.NoError(t, err)
	_, err = spanWriter.Write([]byte("78961")) // goes into 3rd
	require.NoError(t, err)

	_, err = spanWriter.Write([]byte("writes > maxBytes will always fail")) // goes into nothing
	require.Error(t, err, eio.ErrTooLargeWrite.Error())

	// Close should always complete
	require.NoError(t, spanWriter.Close())

	assert.EqualValues(t, map[string]*bytes.Buffer{
		"00000000.data": bytes.NewBufferString("1234"),
		"00000001.data": bytes.NewBufferString("5656"),
		"00000002.data": bytes.NewBufferString("78961"),
	}, db)

	_, err = spanWriter.Write([]byte("write to close stream should fail")) // goes into 3rd
	require.Error(t, err, eio.ErrAlreadyClosed.Error())
}
