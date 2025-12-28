package main

import (
	"log"

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

	// Create a stream (idempotent)
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"orders.created"},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return err
	}

	//Publish Message
	msg := []byte(`{"order_id":"123","amount":99.99}`)
	_, err = js.Publish("orders.created", msg)
	if err != nil {
		return err
	}

	log.Println("Order event published")
	return nil
}
