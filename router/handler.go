package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/logica0419/scheduled-messenger-bot/model"
	"github.com/logica0419/scheduled-messenger-bot/service"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

const (
	pingEvent = "PING"   // PING イベント
	joinEvent = "JOINED" // JOIN イベント
)

// Botのハンドラ (ヘッダーの "X-TRAQ-BOT-EVENT" を見てイベントごとにハンドラを割り振る)
func botEventHandler(c echo.Context) error {
	err := func() error {
		switch c.Request().Header.Get(botEventHeader) {
		case pingEvent:
			return pingHandler(c)
		case joinEvent:
			return joinHandler(c)
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

// JOIN イベントハンドラ
func joinHandler(c echo.Context) error {
	// リクエストボディの取得
	req := &model.RoomReq{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to get request body: %s", err)})
	}

	// メッセージの生成
	chanPath := req.GetChannelPath()
	mes := service.CreateJoinMessage(chanPath)

	// メッセージを JOIN したチャンネルに送信
	chanID := req.GetChannelID()
	err = api.SendMessage(chanID, mes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to send join message: %s", err)})
	}

	return c.NoContent(http.StatusNoContent)
}
