package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatMessageWithAgentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewChatMessageWithAgentRepositoryImpl(ex interfaces.SQLExecuter) usecase.ChatMessageWithAgentRepository {
	return &ChatMessageWithAgentRepositoryImpl{
		Name:     "ChatMessageWithAgentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *ChatMessageWithAgentRepositoryImpl) Create(chatMessageWithAgent *entity.ChatMessageWithAgent) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO chat_message_with_agents (
			uuid,
			thread_id,
			agent_staff_id,
			message,
			file_url,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?,?
			)`,
		utility.CreateUUID(),
		chatMessageWithAgent.ThreadID,
		chatMessageWithAgent.AgentStaffID,
		chatMessageWithAgent.Message,
		chatMessageWithAgent.FileURL,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	chatMessageWithAgent.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 複数取得
//
func (repo *ChatMessageWithAgentRepositoryImpl) GetByThreadID(threadID uint) ([]*entity.ChatMessageWithAgent, error) {
	var (
		chatMessageWithAgent []*entity.ChatMessageWithAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetByThreadID",
		&chatMessageWithAgent, `
		SELECT 
			chat_message.*,
			staffs.staff_name AS staff_name
		FROM 
			chat_message_with_agents AS chat_message
		INNER JOIN 
			agent_staffs AS staffs
		ON
			chat_message.agent_staff_id = staffs.id
		WHERE
		chat_message.thread_id = ?
		ORDER BY chat_message.created_at ASC
		`,
		threadID)

	if err != nil {
		return nil, err
	}

	return chatMessageWithAgent, nil
}

func (repo *ChatMessageWithAgentRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.ChatMessageWithAgent, error) {
	var (
		chatMessageWithAgent []*entity.ChatMessageWithAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&chatMessageWithAgent, `
		SELECT 
			chat_message.*,
			staffs.staff_name AS staff_name
		FROM chat_message_with_agents AS chat_message
		INNER JOIN 
			agent_staffs AS staffs
		ON
			chat_message.agent_staff_id = staffs.id
		WHERE
			chat_message.thread_id IN (
				SELECT id
				FROM chat_thread_with_agents
				WHERE group_id IN (
					SELECT id
					FROM chat_group_with_agents
					WHERE agent1_id = ? OR agent2_id = ?
				)
			)
		ORDER BY chat_message.created_at ASC
		`,
		agentID, agentID)

	if err != nil {
		return nil, err
	}

	return chatMessageWithAgent, nil
}

// All
func (repo *ChatMessageWithAgentRepositoryImpl) All() ([]*entity.ChatMessageWithAgent, error) {
	var (
		chatMessageWithAgent []*entity.ChatMessageWithAgent
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&chatMessageWithAgent, `
		SELECT 
			chat_message.*,
			staffs.staff_name AS staff_name
		FROM chat_message_with_agents AS chat_message
		INNER JOIN 
			agent_staffs AS staffs
		ON
			chat_message.agent_staff_id = staffs.id
		ORDER BY chat_message.created_at ASC
		`)

	if err != nil {
		return nil, err
	}

	return chatMessageWithAgent, nil
}
