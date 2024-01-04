BINARY_NAME=example-bp-engine.out

all: build

build:
	go build -o ${BINARY_NAME} main.go

dev: export CONFIG_FILE=./config.json
dev: export DB_URL=/home/alex/Workspace/go/bp-engine/example/bp-engine.db
dev:
	go run main.go


watch:
	nodemon --watch './**/*' --ext 'go,html,js,css' --signal SIGTERM --exec 'make dev'

clean:
	go clean
	rm ${BINARY_NAME}