package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerLanguageSkillRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerLanguageSkillRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerLanguageSkillRepository {
	return &JobSeekerLanguageSkillRepositoryImpl{
		Name:     "JobSeekerLanguageSkillRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerLanguageSkillRepositoryImpl) Create(languageSkill *entity.JobSeekerLanguageSkill) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_language_skills (
				job_seeker_id,
				language_type,
				language_level,
				toeic,
				toeic_examination_year,
				toefl_ibt,
				toefl_ibt_examination_year,
				toefl_pbt,
				toefl_pbt_examination_year,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?
			)
		`,
		languageSkill.JobSeekerID,
		languageSkill.LanguageType,
		languageSkill.LanguageLevel,
		languageSkill.Toeic,
		languageSkill.ToeicExaminationYear,
		languageSkill.ToeflIBT,
		languageSkill.ToeflIBTExaminationYear,
		languageSkill.ToeflPBT,
		languageSkill.ToeflPBTExaminationYear,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	languageSkill.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerLanguageSkillRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM job_seeker_language_skills
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
func (repo *JobSeekerLanguageSkillRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerLanguageSkill, error) {
	var languageSkillList []*entity.JobSeekerLanguageSkill

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&languageSkillList, `
		SELECT *
		FROM job_seeker_language_skills
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return languageSkillList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerLanguageSkillRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerLanguageSkill, error) {
	var languageSkillList []*entity.JobSeekerLanguageSkill

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&languageSkillList, `
			SELECT 
				jsls.*
			FROM 
				job_seeker_language_skills AS jsls
			INNER JOIN 
				job_seekers AS js
			ON
				jsls.job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return languageSkillList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerLanguageSkillRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerLanguageSkill, error) {
	var languageSkillList []*entity.JobSeekerLanguageSkill

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&languageSkillList, `
			SELECT 
				jsls.*
			FROM 
				job_seeker_language_skills AS jsls
			INNER JOIN 
				job_seekers AS js
			ON
				jsls.job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return languageSkillList, nil
}

// 求職者リストから所持言語スキルを取得
func (repo *JobSeekerLanguageSkillRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerLanguageSkill, error) {
	var languageSkillList []*entity.JobSeekerLanguageSkill

	if len(idList) == 0 {
		return languageSkillList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_language_skills
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&languageSkillList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return languageSkillList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerLanguageSkillRepositoryImpl) All() ([]*entity.JobSeekerLanguageSkill, error) {
	var languageSkillList []*entity.JobSeekerLanguageSkill

	err := repo.executer.Select(
		repo.Name+".All",
		&languageSkillList, `
							SELECT *
							FROM job_seeker_language_skills
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return languageSkillList, nil
}
