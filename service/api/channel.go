package api

import (
	"fmt"
)

type ActionBody struct {
	ChannelID string `json:"channelId,omitempty"`
}

func (api *API) ChannelAction(cmd string, chanID string) error {
	// URL を生成
	url := fmt.Sprintf("%s/bots/%s/actions/%s", baseUrl, api.config.Bot_ID, cmd)

	// ボディを作成
	body := ActionBody{ChannelID: chanID}

	// リクエストを送信
	err := api.Post(url, body)
	if err != nil {
		return err
	}

	return nil
}
