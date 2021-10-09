package api

import (
	"fmt"
	"net/http"

	"github.com/logica0419/scheduled-messenger-bot/config"
)

// 共通 API ベース URL
const baseUrl = "https://q.trap.jp/api/v3"

// API クライアント構造体
type API struct {
	client *http.Client
	config *config.Config
}

// 共通 API クライアント
var api *API

func SetUpApi() {
	client := new(http.Client)
	api = &API{client: client, config: config.C}
}

// リクエストのヘッダにトークンを付与
func (api *API) setTokenHeader(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.config.Bot_Access_Token))
}

// リクエストのヘッダに JSON の Content-Type を付与
func (api *API) setJsonHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
}
