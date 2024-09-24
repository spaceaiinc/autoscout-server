package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type MessageTemplateHandler interface {
	// 汎用系 API
	CreateMessageTemplate(param entity.CreateOrUpdateMessageTemplateParam) (presenter.Presenter, error)
	UpdateMessageTemplate(param entity.CreateOrUpdateMessageTemplateParam, templateID uint) (presenter.Presenter, error)
	DeleteMessageTemplate(templateID uint) (presenter.Presenter, error)
	GetMessageTemplateByID(templateID uint) (presenter.Presenter, error)
	GetMessageTemplateListByAgentStaffID(agentStaffID uint) (presenter.Presenter, error)
	GetMessageTemplateListByAgentStaffIDAndSendScene(agentStaffID, sendScene uint) (presenter.Presenter, error)
	GetMessageTemplateListByAgentID(messageTemplateGroupID uint) (presenter.Presenter, error)

	// Admin API
}

type MessageTemplateHandlerImpl struct {
	messageTemplateInteractor interactor.MessageTemplateInteractor
}

func NewMessageTemplateHandlerImpl(mtI interactor.MessageTemplateInteractor) MessageTemplateHandler {
	return &MessageTemplateHandlerImpl{
		messageTemplateInteractor: mtI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *MessageTemplateHandlerImpl) CreateMessageTemplate(param entity.CreateOrUpdateMessageTemplateParam) (presenter.Presenter, error) {
	output, err := h.messageTemplateInteractor.CreateMessageTemplate(interactor.CreateMessageTemplateInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMessageTemplateJSONPresenter(responses.NewMessageTemplate(output.MessageTemplate)), nil
}

func (h *MessageTemplateHandlerImpl) UpdateMessageTemplate(param entity.CreateOrUpdateMessageTemplateParam, templateID uint) (presenter.Presenter, error) {
	output, err := h.messageTemplateInteractor.UpdateMessageTemplate(interactor.UpdateMessageTemplateInput{
		UpdateParam: param,
		TemplateID:  templateID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMessageTemplateJSONPresenter(responses.NewMessageTemplate(output.MessageTemplate)), nil
}

func (h *MessageTemplateHandlerImpl) DeleteMessageTemplate(templateID uint) (presenter.Presenter, error) {
	output, err := h.messageTemplateInteractor.DeleteMessageTemplate(interactor.DeleteMessageTemplateInput{
		TemplateID: templateID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// メッセージテンプレートIDからメッセージテンプレートを取得する
func (h *MessageTemplateHandlerImpl) GetMessageTemplateByID(templateID uint) (presenter.Presenter, error) {
	output, err := h.messageTemplateInteractor.GetMessageTemplateByID(interactor.GetMessageTemplateByIDInput{
		TemplateID: templateID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMessageTemplateJSONPresenter(responses.NewMessageTemplate(output.MessageTemplate)), nil
}

// メッセージテンプレートグループの一覧取得（自分が関わっているメッセージテンプレート）
func (h *MessageTemplateHandlerImpl) GetMessageTemplateListByAgentStaffID(agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.messageTemplateInteractor.GetMessageTemplateListByAgentStaffID(interactor.GetMessageTemplateListByAgentStaffIDInput{
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMessageTemplateListJSONPresenter(responses.NewMessageTemplateList(output.MessageTemplateList)), nil
}

// メッセージテンプレートグループの一覧取得（自分が関わっているメッセージテンプレート）
func (h *MessageTemplateHandlerImpl) GetMessageTemplateListByAgentStaffIDAndSendScene(agentStaffID, sendScene uint) (presenter.Presenter, error) {
	output, err := h.messageTemplateInteractor.GetMessageTemplateListByAgentStaffIDAndSendScene(interactor.GetMessageTemplateListByAgentStaffIDAndSendSceneInput{
		AgentStaffID: agentStaffID,
		SendScene:    sendScene,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMessageTemplateListJSONPresenter(responses.NewMessageTemplateList(output.MessageTemplateList)), nil
}

// メッセージテンプレートグループの一覧取得（自分が関わっているメッセージテンプレート）
func (h *MessageTemplateHandlerImpl) GetMessageTemplateListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.messageTemplateInteractor.GetMessageTemplateListByAgentID(interactor.GetMessageTemplateListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMessageTemplateListJSONPresenter(responses.NewMessageTemplateList(output.MessageTemplateList)), nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API

/****************************************************************************************/
