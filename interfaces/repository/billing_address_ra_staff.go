package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type BillingAddressRAStaffRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewBillingAddressRAStaffRepositoryImpl(ex interfaces.SQLExecuter) usecase.BillingAddressRAStaffRepository {
	return &BillingAddressRAStaffRepositoryImpl{
		Name:     "BillingAddressRAStaffRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *BillingAddressRAStaffRepositoryImpl) Create(raStaff *entity.BillingAddressRAStaff) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO billing_address_ra_staffs (
				billing_address_id,
				billing_address_staff_name,
				billing_address_staff_email,
				billing_address_staff_phone_number,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)
		`,
		raStaff.BillingAddressID,
		raStaff.BillingAddressStaffName,
		raStaff.BillingAddressStaffEmail,
		raStaff.BillingAddressStaffPhoneNumber,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	raStaff.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除 API
//
// 請求先の人事担当者の削除
func (repo *BillingAddressRAStaffRepositoryImpl) DeleteByBillingAddressID(billingAddressID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM billing_address_ra_staffs
		WHERE billing_address_id = ?
		`, billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 複数取得 API
//
// 請求先IDから請求先の人事担当者を取得
func (repo *BillingAddressRAStaffRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.BillingAddressRAStaff, error) {
	var (
		raStaffList []*entity.BillingAddressRAStaff
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&raStaffList, `
		SELECT *
		FROM billing_address_ra_staffs
		WHERE
			billing_address_id = ?
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return raStaffList, err
	}

	return raStaffList, nil
}

// 求人企業IDから請求書一覧を取得
func (repo *BillingAddressRAStaffRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.BillingAddressRAStaff, error) {
	var (
		raStaffList []*entity.BillingAddressRAStaff
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&raStaffList, `
		SELECT *
		FROM 
		billing_address_ra_staffs
		WHERE
		billing_address_id IN (
			SELECT id
			FROM billing_addresses
			WHERE 
			enterprise_id = ?
			)
						`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return raStaffList, nil
}

// 全ての請求先の人事担当者を取得
func (repo *BillingAddressRAStaffRepositoryImpl) All() ([]*entity.BillingAddressRAStaff, error) {
	var (
		raStaffList []*entity.BillingAddressRAStaff
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&raStaffList, `
		SELECT *
		FROM billing_address_ra_staffs
		`,
	)

	if err != nil {
		fmt.Println(err)
		return raStaffList, err
	}

	return raStaffList, nil
}
