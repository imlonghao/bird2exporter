package main

import (
	"os"
	"strings"

	"github.com/imlonghao/bird2exporter/bird"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var b *bird.Bird

func main() {
	birdCtl := "/run/bird.ctl"
	if birdCtlEnv := os.Getenv("BIRDCTL"); birdCtlEnv != "" {
		birdCtl = birdCtlEnv
	}
	b = bird.New(birdCtl)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	if origins := os.Getenv("ORIGINS"); origins != "" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: strings.Split(origins, ","),
		}))
	}
	e.GET("/", index)
	e.GET("/protocol", protocol)
	e.GET("/protocol/all", protocolAll)
	e.POST("/exec", exec)
	e.Logger.Fatal(e.Start(":1919"))
}
