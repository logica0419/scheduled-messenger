package model

import (
	"time"

	"github.com/google/uuid"
)

// スケジュールドメッセージのモデル定義
type SchMes struct {
	// メッセージの ID
	ID uuid.UUID `gorm:"type:char(36);not null;primaryKey"`

	// 予約ユーザーの traQ ID
	UserID string `gorm:"type:varchar(32);not null;index"`

	// 投稿時間
	Time time.Time `gorm:"not null;index"`

	// 投稿先チャンネルの ID
	ChannelID uuid.UUID `gorm:"type:char(36);not null"`

	// メッセージ本文
	Body string `gorm:"not null"`
}
