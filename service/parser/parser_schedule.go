package parser

import (
	"fmt"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
)

// 予約作成コマンドから要素を抽出
func argparseScheduleCommand(command []string) (*string, *string, *string, *int, error) {
	// パーサーを定義
	parser := argparse.NewParser("schedule", "メッセージを登録する")

	// argumentを定義
	channel := parser.String("c", "channel", &argparse.Options{Help: "メッセージを送るチャンネル \"#\"からフルパスを記述してください 省略した場合は予約メッセージを送信したチャンネルに送ります"})
	postTime := parser.String("t", "time", &argparse.Options{Required: true,
		Help: "メッセージを送る時間 予約投稿フォーマット:`yyyy/mm/dd/hh:mm` 定期投稿フォーマット:`yyyy/mm/dd/hh:mm/d(曜日)` 曜日(オプション)は0が日曜、6が土曜に対応する1桁の整数で記入して下さい 曜日を&で繋げると複数曜日の同じ時間を同時指定できます"})
	body := parser.String("b", "body", &argparse.Options{Required: true,
		Help: "送るメッセージ スペースが入るときは\"\"や''でくくって下さい 予約メッセージではメンションせずにbodyにメンションを入れたい場合、 \"@.{id}\"と打てば自動で\"@{id}\"に変換されます 改行したい場合は\"\"や''でくくった上で、改行したい箇所に\\nを挿入するか 通常通り改行して下さい"})
	repeat := parser.Int("r", "repeat", &argparse.Options{Help: "(定期投稿のみ) 投稿を繰り返す回数 1以上の回数を指定してください 0か省略で無限に繰り返します"})
	parser.DisableHelp()

	// パース
	err := parser.Parse(command)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf(parser.Usage(err))
	}

	// 不可欠な argument 以外がデフォルト値の場合、 nil に変換
	if *channel == "" {
		channel = nil
	}
	if *repeat == 0 {
		repeat = nil
	}

	// ボディをパース
	parsedBody := bodyParse(body)

	return postTime, channel, parsedBody, repeat, nil
}

// 予約作成コマンドをパース
func ParseScheduleCommand(req *event.MessageEvent) (*string, *string, *string, *string, *int, error) {
	// メッセージを配列に
	mesList, err := argvParse(req.GetText())
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to parse argv: %s", err)
	}

	// 冒頭にメンションがついていた場合要素をドロップ
	if strings.Contains(mesList[0], "@") {
		mesList = mesList[1:]
	}

	// メッセージをパース
	parsedTime, distChannel, body, repeat, err := argparseScheduleCommand(mesList)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	var distChannelID *string
	// チャンネルが指定されてない場合、ID と名前を取得
	if distChannel == nil {
		_distChannel := "このチャンネル"
		distChannel = &_distChannel
		_distChannelID := req.GetChannelID()
		distChannelID = &_distChannelID
	} else {
		// 指定されている場合、distChannel の ID を embedded から取得
		for _, v := range req.GetEmbeddedList() {
			if v.Raw == *distChannel && v.Type == "channel" {
				distChannelID = &v.ID
				break
			}
		}

		// embedded に ID が見つからなかった場合、エラーメッセージを送信
		if distChannelID == nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("無効なチャンネルです\n")
		}
	}

	return parsedTime, distChannel, distChannelID, body, repeat, err
}
