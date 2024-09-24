package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type EnterpriseProfileRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewEnterpriseProfileRepositoryImpl(ex interfaces.SQLExecuter) usecase.EnterpriseProfileRepository {
	return &EnterpriseProfileRepositoryImpl{
		Name:     "EnterpriseProfileRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//

func (repo *EnterpriseProfileRepositoryImpl) Create(enterprise *entity.EnterpriseProfile) error {
	now := time.Now().In(time.UTC)
	enterprise.UUID = utility.CreateUUID()

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO enterprise_profiles (
			uuid,
			company_name,
			agent_staff_id,
			corporate_site_url,
			representative,
			establishment,
			post_code,
			office_location,
			employee_number_single,
			employee_number_group,
			capital,
			public_offering,
			earnings_year,
			earnings,
			business_detail,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		enterprise.UUID,
		enterprise.CompanyName,
		enterprise.AgentStaffID,
		enterprise.CorporateSiteURL,
		enterprise.Representative,
		enterprise.Establishment,
		enterprise.PostCode,
		enterprise.OfficeLocation,
		enterprise.EmployeeNumberSingle,
		enterprise.EmployeeNumberGroup,
		enterprise.Capital,
		enterprise.PublicOffering,
		enterprise.EarningsYear,
		enterprise.Earnings,
		enterprise.BusinessDetail,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	enterprise.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
func (repo *EnterpriseProfileRepositoryImpl) Update(id uint, enterprise *entity.EnterpriseProfile) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE enterprise_profiles 
		SET
			company_name 						= ?,
			agent_staff_id 					= ?,
			corporate_site_url 			= ?,
			representative 					= ?,
			establishment 					= ?,
			post_code 							= ?,
			office_location 				= ?,
			employee_number_single 	= ?,
			employee_number_group 	= ?,
			capital 								= ?,
			public_offering 				= ?,
			earnings_year 					= ?,
			earnings 								= ?,
			business_detail 				= ?,
			updated_at 							= ?
		WHERE 
			id = ?
		`,
		enterprise.CompanyName,
		enterprise.AgentStaffID,
		enterprise.CorporateSiteURL,
		enterprise.Representative,
		enterprise.Establishment,
		enterprise.PostCode,
		enterprise.OfficeLocation,
		enterprise.EmployeeNumberSingle,
		enterprise.EmployeeNumberGroup,
		enterprise.Capital,
		enterprise.PublicOffering,
		enterprise.EarningsYear,
		enterprise.Earnings,
		enterprise.BusinessDetail,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *EnterpriseProfileRepositoryImpl) UpdateAgentStaffID(id, agentStaffID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE enterprise_profiles 
		SET
			agent_staff_id 					= ?,
			updated_at 							= ?
		WHERE 
			id = ?
		`,
		agentStaffID,
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

func (repo *EnterpriseProfileRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM enterprise_profiles
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
/// 単数取得 API
//

func (repo *EnterpriseProfileRepositoryImpl) FindByID(id uint) (*entity.EnterpriseProfile, error) {
	var (
		enterprise entity.EnterpriseProfile
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&enterprise, `
		SELECT 
			enterprise.*, staff.staff_name, staff.agent_id
		FROM 
			enterprise_profiles AS enterprise
		INNER JOIN
			agent_staffs AS staff
		ON
			enterprise.agent_staff_id = staff.id
		WHERE
			enterprise.id = ?
		LIMIT 1
		`, id)

	if err != nil {
		return nil, err
	}

	return &enterprise, nil
}

// 求人IDから企業情報を取得(求人票作成用)
func (repo *EnterpriseProfileRepositoryImpl) FindByJobInformationID(jobInformationID uint) (*entity.EnterpriseProfile, error) {
	var (
		enterprise entity.EnterpriseProfile
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobInformationID",
		&enterprise, `
		SELECT 
			enterprise.*, staff.staff_name, staff.agent_id
		FROM 
			enterprise_profiles AS enterprise
		INNER JOIN
			agent_staffs AS staff
		ON
			enterprise.agent_staff_id = staff.id
		WHERE
			id = (
				SELECT enterprise_id
				FROM billing_addresses
				WHERE id = (
					SELECT billing_address_id
					FROM job_informations AS job_info
					WHERE id = ?
				)
			)
		LIMIT 1
		`, jobInformationID)

	if err != nil {
		return nil, err
	}

	return &enterprise, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// エージェントIDで企業情報を取得
func (repo *EnterpriseProfileRepositoryImpl) GetByAgentStaffID(agentStaffID uint) ([]*entity.EnterpriseProfile, error) {
	var (
		enterpriseList []*entity.EnterpriseProfile
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentStaffID",
		&enterpriseList, `
		SELECT 
			enterprise.*, staff.staff_name, staff.agent_id
		FROM 
			enterprise_profiles AS enterprise
		INNER JOIN
			agent_staffs AS staff
		ON
			enterprise.agent_staff_id = staff.id
		WHERE
			enterprise.agent_staff_id = ?
		`,
		agentStaffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return enterpriseList, nil
}

// agentIDから企業一覧を取得
func (repo *EnterpriseProfileRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.EnterpriseProfile, error) {
	var (
		enterpriseList []*entity.EnterpriseProfile
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&enterpriseList, `
			SELECT 
				enterprise.*, staff.staff_name, staff.agent_id
			FROM 
				enterprise_profiles AS enterprise
			INNER JOIN
				agent_staffs AS staff
			ON
			 	enterprise.agent_staff_id = staff.id
			WHERE
				staff.agent_id = ?
			ORDER BY id DESC
			`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return enterpriseList, nil
}

// 絞り込み検索
func (repo *EnterpriseProfileRepositoryImpl) GetByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.EnterpriseProfile, error) {
	var (
		enterpriseList []*entity.EnterpriseProfile
	)
	if freeWord != "" {
		// フリーワードがある場合
		err := repo.executer.Select(
			repo.Name+".GetByAgentIDAndFreeWord",
			&enterpriseList, `
			SELECT 
				enterprise.* , staff.staff_name, staff.agent_id
			FROM 
				enterprise_profiles AS enterprise
			INNER JOIN
				agent_staffs AS staff
			ON
			 enterprise.agent_staff_id = staff.id
			WHERE
				staff.agent_id = ?
			AND
				(MATCH (enterprise.company_name) AGAINST (? IN BOOLEAN MODE)
			OR
				enterprise.id = ?)
			ORDER BY id DESC
			`,
			agentID, freeWord, freeWord,
		)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

	} else {
		// フリーワードが無い場合
		err := repo.executer.Select(
			repo.Name+".GetByAgentIDAndFreeWord",
			&enterpriseList, `
			SELECT 
				enterprise.*, staff.staff_name, staff.agent_id
			FROM 
				enterprise_profiles AS enterprise
			INNER JOIN
				agent_staffs AS staff
			ON
				enterprise.agent_staff_id = staff.id
			WHERE
				staff.agent_id = ?
			ORDER BY id DESC
			`,
			agentID,
		)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return enterpriseList, nil
}

// すべての企業情報を取得
func (repo *EnterpriseProfileRepositoryImpl) All() ([]*entity.EnterpriseProfile, error) {
	var (
		enterpriseList []*entity.EnterpriseProfile
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&enterpriseList, `
		SELECT 
			enterprise.*, staff.staff_name, staff.agent_id
		FROM 
			enterprise_profiles AS enterprise
		INNER JOIN
			agent_staffs AS staff
		ON
			enterprise.agent_staff_id = staff.id
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return enterpriseList, nil
}

/****************************************************************************************/
/// その他 API
//
// 郵便番号の下4桁チェック
func (repo *EnterpriseProfileRepositoryImpl) CheckPostCode(postCodeUnderNumber string) (*entity.EnterpriseProfile, error) {
	var (
		enterprise entity.EnterpriseProfile
	)

	err := repo.executer.Get(
		repo.Name+".CheckPostCode",
		&enterprise, `
		SELECT 
			enterprise.*
		FROM 
			enterprise_profiles AS enterprise
		WHERE
			SUBSTRING(enterprise.post_code, 5) = ?
		LIMIT 1
		`,
		postCodeUnderNumber,
	)

	if err != nil {
		return nil, err
	}

	return &enterprise, nil
}

/****************************************************************************************/
