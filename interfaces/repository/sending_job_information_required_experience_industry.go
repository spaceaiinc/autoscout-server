package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredExperienceIndustryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredExperienceIndustryRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredExperienceIndustryRepository {
	return &SendingJobInformationRequiredExperienceIndustryRepositoryImpl{
		Name:     "SendingJobInformationRequiredExperienceIndustryRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredExperienceIndustryRepositoryImpl) Create(requiredExperienceIndustry *entity.SendingJobInformationRequiredExperienceIndustry) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_experience_industries (
				experience_job_id,
				experience_industry,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		requiredExperienceIndustry.ExperienceJobID,
		requiredExperienceIndustry.ExperienceIndustry,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredExperienceIndustry.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationRequiredExperienceIndustryRepositoryImpl) GetListBySendingJobInformationID(jobInformationID uint) ([]*entity.SendingJobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.SendingJobInformationRequiredExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&requiredExperienceIndustryList, `
		SELECT *
		FROM sending_job_information_required_experience_industries
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
		return requiredExperienceIndustryList, err
	}

	return requiredExperienceIndustryList, nil
}

func (repo *SendingJobInformationRequiredExperienceIndustryRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.SendingJobInformationRequiredExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&requiredExperienceIndustryList, `
			SELECT *
			FROM sending_job_information_required_experience_industries
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

	return requiredExperienceIndustryList, nil
}

func (repo *SendingJobInformationRequiredExperienceIndustryRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.SendingJobInformationRequiredExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&requiredExperienceIndustryList, `
			SELECT *
			FROM sending_job_information_required_experience_industries
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

	return requiredExperienceIndustryList, nil
}

// 求人リストから必要経験業職種を取得
func (repo *SendingJobInformationRequiredExperienceIndustryRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.SendingJobInformationRequiredExperienceIndustry
	)

	if len(idList) == 0 {
		return requiredExperienceIndustryList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_experience_industries.*
		FROM 
			sending_job_information_required_experience_industries
		INNER JOIN
			sending_job_information_required_experience_jobs AS experience_job
		ON
			experience_job.id = sending_job_information_required_experience_industries.experience_job_id
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
		&requiredExperienceIndustryList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceIndustryList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRequiredExperienceIndustryRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.SendingJobInformationRequiredExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredExperienceIndustryList, `
			SELECT *
			FROM sending_job_information_required_experience_industries
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceIndustryList, nil
}
