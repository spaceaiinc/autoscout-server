package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationWorkCharmPointRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationWorkCharmPointRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationWorkCharmPointRepository {
	return &SendingJobInformationWorkCharmPointRepositoryImpl{
		Name:     "SendingJobInformationWorkCharmPointRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationWorkCharmPointRepositoryImpl) Create(workCharmPoint *entity.SendingJobInformationWorkCharmPoint) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_work_charm_points (
				sending_job_information_id,
				title,
				contents,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		workCharmPoint.SendingJobInformationID,
		workCharmPoint.Title,
		workCharmPoint.Contents,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	workCharmPoint.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationWorkCharmPointRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_work_charm_points
		WHERE	sending_job_information_id = ?
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationWorkCharmPointRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.SendingJobInformationWorkCharmPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&workCharmPointList, `
		SELECT *
		FROM sending_job_information_work_charm_points
		WHERE	sending_job_information_id = ?
		`,
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return workCharmPointList, err
	}

	return workCharmPointList, nil
}

func (repo *SendingJobInformationWorkCharmPointRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.SendingJobInformationWorkCharmPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&workCharmPointList, `
			SELECT *
			FROM sending_job_information_work_charm_points
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

	return workCharmPointList, nil
}

func (repo *SendingJobInformationWorkCharmPointRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.SendingJobInformationWorkCharmPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&workCharmPointList, `
			SELECT *
			FROM sending_job_information_work_charm_points
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

	return workCharmPointList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationWorkCharmPointRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.SendingJobInformationWorkCharmPoint
	)

	if len(idList) == 0 {
		return workCharmPointList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_work_charm_points.*
		FROM 
			sending_job_information_work_charm_points
		WHERE
			sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&workCharmPointList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return workCharmPointList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationWorkCharmPointRepositoryImpl) GetAll() ([]*entity.SendingJobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.SendingJobInformationWorkCharmPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&workCharmPointList, `
			SELECT *
			FROM sending_job_information_work_charm_points
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return workCharmPointList, nil
}
