module example

go 1.22

require (
	github.com/appleboy/graceful v0.0.4
	github.com/golang-queue/nats v0.0.3-0.20210907015837-3e2e4b448b3d
	github.com/golang-queue/queue v0.2.1
)

require (
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/nats-io/nats.go v1.38.0 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/golang-queue/nats => ../../
