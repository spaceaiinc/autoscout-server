package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type SelectionQuestionnaire struct {
	ID                              uint                              `db:"id" json:"id"`
	UUID                            uuid.UUID                         `db:"uuid" json:"uuid"`
	JobSeekerID                     uint                              `db:"job_seeker_id" json:"job_seeker_id"`
	JobInformationID                uint                              `db:"job_information_id" json:"job_information_id"`
	JobInformationUUID              uuid.UUID                         `db:"job_information_uuid" json:"job_information_uuid"`
	CompanyName                     string                            `db:"company_name" json:"company_name"`
	Title                           string                            `db:"title" json:"title"`
	SelectionInformationID          uint                              `db:"selection_information_id" json:"selection_information_id"`
	SelectionType                   null.Int                          `json:"selection_type" db:"selection_type"`
	MyRanking                       null.Int                          `db:"my_ranking" json:"my_ranking"`
	MyRankingReason                 string                            `db:"my_ranking_reason" json:"my_ranking_reason"`
	ConcernPoint                    string                            `db:"concern_point" json:"concern_point"`
	ContinueSelection               null.Int                          `db:"continue_selection" json:"continue_selection"`
	MyRankingDetail                 string                            `db:"my_ranking_detail" json:"my_ranking_detail"`
	SelectionQuestion               string                            `db:"selection_question" json:"selection_question"`
	Remarks                         string                            `db:"remarks" json:"remarks"`
	IsAnswer                        bool                              `db:"is_answer" json:"is_answer"`
	IsSelfIntroduction              bool                              `db:"is_self_introduction" json:"is_self_introduction"`
	IsSelfPR                        bool                              `db:"is_self_pr" json:"is_self_pr"`
	IsRetireReason                  bool                              `db:"is_retire_reason" json:"is_retire_reason"`
	IsJobChangeAxis                 bool                              `db:"is_job_change_axis" json:"is_job_change_axis"`
	IsApplyingReason                bool                              `db:"is_applying_reason" json:"is_applying_reason"`
	IsCareerVision                  bool                              `db:"is_career_vision" json:"is_career_vision"`
	IsReverseQuestion               bool                              `db:"is_reverse_question" json:"is_reverse_question"`
	IntentionToJobOffer             null.Int                          `db:"intention_to_job_offer" json:"intention_to_job_offer"`
	IntentionDetail                 string                            `db:"intention_detail" json:"intention_detail"`
	CreatedAt                       time.Time                         `db:"created_at" json:"created_at"`
	UpdatedAt                       time.Time                         `db:"updated_at" json:"updated_at"`
	SelectionQuestionnaireMyRanking []SelectionQuestionnaireMyRanking `db:"selection_questionnaire_my_ranking" json:"selection_questionnaire_my_ranking"`
}

func NewSelectionQuestionnaire(
	jobSeekerID uint,
	jobInformationID uint,
	selectionInformationID uint,
	myRanking null.Int,
	myRankingReason string,
	concernPoint string,
	continueSelection null.Int,
	myRankingDetail string,
	selectionQuestion string,
	remarks string,
	isAnswer bool,
	isSelfIntroduction bool,
	isSelfPR bool,
	isRetireReason bool,
	isJobChangeAxis bool,
	isApplyingReason bool,
	isCareerVision bool,
	isReverseQuestion bool,
	intentionToJobOffer null.Int,
	intentionDetail string,
) *SelectionQuestionnaire {
	return &SelectionQuestionnaire{
		JobSeekerID:            jobSeekerID,
		JobInformationID:       jobInformationID,
		SelectionInformationID: selectionInformationID,
		MyRanking:              myRanking,
		MyRankingReason:        myRankingReason,
		ConcernPoint:           concernPoint,
		ContinueSelection:      continueSelection,
		MyRankingDetail:        myRankingDetail,
		SelectionQuestion:      selectionQuestion,
		Remarks:                remarks,
		IsAnswer:               isAnswer,
		IsSelfIntroduction:     isSelfIntroduction,
		IsSelfPR:               isSelfPR,
		IsRetireReason:         isRetireReason,
		IsJobChangeAxis:        isJobChangeAxis,
		IsApplyingReason:       isApplyingReason,
		IsCareerVision:         isCareerVision,
		IsReverseQuestion:      isReverseQuestion,
		IntentionToJobOffer:    intentionToJobOffer,
		IntentionDetail:        intentionDetail,
	}
}

type CreateOrUpdateSelectionQuestionnaireParam struct {
	JobSeekerID            uint     `json:"job_seeker_id" validate:"required"`
	JobInformationID       uint     `json:"job_information_id" validate:"required"`
	SelectionInformationID uint     `json:"selection_information_id" validate:"required"`
	SelectionType          null.Int `json:"selection_type"`
	MyRanking              null.Int `json:"my_ranking" validate:"required"`
	MyRankingReason        string   `json:"my_ranking_reason"`
	ConcernPoint           string   `json:"concern_point"`
	ContinueSelection      null.Int `json:"continue_selection" validate:"required"`
	MyRankingDetail        string   `json:"my_ranking_detail"`
	SelectionQuestion      string   `json:"selection_question"`
	Remarks                string   `json:"remarks"`
	IsAnswer               bool     `json:"is_answer" validate:"required"`
	IsSelfIntroduction     bool     `json:"is_self_introduction"`
	IsSelfPR               bool     `json:"is_self_pr"`
	IsRetireReason         bool     `json:"is_retire_reason"`
	IsJobChangeAxis        bool     `json:"is_job_change_axis"`
	IsApplyingReason       bool     `json:"is_applying_reason"`
	IsCareerVision         bool     `json:"is_career_vision"`
	IsReverseQuestion      bool     `json:"is_reverse_question"`
	IntentionToJobOffer    null.Int `json:"intention_to_job_offer" validate:"required"`
	IntentionDetail        string   `json:"intention_detail"`

	SelectionQuestionnaireMyRanking []SelectionQuestionnaireMyRanking `json:"selection_questionnaire_my_ranking"`
}

type DeleteSelectionQuestionnaireParam struct {
	ID uint `json:"id" validate:"required"`
}
