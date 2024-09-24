package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type EnterpriseActivityRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewEnterpriseActivityRepositoryImpl(ex interfaces.SQLExecuter) usecase.EnterpriseActivityRepository {
	return &EnterpriseActivityRepositoryImpl{
		Name:     "EnterpriseActivityRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *EnterpriseActivityRepositoryImpl) Create(activity *entity.EnterpriseActivity) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO enterprise_activities (
				enterprise_id,
				content,
				added_at,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		activity.EnterpriseID,
		activity.Content,
		activity.AddedAt,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	activity.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 複数取得 API
//
func (repo *EnterpriseActivityRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.EnterpriseActivity, error) {
	var activityList []*entity.EnterpriseActivity

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&activityList, `
		SELECT *
		FROM enterprise_activities
		WHERE
			enterprise_id = ?
		ORDER BY added_at DESC, id DESC
		`,
		enterpriseID,
	)
	if err != nil {
		fmt.Println(err)
		return activityList, err
	}

	return activityList, nil
}

// 担当者IDで企業情報を取得
func (repo *EnterpriseActivityRepositoryImpl) GetByAgentStaffID(agentStaffID uint) ([]*entity.EnterpriseActivity, error) {
	var activityList []*entity.EnterpriseActivity

	err := repo.executer.Select(
		repo.Name+".GetByAgentStaffID",
		&activityList, `
		SELECT *
		FROM 
		enterprise_activities
		WHERE
		enterprise_id IN (
			SELECT id
			FROM enterprise_profiles
			WHERE agent_staff_id = ?
		)
		ORDER BY added_at DESC, id DESC
		`,
		agentStaffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return activityList, nil
}

// エージェントIDから企業一覧を取得
func (repo *EnterpriseActivityRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.EnterpriseActivity, error) {
	var activityList []*entity.EnterpriseActivity

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&activityList, `
		SELECT *
		FROM 
		enterprise_activities
		WHERE
		enterprise_id IN (
			SELECT id
			FROM enterprise_profiles
			WHERE 
				agent_staff_id IN (
					SELECT id
					FROM agent_staffs
					WHERE
					agent_id = ?
					)
			)
		ORDER BY added_at DESC, id DESC
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return activityList, nil
}

// 企業IDリストからリストを取得
func (repo *EnterpriseActivityRepositoryImpl) GetByEnterpriseIDList(enterpriseIDList []uint) ([]*entity.EnterpriseActivity, error) {
	var activityList []*entity.EnterpriseActivity

	if len(enterpriseIDList) == 0 {
		return activityList, nil
	}

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseIDList",
		&activityList, `
		SELECT *
		FROM enterprise_activities
		WHERE
			enterprise_id IN (?)
		ORDER BY added_at DESC, id DESC
		`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(enterpriseIDList)), ", "), "[]"),
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return activityList, nil
}

// すべての企業情報を取得
func (repo *EnterpriseActivityRepositoryImpl) All() ([]*entity.EnterpriseActivity, error) {
	var activityList []*entity.EnterpriseActivity

	err := repo.executer.Select(
		repo.Name+".All",
		&activityList,
		`
			SELECT *
			FROM enterprise_activities
			ORDER BY added_at DESC, id DESC
		`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return activityList, nil
}
