package main

import (
	"log"
)

func main() {
	config := readConfig()
	log.Printf("Starting Benchmark Runner")
	conn:= getConnection(config)
	log.Printf("Connected to Messaging Queue")
	defer conn.Close()

	log.Printf("Connecting to Run Queue")
	runQueue(conn, config)
	
}
