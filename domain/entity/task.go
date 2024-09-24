package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type Task struct {
	ID                          uint      `db:"id" json:"id"`
	TaskGroupID                 uint      `db:"task_group_id" json:"task_group_id"`
	RAStaffID                   uint      `db:"ra_staff_id" json:"ra_staff_id"`
	CAStaffID                   uint      `db:"ca_staff_id" json:"ca_staff_id"`
	PhaseCategory               null.Int  `db:"phase_category" json:"phase_category"`
	PhaseSubCategory            null.Int  `db:"phase_sub_category" json:"phase_sub_category"`
	StaffType                   null.Int  `db:"staff_type" json:"staff_type"`
	ExecutedStaffID             uint      `db:"executed_staff_id" json:"executed_staff_id"`
	Remarks                     string    `db:"remarks" json:"remarks"`
	DeadlineDay                 string    `db:"deadline_day" json:"deadline_day"`
	DeadlineTime                null.Int  `db:"deadline_time" json:"deadline_time"`
	TalkAboutInInterview        string    `db:"talk_about_in_interview" json:"talk_about_in_interview"`
	ScheduleCollectionCondition string    `db:"schedule_collection_condition" json:"schedule_collection_condition"`
	ExamGuideContent            string    `db:"exam_guide_content" json:"exam_guide_content"`
	IsCheckDoubleSided          bool      `db:"is_check_double_sided" json:"is_check_double_sided"`
	CreatedAt                   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                   time.Time `db:"updated_at" json:"updated_at"`

	// 送付をお勧めする書類
	IsRecommendDocument TaskIsRecommendDocument `db:"is_recommend_document" json:"is_recommend_document"`

	// 担当者
	RAStaffName string `db:"ra_staff_name" json:"ra_staff_name"`
	CAStaffName string `db:"ca_staff_name" json:"ca_staff_name"`
	RAAgentID   uint   `db:"ra_agent_id" json:"ra_agent_id"`
	CAAgentID   uint   `db:"ca_agent_id" json:"ca_agent_id"`
	RAAgentName string `db:"ra_agent_name" json:"ra_agent_name"`
	CAAgentName string `db:"ca_agent_name" json:"ca_agent_name"`

	// 評価点項目
	GoodPoint string `db:"good_point" json:"good_point"`
	NGPoint   string `db:"ng_point" json:"ng_point"`

	// 面談項目
	// InterviewDate   string   `db:"interview_date" json:"interview_date"`
	// InterviewDateID null.Int `db:"interview_date_id" json:"interview_date_id"` // 面談日時ID　※被りチェック用

	// 書類項目
	UnofficialOfferNoticeURL  string `db:"unofficial_offer_notice_url" json:"unofficial_offer_notice_url"`
	ConditionsNoticeURL       string `db:"conditions_notice_url" json:"conditions_notice_url"`
	UnofficialOfferConsentURL string `db:"unofficial_offer_consent_url" json:"unofficial_offer_consent_url"`

	// 求人項目
	JobInformationUUID      uuid.UUID `db:"job_information_uuid" json:"job_information_uuid"`
	JobInformationID        uint      `db:"job_information_id" json:"job_information_id"`
	BillingAddressID        uint      `db:"billing_address_id" json:"billing_address_id"`
	CompanyName             string    `db:"company_name" json:"company_name"`
	Title                   string    `db:"title" json:"title"`
	ExternalJobListingURL   string    `db:"external_job_listing_url" json:"external_job_listing_url"`   // 外部求人の時に格納される
	RequiredDocumentsDetail string    `db:"required_documents_detail" json:"required_documents_detail"` //推薦時に必要な情報・書類の詳細
	HowToRecommend          string    `db:"how_to_recommend" json:"how_to_recommend"`

	// 求職者項目
	JobSeekerUUID uuid.UUID `db:"job_seeker_uuid" json:"job_seeker_uuid"`
	JobSeekerID   uint      `db:"job_seeker_id" json:"job_seeker_id"`
	LastName      string    `db:"last_name" json:"last_name"`
	FirstName     string    `db:"first_name" json:"first_name"`
	LastFurigana  string    `db:"last_furigana" json:"last_furigana"`
	FirstFurigana string    `db:"first_furigana" json:"first_furigana"`

	// タスクグループ
	TaskGroupUUID     uuid.UUID `db:"task_group_uuid" json:"task_group_uuid"`
	JoiningDate       string    `db:"joining_date" json:"joining_date"`
	RALastRequestAt   time.Time `db:"ra_last_request_at" json:"ra_last_request_at"` // RAの最終依頼時間
	RALastWatchedAt   time.Time `db:"ra_last_watched_at" json:"ra_last_watched_at"` // RAの最終閲覧時間
	CALastRequestAt   time.Time `db:"ca_last_request_at" json:"ca_last_request_at"` // CAの最終依頼時間
	CALastWatchedAt   time.Time `db:"ca_last_watched_at" json:"ca_last_watched_at"` // CAの最終閲覧時間
	IsDoubleSided     bool      `db:"is_double_sided" json:"is_double_sided"`
	IsSelfApplication bool      `db:"is_self_application" json:"is_self_application"` // 自己応募フラグ（true: 自己応募, false: エージェント応募）

	IsUnconfirmed bool `db:"is_unconfirmed" json:"is_unconfirmed"` // 未確認フラグ（true: 未確認, false: 確認済み）

	// GetLatestByAgentIDで使用
	ExternalJobInformationTitle string `db:"external_job_information_title" json:"-"`
	ExternalCompanyName         string `db:"external_company_name" json:"-"`

	// 選考フロー
	SelectionFlowPatternID null.Int `db:"selection_flow_pattern_id" json:"selection_flow_pattern_id"`
	SelectionInformationID null.Int `db:"selection_information_id" json:"selection_information_id"`
	FlowPattern            null.Int `db:"flow_pattern" json:"flow_pattern"`
	IsQuestionnairy        bool     `db:"is_questionnairy" json:"is_questionnairy"`

	// アンケート
	SelectionQuestionnaireList []SelectionQuestionnaire `db:"selection_questionnaire_list" json:"selection_questionnaire_list"`

	// 選考の確定日時・選考の候補日時
	SelectionDate JobSeekerSchedule   `json:"selection_date"`
	PossibleDates []JobSeekerSchedule `json:"possible_dates"`
}

func NewTask(
	taskGroupID uint,
	phaseCategory null.Int,
	phaseSubCategory null.Int,
	staffType null.Int,
	executedStaffID uint,
	remarks string,
	deadlineDay string,
	deadlineTime null.Int,
	talkAboutInInterview string,
	scheduleCollectionCondition string,
	examGuideContent string,
	isCheckDoubleSided bool,
) *Task {
	return &Task{
		TaskGroupID:                 taskGroupID,
		PhaseCategory:               phaseCategory,
		PhaseSubCategory:            phaseSubCategory,
		StaffType:                   staffType,
		ExecutedStaffID:             executedStaffID,
		Remarks:                     remarks,
		DeadlineDay:                 deadlineDay,
		DeadlineTime:                deadlineTime,
		TalkAboutInInterview:        talkAboutInInterview,
		ScheduleCollectionCondition: scheduleCollectionCondition,
		ExamGuideContent:            examGuideContent,
		IsCheckDoubleSided:          isCheckDoubleSided,
	}
}

type TaskListByType struct {
	ExecuteRATaskList []*Task `db:"execute_ra_task_list" json:"execute_ra_task_list"`
	ExecuteCATaskList []*Task `db:"execute_ca_task_list" json:"execute_ca_task_list"`
	RequestRATaskList []*Task `db:"request_ra_task_list" json:"request_ra_task_list"`
	RequestCATaskList []*Task `db:"request_ca_task_list" json:"request_ca_task_list"`
}

type RATaskListByType struct {
	ExecuteRATaskList []*Task `db:"execute_ra_task_list" json:"execute_ra_task_list"`
	RequestRATaskList []*Task `db:"request_ra_task_list" json:"request_ra_task_list"`
}

type CATaskListByType struct {
	ExecuteCATaskList []*Task `db:"execute_ca_task_list" json:"execute_ca_task_list"`
	RequestCATaskList []*Task `db:"request_ca_task_list" json:"request_ca_task_list"`
}

type DeleteTaskParam struct {
	ID          uint `json:"id" validate:"required"`
	TaskGroupID uint `json:"task_group_id" validate:"required"`
}

// 次のタスク
type NextTaskParam struct {
	TaskGroupID                 uint                `json:"task_group_id" validate:"required"`
	JobSeekerID                 uint                `json:"job_seeker_id" validate:"required"`
	JobInformationID            uint                `json:"job_information_id" validate:"required"`
	SelectionFlowPatternID      null.Int            `json:"selection_flow_pattern_id"`
	SelectionInformationID      null.Int            `json:"selection_information_id"`
	RAStaffID                   uint                `json:"ra_staff_id"`
	CAStaffID                   uint                `json:"ca_staff_id"`
	PrevPhaseCategory           null.Int            `json:"prev_phase_category" validate:"required"`
	PrevPhaseSubCategory        null.Int            `json:"prev_phase_sub_category" validate:"required"`
	PrevStaffType               null.Int            `json:"prev_staff_type" validate:"required"`
	PhaseCategory               null.Int            `json:"phase_category" validate:"required"`
	PhaseSubCategory            null.Int            `json:"phase_sub_category" validate:"required"`
	StaffType                   null.Int            `json:"staff_type" validate:"required"`
	ExecutedStaffID             uint                `json:"executed_staff_id" validate:"required"`
	Remarks                     string              `json:"remarks"`
	DeadlineDay                 string              `json:"deadline_day" validate:"required"`
	DeadlineTime                null.Int            `json:"deadline_time" validate:"required"`
	GoodPoint                   string              `json:"good_point"`
	NGPoint                     string              `json:"ng_point"`
	Document1URL                string              `json:"document1_url"`
	Document2URL                string              `json:"document2_url"`
	Document3URL                string              `json:"document3_url"`
	Document4URL                string              `json:"document4_url"`
	AcceptJobOfferConsentURL    string              `json:"accept_job_offer_consent_url"`
	SelectionDate               JobSeekerSchedule   `json:"selection_date"`
	PossibleDates               []JobSeekerSchedule `json:"possible_dates"`
	JoiningDate                 string              `json:"joining_date"`
	TalkAboutInInterview        string              `json:"talk_about_in_interview"`
	ScheduleCollectionCondition string              `json:"schedule_collection_condition"`
	ExamGuideContent            string              `json:"exam_guide_content"`
	IsCheckDoubleSided          bool                `db:"is_check_double_sided" json:"is_check_double_sided"`

	// メール
	AgentStaffID uint        `json:"agent_staff_id"` // 重複しないカラム毎の担当者のid
	Tos          []EmailUser `json:"tos"`            // 送信先の情報
	Subject      string      `json:"subject"`        // メールの件名
	Content      string      `json:"content"`        // メールの本文

	// 求職者へのメッセージ（LINE or メール）
	Message TaskMessageParam `json:"message"`

	// 「0: メッセージを送信する」「1: 次の選考をスキップする」「2: リスケジュールが実行時のアクション」「3: 求職者へのメッセージ（LINE or メール）」
	TaskOption null.Int `json:"task_option"`

	// 送付をお勧めする書類
	IsRecommendDocument TaskIsRecommendDocument `db:"is_recommend_document" json:"is_recommend_document"`
}

// 次のタスク
type CreateTaskInBatchProcessingParam struct {
	TaskGroupID            uint        `json:"task_group_id" validate:"required"`
	RAStaffID              uint        `json:"ra_staff_id" validate:"required"`
	CAStaffID              uint        `json:"ca_staff_id" validate:"required"`
	JobSeekerID            uint        `json:"job_seeker_id" validate:"required"`
	JobInformationID       uint        `json:"job_information_id" validate:"required"`
	SelectionFlowPatternID null.Int    `json:"selection_flow_pattern_id"`
	DeadlineDay            string      `json:"deadline_day" validate:"required"`
	DeadlineTime           null.Int    `json:"deadline_time" validate:"required"`
	PrevStaffType          null.Int    `json:"prev_staff_type" validate:"required"`
	TaskList               []TaskParam `json:"task_list"`
}

type TaskParam struct {
	PrevPhaseCategory           null.Int            `json:"prev_phase_category" validate:"required"`
	PrevPhaseSubCategory        null.Int            `json:"prev_phase_sub_category" validate:"required"`
	PhaseCategory               null.Int            `json:"phase_category" validate:"required"`
	PhaseSubCategory            null.Int            `json:"phase_sub_category" validate:"required"`
	StaffType                   null.Int            `json:"staff_type" validate:"required"`
	SelectionInformationID      null.Int            `json:"selection_information_id"`
	ExecutedStaffID             uint                `json:"executed_staff_id" validate:"required"`
	Remarks                     string              `json:"remarks"`
	GoodPoint                   string              `json:"good_point"`
	NGPoint                     string              `json:"ng_point"`
	Document1URL                string              `json:"document1_url"`
	Document2URL                string              `json:"document2_url"`
	Document3URL                string              `json:"document3_url"`
	Document4URL                string              `json:"document4_url"`
	AcceptJobOfferConsentURL    string              `json:"accept_job_offer_consent_url"`
	SelectionDate               JobSeekerSchedule   `json:"selection_date"`
	PossibleDates               []JobSeekerSchedule `json:"possible_dates"`
	JoiningDate                 string              `json:"joining_date"`
	TalkAboutInInterview        string              `json:"talk_about_in_interview"`
	ScheduleCollectionCondition string              `json:"schedule_collection_condition"`
	ExamGuideContent            string              `json:"exam_guide_content"`
	IsCheckDoubleSided          bool                `db:"is_check_double_sided" json:"is_check_double_sided"`

	// メール
	AgentStaffID uint        `json:"agent_staff_id"` // 重複しないカラム毎の担当者のid
	Tos          []EmailUser `json:"tos"`            // 送信先の情報
	Subject      string      `json:"subject"`        // メールの件名
	Content      string      `json:"content"`        // メールの本文

	// 求職者へのメッセージ（LINE or メール）
	Message TaskMessageParam `json:"message"`

	// 「0: メッセージを送信する」「1: 次の選考をスキップする」「2: リスケジュールが実行時のアクション」「3: 求職者へのメッセージ（LINE or メール）」
	TaskOption null.Int `json:"task_option"`

	// 送付をお勧めする書類
	IsRecommendDocument TaskIsRecommendDocument `db:"is_recommend_document" json:"is_recommend_document"`
}

// 同一タスク処理用
type SameTaskParam struct {
	TaskGroupID      uint `json:"task_group_id" validate:"required"`
	JobSeekerID      uint `json:"job_seeker_id" validate:"required"`
	JobInformationID uint `json:"job_information_id" validate:"required"`
	RAAgentID        uint `json:"ra_agent_id" validate:"required"`
	CAAgentID        uint `json:"ca_agent_id" validate:"required"`
	RAStaffID        uint `json:"ra_staff_id" validate:"required"`
	CAStaffID        uint `json:"ca_staff_id" validate:"required"`
}

type TaskMessageParam struct {
	LineActive bool   `json:"line_active"`
	Subject    string `json:"subject"` // メッセージの件名
	Content    string `json:"content"` // メッセージの本文
	Email      string `json:"email"`   // メッセージの本文
}

// 同一タスクをまとめて処理するapiに使用
type NextSameTaskListParam struct {
	Task         NextTaskParam   `json:"task" validate:"required"`
	SameTaskList []SameTaskParam `json:"same_task_list" validate:"required"`
}

// マイページのマッチ求人からエントリー {job_seeker_uuid, job_information_uuid, type}
type CreateEntryTaskFromMatchingJobParam struct {
	JobSeekerUUID      uuid.UUID `json:"job_seeker_uuid" validate:"required"`
	JobInformationUUID uuid.UUID `json:"job_information_uuid" validate:"required"`
	IsEntry            bool      `json:"is_entry"`
}

const (
	OptionSendMessage             int64 = iota // メッセージを送信する
	OptionSkipSelection                        // 次の選考をスキップする
	OptionReschedule                           // リスケジュールが実行時のアクション
	OptionSendMessageForJobSeeker              // 求職者へのメッセージ（LINE or メール）
	OptionRescheduleAndDelete                  // 新リスケ処理（過去の特定のフェーズまでタスクを削除）
)

type StaffType int64

const (
	CA      StaffType = iota // CA
	RA                       // RA
	CA_Boss                  // CA上長
	RA_Boss                  // RA上長
)

type TaskCategory uint

const (
	Entry             TaskCategory = iota // エントリー
	DocumentSelection                     // 書類選考
	FirstSelection                        // 1次選考
	SecondSelection                       // 2次選考
	ThirdSelection                        // 3次選考
	FourthSelection                       // 4次選考
	FifthSelection                        // 5次選考
	FinalSelection                        // 最終選考
	HoldJobOffer                          // 内定保留
	AcceptJobOffer                        // 内定承諾
	IssueInvoice                          // 請求書発行
	Refund                                // 返金
	Decline           = 99                // 辞退
)

// エントリーのサブカテゴリ
type EntrySub int64

const (
	SoundOutMask                      EntrySub = iota // マスクレジュメ打診の依頼（RAタスク）
	InHouseNG                                         // マスクレジュメ社内NG（CAタスク）
	CollectResultOfMask                               // マスクレジュメの結果回収中（RAタスク）
	EnterpriseNG                                      // マスクレジュメ不合格（CAタスク）
	SoundOutJobInformation                            // 求人打診の依頼（CAタスク）
	ConfirmApplicationIntention                       // 応募意思確認中（CAタスク）
	NotAcceptEntry                                    // 応募辞退
	DeclineEntry                                      // エントリー辞退
	HoldEntry                                         // エントリー保留
	RequestShareJobSeeker                             // 求職者シェア依頼（CAタスク）
	WithoutSoundOut                                   // 打診せず終了（RAタスク）
	RequestShareJobInformation                        // 求人シェア依頼（RAタスク）
	ConfirmPossibility                                // 企業に合格可能性を確認中（RAタスク）
	Unlikely                                          // 合格可能性が低いためNG（CAタスク）
	CloseEntryPhaseForUnlikely        = 90            // 合格可能性が低いためNG終了
	CloseEntryPhaseForWithoutSoundOut = 91            // 打診せず終了
	CloseEntryPhaseForDeclineEntry    = 92            // エントリー辞退
	CloseEntryPhaseForRequestDecline  = 98            // 辞退終了
	CloseEntryPhase                   = 99            // エントリーフェーズタスクの終了
)

// 書類選考のサブカテゴリ
type DocumentSelectionSub int64

const (
	RequestRecommendations                DocumentSelectionSub = iota // 推薦依頼
	CollectResultOfDocumentSelection                                  // 結果回収中
	FailingNotification                                               // 不合格通知
	RequestEntry                                                      // エントリー依頼
	PrepareDocument                                                   // 書類準備（CAタスク）
	NotAcceptDocumentSelection                                        // 応募辞退
	DeclineNotificationForDocumentPhase                               // 終了通知（CAタスク）
	RequestDeclineProcessForDocumentPhase                             // 辞退処理依頼（RAタスク）
	CloseDocumentPhaseForDeclineEntry     = 90                        // エントリー辞退
	CloseDocumentPhaseForFailingNotice    = 91                        // 不合格終了
	CloseDocumentPhaseForDrop             = 97                        // 中断終了
	CloseDocumentPhaseForRequestDecline   = 98                        // 辞退終了
	CloseDocumentPhase                    = 99                        // 書類選考フェーズタスクの終了
)

// ○次選考のサブカテゴリ
type SelectionSub int64

const (
	RequestCollectionOfSchedule                    SelectionSub = iota // 候補日回収依頼（CAタスク）
	CollectingSchedule                                                 // 候補日回収中（CAタスク）
	RequestScheduleAdjustmentForSelection                              // 日程調整依頼（RAタスク）
	EnterpriseAdjustmentToScheduleForSelection                         // 日程調整中（企業）（RAタスク）
	RequestConfirmScheduleForSelection                                 // 日程案内依頼（詳細未案内）（CAタスク）
	JobSeekerConfirmsScheduleForSelection                              // 日程確認中（詳細未案内）（CAタスク）
	ScheduleConfirmedForSelection                                      // 日程確定（詳細未案内）（RAタスク）
	ConfirmDayBefore                                                   // 前日確認依頼（CAタスク）
	RequestGuidance                                                    // 案内依頼（CAタスク）
	CollectSelectionThought                                            // 選考所感の回収（CAタスク）
	RequestCollectSelectionThought                                     // 選考結果の回収依頼（RAタスク）
	JobSeekerSupport                                                   // 求職者対応中（CAタスク）
	FailingNoticeOfSelection                                           // 不合格通知（CAタスク）
	RequestDetailedInformation                                         // 詳細案内依頼（CAタスク）
	RequestConfirmScheduleAndSelectionDetail                           // 日程・詳細案内依頼（CAタスク）
	JobSeekerConfirmsScheduleAndSelectionDetail                        // 日程・詳細確認中（CAタスク）
	ScheduleConfirmedAndGuidedDetail                                   // 日程確定（詳細案内済み）（RAタスク）
	DeclineNotification                                                // 終了通知（CAタスク）
	RequestDeclineProcess                                              // 辞退処理依頼（RAタスク）
	SkipSelectionThoughtQuestionnaireText                              // 選考所感の回収（前日確認メール＋所感回収アンケート（文面）を送信をスキップする）（CAタスク）
	SendSelectionThoughtQuestionnaireText                              // 選考所感の回収（前日確認メール＋所感回収アンケート（文面）を送信する）（CAタスク）
	CloseSelectionPhaseForFailingNoticeOfSelection = 90                // 不合格終了
	CloseSelectionPhaseForDrop                     = 97                // 中断終了
	CloseSelectionPhaseForRequestDecline           = 98                // 辞退終了
	CloseSelectionPhase                            = 99                // 選考フェーズタスクの終了
	SelectionSkip                                  = 999               // 選考スキップ
	RescheduleDeleteOfSelection                    = 1000              // リスケ（タスク削除）
)

// 内定保留のサブカテゴリ
type HoldJobOfferSub int64

const (
	RequestJobOfferNotification                        HoldJobOfferSub = iota // 内定通知依頼（CAタスク）
	AcceptJobOfferMind                                                        // 承諾意思の確認（CAタスク）
	RequestCollectionOfScheduleForHold                                        // オファー面談の候補日回収依頼（CAタスク）
	CollectionOfSchedule                                                      // オファー面談の日程回収中（CAタスク）
	RequestScheduleAdjustment                                                 // オファー面談の日程調整依頼（RAタスク）
	AdjustmentToSchedule                                                      // オファー面談の日程調整中（RAタスク）
	RequestConfirmSchedule                                                    // オファー面談の日程確認依頼（CAタスク）
	ConfirmsSchedule                                                          // オファー面談の日程・詳細確認中（CAタスク）
	ScheduleConfirmed                                                         // オファー面談の日程確定（RAタスク）
	RequestDetailedInformationForHold                                         // オファー面談の詳細案内依頼（CAタスク）
	RequestCollectionOfInformation                                            // 情報回収依頼（RAタスク）
	CollectingInformation                                                     // 情報回収中（RAタスク）
	DeclineJobOffer                                                           // 内定辞退（RAタスク）
	CloseHoldJobOfferPhaseForDeclineJobOffer           = 90                   // タスク終了/内定辞退終了
	CloseHoldJobOfferPhaseForRequestDeclineOfSelection = 98                   // タスク終了/辞退終了
	HoldJobOfferPhaseClose                             = 99                   // 内定保留フェーズのタスク終了
	RescheduleDeleteOfHoldJobOffer                     = 1000                 // リスケ（タスク削除）
)

// 内定承諾のサブカテゴリ
type AcceptJobOfferSub int64

const (
	Accept AcceptJobOfferSub = iota // 内定承諾
	FollowJoiningCompany
	RequestConfirmJoiningDay
	ConfirmingJoiningCompany
	RequestJoiningCompanyAdjustment
	RequestGuidanceForAccept
	ConfirmingWithJobSeeker
	RequestSupport
	Supporting
	ConfirmJoiningCompany
	RequestConfirmJoiningCompany
	InvoiceInvoice
	AcceptJobOfferPhaseClose = 99   // 承諾後辞退終了
	Decision                 = 1001 // 決定終了
)

// 請求書発行
type IssueInvoiceSub int64

// 返金
type RefundSub int64

// 辞退
type DeclineSub int64

const (
	RequestDecline            DeclineSub = 100 + iota // 辞退依頼
	RequestHoldBackDecline                            // 辞退引き止め依頼
	ConfirmDeclineIntention                           // 辞退意思確認中
	RequestReScheduleToCA                             // キャンセル・リスケ依頼（CA）
	ConfirmReScheduleToCA                             // キャンセル・リスケ確認中
	CompleteConfirmReSchedule                         // キャンセル・リスケ確認完了
	RequestReScheduleToRA                             // キャンセル・リスケ依頼（RA）
	ConfirmReScheduleToRA                             // キャンセル・リスケ確認中（RA）
	ContinueSelection         = 1000                  // 選考継続（直前の○次選考のタスクに戻る）
	// FailingNoticeOfDecline               = 1001                  // 不合格通知（直前の○次選考のタスクに戻る）
	// RequestCollectionOfScheduleOfDecline = 1002                  // 不合格通知（直前の○次選考のタスクに戻る）
	// RequestGuidanceOfDecline             = 1003                  // 不合格通知（直前の○次選考のタスクに戻る）
	DeclinePhaseClose = 99 // タスク終了
)

type SearchTask struct {
	JobSeekerFreeWord  string
	EnterpriseFreeWord string
	RAStaffID          string
	CAStaffID          string
	TaskPhaseTypes     []null.Int
}

func NewSearchTask(
	jobSeekerFreeWord string,
	enterpriseFreeWord string,
	raStaffID string,
	caStaffID string,
	taskPhaseTypes []null.Int,
) *SearchTask {
	return &SearchTask{
		JobSeekerFreeWord:  jobSeekerFreeWord,
		EnterpriseFreeWord: enterpriseFreeWord,
		RAStaffID:          raStaffID,
		CAStaffID:          caStaffID,
		TaskPhaseTypes:     taskPhaseTypes,
	}
}

// 求人打診
type SoundOutForJobInformationParam struct {
	AgentID                     uint            `json:"agent_id" validate:"required"`
	AgentStaffID                uint            `json:"agent_staff_id" validate:"required"`
	GroupList                   []SoundOutGroup `json:"group_list" validate:"required"`
	Remarks                     string          `db:"remarks" json:"remarks"`
	DeadlineDay                 string          `json:"deadline_day" validate:"required"`
	DeadlineTime                null.Int        `json:"deadline_time" validate:"required"`
	TalkAboutInInterview        string          `json:"talk_about_in_interview"`
	ScheduleCollectionCondition string          `json:"schedule_collection_condition"`
	ExamGuideContent            string          `json:"exam_guide_content"`
}

type SendJobListingMessages struct {
	JobSeekerID uint   `json:"job_seeker_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Subject     string `json:"subject"`
	Content     string `json:"content" validate:"required"`
}

type SoundOutForJobInformationForSendMessageParam struct {
	AgentID                     uint                     `json:"agent_id" validate:"required"`
	AgentStaffID                uint                     `json:"agent_staff_id" validate:"required"`
	GroupList                   []SoundOutGroup          `json:"group_list" validate:"required"`
	Remarks                     string                   `json:"remarks"`
	DeadlineDay                 string                   `json:"deadline_day" validate:"required"`
	DeadlineTime                null.Int                 `json:"deadline_time" validate:"required"`
	TalkAboutInInterview        string                   `json:"talk_about_in_interview"`
	ScheduleCollectionCondition string                   `json:"schedule_collection_condition"`
	ExamGuideContent            string                   `json:"exam_guide_content"`
	Messages                    []SendJobListingMessages `json:"messages" validate:"required"`
}

type JobSeekerTask struct {
	JobSeekerID       uint      `json:"job_seeker_id"`
	CAStaffID         uint      `json:"ca_staff_id"`
	LastName          string    `json:"last_name"`
	FirstName         string    `json:"first_name"`
	LastFurigana      string    `json:"last_furigana"`
	FirstFurigana     string    `json:"first_furigana"`
	TaskCount         TaskCount `json:"task_count"`
	MyTaskList        []Task    `json:"my_task_list"`      // 自分担当
	RequestTaskList   []Task    `json:"request_task_list"` // 依頼中
	IsSelfApplication bool      `json:"is_self_application"`
	HasUnconfirmed    bool      `json:"has_unconfirmed"` // 未確認タスク所持フラグ
}

type TaskCount struct {
	EntryCount     uint `json:"entry_count" db:"entry_count"`
	DocumentCount  uint `json:"document_count" db:"document_count"`
	SelectionCount uint `json:"selection_count" db:"selection_count"`
	HoldCount      uint `json:"hold_count" db:"hold_count"`
	AcceptCount    uint `json:"accept_count" db:"accept_count"`
	UnwatchCount   uint `json:"unwatch_count" db:"unwatch_count"`
	TodayCount     uint `json:"today_count" db:"today_count"`
	TimeOutCount   uint `json:"time_out_count" db:"time_out_count"`
}
