package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"nhooyr.io/websocket"
)

//go:generate mockery --name conn
type conn interface {
	Read(context.Context) (websocket.MessageType, []byte, error)
	Write(context.Context, websocket.MessageType, []byte) error
	Close(websocket.StatusCode, string) error
}

type routine struct {
	ctx        context.Context
	conn       conn
	datawriter io.WriteCloser
	interval   time.Duration
	id         int
}

func newRoutine(ctx context.Context, url string, id int) (*routine, error) {
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		if c != nil {
			c.Close(websocket.StatusAbnormalClosure, "Something happened...")
		}
		return nil, fmt.Errorf("can't dial ws #%d: %w", id, err)
	}

	datawriter, err := newFileWriter(id)
	if err != nil {
		return nil, fmt.Errorf("cannot create file writer: %w", err)
	}

	return &routine{
		ctx:        ctx,
		conn:       c,
		datawriter: datawriter,
		id:         id,
		interval:   500 * time.Millisecond,
	}, nil
}

func (r *routine) start(messages int) {
	for i := 0; i < messages; i++ {
		select {
		case <-r.ctx.Done():
			fmt.Printf("Canceling routine #%d\n", r.id)
			return
		default:
			msg := "Hello " + fmt.Sprint(i) + " from goroutine #" + fmt.Sprint(r.id)
			err := r.conn.Write(r.ctx, websocket.MessageText, []byte(msg))
			if err != nil {
				fmt.Println("Can't send message: ", msg, " with error: ", err)
				return
			}

			if _, msg, err := r.conn.Read(r.ctx); err == nil {
				_, _ = r.datawriter.Write(msg)
			}

			time.Sleep(r.interval)
		}
	}
}

func (r *routine) close() {
	r.datawriter.Close()
	r.conn.Close(websocket.StatusNormalClosure, "Done!")
}
