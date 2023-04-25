package echoing

import (
	"context"
	"fmt"
	"time"

	"nhooyr.io/websocket"
)

//go:generate mockery --name conn
type conn interface {
	Read(context.Context) (websocket.MessageType, []byte, error)
	Write(context.Context, websocket.MessageType, []byte) error
}

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func Echo(ctx context.Context, c conn) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	typ, msg, err := c.Read(ctx)
	if err != nil {
		return fmt.Errorf("reading from ws: %w", err)
	}

	err = c.Write(ctx, typ, msg)
	if err != nil {
		return fmt.Errorf("writing to ws: %w", err)
	}

	return nil
}
