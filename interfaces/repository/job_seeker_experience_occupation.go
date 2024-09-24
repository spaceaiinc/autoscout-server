package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerExperienceOccupationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerExperienceOccupationRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerExperienceOccupationRepository {
	return &JobSeekerExperienceOccupationRepositoryImpl{
		Name:     "JobSeekerExperienceOccupationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerExperienceOccupationRepositoryImpl) Create(experienceOccupation *entity.JobSeekerExperienceOccupation) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_experience_occupations (
				department_history_id,
				occupation,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		experienceOccupation.DepartmentHistoryID,
		experienceOccupation.Occupation,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	experienceOccupation.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 複数取得
//
func (repo *JobSeekerExperienceOccupationRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerExperienceOccupation, error) {
	var experienceOccupationList []*entity.JobSeekerExperienceOccupation

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&experienceOccupationList, `
		SELECT *
		FROM job_seeker_experience_occupations
		WHERE department_history_id IN (
			SELECT id
			FROM job_seeker_department_histories
			WHERE work_history_id IN (
				SELECT id
				FROM job_seeker_work_histories
				WHERE job_seeker_id = ?
			)
		)
		`,
		jobSeekerID,
	)

	fmt.Println("experienceOccupationList", experienceOccupationList)

	if err != nil {
		fmt.Println(err)
		return experienceOccupationList, err
	}

	return experienceOccupationList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerExperienceOccupationRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerExperienceOccupation, error) {
	var experienceOccupationList []*entity.JobSeekerExperienceOccupation

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&experienceOccupationList, `
			SELECT 
				jseo.*
			FROM 
				job_seeker_experience_occupations AS jseo
			INNER JOIN
				job_seeker_department_histories AS jsdh
			ON
				jseo.department_history_id = jsdh.id
			INNER JOIN
				job_seeker_work_histories AS jswh
			ON
				jsdh.work_history_id = jswh.id
			INNER JOIN
				job_seekers AS js
			ON
				jswh.job_seeker_id = js.id
			WHERE 
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return experienceOccupationList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerExperienceOccupationRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerExperienceOccupation, error) {
	var experienceOccupationList []*entity.JobSeekerExperienceOccupation

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&experienceOccupationList, `
			SELECT 
				jseo.*
			FROM 
				job_seeker_experience_occupations AS jseo
			INNER JOIN
				job_seeker_department_histories AS jsdh
			ON
				jseo.department_history_id = jsdh.id
			INNER JOIN
				job_seeker_work_histories AS jswh
			ON
				jsdh.work_history_id = jswh.id
			INNER JOIN
				job_seekers AS js
			ON
				jswh.job_seeker_id = js.id
			WHERE 
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return experienceOccupationList, nil
}

// 求職者リストから経験職種を取得
func (repo *JobSeekerExperienceOccupationRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerExperienceOccupation, error) {
	var experienceOccupationList []*entity.JobSeekerExperienceOccupation

	if len(idList) == 0 {
		return experienceOccupationList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_experience_occupations
		WHERE department_history_id IN (
			SELECT id
			FROM job_seeker_department_histories
			WHERE work_history_id IN (
				SELECT id
				FROM job_seeker_work_histories
				WHERE job_seeker_id IN (%s)
			)
		)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&experienceOccupationList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return experienceOccupationList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerExperienceOccupationRepositoryImpl) All() ([]*entity.JobSeekerExperienceOccupation, error) {
	var experienceOccupationList []*entity.JobSeekerExperienceOccupation

	err := repo.executer.Select(
		repo.Name+".All",
		&experienceOccupationList, `
							SELECT *
							FROM job_seeker_experience_occupations
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return experienceOccupationList, nil
}
