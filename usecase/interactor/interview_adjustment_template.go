package interactor

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type InterviewAdjustmentTemplateInteractor interface {
	// 汎用系
	CreateInterviewAdjustmentTemplate(input CreateInterviewAdjustmentTemplateInput) (CreateInterviewAdjustmentTemplateOutput, error)                                                                   // 面談調整テンプレートの登録
	UpdateInterviewAdjustmentTemplate(input UpdateInterviewAdjustmentTemplateInput) (UpdateInterviewAdjustmentTemplateOutput, error)                                                                   // 面談調整テンプレートの更新
	DeleteInterviewAdjustmentTemplate(input DeleteInterviewAdjustmentTemplateInput) (DeleteInterviewAdjustmentTemplateOutput, error)                                                                   // 面談調整テンプレートの削除
	GetInterviewAdjustmentTemplateByID(input GetInterviewAdjustmentTemplateByIDInput) (GetInterviewAdjustmentTemplateByIDOutput, error)                                                                // IDから面談調整テンプレートの取得
	GetInterviewAdjustmentTemplateListByAgentID(input GetInterviewAdjustmentTemplateListByAgentIDInput) (GetInterviewAdjustmentTemplateListByAgentIDOutput, error)                                     // エージェントIDから面談調整テンプレートの取得
	GetInterviewAdjustmentTemplateListByAgentIDAndSendScene(input GetInterviewAdjustmentTemplateListByAgentIDAndSendSceneInput) (GetInterviewAdjustmentTemplateListByAgentIDAndSendSceneOutput, error) // 担当者IDと送信シーンから面談調整テンプレートの取得
	// Admin API
}

type InterviewAdjustmentTemplateInteractorImpl struct {
	firebase                              usecase.Firebase
	sendgrid                              config.Sendgrid
	interviewAdjustmentTemplateRepository usecase.InterviewAdjustmentTemplateRepository
	jobSeekerRepository                   usecase.JobSeekerRepository
}

// InterviewAdjustmentTemplateInteractorImpl is an implementation of InterviewAdjustmentTemplateInteractor
func NewInterviewAdjustmentTemplateInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	iatR usecase.InterviewAdjustmentTemplateRepository,
	jsR usecase.JobSeekerRepository,
) InterviewAdjustmentTemplateInteractor {
	return &InterviewAdjustmentTemplateInteractorImpl{
		firebase:                              fb,
		sendgrid:                              sg,
		interviewAdjustmentTemplateRepository: iatR,
		jobSeekerRepository:                   jsR,
	}
}

/****************************************************************************************/
/// 汎用系
//
// 面談調整テンプレートの登録
type CreateInterviewAdjustmentTemplateInput struct {
	CreateParam entity.CreateOrUpdateInterviewAdjustmentTemplateParam
}

type CreateInterviewAdjustmentTemplateOutput struct {
	InterviewAdjustmentTemplate *entity.InterviewAdjustmentTemplate
}

func (i *InterviewAdjustmentTemplateInteractorImpl) CreateInterviewAdjustmentTemplate(input CreateInterviewAdjustmentTemplateInput) (CreateInterviewAdjustmentTemplateOutput, error) {
	var (
		output CreateInterviewAdjustmentTemplateOutput
		err    error
	)

	// タスクの作成
	template := entity.NewInterviewAdjustmentTemplate(
		input.CreateParam.AgentID,
		input.CreateParam.SendScene,
		input.CreateParam.Title,
		input.CreateParam.Subject,
		input.CreateParam.Content,
	)

	err = i.interviewAdjustmentTemplateRepository.Create(template)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InterviewAdjustmentTemplate = template

	return output, nil
}

// 面談調整テンプレートの更新
type UpdateInterviewAdjustmentTemplateInput struct {
	UpdateParam entity.CreateOrUpdateInterviewAdjustmentTemplateParam
	TemplateID  uint
}

type UpdateInterviewAdjustmentTemplateOutput struct {
	InterviewAdjustmentTemplate *entity.InterviewAdjustmentTemplate
}

func (i *InterviewAdjustmentTemplateInteractorImpl) UpdateInterviewAdjustmentTemplate(input UpdateInterviewAdjustmentTemplateInput) (UpdateInterviewAdjustmentTemplateOutput, error) {
	var (
		output UpdateInterviewAdjustmentTemplateOutput
		err    error
	)

	// タスクの作成
	template := entity.NewInterviewAdjustmentTemplate(
		input.UpdateParam.AgentID,
		input.UpdateParam.SendScene,
		input.UpdateParam.Title,
		input.UpdateParam.Subject,
		input.UpdateParam.Content,
	)

	err = i.interviewAdjustmentTemplateRepository.Update(input.TemplateID, template)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InterviewAdjustmentTemplate = template

	return output, nil
}

// 面談調整テンプレートの削除
type DeleteInterviewAdjustmentTemplateInput struct {
	TemplateID uint
}

type DeleteInterviewAdjustmentTemplateOutput struct {
	OK bool
}

func (i *InterviewAdjustmentTemplateInteractorImpl) DeleteInterviewAdjustmentTemplate(input DeleteInterviewAdjustmentTemplateInput) (DeleteInterviewAdjustmentTemplateOutput, error) {
	var (
		output DeleteInterviewAdjustmentTemplateOutput
		err    error
	)

	err = i.interviewAdjustmentTemplateRepository.Delete(input.TemplateID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// IDから面談調整テンプレートの取得
type GetInterviewAdjustmentTemplateByIDInput struct {
	TemplateID uint
}

type GetInterviewAdjustmentTemplateByIDOutput struct {
	InterviewAdjustmentTemplate *entity.InterviewAdjustmentTemplate
}

func (i *InterviewAdjustmentTemplateInteractorImpl) GetInterviewAdjustmentTemplateByID(input GetInterviewAdjustmentTemplateByIDInput) (GetInterviewAdjustmentTemplateByIDOutput, error) {
	var (
		output GetInterviewAdjustmentTemplateByIDOutput
		err    error
	)

	template, err := i.interviewAdjustmentTemplateRepository.FindByID(input.TemplateID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InterviewAdjustmentTemplate = template

	return output, nil
}

// エージェントIDから面談調整テンプレートの取得
type GetInterviewAdjustmentTemplateListByAgentIDInput struct {
	AgentID uint
}

type GetInterviewAdjustmentTemplateListByAgentIDOutput struct {
	InterviewAdjustmentTemplateList []*entity.InterviewAdjustmentTemplate
}

func (i *InterviewAdjustmentTemplateInteractorImpl) GetInterviewAdjustmentTemplateListByAgentID(input GetInterviewAdjustmentTemplateListByAgentIDInput) (GetInterviewAdjustmentTemplateListByAgentIDOutput, error) {
	var (
		output GetInterviewAdjustmentTemplateListByAgentIDOutput
		err    error
	)

	templateList, err := i.interviewAdjustmentTemplateRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InterviewAdjustmentTemplateList = templateList

	return output, nil
}

// 担当者IDと送信シーンから面談調整テンプレートの取得
type GetInterviewAdjustmentTemplateListByAgentIDAndSendSceneInput struct {
	AgentID   uint
	SendScene uint
}

type GetInterviewAdjustmentTemplateListByAgentIDAndSendSceneOutput struct {
	InterviewAdjustmentTemplateList []*entity.InterviewAdjustmentTemplate
}

func (i *InterviewAdjustmentTemplateInteractorImpl) GetInterviewAdjustmentTemplateListByAgentIDAndSendScene(input GetInterviewAdjustmentTemplateListByAgentIDAndSendSceneInput) (GetInterviewAdjustmentTemplateListByAgentIDAndSendSceneOutput, error) {
	var (
		output GetInterviewAdjustmentTemplateListByAgentIDAndSendSceneOutput
		err    error
	)

	templateList, err := i.interviewAdjustmentTemplateRepository.GetByAgentIDAndSendScene(input.AgentID, input.SendScene)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InterviewAdjustmentTemplateList = templateList

	return output, nil
}
