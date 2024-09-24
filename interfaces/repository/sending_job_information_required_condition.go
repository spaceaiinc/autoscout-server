package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredConditionRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredConditionRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredConditionRepository {
	return &SendingJobInformationRequiredConditionRepositoryImpl{
		Name:     "SendingJobInformationRequiredConditionRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredConditionRepositoryImpl) Create(feature *entity.SendingJobInformationRequiredCondition) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_conditions (
				sending_job_information_id,
				is_common,
				required_management,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		feature.SendingJobInformationID,
		feature.IsCommon,
		feature.RequiredManagement,
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

func (repo *SendingJobInformationRequiredConditionRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_required_conditions
		WHERE	sending_job_information_id = ?
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationRequiredConditionRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationRequiredCondition, error) {
	var (
		featureList []*entity.SendingJobInformationRequiredCondition
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&featureList, `
		SELECT *
		FROM sending_job_information_required_conditions
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

func (repo *SendingJobInformationRequiredConditionRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredCondition, error) {
	var (
		featureList []*entity.SendingJobInformationRequiredCondition
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&featureList, `
			SELECT *
			FROM sending_job_information_required_conditions
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

func (repo *SendingJobInformationRequiredConditionRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredCondition, error) {
	var (
		featureList []*entity.SendingJobInformationRequiredCondition
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&featureList, `
			SELECT *
			FROM sending_job_information_required_conditions
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
func (repo *SendingJobInformationRequiredConditionRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredCondition, error) {
	var (
		featureList []*entity.SendingJobInformationRequiredCondition
	)

	if len(idList) == 0 {
		return featureList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_conditions.*
		FROM 
			sending_job_information_required_conditions
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
func (repo *SendingJobInformationRequiredConditionRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredCondition, error) {
	var (
		featureList []*entity.SendingJobInformationRequiredCondition
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&featureList, `
			SELECT *
			FROM sending_job_information_required_conditions
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return featureList, nil
}
