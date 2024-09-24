package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerExperienceIndustryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerExperienceIndustryRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerExperienceIndustryRepository {
	return &JobSeekerExperienceIndustryRepositoryImpl{
		Name:     "JobSeekerExperienceIndustryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerExperienceIndustryRepositoryImpl) Create(experienceIndustry *entity.JobSeekerExperienceIndustry) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_experience_industries (
				work_history_id,
				industry,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		experienceIndustry.WorkHistoryID,
		experienceIndustry.Industry,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	experienceIndustry.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 複数取得
//
func (repo *JobSeekerExperienceIndustryRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerExperienceIndustry, error) {
	var experienceIndustryList []*entity.JobSeekerExperienceIndustry

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&experienceIndustryList, `
		SELECT *
		FROM job_seeker_experience_industries
		WHERE
			work_history_id IN (
				SELECT id
				FROM job_seeker_work_histories
				WHERE job_seeker_id = ?
			)
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return experienceIndustryList, err
	}

	return experienceIndustryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerExperienceIndustryRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerExperienceIndustry, error) {
	var experienceIndustryList []*entity.JobSeekerExperienceIndustry

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&experienceIndustryList, `
			SELECT 
				jsei.*
			FROM 
				job_seeker_experience_industries AS jsei
			INNER JOIN
				job_seeker_work_histories AS jswh
			ON
				jsei.work_history_id = jswh.id
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

	return experienceIndustryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerExperienceIndustryRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerExperienceIndustry, error) {
	var experienceIndustryList []*entity.JobSeekerExperienceIndustry

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&experienceIndustryList, `
			SELECT 
				jsei.*
			FROM 
				job_seeker_experience_industries AS jsei
			INNER JOIN
				job_seeker_work_histories AS jswh
			ON
				jsei.work_history_id = jswh.id
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

	return experienceIndustryList, nil
}

// 求職者リストから経験業界を取得
func (repo *JobSeekerExperienceIndustryRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerExperienceIndustry, error) {
	var experienceIndustryList []*entity.JobSeekerExperienceIndustry

	if len(idList) == 0 {
		return experienceIndustryList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_experience_industries
		WHERE
			work_history_id IN (
				SELECT id
				FROM job_seeker_work_histories
				WHERE job_seeker_id IN (%s)
			)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&experienceIndustryList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return experienceIndustryList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerExperienceIndustryRepositoryImpl) All() ([]*entity.JobSeekerExperienceIndustry, error) {
	var experienceIndustryList []*entity.JobSeekerExperienceIndustry

	err := repo.executer.Select(
		repo.Name+".All",
		&experienceIndustryList, `
							SELECT *
							FROM job_seeker_experience_industries
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return experienceIndustryList, nil
}
