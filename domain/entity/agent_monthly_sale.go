package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type AgentMonthlySale struct {
	ID                                      uint       `json:"id" db:"id"`
	ManagementID                            uint       `json:"management_id" db:"management_id"`
	SalesMonth                              string     `json:"sales_month" db:"sales_month"`
	OrderSalesBudget                        null.Float `json:"order_sales_budget" db:"order_sales_budget"`
	OrderCostBudget                         null.Float `json:"order_cost_budget" db:"order_cost_budget"`
	OrderGrossProfitBudget                  null.Float `json:"order_gross_profit_budget" db:"order_gross_profit_budget"`
	OrderAssumedUnitPrice                   null.Float `json:"order_assumed_unit_price" db:"order_assumed_unit_price"`
	OrderExpectedOfferAcceptance            null.Float `json:"order_expected_offer_acceptance" db:"order_expected_offer_acceptance"`
	ClaimSalesRevenueBudget                 null.Float `json:"claim_sales_revenue_budget" db:"claim_sales_revenue_budget"`
	ClaimCostBudget                         null.Float `json:"claim_cost_budget" db:"claim_cost_budget"`
	ClaimGrossMarginBudget                  null.Float `json:"claim_gross_margin_budget" db:"claim_gross_margin_budget"`
	ClaimAssumedUnitPrice                   null.Float `json:"claim_assumed_unit_price" db:"claim_assumed_unit_price"`
	ClaimExpectedNewEmployeeNumber          null.Float `json:"claim_expected_new_employee_number" db:"claim_expected_new_employee_number"`
	SeekerOfferAcceptanceTarget             null.Float `json:"seeker_offer_acceptance_target" db:"seeker_offer_acceptance_target"`
	SeekerOfferTarget                       null.Float `json:"seeker_offer_target" db:"seeker_offer_target"`
	SeekerFinalSelectionTarget              null.Float `json:"seeker_final_selection_target" db:"seeker_final_selection_target"`
	SeekerSelectionTarget                   null.Float `json:"seeker_selection_target" db:"seeker_selection_target"`
	SeekerRecommendationCompletionTarget    null.Float `json:"seeker_recommendation_completion_target" db:"seeker_recommendation_completion_target"`
	SeekerJobIntroductionTarget             null.Float `json:"seeker_job_introduction_target" db:"seeker_job_introduction_target"`
	SeekerInterviewTarget                   null.Float `json:"seeker_interview_target" db:"seeker_interview_target"`
	ActiveOfferTarget                       null.Float `json:"active_offer_target" db:"active_offer_target"`
	ActiveFinalSelectionTarget              null.Float `json:"active_final_selection_target" db:"active_final_selection_target"`
	ActiveSelectionTarget                   null.Float `json:"active_selection_target" db:"active_selection_target"`
	ActiveRecommendationCompletionTarget    null.Float `json:"active_recommendation_completion_target" db:"active_recommendation_completion_target"`
	ActiveJobIntroductionTarget             null.Float `json:"active_job_introduction_target" db:"active_job_introduction_target"`
	InterviewOfferAcceptanceTarget          null.Float `json:"interview_offer_acceptance_target" db:"interview_offer_acceptance_target"`
	InterviewOfferTarget                    null.Float `json:"interview_offer_target" db:"interview_offer_target"`
	InterviewFinalSelectionTarget           null.Float `json:"interview_final_selection_target" db:"interview_final_selection_target"`
	InterviewSelectionTarget                null.Float `json:"interview_selection_target" db:"interview_selection_target"`
	InterviewRecommendationCompletionTarget null.Float `json:"interview_recommendation_completion_target" db:"interview_recommendation_completion_target"`
	InterviewJobIntroductionTarget          null.Float `json:"interview_job_introduction_target" db:"interview_job_introduction_target"`
	InterviewInterviewTarget                null.Float `json:"interview_interview_target" db:"interview_interview_target"`

	// 計算で取得する項目
	OrderSalesPerformance                        null.Float               `json:"order_sales_performance" db:"order_sales_performance"`                                                 // 受注実績
	OrderSalesResultRate                         null.Float               `json:"order_sales_result_rate" db:"order_sales_result_rate"`                                                 // 予実比
	OrderCostPerformance                         null.Float               `json:"order_cost_performance" db:"order_cost_performance"`                                                   // 原価実績
	OrderGrossProfitPerformance                  null.Float               `json:"order_gross_profit_performance" db:"order_gross_profit_performance"`                                   // 受注粗利実績
	OrderGrossProfitResultRate                   null.Float               `json:"order_gross_profit_result_rate" db:"order_gross_profit_result_rate"`                                   // 予実比
	OrderUnitPricePerformance                    null.Float               `json:"order_unit_price_performance" db:"order_unit_price_performance"`                                       // 単価実績
	OrderOfferAcceptancePerformance              null.Float               `json:"order_offer_acceptance_performance" db:"order_offer_acceptance_performance"`                           // 内定承諾実績
	OrderAccuracyA                               null.Float               `json:"order_accuracy_a" db:"order_accuracy_a"`                                                               // Aヨミ
	OrderAccuracyB                               null.Float               `json:"order_accuracy_b" db:"order_accuracy_b"`                                                               // Bヨミ
	OrderAccuracyC                               null.Float               `json:"order_accuracy_c" db:"order_accuracy_c"`                                                               // Cヨミ
	OrderAccuracyTopic                           null.Float               `json:"order_accuracy_topic" db:"order_accuracy_topic"`                                                       // ネタ
	OrderAccuracyIDList                          []uint                   `json:"order_accuracy_id_list" db:"order_accuracy_id_list"`                                                   // ヨミのIDリスト
	ClaimSalesPerformance                        null.Float               `json:"claim_sales_performance" db:"claim_sales_performance"`                                                 // 請求実績
	ClaimSalesResultRate                         null.Float               `json:"claim_sales_result_rate" db:"claim_sales_result_rate"`                                                 // 予実比
	ClaimCostPerformance                         null.Float               `json:"claim_cost_performance" db:"claim_cost_performance"`                                                   // 原価実績
	ClaimGrossProfitPerformance                  null.Float               `json:"claim_gross_profit_performance" db:"claim_gross_profit_performance"`                                   // 請求粗利実績
	ClaimGrossProfitResultRate                   null.Float               `json:"claim_gross_profit_result_rate" db:"claim_gross_profit_result_rate"`                                   // 予実比
	ClaimUnitPricePerformance                    null.Float               `json:"claim_unit_price_performance" db:"claim_unit_price_performance"`                                       // 単価実績
	ClaimNewEmployeeNumbererformance             null.Float               `json:"claim_new_employee_number_performance" db:"claim_new_employee_number_performance"`                     // 入社人数実績
	ClaimAccuracyA                               null.Float               `json:"claim_accuracy_a" db:"claim_accuracy_a"`                                                               // Aヨミ
	ClaimAccuracyB                               null.Float               `json:"claim_accuracy_b" db:"claim_accuracy_b"`                                                               // Bヨミ
	ClaimAccuracyC                               null.Float               `json:"claim_accuracy_c" db:"claim_accuracy_c"`                                                               // Cヨミ
	ClaimAccuracyTopic                           null.Float               `json:"claim_accuracy_topic" db:"claim_accuracy_topic"`                                                       // ネタ
	ClaimAccuracyIDList                          []uint                   `json:"claim_accuracy_id_list" db:"claim_accuracy_id_list"`                                                   // ヨミのIDリスト
	SeekerOfferAcceptancePerformance             null.Float               `json:"seeker_offer_acceptance_performance" db:"seeker_offer_acceptance_performance"`                         // 内定承諾数
	SeekerOfferAcceptanceDifference              null.Float               `json:"seeker_offer_acceptance_difference" db:"seeker_offer_acceptance_difference"`                           // 内定承諾の差異
	SeekerOfferAcceptanceRate                    null.Float               `json:"seeker_offer_acceptance_rate" db:"seeker_offer_acceptance_rate"`                                       // 承諾率
	SeekerOfferPerformance                       null.Float               `json:"seeker_offer_performance" db:"seeker_offer_performance"`                                               // 内定数
	SeekerOfferDifference                        null.Float               `json:"seeker_offer_difference" db:"seeker_offer_difference"`                                                 // 内定の差異
	SeekerOfferRate                              null.Float               `json:"seeker_offer_rate" db:"seeker_offer_rate"`                                                             // 内定率
	SeekerFinalSelectionPerformance              null.Float               `json:"seeker_final_selection_performance" db:"seeker_final_selection_performance"`                           // 最終選考数
	SeekerFinalSelectionDifference               null.Float               `json:"seeker_final_selection_difference" db:"seeker_final_selection_difference"`                             // 最終選考の差異
	SeekerFinalSelectionRate                     null.Float               `json:"seeker_final_selection_rate" db:"seeker_final_selection_rate"`                                         // 最終選考実施率
	SeekerSelectionPerformance                   null.Float               `json:"seeker_selection_performance" db:"seeker_selection_performance"`                                       // 選考数
	SeekerSelectionDifference                    null.Float               `json:"seeker_selection_difference" db:"seeker_selection_difference"`                                         // 選考の差異
	SeekerSelectionRate                          null.Float               `json:"seeker_selection_rate" db:"seeker_selection_rate"`                                                     // 選考実施率
	SeekerRecommendationCompletionPerformance    null.Float               `json:"seeker_recommendation_completion_performance" db:"seeker_recommendation_completion_performance"`       // 推薦完了数
	SeekerRecommendationCompletionDifference     null.Float               `json:"seeker_recommendation_completion_difference" db:"seeker_recommendation_completion_difference"`         // 推薦完了の差異
	SeekerRecommendationCompletionRate           null.Float               `json:"seeker_recommendation_completion_rate" db:"seeker_recommendation_completion_rate"`                     // 応諾率
	SeekerJobIntroductionPerformance             null.Float               `json:"seeker_job_introduction_performance" db:"seeker_job_introduction_performance"`                         // 求人紹介数
	SeekerJobIntroductionDifference              null.Float               `json:"seeker_job_introduction_difference" db:"seeker_job_introduction_difference"`                           // 求人紹介の差異
	SeekerJobIntroductionRate                    null.Float               `json:"seeker_job_introduction_rate" db:"seeker_job_introduction_rate"`                                       // 紹介稼働率
	SeekerInterviewPerformance                   null.Float               `json:"seeker_interview_performance" db:"seeker_interview_performance"`                                       // 面談数
	SeekerInterviewDifference                    null.Float               `json:"seeker_interview_difference" db:"seeker_interview_difference"`                                         // 面談の差異
	InterviewOfferAcceptancePerformance          null.Float               `json:"interview_offer_acceptance_performance" db:"interview_offer_acceptance_performance"`                   // 内定承諾数
	InterviewOfferAcceptanceDifference           null.Float               `json:"interview_offer_acceptance_difference" db:"interview_offer_acceptance_difference"`                     // 内定承諾の差異
	InterviewOfferAcceptanceRate                 null.Float               `json:"interview_offer_acceptance_rate" db:"interview_offer_acceptance_rate"`                                 // 承諾率
	InterviewOfferPerformance                    null.Float               `json:"interview_offer_performance" db:"interview_offer_performance"`                                         // 内定数
	InterviewOfferDifference                     null.Float               `json:"interview_offer_difference" db:"interview_offer_difference"`                                           // 内定の差異
	InterviewOfferRate                           null.Float               `json:"interview_offer_rate" db:"interview_offer_rate"`                                                       // 内定率
	InterviewFinalSelectionPerformance           null.Float               `json:"interview_final_selection_performance" db:"interview_final_selection_performance"`                     // 最終選考数
	InterviewFinalSelectionDifference            null.Float               `json:"interview_final_selection_difference" db:"interview_final_selection_difference"`                       // 最終選考の差異
	InterviewFinalSelectionRate                  null.Float               `json:"interview_final_selection_rate" db:"interview_final_selection_rate"`                                   // 最終選考実施率
	InterviewSelectionPerformance                null.Float               `json:"interview_selection_performance" db:"interview_selection_performance"`                                 // 選考数
	InterviewSelectionDifference                 null.Float               `json:"interview_selection_difference" db:"interview_selection_difference"`                                   // 選考の差異
	InterviewSelectionRate                       null.Float               `json:"interview_selection_rate" db:"interview_selection_rate"`                                               // 選考実施率
	InterviewRecommendationCompletionPerformance null.Float               `json:"interview_recommendation_completion_performance" db:"interview_recommendation_completion_performance"` // 推薦完了数
	InterviewRecommendationCompletionDifference  null.Float               `json:"interview_recommendation_completion_difference" db:"interview_recommendation_completion_difference"`   // 推薦完了の差異
	InterviewRecommendationCompletionRate        null.Float               `json:"interview_recommendation_completion_rate" db:"interview_recommendation_completion_rate"`               // 応諾率
	InterviewJobIntroductionPerformance          null.Float               `json:"interview_job_introduction_performance" db:"interview_job_introduction_performance"`                   // 求人紹介数
	InterviewJobIntroductionDifference           null.Float               `json:"interview_job_introduction_difference" db:"interview_job_introduction_difference"`                     // 求人紹介の差異
	InterviewJobIntroductionRate                 null.Float               `json:"interview_job_introduction_rate" db:"interview_job_introduction_rate"`                                 // 紹介稼働率
	InterviewInterviewPerformance                null.Float               `json:"interview_interview_performance" db:"interview_interview_performance"`                                 // 面談数
	InterviewInterviewDifference                 null.Float               `json:"interview_interview_difference" db:"interview_interview_difference"`                                   // 面談の差異
	CreatedAt                                    time.Time                `db:"created_at" json:"-"`
	UpdatedAt                                    time.Time                `db:"updated_at" json:"-"`
	AgentStaffMonthlySaleList                    []*AgentStaffMonthlySale `json:"agent_staff_sale_target_management_list"`
}

func NewAgentMonthlySale(
	managementID uint,
	salesMonth string,
	orderSalesBudget null.Float,
	orderCostBudget null.Float,
	orderGrossProfitBudget null.Float,
	orderAssumedUnitPrice null.Float,
	orderExpectedOfferAcceptance null.Float,
	claimSalesRevenueBudget null.Float,
	claimCostBudget null.Float,
	claimGrossMarginBudget null.Float,
	claimAssumedUnitPrice null.Float,
	claimExpectedNewEmployeeNumber null.Float,
	seekerOfferAcceptanceTarget null.Float,
	seekerOfferTarget null.Float,
	seekerFinalSelectionTarget null.Float,
	seekerSelectionTarget null.Float,
	seekerRecommendationCompletionTarget null.Float,
	seekerJobIntroductionTarget null.Float,
	seekerInterviewTarget null.Float,
	activeOfferTarget null.Float,
	activeFinalSelectionTarget null.Float,
	activeSelectionTarget null.Float,
	activeRecommendationCompletionTarget null.Float,
	activeJobIntroductionTarget null.Float,
	interviewOfferAcceptanceTarget null.Float,
	interviewOfferTarget null.Float,
	interviewFinalSelectionTarget null.Float,
	interviewSelectionTarget null.Float,
	interviewRecommendationCompletionTarget null.Float,
	interviewJobIntroductionTarget null.Float,
	interviewInterviewTarget null.Float,
) *AgentMonthlySale {
	return &AgentMonthlySale{
		ManagementID:                            managementID,
		SalesMonth:                              salesMonth,
		OrderSalesBudget:                        orderSalesBudget,
		OrderCostBudget:                         orderCostBudget,
		OrderGrossProfitBudget:                  orderGrossProfitBudget,
		OrderAssumedUnitPrice:                   orderAssumedUnitPrice,
		OrderExpectedOfferAcceptance:            orderExpectedOfferAcceptance,
		ClaimSalesRevenueBudget:                 claimSalesRevenueBudget,
		ClaimCostBudget:                         claimCostBudget,
		ClaimGrossMarginBudget:                  claimGrossMarginBudget,
		ClaimAssumedUnitPrice:                   claimAssumedUnitPrice,
		ClaimExpectedNewEmployeeNumber:          claimExpectedNewEmployeeNumber,
		SeekerOfferAcceptanceTarget:             seekerOfferAcceptanceTarget,
		SeekerOfferTarget:                       seekerOfferTarget,
		SeekerFinalSelectionTarget:              seekerFinalSelectionTarget,
		SeekerSelectionTarget:                   seekerSelectionTarget,
		SeekerRecommendationCompletionTarget:    seekerRecommendationCompletionTarget,
		SeekerJobIntroductionTarget:             seekerJobIntroductionTarget,
		SeekerInterviewTarget:                   seekerInterviewTarget,
		ActiveOfferTarget:                       activeOfferTarget,
		ActiveFinalSelectionTarget:              activeFinalSelectionTarget,
		ActiveSelectionTarget:                   activeSelectionTarget,
		ActiveRecommendationCompletionTarget:    activeRecommendationCompletionTarget,
		ActiveJobIntroductionTarget:             activeJobIntroductionTarget,
		InterviewOfferAcceptanceTarget:          interviewOfferAcceptanceTarget,
		InterviewOfferTarget:                    interviewOfferTarget,
		InterviewFinalSelectionTarget:           interviewFinalSelectionTarget,
		InterviewSelectionTarget:                interviewSelectionTarget,
		InterviewRecommendationCompletionTarget: interviewRecommendationCompletionTarget,
		InterviewJobIntroductionTarget:          interviewJobIntroductionTarget,
		InterviewInterviewTarget:                interviewInterviewTarget,
	}
}

type CreateOrUpdateAgentMonthlyManagementParam struct {
	ID                uint               `json:"id"`
	AgentID           uint               `json:"agent_id" validate:"required"`
	FiscalYear        string             `json:"fiscal_year" validate:"required"`
	IsOpen            bool               `json:"is_open" db:"is_open"`
	AgentMonthlySales []AgentMonthlySale `json:"agent_monthly_sales" validate:"required"`
}
