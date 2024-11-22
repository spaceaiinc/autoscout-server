package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// 求人票
type JobListing struct {
	JobInformationID        uint      `db:"job_information_id" json:"job_information_id"`                   // 重複しないカラム毎の求人id
	JobInformationUUID      uuid.UUID `db:"job_information_uuid" json:"job_information_uuid"`               // 求人のUUID
	BillingAddressID        uint      `db:"billing_address_id" json:"billing_address_id"`                   // 重複しないID
	AgentStaffID            uint      `db:"agent_staff_id" json:"agent_staff_id"`                           // 重複しないカラム毎のRA担当者のid
	AgentID                 uint      `db:"agent_id" json:"agent_id,omitempty"`                             // エージェントID
	CompanyName             string    `db:"company_name" json:"company_name"`                               // 会社名
	CorporateSiteURL        string    `db:"corporate_site_url" json:"corporate_site_url"`                   // ホームページ
	PostCode                string    `db:"post_code" json:"post_code"`                                     // 郵便番号
	OfficeLocation          string    `db:"office_location" json:"office_location"`                         // 本社住所
	EmployeeNumberSingle    null.Int  `db:"employee_number_single" json:"employee_number_single"`           // 従業員数（単体）
	EmployeeNumberGroup     null.Int  `db:"employee_number_group" json:"employee_number_group"`             // 従業員数（連結）
	Establishment           string    `db:"establishment" json:"establishment"`                             // 設立（年月）
	PublicOffering          null.Int  `db:"public_offering" json:"public_offering"`                         // 株式公開
	Earnings                string    `db:"earnings" json:"earnings"`                                       // 売上高
	EarningsYear            null.Int  `db:"earnings_year" json:"earnings_year"`                             // 売上年度
	BusinessDetail          string    `db:"business_detail" json:"business_detail"`                         // 事業内容
	Title                   string    `db:"title" json:"title"`                                             // 求人タイトル
	WorkDetail              string    `db:"work_detail" json:"work_detail"`                                 // 仕事内容
	WorkLocation            string    `db:"work_location" json:"work_location"`                             // 勤務地（雇入れ直後）
	Transfer                null.Int  `db:"transfer" json:"transfer"`                                       // 転勤
	TransferDetail          string    `db:"transfer_detail" json:"transfer_detail"`                         // 変更の範囲
	UnderIncome             null.Int  `db:"under_income" json:"under_income"`                               // 年収下限
	OverIncome              null.Int  `db:"over_income" json:"over_income"`                                 // 年収上限
	Salary                  string    `db:"salary" json:"salary"`                                           // 給与詳細・昇給賞与
	Insurance               string    `db:"insurance" json:"insurance"`                                     // 諸手当・福利厚生
	WorkTime                string    `db:"work_time" json:"work_time"`                                     // 勤務時間
	OvertimeAverage         string    `db:"overtime_average" json:"overtime_average"`                       // 平均残業時間
	FixedOvertimePayment    null.Int  `db:"fixed_overtime_payment" json:"fixed_overtime_payment"`           // 固定残業代超過分の支払い有無
	FixedOvertimeDetail     string    `db:"fixed_overtime_detail" json:"fixed_overtime_detail"`             // 固定残業代の詳細
	TrialPeriod             null.Int  `db:"trial_period" json:"trial_period"`                               // 試用期間有無
	TrialPeriodDetail       string    `db:"trial_period_detail" json:"trial_period_detail"`                 // 試用期間詳細
	EmploymentPeriod        null.Int  `db:"employment_period" json:"employment_period"`                     // 雇用期間有無
	EmploymentPeriodDetail  string    `db:"employment_period_detail" json:"employment_period_detail"`       // 雇用期間詳細
	HolidayDetail           string    `db:"holiday_detail" json:"holiday_detail"`                           // 休日休暇
	PassiveSmoking          null.Int  `db:"passive_smoking" json:"passive_smoking"`                         // 受動喫煙対策の有無
	SelectionFlow           string    `db:"selection_flow" json:"selection_flow"`                           // 選考フロー
	EmploymentInsurance     bool      `db:"employment_insurance" json:"employment_insurance"`               // 雇用保険の有無
	AccidentInsurance       bool      `db:"accident_insurance" json:"accident_insurance"`                   // 労災保険の有無
	HealthInsurance         bool      `db:"health_insurance" json:"health_insurance"`                       // 健康保険の有無
	PensionInsurance        bool      `db:"pension_insurance" json:"pension_insurance"`                     // 厚生年金保険の有無
	IsExternal              bool      `db:"is_external" json:"is_external"`                                 // 外部求人フラグ true: 外部求人, false: autoscout内求人
	WorkDetailAfterHiring   string    `db:"work_detail_after_hiring" json:"work_detail_after_hiring"`       // 仕事内容（雇入れ直後）
	WorkDetailScopeOfChange string    `db:"work_detail_scope_of_change" json:"work_detail_scope_of_change"` // 仕事内容（変更の範囲）
	IsGuaranteedInterview   bool      `db:"is_guaranteed_interview" json:"is_guaranteed_interview"`         // 面接確約フラグ
	CreatedAt               time.Time `db:"created_at" json:"created_at"`                                   // 作成日時
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at"`                                   // 更新日時

	// 関連テーブル
	Industries         []EnterpriseIndustry             `db:"industries" json:"industries"`
	Prefectures        []JobInformationPrefecture       `db:"prefectures" json:"prefectures"`
	EmploymentStatuses []JobInformationEmploymentStatus `db:"employment_statuses" json:"employment_statuses"`
	WorkCharmPoints    []JobInformationWorkCharmPoint   `db:"work_charm_points" json:"work_charm_points"`
	Features           []JobInformationFeature          `db:"features" json:"features"`
	Occupations        []JobInformationOccupation       `db:"occupations" json:"occupations"`

	// 求人企業のログインに使用
	PhoneNumber string `db:"phone_number" json:"phone_number"` // 電話番号

	// 最新のタスク情報（求職者のマイページで使用）
	LatestTask TaskForJobListing `db:"latest_task" json:"latest_task"`
}

func NewJobListing(
	jobInformationID uint,
	agentStaffID uint,
	companyName string,
	corporateSiteURL string,
	postCode string,
	officeLocation string,
	employeeNumberSingle null.Int,
	employeeNumberGroup null.Int,
	establishment string,
	publicOffering null.Int,
	earnings string,
	earningsYear null.Int,
	businessDetail string,
	title string,
	workDetail string,
	workLocation string,
	transfer null.Int,
	transferDetail string,
	underIncome null.Int,
	overIncome null.Int,
	salary string,
	insurance string,
	workTime string,
	overtimeAverage string,
	fixedOvertimePayment null.Int,
	fixedOvertimeDetail string,
	trialPeriod null.Int,
	trialPeriodDetail string,
	employmentPeriod null.Int,
	employmentPeriodDetail string,
	holidayDetail string,
	passiveSmoking null.Int,
	selectionFlow string,
	employmentInsurance bool,
	accidentInsurance bool,
	healthInsurance bool,
	pensionInsurance bool,
	isExternal bool,
	workDetailAfterHiring string,
	workDetailScopeOfChange string,
) *JobListing {
	return &JobListing{
		JobInformationID:        jobInformationID,
		AgentStaffID:            agentStaffID,
		CompanyName:             companyName,
		CorporateSiteURL:        corporateSiteURL,
		PostCode:                postCode,
		OfficeLocation:          officeLocation,
		EmployeeNumberSingle:    employeeNumberSingle,
		EmployeeNumberGroup:     employeeNumberGroup,
		Establishment:           establishment,
		PublicOffering:          publicOffering,
		Earnings:                earnings,
		EarningsYear:            earningsYear,
		BusinessDetail:          businessDetail,
		Title:                   title,
		WorkDetail:              workDetail,
		WorkLocation:            workLocation,
		Transfer:                transfer,
		TransferDetail:          transferDetail,
		UnderIncome:             underIncome,
		OverIncome:              overIncome,
		Salary:                  salary,
		Insurance:               insurance,
		WorkTime:                workTime,
		OvertimeAverage:         overtimeAverage,
		FixedOvertimePayment:    fixedOvertimePayment,
		FixedOvertimeDetail:     fixedOvertimeDetail,
		TrialPeriod:             trialPeriod,
		TrialPeriodDetail:       trialPeriodDetail,
		EmploymentPeriod:        employmentPeriod,
		EmploymentPeriodDetail:  employmentPeriodDetail,
		HolidayDetail:           holidayDetail,
		PassiveSmoking:          passiveSmoking,
		SelectionFlow:           selectionFlow,
		EmploymentInsurance:     employmentInsurance,
		AccidentInsurance:       accidentInsurance,
		HealthInsurance:         healthInsurance,
		PensionInsurance:        pensionInsurance,
		IsExternal:              isExternal,
		WorkDetailAfterHiring:   workDetailAfterHiring,
		WorkDetailScopeOfChange: workDetailScopeOfChange,
	}
}

type JobListingForSending struct {
	ID                     uint      `db:"id" json:"id"`                                             // 重複しないID
	UUID                   uuid.UUID `db:"uuid" json:"uuid"`                                         // UUID
	JobInformationID       uint      `db:"job_information_id" json:"job_information_id"`             // 重複しないカラム毎の求人id
	JobInformationUUID     uuid.UUID `db:"job_information_uuid" json:"job_information_uuid"`         // 求人のUUID
	BillingAddressID       uint      `db:"billing_address_id" json:"billing_address_id"`             // 重複しないID
	AgentStaffID           uint      `db:"agent_staff_id" json:"agent_staff_id"`                     // 重複しないカラム毎のRA担当者のid
	AgentID                uint      `db:"agent_id" json:"agent_id,omitempty"`                       // エージェントID
	CompanyName            string    `db:"company_name" json:"company_name"`                         // 会社名
	CorporateSiteURL       string    `db:"corporate_site_url" json:"corporate_site_url"`             // ホームページ
	PostCode               string    `db:"post_code" json:"post_code"`                               // 郵便番号
	OfficeLocation         string    `db:"office_location" json:"office_location"`                   // 本社住所
	EmployeeNumberSingle   null.Int  `db:"employee_number_single" json:"employee_number_single"`     // 従業員数（単体）
	EmployeeNumberGroup    null.Int  `db:"employee_number_group" json:"employee_number_group"`       // 従業員数（連結）
	Establishment          string    `db:"establishment" json:"establishment"`                       // 設立（年月）
	PublicOffering         null.Int  `db:"public_offering" json:"public_offering"`                   // 株式公開
	Earnings               string    `db:"earnings" json:"earnings"`                                 // 売上高
	EarningsYear           null.Int  `db:"earnings_year" json:"earnings_year"`                       // 売上年度
	BusinessDetail         string    `db:"business_detail" json:"business_detail"`                   // 事業内容
	Title                  string    `db:"title" json:"title"`                                       // 求人タイトル
	Background             null.Int  `db:"background" json:"background"`                             // 募集背景
	WorkDetail             string    `db:"work_detail" json:"work_detail"`                           // 仕事内容
	WorkLocation           string    `db:"work_location" json:"work_location"`                       // 勤務地（雇入れ直後）
	Transfer               null.Int  `db:"transfer" json:"transfer"`                                 // 転勤
	TransferDetail         string    `db:"transfer_detail" json:"transfer_detail"`                   // 変更の範囲
	UnderIncome            null.Int  `db:"under_income" json:"under_income"`                         // 年収下限
	OverIncome             null.Int  `db:"over_income" json:"over_income"`                           // 年収上限
	Salary                 string    `db:"salary" json:"salary"`                                     // 給与詳細・昇給賞与
	Insurance              string    `db:"insurance" json:"insurance"`                               // 諸手当・福利厚生
	WorkTime               string    `db:"work_time" json:"work_time"`                               // 勤務時間
	Overtime               null.Int  `db:"overtime" json:"overtime"`                                 // 残業有無
	OvertimeAverage        string    `db:"overtime_average" json:"overtime_average"`                 // 平均残業時間
	FixedOvertime          null.Int  `db:"fixed_overtime" json:"fixed_overtime"`                     // 固定残業代有無
	FixedOvertimePayment   null.Int  `db:"fixed_overtime_payment" json:"fixed_overtime_payment"`     // 固定残業代超過分の支払い有無
	FixedOvertimeDetail    string    `db:"fixed_overtime_detail" json:"fixed_overtime_detail"`       // 固定残業代の詳細
	TrialPeriod            null.Int  `db:"trial_period" json:"trial_period"`                         // 試用期間有無
	TrialPeriodDetail      string    `db:"trial_period_detail" json:"trial_period_detail"`           // 試用期間詳細
	EmploymentPeriod       null.Int  `db:"employment_period" json:"employment_period"`               // 雇用期間有無
	EmploymentPeriodDetail string    `db:"employment_period_detail" json:"employment_period_detail"` // 雇用期間詳細
	HolidayDetail          string    `db:"holiday_detail" json:"holiday_detail"`                     // 休日休暇
	PassiveSmoking         null.Int  `db:"passive_smoking" json:"passive_smoking"`                   // 受動喫煙対策の有無
	SelectionFlow          string    `db:"selection_flow" json:"selection_flow"`                     // 選考フロー
	EmploymentInsurance    bool      `db:"employment_insurance" json:"employment_insurance"`         // 雇用保険の有無
	AccidentInsurance      bool      `db:"accident_insurance" json:"accident_insurance"`             // 労災保険の有無
	HealthInsurance        bool      `db:"health_insurance" json:"health_insurance"`                 // 健康保険の有無
	PensionInsurance       bool      `db:"pension_insurance" json:"pension_insurance"`               // 厚生年金保険の有無
	CreatedAt              time.Time `db:"created_at" json:"created_at"`                             // 作成日時
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`                             // 更新日時

	// 求人企業のログインに使用
	PhoneNumber string `db:"phone_number" json:"phone_number"` // 電話番号
}

func NewJobListingForSending(
	jobInformationID uint,
	agentStaffID uint,
	companyName string,
	corporateSiteURL string,
	postCode string,
	officeLocation string,
	employeeNumberSingle null.Int,
	employeeNumberGroup null.Int,
	establishment string,
	publicOffering null.Int,
	earnings string,
	earningsYear null.Int,
	businessDetail string,
	title string,
	background null.Int,
	workDetail string,
	workLocation string,
	transfer null.Int,
	transferDetail string,
	underIncome null.Int,
	overIncome null.Int,
	salary string,
	insurance string,
	workTime string,
	overtime null.Int,
	overtimeAverage string,
	fixedOvertime null.Int,
	fixedOvertimePayment null.Int,
	fixedOvertimeDetail string,
	trialPeriod null.Int,
	trialPeriodDetail string,
	employmentPeriod null.Int,
	employmentPeriodDetail string,
	holidayDetail string,
	passiveSmoking null.Int,
	selectionFlow string,
	employmentInsurance bool,
	accidentInsurance bool,
	healthInsurance bool,
	pensionInsurance bool,
) *JobListingForSending {
	return &JobListingForSending{
		JobInformationID:       jobInformationID,
		AgentStaffID:           agentStaffID,
		CompanyName:            companyName,
		CorporateSiteURL:       corporateSiteURL,
		PostCode:               postCode,
		OfficeLocation:         officeLocation,
		EmployeeNumberSingle:   employeeNumberSingle,
		EmployeeNumberGroup:    employeeNumberGroup,
		Establishment:          establishment,
		PublicOffering:         publicOffering,
		Earnings:               earnings,
		EarningsYear:           earningsYear,
		BusinessDetail:         businessDetail,
		Title:                  title,
		Background:             background,
		WorkDetail:             workDetail,
		WorkLocation:           workLocation,
		Transfer:               transfer,
		TransferDetail:         transferDetail,
		UnderIncome:            underIncome,
		OverIncome:             overIncome,
		Salary:                 salary,
		Insurance:              insurance,
		WorkTime:               workTime,
		Overtime:               overtime,
		OvertimeAverage:        overtimeAverage,
		FixedOvertime:          fixedOvertime,
		FixedOvertimePayment:   fixedOvertimePayment,
		FixedOvertimeDetail:    fixedOvertimeDetail,
		TrialPeriod:            trialPeriod,
		TrialPeriodDetail:      trialPeriodDetail,
		EmploymentPeriod:       employmentPeriod,
		EmploymentPeriodDetail: employmentPeriodDetail,
		HolidayDetail:          holidayDetail,
		PassiveSmoking:         passiveSmoking,
		SelectionFlow:          selectionFlow,
		EmploymentInsurance:    employmentInsurance,
		AccidentInsurance:      accidentInsurance,
		HealthInsurance:        healthInsurance,
		PensionInsurance:       pensionInsurance,
	}
}

// 求人票用のタスク情報
type TaskForJobListing struct {
	TaskGroupID           uint     `db:"task_group_id" json:"task_group_id"`
	JobInformationID      uint     `db:"job_information_id" json:"job_information_id"`
	JobSeekerID           uint     `db:"job_seeker_id" json:"job_seeker_id"`
	RAStaffID             uint     `db:"ra_staff_id" json:"ra_staff_id"`
	CAStaffID             uint     `db:"ca_staff_id" json:"ca_staff_id"`
	PhaseCategory         null.Int `db:"phase_category" json:"phase_category"`
	PhaseSubCategory      null.Int `db:"phase_sub_category" json:"phase_sub_category"`
	StaffType             null.Int `db:"staff_type" json:"staff_type"`
	ExecutedStaffID       uint     `db:"executed_staff_id" json:"executed_staff_id"`
	ExternalJobListingURL string   `db:"external_job_listing_url" json:"external_job_listing_url"` // 外部求人の時に格納される
}

func NewTaskForJobListing(
	taskGroupID uint,
	jobInformationsID uint,
	jobSeekerID uint,
	RAStaffID uint,
	CAStaffID uint,
	phaseCategory null.Int,
	phaseSubCategory null.Int,
	staffType null.Int,
	executedStaffID uint,
	externalJobListingURL string,
) *TaskForJobListing {
	return &TaskForJobListing{
		TaskGroupID:           taskGroupID,
		JobInformationID:      jobInformationsID,
		JobSeekerID:           jobSeekerID,
		RAStaffID:             RAStaffID,
		CAStaffID:             CAStaffID,
		PhaseCategory:         phaseCategory,
		PhaseSubCategory:      phaseSubCategory,
		StaffType:             staffType,
		ExecutedStaffID:       executedStaffID,
		ExternalJobListingURL: externalJobListingURL,
	}
}
