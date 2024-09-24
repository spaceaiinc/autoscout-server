package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type TaskRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewTaskRepositoryImpl(ex interfaces.SQLExecuter) usecase.TaskRepository {
	return &TaskRepositoryImpl{
		Name:     "TaskRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// タスクの作成
func (repo *TaskRepositoryImpl) Create(task *entity.Task) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO tasks (
			task_group_id,
			phase_category,
			phase_sub_category,
			staff_type,
			executed_staff_id,
			remarks,
			deadline_day,
			deadline_time,
			talk_about_in_interview,
			schedule_collection_condition,
			exam_guide_content,
			is_check_double_sided,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, 
				?, ?, ?, ?
			)`,
		task.TaskGroupID,
		task.PhaseCategory,
		task.PhaseSubCategory,
		task.StaffType,
		task.ExecutedStaffID,
		task.Remarks,
		task.DeadlineDay,
		task.DeadlineTime,
		task.TalkAboutInInterview,
		task.ScheduleCollectionCondition,
		task.ExamGuideContent,
		task.IsCheckDoubleSided,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	task.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
func (repo *TaskRepositoryImpl) Update(id uint, task *entity.Task) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE tasks
		SET
			phase_category = ?,
			phase_sub_category = ?,
			staff_type = ?,
			executed_staff_id = ?,
			remarks = ?,
			deadline_day = ?,
			deadline_time = ?,
			talk_about_in_interview = ?,
			schedule_collection_condition = ?,
			exam_guide_content = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		task.PhaseCategory,
		task.PhaseSubCategory,
		task.StaffType,
		task.ExecutedStaffID,
		task.Remarks,
		task.DeadlineDay,
		task.DeadlineTime,
		task.TalkAboutInInterview,
		task.ScheduleCollectionCondition,
		task.ExamGuideContent,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 削除
//
func (repo *TaskRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM tasks
		WHERE id = ?
		`, id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 単数取得
//
func (repo *TaskRepositoryImpl) FindByID(id uint) (*entity.Task, error) {
	var (
		task entity.Task
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&task, `
		SELECT 
			task.*, task_group.uuid AS task_group_uuid, task_group.is_double_sided
		FROM 
			tasks AS task
		INNER JOIN
			task_groups AS task_group
		ON
			task.task_group_id = task_group.id
		WHERE
			task.id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (repo *TaskRepositoryImpl) FindByJobSeekerUUIDAndJobInformationUUID(jobseekerUUID, jobInformationUUID uuid.UUID) (*entity.Task, error) {
	var (
		task entity.Task
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerUUIDAndJobInformationUUID",
		&task, `
		SELECT 
			tasks.*, seeker.agent_staff_id AS ca_staff_id, billing.agent_staff_id AS ra_staff_id
		FROM 
			tasks
		INNER JOIN 
			task_groups
		ON 
			tasks.task_group_id = task_groups.id
		INNER JOIN
			job_seekers AS seeker
		ON
			task_groups.job_seeker_id = seeker.id
		INNER JOIN
			job_informations AS job_info
		ON
			task_groups.job_information_id = job_info.id
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		WHERE
			seeker.uuid = ?
		AND
			job_info.uuid = ?
		ORDER BY 
			id DESC
		LIMIT 1
		`,
		jobseekerUUID, jobInformationUUID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &task, nil
}

func (repo *TaskRepositoryImpl) FindLatestWithRelatedByJobSeekerUUIDAndJobInformationUUID(jobSeekerUUID, jobInformationUUID uuid.UUID) (*entity.Task, error) {
	var (
		task entity.Task
	)

	err := repo.executer.Get(
		repo.Name+".FindLatestWithRelatedByJobSeekerUUIDAndJobInformationUUID",
		&task, `
			SELECT 
				task.*,
				task_group.uuid AS task_group_uuid,	task_group.is_double_sided,
				task_group.selection_flow_pattern_id, task_group.joining_date,
				task_group.ra_last_request_at, task_group.ra_last_watched_at, 
				task_group.ca_last_request_at, task_group.ca_last_watched_at, 
				seeker.id AS job_seeker_id, seeker.uuid AS job_seeker_uuid,
				seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
				billing.how_to_recommend, 
				job_info.id AS job_information_id, job_info.uuid AS job_information_uuid, job_info.billing_address_id,
				ca_staff.id AS ca_staff_id, ca_staff.staff_name AS ca_staff_name, ca_agent.agent_name AS ca_agent_name, ca_staff.agent_id AS ca_agent_id,
				ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_agent.agent_name AS ra_agent_name, ra_staff.agent_id AS ra_agent_id,
				flow_pattern.flow_pattern,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM
				tasks AS task
				INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				task_group.job_seeker_id = seeker.id
			INNER JOIN
				job_informations AS job_info
			ON
				task_group.job_information_id = job_info.id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS ra_staff
			ON
				billing.agent_staff_id = ra_staff.id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			INNER JOIN
				agents AS ra_agent
			ON
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				agents AS ca_agent
			ON
				ca_staff.agent_id = ca_agent.id
			LEFT OUTER JOIN
				job_information_selection_flow_patterns AS flow_pattern
			ON
				task_group.selection_flow_pattern_id = flow_pattern.id
			WHERE
				task.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)		
			AND
				job_info.uuid = ?	
			AND
				seeker.uuid = ?
			LIMIT 1
		`, jobInformationUUID, jobSeekerUUID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &task, nil
}

func (repo *TaskRepositoryImpl) FindLatestByGroupID(taskGroupID uint) (*entity.Task, error) {
	var (
		task entity.Task
	)

	err := repo.executer.Get(
		repo.Name+".FindLatestByGroupID",
		&task, `
		SELECT *
		FROM tasks
		WHERE
			task_group_id = ?
		ORDER BY 
			id DESC
		LIMIT 1
		`,
		taskGroupID)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// 直近の「書類選考/結果回収中」or「○次選考/選考結果の回収依頼」を取得
func (repo *TaskRepositoryImpl) FindLatestByGroupIDAndCollectResultPhase(taskGroupID uint) (*entity.Task, error) {
	var (
		task entity.Task
	)

	err := repo.executer.Get(
		repo.Name+".FindLatestByGroupIDAndCollectResultPhase",
		&task, `
		SELECT *
		FROM tasks
		WHERE
			task_group_id = ?
		AND (
			phase_category = ? AND
			phase_sub_category = ?
		) OR (
			phase_category IN(?, ?, ?, ?, ?, ?) AND
			phase_sub_category = ?
		)
		ORDER BY 
			id DESC
		LIMIT 1
		`,
		taskGroupID,
		entity.DocumentSelection,
		entity.CollectResultOfDocumentSelection,
		entity.FirstSelection,
		entity.SecondSelection,
		entity.ThirdSelection,
		entity.FourthSelection,
		entity.FifthSelection,
		entity.FinalSelection,
		entity.RequestCollectSelectionThought,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &task, nil
}

// 「phase_sub_category」が辞退処理する前（100未満）のタスクを取得
func (repo *TaskRepositoryImpl) FindLatestForContinue(taskGroupID uint, phase null.Int) (*entity.Task, error) {
	var (
		task entity.Task
	)

	err := repo.executer.Get(
		repo.Name+".FindLatestForContinue",
		&task, `
		SELECT *
		FROM tasks
		WHERE
			task_group_id = ?
		AND
			phase_category = ?
		AND
			phase_sub_category < 100
		ORDER BY 
			created_at DESC
		LIMIT 1
		`,
		taskGroupID, phase)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// 求職者IDと求人IDを使って最新タスクを取得
func (repo *TaskRepositoryImpl) FindLatestByJobSeekerIDAndJobInformationID(jobSeekerID, jobInformationID uint) (*entity.Task, error) {
	var (
		task entity.Task
	)

	/**
		job_information_selection_flow_patternsのレコードを外部結合で取得
	**/
	err := repo.executer.Get(
		repo.Name+".FindLatestByJobSeekerIDAndJobInformationID",
		&task, `
			SELECT 
				task.*,
				task_group.uuid AS task_group_uuid,	task_group.is_double_sided,
				task_group.selection_flow_pattern_id, task_group.joining_date,
				task_group.ra_last_request_at, task_group.ra_last_watched_at, 
				task_group.ca_last_request_at, task_group.ca_last_watched_at, 
				seeker.id AS job_seeker_id, seeker.uuid AS job_seeker_uuid,
				seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
				billing.how_to_recommend, 
				job_info.id AS job_information_id, job_info.uuid AS job_information_uuid, job_info.billing_address_id,
				ca_staff.id AS ca_staff_id, ca_staff.staff_name AS ca_staff_name, ca_agent.agent_name AS ca_agent_name, ca_staff.agent_id AS ca_agent_id,
				ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_agent.agent_name AS ra_agent_name, ra_staff.agent_id AS ra_agent_id,
				flow_pattern.flow_pattern,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				task_group.job_seeker_id = seeker.id
			INNER JOIN
				job_informations AS job_info
			ON
				task_group.job_information_id = job_info.id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS ra_staff
			ON
				billing.agent_staff_id = ra_staff.id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			INNER JOIN
				agents AS ra_agent
			ON
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				agents AS ca_agent
			ON
				ca_staff.agent_id = ca_agent.id
			LEFT OUTER JOIN
				job_information_selection_flow_patterns AS flow_pattern
			ON
				task_group.selection_flow_pattern_id = flow_pattern.id
			WHERE 
				task_group.job_seeker_id = ?
			AND
				task_group.job_information_id = ?
			AND
				task.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			LIMIT 1
		`,
		jobSeekerID, jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &task, nil
}

/****************************************************************************************/
/// 複数取得
//
// taskGroupIDを使ってタスクグループの一覧を取得
func (repo *TaskRepositoryImpl) GetByTaskGroupID(taskGroupID uint) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	err := repo.executer.Select(
		repo.Name+".GetByTaskGroupID",
		&taskList, `
			SELECT 
				task.*
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			WHERE
				task.task_group_id = ? 
			ORDER BY 
				task.id DESC
		`,
		taskGroupID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// taskGroupIDを使ってタスクグループの一覧を取得 *面談日時と評価ポイントも含む
func (repo *TaskRepositoryImpl) GetWithRelatedByTaskGroupID(taskGroupID uint) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	err := repo.executer.Select(
		repo.Name+".GetWithRelatedByTaskGroupID",
		&taskList, `
			SELECT 
				task.*,
				task_group.uuid AS task_group_uuid,	task_group.is_double_sided,
				task_group.selection_flow_pattern_id, task_group.joining_date,
				task_group.ra_last_request_at, task_group.ra_last_watched_at, 
				task_group.ca_last_request_at, task_group.ca_last_watched_at, 
				seeker.id AS job_seeker_id, seeker.uuid AS job_seeker_uuid,
				seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
				billing.how_to_recommend,
				job_info.id AS job_information_id, job_info.uuid AS job_information_uuid, job_info.billing_address_id, job_info.required_documents_detail,
				ca_staff.id AS ca_staff_id, ca_staff.staff_name AS ca_staff_name, ca_agent.agent_name AS ca_agent_name, ca_staff.agent_id AS ca_agent_id,
				ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_agent.agent_name AS ra_agent_name, ra_staff.agent_id AS ra_agent_id,
				flow_pattern.flow_pattern, selection_info.id AS selection_information_id, IFNULL(selection_info.is_questionnairy, FALSE) AS is_questionnairy,
				IFNULL(evaluation.good_point, '') AS good_point, IFNULL(evaluation.ng_point, '') AS ng_point,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				task_group.job_seeker_id = seeker.id
			INNER JOIN
				job_informations AS job_info
			ON
				task_group.job_information_id = job_info.id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS ra_staff
			ON
				billing.agent_staff_id = ra_staff.id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			INNER JOIN
				agents AS ra_agent
			ON
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				agents AS ca_agent
			ON
				ca_staff.agent_id = ca_agent.id
			LEFT OUTER JOIN
				job_information_selection_flow_patterns AS flow_pattern
			ON
				task_group.selection_flow_pattern_id = flow_pattern.id
			LEFT OUTER JOIN
				job_information_selection_informations AS selection_info
			ON
				flow_pattern.id = selection_info.selection_flow_id
			AND
				task.phase_category = selection_info.selection_type
			LEFT OUTER JOIN
				evaluation_points AS evaluation
			ON
				task.id = evaluation.task_id
			WHERE
				task.task_group_id = ? 
			ORDER BY 
				task.id DESC
		`,
		taskGroupID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// agentIDを使ってタスクグループの一覧を取得
func (repo *TaskRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&taskList, `
			SELECT 
				task.*, task_group.job_seeker_id, task_group.job_information_id
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = task_group.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			INNER JOIN
				agent_staffs AS ra_staff
			ON
				billing.agent_staff_id = ra_staff.id
			WHERE
				ca_staff.agent_id = ? 
			OR
				ra_staff.agent_id = ? 
			ORDER BY 
				task.id DESC
		`,
		agentID, agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// agentStaffIDを使って最新タスクの一覧を取得
func (repo *TaskRepositoryImpl) GetLatestByStaffID(staffID uint) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	// job_information_selection_flow_patternsのレコードを外部結合で取得
	// phase_sub_categoryが「90 ~ 99」は終了
	// phase_categoryが「9」でphase_sub_categoryが「1001」は決定終了

	err := repo.executer.Select(
		repo.Name+".GetLatestByStaffID",
		&taskList, `
			SELECT 
				task.*,
				task_group.uuid AS task_group_uuid,	task_group.is_double_sided,
				task_group.selection_flow_pattern_id, task_group.joining_date,
				task_group.ra_last_request_at, task_group.ra_last_watched_at, 
				task_group.ca_last_request_at, task_group.ca_last_watched_at, 
				task_group.is_self_application,
				seeker.id AS job_seeker_id, seeker.uuid AS job_seeker_uuid,
				seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
				billing.how_to_recommend, 
				job_info.id AS job_information_id, job_info.uuid AS job_information_uuid, job_info.billing_address_id, job_info.required_documents_detail,
				ca_staff.id AS ca_staff_id, ca_staff.staff_name AS ca_staff_name, ca_agent.agent_name AS ca_agent_name, ca_staff.agent_id AS ca_agent_id,
				ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_agent.agent_name AS ra_agent_name, ra_staff.agent_id AS ra_agent_id,
				flow_pattern.flow_pattern, selection_info.id AS selection_information_id, IFNULL(selection_info.is_questionnairy, FALSE) AS is_questionnairy,
				IFNULL(evaluation.good_point, '') AS good_point, IFNULL(evaluation.ng_point, '') AS ng_point,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != ''
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				task_group.job_seeker_id = seeker.id
			INNER JOIN
				job_informations AS job_info
			ON
				task_group.job_information_id = job_info.id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS ra_staff
			ON
				billing.agent_staff_id = ra_staff.id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			INNER JOIN
				agents AS ra_agent
			ON
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				agents AS ca_agent
			ON
				ca_staff.agent_id = ca_agent.id
			LEFT OUTER JOIN
				job_information_selection_flow_patterns AS flow_pattern
			ON
				task_group.selection_flow_pattern_id = flow_pattern.id
			LEFT OUTER JOIN
				job_information_selection_informations AS selection_info
			ON
				flow_pattern.id = selection_info.selection_flow_id
			AND
				task.phase_category = selection_info.selection_type
			LEFT OUTER JOIN
				evaluation_points AS evaluation
			ON
				task.id = evaluation.task_id
			WHERE (
				ca_staff.id = ? OR
				(!task_group.is_double_sided && ra_staff.id = ?)
			)
			AND
				task.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			AND (
				(
					task.phase_category IN(0, 1, 2, 3, 4, 5, 6, 7, 8) &&
					task.phase_sub_category NOT IN(90, 91, 92, 93, 94, 95, 96, 97, 98, 99) 
				) OR (
					task.phase_category = 9 &&
					task.phase_sub_category NOT IN(90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 1001) 
				)
			)
			ORDER BY 
				deadline_day ASC,
				deadline_time ASC
		`,
		staffID, staffID,
	)

	if err != nil {
		fmt.Println("GetLatestByStaffID error")
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// 「同一求職者」で「同一フェーズ」のタスクのリストを取得
func (repo *TaskRepositoryImpl) GetLatestSameByJobSeekerID(jobSeekerID, id uint, phase, phaseSub null.Int) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	err := repo.executer.Select(
		repo.Name+".GetLatestSameByJobSeekerID",
		&taskList, `
			SELECT 
				task.*,
				task_group.uuid AS task_group_uuid,	task_group.is_double_sided,
				task_group.selection_flow_pattern_id, task_group.joining_date,
				task_group.ra_last_request_at, task_group.ra_last_watched_at, 
				task_group.ca_last_request_at, task_group.ca_last_watched_at, 
				seeker.id AS job_seeker_id, seeker.uuid AS job_seeker_uuid,
				seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
				billing.how_to_recommend, 
				job_info.id AS job_information_id, job_info.uuid AS job_information_uuid, job_info.billing_address_id, job_info.required_documents_detail,
				ca_staff.id AS ca_staff_id, ca_staff.staff_name AS ca_staff_name, ca_agent.agent_name AS ca_agent_name, ca_staff.agent_id AS ca_agent_id,
				ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_agent.agent_name AS ra_agent_name, ra_staff.agent_id AS ra_agent_id,
				flow_pattern.flow_pattern, selection_info.id AS selection_information_id, IFNULL(selection_info.is_questionnairy, FALSE) AS is_questionnairy,
				IFNULL(evaluation.good_point, '') AS good_point, IFNULL(evaluation.ng_point, '') AS ng_point,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				task_group.job_seeker_id = seeker.id
			INNER JOIN
				job_informations AS job_info
			ON
				task_group.job_information_id = job_info.id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS ra_staff
			ON
				billing.agent_staff_id = ra_staff.id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			INNER JOIN
				agents AS ra_agent
			ON
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				agents AS ca_agent
			ON
				ca_staff.agent_id = ca_agent.id
			LEFT OUTER JOIN
				job_information_selection_flow_patterns AS flow_pattern
			ON
				task_group.selection_flow_pattern_id = flow_pattern.id
			LEFT OUTER JOIN
				job_information_selection_informations AS selection_info
			ON
				flow_pattern.id = selection_info.selection_flow_id
			AND
				task.phase_category = selection_info.selection_type
			LEFT OUTER JOIN
				evaluation_points AS evaluation
			ON
				task.id = evaluation.task_id
			WHERE 
				seeker.id = ?
			AND 
				task.id != ?
			AND
				task.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			AND
				task.phase_category = ?
			AND
				task.phase_sub_category = ?
			ORDER BY 
				deadline_day ASC,
				deadline_time ASC
		`,
		jobSeekerID, id, phase, phaseSub,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// agentStaffIDを使って最新タスクの一覧を取得
func (repo *TaskRepositoryImpl) GetLatestByGroupIDList(groupIDList []uint) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)
	if len(groupIDList) == 0 {
		return taskList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			task.*,
			IFNULL(evaluation.good_point, '') AS good_point, 
			IFNULL(evaluation.ng_point, '') AS ng_point
		FROM 
			tasks AS task
		LEFT OUTER JOIN
			evaluation_points AS evaluation
		ON
			task.id = evaluation.task_id
		INNER JOIN (
				SELECT 
						task_group_id, 
						MAX(id) AS latest_task_id
				FROM 
						tasks
				WHERE 
						task_group_id IN (%s)
				GROUP BY 
						task_group_id
		) AS latest_task
		ON 
				task.id = latest_task.latest_task_id
		ORDER BY 
			task.id DESC
		`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(groupIDList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetLatestByGroupIDList",
		&taskList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// 書類選考以降のタスク情報取得
func (repo *TaskRepositoryImpl) GetLatestAfterSelectPhaseByJobSeekerID(jobSeekerID uint, phase entity.TaskCategory) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	/**
		job_information_selection_flow_patternsのレコードを外部結合で取得
	**/

	err := repo.executer.Select(
		repo.Name+".GetLatestAfterSelectPhaseByJobSeekerID",
		&taskList, `
			SELECT 
				task.*,
				task_group.uuid AS task_group_uuid,	task_group.is_double_sided,
				task_group.selection_flow_pattern_id, task_group.joining_date,
				task_group.ra_last_request_at, task_group.ra_last_watched_at, 
				task_group.ca_last_request_at, task_group.ca_last_watched_at, 
				seeker.id AS job_seeker_id, seeker.uuid AS job_seeker_uuid,
				seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
				billing.how_to_recommend, 
				job_info.id AS job_information_id, job_info.uuid AS job_information_uuid, job_info.billing_address_id, job_info.required_documents_detail,
				ca_staff.id AS ca_staff_id, ca_staff.staff_name AS ca_staff_name, ca_agent.agent_name AS ca_agent_name, ca_staff.agent_id AS ca_agent_id,
				ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_agent.agent_name AS ra_agent_name, ra_staff.agent_id AS ra_agent_id,
				flow_pattern.flow_pattern, selection_info.id AS selection_information_id, IFNULL(selection_info.is_questionnairy, FALSE) AS is_questionnairy,
				IFNULL(evaluation.good_point, '') AS good_point, IFNULL(evaluation.ng_point, '') AS ng_point,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				task_group.job_seeker_id = seeker.id
			INNER JOIN
				job_informations AS job_info
			ON
				task_group.job_information_id = job_info.id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS ra_staff
			ON
				billing.agent_staff_id = ra_staff.id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			INNER JOIN
				agents AS ra_agent
			ON
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				agents AS ca_agent
			ON
				ca_staff.agent_id = ca_agent.id
			LEFT OUTER JOIN
				job_information_selection_flow_patterns AS flow_pattern
			ON
				task_group.selection_flow_pattern_id = flow_pattern.id
			LEFT OUTER JOIN
				job_information_selection_informations AS selection_info
			ON
				flow_pattern.id = selection_info.selection_flow_id
			AND
				task.phase_category = selection_info.selection_type
			LEFT OUTER JOIN
				evaluation_points AS evaluation
			ON
				task.id = evaluation.task_id
			WHERE 
				seeker.id = ?
			AND
				task.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE 
						task.task_group_id = task_join.task_group_id
					AND 
						task_join.phase_category >= ?
					ORDER BY 
						task_join.id DESC
					LIMIT 1
				)
			ORDER BY 
				deadline_day ASC,
				deadline_time ASC
		`,
		jobSeekerID, phase,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// agentIDを使ってタスクグループの一覧を取得
func (repo *TaskRepositoryImpl) GetLatestByAgentIDAndFreeWord(agentID uint, jobSeekerFreeWord string, enterpriseFreeWord string) ([]*entity.Task, error) {
	var (
		taskList                []*entity.Task
		jobSeekerFreeWordQuery  string
		enterpriseFreeWordQuery string
	)

	// 検索ワードがある場合
	if jobSeekerFreeWord != "" {
		jobSeekerFreeWord = "%" + jobSeekerFreeWord + "%"
		jobSeekerFreeWordQuery = fmt.Sprintf(`
			AND (	
				seeker.last_name LIKE '%s' 	
				OR  seeker.first_name LIKE '%s'	
				OR  seeker.last_furigana LIKE '%s'  	
				OR  seeker.first_furigana LIKE '%s' 
				)`, jobSeekerFreeWord, jobSeekerFreeWord, jobSeekerFreeWord, jobSeekerFreeWord)
	}

	// 検索ワードがある場合
	if enterpriseFreeWord != "" {
		enterpriseFreeWordQuery = fmt.Sprintf(`
			AND (
				job_info.is_external = TRUE AND (
					MATCH(task_group.external_company_name) AGAINST('%s' IN BOOLEAN MODE) OR 
					MATCH(task_group.external_job_information_title) AGAINST('%s' IN BOOLEAN MODE)
				) OR
				job_info.is_external = FALSE AND (
					MATCH(enterprise.company_name) AGAINST('%s' IN BOOLEAN MODE) OR 
					MATCH(job_info.title) AGAINST('%s' IN BOOLEAN MODE) 
				)
			)`, enterpriseFreeWord, enterpriseFreeWord, enterpriseFreeWord, enterpriseFreeWord)
	}

	// クエリをまとめる
	query := fmt.Sprintf(`
	SELECT 
		task.*,
		task_group.uuid AS task_group_uuid,	task_group.is_double_sided,
		task_group.selection_flow_pattern_id, task_group.joining_date,
		task_group.ra_last_request_at, task_group.ra_last_watched_at, 
		task_group.ca_last_request_at, task_group.ca_last_watched_at, 
		seeker.id AS job_seeker_id, seeker.uuid AS job_seeker_uuid,
		seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
		billing.how_to_recommend, 
		job_info.id AS job_information_id, job_info.uuid AS job_information_uuid, job_info.billing_address_id, job_info.required_documents_detail,
		ca_staff.id AS ca_staff_id, ca_staff.staff_name AS ca_staff_name, ca_agent.agent_name AS ca_agent_name, ca_staff.agent_id AS ca_agent_id,
		ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_agent.agent_name AS ra_agent_name, ra_staff.agent_id AS ra_agent_id,
		IFNULL(evaluation.good_point, '') AS good_point, IFNULL(evaluation.ng_point, '') AS ng_point,
		CASE 
			WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
			THEN task_group.external_job_information_title
			ELSE job_info.title
		END AS title,
		CASE 
			WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
			THEN task_group.external_company_name
			ELSE enterprise.company_name
		END AS company_name,
		task_group.external_job_listing_url
	FROM 
		tasks AS task
	INNER JOIN
		task_groups AS task_group
	ON
		task.task_group_id = task_group.id
	INNER JOIN
		job_seekers AS seeker
	ON
		task_group.job_seeker_id = seeker.id
	INNER JOIN
		job_informations AS job_info
	ON
		task_group.job_information_id = job_info.id
	INNER JOIN
		billing_addresses AS billing
	ON
		job_info.billing_address_id = billing.id
	INNER JOIN
		enterprise_profiles AS enterprise
	ON
		billing.enterprise_id = enterprise.id
	INNER JOIN
		agent_staffs AS ra_staff
	ON
		billing.agent_staff_id = ra_staff.id
	INNER JOIN
		agent_staffs AS ca_staff
	ON
		seeker.agent_staff_id = ca_staff.id
	INNER JOIN
		agents AS ra_agent
	ON
		ra_staff.agent_id = ra_agent.id
	INNER JOIN
		agents AS ca_agent
	ON
		ca_staff.agent_id = ca_agent.id
	LEFT OUTER JOIN
		evaluation_points AS evaluation
	ON
		task.id = evaluation.task_id
	WHERE (
		task_group.ra_agent_id = %v OR
		task_group.ca_agent_id = %v
	)
	AND
		task.id = (
			SELECT task_join.id
			FROM tasks AS task_join 
			WHERE task.task_group_id = task_join.task_group_id
			ORDER BY task_join.id DESC
			LIMIT 1
		)
	%s
	%s
	ORDER BY 
		task.id DESC
	`, agentID, agentID,
		jobSeekerFreeWordQuery, enterpriseFreeWordQuery)

	// クエリを実行
	err := repo.executer.Select(
		repo.Name+".GetLatestByAgentIDAndFreeWord",
		&taskList,
		query,
	)
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

/****************************************************************************************/
/****************************************************************************************/
// アライアンス解除時のタスク確認用
//
// 各エージェントIDがCA/RAに該当するタスクを取得　*アライアンス解除時のタスク確認用
func (repo *TaskRepositoryImpl) GetByEachAgentID(agent1ID, agent2ID uint) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	err := repo.executer.Select(
		repo.Name+".GetByEachAgentID",
		&taskList, `
			SELECT 
				task.*, task_group.job_seeker_id, task_group.job_information_id,
				seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = task_group.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			INNER JOIN
				agent_staffs AS ra_staff
			ON
				billing.agent_staff_id = ra_staff.id
			WHERE
				(
					(ca_staff.agent_id = ? AND ra_staff.agent_id = ?)
			OR
					(ra_staff.agent_id = ? AND ca_staff.agent_id = ?)
				)
			AND
				task.phase_sub_category NOT IN(90, 91, 92, 93, 94, 95, 96, 97, 98, 99)
			AND
				task.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			ORDER BY 
				task.id DESC
		`,
		agent1ID, agent2ID,
		agent1ID, agent2ID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// 全ての最新タスクを取得
func (repo *TaskRepositoryImpl) GetLatest() ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	/**
		job_information_selection_flow_patternsのレコードを外部結合で取得
	**/
	err := repo.executer.Select(
		repo.Name+".GetLatest",
		&taskList, `
			SELECT 
				task.*,
				task_group.uuid AS task_group_uuid,	task_group.is_double_sided,
				task_group.selection_flow_pattern_id, task_group.joining_date,
				task_group.ra_last_request_at, task_group.ra_last_watched_at, 
				task_group.ca_last_request_at, task_group.ca_last_watched_at, 
				seeker.id AS job_seeker_id,
				seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
				billing.how_to_recommend, 
				job_info.id AS job_information_id, job_info.billing_address_id, job_info.required_documents_detail,
				ca_staff.id AS ca_staff_id, ca_staff.staff_name AS ca_staff_name, ca_agent.agent_name AS ca_agent_name, ca_staff.agent_id AS ca_agent_id,
				ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_agent.agent_name AS ra_agent_name, ra_staff.agent_id AS ra_agent_id,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				task_group.job_seeker_id = seeker.id
			INNER JOIN
				job_informations AS job_info
			ON
				task_group.job_information_id = job_info.id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS ra_staff
			ON
				billing.agent_staff_id = ra_staff.id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			INNER JOIN
				agents AS ra_agent
			ON
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				agents AS ca_agent
			ON
				ca_staff.agent_id = ca_agent.id
			WHERE
				task.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			AND
				task.phase_sub_category NOT IN(90, 91, 92, 93, 94, 95, 96, 97, 98, 99)
			ORDER BY 
				deadline_day ASC,
				deadline_time ASC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

func (repo *TaskRepositoryImpl) GetLatestByJobSeekerUUIDAndJobInformationIDList(jobSeekerUUID uuid.UUID, jobInformationIDList []uint) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	jobInformationIDListStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]")

	query := fmt.Sprintf(`
		SELECT 
			task.*, task_group.job_information_id, task_group.job_seeker_id,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
				THEN task_group.external_job_information_title
				ELSE job_info.title
			END AS title,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
				THEN task_group.external_company_name
				ELSE enterprise.company_name
			END AS company_name,
			task_group.external_job_listing_url
		FROM
			tasks AS task
		INNER JOIN
			task_groups AS task_group
		ON
			task.task_group_id = task_group.id
		INNER JOIN
			job_seekers AS seeker
		ON
			task_group.job_seeker_id = seeker.id
		INNER JOIN
			job_informations AS job_info
		ON
			task_group.job_information_id = job_info.id
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		WHERE
			task.id = (
				SELECT task_join.id
				FROM tasks AS task_join 
				WHERE task.task_group_id = task_join.task_group_id
				ORDER BY task_join.id DESC
				LIMIT 1
			)			
		AND
			seeker.uuid = ?
		AND
			task_group.job_information_id IN (%s)
		ORDER BY 
			task_group.id DESC
	`, jobInformationIDListStr)

	err := repo.executer.Select(
		repo.Name+".GetLatestByJobSeekerUUIDAndJobInformationIDList",
		&taskList, query, jobSeekerUUID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// アクティブなタスクの一覧を取得（請求先ID）
func (repo *TaskRepositoryImpl) GetActiveByBillingAddressID(billingAddressID uint) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	err := repo.executer.Select(
		repo.Name+".GetActiveByBillingAddressID",
		&taskList, `
			SELECT 
				task.*
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_informations AS job_info
			ON
				task_group.job_information_id = job_info.id
			WHERE 
				job_info.billing_address_id = ?
			AND
				task.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			AND
				task.phase_sub_category NOT IN(90, 91, 92, 93, 94, 95, 96, 97, 98, 99)
			ORDER BY 
				deadline_day ASC,
				deadline_time ASC
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// アクティブなタスクの一覧を取得（求人ID）
func (repo *TaskRepositoryImpl) GetActiveByJobInformationID(jobInformationID uint) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	err := repo.executer.Select(
		repo.Name+".GetActiveByJobInformationID",
		&taskList, `
			SELECT 
				task.*
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			WHERE 
				task_group.job_information_id = ?
			AND
				task.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			AND
				task.phase_sub_category NOT IN(90, 91, 92, 93, 94, 95, 96, 97, 98, 99)
			ORDER BY 
				deadline_day ASC,
				deadline_time ASC
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

// アクティブなタスクの一覧を取得（選考ID）
func (repo *TaskRepositoryImpl) GetActiveBySelectionFlowPatternID(selectionID uint) ([]*entity.Task, error) {
	var (
		taskList []*entity.Task
	)

	err := repo.executer.Select(
		repo.Name+".GetActiveBySelectionFlowPatternID",
		&taskList, `
			SELECT 
				task.*
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			WHERE 
				task_group.selection_flow_pattern_id = ?
			AND
				task.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			ORDER BY 
				deadline_day ASC,
				deadline_time ASC
		`,
		selectionID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskList, nil
}

/****************************************************************************************/
// 売上管理のカウント関連
// 当月カウントベース
//

// 内定数
func (repo *TaskRepositoryImpl) GetSeekerOfferPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetSeekerOfferPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND
				DATE_FORMAT(CONVERT_TZ(task.created_at, '+00:00','+09:00'), '%Y-%m') = ?
			AND
				task.phase_category = 8
			ORDER BY 
				task.id ASC
		`,
		agentID, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 最終選考
func (repo *TaskRepositoryImpl) GetSeekerFinalSelectionPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetSeekerFinalSelectionPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND
				DATE_FORMAT(CONVERT_TZ(task.created_at, '+00:00','+09:00'), '%Y-%m') = ?
			AND
				task.phase_category = 7
			ORDER BY 
				task.id ASC
		`,
		agentID, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 選考数
func (repo *TaskRepositoryImpl) GetSeekerSelectionPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetSeekerSelectionPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND
				DATE_FORMAT(CONVERT_TZ(task.created_at, '+00:00','+09:00'), '%Y-%m') = ?
			AND
				task.phase_category IN(2, 7) AND
				task.phase_sub_category IN (0, 4, 8)
			AND EXISTS (
				SELECT 1
				FROM 
					tasks AS task_join
				WHERE 
					task_join.task_group_id = task.task_group_id
				AND 
					task_join.id < task.id
				AND 
					task_join.phase_category = 1 AND 
					task_join.phase_sub_category = 1
				ORDER BY 
					task_join.id DESC
				LIMIT 1
			)
		`,
		agentID, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 推薦完了数
func (repo *TaskRepositoryImpl) GetSeekerRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	// 「phase_category」 2: 一次選考, 7: 最終選考
	// 「phase_sub_category」 0: 候補日回収依頼, 4: 日程案内依頼（詳細未案内）, 8: 案内依頼, 14: 日程・詳細案内依頼
	err := repo.executer.Get(
		repo.Name+".GetSeekerRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND
				DATE_FORMAT(CONVERT_TZ(task.created_at, '+00:00','+09:00'), '%Y-%m') = ?
			AND (
				task.phase_category = 2 AND (
					task.phase_sub_category = 0 OR
					task.phase_sub_category = 4 OR
					task.phase_sub_category = 8 OR
					task.phase_sub_category = 14
				) OR
				task.phase_category = 7 AND (
					task.phase_sub_category = 0 OR
					task.phase_sub_category = 4 OR
					task.phase_sub_category = 8 OR
					task.phase_sub_category = 14
				)
			)
			ORDER BY 
				task.id ASC
		`,
		agentID, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 求人紹介数
func (repo *TaskRepositoryImpl) GetSeekerJobIntroductionPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetSeekerJobIntroductionPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND
				DATE_FORMAT(CONVERT_TZ(task.created_at, '+00:00','+09:00'), '%Y-%m') = ?
			AND
				task.phase_category = 0 AND
				task.phase_sub_category = 5
			ORDER BY 
				task.id ASC
		`,
		agentID, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

/****************************************************************************************/
// 売上管理のカウント関連
// 面談実施月カウントベース
//

// 内定数
func (repo *TaskRepositoryImpl) GetInterviewOfferPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewOfferPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND (
				(interview_group.first_interview_date <> '' AND LEFT(interview_group.first_interview_date, 7) = ?)
				OR (interview_group.first_interview_date = '' AND LEFT(interview_group.interview_date, 7) = ?)
			)
			AND
				task.phase_category = 8
			ORDER BY 
				task.id ASC
		`,
		agentID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 最終選考
func (repo *TaskRepositoryImpl) GetInterviewFinalSelectionPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewFinalSelectionPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND (
				(interview_group.first_interview_date <> '' AND LEFT(interview_group.first_interview_date, 7) = ?)
				OR (interview_group.first_interview_date = '' AND LEFT(interview_group.interview_date, 7) = ?)
			)
			AND
				task.phase_category = 7
			ORDER BY 
				task.id ASC
		`,
		agentID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 選考
func (repo *TaskRepositoryImpl) GetInterviewSelectionPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewSelectionPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND (
				(interview_group.first_interview_date <> '' AND LEFT(interview_group.first_interview_date, 7) = ?)
				OR (interview_group.first_interview_date = '' AND LEFT(interview_group.interview_date, 7) = ?)
			)
			AND
				task.phase_category IN(2, 7) AND
				task.phase_sub_category IN (0, 4, 8)
			AND EXISTS (
				SELECT 1
				FROM 
					tasks AS task_join
				WHERE 
					task_join.task_group_id = task.task_group_id
				AND 
					task_join.id < task.id
				AND 
					task_join.phase_category = 1 AND 
					task_join.phase_sub_category = 1
				ORDER BY 
					task_join.id DESC
				LIMIT 1
			)
		`,
		agentID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 推薦完了数
func (repo *TaskRepositoryImpl) GetInterviewRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	// 書類選考通過（推薦完了）
	// 「phase_category」 2: 一次選考, 7: 最終選考
	// 「phase_sub_category」 0: 候補日回収依頼, 4: 日程案内依頼（詳細未案内）, 8: 案内依頼, 14: 日程・詳細案内依頼
	err := repo.executer.Get(
		repo.Name+".GetInterviewRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND (
				(interview_group.first_interview_date <> '' AND LEFT(interview_group.first_interview_date, 7) = ?)
				OR (interview_group.first_interview_date = '' AND LEFT(interview_group.interview_date, 7) = ?)
			)
			AND (
				task.phase_category = 2 AND (
					task.phase_sub_category = 0 OR
					task.phase_sub_category = 4 OR
					task.phase_sub_category = 8 OR
					task.phase_sub_category = 14
				) OR
				task.phase_category = 7 AND (
					task.phase_sub_category = 0 OR
					task.phase_sub_category = 4 OR
					task.phase_sub_category = 8 OR
					task.phase_sub_category = 14
				)
			)
			ORDER BY 
				task.id ASC
		`,
		agentID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 求人紹介数
func (repo *TaskRepositoryImpl) GetInterviewJobIntroductionPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewJobIntroductionPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND (
				(interview_group.first_interview_date <> '' AND LEFT(interview_group.first_interview_date, 7) = ?)
				OR (interview_group.first_interview_date = '' AND LEFT(interview_group.interview_date, 7) = ?)
			)
			AND
				task.phase_category = 0 AND
				task.phase_sub_category = 5
			ORDER BY 
				task.id ASC
		`,
		agentID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

/****************************************************************************************/
// 売上管理のカウント関連（個人）
// 当月カウントベース
//

// 内定数
func (repo *TaskRepositoryImpl) GetSeekerOfferPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetSeekerOfferPerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				task.staff_type = 0
			AND
				task.executed_staff_id = ?
			AND
				DATE_FORMAT(CONVERT_TZ(task.created_at, '+00:00','+09:00'), '%Y-%m') = ?
			AND
				task.phase_category = 8
			ORDER BY 
				task.id ASC
		`,
		staffID, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 最終選考
func (repo *TaskRepositoryImpl) GetSeekerFinalSelectionPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetSeekerFinalSelectionPerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				task.staff_type = 0
			AND
				task.executed_staff_id = ?
			AND
				DATE_FORMAT(CONVERT_TZ(task.created_at, '+00:00','+09:00'), '%Y-%m') = ?
			AND
				task.phase_category = 7
			ORDER BY 
				task.id ASC
		`,
		staffID, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 選考数
func (repo *TaskRepositoryImpl) GetSeekerSelectionPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetSeekerSelectionPerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				task.staff_type = 0
			AND
				task.executed_staff_id = ?
			AND
				DATE_FORMAT(CONVERT_TZ(task.created_at, '+00:00','+09:00'), '%Y-%m') = ?
			AND
				task.phase_category IN(2, 7) AND
				task.phase_sub_category IN (0, 4, 8)
			AND EXISTS (
				SELECT 1
				FROM 
					tasks AS task_join
				WHERE 
					task_join.task_group_id = task.task_group_id
				AND 
					task_join.id < task.id
				AND 
					task_join.phase_category = 1 AND 
					task_join.phase_sub_category = 1
				ORDER BY 
					task_join.id DESC
				LIMIT 1
			)
		`,
		staffID, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 推薦完了数
func (repo *TaskRepositoryImpl) GetSeekerRecommendationCompletionPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	// 書類選考通過（推薦完了）
	// 「phase_category」 2: 一次選考, 7: 最終選考
	// 「phase_sub_category」 0: 候補日回収依頼, 4: 日程案内依頼（詳細未案内）, 8: 案内依頼, 14: 日程・詳細案内依頼
	err := repo.executer.Get(
		repo.Name+".GetSeekerRecommendationCompletionPerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				task.staff_type = 0
			AND
				task.executed_staff_id = ?
			AND
				DATE_FORMAT(CONVERT_TZ(task.created_at, '+00:00','+09:00'), '%Y-%m') = ?
			AND (
				task.phase_category = 2 AND (
					task.phase_sub_category = 0 OR
					task.phase_sub_category = 4 OR
					task.phase_sub_category = 8 OR
					task.phase_sub_category = 14
				) OR
				task.phase_category = 7 AND (
					task.phase_sub_category = 0 OR
					task.phase_sub_category = 4 OR
					task.phase_sub_category = 8 OR
					task.phase_sub_category = 14
				)
			)
			ORDER BY 
				task.id ASC
		`,
		staffID, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 求人紹介数
func (repo *TaskRepositoryImpl) GetSeekerJobIntroductionPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetSeekerJobIntroductionPerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				task.staff_type = 0
			AND
				task.executed_staff_id = ?
			AND
				DATE_FORMAT(CONVERT_TZ(task.created_at, '+00:00','+09:00'), '%Y-%m') = ?
			AND
				task.phase_category = 0 AND
				task.phase_sub_category = 5
			ORDER BY 
				task.id ASC
		`,
		staffID, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

/****************************************************************************************/
// 売上管理のカウント関連（個人）
// 面談実施月カウントベース
//

// 内定数
func (repo *TaskRepositoryImpl) GetInterviewOfferPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewOfferPerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = interview_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(interview_group.first_interview_date <> '' AND LEFT(interview_group.first_interview_date, 7) = ?)
				OR (interview_group.first_interview_date = '' AND LEFT(interview_group.interview_date, 7) = ?)
			)
			AND
				task.phase_category = 8
			ORDER BY 
				task.id ASC
		`,
		staffID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 最終選考
func (repo *TaskRepositoryImpl) GetInterviewFinalSelectionPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewFinalSelectionPerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = interview_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(interview_group.first_interview_date <> '' AND LEFT(interview_group.first_interview_date, 7) = ?)
				OR (interview_group.first_interview_date = '' AND LEFT(interview_group.interview_date, 7) = ?)
			)
			AND
				task.phase_category = 7
			ORDER BY 
				task.id ASC
		`,
		staffID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 選考
func (repo *TaskRepositoryImpl) GetInterviewSelectionPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewSelectionPerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = interview_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(interview_group.first_interview_date <> '' AND LEFT(interview_group.first_interview_date, 7) = ?)
				OR (interview_group.first_interview_date = '' AND LEFT(interview_group.interview_date, 7) = ?)
			)
			AND
				task.phase_category IN(2, 7) AND
				task.phase_sub_category IN (0, 4, 8)
			AND EXISTS (
				SELECT 1
				FROM 
					tasks AS task_join
				WHERE 
					task_join.task_group_id = task.task_group_id
				AND 
					task_join.id < task.id
				AND 
					task_join.phase_category = 1 AND 
					task_join.phase_sub_category = 1
				ORDER BY 
					task_join.id DESC
				LIMIT 1
			)
		`,
		staffID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 推薦完了数
func (repo *TaskRepositoryImpl) GetInterviewRecommendationCompletionPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	// 書類選考通過（推薦完了）
	// 「phase_category」 2: 一次選考, 7: 最終選考
	// 「phase_sub_category」 0: 候補日回収依頼, 4: 日程案内依頼（詳細未案内）, 8: 案内依頼, 14: 日程・詳細案内依頼
	err := repo.executer.Get(
		repo.Name+".GetInterviewRecommendationCompletionPerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = interview_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(interview_group.first_interview_date <> '' AND LEFT(interview_group.first_interview_date, 7) = ?)
				OR (interview_group.first_interview_date = '' AND LEFT(interview_group.interview_date, 7) = ?)
			)
			AND (
				task.phase_category = 2 AND (
					task.phase_sub_category = 0 OR
					task.phase_sub_category = 4 OR
					task.phase_sub_category = 8 OR
					task.phase_sub_category = 14
				) OR
				task.phase_category = 7 AND (
					task.phase_sub_category = 0 OR
					task.phase_sub_category = 4 OR
					task.phase_sub_category = 8 OR
					task.phase_sub_category = 14
				)
			)
			ORDER BY 
				task.id ASC
		`,
		staffID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 求人紹介数
func (repo *TaskRepositoryImpl) GetInterviewJobIntroductionPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewJobIntroductionPerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = interview_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(interview_group.first_interview_date <> '' AND LEFT(interview_group.first_interview_date, 7) = ?)
				OR (interview_group.first_interview_date = '' AND LEFT(interview_group.interview_date, 7) = ?)
			)
			AND
				task.phase_category = 0 AND
				task.phase_sub_category = 5
			ORDER BY 
				task.id ASC
		`,
		staffID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

/****************************************************************************************/
// ダッシュボード 売上管理のカウント関連
// 面談実施月カウントベース
//

// 内定数
func (repo *TaskRepositoryImpl) GetInterviewOfferPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewOfferPerformanceCountByAgentIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND (
				(
					interview_group.first_interview_date <> '' AND 
					LEFT(interview_group.first_interview_date, 7) >= ? AND
					LEFT(interview_group.first_interview_date, 7) <= ?
				) OR 
				(
					interview_group.first_interview_date = '' AND 
					LEFT(interview_group.interview_date, 7) >= ? AND
					LEFT(interview_group.interview_date, 7) <= ?
				)
			)
			AND
				task.phase_category = 8
			ORDER BY 
				task.id ASC
		`,
		agentID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 最終選考
func (repo *TaskRepositoryImpl) GetInterviewFinalSelectionPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewFinalSelectionPerformanceCountByAgentIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND (
				(
					interview_group.first_interview_date <> '' AND 
					LEFT(interview_group.first_interview_date, 7) >= ? AND
					LEFT(interview_group.first_interview_date, 7) <= ?
				) OR 
				(
					interview_group.first_interview_date = '' AND 
					LEFT(interview_group.interview_date, 7) >= ? AND
					LEFT(interview_group.interview_date, 7) <= ?
				)
			)
			AND
				task.phase_category = 7
			ORDER BY 
				task.id ASC
		`,
		agentID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 選考
func (repo *TaskRepositoryImpl) GetInterviewSelectionPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewSelectionPerformanceCountByAgentIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND (
				(
					interview_group.first_interview_date <> '' AND 
					LEFT(interview_group.first_interview_date, 7) >= ? AND
					LEFT(interview_group.first_interview_date, 7) <= ?
				) OR 
				(
					interview_group.first_interview_date = '' AND 
					LEFT(interview_group.interview_date, 7) >= ? AND
					LEFT(interview_group.interview_date, 7) <= ?
				)
			)
			AND
				task.phase_category IN(2, 7) AND
				task.phase_sub_category IN (0, 4, 8)
			AND EXISTS (
				SELECT 1
				FROM 
					tasks AS task_join
				WHERE 
					task_join.task_group_id = task.task_group_id
				AND 
					task_join.id < task.id
				AND 
					task_join.phase_category = 1 AND 
					task_join.phase_sub_category = 1
				ORDER BY 
					task_join.id DESC
				LIMIT 1
			)
		`,
		agentID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 推薦完了数
func (repo *TaskRepositoryImpl) GetInterviewRecommendationCompletionPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewRecommendationCompletionPerformanceCountByAgentIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND (
				(
					interview_group.first_interview_date <> '' AND 
					LEFT(interview_group.first_interview_date, 7) >= ? AND
					LEFT(interview_group.first_interview_date, 7) <= ?
				) OR 
				(
					interview_group.first_interview_date = '' AND 
					LEFT(interview_group.interview_date, 7) >= ? AND
					LEFT(interview_group.interview_date, 7) <= ?
				)
			)
			AND
				task.phase_category = 1 AND
				task.phase_sub_category = 1
			ORDER BY 
				task.id ASC
		`,
		agentID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 求人紹介数
func (repo *TaskRepositoryImpl) GetInterviewJobIntroductionPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewJobIntroductionPerformanceCountByAgentIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE
				ca_staff.agent_id = ?
			AND (
				(
					interview_group.first_interview_date <> '' AND 
					LEFT(interview_group.first_interview_date, 7) >= ? AND
					LEFT(interview_group.first_interview_date, 7) <= ?
				) OR 
				(
					interview_group.first_interview_date = '' AND 
					LEFT(interview_group.interview_date, 7) >= ? AND
					LEFT(interview_group.interview_date, 7) <= ?
				)
			)
			AND
				task.phase_category = 0 AND
				task.phase_sub_category = 5
			ORDER BY 
				task.id ASC
		`,
		agentID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

/****************************************************************************************/
// ダッシュボード 売上管理のカウント関連（個人）
// 面談実施月カウントベース
//

// 内定数
func (repo *TaskRepositoryImpl) GetInterviewOfferPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewOfferPerformanceCountByStaffIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = interview_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(
					interview_group.first_interview_date <> '' AND 
					LEFT(interview_group.first_interview_date, 7) >= ? AND
					LEFT(interview_group.first_interview_date, 7) <= ?
				) OR 
				(
					interview_group.first_interview_date = '' AND 
					LEFT(interview_group.interview_date, 7) >= ? AND
					LEFT(interview_group.interview_date, 7) <= ?
				)
			)
			AND
				task.phase_category = 8
			ORDER BY 
				task.id ASC
		`,
		staffID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 最終選考
func (repo *TaskRepositoryImpl) GetInterviewFinalSelectionPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewFinalSelectionPerformanceCountByStaffIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = interview_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(
					interview_group.first_interview_date <> '' AND 
					LEFT(interview_group.first_interview_date, 7) >= ? AND
					LEFT(interview_group.first_interview_date, 7) <= ?
				) OR 
				(
					interview_group.first_interview_date = '' AND 
					LEFT(interview_group.interview_date, 7) >= ? AND
					LEFT(interview_group.interview_date, 7) <= ?
				)
			)
			AND
				task.phase_category = 7
			ORDER BY 
				task.id ASC
		`,
		staffID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 選考
func (repo *TaskRepositoryImpl) GetInterviewSelectionPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewSelectionPerformanceCountByStaffIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = interview_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(
					interview_group.first_interview_date <> '' AND 
					LEFT(interview_group.first_interview_date, 7) >= ? AND
					LEFT(interview_group.first_interview_date, 7) <= ?
				) OR 
				(
					interview_group.first_interview_date = '' AND 
					LEFT(interview_group.interview_date, 7) >= ? AND
					LEFT(interview_group.interview_date, 7) <= ?
				)
			)
			AND
				task.phase_category IN(2, 7) AND
				task.phase_sub_category IN (0, 4, 8)
			AND EXISTS (
				SELECT 1
				FROM 
					tasks AS task_join
				WHERE 
					task_join.task_group_id = task.task_group_id
				AND 
					task_join.id < task.id
				AND 
					task_join.phase_category = 1 AND 
					task_join.phase_sub_category = 1
				ORDER BY 
					task_join.id DESC
				LIMIT 1
			)
		`,
		staffID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 推薦完了数
func (repo *TaskRepositoryImpl) GetInterviewRecommendationCompletionPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewRecommendationCompletionPerformanceCountByStaffIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = interview_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(
					interview_group.first_interview_date <> '' AND 
					LEFT(interview_group.first_interview_date, 7) >= ? AND
					LEFT(interview_group.first_interview_date, 7) <= ?
				) OR 
				(
					interview_group.first_interview_date = '' AND 
					LEFT(interview_group.interview_date, 7) >= ? AND
					LEFT(interview_group.interview_date, 7) <= ?
				)
			)
			AND
				task.phase_category = 1 AND
				task.phase_sub_category = 1
			ORDER BY 
				task.id ASC
		`,
		staffID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 求人紹介数
func (repo *TaskRepositoryImpl) GetInterviewJobIntroductionPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewJobIntroductionPerformanceCountByStaffIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				tasks AS task
			INNER JOIN
				task_groups AS task_group
			ON
				task.task_group_id = task_group.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_group.job_seeker_id
			INNER JOIN 
				interview_task_groups AS interview_group
			ON 
				task_group.job_seeker_id = interview_group.job_seeker_id
			INNER JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = interview_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(
					interview_group.first_interview_date <> '' AND 
					LEFT(interview_group.first_interview_date, 7) >= ? AND
					LEFT(interview_group.first_interview_date, 7) <= ?
				) OR 
				(
					interview_group.first_interview_date = '' AND 
					LEFT(interview_group.interview_date, 7) >= ? AND
					LEFT(interview_group.interview_date, 7) <= ?
				)
			)
			AND
				task.phase_category = 0 AND
				task.phase_sub_category = 5
			ORDER BY 
				task.id ASC
		`,
		staffID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}
