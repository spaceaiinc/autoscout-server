package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationRequiredExperienceIndustryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredExperienceIndustryRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredExperienceIndustryRepository {
	return &JobInformationRequiredExperienceIndustryRepositoryImpl{
		Name:     "JobInformationRequiredExperienceIndustryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 必要経験業種を作成
func (repo *JobInformationRequiredExperienceIndustryRepositoryImpl) Create(requiredExperienceIndustry *entity.JobInformationRequiredExperienceIndustry) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_experience_industries (
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
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredExperienceIndustry.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 複数取得 API
//
// 指定求人IDの必要経験業種を取得
func (repo *JobInformationRequiredExperienceIndustryRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.JobInformationRequiredExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&requiredExperienceIndustryList, `
		SELECT *
		FROM job_information_required_experience_industries
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
		return requiredExperienceIndustryList, err
	}

	return requiredExperienceIndustryList, nil
}

func (repo *JobInformationRequiredExperienceIndustryRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.JobInformationRequiredExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredExperienceIndustryList, `
			SELECT *
			FROM job_information_required_experience_industries
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

	return requiredExperienceIndustryList, nil
}

func (repo *JobInformationRequiredExperienceIndustryRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.JobInformationRequiredExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredExperienceIndustryList, `
			SELECT *
			FROM job_information_required_experience_industries
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

	return requiredExperienceIndustryList, nil
}

func (repo *JobInformationRequiredExperienceIndustryRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.JobInformationRequiredExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&requiredExperienceIndustryList, `
			SELECT *
			FROM job_information_required_experience_industries
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

	return requiredExperienceIndustryList, nil
}

// 求人リストから必要経験業職種を取得
func (repo *JobInformationRequiredExperienceIndustryRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.JobInformationRequiredExperienceIndustry
	)

	if len(jobInformationIDList) == 0 {
		return requiredExperienceIndustryList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_required_experience_industries.*
		FROM 
			job_information_required_experience_industries
		INNER JOIN
			job_information_required_experience_jobs AS experience_job
		ON
			experience_job.id = job_information_required_experience_industries.experience_job_id
		INNER JOIN
			job_information_required_conditions AS required_condition
		ON
			required_condition.id = experience_job.condition_id
		WHERE 
			required_condition.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredExperienceIndustryList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceIndustryList, nil
}

// 指定求人IDリストと業界の必要業界経験を取得
func (repo *JobInformationRequiredExperienceIndustryRepositoryImpl) GetByJobInformationIDListAndIndustryList(jobInformationIDList []uint, industryList []entity.ExperienceIndustry) ([]*entity.JobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.JobInformationRequiredExperienceIndustry
		industryListStr                string
		lenIndustryList                = len(industryList)
	)

	if len(jobInformationIDList) == 0 || lenIndustryList == 0 {
		return requiredExperienceIndustryList, nil
	}

	// 業界リストを文字列に変換
	for i, industry := range industryList {
		if industry.Industry.Valid {
			industryListStr += fmt.Sprint(industry.Industry.Int64)
			// 最後の要素以外はカンマを追加
			if i != lenIndustryList-1 {
				industryListStr += ", "
			}
		}
	}

	query := fmt.Sprintf(`
			SELECT 
				job_information_required_experience_industries.*
			FROM 
				job_information_required_experience_industries
			INNER JOIN
				job_information_required_experience_jobs AS experience_job
			ON
				experience_job.id = job_information_required_experience_industries.experience_job_id
			INNER JOIN
				job_information_required_conditions AS required_condition
			ON
				required_condition.id = experience_job.condition_id
			WHERE 
				required_condition.job_information_id IN (%s)
			AND 
				experience_industry IN (%s)
		`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"),
		industryListStr,
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDListAndIndustryList",
		&requiredExperienceIndustryList,
		query,
	)

	for _, requiredExperienceIndustry := range requiredExperienceIndustryList {
		fmt.Println(requiredExperienceIndustry.ExperienceIndustry)
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceIndustryList, nil
}

func (repo *JobInformationRequiredExperienceIndustryRepositoryImpl) All() ([]*entity.JobInformationRequiredExperienceIndustry, error) {
	var (
		requiredExperienceIndustryList []*entity.JobInformationRequiredExperienceIndustry
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&requiredExperienceIndustryList, `
			SELECT *
			FROM job_information_required_experience_industries
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredExperienceIndustryList, nil
}
