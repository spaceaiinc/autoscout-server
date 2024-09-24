package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type JobInformationRequiredLicenseRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredLicenseRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredLicenseRepository {
	return &JobInformationRequiredLicenseRepositoryImpl{
		Name:     "JobInformationRequiredLicenseRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 必要資格を作成
func (repo *JobInformationRequiredLicenseRepositoryImpl) Create(requiredLicense *entity.JobInformationRequiredLicense) error {
	now := time.Now().In(time.UTC)
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_licenses (
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

/****************************************************************************************/
// 複数取得 API
//
func (repo *JobInformationRequiredLicenseRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.JobInformationRequiredLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&requiredLicenseList, `
		SELECT *
		FROM job_information_required_licenses
		WHERE
			condition_id IN (
				SELECT id
				FROM job_information_required_conditions
				WHERE	job_information_id = ?
			)
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return requiredLicenseList, err
	}

	return requiredLicenseList, nil
}

func (repo *JobInformationRequiredLicenseRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.JobInformationRequiredLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredLicenseList, `
			SELECT *
			FROM job_information_required_licenses
			WHERE
				condition_id IN (
					SELECT id
					FROM job_information_required_conditions
					WHERE job_information_id IN (
						SELECT id
						FROM job_informations
						WHERE
						billing_address_id = ?
					)
				)
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLicenseList, nil
}

func (repo *JobInformationRequiredLicenseRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.JobInformationRequiredLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredLicenseList, `
			SELECT *
			FROM job_information_required_licenses
			WHERE
				condition_id IN (
					SELECT id
					FROM job_information_required_conditions
					WHERE job_information_id IN (
						SELECT id
						FROM job_informations
						WHERE 
							billing_address_id IN (
								SELECT id
								FROM billing_addresses
								WHERE enterprise_id = ?
							) 
						)
					)
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLicenseList, nil
}

func (repo *JobInformationRequiredLicenseRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.JobInformationRequiredLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&requiredLicenseList, `
			SELECT *
			FROM job_information_required_licenses
			WHERE
				condition_id IN (
					SELECT id
					FROM job_information_required_conditions
					WHERE job_information_id IN (
						SELECT id
						FROM job_informations
						WHERE billing_address_id IN (
							SELECT id
							FROM billing_addresses
							WHERE enterprise_id IN (
								SELECT id
								FROM enterprise_profiles
								WHERE agent_staff_id IN (
									SELECT id
									FROM agent_staffs
									WHERE agent_id = ?
								)
							)
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

	return requiredLicenseList, nil
}

// 求人リストから必要資格を取得
func (repo *JobInformationRequiredLicenseRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.JobInformationRequiredLicense
	)

	if len(jobInformationIDList) == 0 {
		return requiredLicenseList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_required_licenses.*
		FROM 
			job_information_required_licenses
		INNER JOIN
			job_information_required_conditions AS required_condition
		ON
			required_condition.id = job_information_required_licenses.condition_id
		WHERE
			required_condition.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredLicenseList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLicenseList, nil
}

// 求人リストと資格タイプから必要資格を取得
func (repo *JobInformationRequiredLicenseRepositoryImpl) GetByJobInformationIDListAndLicenceTypeList(jobInformationIDList []uint, licenceTypeList []null.Int) ([]*entity.JobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.JobInformationRequiredLicense
		licenseTypeListStr  string
		lenLicenseTypeList  = len(licenceTypeList)
	)

	if len(jobInformationIDList) == 0 || lenLicenseTypeList == 0 {
		return requiredLicenseList, nil
	}

	for i, licenceType := range licenceTypeList {
		if licenceType.Valid {
			licenseTypeListStr += fmt.Sprintf("%d", licenceType.Int64)
			if i < lenLicenseTypeList-1 {
				licenseTypeListStr += ", "
			}
		}
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_required_licenses.*
		FROM 
			job_information_required_licenses
		INNER JOIN
			job_information_required_conditions AS required_condition
		ON
			required_condition.id = job_information_required_licenses.condition_id
		WHERE
			required_condition.job_information_id IN (%s)
		AND
			job_information_required_licenses.license IN (%s)
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"),
		licenseTypeListStr,
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDListAndLicenceTypeList",
		&requiredLicenseList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLicenseList, nil
}

func (repo *JobInformationRequiredLicenseRepositoryImpl) All() ([]*entity.JobInformationRequiredLicense, error) {
	var (
		requiredLicenseList []*entity.JobInformationRequiredLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredLicenseList, `
			SELECT *
			FROM job_information_required_licenses
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLicenseList, nil
}
