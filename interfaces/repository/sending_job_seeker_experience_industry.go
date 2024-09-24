package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerExperienceIndustryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerExperienceIndustryRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerExperienceIndustryRepository {
	return &SendingJobSeekerExperienceIndustryRepositoryImpl{
		Name:     "SendingJobSeekerExperienceIndustryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerExperienceIndustryRepositoryImpl) Create(experienceIndustry *entity.SendingJobSeekerExperienceIndustry) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_experience_industries (
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

func (repo *SendingJobSeekerExperienceIndustryRepositoryImpl) Update(industry *entity.SendingJobSeekerExperienceIndustry) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE sending_job_seeker_experience_industries
		SET
			work_history_id = ?,
			industry = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		industry.WorkHistoryID,
		industry.Industry,
		time.Now().In(time.UTC),
		industry.ID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerExperienceIndustryRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerExperienceIndustry, error) {
	var (
		experienceIndustryList []*entity.SendingJobSeekerExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&experienceIndustryList, `
		SELECT *
		FROM sending_job_seeker_experience_industries
		WHERE
			work_history_id IN (
				SELECT id
				FROM sending_job_seeker_work_histories
				WHERE sending_job_seeker_id = ?
			)
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return experienceIndustryList, err
	}

	return experienceIndustryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerExperienceIndustryRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerExperienceIndustry, error) {
	var (
		experienceIndustryList []*entity.SendingJobSeekerExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&experienceIndustryList, `
			SELECT 
				jsei.*
			FROM 
				sending_job_seeker_experience_industries AS jsei
			INNER JOIN
				sending_job_seeker_work_histories AS jswh
			ON
				jsei.work_history_id = jswh.id
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

	return experienceIndustryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerExperienceIndustryRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerExperienceIndustry, error) {
	var (
		experienceIndustryList []*entity.SendingJobSeekerExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&experienceIndustryList, `
			SELECT 
				jsei.*
			FROM 
				sending_job_seeker_experience_industries AS jsei
			INNER JOIN
				sending_job_seeker_work_histories AS jswh
			ON
				jsei.work_history_id = jswh.id
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

	return experienceIndustryList, nil
}

// 求職者リストから経験業界を取得
func (repo *SendingJobSeekerExperienceIndustryRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerExperienceIndustry, error) {
	var (
		experienceIndustryList []*entity.SendingJobSeekerExperienceIndustry
	)

	if len(idList) == 0 {
		return experienceIndustryList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_experience_industries
		WHERE
			work_history_id IN (
				SELECT id
				FROM sending_job_seeker_work_histories
				WHERE sending_job_seeker_id IN (%s)
			)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&experienceIndustryList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return experienceIndustryList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerExperienceIndustryRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerExperienceIndustry, error) {
	var (
		experienceIndustryList []*entity.SendingJobSeekerExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&experienceIndustryList, `
							SELECT *
							FROM sending_job_seeker_experience_industries
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return experienceIndustryList, nil
}
