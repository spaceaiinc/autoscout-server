package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerDepartmentHistoryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerDepartmentHistoryRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerDepartmentHistoryRepository {
	return &JobSeekerDepartmentHistoryRepositoryImpl{
		Name:     "JobSeekerDepartmentHistoryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerDepartmentHistoryRepositoryImpl) Create(departmentHistory *entity.JobSeekerDepartmentHistory) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_department_histories (
				work_history_id,
				department,
				management_number,
				management_detail,
				job_description,
				start_year,
				end_year,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?
			)
		`,
		departmentHistory.WorkHistoryID,
		departmentHistory.Department,
		departmentHistory.ManagementNumber,
		departmentHistory.ManagementDetail,
		departmentHistory.JobDescription,
		departmentHistory.StartYear,
		departmentHistory.EndYear,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	departmentHistory.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 複数取得
//
func (repo *JobSeekerDepartmentHistoryRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDepartmentHistory, error) {
	var (
		departmentHistoryList []*entity.JobSeekerDepartmentHistory
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&departmentHistoryList, `
		SELECT *
		FROM job_seeker_department_histories
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
		return departmentHistoryList, err
	}

	return departmentHistoryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerDepartmentHistoryRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerDepartmentHistory, error) {
	var departmentHistoryList []*entity.JobSeekerDepartmentHistory

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&departmentHistoryList, `
		SELECT 
			jsdh.*
		FROM 
			job_seeker_department_histories AS jsdh
		INNER JOIN
			job_seeker_work_histories AS jswh
		ON 
			jsdh.work_history_id = jswh.id
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

	return departmentHistoryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerDepartmentHistoryRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerDepartmentHistory, error) {
	var departmentHistoryList []*entity.JobSeekerDepartmentHistory

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&departmentHistoryList, `
		SELECT 
			jsdh.*
		FROM 
			job_seeker_department_histories AS jsdh
		INNER JOIN
			job_seeker_work_histories AS jswh
		ON 
			jsdh.work_history_id = jswh.id
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

	return departmentHistoryList, nil
}

// 求職者リストから経験職種を取得
func (repo *JobSeekerDepartmentHistoryRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerDepartmentHistory, error) {
	var departmentHistoryList []*entity.JobSeekerDepartmentHistory

	if len(idList) == 0 {
		return departmentHistoryList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_department_histories
		WHERE
			work_history_id IN (
				SELECT id
				FROM job_seeker_work_histories
				WHERE job_seeker_id IN (%s)
			)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&departmentHistoryList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return departmentHistoryList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerDepartmentHistoryRepositoryImpl) All() ([]*entity.JobSeekerDepartmentHistory, error) {
	var departmentHistoryList []*entity.JobSeekerDepartmentHistory

	err := repo.executer.Select(
		repo.Name+".All",
		&departmentHistoryList, `
							SELECT *
							FROM job_seeker_department_histories
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return departmentHistoryList, nil
}
