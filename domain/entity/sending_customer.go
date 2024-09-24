package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingCustomer struct {
	ID                   uint      `db:"id" json:"id"`
	AgentID              uint      `db:"agent_id" json:"agent_id"`                           // エージェントのID
	LastName             string    `db:"last_name" json:"last_name"`                         // 名前（姓）
	FirstName            string    `db:"first_name" json:"first_name"`                       // 名前（名）
	LastFurigana         string    `db:"last_furigana" json:"last_furigana"`                 // フリガナ（セイ）
	FirstFurigana        string    `db:"first_furigana" json:"first_furigana"`               // フリガナ（メイ）
	PhoneNumber          string    `db:"phone_number" json:"phone_number"`                   // 電話番号
	Email                string    `db:"email" json:"email"`                                 // メールアドレス
	ResumeOriginURL      string    `db:"resume_origin_url" json:"resume_origin_url"`         // 履歴書原本のURL（Word or Excel）
	ResumePDFURL         string    `db:"resume_pdf_url" json:"resume_pdf_url"`               // 履歴書のURL（PDF）
	CVOriginURL          string    `db:"cv_origin_url" json:"cv_origin_url"`                 // 職務経歴書原本のURL（Word or Excel）
	CVPDFURL             string    `db:"cv_pdf_url" json:"cv_pdf_url"`                       // 職務経歴書のURL（PDF）
	InterviewDate        time.Time `db:"interview_date" json:"interview_date"`               // 面談日時
	InterviewInformation string    `db:"interview_information" json:"interview_information"` // 面談情報
	Remarks              string    `db:"remarks" json:"remarks"`                             // 備考
	Gender               null.Int  `db:"gender" json:"gender"`                               // 性別
	Nationality          null.Int  `db:"nationality" json:"nationality"`                     // 国籍
	NationalityRemarks   string    `db:"nationality_remarks" json:"nationality_remarks"`     // 国籍 外国籍を選択→国籍備考（フリーテキスト）を表示
	Birthday             string    `db:"birthday" json:"birthday"`                           // 生年月日
	PostCode             string    `db:"post_code" json:"post_code"`                         // 郵便番号
	Prefecture           null.Int  `db:"prefecture" json:"prefecture"`                       // 都道府県
	Address              string    `db:"address" json:"address"`                             // 住所詳細（市町村 番地 建物名 部屋番号）
	AddressFurigana      string    `db:"address_furigana" json:"address_furigana"`           // 住所詳細（フリガナ）
	SchoolName           string    `db:"school_name" json:"school_name"`                     // 学校名（最終学歴）
	Subject              string    `db:"subject" json:"subject"`                             // 学部・学科・コース（最終学歴）
	EntranceYear         string    `db:"entrance_year" json:"entrance_year"`                 // 入学年月（最終学歴）
	GraduationYear       string    `db:"graduation_year" json:"graduation_year"`             // 卒業年月（最終学歴）
	StateOfEmployment    null.Int  `db:"state_of_employment" json:"state_of_employment"`     // 就業状況
	JobChange            null.Int  `db:"job_change" json:"job_change"`                       // 転職回数
	JobSummary           string    `db:"job_summary" json:"job_summary"`                     // 職務要約
	CompanyName          string    `db:"company_name" json:"company_name"`                   // 会社名（直近の就業先）
	JoiningYear          string    `db:"joining_year" json:"joining_year"`                   // 入社年月（直近の就業先）
	RetireYear           string    `db:"retire_year" json:"retire_year"`                     // 退職年月（直近の就業先）
	FirstStatus          null.Int  `db:"first_status" json:"first_status"`                   // 開始ステータス（入社, 入行, 入局など）（直近の就業先）
	LastStatus           null.Int  `db:"last_status" json:"last_status"`                     // 終了ステータス（一身上の都合により退職, 派遣期間満了につき退職など）（直近の就業先）
	JobDescription       string    `db:"job_description" json:"job_description"`             // 職務内容（直近の就業先）
	HistorySupplement    string    `db:"history_supplement" json:"history_supplement"`       // 経歴補足
	Phase                null.Int  `db:"phase" json:"phase"`                                 // 相談状況（0: 面談実施待ち, 1: 送客応諾, 2: 送客完了, 3: 送客なし/終了）
	SendingAt            time.Time `db:"sending_at" json:"sending_at"`                       // 送客の実行時間
	CreatedAt            time.Time `db:"created_at" json:"created_at"`                       // 作成日時
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`                       // 最終更新日時

	JobSeekerID  uint   `db:"job_seeker_id" json:"job_seeker_id"`
	AgentName    string `db:"agent_name" json:"agent_name"`
	SendingCount uint   `db:"sending_count" json:"sending_count"`
}

func NewSendingCustomer(
	agentID uint,
	lastName string,
	firstName string,
	lastFurigana string,
	firstFurigana string,
	phoneNumber string,
	email string,
	resumeOriginURL string,
	resumePDFURL string,
	cvOriginURL string,
	cvPDFURL string,
	interviewDate time.Time,
	interviewInformation string,
	remarks string,
	gender null.Int,
	nationality null.Int,
	nationalityRemarks string,
	birthday string,
	postCode string,
	prefecture null.Int,
	address string,
	addressFurigana string,
	schoolName string,
	subject string,
	entranceYear string,
	graduationYear string,
	stateOfEmployment null.Int,
	jobChange null.Int,
	jobSummary string,
	companyName string,
	joiningYear string,
	retireYear string,
	firstStatus null.Int,
	lastStatus null.Int,
	jobDescription string,
	historySupplement string,
	phase null.Int,
) *SendingCustomer {
	return &SendingCustomer{
		AgentID:              agentID,
		LastName:             lastName,
		FirstName:            firstName,
		LastFurigana:         lastFurigana,
		FirstFurigana:        firstFurigana,
		PhoneNumber:          phoneNumber,
		Email:                email,
		ResumeOriginURL:      resumeOriginURL,
		ResumePDFURL:         resumePDFURL,
		CVOriginURL:          cvOriginURL,
		CVPDFURL:             cvPDFURL,
		InterviewDate:        interviewDate,
		InterviewInformation: interviewInformation,
		Remarks:              remarks,
		Gender:               gender,
		Nationality:          nationality,
		NationalityRemarks:   nationalityRemarks,
		Birthday:             birthday,
		PostCode:             postCode,
		Prefecture:           prefecture,
		Address:              address,
		AddressFurigana:      addressFurigana,
		SchoolName:           schoolName,
		Subject:              subject,
		EntranceYear:         entranceYear,
		GraduationYear:       graduationYear,
		StateOfEmployment:    stateOfEmployment,
		JobChange:            jobChange,
		JobSummary:           jobSummary,
		CompanyName:          companyName,
		JoiningYear:          joiningYear,
		RetireYear:           retireYear,
		FirstStatus:          firstStatus,
		LastStatus:           lastStatus,
		JobDescription:       jobDescription,
		HistorySupplement:    historySupplement,
		Phase:                phase,
	}
}

// 面談日時を更新
type UpdateSendingCustomerInterviewDateParam struct {
	InterviewDate     time.Time `json:"interview_date" validate:"required"`
	SendingCustomerID uint      `json:"sending_customer_id" validate:"required"`
	Name              string    `json:"name" validate:"required"`
}
