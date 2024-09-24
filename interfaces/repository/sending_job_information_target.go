package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationTargetRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationTargetRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationTargetRepository {
	return &SendingJobInformationTargetRepositoryImpl{
		Name:     "SendingJobInformationTargetRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationTargetRepositoryImpl) Create(target *entity.SendingJobInformationTarget) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_targets (
				sending_job_information_id,
				target,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		target.SendingJobInformationID,
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

func (repo *SendingJobInformationTargetRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_targets
		WHERE 
			sending_job_information_id = ?
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationTargetRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationTarget, error) {
	var (
		targetList []*entity.SendingJobInformationTarget
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&targetList, `
		SELECT *
		FROM sending_job_information_targets
		WHERE	sending_job_information_id = ?
		`,
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return targetList, err
	}

	return targetList, nil
}

func (repo *SendingJobInformationTargetRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationTarget, error) {
	var (
		targetList []*entity.SendingJobInformationTarget
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&targetList, `
			SELECT *
			FROM sending_job_information_targets
			WHERE
				sending_job_information_id IN (
					SELECT id
					FROM sending_job_informations
					WHERE
					sending_billing_address_id = ?
				)
		`,
		sendingBillingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return targetList, nil
}

func (repo *SendingJobInformationTargetRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationTarget, error) {
	var (
		targetList []*entity.SendingJobInformationTarget
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&targetList, `
			SELECT *
			FROM sending_job_information_targets
			WHERE
				sending_job_information_id IN (
					SELECT id
					FROM sending_job_informations
					WHERE 
						sending_billing_address_id IN (
							SELECT id
							FROM sending_billing_addresses
							WHERE sending_enterprise_id = ?
						) 
					)
		`,
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return targetList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationTargetRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationTarget, error) {
	var (
		targetList []*entity.SendingJobInformationTarget
	)

	if len(idList) == 0 {
		return targetList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_targets.*
		FROM 
			sending_job_information_targets
		WHERE
			sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&targetList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return targetList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationTargetRepositoryImpl) GetAll() ([]*entity.SendingJobInformationTarget, error) {
	var (
		targetList []*entity.SendingJobInformationTarget
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&targetList, `
			SELECT *
			FROM sending_job_information_targets
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return targetList, nil
}
