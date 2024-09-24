package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerExperienceOccupationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerExperienceOccupationRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerExperienceOccupationRepository {
	return &SendingJobSeekerExperienceOccupationRepositoryImpl{
		Name:     "SendingJobSeekerExperienceOccupationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerExperienceOccupationRepositoryImpl) Create(experienceOccupation *entity.SendingJobSeekerExperienceOccupation) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_experience_occupations (
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

func (repo *SendingJobSeekerExperienceOccupationRepositoryImpl) DeleteByID(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_experience_occupations
		WHERE id = ?
		`, id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerExperienceOccupationRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerExperienceOccupation, error) {
	var (
		experienceOccupationList []*entity.SendingJobSeekerExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&experienceOccupationList, `
		SELECT *
		FROM sending_job_seeker_experience_occupations
		WHERE department_history_id IN (
			SELECT id
			FROM sending_job_seeker_department_histories
			WHERE work_history_id IN (
				SELECT id
				FROM sending_job_seeker_work_histories
				WHERE sending_job_seeker_id = ?
			)
		)
		`,
		sendingJobSeekerID,
	)

	fmt.Println("experienceOccupationList", experienceOccupationList)

	if err != nil {
		fmt.Println(err)
		return experienceOccupationList, err
	}

	return experienceOccupationList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerExperienceOccupationRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerExperienceOccupation, error) {
	var (
		experienceOccupationList []*entity.SendingJobSeekerExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&experienceOccupationList, `
			SELECT 
				jseo.*
			FROM 
				sending_job_seeker_experience_occupations AS jseo
			INNER JOIN
				sending_job_seeker_department_histories AS jsdh
			ON
				jseo.department_history_id = jsdh.id
			INNER JOIN
				sending_job_seeker_work_histories AS jswh
			ON
				jsdh.work_history_id = jswh.id
			INNER JOIN
				sending_job_seekers AS js
			ON
				jswh.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerExperienceOccupationRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerExperienceOccupation, error) {
	var (
		experienceOccupationList []*entity.SendingJobSeekerExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&experienceOccupationList, `
			SELECT 
				jseo.*
			FROM 
				sending_job_seeker_experience_occupations AS jseo
			INNER JOIN
				sending_job_seeker_department_histories AS jsdh
			ON
				jseo.department_history_id = jsdh.id
			INNER JOIN
				sending_job_seeker_work_histories AS jswh
			ON
				jsdh.work_history_id = jswh.id
			INNER JOIN
				sending_job_seekers AS js
			ON
				jswh.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerExperienceOccupationRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerExperienceOccupation, error) {
	var (
		experienceOccupationList []*entity.SendingJobSeekerExperienceOccupation
	)

	if len(idList) == 0 {
		return experienceOccupationList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_experience_occupations
		WHERE department_history_id IN (
			SELECT id
			FROM sending_job_seeker_department_histories
			WHERE work_history_id IN (
				SELECT id
				FROM sending_job_seeker_work_histories
				WHERE sending_job_seeker_id IN (%s)
			)
		)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&experienceOccupationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return experienceOccupationList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerExperienceOccupationRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerExperienceOccupation, error) {
	var (
		experienceOccupationList []*entity.SendingJobSeekerExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&experienceOccupationList, `
							SELECT *
							FROM sending_job_seeker_experience_occupations
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return experienceOccupationList, nil
}
