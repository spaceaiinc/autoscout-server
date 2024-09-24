package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationRequiredLanguageTypeRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredLanguageTypeRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredLanguageTypeRepository {
	return &JobInformationRequiredLanguageTypeRepositoryImpl{
		Name:     "JobInformationRequiredLanguageTypeRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 必要言語を作成
func (repo *JobInformationRequiredLanguageTypeRepositoryImpl) Create(requiredLanguageType *entity.JobInformationRequiredLanguageType) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_language_types (
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
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredLanguageType.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 複数取得 API
//
// 指定求人IDの必要言語を取得
func (repo *JobInformationRequiredLanguageTypeRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.JobInformationRequiredLanguageType
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&requiredLanguageTypeList, `
		SELECT *
		FROM job_information_required_language_types
		WHERE
			language_id IN (
				SELECT id
				FROM job_information_required_languages
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
		return requiredLanguageTypeList, err
	}

	return requiredLanguageTypeList, nil
}

func (repo *JobInformationRequiredLanguageTypeRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.JobInformationRequiredLanguageType
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredLanguageTypeList, `
			SELECT *
			FROM job_information_required_language_types
			WHERE
				language_id IN (
					SELECT id
					FROM job_information_required_languages
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

	return requiredLanguageTypeList, nil
}

func (repo *JobInformationRequiredLanguageTypeRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.JobInformationRequiredLanguageType
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredLanguageTypeList, `
			SELECT *
			FROM job_information_required_language_types
			WHERE
				language_id IN (
					SELECT id
					FROM job_information_required_languages
					WHERE condition_id IN (
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

	return requiredLanguageTypeList, nil
}

func (repo *JobInformationRequiredLanguageTypeRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.JobInformationRequiredLanguageType
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&requiredLanguageTypeList, `
			SELECT *
			FROM job_information_required_language_types
			WHERE
				language_id IN (
					SELECT id
					FROM job_information_required_languages
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

	return requiredLanguageTypeList, nil
}

// 求人リストから必要言語を取得
func (repo *JobInformationRequiredLanguageTypeRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.JobInformationRequiredLanguageType
	)

	if len(jobInformationIDList) == 0 {
		return requiredLanguageTypeList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_required_language_types.*
		FROM 
			job_information_required_language_types
		INNER JOIN
			job_information_required_languages AS language
		ON
			language.id = job_information_required_language_types.language_id
		INNER JOIN
			job_information_required_conditions AS required_condition
		ON
			required_condition.id = language.condition_id
		WHERE
			required_condition.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredLanguageTypeList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLanguageTypeList, nil
}

// 指定求人IDリストと言語タイプの必要言語スキルの種類を取得
func (repo *JobInformationRequiredLanguageTypeRepositoryImpl) GetByJobInformationIDListAndLanguageTypeList(jobInformationIDList []uint, languageTypeList []entity.Language) ([]*entity.JobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.JobInformationRequiredLanguageType
		languageTypeListStr      string
		lenLanguageType          = len(languageTypeList)
	)

	if len(jobInformationIDList) == 0 || lenLanguageType == 0 {
		return requiredLanguageTypeList, nil
	}

	for i, languageType := range languageTypeList {
		if languageType.LanguageType.Valid {
			languageTypeListStr += fmt.Sprint(languageType.LanguageType.Int64)
			if i != lenLanguageType-1 {
				languageTypeListStr += ", "
			}
		}
	}

	query := fmt.Sprintf(`
			SELECT 
				job_information_required_language_types.*
			FROM 
				job_information_required_language_types
			INNER JOIN
				job_information_required_languages AS language
			ON
				language.id = job_information_required_language_types.language_id
			INNER JOIN
				job_information_required_conditions AS required_condition
			ON
				required_condition.id = language.condition_id
			WHERE
				required_condition.job_information_id IN (%s)
			AND
				language_type IN (%s)
		`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"),
		languageTypeListStr,
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDListAndLanguageTypeList",
		&requiredLanguageTypeList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLanguageTypeList, nil
}

func (repo *JobInformationRequiredLanguageTypeRepositoryImpl) All() ([]*entity.JobInformationRequiredLanguageType, error) {
	var (
		requiredLanguageTypeList []*entity.JobInformationRequiredLanguageType
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&requiredLanguageTypeList, `
			SELECT *
			FROM job_information_required_language_types
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredLanguageTypeList, nil
}
