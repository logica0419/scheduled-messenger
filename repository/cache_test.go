package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/logica0419/scheduled-messenger-bot/model"
	"github.com/logica0419/scheduled-messenger-bot/service/parser"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGormRepository_setSchMesPeriodicCache(t *testing.T) {
	mockCache := cache.New(0, 10*time.Minute)
	timeStr1 := "*/12/*/05:06/6"
	time1, err := parser.TimeParsePeriodic(&timeStr1)
	if err != nil || len(time1) != 1 {
		assert.Fail(t, "failed to parse time", err)
	}
	timeStr2 := "*/01/01/*:06/6&3"
	time2, err := parser.TimeParsePeriodic(&timeStr2)
	if err != nil || len(time2) != 2 {
		assert.Fail(t, "failed to parse time", err)
	}
	repeat2 := 3

	type fields struct {
		c *cache.Cache
	}
	type args struct {
		content []*model.SchMesPeriodic
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "case1",
			fields: fields{c: mockCache},
			args:   args{content: nil},
		},
		{
			name:   "case2",
			fields: fields{c: mockCache},
			args: args{content: []*model.SchMesPeriodic{
				{ID: uuid.New(), UserID: "test1", Time: *time1[0], Repeat: nil, ChannelID: uuid.New(), Body: "test1"},
			}},
		},
		{
			name:   "case3",
			fields: fields{c: mockCache},
			args: args{content: []*model.SchMesPeriodic{
				{ID: uuid.New(), UserID: "test2", Time: *time2[0], Repeat: &repeat2, ChannelID: uuid.New(), Body: "test2"},
				{ID: uuid.New(), UserID: "test2", Time: *time2[1], Repeat: &repeat2, ChannelID: uuid.New(), Body: "test2"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &GormRepository{
				db: new(gorm.DB),
				c:  tt.fields.c,
			}
			repo.setSchMesPeriodicCache(tt.args.content)

			content, found := repo.c.Get(schMesPeriodicCacheID)
			assert.True(t, found, "failed to get cache")
			for i, v := range content.([]*model.SchMesPeriodic) {
				assert.Equal(t, *v, *tt.args.content[i], "cache doesn't match")
			}
		})
	}
}

func TestGormRepository_getSchMesPeriodicCache(t *testing.T) {
	mockCache := cache.New(0, 10*time.Minute)
	timeStr1 := "*/12/*/05:06/6"
	time1, err := parser.TimeParsePeriodic(&timeStr1)
	if err != nil || len(time1) != 1 {
		assert.Fail(t, "failed to parse time", err)
	}
	timeStr2 := "*/01/01/*:06/6&3"
	time2, err := parser.TimeParsePeriodic(&timeStr2)
	if err != nil || len(time2) != 2 {
		assert.Fail(t, "failed to parse time", err)
	}
	repeat2 := 3

	type fields struct {
		c *cache.Cache
	}
	tests := []struct {
		name   string
		fields fields
		want   []*model.SchMesPeriodic
		setup  bool
	}{
		{
			name:   "case1",
			fields: fields{c: mockCache},
			want:   nil,
			setup:  false,
		},
		{
			name:   "case2",
			fields: fields{c: mockCache},
			want:   nil,
			setup:  true,
		},
		{
			name:   "case3",
			fields: fields{c: mockCache},
			want: []*model.SchMesPeriodic{
				{ID: uuid.New(), UserID: "test1", Time: *time1[0], Repeat: nil, ChannelID: uuid.New(), Body: "test1"},
			},
			setup: true,
		},
		{
			name:   "case4",
			fields: fields{c: mockCache},
			want: []*model.SchMesPeriodic{
				{ID: uuid.New(), UserID: "test2", Time: *time2[0], Repeat: &repeat2, ChannelID: uuid.New(), Body: "test2"},
				{ID: uuid.New(), UserID: "test2", Time: *time2[1], Repeat: &repeat2, ChannelID: uuid.New(), Body: "test2"},
			},
			setup: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &GormRepository{
				db: new(gorm.DB),
				c:  tt.fields.c,
			}
			if tt.setup {
				repo.c.Set(schMesPeriodicCacheID, tt.want, cache.NoExpiration)
			}

			assert.Equal(t, tt.want, repo.getSchMesPeriodicCache())
		})
	}
}

func TestGormRepository_deleteSchMesPeriodicCache(t *testing.T) {
	mockCache := cache.New(0, 10*time.Minute)
	timeStr1 := "*/12/*/05:06/6"
	time1, err := parser.TimeParsePeriodic(&timeStr1)
	if err != nil || len(time1) != 1 {
		assert.Fail(t, "failed to parse time", err)
	}
	timeStr2 := "*/01/01/*:06/6&3"
	time2, err := parser.TimeParsePeriodic(&timeStr2)
	if err != nil || len(time2) != 2 {
		assert.Fail(t, "failed to parse time", err)
	}
	repeat2 := 3

	type fields struct {
		c *cache.Cache
	}
	tests := []struct {
		name    string
		fields  fields
		content []*model.SchMesPeriodic
	}{
		{
			name:    "case1",
			fields:  fields{c: mockCache},
			content: nil,
		},
		{
			name:   "case2",
			fields: fields{c: mockCache},
			content: []*model.SchMesPeriodic{
				{ID: uuid.New(), UserID: "test1", Time: *time1[0], Repeat: nil, ChannelID: uuid.New(), Body: "test1"},
			},
		},
		{
			name:   "case3",
			fields: fields{c: mockCache},
			content: []*model.SchMesPeriodic{
				{ID: uuid.New(), UserID: "test2", Time: *time2[0], Repeat: &repeat2, ChannelID: uuid.New(), Body: "test2"},
				{ID: uuid.New(), UserID: "test2", Time: *time2[1], Repeat: &repeat2, ChannelID: uuid.New(), Body: "test2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &GormRepository{
				db: new(gorm.DB),
				c:  tt.fields.c,
			}
			repo.c.Set(schMesPeriodicCacheID, tt.content, cache.NoExpiration)

			repo.deleteSchMesPeriodicCache()

			content, found := repo.c.Get(schMesPeriodicCacheID)
			assert.False(t, found, "cache should be deleted")
			assert.Nil(t, content, "cache should be deleted")
		})
	}
}
