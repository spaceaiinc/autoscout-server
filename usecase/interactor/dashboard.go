package interactor

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type DashboardInteractor interface {
	// デフォルト表示
	GetDashboardForDefaultPreview(input GetDashboardForDefaultPreviewInput) (GetDashboardForDefaultPreviewOutput, error)
	GetSaleListForDefaultPreview(input GetSaleListForDefaultPreviewInput) (GetSaleListForDefaultPreviewOutput, error)
	GetReleaseJobSeekerListForDefaultPreview(input GetReleaseJobSeekerListForDefaultPreviewInput) (GetReleaseJobSeekerListForDefaultPreviewOutput, error)

	// 検索
	GetDashboardForSearch(input GetDashboardForSearchInput) (GetDashboardForSearchOutput, error)
	GetSaleListForSearch(input GetSaleListForSearchInput) (GetSaleListForSearchOutput, error)
	GetReleaseJobSeekerListForSearch(input GetReleaseJobSeekerListForSearchInput) (GetReleaseJobSeekerListForSearchOutput, error)

	// CSV
	ExportAccuracyCSV(input ExportAccuracyCSVInput) (ExportAccuracyCSVOutput, error)
}

type DashboardInteractorImpl struct {
	firebase                                usecase.Firebase
	sendgrid                                config.Sendgrid
	agentMonthlySaleRepository              usecase.AgentMonthlySaleRepository
	agentStaffMonthlySaleRepository         usecase.AgentStaffMonthlySaleRepository
	agentSaleManagementRepository           usecase.AgentSaleManagementRepository
	staffSaleManagementRepository           usecase.AgentStaffSaleManagementRepository
	saleRepository                          usecase.SaleRepository
	taskRepository                          usecase.TaskRepository
	interviewTaskRepository                 usecase.InterviewTaskRepository
	interviewTaskGroupRepository            usecase.InterviewTaskGroupRepository
	jobSeekerRepository                     usecase.JobSeekerRepository
	jobSeekerStudentHistoryRepository       usecase.JobSeekerStudentHistoryRepository
	jobSeekerWorkHistoryRepository          usecase.JobSeekerWorkHistoryRepository
	jobSeekerExperienceIndustryRepository   usecase.JobSeekerExperienceIndustryRepository
	jobSeekerDepartmentHistoryRepository    usecase.JobSeekerDepartmentHistoryRepository
	jobSeekerLicenseRepository              usecase.JobSeekerLicenseRepository
	jobSeekerSelfPromotionRepository        usecase.JobSeekerSelfPromotionRepository
	jobSeekerDocumentRepository             usecase.JobSeekerDocumentRepository
	jobSeekerDesiredIndustryRepository      usecase.JobSeekerDesiredIndustryRepository
	jobSeekerDesiredOccupationRepository    usecase.JobSeekerDesiredOccupationRepository
	jobSeekerDesiredWorkLocationRepository  usecase.JobSeekerDesiredWorkLocationRepository
	jobSeekerDesiredHolidayTypeRepository   usecase.JobSeekerDesiredHolidayTypeRepository
	jobSeekerDevelopmentSkillRepository     usecase.JobSeekerDevelopmentSkillRepository
	jobSeekerLanguageSkillRepository        usecase.JobSeekerLanguageSkillRepository
	jobSeekerPCToolRepository               usecase.JobSeekerPCToolRepository
	jobSeekerHideToAgentRepository          usecase.JobSeekerHideToAgentRepository
	jobSeekerExperienceOccupationRepository usecase.JobSeekerExperienceOccupationRepository
	jobSeekerDesiredCompanyScaleRepository  usecase.JobSeekerDesiredCompanyScaleRepository
}

// DashboardInteractorImpl is an implementation of DashboardInteractor
func NewDashboardInteractorImpl(
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
	jR usecase.JobSeekerRepository,
	jsshR usecase.JobSeekerStudentHistoryRepository,
	jswhR usecase.JobSeekerWorkHistoryRepository,
	jseiR usecase.JobSeekerExperienceIndustryRepository,
	jseoR usecase.JobSeekerExperienceOccupationRepository,
	jslR usecase.JobSeekerLicenseRepository,
	jsspR usecase.JobSeekerSelfPromotionRepository,
	jsdR usecase.JobSeekerDocumentRepository,
	jsdiR usecase.JobSeekerDesiredIndustryRepository,
	jsdoR usecase.JobSeekerDesiredOccupationRepository,
	jsdwlR usecase.JobSeekerDesiredWorkLocationRepository,
	jsdhtR usecase.JobSeekerDesiredHolidayTypeRepository,
	jsdsR usecase.JobSeekerDevelopmentSkillRepository,
	jslsR usecase.JobSeekerLanguageSkillRepository,
	jsptR usecase.JobSeekerPCToolRepository,
	jshR usecase.JobSeekerHideToAgentRepository,
	jsdhR usecase.JobSeekerDepartmentHistoryRepository,
	jsdcsR usecase.JobSeekerDesiredCompanyScaleRepository,
) DashboardInteractor {
	return &DashboardInteractorImpl{
		firebase:                                fb,
		sendgrid:                                sg,
		agentMonthlySaleRepository:              amR,
		agentStaffMonthlySaleRepository:         asmS,
		agentSaleManagementRepository:           asmR,
		staffSaleManagementRepository:           assmR,
		saleRepository:                          sR,
		taskRepository:                          tR,
		interviewTaskRepository:                 itR,
		interviewTaskGroupRepository:            itgR,
		jobSeekerRepository:                     jR,
		jobSeekerStudentHistoryRepository:       jsshR,
		jobSeekerWorkHistoryRepository:          jswhR,
		jobSeekerExperienceIndustryRepository:   jseiR,
		jobSeekerExperienceOccupationRepository: jseoR,
		jobSeekerLicenseRepository:              jslR,
		jobSeekerSelfPromotionRepository:        jsspR,
		jobSeekerDocumentRepository:             jsdR,
		jobSeekerDesiredIndustryRepository:      jsdiR,
		jobSeekerDesiredOccupationRepository:    jsdoR,
		jobSeekerDesiredWorkLocationRepository:  jsdwlR,
		jobSeekerDesiredHolidayTypeRepository:   jsdhtR,
		jobSeekerDevelopmentSkillRepository:     jsdsR,
		jobSeekerLanguageSkillRepository:        jslsR,
		jobSeekerPCToolRepository:               jsptR,
		jobSeekerHideToAgentRepository:          jshR,
		jobSeekerDepartmentHistoryRepository:    jsdhR,
		jobSeekerDesiredCompanyScaleRepository:  jsdcsR,
	}
}

/****************************************************************************************/
/// 汎用系
//

// ダッシュボードのデフォルト表示に必要な情報を取得（当月が含まれる売上期間の受注ベースの個人売上情報を取得する）
type GetDashboardForDefaultPreviewInput struct {
	AgentStaffID uint
}

type GetDashboardForDefaultPreviewOutput struct {
	Dashboard *entity.Dashboard
}

func (i *DashboardInteractorImpl) GetDashboardForDefaultPreview(input GetDashboardForDefaultPreviewInput) (GetDashboardForDefaultPreviewOutput, error) {
	var (
		output              GetDashboardForDefaultPreviewOutput
		err                 error
		staffSaleManagement *entity.AgentStaffSaleManagement
		dashboard           *entity.Dashboard
	)

	/******************   決算期間が当月を含むstaffSaleManagementを全て取得   ******************/

	// 当月を「2006-01」形式で取得
	layout := "2006-01"
	nowTime := time.Now()
	thisMonth := nowTime.Format(layout)

	// 決算期間が当月を含むstaffSaleManagementを全て取得
	staffSaleManagementList, err := i.staffSaleManagementRepository.GetByStaffIDAndThidMonth(input.AgentStaffID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if len(staffSaleManagementList) == 1 {
		// 該当の担当者売上情報が一つだけの場合
		staffSaleManagement = staffSaleManagementList[0]
	} else if len(staffSaleManagementList) > 1 {
		// 該当の担当者売上情報が複数の場合
		for _, ssm := range staffSaleManagementList {
			if ssm.IsOpen {
				staffSaleManagement = ssm
				break // Openになっているのは一つだけのため合致したタイミングでbreakでループを抜ける
			}
		}
	} else {
		// 該当の担当者売上情報が無い場合
		// wrapped := fmt.Errorf("%s:%w", errors.New("あなたの個人売上情報が登録されていません。"), entity.ErrRequestError)
		// return output, wrapped
		return output, nil
	}

	/******************   計算に必要な「担当者の売上情報」と「ヨミ情報」を取得   ******************/

	agentStaffMonthlySale, err := i.agentStaffMonthlySaleRepository.FindByStaffManagementIDAndMonth(staffSaleManagement.ID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	saleList, err := i.saleRepository.GetContractSignedByStaffIDAndMonth(input.AgentStaffID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/******************   計算   ******************/

	var (
		salesPerformance       float64 // 受注実績
		salesResultRate        float64 // 受注売上達成率
		grossProfitPerformance float64 // 受注粗利実績
		grossProfitResultRate  float64 // 受注粗利達成率
		accuracyAccept         float64 // 確定(内定承諾)
		accuracyA              float64 // Aヨミ
		accuracyB              float64 // Bヨミ
		accuracyC              float64 // Cヨミ
		accuracyTopic          float64 // ネタ
		// accuracyList                    []*entity.Sale      // ヨミリスト
		// releaseList                     []*entity.JobSeeker // ヨミリスト

		interviewOfferAcceptancePerformance          float64 // 内定承諾数
		interviewOfferAcceptanceResultRate           float64 // 内定承諾達成率
		interviewOfferPerformance                    float64 // 内定数
		interviewOfferResultRate                     float64 // 内定数達成率
		interviewFinalSelectionPerformance           float64 // 最終選考数
		interviewFinalSelectionResultRate            float64 // 最終選考数達成率
		interviewSelectionPerformance                float64 // 選考数
		interviewSelectionResultRate                 float64 // 選考数達成率
		interviewRecommendationCompletionPerformance float64 // 推薦完了数
		interviewRecommendationCompletionResultRate  float64 // 推薦完了数達成率
		interviewJobIntroductionPerformance          float64 // 求人紹介数
		interviewJobIntroductionResultRate           float64 // 求人紹介数達成率
		interviewPerformance                         float64 // 面談数
		interviewResultRate                          float64 // 面談数達成率
	)

	for _, sale := range saleList {
		// 受注関連の計算
		if agentStaffMonthlySale.SalesMonth == sale.ContractSignedMonth {
			var (
				caSale    float64 = 0 // CAの売上
				raSale    float64 = 0 // RAの売上
				staffSale float64 = 0 // スタッフの売上

				caGrossProfit    float64 = 0 // CAの粗利
				raGrossProfit    float64 = 0 // RAの粗利
				staffGrossProfit float64 = 0 // 担当者の粗利
			)

			// CA売上比率に応じた売上を計算
			if input.AgentStaffID == sale.CAStaffID && sale.CaSalesRatio.Float64 > 0 {
				caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
				caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
			}

			// RA売上比率に応じた売上を計算
			if input.AgentStaffID == sale.RAStaffID && sale.RaSalesRatio.Float64 > 0 {
				raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
				raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
			}

			// エージェントの売上
			staffSale = caSale + raSale
			staffGrossProfit = caGrossProfit + raGrossProfit

			// 内定承諾関連の計算
			if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
				salesPerformance = salesPerformance + staffSale                    // 受注売上実績の計算
				grossProfitPerformance = grossProfitPerformance + staffGrossProfit // 受注粗利実績の計算
				accuracyAccept = accuracyAccept + staffSale                        // 確定(内定承諾)の請求金額合計
			}

			// Aヨミの計算
			if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
				accuracyA = accuracyA + staffSale // Aヨミの請求金額合計
			}

			// Bヨミの計算
			if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
				accuracyB = accuracyB + staffSale // Bヨミの請求金額合計
			}

			// Cヨミの計算
			if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
				accuracyC = accuracyC + staffSale // Cヨミの請求金額合計
			}

			// ネタの計算
			if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
				accuracyTopic = accuracyTopic + staffSale // ネタの請求金額合計
			}
		}
	}

	/********* 受注売上 **********/

	// 受注売上達成率（売上実績 / 売上予算）
	if agentStaffMonthlySale.OrderSalesBudget.Float64 > 0 {
		salesResultRate = math.Floor(salesPerformance / agentStaffMonthlySale.OrderSalesBudget.Float64 * 100)
	}

	// 粗利達成率（粗利実績 / 粗利予算）
	if agentStaffMonthlySale.OrderGrossProfitBudget.Float64 > 0 {
		grossProfitResultRate = math.Floor(grossProfitPerformance / agentStaffMonthlySale.OrderGrossProfitBudget.Float64 * 100)
	}

	/********* 求職者数（面談実施月カウントベース） **********/

	// 内定承諾数
	interviewOfferAcceptancePerformance, err = i.interviewTaskGroupRepository.GetInterviewOfferAcceptancePerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 内定承諾達成率（内定承諾数 / 内定承諾目標 * 100）
	if agentStaffMonthlySale.InterviewOfferAcceptanceTarget.Float64 > 0 {
		interviewOfferAcceptanceResultRate = math.Floor(interviewOfferAcceptancePerformance / agentStaffMonthlySale.InterviewOfferAcceptanceTarget.Float64 * 100)
	}

	// 内定数
	interviewOfferPerformance, err = i.taskRepository.GetInterviewOfferPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 内定数達成率（内定数 / 内定目標 * 100）
	if agentStaffMonthlySale.InterviewOfferTarget.Float64 > 0 {
		interviewOfferResultRate = math.Floor(interviewOfferPerformance / agentStaffMonthlySale.InterviewOfferTarget.Float64 * 100)
	}

	// 最終選考数
	interviewFinalSelectionPerformance, err = i.taskRepository.GetInterviewFinalSelectionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 最終選考数達成率（最終選考数 / 最終選考目標 * 100）
	if agentStaffMonthlySale.InterviewFinalSelectionTarget.Float64 > 0 {
		interviewFinalSelectionResultRate = math.Floor(interviewFinalSelectionPerformance / agentStaffMonthlySale.InterviewFinalSelectionTarget.Float64 * 100)
	}

	// 選考数
	interviewSelectionPerformance, err = i.taskRepository.GetInterviewSelectionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 選考数達成率（選考数 / 選考目標 * 100）
	if agentStaffMonthlySale.InterviewSelectionTarget.Float64 > 0 {
		interviewSelectionResultRate = math.Floor(interviewSelectionPerformance / agentStaffMonthlySale.InterviewSelectionTarget.Float64 * 100)
	}

	// 推薦完了数
	interviewRecommendationCompletionPerformance, err = i.taskRepository.GetInterviewRecommendationCompletionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 推薦完了数達成率（推薦完了数 / 推薦完了数目標 * 100）
	if agentStaffMonthlySale.InterviewRecommendationCompletionTarget.Float64 > 0 {
		interviewRecommendationCompletionResultRate = math.Floor(interviewRecommendationCompletionPerformance / agentStaffMonthlySale.InterviewRecommendationCompletionTarget.Float64 * 100)
	}

	// 求人紹介数
	interviewJobIntroductionPerformance, err = i.taskRepository.GetInterviewJobIntroductionPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求人紹介数達成率（求人紹介数 / 求人紹介数目標 * 100）
	if agentStaffMonthlySale.InterviewJobIntroductionTarget.Float64 > 0 {
		interviewJobIntroductionResultRate = math.Floor(interviewJobIntroductionPerformance / agentStaffMonthlySale.InterviewJobIntroductionTarget.Float64 * 100)
	}

	// 面談数
	interviewPerformance, err = i.interviewTaskGroupRepository.GetInterviewPerformanceCountByStaffIDAndSalesMonth(input.AgentStaffID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 面談数達成率（面談数 / 面談数目標 * 100）
	if agentStaffMonthlySale.InterviewInterviewTarget.Float64 > 0 {
		interviewResultRate = math.Floor(interviewPerformance / agentStaffMonthlySale.InterviewInterviewTarget.Float64 * 100)
	}

	dashboard = entity.NewDashboard(
		staffSaleManagement.ManagementID,
		staffSaleManagement.FiscalYear,
		salesPerformance,
		agentStaffMonthlySale.OrderSalesBudget.Float64,
		salesResultRate,
		grossProfitPerformance,
		agentStaffMonthlySale.OrderGrossProfitBudget.Float64,
		grossProfitResultRate,
		agentStaffMonthlySale.InterviewOfferAcceptanceTarget.Float64,
		interviewOfferAcceptancePerformance,
		interviewOfferAcceptanceResultRate,
		agentStaffMonthlySale.InterviewOfferTarget.Float64,
		interviewOfferPerformance,
		interviewOfferResultRate,
		agentStaffMonthlySale.InterviewFinalSelectionTarget.Float64,
		interviewFinalSelectionPerformance,
		interviewFinalSelectionResultRate,
		agentStaffMonthlySale.InterviewSelectionTarget.Float64,
		interviewSelectionPerformance,
		interviewSelectionResultRate,
		agentStaffMonthlySale.InterviewRecommendationCompletionTarget.Float64,
		interviewRecommendationCompletionPerformance,
		interviewRecommendationCompletionResultRate,
		agentStaffMonthlySale.InterviewJobIntroductionTarget.Float64,
		interviewJobIntroductionPerformance,
		interviewJobIntroductionResultRate,
		agentStaffMonthlySale.InterviewInterviewTarget.Float64,
		interviewPerformance,
		interviewResultRate,
		accuracyAccept,
		accuracyA,
		accuracyB,
		accuracyC,
		accuracyTopic,
	)

	output.Dashboard = dashboard

	return output, nil
}

type GetSaleListForDefaultPreviewInput struct {
	AgentStaffID uint
	PageNumber   uint
}

type GetSaleListForDefaultPreviewOutput struct {
	SaleLlist     []*entity.Sale
	MaxPageNumber uint
	IDList        []uint
}

func (i *DashboardInteractorImpl) GetSaleListForDefaultPreview(input GetSaleListForDefaultPreviewInput) (GetSaleListForDefaultPreviewOutput, error) {
	var (
		output GetSaleListForDefaultPreviewOutput
		err    error
	)

	/******************   決算期間が当月を含むstaffSaleManagementを全て取得   ******************/

	// 当月を「2006-01」形式で取得
	layout := "2006-01"
	nowTime := time.Now()
	thisMonth := nowTime.Format(layout)

	/******************   「ヨミ情報」を取得（デフォルトは受注）   ******************/

	saleList, err := i.saleRepository.GetContractSignedByStaffIDAndMonth(input.AgentStaffID, thisMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sale := range saleList {
		output.IDList = append(output.IDList, sale.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getSaleListMaxPage(saleList)

	// 指定ページのヨミ情報50件を取得
	output.SaleLlist = getSaleListWithPage(saleList, input.PageNumber)

	return output, nil
}

type GetReleaseJobSeekerListForDefaultPreviewInput struct {
	AgentStaffID uint
	PageNumber   uint
}

type GetReleaseJobSeekerListForDefaultPreviewOutput struct {
	JobSeekerList []*entity.JobSeeker
	MaxPageNumber uint
	IDList        []uint
}

func (i *DashboardInteractorImpl) GetReleaseJobSeekerListForDefaultPreview(input GetReleaseJobSeekerListForDefaultPreviewInput) (GetReleaseJobSeekerListForDefaultPreviewOutput, error) {
	var (
		output GetReleaseJobSeekerListForDefaultPreviewOutput
		err    error
	)

	/******************   「リリース求職者情報」を取得   ******************/

	jobSeekerList, err := i.jobSeekerRepository.GetReleaseByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者IDリストから関連テーブル情報を取得
	// selfpromotionとdocumentは使わないので取得しない
	studentHistory, err := i.jobSeekerStudentHistoryRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workHistory, err := i.jobSeekerWorkHistoryRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceIndustry, err := i.jobSeekerExperienceIndustryRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	departmentHistory, err := i.jobSeekerDepartmentHistoryRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceOccupation, err := i.jobSeekerExperienceOccupationRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredCompanyScale, err := i.jobSeekerDesiredCompanyScaleRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	license, err := i.jobSeekerLicenseRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredIndustry, err := i.jobSeekerDesiredIndustryRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredOccupation, err := i.jobSeekerDesiredOccupationRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredWorkLocation, err := i.jobSeekerDesiredWorkLocationRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredHolidayType, err := i.jobSeekerDesiredHolidayTypeRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	developmentSkill, err := i.jobSeekerDevelopmentSkillRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	languageSkill, err := i.jobSeekerLanguageSkillRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	pcTool, err := i.jobSeekerPCToolRepository.GetByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, jobSeeker := range jobSeekerList {
		for _, sh := range studentHistory {
			if jobSeeker.ID == sh.JobSeekerID {
				jobSeeker.StudentHistories = append(jobSeeker.StudentHistories, *sh)
			}
		}

		for _, wh := range workHistory {
			if jobSeeker.ID == wh.JobSeekerID {
				for _, ei := range experienceIndustry {
					if wh.ID == ei.WorkHistoryID {
						wh.ExperienceIndustries = append(wh.ExperienceIndustries, *ei)
					}
				}

				for _, dh := range departmentHistory {
					if wh.ID == dh.WorkHistoryID {
						for _, eo := range experienceOccupation {
							if dh.ID == eo.DepartmentHistoryID {
								dh.ExperienceOccupations = append(dh.ExperienceOccupations, *eo)
							}
						}
						wh.DepartmentHistories = append(wh.DepartmentHistories, *dh)
					}
				}
				jobSeeker.WorkHistories = append(jobSeeker.WorkHistories, *wh)
			}
		}

		for _, dcs := range desiredCompanyScale {
			if jobSeeker.ID == dcs.JobSeekerID {
				jobSeeker.DesiredCompanyScales = append(jobSeeker.DesiredCompanyScales, *dcs)
			}
		}

		for _, l := range license {
			if jobSeeker.ID == l.JobSeekerID {
				jobSeeker.Licenses = append(jobSeeker.Licenses, *l)
			}
		}

		for _, di := range desiredIndustry {
			if jobSeeker.ID == di.JobSeekerID {
				jobSeeker.DesiredIndustries = append(jobSeeker.DesiredIndustries, *di)
			}
		}

		for _, do := range desiredOccupation {
			if jobSeeker.ID == do.JobSeekerID {
				jobSeeker.DesiredOccupations = append(jobSeeker.DesiredOccupations, *do)
			}
		}

		for _, dwl := range desiredWorkLocation {
			if jobSeeker.ID == dwl.JobSeekerID {
				jobSeeker.DesiredWorkLocations = append(jobSeeker.DesiredWorkLocations, *dwl)
			}
		}

		for _, dht := range desiredHolidayType {
			if jobSeeker.ID == dht.JobSeekerID {
				jobSeeker.DesiredHolidayTypes = append(jobSeeker.DesiredHolidayTypes, *dht)
			}
		}

		for _, ds := range developmentSkill {
			if jobSeeker.ID == ds.JobSeekerID {
				jobSeeker.DevelopmentSkills = append(jobSeeker.DevelopmentSkills, *ds)
			}
		}

		for _, ls := range languageSkill {
			if jobSeeker.ID == ls.JobSeekerID {
				jobSeeker.LanguageSkills = append(jobSeeker.LanguageSkills, *ls)
			}
		}

		for _, ps := range pcTool {
			if jobSeeker.ID == ps.JobSeekerID {
				jobSeeker.PCTools = append(jobSeeker.PCTools, *ps)
			}
		}
	}

	for _, jobSeeker := range jobSeekerList {
		output.IDList = append(output.IDList, jobSeeker.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getJobSeekerListMaxPage(jobSeekerList)

	// 指定ページの求職者情報50件を取得
	output.JobSeekerList = getJobSeekerListWithPage(jobSeekerList, input.PageNumber)

	return output, nil
}

/****************************************************************************************/
// 検索
//

type GetDashboardForSearchInput struct {
	SearchParam entity.SearchDashboard
}

type GetDashboardForSearchOutput struct {
	Dashboard *entity.Dashboard
}

func (i *DashboardInteractorImpl) GetDashboardForSearch(input GetDashboardForSearchInput) (GetDashboardForSearchOutput, error) {
	var (
		output    GetDashboardForSearchOutput
		err       error
		dashboard *entity.Dashboard
	)
	/******************   絞り込みの条件   ******************/
	// input.SearchParam.ManagementID
	// input.SearchParam.Period
	// input.SearchParam.Target      // 0: 受注, 1: 請求
	// input.SearchParam.SearchRange //

	// 決算期間が当月を含むstaffSaleManagementを全て取得
	agentSaleManagement, err := i.agentSaleManagementRepository.FindByID(input.SearchParam.ManagementID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 決算期間と期限（period）を元に["2023-04", "2024-04"]のような形式で取得
	period, err := changePeriodValue(agentSaleManagement.FiscalYear, input.SearchParam.Period)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/******************   SearchParamの値に応じてデータを取得する   ******************/

	if input.SearchParam.SearchRange == 0 {
		/*** 会社全体の場合 ***/

		if len(period) == 1 {
			// 指定エージェントの「〇〇年-〇月」のダッシュボード情報を取得
			dashboard, err = searchAgentSaleForTheMonth(
				i, input.SearchParam, agentSaleManagement.FiscalYear, period[0],
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else if len(period) == 2 {
			// 指定エージェントの「〇〇年-〇月 〜 〇〇年-〇月」のダッシュボード情報を取得
			dashboard, err = searchAgentSaleFromStartMonthToEndMonth(
				i, input.SearchParam, agentSaleManagement.FiscalYear, period[0], period[1],
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	} else {
		/*** 担当者の場合 ***/

		if len(period) == 1 {
			// 指定担当者の「〇〇年-〇月」のダッシュボード情報を取得
			dashboard, err = searchStaffSaleForTheMonth(
				i, input.SearchParam, agentSaleManagement.FiscalYear, period[0],
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else if len(period) == 2 {
			// 指定担当者の「〇〇年-〇月 〜 〇〇年-〇月」のダッシュボード情報を取得
			dashboard, err = searchStaffSaleFromStartMonthToEndMonth(
				i, input.SearchParam, agentSaleManagement.FiscalYear, period[0], period[1],
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.Dashboard = dashboard

	return output, nil
}

type GetSaleListForSearchInput struct {
	PageNumber  uint
	SearchParam entity.SearchDashboard
}

type GetSaleListForSearchOutput struct {
	SaleLlist     []*entity.Sale
	MaxPageNumber uint
	IDList        []uint
}

func (i *DashboardInteractorImpl) GetSaleListForSearch(input GetSaleListForSearchInput) (GetSaleListForSearchOutput, error) {
	var (
		output            GetSaleListForSearchOutput
		err               error
		saleList          []*entity.Sale
		filteringSaleList []*entity.Sale
	)

	/******************   基本情報を取得   ******************/

	// 決算期間が当月を含むstaffSaleManagementを全て取得
	agentSaleManagement, err := i.agentSaleManagementRepository.FindByID(input.SearchParam.ManagementID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 決算期間と期限（period）を元に["2023-04", "2024-04"]のような形式で取得
	period, err := changePeriodValue(agentSaleManagement.FiscalYear, input.SearchParam.Period)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/******************   「ヨミ情報」を取得   ******************/

	if input.SearchParam.SearchRange == 0 {
		/*** 会社全体の場合 ***/

		if len(period) == 1 {
			// 指定エージェントの「〇〇年-〇月」のヨミ情報を取得
			if input.SearchParam.Target == 0 {
				// 受注ベース
				saleList, err = i.saleRepository.GetContractSignedByAgentIDAndMonth(input.SearchParam.AgentID, period[0])
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			} else if input.SearchParam.Target == 1 {
				// 請求ベース
				saleList, err = i.saleRepository.GetBillingByAgentIDAndMonth(input.SearchParam.AgentID, period[0])
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		} else if len(period) == 2 {
			// 指定エージェントの「〇〇年-〇月 〜 〇〇年-〇月」のヨミ情報を取得
			if input.SearchParam.Target == 0 {
				// 受注ベース
				saleList, err = i.saleRepository.GetContractSignedByAgentIDAndPeriod(input.SearchParam.AgentID, period[0], period[1])
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			} else if input.SearchParam.Target == 1 {
				// 請求ベース
				saleList, err = i.saleRepository.GetBillingByAgentIDAndPeriod(input.SearchParam.AgentID, period[0], period[1])
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}
	} else {
		/*** 担当者の場合 ***/

		if len(period) == 1 {
			// 指定担当者の「〇〇年-〇月」のヨミ情報を取得

			if input.SearchParam.Target == 0 {
				// 受注ベース
				saleList, err = i.saleRepository.GetContractSignedByStaffIDAndMonth(input.SearchParam.SearchRange, period[0])
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			} else if input.SearchParam.Target == 1 {
				// 請求ベース
				saleList, err = i.saleRepository.GetBillingByStaffIDAndMonth(input.SearchParam.SearchRange, period[0])
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		} else if len(period) == 2 {
			// 指定担当者の「〇〇年-〇月 〜 〇〇年-〇月」のヨミ情報を取得

			if input.SearchParam.Target == 0 {
				// 受注ベース
				saleList, err = i.saleRepository.GetContractSignedByStaffIDAndPeriod(input.SearchParam.SearchRange, period[0], period[1])
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			} else if input.SearchParam.Target == 1 {
				// 請求ベース
				saleList, err = i.saleRepository.GetBillingByStaffIDAndPeriod(input.SearchParam.SearchRange, period[0], period[1])
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}
	}

	// 検索対象のヨミフェーズのsaleのみ残す
	for _, sale := range saleList {
		for _, accuracy := range input.SearchParam.AccuracyTypes {
			if sale.Accuracy == accuracy {
				output.IDList = append(output.IDList, sale.ID)
				filteringSaleList = append(filteringSaleList, sale)
				break
			}
		}
	}

	// ページの最大数を取得
	output.MaxPageNumber = getSaleListMaxPage(filteringSaleList)

	// 指定ページのヨミ情報50件を取得
	output.SaleLlist = getSaleListWithPage(filteringSaleList, input.PageNumber)

	return output, nil
}

type GetReleaseJobSeekerListForSearchInput struct {
	PageNumber  uint
	SearchParam entity.SearchDashboard
}

type GetReleaseJobSeekerListForSearchOutput struct {
	JobSeekerList []*entity.JobSeeker
	MaxPageNumber uint
	IDList        []uint
}

func (i *DashboardInteractorImpl) GetReleaseJobSeekerListForSearch(input GetReleaseJobSeekerListForSearchInput) (GetReleaseJobSeekerListForSearchOutput, error) {
	var (
		output        GetReleaseJobSeekerListForSearchOutput
		err           error
		jobSeekerList []*entity.JobSeeker
	)

	if input.SearchParam.SearchRange == 0 {
		// input.SearchParam.SearchRangeが0の場合はエージェントIDを使用する
		agentID := input.SearchParam.AgentID

		// 会社全体
		jobSeekerList, err = i.jobSeekerRepository.GetReleaseByAgentID(agentID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 子テーブル情報を取得
		jobSeekerList, err = getJobSeekerInformationsByJobSeekerList(i, jobSeekerList)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

	} else {
		// 担当者
		// input.SearchParam.SearchRangeが0でない場合は担当者のIDが入ってる
		staffID := input.SearchParam.SearchRange

		jobSeekerList, err = i.jobSeekerRepository.GetReleaseByStaffID(staffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 子テーブル情報を取得
		jobSeekerList, err = getJobSeekerInformationsByStaffID(i, jobSeekerList, staffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, jobSeeker := range jobSeekerList {
		fmt.Println(jobSeeker.LastName, jobSeeker.FirstName)

		output.IDList = append(output.IDList, jobSeeker.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getJobSeekerListMaxPage(jobSeekerList)

	// 指定ページの求職者情報50件を取得
	output.JobSeekerList = getJobSeekerListWithPage(jobSeekerList, input.PageNumber)

	return output, nil
}

/****************************************************************************************/
// 検索関数
//

// {range: 会社全体, period: 1~11}
func searchAgentSaleForTheMonth(
	i *DashboardInteractorImpl,
	param entity.SearchDashboard, // 検索パラム
	fiscalYear string, // 決算月
	month string, // 検索対象の月
) (*entity.Dashboard, error) {
	var (
		dashboard *entity.Dashboard
		saleList  []*entity.Sale
	)

	/******************   計算に必要な「エージェントの売上情報」と「ヨミ情報」を取得   ******************/

	agentMonthlySale, err := i.agentMonthlySaleRepository.FindByManagementIDAndMonth(param.ManagementID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	if param.Target == 0 {
		// 受注ベース
		saleList, err = i.saleRepository.GetContractSignedByAgentIDAndMonth(param.AgentID, month)
		if err != nil {
			fmt.Println(err)
			return dashboard, err
		}
	} else {
		// 請求ベース
		saleList, err = i.saleRepository.GetBillingByAgentIDAndMonth(param.AgentID, month)
		if err != nil {
			fmt.Println(err)
			return dashboard, err
		}
	}

	/******************   計算   ******************/

	var (
		salesPerformance                             float64 // 売上実績
		salesBudget                                  float64 // 売上予算
		salesResultRate                              float64 // 売上達成率
		grossProfitPerformance                       float64 // 粗利実績
		grossProfitBudget                            float64 // 粗利予算
		grossProfitResultRate                        float64 // 粗利達成率
		accuracyAccept                               float64 // 確定(内定承諾)
		accuracyA                                    float64 // Aヨミ
		accuracyB                                    float64 // Bヨミ
		accuracyC                                    float64 // Cヨミ
		accuracyTopic                                float64 // ネタ
		interviewOfferAcceptancePerformance          float64 // 内定承諾数
		interviewOfferAcceptanceResultRate           float64 // 内定承諾達成率
		interviewOfferPerformance                    float64 // 内定数
		interviewOfferResultRate                     float64 // 内定数達成率
		interviewFinalSelectionPerformance           float64 // 最終選考数
		interviewFinalSelectionResultRate            float64 // 最終選考数達成率
		interviewSelectionPerformance                float64 // 選考数
		interviewSelectionResultRate                 float64 // 選考数達成率
		interviewRecommendationCompletionPerformance float64 // 推薦完了数
		interviewRecommendationCompletionResultRate  float64 // 推薦完了数達成率
		interviewJobIntroductionPerformance          float64 // 求人紹介数
		interviewJobIntroductionResultRate           float64 // 求人紹介数達成率
		interviewPerformance                         float64 // 面談数
		interviewResultRate                          float64 // 面談数達成率
	)

	for _, sale := range saleList {
		if param.Target == 0 {
			// 受注関連の計算

			if agentMonthlySale.SalesMonth == sale.ContractSignedMonth {
				var (
					caSale    float64 = 0 // CAの売上
					raSale    float64 = 0 // RAの売上
					agentSale float64 = 0 // エージェントの売上

					caGrossProfit    float64 = 0 // CAの粗利
					raGrossProfit    float64 = 0 // RAの粗利
					agentGrossProfit float64 = 0 // エージェントの粗利
				)

				// CA売上比率に応じた売上を計算
				if param.AgentID == sale.CAAgentID && sale.CaSalesRatio.Float64 > 0 {
					caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
					caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
				}

				// RA売上比率に応じた売上を計算
				if param.AgentID == sale.RAAgentID && sale.RaSalesRatio.Float64 > 0 {
					raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
					raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
				}

				// エージェントの売上
				agentSale = caSale + raSale
				agentGrossProfit = caGrossProfit + raGrossProfit

				// 内定承諾の計算
				if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
					salesPerformance = salesPerformance + agentSale                    // 受注売上実績の計算
					grossProfitPerformance = grossProfitPerformance + agentGrossProfit // 受注粗利実績の計算
					accuracyAccept = accuracyAccept + agentSale                        // Aヨミの請求金額合計
				}

				// Aヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
					accuracyA = accuracyA + agentSale // Aヨミの請求金額合計
				}

				// Bヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
					accuracyB = accuracyB + agentSale // Bヨミの請求金額合計
				}

				// Cヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
					accuracyC = accuracyC + agentSale // Cヨミの請求金額合計
				}

				// ネタの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
					accuracyTopic = accuracyTopic + agentSale // ネタの請求金額合計
				}
			}
		} else {
			// 請求関連の計算
			if agentMonthlySale.SalesMonth == sale.BillingMonth {
				var (
					caSale    float64 = 0 // CAの売上
					raSale    float64 = 0 // RAの売上
					agentSale float64 = 0 // エージェントの売上

					caGrossProfit    float64 = 0 // CAの粗利
					raGrossProfit    float64 = 0 // RAの粗利
					agentGrossProfit float64 = 0 // エージェントの粗利
				)

				// CA売上比率に応じた売上を計算
				if param.AgentID == sale.CAAgentID && sale.CaSalesRatio.Float64 > 0 {
					caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
					caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
				}

				// RA売上比率に応じた売上を計算
				if param.AgentID == sale.RAAgentID && sale.RaSalesRatio.Float64 > 0 {
					raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
					raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
				}

				// エージェントの売上
				agentSale = caSale + raSale
				agentGrossProfit = caGrossProfit + raGrossProfit

				// 確定(内定承諾)の計算
				if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
					salesPerformance = salesPerformance + agentSale                    // 売上実績の計算
					grossProfitPerformance = grossProfitPerformance + agentGrossProfit // 粗利実績の計算
					accuracyAccept = accuracyAccept + agentSale                        // 確定(内定承諾)の請求金額合計
				}

				// Aヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
					accuracyA = accuracyA + agentSale // Aヨミの請求金額合計
				}

				// Bヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
					accuracyB = accuracyB + agentSale // Bヨミの請求金額合計
				}

				// Cヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
					accuracyC = accuracyC + agentSale // Cヨミの請求金額合計
				}

				// ネタの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
					accuracyTopic = accuracyTopic + agentSale // ネタの請求金額合計
				}
			}
		}
	}

	/********* 受注売上 **********/

	if param.Target == 0 {
		// 受注予算
		salesBudget = agentMonthlySale.OrderSalesBudget.Float64
		grossProfitBudget = agentMonthlySale.OrderGrossProfitBudget.Float64

		// 受注売上達成率（売上実績 / 売上予算）
		if salesBudget > 0 {
			salesResultRate = math.Floor(salesPerformance / salesBudget * 100)
		}

		// 粗利達成率（粗利実績 / 粗利予算）
		if grossProfitBudget > 0 {
			grossProfitResultRate = math.Floor(grossProfitPerformance / grossProfitBudget * 100)
		}
	} else {
		// 請求予算
		salesBudget = agentMonthlySale.ClaimSalesRevenueBudget.Float64
		grossProfitBudget = agentMonthlySale.ClaimGrossMarginBudget.Float64

		// 請求売上達成率（売上実績 / 売上予算）
		if salesBudget > 0 {
			salesResultRate = math.Floor(salesPerformance / salesBudget * 100)
		}

		// 粗利達成率（粗利実績 / 粗利予算）
		if grossProfitBudget > 0 {
			grossProfitResultRate = math.Floor(grossProfitPerformance / grossProfitBudget * 100)
		}
	}

	/********* 求職者数（面談実施月カウントベース） **********/

	// 内定承諾数
	interviewOfferAcceptancePerformance, err = i.interviewTaskGroupRepository.GetInterviewOfferAcceptancePerformanceCountByAgentIDAndSalesMonth(param.AgentID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 内定承諾達成率（内定承諾数 / 内定承諾目標 * 100）
	if agentMonthlySale.InterviewOfferAcceptanceTarget.Float64 > 0 {
		interviewOfferAcceptanceResultRate = math.Floor(interviewOfferAcceptancePerformance / agentMonthlySale.InterviewOfferAcceptanceTarget.Float64 * 100)
	}

	// 内定数
	interviewOfferPerformance, err = i.taskRepository.GetInterviewOfferPerformanceCountByAgentIDAndSalesMonth(param.AgentID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 内定数達成率（内定数 / 内定目標 * 100）
	if agentMonthlySale.InterviewOfferTarget.Float64 > 0 {
		interviewOfferResultRate = math.Floor(interviewOfferPerformance / agentMonthlySale.InterviewOfferTarget.Float64 * 100)
	}

	// 最終選考数
	interviewFinalSelectionPerformance, err = i.taskRepository.GetInterviewFinalSelectionPerformanceCountByAgentIDAndSalesMonth(param.AgentID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 最終選考数達成率（最終選考数 / 最終選考目標 * 100）
	if agentMonthlySale.InterviewFinalSelectionTarget.Float64 > 0 {
		interviewFinalSelectionResultRate = math.Floor(interviewFinalSelectionPerformance / agentMonthlySale.InterviewFinalSelectionTarget.Float64 * 100)
	}

	// 選考数
	interviewSelectionPerformance, err = i.taskRepository.GetInterviewSelectionPerformanceCountByAgentIDAndSalesMonth(param.AgentID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 選考数達成率（選考数 / 選考目標 * 100）
	if agentMonthlySale.InterviewSelectionTarget.Float64 > 0 {
		interviewSelectionResultRate = math.Floor(interviewSelectionPerformance / agentMonthlySale.InterviewSelectionTarget.Float64 * 100)
	}

	// 推薦完了数
	interviewRecommendationCompletionPerformance, err = i.taskRepository.GetInterviewRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth(param.AgentID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 推薦完了数達成率（推薦完了数 / 推薦完了数目標 * 100）
	if agentMonthlySale.InterviewRecommendationCompletionTarget.Float64 > 0 {
		interviewRecommendationCompletionResultRate = math.Floor(interviewRecommendationCompletionPerformance / agentMonthlySale.InterviewRecommendationCompletionTarget.Float64 * 100)
	}

	// 求人紹介数
	interviewJobIntroductionPerformance, err = i.taskRepository.GetInterviewJobIntroductionPerformanceCountByAgentIDAndSalesMonth(param.AgentID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 求人紹介数達成率（求人紹介数 / 求人紹介数目標 * 100）
	if agentMonthlySale.InterviewJobIntroductionTarget.Float64 > 0 {
		interviewJobIntroductionResultRate = math.Floor(interviewJobIntroductionPerformance / agentMonthlySale.InterviewJobIntroductionTarget.Float64 * 100)
	}

	// 面談数
	interviewPerformance, err = i.interviewTaskGroupRepository.GetInterviewPerformanceCountByAgentIDAndSalesMonth(param.AgentID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 面談数達成率（面談数 / 面談数目標 * 100）
	if agentMonthlySale.InterviewInterviewTarget.Float64 > 0 {
		interviewResultRate = math.Floor(interviewPerformance / agentMonthlySale.InterviewInterviewTarget.Float64 * 100)
	}

	dashboard = entity.NewDashboard(
		param.ManagementID,
		fiscalYear,
		salesPerformance,
		salesBudget,
		salesResultRate,
		grossProfitPerformance,
		grossProfitBudget,
		grossProfitResultRate,
		agentMonthlySale.InterviewOfferAcceptanceTarget.Float64,
		interviewOfferAcceptancePerformance,
		interviewOfferAcceptanceResultRate,
		agentMonthlySale.InterviewOfferTarget.Float64,
		interviewOfferPerformance,
		interviewOfferResultRate,
		agentMonthlySale.InterviewFinalSelectionTarget.Float64,
		interviewFinalSelectionPerformance,
		interviewFinalSelectionResultRate,
		agentMonthlySale.InterviewSelectionTarget.Float64,
		interviewSelectionPerformance,
		interviewSelectionResultRate,
		agentMonthlySale.InterviewRecommendationCompletionTarget.Float64,
		interviewRecommendationCompletionPerformance,
		interviewRecommendationCompletionResultRate,
		agentMonthlySale.InterviewJobIntroductionTarget.Float64,
		interviewJobIntroductionPerformance,
		interviewJobIntroductionResultRate,
		agentMonthlySale.InterviewInterviewTarget.Float64,
		interviewPerformance,
		interviewResultRate,
		accuracyAccept,
		accuracyA,
		accuracyB,
		accuracyC,
		accuracyTopic,
	)

	return dashboard, nil
}

// {range: 会社全体, period: 12~16}
func searchAgentSaleFromStartMonthToEndMonth(
	i *DashboardInteractorImpl,
	param entity.SearchDashboard, // 検索パラム
	fiscalYear string, // 決算月
	startMonth string, // 検索対象の開始月
	endMonth string, // 検索対象の終了月
) (*entity.Dashboard, error) {
	var (
		dashboard *entity.Dashboard
		saleList  []*entity.Sale
	)

	/******************   計算に必要な「エージェントの売上情報」と「ヨミ情報」を取得   ******************/

	agentMonthlySaleList, err := i.agentMonthlySaleRepository.GetByManagementIDAndPeriod(param.ManagementID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	if param.Target == 0 {
		// 受注ベース
		saleList, err = i.saleRepository.GetContractSignedByAgentIDAndPeriod(param.AgentID, startMonth, endMonth)
		if err != nil {
			fmt.Println(err)
			return dashboard, err
		}
	} else {
		// 請求ベース
		saleList, err = i.saleRepository.GetBillingByAgentIDAndPeriod(param.AgentID, startMonth, endMonth)
		if err != nil {
			fmt.Println(err)
			return dashboard, err
		}
	}

	/******************   計算   ******************/

	var (
		salesPerformance                             float64 // 売上実績
		salesBudget                                  float64 // 売上予算
		salesResultRate                              float64 // 売上達成率
		grossProfitPerformance                       float64 // 粗利実績
		grossProfitBudget                            float64 // 粗利予算
		grossProfitResultRate                        float64 // 粗利達成率
		accuracyAccept                               float64 // 確定(内定承諾)
		accuracyA                                    float64 // Aヨミ
		accuracyB                                    float64 // Bヨミ
		accuracyC                                    float64 // Cヨミ
		accuracyTopic                                float64 // ネタ
		interviewOfferAcceptancePerformance          float64 // 内定承諾数
		interviewOfferAcceptanceTarget               float64 // 内定承諾目標
		interviewOfferAcceptanceResultRate           float64 // 内定承諾達成率
		interviewOfferPerformance                    float64 // 内定数
		interviewOfferTarget                         float64 // 内定数目標
		interviewOfferResultRate                     float64 // 内定数達成率
		interviewFinalSelectionPerformance           float64 // 最終選考数
		interviewFinalSelectionTarget                float64 // 最終選考数目標
		interviewFinalSelectionResultRate            float64 // 最終選考数達成率
		interviewSelectionPerformance                float64 // 選考数
		interviewSelectionTarget                     float64 // 選考数目標
		interviewSelectionResultRate                 float64 // 選考数達成率
		interviewRecommendationCompletionPerformance float64 // 推薦完了数
		interviewRecommendationCompletionTarget      float64 // 推薦完了目標
		interviewRecommendationCompletionResultRate  float64 // 推薦完了数達成率
		interviewJobIntroductionPerformance          float64 // 求人紹介数
		interviewJobIntroductionTarget               float64 // 求人紹介数目標
		interviewJobIntroductionResultRate           float64 // 求人紹介数達成率
		interviewPerformance                         float64 // 面談数
		interviewTarget                              float64 // 面談数目標
		interviewResultRate                          float64 // 面談数達成率
	)

	for _, monthlySale := range agentMonthlySaleList {
		// 予算の合計値を計算
		if param.Target == 0 {
			// 受注関連の計算
			salesBudget = salesBudget + monthlySale.OrderSalesBudget.Float64                   // 売上
			grossProfitBudget = grossProfitBudget + monthlySale.OrderGrossProfitBudget.Float64 // 粗利
		} else {
			// 請求関連の計算
			salesBudget = salesBudget + monthlySale.ClaimSalesRevenueBudget.Float64            // 売上
			grossProfitBudget = grossProfitBudget + monthlySale.ClaimGrossMarginBudget.Float64 // 粗利
		}

		// 面談実施月ベースの目標の合計値を計算
		interviewOfferAcceptanceTarget = interviewOfferAcceptanceTarget + monthlySale.InterviewOfferAcceptanceTarget.Float64
		interviewOfferTarget = interviewOfferTarget + monthlySale.InterviewOfferTarget.Float64
		interviewFinalSelectionTarget = interviewFinalSelectionTarget + monthlySale.InterviewFinalSelectionTarget.Float64
		interviewSelectionTarget = interviewSelectionTarget + monthlySale.InterviewSelectionTarget.Float64
		interviewRecommendationCompletionTarget = interviewRecommendationCompletionTarget + monthlySale.InterviewRecommendationCompletionTarget.Float64
		interviewJobIntroductionTarget = interviewJobIntroductionTarget + monthlySale.InterviewJobIntroductionTarget.Float64
		interviewTarget = interviewTarget + monthlySale.InterviewInterviewTarget.Float64

		for _, sale := range saleList {
			if param.Target == 0 {
				// 受注関連の計算
				if monthlySale.SalesMonth == sale.ContractSignedMonth {
					var (
						caSale    float64 = 0 // CAの売上
						raSale    float64 = 0 // RAの売上
						agentSale float64 = 0 // エージェントの売上

						caGrossProfit    float64 = 0 // CAの粗利
						raGrossProfit    float64 = 0 // RAの粗利
						agentGrossProfit float64 = 0 // エージェントの粗利
					)

					// CA売上比率に応じた売上を計算
					if param.AgentID == sale.CAAgentID && sale.CaSalesRatio.Float64 > 0 {
						caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
						caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
					}

					// RA売上比率に応じた売上を計算
					if param.AgentID == sale.RAAgentID && sale.RaSalesRatio.Float64 > 0 {
						raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
						raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
					}

					// エージェントの売上
					agentSale = caSale + raSale
					agentGrossProfit = caGrossProfit + raGrossProfit

					// 内定承諾の計算
					if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
						salesPerformance = salesPerformance + agentSale                    // 売上実績の計算
						grossProfitPerformance = grossProfitPerformance + agentGrossProfit // 粗利実績の計算
						accuracyAccept = accuracyAccept + agentSale                        // 確定(内定承諾)の請求金額合計
					}

					// Aヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
						accuracyA = accuracyA + agentSale // Aヨミの請求金額合計
					}

					// Bヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
						accuracyB = accuracyB + agentSale // Bヨミの請求金額合計
					}

					// Cヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
						accuracyC = accuracyC + agentSale // Cヨミの請求金額合計
					}

					// ネタの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
						accuracyTopic = accuracyTopic + agentSale // ネタの請求金額合計
					}
				}
			} else {
				// 請求関連の計算
				if monthlySale.SalesMonth == sale.BillingMonth {
					var (
						caSale    float64 = 0 // CAの売上
						raSale    float64 = 0 // RAの売上
						agentSale float64 = 0 // エージェントの売上

						caGrossProfit    float64 = 0 // CAの粗利
						raGrossProfit    float64 = 0 // RAの粗利
						agentGrossProfit float64 = 0 // エージェントの粗利
					)

					// CA売上比率に応じた売上を計算
					if param.AgentID == sale.CAAgentID && sale.CaSalesRatio.Float64 > 0 {
						caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
						caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
					}

					// RA売上比率に応じた売上を計算
					if param.AgentID == sale.RAAgentID && sale.RaSalesRatio.Float64 > 0 {
						raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
						raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
					}

					// エージェントの売上
					agentSale = caSale + raSale
					agentGrossProfit = caGrossProfit + raGrossProfit

					// 内定承諾の計算
					if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
						salesPerformance = salesPerformance + agentSale                    // 請求売上実績の計算
						grossProfitPerformance = grossProfitPerformance + agentGrossProfit // 請求粗利実績の計算
						accuracyAccept = accuracyAccept + agentSale                        // 確定(内定承諾)の請求金額合計
					}

					// Aヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
						accuracyA = accuracyA + agentSale // Aヨミの請求金額合計
					}

					// Bヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
						accuracyB = accuracyB + agentSale // Bヨミの請求金額合計
					}

					// Cヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
						accuracyC = accuracyC + agentSale // Cヨミの請求金額合計
					}

					// ネタの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
						accuracyTopic = accuracyTopic + agentSale // ネタの請求金額合計
					}
				}
			}
		}
	}

	/********* 受注売上 **********/

	// 受注売上達成率（売上実績 / 売上予算）
	if salesBudget > 0 {
		salesResultRate = math.Floor(salesPerformance / salesBudget * 100)
	}

	// 粗利達成率（粗利実績 / 粗利予算）
	if grossProfitBudget > 0 {
		grossProfitResultRate = math.Floor(grossProfitPerformance / grossProfitBudget * 100)
	}

	/********* 求職者数（面談実施月カウントベース） **********/

	// 内定承諾数
	interviewOfferAcceptancePerformance, err = i.interviewTaskGroupRepository.GetInterviewOfferAcceptancePerformanceCountByAgentIDAndPeriod(param.AgentID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 内定承諾達成率（内定承諾数 / 内定承諾目標 * 100）
	if interviewOfferAcceptanceTarget > 0 {
		interviewOfferAcceptanceResultRate = math.Floor(interviewOfferAcceptancePerformance / interviewOfferAcceptanceTarget * 100)
	}

	// 内定数
	interviewOfferPerformance, err = i.taskRepository.GetInterviewOfferPerformanceCountByAgentIDAndPeriod(param.AgentID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 内定数達成率（内定数 / 内定目標 * 100）
	if interviewOfferTarget > 0 {
		interviewOfferResultRate = math.Floor(interviewOfferPerformance / interviewOfferTarget * 100)
	}

	// 最終選考数
	interviewFinalSelectionPerformance, err = i.taskRepository.GetInterviewFinalSelectionPerformanceCountByAgentIDAndPeriod(param.AgentID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 最終選考数達成率（最終選考数 / 最終選考目標 * 100）
	if interviewFinalSelectionTarget > 0 {
		interviewFinalSelectionResultRate = math.Floor(interviewFinalSelectionPerformance / interviewFinalSelectionTarget * 100)
	}

	// 選考数
	interviewSelectionPerformance, err = i.taskRepository.GetInterviewSelectionPerformanceCountByAgentIDAndPeriod(param.AgentID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 選考数達成率（選考数 / 選考目標 * 100）
	if interviewSelectionTarget > 0 {
		interviewSelectionResultRate = math.Floor(interviewSelectionPerformance / interviewSelectionTarget * 100)
	}

	// 推薦完了数
	interviewRecommendationCompletionPerformance, err = i.taskRepository.GetInterviewRecommendationCompletionPerformanceCountByAgentIDAndPeriod(param.AgentID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 推薦完了数達成率（推薦完了数 / 推薦完了数目標 * 100）
	if interviewRecommendationCompletionTarget > 0 {
		interviewRecommendationCompletionResultRate = math.Floor(interviewRecommendationCompletionPerformance / interviewRecommendationCompletionTarget * 100)
	}

	// 求人紹介数
	interviewJobIntroductionPerformance, err = i.taskRepository.GetInterviewJobIntroductionPerformanceCountByAgentIDAndPeriod(param.AgentID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 求人紹介数達成率（求人紹介数 / 求人紹介数目標 * 100）
	if interviewJobIntroductionTarget > 0 {
		interviewJobIntroductionResultRate = math.Floor(interviewJobIntroductionPerformance / interviewJobIntroductionTarget * 100)
	}

	// 面談数
	interviewPerformance, err = i.interviewTaskGroupRepository.GetInterviewPerformanceCountByAgentIDAndPeriod(param.AgentID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 面談数達成率（面談数 / 面談数目標 * 100）
	if interviewTarget > 0 {
		interviewResultRate = math.Floor(interviewPerformance / interviewTarget * 100)
	}

	dashboard = entity.NewDashboard(
		param.ManagementID,
		fiscalYear,
		salesPerformance,
		salesBudget,
		salesResultRate,
		grossProfitPerformance,
		grossProfitBudget,
		grossProfitResultRate,
		interviewOfferAcceptanceTarget,
		interviewOfferAcceptancePerformance,
		interviewOfferAcceptanceResultRate,
		interviewOfferTarget,
		interviewOfferPerformance,
		interviewOfferResultRate,
		interviewFinalSelectionTarget,
		interviewFinalSelectionPerformance,
		interviewFinalSelectionResultRate,
		interviewSelectionTarget,
		interviewSelectionPerformance,
		interviewSelectionResultRate,
		interviewRecommendationCompletionTarget,
		interviewRecommendationCompletionPerformance,
		interviewRecommendationCompletionResultRate,
		interviewJobIntroductionTarget,
		interviewJobIntroductionPerformance,
		interviewJobIntroductionResultRate,
		interviewTarget,
		interviewPerformance,
		interviewResultRate,
		accuracyAccept,
		accuracyA,
		accuracyB,
		accuracyC,
		accuracyTopic,
	)

	return dashboard, nil
}

// {range: 会社全体, period: 1~11}
func searchStaffSaleForTheMonth(
	i *DashboardInteractorImpl,
	param entity.SearchDashboard, // 検索パラム
	fiscalYear string, // 決算月
	month string, // 検索対象の月
) (*entity.Dashboard, error) {
	var (
		dashboard    *entity.Dashboard
		saleList     []*entity.Sale
		staffID      = param.SearchRange  // 検索対象の担当者のID
		managementID = param.ManagementID // マネジメントID
	)

	/******************   計算に必要な「エージェントの売上情報」と「ヨミ情報」を取得   ******************/

	staffSaleManagement, err := i.staffSaleManagementRepository.FindByManagementIDAndStaffID(managementID, staffID)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	agentStaffMonthlySale, err := i.agentStaffMonthlySaleRepository.FindByStaffManagementIDAndMonth(staffSaleManagement.ID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	if param.Target == 0 {
		// 受注ベース
		saleList, err = i.saleRepository.GetContractSignedByStaffIDAndMonth(staffID, month)
		if err != nil {
			fmt.Println(err)
			return dashboard, err
		}
	} else {
		// 請求ベース
		saleList, err = i.saleRepository.GetBillingByStaffIDAndMonth(staffID, month)
		if err != nil {
			fmt.Println(err)
			return dashboard, err
		}
	}

	/******************   計算   ******************/

	var (
		salesPerformance                             float64 // 売上実績
		salesBudget                                  float64 // 売上予算
		salesResultRate                              float64 // 売上達成率
		grossProfitPerformance                       float64 // 粗利実績
		grossProfitBudget                            float64 // 粗利予算
		grossProfitResultRate                        float64 // 受注粗利達成率
		accuracyAccept                               float64 // 確定(内定承諾)
		accuracyA                                    float64 // Aヨミ
		accuracyB                                    float64 // Bヨミ
		accuracyC                                    float64 // Cヨミ
		accuracyTopic                                float64 // ネタ
		interviewOfferAcceptancePerformance          float64 // 内定承諾数
		interviewOfferAcceptanceResultRate           float64 // 内定承諾達成率
		interviewOfferPerformance                    float64 // 内定数
		interviewOfferResultRate                     float64 // 内定数達成率
		interviewFinalSelectionPerformance           float64 // 最終選考数
		interviewFinalSelectionResultRate            float64 // 最終選考数達成率
		interviewSelectionPerformance                float64 // 選考数
		interviewSelectionResultRate                 float64 // 選考数達成率
		interviewRecommendationCompletionPerformance float64 // 推薦完了数
		interviewRecommendationCompletionResultRate  float64 // 推薦完了数達成率
		interviewJobIntroductionPerformance          float64 // 求人紹介数
		interviewJobIntroductionResultRate           float64 // 求人紹介数達成率
		interviewPerformance                         float64 // 面談数
		interviewResultRate                          float64 // 面談数達成率
	)

	for _, sale := range saleList {
		if param.Target == 0 {
			// 受注関連の計算
			if agentStaffMonthlySale.SalesMonth == sale.ContractSignedMonth {
				var (
					caSale    float64 = 0 // CAの売上
					raSale    float64 = 0 // RAの売上
					staffSale float64 = 0 // スタッフの売上

					caGrossProfit    float64 = 0 // CAの粗利
					raGrossProfit    float64 = 0 // RAの粗利
					staffGrossProfit float64 = 0 // 担当者の粗利
				)

				// CA売上比率に応じた売上を計算
				if staffID == sale.CAStaffID && sale.CaSalesRatio.Float64 > 0 {
					caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
					caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
				}

				// RA売上比率に応じた売上を計算
				if staffID == sale.RAStaffID && sale.RaSalesRatio.Float64 > 0 {
					raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
					raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
				}

				// エージェントの売上
				staffSale = caSale + raSale
				staffGrossProfit = caGrossProfit + raGrossProfit

				// 内定承諾の計算
				if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
					salesPerformance = salesPerformance + staffSale                    // 売上実績の計算
					grossProfitPerformance = grossProfitPerformance + staffGrossProfit // 粗利実績の計算
					accuracyAccept = accuracyAccept + staffSale                        // 確定(内定承諾)の請求金額合計
				}

				// Aヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
					accuracyA = accuracyA + staffSale // Aヨミの請求金額合計
				}

				// Bヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
					accuracyB = accuracyB + staffSale // Bヨミの請求金額合計
				}

				// Cヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
					accuracyC = accuracyC + staffSale // Cヨミの請求金額合計
				}

				// ネタの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
					accuracyTopic = accuracyTopic + staffSale // ネタの請求金額合計
				}
			}
		} else {
			// 請求関連の計算
			if agentStaffMonthlySale.SalesMonth == sale.BillingMonth {
				var (
					caSale    float64 = 0 // CAの売上
					raSale    float64 = 0 // RAの売上
					staffSale float64 = 0 // スタッフの売上

					caGrossProfit    float64 = 0 // CAの粗利
					raGrossProfit    float64 = 0 // RAの粗利
					staffGrossProfit float64 = 0 // 担当者の粗利
				)

				// CA売上比率に応じた売上を計算
				if staffID == sale.CAStaffID && sale.CaSalesRatio.Float64 > 0 {
					caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
					caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
				}

				// RA売上比率に応じた売上を計算
				if staffID == sale.RAStaffID && sale.RaSalesRatio.Float64 > 0 {
					raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
					raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
				}

				// エージェントの売上
				staffSale = caSale + raSale
				staffGrossProfit = caGrossProfit + raGrossProfit

				// 内定承諾の計算
				if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
					salesPerformance = salesPerformance + staffSale                    // 受注売上実績の計算
					grossProfitPerformance = grossProfitPerformance + staffGrossProfit // 受注粗利実績の計算
					accuracyAccept = accuracyAccept + staffSale                        // 確定(内定承諾)の請求金額合計
				}

				// Aヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
					accuracyA = accuracyA + staffSale // Aヨミの請求金額合計
				}

				// Bヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
					accuracyB = accuracyB + staffSale // Bヨミの請求金額合計
				}

				// Cヨミの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
					accuracyC = accuracyC + staffSale // Cヨミの請求金額合計
				}

				// ネタの計算
				if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
					accuracyTopic = accuracyTopic + staffSale // ネタの請求金額合計
				}
			}
		}
	}

	/********* 受注売上 **********/

	if param.Target == 0 {
		// 受注予算
		salesBudget = agentStaffMonthlySale.OrderSalesBudget.Float64
		grossProfitBudget = agentStaffMonthlySale.OrderGrossProfitBudget.Float64

		// 受注売上達成率（売上実績 / 売上予算）
		if salesBudget > 0 {
			salesResultRate = math.Floor(salesPerformance / salesBudget * 100)
		}

		// 粗利達成率（粗利実績 / 粗利予算）
		if grossProfitBudget > 0 {
			grossProfitResultRate = math.Floor(grossProfitPerformance / grossProfitBudget * 100)
		}
	} else {
		// 請求予算
		salesBudget = agentStaffMonthlySale.ClaimSalesRevenueBudget.Float64
		grossProfitBudget = agentStaffMonthlySale.ClaimGrossMarginBudget.Float64

		// 請求売上達成率（売上実績 / 売上予算）
		if salesBudget > 0 {
			salesResultRate = math.Floor(salesPerformance / salesBudget * 100)
		}

		// 粗利達成率（粗利実績 / 粗利予算）
		if grossProfitBudget > 0 {
			grossProfitResultRate = math.Floor(grossProfitPerformance / grossProfitBudget * 100)
		}
	}

	/********* 求職者数（面談実施月カウントベース） **********/

	// 内定承諾数
	interviewOfferAcceptancePerformance, err = i.interviewTaskGroupRepository.GetInterviewOfferAcceptancePerformanceCountByStaffIDAndSalesMonth(staffID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 内定承諾達成率（内定承諾数 / 内定承諾目標 * 100）
	if agentStaffMonthlySale.InterviewOfferAcceptanceTarget.Float64 > 0 {
		interviewOfferAcceptanceResultRate = math.Floor(interviewOfferAcceptancePerformance / agentStaffMonthlySale.InterviewOfferAcceptanceTarget.Float64 * 100)
	}

	// 内定数
	interviewOfferPerformance, err = i.taskRepository.GetInterviewOfferPerformanceCountByStaffIDAndSalesMonth(staffID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 内定数達成率（内定数 / 内定目標 * 100）
	if agentStaffMonthlySale.InterviewOfferTarget.Float64 > 0 {
		interviewOfferResultRate = math.Floor(interviewOfferPerformance / agentStaffMonthlySale.InterviewOfferTarget.Float64 * 100)
	}

	// 最終選考数
	interviewFinalSelectionPerformance, err = i.taskRepository.GetInterviewFinalSelectionPerformanceCountByStaffIDAndSalesMonth(staffID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 最終選考数達成率（最終選考数 / 最終選考目標 * 100）
	if agentStaffMonthlySale.InterviewFinalSelectionTarget.Float64 > 0 {
		interviewFinalSelectionResultRate = math.Floor(interviewFinalSelectionPerformance / agentStaffMonthlySale.InterviewFinalSelectionTarget.Float64 * 100)
	}

	// 選考数
	interviewSelectionPerformance, err = i.taskRepository.GetInterviewSelectionPerformanceCountByStaffIDAndSalesMonth(staffID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 選考数達成率（選考数 / 選考目標 * 100）
	if agentStaffMonthlySale.InterviewSelectionTarget.Float64 > 0 {
		interviewSelectionResultRate = math.Floor(interviewSelectionPerformance / agentStaffMonthlySale.InterviewSelectionTarget.Float64 * 100)
	}

	// 推薦完了数
	interviewRecommendationCompletionPerformance, err = i.taskRepository.GetInterviewRecommendationCompletionPerformanceCountByStaffIDAndSalesMonth(staffID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 推薦完了数達成率（推薦完了数 / 推薦完了数目標 * 100）
	if agentStaffMonthlySale.InterviewRecommendationCompletionTarget.Float64 > 0 {
		interviewRecommendationCompletionResultRate = math.Floor(interviewRecommendationCompletionPerformance / agentStaffMonthlySale.InterviewRecommendationCompletionTarget.Float64 * 100)
	}

	// 求人紹介数
	interviewJobIntroductionPerformance, err = i.taskRepository.GetInterviewJobIntroductionPerformanceCountByStaffIDAndSalesMonth(staffID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 求人紹介数達成率（求人紹介数 / 求人紹介数目標 * 100）
	if agentStaffMonthlySale.InterviewJobIntroductionTarget.Float64 > 0 {
		interviewJobIntroductionResultRate = math.Floor(interviewJobIntroductionPerformance / agentStaffMonthlySale.InterviewJobIntroductionTarget.Float64 * 100)
	}

	// 面談数
	interviewPerformance, err = i.interviewTaskGroupRepository.GetInterviewPerformanceCountByStaffIDAndSalesMonth(staffID, month)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 面談数達成率（面談数 / 面談数目標 * 100）
	if agentStaffMonthlySale.InterviewInterviewTarget.Float64 > 0 {
		interviewResultRate = math.Floor(interviewPerformance / agentStaffMonthlySale.InterviewInterviewTarget.Float64 * 100)
	}

	dashboard = entity.NewDashboard(
		param.ManagementID,
		fiscalYear,
		salesPerformance,
		salesBudget,
		salesResultRate,
		grossProfitPerformance,
		grossProfitBudget,
		grossProfitResultRate,
		agentStaffMonthlySale.InterviewOfferAcceptanceTarget.Float64,
		interviewOfferAcceptancePerformance,
		interviewOfferAcceptanceResultRate,
		agentStaffMonthlySale.InterviewOfferTarget.Float64,
		interviewOfferPerformance,
		interviewOfferResultRate,
		agentStaffMonthlySale.InterviewFinalSelectionTarget.Float64,
		interviewFinalSelectionPerformance,
		interviewFinalSelectionResultRate,
		agentStaffMonthlySale.InterviewSelectionTarget.Float64,
		interviewSelectionPerformance,
		interviewSelectionResultRate,
		agentStaffMonthlySale.InterviewRecommendationCompletionTarget.Float64,
		interviewRecommendationCompletionPerformance,
		interviewRecommendationCompletionResultRate,
		agentStaffMonthlySale.InterviewJobIntroductionTarget.Float64,
		interviewJobIntroductionPerformance,
		interviewJobIntroductionResultRate,
		agentStaffMonthlySale.InterviewInterviewTarget.Float64,
		interviewPerformance,
		interviewResultRate,
		accuracyAccept,
		accuracyA,
		accuracyB,
		accuracyC,
		accuracyTopic,
	)

	return dashboard, nil
}

// {range: 担当者, period: 12~16}
func searchStaffSaleFromStartMonthToEndMonth(
	i *DashboardInteractorImpl,
	param entity.SearchDashboard, // 検索パラム
	fiscalYear string, // 決算月
	startMonth string, // 検索対象の開始月
	endMonth string, // 検索対象の終了月
) (*entity.Dashboard, error) {
	var (
		dashboard    *entity.Dashboard
		saleList     []*entity.Sale
		staffID      = param.SearchRange  // 検索対象の担当者のID
		managementID = param.ManagementID // マネジメントID
	)

	/******************   計算に必要な「エージェントの売上情報」と「ヨミ情報」を取得   ******************/

	staffSaleManagement, err := i.staffSaleManagementRepository.FindByManagementIDAndStaffID(managementID, staffID)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	agentStaffMonthlySaleList, err := i.agentStaffMonthlySaleRepository.GetByStaffManagementIDAndPeriod(staffSaleManagement.ID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	if param.Target == 0 {
		// 受注ベース
		saleList, err = i.saleRepository.GetContractSignedByStaffIDAndPeriod(staffID, startMonth, endMonth)
		if err != nil {
			fmt.Println(err)
			return dashboard, err
		}
	} else {
		// 請求ベース
		saleList, err = i.saleRepository.GetBillingByStaffIDAndPeriod(staffID, startMonth, endMonth)
		if err != nil {
			fmt.Println(err)
			return dashboard, err
		}
	}

	/******************   計算   ******************/

	var (
		salesPerformance                             float64 // 売上実績
		salesBudget                                  float64 // 売上予算
		salesResultRate                              float64 // 売上達成率
		grossProfitPerformance                       float64 // 粗利実績
		grossProfitBudget                            float64 // 粗利予算
		grossProfitResultRate                        float64 // 粗利達成率
		accuracyAccept                               float64 // 確定(内定承諾)
		accuracyA                                    float64 // Aヨミ
		accuracyB                                    float64 // Bヨミ
		accuracyC                                    float64 // Cヨミ
		accuracyTopic                                float64 // ネタ
		interviewOfferAcceptancePerformance          float64 // 内定承諾数
		interviewOfferAcceptanceTarget               float64 // 内定承諾目標
		interviewOfferAcceptanceResultRate           float64 // 内定承諾達成率
		interviewOfferPerformance                    float64 // 内定数
		interviewOfferTarget                         float64 // 内定数目標
		interviewOfferResultRate                     float64 // 内定数達成率
		interviewFinalSelectionPerformance           float64 // 最終選考数
		interviewFinalSelectionTarget                float64 // 最終選考数目標
		interviewFinalSelectionResultRate            float64 // 最終選考数達成率
		interviewSelectionPerformance                float64 // 選考数
		interviewSelectionTarget                     float64 // 選考数目標
		interviewSelectionResultRate                 float64 // 選考数達成率
		interviewRecommendationCompletionPerformance float64 // 推薦完了数
		interviewRecommendationCompletionTarget      float64 // 推薦完了目標
		interviewRecommendationCompletionResultRate  float64 // 推薦完了数達成率
		interviewJobIntroductionPerformance          float64 // 求人紹介数
		interviewJobIntroductionTarget               float64 // 求人紹介数目標
		interviewJobIntroductionResultRate           float64 // 求人紹介数達成率
		interviewPerformance                         float64 // 面談数
		interviewTarget                              float64 // 面談数目標
		interviewResultRate                          float64 // 面談数達成率
	)

	for _, monthlySale := range agentStaffMonthlySaleList {
		// 予算の合計値を計算
		if param.Target == 0 {
			// 受注関連の計算
			salesBudget = salesBudget + monthlySale.OrderSalesBudget.Float64                   // 売上
			grossProfitBudget = grossProfitBudget + monthlySale.OrderGrossProfitBudget.Float64 // 粗利
		} else {
			// 請求関連の計算
			salesBudget = salesBudget + monthlySale.ClaimSalesRevenueBudget.Float64            // 売上
			grossProfitBudget = grossProfitBudget + monthlySale.ClaimGrossMarginBudget.Float64 // 粗利
		}

		// 面談実施月ベースの目標の合計値を計算
		interviewOfferAcceptanceTarget = interviewOfferAcceptanceTarget + monthlySale.InterviewOfferAcceptanceTarget.Float64
		interviewOfferTarget = interviewOfferTarget + monthlySale.InterviewOfferTarget.Float64
		interviewFinalSelectionTarget = interviewFinalSelectionTarget + monthlySale.InterviewFinalSelectionTarget.Float64
		interviewSelectionTarget = interviewSelectionTarget + monthlySale.InterviewSelectionTarget.Float64
		interviewRecommendationCompletionTarget = interviewRecommendationCompletionTarget + monthlySale.InterviewRecommendationCompletionTarget.Float64
		interviewJobIntroductionTarget = interviewJobIntroductionTarget + monthlySale.InterviewJobIntroductionTarget.Float64
		interviewTarget = interviewTarget + monthlySale.InterviewInterviewTarget.Float64

		for _, sale := range saleList {
			if param.Target == 0 {
				// 受注関連の計算
				if monthlySale.SalesMonth == sale.ContractSignedMonth {
					var (
						caSale    float64 = 0 // CAの売上
						raSale    float64 = 0 // RAの売上
						staffSale float64 = 0 // スタッフの売上

						caGrossProfit    float64 = 0 // CAの粗利
						raGrossProfit    float64 = 0 // RAの粗利
						staffGrossProfit float64 = 0 // 担当者の粗利
					)

					// CA売上比率に応じた売上を計算
					if staffID == sale.CAStaffID && sale.CaSalesRatio.Float64 > 0 {
						caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
						caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
					}

					// RA売上比率に応じた売上を計算
					if staffID == sale.RAStaffID && sale.RaSalesRatio.Float64 > 0 {
						raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
						raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
					}

					// エージェントの売上
					staffSale = caSale + raSale
					staffGrossProfit = caGrossProfit + raGrossProfit

					// 確定(内定承諾)の計算
					if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
						salesPerformance = salesPerformance + staffSale                    // 受注売上実績の計算
						grossProfitPerformance = grossProfitPerformance + staffGrossProfit // 受注粗利実績の計算
						accuracyAccept = accuracyAccept + staffSale                        // 確定(内定承諾)の請求金額合計
					}

					// Aヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
						accuracyA = accuracyA + staffSale // Aヨミの請求金額合計
					}

					// Bヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
						accuracyB = accuracyB + staffSale // Bヨミの請求金額合計
					}

					// Cヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
						accuracyC = accuracyC + staffSale // Cヨミの請求金額合計
					}

					// ネタの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
						accuracyTopic = accuracyTopic + staffSale // ネタの請求金額合計
					}
				}
			} else {
				// 請求関連の計算
				if monthlySale.SalesMonth == sale.BillingMonth {
					var (
						caSale    float64 = 0 // CAの売上
						raSale    float64 = 0 // RAの売上
						staffSale float64 = 0 // スタッフの売上

						caGrossProfit    float64 = 0 // CAの粗利
						raGrossProfit    float64 = 0 // RAの粗利
						staffGrossProfit float64 = 0 // 担当者の粗利
					)

					// CA売上比率に応じた売上を計算
					if staffID == sale.CAStaffID && sale.CaSalesRatio.Float64 > 0 {
						caSale = sale.BillingAmount.Float64 * sale.CaSalesRatio.Float64 / 100
						caGrossProfit = sale.GrossProfit.Float64 * sale.CaSalesRatio.Float64 / 100
					}

					// RA売上比率に応じた売上を計算
					if staffID == sale.RAStaffID && sale.RaSalesRatio.Float64 > 0 {
						raSale = sale.BillingAmount.Float64 * sale.RaSalesRatio.Float64 / 100
						raGrossProfit = sale.GrossProfit.Float64 * sale.RaSalesRatio.Float64 / 100
					}

					// エージェントの売上
					staffSale = caSale + raSale
					staffGrossProfit = caGrossProfit + raGrossProfit

					// 内定承諾の計算
					if sale.Accuracy == null.NewInt(entity.AccuracyAccept, true) {
						salesPerformance = salesPerformance + staffSale                    // 請求売上実績の計算
						grossProfitPerformance = grossProfitPerformance + staffGrossProfit // 請求粗利実績の計算
						accuracyAccept = accuracyAccept + staffSale                        // 確定(内定承諾)の請求金額合計
					}

					// Aヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyA, true) {
						accuracyA = accuracyA + staffSale // Aヨミの請求金額合計
					}

					// Bヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyB, true) {
						accuracyB = accuracyB + staffSale // Bヨミの請求金額合計
					}

					// Cヨミの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyC, true) {
						accuracyC = accuracyC + staffSale // Cヨミの請求金額合計
					}

					// ネタの計算
					if sale.Accuracy == null.NewInt(entity.AccuracyTopic, true) {
						accuracyTopic = accuracyTopic + staffSale // ネタの請求金額合計
					}
				}
			}
		}
	}

	/********* 受注売上 **********/

	// 受注売上達成率（売上実績 / 売上予算）
	if salesBudget > 0 {
		salesResultRate = math.Floor(salesPerformance / salesBudget * 100)
	}

	// 粗利達成率（粗利実績 / 粗利予算）
	if grossProfitBudget > 0 {
		grossProfitResultRate = math.Floor(grossProfitPerformance / grossProfitBudget * 100)
	}

	/********* 求職者数（面談実施月カウントベース） **********/

	// 内定承諾数
	interviewOfferAcceptancePerformance, err = i.interviewTaskGroupRepository.GetInterviewOfferAcceptancePerformanceCountByStaffIDAndPeriod(staffID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 内定承諾達成率（内定承諾数 / 内定承諾目標 * 100）
	if interviewOfferAcceptanceTarget > 0 {
		interviewOfferAcceptanceResultRate = math.Floor(interviewOfferAcceptancePerformance / interviewOfferAcceptanceTarget * 100)
	}

	// 内定数
	interviewOfferPerformance, err = i.taskRepository.GetInterviewOfferPerformanceCountByStaffIDAndPeriod(staffID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 内定数達成率（内定数 / 内定目標 * 100）
	if interviewOfferTarget > 0 {
		interviewOfferResultRate = math.Floor(interviewOfferPerformance / interviewOfferTarget * 100)
	}

	// 最終選考数
	interviewFinalSelectionPerformance, err = i.taskRepository.GetInterviewFinalSelectionPerformanceCountByStaffIDAndPeriod(staffID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 最終選考数達成率（最終選考数 / 最終選考目標 * 100）
	if interviewFinalSelectionTarget > 0 {
		interviewFinalSelectionResultRate = math.Floor(interviewFinalSelectionPerformance / interviewFinalSelectionTarget * 100)
	}

	// 選考数
	interviewSelectionPerformance, err = i.taskRepository.GetInterviewSelectionPerformanceCountByStaffIDAndPeriod(staffID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 選考数達成率（選考数 / 選考目標 * 100）
	if interviewSelectionTarget > 0 {
		interviewSelectionResultRate = math.Floor(interviewSelectionPerformance / interviewSelectionTarget * 100)
	}

	// 推薦完了数
	interviewRecommendationCompletionPerformance, err = i.taskRepository.GetInterviewRecommendationCompletionPerformanceCountByStaffIDAndPeriod(staffID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 推薦完了数達成率（推薦完了数 / 推薦完了数目標 * 100）
	if interviewRecommendationCompletionTarget > 0 {
		interviewRecommendationCompletionResultRate = math.Floor(interviewRecommendationCompletionPerformance / interviewRecommendationCompletionTarget * 100)
	}

	// 求人紹介数
	interviewJobIntroductionPerformance, err = i.taskRepository.GetInterviewJobIntroductionPerformanceCountByStaffIDAndPeriod(staffID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 求人紹介数達成率（求人紹介数 / 求人紹介数目標 * 100）
	if interviewJobIntroductionTarget > 0 {
		interviewJobIntroductionResultRate = math.Floor(interviewJobIntroductionPerformance / interviewJobIntroductionTarget * 100)
	}

	// 面談数
	interviewPerformance, err = i.interviewTaskGroupRepository.GetInterviewPerformanceCountByStaffIDAndPeriod(staffID, startMonth, endMonth)
	if err != nil {
		fmt.Println(err)
		return dashboard, err
	}

	// 面談数達成率（面談数 / 面談数目標 * 100）
	if interviewTarget > 0 {
		interviewResultRate = math.Floor(interviewPerformance / interviewTarget * 100)
	}

	dashboard = entity.NewDashboard(
		param.ManagementID,
		fiscalYear,
		salesPerformance,
		salesBudget,
		salesResultRate,
		grossProfitPerformance,
		grossProfitBudget,
		grossProfitResultRate,
		interviewOfferAcceptanceTarget,
		interviewOfferAcceptancePerformance,
		interviewOfferAcceptanceResultRate,
		interviewOfferTarget,
		interviewOfferPerformance,
		interviewOfferResultRate,
		interviewFinalSelectionTarget,
		interviewFinalSelectionPerformance,
		interviewFinalSelectionResultRate,
		interviewSelectionTarget,
		interviewSelectionPerformance,
		interviewSelectionResultRate,
		interviewRecommendationCompletionTarget,
		interviewRecommendationCompletionPerformance,
		interviewRecommendationCompletionResultRate,
		interviewJobIntroductionTarget,
		interviewJobIntroductionPerformance,
		interviewJobIntroductionResultRate,
		interviewTarget,
		interviewPerformance,
		interviewResultRate,
		accuracyAccept,
		accuracyA,
		accuracyB,
		accuracyC,
		accuracyTopic,
	)

	return dashboard, nil
}

// エージェントIDで求職者一覧の子テーブル情報を取得する関数
func getJobSeekerInformationsByJobSeekerList(
	i *DashboardInteractorImpl,
	jobSeekerList []*entity.JobSeeker,
) ([]*entity.JobSeeker, error) {
	var (
		lenJobSeekerList = len(jobSeekerList)
		jobSeekerIDList  = make([]uint, 0, lenJobSeekerList)
	)
	if lenJobSeekerList == 0 {
		return jobSeekerList, nil
	}
	for _, jobSeeker := range jobSeekerList {
		jobSeekerIDList = append(jobSeekerIDList, jobSeeker.ID)
	}
	/******************   「リリース求職者情報」を取得   ******************/

	// 求職者IDリストから関連テーブル情報を取得
	// selfpromotionとdocumentは使わないので取得しない
	studentHistory, err := i.jobSeekerStudentHistoryRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	workHistory, err := i.jobSeekerWorkHistoryRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	experienceIndustry, err := i.jobSeekerExperienceIndustryRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	departmentHistory, err := i.jobSeekerDepartmentHistoryRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	experienceOccupation, err := i.jobSeekerExperienceOccupationRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	desiredCompanyScale, err := i.jobSeekerDesiredCompanyScaleRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	license, err := i.jobSeekerLicenseRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	desiredIndustry, err := i.jobSeekerDesiredIndustryRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	desiredOccupation, err := i.jobSeekerDesiredOccupationRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	desiredWorkLocation, err := i.jobSeekerDesiredWorkLocationRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	desiredHolidayType, err := i.jobSeekerDesiredHolidayTypeRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	developmentSkill, err := i.jobSeekerDevelopmentSkillRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	languageSkill, err := i.jobSeekerLanguageSkillRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	pcTool, err := i.jobSeekerPCToolRepository.GetByJobSeekerIDList(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	for _, jobSeeker := range jobSeekerList {
		for _, sh := range studentHistory {
			if jobSeeker.ID == sh.JobSeekerID {
				jobSeeker.StudentHistories = append(jobSeeker.StudentHistories, *sh)
			}
		}

		for _, wh := range workHistory {
			if jobSeeker.ID == wh.JobSeekerID {
				for _, ei := range experienceIndustry {
					if wh.ID == ei.WorkHistoryID {
						wh.ExperienceIndustries = append(wh.ExperienceIndustries, *ei)
					}
				}

				for _, dh := range departmentHistory {
					if wh.ID == dh.WorkHistoryID {
						for _, eo := range experienceOccupation {
							if dh.ID == eo.DepartmentHistoryID {
								dh.ExperienceOccupations = append(dh.ExperienceOccupations, *eo)
							}
						}
						wh.DepartmentHistories = append(wh.DepartmentHistories, *dh)
					}
				}
				jobSeeker.WorkHistories = append(jobSeeker.WorkHistories, *wh)
			}
		}

		for _, dcs := range desiredCompanyScale {
			if jobSeeker.ID == dcs.JobSeekerID {
				jobSeeker.DesiredCompanyScales = append(jobSeeker.DesiredCompanyScales, *dcs)
			}
		}

		for _, l := range license {
			if jobSeeker.ID == l.JobSeekerID {
				jobSeeker.Licenses = append(jobSeeker.Licenses, *l)
			}
		}

		for _, di := range desiredIndustry {
			if jobSeeker.ID == di.JobSeekerID {
				jobSeeker.DesiredIndustries = append(jobSeeker.DesiredIndustries, *di)
			}
		}

		for _, do := range desiredOccupation {
			if jobSeeker.ID == do.JobSeekerID {
				jobSeeker.DesiredOccupations = append(jobSeeker.DesiredOccupations, *do)
			}
		}

		for _, dwl := range desiredWorkLocation {
			if jobSeeker.ID == dwl.JobSeekerID {
				jobSeeker.DesiredWorkLocations = append(jobSeeker.DesiredWorkLocations, *dwl)
			}
		}

		for _, dht := range desiredHolidayType {
			if jobSeeker.ID == dht.JobSeekerID {
				jobSeeker.DesiredHolidayTypes = append(jobSeeker.DesiredHolidayTypes, *dht)
			}
		}

		for _, ds := range developmentSkill {
			if jobSeeker.ID == ds.JobSeekerID {
				jobSeeker.DevelopmentSkills = append(jobSeeker.DevelopmentSkills, *ds)
			}
		}

		for _, ls := range languageSkill {
			if jobSeeker.ID == ls.JobSeekerID {
				jobSeeker.LanguageSkills = append(jobSeeker.LanguageSkills, *ls)
			}
		}

		for _, ps := range pcTool {
			if jobSeeker.ID == ps.JobSeekerID {
				jobSeeker.PCTools = append(jobSeeker.PCTools, *ps)
			}
		}
	}

	return jobSeekerList, nil
}

// 担当者IDで求職者一覧の子テーブル情報を取得する関数
func getJobSeekerInformationsByStaffID(
	i *DashboardInteractorImpl,
	jobSeekerList []*entity.JobSeeker,
	staffID uint,
) ([]*entity.JobSeeker, error) {
	/******************   「リリース求職者情報」を取得   ******************/

	// 求職者IDリストから関連テーブル情報を取得
	// selfpromotionとdocumentは使わないので取得しない
	studentHistory, err := i.jobSeekerStudentHistoryRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	workHistory, err := i.jobSeekerWorkHistoryRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	experienceIndustry, err := i.jobSeekerExperienceIndustryRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	departmentHistory, err := i.jobSeekerDepartmentHistoryRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	experienceOccupation, err := i.jobSeekerExperienceOccupationRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	desiredCompanyScale, err := i.jobSeekerDesiredCompanyScaleRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	license, err := i.jobSeekerLicenseRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	desiredIndustry, err := i.jobSeekerDesiredIndustryRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	desiredOccupation, err := i.jobSeekerDesiredOccupationRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	desiredWorkLocation, err := i.jobSeekerDesiredWorkLocationRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	desiredHolidayType, err := i.jobSeekerDesiredHolidayTypeRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	developmentSkill, err := i.jobSeekerDevelopmentSkillRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	languageSkill, err := i.jobSeekerLanguageSkillRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	pcTool, err := i.jobSeekerPCToolRepository.GetByStaffID(staffID)
	if err != nil {
		fmt.Println(err)
		return jobSeekerList, err
	}

	for _, jobSeeker := range jobSeekerList {
		for _, sh := range studentHistory {
			if jobSeeker.ID == sh.JobSeekerID {
				jobSeeker.StudentHistories = append(jobSeeker.StudentHistories, *sh)
			}
		}

		for _, wh := range workHistory {
			if jobSeeker.ID == wh.JobSeekerID {
				for _, ei := range experienceIndustry {
					if wh.ID == ei.WorkHistoryID {
						wh.ExperienceIndustries = append(wh.ExperienceIndustries, *ei)
					}
				}

				for _, dh := range departmentHistory {
					if wh.ID == dh.WorkHistoryID {
						for _, eo := range experienceOccupation {
							if dh.ID == eo.DepartmentHistoryID {
								dh.ExperienceOccupations = append(dh.ExperienceOccupations, *eo)
							}
						}
						wh.DepartmentHistories = append(wh.DepartmentHistories, *dh)
					}
				}
				jobSeeker.WorkHistories = append(jobSeeker.WorkHistories, *wh)
			}
		}

		for _, dcs := range desiredCompanyScale {
			if jobSeeker.ID == dcs.JobSeekerID {
				jobSeeker.DesiredCompanyScales = append(jobSeeker.DesiredCompanyScales, *dcs)
			}
		}

		for _, l := range license {
			if jobSeeker.ID == l.JobSeekerID {
				jobSeeker.Licenses = append(jobSeeker.Licenses, *l)
			}
		}

		for _, di := range desiredIndustry {
			if jobSeeker.ID == di.JobSeekerID {
				jobSeeker.DesiredIndustries = append(jobSeeker.DesiredIndustries, *di)
			}
		}

		for _, do := range desiredOccupation {
			if jobSeeker.ID == do.JobSeekerID {
				jobSeeker.DesiredOccupations = append(jobSeeker.DesiredOccupations, *do)
			}
		}

		for _, dwl := range desiredWorkLocation {
			if jobSeeker.ID == dwl.JobSeekerID {
				jobSeeker.DesiredWorkLocations = append(jobSeeker.DesiredWorkLocations, *dwl)
			}
		}

		for _, dht := range desiredHolidayType {
			if jobSeeker.ID == dht.JobSeekerID {
				jobSeeker.DesiredHolidayTypes = append(jobSeeker.DesiredHolidayTypes, *dht)
			}
		}

		for _, ds := range developmentSkill {
			if jobSeeker.ID == ds.JobSeekerID {
				jobSeeker.DevelopmentSkills = append(jobSeeker.DevelopmentSkills, *ds)
			}
		}

		for _, ls := range languageSkill {
			if jobSeeker.ID == ls.JobSeekerID {
				jobSeeker.LanguageSkills = append(jobSeeker.LanguageSkills, *ls)
			}
		}

		for _, ps := range pcTool {
			if jobSeeker.ID == ps.JobSeekerID {
				jobSeeker.PCTools = append(jobSeeker.PCTools, *ps)
			}
		}
	}

	return jobSeekerList, nil
}

/****************************************************************************************/
/// CSV API
//
// CSV出力
type ExportAccuracyCSVInput struct {
	AgentID uint
}

type ExportAccuracyCSVOutput struct {
	FilePath *entity.FilePath
}

func (i *DashboardInteractorImpl) ExportAccuracyCSV(input ExportAccuracyCSVInput) (ExportAccuracyCSVOutput, error) {
	var (
		output  ExportAccuracyCSVOutput
		err     error
		records [][]string
		record  []string
	)

	// csvの一行目を作成
	records = append(
		records,
		[]string{
			"ヨミ",
			"会社名/求人タイトル",
			"求職者名",
			"フェーズ",
			"受注月",
			"請求月",
			"請求金額(円)",
			"原価(円)",
			"粗利(円)",
			"担当RA",
			"担当CA",
		},
	)

	/******************   「ヨミ情報」を取得（デフォルトは受注）   ******************/

	saleList, err := i.saleRepository.GetByAgentIDForCSV(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sale := range saleList {
		phase, phaseSub := getStrTaskPhaseAndPhaseSub(sale.PhaseCategory, sale.PhaseSubCategory)
		record = []string{
			getStrAccuracy(sale.Accuracy),
			sale.CompanyName + "/" + sale.Title,
			sale.LastName + sale.FirstName,
			phase + "/" + phaseSub,
			sale.ContractSignedMonth,
			sale.BillingMonth,
			fmt.Sprint(sale.BillingAmount.NullFloat64.Float64),
			fmt.Sprint(sale.Cost.NullFloat64.Float64),
			fmt.Sprint(sale.GrossProfit.NullFloat64.Float64),
			sale.RAAgentName + "/" + sale.RAStaffName,
			sale.CAAgentName + "/" + sale.CAStaffName,
		}
		records = append(records, record)
	}

	// CSVファイルを作成
	filePath := ("./accuracy-" + fmt.Sprint(utility.CreateUUID()) + ".csv")

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	cw := csv.NewWriter(file)

	defer file.Close()
	defer cw.Flush()

	cw.WriteAll(records)

	if err := cw.Error(); err != nil {
		fmt.Println("error writing csv:", err)
		return output, err
	}

	output.FilePath = entity.NewFilePath(filePath)

	return output, nil
}

/****************************************************************************************/
