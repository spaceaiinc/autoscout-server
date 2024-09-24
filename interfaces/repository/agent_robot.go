package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

// エージェントのRPAロボット scout_serviceを子テーブルとして持つ
type AgentRobotRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewAgentRobotRepositoryImpl(ex interfaces.SQLExecuter) usecase.AgentRobotRepository {
	return &AgentRobotRepositoryImpl{
		Name:     "AgentRobotRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
//エージェントロボットの作成
func (repo *AgentRobotRepositoryImpl) Create(agentRobot *entity.AgentRobot) error {
	nowTime := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO agent_robots (
			uuid,
			agent_id,
			name,
			is_entry_active,
			is_scout_active,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?
			)`,
		utility.CreateUUID(),
		agentRobot.AgentID,
		agentRobot.Name,
		agentRobot.IsEntryActive,
		agentRobot.IsScoutActive,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	agentRobot.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
// エージェントロボットの更新
func (repo *AgentRobotRepositoryImpl) Update(id uint, agentRobot *entity.AgentRobot) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE agent_robots
		SET
			name = ?,
			is_entry_active = ?,
			is_scout_active = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		agentRobot.Name,
		agentRobot.IsEntryActive,
		agentRobot.IsScoutActive,
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
/// 削除 API
//
// エージェントロボットの削除
func (repo *AgentRobotRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE 
		FROM agent_robots
		WHERE id = ?
		`,
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
// IDを使ってエージェントロボットを取得
func (repo *AgentRobotRepositoryImpl) FindByID(id uint) (*entity.AgentRobot, error) {
	var (
		agentRobot entity.AgentRobot
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&agentRobot, `
		SELECT 
			robots.*,
			agents.agent_name AS agent_name
		FROM agent_robots AS robots
		INNER JOIN agents AS agents
		ON robots.agent_id = agents.id
		WHERE
			robots.id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &agentRobot, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// agentIDを使ってエージェントロボットの一覧を取得
func (repo *AgentRobotRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.AgentRobot, error) {
	var (
		agentRobotList []*entity.AgentRobot
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&agentRobotList, `
		SELECT 
			robots.*,
			agents.agent_name AS agent_name
		FROM agent_robots AS robots
		INNER JOIN agents AS agents
		ON robots.agent_id = agents.id
		WHERE
			robots.agent_id = ?
			ORDER BY 
				id ASC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentRobotList, nil
}

// 全てのエージェントロボットを取得
func (repo *AgentRobotRepositoryImpl) All() ([]*entity.AgentRobot, error) {
	var (
		agentRobotList []*entity.AgentRobot
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&agentRobotList, `
		SELECT 
			robots.*,
			agents.agent_name AS agent_name
		FROM agent_robots AS robots
		INNER JOIN agents AS agents
		ON robots.agent_id = agents.id
		ORDER BY 
			id DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentRobotList, nil
}
