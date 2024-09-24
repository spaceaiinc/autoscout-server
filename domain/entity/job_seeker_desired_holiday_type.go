package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerDesiredHolidayType struct {
	ID          uint      `db:"id" json:"id"`
	JobSeekerID uint      `db:"job_seeker_id" json:"job_seeker_id"`
	HolidayType null.Int  `db:"holiday_type" json:"holiday_type"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerDesiredHolidayType(
	jobSeekerID uint,
	holidayType null.Int,
) *JobSeekerDesiredHolidayType {
	return &JobSeekerDesiredHolidayType{
		JobSeekerID: jobSeekerID,
		HolidayType: holidayType,
	}
}
