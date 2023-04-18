package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"go-socket/pkg/client"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	if len(os.Args) != 4 {
		fmt.Println("Command line: ./client ip #goroutines #messages")
		os.Exit(1)
	}

	port := os.Args[1]
	routines, _ := strconv.Atoi(os.Args[2])
	messages, _ := strconv.Atoi(os.Args[3])
	url := fmt.Sprintf("ws://localhost:%s", port)

	client.Run(ctx, url, routines, messages)
}
