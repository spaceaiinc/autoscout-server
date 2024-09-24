package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerDevelopmentSkillRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerDevelopmentSkillRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerDevelopmentSkillRepository {
	return &SendingJobSeekerDevelopmentSkillRepositoryImpl{
		Name:     "SendingJobSeekerDevelopmentSkillRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerDevelopmentSkillRepositoryImpl) Create(developmentSkill *entity.SendingJobSeekerDevelopmentSkill) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_development_skills (
				sending_job_seeker_id,
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
		developmentSkill.SendingJobSeekerID,
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

func (repo *SendingJobSeekerDevelopmentSkillRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_development_skills
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerDevelopmentSkillRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerDevelopmentSkill, error) {
	var (
		developmentSkillList []*entity.SendingJobSeekerDevelopmentSkill
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&developmentSkillList, `
		SELECT *
		FROM sending_job_seeker_development_skills
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return developmentSkillList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerDevelopmentSkillRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerDevelopmentSkill, error) {
	var (
		developmentSkillList []*entity.SendingJobSeekerDevelopmentSkill
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&developmentSkillList, `
			SELECT 
				jsds.*
			FROM 
				sending_job_seeker_development_skills AS jsds
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsds.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerDevelopmentSkillRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerDevelopmentSkill, error) {
	var (
		developmentSkillList []*entity.SendingJobSeekerDevelopmentSkill
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&developmentSkillList, `
			SELECT 
				jsds.*
			FROM 
				sending_job_seeker_development_skills AS jsds
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsds.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerDevelopmentSkillRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerDevelopmentSkill, error) {
	var (
		developmentSkillList []*entity.SendingJobSeekerDevelopmentSkill
	)

	if len(idList) == 0 {
		return developmentSkillList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_development_skills
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&developmentSkillList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return developmentSkillList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerDevelopmentSkillRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerDevelopmentSkill, error) {
	var (
		developmentSkillList []*entity.SendingJobSeekerDevelopmentSkill
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&developmentSkillList, `
							SELECT *
							FROM sending_job_seeker_development_skills
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return developmentSkillList, nil
}
