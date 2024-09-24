package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationOccupationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationOccupationRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationOccupationRepository {
	return &JobInformationOccupationRepositoryImpl{
		Name:     "JobInformationOccupationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 募集職種を作成
func (repo *JobInformationOccupationRepositoryImpl) Create(occupation *entity.JobInformationOccupation) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_occupations (
				job_information_id,
				occupation,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		occupation.JobInformationID,
		occupation.Occupation,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	occupation.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 削除 API
//
// 募集職種を削除
func (repo *JobInformationOccupationRepositoryImpl) DeleteByJobInformationID(jobInformationD uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobInformationID",
		`
		DELETE
		FROM job_information_occupations
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
// 指定求人IDの募集職種を取得
func (repo *JobInformationOccupationRepositoryImpl) GetByJobInformationID(jobInformationD uint) ([]*entity.JobInformationOccupation, error) {
	var (
		OccupationList []*entity.JobInformationOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&OccupationList, `
		SELECT *
		FROM job_information_occupations
		WHERE
			job_information_id = ?
		`,
		jobInformationD,
	)

	if err != nil {
		fmt.Println(err)
		return OccupationList, err
	}

	return OccupationList, nil
}

// 指定請求先IDの募集職種を取得
func (repo *JobInformationOccupationRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationOccupation, error) {
	var (
		OccupationList []*entity.JobInformationOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&OccupationList, `
			SELECT *
			FROM job_information_occupations
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

	return OccupationList, nil
}

// 指定企業IDの募集職種を取得
func (repo *JobInformationOccupationRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationOccupation, error) {
	var (
		OccupationList []*entity.JobInformationOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&OccupationList, `
			SELECT *
			FROM job_information_occupations
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

	return OccupationList, nil
}

func (repo *JobInformationOccupationRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationOccupation, error) {
	var (
		OccupationList []*entity.JobInformationOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&OccupationList, `
			SELECT *
			FROM job_information_occupations
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

	return OccupationList, nil
}

// 求人リストから募集対象を取得
func (repo *JobInformationOccupationRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationOccupation, error) {
	var (
		OccupationList []*entity.JobInformationOccupation
	)

	if len(jobInformationIDList) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_information_occupations
		WHERE
			job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&OccupationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return OccupationList, nil
}

func (repo *JobInformationOccupationRepositoryImpl) All() ([]*entity.JobInformationOccupation, error) {
	var (
		OccupationList []*entity.JobInformationOccupation
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&OccupationList, `
			SELECT *
			FROM job_information_occupations
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return OccupationList, nil
}
