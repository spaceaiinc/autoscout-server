package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerDesiredCompanyScaleRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerDesiredCompanyScaleRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerDesiredCompanyScaleRepository {
	return &JobSeekerDesiredCompanyScaleRepositoryImpl{
		Name:     "JobSeekerDesiredCompanyScaleRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerDesiredCompanyScaleRepositoryImpl) Create(desiredCompanyScale *entity.JobSeekerDesiredCompanyScale) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_desired_company_scales (
				job_seeker_id,
				desired_company_scale,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		desiredCompanyScale.JobSeekerID,
		desiredCompanyScale.DesiredCompanyScale,
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
func (repo *JobSeekerDesiredCompanyScaleRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_desired_company_scales
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
func (repo *JobSeekerDesiredCompanyScaleRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDesiredCompanyScale, error) {
	var (
		desiredCompanyScaleList []*entity.JobSeekerDesiredCompanyScale
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&desiredCompanyScaleList, `
		SELECT *
		FROM job_seeker_desired_company_scales
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
func (repo *JobSeekerDesiredCompanyScaleRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerDesiredCompanyScale, error) {
	var (
		desiredCompanyScaleList []*entity.JobSeekerDesiredCompanyScale
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&desiredCompanyScaleList, `
			SELECT 
				jsdcs.*
			FROM 
				job_seeker_desired_company_scales AS jsdcs
			INNER JOIN
				job_seekers AS js
			ON
				jsdcs.job_seeker_id = js.id
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
func (repo *JobSeekerDesiredCompanyScaleRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerDesiredCompanyScale, error) {
	var (
		desiredCompanyScaleList []*entity.JobSeekerDesiredCompanyScale
	)

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&desiredCompanyScaleList, `
			SELECT 
				jsdcs.*
			FROM 
				job_seeker_desired_company_scales AS jsdcs
			INNER JOIN
				job_seekers AS js
			ON
				jsdcs.job_seeker_id = js.id
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
func (repo *JobSeekerDesiredCompanyScaleRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerDesiredCompanyScale, error) {
	var (
		desiredCompanyScaleList []*entity.JobSeekerDesiredCompanyScale
	)

	if len(idList) == 0 {
		return desiredCompanyScaleList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_desired_company_scales
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

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
func (repo *JobSeekerDesiredCompanyScaleRepositoryImpl) All() ([]*entity.JobSeekerDesiredCompanyScale, error) {
	var (
		desiredCompanyScaleList []*entity.JobSeekerDesiredCompanyScale
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&desiredCompanyScaleList, `
							SELECT *
							FROM job_seeker_desired_company_scales
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredCompanyScaleList, nil
}
