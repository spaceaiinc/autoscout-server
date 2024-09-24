package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerDesiredIndustry struct {
	ID              uint      `db:"id" json:"id"`
	JobSeekerID     uint      `db:"job_seeker_id" json:"job_seeker_id"`
	DesiredIndustry null.Int  `db:"desired_industry" json:"desired_industry"`
	DesiredRank     null.Int  `db:"desired_rank" json:"desired_rank"`
	CreatedAt       time.Time `db:"created_at" json:"-"`
	UpdatedAt       time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerDesiredIndustry(
	jobSeekerID uint,
	desiredIndustry null.Int,
	desiredRank null.Int,
) *JobSeekerDesiredIndustry {
	return &JobSeekerDesiredIndustry{
		JobSeekerID:     jobSeekerID,
		DesiredIndustry: desiredIndustry,
		DesiredRank:     desiredRank,
	}
}
