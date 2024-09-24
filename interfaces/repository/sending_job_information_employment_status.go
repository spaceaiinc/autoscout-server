package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationEmploymentStatusRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationEmploymentStatusRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationEmploymentStatusRepository {
	return &SendingJobInformationEmploymentStatusRepositoryImpl{
		Name:     "SendingJobInformationEmploymentStatusRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationEmploymentStatusRepositoryImpl) Create(employmentStatus *entity.SendingJobInformationEmploymentStatus) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_employment_statuses (
				sending_job_information_id,
				employment_status,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		employmentStatus.SendingJobInformationID,
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

func (repo *SendingJobInformationEmploymentStatusRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_employment_statuses
		WHERE	sending_job_information_id = ?
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationEmploymentStatusRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.SendingJobInformationEmploymentStatus
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&employmentStatusList, `
		SELECT *
		FROM sending_job_information_employment_statuses
		WHERE	sending_job_information_id = ?
		`,
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return employmentStatusList, err
	}

	return employmentStatusList, nil
}

func (repo *SendingJobInformationEmploymentStatusRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.SendingJobInformationEmploymentStatus
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&employmentStatusList, `
			SELECT *
			FROM sending_job_information_employment_statuses
			WHERE
				sending_job_information_id IN (
					SELECT id
					FROM sending_job_informations
					WHERE
					sending_billing_address_id = ?
				)
		`,
		sendingBillingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return employmentStatusList, nil
}

func (repo *SendingJobInformationEmploymentStatusRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.SendingJobInformationEmploymentStatus
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&employmentStatusList, `
			SELECT *
			FROM sending_job_information_employment_statuses
			WHERE
				sending_job_information_id IN (
					SELECT id
					FROM sending_job_informations
					WHERE 
						sending_billing_address_id IN (
							SELECT id
							FROM sending_billing_addresses
							WHERE sending_enterprise_id = ?
						) 
					)
		`,
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return employmentStatusList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationEmploymentStatusRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.SendingJobInformationEmploymentStatus
	)

	if len(idList) < 1 {
		return employmentStatusList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_employment_statuses.*
		FROM 
			sending_job_information_employment_statuses
		WHERE
			sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&employmentStatusList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return employmentStatusList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationEmploymentStatusRepositoryImpl) GetAll() ([]*entity.SendingJobInformationEmploymentStatus, error) {
	var (
		employmentStatusList []*entity.SendingJobInformationEmploymentStatus
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&employmentStatusList, `
			SELECT *
			FROM sending_job_information_employment_statuses
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return employmentStatusList, nil
}
