package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	pingEvent = "PING"
)

func botEventHandler(c echo.Context) error {
	err := func() error {
		switch c.Request().Header.Get(botEventHeader) {
		case pingEvent:
			return pingHandler(c)
		default:
			return c.JSON(http.StatusNotImplemented, errorMessage{Message: "not implemented"})
		}
	}()

	return err
}

func pingHandler(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
