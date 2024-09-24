package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SelectionQuestionnaireRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSelectionQuestionnaireRepositoryImpl(ex interfaces.SQLExecuter) usecase.SelectionQuestionnaireRepository {
	return &SelectionQuestionnaireRepositoryImpl{
		Name:     "SelectionQuestionnaireRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// 選考後アンケートの作成
func (repo *SelectionQuestionnaireRepositoryImpl) Create(selectionQuestionnaire *entity.SelectionQuestionnaire) error {
	selectionQuestionnaire.UUID = utility.CreateUUID()

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO selection_questionnaires (
			uuid,
			job_seeker_id,
			job_information_id,
			selection_information_id,
			my_ranking,
			my_ranking_reason,
			concern_point,
			continue_selection,
			my_ranking_detail,
			selection_question,
			remarks,
			is_answer,
			is_self_introduction,
			is_self_pr,
			is_retire_reason,
			is_job_change_axis,
			is_applying_reason,
			is_career_vision,
			is_reverse_question,
			intention_to_job_offer,
			intention_detail,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?,
				?,?
			)`,
		selectionQuestionnaire.UUID,
		selectionQuestionnaire.JobSeekerID,
		selectionQuestionnaire.JobInformationID,
		selectionQuestionnaire.SelectionInformationID,
		selectionQuestionnaire.MyRanking,
		selectionQuestionnaire.MyRankingReason,
		selectionQuestionnaire.ConcernPoint,
		selectionQuestionnaire.ContinueSelection,
		selectionQuestionnaire.MyRankingDetail,
		selectionQuestionnaire.SelectionQuestion,
		selectionQuestionnaire.Remarks,
		selectionQuestionnaire.IsAnswer,
		selectionQuestionnaire.IsSelfIntroduction,
		selectionQuestionnaire.IsSelfPR,
		selectionQuestionnaire.IsRetireReason,
		selectionQuestionnaire.IsJobChangeAxis,
		selectionQuestionnaire.IsApplyingReason,
		selectionQuestionnaire.IsCareerVision,
		selectionQuestionnaire.IsReverseQuestion,
		selectionQuestionnaire.IntentionToJobOffer,
		selectionQuestionnaire.IntentionDetail,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	selectionQuestionnaire.ID = uint(lastID)
	return nil
}

// UUIDから選考後アンケートの作成
func (repo *SelectionQuestionnaireRepositoryImpl) CreateByUUID(selectionQuestionnaire *entity.SelectionQuestionnaire, questionnaireUUID uuid.UUID) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".CreateByUUID",
		`INSERT INTO selection_questionnaires (
			uuid,
			job_seeker_id,
			job_information_id,
			selection_information_id,
			my_ranking,
			my_ranking_reason,
			concern_point,
			continue_selection,
			my_ranking_detail,
			selection_question,
			remarks,
			is_answer,
			is_self_introduction,
			is_self_pr,
			is_retire_reason,
			is_job_change_axis,
			is_applying_reason,
			is_career_vision,
			is_reverse_question,
			intention_to_job_offer,
			intention_detail,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?,
				?, ?
			)`,
		questionnaireUUID,
		selectionQuestionnaire.JobSeekerID,
		selectionQuestionnaire.JobInformationID,
		selectionQuestionnaire.SelectionInformationID,
		selectionQuestionnaire.MyRanking,
		selectionQuestionnaire.MyRankingReason,
		selectionQuestionnaire.ConcernPoint,
		selectionQuestionnaire.ContinueSelection,
		selectionQuestionnaire.MyRankingDetail,
		selectionQuestionnaire.SelectionQuestion,
		selectionQuestionnaire.Remarks,
		selectionQuestionnaire.IsAnswer,
		selectionQuestionnaire.IsSelfIntroduction,
		selectionQuestionnaire.IsSelfPR,
		selectionQuestionnaire.IsRetireReason,
		selectionQuestionnaire.IsJobChangeAxis,
		selectionQuestionnaire.IsApplyingReason,
		selectionQuestionnaire.IsCareerVision,
		selectionQuestionnaire.IsReverseQuestion,
		selectionQuestionnaire.IntentionToJobOffer,
		selectionQuestionnaire.IntentionDetail,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	selectionQuestionnaire.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
// 選考後アンケートの更新
func (repo *SelectionQuestionnaireRepositoryImpl) Update(id uint, selectionQuestionnaire *entity.SelectionQuestionnaire) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE selection_questionnaires
		SET
			my_ranking = ?,
			my_ranking_reason = ?,
			concern_point = ?,
			continue_selection = ?,
			my_ranking_detail = ?,
			selection_question = ?,
			remarks = ?,
			is_answer = ?,
			is_self_introduction=?,
			is_self_pr=?,
			is_retire_reason=?,
			is_job_change_axis=?,
			is_applying_reason=?,
			is_career_vision=?,
			is_reverse_question=?,
			intention_to_job_offer=?,
			intention_detail=?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		selectionQuestionnaire.MyRanking,
		selectionQuestionnaire.MyRankingReason,
		selectionQuestionnaire.ConcernPoint,
		selectionQuestionnaire.ContinueSelection,
		selectionQuestionnaire.MyRankingDetail,
		selectionQuestionnaire.SelectionQuestion,
		selectionQuestionnaire.Remarks,
		selectionQuestionnaire.IsAnswer,
		selectionQuestionnaire.IsSelfIntroduction,
		selectionQuestionnaire.IsSelfPR,
		selectionQuestionnaire.IsRetireReason,
		selectionQuestionnaire.IsJobChangeAxis,
		selectionQuestionnaire.IsApplyingReason,
		selectionQuestionnaire.IsCareerVision,
		selectionQuestionnaire.IsReverseQuestion,
		selectionQuestionnaire.IntentionToJobOffer,
		selectionQuestionnaire.IntentionDetail,
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
/// 単数取得 API
//
// IDから選考後アンケートの取得
func (repo *SelectionQuestionnaireRepositoryImpl) FindByID(id uint) (*entity.SelectionQuestionnaire, error) {
	var (
		selectionQuestionnaire entity.SelectionQuestionnaire
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&selectionQuestionnaire, `
		SELECT 
			question.*, selection_info.selection_type
		FROM 
			selection_questionnaires AS question
		INNER JOIN 
			job_information_selection_informations AS selection_info
		ON
			selection_info.id = question.selection_information_id
		WHERE
			id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &selectionQuestionnaire, nil
}

// UUIDから選考後アンケートの取得
func (repo *SelectionQuestionnaireRepositoryImpl) FindByUUID(uuid uuid.UUID) (*entity.SelectionQuestionnaire, error) {
	var (
		selectionQuestionnaire entity.SelectionQuestionnaire
	)

	err := repo.executer.Get(
		repo.Name+".FindByUUID",
		&selectionQuestionnaire, `
		SELECT 
			question.*, selection_info.selection_type,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
				THEN task_group.external_job_information_title
				ELSE job_info.title
			END AS title,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
				THEN task_group.external_company_name
				ELSE enterprise.company_name
			END AS company_name
		FROM 
			selection_questionnaires AS question
		INNER JOIN 
			job_information_selection_informations AS selection_info
		ON
			selection_info.id = question.selection_information_id
		INNER JOIN 
			job_informations AS job_info
		ON
			job_info.id = question.job_information_id
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			task_groups AS task_group
		ON
			question.job_seeker_id = task_group.job_seeker_id AND
			question.job_information_id = task_group.job_information_id
		WHERE
			question.uuid = ?
		LIMIT 1
		`,
		uuid)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &selectionQuestionnaire, nil
}

// 求職者IDと選考フローIDから選考後アンケートを取得
func (repo *SelectionQuestionnaireRepositoryImpl) FindByJobSeekerIDAndJobInformationIDAndSelectionPhase(jobSeekerID, jobInformationID uint, selectionPhase uint) (*entity.SelectionQuestionnaire, error) {
	var (
		selectionQuestionnaire entity.SelectionQuestionnaire
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerIDAndJobInformationIDAndSelectionInformationID",
		&selectionQuestionnaire, `
		SELECT 
			question.*, selection_info.selection_type
		FROM 
			selection_questionnaires AS question
		INNER JOIN 
			job_information_selection_informations AS selection_info
		ON
			selection_info.id = question.selection_information_id
		INNER JOIN 
			job_informations AS job_info
		ON
			question.job_information_id = job_info.id
		INNER JOIN
			job_seekers AS seeker
		ON
			question.job_seeker_id = seeker.id
		WHERE
			seeker.id = ? AND
			job_info.id = ? AND
			selection_info.selection_type = ?
		LIMIT 1
		`, jobSeekerID, jobInformationID, selectionPhase)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &selectionQuestionnaire, nil
}

// 求職者UUIDと選考フローUUIDから選考後アンケートを取得
func (repo *SelectionQuestionnaireRepositoryImpl) FindByJobSeekerUUIDAndJobInformationUUIDAndSelectionPhase(jobSeekerUUID, jobInformationUUID uuid.UUID, selectionPhase uint) (*entity.SelectionQuestionnaire, error) {
	var (
		selectionQuestionnaire entity.SelectionQuestionnaire
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerUUIDAndJobInformationUUIDAndSelectionPhase",
		&selectionQuestionnaire, `
		SELECT 
			question.*, selection_info.selection_type,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
				THEN task_group.external_job_information_title
				ELSE job_info.title
			END AS title,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
				THEN task_group.external_company_name
				ELSE enterprise.company_name
			END AS company_name
		FROM 
			selection_questionnaires AS question
		INNER JOIN 
			job_information_selection_informations AS selection_info
		ON
			selection_info.id = question.selection_information_id
		INNER JOIN 
			job_informations AS job_info
		ON
			job_info.id = question.job_information_id
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			job_seekers AS seeker
		ON
			question.job_seeker_id = seeker.id
		INNER JOIN
			task_groups AS task_group
		ON
			question.job_seeker_id = task_group.job_seeker_id AND
			question.job_information_id = task_group.job_information_id
		WHERE
			seeker.uuid = ? AND
			job_info.uuid = ? AND
			selection_info.selection_type = ?
		LIMIT 1
		`, jobSeekerUUID, jobInformationUUID, selectionPhase)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &selectionQuestionnaire, nil
}

// 求職者IDから選考後アンケートを取得
func (repo *SelectionQuestionnaireRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.SelectionQuestionnaire, error) {
	var (
		list []*entity.SelectionQuestionnaire
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&list, `
		SELECT 
			question.*, selection_info.selection_type
		FROM 
			selection_questionnaires AS question
		INNER JOIN 
			job_information_selection_informations AS selection_info
		ON
			selection_info.id = question.selection_information_id
		WHERE
			question.job_seeker_id = ?
		`,
		jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}

// 求職者IDと選考フローIDから選考後アンケートを取得
func (repo *SelectionQuestionnaireRepositoryImpl) GetBySelectionFlowID(selectionFlowID uint) ([]*entity.SelectionQuestionnaire, error) {
	var (
		list []*entity.SelectionQuestionnaire
	)

	err := repo.executer.Select(
		repo.Name+".GetBySelectionFlowID",
		&list, `
		SELECT 
			question.*, selection_info.selection_type
		FROM 
			selection_questionnaires AS question
		INNER JOIN 
			job_information_selection_informations AS selection_info
		ON
			selection_info.id = question.selection_information_id
		WHERE
			question.selection_information_id IN (
				SELECT id
				FROM job_information_selection_informations
				WHERE selection_flow_id = ?
			)
		ORDER BY id ASC
		`,
		selectionFlowID,
	)

	if err != nil {
		return nil, err
	}

	return list, nil
}

// 指定uuidに合致するis_answerdがfalseのレコードを取得
func (repo *SelectionQuestionnaireRepositoryImpl) GetFirstListBySelectionFlowID(selectionFlowID uint) ([]*entity.SelectionQuestionnaire, error) {
	var (
		list []*entity.SelectionQuestionnaire
	)

	err := repo.executer.Select(
		repo.Name+".GetFirstListBySelectionFlowID",
		&list, `
		SELECT 
			question.*, selection_info.selection_type
		FROM 
			selection_questionnaires AS question
		INNER JOIN 
			job_information_selection_informations AS selection_info
		ON
			selection_info.id = question.selection_information_id
		WHERE 
			selection_info.selection_flow_id = ?
		AND
			question.id = (
				SELECT 
					question_join.id
				FROM 
					selection_questionnaires AS question_join 
				INNER JOIN 
					job_information_selection_informations AS selection_info_join
				ON
					selection_info_join.id = question_join.selection_information_id
				WHERE 
					selection_info.selection_type = selection_info_join.selection_type
				ORDER BY question_join.created_at ASC
				LIMIT 1
			)
		ORDER BY id ASC
		`,
		selectionFlowID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}

func (repo *SelectionQuestionnaireRepositoryImpl) GetUnanswerdByJobSeekerUUID(jobSeekerUUID uuid.UUID) ([]*entity.SelectionQuestionnaire, error) {
	var (
		list []*entity.SelectionQuestionnaire
	)

	err := repo.executer.Select(
		repo.Name+".GetUnanswerdByJobSeekerUUID",
		&list, `
		SELECT 
			question.*, selection_info.selection_type,
			job_info.uuid AS job_information_uuid,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
				THEN task_group.external_job_information_title
				ELSE job_info.title
			END AS title,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
				THEN task_group.external_company_name
				ELSE enterprise.company_name
			END AS company_name
		FROM 
			selection_questionnaires AS question
		INNER JOIN 
			job_information_selection_informations AS selection_info
		ON
			selection_info.id = question.selection_information_id
		INNER JOIN 
			job_informations AS job_info
		ON
			question.job_information_id = job_info.id
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN 
			job_seekers AS seeker
		ON
			question.job_seeker_id = seeker.id
		INNER JOIN
			task_groups AS task_group
		ON
			question.job_seeker_id = task_group.job_seeker_id AND
			question.job_information_id = task_group.job_information_id
		WHERE
			seeker.uuid = ?
		AND
			question.is_answer = FALSE
		AND
			question.id = (
				SELECT question_join.id
				FROM selection_questionnaires AS question_join 
				WHERE 
					question.job_seeker_id = question_join.job_seeker_id AND
					question.job_information_id = question_join.job_information_id
				ORDER BY question_join.id DESC
				LIMIT 1
			)
		ORDER BY id DESC
		`,
		jobSeekerUUID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}

/****************************************************************************************/
