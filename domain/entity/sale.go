package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Sale struct {
	ID                  uint       `db:"id" json:"id"`
	JobSeekerID         uint       `db:"job_seeker_id" json:"job_seeker_id"`
	JobInformationID    uint       `db:"job_information_id" json:"job_information_id"`
	Accuracy            null.Int   `db:"accuracy" json:"accuracy"`
	ContractSignedMonth string     `db:"contract_signed_month" json:"contract_signed_month"`
	BillingMonth        string     `db:"billing_month" json:"billing_month"`
	BillingAmount       null.Float `db:"billing_amount" json:"billing_amount"`
	Cost                null.Float `db:"cost" json:"cost"`
	GrossProfit         null.Float `db:"gross_profit" json:"gross_profit"`
	RAStaffID           uint       `db:"ra_staff_id" json:"ra_staff_id"`
	CAStaffID           uint       `db:"ca_staff_id" json:"ca_staff_id"`
	RaSalesRatio        null.Float `db:"ra_sales_ratio" json:"ra_sales_ratio"`
	CaSalesRatio        null.Float `db:"ca_sales_ratio" json:"ca_sales_ratio"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updated_at"`

	// 担当者
	RAAgentID   uint   `db:"ra_agent_id" json:"ra_agent_id"`
	CAAgentID   uint   `db:"ca_agent_id" json:"ca_agent_id"`
	RAStaffName string `db:"ra_staff_name" json:"ra_staff_name"`
	CAStaffName string `db:"ca_staff_name" json:"ca_staff_name"`
	RAAgentName string `db:"ra_agent_name" json:"ra_agent_name"`
	CAAgentName string `db:"ca_agent_name" json:"ca_agent_name"`

	// 求職者・求人・企業
	LastName    string `db:"last_name" json:"last_name"`
	FirstName   string `db:"first_name" json:"first_name"`
	Title       string `db:"title" json:"title"`
	CompanyName string `db:"company_name" json:"company_name"`

	// タスク
	TaskGroupID           uint     `db:"task_group_id" json:"task_group_id"`
	PhaseCategory         null.Int `db:"phase_category" json:"phase_category"`
	PhaseSubCategory      null.Int `db:"phase_sub_category" json:"phase_sub_category"`
	ExternalJobListingURL string   `db:"external_job_listing_url" json:"external_job_listing_url"`
}

func NewSale(
	jobSeekerID uint,
	jobInformationID uint,
	accuracy null.Int,
	contractSignedMonth string,
	billingMonth string,
	billingAmount null.Float,
	cost null.Float,
	grossProfit null.Float,
	raStaffID uint,
	raSalesRatio null.Float,
	caStaffID uint,
	caSalesRatio null.Float,
) *Sale {
	return &Sale{
		JobSeekerID:         jobSeekerID,
		JobInformationID:    jobInformationID,
		Accuracy:            accuracy,
		ContractSignedMonth: contractSignedMonth,
		BillingMonth:        billingMonth,
		BillingAmount:       billingAmount,
		Cost:                cost,
		GrossProfit:         grossProfit,
		RAStaffID:           raStaffID,
		RaSalesRatio:        raSalesRatio,
		CAStaffID:           caStaffID,
		CaSalesRatio:        caSalesRatio,
	}
}

type CreateOrUpdateSaleParam struct {
	JobSeekerID         uint       `json:"job_seeker_id" validate:"required"`
	JobInformationID    uint       `json:"job_information_id" validate:"required"`
	Accuracy            null.Int   `db:"accuracy" json:"accuracy"`
	ContractSignedMonth string     `db:"contract_signed_month" json:"contract_signed_month"`
	BillingMonth        string     `db:"billing_month" json:"billing_month"`
	BillingAmount       null.Float `db:"billing_amount" json:"billing_amount"`
	Cost                null.Float `db:"cost" json:"cost"`
	GrossProfit         null.Float `db:"gross_profit" json:"gross_profit"`
	RAStaffID           uint       `db:"ra_staff_id" json:"ra_staff_id"`
	RaSalesRatio        null.Float `db:"ra_sales_ratio" json:"ra_sales_ratio"`
	CAStaffID           uint       `db:"ca_staff_id" json:"ca_staff_id"`
	CaSalesRatio        null.Float `db:"ca_sales_ratio" json:"ca_sales_ratio"`
}

type DeleteSaleParam struct {
	ID uint `json:"id" validate:"required"`
}

const (
	AccuracyAccept  int64 = iota // 内定承諾
	AccuracyA                    // Aヨミ
	AccuracyB                    // Bヨミ
	AccuracyC                    // Cヨミ
	AccuracyTopic                // ネタ
	AccuracyFailure              // 失注
)

type SearchSale struct {
	RAStaffID  string
	CAStaffID  string
	Accuracies []null.Int
}

func NewSearchSale(
	raStaffID string,
	caStaffID string,
	accuracies []null.Int,
) *SearchSale {
	return &SearchSale{
		RAStaffID:  raStaffID,
		CAStaffID:  caStaffID,
		Accuracies: accuracies,
	}
}

type SearchAccuracy struct {
	AgentID                uint
	JobSeekerFreeWord      string
	JobInformationFreeWord string
	ContractSignedMonth    string
	BillingMonth           string
	RAStaffID              string     // interactorの絞り込み時に数値に変換
	CAStaffID              string     // interactorの絞り込み時に数値に変換
	Accuracies             []uint // interactorの絞り込み時に数値に変換
	PageNumber             uint
}

func NewSearchSearchAccuracy(
	agentID uint,
	jobSeekerFreeWord string,
	jobInformationFreeWord string,
	contractSignedMonth string,
	billingMonth string,
	raStaffID string,
	caStaffID string,
	accuracies []uint,
	pageNumber uint,
) *SearchAccuracy {
	return &SearchAccuracy{
		AgentID:                agentID,
		JobSeekerFreeWord:      jobSeekerFreeWord,
		JobInformationFreeWord: jobInformationFreeWord,
		ContractSignedMonth:    contractSignedMonth,
		BillingMonth:           billingMonth,
		RAStaffID:              raStaffID,
		CAStaffID:              caStaffID,
		Accuracies:             accuracies,
		PageNumber:             pageNumber,
	}
}
