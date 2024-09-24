package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredLanguageRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredLanguageRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredLanguageRepository {
	return &SendingJobInformationRequiredLanguageRepositoryImpl{
		Name:     "SendingJobInformationRequiredLanguageRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredLanguageRepositoryImpl) Create(requiredLanguage *entity.SendingJobInformationRequiredLanguage) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_languages (
				condition_id,
				language_level,
				toeic,
				toefl_ibt,
				toefl_pbt,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?
			)
		`,
		requiredLanguage.ConditionID,
		requiredLanguage.LanguageLevel,
		requiredLanguage.Toeic,
		requiredLanguage.ToeflIBT,
		requiredLanguage.ToeflPBT,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredLanguage.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationRequiredLanguageRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_required_languages
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

func (repo *SendingJobInformationRequiredLanguageRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.SendingJobInformationRequiredLanguage
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&requiredLanguageList, `
		SELECT *
		FROM sending_job_information_required_languages
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
		return requiredLanguageList, err
	}

	return requiredLanguageList, nil
}

func (repo *SendingJobInformationRequiredLanguageRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.SendingJobInformationRequiredLanguage
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&requiredLanguageList, `
			SELECT *
			FROM sending_job_information_required_languages
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

	return requiredLanguageList, nil
}

func (repo *SendingJobInformationRequiredLanguageRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.SendingJobInformationRequiredLanguage
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&requiredLanguageList, `
			SELECT *
			FROM sending_job_information_required_languages
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

	return requiredLanguageList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationRequiredLanguageRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.SendingJobInformationRequiredLanguage
	)

	if len(idList) == 0 {
		return requiredLanguageList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_languages.*
		FROM 
			sending_job_information_required_languages
		INNER JOIN
			sending_job_information_required_conditions AS required_condition
		ON
			required_condition.id = sending_job_information_required_languages.condition_id
		WHERE
			required_condition.sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&requiredLanguageList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLanguageList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRequiredLanguageRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.SendingJobInformationRequiredLanguage
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredLanguageList, `
			SELECT *
			FROM sending_job_information_required_languages
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLanguageList, nil
}
