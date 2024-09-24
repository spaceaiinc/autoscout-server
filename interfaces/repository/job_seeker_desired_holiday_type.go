package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerDesiredHolidayTypeRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerDesiredHolidayTypeRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerDesiredHolidayTypeRepository {
	return &JobSeekerDesiredHolidayTypeRepositoryImpl{
		Name:     "JobSeekerDesiredHolidayTypeRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerDesiredHolidayTypeRepositoryImpl) Create(desiredHolidayType *entity.JobSeekerDesiredHolidayType) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_desired_holiday_types (
				job_seeker_id,
				holiday_type,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		desiredHolidayType.JobSeekerID,
		desiredHolidayType.HolidayType,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredHolidayType.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerDesiredHolidayTypeRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_desired_holiday_types
		WHERE job_seeker_id = ?
		`, jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 複数取得
//
func (repo *JobSeekerDesiredHolidayTypeRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDesiredHolidayType, error) {
	var desiredHolidayTypeList []*entity.JobSeekerDesiredHolidayType

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&desiredHolidayTypeList, `
		SELECT *
		FROM job_seeker_desired_holiday_types
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return desiredHolidayTypeList, err
	}

	return desiredHolidayTypeList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerDesiredHolidayTypeRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerDesiredHolidayType, error) {
	var desiredHolidayTypeList []*entity.JobSeekerDesiredHolidayType

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&desiredHolidayTypeList, `
			SELECT 
				jsdht.*
			FROM 
				job_seeker_desired_holiday_types AS jsdht
			INNER JOIN
				job_seekers AS js
			ON
				jsdht.job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredHolidayTypeList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerDesiredHolidayTypeRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerDesiredHolidayType, error) {
	var desiredHolidayTypeList []*entity.JobSeekerDesiredHolidayType

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&desiredHolidayTypeList, `
			SELECT 
				jsdht.*
			FROM 
				job_seeker_desired_holiday_types AS jsdht
			INNER JOIN
				job_seekers AS js
			ON
				jsdht.job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredHolidayTypeList, nil
}

// 求職者リストから希望休日タイプを取得
func (repo *JobSeekerDesiredHolidayTypeRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerDesiredHolidayType, error) {
	var desiredHolidayTypeList []*entity.JobSeekerDesiredHolidayType

	if len(idList) == 0 {
		return desiredHolidayTypeList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_desired_holiday_types
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&desiredHolidayTypeList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredHolidayTypeList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerDesiredHolidayTypeRepositoryImpl) All() ([]*entity.JobSeekerDesiredHolidayType, error) {
	var desiredHolidayTypeList []*entity.JobSeekerDesiredHolidayType

	err := repo.executer.Select(
		repo.Name+".All",
		&desiredHolidayTypeList, `
							SELECT *
							FROM job_seeker_desired_holiday_types
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredHolidayTypeList, nil
}
