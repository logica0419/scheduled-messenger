package api

import (
	"fmt"
	"log"
)

type Message struct {
	Content string `json:"content,omitempty"`
	Embed   bool   `json:"embed,omitempty"`
}

// 指定されたチャンネルに指定されたメッセージを投稿
func (api *API) SendMessage(chanID string, message string) error {
	// 開発モードではコンソールにメッセージを表示するのみ
	if api.config.Dev_Mode {
		log.Printf("Sending %s to %s", message, chanID)
		return nil
	} else {
		// URL を生成
		url := fmt.Sprintf("%s/channels/%s/messages", baseUrl, chanID)

		// ボディを作成
		body := Message{Content: message, Embed: false}

		// リクエストを送信
		err := api.Post(url, body)
		if err != nil {
			return err
		}

		return nil
	}
}
