package parser

import (
	"fmt"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
)

// 予約作成コマンドから要素を抽出
func argparseEditCommand(command []string) (*string, *string, *string, *string, *int, error) {
	// パーサーを定義
	parser := argparse.NewParser("edit", "メッセージを編集する")

	// argumentを定義
	id := parser.String("i", "id", &argparse.Options{Required: true, Help: "編集するメッセージのID"})
	channel := parser.String("c", "channel", &argparse.Options{Help: "メッセージを送るチャンネル \"#\"からフルパスを記述してください 省略した場合は予約メッセージを送信したチャンネルに送ります"})
	postTime := parser.String("t", "time", &argparse.Options{
		Help: "メッセージを送る時間 予約投稿フォーマット:`yyyy/mm/dd/hh:mm` 定期投稿フォーマット:`yyyy/mm/dd/hh:mm/d(曜日)` 曜日はオプションです。曜日は0が日曜、6が土曜に対応する1桁の整数で記入して下さい"})
	body := parser.String("b", "body", &argparse.Options{
		Help: "送るメッセージ スペースが入るときは\"\"や''でくくって下さい 予約メッセージではメンションせずにbodyにメンションを入れたい場合、\"@.{id}\"と打てば自動で\"@{id}\"に変換されます 改行したい場合は、\"\"や''でくくった上で改行したい箇所に \\n を挿入して下さい"})
	repeat := parser.Int("r", "repeat", &argparse.Options{Help: "(定期投稿のみ) 投稿を繰り返す回数 1以上の回数を指定してください -1で繰り返しを取り消します"})
	parser.DisableHelp()

	// パース
	err := parser.Parse(command)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf(parser.Usage(err))
	}

	// 不可欠な argument 以外がデフォルト値の場合、 nil に変換
	if *channel == "" {
		channel = nil
	}
	if *postTime == "" {
		postTime = nil
	}
	if *repeat == 0 {
		repeat = nil
	}

	var parsedBody *string
	// ボディが空でなければパース
	if *body != "" {
		parsedBody = bodyParse(body)
	}

	if channel == nil && postTime == nil && repeat == nil && parsedBody == nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("変更がありません\n")
	}

	return id, postTime, channel, parsedBody, repeat, nil
}

// 予約作成コマンドをパース
func ParseEditCommand(req *event.MessageEvent) (*string, *string, *string, *string, *string, *int, error) {
	// メッセージを配列に
	mesList, err := argvParse(req.GetText())
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to parse argv: %s", err)
	}

	// 冒頭にメンションがついていた場合要素をドロップ
	if strings.Contains(mesList[0], "@") {
		mesList = mesList[1:]
	}

	// メッセージをパース
	id, parsedTime, distChannel, body, repeat, err := argparseEditCommand(mesList)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	var distChannelID *string
	// チャンネルが指定されている場合、distChannel の ID を embedded から取得
	if distChannel != nil {
		for _, v := range req.GetEmbeddedList() {
			if v.Raw == *distChannel && v.Type == "channel" {
				distChannelID = &v.ID
				break
			}
		}

		// embedded に ID が見つからなかった場合、エラーメッセージを送信
		if distChannelID == nil {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("無効なチャンネルです\n")
		}
	}

	return id, parsedTime, distChannel, distChannelID, body, repeat, err
}
