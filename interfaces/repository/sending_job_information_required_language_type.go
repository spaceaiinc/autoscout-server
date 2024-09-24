package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredLanguageTypeRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredLanguageTypeRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredLanguageTypeRepository {
	return &SendingJobInformationRequiredLanguageTypeRepositoryImpl{
		Name:     "SendingJobInformationRequiredLanguageTypeRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredLanguageTypeRepositoryImpl) Create(requiredLanguageType *entity.SendingJobInformationRequiredLanguageType) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_language_types (
				language_id,
				language_type,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		requiredLanguageType.LanguageID,
		requiredLanguageType.LanguageType,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredLanguageType.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationRequiredLanguageTypeRepositoryImpl) Delete(jobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_required_language_types
		WHERE 
			language_id IN (
				SELECT id
				FROM sending_job_information_required_languages
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

func (repo *SendingJobInformationRequiredLanguageTypeRepositoryImpl) GetListBySendingJobInformationID(jobInformationID uint) ([]*entity.SendingJobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.SendingJobInformationRequiredLanguageType
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&requiredLanguageTypeList, `
		SELECT *
		FROM sending_job_information_required_language_types
		WHERE
			language_id IN (
				SELECT id
				FROM sending_job_information_required_languages
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
		return requiredLanguageTypeList, err
	}

	return requiredLanguageTypeList, nil
}

func (repo *SendingJobInformationRequiredLanguageTypeRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.SendingJobInformationRequiredLanguageType
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&requiredLanguageTypeList, `
			SELECT *
			FROM sending_job_information_required_language_types
			WHERE
				language_id IN (
					SELECT id
					FROM sending_job_information_required_languages
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

	return requiredLanguageTypeList, nil
}

func (repo *SendingJobInformationRequiredLanguageTypeRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.SendingJobInformationRequiredLanguageType
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&requiredLanguageTypeList, `
			SELECT *
			FROM sending_job_information_required_language_types
			WHERE
				language_id IN (
					SELECT id
					FROM sending_job_information_required_languages
					WHERE condition_id IN (
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

	return requiredLanguageTypeList, nil
}

// 求人リストから必要言語を取得
func (repo *SendingJobInformationRequiredLanguageTypeRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.SendingJobInformationRequiredLanguageType
	)

	if len(idList) == 0 {
		return requiredLanguageTypeList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_language_types.*
		FROM 
			sending_job_information_required_language_types
		INNER JOIN
			sending_job_information_required_languages AS language
		ON
			language.id = sending_job_information_required_language_types.language_id
		INNER JOIN
			sending_job_information_required_conditions AS required_condition
		ON
			required_condition.id = language.condition_id
		WHERE
			required_condition.sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&requiredLanguageTypeList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLanguageTypeList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRequiredLanguageTypeRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.SendingJobInformationRequiredLanguageType
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredLanguageTypeList, `
			SELECT *
			FROM sending_job_information_required_language_types
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLanguageTypeList, nil
}
