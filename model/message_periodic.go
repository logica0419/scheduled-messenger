package model

import (
	"github.com/google/uuid"
)

// 定期投稿メッセージのモデル定義
type SchMesPeriodic struct {
	ID        uuid.UUID    `gorm:"type:char(36);not null;primaryKey"` // メッセージの ID
	UserID    string       `gorm:"type:varchar(32);not null;index"`   // 予約ユーザーの traQ ID
	Time      PeriodicTime `gorm:"embedded;embeddedPrefix:time_"`     // 定期投稿する時間
	Repeat    *int         `gorm:"type:int"`                          // 定期投稿する回数
	ChannelID uuid.UUID    `gorm:"type:char(36);not null"`            // 投稿先チャンネルの ID
	Body      string       `gorm:"not null"`                          // メッセージ本文
}

// 定期投稿メッセージの時間モデル定義
type PeriodicTime struct {
	Month  *int `gorm:"type:int(2)"` // 月
	Date   *int `gorm:"type:int(2)"` // 日
	Hour   *int `gorm:"type:int(2)"` // 時
	Minute *int `gorm:"type:int(2)"` // 分
	Day    *int `gorm:"type:int(1)"` // 曜日
}
