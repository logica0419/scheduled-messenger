package service

import (
	"fmt"
	"time"
)

// チャンネルに JOINED した際のメッセージを生成
func CreateJoinedMessage(path string) string {
	return fmt.Sprintf("これから Scheduled Messenher Bot は`%s`に投稿されるメッセージをチェックします!\nいつでも呼んで下さい!", path)
}

// チャンネルから LEFT した際のメッセージを生成
func CreateLeftMessage() string {
	return "寂しいですがお別れです...\nScheduled Messenher Bot のご利用、ありがとうございました!"
}

// スケジュール作成時のメッセージを生成
func CreateScheduleCreatedMessage(parsedTime time.Time, distChannel string, body string) string {
	return fmt.Sprintf("`%s`に`%s`、以下の内容を投稿します。\n`%s`", distChannel, parsedTime.Format("2006年01月02日 15:04"), body)
}
