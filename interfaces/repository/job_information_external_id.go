package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationExternalIDRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationExternalIDRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationExternalIDRepository {
	return &JobInformationExternalIDRepositoryImpl{
		Name:     "JobInformationExternalIDRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 他媒体でのIDを作成
func (repo *JobInformationExternalIDRepositoryImpl) Create(externalID *entity.JobInformationExternalID) error {
	now := time.Now().In(time.UTC)
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_external_ids (
				job_information_id,
				external_type,
				external_id,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		externalID.JobInformationID,
		externalID.ExternalType,
		externalID.ExternalID,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	externalID.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 削除 API
//
// 指定求人IDの他媒体でのIDを削除する
func (repo *JobInformationExternalIDRepositoryImpl) DeleteByJobInformationID(jobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobInformationID",
		`
		DELETE
		FROM job_information_external_ids
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
func (repo *JobInformationExternalIDRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationExternalID, error) {
	var (
		externalIDList []*entity.JobInformationExternalID
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&externalIDList, `
		SELECT *
		FROM job_information_external_ids
		WHERE job_information_id = ?
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return externalIDList, err
	}

	return externalIDList, nil
}

func (repo *JobInformationExternalIDRepositoryImpl) GetByAgentIDAndExternalType(agentID uint, externalType entity.JobInformatinoExternalType) ([]*entity.JobInformationExternalID, error) {
	var (
		externalIDList []*entity.JobInformationExternalID
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&externalIDList, `
			SELECT *
			FROM job_information_external_ids
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
			AND
				external_type = ?
		`,
		agentID,
		externalType,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return externalIDList, nil
}
