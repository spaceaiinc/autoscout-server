package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredExperienceJobRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredExperienceJobRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredExperienceJobRepository {
	return &SendingJobInformationRequiredExperienceJobRepositoryImpl{
		Name:     "SendingJobInformationRequiredExperienceJobRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredExperienceJobRepositoryImpl) Create(requiredExperienceJob *entity.SendingJobInformationRequiredExperienceJob) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_experience_jobs (
				condition_id,
				experience_year,
				experience_month,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		requiredExperienceJob.ConditionID,
		requiredExperienceJob.ExperienceYear,
		requiredExperienceJob.ExperienceMonth,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredExperienceJob.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationRequiredExperienceJobRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_required_experience_jobs
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

func (repo *SendingJobInformationRequiredExperienceJobRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.SendingJobInformationRequiredExperienceJob
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&requiredExperienceJobList, `
		SELECT *
		FROM sending_job_information_required_experience_jobs
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
		return requiredExperienceJobList, err
	}

	return requiredExperienceJobList, nil
}

func (repo *SendingJobInformationRequiredExperienceJobRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.SendingJobInformationRequiredExperienceJob
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&requiredExperienceJobList, `
			SELECT *
			FROM sending_job_information_required_experience_jobs
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

	return requiredExperienceJobList, nil
}

func (repo *SendingJobInformationRequiredExperienceJobRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.SendingJobInformationRequiredExperienceJob
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&requiredExperienceJobList, `
			SELECT *
			FROM sending_job_information_required_experience_jobs
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

	return requiredExperienceJobList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationRequiredExperienceJobRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.SendingJobInformationRequiredExperienceJob
	)

	if len(idList) == 0 {
		return requiredExperienceJobList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_experience_jobs.*
		FROM 
			sending_job_information_required_experience_jobs
		INNER JOIN
			sending_job_information_required_conditions AS required_condition
		ON
			required_condition.id = sending_job_information_required_experience_jobs.condition_id
		WHERE
			required_condition.sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&requiredExperienceJobList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceJobList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRequiredExperienceJobRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.SendingJobInformationRequiredExperienceJob
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredExperienceJobList, `
			SELECT *
			FROM sending_job_information_required_experience_jobs
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceJobList, nil
}
