package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type AgentStaffHandler interface {
	// 汎用系 API
	GetAgentStaffListByAgentID(token string, agentID uint) (presenter.Presenter, error)
	GetAgentStaffListByAgentIDOrderByID(agentID, agentStaffID uint) (presenter.Presenter, error)
	UpdateAgentStaff(agentStaffID uint, param entity.AgentStaffUpdateParam) (presenter.Presenter, error)
	UpdateAgentStaffEmail(param entity.AgentStaffEmailUpdateParam, agentStaffID uint) (presenter.Presenter, error)
	UpdateAgentStaffPassword(param entity.AgentStaffPasswordUpdateParam, agentStaffID uint) (presenter.Presenter, error)
	UpdateAgentStaffUsageEnd(param entity.DeleteAgentStaffParam) (presenter.Presenter, error)
	UpdateAgentStaffUsageReStart(agentStaffID uint) (presenter.Presenter, error)
	UpdateAgentStaffNotificationJobSeeker(param entity.UpdateAgentStaffNotificationJobSeekerParam) (presenter.Presenter, error)
	UpdateAgentStaffNotificationUnwatched(param entity.UpdateAgentStaffNotificationUnwatchedParam) (presenter.Presenter, error)
	UpdateAgentStaffAuthority(param entity.UpdateAgentStaffAuthorityParam) (presenter.Presenter, error)
	DeleteAgentStaff(param entity.DeleteAgentStaffParam) (presenter.Presenter, error) // 担当者削除　*firebaseのアイパスを削除&DBのis_deletedをtrueにする
	GetOhterAgentStaffListByAgentIDAndAllianceAgentID(agentID, allianceAgentID, agentStaffID uint) (presenter.Presenter, error)
	GetAgentStaffListWithSaleNotCreated(token string, agentID, managementID uint) (presenter.Presenter, error)
	GetAgentStaffListByAgentIDAndUsageStatusAvailable(token string, agentID uint) (presenter.Presenter, error) // 利用可能の担当者を取得
	GetAgentStaffListByAgentIDAndIsDeletedFalse(token string, agentID uint) (presenter.Presenter, error)       // 未削除の担当者を取得

	// Admin API
	SignUpForAdmin(param entity.AgentStaffSignUpForAdminParam, agentID uint) (presenter.Presenter, error)
	GetAllAgentStaffList() (presenter.Presenter, error) // 全担当者一覧を取得

	// Agent API
	GetAgentStaffMe(token string) (presenter.Presenter, error)
}

type AgentStaffHandlerImpl struct {
	AgentStaffInteractor interactor.AgentStaffInteractor
}

func NewAgentStaffHandlerImpl(acuI interactor.AgentStaffInteractor) AgentStaffHandler {
	return &AgentStaffHandlerImpl{
		AgentStaffInteractor: acuI,
	}
}

/****************************************************************************************/
/// 汎用系 API
func (h *AgentStaffHandlerImpl) GetAgentStaffListByAgentID(token string, agentID uint) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.GetAgentStaffListByAgentID(interactor.GetAgentStaffListByAgentIDInput{
		Token:   token,
		AgentID: agentID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffListJSONPresenter(responses.NewAgentStaffList(output.AgentStaffList)), nil
}

func (h *AgentStaffHandlerImpl) GetAgentStaffListByAgentIDOrderByID(agentID, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.GetAgentStaffListByAgentIDOrderByID(interactor.GetAgentStaffListByAgentIDOrderByIDInput{
		AgentID:      agentID,
		AgentStaffID: agentStaffID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffListJSONPresenter(responses.NewAgentStaffList(output.AgentStaffList)), nil
}

func (h *AgentStaffHandlerImpl) UpdateAgentStaff(agentStaffID uint, param entity.AgentStaffUpdateParam) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.UpdateAgentStaff(interactor.UpdateAgentStaffInput{
		AgentStaffID: agentStaffID,
		UpdateParam:  param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffJSONPresenter(responses.NewAgentStaff(output.AgentStaff)), nil
}

func (h *AgentStaffHandlerImpl) UpdateAgentStaffEmail(param entity.AgentStaffEmailUpdateParam, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.UpdateAgentStaffEmail(interactor.UpdateAgentStaffEmailInput{
		UpdateParam:  param,
		AgentStaffID: agentStaffID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *AgentStaffHandlerImpl) UpdateAgentStaffPassword(param entity.AgentStaffPasswordUpdateParam, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.UpdateAgentStaffPassword(interactor.UpdateAgentStaffPasswordInput{
		UpdateParam:  param,
		AgentStaffID: agentStaffID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *AgentStaffHandlerImpl) UpdateAgentStaffUsageEnd(param entity.DeleteAgentStaffParam) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.UpdateAgentStaffUsageEnd(interactor.UpdateAgentStaffUsageEndInput{
		Param: param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 担当者削除　*firebaseのアイパスを削除&DBのis_deletedをtrueにする
func (h *AgentStaffHandlerImpl) DeleteAgentStaff(param entity.DeleteAgentStaffParam) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.DeleteAgentStaff(interactor.DeleteAgentStaffInput{
		Param: param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *AgentStaffHandlerImpl) UpdateAgentStaffUsageReStart(agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.UpdateAgentStaffUsageReStart(interactor.UpdateAgentStaffUsageReStartInput{
		AgentStaffID: agentStaffID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// メール通知（求職者）の更新 body: {agent_staff_id, notification_job_seeker}
func (h *AgentStaffHandlerImpl) UpdateAgentStaffNotificationJobSeeker(param entity.UpdateAgentStaffNotificationJobSeekerParam) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.UpdateAgentStaffNotificationJobSeeker(interactor.UpdateAgentStaffNotificationJobSeekerInput{
		UpdateParam: param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// メール通知（未処理・未読）の更新 body: {agent_staff_id, notification_unwatched}
func (h *AgentStaffHandlerImpl) UpdateAgentStaffNotificationUnwatched(param entity.UpdateAgentStaffNotificationUnwatchedParam) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.UpdateAgentStaffNotificationUnwatched(interactor.UpdateAgentStaffNotificationUnwatchedInput{
		UpdateParam: param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 管理権限の更新 body: {agent_staff_id, authority}
func (h *AgentStaffHandlerImpl) UpdateAgentStaffAuthority(param entity.UpdateAgentStaffAuthorityParam) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.UpdateAgentStaffAuthority(interactor.UpdateAgentStaffAuthorityInput{
		UpdateParam: param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *AgentStaffHandlerImpl) GetOhterAgentStaffListByAgentIDAndAllianceAgentID(agentID, allianceAgentID, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.GetOhterAgentStaffListByAgentIDAndAllianceAgentID(interactor.GetOhterAgentStaffListByAgentIDAndAllianceAgentIDInput{
		AgentID:         agentID,
		AllianceAgentID: allianceAgentID,
		AgentStaffID:    agentStaffID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffListJSONPresenter(responses.NewAgentStaffList(output.AgentStaffList)), nil
}

func (h *AgentStaffHandlerImpl) GetAgentStaffListWithSaleNotCreated(token string, agentID, managementID uint) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.GetAgentStaffListWithSaleNotCreated(interactor.GetAgentStaffListWithSaleNotCreatedInput{
		Token:        token,
		AgentID:      agentID,
		ManagementID: managementID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffListJSONPresenter(responses.NewAgentStaffList(output.AgentStaffList)), nil
}

// 利用可能の担当者を取得
func (h *AgentStaffHandlerImpl) GetAgentStaffListByAgentIDAndUsageStatusAvailable(token string, agentID uint) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.GetAgentStaffListByAgentIDAndUsageStatusAvailable(interactor.GetAgentStaffListByAgentIDAndUsageStatusAvailableInput{
		Token:   token,
		AgentID: agentID})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffListJSONPresenter(responses.NewAgentStaffList(output.AgentStaffList)), nil
}

// 未削除の担当者を取得
func (h *AgentStaffHandlerImpl) GetAgentStaffListByAgentIDAndIsDeletedFalse(token string, agentID uint) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.GetAgentStaffListByAgentIDAndIsDeletedFalse(interactor.GetAgentStaffListByAgentIDAndIsDeletedFalseInput{
		Token:   token,
		AgentID: agentID})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffListJSONPresenter(responses.NewAgentStaffList(output.AgentStaffList)), nil
}

/****************************************************************************************/
/// Admin API

func (h *AgentStaffHandlerImpl) SignUpForAdmin(param entity.AgentStaffSignUpForAdminParam, agentID uint) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.SignUpForAdmin(interactor.AgentStaffSignUpForAdminInput{
		SignUpParam: param,
		AgentID:     agentID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffJSONPresenter(responses.NewAgentStaff(output.AgentStaff)), nil
}

// 全担当者一覧を取得
func (h *AgentStaffHandlerImpl) GetAllAgentStaffList() (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.GetAllAgentStaffList()

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffListJSONPresenter(responses.NewAgentStaffList(output.AgentStaffList)), nil
}

/****************************************************************************************/
/****************************************************************************************/
/// Agent API

func (h *AgentStaffHandlerImpl) GetAgentStaffMe(token string) (presenter.Presenter, error) {
	output, err := h.AgentStaffInteractor.GetAgentStaffMe(interactor.GetAgentStaffMeInput{Token: token})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffJSONPresenter(responses.NewAgentStaff(output.AgentStaff)), nil
}

/****************************************************************************************/
