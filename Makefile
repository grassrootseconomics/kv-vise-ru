BIN := kv-vise-ru-service
BUILD_CONF := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
BUILD_COMMIT := $(shell git rev-parse --short HEAD 2> /dev/null)
DEBUG := DEV=true

.PHONY: build run clean

clean:
	rm ${BIN}

build:
	${BUILD_CONF} go build -ldflags="-X main.build=${BUILD_COMMIT} -s -w" -o build/${BIN} cmd/service/main.go

run:
	${BUILD_CONF} ${DEBUG} go run cmd/service/main.go