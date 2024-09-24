package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationEmploymentStatusRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationEmploymentStatusRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationEmploymentStatusRepository {
	return &JobInformationEmploymentStatusRepositoryImpl{
		Name:     "JobInformationEmploymentStatusRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 雇用形態を作成する
func (repo *JobInformationEmploymentStatusRepositoryImpl) Create(employmentStatus *entity.JobInformationEmploymentStatus) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_employment_statuses (
				job_information_id,
				employment_status,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		employmentStatus.JobInformationID,
		employmentStatus.EmploymentStatus,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	employmentStatus.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 削除 API
//
// 雇用形態を削除する
func (repo *JobInformationEmploymentStatusRepositoryImpl) DeleteByJobInformationID(jobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobInformationID",
		`
		DELETE
		FROM job_information_employment_statuses
		WHERE job_information_id = ?
		`, jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
// 複数取得 API
//
// 指定求人IDの雇用形態を取得する
func (repo *JobInformationEmploymentStatusRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.JobInformationEmploymentStatus
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&employmentStatusList, `
		SELECT *
		FROM job_information_employment_statuses
		WHERE
			job_information_id = ?
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return employmentStatusList, err
	}

	return employmentStatusList, nil
}

// 指定請求先IDの雇用形態を取得する
func (repo *JobInformationEmploymentStatusRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.JobInformationEmploymentStatus
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&employmentStatusList, `
			SELECT *
			FROM job_information_employment_statuses
			WHERE
				job_information_id IN (
					SELECT id
					FROM job_informations
					WHERE
					billing_address_id = ?
				)
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return employmentStatusList, nil
}

// 指定企業IDの雇用形態を取得する
func (repo *JobInformationEmploymentStatusRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.JobInformationEmploymentStatus
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&employmentStatusList, `
			SELECT *
			FROM job_information_employment_statuses
			WHERE
				job_information_id IN (
					SELECT id
					FROM job_informations
					WHERE 
						billing_address_id IN (
							SELECT id
							FROM billing_addresses
							WHERE enterprise_id = ?
						) 
					)
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return employmentStatusList, nil
}

func (repo *JobInformationEmploymentStatusRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.JobInformationEmploymentStatus
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&employmentStatusList, `
			SELECT *
			FROM job_information_employment_statuses
			WHERE
				job_information_id IN (
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
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return employmentStatusList, nil
}

// 求人リストから雇用形態を取得
func (repo *JobInformationEmploymentStatusRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.JobInformationEmploymentStatus
	)

	if len(jobInformationIDList) == 0 {
		return employmentStatusList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_information_employment_statuses
		WHERE
			job_information_id IN (%s)
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&employmentStatusList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return employmentStatusList, nil
}

func (repo *JobInformationEmploymentStatusRepositoryImpl) All() ([]*entity.JobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.JobInformationEmploymentStatus
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&employmentStatusList, `
			SELECT *
			FROM job_information_employment_statuses
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return employmentStatusList, nil
}
