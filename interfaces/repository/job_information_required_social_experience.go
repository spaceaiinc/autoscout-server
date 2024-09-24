package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationRequiredSocialExperienceRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredSocialExperienceRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredSocialExperienceRepository {
	return &JobInformationRequiredSocialExperienceRepositoryImpl{
		Name:     "JobInformationRequiredSocialExperienceRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
func (repo *JobInformationRequiredSocialExperienceRepositoryImpl) Create(requiredSocialExperience *entity.JobInformationRequiredSocialExperience) error {
	now := time.Now().In(time.UTC)
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_social_experiences (
				job_information_id,
				social_experience_type,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		requiredSocialExperience.JobInformationID,
		requiredSocialExperience.SocialExperienceType,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredSocialExperience.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 削除 API
//
func (repo *JobInformationRequiredSocialExperienceRepositoryImpl) DeleteByJobInformationID(jobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobInformationID",
		`
		DELETE
		FROM job_information_required_social_experiences
		WHERE job_information_id = ?
		`, jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
// 複数取得 API
//
func (repo *JobInformationRequiredSocialExperienceRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.JobInformationRequiredSocialExperience
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&requiredSocialExperienceList, `
		SELECT *
		FROM job_information_required_social_experiences
		WHERE
			job_information_id = ?
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return requiredSocialExperienceList, err
	}

	return requiredSocialExperienceList, nil
}

func (repo *JobInformationRequiredSocialExperienceRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.JobInformationRequiredSocialExperience
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredSocialExperienceList, `
			SELECT *
			FROM job_information_required_social_experiences
			WHERE
				job_information_id IN (
					SELECT id
					FROM job_informations
					WHERE
					billing_address_id = ?
				)
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredSocialExperienceList, nil
}

func (repo *JobInformationRequiredSocialExperienceRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.JobInformationRequiredSocialExperience
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredSocialExperienceList, `
			SELECT *
			FROM job_information_required_social_experiences
			WHERE
				job_information_id IN (
					SELECT id
					FROM job_informations
					WHERE 
						billing_address_id IN (
							SELECT id
							FROM billing_addresses
							WHERE enterprise_id = ?
						) 
					)
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredSocialExperienceList, nil
}

func (repo *JobInformationRequiredSocialExperienceRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.JobInformationRequiredSocialExperience
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&requiredSocialExperienceList, `
			SELECT *
			FROM job_information_required_social_experiences
			WHERE
				job_information_id IN (
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
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredSocialExperienceList, nil
}

// 求人リストから必要社会人経験を取得
func (repo *JobInformationRequiredSocialExperienceRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.JobInformationRequiredSocialExperience
	)

	if len(jobInformationIDList) == 0 {
		return requiredSocialExperienceList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_information_required_social_experiences
		WHERE
			job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredSocialExperienceList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredSocialExperienceList, nil
}

func (repo *JobInformationRequiredSocialExperienceRepositoryImpl) All() ([]*entity.JobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.JobInformationRequiredSocialExperience
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&requiredSocialExperienceList, `
			SELECT *
			FROM job_information_required_social_experiences
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredSocialExperienceList, nil
}
