package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/logica0419/scheduled-messenger-bot/config"
)

const (
	botTokenHeader    = "X-TRAQ-BOT-TOKEN"
	botEventHeader    = "X-TRAQ-BOT-EVENT"
	contentTypeHeader = "Content-Type"
)

type Router struct {
	e      *echo.Echo
	Config *config.Config
}

type errorMessage struct {
	Message string `json:"message,omitempty"`
}

var r *Router

func Setup() *Router {
	r = newRouter()

	r.e.POST("/", botEventHandler, requestVerification)

	return r
}

func newEcho() *echo.Echo {
	e := echo.New()

	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetHeader("${time_rfc3339} ${prefix} ${short_file} ${line} |")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "${time_rfc3339} method = ${method} | uri = ${uri} | status = ${status} ${error}\n"}))

	return e
}

func newRouter() *Router {
	e := newEcho()

	r := Router{e: e, Config: config.C}

	return &r
}

func (r *Router) Start() {
	r.e.Logger.Panic(r.e.Start(":8080"))
}
