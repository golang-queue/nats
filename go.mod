module github.com/golang-queue/nats

go 1.18

require (
	github.com/golang-queue/queue v0.1.4-0.20221230133718-0314ef173f98
	github.com/nats-io/nats.go v1.31.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/goleak v1.2.1
)

require (
	braces.dev/errtrace v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/nats-io/nkeys v0.4.6 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/golang-queue/queue v0.1.4-0.20221230133718-0314ef173f98 => ../queue
