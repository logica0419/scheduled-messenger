package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	pingEvent = "PING" //PING イベント
)

// Botのハンドラ (ヘッダーの "X-TRAQ-BOT-EVENT" を見てイベントごとにハンドラを割り振る)
func botEventHandler(c echo.Context) error {
	err := func() error {
		switch c.Request().Header.Get(botEventHeader) {
		case pingEvent: //PING イベント
			return pingHandler(c)
		default: // 未実装のイベント
			return c.JSON(http.StatusNotImplemented, errorMessage{Message: "not implemented"})
		}
	}()

	return err
}

// PING イベントハンドラ
func pingHandler(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
