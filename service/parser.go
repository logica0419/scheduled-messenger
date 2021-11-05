package service

import (
	"fmt"
	"time"

	"github.com/akamensky/argparse"
	"github.com/cosiner/argv"
)

// プレーンテキストのメッセージを配列に分解
func ArgvParse(message string) ([]string, error) {
	// パース用関数
	var identity = func(s string) (string, error) { return s, nil }

	// パース
	parsed, err := argv.Argv(message, identity, identity)
	if err != nil || len(parsed) == 0 {
		return nil, err
	}

	return parsed[0], nil
}

// メッセージから要素を抽出する
func ParseScheduleMessage(message []string) (time.Time, string, string, error) {
	// パーサーを定義
	parser := argparse.NewParser("schedule", "メッセージを予約する")

	// argumentを定義
	channel := parser.String("c", "channel", &argparse.Options{Default: "", Help: "メッセージを送るチャンネル `#`からフルパスを記述してください 省略した場合は予約メッセージを送信したチャンネルに送ります"})
	postTime := parser.String("t", "time", &argparse.Options{Required: true, Help: "メッセージを送る時間 フォーマット: `yyyy/mm/dd/hh:mm`"})
	body := parser.String("b", "body", &argparse.Options{Required: true, Help: "送るメッセージ スペースが入るときは\"\"や''でくくって下さい 改行はスペースとして扱われます"})
	parser.DisableHelp()

	// パース
	err := parser.Parse(message)
	if err != nil {
		return time.Now(), "", "", fmt.Errorf("メッセージの予約に失敗しました\n```plaintext\n%s```", parser.Usage(err))
	}

	// 時間をパース
	parsedTime, err := timeParse(*postTime)
	if err != nil {
		return time.Now(), "", "", fmt.Errorf("メッセージの予約に失敗しました\n```plaintext\n%s```", parser.Usage(err))
	}

	// 指定された時間が現在時刻より後か確認する
	if time.Now().After(parsedTime) {
		return time.Now(), "", "", fmt.Errorf("メッセージの予約に失敗しました\n```plaintext\n%s```", parser.Usage("現在時刻より後の時間を指定してください"))
	}

	return parsedTime, *channel, *body, nil
}

// 記入された時間をtime.Timeに変換
func timeParse(t string) (time.Time, error) {
	// 記入フォーマットを定義
	const format = "2006/01/02/15:04"

	// フォーマットに従ってパース
	parsed, err := time.Parse(format, t)
	if err != nil {
		return time.Now(), err
	}

	return parsed, nil
}
