package main

import (
  "github.com/streadway/amqp"
  "encoding/json"
)

type Message struct {
    RunnerId   int
    RunId      int
    Status 		 string
}

func send(message Message){
  config:= readConfig()
  conn:=   getConnection(config)

  message.RunnerId = config.RunnerId
  ch, err := conn.Channel()
	failOnError(err, "Failed to open a response-channel")
	defer ch.Close()

  json,err := json.Marshal(message)
  failOnError(err, "Failed to convert to JSON")

  err = ch.Publish(
    "amq.topic",          // exchange
    config.ResponseTopic, // routing key
    false,  // mandatory
    false,  // immediate
    amqp.Publishing {
      ContentType: "text/json",
      Body:        json,
    })

  failOnError(err, "Failed to publish a message")

  ch.Close()
}

func sendStarted(runCommand RunCommand){
   message:= Message{}
   message.RunId = runCommand.RunId
   message.Status = "Started"
   send(message)
}

func sendFinished(runCommand RunCommand){
   message:= Message{}
   message.RunId = runCommand.RunId
   message.Status = "Finished"
   send(message)
}

func sendFailed(runCommand RunCommand){
   message:= Message{}
   message.RunId = runCommand.RunId
   message.Status = "Failed"
   send(message)
}
