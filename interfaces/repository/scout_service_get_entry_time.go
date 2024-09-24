package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ScoutServiceGetEntryTimeRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewScoutServiceGetEntryTimeRepositoryImpl(ex interfaces.SQLExecuter) usecase.ScoutServiceGetEntryTimeRepository {
	return &ScoutServiceGetEntryTimeRepositoryImpl{
		Name:     "ScoutServiceGetEntryTimeRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// エントリー取得時間の作成
func (repo *ScoutServiceGetEntryTimeRepositoryImpl) Create(scoutServiceGetEntryTime *entity.ScoutServiceGetEntryTime) error {
	nowTime := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO scout_service_get_entry_times (
			scout_service_id,
			start_hour,
			start_minute,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)`,
		scoutServiceGetEntryTime.ScoutServiceID,
		scoutServiceGetEntryTime.StartHour,
		scoutServiceGetEntryTime.StartMinute,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	scoutServiceGetEntryTime.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 削除
//
// スカウトサービスIDからエントリー取得時間の削除
func (repo *ScoutServiceGetEntryTimeRepositoryImpl) DeleteByScoutServiceID(scoutServiceID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE 
		FROM scout_service_get_entry_times
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
// IDを使ってエントリー取得時間の取得
func (repo *ScoutServiceGetEntryTimeRepositoryImpl) FindByID(id uint) (*entity.ScoutServiceGetEntryTime, error) {
	var (
		scoutServiceGetEntryTime entity.ScoutServiceGetEntryTime
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&scoutServiceGetEntryTime, `
		SELECT 
			templates.*,
			services.login_id,
			services.password,
			services.service_type
		FROM scout_service_get_entry_times AS templates
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

	return &scoutServiceGetEntryTime, nil
}

/****************************************************************************************/
/// 複数取得
//
// agentIDを使ってエントリー取得時間の一覧を取得
func (repo *ScoutServiceGetEntryTimeRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.ScoutServiceGetEntryTime, error) {
	var (
		scoutServiceGetEntryTimeList []*entity.ScoutServiceGetEntryTime
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&scoutServiceGetEntryTimeList, `
			SELECT *
			FROM scout_service_get_entry_times
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
				start_hour DESC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return scoutServiceGetEntryTimeList, nil
}

// scoutServiceIDを使ってエントリー取得時間の一覧を取得
func (repo *ScoutServiceGetEntryTimeRepositoryImpl) GetByScoutServiceID(scoutServiceID uint) ([]*entity.ScoutServiceGetEntryTime, error) {
	var (
		scoutServiceGetEntryTimeList []*entity.ScoutServiceGetEntryTime
	)

	err := repo.executer.Select(
		repo.Name+".GetByScoutServiceID",
		&scoutServiceGetEntryTimeList, `
			SELECT *
			FROM scout_service_get_entry_times
			WHERE
				scout_service_id = ?
			ORDER BY 
				start_hour DESC
		`,
		scoutServiceID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return scoutServiceGetEntryTimeList, nil
}

// agentRobotIDを使ってエントリー取得時間の一覧を取得
func (repo *ScoutServiceGetEntryTimeRepositoryImpl) GetByAgentRobotID(agentRobotID uint) ([]*entity.ScoutServiceGetEntryTime, error) {
	var (
		scoutServiceGetEntryTimeList []*entity.ScoutServiceGetEntryTime
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentRobotID",
		&scoutServiceGetEntryTimeList, `
				SELECT *
				FROM scout_service_get_entry_times
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

	return scoutServiceGetEntryTimeList, nil
}
