package api

import (
	"fmt"
	"log"
)

// メッセージ投稿リクエストボディ
type Message struct {
	Content string `json:"content,omitempty"`
	Embed   bool   `json:"embed,omitempty"`
}

// 指定されたチャンネルに指定されたメッセージを投稿
func (api *API) SendMessage(chanID string, message string) error {
	// 開発モードではコンソールにメッセージを表示するのみ
	if api.config.Dev_Mode {
		log.Printf("Sending\n%s\nto %s", message, chanID)
		return nil
	} else {
		// URL を生成
		url := fmt.Sprintf("%s/channels/%s/messages", baseUrl, chanID)

		// ボディを作成
		body := Message{Content: message, Embed: true}

		// リクエストを送信
		err := api.post(url, body)
		if err != nil {
			return err
		}

		return nil
	}
}

// デプロイ完了を config で設定したチャンネルに通知
func (api *API) NotifyDeployed() {
	_ = api.SendMessage(api.config.Log_Chan_ID, "Log: The new version of Scheduled Messenger is deployed.")
}
