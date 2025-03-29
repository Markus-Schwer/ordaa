generate:
	go generate ./...

fmt:
	go fmt ./...

lint:
	go tool golangci-lint run

lint-fix:
	go tool golangci-lint run --fix

test:
	go test -v ./...

build:
	mkdir -p out
	go build -o ./out/ordaa
