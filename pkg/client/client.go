package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

func Run(ctx context.Context, url string, routines, messages int) {
	var wg sync.WaitGroup
	wg.Add(routines)

	for j := 0; j < routines; j++ {
		time.Sleep(10 * time.Millisecond)
		go func(r int) {
			defer wg.Done()
			startRoutine(ctx, url, r, messages)
		}(j)
	}

	wg.Wait()
}

func startRoutine(ctx context.Context, url string, routine, messages int) {
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		fmt.Printf("Can't dial ws #%d, err: %s\n", routine, err)
		if c != nil {
			c.Close(websocket.StatusAbnormalClosure, "Something happened...")
		}
		return
	}
	defer c.Close(websocket.StatusNormalClosure, "Done!")

	filename := fmt.Sprintf("test/file%d.txt", routine)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("failed creating file: %s", err)
		return
	}
	defer file.Close()

	datawriter := bufio.NewWriter(file)
	defer datawriter.Flush()

	for i := 0; i < messages; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("Canceling routine #%d\n", routine)
			return
		default:
			msg := "Hello " + fmt.Sprint(i) + " from goroutine #" + fmt.Sprint(routine)
			err = c.Write(ctx, websocket.MessageText, []byte(msg))
			if err != nil {
				fmt.Println("Can't send message: ", msg, " with error: ", err)
				return
			}

			if _, msg, err := c.Read(ctx); err == nil {
				_, _ = datawriter.WriteString(string(msg) + "\n")
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}
