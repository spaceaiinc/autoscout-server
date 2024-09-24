package interactor

import (
	"log"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentRobotInteractor interface {
	// 汎用系 API
	UpdateAgentRobot(input UpdateAgentRobotInput) (UpdateAgentRobotByIDOutput, error)
	GetAgentRobotByID(input GetAgentRobotByIDInput) (GetAgentRobotByIDOutput, error)
	GetAgentRobotListByAgentID(input GetAgentRobotListByAgentIDInput) (GetAgentRobotListByAgentIDOutput, error)

	// Admin API
	CreateAgentRobot(input CreateAgentRobotInput) (CreateAgentRobotByIDOutput, error)
	DeleteAgentRobot(input DeleteAgentRobotInput) (DeleteAgentRobotByIDOutput, error)
}

type AgentRobotInteractorImpl struct {
	firebase                           usecase.Firebase
	sendgrid                           config.Sendgrid
	oneSignal                          config.OneSignal
	agentRobotRepository               usecase.AgentRobotRepository
	scoutServiceRepository             usecase.ScoutServiceRepository
	scoutServiceGetEntryTimeRepository usecase.ScoutServiceGetEntryTimeRepository
	scoutServiceTemplateRepository     usecase.ScoutServiceTemplateRepository
}

// AgentRobotInteractorImpl is an implementation of AgentRobotInteractor
func NewAgentRobotInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	arR usecase.AgentRobotRepository,
	ssR usecase.ScoutServiceRepository,
	ssgetR usecase.ScoutServiceGetEntryTimeRepository,
	sstR usecase.ScoutServiceTemplateRepository,
) AgentRobotInteractor {
	return &AgentRobotInteractorImpl{
		firebase:                           fb,
		sendgrid:                           sg,
		oneSignal:                          os,
		agentRobotRepository:               arR,
		scoutServiceRepository:             ssR,
		scoutServiceGetEntryTimeRepository: ssgetR,
		scoutServiceTemplateRepository:     sstR,
	}
}

/****************************************************************************************/
// 汎用系API
//
// エージェントロボットを更新する
type UpdateAgentRobotInput struct {
	AgentRobotID uint
	UpdateParam  entity.CreateOrUpdateAgentRobotParam
}

type UpdateAgentRobotByIDOutput struct {
	AgentRobot *entity.AgentRobot
}

func (i *AgentRobotInteractorImpl) UpdateAgentRobot(input UpdateAgentRobotInput) (UpdateAgentRobotByIDOutput, error) {
	var (
		output     UpdateAgentRobotByIDOutput
		agentRobot *entity.AgentRobot
		err        error
	)

	agentRobot = entity.NewAgentRobot(
		input.UpdateParam.AgentID,
		input.UpdateParam.Name,
		input.UpdateParam.IsEntryActive,
		input.UpdateParam.IsScoutActive,
	)

	// エージェントロボットを更新
	err = i.agentRobotRepository.Update(input.AgentRobotID, agentRobot)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.AgentRobot = agentRobot

	return output, nil
}

// エージェントロボットを取得する
type GetAgentRobotByIDInput struct {
	AgentRobotID uint
}

type GetAgentRobotByIDOutput struct {
	AgentRobot *entity.AgentRobot
}

func (i *AgentRobotInteractorImpl) GetAgentRobotByID(input GetAgentRobotByIDInput) (GetAgentRobotByIDOutput, error) {
	var (
		output     GetAgentRobotByIDOutput
		agentRobot *entity.AgentRobot
		err        error
	)

	// エージェントロボットを取得
	agentRobot, err = i.agentRobotRepository.FindByID(input.AgentRobotID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	scoutServices, err := i.scoutServiceRepository.GetByAgentRobotID(input.AgentRobotID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// エントリー取得時間を取得
	scoutServiceGetEntryTimes, err := i.scoutServiceGetEntryTimeRepository.GetByAgentRobotID(input.AgentRobotID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// テンプレートを取得
	scoutServiceTemplates, err := i.scoutServiceTemplateRepository.GetByAgentRobotID(input.AgentRobotID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	for _, scoutService := range scoutServices {
		// スカウトサービスにエントリー取得時間をマッピング
		for _, scoutServiceGetEntryTime := range scoutServiceGetEntryTimes {
			if scoutService.ID == scoutServiceGetEntryTime.ScoutServiceID {
				scoutService.GetEntryTimes = append(scoutService.GetEntryTimes, *scoutServiceGetEntryTime)
			}
		}

		// スカウトサービスにテンプレートをマッピング
		for _, scoutServiceTemplate := range scoutServiceTemplates {
			if scoutService.ID == scoutServiceTemplate.ScoutServiceID {
				scoutService.Templates = append(scoutService.Templates, *scoutServiceTemplate)
			}
		}

		// エージェントロボットにスカウトサービスをマッピング
		agentRobot.ScoutServices = append(agentRobot.ScoutServices, *scoutService)
	}

	output.AgentRobot = agentRobot

	return output, nil
}

// エージェントIDからエージェントロボットを取得する
type GetAgentRobotListByAgentIDInput struct {
	AgentID uint
}

type GetAgentRobotListByAgentIDOutput struct {
	AgentRobotList []*entity.AgentRobot
}

func (i *AgentRobotInteractorImpl) GetAgentRobotListByAgentID(input GetAgentRobotListByAgentIDInput) (GetAgentRobotListByAgentIDOutput, error) {
	var (
		output         GetAgentRobotListByAgentIDOutput
		agentRobotList []*entity.AgentRobot
		err            error
	)

	// エージェントロボットを取得
	agentRobotList, err = i.agentRobotRepository.GetByAgentID(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	scoutServices, err := i.scoutServiceRepository.GetByAgentID(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// エントリー取得時間を取得
	scoutServiceGetEntryTimes, err := i.scoutServiceGetEntryTimeRepository.GetByAgentID(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// テンプレートを取得
	scoutServiceTemplates, err := i.scoutServiceTemplateRepository.GetByAgentID(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	for _, agentRobot := range agentRobotList {
		for _, scoutService := range scoutServices {
			// スカウトサービスにエントリー取得時間をマッピング
			for _, scoutServiceGetEntryTime := range scoutServiceGetEntryTimes {
				if scoutService.ID == scoutServiceGetEntryTime.ScoutServiceID {
					scoutService.GetEntryTimes = append(scoutService.GetEntryTimes, *scoutServiceGetEntryTime)
				}
			}

			// スカウトサービスにテンプレートをマッピング
			for _, scoutServiceTemplate := range scoutServiceTemplates {
				if scoutService.ID == scoutServiceTemplate.ScoutServiceID {
					scoutService.Templates = append(scoutService.Templates, *scoutServiceTemplate)
				}
			}

			// エージェントロボットにスカウトサービスをマッピング
			agentRobot.ScoutServices = append(agentRobot.ScoutServices, *scoutService)
		}

		output.AgentRobotList = append(output.AgentRobotList, agentRobot)
	}

	return output, nil
}

/****************************************************************************************/
// ADMIN系API
//
// エージェントロボットを作成する
type CreateAgentRobotInput struct {
	CreateParam entity.CreateOrUpdateAgentRobotParam
}

type CreateAgentRobotByIDOutput struct {
	AgentRobot *entity.AgentRobot
}

func (i *AgentRobotInteractorImpl) CreateAgentRobot(input CreateAgentRobotInput) (CreateAgentRobotByIDOutput, error) {
	var (
		output     CreateAgentRobotByIDOutput
		agentRobot *entity.AgentRobot
		err        error
	)

	agentRobot = entity.NewAgentRobot(
		input.CreateParam.AgentID,
		input.CreateParam.Name,
		input.CreateParam.IsEntryActive,
		input.CreateParam.IsScoutActive,
	)

	// エージェントロボットを作成
	err = i.agentRobotRepository.Create(agentRobot)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.AgentRobot = agentRobot

	return output, nil
}

// エージェントロボットを削除する
type DeleteAgentRobotInput struct {
	AgentRobotID uint
}

type DeleteAgentRobotByIDOutput struct {
	OK bool
}

func (i *AgentRobotInteractorImpl) DeleteAgentRobot(input DeleteAgentRobotInput) (DeleteAgentRobotByIDOutput, error) {
	var (
		output DeleteAgentRobotByIDOutput
		err    error
	)

	// エージェントロボットを削除
	err = i.agentRobotRepository.Delete(input.AgentRobotID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}
