package entity

import (
	"time"
)

type JobInformationWorkCharmPoint struct {
	ID               uint      `db:"id" json:"id"`
	JobInformationID uint      `db:"job_information_id" json:"job_information_id"`
	Title            string    `db:"title" json:"title"`
	Contents         string    `db:"contents" json:"contents"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationWorkCharmPoint(
	jobInformationID uint,
	title string,
	contents string,
) *JobInformationWorkCharmPoint {
	return &JobInformationWorkCharmPoint{
		JobInformationID: jobInformationID,
		Title:            title,
		Contents:         contents,
	}
}
