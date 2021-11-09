package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

// 指定された ID の予約投稿メッセージのレコードを取得
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

// 指定された UserID の予約投稿メッセージのレコードを全取得
func (repo *GormRepository) GetSchMesByUserID(userID string) ([]*model.SchMes, error) {
	// 空のメッセージ構造体の変数を作成
	var schMes []*model.SchMes

	// レコードを取得
	res := repo.getTx().Order("time asc").Where("user_id = ?", userID).Find(&schMes)
	if res.Error != nil {
		return nil, res.Error
	}

	return schMes, nil
}

// 指定された時間より前の time を持つ予約投稿メッセージのレコードを全取得
func (repo *GormRepository) GetSchMesByTime(time time.Time) ([]*model.SchMes, error) {
	// 空の予約投稿メッセージ構造体の変数を作成
	var schMes []*model.SchMes

	// レコードを取得
	res := repo.getTx().Where("time <= ?", time).Find(&schMes)
	if res.Error != nil {
		return nil, res.Error
	}

	return schMes, nil
}

// 予約投稿メッセージのレコードを新規作成
func (repo *GormRepository) ResisterSchMes(schMes *model.SchMes) error {
	// レコードを作成
	res := repo.getTx().Create(schMes)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// 指定された ID の予約投稿メッセージのレコードを削除
func (repo *GormRepository) DeleteSchMesByID(mesID uuid.UUID) error {
	// ID のみの予約投稿メッセージ構造体の変数を作成 (primary key 指定のため)
	schMes := model.SchMes{
		ID: mesID,
	}

	// レコードを削除
	res := repo.getTx().Delete(&schMes)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
