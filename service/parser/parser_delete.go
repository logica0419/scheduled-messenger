package parser

import (
	"fmt"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/logica0419/scheduled-messenger-bot/model/event"
)

// 予約削除コマンドから ID を抽出
func argparseDeleteCommand(command []string) (*string, error) {
	// パーサーを定義
	parser := argparse.NewParser("delete", "メッセージを削除する")

	// argumentを定義
	id := parser.String("i", "id", &argparse.Options{Required: true, Help: "削除したいメッセージのUUID"})
	parser.DisableHelp()

	// パース
	err := parser.Parse(command)
	if err != nil {
		return nil, fmt.Errorf("メッセージの削除に失敗しました\n```plaintext\n%s```", parser.Usage(err))
	}

	return id, nil
}

// 予約削除コマンドをパース
func ParseDeleteCommand(req *event.MessageEvent) (*string, error) {
	// メッセージを配列に
	listedReqMes, err := argvParse(req.GetText())
	if err != nil {
		return nil, fmt.Errorf("failed to parse argv: %s", err)
	}

	// 冒頭にメンションがついていた場合要素をドロップ
	if strings.Contains(listedReqMes[0], "@") {
		listedReqMes = listedReqMes[1:]
	}

	// メッセージをパース
	id, err := argparseDeleteCommand(listedReqMes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse argv: %s", err)
	}

	return id, err
}
