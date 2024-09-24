package entity

import (
	"time"
)

type SendingJobSeekerSelfPromotion struct {
	ID          uint      `db:"id" json:"id"`
	SendingJobSeekerID uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	Title       string    `db:"title" json:"title"`
	Contents    string    `db:"contents" json:"contents"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerSelfPromotion(
	sendingJobSeekerID uint,
	title string,
	contents string,
) *SendingJobSeekerSelfPromotion {
	return &SendingJobSeekerSelfPromotion{
		SendingJobSeekerID: sendingJobSeekerID,
		Title:       title,
		Contents:    contents,
	}
}
