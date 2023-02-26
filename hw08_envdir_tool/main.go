package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	envDirPath := args[1]
	commandAndArgs := args[2:]

	env, err := ReadDir(envDirPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(13)
	}

	exitCode := RunCmd(commandAndArgs, env)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
