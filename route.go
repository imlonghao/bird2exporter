package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type message struct {
	Q string `json:"q"`
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "?")
}

func protocol(c echo.Context) error {
	return c.JSON(http.StatusOK, b.ShowProtocol())
}

func protocolAll(c echo.Context) error {
	return c.JSON(http.StatusOK, b.ShowProtocolAll())
}

func exec(c echo.Context) error {
	command := new(message)
	if err := c.Bind(command); err != nil {
		return c.String(http.StatusBadRequest, "?")
	}
	for _, i := range command.Q {
		if i < 32 || i > 126 {
			return c.String(http.StatusBadRequest, "?")
		}
	}
	command.Q = "show " + command.Q
	return c.JSON(http.StatusOK, b.Exec(command.Q))
}
