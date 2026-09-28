module example

go 1.26.8

require (
	github.com/appleboy/graceful v1.3.0
	github.com/golang-queue/nats v0.2.0
	github.com/golang-queue/queue v0.5.0
)

require (
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.18.7 // indirect
	github.com/nats-io/nats.go v1.52.0 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.56.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)

replace github.com/golang-queue/nats => ../../
