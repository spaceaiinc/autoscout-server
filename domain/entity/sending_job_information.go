package entity

import (
	"time"

	"github.com/google/uuid"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformation struct {
	ID                      uint      `db:"id" json:"id"`
	UUID                    uuid.UUID `db:"uuid" json:"uuid"`
	SendingEnterpriseName   string    `db:"sending_enterprise_name" json:"sending_enterprise_name"` // enterprisesのcompany_name
	AgentStaffID            uint      `db:"agent_staff_id" json:"agent_staff_id"`
	StaffName               string    `db:"staff_name" json:"staff_name"`
	SendingBillingAddressID uint      `db:"sending_billing_address_id" json:"sending_billing_address_id"`
	CompanyName             string    `db:"company_name" json:"company_name"`
	Title                   string    `db:"title" json:"title"`
	RecruitmentState        null.Int  `db:"recruitment_state" json:"recruitment_state"`
	ExpirationDate          string    `db:"expiration_date" json:"expiration_date"`
	Background              null.Int  `db:"background" json:"background"`
	WorkDetail              string    `db:"work_detail" json:"work_detail"`
	NumberOfHires           null.Int  `db:"number_of_hires" json:"number_of_hires"`
	WorkLocation            string    `db:"work_location" json:"work_location"`
	Transfer                null.Int  `db:"transfer" json:"transfer"`
	TransferDetail          string    `db:"transfer_detail" json:"transfer_detail"`
	UnderIncome             null.Int  `db:"under_income" json:"under_income"`
	OverIncome              null.Int  `db:"over_income" json:"over_income"`
	Salary                  string    `db:"salary" json:"salary"`
	Insurance               string    `db:"insurance" json:"insurance"`
	WorkTime                string    `db:"work_time" json:"work_time"`
	Overtime                null.Int  `db:"overtime" json:"overtime"`
	OvertimeAverage         string    `db:"overtime_average" json:"overtime_average"`
	FixedOvertime           null.Int  `db:"fixed_overtime" json:"fixed_overtime"`
	FixedOvertimePayment    null.Int  `db:"fixed_overtime_payment" json:"fixed_overtime_payment"`
	FixedOvertimeDetail     string    `db:"fixed_overtime_detail" json:"fixed_overtime_detail"`
	TrialPeriod             null.Int  `db:"trial_period" json:"trial_period"`
	TrialPeriodDetail       string    `db:"trial_period_detail" json:"trial_period_detail"`
	EmploymentPeriod        null.Int  `db:"employment_period" json:"employment_period"`
	EmploymentPeriodDetail  string    `db:"employment_period_detail" json:"employment_period_detail"`
	HolidayType             null.Int  `db:"holiday_type" json:"holiday_type"`
	HolidayDetail           string    `db:"holiday_detail" json:"holiday_detail"`
	PassiveSmoking          null.Int  `db:"passive_smoking" json:"passive_smoking"`
	SelectionFlow           string    `db:"selection_flow" json:"selection_flow"`
	Gender                  null.Int  `db:"gender" json:"gender"`
	Nationality             null.Int  `db:"nationality" json:"nationality"`
	FinalEducation          null.Int  `db:"final_education" json:"final_education"`
	SchoolLevel             null.Int  `db:"school_level" json:"school_level"`
	MedicalHistory          null.Int  `db:"medical_history" json:"medical_history"`
	AgeUnder                null.Int  `db:"age_under" json:"age_under"`
	AgeOver                 null.Int  `db:"age_over" json:"age_over"`
	JobChange               null.Int  `db:"job_change" json:"job_change"`
	ShortResignation        null.Int  `db:"short_resignation" json:"short_resignation"`
	ShortResignationRemarks string    `db:"short_resignation_remarks" json:"short_resignation_remarks"`
	SocialExperienceYear    null.Int  `db:"social_experience_year" json:"social_experience_year"`
	SocialExperienceMonth   null.Int  `db:"social_experience_month" json:"social_experience_month"`
	OtherRequired           string    `db:"other_required" json:"other_required"`
	Appearance              null.Int  `db:"appearance" json:"appearance"`
	Communication           null.Int  `db:"communication" json:"communication"`
	Thinking                null.Int  `db:"thinking" json:"thinking"`
	TargetDetail            string    `db:"target_detail" json:"target_detail"`
	EmploymentInsurance     bool      `db:"employment_insurance" json:"employment_insurance"`     // 雇用保険の有無
	AccidentInsurance       bool      `db:"accident_insurance" json:"accident_insurance"`         // 労災保険の有無
	HealthInsurance         bool      `db:"health_insurance" json:"health_insurance"`             // 健康保険の有無
	PensionInsurance        bool      `db:"pension_insurance" json:"pension_insurance"`           // 厚生年金保険の有無
	RegisterPhase           null.Int  `db:"register_phase" json:"register_phase"`                 // 求人の登録状況（0: 本登録, 1: 仮登録）
	StudyCategory           null.Int  `db:"study_category" json:"study_category"`                 // 専攻学科の大分類(0:　理系尚可, 1: 理系のみ, 99: 不問)
	WordSkill               null.Int  `db:"word_skill" json:"word_skill"`                         // Wordのスキル
	ExcelSkill              null.Int  `db:"excel_skill" json:"excel_skill"`                       // Excelのスキル
	PowerPointSkill         null.Int  `db:"power_point_skill" json:"power_point_skill"`           // PowerPointのスキル
	CorporateSiteURL        string    `db:"corporate_site_url" json:"corporate_site_url"`         // 企業サイトURL
	PostCode                string    `db:"post_code" json:"post_code"`                           // 郵便番号
	OfficeLocation          string    `db:"office_location" json:"office_location"`               // 本社所在地
	Establishment           string    `db:"establishment" json:"establishment"`                   // 設立年月
	EmployeeNumberSingle    null.Int  `db:"employee_number_single" json:"employee_number_single"` // 従業員数（単体）
	EmployeeNumberGroup     null.Int  `db:"employee_number_group" json:"employee_number_group"`   // 従業員数（連結）
	PublicOffering          null.Int  `db:"public_offering" json:"public_offering"`               // 上場区分
	EarningsYear            null.Int  `db:"earnings_year" json:"earnings_year"`                   // 売上年度
	Earnings                string    `db:"earnings" json:"earnings"`                             // 売上高
	BusinessDetail          string    `db:"business_detail" json:"business_detail"`               // 事業内容
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted               bool      `db:"is_deleted" json:"is_deleted"` // 論理削除フラグ false: 有効, true: 削除済み

	RequiredExperienceJobDetail string `db:"required_experience_job_detail" json:"required_experience_job_detail"`

	// 他テーブル
	Targets                   []SendingJobInformationTarget                   `db:"-" json:"targets"`
	Features                  []SendingJobInformationFeature                  `db:"-" json:"features"`
	Prefectures               []SendingJobInformationPrefecture               `db:"-" json:"prefectures"`
	EmploymentStatuses        []SendingJobInformationEmploymentStatus         `db:"-" json:"employment_statuses"`
	WorkCharmPoints           []SendingJobInformationWorkCharmPoint           `db:"-" json:"work_charm_points"`
	RequiredConditions        []SendingJobInformationRequiredCondition        `db:"-" json:"required_conditions"` // 必要条件　複数
	RequiredSocialExperiences []SendingJobInformationRequiredSocialExperience `db:"-" json:"required_social_experiences"`
	Occupations               []SendingJobInformationOccupation               `db:"-" json:"occupations"`
	Industries                []SendingJobInformationIndustry                 `db:"-" json:"industries"`
	ReferenceMaterial         SendingEnterpriseReferenceMaterial              `db:"-" json:"reference_materials"`

	// 他社エージェント同士の求人の重複判定用
	IsDuplicate bool `json:"is_duplicate"`

	SendingEnterpriseID uint `db:"sending_enterprise_id" json:"sending_enterprise_id"`
	AgentID             uint `db:"agent_id" json:"agent_id"`

	// csv
	RecordLine uint `json:"record_line,omitempty"` // レコード行数 csvインポート用
}

func NewSendingJobInformation(
	sendingBillingAddressID uint,
	companyName string,
	title string,
	recruitmentState null.Int,
	expirationDate string,
	background null.Int,
	workDetail string,
	numberOfHires null.Int,
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
	holidayType null.Int,
	holidayDetail string,
	passiveSmoking null.Int,
	selectionFlow string,
	gender null.Int,
	nationality null.Int,
	finalEducation null.Int,
	schoolLevel null.Int,
	medicalHistory null.Int,
	ageUnder null.Int,
	ageOver null.Int,
	jobChange null.Int,
	shortResignation null.Int,
	shortResignationRemarks string,
	socialExperienceYear null.Int,
	socialExperienceMonth null.Int,
	otherRequired string,
	appearance null.Int,
	communication null.Int,
	thinking null.Int,
	targetDetail string,
	requiredExperienceJobDetail string,
	employmentInsurance bool,
	accidentInsurance bool,
	healthInsurance bool,
	pensionInsurance bool,
	registerPhase null.Int,
	studyCategory null.Int,
	wordSkill null.Int,
	excelSkill null.Int,
	powerPointSkill null.Int,
	corporateSiteURL string,
	postCode string,
	officeLocation string,
	establishment string,
	employeeNumberSingle null.Int,
	employeeNumberGroup null.Int,
	publicOffering null.Int,
	earningsYear null.Int,
	earnings string,
	businessDetail string,
) *SendingJobInformation {
	return &SendingJobInformation{
		SendingBillingAddressID:     sendingBillingAddressID,
		CompanyName:                 companyName,
		Title:                       title,
		RecruitmentState:            recruitmentState,
		ExpirationDate:              expirationDate,
		Background:                  background,
		WorkDetail:                  workDetail,
		NumberOfHires:               numberOfHires,
		WorkLocation:                workLocation,
		Transfer:                    transfer,
		TransferDetail:              transferDetail,
		UnderIncome:                 underIncome,
		OverIncome:                  overIncome,
		Salary:                      salary,
		Insurance:                   insurance,
		WorkTime:                    workTime,
		Overtime:                    overtime,
		OvertimeAverage:             overtimeAverage,
		FixedOvertime:               fixedOvertime,
		FixedOvertimePayment:        fixedOvertimePayment,
		FixedOvertimeDetail:         fixedOvertimeDetail,
		TrialPeriod:                 trialPeriod,
		TrialPeriodDetail:           trialPeriodDetail,
		EmploymentPeriod:            employmentPeriod,
		EmploymentPeriodDetail:      employmentPeriodDetail,
		HolidayType:                 holidayType,
		HolidayDetail:               holidayDetail,
		PassiveSmoking:              passiveSmoking,
		SelectionFlow:               selectionFlow,
		Gender:                      gender,
		Nationality:                 nationality,
		FinalEducation:              finalEducation,
		SchoolLevel:                 schoolLevel,
		MedicalHistory:              medicalHistory,
		AgeUnder:                    ageUnder,
		AgeOver:                     ageOver,
		JobChange:                   jobChange,
		ShortResignation:            shortResignation,
		ShortResignationRemarks:     shortResignationRemarks,
		SocialExperienceYear:        socialExperienceYear,
		SocialExperienceMonth:       socialExperienceMonth,
		OtherRequired:               otherRequired,
		Appearance:                  appearance,
		Communication:               communication,
		Thinking:                    thinking,
		TargetDetail:                targetDetail,
		RequiredExperienceJobDetail: requiredExperienceJobDetail,
		EmploymentInsurance:         employmentInsurance,
		AccidentInsurance:           accidentInsurance,
		HealthInsurance:             healthInsurance,
		PensionInsurance:            pensionInsurance,
		RegisterPhase:               registerPhase,
		StudyCategory:               studyCategory,
		WordSkill:                   wordSkill,
		ExcelSkill:                  excelSkill,
		PowerPointSkill:             powerPointSkill,
		CorporateSiteURL:            corporateSiteURL,
		PostCode:                    postCode,
		OfficeLocation:              officeLocation,
		Establishment:               establishment,
		EmployeeNumberSingle:        employeeNumberSingle,
		EmployeeNumberGroup:         employeeNumberGroup,
		PublicOffering:              publicOffering,
		EarningsYear:                earningsYear,
		Earnings:                    earnings,
		BusinessDetail:              businessDetail,
	}
}

type CreateSendingJobInformationParam struct {
	SendingBillingAddressID     null.Int `db:"sending_billing_address_id" json:"sending_billing_address_id" validate:"required"`
	CompanyName                 string   `db:"company_name" json:"company_name"`
	Title                       string   `db:"title" json:"title"`
	RecruitmentState            null.Int `db:"recruitment_state" json:"recruitment_state"`
	ExpirationDate              string   `db:"expiration_date" json:"expiration_date"`
	Background                  null.Int `db:"background" json:"background"`
	WorkDetail                  string   `db:"work_detail" json:"work_detail"`
	NumberOfHires               null.Int `db:"number_of_hires" json:"number_of_hires"`
	WorkLocation                string   `db:"work_location" json:"work_location"`
	Transfer                    null.Int `db:"transfer" json:"transfer"`
	TransferDetail              string   `db:"transfer_detail" json:"transfer_detail"`
	UnderIncome                 null.Int `db:"under_income" json:"under_income"`
	OverIncome                  null.Int `db:"over_income" json:"over_income"`
	Salary                      string   `db:"salary" json:"salary"`
	Insurance                   string   `db:"insurance" json:"insurance"`
	WorkTime                    string   `db:"work_time" json:"work_time"`
	Overtime                    null.Int `db:"overtime" json:"overtime"`
	OvertimeAverage             string   `db:"overtime_average" json:"overtime_average"`
	FixedOvertime               null.Int `db:"fixed_overtime" json:"fixed_overtime"`
	FixedOvertimePayment        null.Int `db:"fixed_overtime_payment" json:"fixed_overtime_payment"`
	FixedOvertimeDetail         string   `db:"fixed_overtime_detail" json:"fixed_overtime_detail"`
	TrialPeriod                 null.Int `db:"trial_period" json:"trial_period"`
	TrialPeriodDetail           string   `db:"trial_period_detail" json:"trial_period_detail"`
	EmploymentPeriod            null.Int `db:"employment_period" json:"employment_period"`
	EmploymentPeriodDetail      string   `db:"employment_period_detail" json:"employment_period_detail"`
	HolidayType                 null.Int `db:"holiday_type" json:"holiday_type"`
	HolidayDetail               string   `db:"holiday_detail" json:"holiday_detail"`
	PassiveSmoking              null.Int `db:"passive_smoking" json:"passive_smoking"`
	SelectionFlow               string   `db:"selection_flow" json:"selection_flow"`
	Gender                      null.Int `db:"gender" json:"gender"`
	Nationality                 null.Int `db:"nationality" json:"nationality"`
	FinalEducation              null.Int `db:"final_education" json:"final_education"`
	SchoolLevel                 null.Int `db:"school_level" json:"school_level"`
	MedicalHistory              null.Int `db:"medical_history" json:"medical_history"`
	AgeUnder                    null.Int `db:"age_under" json:"age_under"`
	AgeOver                     null.Int `db:"age_over" json:"age_over"`
	JobChange                   null.Int `db:"job_change" json:"job_change"`
	ShortResignation            null.Int `db:"short_resignation" json:"short_resignation"`
	ShortResignationRemarks     string   `db:"short_resignation_remarks" json:"short_resignation_remarks"`
	SocialExperienceYear        null.Int `db:"social_experience_year" json:"social_experience_year"`
	SocialExperienceMonth       null.Int `db:"social_experience_month" json:"social_experience_month"`
	OtherRequired               string   `db:"other_required" json:"other_required"`
	Appearance                  null.Int `db:"appearance" json:"appearance"`
	Communication               null.Int `db:"communication" json:"communication"`
	Thinking                    null.Int `db:"thinking" json:"thinking"`
	TargetDetail                string   `db:"target_detail" json:"target_detail"`
	ContactManner               null.Int `db:"contact_manner" json:"contact_manner"`
	RequiredManagement          null.Int `db:"required_management" json:"required_management"` //必要マネジメント経験（なし,1〜5名,6〜10名,11名〜30名,31名〜）
	RequiredExperienceJobDetail string   `db:"required_experience_job_detail" json:"required_experience_job_detail"`
	EmploymentInsurance         bool     `db:"employment_insurance" json:"employment_insurance"`     // 雇用保険の有無
	AccidentInsurance           bool     `db:"accident_insurance" json:"accident_insurance"`         // 労災保険の有無
	HealthInsurance             bool     `db:"health_insurance" json:"health_insurance"`             // 健康保険の有無
	PensionInsurance            bool     `db:"pension_insurance" json:"pension_insurance"`           // 厚生年金保険の有無
	RegisterPhase               null.Int `db:"register_phase" json:"register_phase"`                 // 求人の登録状況（0: 本登録, 1: 仮登録）
	StudyCategory               null.Int `db:"study_category" json:"study_category"`                 // 専攻学科の大分類(0:　理系尚可, 1: 理系のみ, 99: 不問)
	WordSkill                   null.Int `db:"word_skill" json:"word_skill"`                         // Wordのスキル
	ExcelSkill                  null.Int `db:"excel_skill" json:"excel_skill"`                       // Excelのスキル
	PowerPointSkill             null.Int `db:"power_point_skill" json:"power_point_skill"`           // PowerPointのスキル
	CorporateSiteURL            string   `db:"corporate_site_url" json:"corporate_site_url"`         // 企業サイトURL
	PostCode                    string   `db:"post_code" json:"post_code"`                           // 郵便番号
	OfficeLocation              string   `db:"office_location" json:"office_location"`               // 本社所在地
	Establishment               string   `db:"establishment" json:"establishment"`                   // 設立年月
	EmployeeNumberSingle        null.Int `db:"employee_number_single" json:"employee_number_single"` // 従業員数（単体）
	EmployeeNumberGroup         null.Int `db:"employee_number_group" json:"employee_number_group"`   // 従業員数（連結）
	PublicOffering              null.Int `db:"public_offering" json:"public_offering"`               // 上場区分
	EarningsYear                null.Int `db:"earnings_year" json:"earnings_year"`                   // 売上年度
	Earnings                    string   `db:"earnings" json:"earnings"`                             // 売上高
	BusinessDetail              string   `db:"business_detail" json:"business_detail"`               // 事業内容

	// 他テーブル
	Targets                   []SendingJobInformationTarget                   `db:"-" json:"targets"`
	Features                  []SendingJobInformationFeature                  `db:"-" json:"features"`
	Prefectures               []SendingJobInformationPrefecture               `db:"-" json:"prefectures"`
	EmploymentStatuses        []SendingJobInformationEmploymentStatus         `db:"-" json:"employment_statuses"`
	WorkCharmPoints           []SendingJobInformationWorkCharmPoint           `db:"-" json:"work_charm_points"`
	RequiredConditions        []SendingJobInformationRequiredCondition        `db:"-" json:"required_conditions"` // 必要条件　複数
	RequiredSocialExperiences []SendingJobInformationRequiredSocialExperience `db:"-" json:"required_social_experiences"`
	Occupations               []SendingJobInformationOccupation               `db:"-" json:"occupations"`
	Industries                []SendingJobInformationIndustry                 `db:"-" json:"industries"`
}

type UpdateSendingJobInformationParam struct {
	SendingBillingAddressID     uint     `db:"sending_billing_address_id" json:"sending_billing_address_id"`
	CompanyName                 string   `db:"company_name" json:"company_name"`
	Title                       string   `db:"title" json:"title"`
	RecruitmentState            null.Int `db:"recruitment_state" json:"recruitment_state"`
	ExpirationDate              string   `db:"expiration_date" json:"expiration_date"`
	Background                  null.Int `db:"background" json:"background"`
	WorkDetail                  string   `db:"work_detail" json:"work_detail"`
	NumberOfHires               null.Int `db:"number_of_hires" json:"number_of_hires"`
	WorkLocation                string   `db:"work_location" json:"work_location"`
	Transfer                    null.Int `db:"transfer" json:"transfer"`
	TransferDetail              string   `db:"transfer_detail" json:"transfer_detail"`
	UnderIncome                 null.Int `db:"under_income" json:"under_income"`
	OverIncome                  null.Int `db:"over_income" json:"over_income"`
	Salary                      string   `db:"salary" json:"salary"`
	Insurance                   string   `db:"insurance" json:"insurance"`
	WorkTime                    string   `db:"work_time" json:"work_time"`
	Overtime                    null.Int `db:"overtime" json:"overtime"`
	OvertimeAverage             string   `db:"overtime_average" json:"overtime_average"`
	FixedOvertime               null.Int `db:"fixed_overtime" json:"fixed_overtime"`
	FixedOvertimePayment        null.Int `db:"fixed_overtime_payment" json:"fixed_overtime_payment"`
	FixedOvertimeDetail         string   `db:"fixed_overtime_detail" json:"fixed_overtime_detail"`
	TrialPeriod                 null.Int `db:"trial_period" json:"trial_period"`
	TrialPeriodDetail           string   `db:"trial_period_detail" json:"trial_period_detail"`
	EmploymentPeriod            null.Int `db:"employment_period" json:"employment_period"`
	EmploymentPeriodDetail      string   `db:"employment_period_detail" json:"employment_period_detail"`
	HolidayType                 null.Int `db:"holiday_type" json:"holiday_type"`
	HolidayDetail               string   `db:"holiday_detail" json:"holiday_detail"`
	PassiveSmoking              null.Int `db:"passive_smoking" json:"passive_smoking"`
	SelectionFlow               string   `db:"selection_flow" json:"selection_flow"`
	Gender                      null.Int `db:"gender" json:"gender"`
	Nationality                 null.Int `db:"nationality" json:"nationality"`
	FinalEducation              null.Int `db:"final_education" json:"final_education"`
	SchoolLevel                 null.Int `db:"school_level" json:"school_level"`
	MedicalHistory              null.Int `db:"medical_history" json:"medical_history"`
	AgeUnder                    null.Int `db:"age_under" json:"age_under"`
	AgeOver                     null.Int `db:"age_over" json:"age_over"`
	JobChange                   null.Int `db:"job_change" json:"job_change"`
	ShortResignation            null.Int `db:"short_resignation" json:"short_resignation"`
	ShortResignationRemarks     string   `db:"short_resignation_remarks" json:"short_resignation_remarks"`
	SocialExperienceYear        null.Int `db:"social_experience_year" json:"social_experience_year"`
	SocialExperienceMonth       null.Int `db:"social_experience_month" json:"social_experience_month"`
	OtherRequired               string   `db:"other_required" json:"other_required"`
	Appearance                  null.Int `db:"appearance" json:"appearance"`
	Communication               null.Int `db:"communication" json:"communication"`
	Thinking                    null.Int `db:"thinking" json:"thinking"`
	TargetDetail                string   `db:"target_detail" json:"target_detail"`
	RequiredExperienceJobDetail string   `db:"required_experience_job_detail" json:"required_experience_job_detail"`
	EmploymentInsurance         bool     `db:"employment_insurance" json:"employment_insurance"`     // 雇用保険の有無
	AccidentInsurance           bool     `db:"accident_insurance" json:"accident_insurance"`         // 労災保険の有無
	HealthInsurance             bool     `db:"health_insurance" json:"health_insurance"`             // 健康保険の有無
	PensionInsurance            bool     `db:"pension_insurance" json:"pension_insurance"`           // 厚生年金保険の有無
	RegisterPhase               null.Int `db:"register_phase" json:"register_phase"`                 // 求人の登録状況（0: 本登録, 1: 仮登録）
	StudyCategory               null.Int `db:"study_category" json:"study_category"`                 // 専攻学科の大分類(0:　理系尚可, 1: 理系のみ, 99: 不問)
	WordSkill                   null.Int `db:"word_skill" json:"word_skill"`                         // Wordのスキル
	ExcelSkill                  null.Int `db:"excel_skill" json:"excel_skill"`                       // Excelのスキル
	PowerPointSkill             null.Int `db:"power_point_skill" json:"power_point_skill"`           // PowerPointのスキル
	CorporateSiteURL            string   `db:"corporate_site_url" json:"corporate_site_url"`         // 企業サイトURL
	PostCode                    string   `db:"post_code" json:"post_code"`                           // 郵便番号
	OfficeLocation              string   `db:"office_location" json:"office_location"`               // 本社所在地
	Establishment               string   `db:"establishment" json:"establishment"`                   // 設立年月
	EmployeeNumberSingle        null.Int `db:"employee_number_single" json:"employee_number_single"` // 従業員数（単体）
	EmployeeNumberGroup         null.Int `db:"employee_number_group" json:"employee_number_group"`   // 従業員数（連結）
	PublicOffering              null.Int `db:"public_offering" json:"public_offering"`               // 上場区分
	EarningsYear                null.Int `db:"earnings_year" json:"earnings_year"`                   // 売上年度
	Earnings                    string   `db:"earnings" json:"earnings"`                             // 売上高
	BusinessDetail              string   `db:"business_detail" json:"business_detail"`               // 事業内容

	// 他テーブル
	Targets                   []SendingJobInformationTarget                   `db:"-" json:"targets"`
	Features                  []SendingJobInformationFeature                  `db:"-" json:"features"`
	Prefectures               []SendingJobInformationPrefecture               `db:"-" json:"prefectures"`
	EmploymentStatuses        []SendingJobInformationEmploymentStatus         `db:"-" json:"employment_statuses"`
	WorkCharmPoints           []SendingJobInformationWorkCharmPoint           `db:"-" json:"work_charm_points"`
	RequiredConditions        []SendingJobInformationRequiredCondition        `db:"-" json:"required_conditions"` // 必要条件　複数
	RequiredSocialExperiences []SendingJobInformationRequiredSocialExperience `db:"-" json:"required_social_experiences"`
	Occupations               []SendingJobInformationOccupation               `db:"-" json:"occupations"`
	Industries                []SendingJobInformationIndustry                 `db:"-" json:"industries"`
}

type DeleteSendingJobInformationParam struct {
	ID uint `json:"id" validate:"required"`
}

type SearchSendingJobInformation struct {
	FreeWord              string
	AgentStaffID          string
	Industries            []null.Int
	Occupations           []null.Int
	Employments           []null.Int
	Prefectures           []null.Int
	UnderIncome           string
	OverIncome            string
	GenderTypes           []null.Int
	AgeUnder              string
	AgeOver               string
	FinalEducationTypes   []null.Int
	SchoolLevelTypes      []null.Int
	StudyCategoryTypes    []null.Int
	NationalityTypes      []null.Int
	MedicalHistoryTypes   []null.Int
	JobChangeTypes        []null.Int
	ShortResignationTypes []null.Int
	AppearanceTypes       []null.Int
	CommunicationTypes    []null.Int
	ThinkingTypes         []null.Int

	RequiredExperienceIndustries  []null.Int
	RequiredExperienceOccupations []null.Int
	RequiredSocialExperienceTypes []null.Int
	RequiredSocialExperienceYear  string
	RequiredSocialExperienceMonth string
	RequiredManagement            string
	RequiredLicenses              []null.Int
	RequiredLanguages             []null.Int
	RequiredLanguageLevels        []null.Int
	RequiredExcelSkills           []null.Int
	RequiredWordSkills            []null.Int
	RequiredPowerPointSkills      []null.Int
	RequiredAnotherPCSkills       []null.Int
	RequiredDevelopmentLanguages  []null.Int
	RequiredDevelopmentOS         []null.Int

	TransferTypes     []null.Int
	HolidayTypes      []null.Int
	CompanyScaleTypes []null.Int
}

func NewSearchSendingJobInformation(
	freeword string,
	agentStaffID string,
	industries []null.Int,
	occupations []null.Int,
	employments []null.Int,
	prefectures []null.Int,
	underIncome string,
	overIncome string,
	genderTypes []null.Int,
	ageUnder string,
	ageOver string,
	finalEducationTypes []null.Int,
	schoolLevelTypes []null.Int,
	studyCategoryTypes []null.Int,
	nationalityTypes []null.Int,
	medicalHistoryTypes []null.Int,
	jobChangeTypes []null.Int,
	shortResignationTypes []null.Int,
	appearanceTypes []null.Int,
	communicationTypes []null.Int,
	thinkingTypes []null.Int,
	requiredExperienceIndustries []null.Int,
	requiredExperienceOccupations []null.Int,
	requiredSocialExperienceTypes []null.Int,
	requiredSocialExperienceYear string,
	requiredSocialExperienceMonth string,
	requiredManagement string,
	requiredLicenses []null.Int,
	requiredLanguages []null.Int,
	requiredLanguageLevels []null.Int,
	requiredExcelSkills []null.Int,
	requiredWordSkills []null.Int,
	requiredPowerPointSkills []null.Int,
	requiredAnotherPCSkills []null.Int,
	requiredDevelopmentLanguages []null.Int,
	requiredDevelopmentOS []null.Int,
	transferTypes []null.Int,
	holidayTypes []null.Int,
	companyScaleTypes []null.Int,
) *SearchSendingJobInformation {
	return &SearchSendingJobInformation{
		FreeWord:                      freeword,
		AgentStaffID:                  agentStaffID,
		Industries:                    industries,
		Occupations:                   occupations,
		Employments:                   employments,
		Prefectures:                   prefectures,
		UnderIncome:                   underIncome,
		OverIncome:                    overIncome,
		GenderTypes:                   genderTypes,
		AgeUnder:                      ageUnder,
		AgeOver:                       ageOver,
		FinalEducationTypes:           finalEducationTypes,
		SchoolLevelTypes:              schoolLevelTypes,
		StudyCategoryTypes:            studyCategoryTypes,
		NationalityTypes:              nationalityTypes,
		MedicalHistoryTypes:           medicalHistoryTypes,
		JobChangeTypes:                jobChangeTypes,
		ShortResignationTypes:         shortResignationTypes,
		AppearanceTypes:               appearanceTypes,
		CommunicationTypes:            communicationTypes,
		ThinkingTypes:                 thinkingTypes,
		RequiredExperienceIndustries:  requiredExperienceIndustries,
		RequiredExperienceOccupations: requiredExperienceOccupations,
		RequiredSocialExperienceTypes: requiredSocialExperienceTypes,
		RequiredSocialExperienceYear:  requiredSocialExperienceYear,
		RequiredSocialExperienceMonth: requiredSocialExperienceMonth,
		RequiredManagement:            requiredManagement,
		RequiredLicenses:              requiredLicenses,
		RequiredLanguages:             requiredLanguages,
		RequiredLanguageLevels:        requiredLanguageLevels,
		RequiredExcelSkills:           requiredExcelSkills,
		RequiredWordSkills:            requiredWordSkills,
		RequiredPowerPointSkills:      requiredPowerPointSkills,
		RequiredAnotherPCSkills:       requiredAnotherPCSkills,
		RequiredDevelopmentLanguages:  requiredDevelopmentLanguages,
		RequiredDevelopmentOS:         requiredDevelopmentOS,
		TransferTypes:                 transferTypes,
		HolidayTypes:                  holidayTypes,
		CompanyScaleTypes:             companyScaleTypes,
	}
}
