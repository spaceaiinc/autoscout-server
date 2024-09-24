package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatGroupWithJobSeekerRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewChatGroupWithJobSeekerRepositoryImpl(ex interfaces.SQLExecuter) usecase.ChatGroupWithJobSeekerRepository {
	return &ChatGroupWithJobSeekerRepositoryImpl{
		Name:     "ChatGroupWithJobSeekerRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
//エージェントと求職者のチャットグループの作成
func (repo *ChatGroupWithJobSeekerRepositoryImpl) Create(chatGroupWithJobSeeker *entity.ChatGroupWithJobSeeker) error {
	nowTime := time.Now().In(time.UTC)
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO chat_group_with_job_seekers (
			agent_id,
			job_seeker_id,
			agent_last_send_at,	                        -- エージェントの最終送信時間
			agent_last_watched_at,	                    -- エージェントの最終閲覧時間
			job_seeker_last_send_at,	                    -- 求職者の最終送信時間
			job_seeker_last_watched_at,	
			line_active,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?
			)`,
		chatGroupWithJobSeeker.AgentID,
		chatGroupWithJobSeeker.JobSeekerID,
		nowTime,
		nowTime,
		nowTime,
		nowTime,
		false, // 作成の時はLINE連携してないからfalse
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	chatGroupWithJobSeeker.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
// エージェントの最終閲覧時間を更新
func (repo *ChatGroupWithJobSeekerRepositoryImpl) UpdateAgentLastWatchedAt(id uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgentLastWatchedAt",
		`
		UPDATE 
			chat_group_with_job_seekers 
		SET
			agent_last_watched_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime,
		nowTime,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// エージェントの最終送信時間を更新
func (repo *ChatGroupWithJobSeekerRepositoryImpl) UpdateAgentLastSendAt(id uint) error {
	nowTime := time.Now().In(time.UTC)
	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgentLastSendAt",
		`
			UPDATE 
				chat_group_with_job_seekers 
			SET
				agent_last_send_at = ?,
				updated_at = ?
			WHERE 
				id = ?
		`,
		nowTime,
		nowTime,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 求職者の最終閲覧時間と最終送信時間を更新
func (repo *ChatGroupWithJobSeekerRepositoryImpl) UpdateJobSeekerLastWatchedAtAndSendAt(id uint, sendAt time.Time) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateJobSeekerLastWatchedAtAndSendAt",
		`
			UPDATE 
				chat_group_with_job_seekers 
			SET
				job_seeker_last_send_at = ?,
				job_seeker_last_watched_at = ?,	
				updated_at = ?
			WHERE 
				id = ?
		`,
		sendAt,
		sendAt,
		sendAt,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 求職者とのチャットグループのLINEの有効/無効を更新
func (repo *ChatGroupWithJobSeekerRepositoryImpl) UpdateJobSeekerLineActive(jobSeekerID uint, isActive bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateJobSeekerLineActive",
		`
			UPDATE 
				chat_group_with_job_seekers 
			SET
				line_active = ?,
				updated_at = ?
			WHERE 
				job_seeker_id = ?
		`,
		isActive,
		time.Now().In(time.UTC),
		jobSeekerID,
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
// IDを使ってエージェントと求職者のチャットグループを取得
func (repo *ChatGroupWithJobSeekerRepositoryImpl) FindByID(id uint) (*entity.ChatGroupWithJobSeeker, error) {
	var (
		chatGroupWithJobSeeker entity.ChatGroupWithJobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&chatGroupWithJobSeeker, `
		SELECT 
			chat_group.*,
			seeker.last_name, seeker.first_name, 
			seeker.last_furigana, seeker.first_furigana
		FROM 
			chat_group_with_job_seekers AS chat_group
		INNER JOIN
			job_seekers AS seeker
		ON
			chat_group.job_seeker_id = seeker.id
		WHERE
			chat_group.id = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &chatGroupWithJobSeeker, nil
}

// agentIDを使ってエージェントと求職者のチャットグループを取得
func (repo *ChatGroupWithJobSeekerRepositoryImpl) FindByAgentIDAndJobSeekerID(agentID, jobSeekerID uint) (*entity.ChatGroupWithJobSeeker, error) {
	var (
		chatGroupWithJobSeeker entity.ChatGroupWithJobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByAgentIDAndJobSeekerID",
		&chatGroupWithJobSeeker, `
		SELECT 
			*
		FROM 
			chat_group_with_job_seekers
		WHERE
			agent_id = ?
		AND
			job_seeker_id = ?
		LIMIT 1
		`,
		agentID, jobSeekerID,
	)

	if err != nil {
		return nil, err
	}

	return &chatGroupWithJobSeeker, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// agentIDを使ってエージェントと求職者のチャットグループグループの一覧を取得 *LINE連携していない求職者も含むよう変更し、メールアドレスを取得する(2023/08/08)
func (repo *ChatGroupWithJobSeekerRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.ChatGroupWithJobSeeker, error) {
	var (
		chatGroupWithJobSeekerList []*entity.ChatGroupWithJobSeeker
	)

	err := repo.executer.Select(
		repo.Name+"GetByAgentID",
		&chatGroupWithJobSeekerList, `
			SELECT 
				chat_group.*,
				seeker.line_id, seeker.last_name, seeker.first_name, 
				seeker.last_furigana, seeker.first_furigana, seeker.agent_staff_id AS ca_staff_id,
				IFNULL(staff.staff_name, '') AS staff_name,
				seeker.email,
				seeker.phase
			FROM 
				chat_group_with_job_seekers AS chat_group
			INNER JOIN
				job_seekers AS seeker
			ON
				chat_group.job_seeker_id = seeker.id
			LEFT OUTER JOIN
				agent_staffs AS staff
			ON
				seeker.agent_staff_id = staff.id
			WHERE
				chat_group.agent_id = ?
			AND
				seeker.email IS NOT NULL 
			AND 
			  seeker.email != ''
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return chatGroupWithJobSeekerList, nil
}

// agentIDとfreeWordを使ってエージェントと求職者のチャットグループグループの一覧を取得
func (repo *ChatGroupWithJobSeekerRepositoryImpl) GetByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.ChatGroupWithJobSeeker, error) {
	var (
		chatGroupWithJobSeekerList []*entity.ChatGroupWithJobSeeker
		freeWordQuery              string
	)

	if freeWord != "" {
		freeWordForLike := "%" + freeWord + "%"

		freeWordQuery = fmt.Sprintf(`
			AND (	
				CONCAT(seeker.last_name, seeker.first_name) LIKE '%s' OR 
				CONCAT(seeker.last_furigana, seeker.first_furigana) LIKE '%s'
			)
		`, freeWordForLike, freeWordForLike)
	}

	query := fmt.Sprintf(`
		SELECT 
			chat_group.*,
			seeker.line_id, seeker.last_name, seeker.first_name, 
			seeker.last_furigana, seeker.first_furigana, seeker.agent_staff_id AS ca_staff_id,
			IFNULL(staff.staff_name, '') AS staff_name, document.id_photo_url,
			seeker.email,
			seeker.phase
		FROM 
			chat_group_with_job_seekers AS chat_group
		INNER JOIN
			job_seekers AS seeker
		ON
			chat_group.job_seeker_id = seeker.id
		INNER JOIN
			job_seeker_documents AS document
		ON
			seeker.id = document.job_seeker_id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		WHERE
			chat_group.agent_id = %v
		AND
			seeker.email IS NOT NULL 
		AND 
			seeker.email != ''
		%s
		ORDER BY
			CASE
				WHEN chat_group.agent_last_send_at > chat_group.job_seeker_last_send_at
				THEN chat_group.agent_last_send_at
				ELSE chat_group.job_seeker_last_send_at
			END DESC,
			chat_group.job_seeker_id DESC
	`, agentID, freeWordQuery)

	err := repo.executer.Select(
		repo.Name+"GetByAgentIDAndFreeWord",
		&chatGroupWithJobSeekerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return chatGroupWithJobSeekerList, nil
}

// エージェントIDからアクティブなチャットグループを取得 ※管理側一斉送信用
func (repo *ChatGroupWithJobSeekerRepositoryImpl) GetActiveJobSeekerByAgentID(agentID uint) ([]*entity.ChatGroupWithJobSeeker, error) {
	var (
		chatGroupWithJobSeekerList []*entity.ChatGroupWithJobSeeker

		Hamada    = 2
		Okuya     = 3
		Okamoto   = 54
		Tadashige = 55
	)

	err := repo.executer.Select(
		repo.Name+"GetActiveJobSeekerByAgentID",
		&chatGroupWithJobSeekerList, `
			SELECT 
				chat_group.*,
				seeker.line_id, seeker.last_name, seeker.first_name, seeker.uuid AS job_seeker_uuid,
				seeker.last_furigana, seeker.first_furigana, seeker.agent_staff_id AS ca_staff_id,
				IFNULL(staff.staff_name, '') AS staff_name,
				seeker.email,
				seeker.phase
			FROM 
				chat_group_with_job_seekers AS chat_group
			INNER JOIN
				job_seekers AS seeker
			ON
				chat_group.job_seeker_id = seeker.id
			LEFT OUTER JOIN
				agent_staffs AS staff
			ON
				seeker.agent_staff_id = staff.id
			WHERE
				chat_group.agent_id = ?
			AND
				seeker.agent_staff_id IN(?, ?, ?, ?)
			AND
				seeker.phase IN(4, 5, 6)
			AND
				seeker.line_id != ''
		`,
		agentID,
		Hamada, Okuya, Okamoto, Tadashige,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return chatGroupWithJobSeekerList, nil
}

// All *未読通知判定に使用
func (repo *ChatGroupWithJobSeekerRepositoryImpl) All() ([]*entity.ChatGroupWithJobSeeker, error) {
	var (
		chatGroupWithJobSeekerList []*entity.ChatGroupWithJobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&chatGroupWithJobSeekerList, `
		SELECT 
			chat_group.*,
			seeker.last_name, seeker.first_name, 
			seeker.last_furigana, seeker.first_furigana, 
			seeker.agent_staff_id AS ca_staff_id
		FROM 
			chat_group_with_job_seekers AS chat_group
		INNER JOIN
			job_seekers AS seeker
		ON
			chat_group.job_seeker_id = seeker.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		WHERE
			line_active = true
		ORDER BY
			id DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return chatGroupWithJobSeekerList, nil
}

/****************************************************************************************/
/// その他 API
//
// エージェントの未読通知数を取得
func (repo *ChatGroupWithJobSeekerRepositoryImpl) CountNotificationByAgentID(agentID uint) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+"CountNotificationByAgentID",
		&result, `
            SELECT 
                COUNT(*) AS count
            FROM 
                chat_group_with_job_seekers AS chat_group
			INNER JOIN
				job_seekers AS seeker
			ON
				chat_group.job_seeker_id = seeker.id
			WHERE
				seeker.agent_id = ?
			AND
                chat_group.job_seeker_last_send_at > chat_group.agent_last_watched_at
        `,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}
