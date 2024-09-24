package interactor

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type MessageTemplateInteractor interface {
	// 汎用系 API
	CreateMessageTemplate(input CreateMessageTemplateInput) (CreateMessageTemplateOutput, error)
	UpdateMessageTemplate(input UpdateMessageTemplateInput) (UpdateMessageTemplateOutput, error)
	DeleteMessageTemplate(input DeleteMessageTemplateInput) (DeleteMessageTemplateOutput, error)
	GetMessageTemplateByID(input GetMessageTemplateByIDInput) (GetMessageTemplateByIDOutput, error)
	GetMessageTemplateListByAgentStaffID(input GetMessageTemplateListByAgentStaffIDInput) (GetMessageTemplateListByAgentStaffIDOutput, error)
	GetMessageTemplateListByAgentStaffIDAndSendScene(input GetMessageTemplateListByAgentStaffIDAndSendSceneInput) (GetMessageTemplateListByAgentStaffIDAndSendSceneOutput, error)
	GetMessageTemplateListByAgentID(input GetMessageTemplateListByAgentIDInput) (GetMessageTemplateListByAgentIDOutput, error)
	// Admin API
}

type MessageTemplateInteractorImpl struct {
	firebase                  usecase.Firebase
	sendgrid                  config.Sendgrid
	messageTemplateRepository usecase.MessageTemplateRepository
}

// MessageTemplateInteractorImpl is an implementation of MessageTemplateInteractor
func NewMessageTemplateInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	mt usecase.MessageTemplateRepository,
) MessageTemplateInteractor {
	return &MessageTemplateInteractorImpl{
		firebase:                  fb,
		sendgrid:                  sg,
		messageTemplateRepository: mt,
	}
}

/****************************************************************************************/
/// 汎用系API
//
//タスクの作成
type CreateMessageTemplateInput struct {
	CreateParam entity.CreateOrUpdateMessageTemplateParam
}

type CreateMessageTemplateOutput struct {
	MessageTemplate *entity.MessageTemplate
}

func (i *MessageTemplateInteractorImpl) CreateMessageTemplate(input CreateMessageTemplateInput) (CreateMessageTemplateOutput, error) {
	var (
		output CreateMessageTemplateOutput
		err    error
	)

	// 次にタスクを作成
	messageTemplate := entity.NewMessageTemplate(
		input.CreateParam.AgentStaffID,
		input.CreateParam.SendScene,
		input.CreateParam.Title,
		input.CreateParam.Subject,
		input.CreateParam.Content,
	)

	err = i.messageTemplateRepository.Create(messageTemplate)
	if err != nil {
		return output, err
	}

	output.MessageTemplate = messageTemplate

	return output, nil
}

// タスクの更新
type UpdateMessageTemplateInput struct {
	UpdateParam entity.CreateOrUpdateMessageTemplateParam
	TemplateID  uint
}

type UpdateMessageTemplateOutput struct {
	MessageTemplate *entity.MessageTemplate
}

func (i *MessageTemplateInteractorImpl) UpdateMessageTemplate(input UpdateMessageTemplateInput) (UpdateMessageTemplateOutput, error) {
	var (
		output UpdateMessageTemplateOutput
		err    error
	)

	messageTemplate := entity.NewMessageTemplate(
		input.UpdateParam.AgentStaffID,
		input.UpdateParam.SendScene,
		input.UpdateParam.Title,
		input.UpdateParam.Subject,
		input.UpdateParam.Content,
	)

	err = i.messageTemplateRepository.Update(input.TemplateID, messageTemplate)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	messageTemplate.ID = input.TemplateID
	output.MessageTemplate = messageTemplate

	return output, nil
}

// メッセージテンプレートの削除
type DeleteMessageTemplateInput struct {
	TemplateID uint
}

type DeleteMessageTemplateOutput struct {
	OK bool
}

func (i *MessageTemplateInteractorImpl) DeleteMessageTemplate(input DeleteMessageTemplateInput) (DeleteMessageTemplateOutput, error) {
	var (
		output DeleteMessageTemplateOutput
	)

	err := i.messageTemplateRepository.Delete(input.TemplateID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

// IDを使ってメッセージテンプレートを取得する
type GetMessageTemplateByIDInput struct {
	TemplateID uint
}

type GetMessageTemplateByIDOutput struct {
	MessageTemplate *entity.MessageTemplate
}

func (i *MessageTemplateInteractorImpl) GetMessageTemplateByID(input GetMessageTemplateByIDInput) (GetMessageTemplateByIDOutput, error) {
	var (
		output GetMessageTemplateByIDOutput
		err    error
	)

	messageTemplate, err := i.messageTemplateRepository.FindByID(input.TemplateID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.MessageTemplate = messageTemplate

	return output, nil
}

// タスクグループの一覧取得（自分が関わっているタスク）
type GetMessageTemplateListByAgentStaffIDInput struct {
	AgentStaffID uint
}

type GetMessageTemplateListByAgentStaffIDOutput struct {
	MessageTemplateList []*entity.MessageTemplate
}

func (i *MessageTemplateInteractorImpl) GetMessageTemplateListByAgentStaffID(input GetMessageTemplateListByAgentStaffIDInput) (GetMessageTemplateListByAgentStaffIDOutput, error) {
	var (
		output GetMessageTemplateListByAgentStaffIDOutput
		err    error
	)

	messageTemplateList, err := i.messageTemplateRepository.GetByAgentStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.MessageTemplateList = messageTemplateList

	return output, nil
}

// タスクグループの一覧取得（自分が関わっているタスク）
type GetMessageTemplateListByAgentStaffIDAndSendSceneInput struct {
	AgentStaffID uint
	SendScene    uint
}

type GetMessageTemplateListByAgentStaffIDAndSendSceneOutput struct {
	MessageTemplateList []*entity.MessageTemplate
}

func (i *MessageTemplateInteractorImpl) GetMessageTemplateListByAgentStaffIDAndSendScene(input GetMessageTemplateListByAgentStaffIDAndSendSceneInput) (GetMessageTemplateListByAgentStaffIDAndSendSceneOutput, error) {
	var (
		output GetMessageTemplateListByAgentStaffIDAndSendSceneOutput
		err    error
	)

	messageTemplateList, err := i.messageTemplateRepository.GetByAgentStaffIDAndSendScene(input.AgentStaffID, input.SendScene)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.MessageTemplateList = messageTemplateList

	return output, nil
}

// タスクグループの一覧取得（自分が関わっているタスク）
type GetMessageTemplateListByAgentIDInput struct {
	AgentID uint
}

type GetMessageTemplateListByAgentIDOutput struct {
	MessageTemplateList []*entity.MessageTemplate
}

func (i *MessageTemplateInteractorImpl) GetMessageTemplateListByAgentID(input GetMessageTemplateListByAgentIDInput) (GetMessageTemplateListByAgentIDOutput, error) {
	var (
		output GetMessageTemplateListByAgentIDOutput
		err    error
	)

	messageTemplateList, err := i.messageTemplateRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.MessageTemplateList = messageTemplateList

	return output, nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API
//

/****************************************************************************************/
