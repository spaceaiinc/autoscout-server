package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingEnterpriseRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingEnterpriseRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingEnterpriseRepository {
	return &SendingEnterpriseRepositoryImpl{
		Name:     "SendingEnterpriseRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *SendingEnterpriseRepositoryImpl) Create(sendingEnterprise *entity.SendingEnterprise) error {

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO sending_enterprises (
			uuid,
			company_name,
			agent_staff_id,
			corporate_site_url,
			representative,
			establishment,
			post_code,
			office_location,
			employee_number_single,
			public_offering,
			sending_target,
			sending_required_info,
			remarks,
			is_deleted,
			password,
			created_at,
			updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?, ?
		)
		`,
		utility.CreateUUID(),
		sendingEnterprise.CompanyName,
		sendingEnterprise.AgentStaffID,
		sendingEnterprise.CorporateSiteURL,
		sendingEnterprise.Representative,
		sendingEnterprise.Establishment,
		sendingEnterprise.PostCode,
		sendingEnterprise.OfficeLocation,
		sendingEnterprise.EmployeeNumberSingle,
		sendingEnterprise.PublicOffering,
		sendingEnterprise.SendingTarget,
		sendingEnterprise.SendingRequiredInfo,
		sendingEnterprise.Remarks,
		false, // is_deleted
		sendingEnterprise.Password,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	sendingEnterprise.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
func (repo *SendingEnterpriseRepositoryImpl) Update(sendingEnterprise *entity.SendingEnterprise, sendingEnterpriseID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE sending_enterprises 
		SET
			company_name 						= ?,
			agent_staff_id 					= ?,
			corporate_site_url 			= ?,
			representative 					= ?,
			establishment 					= ?,
			post_code 							= ?,
			office_location 				= ?,
			employee_number_single 	= ?,
			public_offering 				= ?,
			sending_target          = ?,
			sending_required_info   = ?,
			remarks                 = ?,
			updated_at 							= ?
		WHERE 
			id = ?
		`,
		sendingEnterprise.CompanyName,
		sendingEnterprise.AgentStaffID,
		sendingEnterprise.CorporateSiteURL,
		sendingEnterprise.Representative,
		sendingEnterprise.Establishment,
		sendingEnterprise.PostCode,
		sendingEnterprise.OfficeLocation,
		sendingEnterprise.EmployeeNumberSingle,
		sendingEnterprise.PublicOffering,
		sendingEnterprise.SendingTarget,
		sendingEnterprise.SendingRequiredInfo,
		sendingEnterprise.Remarks,
		time.Now().In(time.UTC),
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingEnterpriseRepositoryImpl) UpdatePassword(sendingEnterpriseID uint, password string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdatePassword",
		`
		UPDATE sending_enterprises 
		SET
			password = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		password,
		time.Now().In(time.UTC),
		sendingEnterpriseID,
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
func (repo *SendingEnterpriseRepositoryImpl) Delete(sendingEnterpriseID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		SET
			is_deleted = true,
		FROM 
			sending_enterprises
		WHERE 
			id = ?
		`, sendingEnterpriseID,
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
func (repo *SendingEnterpriseRepositoryImpl) FindByID(sendingEnterpriseID uint) (*entity.SendingEnterprise, error) {
	var (
		sendingEnterprise entity.SendingEnterprise
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&sendingEnterprise, `
		SELECT 
			sending_enterprise.*, staff.staff_name, staff.agent_id
		FROM 
			sending_enterprises AS sending_enterprise
		INNER JOIN
			agent_staffs AS staff
		ON
			sending_enterprise.agent_staff_id = staff.id
		WHERE
			sending_enterprise.id = ?
		LIMIT 1
		`, sendingEnterpriseID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &sendingEnterprise, nil
}

func (repo *SendingEnterpriseRepositoryImpl) FindByUUID(uuid uuid.UUID) (*entity.SendingEnterprise, error) {
	var (
		sendingEnterprise entity.SendingEnterprise
	)

	err := repo.executer.Get(
		repo.Name+".FindByUUID",
		&sendingEnterprise, `
		SELECT 
			sending_enterprise.*
		FROM 
			sending_enterprises AS sending_enterprise
		WHERE
			sending_enterprise.uuid = ?
		LIMIT 1
		`, uuid)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &sendingEnterprise, nil
}

/****************************************************************************************/
/// 複数取得 API
//
func (repo *SendingEnterpriseRepositoryImpl) GetByIDList(enterpriseIDList []uint) ([]*entity.SendingEnterprise, error) {
	var (
		sendingEnterpriseList []*entity.SendingEnterprise
	)

	if len(enterpriseIDList) < 1 {
		return sendingEnterpriseList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_enterprise.*, staff.staff_name, staff.agent_id
		FROM 
			sending_enterprises AS sending_enterprise
		INNER JOIN
			agent_staffs AS staff
		ON
			sending_enterprise.agent_staff_id = staff.id
		WHERE
			sending_enterprise.id IN(%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(enterpriseIDList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByIDList",
		&sendingEnterpriseList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingEnterpriseList, nil
}

// 指定IDの求職者をまだ送客していない企業の一覧を取得
func (repo *SendingEnterpriseRepositoryImpl) GetSendingEnterpriseAndJobInfoCountByHaveNotSentYetBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingEnterprise, error) {
	var (
		sendingEnterpriseList []*entity.SendingEnterprise
	)

	err := repo.executer.Select(
		repo.Name+".GetSendingEnterpriseAndJobInfoCountByHaveNotSentYetBySendingJobSeekerID",
		&sendingEnterpriseList, `
		SELECT 
			enterprise.*, staff.staff_name, staff.agent_id, COUNT(job_info.id) AS job_information_count
		FROM 
			sending_enterprises AS enterprise
		INNER JOIN
			agent_staffs AS staff
		ON
			enterprise.agent_staff_id = staff.id
		LEFT JOIN 
			sending_billing_addresses AS billing
		ON 
			enterprise.id = billing.sending_enterprise_id
		LEFT JOIN 
			sending_job_informations AS job_info
		ON 
			billing.id = job_info.sending_billing_address_id
		LEFT JOIN 
			sending_phases AS phase
		ON 
			enterprise.id = phase.sending_enterprise_id AND
			? = phase.sending_job_seeker_id
		WHERE 
			phase.id IS NULL
		GROUP BY
			enterprise.id,
			enterprise.uuid,
			enterprise.company_name,
			enterprise.agent_staff_id,
			enterprise.corporate_site_url,
			enterprise.representative,
			enterprise.establishment,
			enterprise.post_code,
			enterprise.office_location,
			enterprise.employee_number_single,
			enterprise.public_offering,
			enterprise.is_deleted,
			staff.staff_name, staff.agent_id
		ORDER BY
			enterprise.id ASC
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingEnterpriseList, nil
}

// 指定IDの求職者をまだ送客していない企業の一覧を取得
func (repo *SendingEnterpriseRepositoryImpl) GetSendingEnterpriseByHaveNotSentYetBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingEnterprise, error) {
	var (
		sendingEnterpriseList []*entity.SendingEnterprise
	)

	err := repo.executer.Select(
		repo.Name+".GetSendingEnterpriseByHaveNotSentYetBySendingJobSeekerID",
		&sendingEnterpriseList, `
		SELECT 
			enterprise.*, staff.staff_name, staff.agent_id
		FROM 
			sending_enterprises AS enterprise
		INNER JOIN
			agent_staffs AS staff
		ON
			enterprise.agent_staff_id = staff.id
		LEFT JOIN 
			sending_phases AS phase
		ON 
			enterprise.id = phase.sending_enterprise_id AND
			? = phase.sending_job_seeker_id
		WHERE 
			phase.id IS NULL
		ORDER BY
			enterprise.id ASC
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingEnterpriseList, nil
}

// すべての企業情報を取得
func (repo *SendingEnterpriseRepositoryImpl) All() ([]*entity.SendingEnterprise, error) {
	var (
		sendingEnterpriseList []*entity.SendingEnterprise
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&sendingEnterpriseList, `
		SELECT 
			sending_enterprise.*, staff.staff_name, staff.agent_id
		FROM 
			sending_enterprises AS sending_enterprise
		INNER JOIN
			agent_staffs AS staff
		ON
			sending_enterprise.agent_staff_id = staff.id
		ORDER BY
			sending_enterprise.id DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingEnterpriseList, nil
}
