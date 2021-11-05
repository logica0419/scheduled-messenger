package model

import (
	"time"

	"github.com/google/uuid"
)

type ScheduledMes struct {
	Id      uuid.UUID
	User    string
	Time    time.Time
	Channel string
	Body    string
}
