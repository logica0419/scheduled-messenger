package repository

import (
	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

type SchMesRepository interface {
	GetSchMesByID(mesID uuid.UUID) (*model.SchMes, error)
	GetSchMesByUserID(userID string) ([]*model.SchMes, error)
	ResisterSchMes(mes *model.SchMes) error
	DeleteSchMesByID(mesID uuid.UUID) error
}
