package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"nhooyr.io/websocket"
)

type fakeWriteCloser struct{}

func (f fakeWriteCloser) Write([]byte) (int, error) {
	return 0, nil
}

func (f fakeWriteCloser) Close() error {
	return nil
}

func TestWriteAndReadCorrectNumberOfMessages(t *testing.T) {
	conn := NewMockconn(t)
	conn.EXPECT().
		Write(mock.Anything, mock.Anything, mock.Anything).
		Times(5).
		Return(nil)
	conn.EXPECT().
		Read(mock.Anything).
		Times(5).
		Return(websocket.MessageText, []byte(""), nil)

	r := routine{ctx: context.Background(), conn: conn, datawriter: fakeWriteCloser{}}
	r.start(5)
}

func TestRoutineIsClose(t *testing.T) {
	conn := NewMockconn(t)
	conn.EXPECT().Close(websocket.StatusNormalClosure, "Done!").Return(nil)

	r := routine{ctx: context.Background(), conn: conn, datawriter: fakeWriteCloser{}}
	r.close()
}
