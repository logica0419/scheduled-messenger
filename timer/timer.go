package timer

import (
	"github.com/logica0419/scheduled-messenger-bot/config"
	"github.com/logica0419/scheduled-messenger-bot/repository"
	"github.com/logica0419/scheduled-messenger-bot/service/api"
	"github.com/robfig/cron/v3"
)

type timerFunc struct {
	schedule string
	handler  func()
}

type Timer struct {
	cron *cron.Cron
	c    *config.Config
	api  *api.API
	repo repository.Repository
}

func Setup(c *config.Config, api *api.API, repo repository.Repository) (*Timer, error) {
	cron := cron.New()

	t := &Timer{cron: cron, c: c, api: api, repo: repo}

	err := t.addFuncs()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Timer) addFuncs() error {
	var timerFuncs = []timerFunc{
		{schedule: "* * * * *", handler: t.normalMesHandler},
	}

	for _, v := range timerFuncs {
		_, err := t.cron.AddFunc(v.schedule, v.handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Timer) Start() {
	t.cron.Start()
}
