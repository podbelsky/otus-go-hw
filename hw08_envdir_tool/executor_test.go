package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("good cmd", func(t *testing.T) {
		rc := RunCmd([]string{"bash", "-c", "exit 7"}, Environment{})

		require.Equal(t, 7, rc)
	})

	t.Run("bad cmd", func(t *testing.T) {
		code := RunCmd([]string{}, Environment{})
		require.Equalf(t, 1, code, "should be error code 1")
	})
}
