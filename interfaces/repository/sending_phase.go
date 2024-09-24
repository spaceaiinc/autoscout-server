package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingPhaseRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingPhaseRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingPhaseRepository {
	return &SendingPhaseRepositoryImpl{
		Name:     "SendingPhaseRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingPhaseRepositoryImpl) Create(phase *entity.SendingPhase) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_phases (
				uuid,
				sending_job_seeker_id,
				sending_enterprise_id,
				phase,
				sending_date,
				is_attended,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?
			)
		`,
		utility.CreateUUID(),
		phase.SendingJobSeekerID,
		phase.SendingEnterpriseID,
		phase.Phase,
		phase.SendingDate,
		phase.IsAttended,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	phase.ID = uint(lastID)

	return nil
}

// 進捗の更新
func (repo *SendingPhaseRepositoryImpl) UpdatePhase(id, phase uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdatePhase",
		`
			UPDATE
				sending_phases
			SET
				phase = ?
			WHERE
				id = ?
		`,
		phase,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 進捗の更新
func (repo *SendingPhaseRepositoryImpl) UpdatePhaseBySendingJobSeekerIDAndSendingEnterpriseID(sendingJobSeekerID, sendingEnterpriseID, phase uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdatePhaseBySendingJobSeekerIDAndSendingEnterpriseID",
		`
			UPDATE
				sending_phases
			SET
				phase = ?
			WHERE
				sending_job_seeker_id = ? AND
				sending_enterprise_id = ?
		`,
		phase, sendingJobSeekerID, sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 送客予定日時の更新
func (repo *SendingPhaseRepositoryImpl) UpdateSendingDate(id uint, sendingDate time.Time) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateSendingDate",
		`
			UPDATE
				sending_phases
			SET
				sending_date = ?
			WHERE
				id = ?
		`,
		sendingDate,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 参加確認有無の更新
func (repo *SendingPhaseRepositoryImpl) UpdateIsAttended(id uint, isAttended bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsAttended",
		`
			UPDATE
				sending_phases
			SET
				is_attended = ?
			WHERE
				id = ?
		`,
		isAttended,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (repo *SendingPhaseRepositoryImpl) FindByID(id uint) (*entity.SendingPhase, error) {
	var (
		phase entity.SendingPhase
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&phase, `
		SELECT
			phase.*,
			seeker.last_name, seeker.first_name,
			seeker.last_furigana, seeker.first_furigana,
			enterprise.company_name AS agent_name,
			billing.schedule_adjustment_url, billing.commission
		FROM
			sending_phases AS phase
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			phase.sending_job_seeker_id = seeker.id
		INNER JOIN
			sending_enterprises AS enterprise
		ON
			phase.sending_enterprise_id = enterprise.id
		INNER JOIN
			sending_billing_addresses AS billing
		ON
			enterprise.id = billing.sending_enterprise_id
		WHERE
			phase.id = ?
		`,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &phase, nil
}

func (repo *SendingPhaseRepositoryImpl) FindBySendingJobSeekerIDAndSendingEnterpriseID(sendingJobSeekerID, sendigEnterpriseID uint) (*entity.SendingPhase, error) {
	var (
		phase entity.SendingPhase
	)

	err := repo.executer.Get(
		repo.Name+".FindBySendingJobSeekerIDAndSendingEnterpriseID",
		&phase, `
		SELECT 
		  *
		FROM 
		  sending_phases
		WHERE
			sending_job_seeker_id = ?
		AND
			sending_enterprise_id = ?
		`,
		sendingJobSeekerID, sendigEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &phase, nil
}

func (repo *SendingPhaseRepositoryImpl) GetSearchListForSendingJobSeekerManagement(agentID uint, freeWord string, staffIDList, senderIDList, sendAgentIDList []uint) ([]*entity.SendingPhase, error) {
	var (
		phaseList      []*entity.SendingPhase
		freeWordQuery  string
		staffQuery     string
		senderQuery    string
		sendAgentQuery string
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

	// 送客先
	if len(sendAgentIDList) > 0 {
		sendAgentQuery = fmt.Sprintf(`
			AND enterprise.id IN(%s)
		`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(sendAgentIDList)), ", "), "[]"))
	}

	query := fmt.Sprintf(`
		SELECT
			phase.*,
			seeker.last_name, seeker.first_name,
			seeker.last_furigana, seeker.first_furigana,
			seeker.agent_staff_id, IFNULL(staff.staff_name, '') AS staff_name,
			seeker.interview_date,
			enterprise.company_name AS agent_name,
			sender_agent.agent_name AS sender_agent_name
		FROM
			sending_phases AS phase
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			phase.sending_job_seeker_id = seeker.id
		INNER JOIN
			sending_customers AS customer
		ON
			seeker.sending_customer_id = customer.id
		INNER JOIN
			agents AS sender_agent
		ON
			customer.agent_id = sender_agent.id
		INNER JOIN
			sending_enterprises AS enterprise
		ON
			phase.sending_enterprise_id = enterprise.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		WHERE
			phase.phase IN(?, ?, ?)
		AND
			seeker.agent_id = ?
			%s %s %s %s
		ORDER BY phase.id DESC
	`, freeWordQuery, staffQuery, senderQuery, sendAgentQuery)

	err = repo.executer.Select(
		repo.Name+".GetSearchListForSendingJobSeekerManagement",
		&phaseList,
		query,
		uint(entity.AcceptSending),   // 送客応諾
		uint(entity.CompleteSending), // 送客完了
		uint(entity.CloseSending),    // 終了
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return phaseList, nil
}

// 求職者IDから進捗一覧を取得
func (repo *SendingPhaseRepositoryImpl) GetListBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingPhase, error) {
	var (
		phaseList []*entity.SendingPhase
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobSeekerID",
		&phaseList, `
			SELECT
				phase.*,
				enterprise.company_name AS agent_name,
				IFNULL(end_status.end_reason, '') AS end_reason,
				end_status.end_status AS end_status
			FROM
				sending_phases AS phase
			INNER JOIN
				sending_enterprises AS enterprise
			ON
				phase.sending_enterprise_id = enterprise.id
			LEFT JOIN
				sending_phase_end_statuses AS end_status
			ON
				phase.id = end_status.sending_phase_id
			WHERE
				phase.sending_job_seeker_id = ?
			ORDER BY phase.id DESC
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return phaseList, nil
}

func (repo *SendingPhaseRepositoryImpl) GetListByPhaseAndEnterpriseIDList(phase uint, sendingEnterpriseIDList []uint) ([]*entity.SendingPhase, error) {
	var (
		phaseList []*entity.SendingPhase
	)

	if len(sendingEnterpriseIDList) == 0 {
		return phaseList, nil
	}

	query := fmt.Sprintf(`
		SELECT
			phase.*,
			seeker.last_name, seeker.first_name,
			seeker.last_furigana, seeker.first_furigana,
			enterprise.company_name AS agent_name
		FROM
			sending_phases AS phase
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			phase.sending_job_seeker_id = seeker.id
		INNER JOIN
			sending_enterprises AS enterprise
		ON
			phase.sending_enterprise_id = enterprise.id
		WHERE
			phase.sending_enterprise_id IN(%s)
		AND
			phase.phase = ?
		ORDER BY phase.id DESC
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(sendingEnterpriseIDList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByPhaseAndEnterpriseIDList",
		&phaseList, query, phase,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return phaseList, nil
}

func (repo *SendingPhaseRepositoryImpl) GetListByAgentIDAndPhaseAndEnterpriseIDList(agentID, phase uint, sendingEnterpriseIDList []uint) ([]*entity.SendingPhase, error) {
	var (
		phaseList []*entity.SendingPhase
	)

	if len(sendingEnterpriseIDList) == 0 {
		return phaseList, nil
	}

	query := fmt.Sprintf(`
		SELECT
			phase.*,
			seeker.last_name, seeker.first_name,
			seeker.last_furigana, seeker.first_furigana,
			enterprise.company_name AS agent_name
		FROM
			sending_phases AS phase
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			phase.sending_job_seeker_id = seeker.id
		INNER JOIN
			sending_enterprises AS enterprise
		ON
			phase.sending_enterprise_id = enterprise.id
		WHERE
			phase.sending_enterprise_id IN(%s)
		AND
			phase.phase = ?
		AND
			seeker.agent_id = ?
		ORDER BY phase.id DESC
		`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(sendingEnterpriseIDList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentIDAndPhaseAndEnterpriseIDList",
		&phaseList, query, phase, agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return phaseList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての進捗情報を取得
func (repo *SendingPhaseRepositoryImpl) GetAll() ([]*entity.SendingPhase, error) {
	var (
		phaseList []*entity.SendingPhase
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&phaseList,
		`
			SELECT
				phase.*,
				seeker.last_name, seeker.first_name,
				seeker.last_furigana, seeker.first_furigana,
				enterprise.company_name AS agent_name
			FROM
				sending_phases AS phase
			INNER JOIN
				sending_job_seekers AS seeker
			ON
				phase.sending_job_seeker_id = seeker.id
			INNER JOIN
				sending_enterprises AS enterprise
			ON
				phase.sending_enterprise_id = enterprise.id
			WHERE
				phase.phase = ?
			ORDER BY phase.id DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return phaseList, nil
}
