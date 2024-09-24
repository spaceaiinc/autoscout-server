package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerDesiredOccupationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerDesiredOccupationRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerDesiredOccupationRepository {
	return &JobSeekerDesiredOccupationRepositoryImpl{
		Name:     "JobSeekerDesiredOccupationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerDesiredOccupationRepositoryImpl) Create(desiredOccupation *entity.JobSeekerDesiredOccupation) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_desired_occupations (
				job_seeker_id,
				desired_occupation,
				desired_rank,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		desiredOccupation.JobSeekerID,
		desiredOccupation.DesiredOccupation,
		desiredOccupation.DesiredRank,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredOccupation.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerDesiredOccupationRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_desired_occupations
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
func (repo *JobSeekerDesiredOccupationRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDesiredOccupation, error) {
	var desiredOccupationList []*entity.JobSeekerDesiredOccupation

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&desiredOccupationList, `
		SELECT *
		FROM job_seeker_desired_occupations
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return desiredOccupationList, err
	}

	return desiredOccupationList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerDesiredOccupationRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerDesiredOccupation, error) {
	var desiredOccupationList []*entity.JobSeekerDesiredOccupation

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&desiredOccupationList, `
			SELECT 
				jsdo.*
			FROM 
				job_seeker_desired_occupations AS jsdo
			INNER JOIN
				job_seekers AS js
			ON
				jsdo.job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredOccupationList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerDesiredOccupationRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerDesiredOccupation, error) {
	var desiredOccupationList []*entity.JobSeekerDesiredOccupation

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&desiredOccupationList, `
			SELECT 
				jsdo.*
			FROM 
				job_seeker_desired_occupations AS jsdo
			INNER JOIN
				job_seekers AS js
			ON
				jsdo.job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredOccupationList, nil
}

// 求職者リストから希望職種を取得
func (repo *JobSeekerDesiredOccupationRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerDesiredOccupation, error) {
	var desiredOccupationList []*entity.JobSeekerDesiredOccupation

	if len(idList) == 0 {
		return desiredOccupationList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_desired_occupations
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&desiredOccupationList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredOccupationList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerDesiredOccupationRepositoryImpl) All() ([]*entity.JobSeekerDesiredOccupation, error) {
	var desiredOccupationList []*entity.JobSeekerDesiredOccupation

	err := repo.executer.Select(
		repo.Name+".All",
		&desiredOccupationList, `
							SELECT *
							FROM job_seeker_desired_occupations
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredOccupationList, nil
}
