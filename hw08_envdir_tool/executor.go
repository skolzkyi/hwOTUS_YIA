package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var ErrBadCmd = errors.New("bad cmd data")

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var err error

	if len(cmd) < 2 {
		err = ErrBadCmd
		fmt.Println(err.Error())
		return 13
	}

	name := cmd[0]
	arguments := cmd[1:]

	comm := exec.Command(name, arguments...)

	comm.Stdin = os.Stdin
	comm.Stdout = os.Stdout
	comm.Stderr = os.Stderr

	for tag, curEnv := range env {
		_, ok := os.LookupEnv(tag)
		if ok {
			os.Unsetenv(tag)
		}
		if !curEnv.NeedRemove {
			os.Setenv(tag, curEnv.Value)
		}
	}

	comm.Env = os.Environ()

	err = comm.Run()
	if err != nil {
		returnCode = comm.ProcessState.ExitCode()
	}

	return
}
