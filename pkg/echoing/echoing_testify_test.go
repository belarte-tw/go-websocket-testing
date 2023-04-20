package echoing_test

import (
	"context"
	"errors"
	"testing"

	"go-socket/pkg/echoing"

	"github.com/stretchr/testify/assert"
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
			m := &mockConn{
				writer: mockWriteCloser{},
				reader: mockReader{msg: test.msg},
			}

			err := echoing.Echo(context.TODO(), m, l)

			assert.Nil(t, err, "Cannot echo")
			assert.Equal(t, string(test.msg), string(m.writer.msg), "Message should be equal")
		})
	}
}

func TestIfyEchoingFailToCloseWriter(t *testing.T) {
	want := "something bad happened"
	m := &mockConn{
		writer: mockWriteCloser{err: errors.New(want)},
	}

	err := echoing.Echo(context.TODO(), m, l)

	assert.NotNil(t, err)
	assert.Equal(t, want, err.Error(), "Wrong error message")
}
