module example

go 1.26.10

require (
	github.com/appleboy/graceful v1.3.0
	github.com/golang-queue/nats v0.3.0
	github.com/golang-queue/queue v0.6.2
)

require (
	github.com/jpillora/backoff v1.1.0 // indirect
	github.com/klauspost/compress v1.20.4 // indirect
	github.com/nats-io/nats.go v1.55.0 // indirect
	github.com/nats-io/nkeys v0.4.17 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.59.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
)

replace github.com/golang-queue/nats => ../../
