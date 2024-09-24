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

type SendingJobSeekerRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerRepository {
	return &SendingJobSeekerRepositoryImpl{
		Name:     "SendingJobSeekerRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
// 送客求職者の作成
func (repo *SendingJobSeekerRepositoryImpl) Create(sendingJobSeeker *entity.SendingJobSeeker) error {
	sendingJobSeeker.UUID = utility.CreateUUID()

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO sending_job_seekers (
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
			appearance_detail_of_truth,
			appearance_detail,
			communication,
			communication_detail_of_truth,
			communication_detail,
			thinking,
			thinking_detail_of_truth,
			thinking_detail,
			secret_memo,
			job_hunting_state,
			recommend_reason,
			phase,
			interview_date,
			agreement,
			study_category,
			word_skill,
			excel_skill,
			power_point_skill,
			public_memo,
			nationality_remarks,
			medical_history_remarks,
			acceptance_points,
			question,
			sending_customer_id,
			activity_memo,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?
		)`,
		sendingJobSeeker.UUID,
		sendingJobSeeker.AgentID,
		sendingJobSeeker.AgentStaffID,
		sendingJobSeeker.LineID,
		sendingJobSeeker.UserStatus,
		sendingJobSeeker.LastName,
		sendingJobSeeker.FirstName,
		sendingJobSeeker.LastFurigana,
		sendingJobSeeker.FirstFurigana,
		sendingJobSeeker.Gender,
		sendingJobSeeker.GenderRemarks,
		sendingJobSeeker.Birthday,
		sendingJobSeeker.Spouse,
		sendingJobSeeker.SupportObligation,
		sendingJobSeeker.Dependents,
		sendingJobSeeker.PhoneNumber,
		sendingJobSeeker.Email,
		sendingJobSeeker.EmergencyPhoneNumber,
		sendingJobSeeker.PostCode,
		sendingJobSeeker.Prefecture,
		sendingJobSeeker.Address,
		sendingJobSeeker.AddressFurigana,
		sendingJobSeeker.StateOfEmployment,
		sendingJobSeeker.JobSummary,
		sendingJobSeeker.HistorySupplement,
		sendingJobSeeker.ResearchContent,
		sendingJobSeeker.JoinCompanyPeriod,
		sendingJobSeeker.JobChange,
		sendingJobSeeker.AnnualIncome,
		sendingJobSeeker.DesiredAnnualIncome,
		sendingJobSeeker.Transfer,
		sendingJobSeeker.TransferRequirement,
		sendingJobSeeker.ShortResignation,
		sendingJobSeeker.ShortResignationRemarks,
		sendingJobSeeker.MedicalHistory,
		sendingJobSeeker.Nationality,
		sendingJobSeeker.Appearance,
		sendingJobSeeker.AppearanceDetailOfTruth,
		sendingJobSeeker.AppearanceDetail,
		sendingJobSeeker.Communication,
		sendingJobSeeker.CommunicationDetailOfTruth,
		sendingJobSeeker.CommunicationDetail,
		sendingJobSeeker.Thinking,
		sendingJobSeeker.ThinkingDetailOfTruth,
		sendingJobSeeker.ThinkingDetail,
		sendingJobSeeker.SecretMemo,
		sendingJobSeeker.JobHuntingState,
		sendingJobSeeker.RecommendReason,
		sendingJobSeeker.Phase,
		sendingJobSeeker.InterviewDate,
		sendingJobSeeker.Agreement,
		sendingJobSeeker.StudyCategory,
		sendingJobSeeker.WordSkill,
		sendingJobSeeker.ExcelSkill,
		sendingJobSeeker.PowerPointSkill,
		sendingJobSeeker.PublicMemo,
		sendingJobSeeker.NationalityRemarks,
		sendingJobSeeker.MedicalHistoryRemarks,
		sendingJobSeeker.AcceptancePoints,
		"", // 初回アンケートの質問・要望（Createの時は空文字）
		sendingJobSeeker.SendingCustomerID,
		"",
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	sendingJobSeeker.ID = uint(lastID)
	return nil
}

// agreementは更新しない
func (repo *SendingJobSeekerRepositoryImpl) Update(sendingJobSeeker *entity.SendingJobSeeker, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE sending_job_seekers
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
			appearance_detail_of_truth = ?,
			appearance_detail = ?,
			communication = ?,
			communication_detail_of_truth = ?,
			communication_detail = ?,
			thinking = ?,
			thinking_detail_of_truth = ?,
			thinking_detail = ?,
			secret_memo = ?,
			job_hunting_state = ?,
			recommend_reason = ?,
			phase = ?,
			study_category = ?,
			word_skill = ?,
			excel_skill = ?,
			power_point_skill = ?,
			public_memo = ?,
			nationality_remarks = ?,
			medical_history_remarks = ?,
			acceptance_points = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		sendingJobSeeker.AgentStaffID,
		sendingJobSeeker.UserStatus,
		sendingJobSeeker.LastName,
		sendingJobSeeker.FirstName,
		sendingJobSeeker.LastFurigana,
		sendingJobSeeker.FirstFurigana,
		sendingJobSeeker.Gender,
		sendingJobSeeker.GenderRemarks,
		sendingJobSeeker.Birthday,
		sendingJobSeeker.Spouse,
		sendingJobSeeker.SupportObligation,
		sendingJobSeeker.Dependents,
		sendingJobSeeker.PhoneNumber,
		sendingJobSeeker.Email,
		sendingJobSeeker.EmergencyPhoneNumber,
		sendingJobSeeker.PostCode,
		sendingJobSeeker.Prefecture,
		sendingJobSeeker.Address,
		sendingJobSeeker.AddressFurigana,
		sendingJobSeeker.StateOfEmployment,
		sendingJobSeeker.JobSummary,
		sendingJobSeeker.HistorySupplement,
		sendingJobSeeker.ResearchContent,
		sendingJobSeeker.JoinCompanyPeriod,
		sendingJobSeeker.JobChange,
		sendingJobSeeker.AnnualIncome,
		sendingJobSeeker.DesiredAnnualIncome,
		sendingJobSeeker.Transfer,
		sendingJobSeeker.TransferRequirement,
		sendingJobSeeker.ShortResignation,
		sendingJobSeeker.ShortResignationRemarks,
		sendingJobSeeker.MedicalHistory,
		sendingJobSeeker.Nationality,
		sendingJobSeeker.Appearance,
		sendingJobSeeker.AppearanceDetailOfTruth,
		sendingJobSeeker.AppearanceDetail,
		sendingJobSeeker.Communication,
		sendingJobSeeker.CommunicationDetailOfTruth,
		sendingJobSeeker.CommunicationDetail,
		sendingJobSeeker.Thinking,
		sendingJobSeeker.ThinkingDetailOfTruth,
		sendingJobSeeker.ThinkingDetail,
		sendingJobSeeker.SecretMemo,
		sendingJobSeeker.JobHuntingState,
		sendingJobSeeker.RecommendReason,
		sendingJobSeeker.Phase,
		// sendingJobSeeker.InterviewDate, // 面談調整時間は更新しない　将来的には送客求職者テーブルから面談日時カラムを消す
		sendingJobSeeker.StudyCategory,
		sendingJobSeeker.WordSkill,
		sendingJobSeeker.ExcelSkill,
		sendingJobSeeker.PowerPointSkill,
		sendingJobSeeker.PublicMemo,
		sendingJobSeeker.NationalityRemarks,
		sendingJobSeeker.MedicalHistoryRemarks,
		sendingJobSeeker.AcceptancePoints,
		// sendingJobSeeker.SendCustomerID, // 送客IDは更新しない
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 開発環境で実行時にテストユーザーの名前で更新する
func (repo *SendingJobSeekerRepositoryImpl) UpdateForDev(sendingJobSeeker *entity.SendingJobSeeker) error {
	sendingJobSeeker.LastName = "テスト"
	sendingJobSeeker.FirstName = fmt.Sprintf("ユーザー%v", sendingJobSeeker.ID)
	sendingJobSeeker.LastFurigana = "テスト"
	sendingJobSeeker.FirstFurigana = fmt.Sprintf("ユーザー%v", sendingJobSeeker.ID)
	sendingJobSeeker.PhoneNumber = "080-0000-0000"

	_, err := repo.executer.Exec(
		repo.Name+".UpdateForDev",
		`
		UPDATE sending_job_seekers
		SET
			last_name = ?,
			first_name = ?,
			last_furigana = ?,
			first_furigana = ?,
			phone_number = ?
		WHERE 
			id = ?
		`,
		sendingJobSeeker.LastName,
		sendingJobSeeker.FirstName,
		sendingJobSeeker.LastFurigana,
		sendingJobSeeker.FirstFurigana,
		sendingJobSeeker.PhoneNumber,
		sendingJobSeeker.ID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// phaseの更新
func (repo *SendingJobSeekerRepositoryImpl) UpdatePhase(sendingJobSeekerID uint, phase null.Int) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdatePhase",
		`
		UPDATE sending_job_seekers
		SET
			phase = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		phase,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 面談日時の更新
func (repo *SendingJobSeekerRepositoryImpl) UpdateInterviewDateByCustomerID(sendingCustomerID uint, interviewDate time.Time) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateInterviewDateByCustomerID",
		`
		UPDATE sending_job_seekers
		SET
			interview_date = ?,
			updated_at = ?
		WHERE 
			sending_customer_id = ?
		`,
		interviewDate,
		time.Now().In(time.UTC),
		sendingCustomerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// アクティビティメモを更新
func (repo *SendingJobSeekerRepositoryImpl) UpdateActivityMemo(sendingJobSeekeID uint, activityMemo string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateActivityMemo",
		`
		UPDATE sending_job_seekers
		SET
			activity_memo = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		activityMemo,
		time.Now().In(time.UTC),
		sendingJobSeekeID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// questionの更新
func (repo *SendingJobSeekerRepositoryImpl) UpdateQuestion(sendingJobSeekerID uint, question string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateQuestion",
		`
		UPDATE sending_job_seekers
		SET
			question = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		question,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// agreementの更新
func (repo *SendingJobSeekerRepositoryImpl) UpdateAgreement(sendingJobSeekerID uint, agreement bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgreement",
		`
		UPDATE sending_job_seekers
		SET
			agreement = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		agreement,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerRepositoryImpl) UpdateAgentStaffIDBySendingCustomerID(sendingCustomerID, agentStaffID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgentStaffIDBySendingCustomerID",
		`
		UPDATE sending_job_seekers 
		SET
			agent_staff_id = ?,
			updated_at = ?
		WHERE 
			sending_customer_id = ?
		`,
		agentStaffID,
		time.Now().In(time.UTC),
		sendingCustomerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seekers
		WHERE id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerRepositoryImpl) FindByID(sendingJobSeekerID uint) (*entity.SendingJobSeeker, error) {
	var (
		sendingJobSeeker entity.SendingJobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&sendingJobSeeker, `
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email,
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number, 
			agent.agent_name,
			chat_group.line_active,
			sender_agent.agent_name AS sender_agent_name,
			IFNULL(end_status.end_reason, '') AS end_reason,
			end_status.end_status AS end_status
		FROM 
			sending_job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		INNER JOIN
			chat_group_with_sending_job_seekers AS chat_group
		ON
			seeker.id = chat_group.sending_job_seeker_id
		INNER JOIN
			sending_customers AS customer
		ON
			seeker.sending_customer_id = customer.id
		INNER JOIN
			agents AS sender_agent
		ON
			customer.agent_id = sender_agent.id
		LEFT OUTER JOIN
			sending_job_seeker_end_statuses AS end_status
		ON
			seeker.id = end_status.sending_job_seeker_id
		WHERE
			seeker.id = ?
		LIMIT 1
		`, sendingJobSeekerID)
	if err != nil {
		return nil, err
	}

	return &sendingJobSeeker, nil
}

func (repo *SendingJobSeekerRepositoryImpl) FindByUUID(sendingJobSeekerUUID uuid.UUID) (*entity.SendingJobSeeker, error) {
	var (
		sendingJobSeeker entity.SendingJobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByUUID",
		&sendingJobSeeker, `
		SELECT 
			*
		FROM 
			sending_job_seekers
		WHERE
			uuid = ?
		LIMIT 1
		`,
		sendingJobSeekerUUID,
	)

	if err != nil {
		return nil, err
	}

	return &sendingJobSeeker, nil
}

// 送客IDから送客求職者を取得
func (repo *SendingJobSeekerRepositoryImpl) FindBySendingCustomerID(sendingCustomerID uint) (*entity.SendingJobSeeker, error) {
	var (
		sendingJobSeeker entity.SendingJobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindBySendingCustomerID",
		&sendingJobSeeker, `
		SELECT 
			seeker.*
		FROM 
			sending_job_seekers AS seeker
		WHERE
			seeker.sending_customer_id = ?
		LIMIT 1
		`,
		sendingCustomerID,
	)

	if err != nil {
		return nil, err
	}

	return &sendingJobSeeker, nil
}

// Webhookで送られてきたLineIDから求職者を取得
func (repo *SendingJobSeekerRepositoryImpl) FindByLineID(sendingJobSeekerLineID string) (*entity.SendingJobSeeker, error) {
	var (
		sendingJobSeeker entity.SendingJobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByLineID",
		&sendingJobSeeker, `
		SELECT 
			seeker.*
		FROM 
			sending_job_seekers AS seeker
		WHERE
			seeker.line_id = ?
		LIMIT 1
		`,
		sendingJobSeekerLineID,
	)

	if err != nil {
		return nil, err
	}

	return &sendingJobSeeker, nil
}

// LINEIDを更新
func (repo *SendingJobSeekerRepositoryImpl) UpdateLineID(jobSeekeUUID uuid.UUID, lineID string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateLineID",
		`
			UPDATE
				sending_job_seekers
			SET
				line_id = ?
			WHERE
				uuid = ?
		`,
		lineID,
		jobSeekeUUID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 送客の重複登録判定
func (repo *SendingJobSeekerRepositoryImpl) FindByNameAndPhoneNumber(firstName, lastName, firstFurigana, lastFurigana, phoneNumber string) (*entity.SendingJobSeeker, error) {
	var (
		sendingJobSeeker entity.SendingJobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByNameAndPhoneNumber",
		&sendingJobSeeker, `
		SELECT 
			seeker.*
		FROM 
			sending_job_seekers AS seeker
		WHERE
		seeker.first_name = ?
		AND
		seeker.last_name = ?
		AND
		seeker.first_furigana = ?
		AND
		seeker.last_furigana = ?
		AND
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

	return &sendingJobSeeker, nil
}

// agentIDから送客求職者一覧を取得
func (repo *SendingJobSeekerRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeeker, error) {
	var (
		sendingJobSeekerList []*entity.SendingJobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&sendingJobSeekerList, `
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
			agent.agent_name
		FROM 
			sending_job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		WHERE
			seeker.agent_id = ?
		ORDER BY seeker.id DESC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobSeekerList, nil
}

// フェーズから送客求職者一覧を取得
func (repo *SendingJobSeekerRepositoryImpl) GetListByPhase(phase []uint) ([]*entity.SendingJobSeeker, error) {
	var (
		sendingJobSeekerList []*entity.SendingJobSeeker
	)

	query := fmt.Sprintf(`
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
			agent.agent_name AS agent_name
		FROM 
			sending_job_seekers AS seeker
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		WHERE
			seeker.phase IN (%v)
		ORDER BY seeker.id DESC
		`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(phase)), ","), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&sendingJobSeekerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobSeekerList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての送客求職者情報を取得
func (repo *SendingJobSeekerRepositoryImpl) GetAll() ([]*entity.SendingJobSeeker, error) {
	var (
		sendingJobSeekerList []*entity.SendingJobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&sendingJobSeekerList, `
			SELECT 
				seeker.*, 
				IFNULL(staff.staff_name, '') AS staff_name, 
				IFNULL(staff.email, '') AS staff_email, 
				IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
				agent.agent_name, 
				sender_agent.agent_name AS sender_agent_name, sender_agent.id AS sender_agent_id
			FROM 
				sending_job_seekers AS seeker
			INNER JOIN
				sending_customer AS customer
			ON
				seeker.sending_customer_id = customer.id
			INNER JOIN 
				agents AS sender_agent
			ON
				customer.agent_id = sender_agent.id
			LEFT OUTER JOIN
				agent_staffs AS staff
			ON
				seeker.agent_staff_id = staff.id
			INNER JOIN
				agents AS agent
			ON
				seeker.agent_id = agent.id
			ORDER BY seeker.id DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobSeekerList, nil
}

func (repo *SendingJobSeekerRepositoryImpl) GetSearchListForSendingJobSeekerManagement(agentID uint, freeWord string, staffIDList, senderIDList []uint) ([]*entity.SendingJobSeeker, error) {
	var (
		sendingJobSeekerList []*entity.SendingJobSeeker
		freeWordQuery        string
		staffQuery           string
		senderQuery          string
	)

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {
			freeWordForLike := "%" + freeWord + "%"

			freeWordQuery = fmt.Sprintf(`
				AND (	
					CONCAT(seeker.last_name, seeker.first_name) LIKE '%s' OR 
					CONCAT(seeker.last_furigana, seeker.first_furigana) LIKE '%s'
				)
			`, freeWordForLike, freeWordForLike)
		}
	} else {
		// ハイフンなしの場合数字として扱われるためハイフンなしはこちらに入れる
		freeWordQuery = fmt.Sprintf(`
			AND seeker.id = %d
		`, freeWordInt)
	}

	// 担当者
	if len(staffIDList) > 0 {
		staffQuery = fmt.Sprintf(`
				AND seeker.agent_staff_id IN(%s)
			`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(staffIDList)), ", "), "[]"))
	}

	// 送客元
	if len(senderIDList) > 0 {
		senderQuery = fmt.Sprintf(`
				AND customer.agent_id IN(%s)
			`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(senderIDList)), ", "), "[]"))
	}

	query := fmt.Sprintf(`
		SELECT 
			seeker.*, 
			IFNULL(staff.staff_name, '') AS staff_name, 
			IFNULL(staff.email, '') AS staff_email, 
			IFNULL(staff.staff_phone_number, '') AS staff_phone_number,
			is_view.is_not_waiting_viewed, is_view.is_not_unregister_viewed,
			agent.agent_name AS agent_name,
			sender_agent.agent_name AS sender_agent_name, sender_agent.id AS sender_agent_id
		FROM 
			sending_job_seekers AS seeker
		INNER JOIN
			sending_customers AS customer
		ON
			seeker.sending_customer_id = customer.id
		INNER JOIN 
			agents AS sender_agent
		ON
			customer.agent_id = sender_agent.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			seeker.agent_id = agent.id
		INNER JOIN
			sending_job_seeker_is_views AS is_view
		ON
			seeker.id = is_view.sending_job_seeker_id
		WHERE
			seeker.agent_id = ?
		%s %s %s
		ORDER BY seeker.id DESC
	`, freeWordQuery, staffQuery, senderQuery)

	err = repo.executer.Select(
		repo.Name+".GetSearchListForSendingJobSeekerManagement",
		&sendingJobSeekerList,
		query,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobSeekerList, nil
}

/****************************************************************************************/
