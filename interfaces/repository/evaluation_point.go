package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type EvaluationPointRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewEvaluationPointRepositoryImpl(ex interfaces.SQLExecuter) usecase.EvaluationPointRepository {
	return &EvaluationPointRepositoryImpl{
		Name:     "EvaluationPointRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *EvaluationPointRepositoryImpl) Create(evaluationPoint *entity.EvaluationPoint) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO evaluation_points (
			task_id,
			selection_information_id,
			job_seeker_id,
			job_information_id,
			good_point,
			ng_point,
			is_passed,
			is_re_interview,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)`,
		evaluationPoint.TaskID,
		evaluationPoint.SelectionInformationID,
		evaluationPoint.JobSeekerID,
		evaluationPoint.JobInformationID,
		evaluationPoint.GoodPoint,
		evaluationPoint.NGPoint,
		evaluationPoint.IsPassed,
		evaluationPoint.IsReInterview,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	evaluationPoint.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
func (repo *EvaluationPointRepositoryImpl) Update(id uint, evaluationPoint *entity.EvaluationPoint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE evaluation_points
		SET
			good_point = ?,
			ng_point = ?,
			is_passed = ?,
			is_re_interview = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		evaluationPoint.GoodPoint,
		evaluationPoint.NGPoint,
		evaluationPoint.IsPassed,
		evaluationPoint.IsReInterview,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *EvaluationPointRepositoryImpl) FindByID(id uint) (*entity.EvaluationPoint, error) {
	var (
		evaluationPoint entity.EvaluationPoint
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&evaluationPoint, `
		SELECT *
		FROM evaluation_points
		WHERE
			id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &evaluationPoint, nil
}

/****************************************************************************************/
/// 単数取得
//
func (repo *EvaluationPointRepositoryImpl) FindByTaskID(taskID uint) (*entity.EvaluationPoint, error) {
	var (
		evaluationPoint entity.EvaluationPoint
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&evaluationPoint, `
		SELECT *
		FROM evaluation_points
		WHERE
			task_id = ?
		LIMIT 1
		`,
		taskID)

	if err != nil {
		return nil, err
	}

	return &evaluationPoint, nil
}

/****************************************************************************************/
/// 複数取得
//
// agentStaffIDを使ってタスクグループの一覧を取得
func (repo *EvaluationPointRepositoryImpl) GetByJobSeekerIDAndJobInformationID(jobSeekerID, jobInformationID uint) ([]*entity.EvaluationPoint, error) {
	var (
		evaluationPointList []*entity.EvaluationPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDAndJobInformationID",
		&evaluationPointList, `
			SELECT *
			FROM evaluation_points
			WHERE
				job_seeker_id = ? 
			AND
				job_information_id = ?
			ORDER BY 
				created_at DESC
		`,
		jobSeekerID, jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return evaluationPointList, nil
}

func (repo *EvaluationPointRepositoryImpl) GetBySelectionFlowID(selectionFlowID uint) ([]*entity.EvaluationPoint, error) {
	var (
		evaluationPointList []*entity.EvaluationPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetBySelectionFlowID",
		&evaluationPointList, `
			SELECT 
				point.*,
				seeker.first_name AS first_name,
				seeker.last_name AS last_name,
				seeker.first_furigana AS first_furigana,
				seeker.last_furigana AS last_furigana,
				seeker.agent_id AS ca_agent_id
			FROM
				evaluation_points AS point
			INNER JOIN 
				job_seekers AS seeker
			ON
				point.job_seeker_id = seeker.id
			WHERE
				selection_information_id IN (
					SELECT id 
					FROM job_information_selection_informations 
					WHERE selection_flow_id = ?
				)
			ORDER BY id ASC
		`,
		selectionFlowID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return evaluationPointList, nil
}
