package main

import (
	"os/exec"
)

func exectuteComand(cmdName string, params ...string) <-chan string {
	var (
		cmdOut []byte
		err    error
	)
	messages := make(chan string, 500)
	go func() {
		exec := exec.Command(cmdName, params...)
		if cmdOut, err = exec.Output(); err != nil {
	    //failOnError(err, "Error executing command "+ cmdName)
			messages <- "failed"
			close(messages)
		} else {
			messages <- string(cmdOut)
			close(messages)
		}
	}()

	return messages
}
