all: build

test:
	go test ./...

build:
	go build ./cmd/go-traphandle

run: build
	./go-traphandle -config _files/config.yml -server 127.0.0.1:9999
