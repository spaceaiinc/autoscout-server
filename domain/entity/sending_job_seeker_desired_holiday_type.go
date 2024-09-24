package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerDesiredHolidayType struct {
	ID                 uint      `db:"id" json:"id"`
	SendingJobSeekerID uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	HolidayType        null.Int  `db:"holiday_type" json:"holiday_type"`
	CreatedAt          time.Time `db:"created_at" json:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerDesiredHolidayType(
	sendingJobSeekerID uint,
	holidayType null.Int,
) *SendingJobSeekerDesiredHolidayType {
	return &SendingJobSeekerDesiredHolidayType{
		SendingJobSeekerID: sendingJobSeekerID,
		HolidayType:        holidayType,
	}
}
