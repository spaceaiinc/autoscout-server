package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredExperienceOccupationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredExperienceOccupationRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredExperienceOccupationRepository {
	return &SendingJobInformationRequiredExperienceOccupationRepositoryImpl{
		Name:     "SendingJobInformationRequiredExperienceOccupationRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredExperienceOccupationRepositoryImpl) Create(requiredExperienceOccupation *entity.SendingJobInformationRequiredExperienceOccupation) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_experience_occupations (
				experience_job_id,
				experience_occupation,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		requiredExperienceOccupation.ExperienceJobID,
		requiredExperienceOccupation.ExperienceOccupation,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredExperienceOccupation.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationRequiredExperienceOccupationRepositoryImpl) GetListBySendingJobInformationID(jobInformationID uint) ([]*entity.SendingJobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.SendingJobInformationRequiredExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&requiredExperienceOccupationList, `
		SELECT *
		FROM sending_job_information_required_experience_occupations
		WHERE
			experience_job_id IN (
				SELECT id
				FROM sending_job_information_required_experience_jobs
				WHERE condition_id IN (
					SELECT id
					FROM sending_job_information_required_conditions
					WHERE sending_job_information_id = ?
				)
			)
	`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return requiredExperienceOccupationList, err
	}

	return requiredExperienceOccupationList, nil
}

func (repo *SendingJobInformationRequiredExperienceOccupationRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.SendingJobInformationRequiredExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&requiredExperienceOccupationList, `
			SELECT *
			FROM sending_job_information_required_experience_occupations
			WHERE experience_job_id IN (
				SELECT id
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
				)
		`,
		sendingBillingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceOccupationList, nil
}

func (repo *SendingJobInformationRequiredExperienceOccupationRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.SendingJobInformationRequiredExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&requiredExperienceOccupationList, `
			SELECT *
			FROM sending_job_information_required_experience_occupations
			WHERE experience_job_id IN (
				SELECT id
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
				)
		`,
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceOccupationList, nil
}

// 求人リストから必要経験業職種を取得
func (repo *SendingJobInformationRequiredExperienceOccupationRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.SendingJobInformationRequiredExperienceOccupation
	)
	if len(idList) == 0 {
		return requiredExperienceOccupationList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_experience_occupations.*
		FROM 
			sending_job_information_required_experience_occupations
		INNER JOIN
			sending_job_information_required_experience_jobs AS experience_job
		ON
			experience_job.id = sending_job_information_required_experience_occupations.experience_job_id
		INNER JOIN
			sending_job_information_required_conditions AS required_condition
		ON
			required_condition.id = experience_job.condition_id
		WHERE 
			required_condition.sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&requiredExperienceOccupationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceOccupationList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRequiredExperienceOccupationRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.SendingJobInformationRequiredExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredExperienceOccupationList, `
			SELECT *
			FROM sending_job_information_required_experience_occupations
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceOccupationList, nil
}
