package interactor

import (
	"errors"
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type DeploymentReflectionInteractor interface {
	// 汎用系
	CreateDeploymentInformation(input CreateDeploymentInformationInput) (CreateDeploymentInformationOutput, error)
	UpdateDeploymentInformation(input UpdateDeploymentInformationInput) (UpdateDeploymentInformationOutput, error)
	GetAllDeploymentInformations(input GetAllDeploymentInformationsInput) (GetAllDeploymentInformationsOutput, error)
	UpdateDeployMenConfirmStatusByStaffID(input UpdateDeployMenConfirmStatusByStaffIDInput) (UpdateDeployMenConfirmStatusByStaffIDOutput, error)
	CheckDeploymentReflection(input CheckDeploymentReflectionInput) (CheckDeploymentReflectionOutput, error)
}

type DeploymentReflectionInteractorImpl struct {
	firebase                        usecase.Firebase
	sendgrid                        config.Sendgrid
	oneSignal                       config.OneSignal
	deploymentInformationRepository usecase.DeploymentInformationRepository
	deploymentReflectionRepository  usecase.DeploymentReflectionRepository
	agentStaffRepository            usecase.AgentStaffRepository
}

// DeploymentReflectionInteractorImpl is an implementation of DeploymentReflectionInteractor
func NewDeploymentReflectionInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	diR usecase.DeploymentInformationRepository,
	drR usecase.DeploymentReflectionRepository,
	asR usecase.AgentStaffRepository,
) DeploymentReflectionInteractor {
	return &DeploymentReflectionInteractorImpl{
		firebase:                        fb,
		sendgrid:                        sg,
		oneSignal:                       os,
		deploymentInformationRepository: diR,
		deploymentReflectionRepository:  drR,
		agentStaffRepository:            asR,
	}
}

/****************************************************************************************/
/// 汎用系
//

// デプロイ反映の情報を一括作成
type CreateDeploymentInformationInput struct {
	CreateParam entity.CreateOrUpdateDeploymentInformationParam
}

type CreateDeploymentInformationOutput struct {
	DeploymentInformation *entity.DeploymentInformation
}

func (i *DeploymentReflectionInteractorImpl) CreateDeploymentInformation(input CreateDeploymentInformationInput) (CreateDeploymentInformationOutput, error) {
	var (
		output CreateDeploymentInformationOutput
		err    error
	)

	/************ 1. デプロイ情報の登録 **************/

	deployInfo := entity.NewDeploymentInformation(
		input.CreateParam.BeVer, input.CreateParam.FeVer,
		input.CreateParam.BeDetail, input.CreateParam.FeDetail,
	)

	err = i.deploymentInformationRepository.Create(deployInfo)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 2. デプロイの反映状況の情報を登録 **************/

	agentStaffList, err := i.agentStaffRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	err = i.deploymentReflectionRepository.CreateMulti(deployInfo.ID, agentStaffList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.DeploymentInformation = deployInfo

	return output, nil
}

type GetAllDeploymentInformationsInput struct {
	PageNumber uint
}

type GetAllDeploymentInformationsOutput struct {
	MaxPageNumber             uint
	DeploymentInformationList []*entity.DeploymentInformation
}

func (i *DeploymentReflectionInteractorImpl) GetAllDeploymentInformations(input GetAllDeploymentInformationsInput) (GetAllDeploymentInformationsOutput, error) {
	var (
		output GetAllDeploymentInformationsOutput
	)

	deploymentList, err := i.deploymentInformationRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// ページの最大数を取得
	output.MaxPageNumber = getDeploymentMaxPage(deploymentList)

	// 全て反映ずみの場合はTrue
	output.DeploymentInformationList = getDeploymentListWithPage(deploymentList, input.PageNumber)

	return output, nil
}

// デプロイ反映の情報を一括作成
type UpdateDeploymentInformationInput struct {
	DeploymentID uint
	UpdateParam  entity.CreateOrUpdateDeploymentInformationParam
}

type UpdateDeploymentInformationOutput struct {
	DeploymentInformation *entity.DeploymentInformation
}

func (i *DeploymentReflectionInteractorImpl) UpdateDeploymentInformation(input UpdateDeploymentInformationInput) (UpdateDeploymentInformationOutput, error) {
	var (
		output UpdateDeploymentInformationOutput
		err    error
	)

	/************ 1. デプロイ情報の更新 **************/

	deployInfo := entity.NewDeploymentInformation(
		input.UpdateParam.BeVer, input.UpdateParam.FeVer,
		input.UpdateParam.BeDetail, input.UpdateParam.FeDetail,
	)

	err = i.deploymentInformationRepository.Update(input.DeploymentID, deployInfo)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.DeploymentInformation = deployInfo

	return output, nil
}

// Systemからのお知らせを確認したユーザーを作成P
type UpdateDeployMenConfirmStatusByStaffIDInput struct {
	AgentStaffID uint
}

type UpdateDeployMenConfirmStatusByStaffIDOutput struct {
	OK bool
}

func (i *DeploymentReflectionInteractorImpl) UpdateDeployMenConfirmStatusByStaffID(input UpdateDeployMenConfirmStatusByStaffIDInput) (UpdateDeployMenConfirmStatusByStaffIDOutput, error) {
	var (
		output UpdateDeployMenConfirmStatusByStaffIDOutput
	)

	err := i.deploymentReflectionRepository.UpdateIsReflectedByAgentStaffID(input.AgentStaffID, true)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 現在時刻
	now := time.Now().In(time.UTC)

	err = i.agentStaffRepository.UpdateLastLogin(input.AgentStaffID, now)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 担当者IDからお知らせの未読件数を取得
type CheckDeploymentReflectionInput struct {
	AgentStaffID uint
}

type CheckDeploymentReflectionOutput struct {
	OK bool
}

func (i *DeploymentReflectionInteractorImpl) CheckDeploymentReflection(input CheckDeploymentReflectionInput) (CheckDeploymentReflectionOutput, error) {
	var (
		output CheckDeploymentReflectionOutput
	)

	_, err := i.deploymentReflectionRepository.FindNotReflectedByAgentStaffID(input.AgentStaffID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			// レコード取得できない場合はデプロイ情報反映済みのため、Trueを返す
			output.OK = true
			return output, nil
		} else {
			fmt.Println(err)
			return output, err
		}
	} else {
		// レコード取得した場合はデプロイ情報が未反映のため、Falseを返す
		output.OK = false
		return output, nil
	}
}
