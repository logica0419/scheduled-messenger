package timer

import (
	"github.com/logica0419/scheduled-messenger-bot/config"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
	"github.com/robfig/cron/v3"
)

// メッセージ送信タイマー構造体
type Timer struct {
	cron *cron.Cron
	c    *config.Config
	api  *api.API
	repo repository.Repository
}

// メッセージ送信タイマーのセットアップと取得
func Setup(c *config.Config, api *api.API, repo repository.Repository) (*Timer, error) {
	// cron インスタンスを取得
	cron := cron.New()

	// タイマー構造体を作成
	t := &Timer{cron: cron, c: c, api: api, repo: repo}

	// 定期実行関数を追加
	err := t.addFuncs()
	if err != nil {
		return nil, err
	}

	return t, nil
}

// 定期実行関数をタイマーに追加
func (t *Timer) addFuncs() error {
	// 予約投稿ハンドラを追加
	_, err := t.cron.AddFunc("* * * * *", t.schMesHandler)
	if err != nil {
		return err
	}

	// 定期投稿ハンドラを追加
	_, err = t.cron.AddFunc("* * * * *", t.schMesPeriodicHandler)
	if err != nil {
		return err
	}

	return nil
}

// タイマーをスタート
func (t *Timer) Start() {
	t.cron.Start()
}
