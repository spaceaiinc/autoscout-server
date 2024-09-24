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

type SendingBillingAddressRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingBillingAddressRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingBillingAddressRepository {
	return &SendingBillingAddressRepositoryImpl{
		Name:     "SendingBillingAddressRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
//求職者の作成
func (repo *SendingBillingAddressRepositoryImpl) Create(sendingBillingAddress *entity.SendingBillingAddress) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO sending_billing_addresses (
			uuid,
			sending_enterprise_id,
			agent_staff_id,
			contract_phase,
			contract_date,

			payment_policy,
			company_name,
			address,
			title,
			schedule_adjustment_url,

			commission,
			is_deleted,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, 
				?, ?, ?, ?
		)`,
		utility.CreateUUID(),
		sendingBillingAddress.SendingEnterpriseID,
		sendingBillingAddress.AgentStaffID,
		sendingBillingAddress.ContractPhase,
		sendingBillingAddress.ContractDate,
		sendingBillingAddress.PaymentPolicy,
		sendingBillingAddress.CompanyName,
		sendingBillingAddress.Address,
		sendingBillingAddress.Title,
		sendingBillingAddress.ScheduleAdjustmentURL,
		sendingBillingAddress.Commission,
		false, // is_deleted
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	sendingBillingAddress.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
func (repo *SendingBillingAddressRepositoryImpl) Update(sendingBillingAddress *entity.SendingBillingAddress, sendingBillingAddressID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE 
			sending_billing_addresses
		SET
			agent_staff_id = ?,
			contract_phase = ?,
			contract_date	= ?,
			payment_policy = ?,
			company_name = ?,
			address = ?,
			title = ?,
			schedule_adjustment_url = ?,
			commission = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		sendingBillingAddress.AgentStaffID,
		sendingBillingAddress.ContractPhase,
		sendingBillingAddress.ContractDate,
		sendingBillingAddress.PaymentPolicy,
		sendingBillingAddress.CompanyName,
		sendingBillingAddress.Address,
		sendingBillingAddress.Title,
		sendingBillingAddress.ScheduleAdjustmentURL,
		sendingBillingAddress.Commission,
		time.Now().In(time.UTC),
		sendingBillingAddressID,
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
func (repo *SendingBillingAddressRepositoryImpl) Delete(sendingBillingAddressID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		UPDATE 
			sending_billing_addresses
		SET
			is_deleted = true,
			updated_at = ?
		WHERE 
			id = ?
		`,
		time.Now().In(time.UTC),
		sendingBillingAddressID,
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
func (repo *SendingBillingAddressRepositoryImpl) FindByID(sendingBillingAddressID uint) (*entity.SendingBillingAddress, error) {
	var (
		sendingBillingAddress entity.SendingBillingAddress
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&sendingBillingAddress, `
		SELECT 
			billing.*, staff.staff_name, staff.agent_id
		FROM 
			sending_billing_addresses AS billing
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		WHERE
			billing.id = ?
		LIMIT 1
		`, sendingBillingAddressID)
	if err != nil {
		return nil, err
	}

	return &sendingBillingAddress, nil
}

// 求人企業IDから企業一覧を取得
func (repo *SendingBillingAddressRepositoryImpl) FindBySendingEnterpriseID(enterpriseID uint) (*entity.SendingBillingAddress, error) {
	var (
		sendingBillingAddress entity.SendingBillingAddress
	)

	err := repo.executer.Get(
		repo.Name+".FindBySendingEnterpriseID",
		&sendingBillingAddress,
		`
			SELECT 
				billing.*, staff.staff_name, staff.agent_id
			FROM 
				sending_billing_addresses AS billing
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			WHERE
				billing.sending_enterprise_id = ?
			LIMIT 1
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &sendingBillingAddress, nil
}

/****************************************************************************************/
/// 複数取得 API
//
func (repo *SendingBillingAddressRepositoryImpl) GetByIDList(billingAddressIDList []uint) ([]*entity.SendingBillingAddress, error) {
	var (
		sendingBillingAddressList []*entity.SendingBillingAddress
	)

	if len(billingAddressIDList) < 1 {
		return sendingBillingAddressList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			billing.*, staff.staff_name, staff.agent_id
		FROM 
			sending_billing_addresses AS billing
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		WHERE
			billing.id IN(%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(billingAddressIDList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByIDList",
		&sendingBillingAddressList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingBillingAddressList, nil
}

func (repo *SendingBillingAddressRepositoryImpl) GetBySendingEnterpriseIDList(enterpriseIDList []uint) ([]*entity.SendingBillingAddress, error) {
	var (
		sendingBillingAddressList []*entity.SendingBillingAddress
	)

	if len(enterpriseIDList) < 1 {
		return sendingBillingAddressList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			billing.*, staff.staff_name, staff.agent_id
		FROM 
			sending_billing_addresses AS billing
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		WHERE
			billing.sending_enterprise_id IN(%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(enterpriseIDList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetBySendingEnterpriseIDList",
		&sendingBillingAddressList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingBillingAddressList, nil
}

// エージェントIDから企業一覧を取得
func (repo *SendingBillingAddressRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingBillingAddress, error) {
	var (
		sendingBillingAddressList []*entity.SendingBillingAddress
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&sendingBillingAddressList, `
			SELECT 
				billing.*, staff.staff_name, staff.agent_id
			FROM 
				sending_billing_addresses AS billing
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			WHERE
				billing.sending_enterprise_id IN (
					SELECT id
					FROM sending_enterprises
					WHERE agent_staff_id IN (
						SELECT id
						FROM agent_staffs
						WHERE agent_id = ?
					)
				)
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingBillingAddressList, nil
}

// All
func (repo *SendingBillingAddressRepositoryImpl) All() ([]*entity.SendingBillingAddress, error) {
	var (
		sendingBillingAddressList []*entity.SendingBillingAddress
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&sendingBillingAddressList, `
		SELECT 
			billing.*, staff.staff_name, staff.agent_id
		FROM 
			sending_billing_addresses AS billing
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingBillingAddressList, nil
}

/****************************************************************************************/
