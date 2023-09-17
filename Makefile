BINARY_NAME=task
.DEFAULT_GOAL := run

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go test ./...

testv:
	go test -v ./...

