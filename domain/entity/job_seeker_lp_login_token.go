package entity

import (
	"time"

	"github.com/google/uuid"
)

// 求職者の経験職種（LPからの登録時に使用）
type JobSeekerLPLoginToken struct {
	ID             uint      `db:"id" json:"id"`
	JobSeekerID    uint      `db:"job_seeker_id" json:"job_seeker_id"`
	LoginToken     uuid.UUID `db:"login_token" json:"login_token"`
	CreatedAt      time.Time `db:"created_at" json:"-"`
	UpdatedAt      time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerLPLoginToken(
	jobSeekerID uint,
	loginToken uuid.UUID,
) *JobSeekerLPLoginToken {
	return &JobSeekerLPLoginToken{
		JobSeekerID:    jobSeekerID,
		LoginToken:     loginToken,
	}
}
