package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationTargetRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationTargetRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationTargetRepository {
	return &JobInformationTargetRepositoryImpl{
		Name:     "JobInformationTargetRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 募集対象を作成する
func (repo *JobInformationTargetRepositoryImpl) Create(target *entity.JobInformationTarget) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_targets (
				job_information_id,
				target,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		target.JobInformationID,
		target.Target,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	target.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 削除 API
//
// 指定求人IDの募集対象を削除する
func (repo *JobInformationTargetRepositoryImpl) DeleteByJobInformationID(jobInformationD uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobInformationID",
		`
		DELETE
		FROM job_information_targets
		WHERE job_information_id = ?
		`, jobInformationD,
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
// 指定求人IDの募集対象を取得する
func (repo *JobInformationTargetRepositoryImpl) GetByJobInformationID(jobInformationD uint) ([]*entity.JobInformationTarget, error) {
	var (
		targetList []*entity.JobInformationTarget
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&targetList, `
		SELECT *
		FROM job_information_targets
		WHERE
			job_information_id = ?
		`,
		jobInformationD,
	)

	if err != nil {
		fmt.Println(err)
		return targetList, err
	}

	return targetList, nil
}

// 指定請求先IDの募集対象を取得する
func (repo *JobInformationTargetRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationTarget, error) {
	var (
		targetList []*entity.JobInformationTarget
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&targetList, `
			SELECT *
			FROM job_information_targets
			WHERE job_information_id IN (
				SELECT id
				FROM job_informations
				WHERE billing_address_id = ?
			)
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return targetList, nil
}

// 指定企業IDの募集対象を取得する
func (repo *JobInformationTargetRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationTarget, error) {
	var (
		targetList []*entity.JobInformationTarget
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&targetList, `
			SELECT *
			FROM job_information_targets
			WHERE job_information_id IN (
				SELECT id
				FROM job_informations
				WHERE billing_address_id IN (
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

	return targetList, nil
}

// エージェントIDから募集対象を取得する
func (repo *JobInformationTargetRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationTarget, error) {
	var (
		targetList []*entity.JobInformationTarget
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&targetList, `
			SELECT *
			FROM job_information_targets
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

	return targetList, nil
}

// 求人IDリストから募集対象を取得
func (repo *JobInformationTargetRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationTarget, error) {
	var (
		targetList []*entity.JobInformationTarget
	)

	if len(jobInformationIDList) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_information_targets
		WHERE
			job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationList",
		&targetList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return targetList, nil
}

// 全ての募集対象を取得する
func (repo *JobInformationTargetRepositoryImpl) All() ([]*entity.JobInformationTarget, error) {
	var (
		targetList []*entity.JobInformationTarget
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&targetList, `
			SELECT *
			FROM job_information_targets
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return targetList, nil
}
