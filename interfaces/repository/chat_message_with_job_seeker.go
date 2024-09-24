package repository

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatMessageWithJobSeekerRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewChatMessageWithJobSeekerRepositoryImpl(ex interfaces.SQLExecuter) usecase.ChatMessageWithJobSeekerRepository {
	return &ChatMessageWithJobSeekerRepositoryImpl{
		Name:     "ChatMessageWithJobSeekerRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// チャットメッセージの作成
func (repo *ChatMessageWithJobSeekerRepositoryImpl) Create(chatMessageWithJobSeeker *entity.ChatMessageWithJobSeeker) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO chat_message_with_job_seekers (
			group_id,
			user_type,
			message,
			package_id,
			sticker_id,
			original_content_url,
			preview_image_url,
			duration,
			line_message_id,
			line_message_type,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)`,
		chatMessageWithJobSeeker.GroupID,
		chatMessageWithJobSeeker.UserType,
		chatMessageWithJobSeeker.Message,
		chatMessageWithJobSeeker.PackageID,
		chatMessageWithJobSeeker.StickerID,
		chatMessageWithJobSeeker.OriginalContentURL,
		chatMessageWithJobSeeker.PreviewImageURL,
		chatMessageWithJobSeeker.Duration,
		chatMessageWithJobSeeker.LineMessageID,
		chatMessageWithJobSeeker.LineMessageType,
		chatMessageWithJobSeeker.CreatedAt, // WebHookEventのtimestampが入ってくる
		chatMessageWithJobSeeker.UpdatedAt, // WebHookEventのtimestampが入ってくる
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	chatMessageWithJobSeeker.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 単数取得 API
//
func (repo *ChatMessageWithJobSeekerRepositoryImpl) FindByLINEMessageID(messageID string) (*entity.ChatMessageWithJobSeeker, error) {
	var (
		chatMessageWithJobSeeker entity.ChatMessageWithJobSeeker
	)

	err := repo.executer.Get(
		repo.Name+".FindByLINEMessageID",
		&chatMessageWithJobSeeker, `
		SELECT *
		FROM chat_message_with_job_seekers
		WHERE
			line_message_id = ?
		LIMIT 1
		`,
		messageID,
	)

	if err != nil {
		return nil, err
	}

	return &chatMessageWithJobSeeker, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// グループIDからチャットメッセージの取得
func (repo *ChatMessageWithJobSeekerRepositoryImpl) GetByGroupID(groupID uint) ([]*entity.ChatMessageWithJobSeeker, error) {
	var (
		chatMessageWithJobSeeker []*entity.ChatMessageWithJobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetByGroupID",
		&chatMessageWithJobSeeker, `
		SELECT *
		FROM chat_message_with_job_seekers
		WHERE
			group_id = ?
		`,
		groupID)

	if err != nil {
		return nil, err
	}

	return chatMessageWithJobSeeker, nil
}
