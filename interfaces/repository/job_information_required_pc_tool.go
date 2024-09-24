package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationRequiredPCToolRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRequiredPCToolRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRequiredPCToolRepository {
	return &JobInformationRequiredPCToolRepositoryImpl{
		Name:     "JobInformationRequiredPCToolRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
func (repo *JobInformationRequiredPCToolRepositoryImpl) Create(requiredPCTool *entity.JobInformationRequiredPCTool) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_required_pc_tools (
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

/****************************************************************************************/
// 複数取得 API
//
func (repo *JobInformationRequiredPCToolRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.JobInformationRequiredPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&requiredPCToolList, `
		SELECT *
		FROM job_information_required_pc_tools
		WHERE
			condition_id IN (
				SELECT id
				FROM job_information_required_conditions
				WHERE	job_information_id = ?
			)
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return requiredPCToolList, err
	}

	return requiredPCToolList, nil
}

func (repo *JobInformationRequiredPCToolRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.JobInformationRequiredPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&requiredPCToolList, `
			SELECT *
			FROM job_information_required_pc_tools
			WHERE
				condition_id IN (
					SELECT id
					FROM job_information_required_conditions
					WHERE	job_information_id IN (
						SELECT id
						FROM job_informations
						WHERE
						billing_address_id = ?
					)
				)
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredPCToolList, nil
}

func (repo *JobInformationRequiredPCToolRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.JobInformationRequiredPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&requiredPCToolList, `
			SELECT *
			FROM job_information_required_pc_tools
			WHERE
				condition_id IN (
					SELECT id
					FROM job_information_required_conditions
					WHERE	job_information_id IN (
						SELECT id
						FROM job_informations
						WHERE 
							billing_address_id IN (
								SELECT id
								FROM billing_addresses
								WHERE enterprise_id = ?
							) 
						)
					)
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredPCToolList, nil
}

func (repo *JobInformationRequiredPCToolRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.JobInformationRequiredPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentsID",
		&requiredPCToolList, `
			SELECT *
			FROM job_information_required_pc_tools
			WHERE
				condition_id IN (
					SELECT id
					FROM job_information_required_conditions
					WHERE job_information_id IN (
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
				)
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredPCToolList, nil
}

// 求人リストから必要PCスキルを取得
func (repo *JobInformationRequiredPCToolRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.JobInformationRequiredPCTool
	)

	if len(jobInformationIDList) == 0 {
		return requiredPCToolList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_required_pc_tools.*
		FROM 
			job_information_required_pc_tools
		INNER JOIN
			job_information_required_conditions AS required_condition
		ON
			required_condition.id = job_information_required_pc_tools.condition_id
		WHERE
			required_condition.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&requiredPCToolList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredPCToolList, nil
}

func (repo *JobInformationRequiredPCToolRepositoryImpl) All() ([]*entity.JobInformationRequiredPCTool, error) {
	var (
		requiredPCToolList []*entity.JobInformationRequiredPCTool
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&requiredPCToolList, `
			SELECT *
			FROM job_information_required_pc_tools
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredPCToolList, nil
}
