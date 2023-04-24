package main

import (
	"go-websocket-testing/pkg/echoing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", echoing.Handler)
	e.Logger.Fatal(e.Start("localhost:1323"))
}
