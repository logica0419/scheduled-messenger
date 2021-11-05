package router

import (
	"fmt"
	"log"
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

	// コマンドが含まれているか確認
	for _, cmd := range service.Commands {
		if strings.Contains(req.GetText(), cmd) {
			switch cmd {
			case service.Commands["schedule"]:
				// メッセージをパースし、要素を取得
				parsedTime, distChannel, distChannelID, body, err := service.ParseScheduleCommand(api, req)
				if err != nil {
					return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
				}

				log.Print(distChannelID)

				// 確認メッセージを送信
				mes := service.CreateScheduleCreatedMessage(parsedTime, distChannel, body)
				err = api.SendMessage(req.GetChannelID(), mes)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to send message: %s", err)})
				}
			case service.Commands["join"]:
				// チャンネルに JOIN する
				err = api.ChannelAction("join", req.GetChannelID())
				if err != nil {
					return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to join the channel: %s", err)})
				}
			case service.Commands["leave"]:
				// チャンネルから LEAVE する
				err = api.ChannelAction("leave", req.GetChannelID())
				if err != nil {
					return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to leave the channel: %s", err)})
				}
			}
		}
	}

	return c.NoContent(http.StatusNoContent)
}
