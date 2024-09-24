package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type MessageTemplateRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewMessageTemplateRepositoryImpl(ex interfaces.SQLExecuter) usecase.MessageTemplateRepository {
	return &MessageTemplateRepositoryImpl{
		Name:     "MessageTemplateRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// メッセージテンプレートの作成
func (repo *MessageTemplateRepositoryImpl) Create(messageTemplate *entity.MessageTemplate) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO message_templates (
			agent_staff_id,
			send_scene,
			title,
			subject,
			content,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?
			)`,
		messageTemplate.AgentStaffID,
		messageTemplate.SendScene,
		messageTemplate.Title,
		messageTemplate.Subject,
		messageTemplate.Content,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	messageTemplate.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
// メッセージテンプレートの更新
func (repo *MessageTemplateRepositoryImpl) Update(id uint, messageTemplate *entity.MessageTemplate) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE message_templates
		SET
			send_scene = ?,
			title = ?,
			subject = ?,
			content = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		messageTemplate.SendScene,
		messageTemplate.Title,
		messageTemplate.Subject,
		messageTemplate.Content,
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
// メッセージテンプレートの削除
func (repo *MessageTemplateRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM message_templates
		WHERE id = ?
		`, id,
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
// メッセージテンプレートの取得
func (repo *MessageTemplateRepositoryImpl) FindByID(id uint) (*entity.MessageTemplate, error) {
	var (
		messageTemplate entity.MessageTemplate
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&messageTemplate, `
		SELECT *
		FROM message_templates
		WHERE
			id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &messageTemplate, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// 担当者IDからメッセージテンプレートを取得
func (repo *MessageTemplateRepositoryImpl) GetByAgentStaffID(agentStaffID uint) ([]*entity.MessageTemplate, error) {
	var (
		list []*entity.MessageTemplate
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentStaffID",
		&list, `
			SELECT *
			FROM message_templates
			WHERE
				agent_staff_id = ?
			ORDER BY id ASC
		`,
		agentStaffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}

// 担当者IDと送信シーンからメッセージテンプレートを取得
func (repo *MessageTemplateRepositoryImpl) GetByAgentStaffIDAndSendScene(agentStaffID, sendScene uint) ([]*entity.MessageTemplate, error) {
	var (
		list []*entity.MessageTemplate
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&list, `
		SELECT *
		FROM message_templates
		WHERE
			agent_staff_id = ?
		AND
			send_scene = ?
		`,
		agentStaffID, sendScene)

	if err != nil {
		return nil, err
	}

	return list, nil
}

// エージェントIDからメッセージテンプレートを取得
func (repo *MessageTemplateRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.MessageTemplate, error) {
	var (
		list []*entity.MessageTemplate
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&list, `
		SELECT *
		FROM message_templates
		WHERE
			agent_staff_id IN (
				SELECT id
				FROM agent_staffs
				WHERE agent_id = ?
			)
		ORDER BY id ASC
		`,
		agentID)

	if err != nil {
		return nil, err
	}

	return list, nil
}
