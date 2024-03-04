BINARY_NAME=Saudode

SRC_DIR=./internal

all: build

build:
	go build -o ${BINARY_NAME} ${SRC_DIR}

run:
	go run ${SRC_DIR}

test:
	go test -v ${SRC_DIR}

clean:
	go clean
	rm -f ${BINARY_NAME}

deps:
	go mod download

