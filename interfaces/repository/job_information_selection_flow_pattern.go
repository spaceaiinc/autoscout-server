package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationSelectionFlowPatternRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationSelectionFlowPatternRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationSelectionFlowPatternRepository {
	return &JobInformationSelectionFlowPatternRepositoryImpl{
		Name:     "JobInformationSelectionFlowPatternRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
func (repo *JobInformationSelectionFlowPatternRepositoryImpl) Create(selectionFlowPattern *entity.JobInformationSelectionFlowPattern) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_selection_flow_patterns (
				job_information_id,
				public_status,
				flow_title,
				flow_pattern,
				is_deleted,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?,?
			)
		`,
		selectionFlowPattern.JobInformationID,
		selectionFlowPattern.PublicStatus,
		selectionFlowPattern.FlowTitle,
		selectionFlowPattern.FlowPattern,
		selectionFlowPattern.IsDeleted,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	selectionFlowPattern.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 更新 API
//
func (repo *JobInformationSelectionFlowPatternRepositoryImpl) Update(id uint, selectionFlowPattern *entity.JobInformationSelectionFlowPattern) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE 
			job_information_selection_flow_patterns 
		SET
			public_status = ?,
			flow_title = ?,
			flow_pattern = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		selectionFlowPattern.PublicStatus,
		selectionFlowPattern.FlowTitle,
		selectionFlowPattern.FlowPattern,
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
// 削除 API
//
func (repo *JobInformationSelectionFlowPatternRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		UPDATE 
			job_information_selection_flow_patterns 
		SET
			is_deleted= ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		true,
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
func (repo *JobInformationSelectionFlowPatternRepositoryImpl) FindByID(id uint) (*entity.JobInformationSelectionFlowPattern, error) {
	var (
		selectionFlowPattern entity.JobInformationSelectionFlowPattern
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&selectionFlowPattern, `
		SELECT 
			flow_pattern.*, staff.agent_id
		FROM 
			job_information_selection_flow_patterns AS flow_pattern
		INNER JOIN
			job_informations AS job_info
		ON
			flow_pattern.job_information_id = job_info.id
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		WHERE
			flow_pattern.id = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &selectionFlowPattern, nil
}

/****************************************************************************************/
// 複数取得 API
//
func (repo *JobInformationSelectionFlowPatternRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationSelectionFlowPattern, error) {
	var (
		selectionFlowPatternList []*entity.JobInformationSelectionFlowPattern
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&selectionFlowPatternList, `
		SELECT *
		FROM job_information_selection_flow_patterns
		WHERE
			job_information_id = ?
		AND
			is_deleted = false
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return selectionFlowPatternList, err
	}

	return selectionFlowPatternList, nil
}

// public_statusが0のレコードのみ取得
func (repo *JobInformationSelectionFlowPatternRepositoryImpl) GetOpenByJobInformationID(jobInformationID uint) ([]*entity.JobInformationSelectionFlowPattern, error) {
	var (
		selectionFlowPatternList []*entity.JobInformationSelectionFlowPattern
	)

	err := repo.executer.Select(
		repo.Name+".GetOpenByJobInformationID",
		&selectionFlowPatternList, `
		SELECT *
		FROM job_information_selection_flow_patterns
		WHERE
			job_information_id = ?
		AND
			public_status = 0
		AND
			is_deleted = false
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return selectionFlowPatternList, err
	}

	return selectionFlowPatternList, nil
}

func (repo *JobInformationSelectionFlowPatternRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationSelectionFlowPattern, error) {
	var (
		selectionFlowPatternList []*entity.JobInformationSelectionFlowPattern
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&selectionFlowPatternList, `
			SELECT *
			FROM job_information_selection_flow_patterns
			WHERE 
				job_information_id IN (
					SELECT id
					FROM job_informations
					WHERE 
						billing_address_id = ?
				)
				AND
					is_deleted = false
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionFlowPatternList, nil
}

func (repo *JobInformationSelectionFlowPatternRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationSelectionFlowPattern, error) {
	var (
		selectionFlowPatternList []*entity.JobInformationSelectionFlowPattern
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&selectionFlowPatternList, `
			SELECT *
			FROM job_information_selection_flow_patterns
			WHERE job_information_id IN (
				SELECT id
				FROM job_informations
				WHERE 
					billing_address_id IN (
						SELECT id
						FROM billing_addresses
						WHERE enterprise_id = ?
					)
				AND
					is_deleted = false
			)
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionFlowPatternList, nil
}

func (repo *JobInformationSelectionFlowPatternRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationSelectionFlowPattern, error) {
	var (
		selectionFlowPatternList []*entity.JobInformationSelectionFlowPattern
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&selectionFlowPatternList, `
			SELECT *
			FROM job_information_selection_flow_patterns
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
				is_deleted = false
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionFlowPatternList, nil
}

// idリストから選考フローパターンを取得
func (repo *JobInformationSelectionFlowPatternRepositoryImpl) GetByIDList(idList []uint) ([]*entity.JobInformationSelectionFlowPattern, error) {
	var (
		selectionFlowPatternList []*entity.JobInformationSelectionFlowPattern
	)
	if len(idList) == 0 {
		return selectionFlowPatternList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			flow_pattern.id, 
			flow_pattern.flow_pattern, 
			selection_info.id AS selection_information_id, 
			selection_info.selection_type,
			IFNULL(selection_info.is_questionnairy, FALSE) AS is_questionnairy
		FROM 
			job_information_selection_flow_patterns AS flow_pattern
		LEFT OUTER JOIN
			job_information_selection_informations AS selection_info
		ON
			flow_pattern.id = selection_info.selection_flow_id
		WHERE
			flow_pattern.id IN (%s)
		AND
			flow_pattern.is_deleted = false
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByIDList",
		&selectionFlowPatternList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionFlowPatternList, nil
}

// 求人リストから選考フローパターンを取得
func (repo *JobInformationSelectionFlowPatternRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationSelectionFlowPattern, error) {
	var (
		selectionFlowPatternList []*entity.JobInformationSelectionFlowPattern
	)
	if len(jobInformationIDList) == 0 {
		return selectionFlowPatternList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_information_selection_flow_patterns
		WHERE
			job_information_id IN (%s)
		AND
			is_deleted = false
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&selectionFlowPatternList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionFlowPatternList, nil
}

func (repo *JobInformationSelectionFlowPatternRepositoryImpl) All() ([]*entity.JobInformationSelectionFlowPattern, error) {
	var (
		selectionFlowPatternList []*entity.JobInformationSelectionFlowPattern
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&selectionFlowPatternList, `
			SELECT *
			FROM job_information_selection_flow_patterns
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selectionFlowPatternList, nil
}
