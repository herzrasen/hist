bin = hist

all: init test build

clean:
	rm -f $(bin)
	rm -f coverage.out

init:
	go get

test:
	go test -coverprofile=coverage.out -v ./...
	go tool cover -func coverage.out

build:
	go build -o $(bin)

install:
	cp $(bin) /usr/local/bin/