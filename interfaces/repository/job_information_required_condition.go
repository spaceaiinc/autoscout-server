package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationRequiredConditionRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredConditionRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredConditionRepository {
	return &JobInformationRequiredConditionRepositoryImpl{
		Name:     "JobInformationRequiredConditionRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 必要条件を作成
func (repo *JobInformationRequiredConditionRepositoryImpl) Create(requiredCondition *entity.JobInformationRequiredCondition) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_conditions (
				job_information_id,
				is_common,
				required_management,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		requiredCondition.JobInformationID,
		requiredCondition.IsCommon,
		requiredCondition.RequiredManagement,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredCondition.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 削除 API
//
// 必要条件を削除
func (repo *JobInformationRequiredConditionRepositoryImpl) DeleteByJobInformationID(jobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobInformationID",
		`
		DELETE
		FROM job_information_required_conditions
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
// 指定求人IDの必要条件を取得
func (repo *JobInformationRequiredConditionRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredCondition, error) {
	var (
		requiredConditionList []*entity.JobInformationRequiredCondition
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&requiredConditionList, `
		SELECT *
		FROM job_information_required_conditions
		WHERE
			job_information_id = ?
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return requiredConditionList, err
	}

	return requiredConditionList, nil
}

// 指定請求先IDの必要条件を取得
func (repo *JobInformationRequiredConditionRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredCondition, error) {
	var (
		requiredConditionList []*entity.JobInformationRequiredCondition
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredConditionList, `
			SELECT *
			FROM job_information_required_conditions
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

	return requiredConditionList, nil
}

// 指定企業IDの必要条件を取得
func (repo *JobInformationRequiredConditionRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredCondition, error) {
	var (
		requiredConditionList []*entity.JobInformationRequiredCondition
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredConditionList, `
			SELECT *
			FROM job_information_required_conditions
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

	return requiredConditionList, nil
}

func (repo *JobInformationRequiredConditionRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredCondition, error) {
	var (
		requiredConditionList []*entity.JobInformationRequiredCondition
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&requiredConditionList, `
			SELECT *
			FROM job_information_required_conditions
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

	return requiredConditionList, nil
}

// 求人リストから必要資格を取得
func (repo *JobInformationRequiredConditionRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredCondition, error) {
	var (
		requiredConditionList []*entity.JobInformationRequiredCondition
	)

	if len(jobInformationIDList) == 0 {
		return requiredConditionList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_information_required_conditions
		WHERE
			job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredConditionList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredConditionList, nil
}

func (repo *JobInformationRequiredConditionRepositoryImpl) All() ([]*entity.JobInformationRequiredCondition, error) {
	var (
		requiredConditionList []*entity.JobInformationRequiredCondition
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&requiredConditionList, `
			SELECT *
			FROM job_information_required_conditions
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredConditionList, nil
}
