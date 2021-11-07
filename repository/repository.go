package repository

import "gorm.io/gorm"

// DB 操作リポジトリ
type Repository interface {
	// マイグレーションを実行
	migrate() error

	// 新たなトランザクション用 DB セッションを取得
	getTx() *gorm.DB

	// sch_mes テーブル操作リポジトリ
	SchMesRepository
}
