package entity

import (
	"time"
)

type JobSeekerSelfPromotion struct {
	ID          uint      `db:"id" json:"id"`
	JobSeekerID uint      `db:"job_seeker_id" json:"job_seeker_id"`
	Title       string    `db:"title" json:"title"`
	Contents    string    `db:"contents" json:"contents"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerSelfPromotion(
	jobSeekerID uint,
	title string,
	contents string,
) *JobSeekerSelfPromotion {
	return &JobSeekerSelfPromotion{
		JobSeekerID: jobSeekerID,
		Title:       title,
		Contents:    contents,
	}
}
