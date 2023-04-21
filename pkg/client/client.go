package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
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

type routine struct {
	ctx        context.Context
	conn       *websocket.Conn
	datawriter io.WriteCloser
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

	return &routine{ctx: ctx, conn: c, datawriter: datawriter, id: id}, nil
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

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (r *routine) close() {
	r.datawriter.Close()
	r.conn.Close(websocket.StatusNormalClosure, "Done!")
}

type fileWriter struct {
	dataWriter *bufio.Writer
	file       *os.File
}

func newFileWriter(routine int) (io.WriteCloser, error) {
	filename := fmt.Sprintf("output/file%d.txt", routine)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed creating file: %w", err)
	}

	return &fileWriter{file: file, dataWriter: bufio.NewWriter(file)}, nil
}

func (f *fileWriter) Write(m []byte) (int, error) {
	n, err := f.dataWriter.Write(m)
	if err != nil {
		return n, err
	}

	if err = f.dataWriter.WriteByte('\n'); err != nil {
		return n, err
	}
	return n + 1, nil
}

func (f *fileWriter) Close() error {
	f.dataWriter.Flush()
	return f.file.Close()
}
