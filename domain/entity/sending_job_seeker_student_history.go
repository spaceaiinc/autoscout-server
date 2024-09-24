package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerStudentHistory struct {
	ID                 uint      `db:"id" json:"id"`
	SendingJobSeekerID uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	SchoolCategory     null.Int  `db:"school_category" json:"school_category"`
	SchoolName         string    `db:"school_name" json:"school_name"`
	SchoolLevel        null.Int  `db:"school_level" json:"school_level"`
	Subject            string    `db:"subject" json:"subject"`
	EntranceYear       string    `db:"entrance_year" json:"entrance_year"`
	FirstStatus        null.Int  `db:"first_status" json:"first_status"`
	GraduationYear     string    `db:"graduation_year" json:"graduation_year"`
	LastStatus         null.Int  `db:"last_status" json:"last_status"`
	CreatedAt          time.Time `db:"created_at" json:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerStudentHistory(
	sendingJobSeekerID uint,
	schoolCategory null.Int,
	schoolName string,
	schoolLevel null.Int,
	subject string,
	entranceYear string,
	firstStatus null.Int,
	graduationYear string,
	lastStatus null.Int,
) *SendingJobSeekerStudentHistory {
	return &SendingJobSeekerStudentHistory{
		SendingJobSeekerID: sendingJobSeekerID,
		SchoolCategory:     schoolCategory,
		SchoolName:         schoolName,
		SchoolLevel:        schoolLevel,
		Subject:            subject,
		EntranceYear:       entranceYear,
		FirstStatus:        firstStatus,
		GraduationYear:     graduationYear,
		LastStatus:         lastStatus,
	}
}
