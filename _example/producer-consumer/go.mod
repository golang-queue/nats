module example

go 1.18

require (
	github.com/appleboy/graceful v0.0.4
	github.com/golang-queue/nats v0.0.3-0.20210907015837-3e2e4b448b3d
	github.com/golang-queue/queue v0.1.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.2.0 // indirect
	github.com/nats-io/nats.go v1.13.1-0.20220308171302-2f2f6968e98d // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	golang.org/x/time v0.0.0-20220224211638-0e9765cccd65 // indirect
)

replace github.com/golang-queue/nats => ../../
