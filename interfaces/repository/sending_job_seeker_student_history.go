package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerStudentHistoryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerStudentHistoryRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerStudentHistoryRepository {
	return &SendingJobSeekerStudentHistoryRepositoryImpl{
		Name:     "SendingJobSeekerStudentHistoryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerStudentHistoryRepositoryImpl) Create(studentHistory *entity.SendingJobSeekerStudentHistory) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_student_histories (
				sending_job_seeker_id,
				school_category,
				school_name,
				school_level,
				subject,
				entrance_year,
				first_status,
				graduation_year,
				last_status,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?
			)
		`,
		studentHistory.SendingJobSeekerID,
		studentHistory.SchoolCategory,
		studentHistory.SchoolName,
		studentHistory.SchoolLevel,
		studentHistory.Subject,
		studentHistory.EntranceYear,
		studentHistory.FirstStatus,
		studentHistory.GraduationYear,
		studentHistory.LastStatus,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	studentHistory.ID = uint(lastID)

	return nil
}

func (repo *SendingJobSeekerStudentHistoryRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_student_histories
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerStudentHistoryRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerStudentHistory, error) {
	var (
		studentHistoryList []*entity.SendingJobSeekerStudentHistory
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&studentHistoryList, `
		SELECT *
		FROM sending_job_seeker_student_histories
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return studentHistoryList, err
	}

	return studentHistoryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerStudentHistoryRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerStudentHistory, error) {
	var (
		studentHistoryList []*entity.SendingJobSeekerStudentHistory
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&studentHistoryList, `
			SELECT 
				jssh.*
			FROM 
				sending_job_seeker_student_histories AS jssh
			INNER JOIN
				sending_job_seekers AS js
			ON
				jssh.sending_job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return studentHistoryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerStudentHistoryRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerStudentHistory, error) {
	var (
		studentHistoryList []*entity.SendingJobSeekerStudentHistory
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&studentHistoryList, `
			SELECT 
				jssh.*
			FROM 
				sending_job_seeker_student_histories AS jssh
			INNER JOIN
				sending_job_seekers AS js
			ON
				jssh.sending_job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return studentHistoryList, nil
}

// 求職者リストから学歴情報を取得
func (repo *SendingJobSeekerStudentHistoryRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerStudentHistory, error) {
	var (
		studentHistoryList []*entity.SendingJobSeekerStudentHistory
	)

	if len(idList) == 0 {
		return studentHistoryList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_student_histories
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&studentHistoryList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return studentHistoryList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerStudentHistoryRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerStudentHistory, error) {
	var (
		studentHistoryList []*entity.SendingJobSeekerStudentHistory
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&studentHistoryList, `
							SELECT *
							FROM sending_job_seeker_student_histories
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return studentHistoryList, nil
}
