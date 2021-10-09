package api

import (
	"encoding/json"
	"fmt"
	"io"
)

type Message struct {
	Content string `json:"content,omitempty"`
	Embed   bool   `json:"embed,omitempty"`
}

// 指定されたチャンネルに指定されたメッセージを投稿
func SendMessage(chanID string, message string) error {
	// URL を生成
	url := fmt.Sprintf("%s/channels/%s/messages", api, chanID)

	// ボディを作り、io.Pipe でリーダーと繋いで Goroutine 内でエンコード
	body := Message{Content: message, Embed: false}
	pr, pw := io.Pipe()
	go func() {
		_ = json.NewEncoder(pw).Encode(&body)
		defer pw.Close()
	}()

	// io.Pipe で受け取ったボディを載せて POST リクエスト
	_, err := client.Post(url, "application/json", pr)
	if err != nil {
		return err
	}

	return nil
}
