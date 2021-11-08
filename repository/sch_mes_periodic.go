package repository

import (
	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

// 定期投稿テーブル操作リポジトリ
type SchMesPeriodicRepository interface {
	// 指定された ID の定期投稿メッセージのレコードを取得
	GetSchMesPeriodicByID(mesID uuid.UUID) (*model.SchMesPeriodic, error)

	// 定期投稿メッセージのレコードを新規作成
	ResisterSchMesPeriodic(schMesPeriodic *model.SchMesPeriodic) error
}
