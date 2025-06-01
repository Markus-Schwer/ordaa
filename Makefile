.PHONY: generate
generate:
	go generate ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	go tool golangci-lint run

.PHONY: lint/fix
lint/fix:
	go tool golangci-lint run --fix

.PHONY: test
test:
	go test -count 1 -v ./...

.PHONY: build
build:
	mkdir -p out
	go build -o ./out/ordaa
