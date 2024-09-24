package handler

import (
	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SelectionQuestionnaireHandler interface {
	// 汎用系 API
	CreateSelectionQuestionnaire(param entity.CreateOrUpdateSelectionQuestionnaireParam, questionnaireUUID uuid.UUID) (presenter.Presenter, error)
	UpdateSelectionQuestionnaire(param entity.CreateOrUpdateSelectionQuestionnaireParam, questionnaireID uint) (presenter.Presenter, error)
	GetSelectionQuestionnaireOrNullByUUID(questionnaireUUID uuid.UUID) (presenter.Presenter, error)
	GenerateSelectionQuestionnaireByUUID() (presenter.Presenter, error)
	GetUnansweredQuestionnaireListByJobSeekerUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	GetQuestionnaireForJobSeeker(jobSeekerUUID, jobInformationUUID uuid.UUID, selectionPhase uint) (presenter.Presenter, error)
}

type SelectionQuestionnaireHandlerImpl struct {
	selectionQuestionnaireInteractor interactor.SelectionQuestionnaireInteractor
}

func NewSelectionQuestionnaireHandlerImpl(mtI interactor.SelectionQuestionnaireInteractor) SelectionQuestionnaireHandler {
	return &SelectionQuestionnaireHandlerImpl{
		selectionQuestionnaireInteractor: mtI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//

func (h *SelectionQuestionnaireHandlerImpl) CreateSelectionQuestionnaire(param entity.CreateOrUpdateSelectionQuestionnaireParam, questionnaireUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.selectionQuestionnaireInteractor.CreateSelectionQuestionnaire(interactor.CreateSelectionQuestionnaireInput{
		CreateParam:       param,
		QuestionnaireUUID: questionnaireUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionQuestionnaireJSONPresenter(responses.NewSelectionQuestionnaire(output.SelectionQuestionnaire)), nil
}

func (h *SelectionQuestionnaireHandlerImpl) UpdateSelectionQuestionnaire(param entity.CreateOrUpdateSelectionQuestionnaireParam, questionnaireID uint) (presenter.Presenter, error) {
	output, err := h.selectionQuestionnaireInteractor.UpdateSelectionQuestionnaire(interactor.UpdateSelectionQuestionnaireInput{
		UpdateParam:     param,
		QuestionnaireID: questionnaireID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionQuestionnaireJSONPresenter(responses.NewSelectionQuestionnaire(output.SelectionQuestionnaire)), nil
}

func (h *SelectionQuestionnaireHandlerImpl) GetSelectionQuestionnaireOrNullByUUID(questionnaireUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.selectionQuestionnaireInteractor.GetSelectionQuestionnaireOrNullByUUID(interactor.GetSelectionQuestionnaireOrNullByUUIDInput{
		QuestionnaireUUID: questionnaireUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionQuestionnaireJSONPresenter(responses.NewSelectionQuestionnaire(output.SelectionQuestionnaire)), nil
}

func (h *SelectionQuestionnaireHandlerImpl) GenerateSelectionQuestionnaireByUUID() (presenter.Presenter, error) {
	output, err := h.selectionQuestionnaireInteractor.GenerateSelectionQuestionnaireByUUID()

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionQuestionnaireUUIDJSONPresenter(responses.NewSelectionQuestionnaireUUID(output.QuestionnaireUUID)), nil
}

func (h *SelectionQuestionnaireHandlerImpl) GetUnansweredQuestionnaireListByJobSeekerUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.selectionQuestionnaireInteractor.GetUnansweredQuestionnaireListByJobSeekerUUID(interactor.GetUnansweredQuestionnaireListByJobSeekerUUIDInput{
		JobSeekerUUID: jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionQuestionnaireListJSONPresenter(responses.NewSelectionQuestionnaireList(output.SelectionQuestionnaireList)), nil
}

// 求職者のアンケート情報を取得する（queryParam: job_seeker_uuid, job_information_uuid, selection_information_id）
func (h *SelectionQuestionnaireHandlerImpl) GetQuestionnaireForJobSeeker(jobSeekerUUID, jobInformationUUID uuid.UUID, selectionPhase uint) (presenter.Presenter, error) {
	output, err := h.selectionQuestionnaireInteractor.GetQuestionnaireForJobSeeker(interactor.GetQuestionnaireForJobSeekerInput{
		JobSeekerUUID:      jobSeekerUUID,
		JobInformationUUID: jobInformationUUID,
		SelectionPhase:     selectionPhase,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionQuestionnaireJSONPresenter(responses.NewSelectionQuestionnaire(output.SelectionQuestionnaire)), nil
}

/****************************************************************************************/
