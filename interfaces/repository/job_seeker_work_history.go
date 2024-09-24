package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerWorkHistoryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerWorkHistoryRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerWorkHistoryRepository {
	return &JobSeekerWorkHistoryRepositoryImpl{
		Name:     "JobSeekerWorkHistoryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerWorkHistoryRepositoryImpl) Create(workHistory *entity.JobSeekerWorkHistory) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_work_histories (
				job_seeker_id,
				company_name,
				employee_number_single,
				employee_number_group,
				public_offering,
				joining_year,
				employment_status,
				retire_reason_of_truth,
				retire_reason_of_public,
				retire_year,
				first_status,
				last_status,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)
		`,
		workHistory.JobSeekerID,
		workHistory.CompanyName,
		workHistory.EmployeeNumberSingle,
		workHistory.EmployeeNumberGroup,
		workHistory.PublicOffering,
		workHistory.JoiningYear,
		workHistory.EmploymentStatus,
		workHistory.RetireReasonOfTruth,
		workHistory.RetireReasonOfPublic,
		workHistory.RetireYear,
		workHistory.FirstStatus,
		workHistory.LastStatus,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	workHistory.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerWorkHistoryRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_work_histories
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
func (repo *JobSeekerWorkHistoryRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerWorkHistory, error) {
	var workHistoryList []*entity.JobSeekerWorkHistory

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&workHistoryList, `
		SELECT *
		FROM job_seeker_work_histories
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return workHistoryList, err
	}

	return workHistoryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerWorkHistoryRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerWorkHistory, error) {
	var workHistoryList []*entity.JobSeekerWorkHistory

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&workHistoryList, `
			SELECT 
				jswh.*
			FROM 
				job_seeker_work_histories AS jswh
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

	return workHistoryList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerWorkHistoryRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerWorkHistory, error) {
	var workHistoryList []*entity.JobSeekerWorkHistory

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&workHistoryList, `
			SELECT 
				jswh.*
			FROM 
				job_seeker_work_histories AS jswh
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

	return workHistoryList, nil
}

// 求職者リストから職歴情報を取得
func (repo *JobSeekerWorkHistoryRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerWorkHistory, error) {
	var workHistoryList []*entity.JobSeekerWorkHistory

	if len(idList) == 0 {
		return workHistoryList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_work_histories
		WHERE job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&workHistoryList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return workHistoryList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerWorkHistoryRepositoryImpl) All() ([]*entity.JobSeekerWorkHistory, error) {
	var workHistoryList []*entity.JobSeekerWorkHistory

	err := repo.executer.Select(
		repo.Name+".All",
		&workHistoryList, `
							SELECT *
							FROM job_seeker_work_histories
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return workHistoryList, nil
}
