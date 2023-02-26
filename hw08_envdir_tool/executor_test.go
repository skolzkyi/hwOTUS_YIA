package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// positive cases
	t.Run("PositiveSet", func(t *testing.T) {
		name := "NEWENV1"
		value := "NEWVALUE1"

		_, ok := os.LookupEnv(name)
		require.False(t, ok)

		envir := make(Environment)
		envir[name] = EnvValue{
			Value:      value,
			NeedRemove: false,
		}

		exitCode := RunCmd([]string{"echo", "test"}, envir)
		require.Equal(t, 0, exitCode)

		curValue, ok := os.LookupEnv(name)
		require.True(t, ok)
		require.Equal(t, curValue, value)
	})
	t.Run("PositiveUnset", func(t *testing.T) {
		name := "NEWENV2"
		value := "NEWVALUE2"
		os.Setenv(name, value)
		require.Equal(t, value, os.Getenv(name))

		envir := make(Environment)
		envir[name] = EnvValue{
			Value:      value,
			NeedRemove: true,
		}

		exitCode := RunCmd([]string{"echo", "test"}, envir)
		require.Equal(t, 0, exitCode)

		_, ok := os.LookupEnv(name)
		require.False(t, ok)
	})

	t.Run("PositiveExecCommAmdArgs", func(t *testing.T) {
		out := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		exitCode := RunCmd([]string{"echo", "Command test"}, nil)
		require.Equal(t, 0, exitCode)

		require.NoError(t, w.Close())
		os.Stdout = out

		var buffer bytes.Buffer
		_, err := io.Copy(&buffer, r)
		require.NoError(t, err)

		require.Equal(t, "Command test\n", buffer.String())
	})
	// negative cases
	t.Run("NegativeBadInputCmd", func(t *testing.T) {
		exitCode := RunCmd([]string{}, nil)
		require.Equal(t, 13, exitCode)
	})
}
