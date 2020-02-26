.PHONY: clean build

build: format clean
	go build ./...

test: get
	go test -v .

get:
	go get -t -v ./...

format:
	find . -name \*.go -type f -exec gofmt -w {} \;

clean:
