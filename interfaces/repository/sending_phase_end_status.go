package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingPhaseEndStatusRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingPhaseEndStatusRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingPhaseEndStatusRepository {
	return &SendingPhaseEndStatusRepositoryImpl{
		Name:     "SendingPhaseEndStatusRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingPhaseEndStatusRepositoryImpl) Create(endStatus *entity.SendingPhaseEndStatus) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_phase_end_statuses (
				sending_phase_id,
				end_reason,
				end_status,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		endStatus.SendingPhaseID,
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
