package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type EmailWithSendingJobSeekerRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewEmailWithSendingJobSeekerRepositoryImpl(ex interfaces.SQLExecuter) usecase.EmailWithSendingJobSeekerRepository {
	return &EmailWithSendingJobSeekerRepositoryImpl{
		Name:     "EmailWithSendingJobSeeker",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *EmailWithSendingJobSeekerRepositoryImpl) Create(emailWithSendingJobSeeker *entity.EmailWithSendingJobSeeker) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO email_with_sending_job_seekers (
			sending_job_seeker_id,	     
			subject,
			content,
			file_name,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)`,
		emailWithSendingJobSeeker.SendingJobSeekerID,
		emailWithSendingJobSeeker.Subject,
		emailWithSendingJobSeeker.Content,
		emailWithSendingJobSeeker.FileName,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	emailWithSendingJobSeeker.ID = uint(lastID)
	return nil
}

func (repo *EmailWithSendingJobSeekerRepositoryImpl) GetListBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.EmailWithSendingJobSeeker, error) {
	var (
		emailWithSendingJobSeekerList []*entity.EmailWithSendingJobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobSeekerID",
		&emailWithSendingJobSeekerList, `
		SELECT 
			email.*
		FROM 
			email_with_sending_job_seekers AS email
		WHERE
			email.sending_job_seeker_id = ?
		ORDER BY
			email.id DESC
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return emailWithSendingJobSeekerList, nil
}
