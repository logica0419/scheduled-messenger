package repository

import (
	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

// 指定された ID のスケジュールドメッセージのレコードを取得
func (repo *GormRepository) GetSchMesByID(mesID uuid.UUID) (*model.SchMes, error) {
	// 空のメッセージ構造体の変数を作成
	var schMes *model.SchMes

	// レコードを取得
	res := repo.getTx().Where("id = ?", mesID).Take(&schMes)
	if res.Error != nil {
		return nil, res.Error
	}

	return schMes, nil
}

// 指定された UserID のスケジュールドメッセージのレコードを全取得
func (repo *GormRepository) GetSchMesByUserID(userID string) ([]*model.SchMes, error) {
	// 空のメッセージ構造体の変数を作成
	var schMes []*model.SchMes

	// レコードを取得
	res := repo.getTx().Where("user_id = ?", userID).Find(&schMes)
	if res.Error != nil {
		return nil, res.Error
	}

	return schMes, nil
}

// スケジュールドメッセージのレコードを新規作成
func (repo *GormRepository) ResisterSchMes(schMes *model.SchMes) error {
	// レコードを作成
	res := repo.getTx().Create(schMes)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// 指定された ID のスケジュールドメッセージのレコードを削除
func (repo *GormRepository) DeleteSchMesByID(mesID uuid.UUID) error {
	// 指定された ID のレコードを検索 (存在しない ID の検証)
	_, err := repo.GetSchMesByID(mesID)
	if err != nil {
		return err
	}

	// ID のみのメッセージ構造体の変数を作成 (primary key 指定のため)
	schMes := model.SchMes{
		ID: mesID,
	}

	// 指定された ID のレコードを削除
	res := repo.getTx().Delete(&schMes)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
