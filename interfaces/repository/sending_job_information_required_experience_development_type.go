package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredExperienceDevelopmentTypeRepository {
	return &SendingJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl{
		Name:     "SendingJobInformationRequiredExperienceDevelopmentTypeRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) Create(requiredExperienceDevelopmentType *entity.SendingJobInformationRequiredExperienceDevelopmentType) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_experience_development_types (
				experience_development_id,
				development_type,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		requiredExperienceDevelopmentType.ExperienceDevelopmentID,
		requiredExperienceDevelopmentType.DevelopmentType,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredExperienceDevelopmentType.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) Delete(jobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_required_experience_development_types
		WHERE experience_development_id IN (
			SELECT id
			FROM sending_job_information_required_experience_developments
			WHERE condition_id IN (
				SELECT id
				FROM sending_job_information_required_conditions
				WHERE sending_job_information_id = ?
			)
		)
		`, jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) GetListBySendingJobInformationID(jobInformationID uint) ([]*entity.SendingJobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.SendingJobInformationRequiredExperienceDevelopmentType
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&requiredExperienceDevelopmentTypeList, `
		SELECT *
		FROM sending_job_information_required_experience_development_types
		WHERE
			experience_development_id IN (
				SELECT id
				FROM sending_job_information_required_experience_developments
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
		return requiredExperienceDevelopmentTypeList, err
	}

	return requiredExperienceDevelopmentTypeList, nil
}

func (repo *SendingJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.SendingJobInformationRequiredExperienceDevelopmentType
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&requiredExperienceDevelopmentTypeList, `
			SELECT *
			FROM sending_job_information_required_experience_development_types
			WHERE
				experience_development_id IN (
					SELECT id
					FROM sending_job_information_required_experience_developments
					WHERE condition_id IN (
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

	return requiredExperienceDevelopmentTypeList, nil
}

func (repo *SendingJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.SendingJobInformationRequiredExperienceDevelopmentType
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&requiredExperienceDevelopmentTypeList, `
			SELECT *
			FROM sending_job_information_required_experience_development_types
			WHERE
				experience_development_id IN (
					SELECT id
					FROM sending_job_information_required_experience_developments
					WHERE condition_id IN (
						SELECT id
						FROM sending_job_information_required_conditions
						WHERE sending_job_information_id IN (
							SELECT id
							FROM sending_job_informations
							WHERE sending_billing_address_id IN (
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

	return requiredExperienceDevelopmentTypeList, nil
}

// 求人リストから必要開発経験を取得
func (repo *SendingJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.SendingJobInformationRequiredExperienceDevelopmentType
	)

	if len(idList) == 0 {
		return requiredExperienceDevelopmentTypeList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_experience_development_types.*
		FROM 
			sending_job_information_required_experience_development_types
		INNER JOIN
			sending_job_information_required_experience_developments AS experience_development
		ON
			experience_development.id = sending_job_information_required_experience_development_types.experience_development_id
		INNER JOIN
			sending_job_information_required_conditions AS required_condition
		ON
			required_condition.id = experience_development.condition_id
		WHERE
			required_condition.sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&requiredExperienceDevelopmentTypeList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentTypeList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.SendingJobInformationRequiredExperienceDevelopmentType
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredExperienceDevelopmentTypeList, `
			SELECT *
			FROM sending_job_information_required_experience_development_types
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentTypeList, nil
}
