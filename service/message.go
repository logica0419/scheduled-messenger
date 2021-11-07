package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
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
		"%s に`%s`、以下の内容を投稿します。\n```plaintext\n%s\n```\n予約を取り消したい場合は次のコマンドを Scheduled Messenger に送信して下さい。\n`!delete -i %s`",
		distChannel,
		parsedTime.Format("2006年01月02日 15:04"),
		body,
		id.String(),
	)
}

// スケジュール削除時のメッセージを生成
func CreateScheduleDeletedMessage(id string) string {
	return fmt.Sprintf("ID:`%s`のメッセージ送信予約を取り消しました。", id)
}

// スケジュールリストの表 (MD) を生成
func CreateScheduleListMessage(mesList []*model.SchMes) string {
	var result string

	// メッセージがない場合はその旨を伝える
	if len(mesList) == 0 {
		result = "あなたが予約済みのメッセージはありません。"
	} else {
		// ヘッダー
		result = "|メッセージID|予約時刻|投稿先チャンネルID|本文|\n|----|----|----|----|"

		// メッセージごとに行を追加
		for _, mes := range mesList {
			// 改行記号を string として表示できるよう変換
			replacedBody := strings.Replace(mes.Body, "\n", "`\\n`", -1)

			result += fmt.Sprintf("\n|%s|%s|%s|%s|", mes.ID, mes.Time.Format("2006年01月02日 15:04"), mes.ChannelID, replacedBody)
		}
	}

	return result
}

// DB のレコードから実際に送るメッセージを生成
func CreateScheduledMessage(mes *model.SchMes) string {
	return fmt.Sprintf("### @%s さんからのメッセージ:\n%s", mes.UserID, mes.Body)
}
