package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	port := os.Args[1]
	url := fmt.Sprintf("ws://localhost:%s", port)

	var wg sync.WaitGroup
	wg.Add(3)

	for j := 0; j < 3; j++ {
		go func(r int) {
			c, _, err := websocket.Dial(ctx, url, nil)
			if err != nil {
				fmt.Println("Can't dial...")
				c.Close(websocket.StatusAbnormalClosure, "Something happened...")
				return
			}
			defer c.Close(websocket.StatusNormalClosure, "Done!")
			defer wg.Done()

			for i := 0; i < 5; i++ {
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
