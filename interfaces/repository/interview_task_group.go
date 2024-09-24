package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type InterviewTaskGroupRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewInterviewTaskGroupRepositoryImpl(ex interfaces.SQLExecuter) usecase.InterviewTaskGroupRepository {
	return &InterviewTaskGroupRepositoryImpl{
		Name:     "InterviewTaskGroupRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// 面談調整のタスクグループを作成
func (repo *InterviewTaskGroupRepositoryImpl) Create(interviewTaskGroup *entity.InterviewTaskGroup) error {
	nowTime := time.Now().In(time.UTC)
	interviewTaskGroup.UUID = utility.CreateUUID()
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO interview_task_groups (
			uuid,
			agent_id,
			job_seeker_id,
			interview_date,
			first_interview_date,
			last_request_at,
			last_watched_at,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?
			)
		`,
		interviewTaskGroup.UUID,
		interviewTaskGroup.AgentID,
		interviewTaskGroup.JobSeekerID,
		interviewTaskGroup.InterviewDate,
		interviewTaskGroup.FirstInterviewDate,
		nowTime,
		nowTime,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	interviewTaskGroup.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
// 最終リクエスト時間を更新
func (repo *InterviewTaskGroupRepositoryImpl) UpdateLastRequestAt(id uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateLastRequestAt",
		`
		UPDATE 
			interview_task_groups 
		SET
			last_request_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime,
		nowTime,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 最終閲覧時間を更新
func (repo *InterviewTaskGroupRepositoryImpl) UpdateLastWatchedAt(id uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateLastWatchedAt",
		`
		UPDATE 
			interview_task_groups 
		SET
			last_watched_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime,
		nowTime,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 面談日時を更新
func (repo *InterviewTaskGroupRepositoryImpl) UpdateInterviewDate(id uint, interviewDate time.Time) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateInterviewDate",
		`
		UPDATE 
			interview_task_groups 
		SET
			interview_date = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		interviewDate,
		nowTime,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 初回面談日時(KPIの計算の基準になる)と面談日時を更新
func (repo *InterviewTaskGroupRepositoryImpl) UpdateNormalAndFirstInterviewDate(id uint, normalInterviewDate, firstInterviewDate time.Time) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateNormalAndFirstInterviewDate",
		`
		UPDATE 
			interview_task_groups 
		SET
			interview_date = ?,
			first_interview_date = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		normalInterviewDate,
		firstInterviewDate,
		nowTime,
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
func (repo *InterviewTaskGroupRepositoryImpl) FindByID(id uint) (*entity.InterviewTaskGroup, error) {
	var (
		InterviewTaskGroup entity.InterviewTaskGroup
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&InterviewTaskGroup, `
		SELECT 
			task_group.*, 
			seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana, seeker.phase,
			staff.id AS agent_staff_id, staff.staff_name AS staff_name
		FROM 
			interview_task_groups AS task_group
		INNER JOIN
			interview_tasks AS task
		ON
			task.interview_task_group_id = task_group.id
		INNER JOIN
			job_seekers AS seeker
		ON
			seeker.id = task_group.job_seeker_id
		INNER JOIN
			agent_staffs AS staff
		ON
			task.agent_staff_id = staff.id
		WHERE
			task_group.id = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &InterviewTaskGroup, nil
}

// IDと求職者IDから面談調整のタスクグループを取得
func (repo *InterviewTaskGroupRepositoryImpl) FindByAgentIDAndJobSeekerID(agentID, jobSeekerID uint) (*entity.InterviewTaskGroup, error) {
	var (
		InterviewTaskGroup entity.InterviewTaskGroup
	)

	err := repo.executer.Get(
		repo.Name+".FindByAgentIDAndJobSeekerID",
		&InterviewTaskGroup, `
		SELECT
			task_group.*
		FROM
			interview_task_groups AS task_group
		WHERE
			task_group.agent_id = ? 
		AND 
			task_group.job_seeker_id = ?
		LIMIT 1
		`,
		agentID,
		jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &InterviewTaskGroup, nil
}

/****************************************************************************************/
// 複数取得

/****************************************************************************************/
// カレンダー関連
//
func (repo *InterviewTaskGroupRepositoryImpl) GetByStaffIDAndPeriod(staffID uint, startDate, endDate time.Time) ([]*entity.InterviewTaskGroup, error) {
	var (
		InterviewTaskGroupList []*entity.InterviewTaskGroup
	)

	err := repo.executer.Select(
		repo.Name+".GetByStaffIDAndPeriod",
		&InterviewTaskGroupList, `
			SELECT 
				interview_group.*, seeker.last_name, seeker.first_name, seeker.agent_staff_id, IFNULL(staff.staff_name, '') AS staff_name
			FROM 
				interview_task_groups AS interview_group
			INNER JOIN 
				job_seekers AS seeker
			ON
				interview_group.job_seeker_id = seeker.id
			LEFT OUTER JOIN
				agent_staffs AS staff
			ON
				seeker.agent_staff_id = staff.id
			WHERE
				seeker.agent_staff_id = ? 
			AND (
				interview_group.interview_date > ? AND 
				interview_group.interview_date < ?
			)
		`,
		staffID, startDate, endDate,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return InterviewTaskGroupList, nil
}

func (repo *InterviewTaskGroupRepositoryImpl) GetByStaffIDAndPeriodByStaffIDList(idList []uint, startDate, endDate time.Time) ([]*entity.InterviewTaskGroup, error) {
	var (
		InterviewTaskGroupList []*entity.InterviewTaskGroup
	)

	if len(idList) == 0 {
		return InterviewTaskGroupList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			interview_group.*, 
			seeker.last_name, seeker.first_name, seeker.agent_staff_id, 
			IFNULL(staff.staff_name, '') AS staff_name,
			IFNULL(agents.agent_name, '') AS agent_name
		FROM 
			interview_task_groups AS interview_group
		INNER JOIN 
			job_seekers AS seeker
		ON
			interview_group.job_seeker_id = seeker.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		LEFT OUTER JOIN
			agents
		ON
			staff.agent_id = agents.id
		WHERE
			seeker.agent_staff_id IN(%s)
		AND (
			interview_group.first_interview_date <> '0001-01-01 01:00:00' AND (
				interview_group.interview_date > ? AND 
				interview_group.interview_date < ?
			) OR 
			interview_group.first_interview_date = '0001-01-01 01:00:00' AND (
				interview_group.interview_date > ? AND 
				interview_group.interview_date < ?
			)
		)
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ","), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByStaffIDAndPeriodByStaffIDList",
		&InterviewTaskGroupList, query, startDate, endDate, startDate, endDate,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return InterviewTaskGroupList, nil
}

/****************************************************************************************/
// 売上管理のカウント関連
// 面談実施月カウントベース
//

// 内定承諾数
func (repo *InterviewTaskGroupRepositoryImpl) GetInterviewOfferAcceptancePerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewOfferAcceptancePerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				interview_task_groups AS task_group
			LEFT OUTER JOIN
				sales AS sale
			ON
				task_group.job_seeker_id  = sale.job_seeker_id
			WHERE
				task_group.agent_id = ?
			AND (
				(task_group.first_interview_date <> '' AND LEFT(task_group.first_interview_date, 7) = ?)
				OR (task_group.first_interview_date = '' AND LEFT(task_group.interview_date, 7) = ?)
			)
			AND
				sale.accuracy = 0
		`,
		agentID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 面談数
func (repo *InterviewTaskGroupRepositoryImpl) GetInterviewPerformanceCountByAgentIDAndSalesMonth(agentID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				interview_task_groups AS task_group
			WHERE
				task_group.agent_id = ?
			AND (
				(task_group.first_interview_date <> '' AND LEFT(task_group.first_interview_date, 7) = ?)
				OR (task_group.first_interview_date = '' AND LEFT(task_group.interview_date, 7) = ?)
			)
			AND EXISTS (
				SELECT 1
				FROM interview_tasks AS task
				WHERE task.interview_task_group_id = task_group.id
				AND task.phase_category IN (4, 5, 6)
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

/****************************************************************************************/
// 売上管理のカウント関連（個人）
// 面談実施月カウントベース
//

// 内定承諾数
func (repo *InterviewTaskGroupRepositoryImpl) GetInterviewOfferAcceptancePerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewOfferAcceptancePerformanceCountByStaffIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT interview_group.job_seeker_id) AS count
			FROM 
				interview_task_groups AS interview_group
			INNER JOIN 
				job_seekers AS seeker
			ON
				interview_group.job_seeker_id = seeker.id
			LEFT OUTER JOIN
				sales AS sale
			ON
			interview_group.job_seeker_id  = sale.job_seeker_id
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
				sale.accuracy = 0
		`,
		staffID, salesMonth, salesMonth,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// 面談数
func (repo *InterviewTaskGroupRepositoryImpl) GetInterviewPerformanceCountByStaffIDAndSalesMonth(staffID uint, salesMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewPerformanceCountByAgentIDAndSalesMonth",
		&result, `
			SELECT 
				COUNT(DISTINCT interview_group.job_seeker_id) AS count
			FROM 
				interview_task_groups AS interview_group
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
// ダッシュボード売上管理のカウント関連
// 面談実施月カウントベース
//

// 内定承諾数
func (repo *InterviewTaskGroupRepositoryImpl) GetInterviewOfferAcceptancePerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewOfferAcceptancePerformanceCountByAgentIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				interview_task_groups AS task_group
			LEFT OUTER JOIN
				sales AS sale
			ON
				task_group.job_seeker_id  = sale.job_seeker_id
			WHERE
				task_group.agent_id = ?
			AND (
				(
					task_group.first_interview_date <> '' AND 
					LEFT(task_group.first_interview_date, 7) >= ? AND
					LEFT(task_group.first_interview_date, 7) <= ? 
				) OR
				(
					task_group.first_interview_date = '' AND 
					LEFT(task_group.interview_date, 7) >= ? AND
					LEFT(task_group.interview_date, 7) <= ?
				)
			)
			AND
				sale.accuracy = 0
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

// 面談数
func (repo *InterviewTaskGroupRepositoryImpl) GetInterviewPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewPerformanceCountByAgentIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				interview_task_groups AS task_group
			WHERE
				task_group.agent_id = ?
			AND (
				(
					task_group.first_interview_date <> '' AND 
					LEFT(task_group.first_interview_date, 7) >= ? AND
					LEFT(task_group.first_interview_date, 7) <= ?
				) OR 
				(
					task_group.first_interview_date = '' AND 
					LEFT(task_group.interview_date, 7) >= ? AND
					LEFT(task_group.interview_date, 7) <= ?
				)
			)
			AND EXISTS (
				SELECT 1
				FROM interview_tasks AS task
				WHERE task.interview_task_group_id = task_group.id
				AND task.phase_category IN (4, 5, 6)
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

/****************************************************************************************/
// ダッシュボード売上管理のカウント関連（個人）
// 面談実施月カウントベース
//

// 内定承諾数
func (repo *InterviewTaskGroupRepositoryImpl) GetInterviewOfferAcceptancePerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewOfferAcceptancePerformanceCountByStaffIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				interview_task_groups AS task_group
			INNER JOIN 
				job_seekers AS seeker
			ON
				task_group.job_seeker_id = seeker.id
			LEFT OUTER JOIN
				sales AS sale
			ON
				task_group.job_seeker_id  = sale.job_seeker_id
			WHERE 
				EXISTS (
					SELECT 1
					FROM interview_tasks AS interview_task
					WHERE interview_task.interview_task_group_id = task_group.id
					AND interview_task.phase_category IN (4, 5, 6)
					AND interview_task.ca_staff_id = ?
				) 
			AND (
				(
					task_group.first_interview_date <> '' AND 
					LEFT(task_group.first_interview_date, 7) >= ? AND
					LEFT(task_group.first_interview_date, 7) <= ? 
				) OR
				(
					task_group.first_interview_date = '' AND 
					LEFT(task_group.interview_date, 7) >= ? AND
					LEFT(task_group.interview_date, 7) <= ?
				)
			)
			AND
				sale.accuracy = 0
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

// 面談数
func (repo *InterviewTaskGroupRepositoryImpl) GetInterviewPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetInterviewPerformanceCountByStaffIDAndPeriod",
		&result, `
			SELECT 
				COUNT(DISTINCT task_group.job_seeker_id) AS count
			FROM 
				interview_task_groups AS task_group
			WHERE EXISTS (
				SELECT 1
				FROM interview_tasks AS interview_task
				WHERE interview_task.interview_task_group_id = task_group.id
				AND interview_task.phase_category IN (4, 5, 6)
				AND interview_task.ca_staff_id = ?
			)
			AND (
				(
					task_group.first_interview_date <> '' AND 
					LEFT(task_group.first_interview_date, 7) >= ? AND
					LEFT(task_group.first_interview_date, 7) <= ?
				) OR 
				(
					task_group.first_interview_date = '' AND 
					LEFT(task_group.interview_date, 7) >= ? AND
					LEFT(task_group.interview_date, 7) <= ?
				)
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
