package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	js, _ := nc.JetStream()

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

		_, err := js.PublishMsg(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Events published")
}
