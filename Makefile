clean:
	rm -rf bin/

build: clean
	mkdir bin
	go build -o bin ./...

run-server:
	go run cmd/echo/main.go

run-client:
	go run cmd/client/main.go 1323 3 5

build-docker:
	docker build -t go-websocket-server .

run-docker:
	docker run -p 1323:1323 go-websocket-server

docker: build-docker run-docker