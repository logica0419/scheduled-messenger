package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

// sch_mes テーブル操作リポジトリ
type SchMesRepository interface {
	GetSchMesByID(mesID uuid.UUID) (*model.SchMes, error)
	GetSchMesByUserID(userID string) ([]*model.SchMes, error)
	GetSchMesByTime(time time.Time) ([]*model.SchMes, error)
	ResisterSchMes(mes *model.SchMes) error
	DeleteSchMesByID(mesID uuid.UUID) error
}
