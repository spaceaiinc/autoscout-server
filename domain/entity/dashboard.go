package entity

import "gopkg.in/guregu/null.v4"

type Dashboard struct {
	ManagementID                                 uint         `json:"management_id" db:"management_id"`
	FiscalYear                                   string       `json:"fiscal_year" db:"fiscal_year"`
	SalesPerformance                             float64      `json:"sales_performance" db:"sales_performance"`
	SalesBudget                                  float64      `json:"sales_budget" db:"sales_budget"`
	SalesResultRate                              float64      `json:"sales_result_rate" db:"sales_result_rate"`
	GrossProfitPerformance                       float64      `json:"gross_profit_performance" db:"gross_profit_performance"`
	GrossProfitBudget                            float64      `json:"gross_profit_budget" db:"gross_profit_budget"`
	GrossProfitResultRate                        float64      `json:"gross_profit_result_rate" db:"gross_profit_result_rate"`
	InterviewOfferAcceptanceTarget               float64      `json:"interview_offer_acceptance_target" db:"interview_offer_acceptance_target"`
	InterviewOfferAcceptancePerformance          float64      `json:"interview_offer_acceptance_performance" db:"interview_offer_acceptance_performance"`
	InterviewOfferAcceptanceResultRate           float64      `json:"interview_offer_acceptance_result_rate" db:"interview_offer_acceptance_result_rate"`
	InterviewOfferTarget                         float64      `json:"interview_offer_target" db:"interview_offer_target"`
	InterviewOfferPerformance                    float64      `json:"interview_offer_performance" db:"interview_offer_performance"`
	InterviewOfferResultRate                     float64      `json:"interview_offer_result_rate" db:"interview_offer_result_rate"`
	InterviewFinalSelectionTarget                float64      `json:"interview_final_selection_target" db:"interview_final_selection_target"`
	InterviewFinalSelectionPerformance           float64      `json:"interview_final_selection_performance" db:"interview_final_selection_performance"`
	InterviewFinalSelectionResultRate            float64      `json:"interview_final_selection_result_rate" db:"interview_final_selection_result_rate"`
	InterviewSelectionTarget                     float64      `json:"interview_selection_target" db:"interview_selection_target"`
	InterviewSelectionPerformance                float64      `json:"interview_selection_performance" db:"interview_selection_performance"`
	InterviewSelectionResultRate                 float64      `json:"interview_selection_result_rate" db:"interview_selection_result_rate"`
	InterviewRecommendationCompletionTarget      float64      `json:"interview_recommendation_completion_target" db:"interview_recommendation_completion_target"`
	InterviewRecommendationCompletionPerformance float64      `json:"interview_recommendation_completion_performance" db:"interview_recommendation_completion_performance"`
	InterviewRecommendationCompletionResultRate  float64      `json:"interview_recommendation_completion_result_rate" db:"interview_recommendation_completion_result_rate"`
	InterviewJobIntroductionTarget               float64      `json:"interview_job_introduction_target" db:"interview_job_introduction_target"`
	InterviewJobIntroductionPerformance          float64      `json:"interview_job_introduction_performance" db:"interview_job_introduction_performance"`
	InterviewJobIntroductionResultRate           float64      `json:"interview_job_introduction_result_rate" db:"interview_job_introduction_result_rate"`
	InterviewInterviewTarget                     float64      `json:"interview_interview_target" db:"interview_interview_target"`
	InterviewInterviewPerformance                float64      `json:"interview_interview_performance" db:"interview_interview_performance"`
	InterviewInterviewResultRate                 float64      `json:"interview_interview_result_rate" db:"interview_interview_result_rate"`
	AccuracyAccept                               float64      `json:"accuracy_accept" db:"accuracy_accept"`
	AccuracyA                                    float64      `json:"accuracy_a" db:"accuracy_a"`
	AccuracyB                                    float64      `json:"accuracy_b" db:"accuracy_b"`
	AccuracyC                                    float64      `json:"accuracy_c" db:"accuracy_c"`
	AccuracyTopic                                float64      `json:"accuracy_topic" db:"accuracy_topic"`
	AccuracyList                                 []*Sale      `db:"accuracy_list" json:"accuracy_list"`
	ReleaseList                                  []*JobSeeker `db:"release_list" json:"release_list"`
}

func NewDashboard(
	managementID uint,
	fiscalYear string,
	salesPerformance float64,
	salesBudget float64,
	salesResultRate float64,
	grossProfitPerformance float64,
	grossProfitBudget float64,
	grossProfitResultRate float64,
	interviewOfferAcceptanceTarget float64,
	interviewOfferAcceptancePerformance float64,
	interviewOfferAcceptanceResultRate float64,
	interviewOfferTarget float64,
	interviewOfferPerformance float64,
	interviewOfferResultRate float64,
	interviewFinalSelectionTarget float64,
	interviewFinalSelectionPerformance float64,
	interviewFinalSelectionResultRate float64,
	interviewSelectionTarget float64,
	interviewSelectionPerformance float64,
	interviewSelectionResultRate float64,
	interviewRecommendationCompletionTarget float64,
	interviewRecommendationCompletionPerformance float64,
	interviewRecommendationCompletionResultRate float64,
	interviewJobIntroductionTarget float64,
	interviewJobIntroductionPerformance float64,
	interviewJobIntroductionResultRate float64,
	interviewInterviewTarget float64,
	interviewInterviewPerformance float64,
	interviewInterviewResultRate float64,
	accuracyAccept float64,
	accuracyA float64,
	accuracyB float64,
	accuracyC float64,
	accuracyTopic float64,
) *Dashboard {
	return &Dashboard{
		ManagementID:                                 managementID,
		FiscalYear:                                   fiscalYear,
		SalesPerformance:                             salesPerformance,
		SalesBudget:                                  salesBudget,
		SalesResultRate:                              salesResultRate,
		GrossProfitPerformance:                       grossProfitPerformance,
		GrossProfitBudget:                            grossProfitBudget,
		GrossProfitResultRate:                        grossProfitResultRate,
		InterviewOfferAcceptanceTarget:               interviewOfferAcceptanceTarget,
		InterviewOfferAcceptancePerformance:          interviewOfferAcceptancePerformance,
		InterviewOfferAcceptanceResultRate:           interviewOfferAcceptanceResultRate,
		InterviewOfferTarget:                         interviewOfferTarget,
		InterviewOfferPerformance:                    interviewOfferPerformance,
		InterviewOfferResultRate:                     interviewOfferResultRate,
		InterviewFinalSelectionTarget:                interviewFinalSelectionTarget,
		InterviewFinalSelectionPerformance:           interviewFinalSelectionPerformance,
		InterviewFinalSelectionResultRate:            interviewFinalSelectionResultRate,
		InterviewSelectionTarget:                     interviewSelectionTarget,
		InterviewSelectionPerformance:                interviewSelectionPerformance,
		InterviewSelectionResultRate:                 interviewSelectionResultRate,
		InterviewRecommendationCompletionTarget:      interviewRecommendationCompletionTarget,
		InterviewRecommendationCompletionPerformance: interviewRecommendationCompletionPerformance,
		InterviewRecommendationCompletionResultRate:  interviewRecommendationCompletionResultRate,
		InterviewJobIntroductionTarget:               interviewJobIntroductionTarget,
		InterviewJobIntroductionPerformance:          interviewJobIntroductionPerformance,
		InterviewJobIntroductionResultRate:           interviewJobIntroductionResultRate,
		InterviewInterviewTarget:                     interviewInterviewTarget,
		InterviewInterviewPerformance:                interviewInterviewPerformance,
		InterviewInterviewResultRate:                 interviewInterviewResultRate,
		AccuracyAccept:                               accuracyAccept,
		AccuracyA:                                    accuracyA,
		AccuracyB:                                    accuracyB,
		AccuracyC:                                    accuracyC,
		AccuracyTopic:                                accuracyTopic,
	}
}

type SearchDashboard struct {
	AgentID       uint
	AgentStaffID  uint
	ManagementID  uint
	Period        uint
	Target        uint
	SearchRange   uint
	AccuracyTypes []null.Int
}

func NewSearchDashboard(
	agentID uint,
	agentStaffID uint,
	managementID uint,
	period uint,
	target uint,
	searchRange uint,
	accuracyTypes []null.Int,
) *SearchDashboard {
	return &SearchDashboard{
		AgentID:       agentID,
		AgentStaffID:  agentStaffID,
		ManagementID:  managementID,
		Period:        period,
		Target:        target,
		SearchRange:   searchRange,
		AccuracyTypes: accuracyTypes,
	}
}
