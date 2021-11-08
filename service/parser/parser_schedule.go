package parser

import (
	"fmt"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/cosiner/argv"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

// プレーンテキストのメッセージを配列に分解
func argvParse(message string) ([]string, error) {
	// パース用関数を定義
	var identity = func(s string) (string, error) { return s, nil }

	// パース
	parsed, err := argv.Argv(message, identity, identity)
	if err != nil || len(parsed) == 0 {
		return nil, err
	}

	return parsed[0], nil
}

// 予約作成コマンドから要素を抽出
func argparseScheduleCommand(command []string) (*string, *string, *string, *int, error) {
	// パーサーを定義
	parser := argparse.NewParser("schedule", "メッセージを予約する")

	// argumentを定義
	channel := parser.String("c", "channel", &argparse.Options{Help: "メッセージを送るチャンネル \"#\"からフルパスを記述してください 省略した場合は予約メッセージを送信したチャンネルに送ります"})
	postTime := parser.String("t", "time",
		&argparse.Options{Required: true, Help: "メッセージを送る時間 フォーマット:`yyyy/mm/dd/hh:mm`"})
	body := parser.String("b", "body",
		&argparse.Options{Required: true,
			Help: "送るメッセージ スペースが入るときは\"\"や''でくくって下さい 予約メッセージではメンションせずにbodyにメンションを入れたい場合、\"@.{id}\"と打てば自動で\"@{id}\"に変換されます 改行したい場合は、\"\"や''でくくった上で改行したい箇所に \\n を挿入して下さい"})
	repeat := parser.Int("r", "repeat", &argparse.Options{Help: "(定期投稿のみ) 投稿を繰り返す回数"})
	parser.DisableHelp()

	// パース
	err := parser.Parse(command)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("メッセージの予約に失敗しました\n```plaintext\n%s```", parser.Usage(err))
	}

	// ボディをパース
	parsedBody := bodyParse(body)

	return postTime, channel, parsedBody, repeat, nil
}

// body から特定のルールにマッチする文字列を変換
func bodyParse(body *string) *string {
	// メンションよけのパース (@.{id} を @{id} に変換)
	replacedBody := strings.Replace(*body, "@.", "@", -1)

	return &replacedBody
}

// 予約作成コマンドをパース
func ParseScheduleCommand(api *api.API, req *event.MessageEvent) (*string, *string, *string, *string, *int, error) {
	// メッセージを配列に
	listedReqMes, err := argvParse(req.GetText())
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to parse argv: %s", err)
	}

	// 冒頭にメンションがついていた場合要素をドロップ
	if strings.Contains(listedReqMes[0], "@") {
		listedReqMes = listedReqMes[1:]
	}

	// メッセージをパース
	parsedTime, distChannel, body, repeat, err := argparseScheduleCommand(listedReqMes)
	if err != nil {
		_ = api.SendMessage(req.GetChannelID(), err.Error())
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to parse argv: %s", err)
	}

	// チャンネルが指定されてないとき、ID と名前を取得
	var distChannelID *string
	if distChannel == nil {
		*distChannel = "このチャンネル"
		*distChannelID = req.GetChannelID()
	} else {
		// 指定されている場合 distChannel の ID を embedded から取得
		for _, v := range req.GetEmbeddedList() {
			if v.Raw == *distChannel && v.Type == "channel" {
				*distChannelID = v.ID
				break
			}
		}

		// embedded に ID が見つからなかった場合、エラーメッセージを送信
		if distChannelID == nil {
			_ = api.SendMessage(req.GetChannelID(), "メッセージの予約に失敗しました\n```plaintext\n無効なチャンネルです\n```")
			return nil, nil, nil, nil, nil, fmt.Errorf("failed to get channel ID: %s", err)
		}
	}

	return parsedTime, distChannel, distChannelID, body, repeat, err
}
