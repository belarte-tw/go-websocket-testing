# Build stage
FROM golang:1.20.3-alpine3.17 AS BuildStage

WORKDIR /work
COPY . .
RUN go mod download
RUN go build -o /bin ./...

# Deploy stage
FROM alpine:latest

WORKDIR /
COPY --from=BuildStage /bin/echo /server

EXPOSE 1323

RUN addgroup -g 1000 -S username && \
    adduser -u 1000 -S username -G username
USER username

ENTRYPOINT ["/server"]