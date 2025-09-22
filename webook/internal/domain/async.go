package domain

import (
	"time"
)

type FailedSMS struct {
	ID     int64
	Biz    string
	Args   []string
	Number string
	Status string

	Ctime      time.Time
	UpdateTime time.Time
}
