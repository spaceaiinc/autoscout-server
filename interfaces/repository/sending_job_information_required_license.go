package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredLicenseRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredLicenseRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredLicenseRepository {
	return &SendingJobInformationRequiredLicenseRepositoryImpl{
		Name:     "SendingJobInformationRequiredLicenseRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredLicenseRepositoryImpl) Create(requiredLicense *entity.SendingJobInformationRequiredLicense) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_licenses (
				condition_id,
				license,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		requiredLicense.ConditionID,
		requiredLicense.License,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredLicense.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationRequiredLicenseRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_required_licenses
		WHERE 
			condition_id IN (
				SELECT id
				FROM sending_job_information_required_conditions
				WHERE	sending_job_information_id = ?
			)
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationRequiredLicenseRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.SendingJobInformationRequiredLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&requiredLicenseList, `
		SELECT *
		FROM sending_job_information_required_licenses
		WHERE
			condition_id IN (
				SELECT id
				FROM sending_job_information_required_conditions
				WHERE	sending_job_information_id = ?
			)
		`,
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return requiredLicenseList, err
	}

	return requiredLicenseList, nil
}

func (repo *SendingJobInformationRequiredLicenseRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.SendingJobInformationRequiredLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&requiredLicenseList, `
			SELECT *
			FROM sending_job_information_required_licenses
			WHERE
				condition_id IN (
					SELECT id
					FROM sending_job_information_required_conditions
					WHERE sending_job_information_id IN (
						SELECT id
						FROM sending_job_informations
						WHERE
						sending_billing_address_id = ?
					)
				)
		`,
		sendingBillingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLicenseList, nil
}

func (repo *SendingJobInformationRequiredLicenseRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.SendingJobInformationRequiredLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&requiredLicenseList, `
			SELECT *
			FROM sending_job_information_required_licenses
			WHERE
				condition_id IN (
					SELECT id
					FROM sending_job_information_required_conditions
					WHERE sending_job_information_id IN (
						SELECT id
						FROM sending_job_informations
						WHERE 
							sending_billing_address_id IN (
								SELECT id
								FROM sending_billing_addresses
								WHERE sending_enterprise_id = ?
							) 
						)
					)
		`,
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLicenseList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationRequiredLicenseRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.SendingJobInformationRequiredLicense
	)

	if len(idList) == 0 {
		return requiredLicenseList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_licenses.*
		FROM 
			sending_job_information_required_licenses
		INNER JOIN
			sending_job_information_required_conditions AS required_condition
		ON
			required_condition.id = sending_job_information_required_licenses.condition_id
		WHERE
			required_condition.sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&requiredLicenseList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLicenseList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRequiredLicenseRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.SendingJobInformationRequiredLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredLicenseList, `
			SELECT *
			FROM sending_job_information_required_licenses
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLicenseList, nil
}
