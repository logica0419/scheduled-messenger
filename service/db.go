package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

// 新たなメッセージを生成し、DB に登録
func ResisterSchMes(repo repository.Repository, userID string, time time.Time, channelID string, body string) (*model.SchMes, error) {
	// チャンネル ID を UUID に変換
	channelUUID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, err
	}

	// 新たな SchMes 構造体型変数を生成
	schMes, err := generateSchMes(userID, time, channelUUID, body)
	if err != nil {
		return nil, err
	}

	// DB に登録
	err = repo.ResisterSchMes(schMes)
	if err != nil {
		return nil, err
	}

	return schMes, nil
}

// 新たな SchMes 構造体型変数を生成
func generateSchMes(userID string, time time.Time, channelID uuid.UUID, body string) (*model.SchMes, error) {
	// ID を生成
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// SchMes 構造体型変数を生成
	return &model.SchMes{
		ID:        id,
		UserID:    userID,
		Time:      time,
		ChannelID: channelID,
		Body:      body,
	}, nil
}

// 指定された ID のメッセージを DB から削除
func DeleteSchMesByID(repo repository.Repository, api *api.API, mesID string, userID string) error {
	// ID を UUID に変換
	mesUUID, err := uuid.Parse(mesID)
	if err != nil {
		return err
	}

	// 指定された ID のレコードを検索 (存在しない ID の検証)
	mes, err := repo.GetSchMesByID(mesUUID)
	if err != nil {
		return err
	}

	// 予約したユーザーと削除を試みたユーザーが一致するか検証
	if mes.UserID != userID {
		return fmt.Errorf("access forbidden")
	}

	// DB から削除
	err = repo.DeleteSchMesByID(mesUUID)
	if err != nil {
		return err
	}

	return nil
}
