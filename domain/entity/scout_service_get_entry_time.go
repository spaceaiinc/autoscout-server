package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type ScoutServiceGetEntryTime struct {
	ID             uint      `db:"id" json:"id"`
	ScoutServiceID uint      `db:"scout_service_id" json:"scout_service_id"`
	StartHour      null.Int  `db:"start_hour" json:"start_hour"`     // スカウト開始時間(媒体共通/0:0時, 1:1時, 2:2時, ..., 23:23時)
	StartMinute    null.Int  `db:"start_minute" json:"start_minute"` // スカウト開始分(媒体共通/0:0分, 1:1分, 2:2分, ..., 59:59分)
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

func NewScoutServiceGetEntryTime(
	scoutServiceID uint,
	startHour null.Int,
	startMinute null.Int,
) *ScoutServiceGetEntryTime {
	return &ScoutServiceGetEntryTime{
		ScoutServiceID: scoutServiceID,
		StartHour:      startHour,
		StartMinute:    startMinute,
	}
}
