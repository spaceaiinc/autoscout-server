package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerDesiredWorkLocationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerDesiredWorkLocationRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerDesiredWorkLocationRepository {
	return &JobSeekerDesiredWorkLocationRepositoryImpl{
		Name:     "JobSeekerDesiredWorkLocationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerDesiredWorkLocationRepositoryImpl) Create(desiredWorkLocation *entity.JobSeekerDesiredWorkLocation) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_desired_work_locations (
				job_seeker_id,
				desired_work_location,
				desired_rank,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		desiredWorkLocation.JobSeekerID,
		desiredWorkLocation.DesiredWorkLocation,
		desiredWorkLocation.DesiredRank,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredWorkLocation.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerDesiredWorkLocationRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_desired_work_locations
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
func (repo *JobSeekerDesiredWorkLocationRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDesiredWorkLocation, error) {
	var desiredWorkLocationList []*entity.JobSeekerDesiredWorkLocation

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&desiredWorkLocationList, `
		SELECT *
		FROM job_seeker_desired_work_locations
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return desiredWorkLocationList, err
	}

	return desiredWorkLocationList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerDesiredWorkLocationRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerDesiredWorkLocation, error) {
	var desiredWorkLocationList []*entity.JobSeekerDesiredWorkLocation

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&desiredWorkLocationList, `
			SELECT 
				jsdwl.*
			FROM 
				job_seeker_desired_work_locations AS jsdwl
			INNER JOIN
				job_seekers AS js
			ON
				jsdwl.job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredWorkLocationList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerDesiredWorkLocationRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerDesiredWorkLocation, error) {
	var desiredWorkLocationList []*entity.JobSeekerDesiredWorkLocation

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&desiredWorkLocationList, `
			SELECT 
				jsdwl.*
			FROM 
				job_seeker_desired_work_locations AS jsdwl
			INNER JOIN
				job_seekers AS js
			ON
				jsdwl.job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredWorkLocationList, nil
}

// 求職者リストから希望勤務地を取得
func (repo *JobSeekerDesiredWorkLocationRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerDesiredWorkLocation, error) {
	var desiredWorkLocationList []*entity.JobSeekerDesiredWorkLocation

	if len(idList) == 0 {
		return desiredWorkLocationList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_desired_work_locations
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&desiredWorkLocationList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredWorkLocationList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerDesiredWorkLocationRepositoryImpl) All() ([]*entity.JobSeekerDesiredWorkLocation, error) {
	var desiredWorkLocationList []*entity.JobSeekerDesiredWorkLocation

	err := repo.executer.Select(
		repo.Name+".All",
		&desiredWorkLocationList, `
							SELECT *
							FROM job_seeker_desired_work_locations
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredWorkLocationList, nil
}
