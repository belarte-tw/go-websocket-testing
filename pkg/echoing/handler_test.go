//go:build integration

package echoing_test

import (
	"context"
	"go-websocket-testing/pkg/echoing"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"nhooyr.io/websocket"
)

func TestHandler(t *testing.T) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", echoing.Handler)

	go func() {
		e.Logger.Fatal(e.Start("localhost:6789"))
	}()
	time.Sleep(1 * time.Second)

	ctx := context.Background()
	c, _, err := websocket.Dial(ctx, "ws://localhost:6789", nil)
	assert.Nil(t, err)
	defer func() {
		assert.Nil(t, c.Close(websocket.StatusNormalClosure, "Done!"))
	}()

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
			err = c.Write(ctx, websocket.MessageText, test.msg)
			assert.Nil(t, err)

			typ, msg, err := c.Read(ctx)
			assert.Nil(t, err)
			assert.Equal(t, websocket.MessageText, typ)
			assert.Equal(t, string(test.msg), string(msg))
		})
	}
}
