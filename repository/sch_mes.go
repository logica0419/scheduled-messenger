package repository

import "github.com/logica0419/scheduled-messenger-bot/model"

type SchMesRepository interface {
	ResisterSchMes(mes *model.SchMes) error
	DeleteSchMes(mes *model.SchMes) error
}
