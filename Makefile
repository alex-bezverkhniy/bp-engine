BINARY_NAME=bp-engine.out

all: build

build: test
	go build -o ${BINARY_NAME} engine.go

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go test ./pkg/... ./internal/...