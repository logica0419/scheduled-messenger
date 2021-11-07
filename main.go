package main

import (
	"log"
	"time"

	"github.com/logica0419/scheduled-messenger-bot/config"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/router"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
	"github.com/logica0419/scheduled-messenger-bot/timer"
)

func init() {
	// タイムゾーンの設定
	const location = "Asia/Tokyo"

	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}

	time.Local = loc
}

func main() {
	// ログを出力
	log.Print("Initializing Scheduled Messenger Bot...")

	// 設定を読み込む
	c, err := config.GetConfig()
	if err != nil {
		log.Panicf("Error: failed to get config - %s", err)
	}

	// API クライアントを取得
	api := api.GetApi(c)

	// リポジトリを取得
	repo, err := repository.GetRepository(c)
	if err != nil {
		log.Panicf("Error: failed to initialize DB - %s", err)
	}

	// メッセージ送信タイマーをセットアップ
	t, err := timer.Setup(c, api, repo)
	if err != nil {
		log.Panicf("Error: failed to initialize mes-sending timer - %s", err)
	}

	// メッセージ送信タイマーをスタート
	t.Start()

	// ルーターをセットアップ
	r := router.Setup(c, api, repo)

	// ルーターをスタート
	r.Start()
}
