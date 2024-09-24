package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatMessageToUserWithAgentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewChatMessageToUserWithAgentRepositoryImpl(ex interfaces.SQLExecuter) usecase.ChatMessageToUserWithAgentRepository {
	return &ChatMessageToUserWithAgentRepositoryImpl{
		Name:     "ChatMessageToUserWithAgentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *ChatMessageToUserWithAgentRepositoryImpl) Create(chatMessageToUserWithAgent *entity.ChatMessageToUserWithAgent) error {
	nowTime := time.Now().In(time.UTC)
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO chat_message_to_user_with_agents (
			message_id,
			agent_staff_id,
			send_at,
			watched_at,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)`,
		chatMessageToUserWithAgent.MessageID,
		chatMessageToUserWithAgent.AgentStaffID,
		nowTime,
		nowTime,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	chatMessageToUserWithAgent.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
func (repo *ChatMessageToUserWithAgentRepositoryImpl) UpdateWatchedAtByThreadID(threadID, agentStaffID uint) error {
	nowTime := time.Now().In(time.UTC)
	_, err := repo.executer.Exec(
		repo.Name+".UpdateWatchedAtByThreadID",
		`
		UPDATE 
			chat_message_to_user_with_agents 
		SET
			watched_at = ?,
			updated_at = ?
		WHERE message_id IN (
			SELECT id
			FROM chat_message_with_agents
			WHERE thread_id = ?
		)
		AND
			agent_staff_id = ?
		`,
		nowTime,
		nowTime,
		threadID,
		agentStaffID,
	)

	if err != nil {
		return err
	}

	return nil
}

/****************************************************************************************/
/// 複数取得
//
func (repo *ChatMessageToUserWithAgentRepositoryImpl) GetByThreadID(threadID uint) ([]*entity.ChatMessageToUserWithAgent, error) {
	var (
		chatMessageToUserWithAgent []*entity.ChatMessageToUserWithAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetByThreadID",
		&chatMessageToUserWithAgent, `
		SELECT *
		FROM chat_message_to_user_with_agents
		WHERE message_id IN (
			SELECT id
			FROM chat_message_with_agents
			WHERE thread_id = ?
		)
		ORDER BY created_at DESC
		`,
		threadID,
	)

	if err != nil {
		return nil, err
	}

	return chatMessageToUserWithAgent, nil
}

func (repo *ChatMessageToUserWithAgentRepositoryImpl) GetByGroupIDAndUnwatched(groupID, agentStaffID uint) ([]*entity.ChatMessageToUserWithAgent, error) {
	var (
		chatMessageToUserWithAgent []*entity.ChatMessageToUserWithAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetByGroupIDAndUnwatched",
		&chatMessageToUserWithAgent, `
		SELECT 
			chat_message_to_user.*, 
			chat_thread.id AS thread_id
		FROM 
			chat_message_to_user_with_agents AS chat_message_to_user
		INNER JOIN 
			chat_message_with_agents AS chat_message
		ON 
			chat_message_to_user.message_id = chat_message.id
		INNER JOIN 
			chat_thread_with_agents AS chat_thread
		ON 
			chat_message.thread_id = chat_thread.id
		WHERE 
			chat_thread.group_id = ?
		AND
			chat_message_to_user.agent_staff_id = ? 
		AND 
			(watched_at <= send_at OR watched_at IS NULL)
		ORDER BY created_at DESC
		`,
		groupID,
		agentStaffID,
	)

	if err != nil {
		return nil, err
	}

	return chatMessageToUserWithAgent, nil
}
func (repo *ChatMessageToUserWithAgentRepositoryImpl) GetByAgentStaffIDAndUnwatched(agentStaffID uint) ([]*entity.ChatMessageToUserWithAgent, error) {
	var (
		chatMessageToUserWithAgent []*entity.ChatMessageToUserWithAgent
	)
	err := repo.executer.Select(
		repo.Name+".GetByAgentStaffIDAndUnwatched",
		&chatMessageToUserWithAgent, `
		SELECT 
			chat_message_to_user.*, 
			chat_thread.group_id AS group_id
		FROM 
			chat_message_to_user_with_agents AS chat_message_to_user
		INNER JOIN 
			chat_message_with_agents AS chat_message
		ON 
			chat_message_to_user.message_id = chat_message.id
		INNER JOIN 
			chat_thread_with_agents AS chat_thread
		ON 
			chat_message.thread_id = chat_thread.id
		WHERE 
			chat_message_to_user.agent_staff_id = ? 
		AND 
			(watched_at <= send_at OR watched_at IS NULL)
		ORDER BY created_at DESC
		`,
		agentStaffID,
	)
	if err != nil {
		return nil, err
	}
	return chatMessageToUserWithAgent, nil
}

func (repo *ChatMessageToUserWithAgentRepositoryImpl) GetByMessageIDList(idList []uint) ([]*entity.ChatMessageToUserWithAgent, error) {
	var (
		toUserList []*entity.ChatMessageToUserWithAgent
	)

	if len(idList) == 0 {
		return toUserList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			chat_message_to_user.*,
			staff.staff_name
		FROM 
			chat_message_to_user_with_agents AS chat_message_to_user
		INNER JOIN 
			agent_staffs AS staff
		ON 
			chat_message_to_user.agent_staff_id = staff.id
		WHERE 
			chat_message_to_user.message_id IN (%s)
		ORDER BY id ASC
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ","), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDList",
		&toUserList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return toUserList, nil
}

// GetAll
func (repo *ChatMessageToUserWithAgentRepositoryImpl) All() ([]*entity.ChatMessageToUserWithAgent, error) {
	var (
		chatMessageToUserWithAgent []*entity.ChatMessageToUserWithAgent
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&chatMessageToUserWithAgent, `
		SELECT *
		FROM chat_message_to_user_with_agents
		ORDER BY created_at DESC
		`,
	)

	if err != nil {
		return nil, err
	}

	return chatMessageToUserWithAgent, nil
}
