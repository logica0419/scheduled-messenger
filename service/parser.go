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
	parser := argparse.NewParser("schedule", "Create scheduled message")

	// argumentを定義
	channel := parser.String("c", "channel", &argparse.Options{Required: true, Help: "Channel to send message to, beginning with`#`"})
	postTime := parser.String("t", "time", &argparse.Options{Required: true, Help: "Time to send message at, format:`yy/mm/dd/hh:mm`"})
	body := parser.String("b", "body", &argparse.Options{Required: true, Help: "Message to send"})

	// パース
	err := parser.Parse(message)
	if err != nil {
		return time.Now(), "", "", fmt.Errorf("```plaintext\n%s```", parser.Usage(err))
	}

	// 時間をパース
	parsedTime, err := timeParse(*postTime)
	if err != nil {
		return time.Now(), "", "", fmt.Errorf("```plaintext\n%s```", parser.Usage(err))
	}

	// 指定された時間が現在時刻より後か確認する
	if time.Now().After(parsedTime) {
		return time.Now(), "", "", fmt.Errorf("```plaintext\n%s```", parser.Usage("Error: invalid time - Specify the time later than now."))
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
