package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerStudentHistoryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerStudentHistoryRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerStudentHistoryRepository {
	return &JobSeekerStudentHistoryRepositoryImpl{
		Name:     "JobSeekerStudentHistoryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerStudentHistoryRepositoryImpl) Create(studentHistory *entity.JobSeekerStudentHistory) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_student_histories (
				job_seeker_id,
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
		studentHistory.JobSeekerID,
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

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerStudentHistoryRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_student_histories
		WHERE job_seeker_id = ?
		`, jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *JobSeekerStudentHistoryRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerStudentHistory, error) {
	var studentHistoryList []*entity.JobSeekerStudentHistory

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&studentHistoryList, `
		SELECT *
		FROM job_seeker_student_histories
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return studentHistoryList, err
	}

	return studentHistoryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerStudentHistoryRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerStudentHistory, error) {
	var studentHistoryList []*entity.JobSeekerStudentHistory

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&studentHistoryList, `
			SELECT 
				jssh.*
			FROM 
				job_seeker_student_histories AS jssh
			INNER JOIN
				job_seekers AS js
			ON
				jssh.job_seeker_id = js.id
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
func (repo *JobSeekerStudentHistoryRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerStudentHistory, error) {
	var studentHistoryList []*entity.JobSeekerStudentHistory

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&studentHistoryList, `
			SELECT 
				jssh.*
			FROM 
				job_seeker_student_histories AS jssh
			INNER JOIN
				job_seekers AS js
			ON
				jssh.job_seeker_id = js.id
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
func (repo *JobSeekerStudentHistoryRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerStudentHistory, error) {
	var studentHistoryList []*entity.JobSeekerStudentHistory

	if len(idList) == 0 {
		return studentHistoryList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_student_histories
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&studentHistoryList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return studentHistoryList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerStudentHistoryRepositoryImpl) All() ([]*entity.JobSeekerStudentHistory, error) {
	var studentHistoryList []*entity.JobSeekerStudentHistory

	err := repo.executer.Select(
		repo.Name+".All",
		&studentHistoryList, `
							SELECT *
							FROM job_seeker_student_histories
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return studentHistoryList, nil
}
