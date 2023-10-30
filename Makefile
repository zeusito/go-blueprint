BINARY_NAME=binary

test:
	go test -v ./...

lint:
	golangci-lint run --config=.golangci.yaml

build:
	CGO_ENABLED=0 go build -o ./out/${BINARY_NAME} ./cmd/main.go

run:
	./out/${BINARY_NAME}

clean:
	go clean
	rm -rf ./out
