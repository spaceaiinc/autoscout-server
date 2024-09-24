package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type InterviewTaskRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewInterviewTaskRepositoryImpl(ex interfaces.SQLExecuter) usecase.InterviewTaskRepository {
	return &InterviewTaskRepositoryImpl{
		Name:     "InterviewTaskRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// 面談調整のタスクを作成
func (repo *InterviewTaskRepositoryImpl) Create(interviewTask *entity.InterviewTask) error {
	nowTime := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO interview_tasks (
			interview_task_group_id,
			agent_staff_id,
			ca_staff_id,
			phase_category,
			phase_sub_category,
			remarks,
			deadline_day,
			deadline_time,
			select_action_label,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)
		`,
		interviewTask.InterviewTaskGroupID,
		interviewTask.AgentStaffID,
		interviewTask.CAStaffID,
		interviewTask.PhaseCategory,
		interviewTask.PhaseSubCategory,
		interviewTask.Remarks,
		interviewTask.DeadlineDay,
		interviewTask.DeadlineTime,
		interviewTask.SelectActionLabel,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	interviewTask.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 削除
//
// 面談調整タスクを削除
func (repo *InterviewTaskRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`DELETE FROM interview_tasks WHERE id = ?`,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/****************************************************************************************/
/// 単数取得
//
// IDから面談調整のタスクグループを取得
func (repo *InterviewTaskRepositoryImpl) FindByID(id uint) (*entity.InterviewTask, error) {
	var (
		interviewTask entity.InterviewTask
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&interviewTask, `
		SELECT 
			task.*, ca_staff.id AS ca_staff_id, IFNULL(ca_staff.staff_name, '') AS ca_staff_name,
			task_group.agent_id, task_group.job_seeker_id,
			task_group.last_request_at, task_group.last_watched_at, 
			task_group.interview_date, task_group.first_interview_date,
			agent.uuid AS agent_uuid, agent.agent_name, agent.interview_adjustment_email,
			seeker.uuid AS job_seeker_uuid, seeker.last_name, seeker.first_name, 
			seeker.last_furigana, seeker.first_furigana,
			seeker.phase, seeker.email, chat_group.line_active, IFNULL(staff.staff_name, '') AS staff_name
		FROM 
			interview_tasks AS task 
		INNER JOIN
			interview_task_groups AS task_group
		ON
			task_group.id = task.interview_task_group_id
		INNER JOIN
			agents AS agent
		ON
			agent.id = task_group.agent_id
		INNER JOIN
			job_seekers AS seeker
		ON
			seeker.id = task_group.job_seeker_id
		INNER JOIN
			chat_group_with_job_seekers AS chat_group
		ON
			chat_group.job_seeker_id = seeker.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			task.agent_staff_id = staff.id
		LEFT OUTER JOIN
			agent_staffs AS ca_staff
		ON
			seeker.agent_staff_id = ca_staff.id
		WHERE
			task.id = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &interviewTask, nil
}

// エージェントIDと求職者IDから最新の面談調整を取得
func (repo *InterviewTaskRepositoryImpl) FindLatestByAgentIDAndJobSeekerID(agentID, jobSeekerID uint) (*entity.InterviewTask, error) {
	var (
		interviewTask entity.InterviewTask
	)

	err := repo.executer.Get(
		repo.Name+".FindLatestByAgentIDAndJobSeekerID",
		&interviewTask, `
		SELECT 
			task.*
		FROM 
			interview_tasks AS task 
		INNER JOIN
			interview_task_groups AS task_group
		ON
			task_group.id = task.interview_task_group_id
		WHERE
			task_group.agent_id = ?
		AND
			task_group.job_seeker_id = ?
		ORDER BY 
			task.id DESC
		LIMIT 1
		`,
		agentID, jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &interviewTask, nil
}

/****************************************************************************************/
/// 複数取得
//
// グループIDから面談調整のタスクグループを取得
func (repo *InterviewTaskRepositoryImpl) GetByGroupID(groupID uint) ([]*entity.InterviewTask, error) {
	var (
		interviewTaskList []*entity.InterviewTask
	)

	err := repo.executer.Select(
		repo.Name+".GetByGroupID",
		&interviewTaskList, `
		SELECT 
			task.*, 
			task_group.last_request_at, task_group.last_watched_at,
			seeker.id AS job_seeker_id, seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana, seeker.phase,
			staff.staff_name
		FROM 
			interview_tasks AS task 
		INNER JOIN
			interview_task_groups AS task_group
		ON
			task_group.id = task.interview_task_group_id
		INNER JOIN
			job_seekers AS seeker
		ON
			seeker.id = task_group.job_seeker_id
		INNER JOIN
			agent_staffs AS staff
		ON
			task.agent_staff_id = staff.id
		WHERE
			task.interview_task_group_id = ? 
		ORDER BY
			task.id DESC
		`,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return interviewTaskList, nil
}

// エージェントIDから最新の面談調整のタスクを取得（phaseが「0 or 1」）
func (repo *InterviewTaskRepositoryImpl) GetLatestAdjustmentByAgentID(agentID uint) ([]*entity.InterviewTask, error) {
	var (
		interviewTaskList []*entity.InterviewTask
	)

	// 求職者のagent_staff_idをca_staff_idとして取得する

	err := repo.executer.Select(
		repo.Name+".GetLatestAdjustmentByAgentID",
		&interviewTaskList, `
		SELECT 
			task.*, ca_staff.id AS ca_staff_id, IFNULL(ca_staff.staff_name, '') AS ca_staff_name,
			task_group.agent_id, task_group.job_seeker_id,
			task_group.last_request_at, task_group.last_watched_at, 
			task_group.interview_date, task_group.first_interview_date,
			agent.uuid AS agent_uuid, agent.agent_name, agent.interview_adjustment_email,
			seeker.uuid AS job_seeker_uuid, seeker.last_name, seeker.first_name, 
			seeker.last_furigana, seeker.first_furigana,
			seeker.phase, seeker.email, chat_group.line_active, IFNULL(staff.staff_name, '') AS staff_name,
			seeker.phone_number,
			seeker.created_at AS job_seeker_created_at
		FROM 
			interview_tasks AS task 
		INNER JOIN
			interview_task_groups AS task_group
		ON
			task_group.id = task.interview_task_group_id
		INNER JOIN
			agents AS agent
		ON
			agent.id = task_group.agent_id
		INNER JOIN
			job_seekers AS seeker
		ON
			seeker.id = task_group.job_seeker_id
		INNER JOIN
			chat_group_with_job_seekers AS chat_group
		ON
			chat_group.job_seeker_id = seeker.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			task.agent_staff_id = staff.id
		LEFT OUTER JOIN
			agent_staffs AS ca_staff
		ON
			seeker.agent_staff_id = ca_staff.id
		WHERE
			task_group.agent_id = ? 
		AND 
			task.id = (
				SELECT task_join.id
				FROM interview_tasks AS task_join 
				WHERE 
					task.interview_task_group_id = task_join.interview_task_group_id
				AND
					task.phase_category IN (0, 1)
				ORDER BY task_join.created_at DESC
				LIMIT 1
			)
		ORDER BY task.phase_category ASC, task.deadline_day ASC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return interviewTaskList, nil
}

// エージェントIDから最新の面談調整のタスクを取得（phaseが「2 or 3」）
func (repo *InterviewTaskRepositoryImpl) GetLatestConfirmationByAgentID(agentID uint) ([]*entity.InterviewTask, error) {
	var (
		interviewTaskList []*entity.InterviewTask
	)

	err := repo.executer.Select(
		repo.Name+".GetLatestConfirmationByAgentID",
		&interviewTaskList, `
		SELECT 
			task.*, ca_staff.id AS ca_staff_id, IFNULL(ca_staff.staff_name, '') AS ca_staff_name,
			task_group.agent_id, task_group.job_seeker_id,
			task_group.last_request_at, task_group.last_watched_at, 
			task_group.interview_date, task_group.first_interview_date,
			agent.uuid AS agent_uuid, agent.agent_name, agent.interview_adjustment_email,
			seeker.uuid AS job_seeker_uuid, seeker.last_name, seeker.first_name, 
			seeker.last_furigana, seeker.first_furigana,
			seeker.phase, seeker.email, chat_group.line_active, IFNULL(staff.staff_name, '') AS staff_name,
			seeker.phone_number,
			seeker.created_at AS job_seeker_created_at
		FROM 
			interview_tasks AS task 
		INNER JOIN
			interview_task_groups AS task_group
		ON
			task_group.id = task.interview_task_group_id
		INNER JOIN
			agents AS agent
		ON
			agent.id = task_group.agent_id
		INNER JOIN
			job_seekers AS seeker
		ON
			seeker.id = task_group.job_seeker_id
		INNER JOIN
			chat_group_with_job_seekers AS chat_group
		ON
			chat_group.job_seeker_id = seeker.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			task.agent_staff_id = staff.id
		LEFT OUTER JOIN
			agent_staffs AS ca_staff
		ON
			seeker.agent_staff_id = ca_staff.id
		WHERE
			task_group.agent_id = ? 
		AND 
			task.id = (
				SELECT task_join.id
				FROM interview_tasks AS task_join 
				WHERE 
					task.interview_task_group_id = task_join.interview_task_group_id
				AND
					task.phase_category IN (2, 3)
				ORDER BY 
					task_join.created_at DESC
				LIMIT 1
			)
		ORDER BY task.phase_category ASC, task.deadline_day ASC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return interviewTaskList, nil
}
