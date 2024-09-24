package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerExperienceJobRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerExperienceJobRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerExperienceJobRepository {
	return &JobSeekerExperienceJobRepositoryImpl{
		Name:     "JobSeekerExperienceJobRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerExperienceJobRepositoryImpl) Create(desiredCompanyScale *entity.JobSeekerExperienceJob) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_experience_jobs (
				job_seeker_id,
				occupation,
				experience_year,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		desiredCompanyScale.JobSeekerID,
		desiredCompanyScale.Occupation,
		desiredCompanyScale.ExperienceYear,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredCompanyScale.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerExperienceJobRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_experience_jobs
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
func (repo *JobSeekerExperienceJobRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerExperienceJob, error) {
	var desiredCompanyScaleList []*entity.JobSeekerExperienceJob

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&desiredCompanyScaleList, `
		SELECT *
		FROM job_seeker_experience_jobs
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return desiredCompanyScaleList, err
	}

	return desiredCompanyScaleList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerExperienceJobRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerExperienceJob, error) {
	var desiredCompanyScaleList []*entity.JobSeekerExperienceJob

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&desiredCompanyScaleList, `
			SELECT 
				jsej.*
			FROM 
				job_seeker_experience_jobs AS jsej
			INNER JOIN
				job_seekers AS js
			ON
				jsej.job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredCompanyScaleList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerExperienceJobRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerExperienceJob, error) {
	var desiredCompanyScaleList []*entity.JobSeekerExperienceJob

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&desiredCompanyScaleList, `
			SELECT 
				jsej.*
			FROM 
				job_seeker_experience_jobs AS jsej
			INNER JOIN
				job_seekers AS js
			ON
				jsej.job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredCompanyScaleList, nil
}

// 求職者リストから希望休日タイプを取得
func (repo *JobSeekerExperienceJobRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerExperienceJob, error) {
	var desiredCompanyScaleList []*entity.JobSeekerExperienceJob

	if len(idList) == 0 {
		return desiredCompanyScaleList, nil
	}

	commaJoinIDs := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]")

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_experience_jobs
		WHERE
			job_seeker_id IN (%s)
	`, commaJoinIDs)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&desiredCompanyScaleList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredCompanyScaleList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerExperienceJobRepositoryImpl) All() ([]*entity.JobSeekerExperienceJob, error) {
	var desiredCompanyScaleList []*entity.JobSeekerExperienceJob

	err := repo.executer.Select(
		repo.Name+".All",
		&desiredCompanyScaleList, `
			SELECT * FROM job_seeker_experience_jobs
		`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredCompanyScaleList, nil
}
