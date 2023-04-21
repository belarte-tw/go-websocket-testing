package client

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func Run(ctx context.Context, url string, routines, messages int) {
	var wg sync.WaitGroup
	wg.Add(routines)

	for j := 0; j < routines; j++ {
		time.Sleep(10 * time.Millisecond)

		go func(r int) {
			defer wg.Done()
			routine, err := newRoutine(ctx, url, r)
			if err != nil {
				fmt.Printf("Failed to create routine #%d: %s", r, err.Error())
				return
			}
			routine.start(messages)
			routine.close()
		}(j)
	}

	wg.Wait()
}
