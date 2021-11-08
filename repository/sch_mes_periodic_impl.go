package repository

import (
	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

// 指定された ID の定期投稿メッセージのレコードを取得
func (repo *GormRepository) GetSchMesPeriodicByID(mesID uuid.UUID) (*model.SchMesPeriodic, error) {
	// 空のメッセージ構造体の変数を作成
	var schMesPeriodic *model.SchMesPeriodic

	// レコードを取得
	res := repo.getTx().Where("id = ?", mesID).Take(&schMesPeriodic)
	if res.Error != nil {
		return nil, res.Error
	}

	return schMesPeriodic, nil
}

// 定期投稿メッセージのレコードを新規作成
func (repo *GormRepository) ResisterSchMesPeriodic(schMesPeriodic *model.SchMesPeriodic) error {
	// レコードを作成
	res := repo.getTx().Create(schMesPeriodic)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
