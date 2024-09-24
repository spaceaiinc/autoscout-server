package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerLanguageSkillRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerLanguageSkillRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerLanguageSkillRepository {
	return &SendingJobSeekerLanguageSkillRepositoryImpl{
		Name:     "SendingJobSeekerLanguageSkillRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerLanguageSkillRepositoryImpl) Create(languageSkill *entity.SendingJobSeekerLanguageSkill) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_language_skills (
				sending_job_seeker_id,
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
		languageSkill.SendingJobSeekerID,
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

func (repo *SendingJobSeekerLanguageSkillRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_language_skills
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerLanguageSkillRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerLanguageSkill, error) {
	var (
		languageSkillList []*entity.SendingJobSeekerLanguageSkill
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&languageSkillList, `
		SELECT *
		FROM sending_job_seeker_language_skills
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return languageSkillList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerLanguageSkillRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerLanguageSkill, error) {
	var (
		languageSkillList []*entity.SendingJobSeekerLanguageSkill
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&languageSkillList, `
			SELECT 
				jsls.*
			FROM 
				sending_job_seeker_language_skills AS jsls
			INNER JOIN 
				sending_job_seekers AS js
			ON
				jsls.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerLanguageSkillRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerLanguageSkill, error) {
	var (
		languageSkillList []*entity.SendingJobSeekerLanguageSkill
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&languageSkillList, `
			SELECT 
				jsls.*
			FROM 
				sending_job_seeker_language_skills AS jsls
			INNER JOIN 
				sending_job_seekers AS js
			ON
				jsls.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerLanguageSkillRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerLanguageSkill, error) {
	var (
		languageSkillList []*entity.SendingJobSeekerLanguageSkill
	)

	if len(idList) == 0 {
		return languageSkillList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_language_skills
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&languageSkillList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return languageSkillList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerLanguageSkillRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerLanguageSkill, error) {
	var (
		languageSkillList []*entity.SendingJobSeekerLanguageSkill
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&languageSkillList, `
							SELECT *
							FROM sending_job_seeker_language_skills
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return languageSkillList, nil
}
