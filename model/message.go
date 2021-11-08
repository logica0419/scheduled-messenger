package model

import (
	"time"

	"github.com/google/uuid"
)

// 予約投稿メッセージのモデル定義
type SchMes struct {
	ID        uuid.UUID `gorm:"type:char(36);not null;primaryKey"` // メッセージの ID
	UserID    string    `gorm:"type:varchar(32);not null;index"`   // 予約ユーザーの traQ ID
	Time      time.Time `gorm:"not null;index"`                    // 投稿時間
	ChannelID uuid.UUID `gorm:"type:char(36);not null"`            // 投稿先チャンネルの ID
	Body      string    `gorm:"not null"`                          // メッセージ本文
}
