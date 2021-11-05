package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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
func CreateScheduleCreatedMessage(parsedTime time.Time, distChannel string, body string, id uuid.UUID) string {
	return fmt.Sprintf(
		"`%s`に`%s`、以下の内容を投稿します。\n```plaintext\n%s\n```\n予約を取り消したい場合は次のコマンドを Scheduled Messenger に送信して下さい\n`!delete -i %s`",
		distChannel,
		parsedTime.Format("2006年01月02日 15:04"),
		body,
		id.String(),
	)
}
