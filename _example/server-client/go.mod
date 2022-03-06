module example

go 1.16

require (
	github.com/golang-queue/nats v0.0.3-0.20210907015837-3e2e4b448b3d // indirect
	github.com/golang-queue/queue v0.0.11
	github.com/golang-queue/redisdb v0.0.3-0.20210905082752-77eb8e774b71
)

replace github.com/golang-queue/nats => ../../
