package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type EmailWithJobSeekerRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewEmailWithJobSeekerRepositoryImpl(ex interfaces.SQLExecuter) usecase.EmailWithJobSeekerRepository {
	return &EmailWithJobSeekerRepositoryImpl{
		Name:     "EmailWithJobSeeker",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *EmailWithJobSeekerRepositoryImpl) Create(emailWithJobSeeker *entity.EmailWithJobSeeker) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO email_with_job_seekers (
			job_seeker_id,	     
			subject,
			content,
			file_name,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)`,
		emailWithJobSeeker.JobSeekerID,
		emailWithJobSeeker.Subject,
		emailWithJobSeeker.Content,
		emailWithJobSeeker.FileName,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	emailWithJobSeeker.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 取得 API
//
func (repo *EmailWithJobSeekerRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.EmailWithJobSeeker, error) {
	var (
		emailWithJobSeekerList []*entity.EmailWithJobSeeker
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&emailWithJobSeekerList, `
		SELECT 
			email.*
		FROM 
			email_with_job_seekers AS email
		WHERE
			email.job_seeker_id = ?
		ORDER BY
			email.id DESC
		`,
		jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return emailWithJobSeekerList, nil
}
