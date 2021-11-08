package model

import (
	"fmt"

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

// 曜日の対応表
var days = [7]string{"日", "月", "火", "水", "木", "金", "土"}

// 時間をメッセージに使えるようフォーマット
func (time PeriodicTime) Format() string {
	// 年
	mes := "毎年"

	// 月
	if time.Month != nil {
		mes += fmt.Sprintf(" %02d 月", *time.Month)
	} else {
		mes += "毎月"
	}

	// 日
	if time.Date != nil {
		mes += fmt.Sprintf(" %02d 日", *time.Date)
	} else {
		mes += "毎日"
	}

	mes += " "

	// 時
	if time.Hour != nil {
		mes += fmt.Sprintf(" %02d 時", *time.Hour)
	} else {
		mes += "毎時"
	}

	// 分
	if time.Minute != nil {
		mes += fmt.Sprintf(" %02d 分", *time.Minute)
	} else {
		mes += "毎分"
	}

	mes += " "

	// 曜日
	if time.Day != nil {
		mes += fmt.Sprintf("毎%s曜日", days[*time.Day])
	}

	return mes
}
