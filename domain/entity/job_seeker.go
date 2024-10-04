package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// 求職者タイプ
type JobSeekerType int64

const (
	TypeAllJobSeeker      JobSeekerType = iota // すべての求職者
	TypeOwnJobSeeker                           // 自社求職者
	TypeAllianceJobSeeker                      // アライアンス求職者
	// TypeHelpJobSeeker                          // お助け求職者
)

type JobSeeker struct {
	ID                      uint      `db:"id" json:"id"`
	UUID                    uuid.UUID `db:"uuid" json:"uuid"`
	AgentID                 uint      `db:"agent_id" json:"agent_id"`
	AgentStaffID            null.Int  `db:"agent_staff_id" json:"agent_staff_id"`
	LineID                  string    `db:"line_id" json:"-"`
	LineActive              bool      `db:"line_active" json:"line_active"`
	AgentName               string    `db:"agent_name" json:"agent_name"`
	StaffName               string    `db:"staff_name" json:"staff_name"`
	StaffEmail              string    `db:"staff_email" json:"staff_email"`
	StaffPhoneNumber        string    `db:"staff_phone_number" json:"staff_phone_number"`
	UserStatus              null.Int  `db:"user_status" json:"user_status"`
	LastName                string    `db:"last_name" json:"last_name"`
	FirstName               string    `db:"first_name" json:"first_name"`
	LastFurigana            string    `db:"last_furigana" json:"last_furigana"`
	FirstFurigana           string    `db:"first_furigana" json:"first_furigana"`
	Gender                  null.Int  `db:"gender" json:"gender"`
	GenderRemarks           string    `db:"gender_remarks" json:"gender_remarks"`
	Birthday                string    `db:"birthday" json:"birthday"`
	Spouse                  null.Int  `db:"spouse" json:"spouse"`
	SupportObligation       null.Int  `db:"support_obligation" json:"support_obligation"`
	Dependents              null.Int  `db:"dependents" json:"dependents"`
	PhoneNumber             string    `db:"phone_number" json:"phone_number"`
	Email                   string    `db:"email" json:"email"`
	EmergencyPhoneNumber    string    `db:"emergency_phone_number" json:"emergency_phone_number"`
	PostCode                string    `db:"post_code" json:"post_code"`
	Prefecture              null.Int  `db:"prefecture" json:"prefecture"`
	Address                 string    `db:"address" json:"address"`
	AddressFurigana         string    `db:"address_furigana" json:"address_furigana"`
	StateOfEmployment       null.Int  `db:"state_of_employment" json:"state_of_employment"`
	JobSummary              string    `db:"job_summary" json:"job_summary"`
	HistorySupplement       string    `db:"history_supplement" json:"history_supplement"`
	ResearchContent         string    `db:"research_content" json:"research_content"`
	JoinCompanyPeriod       null.Int  `db:"join_company_period" json:"join_company_period"`
	JobChange               null.Int  `db:"job_change" json:"job_change"`
	AnnualIncome            null.Int  `db:"annual_income" json:"annual_income"`
	DesiredAnnualIncome     null.Int  `db:"desired_annual_income" json:"desired_annual_income"`
	Transfer                null.Int  `db:"transfer" json:"transfer"`
	TransferRequirement     string    `db:"transfer_requirement" json:"transfer_requirement"`
	ShortResignation        null.Int  `db:"short_resignation" json:"short_resignation"`
	ShortResignationRemarks string    `db:"short_resignation_remarks" json:"short_resignation_remarks"`
	MedicalHistory          null.Int  `db:"medical_history" json:"medical_history"`
	Nationality             null.Int  `db:"nationality" json:"nationality"`
	Appearance              null.Int  `db:"appearance" json:"appearance"`
	Communication           null.Int  `db:"communication" json:"communication"`
	Thinking                null.Int  `db:"thinking" json:"thinking"`
	RecommendationProfile   string    `db:"recommendation_profile" json:"recommendation_profile"`
	CandidProfile           string    `db:"candid_profile" json:"candid_profile"`
	SecretMemo              string    `db:"secret_memo" json:"secret_memo"`
	JobHuntingState         null.Int  `db:"job_hunting_state" json:"job_hunting_state"`
	RecommendReason         string    `db:"recommend_reason" json:"recommend_reason"`
	Phase                   null.Int  `db:"phase" json:"phase"`
	CanViewMatchingJob      bool      `db:"can_view_matching_job" json:"can_view_matching_job"`

	// 面談調整タスク
	InterviewDate      time.Time `db:"interview_date" json:"interview_date"`
	FirstInterviewDate time.Time `db:"first_interview_date" json:"first_interview_date"`

	Agreement       bool     `db:"agreement" json:"agreement"`
	RegisterPhase   null.Int `db:"register_phase" json:"register_phase"`       // 求職者の登録状況（0: 本登録, 1: 仮登録）
	StudyCategory   null.Int `db:"study_category" json:"study_category"`       // 専攻学科の大分類(0:　理系 1: 文系)
	WordSkill       null.Int `db:"word_skill" json:"word_skill"`               // Wordのスキル
	ExcelSkill      null.Int `db:"excel_skill" json:"excel_skill"`             // Excelのスキル
	PowerPointSkill null.Int `db:"power_point_skill" json:"power_point_skill"` // PowerPointのスキル
	ActivityMemo    string   `db:"activity_memo" json:"activity_memo"`

	// 流入経路
	InflowChannelID null.Int `db:"inflow_channel_id" json:"inflow_channel_id"` // agent_inflow_channel_optionsテーブルのID
	ChannelName     string   `db:"channel_name" json:"channel_name"`           // 流入経路の名前 *求職者の詳細画面で使用

	NationalityRemarks    string `db:"nationality_remarks" json:"nationality_remarks"`         // 国籍 外国籍を選択→国籍備考（フリーテキスト）を表示
	MedicalHistoryRemarks string `db:"medical_history_remarks" json:"medical_history_remarks"` // 既往歴 ありを選択→既往歴備考（フリーテキスト）を表示
	AcceptancePoints      string `db:"acceptance_points" json:"acceptance_points"`             // 応募承諾のポイント

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// 他テーブル
	StudentHistories     []JobSeekerStudentHistory      `json:"student_histories"`      // 求職者の学歴情報
	WorkHistories        []JobSeekerWorkHistory         `json:"work_histories"`         // 求職者の職歴情報
	Licenses             []JobSeekerLicense             `json:"licenses"`               // 求職者の所持資格
	SelfPromotions       []JobSeekerSelfPromotion       `json:"self_promotions"`        // 求職者の自己PR
	Documents            JobSeekerDocument              `json:"documents"`              // 求職者の資料
	DesiredIndustries    []JobSeekerDesiredIndustry     `json:"desired_industries"`     // 求職者の希望業界
	DesiredOccupations   []JobSeekerDesiredOccupation   `json:"desired_occupations"`    // 求職者の希望職種
	DesiredWorkLocations []JobSeekerDesiredWorkLocation `json:"desired_work_locations"` // 求職者の希望勤務地
	DevelopmentSkills    []JobSeekerDevelopmentSkill    `json:"development_skills"`     // 求職者の開発スキル
	LanguageSkills       []JobSeekerLanguageSkill       `json:"language_skills"`        // 求職者の言語スキル
	PCTools              []JobSeekerPCTool              `json:"pc_tools"`               // 求職者のPCツール
	DesiredHolidayTypes  []JobSeekerDesiredHolidayType  `json:"desired_holiday_types"`  // 求職者の休日タイプ
	HideToAgents         []JobSeekerHideToAgent         `json:"hide_to_agents"`         // 求職者の非表示エージェント
	DesiredCompanyScales []JobSeekerDesiredCompanyScale `json:"desired_company_scales"` // 求職者の希望企業規模

	// 他社エージェント同士の求職者の重複判定用
	IsDuplicate bool `json:"is_duplicate"`

	Question string `db:"question" json:"question"`

	// ゲストページ用の認証
	Password           string `db:"password" json:"-"`             // パスワード
	ResetPasswordToken string `db:"reset_password_token" json:"-"` // パスワードリセット用のトークン メールごとに変更

	ExternalID string `db:"external_id" json:"external_id"` // 外部ID

	VincereJobTitle1 string `db:"-" json:"-"` // Vincereの求人タイトル1
	VincereJobTitle2 string `db:"-" json:"-"` // Vincereの求人タイトル2
	VincereJobTitle3 string `db:"-" json:"-"` // Vincereの求人タイトル3
}

func NewJobSeeker(
	agentID uint,
	agentStaffID null.Int,
	userStatus null.Int,
	lastName string,
	firstName string,
	lastFurigana string,
	firstFurigana string,
	gender null.Int,
	genderRemarks string,
	birthday string,
	spouse null.Int,
	supportObligation null.Int,
	dependents null.Int,
	phoneNumber string,
	email string,
	emergencyPhoneNumber string,
	postCode string,
	prefecture null.Int,
	address string,
	addressFurigana string,
	stateOfEmployment null.Int,
	jobSummary string,
	historySupplement string,
	researchContent string,
	joinCompanyPeriod null.Int,
	jobChange null.Int,
	annualIncome null.Int,
	desiredAnnualIncome null.Int,
	transfer null.Int,
	transferRequirement string,
	shortResignation null.Int,
	shortResignationRemarks string,
	medicalHistory null.Int,
	nationality null.Int,
	appearance null.Int,
	communication null.Int,
	thinking null.Int,
	recommendationProfile string,
	candidProfile string,
	secretMemo string,
	jobHuntingState null.Int,
	recommendReason string,
	phase null.Int,
	interviewDate time.Time,
	registerPhase null.Int,
	studyCategory null.Int,
	wordSkill null.Int,
	excelSkill null.Int,
	powerPointSkill null.Int,
	inflowChannelID null.Int,
	nationalityRemarks string,
	medicalHistoryRemarks string,
	acceptancePoints string,
) *JobSeeker {
	return &JobSeeker{
		AgentID:                 agentID,
		AgentStaffID:            agentStaffID,
		UserStatus:              userStatus,
		LastName:                lastName,
		FirstName:               firstName,
		LastFurigana:            lastFurigana,
		FirstFurigana:           firstFurigana,
		Gender:                  gender,
		GenderRemarks:           genderRemarks,
		Birthday:                birthday,
		Spouse:                  spouse,
		SupportObligation:       supportObligation,
		Dependents:              dependents,
		PhoneNumber:             phoneNumber,
		Email:                   email,
		EmergencyPhoneNumber:    emergencyPhoneNumber,
		PostCode:                postCode,
		Prefecture:              prefecture,
		Address:                 address,
		AddressFurigana:         addressFurigana,
		StateOfEmployment:       stateOfEmployment,
		JobSummary:              jobSummary,
		HistorySupplement:       historySupplement,
		ResearchContent:         researchContent,
		JoinCompanyPeriod:       joinCompanyPeriod,
		JobChange:               jobChange,
		AnnualIncome:            annualIncome,
		DesiredAnnualIncome:     desiredAnnualIncome,
		Transfer:                transfer,
		TransferRequirement:     transferRequirement,
		ShortResignation:        shortResignation,
		ShortResignationRemarks: shortResignationRemarks,
		MedicalHistory:          medicalHistory,
		Nationality:             nationality,
		Appearance:              appearance,
		Communication:           communication,
		Thinking:                thinking,
		RecommendationProfile:   recommendationProfile,
		CandidProfile:           candidProfile,
		SecretMemo:              secretMemo,
		JobHuntingState:         jobHuntingState,
		RecommendReason:         recommendReason,
		Phase:                   phase,
		InterviewDate:           interviewDate,
		RegisterPhase:           registerPhase,
		StudyCategory:           studyCategory,
		WordSkill:               wordSkill,
		ExcelSkill:              excelSkill,
		PowerPointSkill:         powerPointSkill,
		InflowChannelID:         inflowChannelID,
		NationalityRemarks:      nationalityRemarks,
		MedicalHistoryRemarks:   medicalHistoryRemarks,
		AcceptancePoints:        acceptancePoints,
	}
}

type CreateOrUpdateJobSeekerParam struct {
	AgentID                 uint      `db:"agent_id" json:"agent_id" validate:"required"`
	AgentStaffID            null.Int  `db:"agent_staff_id" json:"agent_staff_id"`
	UserStatus              null.Int  `db:"user_status" json:"user_status"`
	LastName                string    `db:"last_name" json:"last_name"`
	FirstName               string    `db:"first_name" json:"first_name"`
	LastFurigana            string    `db:"last_furigana" json:"last_furigana"`
	FirstFurigana           string    `db:"first_furigana" json:"first_furigana"`
	Gender                  null.Int  `db:"gender" json:"gender"`
	GenderRemarks           string    `db:"gender_remarks" json:"gender_remarks"`
	Birthday                string    `db:"birthday" json:"birthday"`
	Spouse                  null.Int  `db:"spouse" json:"spouse"`
	SupportObligation       null.Int  `db:"support_obligation" json:"support_obligation"`
	Dependents              null.Int  `db:"dependents" json:"dependents"`
	PhoneNumber             string    `db:"phone_number" json:"phone_number"`
	Email                   string    `db:"email" json:"email"`
	EmergencyPhoneNumber    string    `db:"emergency_phone_number" json:"emergency_phone_number"`
	PostCode                string    `db:"post_code" json:"post_code"`
	Prefecture              null.Int  `db:"prefecture" json:"prefecture"`
	Address                 string    `db:"address" json:"address"`
	AddressFurigana         string    `db:"address_furigana" json:"address_furigana"`
	StateOfEmployment       null.Int  `db:"state_of_employment" json:"state_of_employment"`
	JobSummary              string    `db:"job_summary" json:"job_summary"`
	HistorySupplement       string    `db:"history_supplement" json:"history_supplement"`
	ResearchContent         string    `db:"research_content" json:"research_content"`
	JoinCompanyPeriod       null.Int  `db:"join_company_period" json:"join_company_period"`
	JobChange               null.Int  `db:"job_change" json:"job_change"`
	AnnualIncome            null.Int  `db:"annual_income" json:"annual_income"`
	DesiredAnnualIncome     null.Int  `db:"desired_annual_income" json:"desired_annual_income"`
	Transfer                null.Int  `db:"transfer" json:"transfer"`
	TransferRequirement     string    `db:"transfer_requirement" json:"transfer_requirement"`
	ShortResignation        null.Int  `db:"short_resignation" json:"short_resignation"`
	ShortResignationRemarks string    `db:"short_resignation_remarks" json:"short_resignation_remarks"`
	MedicalHistory          null.Int  `db:"medical_history" json:"medical_history"`
	Nationality             null.Int  `db:"nationality" json:"nationality"`
	Appearance              null.Int  `db:"appearance" json:"appearance"`
	Communication           null.Int  `db:"communication" json:"communication"`
	Thinking                null.Int  `db:"thinking" json:"thinking"`
	RecommendationProfile   string    `db:"recommendation_profile" json:"recommendation_profile"`
	CandidProfile           string    `db:"candid_profile" json:"candid_profile"`
	SecretMemo              string    `db:"secret_memo" json:"secret_memo"`
	JobHuntingState         null.Int  `db:"job_hunting_state" json:"job_hunting_state"`
	RecommendReason         string    `db:"recommend_reason" json:"recommend_reason"`
	Phase                   null.Int  `db:"phase" json:"phase"`
	InterviewDate           time.Time `db:"interview_date" json:"interview_date"`
	// Agreement                  bool      `db:"agreement" json:"agreement"`
	RegisterPhase         null.Int  `db:"register_phase" json:"register_phase"`       // 求職者の登録状況（0: 本登録, 1: 仮登録）
	StudyCategory         null.Int  `db:"study_category" json:"study_category"`       // 専攻学科の大分類(0:　理系 1: 文系)
	WordSkill             null.Int  `db:"word_skill" json:"word_skill"`               // Wordのスキル
	ExcelSkill            null.Int  `db:"excel_skill" json:"excel_skill"`             // Excelのスキル
	PowerPointSkill       null.Int  `db:"power_point_skill" json:"power_point_skill"` // PowerPointのスキル
	UUID                  uuid.UUID `db:"uuid" json:"uuid"`
	InflowChannelID       null.Int  `db:"inflow_channel_id" json:"inflow_channel_id"`             // agent_inflow_channel_optionsテーブルのID
	NationalityRemarks    string    `db:"nationality_remarks" json:"nationality_remarks"`         // 国籍 外国籍を選択→国籍備考（フリーテキスト）を表示
	MedicalHistoryRemarks string    `db:"medical_history_remarks" json:"medical_history_remarks"` // 既往歴 ありを選択→既往歴備考（フリーテキスト）を表示
	AcceptancePoints      string    `db:"acceptance_points" json:"acceptance_points"`             // 応募承諾のポイント

	// 他テーブル
	StudentHistories     []JobSeekerStudentHistory      `json:"student_histories"`      // 求職者の学歴情報
	WorkHistories        []JobSeekerWorkHistory         `json:"work_histories"`         // 求職者の職歴情報
	Licenses             []JobSeekerLicense             `json:"licenses"`               // 求職者の所持資格
	SelfPromotions       []JobSeekerSelfPromotion       `json:"self_promotions"`        // 求職者の自己PR
	DesiredIndustries    []JobSeekerDesiredIndustry     `json:"desired_industries"`     // 求職者の希望業界
	DesiredOccupations   []JobSeekerDesiredOccupation   `json:"desired_occupations"`    // 求職者の希望職種
	DesiredWorkLocations []JobSeekerDesiredWorkLocation `json:"desired_work_locations"` // 求職者の希望勤務地
	DevelopmentSkills    []JobSeekerDevelopmentSkill    `json:"development_skills"`     // 求職者の開発スキル
	LanguageSkills       []JobSeekerLanguageSkill       `json:"language_skills"`        // 求職者の言語スキル
	PCTools              []JobSeekerPCTool              `json:"pc_tools"`               // 求職者のPCツール
	DesiredHolidayTypes  []JobSeekerDesiredHolidayType  `json:"desired_holiday_types"`  // 求職者の休日タイプ
	HideToAgents         []JobSeekerHideToAgent         `json:"hide_to_agents"`         // 求職者の非表示エージェント
	DesiredCompanyScales []JobSeekerDesiredCompanyScale `json:"desired_company_scales"` // 求職者の希望企業規模
}

type DeleteJobSeekerParam struct {
	ID uint `json:"id" validate:"required"`
}

type SearchJobSeeker struct {
	FreeWord                 string
	AgentStaffID             string
	PhaseTypes               []null.Int
	GenderTypes              []null.Int
	AgeUnder                 string
	AgeOver                  string
	DesiredIndustries        []null.Int
	DesiredOccupations       []null.Int
	DesiredWorkLocations     []null.Int
	FinalEducationTypes      []null.Int
	StudyCategoryTypes       []null.Int
	SchoolLevelTypes         []null.Int
	NationalityTypes         []null.Int
	JobChangeTypes           []null.Int
	ShortResignationTypes    []null.Int
	UnderIncome              string
	OverIncome               string
	DesiredTransferTypes     []null.Int
	DesiredHolidayTypes      []null.Int
	DesiredCompanyScaleTypes []null.Int
	ExperienceIndustries     []null.Int
	ExperienceOccupations    []null.Int
	SocialExperiences        []null.Int
	Management               string
	Licenses                 []null.Int
	Languages                []null.Int
	ExcelSkills              []null.Int
	WordSkills               []null.Int
	PowerPointSkills         []null.Int
	AnotherPCSkills          []null.Int
	DevelopmentLanguages     []null.Int
	DevelopmentOS            []null.Int
	AppearanceTypes          []null.Int
	CommunicationTypes       []null.Int
	ThinkingTypes            []null.Int
}

func NewSearchJobSeeker(
	freeword string,
	agentStaffID string,
	phaseTypes []null.Int,
	genderTypes []null.Int,
	ageUnder string,
	ageOver string,
	desiredIndustries []null.Int,
	desiredOccupations []null.Int,
	desiredWorkLocations []null.Int,
	finalEducationTypes []null.Int,
	studyCategoryTypes []null.Int,
	schoolLevelTypes []null.Int,
	nationalityTypes []null.Int,
	jobChangeTypes []null.Int,
	shortResignationTypes []null.Int,
	underIncome string,
	overIncome string,
	desiredTransferTypes []null.Int,
	desiredHolidayTypes []null.Int,
	desiredCompanyScaleTypes []null.Int,
	experienceIndustries []null.Int,
	experienceOccupations []null.Int,
	socialExperiences []null.Int,
	management string,
	licenses []null.Int,
	languages []null.Int,
	excelSkills []null.Int,
	wordSkills []null.Int,
	powerPointSkills []null.Int,
	anotherPCSkills []null.Int,
	developmentLanguages []null.Int,
	developmentOS []null.Int,
	appearanceTypes []null.Int,
	communicationTypes []null.Int,
	thinkingTypes []null.Int,
) *SearchJobSeeker {
	return &SearchJobSeeker{
		FreeWord:                 freeword,
		AgentStaffID:             agentStaffID,
		PhaseTypes:               phaseTypes,
		GenderTypes:              genderTypes,
		AgeUnder:                 ageUnder,
		AgeOver:                  ageOver,
		DesiredIndustries:        desiredIndustries,
		DesiredOccupations:       desiredOccupations,
		DesiredWorkLocations:     desiredWorkLocations,
		FinalEducationTypes:      finalEducationTypes,
		StudyCategoryTypes:       studyCategoryTypes,
		SchoolLevelTypes:         schoolLevelTypes,
		NationalityTypes:         nationalityTypes,
		JobChangeTypes:           jobChangeTypes,
		ShortResignationTypes:    shortResignationTypes,
		UnderIncome:              underIncome,
		OverIncome:               overIncome,
		Management:               management,
		DesiredTransferTypes:     desiredTransferTypes,
		DesiredHolidayTypes:      desiredHolidayTypes,
		DesiredCompanyScaleTypes: desiredCompanyScaleTypes,
		ExperienceIndustries:     experienceIndustries,
		ExperienceOccupations:    experienceOccupations,
		SocialExperiences:        socialExperiences,
		Licenses:                 licenses,
		Languages:                languages,
		ExcelSkills:              excelSkills,
		WordSkills:               wordSkills,
		PowerPointSkills:         powerPointSkills,
		AnotherPCSkills:          anotherPCSkills,
		DevelopmentLanguages:     developmentLanguages,
		DevelopmentOS:            developmentOS,
		AppearanceTypes:          appearanceTypes,
		CommunicationTypes:       communicationTypes,
		ThinkingTypes:            thinkingTypes,
	}
}

type UpdateJobSeekerLineIDParam struct {
	JobSeekerUUID uuid.UUID `json:"job_seeker_uuid" validate:"required"`
	AgentUUID     uuid.UUID `json:"agent_uuid"  validate:"required"`
	Code          string    `json:"code"  validate:"required"`
}

type ActivityMemoParam struct {
	ActivityMemo string `json:"activity_memo"`
}

type UpdatePhaseParam struct {
	JobSeekerID uint     `db:"job_seeker_id" json:"job_seeker_id" validate:"required"`
	Phase       null.Int `db:"phase" json:"phase"`
}

type UpdateJobSeekerCanViewMatchingJobParam struct {
	JobSeekerID        uint `db:"job_seeker_id" json:"job_seeker_id" validate:"required"`
	CanViewMatchingJob bool `db:"can_view_matching_job" json:"can_view_matching_job"`
}

// 面談フェーズ
type Phase int64

const (
	EntryInterview           Phase = iota // エントリー
	InvitationInterview                   // 面談案内済み（面談調整中）
	ReservationInterview                  // 面談予約完了
	WaitingInterview                      // 面談実施待ち
	PreparingAfterInterview               // 面談実施済み（準備中）　※下書き可能
	OperatingAfterInterview               // 面談実施済み（稼働中）　※本登録のみ
	ReleasedAfterInterview                // 面談実施済み（リリース状態）　※本登録のみ
	OfferedAfterInterview                 // サービス終了/決定者
	ContinuingAfterInterview              // サービス終了/今後継続連絡
	QuitedAfterInterview                  // サービス終了/転職活動終了
)

// 求職者・求人で共通
type SearchType string

const (
	SearchAll      SearchType = "all"
	SearchOwn      SearchType = "own"
	SearchAlliance SearchType = "alliance"
	// SearchHelp     SearchType = "help"
)

// ゲストページ用の求職者情報
type JobSeekerForGuest struct {
	ID         uint `json:"id"`
	Agreement  bool `json:"agreement"`
	LineActive bool `json:"line_active"`
	IsBlocked  bool `json:"is_blocked"`
}

// ゲストページ用の求職者情報
type JobSeekerDesiredForGuest struct {
	ID                   uint       `json:"id"`
	Phase                null.Int   `json:"phase"`
	DesiredAnnualIncome  null.Int   `json:"desired_annual_income"`
	DesiredWorkLocations []null.Int `json:"desired_work_locations"`
	DesiredOccupations   []null.Int `json:"desired_occupations"`
	DesiredIndustries    []null.Int `json:"desired_industries"`
}

type CheckJobSeekerByUUIDAndNameParam struct {
	UUID uuid.UUID `json:"uuid" validate:"required"`
	Name string    `json:"name" validate:"required"`
}

// パスワードリセット
type UpdateJobSeekerPasswordParam struct {
	JobSeekerUUID      uuid.UUID `json:"job_seeker_uuid" validate:"required"`
	Password           string    `json:"password" validate:"required"`
	ResetPasswordToken string    `json:"reset_password_token" validate:"required"`
}

// パスワードリセットメール送信
type SendJobSeekerResetPasswordEmailParam struct {
	JobSeekerUUID uuid.UUID `json:"job_seeker_uuid" validate:"required"`
	Email         string    `json:"email" validate:"required"`
}

// 求職者のお問い合わせ
type SendJobSeekerContactParam struct {
	Email   string `json:"email" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// 求職者のお問い合わせ
type UpdateJobSeekerInterviewDateFromGestPageParam struct {
	InterviewDate time.Time `json:"interview_date" validate:"required"`
	JobSeekerID   uint      `json:"job_seeker_id" validate:"required"`
	StaffName     string    `json:"staff_name" validate:"required"`
}

type JobSeekerRegisterStatus struct {
	ID                     uint      `db:"id" json:"id"`
	UUID                   uuid.UUID `db:"uuid" json:"uuid"`
	IsCompletedRegister    bool      `json:"is_completed_register"`
	IsCompletedPhoneNumber bool      `json:"is_completed_phone_number"`
	IsCompletedDesired     bool      `json:"is_completed_desired"`
}
