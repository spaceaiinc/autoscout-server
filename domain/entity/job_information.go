package entity

import (
	"time"

	"github.com/google/uuid"

	"gopkg.in/guregu/null.v4"
)

// 求人タイプ
type JobInformationType int64

const (
	TypeAllJobInformation      JobInformationType = iota // 全ての求人
	TypeOwnJobInformation                                // 自社求人
	TypeAllianceJobInformation                           // アライアンス求人
	// TypeHelpJobInformation                               // お助け求人
)

type JobInformation struct {
	ID                         uint      `db:"id" json:"id"`
	UUID                       uuid.UUID `db:"uuid" json:"uuid"`
	EnterpriseID               uint      `db:"enterprise_id" json:"enterprise_id"`
	BillingAddressID           uint      `db:"billing_address_id" json:"billing_address_id"`
	AgentStaffID               uint      `db:"agent_staff_id" json:"agent_staff_id"`
	StaffName                  string    `db:"staff_name" json:"staff_name"`
	AgentName                  string    `db:"agent_name" json:"agent_name"`
	Title                      string    `db:"title" json:"title"`
	RecruitmentState           null.Int  `db:"recruitment_state" json:"recruitment_state"`
	ExpirationDate             string    `db:"expiration_date" json:"expiration_date"`
	WorkDetail                 string    `db:"work_detail" json:"work_detail"`
	NumberOfHires              null.Int  `db:"number_of_hires" json:"number_of_hires"`
	WorkLocation               string    `db:"work_location" json:"work_location"`
	Transfer                   null.Int  `db:"transfer" json:"transfer"`
	TransferDetail             string    `db:"transfer_detail" json:"transfer_detail"`
	UnderIncome                null.Int  `db:"under_income" json:"under_income"`
	OverIncome                 null.Int  `db:"over_income" json:"over_income"`
	Salary                     string    `db:"salary" json:"salary"`
	Insurance                  string    `db:"insurance" json:"insurance"`
	WorkTime                   string    `db:"work_time" json:"work_time"`
	OvertimeAverage            string    `db:"overtime_average" json:"overtime_average"`
	FixedOvertimePayment       null.Int  `db:"fixed_overtime_payment" json:"fixed_overtime_payment"`
	FixedOvertimeDetail        string    `db:"fixed_overtime_detail" json:"fixed_overtime_detail"`
	TrialPeriod                null.Int  `db:"trial_period" json:"trial_period"`
	TrialPeriodDetail          string    `db:"trial_period_detail" json:"trial_period_detail"`
	EmploymentPeriod           null.Int  `db:"employment_period" json:"employment_period"`
	EmploymentPeriodDetail     string    `db:"employment_period_detail" json:"employment_period_detail"`
	HolidayType                null.Int  `db:"holiday_type" json:"holiday_type"`
	HolidayDetail              string    `db:"holiday_detail" json:"holiday_detail"`
	PassiveSmoking             null.Int  `db:"passive_smoking" json:"passive_smoking"`
	SelectionFlow              string    `db:"selection_flow" json:"selection_flow"`
	Gender                     null.Int  `db:"gender" json:"gender"`
	Nationality                null.Int  `db:"nationality" json:"nationality"`
	FinalEducation             null.Int  `db:"final_education" json:"final_education"`
	SchoolLevel                null.Int  `db:"school_level" json:"school_level"`
	MedicalHistory             null.Int  `db:"medical_history" json:"medical_history"`
	AgeUnder                   null.Int  `db:"age_under" json:"age_under"`
	AgeOver                    null.Int  `db:"age_over" json:"age_over"`
	JobChange                  null.Int  `db:"job_change" json:"job_change"`
	ShortResignation           null.Int  `db:"short_resignation" json:"short_resignation"`
	ShortResignationRemarks    string    `db:"short_resignation_remarks" json:"short_resignation_remarks"`
	SocialExperienceYear       null.Int  `db:"social_experience_year" json:"social_experience_year"`
	SocialExperienceMonth      null.Int  `db:"social_experience_month" json:"social_experience_month"`
	Appearance                 null.Int  `db:"appearance" json:"appearance"`
	Communication              null.Int  `db:"communication" json:"communication"`
	Thinking                   null.Int  `db:"thinking" json:"thinking"`
	TargetDetail               string    `db:"target_detail" json:"target_detail"`
	Commission                 null.Int  `db:"commission" json:"commission"`
	CommissionRate             null.Int  `db:"commission_rate" json:"commission_rate"`
	CommissionDetail           string    `db:"commission_detail" json:"commission_detail"`
	RefundPolicy               string    `db:"refund_policy" json:"refund_policy"`
	RequiredManagement         null.Int  `db:"required_management" json:"required_management"`                     //必要マネジメント経験（なし,1〜5名,6〜10名,11名〜30名,31名〜）
	SecretMemo                 string    `db:"secret_memo" json:"secret_memo"`                                     //社内限定メモ
	RequiredDocumentsDetail    string    `db:"required_documents_detail" json:"required_documents_detail"`         //推薦時に必要な情報・書類の詳細
	EmploymentInsurance        bool      `db:"employment_insurance" json:"employment_insurance"`                   // 雇用保険の有無
	AccidentInsurance          bool      `db:"accident_insurance" json:"accident_insurance"`                       // 労災保険の有無
	HealthInsurance            bool      `db:"health_insurance" json:"health_insurance"`                           // 健康保険の有無
	PensionInsurance           bool      `db:"pension_insurance" json:"pension_insurance"`                         // 厚生年金保険の有無
	RegisterPhase              null.Int  `db:"register_phase" json:"register_phase"`                               // 求人の登録状況（0: 本登録, 1: 仮登録）
	IsDeleted                  bool      `db:"is_deleted" json:"is_deleted"`                                       // 論理削除フラグ false: 有効, true: 削除済み
	StudyCategory              null.Int  `db:"study_category" json:"study_category"`                               // 専攻学科の大分類(0:　理系尚可, 1: 理系のみ, 99: 不問)
	DriverLicence              null.Int  `db:"driver_licence" json:"driver_licence"`                               // 普通自動車免許（0: 必須, 1: 入社時までに取得必須, 99: 不要）
	WordSkill                  null.Int  `db:"word_skill" json:"word_skill"`                                       // Wordのスキル
	ExcelSkill                 null.Int  `db:"excel_skill" json:"excel_skill"`                                     // Excelのスキル
	PowerPointSkill            null.Int  `db:"power_point_skill" json:"power_point_skill"`                         // PowerPointのスキル
	IsExternal                 bool      `db:"is_external" json:"is_external"`                                     // 外部求人フラグ true: 外部求人, false: autoscout内求人
	WorkDetailAfterHiring      string    `db:"work_detail_after_hiring" json:"work_detail_after_hiring"`           // 仕事内容（雇入れ直後）
	WorkDetailScopeOfChange    string    `db:"work_detail_scope_of_change" json:"work_detail_scope_of_change"`     // 仕事内容（変更の範囲）
	OfferRate                  null.Int  `db:"offer_rate" json:"offer_rate"`                                       // 内定率
	DocumentPassingRate        null.Int  `db:"document_passing_rate" json:"document_passing_rate"`                 // 書類通過率
	NumberOfRecentApplications null.Int  `db:"number_of_recent_applications" json:"number_of_recent_applications"` // 直近の応募数
	IsGuaranteedInterview      bool      `db:"is_guaranteed_interview" json:"is_guaranteed_interview"`             // 面接確約フラグ
	CreatedAt                  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                  time.Time `db:"updated_at" json:"updated_at"`

	// 企業関連
	CompanyName                 string   `db:"company_name" json:"company_name"`
	PostCode                    string   `db:"post_code" json:"post_code"`
	OfficeLocation              string   `db:"office_location" json:"office_location"`
	EmployeeNumberSingle        null.Int `db:"employee_number_single" json:"employee_number_single"`
	EmployeeNumberGroup         null.Int `db:"employee_number_group" json:"employee_number_group"`
	CorporateSiteURL            string   `db:"corporate_site_url" json:"corporate_site_url"`
	Establishment               string   `db:"establishment" json:"establishment"`
	PublicOffering              null.Int `db:"public_offering" json:"public_offering"`
	Earnings                    string   `db:"earnings" json:"earnings"`
	EarningsYear                null.Int `db:"earnings_year" json:"earnings_year"`
	BusinessDetail              string   `db:"business_detail" json:"business_detail"`
	RequiredExperienceJobDetail string   `db:"required_experience_job_detail" json:"required_experience_job_detail"`

	// 他テーブル
	Targets                        []JobInformationTarget                        `db:"-" json:"targets"`
	Features                       []JobInformationFeature                       `db:"-" json:"features"`
	Prefectures                    []JobInformationPrefecture                    `db:"-" json:"prefectures"`
	EmploymentStatuses             []JobInformationEmploymentStatus              `db:"-" json:"employment_statuses"`
	WorkCharmPoints                []JobInformationWorkCharmPoint                `db:"-" json:"work_charm_points"`
	RequiredConditions             []JobInformationRequiredCondition             `db:"-" json:"required_conditions"`              // 必要条件　複数
	RequiredLicenses               []JobInformationRequiredLicense               `db:"-" json:"required_licenses"`                // 必要資格　複数
	RequiredPCTools                []JobInformationRequiredPCTool                `db:"-" json:"required_pc_tools"`                // 必要PC業務ツール　複数
	RequiredLanguages              []JobInformationRequiredLanguage              `db:"-" json:"required_languages"`               // 必要言語 単数
	RequiredExperienceDevelopments []JobInformationRequiredExperienceDevelopment `db:"-" json:"required_experience_developments"` //必要開発経験　言語・OS各1つずつ
	RequiredExperienceJobs         []JobInformationRequiredExperienceJob         `db:"-" json:"required_experience_jobs"`         // 必要業職種経験　単数
	RequiredSocialExperiences      []JobInformationRequiredSocialExperience      `db:"-" json:"required_social_experiences"`
	SelectionFlowPatterns          []JobInformationSelectionFlowPattern          `db:"-" json:"selection_flow_patterns"`
	HideToAgents                   []JobInformationHideToAgent                   `db:"-" json:"hide_to_agents"` // 非公開エージェント
	Occupations                    []JobInformationOccupation                    `db:"-" json:"occupations"`
	Industries                     []EnterpriseIndustry                          `db:"-" json:"industries"`
	ReferenceMaterial              EnterpriseReferenceMaterial                   `db:"-" json:"reference_materials"`

	// 絞り込み用
	AgentID           uint                              `db:"agent_id" json:"agent_id,omitempty"`
	CommonCondition   JobInformationRequiredCondition   `db:"-" json:"-"` // 共通条件
	PatternConditions []JobInformationRequiredCondition `db:"-" json:"-"` // パターン条件

	// カウント用
	TotalCount null.Int `db:"total_count" json:"total_count"`

	//csv用
	BillingAddressTitle string `json:"billing_address_title,omitempty"` // 請求先タイトル csv該当企業判定用
	RecordLine          uint   `json:"record_line,omitempty"`           // レコード行数 csvインポート用

	// 他社エージェント同士の求人の重複判定用
	IsDuplicate bool `json:"is_duplicate"`

	// 求人検索用
	JobInformationType JobInformationType `db:"job_information_type" json:"job_information_type"`

	// 他媒体でのID(JobInformationExternalID)
	ExternalType null.Int `db:"external_type" json:"-"`
	ExternalID   string   `db:"external_id" json:"-"`

	// タスク用
	HowToRecommend string `db:"how_to_recommend" json:"-"`
}

func NewJobInformation(
	billingAddressID uint,
	title string,
	recruitmentState null.Int,
	expirationDate string,
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
	overtimeAverage string,
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
	appearance null.Int,
	communication null.Int,
	thinking null.Int,
	targetDetail string,
	commission null.Int,
	commissionRate null.Int,
	commissionDetail string,
	refundPolicy string,
	requiredExperienceJobDetail string,
	secretMemo string,
	requiredDocumentsDetail string,
	employmentInsurance bool,
	accidentInsurance bool,
	healthInsurance bool,
	pensionInsurance bool,
	registerPhase null.Int,
	studyCategory null.Int,
	driverLicence null.Int,
	wordSkill null.Int,
	excelSkill null.Int,
	powerPointSkill null.Int,
	isExternal bool,
	workDetailAfterHiring string,
	workDetailScopeOfChange string,
	offerRate null.Int,
	documentPassingRate null.Int,
	numberOfRecentApplications null.Int,
	isGuaranteedInterview bool,
) *JobInformation {
	return &JobInformation{
		BillingAddressID:            billingAddressID,
		Title:                       title,
		RecruitmentState:            recruitmentState,
		ExpirationDate:              expirationDate,
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
		OvertimeAverage:             overtimeAverage,
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
		Appearance:                  appearance,
		Communication:               communication,
		Thinking:                    thinking,
		TargetDetail:                targetDetail,
		Commission:                  commission,
		CommissionRate:              commissionRate,
		CommissionDetail:            commissionDetail,
		RefundPolicy:                refundPolicy,
		RequiredExperienceJobDetail: requiredExperienceJobDetail,
		SecretMemo:                  secretMemo,
		RequiredDocumentsDetail:     requiredDocumentsDetail,
		EmploymentInsurance:         employmentInsurance,
		AccidentInsurance:           accidentInsurance,
		HealthInsurance:             healthInsurance,
		PensionInsurance:            pensionInsurance,
		RegisterPhase:               registerPhase,
		StudyCategory:               studyCategory,
		DriverLicence:               driverLicence,
		WordSkill:                   wordSkill,
		ExcelSkill:                  excelSkill,
		PowerPointSkill:             powerPointSkill,
		IsExternal:                  isExternal,
		WorkDetailAfterHiring:       workDetailAfterHiring,
		WorkDetailScopeOfChange:     workDetailScopeOfChange,
		OfferRate:                   offerRate,
		DocumentPassingRate:         documentPassingRate,
		NumberOfRecentApplications:  numberOfRecentApplications,
		IsGuaranteedInterview:       isGuaranteedInterview,
	}
}

type CreateJobInformationParam struct {
	Title                       string   `db:"title" json:"title"`
	RecruitmentState            null.Int `db:"recruitment_state" json:"recruitment_state"`
	ExpirationDate              string   `db:"expiration_date" json:"expiration_date"`
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
	OvertimeAverage             string   `db:"overtime_average" json:"overtime_average"`
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
	Appearance                  null.Int `db:"appearance" json:"appearance"`
	Communication               null.Int `db:"communication" json:"communication"`
	Thinking                    null.Int `db:"thinking" json:"thinking"`
	TargetDetail                string   `db:"target_detail" json:"target_detail"`
	Commission                  null.Int `db:"commission" json:"commission"`
	CommissionRate              null.Int `db:"commission_rate" json:"commission_rate"`
	CommissionDetail            string   `db:"commission_detail" json:"commission_detail"`
	RefundPolicy                string   `db:"refund_policy" json:"refund_policy"`
	ContactManner               null.Int `db:"contact_manner" json:"contact_manner"`
	RequiredManagement          null.Int `db:"required_management" json:"required_management"` //必要マネジメント経験（なし,1〜5名,6〜10名,11名〜30名,31名〜）
	RequiredExperienceJobDetail string   `db:"required_experience_job_detail" json:"required_experience_job_detail"`
	SecretMemo                  string   `db:"secret_memo" json:"secret_memo"`                                     //社内限定メモ
	RequiredDocumentsDetail     string   `db:"required_documents_detail" json:"required_documents_detail"`         //推薦時に必要な情報・書類の詳細
	EmploymentInsurance         bool     `db:"employment_insurance" json:"employment_insurance"`                   // 雇用保険の有無
	AccidentInsurance           bool     `db:"accident_insurance" json:"accident_insurance"`                       // 労災保険の有無
	HealthInsurance             bool     `db:"health_insurance" json:"health_insurance"`                           // 健康保険の有無
	PensionInsurance            bool     `db:"pension_insurance" json:"pension_insurance"`                         // 厚生年金保険の有無
	RegisterPhase               null.Int `db:"register_phase" json:"register_phase"`                               // 求人の登録状況（0: 本登録, 1: 仮登録）
	StudyCategory               null.Int `db:"study_category" json:"study_category"`                               // 専攻学科の大分類(0:　理系尚可, 1: 理系のみ, 99: 不問)
	DriverLicence               null.Int `db:"driver_licence" json:"driver_licence"`                               // 普通自動車免許（0: 必須, 1: 入社時までに取得必須, 99: 不要）
	WordSkill                   null.Int `db:"word_skill" json:"word_skill"`                                       // Wordのスキル
	ExcelSkill                  null.Int `db:"excel_skill" json:"excel_skill"`                                     // Excelのスキル
	PowerPointSkill             null.Int `db:"power_point_skill" json:"power_point_skill"`                         // PowerPointのスキル
	IsExternal                  bool     `db:"is_external" json:"is_external"`                                     // 外部求人フラグ true: 外部求人, false: autoscout内求人
	WorkDetailAfterHiring       string   `db:"work_detail_after_hiring" json:"work_detail_after_hiring"`           // 仕事内容（雇入れ直後）
	WorkDetailScopeOfChange     string   `db:"work_detail_scope_of_change" json:"work_detail_scope_of_change"`     // 仕事内容（変更の範囲）
	OfferRate                   null.Int `db:"offer_rate" json:"offer_rate"`                                       // 内定率
	DocumentPassingRate         null.Int `db:"document_passing_rate" json:"document_passing_rate"`                 // 書類通過率
	NumberOfRecentApplications  null.Int `db:"number_of_recent_applications" json:"number_of_recent_applications"` // 直近の応募数
	IsGuaranteedInterview       bool     `db:"is_guaranteed_interview" json:"is_guaranteed_interview"`             // 面接確約フラグ

	// 他テーブル
	Targets            []JobInformationTarget            `db:"-" json:"targets"`
	Features           []JobInformationFeature           `db:"-" json:"features"`
	Prefectures        []JobInformationPrefecture        `db:"-" json:"prefectures"`
	EmploymentStatuses []JobInformationEmploymentStatus  `db:"-" json:"employment_statuses"`
	WorkCharmPoints    []JobInformationWorkCharmPoint    `db:"-" json:"work_charm_points"`
	RequiredConditions []JobInformationRequiredCondition `db:"-" json:"required_conditions"` // 必要条件　複数
	// RequiredLicenses               []JobInformationRequiredLicense               `db:"-" json:"required_licenses"`                // 必要資格　複数
	// RequiredPCTools                []JobInformationRequiredPCTool                `db:"-" json:"required_pc_tools"`                // 必要PC業務ツール　複数
	// RequiredLanguages              JobInformationRequiredLanguage                `db:"-" json:"required_languages"`               // 必要言語 単数
	// RequiredExperienceDevelopments []JobInformationRequiredExperienceDevelopment `db:"-" json:"required_experience_developments"` //必要開発経験　言語・OS各1つずつ
	// RequiredExperienceJobs         JobInformationRequiredExperienceJob           `db:"-" json:"required_experience_jobs"`         // 必要業職種経験　単数
	RequiredSocialExperiences []JobInformationRequiredSocialExperience `db:"-" json:"required_social_experiences"`
	SelectionFlowPatterns     []JobInformationSelectionFlowPattern     `db:"-" json:"selection_flow_patterns"`
	HideToAgents              []JobInformationHideToAgent              `db:"-" json:"hide_to_agents"` // 非公開エージェント
	Occupations               []JobInformationOccupation               `db:"-" json:"occupations"`
}

type UpdateJobInformationParam struct {
	BillingAddressID            uint     `db:"billing_address_id" json:"billing_address_id"`
	Title                       string   `db:"title" json:"title"`
	RecruitmentState            null.Int `db:"recruitment_state" json:"recruitment_state"`
	ExpirationDate              string   `db:"expiration_date" json:"expiration_date"`
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
	OvertimeAverage             string   `db:"overtime_average" json:"overtime_average"`
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
	Appearance                  null.Int `db:"appearance" json:"appearance"`
	Communication               null.Int `db:"communication" json:"communication"`
	Thinking                    null.Int `db:"thinking" json:"thinking"`
	TargetDetail                string   `db:"target_detail" json:"target_detail"`
	Commission                  null.Int `db:"commission" json:"commission"`
	CommissionRate              null.Int `db:"commission_rate" json:"commission_rate"`
	CommissionDetail            string   `db:"commission_detail" json:"commission_detail"`
	RefundPolicy                string   `db:"refund_policy" json:"refund_policy"`
	RequiredManagement          null.Int `db:"required_management" json:"required_management"` //必要マネジメント経験（なし,1〜5名,6〜10名,11名〜30名,31名〜）
	RequiredExperienceJobDetail string   `db:"required_experience_job_detail" json:"required_experience_job_detail"`
	SecretMemo                  string   `db:"secret_memo" json:"secret_memo"`                                     //社内限定メモ
	RequiredDocumentsDetail     string   `db:"required_documents_detail" json:"required_documents_detail"`         //推薦時に必要な情報・書類の詳細
	EmploymentInsurance         bool     `db:"employment_insurance" json:"employment_insurance"`                   // 雇用保険の有無
	AccidentInsurance           bool     `db:"accident_insurance" json:"accident_insurance"`                       // 労災保険の有無
	HealthInsurance             bool     `db:"health_insurance" json:"health_insurance"`                           // 健康保険の有無
	PensionInsurance            bool     `db:"pension_insurance" json:"pension_insurance"`                         // 厚生年金保険の有無
	RegisterPhase               null.Int `db:"register_phase" json:"register_phase"`                               // 求人の登録状況（0: 本登録, 1: 仮登録）
	StudyCategory               null.Int `db:"study_category" json:"study_category"`                               // 専攻学科の大分類(0:　理系尚可, 1: 理系のみ, 99: 不問)
	DriverLicence               null.Int `db:"driver_licence" json:"driver_licence"`                               // 普通自動車免許（0: 必須, 1: 入社時までに取得必須, 99: 不要）
	WordSkill                   null.Int `db:"word_skill" json:"word_skill"`                                       // Wordのスキル
	ExcelSkill                  null.Int `db:"excel_skill" json:"excel_skill"`                                     // Excelのスキル
	PowerPointSkill             null.Int `db:"power_point_skill" json:"power_point_skill"`                         // PowerPointのスキル
	IsExternal                  bool     `db:"is_external" json:"is_external"`                                     // 外部求人フラグ true: 外部求人, false: autoscout内求人
	WorkDetailAfterHiring       string   `db:"work_detail_after_hiring" json:"work_detail_after_hiring"`           // 仕事内容（雇入れ直後）
	WorkDetailScopeOfChange     string   `db:"work_detail_scope_of_change" json:"work_detail_scope_of_change"`     // 仕事内容（変更の範囲）
	OfferRate                   null.Int `db:"offer_rate" json:"offer_rate"`                                       // 内定率
	DocumentPassingRate         null.Int `db:"document_passing_rate" json:"document_passing_rate"`                 // 書類通過率
	NumberOfRecentApplications  null.Int `db:"number_of_recent_applications" json:"number_of_recent_applications"` // 直近の応募数
	IsGuaranteedInterview       bool     `db:"is_guaranteed_interview" json:"is_guaranteed_interview"`             // 面接確約フラグ

	// 他テーブル
	Targets            []JobInformationTarget            `db:"-" json:"targets"`
	Features           []JobInformationFeature           `db:"-" json:"features"`
	Prefectures        []JobInformationPrefecture        `db:"-" json:"prefectures"`
	EmploymentStatuses []JobInformationEmploymentStatus  `db:"-" json:"employment_statuses"`
	WorkCharmPoints    []JobInformationWorkCharmPoint    `db:"-" json:"work_charm_points"`
	RequiredConditions []JobInformationRequiredCondition `db:"-" json:"required_conditions"` // 必要条件　複数
	// RequiredLicenses               []JobInformationRequiredLicense               `db:"-" json:"required_licenses"`                // 必要資格　複数
	// RequiredPCTools                []JobInformationRequiredPCTool                `db:"-" json:"required_pc_tools"`                // 必要PC業務ツール　複数
	// RequiredLanguages              JobInformationRequiredLanguage                `db:"-" json:"required_languages"`               // 必要言語 単数
	// RequiredExperienceDevelopments []JobInformationRequiredExperienceDevelopment `db:"-" json:"required_experience_developments"` //必要開発経験　言語・OS各1つずつ
	// RequiredExperienceJobs         JobInformationRequiredExperienceJob           `db:"-" json:"required_experience_jobs"`         // 必要業職種経験　単数
	RequiredSocialExperiences []JobInformationRequiredSocialExperience `db:"-" json:"required_social_experiences"`
	SelectionFlowPatterns     []JobInformationSelectionFlowPattern     `db:"-" json:"selection_flow_patterns"`
	HideToAgents              []JobInformationHideToAgent              `db:"-" json:"hide_to_agents"`
	Occupations               []JobInformationOccupation               `db:"-" json:"occupations"`
}

type DeleteJobInformationParam struct {
	ID uint `json:"id" validate:"required"`
}

type SearchJobInformation struct {
	FreeWord              string
	AgentStaffID          string
	Industries            []null.Int
	Occupations           []null.Int
	Employments           []null.Int
	Prefectures           []null.Int
	UnderIncome           string
	OverIncome            string
	GenderTypes           []null.Int
	Age                   string
	FinalEducationTypes   []null.Int
	SchoolLevelTypes      []null.Int
	StudyCategoryTypes    []null.Int
	NationalityTypes      []null.Int
	MedicalHistoryTypes   []null.Int
	JobChangeTypes        []null.Int
	ShortResignationTypes []null.Int
	DriverLicenceTypes    []null.Int
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
	Features          []null.Int

	OfferRateTypes                  []null.Int
	DocumentPassingRateTypes        []null.Int
	NumberOfRecentApplicationsTypes []null.Int
	IsGuaranteedInterview           bool
}

func NewSearchJobInformation(
	freeword string,
	agentStaffID string,
	industries []null.Int,
	occupations []null.Int,
	employments []null.Int,
	prefectures []null.Int,
	underIncome string,
	overIncome string,
	genderTypes []null.Int,
	age string,
	finalEducationTypes []null.Int,
	schoolLevelTypes []null.Int,
	studyCategoryTypes []null.Int,
	nationalityTypes []null.Int,
	medicalHistoryTypes []null.Int,
	jobChangeTypes []null.Int,
	shortResignationTypes []null.Int,
	driverLicenceTypes []null.Int,
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
	features []null.Int,
	offerRateTypes []null.Int,
	documentPassingRateTypes []null.Int,
	numberOfRecentApplicationsTypes []null.Int,
	isGuaranteedInterview bool,
) *SearchJobInformation {
	return &SearchJobInformation{
		FreeWord:                        freeword,
		AgentStaffID:                    agentStaffID,
		Industries:                      industries,
		Occupations:                     occupations,
		Employments:                     employments,
		Prefectures:                     prefectures,
		UnderIncome:                     underIncome,
		OverIncome:                      overIncome,
		GenderTypes:                     genderTypes,
		Age:                             age,
		FinalEducationTypes:             finalEducationTypes,
		SchoolLevelTypes:                schoolLevelTypes,
		StudyCategoryTypes:              studyCategoryTypes,
		NationalityTypes:                nationalityTypes,
		MedicalHistoryTypes:             medicalHistoryTypes,
		JobChangeTypes:                  jobChangeTypes,
		ShortResignationTypes:           shortResignationTypes,
		DriverLicenceTypes:              driverLicenceTypes,
		AppearanceTypes:                 appearanceTypes,
		CommunicationTypes:              communicationTypes,
		ThinkingTypes:                   thinkingTypes,
		RequiredExperienceIndustries:    requiredExperienceIndustries,
		RequiredExperienceOccupations:   requiredExperienceOccupations,
		RequiredSocialExperienceTypes:   requiredSocialExperienceTypes,
		RequiredSocialExperienceYear:    requiredSocialExperienceYear,
		RequiredSocialExperienceMonth:   requiredSocialExperienceMonth,
		RequiredManagement:              requiredManagement,
		RequiredLicenses:                requiredLicenses,
		RequiredLanguages:               requiredLanguages,
		RequiredLanguageLevels:          requiredLanguageLevels,
		RequiredExcelSkills:             requiredExcelSkills,
		RequiredWordSkills:              requiredWordSkills,
		RequiredPowerPointSkills:        requiredPowerPointSkills,
		RequiredAnotherPCSkills:         requiredAnotherPCSkills,
		RequiredDevelopmentLanguages:    requiredDevelopmentLanguages,
		RequiredDevelopmentOS:           requiredDevelopmentOS,
		TransferTypes:                   transferTypes,
		HolidayTypes:                    holidayTypes,
		CompanyScaleTypes:               companyScaleTypes,
		Features:                        features,
		OfferRateTypes:                  offerRateTypes,
		DocumentPassingRateTypes:        documentPassingRateTypes,
		NumberOfRecentApplicationsTypes: numberOfRecentApplicationsTypes,
		IsGuaranteedInterview:           isGuaranteedInterview,
	}
}
