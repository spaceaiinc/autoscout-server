package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredExperienceDevelopmentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredExperienceDevelopmentRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredExperienceDevelopmentRepository {
	return &SendingJobInformationRequiredExperienceDevelopmentRepositoryImpl{
		Name:     "SendingJobInformationRequiredExperienceDevelopmentRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredExperienceDevelopmentRepositoryImpl) Create(requiredExperienceDevelopment *entity.SendingJobInformationRequiredExperienceDevelopment) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_experience_developments (
				condition_id,
				development_category,
				experience_year,
				experience_month,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)
		`,
		requiredExperienceDevelopment.ConditionID,
		requiredExperienceDevelopment.DevelopmentCategory,
		requiredExperienceDevelopment.ExperienceYear,
		requiredExperienceDevelopment.ExperienceMonth,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredExperienceDevelopment.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationRequiredExperienceDevelopmentRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_required_experience_developments
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

func (repo *SendingJobInformationRequiredExperienceDevelopmentRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.SendingJobInformationRequiredExperienceDevelopment
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&requiredExperienceDevelopmentList, `
		SELECT *
		FROM sending_job_information_required_experience_developments
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
		return requiredExperienceDevelopmentList, err
	}

	return requiredExperienceDevelopmentList, nil
}

func (repo *SendingJobInformationRequiredExperienceDevelopmentRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.SendingJobInformationRequiredExperienceDevelopment
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&requiredExperienceDevelopmentList, `
			SELECT *
			FROM sending_job_information_required_experience_developments
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

	return requiredExperienceDevelopmentList, nil
}

func (repo *SendingJobInformationRequiredExperienceDevelopmentRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.SendingJobInformationRequiredExperienceDevelopment
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&requiredExperienceDevelopmentList, `
			SELECT *
			FROM sending_job_information_required_experience_developments
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

	return requiredExperienceDevelopmentList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationRequiredExperienceDevelopmentRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.SendingJobInformationRequiredExperienceDevelopment
	)

	if len(idList) == 0 {
		return requiredExperienceDevelopmentList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_experience_developments.*
		FROM 
			sending_job_information_required_experience_developments
		INNER JOIN
			sending_job_information_required_conditions AS required_condition
		ON
			required_condition.id = sending_job_information_required_experience_developments.condition_id
		WHERE
			required_condition.sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&requiredExperienceDevelopmentList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRequiredExperienceDevelopmentRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.SendingJobInformationRequiredExperienceDevelopment
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredExperienceDevelopmentList, `
			SELECT *
			FROM sending_job_information_required_experience_developments
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentList, nil
}
