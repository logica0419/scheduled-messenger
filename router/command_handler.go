package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/service"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
	"gorm.io/gorm"
)

var commands = map[string]string{
	"schedule": "!schedule",
	"join":     "!join",
	"leave":    "!leave",
	"delete":   "!delete",
}

// schedule コマンドハンドラー
func scheduleHandler(c echo.Context, api *api.API, repo repository.Repository, req *event.MessageEvent) error {
	// メッセージをパースし、要素を取得
	parsedTime, distChannel, distChannelID, body, err := service.ParseScheduleCommand(api, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
	}

	// スケジュールを DB に登録
	schMes, err := service.ResisterSchMes(repo, req.GetUserID(), parsedTime, distChannelID, body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: err.Error()})
	}

	// 確認メッセージを送信
	mes := service.CreateScheduleCreatedMessage(schMes.Time, distChannel, schMes.Body, schMes.ID)
	err = api.SendMessage(req.GetChannelID(), mes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to send message: %s", err)})
	}

	return c.NoContent(http.StatusNoContent)
}

//delete コマンドハンドラー
func deleteHandler(c echo.Context, api *api.API, repo repository.Repository, req *event.MessageEvent) error {
	// メッセージをパースし、要素を取得
	id, err := service.ParseDeleteCommand(api, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
	}

	// スケジュールを DB から削除
	err = service.DeleteSchMes(repo, api, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = api.SendMessage(req.GetChannelID(), "メッセージの削除に失敗しました\n```plaintext\n存在しないIDです\n```")
		}
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: err.Error()})
	}

	// 確認メッセージを送信
	mes := service.CreateScheduleDeletedMessage(id)
	err = api.SendMessage(req.GetChannelID(), mes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to send message: %s", err)})
	}

	return c.NoContent(http.StatusNoContent)
}

// join コマンドハンドラー
func joinHandler(c echo.Context, api *api.API, req *event.MessageEvent) error {
	// チャンネルに JOIN する
	err := api.ChannelAction("join", req.GetChannelID())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to join the channel: %s", err)})
	}

	return c.NoContent(http.StatusNoContent)
}

// leave コマンドハンドラー
func leaveHandler(c echo.Context, api *api.API, req *event.MessageEvent) error {
	// チャンネルから LEAVE する
	err := api.ChannelAction("leave", req.GetChannelID())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to leave the channel: %s", err)})
	}

	return c.NoContent(http.StatusNoContent)
}
