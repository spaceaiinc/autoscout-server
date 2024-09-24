package repository

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatMessageWithSendingJobSeekerRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewChatMessageWithSendingJobSeekerRepositoryImpl(ex interfaces.SQLExecuter) usecase.ChatMessageWithSendingJobSeekerRepository {
	return &ChatMessageWithSendingJobSeekerRepositoryImpl{
		Name:     "ChatMessageWithSendingJobSeekerRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
//タスクの作成
func (repo *ChatMessageWithSendingJobSeekerRepositoryImpl) Create(chatMessageWithSendingJobSeeker *entity.ChatMessageWithSendingJobSeeker) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO chat_message_with_sending_job_seekers (
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
		chatMessageWithSendingJobSeeker.GroupID,
		chatMessageWithSendingJobSeeker.UserType,
		chatMessageWithSendingJobSeeker.Message,
		chatMessageWithSendingJobSeeker.PackageID,
		chatMessageWithSendingJobSeeker.StickerID,
		chatMessageWithSendingJobSeeker.OriginalContentURL,
		chatMessageWithSendingJobSeeker.PreviewImageURL,
		chatMessageWithSendingJobSeeker.Duration,
		chatMessageWithSendingJobSeeker.LineMessageID,
		chatMessageWithSendingJobSeeker.LineMessageType,
		chatMessageWithSendingJobSeeker.CreatedAt, // WebHookEventのtimestampが入ってくる
		chatMessageWithSendingJobSeeker.UpdatedAt, // WebHookEventのtimestampが入ってくる
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	chatMessageWithSendingJobSeeker.ID = uint(lastID)
	return nil
}

func (repo *ChatMessageWithSendingJobSeekerRepositoryImpl) GetListByGroupID(groupID uint) ([]*entity.ChatMessageWithSendingJobSeeker, error) {
	var (
		chatMessageWithSendingJobSeeker []*entity.ChatMessageWithSendingJobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetListByGroupID",
		&chatMessageWithSendingJobSeeker, `
		SELECT *
		FROM chat_message_with_sending_job_seekers
		WHERE
			group_id = ?
		`,
		groupID)

	if err != nil {
		return nil, err
	}

	return chatMessageWithSendingJobSeeker, nil
}
