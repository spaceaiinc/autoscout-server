package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationRequiredSocialExperience struct {
	ID                   uint      `db:"id" json:"id"`
	JobInformationID     uint      `db:"job_information_id" json:"job_information_id"`
	SocialExperienceType null.Int  `db:"social_experience_type" json:"social_experience_type"`
	CreatedAt            time.Time `db:"created_at" json:"-"`
	UpdatedAt            time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationRequiredSocialExperience(
	jobInformationID uint,
	socialExperienceType null.Int,
) *JobInformationRequiredSocialExperience {
	return &JobInformationRequiredSocialExperience{
		JobInformationID:     jobInformationID,
		SocialExperienceType: socialExperienceType,
	}
}
