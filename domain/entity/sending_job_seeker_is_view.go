package entity

import (
	"time"
)

type SendingJobSeekerIsView struct {
	ID                    uint      `db:"id" json:"id"`
	SendingJobSeekerID    uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	IsNotWaitingViewed    bool      `db:"is_not_waiting_viewed" json:"is_not_waiting_viewed"`
	IsNotUnregisterViewed bool      `db:"is_not_unregister_viewed" json:"is_not_unregister_viewed"`
	CreatedAt             time.Time `db:"created_at" json:"-"`
	UpdatedAt             time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerIsView(
	sendingJobSeekerID uint,
	isNotWaitingViewed bool,
	isNotUnregisterViewed bool,
) *SendingJobSeekerIsView {
	return &SendingJobSeekerIsView{
		SendingJobSeekerID:    sendingJobSeekerID,
		IsNotWaitingViewed:    isNotWaitingViewed,
		IsNotUnregisterViewed: isNotUnregisterViewed,
	}
}

type IsViewCount struct {
	WaitingCount    uint `db:"waiting_count" json:"waiting_count"`
	UnregisterCount uint `db:"unregister_count" json:"unregister_count"`
}
