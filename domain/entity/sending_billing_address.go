package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type SendingBillingAddress struct {
	ID                    uint      `db:"id" json:"id"`
	UUID                  uuid.UUID `db:"uuid" json:"uuid"`
	SendingEnterpriseID   uint      `db:"sending_enterprise_id" json:"sending_enterprise_id"`
	AgentID               uint      `db:"agent_id" json:"agent_id"`
	AgentStaffID          uint      `db:"agent_staff_id" json:"agent_staff_id"`
	StaffName             string    `db:"staff_name" json:"staff_name"`
	ContractPhase         null.Int  `db:"contract_phase" json:"contract_phase"` // 0: リーガルチェック中, 1: リーガルチェック完了, 2: 契約締結済み
	ContractDate          string    `db:"contract_date" json:"contract_date"`
	PaymentPolicy         string    `db:"payment_policy" json:"payment_policy"`
	CompanyName           string    `db:"company_name" json:"company_name"`
	Address               string    `db:"address" json:"address"`
	Title                 string    `db:"title" json:"title"`                                     // 請求先タイトル
	ScheduleAdjustmentURL string    `db:"schedule_adjustment_url" json:"schedule_adjustment_url"` // 日程調整URL
	Commission            null.Int  `db:"commission" json:"commission"`                           // 送客単価
	IsDeleted             bool      `db:"is_deleted" json:"is_deleted"`                           // 論理削除フラグ false: 有効, true: 削除済み
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time `db:"updated_at" json:"updated_at"`

	// 他テーブル
	Staffs                    []SendingBillingAddressStaff `json:"staffs"`
	SendingJobInformationList []SendingJobInformation      `json:"sending_job_information_list"`
}

func NewSendingBillingAddress(
	sendingEnterpriseID uint,
	agentStaffID uint,
	contractPhase null.Int,
	contractDate string,
	paymentPolicy string,
	companyName string,
	address string,
	title string,
	scheduleAdjustmentURL string,
	commission null.Int,
) *SendingBillingAddress {
	return &SendingBillingAddress{
		SendingEnterpriseID:   sendingEnterpriseID,
		AgentStaffID:          agentStaffID,
		ContractPhase:         contractPhase,
		ContractDate:          contractDate,
		PaymentPolicy:         paymentPolicy,
		CompanyName:           companyName,
		Address:               address,
		Title:                 title,
		ScheduleAdjustmentURL: scheduleAdjustmentURL,
		Commission:            commission,
	}
}

type UpdateSendingBillingAddressParam struct {
	SendingEnterpriseID   uint     `json:"sending_enterprise_id"`
	AgentStaffID          uint     `json:"agent_staff_id" validate:"required"`
	ContractPhase         null.Int `json:"contract_phase"`
	ContractDate          string   `json:"contract_date"`
	PaymentPolicy         string   `json:"payment_policy"`
	CompanyName           string   `db:"company_name" json:"company_name"`
	Address               string   `json:"address"`
	Title                 string   `json:"title"`
	ScheduleAdjustmentURL string   `db:"schedule_adjustment_url" json:"schedule_adjustment_url"` // 日程調整URL
	Commission            null.Int `db:"commission" json:"commission"`                           // 送客単価

	// 他テーブル
	Staffs []SendingBillingAddressStaff `json:"staffs"`
}
