package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Content string `json:"content,omitempty"`
	Embed   bool   `json:"embed,omitempty"`
}

// 指定されたチャンネルに指定されたメッセージを投稿
func SendMessage(chanID string, message string) error {
	// URL を生成
	url := fmt.Sprintf("%s/channels/%s/messages", baseUrl, chanID)

	// ボディを作り、io.Pipe でリーダーと繋いで Goroutine 内でエンコード
	body := Message{Content: message, Embed: false}
	pr, pw := io.Pipe()
	go func() {
		_ = json.NewEncoder(pw).Encode(&body)
		defer pw.Close()
	}()

	// io.Pipe で受け取ったボディを載せて POST リクエストを作成
	req, err := http.NewRequest("POST", url, pr)
	if err != nil {
		return err
	}

	// ヘッダーにトークンを設定
	api.setTokenHeader(req)

	// リクエストを送信
	_, err = api.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
