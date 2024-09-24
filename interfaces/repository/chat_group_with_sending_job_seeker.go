package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatGroupWithSendingJobSeekerRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewChatGroupWithSendingJobSeekerRepositoryImpl(ex interfaces.SQLExecuter) usecase.ChatGroupWithSendingJobSeekerRepository {
	return &ChatGroupWithSendingJobSeekerRepositoryImpl{
		Name:     "ChatGroupWithSendingJobSeekerRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
//エージェントと求職者のチャットグループの作成
func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) Create(chatGroupWithSendingJobSeeker *entity.ChatGroupWithSendingJobSeeker) error {
	nowTime := time.Now().In(time.UTC)
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO chat_group_with_sending_job_seekers (
			sending_job_seeker_id,
			agent_last_send_at,	                        -- エージェントの最終送信時間
			agent_last_watched_at,	                    -- エージェントの最終閲覧時間
			sending_job_seeker_last_send_at,	                    -- 求職者の最終送信時間
			sending_job_seeker_last_watched_at,	
			line_active,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?
			)`,
		chatGroupWithSendingJobSeeker.SendingJobSeekerID,
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

	chatGroupWithSendingJobSeeker.ID = uint(lastID)
	return nil
}

func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) UpdateAgentLastWatchedAt(groupID uint) error {
	nowTime := time.Now().In(time.UTC)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgentLastWatchedAt",
		`
		UPDATE 
			chat_group_with_sending_job_seekers 
		SET
			agent_last_watched_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		nowTime,
		nowTime,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) UpdateAgentLastSendAt(groupID uint) error {
	nowTime := time.Now().In(time.UTC)
	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgentLastSendAt",
		`
			UPDATE 
				chat_group_with_sending_job_seekers 
			SET
				agent_last_send_at = ?,
				updated_at = ?
			WHERE 
				id = ?
		`,
		nowTime,
		nowTime,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) UpdateSendingJobSeekerLastWatchedAtAndSendAt(groupID uint, sendAt time.Time) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateSendingJobSeekerLastWatchedAtAndSendAt",
		`
			UPDATE 
				chat_group_with_sending_job_seekers 
			SET
				sending_job_seeker_last_send_at = ?,
				sending_job_seeker_last_watched_at = ?,	
				updated_at = ?
			WHERE 
				id = ?
		`,
		sendAt,
		sendAt,
		sendAt,
		groupID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) UpdateSendingJobSeekerLineActive(isActive bool, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateSendingJobSeekerLineActive",
		`
			UPDATE 
				chat_group_with_sending_job_seekers 
			SET
				line_active = ?,
				updated_at = ?
			WHERE 
				sending_job_seeker_id = ?
		`,
		isActive,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) FindByID(chatGroupWithSendingJobSeekerID uint) (*entity.ChatGroupWithSendingJobSeeker, error) {
	var (
		chatGroupWithSendingJobSeeker entity.ChatGroupWithSendingJobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&chatGroupWithSendingJobSeeker, `
		SELECT 
			chat_group.*,
			seeker.last_name, seeker.first_name, 
			seeker.last_furigana, seeker.first_furigana
		FROM 
			chat_group_with_sending_job_seekers AS chat_group
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			chat_group.sending_job_seeker_id = seeker.id
		WHERE
			chat_group.id = ?
		LIMIT 1
		`,
		chatGroupWithSendingJobSeekerID,
	)

	if err != nil {
		return nil, err
	}

	return &chatGroupWithSendingJobSeeker, nil
}

func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) (*entity.ChatGroupWithSendingJobSeeker, error) {
	var (
		chatGroupWithSendingJobSeeker entity.ChatGroupWithSendingJobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindBySendingJobSeekerID",
		&chatGroupWithSendingJobSeeker, `
		SELECT 
			*
		FROM 
			chat_group_with_sending_job_seekers
		WHERE
			sending_job_seeker_id = ?
		LIMIT 1
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		return nil, err
	}

	return &chatGroupWithSendingJobSeeker, nil
}

// ※sending_job_seekerテーブルと結合させたagent_staffテーブルのagent_idと合致するレコードを取得する
func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.ChatGroupWithSendingJobSeeker, error) {
	var (
		chatGroupWithSendingJobSeekerList []*entity.ChatGroupWithSendingJobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&chatGroupWithSendingJobSeekerList, `
		SELECT 
			chat_group.*,
			seeker.last_name, seeker.first_name, 
			seeker.last_furigana, seeker.first_furigana,
			seeker.email,
			seeker.phase,
			seeker.line_id,
			seeker.agent_staff_id AS ca_staff_id,
			IFNULL(staff.staff_name, '') AS staff_name
		FROM 
			chat_group_with_sending_job_seekers AS chat_group
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			chat_group.sending_job_seeker_id = seeker.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		WHERE
			seeker.agent_id = ?
		ORDER BY
			id DESC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return chatGroupWithSendingJobSeekerList, nil
}

func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) GetNotificationCountByAgentID(agentID uint) (float64, error) {
	result := struct {
		Count float64 `db:"count"`
	}{}

	err := repo.executer.Get(
		repo.Name+"GetNotificationCountByAgentID",
		&result, `
            SELECT 
                COUNT(*) AS count
            FROM 
                chat_group_with_sending_job_seekers AS chat_group
			INNER JOIN
				sending_job_seekers AS seeker
			ON
				chat_group.sending_job_seeker_id = seeker.id
			WHERE
				seeker.agent_id = ?
			AND
                chat_group.sending_job_seeker_last_send_at > chat_group.agent_last_watched_at
        `,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, nil
}

// ※sending_job_seekerテーブルと結合させたagent_staffテーブルのagent_idとFreewordと合致するレコードを取得する
func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) GetListByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.ChatGroupWithSendingJobSeeker, error) {
	var (
		chatGroupWithSendingJobSeekerList []*entity.ChatGroupWithSendingJobSeeker
		freeWordQuery                     string
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
			seeker.last_name, seeker.first_name, 
			seeker.last_furigana, seeker.first_furigana,
			seeker.email,
			seeker.phase,
			seeker.line_id,
			seeker.agent_staff_id AS ca_staff_id,
			IFNULL(staff.staff_name, '') AS staff_name
		FROM 
			chat_group_with_sending_job_seekers AS chat_group
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			chat_group.sending_job_seeker_id = seeker.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		WHERE
			seeker.agent_id = %v
		%s
		ORDER BY
			id DESC
		`,
		agentID, freeWordQuery)

	err := repo.executer.Select(
		repo.Name+"GetListByAgentIDAndFreeWord",
		&chatGroupWithSendingJobSeekerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return chatGroupWithSendingJobSeekerList, nil
}

/****************************************************************************************/
/****************************************************************************************/
/// Admin API
//
// GetAll *未読通知判定に使用
func (repo *ChatGroupWithSendingJobSeekerRepositoryImpl) GetAll() ([]*entity.ChatGroupWithSendingJobSeeker, error) {
	var (
		chatGroupWithSendingJobSeekerList []*entity.ChatGroupWithSendingJobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&chatGroupWithSendingJobSeekerList, `
		SELECT 
			chat_group.*,
			seeker.last_name, seeker.first_name, 
			seeker.last_furigana, seeker.first_furigana,
			seeker.email,
			seeker.phase,
			seeker.line_id,
			seeker.agent_staff_id AS ca_staff_id,
			IFNULL(staff.staff_name, '') AS staff_name
		FROM 
			chat_group_with_sending_job_seekers AS chat_group
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			chat_group.sending_job_seeker_id = seeker.id
		LEFT OUTER JOIN
			agent_staffs AS staff
		ON
			seeker.agent_staff_id = staff.id
		ORDER BY
			id DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return chatGroupWithSendingJobSeekerList, nil
}
