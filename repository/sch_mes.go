package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

// 予約投稿テーブル操作リポジトリ
type SchMesRepository interface {
	// 指定された ID の予約投稿メッセージのレコードを取得
	GetSchMesByID(mesID uuid.UUID) (*model.SchMes, error)

	// 指定された UserID の予約投稿メッセージのレコードを全取得
	GetSchMesByUserID(userID string) ([]*model.SchMes, error)

	// 指定された時間より前の time を持つメッセージのレコードを全取得
	GetSchMesByTime(time time.Time) ([]*model.SchMes, error)

	// 予約投稿メッセージのレコードを新規作成
	ResisterSchMes(mes *model.SchMes) error

	// 指定された ID の予約投稿メッセージのレコードを削除
	DeleteSchMesByID(mesID uuid.UUID) error
}
