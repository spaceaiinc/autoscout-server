package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerDevelopmentSkillRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerDevelopmentSkillRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerDevelopmentSkillRepository {
	return &JobSeekerDevelopmentSkillRepositoryImpl{
		Name:     "JobSeekerDevelopmentSkillRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerDevelopmentSkillRepositoryImpl) Create(developmentSkill *entity.JobSeekerDevelopmentSkill) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_development_skills (
				job_seeker_id,
				development_category,
				development_type,
				experience_year,
				experience_month,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?
			)
		`,
		developmentSkill.JobSeekerID,
		developmentSkill.DevelopmentCategory,
		developmentSkill.DevelopmentType,
		developmentSkill.ExperienceYear,
		developmentSkill.ExperienceMonth,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	developmentSkill.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerDevelopmentSkillRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_development_skills
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
func (repo *JobSeekerDevelopmentSkillRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDevelopmentSkill, error) {
	var developmentSkillList []*entity.JobSeekerDevelopmentSkill

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&developmentSkillList, `
		SELECT *
		FROM job_seeker_development_skills
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return developmentSkillList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerDevelopmentSkillRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerDevelopmentSkill, error) {
	var developmentSkillList []*entity.JobSeekerDevelopmentSkill

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&developmentSkillList, `
			SELECT 
				jsds.*
			FROM 
				job_seeker_development_skills AS jsds
			INNER JOIN
				job_seekers AS js
			ON
				jsds.job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return developmentSkillList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerDevelopmentSkillRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerDevelopmentSkill, error) {
	var developmentSkillList []*entity.JobSeekerDevelopmentSkill

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&developmentSkillList, `
			SELECT 
				jsds.*
			FROM 
				job_seeker_development_skills AS jsds
			INNER JOIN
				job_seekers AS js
			ON
				jsds.job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return developmentSkillList, nil
}

// 求職者リストから希望開発スキルを取得
func (repo *JobSeekerDevelopmentSkillRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerDevelopmentSkill, error) {
	var developmentSkillList []*entity.JobSeekerDevelopmentSkill

	if len(idList) == 0 {
		return developmentSkillList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_development_skills
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&developmentSkillList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return developmentSkillList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerDevelopmentSkillRepositoryImpl) All() ([]*entity.JobSeekerDevelopmentSkill, error) {
	var developmentSkillList []*entity.JobSeekerDevelopmentSkill

	err := repo.executer.Select(
		repo.Name+".All",
		&developmentSkillList, `
							SELECT *
							FROM job_seeker_development_skills
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return developmentSkillList, nil
}
