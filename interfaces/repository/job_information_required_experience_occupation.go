package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationRequiredExperienceOccupationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredExperienceOccupationRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredExperienceOccupationRepository {
	return &JobInformationRequiredExperienceOccupationRepositoryImpl{
		Name:     "JobInformationRequiredExperienceOccupationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 必要職種経験を作成
func (repo *JobInformationRequiredExperienceOccupationRepositoryImpl) Create(requiredExperienceOccupation *entity.JobInformationRequiredExperienceOccupation) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_experience_occupations (
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

func (repo *JobInformationRequiredExperienceOccupationRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.JobInformationRequiredExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&requiredExperienceOccupationList, `
		SELECT *
		FROM job_information_required_experience_occupations
		WHERE
			experience_job_id IN (
				SELECT id
				FROM job_information_required_experience_jobs
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
		return requiredExperienceOccupationList, err
	}

	return requiredExperienceOccupationList, nil
}

func (repo *JobInformationRequiredExperienceOccupationRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.JobInformationRequiredExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredExperienceOccupationList, `
			SELECT *
			FROM job_information_required_experience_occupations
			WHERE experience_job_id IN (
				SELECT id
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
				)
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceOccupationList, nil
}

func (repo *JobInformationRequiredExperienceOccupationRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.JobInformationRequiredExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredExperienceOccupationList, `
			SELECT *
			FROM job_information_required_experience_occupations
			WHERE experience_job_id IN (
				SELECT id
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
				)
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceOccupationList, nil
}

func (repo *JobInformationRequiredExperienceOccupationRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.JobInformationRequiredExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&requiredExperienceOccupationList, `
			SELECT *
			FROM job_information_required_experience_occupations
			WHERE experience_job_id IN (
				SELECT id
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

	return requiredExperienceOccupationList, nil
}

// 求人リストから必要経験業職種を取得
func (repo *JobInformationRequiredExperienceOccupationRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.JobInformationRequiredExperienceOccupation
	)
	if len(jobInformationIDList) == 0 {
		return requiredExperienceOccupationList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_required_experience_occupations.*
		FROM 
			job_information_required_experience_occupations
		INNER JOIN
			job_information_required_experience_jobs AS experience_job
		ON
			experience_job.id = job_information_required_experience_occupations.experience_job_id
		INNER JOIN
			job_information_required_conditions AS required_condition
		ON
			required_condition.id = experience_job.condition_id
		WHERE 
			required_condition.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredExperienceOccupationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceOccupationList, nil
}

// 指定求人IDリストと職種の必要職種経験を取得
func (repo *JobInformationRequiredExperienceOccupationRepositoryImpl) GetByJobInformationIDListAndOccupationList(jobInformationIDList []uint, occupationList []entity.ExperienceOccupation) ([]*entity.JobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.JobInformationRequiredExperienceOccupation
		occupationListStr                string
		lenOccupationList                = len(occupationList)
	)

	if len(jobInformationIDList) == 0 || lenOccupationList == 0 {
		return requiredExperienceOccupationList, nil
	}

	for i, occupation := range occupationList {
		if occupation.Occupation.Valid {
			occupationListStr += fmt.Sprint(occupation.Occupation.Int64)
			if i != lenOccupationList-1 {
				occupationListStr += ", "
			}
		}
	}

	query := fmt.Sprintf(`
			SELECT 
				job_information_required_experience_occupations.*
			FROM 
				job_information_required_experience_occupations
			INNER JOIN
				job_information_required_experience_jobs AS experience_job
			ON
				experience_job.id = job_information_required_experience_occupations.experience_job_id
			INNER JOIN
				job_information_required_conditions AS required_condition
			ON
				required_condition.id = experience_job.condition_id
			WHERE 
				required_condition.job_information_id IN (%s)
			AND
				experience_occupation IN (%s)
		`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"),
		occupationListStr,
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDListAndOccupationList",
		&requiredExperienceOccupationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceOccupationList, nil
}

func (repo *JobInformationRequiredExperienceOccupationRepositoryImpl) All() ([]*entity.JobInformationRequiredExperienceOccupation, error) {
	var (
		requiredExperienceOccupationList []*entity.JobInformationRequiredExperienceOccupation
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&requiredExperienceOccupationList, `
			SELECT *
			FROM job_information_required_experience_occupations
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceOccupationList, nil
}
