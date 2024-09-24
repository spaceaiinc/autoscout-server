package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationFeatureRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationFeatureRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationFeatureRepository {
	return &SendingJobInformationFeatureRepositoryImpl{
		Name:     "SendingJobInformationFeatureRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationFeatureRepositoryImpl) Create(feature *entity.SendingJobInformationFeature) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_features (
				sending_job_information_id,
				feature,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		feature.SendingJobInformationID,
		feature.Feature,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	feature.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationFeatureRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_features
		WHERE	sending_job_information_id = ?
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationFeatureRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationFeature, error) {
	var (
		featureList []*entity.SendingJobInformationFeature
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&featureList, `
		SELECT *
		FROM sending_job_information_features
		WHERE	sending_job_information_id = ?
		`,
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return featureList, err
	}

	return featureList, nil
}

func (repo *SendingJobInformationFeatureRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationFeature, error) {
	var (
		featureList []*entity.SendingJobInformationFeature
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&featureList, `
			SELECT *
			FROM sending_job_information_features
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

	return featureList, nil
}

func (repo *SendingJobInformationFeatureRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationFeature, error) {
	var (
		featureList []*entity.SendingJobInformationFeature
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&featureList, `
			SELECT *
			FROM sending_job_information_features
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

	return featureList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationFeatureRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationFeature, error) {
	var (
		featureList []*entity.SendingJobInformationFeature
	)

	if len(idList) == 0 {
		return featureList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_features.*
		FROM 
			sending_job_information_features
		WHERE
			sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&featureList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return featureList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationFeatureRepositoryImpl) GetAll() ([]*entity.SendingJobInformationFeature, error) {
	var (
		featureList []*entity.SendingJobInformationFeature
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&featureList, `
			SELECT *
			FROM sending_job_information_features
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return featureList, nil
}
