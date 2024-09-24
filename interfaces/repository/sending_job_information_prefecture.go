package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationPrefectureRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationPrefectureRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationPrefectureRepository {
	return &SendingJobInformationPrefectureRepositoryImpl{
		Name:     "SendingJobInformationPrefectureRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationPrefectureRepositoryImpl) Create(prefecture *entity.SendingJobInformationPrefecture) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_prefectures (
				sending_job_information_id,
				prefecture,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		prefecture.SendingJobInformationID,
		prefecture.Prefecture,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	prefecture.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationPrefectureRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_prefectures
		WHERE	sending_job_information_id = ?
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationPrefectureRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationPrefecture, error) {
	var (
		prefectureList []*entity.SendingJobInformationPrefecture
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&prefectureList, `
		SELECT *
		FROM sending_job_information_prefectures
		WHERE	sending_job_information_id = ?
		`,
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return prefectureList, err
	}

	return prefectureList, nil
}

func (repo *SendingJobInformationPrefectureRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationPrefecture, error) {
	var (
		prefectureList []*entity.SendingJobInformationPrefecture
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&prefectureList, `
			SELECT *
			FROM sending_job_information_prefectures
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

	return prefectureList, nil
}

func (repo *SendingJobInformationPrefectureRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationPrefecture, error) {
	var (
		prefectureList []*entity.SendingJobInformationPrefecture
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&prefectureList, `
			SELECT *
			FROM sending_job_information_prefectures
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

	return prefectureList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationPrefectureRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationPrefecture, error) {
	var (
		prefectureList []*entity.SendingJobInformationPrefecture
	)

	if len(idList) == 0 {
		return prefectureList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_prefectures.*
		FROM 
			sending_job_information_prefectures
		WHERE
			sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&prefectureList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return prefectureList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationPrefectureRepositoryImpl) GetAll() ([]*entity.SendingJobInformationPrefecture, error) {
	var (
		prefectureList []*entity.SendingJobInformationPrefecture
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&prefectureList, `
			SELECT *
			FROM sending_job_information_prefectures
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return prefectureList, nil
}
