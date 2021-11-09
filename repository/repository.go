package repository

// DB 操作リポジトリ
type Repository interface {
	// 予約投稿テーブル操作リポジトリ
	SchMesRepository

	// 定期投稿テーブル操作リポジトリ
	SchMesPeriodicRepository
}
