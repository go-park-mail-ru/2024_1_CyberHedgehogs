BINARY_NAME=Saudode

SRC_DIR=./cmd

all: build

build:
	go build -o ${BINARY_NAME} ${SRC_DIR}/main.go

run:
	go run ${SRC_DIR}/main.go

build_and_run: build run

test:
	go test -v ${SRC_DIR}/...

clean:
	go clean
	rm -f ${BINARY_NAME}

deps:
	go mod download

