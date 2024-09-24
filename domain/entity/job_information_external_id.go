package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// 求人の他媒体でのIDを管理するテーブル
type JobInformationExternalID struct {
	ID               uint      `db:"id" json:"-"`
	JobInformationID uint      `db:"job_information_id" json:"job_information_id"`
	ExternalType     null.Int  `db:"external_type" json:"external_type"` // 他媒体のタイプ(0: エージェントバンク)
	ExternalID       string    `db:"external_id" json:"external_id"`     // 他媒体でのID
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationExternalID(
	jobInformationID uint,
	externalType null.Int,
	externalID string,
) *JobInformationExternalID {
	return &JobInformationExternalID{
		JobInformationID: jobInformationID,
		ExternalType:     externalType,
		ExternalID:       externalID,
	}
}

// 求人の他媒体のタイプ
type JobInformatinoExternalType int64

const (
	// エージェントバンク
	JobInformatinoExternalTypeAgentBank JobInformatinoExternalType = iota
)
