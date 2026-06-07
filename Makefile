.PHONY: build run test lint clean

BIN := hariko
SRC := ./cmd/hariko

build:
	go build -o $(BIN) $(SRC)

run: build
	./$(BIN) -f examples/routes.yaml

test:
	go test -v -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BIN) coverage.out coverage.html

coverage: test
	go tool cover -html=coverage.out -o coverage.html
