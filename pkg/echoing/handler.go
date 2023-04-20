package echoing

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

var connections uint64
var disconnections uint64

func Handler(ctx echo.Context) error {
	c, err := websocket.Accept(ctx.Response().Writer, ctx.Request(), nil)
	if err != nil {
		return err
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	atomic.AddUint64(&connections, 1)
	log.Print("Connections: ", connections-disconnections)

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err = Echo(ctx.Request().Context(), c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			atomic.AddUint64(&disconnections, 1)
			log.Printf("Status: connections=%d, disconnections=%d", connections-disconnections, disconnections)
			return nil
		}
		if err != nil {
			return err
		}
	}
}
