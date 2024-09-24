package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type TaskGroupRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewTaskGroupRepositoryImpl(ex interfaces.SQLExecuter) usecase.TaskGroupRepository {
	return &TaskGroupRepositoryImpl{
		Name:     "TaskGroupRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// タスクの作成
func (repo *TaskGroupRepositoryImpl) Create(taskGroup *entity.TaskGroup) error {
	nowTime := time.Now().In(time.UTC)
	taskGroup.UUID = utility.CreateUUID()
	minTime := utility.EarliestTime()

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO task_groups (
			uuid,
			job_seeker_id,
			job_information_id,
			selection_flow_pattern_id,
			ra_last_request_at,
			ra_last_watched_at,
			ca_last_request_at,
			ca_last_watched_at,
			joining_date,
			is_double_sided,
			external_job_information_title,
			external_company_name,
			external_job_listing_url,
			ra_agent_id,
			ca_agent_id,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)
		`,
		taskGroup.UUID,
		taskGroup.JobSeekerID,
		taskGroup.JobInformationID,
		null.NewInt(0, false),
		minTime,
		minTime,
		minTime,
		minTime,
		taskGroup.JoiningDate,
		taskGroup.IsDoubleSided,
		taskGroup.ExternalJobInformationTitle,
		taskGroup.ExternalCompanyName,
		taskGroup.ExternalJobListingURL,
		taskGroup.RAAgentID,
		taskGroup.CAAgentID,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	taskGroup.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
// RAの最終依頼時間
func (repo *TaskGroupRepositoryImpl) UpdateRALastRequestAt(groupID uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateRALastRequestAt",
		`
		UPDATE 
			task_groups 
		SET
			ra_last_request_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime,
		nowTime,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// RAの最終閲覧時間
func (repo *TaskGroupRepositoryImpl) UpdateRALastWatchedAt(groupID uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateRALastWatchedAt",
		`
		UPDATE 
			task_groups 
		SET
			ra_last_watched_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime,
		nowTime,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// CAの最終依頼時間
func (repo *TaskGroupRepositoryImpl) UpdateCALastRequestAt(groupID uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateCALastRequestAt",
		`
		UPDATE 
			task_groups 
		SET
			ca_last_request_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime,
		nowTime,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// CAの最終閲覧時間
func (repo *TaskGroupRepositoryImpl) UpdateCALastWatchedAt(groupID uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateCALastWatchedAt",
		`
		UPDATE 
			task_groups 
		SET
			ca_last_watched_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime,
		nowTime,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// RAとCAの最終依頼時間を更新
func (repo *TaskGroupRepositoryImpl) UpdateLastRequestAt(groupID uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateLastRequestAt",
		`
		UPDATE 
			task_groups 
		SET
			ca_last_request_at = ?,
			ra_last_request_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime,
		nowTime,
		nowTime,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// SelectionFlowIDの更新
func (repo *TaskGroupRepositoryImpl) UpdateSelectionFlowPatternID(groupID uint, selectionFlowID null.Int) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateSelectionFlowPatternID",
		`
		UPDATE 
			task_groups 
		SET
			selection_flow_pattern_id = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		selectionFlowID,
		nowTime,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 入社日の更新
func (repo *TaskGroupRepositoryImpl) UpdateJoiningDate(groupID uint, joiningDate string) error {

	_, err := repo.executer.Exec(
		repo.Name+".UpdateJoiningDate",
		`
		UPDATE 
			task_groups 
		SET
			joining_date = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		joiningDate,
		time.Now().In(time.UTC),
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 両面タスクの有無を更新
func (repo *TaskGroupRepositoryImpl) UpdateIsDoubleSided(groupID uint, isDoubleSided bool) error {

	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsDoubleSided",
		`
		UPDATE 
			task_groups 
		SET
			is_double_sided = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		isDoubleSided,
		time.Now().In(time.UTC),
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 両面タスクの有無を更新
func (repo *TaskGroupRepositoryImpl) UpdateListIsDoubleSided(groupIDList []uint, isDoubleSided bool) error {

	if len(groupIDList) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE 
			task_groups 
		SET
			is_double_sided = ?,
			updated_at = ?
		WHERE 
			id IN(%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(groupIDList)), ", "), "[]"),
	)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateListIsDoubleSided",
		query,
		isDoubleSided,           // 1個目の?
		time.Now().In(time.UTC), // 2個目の?
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 外部求人項目の更新
func (repo *TaskGroupRepositoryImpl) UpdateExternalJob(groupID uint, param entity.ExternalJob) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateExternalJob",
		`
			UPDATE 
				task_groups 
			SET
				external_job_information_title = ?,
				external_company_name = ?,
				external_job_listing_url = ?,
				updated_at = ?
			WHERE 
				id = ?
		`,
		param.ExternalJobInformationTitle,
		param.ExternalCompanyName,
		param.ExternalJobListingURL,
		time.Now().In(time.UTC),
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *TaskGroupRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM task_groups
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
func (repo *TaskGroupRepositoryImpl) FindByID(id uint) (*entity.TaskGroup, error) {
	var (
		taskGroup entity.TaskGroup
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&taskGroup, `
		SELECT 
			task_group.*, 
			seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana, seeker.birthday,
			ca_staff.id AS ca_staff_id, ca_staff.staff_name AS ca_staff_name, ca_staff.agent_id AS ca_agent_id,
			ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_staff.agent_id AS ra_agent_id,
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
			task_groups AS task_group
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
			task_group.id = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &taskGroup, nil
}

/****************************************************************************************/
/// 複数取得
//
// 求人IDの一覧を使ってタスクグループの一覧を取得
func (repo *TaskGroupRepositoryImpl) GetByJobInformationIDList(idList []uint) ([]*entity.TaskGroup, error) {
	var (
		taskGroupList []*entity.TaskGroup
	)

	if len(idList) == 0 {
		return taskGroupList, nil
	}

	query := fmt.Sprintf(`
	SELECT 
		task_group.*, 
		seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
		IFNULL(ca_staff.id, 0) AS ca_staff_id,
		IFNULL(ca_staff.staff_name, '') AS ca_staff_name,
		IFNULL(ca_staff.agent_id, 0) AS ca_agent_id,
		ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_staff.agent_id AS ra_agent_id,
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
		task_groups AS task_group
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
	LEFT JOIN
		agent_staffs AS ca_staff
	ON
		seeker.agent_staff_id = ca_staff.id
	INNER JOIN
		agent_staffs AS ra_staff
	ON
		billing.agent_staff_id = ra_staff.id
	WHERE
		task_group.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&taskGroupList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskGroupList, nil
}

// 求職者IDの一覧を使ってタスクグループの一覧を取得
func (repo *TaskGroupRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.TaskGroup, error) {
	var (
		taskGroupList []*entity.TaskGroup
	)

	if len(idList) == 0 {
		return taskGroupList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			task_group.*, 
			seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
			IFNULL(ca_staff.id, 0) AS ca_staff_id,
			IFNULL(ca_staff.staff_name, '') AS ca_staff_name,
			IFNULL(ca_staff.agent_id, 0) AS ca_agent_id,
			ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_staff.agent_id AS ra_agent_id,
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
			task_groups AS task_group
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
		LEFT JOIN
			agent_staffs AS ca_staff
		ON
			seeker.agent_staff_id = ca_staff.id
		INNER JOIN
			agent_staffs AS ra_staff
		ON
			billing.agent_staff_id = ra_staff.id
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ","), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&taskGroupList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskGroupList, nil
}

// 求職者IDでタスクグループの一覧を取得
func (repo *TaskGroupRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.TaskGroup, error) {
	var (
		taskGroupList []*entity.TaskGroup
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&taskGroupList,
		`
			SELECT 
				task_group.*, 
				seeker.last_name, seeker.first_name, seeker.last_furigana, seeker.first_furigana,
				IFNULL(ca_staff.id, 0) AS ca_staff_id,
				IFNULL(ca_staff.staff_name, '') AS ca_staff_name,
				IFNULL(ca_staff.agent_id, 0) AS ca_agent_id,
				ra_staff.id AS ra_staff_id, ra_staff.staff_name AS ra_staff_name, ra_staff.agent_id AS ra_agent_id,
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
				task_groups AS task_group
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
			LEFT JOIN
				agent_staffs AS ca_staff
			ON
				seeker.agent_staff_id = ca_staff.id
			INNER JOIN
				agent_staffs AS ra_staff
			ON
				billing.agent_staff_id = ra_staff.id
			WHERE
				job_seeker_id = ?
		`, jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskGroupList, nil
}

// is_double_sidedがfalseでra_staff_idとca_staff_idが一致しているtaskGroup
func (repo *TaskGroupRepositoryImpl) GetNotDoubleSidedSameRAAndCAByStaffID(agentStaffID uint) ([]*entity.TaskGroup, error) {
	var (
		taskGroupList []*entity.TaskGroup
	)

	err := repo.executer.Select(
		repo.Name+".GetNotDoubleSidedSameRAAndCAByStaffID",
		&taskGroupList, `
			SELECT 
				task_group.*
			FROM 
				task_groups AS task_group
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
				task_group.is_double_sided = FALSE
			AND 
				ca_staff.id = ? AND
				ra_staff.id = ?
		`, agentStaffID, agentStaffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskGroupList, nil
}

func (repo *TaskGroupRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.TaskGroup, error) {
	var (
		taskGroupList []*entity.TaskGroup
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&taskGroupList, `
		SELECT 
			task_group.id,
			task_group.uuid,	
			task_group.is_double_sided,
			task_group.joining_date,
			task_group.ra_last_request_at, task_group.ra_last_watched_at, 
			task_group.ca_last_request_at, task_group.ca_last_watched_at, 
			task_group.job_information_id,
			task_group.job_seeker_id,
			task_group.external_job_information_title,
			task_group.external_company_name,
			task_group.external_job_listing_url,
			task_group.selection_flow_pattern_id
		FROM 
			task_groups AS task_group
		WHERE 
			task_group.ra_agent_id = ? 
		OR
			task_group.ca_agent_id = ?
		ORDER BY 
			task_group.id DESC
		`,
		agentID,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return taskGroupList, nil
}
