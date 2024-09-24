package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// 他媒体の求人企業を一括インポートするテーブル
type InitialEnterpriseImporter struct {
	ID           uint      `json:"id" db:"id"`                                         // 重複しないID
	UUID         string    `json:"uuid" db:"uuid"`                                     // 重複しないカラム毎のUUID
	AgentStaffID uint      `json:"agent_staff_id" db:"agent_staff_id"`                 // エージェントスタッフID RA担当者
	ServiceType  null.Int  `json:"service_type" db:"service_type" validate:"required"` // サービスタイプ 0: サーカス
	LoginID      string    `json:"login_id" db:"login_id" validate:"required"`         // ログインID
	Password     string    `json:"password" db:"password" validate:"required"`         // パスワード
	StartDate    string    `json:"start_date" db:"start_date" validate:"required"`     // 取得開始日 fmt: yyyy-mm-dd
	StartHour    null.Int  `json:"start_hour" db:"start_hour" validate:"required"`     // 実行開始時間
	MaxCount     null.Int  `json:"max_count" db:"max_count"`                           // 取得件数
	Offset       null.Int  `json:"offset" db:"offset"`                                 // 取得開始位置
	SearchURL    string    `json:"search_url" db:"search_url"`                         // 検索URL サーカスの場合
	IsSuccess    bool      `json:"is_success" db:"is_success"`                         // 成功したかどうか true: 成功 false: 失敗or未実行
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func NewInitialEnterpriseImporter(
	agentStaffID uint,
	serviceType null.Int,
	loginID string,
	password string,
	startDate string,
	startHour null.Int,
	maxCount null.Int,
	offset null.Int,
	searchURL string,
	isSuccess bool,
) *InitialEnterpriseImporter {
	return &InitialEnterpriseImporter{
		AgentStaffID: agentStaffID,
		ServiceType:  serviceType,
		LoginID:      loginID,
		Password:     password,
		StartDate:    startDate,
		StartHour:    startHour,
		MaxCount:     maxCount,
		Offset:       offset,
		SearchURL:    searchURL,
		IsSuccess:    isSuccess,
	}
}

// サービスタイプ
const (
	EnterpriseImporterServiceTypeCircus int64 = iota
	EnterpriseImporterServiceTypeAgentBank
)

var EnterpriseImporterServiceTypeLabel = map[int64]string{
	EnterpriseImporterServiceTypeCircus:    "サーカス",
	EnterpriseImporterServiceTypeAgentBank: "エージェントバンク",
}
