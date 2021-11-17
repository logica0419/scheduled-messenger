package repository

import (
	"github.com/logica0419/scheduled-messenger-bot/model"
	"github.com/patrickmn/go-cache"
)

const schMesPeriodicCacheID = "sch_mes_periodic_all"

func (repo *GormRepository) setSchMesPeriodicCache(content []*model.SchMesPeriodic) {
	repo.c.Set(schMesPeriodicCacheID, content, cache.NoExpiration)
}

func (repo *GormRepository) getSchMesPeriodicCache() []*model.SchMesPeriodic {
	content, found := repo.c.Get(schMesPeriodicCacheID)
	if !found {
		return nil
	}

	return content.([]*model.SchMesPeriodic)
}

func (repo *GormRepository) deleteSchMesPeriodicCache() {
	repo.c.Delete(schMesPeriodicCacheID)
}
