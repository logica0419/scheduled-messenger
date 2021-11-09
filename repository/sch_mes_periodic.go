package repository

import (
	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

// 定期投稿テーブル操作リポジトリ
type SchMesPeriodicRepository interface {
	// 定期投稿メッセージを全取得
	GetSchMesPeriodicAll() ([]*model.SchMesPeriodic, error)

	// 指定された ID の定期投稿メッセージのレコードを取得
	GetSchMesPeriodicByID(mesID uuid.UUID) (*model.SchMesPeriodic, error)

	// 指定された UserID の予約投稿メッセージのレコードを全取得
	GetSchMesPeriodicByUserID(userID string) ([]*model.SchMesPeriodic, error)

	// 定期投稿メッセージのレコードを新規作成
	ResisterSchMesPeriodic(schMesPeriodic *model.SchMesPeriodic) error

	// 指定された ID の定期投稿メッセージのレコードを削除
	DeleteSchMesPeriodicByID(mesID uuid.UUID) error
}
