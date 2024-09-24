package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerDesiredCompanyScale struct {
	ID                  uint     `db:"id" json:"id"`
	JobSeekerID         uint     `db:"job_seeker_id" json:"job_seeker_id"`
	DesiredCompanyScale null.Int `db:"desired_company_scale" json:"desired_company_scale"`
	// DesiredRank         null.Int  `db:"desired_rank" json:"desired_rank"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerDesiredCompanyScale(
	jobSeekerID uint,
	desiredCompanyScale null.Int,
	// desiredRank null.Int,
) *JobSeekerDesiredCompanyScale {
	return &JobSeekerDesiredCompanyScale{
		JobSeekerID:         jobSeekerID,
		DesiredCompanyScale: desiredCompanyScale,
		// DesiredRank:         desiredRank,
	}
}
