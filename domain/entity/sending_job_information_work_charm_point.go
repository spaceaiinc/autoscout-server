package entity

import (
	"time"
)

type SendingJobInformationWorkCharmPoint struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobInformationID uint      `db:"sending_job_information_id" json:"sending_job_information_id"`
	Title                   string    `db:"title" json:"title"`
	Contents                string    `db:"contents" json:"contents"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationWorkCharmPoint(
	sendingJobInformationID uint,
	title string,
	contents string,
) *SendingJobInformationWorkCharmPoint {
	return &SendingJobInformationWorkCharmPoint{
		SendingJobInformationID: sendingJobInformationID,
		Title:                   title,
		Contents:                contents,
	}
}
