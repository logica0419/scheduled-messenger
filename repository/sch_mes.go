package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

// sch_mes テーブル操作リポジトリ
type SchMesRepository interface {
	// 指定された ID のスケジュールドメッセージのレコードを取得
	GetSchMesByID(mesID uuid.UUID) (*model.SchMes, error)

	// 指定された UserID のスケジュールドメッセージのレコードを全取得
	GetSchMesByUserID(userID string) ([]*model.SchMes, error)

	// 指定された時間より前の time を持つメッセージのレコードを全取得
	GetSchMesByTime(time time.Time) ([]*model.SchMes, error)

	// スケジュールドメッセージのレコードを新規作成
	ResisterSchMes(mes *model.SchMes) error

	// 指定された ID のスケジュールドメッセージのレコードを削除
	DeleteSchMesByID(mesID uuid.UUID) error
}
