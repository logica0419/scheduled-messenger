package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	Content string `json:"content,omitempty"`
	Embed   bool   `json:"embed,omitempty"`
}

// 指定されたチャンネルに指定されたメッセージを投稿
func (api *API) SendMessage(chanID string, message string) error {
	// URL を生成
	url := fmt.Sprintf("%s/channels/%s/messages", baseUrl, chanID)

	// ボディを作り、バイト列に変換
	body := Message{Content: message, Embed: false}
	byteBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// 変換したボディを載せて POST リクエストを作成
	req, err := http.NewRequest("POST", url, bytes.NewReader(byteBody))
	if err != nil {
		return err
	}

	// ヘッダーを設定
	setTokenHeader(req, api)
	setJsonHeader(req)

	// リクエストを送信
	res, err := api.client.Do(req)
	log.Println(*res)
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 {
		return fmt.Errorf(res.Status)
	}

	return nil
}
