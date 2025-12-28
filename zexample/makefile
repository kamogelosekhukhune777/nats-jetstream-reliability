run-consumer:
	go run consumer/main.go

run-producer:
	go run producer/main.go

run-all:
	go run consumer/main.go &
	go run producer/main.go

nats-server-run:
	docker run -p 4222:4222 nats:latest -js

run:
	nats-server-run &
	run-all
