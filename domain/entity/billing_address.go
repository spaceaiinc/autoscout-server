package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type BillingAddress struct {
	ID             uint      `db:"id" json:"id"`
	UUID           uuid.UUID `db:"uuid" json:"uuid"`
	EnterpriseID   uint      `db:"enterprise_id" json:"enterprise_id"`
	AgentID        uint      `db:"agent_id" json:"agent_id"`
	AgentStaffID   uint      `db:"agent_staff_id" json:"agent_staff_id"`
	StaffName      string    `db:"staff_name" json:"staff_name"`
	ContractPhase  null.Int  `db:"contract_phase" json:"contract_phase"` // 0: リーガルチェック中, 1: リーガルチェック完了, 2: 契約締結済み
	ContractDate   string    `db:"contract_date" json:"contract_date"`
	PaymentPolicy  string    `db:"payment_policy" json:"payment_policy"`
	CompanyName    string    `db:"company_name" json:"company_name"`
	Address        string    `db:"address" json:"address"`
	HowToRecommend string    `db:"how_to_recommend" json:"how_to_recommend"`
	Title          string    `db:"title" json:"title"` // 請求先タイトル
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted      bool      `db:"is_deleted" json:"is_deleted"` // 論理削除フラグ false: 有効, true: 削除済み

	// 他テーブル
	HRStaffs []BillingAddressHRStaff `json:"hr_staffs"`
	RAStaffs []BillingAddressRAStaff `json:"ra_staffs"`

	EnterpriseCompanyName string `db:"enterprise_company_name" json:"-"`
}

func NewBillingAddress(
	enterpriseID uint,
	agentStaffID uint,
	contractPhase null.Int,
	contractDate string,
	paymentPolicy string,
	companyName string,
	address string,
	howToRecommend string,
	title string,
) *BillingAddress {
	return &BillingAddress{
		EnterpriseID:   enterpriseID,
		AgentStaffID:   agentStaffID,
		ContractPhase:  contractPhase,
		ContractDate:   contractDate,
		PaymentPolicy:  paymentPolicy,
		CompanyName:    companyName,
		Address:        address,
		HowToRecommend: howToRecommend,
		Title:          title,
	}
}

type CreateBillingAddressParam struct {
	AgentStaffID   uint     `json:"agent_staff_id" validate:"required"`
	ContractPhase  null.Int `json:"contract_phase"`
	ContractDate   string   `json:"contract_date"`
	PaymentPolicy  string   `json:"payment_policy"`
	CompanyName    string   `json:"company_name"`
	Address        string   `json:"address"`
	HowToRecommend string   `json:"how_to_recommend"`
	Title          string   `json:"title"`

	// 他テーブル
	HRStaffs []BillingAddressHRStaff `json:"hr_staffs"`
	RAStaffs []BillingAddressRAStaff `json:"ra_staffs"`
}

type UpdateBillingAddressParam struct {
	EnterpriseID   uint     `json:"enterprise_id"`
	AgentStaffID   uint     `json:"agent_staff_id" validate:"required"`
	ContractPhase  null.Int `json:"contract_phase"`
	ContractDate   string   `json:"contract_date"`
	PaymentPolicy  string   `json:"payment_policy"`
	CompanyName    string   `json:"company_name"`
	Address        string   `json:"address"`
	HowToRecommend string   `json:"how_to_recommend"`
	Title          string   `json:"title"`

	// 他テーブル
	HRStaffs []BillingAddressHRStaff `json:"hr_staffs"`
	RAStaffs []BillingAddressRAStaff `json:"ra_staffs"`
}

type DeleteBillingAddressParam struct {
	ID uint `json:"id" validate:"required"`
}

type SearchBillingAddress struct {
	FreeWord     string
	AgentStaffID string
}

func NewSearchBillingAddress(
	freeword string,
	agentStaffID string,
) *SearchBillingAddress {
	return &SearchBillingAddress{
		FreeWord:     freeword,
		AgentStaffID: agentStaffID,
	}
}
