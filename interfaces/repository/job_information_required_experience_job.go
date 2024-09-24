package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationRequiredExperienceJobRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredExperienceJobRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredExperienceJobRepository {
	return &JobInformationRequiredExperienceJobRepositoryImpl{
		Name:     "JobInformationRequiredExperienceJobRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 必要経験業職種を作成
func (repo *JobInformationRequiredExperienceJobRepositoryImpl) Create(requiredExperienceJob *entity.JobInformationRequiredExperienceJob) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_experience_jobs (
				condition_id,
				experience_year,
				experience_month,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?,	?
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

/****************************************************************************************/
// 複数取得 API
//
// 指定求人IDの必要経験業職種を取得
func (repo *JobInformationRequiredExperienceJobRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.JobInformationRequiredExperienceJob
	)

	err := repo.executer.Select(
		repo.Name+".FindByJobInformationID",
		&requiredExperienceJobList, `
		SELECT *
		FROM job_information_required_experience_jobs
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
		return requiredExperienceJobList, err
	}

	return requiredExperienceJobList, nil
}

func (repo *JobInformationRequiredExperienceJobRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.JobInformationRequiredExperienceJob
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredExperienceJobList, `
			SELECT *
			FROM job_information_required_experience_jobs
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

	return requiredExperienceJobList, nil
}

func (repo *JobInformationRequiredExperienceJobRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.JobInformationRequiredExperienceJob
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredExperienceJobList, `
			SELECT *
			FROM job_information_required_experience_jobs
			WHERE
				condition_id IN (
					SELECT id
					FROM job_information_required_conditions
					WHERE job_information_id IN (
						SELECT id
						FROM job_informations
						WHERE 
							billing_address_id IN (
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

	return requiredExperienceJobList, nil
}

func (repo *JobInformationRequiredExperienceJobRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.JobInformationRequiredExperienceJob
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&requiredExperienceJobList, `
			SELECT *
			FROM job_information_required_experience_jobs
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

	return requiredExperienceJobList, nil
}

// 求人リストから必要経験業職種を取得
func (repo *JobInformationRequiredExperienceJobRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.JobInformationRequiredExperienceJob
	)

	if len(jobInformationIDList) == 0 {
		return requiredExperienceJobList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_required_experience_jobs.*
		FROM 
			job_information_required_experience_jobs
		INNER JOIN
			job_information_required_conditions AS required_condition
		ON
			required_condition.id = job_information_required_experience_jobs.condition_id
		WHERE
			required_condition.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredExperienceJobList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceJobList, nil
}

func (repo *JobInformationRequiredExperienceJobRepositoryImpl) All() ([]*entity.JobInformationRequiredExperienceJob, error) {
	var (
		requiredExperienceJobList []*entity.JobInformationRequiredExperienceJob
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&requiredExperienceJobList, `
			SELECT *
			FROM job_information_required_experience_jobs
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceJobList, nil
}
