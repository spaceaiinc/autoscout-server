package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type BillingAddressRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewBillingAddressRepositoryImpl(ex interfaces.SQLExecuter) usecase.BillingAddressRepository {
	return &BillingAddressRepositoryImpl{
		Name:     "BillingAddressRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// 請求先の作成
func (repo *BillingAddressRepositoryImpl) Create(billingAddress *entity.BillingAddress) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO billing_addresses (
			uuid,
			enterprise_id,
			agent_staff_id,
			contract_phase,
			contract_date,
			payment_policy,
			company_name,
			address,
			how_to_recommend,
			title,
			created_at,
			updated_at,
			is_deleted
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
		utility.CreateUUID(),
		billingAddress.EnterpriseID,
		billingAddress.AgentStaffID,
		billingAddress.ContractPhase,
		billingAddress.ContractDate,
		billingAddress.PaymentPolicy,
		billingAddress.CompanyName,
		billingAddress.Address,
		billingAddress.HowToRecommend,
		billingAddress.Title,
		now,
		now,
		false,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	billingAddress.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
// 請求先を更新
func (repo *BillingAddressRepositoryImpl) Update(id uint, billingAddress *entity.BillingAddress) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE 
			billing_addresses
		SET
			agent_staff_id = ?,
			contract_phase = ?,
			contract_date	= ?,
			payment_policy = ?,
			company_name = ?,
			address = ?,
			how_to_recommend = ?,
			title = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		billingAddress.AgentStaffID,
		billingAddress.ContractPhase,
		billingAddress.ContractDate,
		billingAddress.PaymentPolicy,
		billingAddress.CompanyName,
		billingAddress.Address,
		billingAddress.HowToRecommend,
		billingAddress.Title,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// agent_staff_idを更新する。引き継ぎ用
func (repo *BillingAddressRepositoryImpl) UpdateAgentStaffID(id, agentStaffID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE 
			billing_addresses
		SET
			agent_staff_id = ?,
			updated_at = ?
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

// 指定の請求先idリストのagent_staff_idを更新
func (repo *BillingAddressRepositoryImpl) UpdateAgentStaffIDByBillingAddressIDList(idList []uint, agentStaffID uint) error {
	if len(idList) < 1 {
		return nil
	}

	query := fmt.Sprintf(`
	UPDATE 
		billing_addresses
	SET
		agent_staff_id = ?,
		updated_at = ?
	WHERE 
		id IN (%s)
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgentStaffIDByBillingAddressIDList",
		query,
		agentStaffID,
		time.Now().In(time.UTC),
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
// 請求先を削除
func (repo *BillingAddressRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		UPDATE 
			billing_addresses
		SET
			is_deleted = true,
			updated_at = ?
		WHERE 
			id = ?
		`,
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
/// 単数取得 API
//
// 指定IDの請求先を取得
func (repo *BillingAddressRepositoryImpl) FindByID(id uint) (*entity.BillingAddress, error) {
	var (
		billingAddress entity.BillingAddress
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&billingAddress, `
		SELECT 
			billing.*, staff.staff_name, staff.agent_id
		FROM 
			billing_addresses AS billing
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		WHERE
			billing.id = ?
		AND
			billing.is_deleted = false
		LIMIT 1
		`, id)
	if err != nil {
		return nil, err
	}

	return &billingAddress, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// 求人企業IDから企業一覧を取得
func (repo *BillingAddressRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.BillingAddress, error) {
	var (
		billingAddressList []*entity.BillingAddress
	)

	err := repo.executer.Select(
		repo.Name+".GetBillingAddressByEnterpriseID",
		&billingAddressList, `
			SELECT 
				billing.*, staff.staff_name, staff.agent_id
			FROM 
				billing_addresses AS billing
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			WHERE
				billing.enterprise_id = ?
			AND
				billing.is_deleted = false
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return billingAddressList, nil
}

// エージェント担当者IDから企業一覧を取得 *担当引き継ぎ用
func (repo *BillingAddressRepositoryImpl) GetByAgentStaffID(agentStaffID uint) ([]*entity.BillingAddress, error) {
	var (
		billingAddressList []*entity.BillingAddress
	)

	err := repo.executer.Select(
		repo.Name+".GetBillingAddressByAgentID",
		&billingAddressList, `
			SELECT 
				billing.id
			FROM 
				billing_addresses AS billing
			WHERE
				billing.agent_staff_id = ?
			AND
				billing.is_deleted = false
		`,
		agentStaffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return billingAddressList, nil
}

// エージェントIDから企業一覧を取得
func (repo *BillingAddressRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.BillingAddress, error) {
	var (
		billingAddressList []*entity.BillingAddress
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&billingAddressList, `
			SELECT 
				billing.*, 
				staff.staff_name, staff.agent_id,
				enterprise.company_name as enterprise_company_name
			FROM 
				billing_addresses AS billing
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			WHERE
				billing.enterprise_id IN (
					SELECT id
					FROM enterprise_profiles
					WHERE agent_staff_id IN (
						SELECT id
						FROM agent_staffs
						WHERE agent_id = ?
					)
				)
			AND
				billing.is_deleted = false
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return billingAddressList, nil
}

// エージェントIDから企業一覧を取得
func (repo *BillingAddressRepositoryImpl) GetByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.BillingAddress, error) {
	var (
		billingAddressList []*entity.BillingAddress
		freeWordQuery      string
	)

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {
			freeWordForLike := "%" + freeWord + "%"

			freeWordQuery = fmt.Sprintf(`AND (billing_address.company_name LIKE '%s')`, freeWordForLike)
		}
	} else {
		freeWordQuery = fmt.Sprintf(`AND billing_address.id = %v`, freeWordInt)
	}

	query := fmt.Sprintf(`
		SELECT
			billing_address.* , staff.staff_name, staff.agent_id
		FROM
			billing_addresses AS billing_address
		INNER JOIN
			agent_staffs AS staff
		ON
			billing_address.agent_staff_id = staff.id
		WHERE
			staff.agent_id = %v
		%s
			ORDER BY id DESC
	`, agentID, freeWordQuery)

	// フリーワードがある場合
	err = repo.executer.Select(
		repo.Name+".GetByAgentIDAndFreeWord",
		&billingAddressList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return billingAddressList, nil
}

// GetAll
func (repo *BillingAddressRepositoryImpl) All() ([]*entity.BillingAddress, error) {
	var (
		billingAddressList []*entity.BillingAddress
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&billingAddressList, `
		SELECT 
			billing.*, staff.staff_name, staff.agent_id
		FROM 
			billing_addresses AS billing
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
			WHERE billing.is_deleted = false
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return billingAddressList, nil
}

/****************************************************************************************/
