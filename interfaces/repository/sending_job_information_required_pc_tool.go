package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredPCToolRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredPCToolRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredPCToolRepository {
	return &SendingJobInformationRequiredPCToolRepositoryImpl{
		Name:     "SendingJobInformationRequiredPCToolRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredPCToolRepositoryImpl) Create(requiredPCTool *entity.SendingJobInformationRequiredPCTool) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_pc_tools (
				condition_id,
				tool,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		requiredPCTool.ConditionID,
		requiredPCTool.Tool,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredPCTool.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationRequiredPCToolRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_required_pc_tools
		WHERE 
			condition_id IN (
				SELECT id
				FROM sending_job_information_required_conditions
				WHERE	sending_job_information_id = ?
			)
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationRequiredPCToolRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.SendingJobInformationRequiredPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&requiredPCToolList, `
		SELECT *
		FROM sending_job_information_required_pc_tools
		WHERE
			condition_id IN (
				SELECT id
				FROM sending_job_information_required_conditions
				WHERE	sending_job_information_id = ?
			)
		`,
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return requiredPCToolList, err
	}

	return requiredPCToolList, nil
}

func (repo *SendingJobInformationRequiredPCToolRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.SendingJobInformationRequiredPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&requiredPCToolList, `
			SELECT *
			FROM sending_job_information_required_pc_tools
			WHERE
				condition_id IN (
					SELECT id
					FROM sending_job_information_required_conditions
					WHERE sending_job_information_id IN (
						SELECT id
						FROM sending_job_informations
						WHERE
						sending_billing_address_id = ?
					)
				)
		`,
		sendingBillingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredPCToolList, nil
}

func (repo *SendingJobInformationRequiredPCToolRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.SendingJobInformationRequiredPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&requiredPCToolList, `
			SELECT *
			FROM sending_job_information_required_pc_tools
			WHERE
				condition_id IN (
					SELECT id
					FROM sending_job_information_required_conditions
					WHERE sending_job_information_id IN (
						SELECT id
						FROM sending_job_informations
						WHERE 
							sending_billing_address_id IN (
								SELECT id
								FROM sending_billing_addresses
								WHERE sending_enterprise_id = ?
							) 
						)
					)
		`,
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredPCToolList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationRequiredPCToolRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.SendingJobInformationRequiredPCTool
	)

	if len(idList) == 0 {
		return requiredPCToolList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_pc_tools.*
		FROM 
			sending_job_information_required_pc_tools
		INNER JOIN
			sending_job_information_required_conditions AS required_condition
		ON
			required_condition.id = sending_job_information_required_pc_tools.condition_id
		WHERE
			required_condition.sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&requiredPCToolList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredPCToolList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRequiredPCToolRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.SendingJobInformationRequiredPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredPCToolList, `
			SELECT *
			FROM sending_job_information_required_pc_tools
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredPCToolList, nil
}
