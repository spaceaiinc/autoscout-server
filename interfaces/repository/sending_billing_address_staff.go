package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingBillingAddressStaffRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingBillingAddressStaffRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingBillingAddressStaffRepository {
	return &SendingBillingAddressStaffRepositoryImpl{
		Name:     "SendingBillingAddressStaffRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *SendingBillingAddressStaffRepositoryImpl) Create(staff *entity.SendingBillingAddressStaff) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_billing_address_staffs (
				sending_billing_address_id,
				staff_name,
				staff_email,
				staff_phone_number,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)
		`,
		staff.SendingBillingAddressID,
		staff.StaffName,
		staff.StaffEmail,
		staff.StaffPhoneNumber,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	staff.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除 API
//
func (repo *SendingBillingAddressStaffRepositoryImpl) Delete(billingAddressID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_billing_address_staffs
		WHERE sending_billing_address_id = ?
		`, billingAddressID,
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
func (repo *SendingBillingAddressStaffRepositoryImpl) FindBySendingBillingAddressID(billingAddressID uint) ([]*entity.SendingBillingAddressStaff, error) {
	var (
		staffList []*entity.SendingBillingAddressStaff
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingBillingAddressID",
		&staffList, `
		SELECT *
		FROM sending_billing_address_staffs
		WHERE
			sending_billing_address_id = ?
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return staffList, err
	}

	return staffList, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// エージェントIDから請求先一覧を取得
func (repo *SendingBillingAddressStaffRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.SendingBillingAddressStaff, error) {
	var (
		staffList []*entity.SendingBillingAddressStaff
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&staffList, `
		SELECT *
		FROM 
		sending_billing_address_staffs
		WHERE
		sending_billing_address_id IN (
			SELECT id
			FROM sending_billing_addresses
			WHERE 
			enterprise_id IN (
				SELECT id
				FROM sending_enterprises
				WHERE agent_staff_id IN (
					SELECT id
					FROM agent_staffs
					WHERE agent_id = ?
				)
			)
		)
						`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return staffList, nil
}

func (repo *SendingBillingAddressStaffRepositoryImpl) GetByBillingAdressIDList(billingAddressIDList []uint) ([]*entity.SendingBillingAddressStaff, error) {
	var (
		staffList []*entity.SendingBillingAddressStaff
	)

	if len(billingAddressIDList) < 1 {
		return staffList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_billing_address_staffs
		WHERE
			sending_billing_address_id IN(%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(billingAddressIDList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAdressIDList",
		&staffList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return staffList, nil
}

// All
func (repo *SendingBillingAddressStaffRepositoryImpl) All() ([]*entity.SendingBillingAddressStaff, error) {
	var (
		staffList []*entity.SendingBillingAddressStaff
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&staffList, `
		SELECT *
		FROM sending_billing_address_staffs
		`,
	)

	if err != nil {
		fmt.Println(err)
		return staffList, err
	}

	return staffList, nil
}
