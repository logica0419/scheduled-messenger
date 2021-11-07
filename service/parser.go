package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/cosiner/argv"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

// プレーンテキストのメッセージを配列に分解
func argvParse(message string) ([]string, error) {
	// パース用関数
	var identity = func(s string) (string, error) { return s, nil }

	// パース
	parsed, err := argv.Argv(message, identity, identity)
	if err != nil || len(parsed) == 0 {
		return nil, err
	}

	return parsed[0], nil
}

// 予約作成コマンドから要素を抽出
func argparseScheduleCommand(command []string) (time.Time, string, string, error) {
	// パーサーを定義
	parser := argparse.NewParser("schedule", "メッセージを予約する")

	// argumentを定義
	channel := parser.String("c", "channel", &argparse.Options{Default: "", Help: "メッセージを送るチャンネル `#`からフルパスを記述してください 省略した場合は予約メッセージを送信したチャンネルに送ります"})
	postTime := parser.String("t", "time", &argparse.Options{Required: true, Help: "メッセージを送る時間 フォーマット: `yyyy/mm/dd/hh:mm`"})
	body := parser.String("b", "body", &argparse.Options{Required: true, Help: "送るメッセージ スペースが入るときは\"\"や''でくくって下さい 改行したい場合は、したい箇所に`\\n`を挿入して下さい"})
	parser.DisableHelp()

	// パース
	err := parser.Parse(command)
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
	parsed, err := time.ParseInLocation(format, t, time.Local)
	if err != nil {
		return time.Now(), err
	}

	return parsed, nil
}

// 予約作成コマンドをパース
func ParseScheduleCommand(api *api.API, req *event.MessageEvent) (time.Time, string, string, string, error) {
	// メッセージを配列に
	listedReqMes, err := argvParse(req.GetText())
	if err != nil {
		return time.Now(), "", "", "", fmt.Errorf("failed to parse argv: %s", err)
	}

	// 冒頭にメンションがついていた場合要素をドロップ
	if strings.Contains(listedReqMes[0], "@") {
		listedReqMes = listedReqMes[1:]
	}

	// メッセージをパース
	parsedTime, distChannel, body, err := argparseScheduleCommand(listedReqMes)
	if err != nil {
		_ = api.SendMessage(req.GetChannelID(), err.Error())
		return time.Now(), "", "", "", fmt.Errorf("failed to parse argv: %s", err)
	}

	// チャンネルが指定されてないとき、ID と名前を取得
	distChannelID := ""
	if distChannel == "" {
		distChannel = "このチャンネル"
		distChannelID = req.GetChannelID()
	} else {
		// 指定されている場合 distChannel の ID を embedded から取得
		for _, v := range req.Message.Embedded {
			if v.Raw == distChannel && v.Type == "channel" {
				distChannelID = v.ID
				break
			}
		}

		// embedded に ID が見つからなかった場合、エラーメッセージを送信
		if distChannelID == "" {
			_ = api.SendMessage(req.GetChannelID(), "メッセージの予約に失敗しました\n```plaintext\n無効なチャンネルです\n```")
			return time.Now(), "", "", "", fmt.Errorf("failed to get channel ID: %s", err)
		}
	}

	return parsedTime, distChannel, distChannelID, body, err
}

// 予約削除コマンドから ID を抽出
func argparseDeleteCommand(command []string) (string, error) {
	// パーサーを定義
	parser := argparse.NewParser("delete", "メッセージを削除する")

	// argumentを定義
	id := parser.String("i", "id", &argparse.Options{Required: true, Help: "削除したいメッセージのUUID"})
	parser.DisableHelp()

	// パース
	err := parser.Parse(command)
	if err != nil {
		return "", fmt.Errorf("メッセージの削除に失敗しました\n```plaintext\n%s```", parser.Usage(err))
	}

	return *id, nil
}

// 予約削除コマンドをパース
func ParseDeleteCommand(api *api.API, req *event.MessageEvent) (string, error) {
	// メッセージを配列に
	listedReqMes, err := argvParse(req.GetText())
	if err != nil {
		return "", fmt.Errorf("failed to parse argv: %s", err)
	}

	// 冒頭にメンションがついていた場合要素をドロップ
	if strings.Contains(listedReqMes[0], "@") {
		listedReqMes = listedReqMes[1:]
	}

	// メッセージをパース
	id, err := argparseDeleteCommand(listedReqMes)
	if err != nil {
		_ = api.SendMessage(req.GetChannelID(), err.Error())
		return "", fmt.Errorf("failed to parse argv: %s", err)
	}

	return id, err
}
