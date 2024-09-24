package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerEndStatusRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerEndStatusRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerEndStatusRepository {
	return &SendingJobSeekerEndStatusRepositoryImpl{
		Name:     "SendingJobSeekerEndStatusRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerEndStatusRepositoryImpl) Create(endStatus *entity.SendingJobSeekerEndStatus) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_end_statuses (
				sending_job_seeker_id,
				end_reason,
				end_status,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		endStatus.SendingJobSeekerID,
		endStatus.EndReason,
		endStatus.EndStatus,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	endStatus.ID = uint(lastID)

	return nil
}
