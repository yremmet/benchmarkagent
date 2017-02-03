package main

import (
	"log"
  "encoding/json"
  "os"
  "fmt"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Configuration struct {
		RunnerId 	   int
    AmqpUser     string
    AmqpPassword string
    AmqpHost 		 string
    AmqpPort     int
		RunQueue 		 string
		ResponseTopic string
		MongoHost		 string
		Database		 string
		Collection   string
		ResultFolder string
		AuthDatabase string
  	AuthUserName string
  	AuthPassword string
}

func readConfig() (Configuration) {
  file, _ := os.Open("conf.json")
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  err := decoder.Decode(&configuration)
  if err != nil {
    failOnError(err, "Couldn't read Config File")
  }
  return configuration;
}

func getConnection(config Configuration) (*amqp.Connection) {
	connection := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.AmqpUser, config.AmqpPassword, config.AmqpHost, config.AmqpPort)
	conn, err := amqp.Dial(connection)
	failOnError(err, "Failed to connect to Queue")
	return conn
}
