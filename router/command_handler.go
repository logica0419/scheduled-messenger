package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/service"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
	"gorm.io/gorm"
)

// コマンド一覧
var commands = map[string]string{
	"schedule": "!schedule", // 予約メッセージを作成する
	"delete":   "!delete",   // 予約メッセージを削除する
	"list":     "!list",     // 予約メッセージをリスト表示する
	"join":     "!join",     // チャンネルに JOIN する
	"leave":    "!leave",    // チャンネルから LEAVE する
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
	err = service.DeleteSchMesByID(repo, api, id, req.GetUserID())
	if err != nil {
		// 指定した ID のメッセージが存在しない場合エラーメッセージを送信
		if uuid.IsInvalidLengthError(err) || errors.Is(err, gorm.ErrRecordNotFound) {
			_ = api.SendMessage(req.GetChannelID(), "メッセージの削除に失敗しました\n```plaintext\n存在しないIDです\n```")
			return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
		}
		// 予約したユーザーと削除を試みたユーザーが違う場合エラーメッセージを送信
		if errors.Is(err, fmt.Errorf("access forbidden")) {
			_ = api.SendMessage(req.GetChannelID(), "メッセージの削除に失敗しました\n```plaintext\n権限がありません\n```")
			return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
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

// list コマンドハンドラー
func listHandler(c echo.Context, api *api.API, repo repository.Repository, req *event.MessageEvent) error {
	// ユーザー ID を取得
	userID := req.GetUserID()

	// スケジュールを DB から取得
	mesList, err := repo.GetSchMesByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: err.Error()})
	}

	// 予約メッセージリストを送信
	mes := service.CreateScheduleListMessage(mesList)
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
