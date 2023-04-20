package echoing_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"go-socket/pkg/echoing"

	"nhooyr.io/websocket"
)

type mockReader struct {
	msg []byte
}

func (m mockReader) Read(b []byte) (int, error) {
	copy(b, m.msg)
	b = b[:len(m.msg)]
	return len(b), io.EOF
}

type mockWriteCloser struct {
	msg []byte
	err error
}

func (m *mockWriteCloser) Close() error { return m.err }
func (m *mockWriteCloser) Write(b []byte) (n int, err error) {
	m.msg = append([]byte(nil), b...)
	return len(b), nil
}

type mockConn struct {
	reader mockReader
	writer mockWriteCloser
}

func (c *mockConn) Reader(context.Context) (websocket.MessageType, io.Reader, error) {
	return websocket.MessageText, &c.reader, nil
}

func (c *mockConn) Writer(context.Context, websocket.MessageType) (io.WriteCloser, error) {
	return &c.writer, nil
}

func TestEchoing(t *testing.T) {
	tests := map[string]struct {
		msg []byte
	}{
		"empty message":            {msg: []byte("")},
		"simple message":           {msg: []byte("hello world")},
		"multi line message":       {msg: []byte("hello world\nI am a developer!")},
		"message with white space": {msg: []byte("  hello world  ")},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mockConn{
				writer: mockWriteCloser{},
				reader: mockReader{msg: test.msg},
			}

			err := echoing.Echo(context.TODO(), m)
			if err != nil {
				t.Fatal("Cannot echo")
			}

			got := string(m.writer.msg)
			want := string(test.msg)
			if got != want {
				t.Errorf("got '%v' but wanted '%v'", got, want)
			}
		})
	}
}

func TestEchoingFailToCloseWriter(t *testing.T) {
	want := "something bad happened"
	m := &mockConn{
		writer: mockWriteCloser{err: errors.New(want)},
	}

	err := echoing.Echo(context.TODO(), m)
	if err == nil {
		t.Fatal("Should return an error but was nil")
	}

	got := err.Error()
	if got != want {
		t.Errorf("got '%s' but wanted '%s'", got, want)
	}
}
