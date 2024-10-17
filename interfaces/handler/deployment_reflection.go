package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type DeploymentReflectionHandler interface {
	// 汎用系 API
	CreateDeploymentInformation(param entity.CreateOrUpdateDeploymentInformationParam) (presenter.Presenter, error)
	UpdateDeploymentInformation(param entity.CreateOrUpdateDeploymentInformationParam, deploymentID uint) (presenter.Presenter, error)
	GetAllDeploymentInformations(pageNumber uint) (presenter.Presenter, error)
	UpdateDeployMenConfirmStatusByStaffID(agentStaffID uint) (presenter.Presenter, error)
	CheckDeploymentReflection(agentStaffID uint) (presenter.Presenter, error)
}

type DeploymentReflectionHandlerImpl struct {
	notificationForUserInteractor interactor.DeploymentReflectionInteractor
}

func NewDeploymentReflectionHandlerImpl(nfuI interactor.DeploymentReflectionInteractor) DeploymentReflectionHandler {
	return &DeploymentReflectionHandlerImpl{
		notificationForUserInteractor: nfuI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//

// デプロイ情報を作成
func (h *DeploymentReflectionHandlerImpl) CreateDeploymentInformation(param entity.CreateOrUpdateDeploymentInformationParam) (presenter.Presenter, error) {
	output, err := h.notificationForUserInteractor.CreateDeploymentInformation(interactor.CreateDeploymentInformationInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewDeploymentInformationJSONPresenter(responses.NewDeploymentInformation(output.DeploymentInformation)), nil
}

// デプロイ情報を更新
func (h *DeploymentReflectionHandlerImpl) UpdateDeploymentInformation(param entity.CreateOrUpdateDeploymentInformationParam, deploymentID uint) (presenter.Presenter, error) {
	output, err := h.notificationForUserInteractor.UpdateDeploymentInformation(interactor.UpdateDeploymentInformationInput{
		UpdateParam:  param,
		DeploymentID: deploymentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewDeploymentInformationJSONPresenter(responses.NewDeploymentInformation(output.DeploymentInformation)), nil
}

func (h *DeploymentReflectionHandlerImpl) GetAllDeploymentInformations(pageNumber uint) (presenter.Presenter, error) {
	output, err := h.notificationForUserInteractor.GetAllDeploymentInformations(interactor.GetAllDeploymentInformationsInput{
		PageNumber: pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewDeploymentInformationListAndMaxPageJSONPresenter(responses.NewDeploymentInformationListAndMaxPage(output.MaxPageNumber, output.DeploymentInformationList)), nil
}

// Systemからのお知らせを確認したユーザーを作成
func (h *DeploymentReflectionHandlerImpl) UpdateDeployMenConfirmStatusByStaffID(agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.notificationForUserInteractor.UpdateDeployMenConfirmStatusByStaffID(interactor.UpdateDeployMenConfirmStatusByStaffIDInput{
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 担当者IDからお知らせの未読件数を取得
func (h *DeploymentReflectionHandlerImpl) CheckDeploymentReflection(agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.notificationForUserInteractor.CheckDeploymentReflection(interactor.CheckDeploymentReflectionInput{
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}
