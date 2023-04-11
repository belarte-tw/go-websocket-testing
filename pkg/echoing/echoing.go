package echoing

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func Echo(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}
	log.Printf("Received message: %v", typ)

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	n, err := io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}
	log.Printf("Read %d bytes", n)

	err = w.Close()
	return err
}
