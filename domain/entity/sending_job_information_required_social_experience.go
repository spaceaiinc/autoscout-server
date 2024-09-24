package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationRequiredSocialExperience struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobInformationID uint      `db:"sending_job_information_id" json:"sending_job_information_id"`
	SocialExperienceType    null.Int  `db:"social_experience_type" json:"social_experience_type"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationRequiredSocialExperience(
	sendingJobInformationID uint,
	socialExperienceType null.Int,
) *SendingJobInformationRequiredSocialExperience {
	return &SendingJobInformationRequiredSocialExperience{
		SendingJobInformationID: sendingJobInformationID,
		SocialExperienceType:    socialExperienceType,
	}
}
