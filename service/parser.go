package service

import (
	"fmt"
	"log"
	"time"

	"github.com/akamensky/argparse"
	"github.com/cosiner/argv"
)

// プレーンテキストのメッセージを配列に分解
func ArgvParse(message string) []string {
	// パース用関数
	var identity = func(s string) (string, error) { return s, nil }

	// パース
	parsed, err := argv.Argv(message, identity, identity)
	if err != nil || len(parsed) == 0 {
		log.Printf("an error occurred while parsing user command: %s\n", err)
		return nil
	}

	return parsed[0]
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
		return time.Now(), "", "", err
	}
	log.Print("```plaintext\n" + parser.Usage(err) + "```")

	// 時間をパース
	parsedTime, err := timeParse(*postTime)
	if err != nil {
		return time.Now(), "", "", err
	}

	// 指定された時間が現在時刻より後か確認する
	if time.Now().After(parsedTime) {
		return time.Now(), "", "", fmt.Errorf("Error: invalid time: %s", parsedTime)
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
