package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationRequiredLanguageRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredLanguageRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredLanguageRepository {
	return &JobInformationRequiredLanguageRepositoryImpl{
		Name:     "JobInformationRequiredLanguageRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
func (repo *JobInformationRequiredLanguageRepositoryImpl) Create(requiredLanguage *entity.JobInformationRequiredLanguage) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_languages (
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

/****************************************************************************************/
// 複数取得 API
//
// 指定求人IDの必要言語を取得
func (repo *JobInformationRequiredLanguageRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.JobInformationRequiredLanguage
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&requiredLanguageList, `
		SELECT *
		FROM job_information_required_languages
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
		return requiredLanguageList, err
	}

	return requiredLanguageList, nil
}

func (repo *JobInformationRequiredLanguageRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.JobInformationRequiredLanguage
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredLanguageList, `
			SELECT *
			FROM job_information_required_languages
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

	return requiredLanguageList, nil
}

func (repo *JobInformationRequiredLanguageRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.JobInformationRequiredLanguage
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredLanguageList, `
			SELECT *
			FROM job_information_required_languages
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

	return requiredLanguageList, nil
}

func (repo *JobInformationRequiredLanguageRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.JobInformationRequiredLanguage
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&requiredLanguageList, `
			SELECT *
			FROM job_information_required_languages
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

	return requiredLanguageList, nil
}

// 求人リストから必要言語を取得
func (repo *JobInformationRequiredLanguageRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.JobInformationRequiredLanguage
	)

	if len(jobInformationIDList) == 0 {
		return requiredLanguageList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_required_languages.*
		FROM 
			job_information_required_languages
		INNER JOIN
			job_information_required_conditions AS required_condition
		ON
			required_condition.id = job_information_required_languages.condition_id
		WHERE
			required_condition.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredLanguageList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLanguageList, nil
}

func (repo *JobInformationRequiredLanguageRepositoryImpl) All() ([]*entity.JobInformationRequiredLanguage, error) {
	var (
		requiredLanguageList []*entity.JobInformationRequiredLanguage
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&requiredLanguageList, `
			SELECT *
			FROM job_information_required_languages
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLanguageList, nil
}
