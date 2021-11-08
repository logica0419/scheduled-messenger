package router

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/service"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
	"github.com/logica0419/scheduled-messenger-bot/service/parser"
	"gorm.io/gorm"
)

// コマンド一覧
var commands = map[string]string{
	"help":     "help",      // ヘルプを表示する
	"schedule": "!schedule", // 予約メッセージを作成する
	"delete":   "!delete",   // 予約メッセージを削除する
	"list":     "!list",     // 予約メッセージをリスト表示する
	"join":     "!join",     // チャンネルに JOIN する
	"leave":    "!leave",    // チャンネルから LEAVE する
}

// help コマンドハンドラー
func helpHandler(c echo.Context, api *api.API, req *event.MessageEvent) error {
	mes := service.CreateHelpMessage()

	err := api.SendMessage(req.GetChannelID(), mes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to send message: %s", err)})
	}

	return c.NoContent(http.StatusNoContent)
}

// schedule コマンドハンドラー
func scheduleHandler(c echo.Context, api *api.API, repo repository.Repository, req *event.MessageEvent) error {
	// メッセージをパースし、要素を取得
	time, distChannel, distChannelID, body, repeat, err := parser.ParseScheduleCommand(api, req)
	if err != nil {
		// Argv パース以外のところでエラーを吐いたらメッセージを送信
		_err := err
		if !errors.As(fmt.Errorf("failed to parse argv"), &_err) {
			service.SendCreateErrorMessage(api, req.GetChannelID(), fmt.Errorf("メッセージをパースできません\n%s", err))
		}
		return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
	}

	// 確認メッセージ
	var mes string
	// 時間の表記にワイルドカードが含まれているかで処理を分岐
	if strings.Contains(*time, "*") { // 定期投稿
		// 時間をパース
		parsedTime, err := parser.TimeParsePeriodic(time)
		if err != nil {
			service.SendCreateErrorMessage(api, req.GetChannelID(), fmt.Errorf("無効な時間表記です\n%s", err))
			return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
		}

		// スケジュールをDB に 登録
		schMesPeriodic, err := service.ResisterSchMesPeriodic(repo, req.GetUserID(), *parsedTime, *distChannelID, *body, repeat)
		if err != nil {
			service.SendCreateErrorMessage(api, req.GetChannelID(), fmt.Errorf("DB エラーです\n%s", err))
			return c.JSON(http.StatusInternalServerError, errorMessage{Message: err.Error()})
		}

		// 確認メッセージを生成
		mes = service.CreateSchMesPeriodicCreatedMessage(schMesPeriodic.Time, *distChannel, schMesPeriodic.Body, schMesPeriodic.ID, schMesPeriodic.Repeat)

	} else { // 予約投稿
		// repeat が入力されていたらエラーメッセージを送る
		if repeat != nil {
			service.SendCreateErrorMessage(api, req.GetChannelID(), fmt.Errorf("予約投稿でリピートは使用できません"))
			return c.JSON(http.StatusBadRequest, errorMessage{Message: "予約投稿でリピートは使用できません"})
		}

		// 時間をパース
		parsedTime, err := parser.TimeParse(time)
		if err != nil {
			service.SendCreateErrorMessage(api, req.GetChannelID(), fmt.Errorf("無効な時間表記です\n%s", err))
			return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
		}

		// スケジュールを DB に登録
		schMes, err := service.ResisterSchMes(repo, req.GetUserID(), *parsedTime, *distChannelID, *body)
		if err != nil {
			service.SendCreateErrorMessage(api, req.GetChannelID(), fmt.Errorf("DB エラーです\n%s", err))
			return c.JSON(http.StatusInternalServerError, errorMessage{Message: err.Error()})
		}

		// 確認メッセージを生成
		mes = service.CreateSchMesCreatedMessage(schMes.Time, *distChannel, schMes.Body, schMes.ID)
	}

	// 確認メッセージを送信
	err = api.SendMessage(req.GetChannelID(), mes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: fmt.Sprintf("failed to send message: %s", err)})
	}

	return c.NoContent(http.StatusNoContent)
}

//delete コマンドハンドラー
func deleteHandler(c echo.Context, api *api.API, repo repository.Repository, req *event.MessageEvent) error {
	// メッセージをパースし、要素を取得
	id, err := parser.ParseDeleteCommand(api, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
	}

	// 予約投稿スケジュールを DB から削除
	err = service.DeleteSchMesByID(repo, api, *id, req.GetUserID())
	if err != nil {
		// 指定した ID のメッセージが存在しない場合定期投稿の削除を試みる
		if errors.Is(err, gorm.ErrRecordNotFound) {
			goto periodic
		}
		// 指定した ID が無効な場合エラーメッセージを送信
		if uuid.IsInvalidLengthError(err) {
			service.SendDeleteErrorMessage(api, req.GetChannelID(), fmt.Errorf("存在しないIDです"))
			return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
		}
		// 予約したユーザーと削除を試みたユーザーが違う場合エラーメッセージを送信
		if errors.Is(err, service.ErrUserNotMatch) {
			service.SendDeleteErrorMessage(api, req.GetChannelID(), fmt.Errorf("予約メッセージは予約したユーザーしか削除できません"))
			return c.JSON(http.StatusForbidden, errorMessage{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: err.Error()})
	}

	goto message

periodic: // 定期投稿スケジュールを DB から削除
	err = service.DeleteSchMesPeriodicByID(repo, api, *id, req.GetUserID())
	if err != nil {
		// 指定した ID のメッセージが存在しない場合エラーメッセージを送信
		if errors.Is(err, gorm.ErrRecordNotFound) {
			service.SendDeleteErrorMessage(api, req.GetChannelID(), fmt.Errorf("存在しないIDです"))
			return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
		}
		// 予約したユーザーと削除を試みたユーザーが違う場合エラーメッセージを送信
		if errors.Is(err, service.ErrUserNotMatch) {
			service.SendDeleteErrorMessage(api, req.GetChannelID(), fmt.Errorf("予約メッセージは予約したユーザーしか削除できません"))
			return c.JSON(http.StatusForbidden, errorMessage{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, errorMessage{Message: err.Error()})
	}

	goto message

message: // 確認メッセージを送信
	mes := service.CreateSchMesDeletedMessage(*id)
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
