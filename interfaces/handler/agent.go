package handler

import (
	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type AgentHandler interface {
	// 汎用系 API
	GetAgentByAgentID(agentID uint) (presenter.Presenter, error)
	GetAgentByAgentUUID(uuid uuid.UUID) (presenter.Presenter, error)
	UpdateAgent(agentID uint, param entity.CreateOrUpdateAgentParam) (presenter.Presenter, error)
	UpdateAgentAgreementFileURL(param entity.AgentAgreementFileURLParam) (presenter.Presenter, error)
	UpdateAgentForAdmin(param entity.AgentForAdminParam) (presenter.Presenter, error)
	GetAllianceAgentListByAgentID(agentID uint) (presenter.Presenter, error)
	GetAllianceAgentListByAgentIDForSelect(agentID uint) (presenter.Presenter, error)
	GetAgreementFileURL(firebaseToken string) (presenter.Presenter, error)
	GetAgentClaimList(firebaseToken string, pageNumber uint) (presenter.Presenter, error)

	// Admin API
	AgentSignUp(param entity.CreateOrUpdateAgentParam) (presenter.Presenter, error)
	AgentAndAgentStaffSignUp(param entity.CreateAgentAndAgentStaffParam) (presenter.Presenter, error)
	GetAllAgentList() (presenter.Presenter, error)

	// Agent API
	// CreateAgentLineInformation(param entity.CreateOrUpdateAgentLineInformation) (presenter.Presenter, error)

	// LINE 関連
	UpdateAgentLineChannel(param entity.AgentLineChannelParam) (presenter.Presenter, error)
	GetAgentLineChannelByAgentID(agentID uint) (presenter.Presenter, error)
	GetAgentLineLoginChannelIDByAgentUUID(agentUUID uuid.UUID) (presenter.Presenter, error)
}

type AgentHandlerImpl struct {
	AgentInteractor interactor.AgentInteractor
}

func NewAgentHandlerImpl(acaI interactor.AgentInteractor) AgentHandler {
	return &AgentHandlerImpl{
		AgentInteractor: acaI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *AgentHandlerImpl) GetAgentByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.GetAgentByAgentID(interactor.GetAgentByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentJSONPresenter(responses.NewAgent(output.Agent)), nil
}

func (h *AgentHandlerImpl) GetAgentByAgentUUID(uuid uuid.UUID) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.GetAgentByAgentUUID(interactor.GetAgentByAgentUUIDInput{
		UUID: uuid,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentJSONPresenter(responses.NewAgent(output.Agent)), nil
}

func (h *AgentHandlerImpl) UpdateAgent(agentID uint, param entity.CreateOrUpdateAgentParam) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.UpdateAgent(interactor.UpdateAgentInput{
		AgentID:     agentID,
		UpdateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentJSONPresenter(responses.NewAgent(output.Agent)), nil
}

func (h *AgentHandlerImpl) UpdateAgentAgreementFileURL(param entity.AgentAgreementFileURLParam) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.UpdateAgentAgreementFileURL(interactor.UpdateAgentAgreementFileURLInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *AgentHandlerImpl) UpdateAgentForAdmin(param entity.AgentForAdminParam) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.UpdateAgentForAdmin(interactor.UpdateAgentForAdminInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *AgentHandlerImpl) GetAllianceAgentListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.GetAllianceAgentListByAgentID(interactor.GetAllianceAgentListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentListJSONPresenter(responses.NewAgentList(output.AgentList)), nil
}

func (h *AgentHandlerImpl) GetAllianceAgentListByAgentIDForSelect(agentID uint) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.GetAllianceAgentListByAgentIDForSelect(interactor.GetAllianceAgentListByAgentIDForSelectInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentListJSONPresenter(responses.NewAgentList(output.AgentList)), nil
}

func (h *AgentHandlerImpl) GetAgreementFileURL(firebaseToken string) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.GetAgreementFileURL(interactor.GetAgreementFileURLInput{
		Token: firebaseToken,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentJSONPresenter(responses.NewAgent(output.Agent)), nil
}

func (h *AgentHandlerImpl) GetAgentClaimList(firebaseToken string, pageNumber uint) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.GetAgentClaimList(interactor.GetAgentClaimListInput{
		Token:      firebaseToken,
		PageNumber: pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentListAndMaxPageAndIDListJSONPresenter(responses.NewAgentListAndMaxPageAndIDList(output.AgentList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API
//
func (h *AgentHandlerImpl) AgentSignUp(param entity.CreateOrUpdateAgentParam) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.AgentSignUp(interactor.AgentSignUpInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentJSONPresenter(responses.NewAgent(output.Agent)), nil
}

func (h *AgentHandlerImpl) AgentAndAgentStaffSignUp(param entity.CreateAgentAndAgentStaffParam) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.AgentAndAgentStaffSignUp(interactor.AgentAndAgentStaffSignUpInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentJSONPresenter(responses.NewAgent(output.Agent)), nil
}

func (h *AgentHandlerImpl) GetAllAgentList() (presenter.Presenter, error) {
	output, err := h.AgentInteractor.GetAllAgentList()

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentListJSONPresenter(responses.NewAgentList(output.AgentList)), nil
}

/****************************************************************************************/

func (h *AgentHandlerImpl) UpdateAgentLineChannel(param entity.AgentLineChannelParam) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.UpdateAgentLineChannel(interactor.UpdateAgentLineChannelInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentLineJSONPresenter(responses.NewAgentLine(output.AgentLine)), nil
}

func (h *AgentHandlerImpl) GetAgentLineChannelByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.GetAgentLineChannelByAgentID(interactor.GetAgentLineChannelByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentLineJSONPresenter(responses.NewAgentLine(output.AgentLine)), nil
}

func (h *AgentHandlerImpl) GetAgentLineLoginChannelIDByAgentUUID(agentUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.AgentInteractor.GetAgentLineLoginChannelIDByAgentUUID(interactor.GetAgentLineLoginChannelIDByAgentUUIDInput{
		AgentUUID: agentUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentLineLoginChannelIDJSONPresenter(responses.NewAgentLineLoginChannelID(output.LineLoginChannelID)), nil
}
