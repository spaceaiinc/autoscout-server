package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type EnterpriseProfileHandler interface {
	// 汎用系 API
	CreateEnterpriseProfile(param entity.CreateOrUpdateEnterpriseProfileParam) (presenter.Presenter, error)
	UpdateEnterpriseProfile(param entity.CreateOrUpdateEnterpriseProfileParam, enterpriseID uint) (presenter.Presenter, error)
	DeleteEnterpriseProfile(param entity.DeleteEnterpriseProfileParam) (presenter.Presenter, error)
	GetEnterpriseByID(enterpriseID uint) (presenter.Presenter, error)
	GetEnterpriseListByAgentID(agentID uint) (presenter.Presenter, error)
	GetEnterpriseListByAgentStaffID(agentStaffID uint) (presenter.Presenter, error)
	GetSearchEnterpriseListByAgentID(agentID, pageNumber uint, searchParam entity.SearchEnterprise) (presenter.Presenter, error)
	GetEnterpriseListByAgentIDAndPage(agentID, pageNumber uint) (presenter.Presenter, error) //エージェントIDとページ番号で企業一覧を取得する
	DeleteEnterpriseReferenceMaterial(referenceMaterialID, materialType uint) (presenter.Presenter, error)

	// 企業資料系 API
	CreateEnterpriseReferenceMaterial(param entity.CreateOrUpdateEnterpriseReferenceMaterialParam) (presenter.Presenter, error)
	UpdateEnterpriseReferenceMaterial(param entity.CreateOrUpdateEnterpriseReferenceMaterialParam) (presenter.Presenter, error)
	GetEnterpriseReferenceMaterialByEnterpriseID(enterpriseID uint) (presenter.Presenter, error)

	// CSV操作 API
	ImportEnterpriseCSV(param []*entity.EnterpriseAndBillingAddress, missedRecords []uint, agentID uint) (presenter.Presenter, error)
	ImportJobInformationCSV(param []*entity.JobInformation, missedRecords []uint, agentID uint) (presenter.Presenter, error)
	ImportEnterpriseCSVForCircus(param []*entity.EnterpriseAndJobInformation, missedRecords []uint, agentID uint) (presenter.Presenter, error)
	ImportEnterpriseCSVForAgentBank(param []*entity.EnterpriseAndJobInformation, missedRecords []uint, agentID uint) (presenter.Presenter, error)
	ExportEnterpriseCSV(agentID uint) (string, error)

	// 企業と求人のリストから作成
	ImportEnterpriseJSON(param []entity.EnterpriseAndJobInformation, agentStaffID uint) (presenter.Presenter, error)

	// 企業の追加情報
	CreateEnterpriseActivity(param entity.CreateEnterpriseActivityParam) (presenter.Presenter, error)

	// Admin API
	GetInitialEnterprise(pageNumber uint) (presenter.Presenter, error)
}

type EnterpriseProfileHandlerImpl struct {
	enterpriseProfileInteractor interactor.EnterpriseProfileInteractor
}

func NewEnterpriseProfileHandlerImpl(epI interactor.EnterpriseProfileInteractor) EnterpriseProfileHandler {
	return &EnterpriseProfileHandlerImpl{
		enterpriseProfileInteractor: epI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *EnterpriseProfileHandlerImpl) CreateEnterpriseProfile(param entity.CreateOrUpdateEnterpriseProfileParam) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.CreateEnterpriseProfile(interactor.CreateEnterpriseProfileInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseProfileJSONPresenter(responses.NewEnterpriseProfile(output.EnterpriseProfile)), nil
}

func (h *EnterpriseProfileHandlerImpl) UpdateEnterpriseProfile(param entity.CreateOrUpdateEnterpriseProfileParam, enterpriseID uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.UpdateEnterpriseProfile(interactor.UpdateEnterpriseProfileInput{
		UpdateParam:  param,
		EnterpriseID: enterpriseID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseProfileJSONPresenter(responses.NewEnterpriseProfile(output.EnterpriseProfile)), nil
}

func (h *EnterpriseProfileHandlerImpl) DeleteEnterpriseProfile(param entity.DeleteEnterpriseProfileParam) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.DeleteEnterpriseProfile(interactor.DeleteEnterpriseProfileInput{
		DeleteParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *EnterpriseProfileHandlerImpl) DeleteEnterpriseReferenceMaterial(referenceMaterialID, materialType uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.DeleteEnterpriseReferenceMaterial(interactor.DeleteEnterpriseReferenceMaterialInput{
		ReferenceMaterialID: referenceMaterialID,
		MaterialType:        materialType,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 企業IDから企業情報を取得する
func (h *EnterpriseProfileHandlerImpl) GetEnterpriseByID(enterpriseID uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.GetEnterpriseByID(interactor.GetEnterpriseByIDInput{
		EnterpriseID: enterpriseID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseProfileJSONPresenter(responses.NewEnterpriseProfile(output.EnterpriseProfile)), nil
}

// 担当者IDから企業情報一覧を取得する
func (h *EnterpriseProfileHandlerImpl) GetEnterpriseListByAgentStaffID(agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.GetEnterpriseListByAgentStaffID(interactor.GetEnterpriseListByAgentStaffIDInput{
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseProfileListJSONPresenter(responses.NewEnterpriseProfileList(output.EnterpriseProfileList)), nil
}

// エージェントIDから企業一覧を取得する
func (h *EnterpriseProfileHandlerImpl) GetEnterpriseListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.GetEnterpriseListByAgentID(interactor.GetEnterpriseListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseProfileListJSONPresenter(responses.NewEnterpriseProfileList(output.EnterpriseProfileList)), nil
}

// エージェントIDとクエリパラムで企業一覧を絞り込み
func (h *EnterpriseProfileHandlerImpl) GetSearchEnterpriseListByAgentID(agentID, pageNumber uint, searchParam entity.SearchEnterprise) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.GetSearchEnterpriseListByAgentID(interactor.GetSearchEnterpriseListByAgentIDInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseProfilePageListAndMaxPageAndIDListJSONPresenter(responses.NewEnterpriseProfileListAndMaxPageAndIDList(output.EnterpriseProfileList, output.MaxPageNumber, output.IDList)), nil
}

// エージェントIDとページ番号で企業一覧を取得する
func (h *EnterpriseProfileHandlerImpl) GetEnterpriseListByAgentIDAndPage(agentID, pageNumber uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.GetEnterpriseListByAgentIDAndPage(interactor.GetEnterpriseListByAgentIDAndPageInput{
		AgentID:    agentID,
		PageNumber: pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseProfilePageListAndMaxPageAndIDListJSONPresenter(responses.NewEnterpriseProfileListAndMaxPageAndIDList(output.EnterpriseProfileList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/
/// 企業資料関連
//
func (h *EnterpriseProfileHandlerImpl) CreateEnterpriseReferenceMaterial(param entity.CreateOrUpdateEnterpriseReferenceMaterialParam) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.CreateEnterpriseReferenceMaterial(interactor.CreateEnterpriseReferenceMaterialInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseReferenceMaterialJSONPresenter(responses.NewEnterpriseReferenceMaterial(output.EnterpriseReferenceMaterial)), nil
}

func (h *EnterpriseProfileHandlerImpl) UpdateEnterpriseReferenceMaterial(param entity.CreateOrUpdateEnterpriseReferenceMaterialParam) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.UpdateEnterpriseReferenceMaterial(interactor.UpdateEnterpriseReferenceMaterialInput{
		UpdateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseReferenceMaterialJSONPresenter(responses.NewEnterpriseReferenceMaterial(output.EnterpriseReferenceMaterial)), nil
}

// 企業IDから企業資料を取得する
func (h *EnterpriseProfileHandlerImpl) GetEnterpriseReferenceMaterialByEnterpriseID(enterpriseID uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.GetEnterpriseReferenceMaterialByEnterpriseID(interactor.GetEnterpriseReferenceMaterialByEnterpriseIDInput{
		EnterpriseID: enterpriseID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseReferenceMaterialJSONPresenter(responses.NewEnterpriseReferenceMaterial(output.EnterpriseReferenceMaterial)), nil
}

/****************************************************************************************/
/****************************************************************************************/
/// 企業の追加情報 API
func (h *EnterpriseProfileHandlerImpl) CreateEnterpriseActivity(param entity.CreateEnterpriseActivityParam) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.CreateEnterpriseActivity(interactor.CreateEnterpriseActivityInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
/****************************************************************************************/
/// CSV操作 API
//csvファイルを読み込む 企業・請求先
func (h *EnterpriseProfileHandlerImpl) ImportEnterpriseCSV(param []*entity.EnterpriseAndBillingAddress, missedRecords []uint, agentID uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.ImportEnterpriseCSV(interactor.ImportEnterpriseCSVInput{
		CreateParam:   param,
		MissedRecords: missedRecords,
		AgentID:       agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMissedRecordsAndOKJSONPresenter(responses.NewMissedRecordsAndOK(output.MissedRecords, output.OK)), nil
}

// csvファイルを読み込む 求人
func (h *EnterpriseProfileHandlerImpl) ImportJobInformationCSV(param []*entity.JobInformation, missedRecords []uint, agentID uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.ImportJobInformationCSV(interactor.ImportJobInformationCSVInput{
		CreateParam:   param,
		MissedRecords: missedRecords,
		AgentID:       agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMissedRecordsAndOKJSONPresenter(responses.NewMissedRecordsAndOK(output.MissedRecords, output.OK)), nil
}

func (h *EnterpriseProfileHandlerImpl) ImportEnterpriseCSVForCircus(param []*entity.EnterpriseAndJobInformation, missedRecords []uint, agentID uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.ImportEnterpriseCSVForCircus(interactor.ImportEnterpriseCSVForCircusInput{
		CreateParam:   param,
		MissedRecords: missedRecords,
		AgentID:       agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMissedRecordsAndOKJSONPresenter(responses.NewMissedRecordsAndOK(output.MissedRecords, output.OK)), nil
}

func (h *EnterpriseProfileHandlerImpl) ImportEnterpriseCSVForAgentBank(param []*entity.EnterpriseAndJobInformation, missedRecords []uint, agentID uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.ImportEnterpriseCSVForAgentBank(interactor.ImportEnterpriseCSVForAgentBankInput{
		CreateParam:   param,
		MissedRecords: missedRecords,
		AgentID:       agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMissedRecordsAndOKJSONPresenter(responses.NewMissedRecordsAndOK(output.MissedRecords, output.OK)), nil
}

// csvファイルを出力する
func (h *EnterpriseProfileHandlerImpl) ExportEnterpriseCSV(agentID uint) (string, error) {
	output, err := h.enterpriseProfileInteractor.ExportEnterpriseCSV(interactor.ExportEnterpriseCSVInput{
		AgentID: agentID,
	})

	if err != nil {
		return "", err
	}

	return output.FilePath.PathName, nil
}

/****************************************************************************************/
/****************************************************************************************/
/// Admin API

// すべての企業情報を取得する
func (h *EnterpriseProfileHandlerImpl) GetInitialEnterprise(pageNumber uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.GetInitialEnterprise(interactor.GetInitialEnterpriseInput{
		PageNumber: pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEnterpriseProfilePageListAndMaxPageAndIDListJSONPresenter(responses.NewEnterpriseProfileListAndMaxPageAndIDList(output.EnterpriseProfileList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/
/****************************************************************************************/
// 求人取得 API
//
// 企業と求人のリストから作成
func (h *EnterpriseProfileHandlerImpl) ImportEnterpriseJSON(param []entity.EnterpriseAndJobInformation, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.enterpriseProfileInteractor.ImportEnterpriseJSON(interactor.ImportEnterpriseJSONInput{
		CreateParam:  param,
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
