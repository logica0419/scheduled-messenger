package service

import (
	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

// 新たな定期投稿メッセージを生成し、DB に登録
func ResisterSchMesPeriodic(repo repository.Repository, userID string, time model.PeriodicTime, channelID string, body string, repeat *int) (*model.SchMesPeriodic, error) {
	// チャンネル ID を UUID に変換
	channelUUID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, err
	}

	// 新たな SchMesPeriodic 構造体型変数を生成
	schMesPeriodic, err := generateSchMesPeriodic(userID, time, channelUUID, body, repeat)
	if err != nil {
		return nil, err
	}

	// DB に登録
	err = repo.ResisterSchMesPeriodic(schMesPeriodic)
	if err != nil {
		return nil, err
	}

	return schMesPeriodic, nil
}

// 新たな SchMesPeriodic 構造体型変数を生成
func generateSchMesPeriodic(userID string, time model.PeriodicTime, channelID uuid.UUID, body string, repeat *int) (*model.SchMesPeriodic, error) {
	// ID を生成
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// SchMes 構造体型変数を生成
	return &model.SchMesPeriodic{
		ID:        id,
		UserID:    userID,
		Time:      time,
		Repeat:    repeat,
		ChannelID: channelID,
		Body:      body,
	}, nil
}

// 指定された ID の定期投稿メッセージを DB から削除
func DeleteSchMesByIDPeriodic(id string, repo repository.Repository, api *api.API, mesID string, userID string) error {
	// ID を UUID に変換
	mesUUID, err := uuid.Parse(mesID)
	if err != nil {
		return err
	}

	// 指定された ID のレコードを検索 (存在しない ID の検証)
	mes, err := repo.GetSchMesPeriodicByID(mesUUID)
	if err != nil {
		return err
	}

	// 予約したユーザーと削除を試みたユーザーが一致するか検証
	if mes.UserID != userID {
		return ErrUserNotMatch
	}

	// DB から削除
	err = repo.DeleteSchMesPeriodicByID(mesUUID)
	if err != nil {
		return err
	}

	return nil
}
