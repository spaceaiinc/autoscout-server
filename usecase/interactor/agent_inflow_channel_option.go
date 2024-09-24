package interactor

import (
	"log"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentInflowChannelOptionInteractor interface {
	// 汎用系 API
	CreateAgentInflowChannelOption(input CreateAgentInflowChannelOptionInput) (CreateAgentInflowChannelOptionByIDOutput, error)
	UpdateAgentInflowChannelOption(input UpdateAgentInflowChannelOptionInput) (UpdateAgentInflowChannelOptionByIDOutput, error)
	GetAgentInflowChannelOptionByID(input GetAgentInflowChannelOptionByIDInput) (GetAgentInflowChannelOptionByIDOutput, error)
	GetAgentInflowChannelOptionListByAgentID(input GetAgentInflowChannelOptionListByAgentIDInput) (GetAgentInflowChannelOptionListByAgentIDOutput, error)
}

type AgentInflowChannelOptionInteractorImpl struct {
	firebase                           usecase.Firebase
	sendgrid                           config.Sendgrid
	oneSignal                          config.OneSignal
	agentInflowChannelOptionRepository usecase.AgentInflowChannelOptionRepository
}

// AgentInflowChannelOptionInteractorImpl is an implementation of AgentInflowChannelOptionInteractor
func NewAgentInflowChannelOptionInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	arR usecase.AgentInflowChannelOptionRepository,
) AgentInflowChannelOptionInteractor {
	return &AgentInflowChannelOptionInteractorImpl{
		firebase:                           fb,
		sendgrid:                           sg,
		oneSignal:                          os,
		agentInflowChannelOptionRepository: arR,
	}
}

/****************************************************************************************/
// 汎用系API
//
// エージェントの流入経路マスタを作成する
type CreateAgentInflowChannelOptionInput struct {
	CreateParam entity.AgentInflowChannelOption
}

type CreateAgentInflowChannelOptionByIDOutput struct {
	AgentInflowChannelOption *entity.AgentInflowChannelOption
}

func (i *AgentInflowChannelOptionInteractorImpl) CreateAgentInflowChannelOption(input CreateAgentInflowChannelOptionInput) (CreateAgentInflowChannelOptionByIDOutput, error) {
	var (
		output                   CreateAgentInflowChannelOptionByIDOutput
		agentInflowChannelOption *entity.AgentInflowChannelOption
		err                      error
	)

	agentInflowChannelOption = entity.NewAgentInflowChannelOption(
		input.CreateParam.AgentID,
		input.CreateParam.ChannelName,
		input.CreateParam.IsOpen,
	)

	// エージェントの流入経路マスタを作成
	err = i.agentInflowChannelOptionRepository.Create(agentInflowChannelOption)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.AgentInflowChannelOption = agentInflowChannelOption

	return output, nil
}

// エージェントの流入経路マスタを更新する
type UpdateAgentInflowChannelOptionInput struct {
	UpdateParam                entity.AgentInflowChannelOption
	AgentInflowChannelOptionID uint
}

type UpdateAgentInflowChannelOptionByIDOutput struct {
	AgentInflowChannelOption *entity.AgentInflowChannelOption
}

func (i *AgentInflowChannelOptionInteractorImpl) UpdateAgentInflowChannelOption(input UpdateAgentInflowChannelOptionInput) (UpdateAgentInflowChannelOptionByIDOutput, error) {
	var (
		output                   UpdateAgentInflowChannelOptionByIDOutput
		agentInflowChannelOption *entity.AgentInflowChannelOption
		err                      error
	)

	agentInflowChannelOption = entity.NewAgentInflowChannelOption(
		input.UpdateParam.AgentID,
		input.UpdateParam.ChannelName,
		input.UpdateParam.IsOpen,
	)

	// エージェントの流入経路マスタを更新
	err = i.agentInflowChannelOptionRepository.Update(input.AgentInflowChannelOptionID, agentInflowChannelOption)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.AgentInflowChannelOption = agentInflowChannelOption
	// 更新後にリストの対応レコードも更新するためにIDも返す
	output.AgentInflowChannelOption.ID = input.AgentInflowChannelOptionID

	return output, nil
}

// エージェントの流入経路マスタを取得する
type GetAgentInflowChannelOptionByIDInput struct {
	AgentInflowChannelOptionID uint
}

type GetAgentInflowChannelOptionByIDOutput struct {
	AgentInflowChannelOption *entity.AgentInflowChannelOption
}

func (i *AgentInflowChannelOptionInteractorImpl) GetAgentInflowChannelOptionByID(input GetAgentInflowChannelOptionByIDInput) (GetAgentInflowChannelOptionByIDOutput, error) {
	var (
		output                   GetAgentInflowChannelOptionByIDOutput
		agentInflowChannelOption *entity.AgentInflowChannelOption
		err                      error
	)

	// エージェントの流入経路マスタを取得
	agentInflowChannelOption, err = i.agentInflowChannelOptionRepository.FindByID(input.AgentInflowChannelOptionID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.AgentInflowChannelOption = agentInflowChannelOption

	return output, nil
}

// エージェントIDからエージェントの流入経路マスタを取得する
type GetAgentInflowChannelOptionListByAgentIDInput struct {
	AgentID uint
}

type GetAgentInflowChannelOptionListByAgentIDOutput struct {
	AgentInflowChannelOptionList []*entity.AgentInflowChannelOption
}

func (i *AgentInflowChannelOptionInteractorImpl) GetAgentInflowChannelOptionListByAgentID(input GetAgentInflowChannelOptionListByAgentIDInput) (GetAgentInflowChannelOptionListByAgentIDOutput, error) {
	var (
		output                       GetAgentInflowChannelOptionListByAgentIDOutput
		agentInflowChannelOptionList []*entity.AgentInflowChannelOption
		err                          error
	)

	// エージェントの流入経路マスタを取得
	agentInflowChannelOptionList, err = i.agentInflowChannelOptionRepository.GetByAgentID(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.AgentInflowChannelOptionList = agentInflowChannelOptionList

	return output, nil
}
