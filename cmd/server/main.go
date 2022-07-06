package main

import (
	"fmt"
	"log"

	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "client-123", stan.NatsURL("nats://127.0.0.1:4222"))
	if err != nil {
		log.Fatalln(err)
	}

	// Subscribe with durable name
	_, err = sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	}, stan.DurableName("my-durable"))
	if err != nil {
		log.Fatalln(err)
	}
	for {

	}
}
