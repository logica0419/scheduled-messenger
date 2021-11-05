package repository

import "gorm.io/gorm"

type Repository interface {
	migrate() error
	getTx() *gorm.DB
	SchMesRepository
}
