package repository

import (
	"github.com/logica0419/scheduled-messenger-bot/model"
	"gorm.io/gorm"
)

// リポジトリ実装
type GormRepository struct {
	db *gorm.DB
}

// 新たなトランザクション用 DB セッションを取得
func (repo *GormRepository) getTx() *gorm.DB {
	return repo.db.Session(&gorm.Session{})
}

// 現在の全モデル
var models = []interface{}{
	&model.SchMes{},
	&model.SchMesPeriodic{},
}

// DB マイグレーションを実行
func (repo *GormRepository) migrate() error {
	err := repo.db.AutoMigrate(models...)
	if err != nil {
		return err
	}

	return nil
}
