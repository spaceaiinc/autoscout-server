package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationRequiredExperienceDevelopmentTypeRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredExperienceDevelopmentTypeRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredExperienceDevelopmentTypeRepository {
	return &JobInformationRequiredExperienceDevelopmentTypeRepositoryImpl{
		Name:     "JobInformationRequiredExperienceDevelopmentTypeRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 必要開発経験を作成
func (repo *JobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) Create(requiredExperienceDevelopmentType *entity.JobInformationRequiredExperienceDevelopmentType) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_experience_development_types (
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
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredExperienceDevelopmentType.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 複数取得 API
//
// 指定求人IDの必要開発経験を取得
func (repo *JobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.JobInformationRequiredExperienceDevelopmentType
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&requiredExperienceDevelopmentTypeList, `
		SELECT *
		FROM job_information_required_experience_development_types
		WHERE
			experience_development_id IN (
				SELECT id
				FROM job_information_required_experience_developments
				WHERE condition_id IN (
					SELECT id
					FROM job_information_required_conditions
					WHERE job_information_id = ?
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

// 指定請求先IDの必要開発経験を取得
func (repo *JobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.JobInformationRequiredExperienceDevelopmentType
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredExperienceDevelopmentTypeList, `
			SELECT *
			FROM job_information_required_experience_development_types
			WHERE
				experience_development_id IN (
					SELECT id
					FROM job_information_required_experience_developments
					WHERE condition_id IN (
						SELECT id
						FROM job_information_required_conditions
						WHERE job_information_id IN (
							SELECT id
							FROM job_informations
							WHERE
							billing_address_id = ?
					)
				)
			)
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentTypeList, nil
}

// 指定企業IDの必要開発経験を取得
func (repo *JobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.JobInformationRequiredExperienceDevelopmentType
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredExperienceDevelopmentTypeList, `
			SELECT *
			FROM job_information_required_experience_development_types
			WHERE
				experience_development_id IN (
					SELECT id
					FROM job_information_required_experience_developments
					WHERE condition_id IN (
						SELECT id
						FROM job_information_required_conditions
						WHERE job_information_id IN (
							SELECT id
							FROM job_informations
							WHERE billing_address_id IN (
								SELECT id
								FROM billing_addresses
								WHERE enterprise_id = ?
							) 
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

	return requiredExperienceDevelopmentTypeList, nil
}

func (repo *JobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.JobInformationRequiredExperienceDevelopmentType
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&requiredExperienceDevelopmentTypeList, `
			SELECT *
			FROM job_information_required_experience_development_types
			WHERE
				experience_development_id IN (
					SELECT id
					FROM job_information_required_experience_developments
					WHERE condition_id IN (
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
				)
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentTypeList, nil
}

// 求人リストから必要開発経験を取得
func (repo *JobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.JobInformationRequiredExperienceDevelopmentType
	)

	if len(jobInformationIDList) == 0 {
		return requiredExperienceDevelopmentTypeList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_required_experience_development_types.*
		FROM 
			job_information_required_experience_development_types
		INNER JOIN
			job_information_required_experience_developments AS experience_development
		ON
			experience_development.id = job_information_required_experience_development_types.experience_development_id
		INNER JOIN
			job_information_required_conditions AS required_condition
		ON
			required_condition.id = experience_development.condition_id
		WHERE
			required_condition.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredExperienceDevelopmentTypeList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentTypeList, nil
}

func (repo *JobInformationRequiredExperienceDevelopmentTypeRepositoryImpl) All() ([]*entity.JobInformationRequiredExperienceDevelopmentType, error) {
	var (
		requiredExperienceDevelopmentTypeList []*entity.JobInformationRequiredExperienceDevelopmentType
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&requiredExperienceDevelopmentTypeList, `
			SELECT *
			FROM job_information_required_experience_development_types
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentTypeList, nil
}
