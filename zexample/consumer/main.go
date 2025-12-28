package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
		//os.Exit(1)
	}
}

func run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	// Create or bind to a durable consumer
	sub, err := js.PullSubscribe(
		"orders.created",
		"ORDER_PROCESSOR",
		nats.PullMaxWaiting(128),
	)
	if err != nil {
		return err
	}

	log.Println("waiting for messages ....")

	for {
		msgs, err := sub.Fetch(5, nats.MaxWait(10*time.Second))
		if err != nil {
			continue
		}

		for _, msg := range msgs {
			log.Printf("Processing message: %s", string(msg.Data))

			//simulating processing logic
			if err := processOrder(msg.Data); err != nil {
				msg.Nak() //retry later
				continue
			}
			msg.Ack() //successfull processing
		}
	}
}

func processOrder(data []byte) error {
	fmt.Println(string(data))
	return nil
}
