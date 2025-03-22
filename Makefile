generate:
	go generate ./...

fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	mkdir -p out
	go build -o ./out/ordaa
