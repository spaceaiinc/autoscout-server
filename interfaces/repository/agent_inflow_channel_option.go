package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentInflowChannelOptionRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewAgentInflowChannelOptionRepositoryImpl(ex interfaces.SQLExecuter) usecase.AgentInflowChannelOptionRepository {
	return &AgentInflowChannelOptionRepositoryImpl{
		Name:     "inflowChannelOptionRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
//エージェントの流入経路マスタの作成
func (repo *AgentInflowChannelOptionRepositoryImpl) Create(inflowChannelOption *entity.AgentInflowChannelOption) error {
	nowTime := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO agent_inflow_channel_options (
			agent_id,
			channel_name,
			is_open,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)`,
		inflowChannelOption.AgentID,
		inflowChannelOption.ChannelName,
		inflowChannelOption.IsOpen,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	inflowChannelOption.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
// エージェントの流入経路マスタの更新
func (repo *AgentInflowChannelOptionRepositoryImpl) Update(id uint, inflowChannelOption *entity.AgentInflowChannelOption) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE agent_inflow_channel_options
		SET
			channel_name = ?,
			is_open = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		inflowChannelOption.ChannelName,
		inflowChannelOption.IsOpen,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 単数取得 API
//
// IDからエージェントの流入経路マスタの取得
func (repo *AgentInflowChannelOptionRepositoryImpl) FindByID(id uint) (*entity.AgentInflowChannelOption, error) {
	var (
		inflowChannelOption entity.AgentInflowChannelOption
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&inflowChannelOption, `
		SELECT 
			*
		FROM 
		  agent_inflow_channel_options
		WHERE
			id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &inflowChannelOption, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// agentIDを使ってエージェントの流入経路マスタの一覧を取得
func (repo *AgentInflowChannelOptionRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.AgentInflowChannelOption, error) {
	var (
		inflowChannelOptionList []*entity.AgentInflowChannelOption
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&inflowChannelOptionList, `
		SELECT 
			*
		FROM 
		  agent_inflow_channel_options
		WHERE
			agent_id = ?
		ORDER BY 
			id ASC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return inflowChannelOptionList, nil
}
