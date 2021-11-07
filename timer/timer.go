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

// 定期実行関数構造体
type timerFunc struct {
	schedule string
	handler  func()
}

// 定期実行関数リスト
var timerFuncs []timerFunc

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
	// 通常の予約投稿ハンドラを追加
	timerFuncs = append(timerFuncs, timerFunc{schedule: "* * * * *", handler: t.normalMesHandler})

	// timerFuncs に登録された関数を全て追加
	for _, v := range timerFuncs {
		_, err := t.cron.AddFunc(v.schedule, v.handler)
		if err != nil {
			return err
		}
	}

	return nil
}

// タイマーをスタート
func (t *Timer) Start() {
	t.cron.Start()
}
