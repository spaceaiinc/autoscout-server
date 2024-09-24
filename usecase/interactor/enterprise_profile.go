package interactor

import (
	"fmt"
	"strconv"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type EnterpriseProfileInteractor interface {
	// 汎用系 API
	CreateEnterpriseProfile(input CreateEnterpriseProfileInput) (CreateEnterpriseProfileOutput, error)
	UpdateEnterpriseProfile(input UpdateEnterpriseProfileInput) (UpdateEnterpriseProfileOutput, error)
	DeleteEnterpriseProfile(input DeleteEnterpriseProfileInput) (DeleteEnterpriseProfileOutput, error)
	GetEnterpriseByID(input GetEnterpriseByIDInput) (GetEnterpriseByIDOutput, error)                                           // 指定IDの企業情報を削除する関数
	GetEnterpriseListByAgentStaffID(input GetEnterpriseListByAgentStaffIDInput) (GetEnterpriseListByAgentStaffIDOutput, error) // agentStaffのIDを使って企業一覧情報取得する関数
	GetEnterpriseListByAgentID(input GetEnterpriseListByAgentIDInput) (GetEnterpriseListByAgentIDOutput, error)                // agentIDを使って企業一覧情報取得する関数
	DeleteEnterpriseReferenceMaterial(input DeleteEnterpriseReferenceMaterialInput) (DeleteEnterpriseReferenceMaterialOutput, error)

	// ページネーション用
	GetEnterpriseListByAgentIDAndPage(input GetEnterpriseListByAgentIDAndPageInput) (GetEnterpriseListByAgentIDAndPageOutput, error) // agentIDを使って企業一覧情報取得する関数

	// 企業の絞り込み検索
	GetSearchEnterpriseListByAgentID(input GetSearchEnterpriseListByAgentIDInput) (GetSearchEnterpriseListByAgentIDOutput, error)

	// 企業資料関連
	CreateEnterpriseReferenceMaterial(input CreateEnterpriseReferenceMaterialInput) (CreateEnterpriseReferenceMaterialOutput, error)
	UpdateEnterpriseReferenceMaterial(input UpdateEnterpriseReferenceMaterialInput) (UpdateEnterpriseReferenceMaterialOutput, error)
	GetEnterpriseReferenceMaterialByEnterpriseID(input GetEnterpriseReferenceMaterialByEnterpriseIDInput) (GetEnterpriseReferenceMaterialByEnterpriseIDOutput, error)

	// CSV操作系 API
	ImportEnterpriseCSV(input ImportEnterpriseCSVInput) (ImportEnterpriseCSVOutput, error) // CSVファイルを読み込んで企業情報を登録する関数
	ImportJobInformationCSV(input ImportJobInformationCSVInput) (ImportJobInformationCSVOutput, error)
	ImportEnterpriseCSVForCircus(input ImportEnterpriseCSVForCircusInput) (ImportEnterpriseCSVForCircusOutput, error)
	ImportEnterpriseCSVForAgentBank(input ImportEnterpriseCSVForAgentBankInput) (ImportEnterpriseCSVForAgentBankOutput, error)
	ExportEnterpriseCSV(input ExportEnterpriseCSVInput) (ExportEnterpriseCSVOutput, error) // CSVファイルを出力する関数

	// リストで企業情報を作成する
	ImportEnterpriseJSON(input ImportEnterpriseJSONInput) (ImportEnterpriseJSONOutput, error)

	// 求人企業の追加情報 API
	CreateEnterpriseActivity(input CreateEnterpriseActivityInput) (CreateEnterpriseActivityOutput, error)

	// Admin API
	GetInitialEnterprise(input GetInitialEnterpriseInput) (GetInitialEnterpriseOutput, error) // 企業一覧（すべて）
}

type EnterpriseProfileInteractorImpl struct {
	firebase                                           usecase.Firebase
	sendgrid                                           config.Sendgrid
	enterpriseProfileRepository                        usecase.EnterpriseProfileRepository
	enterpriseIndustryRepository                       usecase.EnterpriseIndustryRepository
	enterpriseReferenceMaterialRepository              usecase.EnterpriseReferenceMaterialRepository
	enterpriseActivityRepository                       usecase.EnterpriseActivityRepository
	billingAddressRepository                           usecase.BillingAddressRepository
	billingAddressHRStaffRepository                    usecase.BillingAddressHRStaffRepository
	billingAddressRAStaffRepository                    usecase.BillingAddressRAStaffRepository
	jobInformationRepository                           usecase.JobInformationRepository
	jobInfoTargetRepository                            usecase.JobInformationTargetRepository
	jobInfoFeatureRepository                           usecase.JobInformationFeatureRepository
	jobInfoPrefectureRepository                        usecase.JobInformationPrefectureRepository
	jobInfoWorkCharmPointRepository                    usecase.JobInformationWorkCharmPointRepository
	jobInfoEmploymentStatusRepository                  usecase.JobInformationEmploymentStatusRepository
	jobInfoRequiredConditionRepository                 usecase.JobInformationRequiredConditionRepository
	jobInfoRequiredLicenseRepository                   usecase.JobInformationRequiredLicenseRepository
	jobInfoRequiredPCToolRepository                    usecase.JobInformationRequiredPCToolRepository
	jobInfoRequiredLanguageRepository                  usecase.JobInformationRequiredLanguageRepository
	jobInfoRequiredLanguageTypeRepository              usecase.JobInformationRequiredLanguageTypeRepository
	jobInfoRequiredExperienceDevelopmentRepository     usecase.JobInformationRequiredExperienceDevelopmentRepository
	jobInfoRequiredExperienceDevelopmentTypeRepository usecase.JobInformationRequiredExperienceDevelopmentTypeRepository
	jobInfoRequiredExperienceJobRepository             usecase.JobInformationRequiredExperienceJobRepository
	jobInfoRequiredExperienceIndustryRepository        usecase.JobInformationRequiredExperienceIndustryRepository
	jobInfoRequiredExperienceOccupationRepository      usecase.JobInformationRequiredExperienceOccupationRepository
	jobInfoRequiredSocialExperienceRepository          usecase.JobInformationRequiredSocialExperienceRepository
	jobInfoSelectionFlowPatternRepository              usecase.JobInformationSelectionFlowPatternRepository
	jobInfoSelectionInformationRepository              usecase.JobInformationSelectionInformationRepository
	jobInfoOccupationRepository                        usecase.JobInformationOccupationRepository
	agentRepository                                    usecase.AgentRepository
	agentStaffRepository                               usecase.AgentStaffRepository
	jobInfoHideToAgentRepository                       usecase.JobInformationHideToAgentRepository
	jobInformationExternalIDRepository                 usecase.JobInformationExternalIDRepository
}

// EnterpriseProfileInteractorImpl is an implementation of EnterpriseProfileInteractor
func NewEnterpriseProfileInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	epR usecase.EnterpriseProfileRepository,
	eiR usecase.EnterpriseIndustryRepository,
	ermR usecase.EnterpriseReferenceMaterialRepository,
	eaR usecase.EnterpriseActivityRepository,
	baR usecase.BillingAddressRepository,
	bahsR usecase.BillingAddressHRStaffRepository,
	barsR usecase.BillingAddressRAStaffRepository,
	jR usecase.JobInformationRepository,
	jtR usecase.JobInformationTargetRepository,
	jfR usecase.JobInformationFeatureRepository,
	jpR usecase.JobInformationPrefectureRepository,
	jwcpR usecase.JobInformationWorkCharmPointRepository,
	jesR usecase.JobInformationEmploymentStatusRepository,
	jrcR usecase.JobInformationRequiredConditionRepository,
	jrlR usecase.JobInformationRequiredLicenseRepository,
	jrptR usecase.JobInformationRequiredPCToolRepository,
	jrlgR usecase.JobInformationRequiredLanguageRepository,
	jrlgtR usecase.JobInformationRequiredLanguageTypeRepository,
	jredR usecase.JobInformationRequiredExperienceDevelopmentRepository,
	jredtR usecase.JobInformationRequiredExperienceDevelopmentTypeRepository,
	jrejR usecase.JobInformationRequiredExperienceJobRepository,
	jreiR usecase.JobInformationRequiredExperienceIndustryRepository,
	jreoR usecase.JobInformationRequiredExperienceOccupationRepository,
	jrseR usecase.JobInformationRequiredSocialExperienceRepository,
	jsfpR usecase.JobInformationSelectionFlowPatternRepository,
	jsiR usecase.JobInformationSelectionInformationRepository,
	joR usecase.JobInformationOccupationRepository,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
	jhtaR usecase.JobInformationHideToAgentRepository,
	jeiR usecase.JobInformationExternalIDRepository,
) EnterpriseProfileInteractor {
	return &EnterpriseProfileInteractorImpl{
		firebase:                                           fb,
		sendgrid:                                           sg,
		enterpriseProfileRepository:                        epR,
		enterpriseIndustryRepository:                       eiR,
		enterpriseReferenceMaterialRepository:              ermR,
		enterpriseActivityRepository:                       eaR,
		billingAddressRepository:                           baR,
		billingAddressHRStaffRepository:                    bahsR,
		billingAddressRAStaffRepository:                    barsR,
		jobInformationRepository:                           jR,
		jobInfoTargetRepository:                            jtR,
		jobInfoFeatureRepository:                           jfR,
		jobInfoPrefectureRepository:                        jpR,
		jobInfoWorkCharmPointRepository:                    jwcpR,
		jobInfoEmploymentStatusRepository:                  jesR,
		jobInfoRequiredConditionRepository:                 jrcR,
		jobInfoRequiredLicenseRepository:                   jrlR,
		jobInfoRequiredPCToolRepository:                    jrptR,
		jobInfoRequiredLanguageRepository:                  jrlgR,
		jobInfoRequiredLanguageTypeRepository:              jrlgtR,
		jobInfoRequiredExperienceDevelopmentRepository:     jredR,
		jobInfoRequiredExperienceDevelopmentTypeRepository: jredtR,
		jobInfoRequiredExperienceJobRepository:             jrejR,
		jobInfoRequiredExperienceIndustryRepository:        jreiR,
		jobInfoRequiredExperienceOccupationRepository:      jreoR,
		jobInfoRequiredSocialExperienceRepository:          jrseR,
		jobInfoSelectionFlowPatternRepository:              jsfpR,
		jobInfoSelectionInformationRepository:              jsiR,
		jobInfoOccupationRepository:                        joR,
		agentRepository:                                    aR,
		agentStaffRepository:                               asR,
		jobInfoHideToAgentRepository:                       jhtaR,
		jobInformationExternalIDRepository:                 jeiR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
//求人企業の作成
type CreateEnterpriseProfileInput struct {
	CreateParam entity.CreateOrUpdateEnterpriseProfileParam
}

type CreateEnterpriseProfileOutput struct {
	EnterpriseProfile *entity.EnterpriseProfile
}

func (i *EnterpriseProfileInteractorImpl) CreateEnterpriseProfile(input CreateEnterpriseProfileInput) (CreateEnterpriseProfileOutput, error) {
	var (
		output CreateEnterpriseProfileOutput
		err    error
	)

	enterpriseProfile := entity.NewEnterpriseProfile(
		input.CreateParam.CompanyName,
		input.CreateParam.AgentStaffID,
		input.CreateParam.CorporateSiteURL,
		input.CreateParam.Representative,
		input.CreateParam.Establishment,
		input.CreateParam.PostCode,
		input.CreateParam.OfficeLocation,
		input.CreateParam.EmployeeNumberSingle,
		input.CreateParam.EmployeeNumberGroup,
		input.CreateParam.Capital,
		input.CreateParam.PublicOffering,
		input.CreateParam.EarningsYear,
		input.CreateParam.Earnings,
		input.CreateParam.BusinessDetail,
	)

	err = i.enterpriseProfileRepository.Create(enterpriseProfile)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, industry := range input.CreateParam.Industries {
		industry := entity.NewEnterpriseIndustry(
			enterpriseProfile.ID,
			industry,
		)

		err = i.enterpriseIndustryRepository.Create(industry)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.EnterpriseProfile = enterpriseProfile
	output.EnterpriseProfile.Industries = input.CreateParam.Industries

	return output, nil
}

// 求人企業の更新
type UpdateEnterpriseProfileInput struct {
	UpdateParam  entity.CreateOrUpdateEnterpriseProfileParam
	EnterpriseID uint
}

type UpdateEnterpriseProfileOutput struct {
	EnterpriseProfile *entity.EnterpriseProfile
}

func (i *EnterpriseProfileInteractorImpl) UpdateEnterpriseProfile(input UpdateEnterpriseProfileInput) (UpdateEnterpriseProfileOutput, error) {
	var (
		output UpdateEnterpriseProfileOutput
		err    error
	)

	enterpriseProfile := entity.NewEnterpriseProfile(
		input.UpdateParam.CompanyName,
		input.UpdateParam.AgentStaffID,
		input.UpdateParam.CorporateSiteURL,
		input.UpdateParam.Representative,
		input.UpdateParam.Establishment,
		input.UpdateParam.PostCode,
		input.UpdateParam.OfficeLocation,
		input.UpdateParam.EmployeeNumberSingle,
		input.UpdateParam.EmployeeNumberGroup,
		input.UpdateParam.Capital,
		input.UpdateParam.PublicOffering,
		input.UpdateParam.EarningsYear,
		input.UpdateParam.Earnings,
		input.UpdateParam.BusinessDetail,
	)

	err = i.enterpriseProfileRepository.Update(input.EnterpriseID, enterpriseProfile)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	err = i.enterpriseIndustryRepository.DeleteByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, industry := range input.UpdateParam.Industries {
		industry := entity.NewEnterpriseIndustry(
			input.EnterpriseID,
			industry,
		)

		err = i.enterpriseIndustryRepository.Create(industry)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.EnterpriseProfile = enterpriseProfile
	output.EnterpriseProfile.Industries = input.UpdateParam.Industries

	return output, nil
}

// 求人企業の削除
type DeleteEnterpriseProfileInput struct {
	DeleteParam entity.DeleteEnterpriseProfileParam
}

type DeleteEnterpriseProfileOutput struct {
	OK bool
}

func (i *EnterpriseProfileInteractorImpl) DeleteEnterpriseProfile(input DeleteEnterpriseProfileInput) (DeleteEnterpriseProfileOutput, error) {
	var output DeleteEnterpriseProfileOutput

	err := i.enterpriseProfileRepository.Delete(input.DeleteParam.ID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

// 求人企業の削除
type DeleteEnterpriseReferenceMaterialInput struct {
	ReferenceMaterialID uint
	MaterialType        uint
}

type DeleteEnterpriseReferenceMaterialOutput struct {
	OK bool
}

func (i *EnterpriseProfileInteractorImpl) DeleteEnterpriseReferenceMaterial(input DeleteEnterpriseReferenceMaterialInput) (DeleteEnterpriseReferenceMaterialOutput, error) {
	var output DeleteEnterpriseReferenceMaterialOutput

	err := i.enterpriseReferenceMaterialRepository.UpdateMaterialTypeByID(input.ReferenceMaterialID, input.MaterialType)
	if err != nil {
		return output, err
	}
	output.OK = true

	return output, nil
}

// 企業IDを使って企業情報を取得する
type GetEnterpriseByIDInput struct {
	EnterpriseID uint
}

type GetEnterpriseByIDOutput struct {
	EnterpriseProfile *entity.EnterpriseProfile
}

func (i *EnterpriseProfileInteractorImpl) GetEnterpriseByID(input GetEnterpriseByIDInput) (GetEnterpriseByIDOutput, error) {
	var (
		output GetEnterpriseByIDOutput
		err    error
	)

	enterpriseProfile, err := i.enterpriseProfileRepository.FindByID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseIndustry, err := i.enterpriseIndustryRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseReferenceMaterial, err := i.enterpriseReferenceMaterialRepository.FindByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseActivityList, err := i.enterpriseActivityRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, industry := range enterpriseIndustry {
		enterpriseProfile.Industries = append(enterpriseProfile.Industries, industry.Industry)
	}

	enterpriseProfile.ReferenceMaterialID = enterpriseReferenceMaterial.ID

	enterpriseProfile.ReferenceMaterial = entity.EnterpriseReferenceMaterial{
		EnterpriseID:     input.EnterpriseID,
		Reference1PDFURL: enterpriseReferenceMaterial.Reference1PDFURL,
		Reference2PDFURL: enterpriseReferenceMaterial.Reference2PDFURL,
	}

	for _, activity := range enterpriseActivityList {
		enterpriseProfile.Activities = append(enterpriseProfile.Activities, *activity)
	}

	output.EnterpriseProfile = enterpriseProfile

	return output, nil
}

// 担当者IDから企業情報一覧を取得する
type GetEnterpriseListByAgentStaffIDInput struct {
	AgentStaffID uint
}

type GetEnterpriseListByAgentStaffIDOutput struct {
	EnterpriseProfileList []*entity.EnterpriseProfile
}

func (i *EnterpriseProfileInteractorImpl) GetEnterpriseListByAgentStaffID(input GetEnterpriseListByAgentStaffIDInput) (GetEnterpriseListByAgentStaffIDOutput, error) {
	var (
		output GetEnterpriseListByAgentStaffIDOutput
		err    error
	)

	enterpriseProfileList, err := i.enterpriseProfileRepository.GetByAgentStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseProfileIDList := make([]uint, 0, len(enterpriseProfileList))
	for _, profile := range enterpriseProfileList {
		enterpriseProfileIDList = append(enterpriseProfileIDList, profile.ID)
	}

	enterpriseIndustryList, err := i.enterpriseIndustryRepository.GetByEnterpriseIDList(enterpriseProfileIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseReferenceMaterialList, err := i.enterpriseReferenceMaterialRepository.GetByEnterpriseIDList(enterpriseProfileIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseActivityList, err := i.enterpriseActivityRepository.GetByEnterpriseIDList(enterpriseProfileIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, profile := range enterpriseProfileList {
		for _, industry := range enterpriseIndustryList {
			if profile.ID == industry.EnterpriseID {
				profile.Industries = append(profile.Industries, industry.Industry)
			}
		}

		for _, referenceMaterial := range enterpriseReferenceMaterialList {
			if profile.ID == referenceMaterial.EnterpriseID {
				profile.ReferenceMaterial = entity.EnterpriseReferenceMaterial{
					Reference1PDFURL: referenceMaterial.Reference1PDFURL,
					Reference2PDFURL: referenceMaterial.Reference2PDFURL,
				}
			}
		}

		for _, activity := range enterpriseActivityList {
			if profile.ID == activity.EnterpriseID {
				profile.Activities = append(profile.Activities, *activity)
			}
		}
	}

	output.EnterpriseProfileList = enterpriseProfileList

	return output, nil
}

// エージェントIDから企業情報一覧を取得する
type GetEnterpriseListByAgentIDInput struct {
	AgentID uint
}

type GetEnterpriseListByAgentIDOutput struct {
	EnterpriseProfileList []*entity.EnterpriseProfile
}

func (i *EnterpriseProfileInteractorImpl) GetEnterpriseListByAgentID(input GetEnterpriseListByAgentIDInput) (GetEnterpriseListByAgentIDOutput, error) {
	var (
		output GetEnterpriseListByAgentIDOutput
		err    error
	)

	enterpriseProfileList, err := i.enterpriseProfileRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseProfileIDList := make([]uint, 0, len(enterpriseProfileList))
	for _, profile := range enterpriseProfileList {
		enterpriseProfileIDList = append(enterpriseProfileIDList, profile.ID)
	}

	enterpriseIndustryList, err := i.enterpriseIndustryRepository.GetByEnterpriseIDList(enterpriseProfileIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseReferenceMaterialList, err := i.enterpriseReferenceMaterialRepository.GetByEnterpriseIDList(enterpriseProfileIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseActivityList, err := i.enterpriseActivityRepository.GetByEnterpriseIDList(enterpriseProfileIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, profile := range enterpriseProfileList {

		for _, industry := range enterpriseIndustryList {
			if profile.ID == industry.EnterpriseID {
				profile.Industries = append(profile.Industries, industry.Industry)
			}
		}

		for _, referenceMaterial := range enterpriseReferenceMaterialList {
			if profile.ID == referenceMaterial.EnterpriseID {
				profile.ReferenceMaterial = entity.EnterpriseReferenceMaterial{
					Reference1PDFURL: referenceMaterial.Reference1PDFURL,
					Reference2PDFURL: referenceMaterial.Reference2PDFURL,
				}
			}
		}

		for _, activity := range enterpriseActivityList {
			if profile.ID == activity.EnterpriseID {
				profile.Activities = append(profile.Activities, *activity)
			}
		}

	}

	output.EnterpriseProfileList = enterpriseProfileList

	return output, nil
}

type GetEnterpriseListByAgentIDAndPageInput struct {
	AgentID    uint
	PageNumber uint
}

type GetEnterpriseListByAgentIDAndPageOutput struct {
	EnterpriseProfileList []*entity.EnterpriseProfile
	MaxPageNumber         uint
	IDList                []uint
}

func (i *EnterpriseProfileInteractorImpl) GetEnterpriseListByAgentIDAndPage(input GetEnterpriseListByAgentIDAndPageInput) (GetEnterpriseListByAgentIDAndPageOutput, error) {
	var (
		output GetEnterpriseListByAgentIDAndPageOutput
		err    error
	)

	enterpriseProfileList, err := i.enterpriseProfileRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseProfileIDList := make([]uint, 0, len(enterpriseProfileList))
	for _, profile := range enterpriseProfileList {
		enterpriseProfileIDList = append(enterpriseProfileIDList, profile.ID)
	}

	enterpriseIndustryList, err := i.enterpriseIndustryRepository.GetByEnterpriseIDList(enterpriseProfileIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseReferenceMaterialList, err := i.enterpriseReferenceMaterialRepository.GetByEnterpriseIDList(enterpriseProfileIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// IDListを返す
	for _, enterprise := range enterpriseProfileList {
		output.IDList = append(output.IDList, enterprise.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getEnterpriseListMaxPage(enterpriseProfileList)

	// 指定ページの企業20件を取得（本番実装までは1ページあたり5件）
	enterpriseProfileList20 := getEnterpriseListWithPage(enterpriseProfileList, input.PageNumber)

	for i := range enterpriseProfileList20 {
		for _, industry := range enterpriseIndustryList {
			if enterpriseProfileList20[i].ID == industry.EnterpriseID {
				enterpriseProfileList20[i].Industries = append(enterpriseProfileList20[i].Industries, industry.Industry)
			}
		}

		for _, file := range enterpriseReferenceMaterialList {
			if enterpriseProfileList20[i].ID == file.EnterpriseID {
				enterpriseProfileList20[i].ReferenceMaterial.Reference1PDFURL = file.Reference1PDFURL
				enterpriseProfileList20[i].ReferenceMaterial.Reference2PDFURL = file.Reference2PDFURL
			}
		}
	}

	output.EnterpriseProfileList = enterpriseProfileList20

	return output, nil
}

// エージェントIDから企業名一覧を取得する
type GetSearchEnterpriseListByAgentIDInput struct {
	AgentID     uint
	PageNumber  uint
	SearchParam entity.SearchEnterprise
}

type GetSearchEnterpriseListByAgentIDOutput struct {
	EnterpriseProfileList []*entity.EnterpriseProfile
	MaxPageNumber         uint
	IDList                []uint
}

func (i *EnterpriseProfileInteractorImpl) GetSearchEnterpriseListByAgentID(input GetSearchEnterpriseListByAgentIDInput) (GetSearchEnterpriseListByAgentIDOutput, error) {
	var (
		output GetSearchEnterpriseListByAgentIDOutput
		err    error
	)

	/**
	GetEnterpriseListByAgentIDAndFreeWordは
	フリーワードの有無で処理を分岐

	フリーワードは社名のみ
	*/
	enterpriseProfileList, err := i.enterpriseProfileRepository.GetByAgentIDAndFreeWord(input.AgentID, input.SearchParam.FreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseProfileIDList := make([]uint, 0, len(enterpriseProfileList))
	for _, profile := range enterpriseProfileList {
		enterpriseProfileIDList = append(enterpriseProfileIDList, profile.ID)
	}

	enterpriseIndustryList, err := i.enterpriseIndustryRepository.GetByEnterpriseIDList(enterpriseProfileIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseReferenceMaterialList, err := i.enterpriseReferenceMaterialRepository.GetByEnterpriseIDList(enterpriseProfileIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 縦持ちテーブルの処理
	for _, profile := range enterpriseProfileList {
		for _, industry := range enterpriseIndustryList {
			if profile.ID == industry.EnterpriseID {
				profile.Industries = append(profile.Industries, industry.Industry)
			}
		}

		for _, file := range enterpriseReferenceMaterialList {
			if profile.ID == file.EnterpriseID {
				profile.ReferenceMaterial.Reference1PDFURL = file.Reference1PDFURL
				profile.ReferenceMaterial.Reference2PDFURL = file.Reference2PDFURL
			}
		}
	}

	// 絞り込み項目の結果を代入するための変数を用意
	var (
		enterpriseProfileListWithAgentStaffID []*entity.EnterpriseProfile
		enterpriseProfileListWithIndustry     []*entity.EnterpriseProfile
		// enterpriseProfileListWithPrefecture   []*entity.EnterpriseProfile
		enterpriseProfileListWithCompanyScale []*entity.EnterpriseProfile
	)

	// 営業担当者IDがある場合
	agentStaffID, err := strconv.Atoi(input.SearchParam.AgentStaffID)
	if !(err != nil || agentStaffID == 0) {
		for _, enterprise := range enterpriseProfileList {
			if enterprise.AgentStaffID != uint(agentStaffID) {
				continue
			}
			enterpriseProfileListWithAgentStaffID = append(enterpriseProfileListWithAgentStaffID, enterprise)
		}
	}

	// 営業担当者IDが無い場合
	if err != nil || agentStaffID == 0 {
		enterpriseProfileListWithAgentStaffID = enterpriseProfileList
	}

	fmt.Println("RA担当者: ", enterpriseProfileListWithAgentStaffID)

	// NOTE: 業界
	// 業界のいずれかが入っている場合
	if !(len(input.SearchParam.Industries) == 0) {
	industryLoop:
		for _, enterprise := range enterpriseProfileListWithAgentStaffID {
			/**
			合致した時にtrueに変える
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			for _, industry := range input.SearchParam.Industries {
				if !industry.Valid {
					continue
				}
				if len(enterprise.Industries) == 0 {
					continue
				}

				for _, enterpriseIndustry := range enterprise.Industries {
					if industry == enterpriseIndustry {
						/**
						企業の業界が合致したタイミングでcontinueで次の企業へ
						ダブりが発生しないように制御
						*/
						enterpriseProfileListWithIndustry = append(enterpriseProfileListWithIndustry, enterprise)
						continue industryLoop
					}
				}
			}
		}
	}

	// 業界のいずれかも入っていない場合
	if len(input.SearchParam.Industries) == 0 {
		enterpriseProfileListWithIndustry = enterpriseProfileListWithAgentStaffID
	}

	fmt.Println("業界: ", enterpriseProfileListWithIndustry)

	// Note: 企業規模の絞り込みは企業の従業員数（単体）と比較する
	// 企業規模がある場合
	if !(len(input.SearchParam.CompanyScaleTypes) == 0) {
	companyScaleLoop:
		for _, enterprise := range enterpriseProfileListWithIndustry {
			if enterprise.EmployeeNumberSingle == null.NewInt(0, false) {
				continue
			}
			for _, companyScale := range input.SearchParam.CompanyScaleTypes {
				if !companyScale.Valid {
					continue
				}

				if companyScale == null.NewInt(0, true) {
					// 10名未満の場合
					if enterprise.EmployeeNumberSingle.Int64 < 10 {
						enterpriseProfileListWithCompanyScale = append(enterpriseProfileListWithCompanyScale, enterprise)
						continue companyScaleLoop
					}
				} else if companyScale == null.NewInt(1, true) {
					// 10名以上100名未満の場合
					if enterprise.EmployeeNumberSingle.Int64 >= 10 && enterprise.EmployeeNumberSingle.Int64 < 100 {
						enterpriseProfileListWithCompanyScale = append(enterpriseProfileListWithCompanyScale, enterprise)
						continue companyScaleLoop
					}
				} else if companyScale == null.NewInt(2, true) {
					// 100名以上200名未満の場合
					if enterprise.EmployeeNumberSingle.Int64 >= 100 && enterprise.EmployeeNumberSingle.Int64 < 200 {
						enterpriseProfileListWithCompanyScale = append(enterpriseProfileListWithCompanyScale, enterprise)
						continue companyScaleLoop
					}
				} else if companyScale == null.NewInt(3, true) {
					// 200名以上1000名未満の場合
					if enterprise.EmployeeNumberSingle.Int64 >= 200 && enterprise.EmployeeNumberSingle.Int64 < 1000 {
						enterpriseProfileListWithCompanyScale = append(enterpriseProfileListWithCompanyScale, enterprise)
						continue companyScaleLoop
					}
				} else if companyScale == null.NewInt(4, true) {
					// 1000名以上の場合
					if enterprise.EmployeeNumberSingle.Int64 >= 1000 {
						enterpriseProfileListWithCompanyScale = append(enterpriseProfileListWithCompanyScale, enterprise)
						continue companyScaleLoop
					}
				}
			}
		}
	}

	// 企業規模が無い場合
	if len(input.SearchParam.CompanyScaleTypes) == 0 {
		enterpriseProfileListWithCompanyScale = enterpriseProfileListWithIndustry
	}

	fmt.Println("企業規模: ", enterpriseProfileListWithCompanyScale)

	// IDListを返す
	for _, enterprise := range enterpriseProfileListWithCompanyScale {
		output.IDList = append(output.IDList, enterprise.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getEnterpriseListMaxPage(enterpriseProfileListWithCompanyScale)

	// 指定ページの企業20件を取得（本番実装までは1ページあたり5件）
	output.EnterpriseProfileList = getEnterpriseListWithPage(enterpriseProfileListWithCompanyScale, input.PageNumber)

	return output, nil
}

/****************************************************************************************/
// 企業資料関連
//
// 求人企業資料の作成
type CreateEnterpriseReferenceMaterialInput struct {
	CreateParam entity.CreateOrUpdateEnterpriseReferenceMaterialParam
}

type CreateEnterpriseReferenceMaterialOutput struct {
	EnterpriseReferenceMaterial *entity.EnterpriseReferenceMaterial
}

func (i *EnterpriseProfileInteractorImpl) CreateEnterpriseReferenceMaterial(input CreateEnterpriseReferenceMaterialInput) (CreateEnterpriseReferenceMaterialOutput, error) {
	var (
		output CreateEnterpriseReferenceMaterialOutput
		err    error
	)

	enterpriseReferenceMaterial := entity.NewEnterpriseReferenceMaterial(
		input.CreateParam.EnterpriseID,
		input.CreateParam.Reference1PDFURL,
		input.CreateParam.Reference2PDFURL,
	)

	err = i.enterpriseReferenceMaterialRepository.Create(enterpriseReferenceMaterial)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.EnterpriseReferenceMaterial = enterpriseReferenceMaterial

	return output, nil
}

// 求人企業資料の更新
type UpdateEnterpriseReferenceMaterialInput struct {
	UpdateParam entity.CreateOrUpdateEnterpriseReferenceMaterialParam
}

type UpdateEnterpriseReferenceMaterialOutput struct {
	EnterpriseReferenceMaterial *entity.EnterpriseReferenceMaterial
}

func (i *EnterpriseProfileInteractorImpl) UpdateEnterpriseReferenceMaterial(input UpdateEnterpriseReferenceMaterialInput) (UpdateEnterpriseReferenceMaterialOutput, error) {
	var (
		output UpdateEnterpriseReferenceMaterialOutput
		err    error
	)

	enterpriseReferenceMaterial := entity.NewEnterpriseReferenceMaterial(
		input.UpdateParam.EnterpriseID,
		input.UpdateParam.Reference1PDFURL,
		input.UpdateParam.Reference2PDFURL,
	)

	err = i.enterpriseReferenceMaterialRepository.UpdateByEnterpriseID(input.UpdateParam.EnterpriseID, enterpriseReferenceMaterial)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.EnterpriseReferenceMaterial = enterpriseReferenceMaterial

	return output, nil
}

// 求人企業資料の取得
type GetEnterpriseReferenceMaterialByEnterpriseIDInput struct {
	EnterpriseID uint
}

type GetEnterpriseReferenceMaterialByEnterpriseIDOutput struct {
	EnterpriseReferenceMaterial *entity.EnterpriseReferenceMaterial
}

func (i *EnterpriseProfileInteractorImpl) GetEnterpriseReferenceMaterialByEnterpriseID(input GetEnterpriseReferenceMaterialByEnterpriseIDInput) (GetEnterpriseReferenceMaterialByEnterpriseIDOutput, error) {
	var (
		output GetEnterpriseReferenceMaterialByEnterpriseIDOutput
		err    error
	)

	enterpriseReferenceMaterial, err := i.enterpriseReferenceMaterialRepository.FindByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.EnterpriseReferenceMaterial = enterpriseReferenceMaterial

	return output, nil
}

/****************************************************************************************/
/// 求人企業の追加情報 API
//
// 求人企業の追加情報を作成
type CreateEnterpriseActivityInput struct {
	CreateParam entity.CreateEnterpriseActivityParam
}

type CreateEnterpriseActivityOutput struct {
	OK bool
}

func (i *EnterpriseProfileInteractorImpl) CreateEnterpriseActivity(input CreateEnterpriseActivityInput) (CreateEnterpriseActivityOutput, error) {
	var (
		output CreateEnterpriseActivityOutput
		err    error
	)

	enterpriseActivity := entity.NewEnterpriseActivity(
		input.CreateParam.EnterpriseID,
		input.CreateParam.Content,
		input.CreateParam.AddedAt,
	)

	err = i.enterpriseActivityRepository.Create(enterpriseActivity)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/
/****************************************************************************************/
/// Admin API
//

// すべての企業情報を取得
type GetInitialEnterpriseInput struct {
	PageNumber uint
}

type GetInitialEnterpriseOutput struct {
	EnterpriseProfileList []*entity.EnterpriseProfile
	MaxPageNumber         uint
	IDList                []uint
}

func (i *EnterpriseProfileInteractorImpl) GetInitialEnterprise(input GetInitialEnterpriseInput) (GetInitialEnterpriseOutput, error) {
	var (
		output GetInitialEnterpriseOutput
		err    error
	)

	enterpriseProfileList, err := i.enterpriseProfileRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseIndustryList, err := i.enterpriseIndustryRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseReferenceMaterialList, err := i.enterpriseReferenceMaterialRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// ページの最大数を取得
	output.MaxPageNumber = getEnterpriseListMaxPage(enterpriseProfileList)

	// 指定ページの企業50件を取得（本番実装までは1ページあたり5件）
	enterpriseProfileList50 := getEnterpriseListWithPage(enterpriseProfileList, input.PageNumber)

	for i := range enterpriseProfileList50 {
		for _, industry := range enterpriseIndustryList {
			if enterpriseProfileList50[i].ID == industry.EnterpriseID {
				enterpriseProfileList50[i].Industries = append(enterpriseProfileList50[i].Industries, industry.Industry)
			}
		}

		for _, file := range enterpriseReferenceMaterialList {
			if enterpriseProfileList50[i].ID == file.EnterpriseID {
				enterpriseProfileList50[i].ReferenceMaterial.Reference1PDFURL = file.Reference1PDFURL
				enterpriseProfileList50[i].ReferenceMaterial.Reference2PDFURL = file.Reference2PDFURL
			}
		}
	}

	for _, profile := range enterpriseProfileList {
		for _, industry := range enterpriseIndustryList {
			if profile.ID == industry.EnterpriseID {
				profile.Industries = append(profile.Industries, industry.Industry)
			}
		}

		for _, referenceMaterial := range enterpriseReferenceMaterialList {
			if profile.ID == referenceMaterial.EnterpriseID {
				profile.ReferenceMaterial = entity.EnterpriseReferenceMaterial{
					Reference1PDFURL: referenceMaterial.Reference1PDFURL,
					Reference2PDFURL: referenceMaterial.Reference2PDFURL,
				}
			}
		}
	}

	// IDListを返す
	for _, enterprise := range enterpriseProfileList {
		output.IDList = append(output.IDList, enterprise.ID)
	}

	output.EnterpriseProfileList = enterpriseProfileList50

	return output, nil
}

/****************************************************************************************/
