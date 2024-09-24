package repository

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentStaffRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewAgentStaffRepositoryImpl(ex interfaces.SQLExecuter) usecase.AgentStaffRepository {
	return &AgentStaffRepositoryImpl{
		Name:     "AgentStaffRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *AgentStaffRepositoryImpl) Create(agentStaff *entity.AgentStaff) error {
	timeTime := time.Time{}
	minTime := timeTime.Format("2006-01-02")

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO agent_staffs (
			uuid,
			agent_id,
			firebase_id,
			authority,
			staff_name,
			furigana,
			email,
			staff_phone_number,
			department,
			position,
			remarks,
			usage_status,
			notification,
			notification_job_seeker,
			notification_unwatched,
			last_login,
			usage_start_date,
			usage_end_date,
			is_deleted,
			created_at,
			updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
		`,
		utility.CreateUUID(),
		agentStaff.AgentID,
		agentStaff.FirebaseID,
		agentStaff.Authority,
		agentStaff.StaffName,
		agentStaff.Furigana,
		agentStaff.Email,
		agentStaff.StaffPhoneNumber,
		agentStaff.Department,
		agentStaff.Position,
		agentStaff.Remarks,
		agentStaff.UsageStatus,
		null.NewInt(0, true), // notification *未使用だがnot nullなので
		agentStaff.NotificationJobSeeker,
		agentStaff.NotificationUnwatched,
		time.Now().In(time.UTC),
		time.Now().In(utility.Tokyo),
		minTime,
		false, // is_deleted
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	agentStaff.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
func (repo *AgentStaffRepositoryImpl) Update(id uint, agentStaff *entity.AgentStaff) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE agent_staffs 
		SET
			authority = ?,
			staff_name = ?,
			furigana = ?,
			staff_phone_number = ?,
			department = ?,
			position = ?,
			remarks = ?,
			usage_status = ?,
			notification = ?,
			notification_job_seeker = ?,
			notification_unwatched = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		agentStaff.Authority,
		agentStaff.StaffName,
		agentStaff.Furigana,
		agentStaff.StaffPhoneNumber,
		agentStaff.Department,
		agentStaff.Position,
		agentStaff.Remarks,
		agentStaff.UsageStatus,
		null.NewInt(0, true), // notification *未使用だがnot nullなので
		agentStaff.NotificationJobSeeker,
		agentStaff.NotificationUnwatched,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *AgentStaffRepositoryImpl) UpdateEmail(id uint, email string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateEmail",
		`
		UPDATE agent_staffs 
		SET
			email = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		email,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}
func (repo *AgentStaffRepositoryImpl) UpdateUsageEnd(id uint) error {
	nowTime := time.Now()

	_, err := repo.executer.Exec(
		repo.Name+".UpdateUsageEnd",
		`
		UPDATE agent_staffs 
		SET
			usage_status = 1,
			usage_end_date = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime.In(utility.Tokyo), // JST (日本標準時) に変換
		nowTime.In(time.UTC),
		id,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *AgentStaffRepositoryImpl) UpdateUsageStart(id uint) error {
	nowTime := time.Now()

	_, err := repo.executer.Exec(
		repo.Name+".UpdateUsageStart",
		`
		UPDATE agent_staffs 
		SET
			usage_status = 0,
			usage_start_date = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime.In(utility.Tokyo), // JST (日本標準時) に変換
		nowTime.In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *AgentStaffRepositoryImpl) UpdateNotificationJobSeeker(id uint, notificationJobSeeker bool) error {

	_, err := repo.executer.Exec(
		repo.Name+".UpdateNotificationJobSeeker",
		`
		UPDATE agent_staffs 
		SET
			notification_job_seeker = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		notificationJobSeeker,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *AgentStaffRepositoryImpl) UpdateNotificationUnwatched(id uint, notificationUnwatched bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateNotificationUnwatched",
		`
		UPDATE agent_staffs 
		SET
			notification_unwatched = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		notificationUnwatched,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *AgentStaffRepositoryImpl) UpdateAuthority(id uint, authority uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateAuthority",
		`
		UPDATE agent_staffs 
		SET
			authority = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		authority,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *AgentStaffRepositoryImpl) UpdateLastLogin(id uint, lastLoginTime time.Time) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateLastLogin",
		`
		UPDATE agent_staffs 
		SET
			last_login = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		lastLoginTime,
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
/// 削除 API
//
// usage_status と is_deleted を更新
func (repo *AgentStaffRepositoryImpl) Delete(id uint) error {
	nowTime := time.Now()

	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		UPDATE agent_staffs 
		SET
			usage_status = 1,
			is_deleted = true,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime.In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 単数取得 API
//
func (repo *AgentStaffRepositoryImpl) FindByID(userID uint) (*entity.AgentStaff, error) {
	var (
		agentStaff entity.AgentStaff
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&agentStaff, `
		SELECT 
			staff.*, 
			agent.agent_name, agent.uuid AS agent_uuid
		FROM 
			agent_staffs AS staff
		INNER JOIN 
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE
			staff.id = ?
		LIMIT 1
		`,
		userID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &agentStaff, nil
}

func (repo *AgentStaffRepositoryImpl) FindByFirebaseID(firebaseID string) (*entity.AgentStaff, error) {
	var (
		agentStaff entity.AgentStaff
	)

	err := repo.executer.Get(
		repo.Name+".FindByFirebaseID",
		&agentStaff, `
		SELECT 
			staff.*, 
			agent.agent_name, agent.line_bot_id, agent.uuid AS agent_uuid, 
			agent.is_crm_active, agent.is_alliance_active, agent.is_sending_active, agent.sending_type
		FROM 
			agent_staffs AS staff
		INNER JOIN 
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE
			firebase_id = ?
		LIMIT 1
		`, firebaseID)
	if err != nil {
		return nil, err
	}

	return &agentStaff, nil
}

func (repo *AgentStaffRepositoryImpl) FindByJobSeekerID(jobSeekerID uint) (*entity.AgentStaff, error) {
	var (
		agentStaff entity.AgentStaff
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerID",
		&agentStaff, `
		SELECT 
			staff.*, 
			agent.agent_name
		FROM 
			agent_staffs AS staff
		INNER JOIN 
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE
			staff.id = (
				SELECT agent_staff_id
				FROM job_seekers
				WHERE id = ?
			)
		LIMIT 1
		`,
		jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &agentStaff, nil
}

func (repo *AgentStaffRepositoryImpl) FindByName(name string) (*entity.AgentStaff, error) {
	var (
		agentStaff entity.AgentStaff
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerID",
		&agentStaff, `
			SELECT *
			FROM agent_staffs
			WHERE staff_name = ?
			LIMIT 1
		`,
		name,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &agentStaff, nil
}

func (repo *AgentStaffRepositoryImpl) FindStaffAndAgentLine(id uint) (*entity.AgentStaff, error) {
	var (
		agentStaff entity.AgentStaff
	)

	err := repo.executer.Get(
		repo.Name+".FindStaffAndAgentLine",
		&agentStaff, `
		SELECT 
			staff.*, 
			agent.agent_name, 
			agent.line_messaging_channel_secret,
			agent.line_messaging_channel_access_token
		FROM 
			agent_staffs AS staff
		INNER JOIN 
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE
			staff.id = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &agentStaff, nil
}

func (repo *AgentStaffRepositoryImpl) FindStaffNameByJobSeekerID(jobSeekerID uint) (string, error) {
	var (
		agentStaffName string
	)

	err := repo.executer.Get(
		repo.Name+".FindStaffNameByJobSeekerID",
		&agentStaffName, `
		SELECT 
			staff.staff_name
		FROM 
			agent_staffs AS staff
		INNER JOIN 
			job_seekers AS job_seeker
		ON
			staff.id = job_seeker.agent_staff_id
		WHERE 
			job_seeker.id = ?
		LIMIT 1
		`,
		jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return agentStaffName, nil
}

/****************************************************************************************/
/// 複数取得 API
//

func (repo *AgentStaffRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.AgentStaff, error) {
	var (
		agentStaffUsers []*entity.AgentStaff
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&agentStaffUsers, `
		SELECT 
			staff.*, agent.agent_name
		FROM 
			agent_staffs AS staff
		INNER JOIN 
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE 
			staff.agent_id = ?
		ORDER BY 
			last_login DESC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentStaffUsers, nil
}

// 利用可能な担当者の一覧を取得
func (repo *AgentStaffRepositoryImpl) GetByNotIDAndAgentIDAndAllianceAgentID(id, agentID, allianceAgentID uint) ([]*entity.AgentStaff, error) {
	var (
		agentStaffUsers []*entity.AgentStaff
	)

	// エージェントチャットで使用する項目のみ
	err := repo.executer.Select(
		repo.Name+".GetByAgentNotIDAndAllianceAgentID",
		&agentStaffUsers, `
		SELECT 
			id,
			agent_id,
			staff_name
		FROM 
			agent_staffs
		WHERE 
			agent_id IN (?, ?)
		AND
			NOT(id = ?)
		AND
			usage_status = 0
		ORDER BY
			last_login DESC
		`,
		agentID, allianceAgentID, id,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentStaffUsers, nil
}

func (repo *AgentStaffRepositoryImpl) GetByAgentIDAndNotManagementID(agentID, managementID uint) ([]*entity.AgentStaff, error) {
	var (
		agentStaffList []*entity.AgentStaff
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDAndNotManagementID",
		&agentStaffList,
		`SELECT
			staff.*
		FROM
			agent_staffs AS staff
		WHERE
			staff.agent_id = ?
		AND
			staff.id NOT IN (
				SELECT
					agent_staff_id
				FROM
					agent_staff_sale_managements
				WHERE
					management_id = ?
			)
		ORDER BY 
			staff.last_login DESC
		`,
		agentID, managementID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentStaffList, nil
}

// usage_status == 利用可能(0) の担当者一覧を取得
func (repo *AgentStaffRepositoryImpl) GetByAgentIDAndUsageStatusAvailable(agentID uint) ([]*entity.AgentStaff, error) {
	var (
		agentStaffUsers []*entity.AgentStaff
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDAndUsageStatusAvailable",
		&agentStaffUsers, `
		SELECT 
			staff.*, agent.agent_name
		FROM 
			agent_staffs AS staff
		INNER JOIN 
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE 
			staff.agent_id = ?
		AND
		  staff.usage_status = 0
		AND
			staff.is_deleted = false
		ORDER BY 
			last_login DESC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentStaffUsers, nil
}

// ※削除されていない担当者の一覧を取得
func (repo *AgentStaffRepositoryImpl) GetByAgentIDAndIsDeletedFalse(agentID uint) ([]*entity.AgentStaff, error) {
	var (
		agentStaffUsers []*entity.AgentStaff
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDAndIsDeletedFalse",
		&agentStaffUsers, `
		SELECT 
			staff.*, agent.agent_name
		FROM 
			agent_staffs AS staff
		INNER JOIN 
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE 
			staff.agent_id = ?
		AND
			staff.is_deleted = false
		ORDER BY 
			last_login DESC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentStaffUsers, nil
}

// 全ての担当者一覧を取得
func (repo *AgentStaffRepositoryImpl) All() ([]*entity.AgentStaff, error) {
	var (
		agentStaffUsers []*entity.AgentStaff
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&agentStaffUsers, `
			SELECT 
				staff.*, agent.agent_name
			FROM 
				agent_staffs AS staff
			INNER JOIN 
				agents AS agent
			ON
				staff.agent_id = agent.id
			ORDER BY id ASC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentStaffUsers, nil
}
