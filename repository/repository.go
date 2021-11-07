package repository

import "gorm.io/gorm"

// DB 操作リポジトリ
type Repository interface {
	migrate() error
	getTx() *gorm.DB
	SchMesRepository
}
