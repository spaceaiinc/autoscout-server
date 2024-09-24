package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerDesiredIndustryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerDesiredIndustryRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerDesiredIndustryRepository {
	return &JobSeekerDesiredIndustryRepositoryImpl{
		Name:     "JobSeekerDesiredIndustryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerDesiredIndustryRepositoryImpl) Create(desiredIndustry *entity.JobSeekerDesiredIndustry) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_desired_industries (
				job_seeker_id,
				desired_industry,
				desired_rank,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		desiredIndustry.JobSeekerID,
		desiredIndustry.DesiredIndustry,
		desiredIndustry.DesiredRank,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredIndustry.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerDesiredIndustryRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_desired_industries
		WHERE job_seeker_id = ?
		`, jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *JobSeekerDesiredIndustryRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDesiredIndustry, error) {
	var desiredIndustryList []*entity.JobSeekerDesiredIndustry

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&desiredIndustryList, `
		SELECT
			*
		FROM 
			job_seeker_desired_industries
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return desiredIndustryList, err
	}

	return desiredIndustryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerDesiredIndustryRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerDesiredIndustry, error) {
	var desiredIndustryList []*entity.JobSeekerDesiredIndustry

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&desiredIndustryList, `
			SELECT 
				jsdi.*
			FROM 
				job_seeker_desired_industries AS jsdi
			INNER JOIN
				job_seekers AS js
			ON
				jsdi.job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredIndustryList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerDesiredIndustryRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerDesiredIndustry, error) {
	var desiredIndustryList []*entity.JobSeekerDesiredIndustry

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&desiredIndustryList, `
			SELECT 
				jsdi.*
			FROM 
				job_seeker_desired_industries AS jsdi
			INNER JOIN
				job_seekers AS js
			ON
				jsdi.job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredIndustryList, nil
}

// 求職者リストから希望業界を取得
func (repo *JobSeekerDesiredIndustryRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerDesiredIndustry, error) {
	var desiredIndustryList []*entity.JobSeekerDesiredIndustry

	if len(idList) == 0 {
		return desiredIndustryList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_desired_industries
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&desiredIndustryList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredIndustryList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerDesiredIndustryRepositoryImpl) All() ([]*entity.JobSeekerDesiredIndustry, error) {
	var desiredIndustryList []*entity.JobSeekerDesiredIndustry

	err := repo.executer.Select(
		repo.Name+".All",
		&desiredIndustryList, `
							SELECT *
							FROM job_seeker_desired_industries
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredIndustryList, nil
}
