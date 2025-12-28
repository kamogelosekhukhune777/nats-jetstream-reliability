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
	}
}

func run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS server: %w", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Primary stream
	js.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",
		Subjects: []string{"events.process"},
		Storage:  nats.FileStorage,
		// Enables exactly-once semantics via message deduplication
		// based on Msg-Id headers
		Duplicates: time.Minute,
	})

	// DLQ stream
	js.AddStream(&nats.StreamConfig{
		Name:     "DLQ",
		Subjects: []string{"events.dlq"},
		Storage:  nats.FileStorage,
	})

	for i := 0; i < 10; i++ {
		msg := nats.NewMsg("events.process")

		// 6. Exactly-once semantics via deduplication headers
		msg.Header.Set(nats.MsgIdHdr, "event-id-123") // stable idempotency key

		msg.Data = []byte(`{"event":"order_created"}`)

		// 7. Publish message with deduplication headers
		_, err := js.PublishMsg(msg)
		if err != nil {
			return err
		}

	}

	log.Println("Events published")
	return nil
}
