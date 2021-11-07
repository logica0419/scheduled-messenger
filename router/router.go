package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/logica0419/scheduled-messenger-bot/config"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
)

const (
	botTokenHeader    = "X-TRAQ-BOT-TOKEN"
	botEventHeader    = "X-TRAQ-BOT-EVENT"
	contentTypeHeader = "Content-Type"
)

// ルーター
type Router struct {
	e      *echo.Echo
	Config *config.Config
	Api    *api.API
	Repo   repository.Repository
}

// エラーレスポンスのボディ構造体
type errorMessage struct {
	Message string `json:"message,omitempty"`
}

// ルーターのセットアップと取得
func Setup(c *config.Config, api *api.API, repo repository.Repository) *Router {
	// ルーターを取得
	r := newRouter(c, api, repo)

	// ハンドラを追加
	r.e.POST("/", r.botEventHandler, r.requestVerification)

	return r
}

// 新しい Echo インスタンスを取得
func newEcho() *echo.Echo {
	e := echo.New()

	// ログの設定
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetHeader("${time_rfc3339} ${prefix} ${short_file} ${line} |")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "${time_rfc3339} method = ${method} | uri = ${uri} | status = ${status} ${error}\n"}))

	return e
}

// 新しいルーターを取得
func newRouter(c *config.Config, api *api.API, repo repository.Repository) *Router {
	// Echo インスタンスを取得
	e := newEcho()

	// ルーターを作る
	r := Router{e: e, Config: c, Api: api, Repo: repo}

	return &r
}

// ルーターを起動
func (r *Router) Start() {
	r.e.Logger.Panic(r.e.Start(":8080"))
}
