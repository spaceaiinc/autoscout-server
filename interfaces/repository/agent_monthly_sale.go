package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentMonthlySaleRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewAgentMonthlySaleRepositoryImpl(ex interfaces.SQLExecuter) usecase.AgentMonthlySaleRepository {
	return &AgentMonthlySaleRepositoryImpl{
		Name:     "AgentMonthlySaleRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *AgentMonthlySaleRepositoryImpl) Create(agentMonthlySale *entity.AgentMonthlySale) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO agent_monthly_sales (
			management_id,
			sales_month,
			order_sales_budget,
			order_cost_budget,
			order_gross_profit_budget,
			order_assumed_unit_price,
			order_expected_offer_acceptance,
			claim_sales_revenue_budget,
			claim_cost_budget,
			claim_gross_margin_budget,
			claim_assumed_unit_price,
			claim_expected_new_employee_number,
			seeker_offer_acceptance_target,
			seeker_offer_target,
			seeker_final_selection_target,
			seeker_selection_target,
			seeker_recommendation_completion_target,
			seeker_job_introduction_target,
			seeker_interview_target,
			active_offer_target,
			active_final_selection_target,
			active_selection_target,
			active_recommendation_completion_target,
			active_job_introduction_target,
			interview_offer_acceptance_target,           	
			interview_offer_target,                      
			interview_final_selection_target,            	
			interview_selection_target,                  
			interview_recommendation_completion_target,  	
			interview_job_introduction_target,           	
			interview_interview_target,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?
			)`,
		agentMonthlySale.ManagementID,
		agentMonthlySale.SalesMonth,
		agentMonthlySale.OrderSalesBudget,
		agentMonthlySale.OrderCostBudget,
		agentMonthlySale.OrderGrossProfitBudget,
		agentMonthlySale.OrderAssumedUnitPrice,
		agentMonthlySale.OrderExpectedOfferAcceptance,
		agentMonthlySale.ClaimSalesRevenueBudget,
		agentMonthlySale.ClaimCostBudget,
		agentMonthlySale.ClaimGrossMarginBudget,
		agentMonthlySale.ClaimAssumedUnitPrice,
		agentMonthlySale.ClaimExpectedNewEmployeeNumber,
		agentMonthlySale.SeekerOfferAcceptanceTarget,
		agentMonthlySale.SeekerOfferTarget,
		agentMonthlySale.SeekerFinalSelectionTarget,
		agentMonthlySale.SeekerSelectionTarget,
		agentMonthlySale.SeekerRecommendationCompletionTarget,
		agentMonthlySale.SeekerJobIntroductionTarget,
		agentMonthlySale.SeekerInterviewTarget,
		agentMonthlySale.ActiveOfferTarget,
		agentMonthlySale.ActiveFinalSelectionTarget,
		agentMonthlySale.ActiveSelectionTarget,
		agentMonthlySale.ActiveRecommendationCompletionTarget,
		agentMonthlySale.ActiveJobIntroductionTarget,
		agentMonthlySale.InterviewOfferAcceptanceTarget,
		agentMonthlySale.InterviewOfferTarget,
		agentMonthlySale.InterviewFinalSelectionTarget,
		agentMonthlySale.InterviewSelectionTarget,
		agentMonthlySale.InterviewRecommendationCompletionTarget,
		agentMonthlySale.InterviewJobIntroductionTarget,
		agentMonthlySale.InterviewInterviewTarget,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	agentMonthlySale.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
func (repo *AgentMonthlySaleRepositoryImpl) Update(agentMonthlySale *entity.AgentMonthlySale, agentMonthlySaleID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE agent_monthly_sales
		SET
			sales_month = ?,
			order_sales_budget = ?,
			order_cost_budget = ?,
			order_gross_profit_budget = ?,
			order_assumed_unit_price = ?,
			order_expected_offer_acceptance = ?,
			claim_sales_revenue_budget = ?,
			claim_cost_budget = ?,
			claim_gross_margin_budget = ?,
			claim_assumed_unit_price = ?,
			claim_expected_new_employee_number = ?,
			seeker_offer_acceptance_target = ?,
			seeker_offer_target = ?,
			seeker_final_selection_target = ?,
			seeker_selection_target = ?,
			seeker_recommendation_completion_target = ?,
			seeker_job_introduction_target = ?,
			seeker_interview_target = ?,
			active_offer_target = ?,
			active_final_selection_target = ?,
			active_selection_target = ?,
			active_recommendation_completion_target = ?,
			active_job_introduction_target = ?,
			interview_offer_acceptance_target = ?,           	
			interview_offer_target = ?,                      
			interview_final_selection_target = ?,            	
			interview_selection_target = ?,                  
			interview_recommendation_completion_target = ?,  	
			interview_job_introduction_target = ?,           	
			interview_interview_target = ?,
			updated_at = ?
		WHERE
			id = ?
		`,
		agentMonthlySale.SalesMonth,
		agentMonthlySale.OrderSalesBudget,
		agentMonthlySale.OrderCostBudget,
		agentMonthlySale.OrderGrossProfitBudget,
		agentMonthlySale.OrderAssumedUnitPrice,
		agentMonthlySale.OrderExpectedOfferAcceptance,
		agentMonthlySale.ClaimSalesRevenueBudget,
		agentMonthlySale.ClaimCostBudget,
		agentMonthlySale.ClaimGrossMarginBudget,
		agentMonthlySale.ClaimAssumedUnitPrice,
		agentMonthlySale.ClaimExpectedNewEmployeeNumber,
		agentMonthlySale.SeekerOfferAcceptanceTarget,
		agentMonthlySale.SeekerOfferTarget,
		agentMonthlySale.SeekerFinalSelectionTarget,
		agentMonthlySale.SeekerSelectionTarget,
		agentMonthlySale.SeekerRecommendationCompletionTarget,
		agentMonthlySale.SeekerJobIntroductionTarget,
		agentMonthlySale.SeekerInterviewTarget,
		agentMonthlySale.ActiveOfferTarget,
		agentMonthlySale.ActiveFinalSelectionTarget,
		agentMonthlySale.ActiveSelectionTarget,
		agentMonthlySale.ActiveRecommendationCompletionTarget,
		agentMonthlySale.ActiveJobIntroductionTarget,
		agentMonthlySale.InterviewOfferAcceptanceTarget,
		agentMonthlySale.InterviewOfferTarget,
		agentMonthlySale.InterviewFinalSelectionTarget,
		agentMonthlySale.InterviewSelectionTarget,
		agentMonthlySale.InterviewRecommendationCompletionTarget,
		agentMonthlySale.InterviewJobIntroductionTarget,
		agentMonthlySale.InterviewInterviewTarget,
		time.Now().In(time.UTC),
		agentMonthlySaleID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 単数取得 API
//
func (repo *AgentMonthlySaleRepositoryImpl) FindByManagementIDAndMonth(managementID uint, month string) (*entity.AgentMonthlySale, error) {
	var (
		agentMonthlySale entity.AgentMonthlySale
	)
	err := repo.executer.Get(
		repo.Name+".FindByManagementIDAndMonth",
		&agentMonthlySale, `
		SELECT 
			agent_sale.*, 
			DATE_FORMAT(agent_sale.sales_month, '%Y-%m') AS sales_month
		FROM 
			agent_monthly_sales AS agent_sale
		WHERE
			agent_sale.management_id = ?
		AND
			DATE_FORMAT(agent_sale.sales_month, '%Y-%m') = ?
		LIMIT 1
	`,
		managementID, month)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &agentMonthlySale, nil
}

/****************************************************************************************/
/// 複数取得 API
//
func (repo *AgentMonthlySaleRepositoryImpl) GetByManagementID(managementID uint) ([]*entity.AgentMonthlySale, error) {
	var (
		saleTargetList []*entity.AgentMonthlySale
	)

	err := repo.executer.Select(
		repo.Name+".GetByManagementID",
		&saleTargetList, `
		SELECT 
			agent_sale.*, DATE_FORMAT(agent_sale.sales_month, '%Y-%m') AS sales_month
		FROM 
			agent_monthly_sales AS agent_sale
		WHERE
			management_id = ?
		ORDER BY id ASC
		`, managementID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleTargetList, nil
}

func (repo *AgentMonthlySaleRepositoryImpl) GetByManagementIDAndPeriod(managementID uint, startMonth, endMonth string) ([]*entity.AgentMonthlySale, error) {
	var (
		agentMonthlySaleList []*entity.AgentMonthlySale
	)

	err := repo.executer.Select(
		repo.Name+".GetByManagementIDAndPeriod",
		&agentMonthlySaleList, `
			SELECT 
				agent_sale.*, 
				DATE_FORMAT(agent_sale.sales_month, '%Y-%m') AS sales_month
			FROM 
				agent_monthly_sales AS agent_sale
			WHERE
				agent_sale.management_id = ?
			AND 
				DATE_FORMAT(agent_sale.sales_month, '%Y-%m') >= ?  AND
				DATE_FORMAT(agent_sale.sales_month, '%Y-%m') <= ?
			ORDER BY id ASC
		`, managementID, startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentMonthlySaleList, nil
}
