package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerDepartmentHistoryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerDepartmentHistoryRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerDepartmentHistoryRepository {
	return &SendingJobSeekerDepartmentHistoryRepositoryImpl{
		Name:     "SendingJobSeekerDepartmentHistoryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerDepartmentHistoryRepositoryImpl) Create(departmentHistory *entity.SendingJobSeekerDepartmentHistory) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_department_histories (
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

func (repo *SendingJobSeekerDepartmentHistoryRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerDepartmentHistory, error) {
	var (
		departmentHistoryList []*entity.SendingJobSeekerDepartmentHistory
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&departmentHistoryList, `
		SELECT *
		FROM sending_job_seeker_department_histories
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
		return departmentHistoryList, err
	}

	return departmentHistoryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerDepartmentHistoryRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerDepartmentHistory, error) {
	var (
		departmentHistoryList []*entity.SendingJobSeekerDepartmentHistory
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&departmentHistoryList, `
		SELECT 
			jsdh.*
		FROM 
			sending_job_seeker_department_histories AS jsdh
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

	return departmentHistoryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerDepartmentHistoryRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerDepartmentHistory, error) {
	var (
		departmentHistoryList []*entity.SendingJobSeekerDepartmentHistory
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&departmentHistoryList, `
		SELECT 
			jsdh.*
		FROM 
			sending_job_seeker_department_histories AS jsdh
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

	return departmentHistoryList, nil
}

// 求職者リストから経験職種を取得
func (repo *SendingJobSeekerDepartmentHistoryRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerDepartmentHistory, error) {
	var (
		departmentHistoryList []*entity.SendingJobSeekerDepartmentHistory
	)

	if len(idList) == 0 {
		return departmentHistoryList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_department_histories
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
		&departmentHistoryList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return departmentHistoryList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerDepartmentHistoryRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerDepartmentHistory, error) {
	var (
		departmentHistoryList []*entity.SendingJobSeekerDepartmentHistory
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&departmentHistoryList, `
							SELECT *
							FROM sending_job_seeker_department_histories
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return departmentHistoryList, nil
}
