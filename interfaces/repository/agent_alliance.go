package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentAllianceRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewAgentAllianceRepositoryImpl(ex interfaces.SQLExecuter) usecase.AgentAllianceRepository {
	return &AgentAllianceRepositoryImpl{
		Name:     "AgentAllianceRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//

func (repo *AgentAllianceRepositoryImpl) Create(agentAlliance *entity.AgentAlliance) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO agent_alliances (
			agent1_id,
			agent2_id,
			agent1_request,
			agent2_request,
			agent1_cancel_request,
			agent2_cancel_request,
			created_at,
			updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)
		`,
		agentAlliance.Agent1ID,
		agentAlliance.Agent2ID,
		agentAlliance.Agent1Request,
		agentAlliance.Agent2Request,
		false,
		false,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	agentAlliance.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//

func (repo *AgentAllianceRepositoryImpl) Update(id uint, agentAlliance *entity.AgentAlliance) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`UPDATE 
			agent_alliances 
		SET
			agent1_request = ?,
			agent2_request = ?,
			agent1_cancel_request = ?,
			agent2_cancel_request = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		agentAlliance.Agent1Request,
		agentAlliance.Agent2Request,
		agentAlliance.Agent1CancelRequest,
		agentAlliance.Agent2CancelRequest,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 指定IDのレコードの指定のエージェントIDと合致する申請状況を更新
func (repo *AgentAllianceRepositoryImpl) UpdateAgentRequest(id uint, agentID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgentRequest",
		`UPDATE 
			agent_alliances 
		SET
			agent1_request = (
				CASE agent1_id WHEN ? 
					THEN ? 
					ELSE agent1_request 
					END
			),
			agent2_request = (
				CASE agent2_id WHEN ? 
					THEN ? 
					ELSE agent2_request 
					END
			),
			updated_at = ?
		WHERE 
			id = ?
		`,
		agentID,
		true,
		agentID,
		true,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/****************************************************************************************/
/// 単数取得 API
//

func (repo *AgentAllianceRepositoryImpl) FindByID(id uint) (*entity.AgentAlliance, error) {
	var (
		agentAlliance entity.AgentAlliance
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&agentAlliance, `
		SELECT *
		FROM 
			agent_alliances
		WHERE
			id = ?
		`, id,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &agentAlliance, nil
}

func (repo *AgentAllianceRepositoryImpl) FindByAgentID(agent1ID uint, agent2ID uint) (*entity.AgentAlliance, error) {
	var (
		agentAlliance entity.AgentAlliance
	)

	err := repo.executer.Get(
		repo.Name+".FindByAgentID",
		&agentAlliance, `
		SELECT *
		FROM 
			agent_alliances
		WHERE (
			agent1_id = ? AND
			agent2_id = ?
		) OR (
			agent2_id = ? AND
			agent1_id = ?
		)
		`,
		agent1ID, agent2ID,
		agent1ID, agent2ID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &agentAlliance, nil
}

/****************************************************************************************/
/// 複数取得 API
//

func (repo *AgentAllianceRepositoryImpl) GetByAgentIDAndRequestDone(agentID uint) ([]*entity.AgentAlliance, error) {
	var (
		agentAlliance []*entity.AgentAlliance
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDAndRequestDone",
		&agentAlliance, `
		SELECT 
		  alliance.*,
			agents.agent_name,
			agents.office_location,
			agents.representative,
			agents.establish
		FROM agent_alliances AS alliance
		INNER JOIN agents
		ON (
			(
				alliance.agent1_id = ? AND
				alliance.agent2_id = agents.id
			) OR (
				alliance.agent2_id = ? AND
				alliance.agent1_id = agents.id
			)
		)
		WHERE 
			(agent1_id = ? OR agent2_id = ?)
		AND 
			agent1_request = TRUE
    	AND 
			agent2_request = TRUE
		`,
		agentID, agentID,
		agentID, agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentAlliance, nil
}

func (repo *AgentAllianceRepositoryImpl) GetByMyAgentIDAndOtherIDList(myAgentID uint, otherAgentIDList []uint) ([]*entity.AgentAlliance, error) {
	var (
		allianceList []*entity.AgentAlliance
	)

	if len(otherAgentIDList) < 1 {
		return allianceList, nil
	}

	otherAgentIDListStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(otherAgentIDList)), ", "), "[]")

	query := fmt.Sprintf(`
		SELECT *
		FROM agent_alliances
		WHERE (
			agent1_id IN(%s) AND
			agent2_id = ?
		) OR (
			agent1_id = ? AND
			agent2_id IN(%s)
		)
	`, otherAgentIDListStr, otherAgentIDListStr)

	err := repo.executer.Select(
		repo.Name+".GetByMyAgentIDAndOtherIDList",
		&allianceList, query,
		myAgentID, myAgentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return allianceList, nil
}

func (repo *AgentAllianceRepositoryImpl) All() ([]*entity.AgentAlliance, error) {
	var (
		agentAlliance []*entity.AgentAlliance
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&agentAlliance, `
		SELECT *
		FROM agent_alliances
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentAlliance, nil
}
