run-consumer:
	go run cmd/consumer/main.go

run-producer:
	go run cmd/producer/main.go

nats-server-run:
	docker run -p 4222:4222 nats:latest -js
