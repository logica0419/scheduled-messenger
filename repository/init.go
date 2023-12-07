package repository

import (
	"fmt"
	"time"

	"github.com/logica0419/scheduled-messenger-bot/config"
	"github.com/patrickmn/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 新たなリポジトリを生成
func GetRepository(c *config.Config) (Repository, error) {
	// DB を取得
	db, err := getDB(c)
	if err != nil {
		return nil, err
	}

	// 新規キャッシュを作成
	cache := cache.New(0, 10*time.Minute)

	// リポジトリを作成して Migration を実行
	repo := &GormRepository{
		db: db,
		c:  cache,
	}
	err = repo.migrate()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// DB を取得
func getDB(c *config.Config) (*gorm.DB, error) {
	// DSN を生成
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.NS_MariaDB_Username,
		c.NS_MariaDB_Password,
		c.NS_MariaDB_Hostname,
		c.NS_MariaDB_Database,
	)

	// DB に接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
