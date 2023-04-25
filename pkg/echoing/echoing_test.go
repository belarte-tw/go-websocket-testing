package echoing_test

import (
	"context"
	"errors"
	"go-websocket-testing/pkg/echoing"
	"testing"

	"nhooyr.io/websocket"
)

type mockConn struct {
	input      []byte
	output     []byte
	writeError error
}

func (c *mockConn) Read(context.Context) (websocket.MessageType, []byte, error) {
	return websocket.MessageText, c.input, nil
}

func (c *mockConn) Write(ctx context.Context, t websocket.MessageType, b []byte) error {
	c.output = append([]byte(nil), b...)
	return c.writeError
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
				input: test.msg,
			}

			err := echoing.Echo(context.TODO(), m)
			if err != nil {
				t.Fatal("Cannot echo")
			}

			got := string(m.output)
			want := string(test.msg)
			if got != want {
				t.Errorf("got '%v' but wanted '%v'", got, want)
			}
		})
	}
}

func TestEchoingFailToWrite(t *testing.T) {
	want := errors.New("something bad happened")
	m := &mockConn{
		writeError: want,
	}

	err := echoing.Echo(context.TODO(), m)
	if err == nil {
		t.Fatal("Should return an error but was nil")
	}

	if !errors.Is(err, want) {
		t.Errorf("got '%s' but wanted '%s'", err, want)
	}
}
