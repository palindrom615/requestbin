all: requestbin

requestbin:
	GOARCH=amd64 GOOS=linux go build -o build/requestbin ./cmd/requestbin
	zip -j build/requestbin.zip build/requestbin

clean:
	go clean
	rm -rf ./build
