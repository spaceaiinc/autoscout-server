package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ScoutServiceRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewScoutServiceRepositoryImpl(ex interfaces.SQLExecuter) usecase.ScoutServiceRepository {
	return &ScoutServiceRepositoryImpl{
		Name:     "ScoutServiceRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// スカウトサービスの作成
func (repo *ScoutServiceRepositoryImpl) Create(scoutService *entity.ScoutService) error {
	nowTime := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO scout_services (
			uuid,
			agent_robot_id,
			agent_staff_id,
			login_id,
			password,
			service_type,
			is_active,
			memo,
			template_title_for_employed,
			template_title_for_unemployed,
			interview_adjustment_template_id,
			inflow_channel_id,
			last_send_count,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?
			)`,
		utility.CreateUUID(),
		scoutService.AgentRobotID,
		scoutService.AgentStaffID,
		scoutService.LoginID,
		scoutService.Password,
		scoutService.ServiceType,
		scoutService.IsActive,
		scoutService.Memo,
		scoutService.TemplateTitleForEmployed,
		scoutService.TemplateTitleForUnemployed,
		scoutService.InterviewAdjustmentTemplateID,
		scoutService.InflowChannelID,
		0, // last_send_count
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	scoutService.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
// スカウトサービスの更新
func (repo *ScoutServiceRepositoryImpl) Update(id uint, scoutService *entity.ScoutService) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE scout_services
		SET
			agent_staff_id = ?,
			login_id = ?,
			service_type = ?,
			is_active = ?,
			memo = ?,
			template_title_for_employed = ?,
			template_title_for_unemployed = ?,
			interview_adjustment_template_id = ?,
			inflow_channel_id = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		scoutService.AgentStaffID,
		scoutService.LoginID,
		scoutService.ServiceType,
		scoutService.IsActive,
		scoutService.Memo,
		scoutService.TemplateTitleForEmployed,
		scoutService.TemplateTitleForUnemployed,
		scoutService.InterviewAdjustmentTemplateID,
		scoutService.InflowChannelID,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// スカウトサービスのパスワード更新
func (repo *ScoutServiceRepositoryImpl) UpdatePassword(param entity.UpdateScoutServicePasswordParam) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdatePassword",
		`
		UPDATE scout_services
		SET
			password = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		param.Password,
		time.Now().In(time.UTC),
		param.ScoutServiceID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// スカウトサービスの最終送信求職者のIDを更新
func (repo *ScoutServiceRepositoryImpl) UpdateLastSendCount(id, lastSendCount uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateLastSendCount",
		`
		UPDATE scout_services
		SET
			last_send_count = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		lastSendCount,
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
/// 削除
//
// スカウトサービスの削除
func (repo *ScoutServiceRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE 
		FROM scout_services
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
/// 単数取得
//
// IDを使ってスカウトサービスを取得
func (repo *ScoutServiceRepositoryImpl) FindByID(id uint) (*entity.ScoutService, error) {
	var (
		scoutService entity.ScoutService
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&scoutService, `
		SELECT 
			service.*,
			IFNULL(staff.staff_name, '') AS staff_name,
			IFNULL(robot.name, '') AS robot_name
		FROM 
			scout_services AS service
		INNER JOIN
			agent_robots AS robot
		ON
			service.agent_robot_id = robot.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			service.agent_staff_id = staff.id
		WHERE
			service.id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &scoutService, nil
}

/****************************************************************************************/
/// 複数取得
//
// agentIDを使ってスカウトサービスの一覧を取得
func (repo *ScoutServiceRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.ScoutService, error) {
	var (
		scoutServiceList []*entity.ScoutService
	)

	// passwordは取得しない
	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&scoutServiceList, `
		SELECT 
			service.*,
			IFNULL(staff.staff_name, '') AS staff_name,
			IFNULL(robot.name, '') AS robot_name
		FROM 
			scout_services AS service
		INNER JOIN
			agent_robots AS robot
		ON
			service.agent_robot_id = robot.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			service.agent_staff_id = staff.id
		WHERE
			service.agent_robot_id IN (
				SELECT id
				FROM agent_robots
				WHERE agent_id = ?
			)
		ORDER BY 
			service.id DESC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return scoutServiceList, nil
}

// エージェントロボットIDからスカウトサービスの一覧を取得
func (repo *ScoutServiceRepositoryImpl) GetByAgentRobotID(agentRobotID uint) ([]*entity.ScoutService, error) {
	var (
		scoutServiceList []*entity.ScoutService
	)

	// RPAで使用するためpassword含む
	err := repo.executer.Select(
		repo.Name+".GetByAgentRobotID",
		&scoutServiceList, `
			SELECT *
			FROM scout_services
			WHERE
				agent_robot_id = ?
			ORDER BY 
				id DESC
		`,
		agentRobotID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return scoutServiceList, nil
}
