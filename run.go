package main

import (
	"encoding/json"
	"log"
	"os"
)

type RunCommand struct {
	RunId         int
	Name          string
	Configuration string
	Flags         string
}

func runBenchmark(msg []byte) {
	log.Printf("Received a message: %s", msg)

	runCommand := parseMsg(msg)
	writeConfigFile(runCommand)

	startBenchmark(runCommand)
}

func startBenchmark(runCommand RunCommand) {
	messages := exectuteComand("./pkb.py", "--openstack_network=bechmarks", "--ignore_package_requirements", "--benchmark_config_file="+runCommand.Name, runCommand.Flags)
	log.Printf("Started Running")
	sendStarted(runCommand)
	message := <-messages
	if message == "failed" {
		sendFailed(runCommand)
		moveFailed(runCommand)
	} else {
		collectResults(runCommand)
		sendFinished(runCommand)
		log.Printf("End Running Output:")
		log.Printf(message)
	}
}

func writeConfigFile(runCommand RunCommand) {
	f, err := os.Create(runCommand.Name)
	failOnError(err, "Couldn't read Config File")
	f.Write([]byte(runCommand.Configuration))
}

func parseMsg(msg []byte) RunCommand {
	runCommand := RunCommand{}
	err := json.Unmarshal(msg, &runCommand)

	if err != nil {
		failOnError(err, "Couldn't read Config File")
	}

	return runCommand
}
