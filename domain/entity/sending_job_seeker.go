package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// 送客求職者
// job_seekersからhide_to_agents, inflow_channel_idを削除して、sending_customer_idを追加
type SendingJobSeeker struct {
	ID                         uint      `db:"id" json:"id"`
	UUID                       uuid.UUID `db:"uuid" json:"uuid"`
	AgentID                    uint      `db:"agent_id" json:"agent_id"`
	AgentStaffID               null.Int  `db:"agent_staff_id" json:"agent_staff_id"`
	SenderAgentID              uint      `db:"sender_agent_id" json:"sender_agent_id"` // 送客元のエージェントID
	SenderAgentName            string    `db:"sender_agent_name" json:"sender_agent_name"` // 送客元のエージェント名
	LineID                     string    `db:"line_id" json:"-"`
	LineActive                 bool      `db:"line_active" json:"line_active"`
	AgentName                  string    `db:"agent_name" json:"agent_name"`
	StaffName                  string    `db:"staff_name" json:"staff_name"`
	StaffEmail                 string    `db:"staff_email" json:"staff_email"`
	StaffPhoneNumber           string    `db:"staff_phone_number" json:"staff_phone_number"`
	UserStatus                 null.Int  `db:"user_status" json:"user_status"`
	LastName                   string    `db:"last_name" json:"last_name"`
	FirstName                  string    `db:"first_name" json:"first_name"`
	LastFurigana               string    `db:"last_furigana" json:"last_furigana"`
	FirstFurigana              string    `db:"first_furigana" json:"first_furigana"`
	Gender                     null.Int  `db:"gender" json:"gender"`
	GenderRemarks              string    `db:"gender_remarks" json:"gender_remarks"`
	Birthday                   string    `db:"birthday" json:"birthday"`
	Spouse                     null.Int  `db:"spouse" json:"spouse"`
	SupportObligation          null.Int  `db:"support_obligation" json:"support_obligation"`
	Dependents                 null.Int  `db:"dependents" json:"dependents"`
	PhoneNumber                string    `db:"phone_number" json:"phone_number"`
	Email                      string    `db:"email" json:"email"`
	EmergencyPhoneNumber       string    `db:"emergency_phone_number" json:"emergency_phone_number"`
	PostCode                   string    `db:"post_code" json:"post_code"`
	Prefecture                 null.Int  `db:"prefecture" json:"prefecture"`
	Address                    string    `db:"address" json:"address"`
	AddressFurigana            string    `db:"address_furigana" json:"address_furigana"`
	StateOfEmployment          null.Int  `db:"state_of_employment" json:"state_of_employment"`
	JobSummary                 string    `db:"job_summary" json:"job_summary"`
	HistorySupplement          string    `db:"history_supplement" json:"history_supplement"`
	ResearchContent            string    `db:"research_content" json:"research_content"`
	JoinCompanyPeriod          null.Int  `db:"join_company_period" json:"join_company_period"`
	JobChange                  null.Int  `db:"job_change" json:"job_change"`
	AnnualIncome               null.Int  `db:"annual_income" json:"annual_income"`
	DesiredAnnualIncome        null.Int  `db:"desired_annual_income" json:"desired_annual_income"`
	Transfer                   null.Int  `db:"transfer" json:"transfer"`
	TransferRequirement        string    `db:"transfer_requirement" json:"transfer_requirement"`
	ShortResignation           null.Int  `db:"short_resignation" json:"short_resignation"`
	ShortResignationRemarks    string    `db:"short_resignation_remarks" json:"short_resignation_remarks"`
	MedicalHistory             null.Int  `db:"medical_history" json:"medical_history"`
	Nationality                null.Int  `db:"nationality" json:"nationality"`
	Appearance                 null.Int  `db:"appearance" json:"appearance"`
	AppearanceDetailOfTruth    string    `db:"appearance_detail_of_truth" json:"appearance_detail_of_truth"`
	AppearanceDetail           string    `db:"appearance_detail" json:"appearance_detail"`
	Communication              null.Int  `db:"communication" json:"communication"`
	CommunicationDetailOfTruth string    `db:"communication_detail_of_truth" json:"communication_detail_of_truth"`
	CommunicationDetail        string    `db:"communication_detail" json:"communication_detail"`
	Thinking                   null.Int  `db:"thinking" json:"thinking"`
	ThinkingDetailOfTruth      string    `db:"thinking_detail_of_truth" json:"thinking_detail_of_truth"`
	ThinkingDetail             string    `db:"thinking_detail" json:"thinking_detail"`
	SecretMemo                 string    `db:"secret_memo" json:"secret_memo"`
	JobHuntingState            null.Int  `db:"job_hunting_state" json:"job_hunting_state"`
	RecommendReason            string    `db:"recommend_reason" json:"recommend_reason"`
	Phase                      null.Int  `db:"phase" json:"phase"`
	Question                   string    `db:"question" json:"question"`

	// 面談調整タスク
	InterviewDate time.Time `db:"interview_date" json:"interview_date"`

	Agreement       bool     `db:"agreement" json:"agreement"`
	StudyCategory   null.Int `db:"study_category" json:"study_category"`       // 専攻学科の大分類(0:　理系 1: 文系)
	WordSkill       null.Int `db:"word_skill" json:"word_skill"`               // Wordのスキル
	ExcelSkill      null.Int `db:"excel_skill" json:"excel_skill"`             // Excelのスキル
	PowerPointSkill null.Int `db:"power_point_skill" json:"power_point_skill"` // PowerPointのスキル
	ActivityMemo    string   `db:"activity_memo" json:"activity_memo"`

	PublicMemo            string `db:"public_memo" json:"public_memo"`                         // 他社エージェント向けメモ
	NationalityRemarks    string `db:"nationality_remarks" json:"nationality_remarks"`         // 国籍 外国籍を選択→国籍備考（フリーテキスト）を表示
	MedicalHistoryRemarks string `db:"medical_history_remarks" json:"medical_history_remarks"` // 既往歴 ありを選択→既往歴備考（フリーテキスト）を表示
	AcceptancePoints      string `db:"acceptance_points" json:"acceptance_points"`             // 応募承諾のポイント

	SendingCustomerID uint   `db:"sending_customer_id" json:"sending_customer_id"` // 送客求職者ID(sending_customersのid) *使用箇所:流入経路はエージェントの名前を入れる。
	SendingAgentName  string `db:"sending_agent_name" json:"sending_agent_name"`   // 送客求職者のエージェント名

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Name      string    `db:"name" json:"name"`
	Furigana  string    `db:"furigana" json:"furigana"`

	// 他テーブル
	StudentHistories      []SendingJobSeekerStudentHistory      `json:"student_histories"`      // 求職者の学歴情報
	WorkHistories         []SendingJobSeekerWorkHistory         `json:"work_histories"`         // 求職者の職歴情報
	Licenses              []SendingJobSeekerLicense             `json:"licenses"`               // 求職者の所持資格
	SelfPromotions        []SendingJobSeekerSelfPromotion       `json:"self_promotions"`        // 求職者の自己PR
	Documents             SendingJobSeekerDocument              `json:"documents"`              // 求職者の資料
	DesiredIndustries     []SendingJobSeekerDesiredIndustry     `json:"desired_industries"`     // 求職者の希望業界
	DesiredOccupations    []SendingJobSeekerDesiredOccupation   `json:"desired_occupations"`    // 求職者の希望職種
	DesiredWorkLocations  []SendingJobSeekerDesiredWorkLocation `json:"desired_work_locations"` // 求職者の希望勤務地
	DevelopmentSkills     []SendingJobSeekerDevelopmentSkill    `json:"development_skills"`     // 求職者の開発スキル
	LanguageSkills        []SendingJobSeekerLanguageSkill       `json:"language_skills"`        // 求職者の言語スキル
	PCTools               []SendingJobSeekerPCTool              `json:"pc_tools"`               // 求職者のPCツール
	DesiredHolidayTypes   []SendingJobSeekerDesiredHolidayType  `json:"desired_holiday_types"`  // 求職者の休日タイプ
	DesiredCompanyScales  []SendingJobSeekerDesiredCompanyScale `json:"desired_company_scales"` // 求職者の希望企業規模
	IsView                SendingJobSeekerIsView                `json:"is_view"`
	IsNotWaitingViewed    bool                                  `db:"is_not_waiting_viewed" json:"is_not_waiting_viewed"`
	IsNotUnregisterViewed bool                                  `db:"is_not_unregister_viewed" json:"is_not_unregister_viewed"`

	// 他社エージェント同士の求人の重複判定用
	IsDuplicate bool `json:"is_duplicate"`

	// 終了理由
	EndReason string   `json:"end_reason" db:"end_reason"` // 終了理由
	EndStatus null.Int `json:"end_status" db:"end_status"` // 終了ステータス
}

func NewSendingJobSeeker(
	agentID uint,
	agentStaffID null.Int,
	lineID string,
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
	appearanceDetailOfTruth string,
	appearanceDetail string,
	communication null.Int,
	communicationDetailOfTruth string,
	communicationDetail string,
	thinking null.Int,
	thinkingDetailOfTruth string,
	thinkingDetail string,
	secretMemo string,
	jobHuntingState null.Int,
	recommendReason string,
	phase null.Int,
	interviewDate time.Time,
	agreement bool,
	studyCategory null.Int,
	wordSkill null.Int,
	excelSkill null.Int,
	powerPointSkill null.Int,
	publicMemo string,
	nationalityRemarks string,
	medicalHistoryRemarks string,
	acceptancePoints string,
	sendingCustomerID uint,
) *SendingJobSeeker {
	return &SendingJobSeeker{
		AgentID:                    agentID,
		AgentStaffID:               agentStaffID,
		LineID:                     lineID,
		UserStatus:                 userStatus,
		LastName:                   lastName,
		FirstName:                  firstName,
		LastFurigana:               lastFurigana,
		FirstFurigana:              firstFurigana,
		Gender:                     gender,
		GenderRemarks:              genderRemarks,
		Birthday:                   birthday,
		Spouse:                     spouse,
		SupportObligation:          supportObligation,
		Dependents:                 dependents,
		PhoneNumber:                phoneNumber,
		Email:                      email,
		EmergencyPhoneNumber:       emergencyPhoneNumber,
		PostCode:                   postCode,
		Prefecture:                 prefecture,
		Address:                    address,
		AddressFurigana:            addressFurigana,
		StateOfEmployment:          stateOfEmployment,
		JobSummary:                 jobSummary,
		HistorySupplement:          historySupplement,
		ResearchContent:            researchContent,
		JoinCompanyPeriod:          joinCompanyPeriod,
		JobChange:                  jobChange,
		AnnualIncome:               annualIncome,
		DesiredAnnualIncome:        desiredAnnualIncome,
		Transfer:                   transfer,
		TransferRequirement:        transferRequirement,
		ShortResignation:           shortResignation,
		ShortResignationRemarks:    shortResignationRemarks,
		MedicalHistory:             medicalHistory,
		Nationality:                nationality,
		Appearance:                 appearance,
		AppearanceDetailOfTruth:    appearanceDetailOfTruth,
		AppearanceDetail:           appearanceDetail,
		Communication:              communication,
		CommunicationDetailOfTruth: communicationDetailOfTruth,
		CommunicationDetail:        communicationDetail,
		Thinking:                   thinking,
		ThinkingDetailOfTruth:      thinkingDetailOfTruth,
		ThinkingDetail:             thinkingDetail,
		SecretMemo:                 secretMemo,
		JobHuntingState:            jobHuntingState,
		RecommendReason:            recommendReason,
		Phase:                      phase,
		InterviewDate:              interviewDate,
		Agreement:                  agreement,
		StudyCategory:              studyCategory,
		WordSkill:                  wordSkill,
		ExcelSkill:                 excelSkill,
		PowerPointSkill:            powerPointSkill,
		PublicMemo:                 publicMemo,
		NationalityRemarks:         nationalityRemarks,
		MedicalHistoryRemarks:      medicalHistoryRemarks,
		AcceptancePoints:           acceptancePoints,
		SendingCustomerID:          sendingCustomerID,
	}
}

// 最初の作成 *作成時は名前と連絡先のみで、それ以外の項目は更新時に作成する
type CreateSendingJobSeekerParam struct {
	AgentID       uint   `db:"agent_id" json:"agent_id" validate:"required"`             // エージェントのID
	LastName      string `db:"last_name" json:"last_name" validate:"required"`           // 名前（姓）
	FirstName     string `db:"first_name" json:"first_name" validate:"required"`         // 名前（名）
	LastFurigana  string `db:"last_furigana" json:"last_furigana" validate:"required"`   // フリガナ（セイ）
	FirstFurigana string `db:"first_furigana" json:"first_furigana" validate:"required"` // フリガナ（メイ）
	PhoneNumber   string `db:"phone_number" json:"phone_number" validate:"required"`     // 電話番号
	Email         string `db:"email" json:"email" validate:"required"`                   // メールアドレス
}

// 最初の更新 *作成時は名前と連絡先のみで、それ以外の項目は更新時に作成する
type FirstUpdateSendingJobSeekerParam struct {
	AgentID              uint      `db:"agent_id" json:"agent_id" validate:"required"`             // エージェントのID
	LastName             string    `db:"last_name" json:"last_name" validate:"required"`           // 名前（姓）
	FirstName            string    `db:"first_name" json:"first_name" validate:"required"`         // 名前（名）
	LastFurigana         string    `db:"last_furigana" json:"last_furigana" validate:"required"`   // フリガナ（セイ）
	FirstFurigana        string    `db:"first_furigana" json:"first_furigana" validate:"required"` // フリガナ（メイ）
	PhoneNumber          string    `db:"phone_number" json:"phone_number" validate:"required"`     // 電話番号
	Email                string    `db:"email" json:"email" validate:"required"`                   // メールアドレス
	ResumeOriginURL      string    `db:"resume_origin_url" json:"resume_origin_url"`               // 履歴書原本のURL（Word or Excel）
	ResumePDFURL         string    `db:"resume_pdf_url" json:"resume_pdf_url"`                     // 履歴書のURL（PDF）
	CVOriginURL          string    `db:"cv_origin_url" json:"cv_origin_url"`                       // 職務経歴書原本のURL（Word or Excel）
	CVPDFURL             string    `db:"cv_pdf_url" json:"cv_pdf_url"`                             // 職務経歴書のURL（PDF）
	InterviewDate        time.Time `db:"interview_date" json:"interview_date"`                     // 面談日時
	InterviewInformation string    `db:"interview_information" json:"interview_information"`       // 面談情報
	Remarks              string    `db:"remarks" json:"remarks"`                                   // 備考
	Gender               null.Int  `db:"gender" json:"gender"`                                     // 性別
	Nationality          null.Int  `db:"nationality" json:"nationality"`                           // 国籍
	NationalityRemarks   string    `db:"nationality_remarks" json:"nationality_remarks"`           // 国籍 外国籍を選択→国籍備考（フリーテキスト）を表示
	Birthday             string    `db:"birthday" json:"birthday"`                                 // 生年月日
	PostCode             string    `db:"post_code" json:"post_code"`                               // 郵便番号
	Prefecture           null.Int  `db:"prefecture" json:"prefecture"`                             // 都道府県
	Address              string    `db:"address" json:"address"`                                   // 住所詳細（市町村 番地 建物名 部屋番号）
	AddressFurigana      string    `db:"address_furigana" json:"address_furigana"`                 // 住所詳細（フリガナ）
	SchoolName           string    `db:"school_name" json:"school_name"`                           // 学校名（最終学歴）
	Subject              string    `db:"subject" json:"subject"`                                   // 学部・学科・コース（最終学歴）
	EntranceYear         string    `db:"entrance_year" json:"entrance_year"`                       // 入学年月（最終学歴）
	GraduationYear       string    `db:"graduation_year" json:"graduation_year"`                   // 卒業年月（最終学歴）
	StateOfEmployment    null.Int  `db:"state_of_employment" json:"state_of_employment"`           // 就業状況
	JobChange            null.Int  `db:"job_change" json:"job_change"`                             // 転職回数
	JobSummary           string    `db:"job_summary" json:"job_summary"`                           // 職務要約
	CompanyName          string    `db:"company_name" json:"company_name"`                         // 会社名（直近の就業先）
	JoiningYear          string    `db:"joining_year" json:"joining_year"`                         // 入社年月（直近の就業先）
	RetireYear           string    `db:"retire_year" json:"retire_year"`                           // 退職年月（直近の就業先）
	FirstStatus          null.Int  `db:"first_status" json:"first_status"`                         // 開始ステータス（入社, 入行, 入局など）（直近の就業先）
	LastStatus           null.Int  `db:"last_status" json:"last_status"`                           // 終了ステータス（一身上の都合により退職, 派遣期間満了につき退職など）（直近の就業先）
	JobDescription       string    `db:"job_description" json:"job_description"`                   // 職務内容（直近の就業先）
	HistorySupplement    string    `db:"history_supplement" json:"history_supplement"`             // 経歴補足
	Phase                null.Int  `db:"phase" json:"phase"`                                       // 相談状況（0: 面談実施待ち, 1: 送客応諾, 2: 送客完了, 3: 送客なし/終了）
	SendingAt            time.Time `db:"sending_at" json:"sending_at"`                             // 送客の実行時間
}

// 初回以降の更新
type UpdateSendingJobSeekerParam struct {
	AgentID                    uint      `db:"agent_id" json:"agent_id" validate:"required"`
	AgentStaffID               null.Int  `db:"agent_staff_id" json:"agent_staff_id"`
	LineID                     string    `db:"line_id" json:"-"`
	UserStatus                 null.Int  `db:"user_status" json:"user_status"`
	LastName                   string    `db:"last_name" json:"last_name"`
	FirstName                  string    `db:"first_name" json:"first_name"`
	LastFurigana               string    `db:"last_furigana" json:"last_furigana"`
	FirstFurigana              string    `db:"first_furigana" json:"first_furigana"`
	Gender                     null.Int  `db:"gender" json:"gender"`
	GenderRemarks              string    `db:"gender_remarks" json:"gender_remarks"`
	Birthday                   string    `db:"birthday" json:"birthday"`
	Spouse                     null.Int  `db:"spouse" json:"spouse"`
	SupportObligation          null.Int  `db:"support_obligation" json:"support_obligation"`
	Dependents                 null.Int  `db:"dependents" json:"dependents"`
	PhoneNumber                string    `db:"phone_number" json:"phone_number"`
	Email                      string    `db:"email" json:"email"`
	EmergencyPhoneNumber       string    `db:"emergency_phone_number" json:"emergency_phone_number"`
	PostCode                   string    `db:"post_code" json:"post_code"`
	Prefecture                 null.Int  `db:"prefecture" json:"prefecture"`
	Address                    string    `db:"address" json:"address"`
	AddressFurigana            string    `db:"address_furigana" json:"address_furigana"`
	StateOfEmployment          null.Int  `db:"state_of_employment" json:"state_of_employment"`
	JobSummary                 string    `db:"job_summary" json:"job_summary"`
	HistorySupplement          string    `db:"history_supplement" json:"history_supplement"`
	ResearchContent            string    `db:"research_content" json:"research_content"`
	JoinCompanyPeriod          null.Int  `db:"join_company_period" json:"join_company_period"`
	JobChange                  null.Int  `db:"job_change" json:"job_change"`
	AnnualIncome               null.Int  `db:"annual_income" json:"annual_income"`
	DesiredAnnualIncome        null.Int  `db:"desired_annual_income" json:"desired_annual_income"`
	Transfer                   null.Int  `db:"transfer" json:"transfer"`
	TransferRequirement        string    `db:"transfer_requirement" json:"transfer_requirement"`
	ShortResignation           null.Int  `db:"short_resignation" json:"short_resignation"`
	ShortResignationRemarks    string    `db:"short_resignation_remarks" json:"short_resignation_remarks"`
	MedicalHistory             null.Int  `db:"medical_history" json:"medical_history"`
	Nationality                null.Int  `db:"nationality" json:"nationality"`
	Appearance                 null.Int  `db:"appearance" json:"appearance"`
	AppearanceDetailOfTruth    string    `db:"appearance_detail_of_truth" json:"appearance_detail_of_truth"`
	AppearanceDetail           string    `db:"appearance_detail" json:"appearance_detail"`
	Communication              null.Int  `db:"communication" json:"communication"`
	CommunicationDetailOfTruth string    `db:"communication_detail_of_truth" json:"communication_detail_of_truth"`
	CommunicationDetail        string    `db:"communication_detail" json:"communication_detail"`
	Thinking                   null.Int  `db:"thinking" json:"thinking"`
	ThinkingDetailOfTruth      string    `db:"thinking_detail_of_truth" json:"thinking_detail_of_truth"`
	ThinkingDetail             string    `db:"thinking_detail" json:"thinking_detail"`
	SecretMemo                 string    `db:"secret_memo" json:"secret_memo"`
	JobHuntingState            null.Int  `db:"job_hunting_state" json:"job_hunting_state"`
	RecommendReason            string    `db:"recommend_reason" json:"recommend_reason"`
	Phase                      null.Int  `db:"phase" json:"phase"`
	InterviewDate              time.Time `db:"interview_date" json:"interview_date"`
	Agreement                  bool      `db:"agreement" json:"agreement"`
	StudyCategory              null.Int  `db:"study_category" json:"study_category"`       // 専攻学科の大分類(0:　理系 1: 文系)
	WordSkill                  null.Int  `db:"word_skill" json:"word_skill"`               // Wordのスキル
	ExcelSkill                 null.Int  `db:"excel_skill" json:"excel_skill"`             // Excelのスキル
	PowerPointSkill            null.Int  `db:"power_point_skill" json:"power_point_skill"` // PowerPointのスキル
	UUID                       uuid.UUID `db:"uuid" json:"uuid"`
	InflowChannelID            null.Int  `db:"inflow_channel_id" json:"inflow_channel_id"`             // agent_inflow_channel_optionsテーブルのID
	PublicMemo                 string    `db:"public_memo" json:"public_memo"`                         // 他社エージェント向けメモ
	NationalityRemarks         string    `db:"nationality_remarks" json:"nationality_remarks"`         // 国籍 外国籍を選択→国籍備考（フリーテキスト）を表示
	MedicalHistoryRemarks      string    `db:"medical_history_remarks" json:"medical_history_remarks"` // 既往歴 ありを選択→既往歴備考（フリーテキスト）を表示
	AcceptancePoints           string    `db:"acceptance_points" json:"acceptance_points"`             // 応募承諾のポイント
	SendingCustomerID          uint      `db:"sending_customer_id" json:"sending_customer_id"`         // 送客求職者ID(sending_customersのid) *使用箇所:流入経路はエージェントの名前を入れる。

	// 他テーブル
	StudentHistories     []SendingJobSeekerStudentHistory      `json:"student_histories"`      // 求職者の学歴情報
	WorkHistories        []SendingJobSeekerWorkHistory         `json:"work_histories"`         // 求職者の職歴情報
	Licenses             []SendingJobSeekerLicense             `json:"licenses"`               // 求職者の所持資格
	SelfPromotions       []SendingJobSeekerSelfPromotion       `json:"self_promotions"`        // 求職者の自己PR
	DesiredIndustries    []SendingJobSeekerDesiredIndustry     `json:"desired_industries"`     // 求職者の希望業界
	DesiredOccupations   []SendingJobSeekerDesiredOccupation   `json:"desired_occupations"`    // 求職者の希望職種
	DesiredWorkLocations []SendingJobSeekerDesiredWorkLocation `json:"desired_work_locations"` // 求職者の希望勤務地
	DevelopmentSkills    []SendingJobSeekerDevelopmentSkill    `json:"development_skills"`     // 求職者の開発スキル
	LanguageSkills       []SendingJobSeekerLanguageSkill       `json:"language_skills"`        // 求職者の言語スキル
	PCTools              []SendingJobSeekerPCTool              `json:"pc_tools"`               // 求職者のPCツール
	DesiredHolidayTypes  []SendingJobSeekerDesiredHolidayType  `json:"desired_holiday_types"`  // 求職者の休日タイプ
	DesiredCompanyScales []SendingJobSeekerDesiredCompanyScale `json:"desired_company_scales"` // 求職者の希望企業規模
}

type DeleteSendingJobSeekerParam struct {
	ID uint `json:"id" validate:"required"`
}

type SearchSendingJobSeeker struct {
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

func NewSearchSendingJobSeeker(
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
) *SearchSendingJobSeeker {
	return &SearchSendingJobSeeker{
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

// 面談フェーズ
type SendingJobSeekerPhase int64

const (
	// 日程未登録
	UnregisterSchedule SendingJobSeekerPhase = iota // エントリー

	// 詳細未登録
	UnregisterDetail

	// 面談実施待ち
	WaitingForInterview

	// 送客応諾
	AcceptSending

	// 送客完了
	CompleteSending

	// 送客なし/終了
	CloseSending
)

// 求職者の面談前アンケート登録
// 1. 求職者の同意項目アップデート
// 2. アンケート情報の登録（業界、職種、勤務地、質問要望）
// 3. 求職者情報の更新（業界、職種、勤務地）
// 4. 求職者情報の更新（ファイル（履歴書（原本）、商務経歴書（原本）））
type CreateSendingInitialQuestionnaireParam struct {
	SendingJobSeekerID uint   `json:"sending_job_seeker_id" validate:"required"`
	Question           string `db:"question" json:"question"`

	// 関連テーブル
	DesiredIndustries    []SendingJobSeekerDesiredIndustry     `json:"desired_industries"`     // 求職者の希望業界
	DesiredOccupations   []SendingJobSeekerDesiredOccupation   `json:"desired_occupations"`    // 求職者の希望職種
	DesiredWorkLocations []SendingJobSeekerDesiredWorkLocation `json:"desired_work_locations"` // 求職者の希望勤務地

	// 求職者テーブル
	Agreement bool                     `db:"agreement" json:"agreement"` // 求職者の同意
	Documents SendingJobSeekerDocument `json:"documents"`                // 求職者の資料
}
