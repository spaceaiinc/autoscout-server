package entity

import (
	"time"
)

type SendingJobSeekerDesiredJobInformation struct {
	ID                      uint      `json:"id" db:"id"`
	SendingJobSeekerID      uint      `json:"sending_job_seeker_id" db:"sending_job_seeker_id"`
	SendingJobInformationID uint      `json:"sending_job_information_id" db:"sending_job_information_id"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`

	// db外項目
	SendingEnterpriseID uint   `json:"sending_enterprise_id" db:"sending_enterprise_id"`
	AgentName           string `json:"agent_name" db:"agent_name"`
	CompanyName         string `json:"company_name" db:"company_name"`
	Title               string `json:"title" db:"title"`
}

func NewSendingJobSeekerDesiredJobInformation(
	sendingJobSeekerID uint,
	sendingJobInformationID uint,
) *SendingJobSeekerDesiredJobInformation {
	return &SendingJobSeekerDesiredJobInformation{
		SendingJobSeekerID:      sendingJobSeekerID,
		SendingJobInformationID: sendingJobInformationID,
	}
}

// 複数作成するためのパラメータ
type CreateMultiSendingJobSeekerDesiredJobInformationParam struct {
	SendingJobSeekerID          uint   `json:"sending_job_seeker_id" db:"sending_job_seeker_id"`
	SendingJobInformationIDList []uint `json:"sending_job_information_id_list" db:"sending_job_information_id_list"`
}
