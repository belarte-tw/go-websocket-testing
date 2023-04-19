package echoing_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"go-socket/pkg/echoing"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type mockReader struct{}

func (m mockReader) Read(p []byte) (int, error) { return 0, io.EOF }

type mockWriteCloser struct {
	msg []byte
	err error
}

func (m mockWriteCloser) Close() error                      { return m.err }
func (m mockWriteCloser) Write(b []byte) (n int, err error) { return 0, nil }

type mockConn struct {
	writer mockWriteCloser
}

func (c mockConn) Reader(context.Context) (websocket.MessageType, io.Reader, error) {
	return websocket.MessageBinary, &mockReader{}, nil
}

func (c mockConn) Writer(context.Context, websocket.MessageType) (io.WriteCloser, error) {
	return &c.writer, nil
}

var l = rate.NewLimiter(rate.Every(time.Millisecond*100), 10)

func TestEchoing(t *testing.T) {
	w := mockWriteCloser{}
	m := &mockConn{writer: w}
	err := echoing.Echo(context.TODO(), m, l)
	if err != nil {
		t.Fatal("Cannot echo")
	}

	got := string(w.msg)
	want := ""
	if got != want {
		t.Errorf("got '%s' but wanted '%s'", got, want)
	}
}

func TestEchoingFailToCloseWriter(t *testing.T) {
	want := "something bad happened"
	m := &mockConn{
		writer: mockWriteCloser{err: errors.New(want)},
	}

	err := echoing.Echo(context.TODO(), m, l)
	if err == nil {
		t.Fatal("Should return an error but was nil")
	}

	got := err.Error()
	if got != want {
		t.Errorf("got '%s' but wanted '%s'", got, want)
	}
}
