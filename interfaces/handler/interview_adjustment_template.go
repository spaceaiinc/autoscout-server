package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

// 面談調整テンプレート
type InterviewAdjustmentTemplateHandler interface {
	// 汎用系 API
	CreateInterviewAdjustmentTemplate(param entity.CreateOrUpdateInterviewAdjustmentTemplateParam) (presenter.Presenter, error)                  // 面談調整テンプレートの登録
	UpdateInterviewAdjustmentTemplate(param entity.CreateOrUpdateInterviewAdjustmentTemplateParam, templateID uint) (presenter.Presenter, error) // 面談調整テンプレートの更新
	DeleteInterviewAdjustmentTemplate(templateID uint) (presenter.Presenter, error)                                                              // 面談調整テンプレートの削除
	GetInterviewAdjustmentTemplateByID(templateID uint) (presenter.Presenter, error)                                                             // IDから面談調整テンプレートを取得
	GetInterviewAdjustmentTemplateListByAgentID(agentID uint) (presenter.Presenter, error)                                                       // エージェントIDから面談調整テンプレートを取得
	GetInterviewAdjustmentTemplateListByAgentIDAndSendScene(agentID, sendScene uint) (presenter.Presenter, error)                                // エージェントIDと送信シーンから面談調整テンプレートを取得

	// Admin API
}

type InterviewAdjustmentTemplateHandlerImpl struct {
	interviewAdjustmentTemplateInteractor interactor.InterviewAdjustmentTemplateInteractor
}

func NewInterviewAdjustmentTemplateHandlerImpl(itI interactor.InterviewAdjustmentTemplateInteractor) InterviewAdjustmentTemplateHandler {
	return &InterviewAdjustmentTemplateHandlerImpl{
		interviewAdjustmentTemplateInteractor: itI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *InterviewAdjustmentTemplateHandlerImpl) CreateInterviewAdjustmentTemplate(param entity.CreateOrUpdateInterviewAdjustmentTemplateParam) (presenter.Presenter, error) {
	output, err := h.interviewAdjustmentTemplateInteractor.CreateInterviewAdjustmentTemplate(interactor.CreateInterviewAdjustmentTemplateInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInterviewAdjustmentTemplateJSONPresenter(responses.NewInterviewAdjustmentTemplate(output.InterviewAdjustmentTemplate)), nil
}

func (h *InterviewAdjustmentTemplateHandlerImpl) UpdateInterviewAdjustmentTemplate(param entity.CreateOrUpdateInterviewAdjustmentTemplateParam, templateID uint) (presenter.Presenter, error) {
	output, err := h.interviewAdjustmentTemplateInteractor.UpdateInterviewAdjustmentTemplate(interactor.UpdateInterviewAdjustmentTemplateInput{
		UpdateParam: param,
		TemplateID:  templateID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInterviewAdjustmentTemplateJSONPresenter(responses.NewInterviewAdjustmentTemplate(output.InterviewAdjustmentTemplate)), nil
}

func (h *InterviewAdjustmentTemplateHandlerImpl) DeleteInterviewAdjustmentTemplate(templateID uint) (presenter.Presenter, error) {
	output, err := h.interviewAdjustmentTemplateInteractor.DeleteInterviewAdjustmentTemplate(interactor.DeleteInterviewAdjustmentTemplateInput{
		TemplateID: templateID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *InterviewAdjustmentTemplateHandlerImpl) GetInterviewAdjustmentTemplateByID(templateID uint) (presenter.Presenter, error) {
	output, err := h.interviewAdjustmentTemplateInteractor.GetInterviewAdjustmentTemplateByID(interactor.GetInterviewAdjustmentTemplateByIDInput{
		TemplateID: templateID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInterviewAdjustmentTemplateJSONPresenter(responses.NewInterviewAdjustmentTemplate(output.InterviewAdjustmentTemplate)), nil
}

func (h *InterviewAdjustmentTemplateHandlerImpl) GetInterviewAdjustmentTemplateListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.interviewAdjustmentTemplateInteractor.GetInterviewAdjustmentTemplateListByAgentID(interactor.GetInterviewAdjustmentTemplateListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInterviewAdjustmentTemplateListJSONPresenter(responses.NewInterviewAdjustmentTemplateList(output.InterviewAdjustmentTemplateList)), nil
}

func (h *InterviewAdjustmentTemplateHandlerImpl) GetInterviewAdjustmentTemplateListByAgentIDAndSendScene(agentID, sendScene uint) (presenter.Presenter, error) {
	output, err := h.interviewAdjustmentTemplateInteractor.GetInterviewAdjustmentTemplateListByAgentIDAndSendScene(interactor.GetInterviewAdjustmentTemplateListByAgentIDAndSendSceneInput{
		AgentID:   agentID,
		SendScene: sendScene,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInterviewAdjustmentTemplateListJSONPresenter(responses.NewInterviewAdjustmentTemplateList(output.InterviewAdjustmentTemplateList)), nil
}
