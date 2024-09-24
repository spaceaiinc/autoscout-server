package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerWorkHistoryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerWorkHistoryRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerWorkHistoryRepository {
	return &SendingJobSeekerWorkHistoryRepositoryImpl{
		Name:     "SendingJobSeekerWorkHistoryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerWorkHistoryRepositoryImpl) Create(workHistory *entity.SendingJobSeekerWorkHistory) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_work_histories (
				sending_job_seeker_id,
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
		workHistory.SendingJobSeekerID,
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

func (repo *SendingJobSeekerWorkHistoryRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_work_histories
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerWorkHistoryRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerWorkHistory, error) {
	var (
		workHistoryList []*entity.SendingJobSeekerWorkHistory
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&workHistoryList, `
		SELECT *
		FROM sending_job_seeker_work_histories
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return workHistoryList, err
	}

	return workHistoryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerWorkHistoryRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerWorkHistory, error) {
	var (
		workHistoryList []*entity.SendingJobSeekerWorkHistory
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&workHistoryList, `
			SELECT 
				jswh.*
			FROM 
				sending_job_seeker_work_histories AS jswh
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

	return workHistoryList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *SendingJobSeekerWorkHistoryRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerWorkHistory, error) {
	var (
		workHistoryList []*entity.SendingJobSeekerWorkHistory
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&workHistoryList, `
			SELECT 
				jswh.*
			FROM 
				sending_job_seeker_work_histories AS jswh
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

	return workHistoryList, nil
}

// 求職者リストから職歴情報を取得
func (repo *SendingJobSeekerWorkHistoryRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerWorkHistory, error) {
	var (
		workHistoryList []*entity.SendingJobSeekerWorkHistory
	)

	if len(idList) == 0 {
		return workHistoryList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_work_histories
		WHERE sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&workHistoryList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return workHistoryList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerWorkHistoryRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerWorkHistory, error) {
	var (
		workHistoryList []*entity.SendingJobSeekerWorkHistory
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&workHistoryList, `
							SELECT *
							FROM sending_job_seeker_work_histories
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return workHistoryList, nil
}
