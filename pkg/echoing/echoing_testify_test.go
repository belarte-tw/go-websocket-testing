package echoing_test

import (
	"context"
	"errors"
	"testing"

	"go-websocket-testing/pkg/echoing"

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
			mockConn.EXPECT().
				Read(mock.Anything).
				Return(websocket.MessageText, test.msg, nil)
			mockConn.EXPECT().
				Write(mock.Anything, websocket.MessageText, test.msg).
				Return(nil)

			err := echoing.Echo(context.TODO(), mockConn)

			assert.Nil(t, err, "Cannot echo")
		})
	}
}

func TestIfyEchoingFailToWrite(t *testing.T) {
	want := "something bad happened"
	mockConn := echoing.NewMockconn(t)
	mockConn.EXPECT().
		Write(mock.Anything, websocket.MessageText, mock.Anything).
		Return(errors.New(want))
	mockConn.EXPECT().
		Read(mock.Anything).
		Return(websocket.MessageText, []byte(""), nil)

	err := echoing.Echo(context.TODO(), mockConn)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), want)
}
