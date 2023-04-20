package echoing

import (
	"context"
	"fmt"
	"io"
	"time"

	"nhooyr.io/websocket"
)

//go:generate mockery --name conn
type conn interface {
	Reader(context.Context) (websocket.MessageType, io.Reader, error)
	Writer(context.Context, websocket.MessageType) (io.WriteCloser, error)
}

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func Echo(ctx context.Context, c conn) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}
