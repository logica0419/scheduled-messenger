package repository

import "github.com/logica0419/scheduled-messenger-bot/model"

// スケジュールドメッセージのレコードを新規作成
func (repo *GormRepository) ResisterSchMes(schMes *model.SchMes) error {
	// レコードを作成
	res := repo.getTx().Create(schMes)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo *GormRepository) DeleteSchMes(schMes *model.SchMes) error {
	return repo.getTx().Delete(schMes).Error
}
