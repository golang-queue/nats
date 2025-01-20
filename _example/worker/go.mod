module example

go 1.22

require (
	github.com/golang-queue/nats v0.0.2-0.20210822122542-200fdcf19ebf
	github.com/golang-queue/queue v0.3.0
)

require (
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/nats-io/nats.go v1.38.0 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/golang-queue/nats => ../../
