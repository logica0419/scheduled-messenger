package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
	"github.com/logica0419/scheduled-messenger-bot/service"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

const (
	pingEvent                 = "PING"                   // PING イベント
	joinedEvent               = "JOINED"                 // JOINED イベント
	leftEvent                 = "LEFT"                   // LEFT イベント
	messageCreatedEvent       = "MESSAGE_CREATED"        // MESSAGE_CREATED イベント
	directMessageCreatedEvent = "DIRECT_MESSAGE_CREATED" // DIRECT_MESSAGE_CREATED イベント
)

// Botのハンドラ (ヘッダーの "X-TRAQ-BOT-EVENT" を見てイベントごとにハンドラを割り振る)
func (r *Router) botEventHandler(c echo.Context) error {
	err := func() error {
		switch c.Request().Header.Get(botEventHeader) {
		case pingEvent:
			return pingHandler(c)
		case joinedEvent, leftEvent:
			return systemHandler(c, r.Api)
		case messageCreatedEvent, directMessageCreatedEvent:
			return messageEventHandler(c, r.Api)
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

	// メッセージを JOINED / LEFT したチャンネルに送信
	chanPath := req.GetChannelPath()
	var mes string
	switch c.Request().Header.Get(botEventHeader) {
	case joinedEvent:
		mes = service.CreateJoinedMessage(chanPath)
	case leftEvent:
		mes = service.CreateLeftMessage()
	}

	err = api.SendMessage(req.GetChannelID(), mes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to send message: %s", err)})
	}

	return c.NoContent(http.StatusNoContent)
}

func messageEventHandler(c echo.Context, api *api.API) error {
	// リクエストボディの取得
	req := &event.MessageEvent{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to get request body: %s", err)})
	}

	// Botからのリクエストは無視
	if req.Message.User.Bot {
		return c.JSON(http.StatusForbidden, errorMessage{Message: "message from bot"})
	}

	// コマンドが含まれているか確認
	for _, cmd := range service.Commands {
		if strings.Contains(req.GetText(), cmd) {
			// メッセージを配列に
			listedReqMes, err := service.ArgvParse(req.GetText())
			if err != nil {
				return c.JSON(http.StatusBadRequest, errorMessage{Message: fmt.Sprintf("failed to parse argv: %s", err)})
			}

			// 冒頭にメンションがついていた場合要素をドロップ
			if strings.Contains(listedReqMes[0], "@") {
				listedReqMes = listedReqMes[1:]
			}

			// メッセージをパース
			parsedTime, distChannel, body, err := service.ParseScheduleMessage(listedReqMes)
			if err != nil {
				_ = api.SendMessage(req.GetChannelID(), err.Error())
				return c.JSON(http.StatusBadRequest, errorMessage{Message: fmt.Sprintf("failed to parse schedule message: %s", err)})
			}

			// 確認メッセージを送信
			mes := service.CreateScheduleCreatedMessage(parsedTime, distChannel, body)
			err = api.SendMessage(req.GetChannelID(), mes)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to send message: %s", err)})
			}
		}
	}

	return c.NoContent(http.StatusNoContent)
}
