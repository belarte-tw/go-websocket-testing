build:
	go build -o bin ./...

run-server:
	go run cmd/echo/main.go

run-client:
	go run cmd/client/main.go 1323