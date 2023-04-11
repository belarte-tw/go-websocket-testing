package main

import (
	"go-socket/pkg/echoing"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

func hello(ctx echo.Context) error {
	c, err := websocket.Accept(ctx.Response().Writer, ctx.Request(), nil)
	log.Print("Connecting")
	if err != nil {
		return err
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err = echoing.Echo(ctx.Request().Context(), c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			log.Print("Disconnecting...")
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":1323"))
}
