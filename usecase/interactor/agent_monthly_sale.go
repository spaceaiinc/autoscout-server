package interactor

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type AgentMonthlySaleInteractor interface {
	// 汎用系 API
	CreateAgentMonthlySale(input CreateAgentMonthlySaleInput) (CreateAgentMonthlySaleOutput, error)
	UpdateAgentMonthlySale(input UpdateAgentMonthlySaleInput) (UpdateAgentMonthlySaleOutput, error)
	GetAgentMonthlySaleList(input GetAgentMonthlySaleListInput) (GetAgentMonthlySaleListOutput, error)
	GetDefaultPreviewAgentMonthlySaleListByAgentID(input GetDefaultPreviewAgentMonthlySaleListByAgentIDInput) (GetDefaultPreviewAgentMonthlySaleListByAgentIDOutput, error)

	// sale_managementsテーブル関連
	GetSaleManagementByID(input GetSaleManagementByIDInput) (GetSaleManagementByIDOutput, error)
	GetSaleManagementListByAgentID(input GetSaleManagementListByAgentIDInput) (GetSaleManagementListByAgentIDOutput, error)
	GetSaleManagementListAndStaffManagementByAgentID(input GetSaleManagementListAndStaffManagementByAgentIDInput) (GetSaleManagementListAndStaffManagementByAgentIDOutput, error)
	GetSaleManagementAndAgentMonthlyByID(input GetSaleManagementAndAgentMonthlyByIDInput) (GetSaleManagementAndAgentMonthlyByIDOutput, error)
	GetSumOfStaffMonthlyByManagementID(input GetSumOfStaffMonthlyByManagementIDInput) (GetSumOfStaffMonthlyByManagementIDOutput, error) // 個人売上の合算値を全体売り上げとして取得する
}

type AgentMonthlySaleInteractorImpl struct {
	firebase                        usecase.Firebase
	sendgrid                        config.Sendgrid
	agentMonthlySaleRepository      usecase.AgentMonthlySaleRepository
	agentStaffMonthlySaleRepository usecase.AgentStaffMonthlySaleRepository
	agentSaleManagementRepository   usecase.AgentSaleManagementRepository
	staffSaleManagementRepository   usecase.AgentStaffSaleManagementRepository
	saleRepository                  usecase.SaleRepository
	taskRepository                  usecase.TaskRepository
	interviewTaskRepository         usecase.InterviewTaskRepository
	interviewTaskGroupRepository    usecase.InterviewTaskGroupRepository
}

// AgentMonthlySaleInteractorImpl is an implementation of AgentMonthlySaleInteractor
func NewAgentMonthlySaleInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	amR usecase.AgentMonthlySaleRepository,
	asmS usecase.AgentStaffMonthlySaleRepository,
	asmR usecase.AgentSaleManagementRepository,
	assmR usecase.AgentStaffSaleManagementRepository,
	sR usecase.SaleRepository,
	tR usecase.TaskRepository,
	itR usecase.InterviewTaskRepository,
	itgR usecase.InterviewTaskGroupRepository,
) AgentMonthlySaleInteractor {
	return &AgentMonthlySaleInteractorImpl{
		firebase:                        fb,
		sendgrid:                        sg,
		agentMonthlySaleRepository:      amR,
		agentStaffMonthlySaleRepository: asmS,
		agentSaleManagementRepository:   asmR,
		staffSaleManagementRepository:   assmR,
		saleRepository:                  sR,
		taskRepository:                  tR,
		interviewTaskRepository:         itR,
		interviewTaskGroupRepository:    itgR,
	}
}

/****************************************************************************************/
// 汎用系API
//

// 売上情報の登録
type CreateAgentMonthlySaleInput struct {
	CreateParam  entity.CreateOrUpdateAgentMonthlyManagementParam
	ManagementID uint
}

type CreateAgentMonthlySaleOutput struct {
	OK bool
}

func (i *AgentMonthlySaleInteractorImpl) CreateAgentMonthlySale(input CreateAgentMonthlySaleInput) (CreateAgentMonthlySaleOutput, error) {
	var (
		output CreateAgentMonthlySaleOutput
		err    error
	)

	agentSaleManagement := entity.NewAgentSaleManagement(
		input.CreateParam.AgentID,
		input.CreateParam.FiscalYear,
		input.CreateParam.IsOpen,
	)
	err = i.agentSaleManagementRepository.Create(agentSaleManagement)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 作成する売上情報がOpenの場合はその他をCloseにする
	if input.CreateParam.IsOpen {
		err = i.agentSaleManagementRepository.UpdateIsOpenOtherThanID(agentSaleManagement.AgentID, agentSaleManagement.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, monthlySale := range input.CreateParam.AgentMonthlySales {
		agentMonthlySale := entity.NewAgentMonthlySale(
			agentSaleManagement.ID,
			monthlySale.SalesMonth,
			monthlySale.OrderSalesBudget,
			monthlySale.OrderCostBudget,
			monthlySale.OrderGrossProfitBudget,
			monthlySale.OrderAssumedUnitPrice,
			monthlySale.OrderExpectedOfferAcceptance,
			monthlySale.ClaimSalesRevenueBudget,
			monthlySale.ClaimCostBudget,
			monthlySale.ClaimGrossMarginBudget,
			monthlySale.ClaimAssumedUnitPrice,
			monthlySale.ClaimExpectedNewEmployeeNumber,
			monthlySale.SeekerOfferAcceptanceTarget,
			monthlySale.SeekerOfferTarget,
			monthlySale.SeekerFinalSelectionTarget,
			monthlySale.SeekerSelectionTarget,
			monthlySale.SeekerRecommendationCompletionTarget,
			monthlySale.SeekerJobIntroductionTarget,
			monthlySale.SeekerInterviewTarget,
			monthlySale.ActiveOfferTarget,
			monthlySale.ActiveFinalSelectionTarget,
			monthlySale.ActiveSelectionTarget,
			monthlySale.ActiveRecommendationCompletionTarget,
			monthlySale.ActiveJobIntroductionTarget,
			monthlySale.InterviewOfferAcceptanceTarget,
			monthlySale.InterviewOfferTarget,
			monthlySale.InterviewFinalSelectionTarget,
			monthlySale.InterviewSelectionTarget,
			monthlySale.InterviewRecommendationCompletionTarget,
			monthlySale.InterviewJobIntroductionTarget,
			monthlySale.InterviewInterviewTarget,
		)
		err = i.agentMonthlySaleRepository.Create(agentMonthlySale)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.OK = true
	return output, nil
}

// 売上情報の更新
type UpdateAgentMonthlySaleInput struct {
	UpdateParam entity.CreateOrUpdateAgentMonthlyManagementParam
}

type UpdateAgentMonthlySaleOutput struct {
	OK bool
}

func (i *AgentMonthlySaleInteractorImpl) UpdateAgentMonthlySale(input UpdateAgentMonthlySaleInput) (UpdateAgentMonthlySaleOutput, error) {
	var (
		output UpdateAgentMonthlySaleOutput
		err    error
	)

	saleManagement := entity.NewAgentSaleManagement(
		input.UpdateParam.AgentID,
		input.UpdateParam.FiscalYear,
		input.UpdateParam.IsOpen,
	)
	saleManagement.ID = input.UpdateParam.ID

	// 更新の売上情報がOpenの場合はその他をCloseにする
	if input.UpdateParam.IsOpen {
		err = i.agentSaleManagementRepository.UpdateIsOpenOtherThanID(saleManagement.AgentID, saleManagement.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.agentSaleManagementRepository.Update(saleManagement, saleManagement.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, monthlySale := range input.UpdateParam.AgentMonthlySales {
		err = i.agentMonthlySaleRepository.Update(&monthlySale, monthlySale.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.OK = true

	return output, nil
}

// 売上情報の取得
type GetAgentMonthlySaleListInput struct {
	AgentID      uint
	ManagementID uint
}

type GetAgentMonthlySaleListOutput struct {
	AgentSaleManagement       *entity.AgentSaleManagement
	AgentStaffSaleManagements []*entity.AgentStaffSaleManagement
	AgentMonthlySales         []*entity.AgentMonthlySale
}

func (i *AgentMonthlySaleInteractorImpl) GetAgentMonthlySaleList(input GetAgentMonthlySaleListInput) (GetAgentMonthlySaleListOutput, error) {
	var (
		output GetAgentMonthlySaleListOutput
		err    error
		// monthlyList []entity.AgentMonthlySale
	)

	agentSaleManagement, err := i.agentSaleManagementRepository.FindByIDAndAgentID(input.ManagementID, input.AgentID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			return output, wrapped
		} else {
			fmt.Println(err)
			return output, err
		}
	}

	staffSaleManagements, err := i.staffSaleManagementRepository.GetStaffNameByManagementID(agentSaleManagement.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	layout := "2006-01" // 入力文字列のフォーマット
	t, err := time.Parse(layout, agentSaleManagement.FiscalYear)
	if err != nil {
		fmt.Println("Error:", err)
	}
	previousMonth := t.AddDate(0, -11, 0)
	startMonth := previousMonth.Format("2006-01")

	agentMonthlySaleList, err := i.agentMonthlySaleRepository.GetByManagementID(input.ManagementID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 指定期間の売上情報を取得
	saleList, err := i.saleRepository.GetByAgentIDForMonthly(input.AgentID, startMonth, agentSaleManagement.FiscalYear)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, monthlySale := range agentMonthlySaleList {
		if monthlySale.SalesMonth == "" {
			continue
		}

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
		)

		for _, sale := range saleList {
			// 受注関連の計算
			if monthlySale.SalesMonth == sale.ContractSignedMonth {
				var (
					caSale    float64 = 0 // CAの売上
					raSale    float64 = 0 // RAの売上
					agentSale float64 = 0 // エージェントの売上

					caCost    float64 = 0 // CAの原価
					raCost    float64 = 0 // RAの原価
					agentCost float64 = 0 // エージェントの原価

					caGrossProfit    float64 = 0 // CAの粗利
					raGrossProfit    float64 = 0 // RAの粗利
					agentGrossProfit float64 = 0 // エージェントの粗利
				)

				// CA売上比率に応じた売上を計算
				if input.AgentID == sale.CAAgentID && sale.CaSalesRatio.Float64 > 0 {
					caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
					caCost = sale.Cost.Float64 * sale.CaSalesRatio.Float64 / 100
					caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
				}

				// RA売上比率に応じた売上を計算
				if input.AgentID == sale.RAAgentID && sale.RaSalesRatio.Float64 > 0 {
					raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
					raCost = sale.Cost.Float64 * sale.RaSalesRatio.Float64 / 100
					raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
				}

				// エージェントの売上
				agentSale = caSale + raSale
				agentCost = caCost + raCost
				agentGrossProfit = caGrossProfit + raGrossProfit

				// 内定承諾の計算
				if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
					orderSalesPerformance = orderSalesPerformance + agentSale                    // 受注売上実績の計算
					orderCostPerformance = orderCostPerformance + agentCost                      // 原価実績の計算
					orderGrossProfitPerformance = orderGrossProfitPerformance + agentGrossProfit // 受注粗利実績の計算

					// 自社がCA担当分のみをカウント
					if input.AgentID == sale.CAAgentID {
						orderOfferAcceptancePerformance = orderOfferAcceptancePerformance + 1 // 内定承諾実績
						seekerOfferAcceptancePerformance = seekerOfferAcceptancePerformance + 1
					}
				}

				// Aヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
					orderAccuracyA = orderAccuracyA + agentSale // Aヨミの請求金額合計
				}

				// Bヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
					orderAccuracyB = orderAccuracyB + agentSale // Bヨミの請求金額合計
				}

				// Cヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
					orderAccuracyC = orderAccuracyC + agentSale // Cヨミの請求金額合計
				}

				// ネタの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
					orderAccuracyTopic = orderAccuracyTopic + agentSale // ネタの請求金額合計
				}

				// ヨミのIDリストを作成
				orderAccuracyIDList = append(orderAccuracyIDList, sale.ID)
			}

			// 請求関連の計算
			if monthlySale.SalesMonth == sale.BillingMonth {
				var (
					caSale    float64 = 0 // CAの売上
					raSale    float64 = 0 // RAの売上
					agentSale float64 = 0 // エージェントの売上

					caCost    float64 = 0 // CAの原価
					raCost    float64 = 0 // RAの原価
					agentCost float64 = 0 // エージェントの原価

					caGrossProfit    float64 = 0 // CAの粗利
					raGrossProfit    float64 = 0 // RAの粗利
					agentGrossProfit float64 = 0 // エージェントの粗利
				)

				// CA売上比率に応じた売上を計算
				if input.AgentID == sale.CAAgentID && sale.CaSalesRatio.Float64 > 0 {
					caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
					caCost = sale.Cost.Float64 * sale.CaSalesRatio.Float64 / 100
					caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
				}

				// RA売上比率に応じた売上を計算
				if input.AgentID == sale.RAAgentID && sale.RaSalesRatio.Float64 > 0 {
					raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
					raCost = sale.Cost.Float64 * sale.RaSalesRatio.Float64 / 100
					raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
				}

				// エージェントの売上
				agentSale = caSale + raSale
				agentCost = caCost + raCost
				agentGrossProfit = caGrossProfit + raGrossProfit

				// 内定承諾の計算
				if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
					claimSalesPerformance = claimSalesPerformance + agentSale                    // 請求売上実績の計算
					claimCostPerformance = claimCostPerformance + agentCost                      // 原価実績の計算
					claimGrossProfitPerformance = claimGrossProfitPerformance + agentGrossProfit // 請求粗利実績の計算

					// 自社がCA担当分のみをカウント
					if input.AgentID == sale.CAAgentID {
						claimNewEmployeeNumbererformance = claimNewEmployeeNumbererformance + 1 // 入社人数実績
					}
				}

				// Aヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
					claimAccuracyA = claimAccuracyA + agentSale // Aヨミの請求金額合計
				}

				// Bヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
					claimAccuracyB = claimAccuracyB + agentSale // Bヨミの請求金額合計
				}

				// Cヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
					claimAccuracyC = claimAccuracyC + agentSale // Cヨミの請求金額合計
				}

				// ネタの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
					claimAccuracyTopic = claimAccuracyTopic + agentSale // ネタの請求金額合計
				}

				// ヨミのIDリストを作成
				claimAccuracyIDList = append(claimAccuracyIDList, sale.ID)
			}
		}

		/********* 求職者数（当月カウントベース） **********/

		// 内定承諾数(seekerOfferAcceptancePerformance)
		// 受注売上項目の内定承諾実績と同じになる

		// 内定数
		seekerOfferPerformance, err = i.taskRepository.GetSeekerOfferPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 最終選考数
		seekerFinalSelectionPerformance, err = i.taskRepository.GetSeekerFinalSelectionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 選考数
		seekerSelectionPerformance, err = i.taskRepository.GetSeekerSelectionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 推薦完了数
		seekerRecommendationCompletionPerformance, err = i.taskRepository.GetSeekerRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求人紹介数
		seekerJobIntroductionPerformance, err = i.taskRepository.GetSeekerJobIntroductionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		/********* 求職者数（面談実施月カウントベース） **********/

		// // 内定承諾数
		interviewOfferAcceptancePerformance, err = i.interviewTaskGroupRepository.GetInterviewOfferAcceptancePerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 内定数
		interviewOfferPerformance, err = i.taskRepository.GetInterviewOfferPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 最終選考数
		interviewFinalSelectionPerformance, err = i.taskRepository.GetInterviewFinalSelectionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 選考数
		interviewSelectionPerformance, err = i.taskRepository.GetInterviewSelectionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 推薦完了数
		interviewRecommendationCompletionPerformance, err = i.taskRepository.GetInterviewRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求人紹介数
		interviewJobIntroductionPerformance, err = i.taskRepository.GetInterviewJobIntroductionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 面談数（当月カウントベースと面談実施月ベースで統一）
		interviewPerformance, err = i.interviewTaskGroupRepository.GetInterviewPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
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
	}

	output.AgentSaleManagement = agentSaleManagement
	output.AgentMonthlySales = agentMonthlySaleList
	output.AgentStaffSaleManagements = staffSaleManagements

	return output, nil
}

// 売上情報の取得
type GetDefaultPreviewAgentMonthlySaleListByAgentIDInput struct {
	AgentID uint
}

type GetDefaultPreviewAgentMonthlySaleListByAgentIDOutput struct {
	AgentSaleManagement       *entity.AgentSaleManagement
	AgentStaffSaleManagements []*entity.AgentStaffSaleManagement
	AgentMonthlySales         []*entity.AgentMonthlySale
}

func (i *AgentMonthlySaleInteractorImpl) GetDefaultPreviewAgentMonthlySaleListByAgentID(input GetDefaultPreviewAgentMonthlySaleListByAgentIDInput) (GetDefaultPreviewAgentMonthlySaleListByAgentIDOutput, error) {
	var (
		output GetDefaultPreviewAgentMonthlySaleListByAgentIDOutput
		err    error
		// monthlyList []entity.AgentMonthlySale
	)

	agentSaleManagement, err := i.agentSaleManagementRepository.FindByAgentIDAndIsOpen(input.AgentID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			// Openの情報がない場合は最新のものを取得
			_, err = i.agentSaleManagementRepository.FindLatestByAgentID(input.AgentID)
			if err != nil {
				if errors.Is(err, entity.ErrNotFound) {
					fmt.Println(err)
					return output, nil
				} else {
					fmt.Println(err)
					return output, err
				}
			} else {
				fmt.Println(err)
				return output, err
			}
		} else {
			fmt.Println(err)
			return output, err
		}
	}

	staffSaleManagements, err := i.staffSaleManagementRepository.GetStaffNameByManagementID(agentSaleManagement.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// for _, staffSale := range staffSaleManagements {
	// 	agentSaleManagement.AgentStaffSaleManagements = append(agentSaleManagement.AgentStaffSaleManagements, *staffSale)
	// }

	layout := "2006-01" // 入力文字列のフォーマット
	t, err := time.Parse(layout, agentSaleManagement.FiscalYear)
	if err != nil {
		fmt.Println("Error:", err)
	}
	previousMonth := t.AddDate(0, -11, 0)
	startMonth := previousMonth.Format("2006-01")

	agentMonthlySaleList, err := i.agentMonthlySaleRepository.GetByManagementID(agentSaleManagement.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 指定期間の売上情報を取得
	saleList, err := i.saleRepository.GetByAgentIDForMonthly(input.AgentID, startMonth, agentSaleManagement.FiscalYear)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, monthlySale := range agentMonthlySaleList {
		if monthlySale.SalesMonth == "" {
			continue
		}

		//　計算結果を代入する変数を定美
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
		)

		for _, sale := range saleList {
			// 受注関連の計算
			if monthlySale.SalesMonth == sale.ContractSignedMonth {
				var (
					caSale    float64 = 0 // CAの売上
					raSale    float64 = 0 // RAの売上
					agentSale float64 = 0 // エージェントの売上

					caCost    float64 = 0 // CAの原価
					raCost    float64 = 0 // RAの原価
					agentCost float64 = 0 // エージェントの原価

					caGrossProfit    float64 = 0 // CAの粗利
					raGrossProfit    float64 = 0 // RAの粗利
					agentGrossProfit float64 = 0 // エージェントの粗利
				)

				// CA売上比率に応じた売上を計算
				if input.AgentID == sale.CAAgentID && sale.CaSalesRatio.Float64 > 0 {
					caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
					caCost = sale.Cost.Float64 * sale.CaSalesRatio.Float64 / 100
					caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
				}

				// RA売上比率に応じた売上を計算
				if input.AgentID == sale.RAAgentID && sale.RaSalesRatio.Float64 > 0 {
					raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
					raCost = sale.Cost.Float64 * sale.RaSalesRatio.Float64 / 100
					raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
				}

				// エージェントの売上
				agentSale = caSale + raSale
				agentCost = caCost + raCost
				agentGrossProfit = caGrossProfit + raGrossProfit

				// 内定承諾の計算
				if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
					orderSalesPerformance = orderSalesPerformance + agentSale                    // 受注売上実績の計算
					orderCostPerformance = orderCostPerformance + agentCost                      // 原価実績の計算
					orderGrossProfitPerformance = orderGrossProfitPerformance + agentGrossProfit // 受注粗利実績の計算

					// 自社がCA担当分のみをカウント
					if input.AgentID == sale.CAAgentID {
						orderOfferAcceptancePerformance = orderOfferAcceptancePerformance + 1 // 内定承諾実績
						seekerOfferAcceptancePerformance = seekerOfferAcceptancePerformance + 1
					}
				}

				// Aヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
					orderAccuracyA = orderAccuracyA + agentSale // Aヨミの請求金額合計
				}

				// Bヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
					orderAccuracyB = orderAccuracyB + agentSale // Bヨミの請求金額合計
				}

				// Cヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
					orderAccuracyC = orderAccuracyC + agentSale // Cヨミの請求金額合計
				}

				// ネタの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
					orderAccuracyTopic = orderAccuracyTopic + agentSale // ネタの請求金額合計
				}

				// ヨミのIDリストを作成
				orderAccuracyIDList = append(orderAccuracyIDList, sale.ID)
			}

			// 請求関連の計算
			if monthlySale.SalesMonth == sale.BillingMonth {
				var (
					caSale    float64 = 0 // CAの売上
					raSale    float64 = 0 // RAの売上
					agentSale float64 = 0 // エージェントの売上

					caCost    float64 = 0 // CAの原価
					raCost    float64 = 0 // RAの原価
					agentCost float64 = 0 // エージェントの原価

					caGrossProfit    float64 = 0 // CAの粗利
					raGrossProfit    float64 = 0 // RAの粗利
					agentGrossProfit float64 = 0 // エージェントの粗利
				)

				// CA売上比率に応じた売上を計算
				if input.AgentID == sale.CAAgentID && sale.CaSalesRatio.Float64 > 0 {
					caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
					caCost = sale.Cost.Float64 * sale.CaSalesRatio.Float64 / 100
					caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
				}

				// RA売上比率に応じた売上を計算
				if input.AgentID == sale.RAAgentID && sale.RaSalesRatio.Float64 > 0 {
					raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
					raCost = sale.Cost.Float64 * sale.RaSalesRatio.Float64 / 100
					raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
				}

				// エージェントの売上
				agentSale = caSale + raSale
				agentCost = caCost + raCost
				agentGrossProfit = caGrossProfit + raGrossProfit

				// 内定承諾の計算
				if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
					claimSalesPerformance = claimSalesPerformance + agentSale                    // 請求売上実績の計算
					claimCostPerformance = claimCostPerformance + agentCost                      // 原価実績の計算
					claimGrossProfitPerformance = claimGrossProfitPerformance + agentGrossProfit // 請求粗利実績の計算

					// 自社がCA担当分のみをカウント
					if input.AgentID == sale.CAAgentID {
						claimNewEmployeeNumbererformance = claimNewEmployeeNumbererformance + 1 // 入社人数実績
					}
				}

				// Aヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
					claimAccuracyA = claimAccuracyA + agentSale // Aヨミの請求金額合計
				}

				// Bヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
					claimAccuracyB = claimAccuracyB + agentSale // Bヨミの請求金額合計
				}

				// Cヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
					claimAccuracyC = claimAccuracyC + agentSale // Cヨミの請求金額合計
				}

				// ネタの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
					claimAccuracyTopic = claimAccuracyTopic + agentSale // ネタの請求金額合計
				}

				// ヨミのIDリストを作成
				claimAccuracyIDList = append(claimAccuracyIDList, sale.ID)
			}
		}

		/********* 求職者数（当月カウントベース） **********/

		// 内定承諾数(seekerOfferAcceptancePerformance)
		// 受注売上項目の内定承諾実績と同じになる

		// 内定数
		seekerOfferPerformance, err = i.taskRepository.GetSeekerOfferPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 最終選考数
		seekerFinalSelectionPerformance, err = i.taskRepository.GetSeekerFinalSelectionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 選考数
		seekerSelectionPerformance, err = i.taskRepository.GetSeekerSelectionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 推薦完了数
		seekerRecommendationCompletionPerformance, err = i.taskRepository.GetSeekerRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求人紹介数
		seekerJobIntroductionPerformance, err = i.taskRepository.GetSeekerJobIntroductionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		/********* 求職者数（面談実施月カウントベース） **********/

		// // 内定承諾数
		interviewOfferAcceptancePerformance, err = i.interviewTaskGroupRepository.GetInterviewOfferAcceptancePerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 内定数
		interviewOfferPerformance, err = i.taskRepository.GetInterviewOfferPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 最終選考数
		interviewFinalSelectionPerformance, err = i.taskRepository.GetInterviewFinalSelectionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 選考数
		interviewSelectionPerformance, err = i.taskRepository.GetInterviewSelectionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 推薦完了数
		interviewRecommendationCompletionPerformance, err = i.taskRepository.GetInterviewRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求人紹介数
		interviewJobIntroductionPerformance, err = i.taskRepository.GetInterviewJobIntroductionPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 面談数（当月カウントベースと面談実施月ベースで統一）
		interviewPerformance, err = i.interviewTaskGroupRepository.GetInterviewPerformanceCountByAgentIDAndSalesMonth(input.AgentID, monthlySale.SalesMonth)
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
	}

	output.AgentSaleManagement = agentSaleManagement
	output.AgentMonthlySales = agentMonthlySaleList
	output.AgentStaffSaleManagements = staffSaleManagements

	return output, nil
}

/****************************************************************************************/
// sale_managementsテーブル関連 API
//

// 売上情報一覧の取得
type GetSaleManagementListByAgentIDInput struct {
	AgentID uint
}

type GetSaleManagementListByAgentIDOutput struct {
	AgentSaleManagementList []*entity.AgentSaleManagement
}

func (i *AgentMonthlySaleInteractorImpl) GetSaleManagementListByAgentID(input GetSaleManagementListByAgentIDInput) (GetSaleManagementListByAgentIDOutput, error) {
	var (
		output GetSaleManagementListByAgentIDOutput
		err    error
	)

	agentSaleManagementList, err := i.agentSaleManagementRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.AgentSaleManagementList = agentSaleManagementList

	return output, nil
}

// 売上情報一覧の取得
type GetSaleManagementListAndStaffManagementByAgentIDInput struct {
	AgentID uint
}

type GetSaleManagementListAndStaffManagementByAgentIDOutput struct {
	AgentSaleManagementList []*entity.AgentSaleManagement
}

func (i *AgentMonthlySaleInteractorImpl) GetSaleManagementListAndStaffManagementByAgentID(input GetSaleManagementListAndStaffManagementByAgentIDInput) (GetSaleManagementListAndStaffManagementByAgentIDOutput, error) {
	var (
		output GetSaleManagementListAndStaffManagementByAgentIDOutput
		err    error
	)

	agentSaleManagementList, err := i.agentSaleManagementRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	staffSaleManagementList, err := i.staffSaleManagementRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, agentSale := range agentSaleManagementList {
		for _, staffSale := range staffSaleManagementList {
			if agentSale.ID == staffSale.ManagementID {
				agentSale.AgentStaffSaleManagements = append(agentSale.AgentStaffSaleManagements, *staffSale)
			}
		}
	}

	output.AgentSaleManagementList = agentSaleManagementList

	return output, nil
}

// 売上情報一覧の取得
type GetSaleManagementByIDInput struct {
	ManagementID uint
}

type GetSaleManagementByIDOutput struct {
	AgentSaleManagement *entity.AgentSaleManagement
}

func (i *AgentMonthlySaleInteractorImpl) GetSaleManagementByID(input GetSaleManagementByIDInput) (GetSaleManagementByIDOutput, error) {
	var (
		output GetSaleManagementByIDOutput
		err    error
	)

	saleManagement, err := i.agentSaleManagementRepository.FindByID(input.ManagementID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.AgentSaleManagement = saleManagement

	return output, nil
}

// 売上情報一覧の取得
type GetSaleManagementAndAgentMonthlyByIDInput struct {
	ManagementID uint
}

type GetSaleManagementAndAgentMonthlyByIDOutput struct {
	AgentSaleManagement *entity.AgentSaleManagement
}

func (i *AgentMonthlySaleInteractorImpl) GetSaleManagementAndAgentMonthlyByID(input GetSaleManagementAndAgentMonthlyByIDInput) (GetSaleManagementAndAgentMonthlyByIDOutput, error) {
	var (
		output GetSaleManagementAndAgentMonthlyByIDOutput
		err    error

		monthlySaleList []entity.AgentMonthlySale
	)

	saleManagement, err := i.agentSaleManagementRepository.FindByID(input.ManagementID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agentMonthlySaleList, err := i.agentMonthlySaleRepository.GetByManagementID(input.ManagementID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, monthlySale := range agentMonthlySaleList {
		monthlySaleList = append(monthlySaleList, *monthlySale)
	}

	saleManagement.AgentMonthlySales = monthlySaleList
	output.AgentSaleManagement = saleManagement

	return output, nil
}

// 個人売上の合算値を全体売り上げとして取得する
type GetSumOfStaffMonthlyByManagementIDInput struct {
	ManagementID uint
}

type GetSumOfStaffMonthlyByManagementIDOutput struct {
	AgentSaleManagement *entity.AgentSaleManagement
}

func (i *AgentMonthlySaleInteractorImpl) GetSumOfStaffMonthlyByManagementID(input GetSumOfStaffMonthlyByManagementIDInput) (GetSumOfStaffMonthlyByManagementIDOutput, error) {
	var (
		output GetSumOfStaffMonthlyByManagementIDOutput
		err    error

		monthlySaleList                              []entity.AgentMonthlySale
		orderSalesBudgetTotal                        = make(map[string]float64)
		orderCostBudgetTotal                         = make(map[string]float64)
		orderGrossProfitBudgetTotal                  = make(map[string]float64)
		orderAssumedUnitPriceTotal                   = make(map[string]float64)
		orderExpectedOfferAcceptanceTotal            = make(map[string]float64)
		claimSalesRevenueBudgetTotal                 = make(map[string]float64)
		claimCostBudgetTotal                         = make(map[string]float64)
		claimGrossMarginBudgetTotal                  = make(map[string]float64)
		claimAssumedUnitPriceTotal                   = make(map[string]float64)
		claimExpectedNewEmployeeNumberTotal          = make(map[string]float64)
		seekerOfferAcceptanceTargetTotal             = make(map[string]float64)
		seekerOfferTargetTotal                       = make(map[string]float64)
		seekerFinalSelectionTargetTotal              = make(map[string]float64)
		seekerSelectionTargetTotal                   = make(map[string]float64)
		seekerRecommendationCompletionTargetTotal    = make(map[string]float64)
		seekerJobIntroductionTargetTotal             = make(map[string]float64)
		seekerInterviewTargetTotal                   = make(map[string]float64)
		interviewOfferAcceptanceTargetTotal          = make(map[string]float64)
		interviewOfferTargetTotal                    = make(map[string]float64)
		interviewFinalSelectionTargetTotal           = make(map[string]float64)
		interviewSelectionTargetTotal                = make(map[string]float64)
		interviewRecommendationCompletionTargetTotal = make(map[string]float64)
		interviewJobIntroductionTargetTotal          = make(map[string]float64)
		interviewInterviewTargetTotal                = make(map[string]float64)
	)

	saleManagement, err := i.agentSaleManagementRepository.FindByID(input.ManagementID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agentMonthlySaleList, err := i.agentMonthlySaleRepository.GetByManagementID(input.ManagementID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	staffMonthlySaleList, err := i.agentStaffMonthlySaleRepository.GetByAgentSaleManagementID(input.ManagementID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 個人売上を月別で合算する
	for _, staffMonthly := range staffMonthlySaleList {
		orderSalesBudgetTotal[staffMonthly.SalesMonth] = orderSalesBudgetTotal[staffMonthly.SalesMonth] + staffMonthly.OrderSalesBudget.Float64
		orderCostBudgetTotal[staffMonthly.SalesMonth] = orderCostBudgetTotal[staffMonthly.SalesMonth] + staffMonthly.OrderCostBudget.Float64
		orderGrossProfitBudgetTotal[staffMonthly.SalesMonth] = orderGrossProfitBudgetTotal[staffMonthly.SalesMonth] + staffMonthly.OrderGrossProfitBudget.Float64
		orderAssumedUnitPriceTotal[staffMonthly.SalesMonth] = orderAssumedUnitPriceTotal[staffMonthly.SalesMonth] + staffMonthly.OrderAssumedUnitPrice.Float64
		orderExpectedOfferAcceptanceTotal[staffMonthly.SalesMonth] = orderExpectedOfferAcceptanceTotal[staffMonthly.SalesMonth] + staffMonthly.OrderExpectedOfferAcceptance.Float64
		claimSalesRevenueBudgetTotal[staffMonthly.SalesMonth] = claimSalesRevenueBudgetTotal[staffMonthly.SalesMonth] + staffMonthly.ClaimSalesRevenueBudget.Float64
		claimCostBudgetTotal[staffMonthly.SalesMonth] = claimCostBudgetTotal[staffMonthly.SalesMonth] + staffMonthly.ClaimCostBudget.Float64
		claimGrossMarginBudgetTotal[staffMonthly.SalesMonth] = claimGrossMarginBudgetTotal[staffMonthly.SalesMonth] + staffMonthly.ClaimGrossMarginBudget.Float64
		claimAssumedUnitPriceTotal[staffMonthly.SalesMonth] = claimAssumedUnitPriceTotal[staffMonthly.SalesMonth] + staffMonthly.ClaimAssumedUnitPrice.Float64
		claimExpectedNewEmployeeNumberTotal[staffMonthly.SalesMonth] = claimExpectedNewEmployeeNumberTotal[staffMonthly.SalesMonth] + staffMonthly.ClaimExpectedNewEmployeeNumber.Float64
		seekerOfferAcceptanceTargetTotal[staffMonthly.SalesMonth] = seekerOfferAcceptanceTargetTotal[staffMonthly.SalesMonth] + staffMonthly.SeekerOfferAcceptanceTarget.Float64
		seekerOfferTargetTotal[staffMonthly.SalesMonth] = seekerOfferTargetTotal[staffMonthly.SalesMonth] + staffMonthly.SeekerOfferTarget.Float64
		seekerFinalSelectionTargetTotal[staffMonthly.SalesMonth] = seekerFinalSelectionTargetTotal[staffMonthly.SalesMonth] + staffMonthly.SeekerFinalSelectionTarget.Float64
		seekerSelectionTargetTotal[staffMonthly.SalesMonth] = seekerSelectionTargetTotal[staffMonthly.SalesMonth] + staffMonthly.SeekerSelectionTarget.Float64
		seekerRecommendationCompletionTargetTotal[staffMonthly.SalesMonth] = seekerRecommendationCompletionTargetTotal[staffMonthly.SalesMonth] + staffMonthly.SeekerRecommendationCompletionTarget.Float64
		seekerJobIntroductionTargetTotal[staffMonthly.SalesMonth] = seekerJobIntroductionTargetTotal[staffMonthly.SalesMonth] + staffMonthly.SeekerJobIntroductionTarget.Float64
		seekerInterviewTargetTotal[staffMonthly.SalesMonth] = seekerInterviewTargetTotal[staffMonthly.SalesMonth] + staffMonthly.SeekerInterviewTarget.Float64
		interviewOfferAcceptanceTargetTotal[staffMonthly.SalesMonth] = interviewOfferAcceptanceTargetTotal[staffMonthly.SalesMonth] + staffMonthly.InterviewOfferAcceptanceTarget.Float64
		interviewOfferTargetTotal[staffMonthly.SalesMonth] = interviewOfferTargetTotal[staffMonthly.SalesMonth] + staffMonthly.InterviewOfferTarget.Float64
		interviewFinalSelectionTargetTotal[staffMonthly.SalesMonth] = interviewFinalSelectionTargetTotal[staffMonthly.SalesMonth] + staffMonthly.InterviewFinalSelectionTarget.Float64
		interviewSelectionTargetTotal[staffMonthly.SalesMonth] = interviewSelectionTargetTotal[staffMonthly.SalesMonth] + staffMonthly.InterviewSelectionTarget.Float64
		interviewRecommendationCompletionTargetTotal[staffMonthly.SalesMonth] = interviewRecommendationCompletionTargetTotal[staffMonthly.SalesMonth] + staffMonthly.InterviewRecommendationCompletionTarget.Float64
		interviewJobIntroductionTargetTotal[staffMonthly.SalesMonth] = interviewJobIntroductionTargetTotal[staffMonthly.SalesMonth] + staffMonthly.InterviewJobIntroductionTarget.Float64
		interviewInterviewTargetTotal[staffMonthly.SalesMonth] = interviewInterviewTargetTotal[staffMonthly.SalesMonth] + staffMonthly.InterviewInterviewTarget.Float64
	}

	// 合算値を全体の売上に代入
	for _, agentMonthly := range agentMonthlySaleList {
		agentMonthly.OrderSalesBudget = null.NewFloat(orderSalesBudgetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.OrderCostBudget = null.NewFloat(orderCostBudgetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.OrderGrossProfitBudget = null.NewFloat(orderGrossProfitBudgetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.OrderAssumedUnitPrice = null.NewFloat(orderAssumedUnitPriceTotal[agentMonthly.SalesMonth], true)
		agentMonthly.OrderExpectedOfferAcceptance = null.NewFloat(orderExpectedOfferAcceptanceTotal[agentMonthly.SalesMonth], true)
		agentMonthly.ClaimSalesRevenueBudget = null.NewFloat(claimSalesRevenueBudgetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.ClaimCostBudget = null.NewFloat(claimCostBudgetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.ClaimGrossMarginBudget = null.NewFloat(claimGrossMarginBudgetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.ClaimAssumedUnitPrice = null.NewFloat(claimAssumedUnitPriceTotal[agentMonthly.SalesMonth], true)
		agentMonthly.ClaimExpectedNewEmployeeNumber = null.NewFloat(claimExpectedNewEmployeeNumberTotal[agentMonthly.SalesMonth], true)
		agentMonthly.SeekerOfferAcceptanceTarget = null.NewFloat(seekerOfferAcceptanceTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.SeekerOfferTarget = null.NewFloat(seekerOfferTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.SeekerFinalSelectionTarget = null.NewFloat(seekerFinalSelectionTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.SeekerSelectionTarget = null.NewFloat(seekerSelectionTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.SeekerRecommendationCompletionTarget = null.NewFloat(seekerRecommendationCompletionTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.SeekerJobIntroductionTarget = null.NewFloat(seekerJobIntroductionTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.SeekerInterviewTarget = null.NewFloat(seekerInterviewTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.InterviewOfferAcceptanceTarget = null.NewFloat(interviewOfferAcceptanceTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.InterviewOfferTarget = null.NewFloat(interviewOfferTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.InterviewFinalSelectionTarget = null.NewFloat(interviewFinalSelectionTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.InterviewSelectionTarget = null.NewFloat(interviewSelectionTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.InterviewRecommendationCompletionTarget = null.NewFloat(interviewRecommendationCompletionTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.InterviewJobIntroductionTarget = null.NewFloat(interviewJobIntroductionTargetTotal[agentMonthly.SalesMonth], true)
		agentMonthly.InterviewInterviewTarget = null.NewFloat(interviewInterviewTargetTotal[agentMonthly.SalesMonth], true)

		monthlySaleList = append(monthlySaleList, *agentMonthly)
	}

	saleManagement.AgentMonthlySales = monthlySaleList
	output.AgentSaleManagement = saleManagement

	return output, nil
}
