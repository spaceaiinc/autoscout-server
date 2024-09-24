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

type ScoutServiceTemplateRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewScoutServiceTemplateRepositoryImpl(ex interfaces.SQLExecuter) usecase.ScoutServiceTemplateRepository {
	return &ScoutServiceTemplateRepositoryImpl{
		Name:     "ScoutServiceTemplateRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
//スカウトテンプレートの作成
func (repo *ScoutServiceTemplateRepositoryImpl) Create(scoutServiceTemplate *entity.ScoutServiceTemplate) error {
	nowTime := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO scout_service_templates (
			scout_service_id,
			start_hour,
			start_minute,
			run_on_monday,
			run_on_tuesday,

			run_on_wednesday,
			run_on_thursday,
			run_on_friday,
			run_on_saturday,
			run_on_sunday,

			scout_count,
			search_title,
			message_title,
			job_information_title,
			job_information_id,

			age_limit,
			scout_type,
			auto_remind,
			reply_limit,
			last_send_count,

			last_send_at,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?
			)`,
		scoutServiceTemplate.ScoutServiceID,
		scoutServiceTemplate.StartHour,
		scoutServiceTemplate.StartMinute,
		scoutServiceTemplate.RunOnMonday,
		scoutServiceTemplate.RunOnTuesday,
		scoutServiceTemplate.RunOnWednesday,
		scoutServiceTemplate.RunOnThursday,
		scoutServiceTemplate.RunOnFriday,
		scoutServiceTemplate.RunOnSaturday,
		scoutServiceTemplate.RunOnSunday,
		scoutServiceTemplate.ScoutCount,
		scoutServiceTemplate.SearchTitle,
		scoutServiceTemplate.MessageTitle,
		scoutServiceTemplate.JobInformationTitle,
		scoutServiceTemplate.JobInformationID,
		scoutServiceTemplate.AgeLimit,
		scoutServiceTemplate.ScoutType,
		scoutServiceTemplate.AutoRemind,
		scoutServiceTemplate.ReplyLimit,
		0,
		utility.EarliestTime(),
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	scoutServiceTemplate.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
// 最終送信日時と最終送信件数の更新
func (repo *ScoutServiceTemplateRepositoryImpl) UpdateLastSend(scoutServiceTemplateID, lastSendCount uint, lastSendAt time.Time) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE scout_service_templates
		SET
			last_send_count = ?,
			last_send_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		lastSendCount,
		lastSendAt,
		time.Now().In(time.UTC),
		scoutServiceTemplateID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 削除
//
// スカウトサービスIDからスカウトテンプレートの削除
func (repo *ScoutServiceTemplateRepositoryImpl) DeleteByScoutServiceID(scoutServiceID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE 
		FROM scout_service_templates
		WHERE scout_service_id = ?
		`,
		scoutServiceID,
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
func (repo *ScoutServiceTemplateRepositoryImpl) FindByID(id uint) (*entity.ScoutServiceTemplate, error) {
	var (
		scoutServiceTemplate entity.ScoutServiceTemplate
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&scoutServiceTemplate, `
		SELECT 
			templates.*,
			services.login_id,
			services.password,
			services.service_type
		FROM scout_service_templates AS templates
		INNER JOIN scout_services AS services
		ON templates.scout_service_id = services.id
		WHERE
			templates.id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &scoutServiceTemplate, nil
}

/****************************************************************************************/
/// 複数取得
//
// agentIDを使ってスカウトテンプレートの一覧を取得
func (repo *ScoutServiceTemplateRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.ScoutServiceTemplate, error) {
	var (
		scoutServiceTemplateList []*entity.ScoutServiceTemplate
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&scoutServiceTemplateList, `
			SELECT *
			FROM scout_service_templates
			WHERE
				scout_service_id IN (
					SELECT id
					FROM scout_services
					WHERE agent_robot_id IN (
						SELECT id
						FROM agent_robots
						WHERE agent_id = ?
					)
				)
			ORDER BY 
				id DESC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return scoutServiceTemplateList, nil
}

// scoutServiceIDを使ってスカウトテンプレートの一覧を取得
func (repo *ScoutServiceTemplateRepositoryImpl) GetByScoutServiceID(scoutServiceID uint) ([]*entity.ScoutServiceTemplate, error) {
	var (
		scoutServiceTemplateList []*entity.ScoutServiceTemplate
	)

	err := repo.executer.Select(
		repo.Name+".GetByScoutServiceID",
		&scoutServiceTemplateList, `
			SELECT *
			FROM scout_service_templates
			WHERE
				scout_service_id = ?
			ORDER BY 
				id DESC
		`,
		scoutServiceID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return scoutServiceTemplateList, nil
}

// agentRobotIDを使ってスカウトテンプレートの一覧を取得
func (repo *ScoutServiceTemplateRepositoryImpl) GetByAgentRobotID(agentRobotID uint) ([]*entity.ScoutServiceTemplate, error) {
	var (
		scoutServiceTemplateList []*entity.ScoutServiceTemplate
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentRobotID",
		&scoutServiceTemplateList, `
			SELECT *
			FROM scout_service_templates
			WHERE
				scout_service_id IN (
					SELECT id
					FROM scout_services
					WHERE agent_robot_id = ?
				)
			ORDER BY 
				id DESC
		`,
		agentRobotID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return scoutServiceTemplateList, nil
}

// IDリストを使ってスカウトテンプレートの一覧を取得
func (repo *ScoutServiceTemplateRepositoryImpl) GetByIDList(idList []uint) ([]*entity.ScoutServiceTemplate, error) {
	var (
		scoutServiceTemplateList []*entity.ScoutServiceTemplate
	)

	if len(idList) == 0 {
		return scoutServiceTemplateList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM scout_service_templates
		WHERE
			id IN(%s)
		ORDER BY
			id DESC
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByAgentRobotID",
		&scoutServiceTemplateList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return scoutServiceTemplateList, nil
}
