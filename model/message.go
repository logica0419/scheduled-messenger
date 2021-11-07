package model

import (
	"time"

	"github.com/google/uuid"
)

// スケジュールドメッセージのモデル定義
type SchMes struct {
	// メッセージの ID
	ID uuid.UUID `gorm:"primaryKey"`
	// 予約ユーザーの traQ ID
	UserID string `gorm:"index"`
	// 投稿時間
	Time time.Time `gorm:"index"`
	// 投稿先チャンネルの ID
	ChannelID uuid.UUID
	// メッセージ本文
	Body string
}
