package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
	"github.com/logica0419/scheduled-messenger-bot/repository"
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

// Botのハンドラ
func (r *Router) botEventHandler(c echo.Context) error {
	// ヘッダーの "X-TRAQ-BOT-EVENT" を見てイベントごとにハンドラを割り振る
	switch c.Request().Header.Get(botEventHeader) {
	case pingEvent:
		return pingHandler(c)

	case joinedEvent, leftEvent:
		return systemHandler(c, r.Api)

	case messageCreatedEvent, directMessageCreatedEvent:
		return messageEventHandler(c, r.Api, r.Repo)

	default: // 未実装のイベント
		return c.JSON(http.StatusNotImplemented, errorMessage{Message: "not implemented"})
	}
}

// PING システムイベントハンドラ
func pingHandler(c echo.Context) error {
	// NoContent を返す
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

	// メッセージを作成
	chanPath := req.GetChannelPath()
	var mes string
	switch c.Request().Header.Get(botEventHeader) {
	case joinedEvent:
		mes = service.CreateJoinedMessage(chanPath)
	case leftEvent:
		mes = service.CreateLeftMessage()
	}

	// メッセージを JOINED / LEFT したチャンネルに送信
	err = api.SendMessage(req.GetChannelID(), mes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to send message: %s", err)})
	}

	return c.NoContent(http.StatusNoContent)
}

func messageEventHandler(c echo.Context, api *api.API, repo repository.Repository) error {
	// リクエストボディの取得
	req := &event.MessageEvent{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to get request body: %s", err)})
	}

	// コマンドが含まれていた場合、コマンドハンドラーを呼び出す
	for _, cmd := range commands {
		if strings.Contains(req.GetText(), cmd) {
			switch cmd {
			case commands["help"]:
				return helpHandler(c, api, req)

			case commands["schedule"]:
				return scheduleHandler(c, api, repo, req)

			case commands["delete"]:
				return deleteHandler(c, api, repo, req)

			case commands["list"]:
				return listHandler(c, api, repo, req)

			case commands["join"]:
				return joinHandler(c, api, req)

			case commands["leave"]:
				return leaveHandler(c, api, req)
			}
		}
	}

	return c.NoContent(http.StatusNoContent)
}
