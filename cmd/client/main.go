package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if len(os.Args) != 4 {
		fmt.Println("Command line: ./client ip #goroutines #messages")
		os.Exit(1)
	}

	port := os.Args[1]
	routines, _ := strconv.Atoi(os.Args[2])
	messages, _ := strconv.Atoi(os.Args[3])
	url := fmt.Sprintf("ws://localhost:%s", port)

	var wg sync.WaitGroup
	wg.Add(routines)

	for j := 0; j < routines; j++ {
		go func(r int) {
			c, _, err := websocket.Dial(ctx, url, nil)
			if err != nil {
				fmt.Println("Can't dial...")
				c.Close(websocket.StatusAbnormalClosure, "Something happened...")
				return
			}
			defer c.Close(websocket.StatusNormalClosure, "Done!")
			defer wg.Done()

			for i := 0; i < messages; i++ {
				msg := "Hello " + fmt.Sprint(i) + " from goroutine #" + fmt.Sprint(r)
				err = c.Write(ctx, websocket.MessageText, []byte(msg))
				if err != nil {
					fmt.Println("Can't send message: ", msg, " with error: ", err)
					return
				}

				if typ, msg, err := c.Read(ctx); err == nil {
					fmt.Printf("[%v] %s\n", typ, string(msg))
				}

				time.Sleep(500 * time.Millisecond)
			}
		}(j)
	}

	wg.Wait()
}
