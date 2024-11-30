package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		log.Println("cmd args is empty")
		return 1
	}

	for eVar, eValue := range env {
		if eValue == "" {
			err := os.Unsetenv(eVar)
			if err != nil {
				fmt.Println(err)
			}

			continue
		}

		err := os.Setenv(eVar, eValue)
		if err != nil {
			fmt.Println(err)
		}
	}

	cmdExec := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout

	err := cmdExec.Run()
	if err != nil {
		fmt.Println(err)
	}

	return cmdExec.ProcessState.ExitCode()
}
