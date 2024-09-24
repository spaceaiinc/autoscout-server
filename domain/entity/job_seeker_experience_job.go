package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// 求職者の経験職種（LPからの登録時に使用）
type JobSeekerExperienceJob struct {
	ID             uint      `db:"id" json:"id"`
	JobSeekerID    uint      `db:"job_seeker_id" json:"job_seeker_id"`
	Occupation     null.Int  `db:"occupation" json:"occupation"`
	ExperienceYear null.Int  `db:"experience_year" json:"experience_year"`
	CreatedAt      time.Time `db:"created_at" json:"-"`
	UpdatedAt      time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerExperienceJob(
	workHistoryID uint,
	occupation null.Int,
	experienceYear null.Int,
) *JobSeekerExperienceJob {
	return &JobSeekerExperienceJob{
		JobSeekerID:    workHistoryID,
		Occupation:     occupation,
		ExperienceYear: experienceYear,
	}
}
