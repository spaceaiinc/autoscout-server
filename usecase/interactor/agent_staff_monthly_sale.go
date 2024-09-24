package interactor

import (
	"fmt"
	"math"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type AgentStaffMonthlySaleInteractor interface {
	// 汎用系 API
	CreateAgentStaffMonthlySale(input CreateAgentStaffMonthlySaleInput) (CreateAgentStaffMonthlySaleOutput, error)
	UpdateAgentStaffMonthlySale(input UpdateAgentStaffMonthlySaleInput) (UpdateAgentStaffMonthlySaleOutput, error)
	GetStaffMonthlySaleList(input GetStaffMonthlySaleListInput) (GetStaffMonthlySaleListOutput, error)
	GetStaffSaleManagementAndAgentMonthlyByID(input GetStaffSaleManagementAndAgentMonthlyByIDInput) (GetStaffSaleManagementAndAgentMonthlyByIDOutput, error)
}

type AgentStaffMonthlySaleInteractorImpl struct {
	firebase                           usecase.Firebase
	sendgrid                           config.Sendgrid
	agentStaffMonthlySaleRepository    usecase.AgentStaffMonthlySaleRepository
	agentStaffSaleManagementRepository usecase.AgentStaffSaleManagementRepository
	agentSaleManagementRepository      usecase.AgentSaleManagementRepository
	saleRepository                     usecase.SaleRepository
	taskRepository                     usecase.TaskRepository
	interviewTaskRepository            usecase.InterviewTaskRepository
	interviewTaskGroupRepository       usecase.InterviewTaskGroupRepository
}

// AgentStaffMonthlySaleInteractorImpl is an implementation of AgentStaffMonthlySaleInteractor
func NewAgentStaffMonthlySaleInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	asmsR usecase.AgentStaffMonthlySaleRepository,
	assmR usecase.AgentStaffSaleManagementRepository,
	asmR usecase.AgentSaleManagementRepository,
	sR usecase.SaleRepository,
	tR usecase.TaskRepository,
	itR usecase.InterviewTaskRepository,
	itgR usecase.InterviewTaskGroupRepository,
) AgentStaffMonthlySaleInteractor {
	return &AgentStaffMonthlySaleInteractorImpl{
		firebase:                           fb,
		sendgrid:                           sg,
		agentStaffMonthlySaleRepository:    asmsR,
		agentStaffSaleManagementRepository: assmR,
		agentSaleManagementRepository:      asmR,
		saleRepository:                     sR,
		taskRepository:                     tR,
		interviewTaskRepository:            itR,
		interviewTaskGroupRepository:       itgR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
//求職者の作成
type CreateAgentStaffMonthlySaleInput struct {
	CreateParam entity.CreateOrUpdateStaffMonthlyManagementParam
}

type CreateAgentStaffMonthlySaleOutput struct {
	OK bool
}

func (i *AgentStaffMonthlySaleInteractorImpl) CreateAgentStaffMonthlySale(input CreateAgentStaffMonthlySaleInput) (CreateAgentStaffMonthlySaleOutput, error) {
	var (
		output CreateAgentStaffMonthlySaleOutput
		err    error
	)

	staffSaleManagement := entity.NewAgentStaffSaleManagement(
		input.CreateParam.ManagementID,
		input.CreateParam.AgentStaffID,
	)
	err = i.agentStaffSaleManagementRepository.Create(staffSaleManagement)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, param := range input.CreateParam.StaffMonthlySales {
		agentStaffMonthlySale := entity.NewAgentStaffMonthlySale(
			staffSaleManagement.ID,
			param.SalesMonth,
			param.OrderSalesBudget,
			param.OrderCostBudget,
			param.OrderGrossProfitBudget,
			param.OrderAssumedUnitPrice,
			param.OrderExpectedOfferAcceptance,
			param.ClaimSalesRevenueBudget,
			param.ClaimCostBudget,
			param.ClaimGrossMarginBudget,
			param.ClaimAssumedUnitPrice,
			param.ClaimExpectedNewEmployeeNumber,
			param.SeekerOfferAcceptanceTarget,
			param.SeekerOfferTarget,
			param.SeekerFinalSelectionTarget,
			param.SeekerSelectionTarget,
			param.SeekerRecommendationCompletionTarget,
			param.SeekerJobIntroductionTarget,
			param.SeekerInterviewTarget,
			param.ActiveOfferTarget,
			param.ActiveFinalSelectionTarget,
			param.ActiveSelectionTarget,
			param.ActiveRecommendationCompletionTarget,
			param.ActiveJobIntroductionTarget,
			param.InterviewOfferAcceptanceTarget,
			param.InterviewOfferTarget,
			param.InterviewFinalSelectionTarget,
			param.InterviewSelectionTarget,
			param.InterviewRecommendationCompletionTarget,
			param.InterviewJobIntroductionTarget,
			param.InterviewInterviewTarget,
		)

		err = i.agentStaffMonthlySaleRepository.Create(agentStaffMonthlySale)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.OK = true
	return output, nil
}

// 個人売上の更新
type UpdateAgentStaffMonthlySaleInput struct {
	UpdateParam entity.CreateOrUpdateStaffMonthlyManagementParam
}

type UpdateAgentStaffMonthlySaleOutput struct {
	OK bool
}

func (i *AgentStaffMonthlySaleInteractorImpl) UpdateAgentStaffMonthlySale(input UpdateAgentStaffMonthlySaleInput) (UpdateAgentStaffMonthlySaleOutput, error) {
	var (
		output UpdateAgentStaffMonthlySaleOutput
		err    error
	)

	for _, monthlySale := range input.UpdateParam.StaffMonthlySales {
		err = i.agentStaffMonthlySaleRepository.Update(&monthlySale, monthlySale.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.OK = true
	return output, nil
}

// 求人企業の更新
type GetStaffMonthlySaleListInput struct {
	AgentStaffID uint
	ManagementID uint
}

type GetStaffMonthlySaleListOutput struct {
	StaffSaleManagement *entity.AgentStaffSaleManagement
}

func (i *AgentStaffMonthlySaleInteractorImpl) GetStaffMonthlySaleList(input GetStaffMonthlySaleListInput) (GetStaffMonthlySaleListOutput, error) {
	var (
		output      GetStaffMonthlySaleListOutput
		err         error
		monthlyList []entity.AgentStaffMonthlySale

		// // 当月ベース
		// checkJobSeekerForOffer           = make(map[uint]bool) // 内定の求職者重複チェック
		// checkJobSeekerForFinalSelection  = make(map[uint]bool) // 最終選考の求職者重複チェック
		// checkJobSeekerForSelection       = make(map[uint]bool) // 選考の求職者重複チェック
		// checkJobSeekerForRecommendation  = make(map[uint]bool) // 推薦の求職者重複チェック
		// checkJobSeekerForJobIntroduction = make(map[uint]bool) // 求人紹介の求職者重複チェック
		// checkJobSeekerForInterview       = make(map[uint]bool) // 面談済みの求職者重複チェック

		// // 面談実施月ベース
		// checkInterviewForOfferAcceptance = make(map[uint]bool) // 内定承諾の求職者重複チェック
		// checkInterviewForOffer           = make(map[uint]bool) // 内定の求職者重複チェック
		// checkInterviewForFinalSelection  = make(map[uint]bool) // 最終選考の求職者重複チェック
		// checkInterviewForSelection       = make(map[uint]bool) // 選考の求職者重複チェック
		// checkInterviewForRecommendation  = make(map[uint]bool) // 推薦の求職者重複チェック
		// checkInterviewForJobIntroduction = make(map[uint]bool) // 求人紹介の求職者重複チェック
	)

	staffSaleManagement, err := i.agentStaffSaleManagementRepository.FindByManagementIDAndStaffID(input.ManagementID, input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	layout := "2006-01" // 入力文字列のフォーマット
	t, err := time.Parse(layout, staffSaleManagement.FiscalYear)
	if err != nil {
		fmt.Println("Error:", err)
	}
	previousMonth := t.AddDate(0, -11, 0)
	startMonth := previousMonth.Format("2006-01")

	agentStaffMonthlySaleList, err := i.agentStaffMonthlySaleRepository.GetByStaffManagementID(staffSaleManagement.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	saleList, err := i.saleRepository.GetByStaffIDForMonthly(input.AgentStaffID, startMonth, staffSaleManagement.FiscalYear)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// // 指定期間のタスクを取得
	// taskList, err := i.taskRepository.GetByStaffIDForMonthly(input.AgentStaffID, startMonth, staffSaleManagement.FiscalYear)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// interviewTaskList, err := i.interviewTaskRepository.GetListByStaffIDForMonthly(input.AgentStaffID, startMonth, staffSaleManagement.FiscalYear)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	for _, monthlySale := range agentStaffMonthlySaleList {
		var (
			orderSalesPerformance            float64 // 受注実績
			orderCostPerformance             float64 // 原価実績
			orderGrossProfitPerformance      float64 // 受注粗利実績
			orderOfferAcceptancePerformance  float64 // 内定承諾実績
			orderAccuracyA                   float64 // Aヨミ
			orderAccuracyB                   float64 // Bヨミ
			orderAccuracyC                   float64 // Cヨミ
			orderAccuracyTopic               float64 // ネタ
			orderAccuracyIDList              []uint  // ヨミリスト
			claimSalesPerformance            float64 // 請求実績
			claimCostPerformance             float64 // 原価実績
			claimGrossProfitPerformance      float64 // 請求粗利実績
			claimNewEmployeeNumbererformance float64 // 入社人数実績
			claimAccuracyA                   float64 // Aヨミ
			claimAccuracyB                   float64 // Bヨミ
			claimAccuracyC                   float64 // Cヨミ
			claimAccuracyTopic               float64 // ネタ
			claimAccuracyIDList              []uint  // ヨミリスト

			seekerOfferAcceptancePerformance             float64 // 内定承諾数
			seekerOfferPerformance                       float64 // 内定数
			seekerFinalSelectionPerformance              float64 // 最終選考数
			seekerSelectionPerformance                   float64 // 選考数
			seekerRecommendationCompletionPerformance    float64 // 推薦完了数
			seekerJobIntroductionPerformance             float64 // 求人紹介数
			interviewOfferAcceptancePerformance          float64 // 内定承諾数
			interviewOfferPerformance                    float64 // 内定数
			interviewFinalSelectionPerformance           float64 // 最終選考数
			interviewSelectionPerformance                float64 // 選考数
			interviewRecommendationCompletionPerformance float64 // 推薦完了数
			interviewJobIntroductionPerformance          float64 // 求人紹介数
			interviewPerformance                         float64 // 面談数

			// seekerInterviewPerformance                   float64 // 面談数
			// interviewInterviewPerformance                float64 // 面談数
		)

		for _, sale := range saleList {
			// 受注関連の計算
			if monthlySale.SalesMonth == sale.ContractSignedMonth {
				var (
					caSale    float64 = 0 // CAの売上
					raSale    float64 = 0 // RAの売上
					staffSale float64 = 0 //スタッフの売上

					caCost    float64 = 0 // CAの原価
					raCost    float64 = 0 // RAの原価
					staffCost float64 = 0 // 担当者の原価

					caGrossProfit    float64 = 0 // CAの粗利
					raGrossProfit    float64 = 0 // RAの粗利
					staffGrossProfit float64 = 0 // 担当者の粗利
				)

				// CA売上比率に応じた売上を計算
				if input.AgentStaffID == sale.CAStaffID && sale.CaSalesRatio.Float64 > 0 {
					caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
					caCost = sale.Cost.Float64 * sale.CaSalesRatio.Float64 / 100
					caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
				}

				// RA売上比率に応じた売上を計算
				if input.AgentStaffID == sale.RAStaffID && sale.RaSalesRatio.Float64 > 0 {
					raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
					raCost = sale.Cost.Float64 * sale.RaSalesRatio.Float64 / 100
					raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
				}

				// エージェントの売上
				staffSale = caSale + raSale
				staffCost = caCost + raCost
				staffGrossProfit = caGrossProfit + raGrossProfit

				// 内定承諾の計算
				if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
					orderSalesPerformance = orderSalesPerformance + staffSale                    // 受注売上実績の計算
					orderCostPerformance = orderCostPerformance + staffCost                      // 原価実績の計算
					orderGrossProfitPerformance = orderGrossProfitPerformance + staffGrossProfit // 受注粗利実績の計算

					// 自分がCA担当分のみをカウント
					if input.AgentStaffID == sale.CAStaffID {
						orderOfferAcceptancePerformance = orderOfferAcceptancePerformance + 1 // 内定承諾実績
						seekerOfferAcceptancePerformance = seekerOfferAcceptancePerformance + 1
					}
				}

				// Aヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
					orderAccuracyA = orderAccuracyA + staffSale // Aヨミの請求金額合計
				}

				// Bヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
					orderAccuracyB = orderAccuracyB + staffSale // Aヨミの請求金額合計
				}

				// Cヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
					orderAccuracyC = orderAccuracyC + staffSale // Aヨミの請求金額合計
				}

				// ネタの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
					orderAccuracyTopic = orderAccuracyTopic + staffSale // Aヨミの請求金額合計
				}

				// ヨミのIDリストを作成
				orderAccuracyIDList = append(orderAccuracyIDList, sale.ID)
			}

			// 請求関連の計算
			if monthlySale.SalesMonth == sale.BillingMonth {
				var (
					caSale    float64 = 0 // CAの売上
					raSale    float64 = 0 // RAの売上
					staffSale float64 = 0 //スタッフの売上

					caCost    float64 = 0 // CAの原価
					raCost    float64 = 0 // RAの原価
					staffCost float64 = 0 // 担当者の原価

					caGrossProfit    float64 = 0 // CAの粗利
					raGrossProfit    float64 = 0 // RAの粗利
					staffGrossProfit float64 = 0 // 担当者の粗利
				)

				// CA売上比率に応じた売上を計算
				if input.AgentStaffID == sale.CAStaffID && sale.CaSalesRatio.Float64 > 0 {
					caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
					caCost = sale.Cost.Float64 * sale.CaSalesRatio.Float64 / 100
					caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
				}

				// RA売上比率に応じた売上を計算
				if input.AgentStaffID == sale.RAStaffID && sale.RaSalesRatio.Float64 > 0 {
					raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
					raCost = sale.Cost.Float64 * sale.RaSalesRatio.Float64 / 100
					raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
				}

				// エージェントの売上
				staffSale = caSale + raSale
				staffCost = caCost + raCost
				staffGrossProfit = caGrossProfit + raGrossProfit

				// 内定承諾の計算
				if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
					claimSalesPerformance = claimSalesPerformance + staffSale                    // 請求売上実績の計算
					claimCostPerformance = claimCostPerformance + staffCost                      // 原価実績の計算
					claimGrossProfitPerformance = claimGrossProfitPerformance + staffGrossProfit // 請求粗利実績の計算

					// 自社がCA担当分のみをカウント
					if input.AgentStaffID == sale.CAStaffID {
						claimNewEmployeeNumbererformance = claimNewEmployeeNumbererformance + 1 // 入社人数実績
					}
				}

				// Aヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
					claimAccuracyA = claimAccuracyA + staffSale // Aヨミの請求金額合計
				}

				// Bヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
					claimAccuracyB = claimAccuracyB + staffSale // Aヨミの請求金額合計
				}

				// Cヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
					claimAccuracyC = claimAccuracyC + staffSale // Aヨミの請求金額合計
				}

				// ネタの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
					claimAccuracyTopic = claimAccuracyTopic + staffSale // Aヨミの請求金額合計
				}

				// ヨミのIDリストを作成
				claimAccuracyIDList = append(claimAccuracyIDList, sale.ID)
			}
		}

		/********* 求職者数（当月カウントベース） **********/

		// 内定承諾数(seekerOfferAcceptancePerformance)
		// 受注売上項目の内定承諾実績と同じになる

		// 内定数
		seekerOfferPerformance, err = i.taskRepository.GetSeekerOfferPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 最終選考数
		seekerFinalSelectionPerformance, err = i.taskRepository.GetSeekerFinalSelectionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 選考数
		seekerSelectionPerformance, err = i.taskRepository.GetSeekerSelectionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 推薦完了数
		seekerRecommendationCompletionPerformance, err = i.taskRepository.GetSeekerRecommendationCompletionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求人紹介数
		seekerJobIntroductionPerformance, err = i.taskRepository.GetSeekerJobIntroductionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		/********* 求職者数（面談実施月カウントベース） **********/

		// 内定承諾数
		interviewOfferAcceptancePerformance, err = i.interviewTaskGroupRepository.GetInterviewOfferAcceptancePerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 内定数
		interviewOfferPerformance, err = i.taskRepository.GetInterviewOfferPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 最終選考数
		interviewFinalSelectionPerformance, err = i.taskRepository.GetInterviewFinalSelectionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 選考数
		interviewSelectionPerformance, err = i.taskRepository.GetInterviewSelectionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 推薦完了数
		interviewRecommendationCompletionPerformance, err = i.taskRepository.GetInterviewRecommendationCompletionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求人紹介数
		interviewJobIntroductionPerformance, err = i.taskRepository.GetInterviewJobIntroductionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 面談数（当月カウントベースと面談実施月ベースで統一）
		interviewPerformance, err = i.interviewTaskGroupRepository.GetInterviewPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		/************************ 受注関連 ************************/

		// 受注売上実績
		monthlySale.OrderSalesPerformance = null.NewFloat(orderSalesPerformance, true)

		// 予実比(受注売上)（0で除算をしないために0チェック）
		if monthlySale.OrderSalesBudget.Float64 > 0 {
			result := math.Floor((orderSalesPerformance / monthlySale.OrderSalesBudget.Float64) * 100)
			monthlySale.OrderSalesResultRate = null.NewFloat(result, true)
		} else {
			monthlySale.OrderSalesResultRate = null.NewFloat(0, true)
		}

		// 原価実績
		monthlySale.OrderCostPerformance = null.NewFloat(orderCostPerformance, true)

		// 受注粗利実績
		monthlySale.OrderGrossProfitPerformance = null.NewFloat(orderGrossProfitPerformance, true)

		// 予実比(受注粗利)（0で除算をしないために0チェック）
		if monthlySale.OrderGrossProfitBudget.Float64 > 0 {
			result := math.Floor((orderGrossProfitPerformance / monthlySale.OrderGrossProfitBudget.Float64) * 100)
			monthlySale.OrderGrossProfitResultRate = null.NewFloat(result, true)
		} else {
			monthlySale.OrderGrossProfitResultRate = null.NewFloat(0, true)
		}

		// 単価実績 （0で除算をしないために0チェック）
		if orderOfferAcceptancePerformance > 0 {
			result := math.Floor(orderSalesPerformance / orderOfferAcceptancePerformance)
			monthlySale.OrderUnitPricePerformance = null.NewFloat(result, true)
		} else {
			monthlySale.OrderUnitPricePerformance = null.NewFloat(0, true)
		}

		// 内定承諾実績
		monthlySale.OrderOfferAcceptancePerformance = null.NewFloat(orderOfferAcceptancePerformance, true)

		// ヨミ関連
		monthlySale.OrderAccuracyA = null.NewFloat(orderAccuracyA, true)
		monthlySale.OrderAccuracyB = null.NewFloat(orderAccuracyB, true)
		monthlySale.OrderAccuracyC = null.NewFloat(orderAccuracyC, true)
		monthlySale.OrderAccuracyTopic = null.NewFloat(orderAccuracyTopic, true)
		monthlySale.OrderAccuracyIDList = orderAccuracyIDList

		/************************ 請求関連 ************************/

		// 請求売上実績
		monthlySale.ClaimSalesPerformance = null.NewFloat(claimSalesPerformance, true)

		// 予実比(請求売上)（0で除算をしないために0チェック）
		if monthlySale.ClaimSalesRevenueBudget.Float64 > 0 {
			result := math.Floor((claimSalesPerformance / monthlySale.ClaimSalesRevenueBudget.Float64) * 100)
			monthlySale.ClaimSalesResultRate = null.NewFloat(result, true)
		} else {
			monthlySale.ClaimSalesResultRate = null.NewFloat(0, true)
		}

		// 原価実績
		monthlySale.ClaimCostPerformance = null.NewFloat(claimCostPerformance, true)

		// 受注粗利実績
		monthlySale.ClaimGrossProfitPerformance = null.NewFloat(claimGrossProfitPerformance, true)

		// 予実比(請求粗利)（0で除算をしないために0チェック）
		if monthlySale.ClaimGrossMarginBudget.Float64 > 0 {
			result := math.Floor((claimGrossProfitPerformance / monthlySale.ClaimGrossMarginBudget.Float64) * 100)
			monthlySale.ClaimGrossProfitResultRate = null.NewFloat(result, true)
		} else {
			monthlySale.ClaimGrossProfitResultRate = null.NewFloat(0, true)
		}

		// 単価実績 （0で除算をしないために0チェック）
		if claimNewEmployeeNumbererformance > 0 {
			result := math.Floor(claimSalesPerformance / claimNewEmployeeNumbererformance)
			monthlySale.ClaimUnitPricePerformance = null.NewFloat(result, true)
		} else {
			monthlySale.ClaimUnitPricePerformance = null.NewFloat(0, true)
		}

		// 入社人数実績
		monthlySale.ClaimNewEmployeeNumbererformance = null.NewFloat(claimNewEmployeeNumbererformance, true)

		// ヨミ関連
		monthlySale.ClaimAccuracyA = null.NewFloat(claimAccuracyA, true)
		monthlySale.ClaimAccuracyB = null.NewFloat(claimAccuracyB, true)
		monthlySale.ClaimAccuracyC = null.NewFloat(claimAccuracyC, true)
		monthlySale.ClaimAccuracyTopic = null.NewFloat(claimAccuracyTopic, true)
		monthlySale.ClaimAccuracyIDList = claimAccuracyIDList

		/************************ 求職者数ベース（当月カウントベース） ************************/

		// 内定承諾(実績)
		monthlySale.SeekerOfferAcceptancePerformance = null.NewFloat(seekerOfferAcceptancePerformance, true)
		// 内定承諾(差異)
		monthlySale.SeekerOfferAcceptanceDifference = null.NewFloat((seekerOfferAcceptancePerformance - monthlySale.SeekerOfferAcceptanceTarget.Float64), true)
		// 承諾率
		if seekerOfferPerformance > 0 {
			result := math.Floor((seekerOfferAcceptancePerformance / seekerOfferPerformance) * 100)
			monthlySale.SeekerOfferAcceptanceRate = null.NewFloat(result, true)
		} else {
			monthlySale.SeekerOfferAcceptanceRate = null.NewFloat(0, true)
		}

		// 内定(実績)
		monthlySale.SeekerOfferPerformance = null.NewFloat(seekerOfferPerformance, true)
		// 内定(差異)
		monthlySale.SeekerOfferDifference = null.NewFloat((seekerOfferPerformance - monthlySale.SeekerOfferTarget.Float64), true)
		// 内定率
		if seekerFinalSelectionPerformance > 0 {
			result := math.Floor((seekerOfferPerformance / seekerFinalSelectionPerformance) * 100)
			monthlySale.SeekerOfferRate = null.NewFloat(result, true)
		} else {
			monthlySale.SeekerOfferRate = null.NewFloat(0, true)
		}

		// 最終選考(実績)
		monthlySale.SeekerFinalSelectionPerformance = null.NewFloat(seekerFinalSelectionPerformance, true)
		// 最終選考(差異)
		monthlySale.SeekerFinalSelectionDifference = null.NewFloat((seekerFinalSelectionPerformance - monthlySale.SeekerFinalSelectionTarget.Float64), true)
		// 最終選考実施率
		if seekerSelectionPerformance > 0 {
			result := math.Floor((seekerFinalSelectionPerformance / seekerSelectionPerformance) * 100)
			monthlySale.SeekerFinalSelectionRate = null.NewFloat(result, true)
		} else {
			monthlySale.SeekerFinalSelectionRate = null.NewFloat(0, true)
		}

		// 選考(実績)
		monthlySale.SeekerSelectionPerformance = null.NewFloat(seekerSelectionPerformance, true)
		// 選考(差異)
		monthlySale.SeekerSelectionDifference = null.NewFloat((seekerSelectionPerformance - monthlySale.SeekerSelectionTarget.Float64), true)
		// 選考実施率
		if seekerRecommendationCompletionPerformance > 0 {
			result := math.Floor((seekerSelectionPerformance / seekerRecommendationCompletionPerformance) * 100)
			monthlySale.SeekerSelectionRate = null.NewFloat(result, true)
		} else {
			monthlySale.SeekerSelectionRate = null.NewFloat(0, true)
		}

		// 推薦完了(実績)
		monthlySale.SeekerRecommendationCompletionPerformance = null.NewFloat(seekerRecommendationCompletionPerformance, true)
		// 推薦完了(差異)
		monthlySale.SeekerRecommendationCompletionDifference = null.NewFloat((seekerRecommendationCompletionPerformance - monthlySale.SeekerRecommendationCompletionTarget.Float64), true)
		// 応諾率
		if seekerJobIntroductionPerformance > 0 {
			result := math.Floor((seekerRecommendationCompletionPerformance / seekerJobIntroductionPerformance) * 100)
			monthlySale.SeekerRecommendationCompletionRate = null.NewFloat(result, true)
		} else {
			monthlySale.SeekerRecommendationCompletionRate = null.NewFloat(0, true)
		}

		// 求人紹介(実績)
		monthlySale.SeekerJobIntroductionPerformance = null.NewFloat(seekerJobIntroductionPerformance, true)
		// 求人紹介(差異)
		monthlySale.SeekerJobIntroductionDifference = null.NewFloat((seekerJobIntroductionPerformance - monthlySale.SeekerJobIntroductionTarget.Float64), true)
		// 紹介稼働率
		if interviewPerformance > 0 {
			result := math.Floor((seekerJobIntroductionPerformance / interviewPerformance) * 100)
			monthlySale.SeekerJobIntroductionRate = null.NewFloat(result, true)
		} else {
			monthlySale.SeekerJobIntroductionRate = null.NewFloat(0, true)
		}

		// 面談(実績)
		monthlySale.SeekerInterviewPerformance = null.NewFloat(interviewPerformance, true)
		// 面談(差異)
		monthlySale.SeekerInterviewDifference = null.NewFloat((interviewPerformance - monthlySale.SeekerInterviewTarget.Float64), true)

		/************************ 求職者数ベース（面談実施月ベース） ************************/

		// 内定承諾(実績)
		monthlySale.InterviewOfferAcceptancePerformance = null.NewFloat(interviewOfferAcceptancePerformance, true)

		// 内定承諾(差異)
		monthlySale.InterviewOfferAcceptanceDifference = null.NewFloat((interviewOfferAcceptancePerformance - monthlySale.InterviewOfferAcceptanceTarget.Float64), true)
		// 承諾率
		if interviewOfferPerformance > 0 {
			result := math.Floor((interviewOfferAcceptancePerformance / interviewOfferPerformance) * 100)
			monthlySale.InterviewOfferAcceptanceRate = null.NewFloat(result, true)
		} else {
			monthlySale.InterviewOfferAcceptanceRate = null.NewFloat(0, true)
		}

		// 内定(実績)
		monthlySale.InterviewOfferPerformance = null.NewFloat(interviewOfferPerformance, true)
		// 内定(差異)
		monthlySale.InterviewOfferDifference = null.NewFloat((interviewOfferPerformance - monthlySale.InterviewOfferTarget.Float64), true)
		// 内定率
		if interviewFinalSelectionPerformance > 0 {
			result := math.Floor((interviewOfferPerformance / interviewFinalSelectionPerformance) * 100)
			monthlySale.InterviewOfferRate = null.NewFloat(result, true)
		} else {
			monthlySale.InterviewOfferRate = null.NewFloat(0, true)
		}

		// 最終選考(実績)
		monthlySale.InterviewFinalSelectionPerformance = null.NewFloat(interviewFinalSelectionPerformance, true)
		// 最終選考(差異)
		monthlySale.InterviewFinalSelectionDifference = null.NewFloat((interviewFinalSelectionPerformance - monthlySale.InterviewFinalSelectionTarget.Float64), true)
		// 最終選考実施率
		if interviewSelectionPerformance > 0 {
			result := math.Floor((interviewFinalSelectionPerformance / interviewSelectionPerformance) * 100)
			monthlySale.InterviewFinalSelectionRate = null.NewFloat(result, true)
		} else {
			monthlySale.InterviewFinalSelectionRate = null.NewFloat(0, true)
		}

		// 選考(実績)
		monthlySale.InterviewSelectionPerformance = null.NewFloat(interviewSelectionPerformance, true)
		// 選考(差異)
		monthlySale.InterviewSelectionDifference = null.NewFloat((interviewSelectionPerformance - monthlySale.InterviewSelectionTarget.Float64), true)
		// 選考実施率
		if interviewRecommendationCompletionPerformance > 0 {
			result := math.Floor((interviewSelectionPerformance / interviewRecommendationCompletionPerformance) * 100)
			monthlySale.InterviewSelectionRate = null.NewFloat(result, true)
		} else {
			monthlySale.InterviewSelectionRate = null.NewFloat(0, true)
		}

		// 推薦完了(実績)
		monthlySale.InterviewRecommendationCompletionPerformance = null.NewFloat(interviewRecommendationCompletionPerformance, true)
		// 推薦完了(差異)
		monthlySale.InterviewRecommendationCompletionDifference = null.NewFloat((interviewRecommendationCompletionPerformance - monthlySale.InterviewRecommendationCompletionTarget.Float64), true)
		// 応諾率
		if interviewJobIntroductionPerformance > 0 {
			result := math.Floor((interviewRecommendationCompletionPerformance / interviewJobIntroductionPerformance) * 100)
			monthlySale.InterviewRecommendationCompletionRate = null.NewFloat(result, true)
		} else {
			monthlySale.InterviewRecommendationCompletionRate = null.NewFloat(0, true)
		}

		// 求人紹介(実績)
		monthlySale.InterviewJobIntroductionPerformance = null.NewFloat(interviewJobIntroductionPerformance, true)
		// 求人紹介(差異)
		monthlySale.InterviewJobIntroductionDifference = null.NewFloat((interviewJobIntroductionPerformance - monthlySale.InterviewJobIntroductionTarget.Float64), true)
		// 紹介稼働率
		if interviewPerformance > 0 {
			result := math.Floor((interviewJobIntroductionPerformance / interviewPerformance) * 100)
			monthlySale.InterviewJobIntroductionRate = null.NewFloat(result, true)
		} else {
			monthlySale.InterviewJobIntroductionRate = null.NewFloat(0, true)
		}

		// 面談(実績)
		monthlySale.InterviewInterviewPerformance = null.NewFloat(interviewPerformance, true)
		// 面談(差異)
		monthlySale.InterviewInterviewDifference = null.NewFloat((interviewPerformance - monthlySale.InterviewInterviewTarget.Float64), true)

		monthlyList = append(monthlyList, *monthlySale)
	}

	staffSaleManagement.StaffMonthlySales = monthlyList
	output.StaffSaleManagement = staffSaleManagement

	return output, nil
}

// 求人企業の更新
type GetStaffSaleManagementAndAgentMonthlyByIDInput struct {
	AgentStaffID uint
	ManagementID uint
}

type GetStaffSaleManagementAndAgentMonthlyByIDOutput struct {
	StaffSaleManagement *entity.AgentStaffSaleManagement
}

func (i *AgentStaffMonthlySaleInteractorImpl) GetStaffSaleManagementAndAgentMonthlyByID(input GetStaffSaleManagementAndAgentMonthlyByIDInput) (GetStaffSaleManagementAndAgentMonthlyByIDOutput, error) {
	var (
		output GetStaffSaleManagementAndAgentMonthlyByIDOutput
		err    error
	)

	staffSaleManagement, err := i.agentStaffSaleManagementRepository.FindByManagementIDAndStaffID(input.ManagementID, input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agentStaffMonthlySaleList, err := i.agentStaffMonthlySaleRepository.GetByStaffManagementID(staffSaleManagement.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, monthly := range agentStaffMonthlySaleList {
		staffSaleManagement.StaffMonthlySales = append(staffSaleManagement.StaffMonthlySales, *monthly)
	}

	output.StaffSaleManagement = staffSaleManagement

	return output, nil
}
