package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationRequiredExperienceDevelopmentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredExperienceDevelopmentRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredExperienceDevelopmentRepository {
	return &JobInformationRequiredExperienceDevelopmentRepositoryImpl{
		Name:     "JobInformationRequiredExperienceDevelopmentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 必要開発経験を作成
func (repo *JobInformationRequiredExperienceDevelopmentRepositoryImpl) Create(requiredExperienceDevelopment *entity.JobInformationRequiredExperienceDevelopment) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_experience_developments (
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

/****************************************************************************************/
// 複数取得 API
//
// 指定求人IDの必要開発経験を取得
func (repo *JobInformationRequiredExperienceDevelopmentRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.JobInformationRequiredExperienceDevelopment
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&requiredExperienceDevelopmentList, `
		SELECT *
		FROM job_information_required_experience_developments
		WHERE
			condition_id IN (
				SELECT id
				FROM job_information_required_conditions
				WHERE job_information_id = ?
			)
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return requiredExperienceDevelopmentList, err
	}

	return requiredExperienceDevelopmentList, nil
}

// 指定請求先IDの必要開発経験を取得
func (repo *JobInformationRequiredExperienceDevelopmentRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.JobInformationRequiredExperienceDevelopment
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredExperienceDevelopmentList, `
			SELECT *
			FROM job_information_required_experience_developments
			WHERE
				condition_id IN (
					SELECT id
					FROM job_information_required_conditions
					WHERE job_information_id IN (
						SELECT id
						FROM job_informations
						WHERE
						billing_address_id = ?
				)
			)
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentList, nil
}

// 指定企業IDの必要開発経験を取得
func (repo *JobInformationRequiredExperienceDevelopmentRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.JobInformationRequiredExperienceDevelopment
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredExperienceDevelopmentList, `
			SELECT *
			FROM job_information_required_experience_developments
			WHERE
				condition_id IN (
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
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentList, nil
}

func (repo *JobInformationRequiredExperienceDevelopmentRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.JobInformationRequiredExperienceDevelopment
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&requiredExperienceDevelopmentList, `
			SELECT *
			FROM job_information_required_experience_developments
			WHERE
				condition_id IN (
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
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentList, nil
}

// 求人リストから必要開発経験を取得
func (repo *JobInformationRequiredExperienceDevelopmentRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.JobInformationRequiredExperienceDevelopment
	)

	if len(jobInformationIDList) == 0 {
		return requiredExperienceDevelopmentList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_required_experience_developments.*
		FROM 	
			job_information_required_experience_developments
		INNER JOIN
			job_information_required_conditions AS required_condition
		ON
			required_condition.id = job_information_required_experience_developments.condition_id
		WHERE
			required_condition.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredExperienceDevelopmentList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentList, nil
}

func (repo *JobInformationRequiredExperienceDevelopmentRepositoryImpl) All() ([]*entity.JobInformationRequiredExperienceDevelopment, error) {
	var (
		requiredExperienceDevelopmentList []*entity.JobInformationRequiredExperienceDevelopment
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&requiredExperienceDevelopmentList, `
			SELECT *
			FROM job_information_required_experience_developments
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceDevelopmentList, nil
}
