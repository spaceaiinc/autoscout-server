package entity

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type DiagnosisParam struct {
	Gender                   null.Int               `json:"gender"`     // 性別 0: 男性, 1: 女性
	Birthyear                string                 `json:"birthyear"`  // 生年月日（年）
	Birthmonth               string                 `json:"birthmonth"` // 生年月日（月）
	Birthday                 string                 `json:"birthday"`   // 生年月日（日）
	Prefecture               null.Int               `json:"prefecture"`
	FirstLanguage            null.Int               `json:"first_language"`
	SchoolCategory           null.Int               `json:"school_category"`
	SchoolName               string                 `json:"school_name"`
	Subject                  string                 `json:"subject"`
	GraduationYear           string                 `json:"graduation_year"`
	GraduationMonth          string                 `json:"graduation_month"`
	CompanyNum               null.Int               `json:"company_num"` // 経験社数 0: 1社, 1: 2社, 2: 3社, 3: 4社, 4: 5社以上
	CompanyName              string                 `json:"company_name"`
	JoiningYear              string                 `json:"joining_year"`
	JoiningMonth             string                 `json:"joining_month"`
	RetireYear               string                 `json:"retire_year"`
	RetireMonth              string                 `json:"retire_month"`
	IsRetire                 bool                   `json:"is_retire"`
	Industries               []ExperienceIndustry   `json:"industries"`
	EmployeeNumber           null.Int               `json:"employee_number"`
	EmployeeStatus           null.Int               `json:"employee_status"`
	Income                   null.Int               `json:"income"`
	ExperienceOccupations    []ExperienceOccupation `json:"experience_occupations"` // 経験職種
	AllExperienceOccupations []ExperienceOccupation `json:"all_experience_occupations"`
	PRPoint                  string                 `json:"pr_point"`
	JobSummary               string                 `json:"job_summary"`
	Languages                []Language             `json:"languages"`
	DriversLicense           null.Int               `json:"drivers_license"`
	Licenses                 []null.Int             `json:"licenses"`

	// LPに診断結果を返す用
	JobInformationList []*JobInformationForDiagnosis `json:"job_information_list"`
}

type JobInformationForDiagnosis struct {
	ID                    uint      `db:"id" json:"id"`
	UUID                  uuid.UUID `db:"uuid" json:"uuid"`
	BillingAddressID      uint      `db:"billing_address_id" json:"billing_address_id"`
	UnderIncome           null.Int  `db:"under_income" json:"under_income"`
	OverIncome            null.Int  `db:"over_income" json:"over_income"`
	Gender                null.Int  `db:"gender" json:"gender"`
	Nationality           null.Int  `db:"nationality" json:"nationality"`
	FinalEducation        null.Int  `db:"final_education" json:"final_education"`
	AgeUnder              null.Int  `db:"age_under" json:"age_under"`
	AgeOver               null.Int  `db:"age_over" json:"age_over"`
	JobChange             null.Int  `db:"job_change" json:"job_change"`
	DriverLicence         null.Int  `db:"driver_licence" json:"driver_licence"`                   // 普通自動車免許（0: 必須, 1: 入社時までに取得必須, 99: 不要）
	IsGuaranteedInterview bool      `db:"is_guaranteed_interview" json:"is_guaranteed_interview"` // 面接確約フラグ

	// 他テーブル
	Features           []JobInformationFeature           `db:"-" json:"features"`
	Prefectures        []JobInformationPrefecture        `db:"-" json:"prefectures"`
	RequiredConditions []JobInformationRequiredCondition `db:"-" json:"required_conditions"` // 必要条件　複数
	RequiredLicenses   []JobInformationRequiredLicense   `db:"-" json:"required_licenses"`   // 必要資格　複数
	RequiredLanguages  []JobInformationRequiredLanguage  `db:"-" json:"required_languages"`  // 必要言語 単数
}

func NewJobInformationForDiagnosis(
	billingAddressID uint,
	underIncome null.Int,
	overIncome null.Int,
	gender null.Int,
	nationality null.Int,
	finalEducation null.Int,
	ageUnder null.Int,
	ageOver null.Int,
	jobChange null.Int,
	driverLicence null.Int,
	isGuaranteedInterview bool,
) *JobInformationForDiagnosis {
	return &JobInformationForDiagnosis{
		BillingAddressID:      billingAddressID,
		UnderIncome:           underIncome,
		OverIncome:            overIncome,
		Gender:                gender,
		Nationality:           nationality,
		FinalEducation:        finalEducation,
		AgeUnder:              ageUnder,
		AgeOver:               ageOver,
		JobChange:             jobChange,
		DriverLicence:         driverLicence,
		IsGuaranteedInterview: isGuaranteedInterview,
	}
}

type ExperienceOccupation struct {
	Occupation     null.Int `json:"occupation" validate:"required"`
	ExperienceYear null.Int `json:"experience_year" validate:"required"`
}

type Language struct {
	LanguageType  null.Int `json:"language_type" validate:"required"`
	LanguageLevel null.Int `json:"language_level" validate:"required"`
}

type ExperienceIndustry struct {
	Industry        null.Int `json:"industry" validate:"required"`
	ExperienceMonth null.Int `json:"experience_month" validate:"required"`
}

type SearchMatchingJobListParam struct {
	JobSeekerUUID         uuid.UUID  `json:"job_seeker_uuid"`         // 求職者UUID
	PageNumber            null.Int   `json:"page_number"`             // ページ番号
	DesiredIndustries     []null.Int `json:"desired_industries"`      // 希望業種
	DesiredOccupations    []null.Int `json:"desired_occupations"`     // 希望職種
	Prefectures           []null.Int `json:"prefectures"`             // 希望勤務地
	Features              []null.Int `json:"features"`                // 特徴
	Income                null.Int   `json:"income"`                  // 希望年収
	IsGuaranteedInterview bool       `json:"is_guaranteed_interview"` // 面接確約

	// 求職者マイページに診断結果を返す用
	JobListingList []*JobListing `json:"job_listing_list"`
}

type InterestedTypeJobListParam struct {
	JobSeekerUUID  uuid.UUID `json:"job_seeker_uuid"` // 求職者UUID
	PageNumber     null.Int  `json:"page_number"`     // ページ番号
	InterestedType null.Int  `json:"interested_type"` // { エントリー希望: 0; 興味あり: 1 }
}

type CreateJobSeekerFromLPParam struct {
	Email                    string                 `json:"email" validate:"required"`
	Password                 string                 `json:"password" validate:"required"`
	LastName                 string                 `json:"last_name" validate:"required"`
	FirstName                string                 `json:"first_name" validate:"required"`
	LastFurigana             string                 `json:"last_furigana" validate:"required"`
	FirstFurigana            string                 `json:"first_furigana" validate:"required"`
	Gender                   null.Int               `json:"gender" validate:"required"`
	Birthyear                string                 `json:"birthyear" validate:"required"`
	Birthmonth               string                 `json:"birthmonth" validate:"required"`
	Birthday                 string                 `json:"birthday" validate:"required"`
	Prefecture               null.Int               `json:"prefecture" validate:"required"`
	FirstLanguage            null.Int               `json:"first_language" validate:"required"`
	SchoolCategory           null.Int               `json:"school_category" validate:"required"`
	SchoolName               string                 `json:"school_name" validate:"required"`
	Subject                  string                 `json:"subject"`
	GraduationYear           string                 `json:"graduation_year" validate:"required"`
	GraduationMonth          string                 `json:"graduation_month" validate:"required"`
	CompanyNum               null.Int               `json:"company_num" validate:"required"`
	CompanyName              string                 `json:"company_name" validate:"required"`
	JoiningYear              string                 `json:"joining_year" validate:"required"`
	JoiningMonth             string                 `json:"joining_month" validate:"required"`
	RetireYear               string                 `json:"retire_year"`
	RetireMonth              string                 `json:"retire_month"`
	IsRetire                 bool                   `json:"is_retire"`
	Industries               []ExperienceIndustry   `json:"industries" validate:"required"`
	EmployeeNumber           null.Int               `json:"employee_number" validate:"required"`
	EmployeeStatus           null.Int               `json:"employee_status" validate:"required"`
	Income                   null.Int               `json:"income" validate:"required"`
	ExperienceOccupations    []ExperienceOccupation `json:"experience_occupations" validate:"required"`
	AllExperienceOccupations []ExperienceOccupation `json:"all_experience_occupations" validate:"required"`
	PRPoint                  string                 `json:"pr_point"`
	JobSummary               string                 `json:"job_summary"`
	Languages                []Language             `json:"languages"`
	DriversLicense           null.Int               `json:"drivers_license" validate:"required"`
	Licenses                 []null.Int             `json:"licenses"`
}

type UpdateJobSeekerPhoneFromLPParam struct {
	UUID        uuid.UUID `json:"uuid" validate:"required"`
	PhoneNumber string    `json:"phone_number" validate:"required"`
}

type UpdateJobSeekerDesiredFromLPParam struct {
	UUID                 uuid.UUID  `json:"uuid" validate:"required"`
	DesiredWorkLocations []null.Int `json:"desired_work_locations" validate:"required"`
	DesiredOccupations   []null.Int `json:"desired_occupations" validate:"required"`
	DesiredIncome        null.Int   `json:"desired_income" validate:"required"`
}

// パスワードリセット
type ResetPasswordFromLPParam struct {
	Password           string `json:"password" validate:"required"`
	ResetPasswordToken string `json:"reset_password_token" validate:"required"`
}

// パスワードリセットメール送信
type SendJobSeekerResetPasswordEmailFromLPParam struct {
	Email string `json:"email" validate:"required"`
}

// LPからのお問い合わせ
type SendContactFromLPParam struct {
	Name          string `json:"name" validate:"required"`
	CompanyName   string `json:"company_name"`
	Email         string `json:"email" validate:"required"`
	ContacMessage string `json:"contact_message" validate:"required"`
}