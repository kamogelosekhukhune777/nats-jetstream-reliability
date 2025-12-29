package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kamogelosekhukhune777/nats-jetstream-reliability/internal/idempotency"
	"github.com/kamogelosekhukhune777/nats-jetstream-reliability/internal/metrics"
	"github.com/kamogelosekhukhune777/nats-jetstream-reliability/internal/workerpool"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	metrics.Register()
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// 3. Consumer max deliveries
	sub, err := js.PullSubscribe(
		"events.process",
		"EVENT_WORKERS",
		nats.MaxDeliver(5),           // retry limit
		nats.AckWait(10*time.Second), // retry delay
	)
	if err != nil {
		return fmt.Errorf("failed to create pull subscription: %w", err)
	}

	store := idempotency.New()
	pool := workerpool.Pool{}

	pool.Run(4, func() {
		for {
			msgs, err := sub.Fetch(1, nats.MaxWait(10*time.Second))
			if err != nil {
				log.Printf("failed to fetch message: %v", err)
				continue
			}
			msg := msgs[0]

			eventID := msg.Header.Get(nats.MsgIdHdr)

			// 4. Idempotency keys
			if store.Seen(eventID) {
				msg.Ack()
				continue
			}

			if err := process(msg.Data); err != nil {
				// 1. Retry with Nak() + delay
				msgDelivered, err := msg.Metadata()
				if err != nil {
					log.Printf("failed to get message metadata: %v", err)
					continue
				}
				if msgDelivered.NumDelivered > 5 {
					// 2. DLQ stream
					js.Publish("events.dlq", msg.Data)
					msg.Ack()
				}
				msg.Nak()
				continue
			}
			msg.Ack()
		}
	})
	pool.Wait()
	return nil
}

func process(data []byte) error {
	log.Println(string(data))
	return nil
}
