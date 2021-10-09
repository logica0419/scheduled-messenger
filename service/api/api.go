package api

import "net/http"

// 共通 API ベース URL
const api = "https://q.trap.jp/api/v3/"

// 共通 HTTP クライアント
var client = new(http.Client)
