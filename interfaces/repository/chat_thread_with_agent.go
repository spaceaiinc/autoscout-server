package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatThreadWithAgentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewChatThreadWithAgentRepositoryImpl(ex interfaces.SQLExecuter) usecase.ChatThreadWithAgentRepository {
	return &ChatThreadWithAgentRepositoryImpl{
		Name:     "ChatThreadWithAgentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *ChatThreadWithAgentRepositoryImpl) Create(chatThreadWithAgent *entity.ChatThreadWithAgent) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO chat_thread_with_agents (
			uuid,
			group_id,
			agent_staff_id,
			title,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)`,
		utility.CreateUUID(),
		chatThreadWithAgent.GroupID,
		chatThreadWithAgent.AgentStaffID,
		chatThreadWithAgent.Title,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	chatThreadWithAgent.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 単数取得
//
func (repo *ChatThreadWithAgentRepositoryImpl) FindByID(threadID uint) (*entity.ChatThreadWithAgent, error) {
	var (
		chatThreadWithAgent entity.ChatThreadWithAgent
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&chatThreadWithAgent, `
		SELECT 
			chat_thread.*,
			staff.staff_name AS staff_name
		FROM 
			chat_thread_with_agents AS chat_thread
		INNER JOIN 
			agent_staffs AS staff
		ON
			chat_thread.agent_staff_id = staff.id
		WHERE chat_thread.id = ?
		`,
		threadID)

	if err != nil {
		return nil, err
	}

	return &chatThreadWithAgent, nil
}

/****************************************************************************************/
/// 複数取得
//
// チャットグループIDからスレッド情報を取得する
func (repo *ChatThreadWithAgentRepositoryImpl) GetByGroupID(groupID uint) ([]*entity.ChatThreadWithAgent, error) {
	var (
		chatThreadWithAgent []*entity.ChatThreadWithAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetByGroupID",
		&chatThreadWithAgent, `
		SELECT 
				chat_thread.*,
				COUNT(chat_message.id) AS message_count,
				-- 最新メッセージの作成日時を返す
				IFNULL(MAX(chat_message.created_at), chat_thread.created_at) AS latest_message_time,
				staff.staff_name AS staff_name
		FROM 
				chat_thread_with_agents AS chat_thread
		LEFT JOIN
				chat_message_with_agents AS chat_message
		ON
				chat_thread.id = chat_message.thread_id
		INNER JOIN 
				agent_staffs AS staff
		ON
				chat_thread.agent_staff_id = staff.id
		WHERE chat_thread.group_id = ?
		GROUP BY chat_thread.id
		ORDER BY chat_thread.created_at DESC
		`,
		groupID)

	if err != nil {
		return nil, err
	}

	return chatThreadWithAgent, nil
}

func (repo *ChatThreadWithAgentRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.ChatThreadWithAgent, error) {
	var (
		chatThreadWithAgent []*entity.ChatThreadWithAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetByGroupID",
		&chatThreadWithAgent, `
		SELECT 
			chat_thread.*,
			staff.staff_name AS staff_name
		FROM 
			chat_thread_with_agents AS chat_thread
		INNER JOIN 
			agent_staffs AS staff
		ON
			chat_thread.agent_staff_id = staff.id
		WHERE
			chat_thread.group_id IN (
				SELECT id
				FROM chat_group_with_agents
				WHERE agent1_id = ? OR agent2_id = ?
			)
		ORDER BY chat_thread.created_at DESC
		`,
		agentID,
		agentID,
	)

	if err != nil {
		return nil, err
	}

	return chatThreadWithAgent, nil
}

// All
func (repo *ChatThreadWithAgentRepositoryImpl) All() ([]*entity.ChatThreadWithAgent, error) {
	var (
		chatThreadWithAgent []*entity.ChatThreadWithAgent
	)

	err := repo.executer.Select(
		repo.Name+" All",
		&chatThreadWithAgent, `
		SELECT 
			chat_thread.*,
			staff.staff_name AS staff_name
		FROM 
			chat_thread_with_agents AS chat_thread
		INNER JOIN 
			agent_staffs AS staff
		ON
			chat_thread.agent_staff_id = staff.id
		ORDER BY chat_thread.created_at DESC
		`,
	)

	if err != nil {
		return nil, err
	}

	return chatThreadWithAgent, nil
}
