package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type JobSeekerInterestedJobListing struct {
	ID               uint      `db:"id" json:"id"`
	UUID             uuid.UUID `db:"uuid" json:"uuid"`
	JobSeekerID      uint      `db:"job_seeker_id" json:"job_seeker_id"`
	JobInformationID uint      `db:"job_information_id" json:"job_information_id"`
	InterestedType   null.Int  `db:"interested_type" json:"interested_type"` // 興味ありの種別(0: エントリー希望, 1: 興味あり)
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerInterestedJobListing(
	jobSeekerID uint,
	jobInformationID uint,
	interestedType null.Int,
) *JobSeekerInterestedJobListing {
	return &JobSeekerInterestedJobListing{
		JobSeekerID:      jobSeekerID,
		JobInformationID: jobInformationID,
		InterestedType:   interestedType,
	}
}

type InterestedType int64

const (
	InterestedTypeEntry      InterestedType = iota // エントリー希望
	InterestedTypeInterested                       // 興味あり
)
