package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentStaffMonthlySaleRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewAgentStaffMonthlySaleRepositoryImpl(ex interfaces.SQLExecuter) usecase.AgentStaffMonthlySaleRepository {
	return &AgentStaffMonthlySaleRepositoryImpl{
		Name:     "AgentStaffMonthlySaleRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *AgentStaffMonthlySaleRepositoryImpl) Create(agentStaffMonthlySale *entity.AgentStaffMonthlySale) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO agent_staff_monthly_sales (
			staff_management_id,
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
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?
			)`,
		agentStaffMonthlySale.StaffManagementID,
		agentStaffMonthlySale.SalesMonth,
		agentStaffMonthlySale.OrderSalesBudget,
		agentStaffMonthlySale.OrderCostBudget,
		agentStaffMonthlySale.OrderGrossProfitBudget,
		agentStaffMonthlySale.OrderAssumedUnitPrice,
		agentStaffMonthlySale.OrderExpectedOfferAcceptance,
		agentStaffMonthlySale.ClaimSalesRevenueBudget,
		agentStaffMonthlySale.ClaimCostBudget,
		agentStaffMonthlySale.ClaimGrossMarginBudget,
		agentStaffMonthlySale.ClaimAssumedUnitPrice,
		agentStaffMonthlySale.ClaimExpectedNewEmployeeNumber,
		agentStaffMonthlySale.SeekerOfferAcceptanceTarget,
		agentStaffMonthlySale.SeekerOfferTarget,
		agentStaffMonthlySale.SeekerFinalSelectionTarget,
		agentStaffMonthlySale.SeekerSelectionTarget,
		agentStaffMonthlySale.SeekerRecommendationCompletionTarget,
		agentStaffMonthlySale.SeekerJobIntroductionTarget,
		agentStaffMonthlySale.SeekerInterviewTarget,
		agentStaffMonthlySale.ActiveOfferTarget,
		agentStaffMonthlySale.ActiveFinalSelectionTarget,
		agentStaffMonthlySale.ActiveSelectionTarget,
		agentStaffMonthlySale.ActiveRecommendationCompletionTarget,
		agentStaffMonthlySale.ActiveJobIntroductionTarget,
		agentStaffMonthlySale.InterviewOfferAcceptanceTarget,
		agentStaffMonthlySale.InterviewOfferTarget,
		agentStaffMonthlySale.InterviewFinalSelectionTarget,
		agentStaffMonthlySale.InterviewSelectionTarget,
		agentStaffMonthlySale.InterviewRecommendationCompletionTarget,
		agentStaffMonthlySale.InterviewJobIntroductionTarget,
		agentStaffMonthlySale.InterviewInterviewTarget,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	agentStaffMonthlySale.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
func (repo *AgentStaffMonthlySaleRepositoryImpl) Update(agentStaffMonthlySale *entity.AgentStaffMonthlySale, id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE agent_staff_monthly_sales
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
		agentStaffMonthlySale.SalesMonth,
		agentStaffMonthlySale.OrderSalesBudget,
		agentStaffMonthlySale.OrderCostBudget,
		agentStaffMonthlySale.OrderGrossProfitBudget,
		agentStaffMonthlySale.OrderAssumedUnitPrice,
		agentStaffMonthlySale.OrderExpectedOfferAcceptance,
		agentStaffMonthlySale.ClaimSalesRevenueBudget,
		agentStaffMonthlySale.ClaimCostBudget,
		agentStaffMonthlySale.ClaimGrossMarginBudget,
		agentStaffMonthlySale.ClaimAssumedUnitPrice,
		agentStaffMonthlySale.ClaimExpectedNewEmployeeNumber,
		agentStaffMonthlySale.SeekerOfferAcceptanceTarget,
		agentStaffMonthlySale.SeekerOfferTarget,
		agentStaffMonthlySale.SeekerFinalSelectionTarget,
		agentStaffMonthlySale.SeekerSelectionTarget,
		agentStaffMonthlySale.SeekerRecommendationCompletionTarget,
		agentStaffMonthlySale.SeekerJobIntroductionTarget,
		agentStaffMonthlySale.SeekerInterviewTarget,
		agentStaffMonthlySale.ActiveOfferTarget,
		agentStaffMonthlySale.ActiveFinalSelectionTarget,
		agentStaffMonthlySale.ActiveSelectionTarget,
		agentStaffMonthlySale.ActiveRecommendationCompletionTarget,
		agentStaffMonthlySale.ActiveJobIntroductionTarget,
		agentStaffMonthlySale.InterviewOfferAcceptanceTarget,
		agentStaffMonthlySale.InterviewOfferTarget,
		agentStaffMonthlySale.InterviewFinalSelectionTarget,
		agentStaffMonthlySale.InterviewSelectionTarget,
		agentStaffMonthlySale.InterviewRecommendationCompletionTarget,
		agentStaffMonthlySale.InterviewJobIntroductionTarget,
		agentStaffMonthlySale.InterviewInterviewTarget,
		time.Now().In(time.UTC),
		id,
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
func (repo *AgentStaffMonthlySaleRepositoryImpl) FindByStaffManagementIDAndMonth(staffManagementID uint, thisMonth string) (*entity.AgentStaffMonthlySale, error) {
	var (
		agentStaffMonthlySale entity.AgentStaffMonthlySale
	)
	err := repo.executer.Get(
		repo.Name+".FindByStaffManagementIDAndMonth",
		&agentStaffMonthlySale, `
		SELECT 
			staff_sale.*, 
			DATE_FORMAT(staff_sale.sales_month, '%Y-%m') AS sales_month
		FROM 
			agent_staff_monthly_sales AS staff_sale
		WHERE
			staff_sale.staff_management_id = ?
		AND
			DATE_FORMAT(staff_sale.sales_month, '%Y-%m') = ?
		LIMIT 1
	`,
		staffManagementID, thisMonth)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &agentStaffMonthlySale, nil
}

/****************************************************************************************/
/// 複数取得 API
//
func (repo *AgentStaffMonthlySaleRepositoryImpl) GetByStaffManagementID(staffManagementID uint) ([]*entity.AgentStaffMonthlySale, error) {
	var (
		list []*entity.AgentStaffMonthlySale
	)
	err := repo.executer.Select(
		repo.Name+".GetByStaffManagementID",
		&list, `
		SELECT 
			*, DATE_FORMAT(sales_month, '%Y-%m') AS sales_month
		FROM 
			agent_staff_monthly_sales
		WHERE
			staff_management_id = ?
	`, staffManagementID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}

func (repo *AgentStaffMonthlySaleRepositoryImpl) GetByAgentSaleManagementID(agentSaleManagementID uint) ([]*entity.AgentStaffMonthlySale, error) {
	var (
		list []*entity.AgentStaffMonthlySale
	)
	err := repo.executer.Select(
		repo.Name+".GetByStaffManagementID",
		&list, `
		SELECT 
			staff_monthly.*, DATE_FORMAT(staff_monthly.sales_month, '%Y-%m') AS sales_month
		FROM 
			agent_staff_monthly_sales AS staff_monthly
		INNER JOIN
			agent_staff_sale_managements AS staff_sale_management
		ON
			staff_monthly.staff_management_id = staff_sale_management.id
		WHERE
			staff_sale_management.management_id = ?
	`, agentSaleManagementID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}

func (repo *AgentStaffMonthlySaleRepositoryImpl) GetByStaffManagementIDAndPeriod(staffManagementID uint, startMonth, endMonth string) ([]*entity.AgentStaffMonthlySale, error) {
	var (
		list []*entity.AgentStaffMonthlySale
	)
	err := repo.executer.Select(
		repo.Name+".GetByStaffManagementIDAndPeriod",
		&list, `
		SELECT 
			staff_sale.*, DATE_FORMAT(staff_sale.sales_month, '%Y-%m') AS sales_month
		FROM 
			agent_staff_monthly_sales AS staff_sale
		WHERE
			staff_sale.staff_management_id = ?
		AND 
			DATE_FORMAT(staff_sale.sales_month, '%Y-%m') >= ?  AND
			DATE_FORMAT(staff_sale.sales_month, '%Y-%m') <= ?
		ORDER BY id ASC
	`,
		staffManagementID,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}
