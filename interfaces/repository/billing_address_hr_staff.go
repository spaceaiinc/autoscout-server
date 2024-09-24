package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type BillingAddressHRStaffRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewBillingAddressHRStaffRepositoryImpl(ex interfaces.SQLExecuter) usecase.BillingAddressHRStaffRepository {
	return &BillingAddressHRStaffRepositoryImpl{
		Name:     "BillingAddressHRStaffRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// 請求先の人事担当者の作成
func (repo *BillingAddressHRStaffRepositoryImpl) Create(hrStaff *entity.BillingAddressHRStaff) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO billing_address_hr_staffs (
				billing_address_id,
				hr_staff_name,
				hr_staff_email,
				hr_staff_phone_number,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)
		`,
		hrStaff.BillingAddressID,
		hrStaff.HRStaffName,
		hrStaff.HRStaffEmail,
		hrStaff.HRStaffPhoneNumber,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	hrStaff.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除 API
//
// 請求先の人事担当者の削除
func (repo *BillingAddressHRStaffRepositoryImpl) DeleteByBillingAddressID(billingAddressID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByBillingAddressID",
		`
		DELETE
		FROM billing_address_hr_staffs
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
// 請求先IDから人事担当者一覧を取得
func (repo *BillingAddressHRStaffRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.BillingAddressHRStaff, error) {
	var (
		hrStaffList []*entity.BillingAddressHRStaff
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&hrStaffList, `
		SELECT *
		FROM billing_address_hr_staffs
		WHERE
			billing_address_id = ?
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return hrStaffList, err
	}

	return hrStaffList, nil
}

// 求人企業IDから請求先一覧を取得
func (repo *BillingAddressHRStaffRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.BillingAddressHRStaff, error) {
	var (
		hrStaffList []*entity.BillingAddressHRStaff
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&hrStaffList, `
		SELECT *
		FROM 
		billing_address_hr_staffs
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

	return hrStaffList, nil
}

// 全てのレコードを取得
func (repo *BillingAddressHRStaffRepositoryImpl) All() ([]*entity.BillingAddressHRStaff, error) {
	var (
		hrStaffList []*entity.BillingAddressHRStaff
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&hrStaffList, `
		SELECT *
		FROM billing_address_hr_staffs
		`,
	)

	if err != nil {
		fmt.Println(err)
		return hrStaffList, err
	}

	return hrStaffList, nil
}
