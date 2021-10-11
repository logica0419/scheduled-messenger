package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
	"github.com/logica0419/scheduled-messenger-bot/service"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

const (
	pingEvent   = "PING"   // PING イベント
	joinedEvent = "JOINED" // JOINED イベント
	leftEvent   = "LEFT"   // LEFT イベント
)

// Botのハンドラ (ヘッダーの "X-TRAQ-BOT-EVENT" を見てイベントごとにハンドラを割り振る)
func (r *Router) botEventHandler(c echo.Context) error {
	err := func() error {
		switch c.Request().Header.Get(botEventHeader) {
		case pingEvent:
			return pingHandler(c)
		case joinedEvent, leftEvent:
			return systemHandler(c, r.Api)
		default: // 未実装のイベント
			return c.JSON(http.StatusNotImplemented, errorMessage{Message: "not implemented"})
		}
	}()

	return err
}

// PING システムイベントハンドラ
func pingHandler(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// JOINED / LEFT システムイベントハンドラ
func systemHandler(c echo.Context, api *api.API) error {
	// リクエストボディの取得
	req := &event.SystemEvent{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to get request body: %s", err)})
	}

	// メッセージの生成
	chanPath := req.GetChannelPath()
	var mes string
	switch c.Request().Header.Get(botEventHeader) {
	case joinedEvent:
		mes = service.CreateJoinedMessage(chanPath)
	case leftEvent:
		mes = service.CreateLeftMessage()
	}

	// メッセージを JOINED / LEFT したチャンネルに送信
	chanID := req.GetChannelID()
	err = api.SendMessage(chanID, mes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to send join message: %s", err)})
	}

	return c.NoContent(http.StatusNoContent)
}
