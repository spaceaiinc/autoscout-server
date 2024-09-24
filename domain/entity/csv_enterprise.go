package entity

import "gopkg.in/guregu/null.v4"

// csvインポート用の企業・請求先情報
type EnterpriseAndBillingAddress struct {
	// 企業情報
	EnterpriseID         uint       `json:"enterprise_id"`
	CompanyName          string     `json:"company_name" validate:"required"`
	AgentStaffID         uint       `json:"agent_staff_id" validate:"required"`
	CorporateSiteURL     string     `json:"corporate_site_url"`
	Representative       string     `json:"representative"`
	Establishment        string     `json:"establishment"`
	PostCode             string     `json:"post_code"`
	OfficeLocation       string     `json:"office_location"`
	EmployeeNumberSingle null.Int   `json:"employee_number_single"`
	EmployeeNumberGroup  null.Int   `json:"employee_number_group"`
	Capital              string     `json:"capital"`
	PublicOffering       null.Int   `json:"public_offering"`
	EarningsYear         null.Int   `json:"earnings_year"`
	Earnings             string     `json:"earnings"`
	BusinessDetail       string     `json:"business_detail"`
	Industries           []null.Int `json:"industries"`

	// 請求書情報
	BillingAddressID          uint     `json:"billing_address_id"`
	ContractPhase             null.Int `json:"contract_phase"`
	ContractDate              string   `json:"contract_date"`
	PaymentPolicy             string   `json:"payment_policy"`
	BillingAddressCompanyName string   `json:"billing_address_company_name"`
	BillingAddressAddress     string   `json:"billing_address_address"`
	HowToRecommend            string   `json:"how_to_recommend"`
	BillingAddressTitle       string   `json:"billing_address_title"`

	AgentStaffIDForBillingAddress uint                    `json:"agent_staff_id_for_billing_address"`
	HRStaffs                      []BillingAddressHRStaff `json:"hr_staff"`
	RAStaffs                      []BillingAddressRAStaff `json:"ra_staff"`

	// db外
	RecordLine uint `db:"-" json:"record_line"`
}

// csvインポート用の企業・請求先・求人情報
type EnterpriseAndJobInformation struct {
	// 企業情報
	EnterpriseID         uint       `json:"enterprise_id"`
	CompanyName          string     `json:"company_name" validate:"required"`
	AgentStaffID         uint       `json:"agent_staff_id" validate:"required"`
	CorporateSiteURL     string     `json:"corporate_site_url"`
	Representative       string     `json:"representative"`
	Establishment        string     `json:"establishment"`
	PostCode             string     `json:"post_code"`
	OfficeLocation       string     `json:"office_location"`
	EmployeeNumberSingle null.Int   `json:"employee_number_single"`
	EmployeeNumberGroup  null.Int   `json:"employee_number_group"`
	Capital              string     `json:"capital"`
	PublicOffering       null.Int   `json:"public_offering"`
	EarningsYear         null.Int   `json:"earnings_year"`
	Earnings             string     `json:"earnings"`
	BusinessDetail       string     `json:"business_detail"`
	Industries           []null.Int `json:"industries"`

	// 請求書情報
	BillingAddressID          uint     `json:"billing_address_id"`
	ContractPhase             null.Int `json:"contract_phase"`
	ContractDate              string   `json:"contract_date"`
	PaymentPolicy             string   `json:"payment_policy"`
	BillingAddressCompanyName string   `json:"billing_address_company_name"`
	BillingAddressAddress     string   `json:"billing_address_address"`
	HowToRecommend            string   `json:"how_to_recommend"`
	BillingAddressTitle       string   `json:"billing_address_title"`

	AgentStaffIDForBillingAddress uint                    `json:"agent_staff_id_for_billing_address"`
	HRStaffs                      []BillingAddressHRStaff `json:"hr_staff"`
	RAStaffs                      []BillingAddressRAStaff `json:"ra_staff"`

	// 求人情報
	JobInformationID            uint     `json:"job_information_id"`
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
	RequiredManagement          null.Int `db:"required_management" json:"required_management"`                 //必要マネジメント経験（なし,1〜5名,6〜10名,11名〜30名,31名〜）
	SecretMemo                  string   `db:"secret_memo" json:"secret_memo"`                                 //社内限定メモ
	RequiredDocumentsDetail     string   `db:"required_documents_detail" json:"required_documents_detail"`     //推薦時に必要な情報・書類の詳細
	EmploymentInsurance         bool     `db:"employment_insurance" json:"employment_insurance"`               // 雇用保険の有無
	AccidentInsurance           bool     `db:"accident_insurance" json:"accident_insurance"`                   // 労災保険の有無
	HealthInsurance             bool     `db:"health_insurance" json:"health_insurance"`                       // 健康保険の有無
	PensionInsurance            bool     `db:"pension_insurance" json:"pension_insurance"`                     // 厚生年金保険の有無
	RegisterPhase               null.Int `db:"register_phase" json:"register_phase"`                           // 求人の登録状況（0: 本登録, 1: 仮登録）
	StudyCategory               null.Int `db:"study_category" json:"study_category"`                           // 専攻学科の大分類(0:　理系尚可, 1: 理系のみ, 99: 不問)
	DriverLicence               null.Int `db:"driver_licence" json:"driver_licence"`                           // 普通自動車免許（0: 必須, 1: 入社時までに取得必須, 99: 不要）
	WordSkill                   null.Int `db:"word_skill" json:"word_skill"`                                   // Wordのスキル
	ExcelSkill                  null.Int `db:"excel_skill" json:"excel_skill"`                                 // Excelのスキル
	PowerPointSkill             null.Int `db:"power_point_skill" json:"power_point_skill"`                     // PowerPointのスキル
	WorkDetailAfterHiring       string   `db:"work_detail_after_hiring" json:"work_detail_after_hiring"`       // 仕事内容（雇入れ直後）
	WorkDetailScopeOfChange     string   `db:"work_detail_scope_of_change" json:"work_detail_scope_of_change"` // 仕事内容（変更の範囲）
	RequiredExperienceJobDetail string   `db:"required_experience_job_detail" json:"required_experience_job_detail"`
	OfferRate                   null.Int `db:"offer_rate" json:"offer_rate"`                                       // 内定率
	DocumentPassingRate         null.Int `db:"document_passing_rate" json:"document_passing_rate"`                 // 書類通過率
	NumberOfRecentApplications  null.Int `db:"number_of_recent_applications" json:"number_of_recent_applications"` // 直近の応募数
	IsGuaranteedInterview       bool     `db:"is_guaranteed_interview" json:"is_guaranteed_interview"`             // 面接確約フラグ

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

	// db外
	RecordLine uint `db:"-" json:"record_line"`

	// サーカス内のID *被り判定用
	CircusEnterpriseID uint `json:"circus_enterprise_id"`

	// 他媒体でのID(JobInformationExternalID)
	ExternalType null.Int `db:"external_type" json:"external_type"`
	ExternalID   string   `db:"external_id" json:"external_id"`
}

// リストで作成する場合 *ImportEnterpriseJSONで使用
type EnterpriseAndJobInformationListParam struct {
	EnterpriseAndJobInformationList []EnterpriseAndJobInformation `json:"enterprise_and_job_information_list"`
}
