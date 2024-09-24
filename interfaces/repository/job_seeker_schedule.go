package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerScheduleRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerScheduleRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerScheduleRepository {
	return &JobSeekerScheduleRepositoryImpl{
		Name:     "JobSeekerScheduleRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerScheduleRepositoryImpl) Create(jobSeekerSchedule *entity.JobSeekerSchedule) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO job_seeker_schedules (
			job_seeker_id,	           
			task_id,	                   
			schedule_type, 
			title,                    
			start_time,      
			end_time,        
			seeker_description,       
			staff_description,   
			is_share,
			repetition_count,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?
			)`,
		jobSeekerSchedule.JobSeekerID,
		jobSeekerSchedule.TaskID,
		jobSeekerSchedule.ScheduleType,
		jobSeekerSchedule.Title,
		jobSeekerSchedule.StartTime,
		jobSeekerSchedule.EndTime,
		jobSeekerSchedule.SeekerDescription,
		jobSeekerSchedule.StaffDescription,
		jobSeekerSchedule.IsShare,
		jobSeekerSchedule.RepetitionCount,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	jobSeekerSchedule.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
func (repo *JobSeekerScheduleRepositoryImpl) Update(id uint, jobSeekerSchedule *entity.JobSeekerSchedule) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE job_seeker_schedules
		SET               
			schedule_type = ?,  
			title = ?,                  
			start_time = ?,      
			end_time = ?,        
			seeker_description = ?,       
			staff_description = ?,     
			is_share = ?,  
			repetition_count = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		jobSeekerSchedule.ScheduleType,
		jobSeekerSchedule.Title,
		jobSeekerSchedule.StartTime,
		jobSeekerSchedule.EndTime,
		jobSeekerSchedule.SeekerDescription,
		jobSeekerSchedule.StaffDescription,
		jobSeekerSchedule.IsShare,
		jobSeekerSchedule.RepetitionCount,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *JobSeekerScheduleRepositoryImpl) UpdateTitle(id uint, title string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateTitle",
		`
		UPDATE job_seeker_schedules
		SET               
			title = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		title,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *JobSeekerScheduleRepositoryImpl) UpdateIsSharedByIDList(idList []uint) error {
	if len(idList) == 0 {
		errMessage := fmt.Sprintf("シェアするスケジュールがありません")
		wrapped := fmt.Errorf("%s:%w", errMessage, entity.ErrRequestError)
		return wrapped
	}

	query := fmt.Sprintf(`
		UPDATE job_seeker_schedules
		SET    
			is_share = true,  
			updated_at = ?
		WHERE 
			id IN(%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ","), "[]"))

	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsSharedByIDList",
		query,
		time.Now().In(time.UTC),
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
func (repo *JobSeekerScheduleRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM job_seeker_schedules
		WHERE id = ?
		`, id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *JobSeekerScheduleRepositoryImpl) DeleteByTaskGroupIDAndScheduleType(taskGroupID, scheduleType uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByTaskGroupIDAndScheduleType",
		`
		DELETE FROM 
			job_seeker_schedules
		WHERE
			task_id IN (
				SELECT task_join.id
				FROM tasks as task_join
				WHERE task_join.task_group_id = ?
			)
		AND
			schedule_type = ?
		`, taskGroupID, scheduleType,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 指定task_group_idでかつ、指定IDまでのタスクを全て削除
func (repo *JobSeekerScheduleRepositoryImpl) DeleteScheduleInTaskGroupAboveTaskID(taskID, taskGroupID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteScheduleInTaskGroupAboveTaskID",
		`
		DELETE FROM 
			job_seeker_schedules
		WHERE
			task_id IN (
				SELECT tasks.id
				FROM tasks
				WHERE tasks.task_group_id = ?
			)
		AND
			task_id > ?
		`, taskGroupID, taskID,
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
func (repo *JobSeekerScheduleRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerSchedule, error) {
	var (
		jobSeekerScheduleList []*entity.JobSeekerSchedule
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&jobSeekerScheduleList, `
		SELECT 
			schedule.*, tasks.task_group_id, tasks.phase_category, reschedule.reschedule_id
		FROM 
			job_seeker_schedules AS schedule
		LEFT OUTER JOIN
			tasks
		ON
			schedule.task_id = tasks.id
		LEFT OUTER JOIN
			job_seeker_reschedules AS reschedule
		ON
			schedule.id = reschedule.reschedule_id
		WHERE 
			schedule.job_seeker_id = ?
		AND
			reschedule.reschedule_id IS NULL
		ORDER BY
			schedule.schedule_type DESC,
			schedule.start_time ASC
		`, jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerScheduleList, nil
}

func (repo *JobSeekerScheduleRepositoryImpl) GetByTaskGroupID(taskGroupID uint) ([]*entity.JobSeekerSchedule, error) {
	var (
		jobSeekerScheduleList []*entity.JobSeekerSchedule
	)

	err := repo.executer.Select(
		repo.Name+".GetByTaskGroupID",
		&jobSeekerScheduleList, `
		SELECT 
			schedule.*, tasks.task_group_id, tasks.phase_category, reschedule.reschedule_id
		FROM 
			job_seeker_schedules AS schedule
		LEFT OUTER JOIN
			tasks
		ON
			schedule.task_id = tasks.id
		LEFT OUTER JOIN
			job_seeker_reschedules AS reschedule
		ON
			schedule.id = reschedule.reschedule_id
		WHERE 
			tasks.task_group_id = ?
		AND
			reschedule.reschedule_id IS NULL
		ORDER BY
			schedule.start_time ASC
		`, taskGroupID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerScheduleList, nil
}

func (repo *JobSeekerScheduleRepositoryImpl) GetByJobSeekerIDAndScheuldType(jobSeekerID, scheduleType uint) ([]*entity.JobSeekerSchedule, error) {
	var (
		jobSeekerScheduleList []*entity.JobSeekerSchedule
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDAndScheuldType",
		&jobSeekerScheduleList, `
		SELECT 
			schedule.*, tasks.task_group_id, tasks.phase_category, reschedule.reschedule_id
		FROM 
			job_seeker_schedules AS schedule
		LEFT OUTER JOIN
			tasks
		ON
			schedule.task_id = tasks.id
		LEFT OUTER JOIN
			job_seeker_reschedules AS reschedule
		ON
			schedule.id = reschedule.reschedule_id
		WHERE 
			schedule.job_seeker_id = ?
		AND
			schedule.schedule_type = ?
		AND
			reschedule.reschedule_id IS NULL
		ORDER BY
			schedule.start_time ASC
		`, jobSeekerID, scheduleType,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerScheduleList, nil
}

func (repo *JobSeekerScheduleRepositoryImpl) GetByTaskGroupIDAndScheduleType(taskGroupID, scheduleType uint) ([]*entity.JobSeekerSchedule, error) {
	var (
		jobSeekerScheduleList []*entity.JobSeekerSchedule
	)

	err := repo.executer.Select(
		repo.Name+".GetByTaskGroupIDAndScheduleType",
		&jobSeekerScheduleList, `
		SELECT 
			schedule.*, tasks.task_group_id, tasks.phase_category, reschedule.reschedule_id
		FROM 
			job_seeker_schedules AS schedule
		LEFT OUTER JOIN
			tasks
		ON
			schedule.task_id = tasks.id
		LEFT OUTER JOIN
			job_seeker_reschedules AS reschedule
		ON
			schedule.id = reschedule.reschedule_id
		WHERE 
			tasks.task_group_id = ?
		AND
			schedule.schedule_type = ?
		AND
			reschedule.reschedule_id IS NULL
		ORDER BY
			schedule.id DESC
		`, taskGroupID, scheduleType,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerScheduleList, nil
}

func (repo *JobSeekerScheduleRepositoryImpl) GetByStaffIDAndPeriod(agentStaffID, scheduleType uint, startDate, endDate time.Time) ([]*entity.JobSeekerSchedule, error) {
	var (
		jobSeekerScheduleList []*entity.JobSeekerSchedule
	)

	err := repo.executer.Select(
		repo.Name+".GetByStaffIDAndPeriod",
		&jobSeekerScheduleList, `
		SELECT 
			schedule.*, tasks.task_group_id, tasks.phase_category, reschedule.reschedule_id,
			seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
			seeker.agent_staff_id AS ca_staff_id,
			ra_staff.staff_name AS ra_staff_name, ca_staff.staff_name AS ca_staff_name,
			ra_agent.agent_name AS ra_agent_name, ca_agent.agent_name AS ca_agent_name,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
				THEN task_group.external_company_name
				ELSE enterprise.company_name
			END AS company_name
		FROM 
			job_seeker_schedules AS schedule
		INNER JOIN
			tasks
		ON
			schedule.task_id = tasks.id
		LEFT OUTER JOIN
			job_seeker_reschedules AS reschedule
		ON
			schedule.id = reschedule.reschedule_id
		INNER JOIN 
			task_groups AS task_group 
		ON
			tasks.task_group_id = task_group.id
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
			seeker.agent_staff_id = ?
		AND
			schedule.schedule_type = ?
		AND
			STR_TO_DATE(schedule.start_time, '%Y-%m-%dT%H:%i') > ?
		AND
			STR_TO_DATE(schedule.start_time, '%Y-%m-%dT%H:%i') < ?
		AND
			reschedule.reschedule_id IS NULL
		ORDER BY
			schedule.start_time ASC
		`, agentStaffID, scheduleType, startDate, endDate,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerScheduleList, nil
}

func (repo *JobSeekerScheduleRepositoryImpl) GetByStaffIDAndPeriodByStaffIDList(idList []uint, scheduleType uint, startDate, endDate time.Time) ([]*entity.JobSeekerSchedule, error) {
	var (
		jobSeekerScheduleList []*entity.JobSeekerSchedule
		format                = "'%Y-%m-%dT%H:%i'"
	)

	if len(idList) == 0 {
		return jobSeekerScheduleList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			schedule.*, tasks.task_group_id, tasks.phase_category, reschedule.reschedule_id,
			seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
			seeker.agent_staff_id AS ca_staff_id,
			ra_staff.staff_name AS ra_staff_name, ca_staff.staff_name AS ca_staff_name,
			ra_agent.agent_name AS ra_agent_name, ca_agent.agent_name AS ca_agent_name,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
				THEN task_group.external_company_name
				ELSE enterprise.company_name
			END AS company_name
		FROM 
			job_seeker_schedules AS schedule
		INNER JOIN
			tasks
		ON
			schedule.task_id = tasks.id
		LEFT OUTER JOIN
			job_seeker_reschedules AS reschedule
		ON
			schedule.id = reschedule.reschedule_id
		INNER JOIN 
			task_groups AS task_group 
		ON
			tasks.task_group_id = task_group.id
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
			seeker.agent_staff_id IN(%s)
		AND
			schedule.schedule_type = ?
		AND
			STR_TO_DATE(schedule.start_time, %s) > ?
		AND
			STR_TO_DATE(schedule.start_time, %s) < ?
		AND
			reschedule.reschedule_id IS NULL
		ORDER BY
			schedule.start_time ASC
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ","), "[]"),
		format, format)

	err := repo.executer.Select(
		repo.Name+".GetByStaffIDAndPeriodByStaffIDList",
		&jobSeekerScheduleList, query, scheduleType, startDate, endDate,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerScheduleList, nil
}

func (repo *JobSeekerScheduleRepositoryImpl) GetByTaskGroupIDList(groupIDList []uint) ([]*entity.JobSeekerSchedule, error) {
	var (
		jobSeekerScheduleList []*entity.JobSeekerSchedule
	)

	if len(groupIDList) == 0 {
		return jobSeekerScheduleList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			schedule.*, tasks.task_group_id, tasks.phase_category, reschedule.reschedule_id
		FROM 
			job_seeker_schedules AS schedule
		LEFT OUTER JOIN
			tasks
		ON
			schedule.task_id = tasks.id
		LEFT OUTER JOIN
			job_seeker_reschedules AS reschedule
		ON
			schedule.id = reschedule.reschedule_id
		WHERE 
			tasks.task_group_id IN(%s)
		AND
			reschedule.reschedule_id IS NULL
		ORDER BY
			schedule.schedule_type DESC,
			schedule.start_time ASC
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(groupIDList)), ","), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByTaskGroupIDList",
		&jobSeekerScheduleList, query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerScheduleList, nil
}
