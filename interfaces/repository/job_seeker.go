package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type JobSeekerRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerRepository {
	return &JobSeekerRepositoryImpl{
		Name:     "JobSeekerRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
//求職者の作成
func (repo *JobSeekerRepositoryImpl) Create(jobSeeker *entity.JobSeeker) error {
	jobSeeker.UUID = utility.CreateUUID()

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO job_seekers (
			uuid,
			agent_id,
			agent_staff_id,
			line_id,
			user_status,
			last_name,
			first_name,
			last_furigana,
			first_furigana,
			gender,
			gender_remarks,
			birthday,
			spouse,
			support_obligation,
			dependents,
			phone_number,
			email,
			emergency_phone_number,
			post_code,
			prefecture,
			address,
			address_furigana,
			state_of_employment,
			job_summary,
			history_supplement,
			research_content,
			join_company_period,
			job_change,
			annual_income,
			desired_annual_income,
			transfer,
			transfer_requirement,
			short_resignation,
			short_resignation_remarks,
			medical_history,
			nationality,
			appearance,
			communication,
			thinking,
			recommendation_profile,
			candid_profile,
			secret_memo,
			job_hunting_state,
			recommend_reason,
			phase,
			interview_date,
			agreement,
			register_phase,
			study_category,
			word_skill,
			excel_skill,
			power_point_skill,
			inflow_channel_id,
			nationality_remarks,
			medical_history_remarks,
			acceptance_points,
			activity_memo,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
		jobSeeker.UUID,
		jobSeeker.AgentID,
		jobSeeker.AgentStaffID,
		"", // LINE IDはブランクで登録
		jobSeeker.UserStatus,
		jobSeeker.LastName,
		jobSeeker.FirstName,
		jobSeeker.LastFurigana,
		jobSeeker.FirstFurigana,
		jobSeeker.Gender,
		jobSeeker.GenderRemarks,
		jobSeeker.Birthday,
		jobSeeker.Spouse,
		jobSeeker.SupportObligation,
		jobSeeker.Dependents,
		jobSeeker.PhoneNumber,
		jobSeeker.Email,
		jobSeeker.EmergencyPhoneNumber,
		jobSeeker.PostCode,
		jobSeeker.Prefecture,
		jobSeeker.Address,
		jobSeeker.AddressFurigana,
		jobSeeker.StateOfEmployment,
		jobSeeker.JobSummary,
		jobSeeker.HistorySupplement,
		jobSeeker.ResearchContent,
		jobSeeker.JoinCompanyPeriod,
		jobSeeker.JobChange,
		jobSeeker.AnnualIncome,
		jobSeeker.DesiredAnnualIncome,
		jobSeeker.Transfer,
		jobSeeker.TransferRequirement,
		jobSeeker.ShortResignation,
		jobSeeker.ShortResignationRemarks,
		jobSeeker.MedicalHistory,
		jobSeeker.Nationality,
		jobSeeker.Appearance,
		jobSeeker.Communication,
		jobSeeker.Thinking,
		jobSeeker.RecommendationProfile,
		jobSeeker.CandidProfile,
		jobSeeker.SecretMemo,
		jobSeeker.JobHuntingState,
		jobSeeker.RecommendReason,
		jobSeeker.Phase,
		jobSeeker.InterviewDate,
		false, // 同意の有無はfalseで登録
		jobSeeker.RegisterPhase,
		jobSeeker.StudyCategory,
		jobSeeker.WordSkill,
		jobSeeker.ExcelSkill,
		jobSeeker.PowerPointSkill,
		jobSeeker.InflowChannelID,
		jobSeeker.NationalityRemarks,
		jobSeeker.MedicalHistoryRemarks,
		jobSeeker.AcceptancePoints,
		"",
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	jobSeeker.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
// agreementは更新しない
func (repo *JobSeekerRepositoryImpl) Update(id uint, jobSeeker *entity.JobSeeker) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE job_seekers
		SET
		agent_staff_id = ?,
		user_status = ?,
		last_name = ?,
		first_name = ?,
		last_furigana = ?,
		first_furigana = ?,
		gender = ?,
		gender_remarks = ?,
		birthday = ?,
		spouse = ?,
		support_obligation = ?,
		dependents = ?,
		phone_number = ?,
		email = ?,
		emergency_phone_number = ?,
		post_code = ?,
		prefecture = ?,
		address	= ?,
		address_furigana = ?,
		state_of_employment = ?,
		job_summary = ?,
		history_supplement = ?,
		research_content = ?,
		join_company_period = ?,
		job_change = ?,
		annual_income = ?,
		desired_annual_income = ?,
		transfer = ?,
		transfer_requirement = ?,
		short_resignation = ?,
		short_resignation_remarks = ?,
		medical_history = ?,
		nationality = ?,
		appearance = ?,
		communication = ?,
		thinking = ?,
		recommendation_profile = ?,
		candid_profile = ?,
		secret_memo = ?,
		job_hunting_state = ?,
		recommend_reason = ?,
		phase = ?,
		register_phase = ?,
		study_category = ?,
		word_skill = ?,
		excel_skill = ?,
		power_point_skill = ?,
		inflow_channel_id = ?,
		nationality_remarks = ?,
		medical_history_remarks = ?,
		acceptance_points = ?,
		updated_at = ?
		WHERE 
			id = ?
		`,
		jobSeeker.AgentStaffID,
		jobSeeker.UserStatus,
		jobSeeker.LastName,
		jobSeeker.FirstName,
		jobSeeker.LastFurigana,
		jobSeeker.FirstFurigana,
		jobSeeker.Gender,
		jobSeeker.GenderRemarks,
		jobSeeker.Birthday,
		jobSeeker.Spouse,
		jobSeeker.SupportObligation,
		jobSeeker.Dependents,
		jobSeeker.PhoneNumber,
		jobSeeker.Email,
		jobSeeker.EmergencyPhoneNumber,
		jobSeeker.PostCode,
		jobSeeker.Prefecture,
		jobSeeker.Address,
		jobSeeker.AddressFurigana,
		jobSeeker.StateOfEmployment,
		jobSeeker.JobSummary,
		jobSeeker.HistorySupplement,
		jobSeeker.ResearchContent,
		jobSeeker.JoinCompanyPeriod,
		jobSeeker.JobChange,
		jobSeeker.AnnualIncome,
		jobSeeker.DesiredAnnualIncome,
		jobSeeker.Transfer,
		jobSeeker.TransferRequirement,
		jobSeeker.ShortResignation,
		jobSeeker.ShortResignationRemarks,
		jobSeeker.MedicalHistory,
		jobSeeker.Nationality,
		jobSeeker.Appearance,
		jobSeeker.Communication,
		jobSeeker.Thinking,
		jobSeeker.RecommendationProfile,
		jobSeeker.CandidProfile,
		jobSeeker.SecretMemo,
		jobSeeker.JobHuntingState,
		jobSeeker.RecommendReason,
		jobSeeker.Phase,
		// jobSeeker.InterviewDate, // 面談調整時間は更新しない　将来的には求職者テーブルから面談日時カラムを消す
		jobSeeker.RegisterPhase,
		jobSeeker.StudyCategory,
		jobSeeker.WordSkill,
		jobSeeker.ExcelSkill,
		jobSeeker.PowerPointSkill,
		jobSeeker.InflowChannelID,
		jobSeeker.NationalityRemarks,
		jobSeeker.MedicalHistoryRemarks,
		jobSeeker.AcceptancePoints,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 開発環境で実行時にテストユーザーの名前で更新する
func (repo *JobSeekerRepositoryImpl) UpdateForDev(id uint) error {
	var (
		lastName      = "テスト"
		firstName     = fmt.Sprintf("ユーザー%v", id)
		lastFurigana  = "テスト"
		firstFurigana = fmt.Sprintf("ユーザー%v", id)
		phoneNumber   = "080-0000-0000"
	)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateForDev",
		`
		UPDATE job_seekers
		SET
			last_name = ?,
			first_name = ?,
			last_furigana = ?,
			first_furigana = ?,
			phone_number = ?
		WHERE 
			id = ?
		`,
		lastName,
		firstName,
		lastFurigana,
		firstFurigana,
		phoneNumber,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 個人情報同意の更新
func (repo *JobSeekerRepositoryImpl) UpdateActivityMemo(id uint, activityMemo string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateActivityMemo",
		`
		UPDATE job_seekers
		SET
		activity_memo  = ?,
		updated_at = ?
		WHERE 
			id = ?
		`,
		activityMemo,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// phaseの更新
func (repo *JobSeekerRepositoryImpl) UpdatePhase(id uint, phase null.Int) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdatePhase",
		`
		UPDATE job_seekers
		SET
			phase  = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		phase,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// phaseの更新
func (repo *JobSeekerRepositoryImpl) UpdateInterviewDate(id uint, interviewDate time.Time) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateInterviewDate",
		`
		UPDATE job_seekers
		SET
			interview_date = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		interviewDate,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// phaseの更新
func (repo *JobSeekerRepositoryImpl) UpdateStaffID(id uint, staffID null.Int) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateStaffID",
		`
		UPDATE job_seekers
		SET
			agent_staff_id = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		staffID,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 個人情報同意の更新
func (repo *JobSeekerRepositoryImpl) UpdateAgreement(id uint, agreement bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgreement",
		`
		UPDATE job_seekers
		SET
		agreement = ?,
		updated_at = ?
		WHERE 
			id = ?
		`,
		agreement,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// created_atを更新する処理 *csvインポートでエントリー日が入力されている場合に使用
func (repo *JobSeekerRepositoryImpl) UpdateCreatedAt(id uint, createdAt time.Time) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateCreatedAt",
		`
		UPDATE job_seekers
		SET
		created_at = ?
		WHERE 
			id = ?
		`,
		createdAt,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 応募承諾のポイントを更新する処理 *求職者相談時に使用
func (repo *JobSeekerRepositoryImpl) UpdateAcceptancePoints(id uint, acceptancePoints string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateCreatedAt",
		`
		UPDATE job_seekers
		SET
		  acceptance_points = ?
		WHERE 
			id = ?
		`,
		acceptancePoints,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// LINEIDを更新
func (repo *JobSeekerRepositoryImpl) UpdateLineIDByUUID(uuid uuid.UUID, lineID string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateLineIDByUUID",
		`
			UPDATE 
				job_seekers
			SET 
				line_id = ?
			WHERE 
				uuid = ?
		`,
		lineID,
		uuid,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// ゲストログイン用のパスワードを更新 リセットトークンを空にする
func (repo *JobSeekerRepositoryImpl) UpdatePassword(uuid uuid.UUID, password string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdatePassword",
		`
			UPDATE 
				job_seekers
			SET 
				password = ?,
				reset_password_token = ""
			WHERE 
				uuid = ?
			`,
		password,
		uuid,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// ゲストログイン用のパスワードリセット時のトークンを更新
func (repo *JobSeekerRepositoryImpl) UpdateResetPasswordToken(uuid uuid.UUID, resetPasswordToken string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateResetPasswordToken",
		`
			UPDATE 
				job_seekers
			SET 
				reset_password_token = ?
			WHERE 
				uuid = ?
			`,
		resetPasswordToken,
		uuid,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// uuidを使って電話番号を更新（LPでのみ使用）
func (repo *JobSeekerRepositoryImpl) UpdatePhoneNumberByUUID(uuid uuid.UUID, phoneNumber string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdatePhoneNumberByUUID",
		`
			UPDATE 
				job_seekers
			SET 
				phone_number = ?,
				updated_at = ?
			WHERE 
				uuid = ?
		`,
		phoneNumber,
		time.Now().In(time.UTC),
		uuid,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// uuidを使って希望年収を更新（LPでのみ使用）
func (repo *JobSeekerRepositoryImpl) UpdateDesiredIncomeByUUID(uuid uuid.UUID, desiredIncome null.Int) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateDesiredIncomeByUUID",
		`
			UPDATE 
				job_seekers
			SET 
				desired_annual_income = ?,
				updated_at = ?
			WHERE 
				uuid = ?
		`,
		desiredIncome,
		time.Now().In(time.UTC),
		uuid,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// おすすめ求人の閲覧権限の更新
func (repo *JobSeekerRepositoryImpl) UpdateCanViewMatchingJob(id uint, canView bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateCanViewMatchingJob",
		`
			UPDATE 
				job_seekers
			SET 
				can_view_matching_job = ?,
				updated_at = ?
			WHERE 
				id = ?
		`,
		canView,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM job_seekers
		WHERE id = ?
		`, id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 単数取得
//
func (repo *JobSeekerRepositoryImpl) FindByID(id uint) (*entity.JobSeeker, error) {
	var (
		jobSeeker entity.JobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&jobSeeker, `
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email,
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number, chat_group.line_active,
			agent.agent_name, task_group.interview_date, task_group.first_interview_date,
			IFNULL(inflow_channel.channel_name, '') AS channel_name,
			IFNULL(questionnaire.question, '') AS question
		FROM 
			job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		INNER JOIN
			interview_task_groups AS task_group
		ON
			seeker.id = task_group.job_seeker_id
		INNER JOIN
			chat_group_with_job_seekers AS chat_group
		ON
			seeker.id = chat_group.job_seeker_id
		LEFT OUTER JOIN
			initial_questionnaires AS questionnaire
		ON
			seeker.id = questionnaire.job_seeker_id
		LEFT OUTER JOIN
			agent_inflow_channel_options AS inflow_channel
		ON
		seeker.inflow_channel_id = inflow_channel.id
		WHERE
			seeker.id = ?
		LIMIT 1
		`, id)
	if err != nil {
		return nil, err
	}

	return &jobSeeker, nil
}

func (repo *JobSeekerRepositoryImpl) FindByUUID(jobSeekerUUID uuid.UUID) (*entity.JobSeeker, error) {
	var (
		jobSeeker entity.JobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByUUID",
		&jobSeeker, `
		SELECT 
			*
		FROM 
			job_seekers
		WHERE
			uuid = ?
		LIMIT 1
		`,
		jobSeekerUUID,
	)

	if err != nil {
		return nil, err
	}

	return &jobSeeker, nil
}

func (repo *JobSeekerRepositoryImpl) FindByTaskGroupUUID(taskGroupUUID uuid.UUID) (*entity.JobSeeker, error) {
	var (
		jobSeeker entity.JobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByTaskGroupUUID",
		&jobSeeker, `
		SELECT 
			seeker.*
		FROM 
			job_seekers AS seeker
		INNER JOIN
			task_groups AS task_group
		ON
			seeker.id = task_group.job_seeker_id
		WHERE
			task_group.uuid = ?
		LIMIT 1
		`,
		taskGroupUUID,
	)

	if err != nil {
		return nil, err
	}

	return &jobSeeker, nil
}

func (repo *JobSeekerRepositoryImpl) FindByAgentIDAndLineID(agentID uint, lineID string) (*entity.JobSeeker, error) {
	var (
		jobSeeker entity.JobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByAgentIDAndLineID",
		&jobSeeker, `
		SELECT 
			*
		FROM 
			job_seekers
		WHERE
			agent_id = ? AND
			line_id = ?
		LIMIT 1
		`, agentID, lineID,
	)

	if err != nil {
		return nil, err
	}

	return &jobSeeker, nil
}

// 送客の重複登録判定（System管理の求職者内で重複チェックするため）
func (repo *JobSeekerRepositoryImpl) FindByNameAndPhoneNumberBySystemAgent(firstName, lastName, firstFurigana, lastFurigana, phoneNumber string) (*entity.JobSeeker, error) {
	var (
		jobSeeker entity.JobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByNameAndPhoneNumberBySystemAgent",
		&jobSeeker, `
		SELECT 
			seeker.*
		FROM 
			job_seekers AS seeker
		WHERE
			agent_id = 1 AND
			seeker.first_name = ? AND
			seeker.last_name = ? AND
			seeker.first_furigana = ? AND
			seeker.last_furigana = ? AND
			seeker.phone_number = ?
		LIMIT 1
		`,
		firstName,
		lastName,
		firstFurigana,
		lastFurigana,
		phoneNumber,
	)

	if err != nil {
		return nil, err
	}

	return &jobSeeker, nil
}

// ハイフンの有無に関わらず取得するためにREPLACEを使用
func (repo *JobSeekerRepositoryImpl) FindByPhoneNumberForLP(phoneNumber string, agentID uint) (*entity.JobSeeker, error) {
	var (
		jobSeeker entity.JobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByPhoneNumberForLP",
		&jobSeeker, `
		SELECT 
			*
		FROM 
			job_seekers
		WHERE
			REPLACE(phone_number, '-', '') = ?
		AND
			agent_id = ?
		ORDER BY 
			id DESC
		LIMIT 1
		`,
		phoneNumber, agentID,
	)

	if err != nil {
		return nil, err
	}

	return &jobSeeker, nil
}

func (repo *JobSeekerRepositoryImpl) FindByEmailForLP(email string, agentID uint) (*entity.JobSeeker, error) {
	var (
		jobSeeker entity.JobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByEmailForLP",
		&jobSeeker, `
		SELECT 
			*
		FROM 
			job_seekers
		WHERE
			email = ?
		AND
			agent_id = ?
		ORDER BY 
			id DESC
		LIMIT 1
		`,
		email, agentID,
	)

	if err != nil {
		return nil, err
	}

	return &jobSeeker, nil
}

func (repo *JobSeekerRepositoryImpl) FindByResetPasswordTokenForLP(resetPasswordToken string, agentID uint) (*entity.JobSeeker, error) {
	var (
		jobSeeker entity.JobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByResetPasswordTokenForLP",
		&jobSeeker, `
		SELECT 
			job_seekers.id, job_seekers.uuid
		FROM 
			job_seekers
		WHERE
			reset_password_token = ?
		AND
			agent_id = ?
		LIMIT 1
		`,
		resetPasswordToken, agentID,
	)

	if err != nil {
		return nil, err
	}

	return &jobSeeker, nil
}

/****************************************************************************************/
/// 複数取得
//
// agentIDから求職者一覧を取得
func (repo *JobSeekerRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&jobSeekerList, `
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
			IFNULL(inflow_channel.channel_name, '') AS channel_name,
			agent.agent_name, IFNULL(questionnaire.question, '') AS question
		FROM 
			job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		LEFT OUTER JOIN
			initial_questionnaires AS questionnaire
		ON
			seeker.id = questionnaire.job_seeker_id
		LEFT OUTER JOIN
			agent_inflow_channel_options AS inflow_channel
		ON
			seeker.inflow_channel_id = inflow_channel.id
		WHERE
			seeker.agent_id = ?
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

// agentIDから3日以内に登録された求職者一覧を取得
func (repo *JobSeekerRepositoryImpl) GetByAgentIDWithinTwoDays(agentID uint) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDWithinTwoDays",
		&jobSeekerList, `
		SELECT
			seeker.*
		FROM
			job_seekers AS seeker
		WHERE
			seeker.agent_id = ?
		AND
			seeker.created_at >= DATE_SUB(NOW(), INTERVAL 2 DAY)
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

// 求職者IDリストに合致する求職者の一覧を取得
func (repo *JobSeekerRepositoryImpl) GetByIDList(idList []uint) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)

	if len(idList) < 1 {
		return jobSeekerList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
			chat_group.line_active,
			agent.agent_name, IFNULL(questionnaire.question, '') AS question
		FROM 
			job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		INNER JOIN
			chat_group_with_job_seekers AS chat_group
		ON
			seeker.id = chat_group.job_seeker_id
		LEFT OUTER JOIN
			initial_questionnaires AS questionnaire
		ON
			seeker.id = questionnaire.job_seeker_id
		WHERE
			seeker.id IN (%s)
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByIDList",
		&jobSeekerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

// 求職者IDリストに合致する求職者の一覧を取得
func (repo *JobSeekerRepositoryImpl) GetByIDListForTask(idList []uint) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)

	if len(idList) < 1 {
		return jobSeekerList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			seeker.id,
			seeker.uuid,
			seeker.last_name, 
			seeker.first_name, 
			seeker.last_furigana, 
			seeker.first_furigana,
			seeker.agent_id,
			seeker.agent_staff_id,
			IFNULL(staff.staff_name, '') AS staff_name,
			agent.agent_name
		FROM 
			job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		WHERE
			seeker.id IN (%s)
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByIDListForTask",
		&jobSeekerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

// agentIDから求職者一覧を取得
func (repo *JobSeekerRepositoryImpl) GetByAgentIDAndPage(agentID uint) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)
	/**
	* キーになるagent_staff_idがnullの場合があるので「LEFT OUTER JOIN」でnullでもテーブル結合できるように
	* nullでテーブル結合させるとstaff_nameもnullになり型エラーになるため、「IFNULL」で空文字に変換
	**/

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDAndPage",
		&jobSeekerList,
		`
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
			agent.agent_name,
			task_group.interview_date, task_group.first_interview_date,
			IFNULL(questionnaire.question, '') AS question
		FROM 
			job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		INNER JOIN
			interview_task_groups AS task_group
		ON
			seeker.id = task_group.job_seeker_id
		LEFT OUTER JOIN
			initial_questionnaires AS questionnaire
		ON
			seeker.id = questionnaire.job_seeker_id
		WHERE
			seeker.agent_id = ?
		ORDER BY seeker.id DESC
	`, agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

// agentIDから求職者一覧を取得　アライアンス含む
func (repo *JobSeekerRepositoryImpl) GetByAgentIDAndFreeWordAndAlliance(agentID uint, freeWord string, agentAllianceList []*entity.AgentAlliance) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList       []*entity.JobSeeker
		allianceAgentIDList []uint
		allianceQuery       string
		freeWordQuery       string
	)

	// アライアンスを締結している場合、申請先or申請元が他社エージェントのIDを取得
	if len(agentAllianceList) > 0 {
		// アライアンス締結済みの他社エージェントのIDを取得
		for _, alliance := range agentAllianceList {
			if alliance.Agent1ID == agentID {
				allianceAgentIDList = append(allianceAgentIDList, alliance.Agent2ID)
			} else if alliance.Agent2ID == agentID {
				allianceAgentIDList = append(allianceAgentIDList, alliance.Agent1ID)
			}
		}

		allianceQuery = fmt.Sprintf(
			"OR ((seeker.phase = 5 OR seeker.phase = 6) AND seeker.register_phase = 0 AND seeker.agent_id IN (%s))",
			strings.Trim(strings.Join(strings.Fields(fmt.Sprint(allianceAgentIDList)), ","), "[]"),
		)
	}

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {
			freeWordForLike := "%" + freeWord + "%"

			freeWordQuery = fmt.Sprintf(`
				AND (	
					CONCAT(seeker.last_name, seeker.first_name) LIKE '%s' OR 
					CONCAT(seeker.last_furigana, seeker.first_furigana) LIKE '%s' OR
					seeker.phone_number = '%s'
				)
			`, freeWordForLike, freeWordForLike, freeWord)
		}

	} else {
		// ハイフンなしの場合数字として扱われるためハイフンなしはこちらに入れる
		freeWordQuery = fmt.Sprintf(`
			AND 	
				seeker.id = %d OR
				REPLACE(seeker.phone_number, '-', '') = '%s'
				`, freeWordInt, freeWord)
	}
	query := fmt.Sprintf(
		`
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
			agent.agent_name,
			task_group.interview_date, task_group.first_interview_date,
			IFNULL(questionnaire.question, '') AS question
		FROM 
			job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		INNER JOIN
			interview_task_groups AS task_group
		ON
			seeker.id = task_group.job_seeker_id
		LEFT OUTER JOIN
			initial_questionnaires AS questionnaire
		ON
			seeker.id = questionnaire.job_seeker_id
		WHERE
		(
			(seeker.agent_id = %v AND seeker.phase IN (3, 4, 5, 6))
			%s
		)
			%s
		ORDER BY seeker.id DESC
	`, agentID,
		allianceQuery,
		freeWordQuery,
	)

	err = repo.executer.Select(
		repo.Name+".GetByAgentIDAndFreeWordAndAlliance",
		&jobSeekerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

// agentStaffIDとサービス終了ではないから求職者一覧を取得
func (repo *JobSeekerRepositoryImpl) GetByAgentStaffIDAndNotRelease(agentID uint) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentStaffIDAndNotRelease",
		&jobSeekerList, `
		SELECT 
			seeker.id
		FROM 
			job_seekers AS seeker
		WHERE
			seeker.agent_staff_id = ?
		AND
			seeker.phase NOT IN (7, 8, 9)
		ORDER BY seeker.id DESC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

// 引数の値で重複する求職者を取得する
func (repo *JobSeekerRepositoryImpl) GetDuplicateByNameAndFuriganaAndEmailAndPhoneNumber(agentID uint, lastName, firstName, lastFurigana, firstFurigana, email, phoneNumber string) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetDuplicateByNameAndFuriganaAndEmailAndPhoneNumber",
		&jobSeekerList, `
		SELECT 
			id, last_name, first_name, last_furigana, first_furigana, email, phone_number
		FROM 
			job_seekers
		WHERE
			agent_id = ? AND
			last_name = ? AND
			first_name = ? AND
			last_furigana = ? AND
			first_furigana = ? AND
			email = ? AND
			phone_number = ?
		ORDER BY id ASC
		`,
		agentID, lastName, firstName, lastFurigana, firstFurigana, email, phoneNumber,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

/****************************************************************************************/
/// 絞り込み検索 API
//

// agentIDから求職者一覧を取得
func (repo *JobSeekerRepositoryImpl) GetByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.JobSeeker, error) {

	var (
		jobSeekerList []*entity.JobSeeker
		freeWordQuery string
	)

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {
			freeWordForLike := "%" + freeWord + "%"

			freeWordQuery = fmt.Sprintf(`
				AND (	
					CONCAT(seeker.last_name, seeker.first_name) LIKE '%s' OR 
					CONCAT(seeker.last_furigana, seeker.first_furigana) LIKE '%s' OR
					seeker.phone_number = '%s'
				)
			`, freeWordForLike, freeWordForLike, freeWord)
		}

	} else {
		// ハイフンなしの場合数字として扱われるためハイフンなしはこちらに入れる
		freeWordQuery = fmt.Sprintf(`
			AND 	
				seeker.id = %d OR
				REPLACE(seeker.phone_number, '-', '') = '%s'
				`, freeWordInt, freeWord)
	}

	// クエリをまとめる
	query := fmt.Sprintf(`
	SELECT 
		seeker.*, 
		IFNULL(staff.staff_name, '') AS staff_name, 
		IFNULL(staff.email, '') AS staff_email, 
		IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
		agent.agent_name,
		task_group.interview_date, task_group.first_interview_date,
		IFNULL(questionnaire.question, '') AS question
	FROM 
		job_seekers AS seeker
	LEFT OUTER JOIN
		agent_staffs AS staff
	ON
		seeker.agent_staff_id = staff.id
	INNER JOIN
		agents AS agent
	ON
		seeker.agent_id = agent.id
	INNER JOIN
		interview_task_groups AS task_group
	ON
		seeker.id = task_group.job_seeker_id
	LEFT OUTER JOIN
		initial_questionnaires AS questionnaire
	ON
		seeker.id = questionnaire.job_seeker_id
	WHERE
		seeker.agent_id = %v
	%s
		ORDER BY seeker.id DESC`,
		agentID,
		freeWordQuery)

	// クエリを実行
	err = repo.executer.Select(
		repo.Name+".GetByAgentIDAndFreeWord",
		&jobSeekerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

// agentIDから求職者一覧を取得
func (repo *JobSeekerRepositoryImpl) GetByOtherAgentIDAndFreeWord(agentID uint, freeWord string, agentAllianceList []*entity.AgentAlliance) ([]*entity.JobSeeker, error) {

	var (
		jobSeekerList       []*entity.JobSeeker
		freeWordQuery       string
		allianceAgentIDList []uint
	)

	// アライアンスエージェントのクエリを作成
	if len(agentAllianceList) < 1 {
		// アライアンスエージェントがいないかつ、他社求人のみの場合は、空を返す
		return nil, nil
	}

	// アライアンス締結済みの他社エージェントのIDを取得
	for _, alliance := range agentAllianceList {
		if alliance.Agent1ID == agentID {
			allianceAgentIDList = append(allianceAgentIDList, alliance.Agent2ID)
		} else if alliance.Agent2ID == agentID {
			allianceAgentIDList = append(allianceAgentIDList, alliance.Agent1ID)
		}
	}

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {
			freeWordForLike := "%" + freeWord + "%"

			freeWordQuery = fmt.Sprintf(`
				AND (	
					CONCAT(seeker.last_name, seeker.first_name) LIKE '%s' OR 
					CONCAT(seeker.last_furigana, seeker.first_furigana) LIKE '%s' OR
					seeker.phone_number = '%s'
				)
			`, freeWordForLike, freeWordForLike, freeWord)
		}

	} else {
		// ハイフンなしの場合数字として扱われるためハイフンなしはこちらに入れる
		freeWordQuery = fmt.Sprintf(`
			AND 	
				seeker.id = %d OR
				REPLACE(seeker.phone_number, '-', '') = '%s'
				`, freeWordInt, freeWord)
	}

	// クエリをまとめる
	query := fmt.Sprintf(`
	SELECT 
		seeker.*, 
		IFNULL(staff.staff_name, '') AS staff_name, 
		IFNULL(staff.email, '') AS staff_email, 
		IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
		agent.agent_name, IFNULL(questionnaire.question, '') AS question
	FROM 
		job_seekers AS seeker
	LEFT OUTER JOIN
		agent_staffs AS staff
	ON
		seeker.agent_staff_id = staff.id
	INNER JOIN
		agents AS agent
	ON
		seeker.agent_id = agent.id
	LEFT OUTER JOIN
		initial_questionnaires AS questionnaire
	ON
		seeker.id = questionnaire.job_seeker_id
	WHERE
		seeker.agent_staff_id IS NOT NULL
	AND
	 	seeker.agent_id IN (%s)
	AND
		(seeker.phase = 5 OR seeker.phase = 6)
	AND 
		seeker.register_phase = 0
		%s
	ORDER BY seeker.id DESC`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(allianceAgentIDList)), ","), "[]"),
		freeWordQuery)

	// クエリを実行
	err = repo.executer.Select(
		repo.Name+".GetByOtherAgentIDAndFreeWord",
		&jobSeekerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

/****************************************************************************************/
// 求人検索→求職者検索 API
//
// 全てのアクティブな求職者
func (repo *JobSeekerRepositoryImpl) GetActiveAllAndFreeWord(agentID uint, freeWord string) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
		freeWordQuery string
	)

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {
			freeWordForLike := "%" + freeWord + "%"

			freeWordQuery = fmt.Sprintf(`
				AND (	
					CONCAT(seeker.last_name, seeker.first_name) LIKE '%s' OR 
					CONCAT(seeker.last_furigana, seeker.first_furigana) LIKE '%s' OR
					seeker.phone_number = '%s'
				)
			`, freeWordForLike, freeWordForLike, freeWord)
		}

	} else {
		// ハイフンなしの場合数字として扱われるためハイフンなしはこちらに入れる
		freeWordQuery = fmt.Sprintf(`
			AND 	
				seeker.id = %d OR
				REPLACE(seeker.phone_number, '-', '') = '%s'
				`, freeWordInt, freeWord)
	}

	/**
	* キーになるagent_staff_idがnullの場合があるので「LEFT OUTER JOIN」でnullでもテーブル結合できるように
	* nullでテーブル結合させるとstaff_nameもnullになり型エラーになるため、「IFNULL」で空文字に変換
	**/

	query := fmt.Sprintf(`
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
			agent.agent_name,
			task_group.interview_date, task_group.first_interview_date,
			IFNULL(questionnaire.question, '') AS question,
			IFNULL(inflow_channel.channel_name, '') AS channel_name
		FROM 
			job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		INNER JOIN
			interview_task_groups AS task_group
		ON
			seeker.id = task_group.job_seeker_id
		LEFT OUTER JOIN
			initial_questionnaires AS questionnaire
		ON
			seeker.id = questionnaire.job_seeker_id
		LEFT OUTER JOIN
			agent_inflow_channel_options AS inflow_channel
		ON
			seeker.inflow_channel_id = inflow_channel.id
		WHERE (
			(
				seeker.agent_id = %v
				AND
				seeker.phase IN (3, 4, 5, 6)
			) OR (
				seeker.agent_id NOT IN(%v)
				AND
				seeker.phase IN (5, 6)
				AND 
				seeker.register_phase = 0
			)
		)
		%s
		ORDER BY 
			seeker.id DESC
	`,
		agentID,
		agentID,
		freeWordQuery,
	)

	err = repo.executer.Select(
		repo.Name+".GetActiveAllAndFreeWord",
		&jobSeekerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

// 自社のアクティブな求職者
func (repo *JobSeekerRepositoryImpl) GetActiveOwnAndFreeWord(agentID uint, freeWord string) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
		freeWordQuery string
	)

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {
			freeWordForLike := "%" + freeWord + "%"

			freeWordQuery = fmt.Sprintf(`
				AND (	
					CONCAT(seeker.last_name, seeker.first_name) LIKE '%s' OR 
					CONCAT(seeker.last_furigana, seeker.first_furigana) LIKE '%s' OR
					seeker.phone_number = '%s'
				)
			`, freeWordForLike, freeWordForLike, freeWord)
		}

	} else {
		// ハイフンなしの場合数字として扱われるためハイフンなしはこちらに入れる
		freeWordQuery = fmt.Sprintf(`
			AND 	
				seeker.id = %d OR
				REPLACE(seeker.phone_number, '-', '') = '%s'
				`, freeWordInt, freeWord)
	}

	/**
	* キーになるagent_staff_idがnullの場合があるので「LEFT OUTER JOIN」でnullでもテーブル結合できるように
	* nullでテーブル結合させるとstaff_nameもnullになり型エラーになるため、「IFNULL」で空文字に変換
	**/

	query := fmt.Sprintf(`
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
			chat_group.line_active,
			agent.agent_name,
			task_group.interview_date, task_group.first_interview_date,
			IFNULL(questionnaire.question, '') AS question,
			IFNULL(inflow_channel.channel_name, '') AS channel_name
		FROM 
			job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		INNER JOIN
			interview_task_groups AS task_group
		ON
			seeker.id = task_group.job_seeker_id
		INNER JOIN
			chat_group_with_job_seekers AS chat_group
		ON
			seeker.id = chat_group.job_seeker_id
		LEFT OUTER JOIN
			initial_questionnaires AS questionnaire
		ON
			seeker.id = questionnaire.job_seeker_id
		LEFT OUTER JOIN
			agent_inflow_channel_options AS inflow_channel
		ON
			seeker.inflow_channel_id = inflow_channel.id
		WHERE 
			seeker.agent_id = %v
			%s
		AND
			seeker.phase IN (3, 4, 5, 6)
		ORDER BY 
			seeker.id DESC
	`, agentID, freeWordQuery)

	err = repo.executer.Select(
		repo.Name+".GetActiveOwnAndFreeWord",
		&jobSeekerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

/****************************************************************************************/
// ダッシュボード関連 API
//

// ・タスクはあるけどヨミ入力がない求職者
// ・書類選考以降の稼働タスクがない求職者
// ・ヨミフェーズが失注になった求職者
// ※フェーズがサービス終了の求職者は表示しない。
func (repo *JobSeekerRepositoryImpl) GetReleaseByStaffID(staffID uint) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetReleaseByStaffID",
		&jobSeekerList, `
		SELECT 
			seeker.*,
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number, 
			agent.agent_name, IFNULL(questionnaire.question, '') AS question
		FROM 
			job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		LEFT OUTER JOIN
			sales AS sale
		ON 
			seeker.id = sale.job_seeker_id
		LEFT OUTER JOIN
			task_groups AS task_group
		ON 
			sale.job_seeker_id = task_group.job_seeker_id
		AND	
			sale.job_information_id = task_group.job_information_id
		LEFT OUTER JOIN
			initial_questionnaires AS questionnaire
		ON
			seeker.id = questionnaire.job_seeker_id
		WHERE
			seeker.agent_staff_id = ?
		AND
			seeker.phase IN(4, 5, 6)
		AND (
			(
				sale.id IS NULL AND
				EXISTS (
					SELECT 1
					FROM tasks
					WHERE tasks.task_group_id = task_group.id
				)
			) OR
			NOT EXISTS (
				SELECT 1
				FROM tasks
				WHERE tasks.task_group_id = task_group.id
				AND tasks.phase_category >= 1
			) OR
			sale.accuracy = 5
		)
		ORDER BY
			FIELD(seeker.phase, 4, 6, 5)
		`, staffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

func (repo *JobSeekerRepositoryImpl) GetReleaseByAgentID(agentID uint) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetReleaseByAgentID",
		&jobSeekerList, `
		SELECT 
			seeker.*,
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number, 
			agent.agent_name, IFNULL(questionnaire.question, '') AS question
		FROM 
			job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		LEFT OUTER JOIN 
			sales AS sale
		ON 
			seeker.id = sale.job_seeker_id
		LEFT OUTER JOIN
			task_groups AS task_group
		ON 
			sale.job_seeker_id = task_group.job_seeker_id
		AND	
			sale.job_information_id = task_group.job_information_id
		LEFT OUTER JOIN
			initial_questionnaires AS questionnaire
		ON
			seeker.id = questionnaire.job_seeker_id
		WHERE
			seeker.agent_id = ?
		AND
			seeker.phase IN(4, 5, 6)
		AND (
			(
				sale.id IS NULL AND
				EXISTS (
					SELECT 1
					FROM tasks
					WHERE tasks.task_group_id = task_group.id
				)
			) OR
			NOT EXISTS (
				SELECT 1
				FROM tasks
				WHERE tasks.task_group_id = task_group.id
				AND tasks.phase_category >= 1
			) OR
			sale.accuracy = 5
		)
		ORDER BY
			FIELD(seeker.phase, 4, 6, 5)
		`, agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerRepositoryImpl) All() ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&jobSeekerList, `
			SELECT *
			FROM job_seekers
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobSeekerList, nil
}

/****************************************************************************************/
/// カウント
//
func (repo *JobSeekerRepositoryImpl) CountByEmail(email string) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+"CountNotificationByAgentID",
		&result, `
			SELECT 
				COUNT(*) AS count
			FROM 
				job_seekers AS seeker
			WHERE
				seeker.email = ?
    `,
		email,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}
