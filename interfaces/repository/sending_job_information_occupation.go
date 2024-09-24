package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationOccupationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationOccupationRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationOccupationRepository {
	return &SendingJobInformationOccupationRepositoryImpl{
		Name:     "SendingJobInformationOccupationRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationOccupationRepositoryImpl) Create(occupation *entity.SendingJobInformationOccupation) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_occupations (
				sending_job_information_id,
				occupation,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		occupation.SendingJobInformationID,
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

func (repo *SendingJobInformationOccupationRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_occupations
		WHERE	sending_job_information_id = ?
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationOccupationRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationOccupation, error) {
	var (
		occupationList []*entity.SendingJobInformationOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&occupationList, `
		SELECT *
		FROM sending_job_information_occupations
		WHERE	sending_job_information_id = ?
		`,
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return occupationList, err
	}

	return occupationList, nil
}

func (repo *SendingJobInformationOccupationRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationOccupation, error) {
	var (
		occupationList []*entity.SendingJobInformationOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&occupationList, `
			SELECT *
			FROM sending_job_information_occupations
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

	return occupationList, nil
}

func (repo *SendingJobInformationOccupationRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationOccupation, error) {
	var (
		occupationList []*entity.SendingJobInformationOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&occupationList, `
			SELECT *
			FROM sending_job_information_occupations
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

	return occupationList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationOccupationRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationOccupation, error) {
	var (
		occupationList []*entity.SendingJobInformationOccupation
	)

	if len(idList) == 0 {
		return occupationList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_occupations.*
		FROM 
			sending_job_information_occupations
		WHERE
			sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&occupationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return occupationList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationOccupationRepositoryImpl) GetAll() ([]*entity.SendingJobInformationOccupation, error) {
	var (
		occupationList []*entity.SendingJobInformationOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&occupationList, `
			SELECT *
			FROM sending_job_information_occupations
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return occupationList, nil
}
