package echoing_test

import (
	"context"
	"errors"
	"testing"

	"go-socket/pkg/echoing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"nhooyr.io/websocket"
)

func TestIfyEchoing(t *testing.T) {
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
			mockConn := echoing.NewMockconn(t)
			wc := &mockWriteCloser{}
			mockConn.EXPECT().
				Reader(mock.Anything).
				Return(websocket.MessageText, &mockReader{msg: test.msg}, nil)
			mockConn.EXPECT().
				Writer(mock.Anything, websocket.MessageText).
				Return(wc, nil)

			err := echoing.Echo(context.TODO(), mockConn, l)

			assert.Nil(t, err, "Cannot echo")
			assert.Equal(t, string(test.msg), string(wc.msg), "Message should be equal")
		})
	}
}

func TestIfyEchoingFailToCloseWriter(t *testing.T) {
	want := "something bad happened"
	mockConn := echoing.NewMockconn(t)
	mockConn.EXPECT().
		Writer(mock.Anything, websocket.MessageText).
		Return(&mockWriteCloser{err: errors.New(want)}, nil)
	mockConn.EXPECT().
		Reader(mock.Anything).
		Return(websocket.MessageText, &mockReader{}, nil)

	err := echoing.Echo(context.TODO(), mockConn, l)

	assert.NotNil(t, err, "Should return an error but was nil")
	assert.Equal(t, want, err.Error(), "got '%s' but wanted '%s'", err.Error(), want)
}
