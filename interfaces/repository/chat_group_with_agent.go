package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatGroupWithAgentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewChatGroupWithAgentRepositoryImpl(ex interfaces.SQLExecuter) usecase.ChatGroupWithAgentRepository {
	return &ChatGroupWithAgentRepositoryImpl{
		Name:     "ChatGroupWithAgentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
//エージェントと求職者のチャットグループの作成
func (repo *ChatGroupWithAgentRepositoryImpl) Create(chatGroupWithAgent *entity.ChatGroupWithAgent) error {
	nowTime := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO chat_group_with_agents (
			uuid,
			agent1_id,
			agent2_id,
			last_send_at,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)`,
		utility.CreateUUID(),
		chatGroupWithAgent.Agent1ID,
		chatGroupWithAgent.Agent2ID,
		nowTime,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	chatGroupWithAgent.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
func (repo *ChatGroupWithAgentRepositoryImpl) UpdateAgentLastWatchedAt(groupID uint, agentID uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgentLastWatchedAt",
		`
		UPDATE 
			chat_group_with_agents 
		SET
			agent1_last_watched_at = (
				CASE agent1_id WHEN ? 
				THEN ? 
				ELSE chat_group_with_agents.agent1_last_watched_at
				END
			),
			agent2_last_watched_at = (
				CASE agent2_id WHEN ? 
				THEN ? 
				ELSE chat_group_with_agents.agent2_last_watched_at
				END
			),
			updated_at = ?
		WHERE 
			id = ?
		`,
		agentID,
		nowTime,
		agentID,
		nowTime,
		nowTime,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// メッセージの送受信ごとにLastSendAtを更新する
func (repo *ChatGroupWithAgentRepositoryImpl) UpdateLastSendAtByThreadID(threadID uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateLastSendAtByThreadID",
		`
			UPDATE 
				chat_group_with_agents 
			SET
				last_send_at = ?,
				updated_at = ?
			WHERE 
				id = (
					SELECT group_id 
					FROM chat_thread_with_agents 
					WHERE id = ?
				)
			`,
		nowTime,
		nowTime,
		threadID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 単数取得
//
func (repo *ChatGroupWithAgentRepositoryImpl) FindByID(id uint) (*entity.ChatGroupWithAgent, error) {
	var (
		chatGroupWithAgent entity.ChatGroupWithAgent
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&chatGroupWithAgent, `
		SELECT 
			chat_group.*
		FROM 
			chat_group_with_agents AS chat_group
		WHERE
			chat_group.id = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &chatGroupWithAgent, nil
}

func (repo *ChatGroupWithAgentRepositoryImpl) FindByIDAndAgentID(id, agentID uint) (*entity.ChatGroupWithAgent, error) {
	var (
		chatGroupWithAgent entity.ChatGroupWithAgent
	)

	err := repo.executer.Get(
		repo.Name+".FindByIDAndAgentID",
		&chatGroupWithAgent, `
		SELECT 
			chat_group.*,
		CASE 
			WHEN 
				agent1.id = ? 
			THEN 
				agent2.agent_name
			WHEN 
				agent2.id = ? 
			THEN 
				agent1.agent_name
			END AS agent_name
		FROM 
			chat_group_with_agents AS chat_group
		INNER JOIN
			agents AS agent1
		ON
			chat_group.agent1_id = agent1.id
		INNER JOIN
			agents AS agent2
		ON
			chat_group.agent2_id = agent2.id
		WHERE
			chat_group.id = ?
		LIMIT 1
		`,
		agentID,
		agentID,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &chatGroupWithAgent, nil
}

func (repo *ChatGroupWithAgentRepositoryImpl) FindByAgentID(agent1ID, agent2ID uint) (*entity.ChatGroupWithAgent, error) {
	var (
		chatGroupWithAgent entity.ChatGroupWithAgent
	)

	err := repo.executer.Get(
		repo.Name+".FindByAgentID",
		&chatGroupWithAgent, `
		SELECT 
			*
		FROM 
			chat_group_with_agents
		WHERE (
			agent1_id = ? AND
			agent2_id = ?
		) OR (
			agent2_id = ? AND
			agent1_id = ?
		)
		LIMIT 1
		`,
		agent1ID, agent2ID,
		agent1ID, agent2ID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &chatGroupWithAgent, nil
}

/****************************************************************************************/
/// 複数取得
//
//agentStaffIDを使ってエージェントと求職者のチャットグループグループの一覧を取得
func (repo *ChatGroupWithAgentRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.ChatGroupWithAgent, error) {
	var (
		chatGroupWithAgentList []*entity.ChatGroupWithAgent
	)

	err := repo.executer.Select(
		repo.Name+"GetByAgentID",
		&chatGroupWithAgentList, `
		SELECT 
		chat_group.*,
		CASE 
			WHEN 
				agent1.id = ? 
			THEN 
				agent2.id
			WHEN 
				agent2.id = ? 
			THEN 
				agent1.id
		END 
			AS agent_id,
		CASE 
			WHEN 
				agent1.id = ? 
			THEN 
				agent2.agent_name
			WHEN 
				agent2.id = ? 
			THEN 
				agent1.agent_name
		END 
			AS agent_name
		FROM 
			chat_group_with_agents AS chat_group
		INNER JOIN
			agents AS agent1
		ON
			chat_group.agent1_id = agent1.id
		INNER JOIN
			agents AS agent2
		ON
			chat_group.agent2_id = agent2.id
		WHERE
			chat_group.agent1_id = ?
		OR
			chat_group.agent2_id = ?
		ORDER BY 
			(chat_group.agent1_id = ? AND chat_group.agent2_id = ?) DESC, 
			chat_group.last_send_at DESC
		`,
		agentID,
		agentID,
		agentID,
		agentID,
		agentID,
		agentID,
		agentID,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return chatGroupWithAgentList, nil
}

func (repo *ChatGroupWithAgentRepositoryImpl) GetByMyAgentIDAndOtherIDList(myAgentID uint, otherAgentIDList []uint) ([]*entity.ChatGroupWithAgent, error) {
	var (
		chatGroupWithAgentList []*entity.ChatGroupWithAgent
	)

	if len(otherAgentIDList) == 0 {
		return chatGroupWithAgentList, nil
	}

	otherAgentIDListStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(otherAgentIDList)), ","), "[]")

	query := fmt.Sprintf(`
		SELECT 
			*
		FROM 
			chat_group_with_agents
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
		&chatGroupWithAgentList, query,
		myAgentID, myAgentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return chatGroupWithAgentList, nil
}

// All
func (repo *ChatGroupWithAgentRepositoryImpl) All() ([]*entity.ChatGroupWithAgent, error) {
	var (
		chatGroupWithAgentList []*entity.ChatGroupWithAgent
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&chatGroupWithAgentList, `
		SELECT 
			chat_group.*,
			agent1.agent_name AS agent1_name,
			agent2.agent_name AS agent2_name
		FROM 
			chat_group_with_agents AS chat_group
		INNER JOIN
			agents AS agent1
		ON
			chat_group.agent1_id = agent1.id
		INNER JOIN
			agents AS agent2
		ON
			chat_group.agent2_id = agent2.id
		ORDER BY
			chat_group.id DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return chatGroupWithAgentList, nil
}
