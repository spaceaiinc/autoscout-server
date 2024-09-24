package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type InterviewAdjustmentTemplateRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewInterviewAdjustmentTemplateRepositoryImpl(ex interfaces.SQLExecuter) usecase.InterviewAdjustmentTemplateRepository {
	return &InterviewAdjustmentTemplateRepositoryImpl{
		Name:     "InterviewAdjustmentTemplateRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// 面談調整テンプレートを作成
func (repo *InterviewAdjustmentTemplateRepositoryImpl) Create(template *entity.InterviewAdjustmentTemplate) error {
	nowTime := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO interview_adjustment_templates (
			agent_id,
			send_scene,
			title,
			subject,
			content,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?
			)
		`,
		template.AgentID,
		template.SendScene,
		template.Title,
		template.Subject,
		template.Content,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	template.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
// 面談調整テンプレートを更新
func (repo *InterviewAdjustmentTemplateRepositoryImpl) Update(id uint, template *entity.InterviewAdjustmentTemplate) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`UPDATE interview_adjustment_templates SET
			agent_id = ?,
			send_scene = ?,
			title = ?,
			subject = ?,
			content = ?,
			updated_at = ?
		WHERE id = ?
		`,
		template.AgentID,
		template.SendScene,
		template.Title,
		template.Subject,
		template.Content,
		nowTime,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/****************************************************************************************/
/// 削除
//
// 面談調整テンプレートを削除
func (repo *InterviewAdjustmentTemplateRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`DELETE FROM interview_adjustment_templates WHERE id = ?`,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/****************************************************************************************/
/// 単数取得
//
// IDから面談調整テンプレートを取得
func (repo *InterviewAdjustmentTemplateRepositoryImpl) FindByID(id uint) (*entity.InterviewAdjustmentTemplate, error) {
	var (
		template entity.InterviewAdjustmentTemplate
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&template, `
		SELECT 
			*
		FROM 
			interview_adjustment_templates
		WHERE
			id = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &template, nil
}

/****************************************************************************************/
/// 複数取得
//

// エージェントIDから面談調整テンプレートを取得
func (repo *InterviewAdjustmentTemplateRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.InterviewAdjustmentTemplate, error) {
	var (
		templateList []*entity.InterviewAdjustmentTemplate
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&templateList, `
		SELECT 
			*
		FROM 
			interview_adjustment_templates
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

	return templateList, nil
}

// グループIDから面談調整テンプレートを取得
func (repo *InterviewAdjustmentTemplateRepositoryImpl) GetByAgentIDAndSendScene(agentID, sendScene uint) ([]*entity.InterviewAdjustmentTemplate, error) {
	var (
		templateList []*entity.InterviewAdjustmentTemplate
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&templateList, `
		SELECT 
			*
		FROM 
			interview_adjustment_templates
		WHERE
			agent_id = ?
		AND 
			send_scene = ?
		ORDER BY
			id ASC
		`,
		agentID,
		sendScene,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return templateList, nil
}
