clean:
	rm -rf bin/

build: clean
	mkdir bin
	go build -o bin ./...

test-generate:
	go generate ./...

test-unit:
	go test ./...

test-integration:
	go test ./... --tags=integration

test: test-generate test-unit test-integration

run-server:
	go run cmd/echo/main.go

run-client:
	go run cmd/client/main.go 1323 3 5

docker-build:
	docker build -t go-websocket-server .

docker-run:
	docker run -p 1323:1323 go-websocket-server

docker: docker-build docker-run