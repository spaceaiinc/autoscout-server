package handler

import (
	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SendingEnterpriseHandler interface {
	// 汎用系 API
	CreateSendingEnterprise(param entity.CreateOrUpdateSendingEnterpriseParam) (presenter.Presenter, error)
	UpdateSendingEnterprise(param entity.CreateOrUpdateSendingEnterpriseParam, id uint) (presenter.Presenter, error)
	UpdateSendingEnterprisePassword(id uint, password string) (presenter.Presenter, error)
	SigninSendingEnterprise(param entity.SigninSendingEnterprisePasswordParam) (presenter.Presenter, error)
	DeleteSendingEnterprise(id uint) (presenter.Presenter, error)
	GetSendingEnterpriseByID(id uint) (presenter.Presenter, error)
	GetSendingEnterpriseAndBillingAddressByID(id uint) (presenter.Presenter, error)
	GetSendingEnterpriseByUUID(uuid uuid.UUID) (presenter.Presenter, error)

	// ページネーション
	GetAllSendingEnterpriseByPageAndFreeWord(pageNumber uint, freeWord string) (presenter.Presenter, error)

	// 固有ページ用
	GetSendingInformationForSendingMail(sendingJobInformationIDList []uint) (presenter.Presenter, error)
	SendSendingMail(param entity.SendSendingMailParam) (presenter.Presenter, error)
	SendMailForRSVP(param entity.SendSendingMailForRSVPParam) (presenter.Presenter, error)
	GetSendingEnterpriseAndAcceptJobSeekerListByAgentID(agentID uint) (presenter.Presenter, error)

	// 企業資料系 API
	CreateSendingEnterpriseReferenceMaterial(param entity.CreateOrUpdateSendingEnterpriseReferenceMaterialParam) (presenter.Presenter, error)
	UpdateSendingEnterpriseReferenceMaterial(param entity.CreateOrUpdateSendingEnterpriseReferenceMaterialParam) (presenter.Presenter, error)
	DeleteSendingEnterpriseReferenceMaterial(sendingEnterpriseID uint, fileType string) (presenter.Presenter, error)
	GetSendingEnterpriseReferenceMaterialBySendingEnterpriseID(id uint) (presenter.Presenter, error)
}

type SendingEnterpriseHandlerImpl struct {
	sendingEnterpriseInteractor interactor.SendingEnterpriseInteractor
}

func NewSendingEnterpriseHandlerImpl(epI interactor.SendingEnterpriseInteractor) SendingEnterpriseHandler {
	return &SendingEnterpriseHandlerImpl{
		sendingEnterpriseInteractor: epI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *SendingEnterpriseHandlerImpl) CreateSendingEnterprise(param entity.CreateOrUpdateSendingEnterpriseParam) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.CreateSendingEnterprise(interactor.CreateSendingEnterpriseInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseJSONPresenter(responses.NewSendingEnterprise(output.SendingEnterprise)), nil
}

func (h *SendingEnterpriseHandlerImpl) UpdateSendingEnterprise(param entity.CreateOrUpdateSendingEnterpriseParam, id uint) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.UpdateSendingEnterprise(interactor.UpdateSendingEnterpriseInput{
		UpdateParam:         param,
		SendingEnterpriseID: id,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseJSONPresenter(responses.NewSendingEnterprise(output.SendingEnterprise)), nil
}

func (h *SendingEnterpriseHandlerImpl) UpdateSendingEnterprisePassword(id uint, password string) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.UpdateSendingEnterprisePassword(interactor.UpdateSendingEnterprisePasswordInput{
		SendingEnterpriseID: id,
		Password:            password,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingEnterpriseHandlerImpl) SigninSendingEnterprise(param entity.SigninSendingEnterprisePasswordParam) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.SigninSendingEnterprise(interactor.SigninSendingEnterpriseInput{
		SigninParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseForSigninJSONPresenter(responses.NewSendingEnterpriseForSignin(output.SendingEnterprise, output.SendingJobSeeker, output.SendingShareDocument)), nil
}

func (h *SendingEnterpriseHandlerImpl) DeleteSendingEnterprise(id uint) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.DeleteSendingEnterprise(interactor.DeleteSendingEnterpriseInput{
		SendingEnterpriseID: id,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 企業IDから企業情報を取得する
func (h *SendingEnterpriseHandlerImpl) GetSendingEnterpriseByID(id uint) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.GetSendingEnterpriseByID(interactor.GetSendingEnterpriseByIDInput{
		SendingEnterpriseID: id,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseJSONPresenter(responses.NewSendingEnterprise(output.SendingEnterprise)), nil
}

// 企業IDから企業情報と請求先情報を取得する
func (h *SendingEnterpriseHandlerImpl) GetSendingEnterpriseAndBillingAddressByID(id uint) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.GetSendingEnterpriseAndBillingAddressByID(interactor.GetSendingEnterpriseAndBillingAddressByIDInput{
		SendingEnterpriseID: id,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseJSONPresenter(responses.NewSendingEnterprise(output.SendingEnterprise)), nil
}

// 企業UUIDから企業情報を取得する
func (h *SendingEnterpriseHandlerImpl) GetSendingEnterpriseByUUID(uuid uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.GetSendingEnterpriseByUUID(interactor.GetSendingEnterpriseByUUIDInput{
		SendingEnterpriseUUID: uuid,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseJSONPresenter(responses.NewSendingEnterprise(output.SendingEnterprise)), nil
}

/****************************************************************************************/
/// ページネーション
//
// すべての企業情報を取得する
func (h *SendingEnterpriseHandlerImpl) GetAllSendingEnterpriseByPageAndFreeWord(pageNumber uint, freeWord string) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.GetAllSendingEnterpriseByPageAndFreeWord(interactor.GetAllSendingEnterpriseByPageAndFreeWordInput{
		PageNumber: pageNumber,
		FreeWord:   freeWord,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterprisePageListAndMaxPageAndIDListJSONPresenter(responses.NewSendingEnterpriseListAndMaxPageAndIDList(output.SendingEnterpriseList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/
/// 企業資料関連
//
func (h *SendingEnterpriseHandlerImpl) CreateSendingEnterpriseReferenceMaterial(param entity.CreateOrUpdateSendingEnterpriseReferenceMaterialParam) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.CreateSendingEnterpriseReferenceMaterial(interactor.CreateSendingEnterpriseReferenceMaterialInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseReferenceMaterialJSONPresenter(responses.NewSendingEnterpriseReferenceMaterial(output.SendingEnterpriseReferenceMaterial)), nil
}

func (h *SendingEnterpriseHandlerImpl) UpdateSendingEnterpriseReferenceMaterial(param entity.CreateOrUpdateSendingEnterpriseReferenceMaterialParam) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.UpdateSendingEnterpriseReferenceMaterial(interactor.UpdateSendingEnterpriseReferenceMaterialInput{
		UpdateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseReferenceMaterialJSONPresenter(responses.NewSendingEnterpriseReferenceMaterial(output.SendingEnterpriseReferenceMaterial)), nil
}

func (h *SendingEnterpriseHandlerImpl) DeleteSendingEnterpriseReferenceMaterial(sendingEntepriseID uint, fileType string) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.DeleteSendingEnterpriseReferenceMaterial(interactor.DeleteSendingEnterpriseReferenceMaterialInput{
		SendingEnterpriseID: sendingEntepriseID,
		FileType:            fileType,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 企業IDから企業資料を取得する
func (h *SendingEnterpriseHandlerImpl) GetSendingEnterpriseReferenceMaterialBySendingEnterpriseID(id uint) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.GetSendingEnterpriseReferenceMaterialBySendingEnterpriseID(interactor.GetSendingEnterpriseReferenceMaterialBySendingEnterpriseIDInput{
		SendingEnterpriseID: id,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseReferenceMaterialJSONPresenter(responses.NewSendingEnterpriseReferenceMaterial(output.SendingEnterpriseReferenceMaterial)), nil
}

/****************************************************************************************/
// 固有ページ用
//
func (h *SendingEnterpriseHandlerImpl) GetSendingInformationForSendingMail(sendingJobInformationIDList []uint) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.GetSendingInformationForSendingMail(interactor.GetSendingInformationForSendingMailInput{
		SendingJobInformationIDList: sendingJobInformationIDList,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseListJSONPresenter(responses.NewSendingEnterpriseList(output.SendingEnterpriseList)), nil
}

func (h *SendingEnterpriseHandlerImpl) SendSendingMail(param entity.SendSendingMailParam) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.SendSendingMail(interactor.SendSendingMailInput{
		SendParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingEnterpriseHandlerImpl) SendMailForRSVP(param entity.SendSendingMailForRSVPParam) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.SendMailForRSVP(interactor.SendMailForRSVPInput{
		SendParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingEnterpriseHandlerImpl) GetSendingEnterpriseAndAcceptJobSeekerListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.sendingEnterpriseInteractor.GetSendingEnterpriseAndAcceptJobSeekerListByAgentID(interactor.GetSendingEnterpriseAndAcceptJobSeekerListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseListJSONPresenter(responses.NewSendingEnterpriseList(output.SendingEnterpriseList)), nil
}
