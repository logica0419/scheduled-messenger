package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
)

// ヘルプメッセージを生成
func CreateHelpMessage() string {
	return "[Wiki](https://wiki.trap.jp/bot/Sch_Mes#head2) の使い方を参照してください！"
}

// チャンネルに JOINED した際のメッセージを生成
func CreateJoinedMessage(path string) string {
	return fmt.Sprintf("これから Scheduled Messenher は`%s`に投稿されるメッセージをチェックします!\nいつでも呼んで下さい!", path)
}

// チャンネルから LEFT した際のメッセージを生成
func CreateLeftMessage() string {
	return "寂しいですがお別れです...\nScheduled Messenger のご利用、ありがとうございました!"
}

// 予約投稿メッセージ作成時 / 編集時のメッセージを生成
func CreateSchMesCreatedEditedMessage(parsedTime time.Time, distChannel *string, body string, id uuid.UUID) string {
	// 空のメッセージを作成
	var mes string

	// チャンネルを追加
	if distChannel != nil {
		mes += fmt.Sprintf("%s に", *distChannel)
	} else {
		mes += "以前の登録と同じチャンネルに"
	}

	// 残りの文字列を追加
	mes += fmt.Sprintf("`%s`、以下の内容を投稿します。\n```plaintext\n%s\n```\n登録を取り消したい場合は次のコマンドを Scheduled Messenger に送信して下さい。\n`!delete -i %s`\n登録したメッセージを編集したい場合は次の Prefix を使って下さい。\n`!edit -i %s`",
		parsedTime.Format("2006年01月02日 15時04分"),
		body,
		id.String(),
		id.String(),
	)

	return mes
}

// 定期投稿メッセージ作成時 / 編集時のメッセージを生成
func CreateSchMesPeriodicCreatedEditedMessage(parsedTime model.PeriodicTime, distChannel *string, body string, id uuid.UUID, repeat *int) string {
	// 空のメッセージを作成
	var mes string

	// チャンネルを追加
	if distChannel != nil {
		mes += fmt.Sprintf("%s に", *distChannel)
	} else {
		mes += "以前の登録と同じチャンネルに"
	}

	// チャンネルと時間を追加
	mes += fmt.Sprintf("`%s`、", parsedTime.Format())

	// リピートの設定がある場合追加
	if repeat != nil {
		mes += fmt.Sprintf("`%d 回`", *repeat)
	}

	// 残りの文字列を追加
	mes += fmt.Sprintf(
		"以下の内容を投稿します。\n```plaintext\n%s\n```\n登録を取り消したい場合は次のコマンドを Scheduled Messenger に送信して下さい。\n`!delete -i %s`\n登録したメッセージを編集したい場合は次の Prefix を使って下さい。\n`!edit -i %s`",
		body,
		id.String(),
		id.String(),
	)

	return mes
}

// 登録メッセージ削除時のメッセージを生成
func CreateSchMesDeletedMessage(id string) string {
	return fmt.Sprintf("ID:`%s`のメッセージ送信登録を取り消しました。", id)
}

// 登録メッセージリストを生成
func CreateScheduleListMessage(mesList []*model.SchMes, mesListPeriodic []*model.SchMesPeriodic) string {
	var result string

	// 予約投稿
	result += "### 予約投稿メッセージ\n"
	// メッセージがない場合はその旨を伝える
	if len(mesList) == 0 {
		result += "登録済みのメッセージはありません。"
	} else {
		// ヘッダー
		result += "|メッセージID|投稿時刻|投稿先チャンネルID|本文|\n|----|----|----|----|"

		// メッセージごとに行を追加
		for _, mes := range mesList {
			// 改行記号を string として表示できるよう変換
			replacedBody := strings.Replace(mes.Body, "\n", "`\\n`", -1)

			// テーブル内からメンションが飛ばないように "@" を変換
			replacedBody = strings.Replace(replacedBody, "@", "`@`", -1)

			result += fmt.Sprintf("\n|%s|%s|%s|%s|", mes.ID, mes.Time.Format("2006年01月02日 15:04"), mes.ChannelID, replacedBody)
		}
	}

	// 定期投稿
	result += "\n### 定期投稿メッセージ\n"
	// メッセージがない場合はその旨を伝える
	if len(mesListPeriodic) == 0 {
		result += "登録済みのメッセージはありません。"
	} else {
		// ヘッダー
		result += "|メッセージID|投稿時刻|残り投稿回数|投稿先チャンネルID|本文|\n|----|----|----|----|----|"

		// メッセージごとに行を追加
		for _, mes := range mesListPeriodic {
			// 改行記号を string として表示できるよう変換
			replacedBody := strings.Replace(mes.Body, "\n", "`\\n`", -1)

			// テーブル内からメンションが飛ばないように "@" を変換
			replacedBody = strings.Replace(replacedBody, "@", "`@`", -1)

			result += fmt.Sprintf("\n|%s|%s|%s|%s|%s|", mes.ID, mes.Time.Format(), formatRepeat(mes.Repeat), mes.ChannelID, replacedBody)
		}
	}

	return result
}

// リピート回数のリスト用変換
func formatRepeat(repeat *int) string {
	if repeat == nil {
		return "∞"
	} else {
		return strconv.Itoa(*repeat)
	}
}

// DB のレコードから実際に送るメッセージを生成
func CreateScheduledMessage(userID string, body string) string {
	return fmt.Sprintf("#### *@%s さんからのメッセージ*\n---\n%s", userID, body)
}
