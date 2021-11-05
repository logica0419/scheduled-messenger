package repository

import (
	"fmt"

	"github.com/logica0419/scheduled-messenger-bot/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// データベースを取得
func getDB(c *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.MariaDB_Username, c.MariaDB_Password, c.MariaDB_Hostname, c.MariaDB_Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 新たなリポジトリを生成
func GetRepository(c *config.Config) (Repository, error) {
	// DB を取得
	db, err := getDB(c)
	if err != nil {
		return nil, err
	}

	// リポジトリを作成して Migration を実行
	repo := &GormRepository{
		db: db,
	}
	err = repo.migrate()
	if err != nil {
		return nil, err
	}

	return repo, nil
}
