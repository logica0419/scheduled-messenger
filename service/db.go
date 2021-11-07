package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

// 新たなメッセージを生成し、DB に登録
func ResisterSchMes(repo repository.Repository, userID string, time time.Time, channelID string, body string) (*model.SchMes, error) {
	channelUUID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, err
	}

	schMes, err := generateSchMes(userID, time, channelUUID, body)
	if err != nil {
		return nil, err
	}

	err = repo.ResisterSchMes(schMes)
	if err != nil {
		return nil, err
	}

	return schMes, nil
}

// 新たな ScheduleMes 構造体型変数を生成
func generateSchMes(userID string, time time.Time, channelID uuid.UUID, body string) (*model.SchMes, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &model.SchMes{
		ID:        id,
		UserID:    userID,
		Time:      time,
		ChannelID: channelID,
		Body:      body,
	}, nil
}

// 指定された ID のメッセージを DB から削除
func DeleteSchMesByID(repo repository.Repository, api *api.API, mesID string) error {
	mesUUID, err := uuid.Parse(mesID)
	if err != nil {
		return err
	}

	err = repo.DeleteSchMesByID(mesUUID)
	if err != nil {
		return err
	}

	return nil
}
