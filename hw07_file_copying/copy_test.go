package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Unsupported", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("Empty", func(t *testing.T) {
		err := Copy("./testdata/empty.txt", "/tmp/", 0, 0)
		require.Equal(t, ErrFromFileIsEmpty, err)
	})

	t.Run("Non-existent", func(t *testing.T) {
		err := Copy("./non-existent-file", "./any-dest", 0, 0)
		require.Error(t, err)
	})

	t.Run("Offset exceeded", func(t *testing.T) {
		tempFile, err := os.CreateTemp(".", "out_")
		if err != nil {
			t.FailNow()
		}

		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		err = Copy("./testdata/input.txt", tempFile.Name(), 10000, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})
}
