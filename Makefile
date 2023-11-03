BINARY_NAME=requestbin

all:
	GOARCH=amd64 GOOS=linux go build -o build/$(BINARY_NAME) cmd/main.go
	zip -j build/$(BINARY_NAME).zip build/$(BINARY_NAME)

clean:
	go clean
	rm -rf ./build