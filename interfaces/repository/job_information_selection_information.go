package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type JobInformationSelectionInformationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationSelectionInformationRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationSelectionInformationRepository {
	return &JobInformationSelectionInformationRepositoryImpl{
		Name:     "JobInformationSelectionInformationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
func (repo *JobInformationSelectionInformationRepositoryImpl) Create(selectionInformation *entity.JobInformationSelectionInformation) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_selection_informations (
				selection_flow_id,
				selection_type,
				selection_point,
				passed_example,
				fail_example,
				passing_rate,
				is_questionnairy,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?
			)
		`,
		selectionInformation.SelectionFlowID,
		selectionInformation.SelectionType,
		selectionInformation.SelectionPoint,
		selectionInformation.PassedExample,
		selectionInformation.FailExample,
		selectionInformation.PassingRate,
		selectionInformation.IsQuestionnairy,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	selectionInformation.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 更新 API
//
func (repo *JobInformationSelectionInformationRepositoryImpl) Update(id uint, selectionInformation *entity.JobInformationSelectionInformation) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE 
		job_information_selection_informations 
		SET
			selection_type = ?,
			selection_point = ?,
			passed_example = ?,
			fail_example = ?,
			passing_rate = ?,
			is_questionnairy = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		selectionInformation.SelectionType,
		selectionInformation.SelectionPoint,
		selectionInformation.PassedExample,
		selectionInformation.FailExample,
		selectionInformation.PassingRate,
		selectionInformation.IsQuestionnairy,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
// 単数取得 API
//
func (repo *JobInformationSelectionInformationRepositoryImpl) FindBySelectionFlowIDAndSelectionType(selectionFlowID, selectionType uint) (*entity.JobInformationSelectionInformation, error) {
	var (
		selectionInformation entity.JobInformationSelectionInformation
	)

	err := repo.executer.Get(
		repo.Name+".FindBySelectionFlowIDAndSelectionType",
		&selectionInformation, `
		SELECT *
		FROM job_information_selection_informations
		WHERE
			selection_flow_id = ?
		AND
			selection_type = ?
		LIMIT 1
		`,
		selectionFlowID, selectionType,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &selectionInformation, nil
}

// 指定求職者の選考情報を取得
func (repo *JobInformationSelectionInformationRepositoryImpl) FindByJobSeekerIDAndJobInformationIDAndSelectionType(jobSeekerID, jobInformationID uint, selectionType null.Int) (*entity.JobInformationSelectionInformation, error) {
	var (
		selectionInformation entity.JobInformationSelectionInformation
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerIDAndJobInformationIDAndSelectionType",
		&selectionInformation, `
		SELECT 
			selection_info.*
		FROM 
			job_information_selection_informations AS selection_info
		INNER JOIN
			job_information_selection_flow_patterns AS flow_patterns
		ON
			selection_info.selection_flow_id = flow_patterns.id
		INNER JOIN
			job_informations AS job_info
		ON
			flow_patterns.job_information_id = job_info.id
		INNER JOIN
			task_groups
		ON
			job_info.id = task_groups.job_information_id
		WHERE
			task_groups.job_seeker_id = ? AND
			task_groups.job_information_id = ? AND
			selection_info.selection_type = ?
		LIMIT 1
		`,
		jobSeekerID, jobInformationID, selectionType,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &selectionInformation, nil
}

/****************************************************************************************/
// 複数取得 API
//
func (repo *JobInformationSelectionInformationRepositoryImpl) GetBySelectionFlowID(selectionFlowID uint) ([]*entity.JobInformationSelectionInformation, error) {
	var (
		selectionInformationList []*entity.JobInformationSelectionInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetBySelectionFlowID",
		&selectionInformationList, `
		SELECT *
		FROM job_information_selection_informations
		WHERE
			selection_flow_id = ?
		ORDER BY 
			selection_type ASC
		`,
		selectionFlowID,
	)

	if err != nil {
		fmt.Println(err)
		return selectionInformationList, err
	}

	return selectionInformationList, nil
}

func (repo *JobInformationSelectionInformationRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationSelectionInformation, error) {
	var (
		selectionInformationList []*entity.JobInformationSelectionInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&selectionInformationList, `
		SELECT *
		FROM job_information_selection_informations
		WHERE
			selection_flow_id IN (
				SELECT id
				FROM job_information_selection_flow_patterns
				WHERE job_information_id= ?
			)
		ORDER BY 
			selection_type ASC
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return selectionInformationList, err
	}

	return selectionInformationList, nil
}

func (repo *JobInformationSelectionInformationRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationSelectionInformation, error) {
	var (
		selectionInformationList []*entity.JobInformationSelectionInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&selectionInformationList, `
			SELECT *
			FROM job_information_selection_informations
			WHERE selection_flow_id IN (
				SELECT id
				FROM job_information_selection_flow_patterns
				WHERE job_information_id IN (
					SELECT id
					FROM job_informations
					WHERE billing_address_id = ?
				)
			)
			ORDER BY 
				selection_type ASC
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionInformationList, nil
}

func (repo *JobInformationSelectionInformationRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationSelectionInformation, error) {
	var (
		selectionInformationList []*entity.JobInformationSelectionInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&selectionInformationList, `
			SELECT *
			FROM job_information_selection_informations
			WHERE selection_flow_id IN (
				SELECT id
				FROM job_information_selection_flow_patterns
				WHERE job_information_id IN (
					SELECT id
					FROM job_informations
					WHERE billing_address_id IN (
						SELECT id
						FROM billing_addresses
						WHERE enterprise_id = ?
					)
				)
			)
			ORDER BY 
				selection_type ASC
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionInformationList, nil
}

func (repo *JobInformationSelectionInformationRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationSelectionInformation, error) {
	var (
		selectionInformationList []*entity.JobInformationSelectionInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&selectionInformationList, `
			SELECT *
			FROM job_information_selection_informations
			WHERE selection_flow_id IN (
				SELECT id
				FROM job_information_selection_flow_patterns
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
			ORDER BY 
				selection_type ASC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionInformationList, nil
}

// 求人リストから選考情報を取得
func (repo *JobInformationSelectionInformationRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationSelectionInformation, error) {
	var (
		selectionInformationList []*entity.JobInformationSelectionInformation
	)

	if len(jobInformationIDList) == 0 {
		return selectionInformationList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_information_selection_informations.*
		FROM 
			job_information_selection_informations
		INNER JOIN
			job_information_selection_flow_patterns AS flow_pattern
		ON
			flow_pattern.id = job_information_selection_informations.selection_flow_id
		WHERE 
			flow_pattern.job_information_id IN (%s)
		ORDER BY 
			selection_type ASC
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&selectionInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionInformationList, nil
}

func (repo *JobInformationSelectionInformationRepositoryImpl) All() ([]*entity.JobInformationSelectionInformation, error) {
	var (
		selectionInformationList []*entity.JobInformationSelectionInformation
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&selectionInformationList, `
			SELECT *
			FROM job_information_selection_informations
			ORDER BY selection_type ASC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionInformationList, nil
}
