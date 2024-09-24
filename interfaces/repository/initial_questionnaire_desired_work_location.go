package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type InitialQuestionnaireDesiredWorkLocationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewInitialQuestionnaireDesiredWorkLocationRepositoryImpl(ex interfaces.SQLExecuter) usecase.InitialQuestionnaireDesiredWorkLocationRepository {
	return &InitialQuestionnaireDesiredWorkLocationRepositoryImpl{
		Name:     "InitialQuestionnaireDesiredWorkLocationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
//求職者の面談前アンケートの希望勤務地を作成
func (repo *InitialQuestionnaireDesiredWorkLocationRepositoryImpl) Create(desiredWorkLocation *entity.InitialQuestionnaireDesiredWorkLocation) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO initial_questionnaire_desired_work_locations (
			initial_questionnaire_id,
			desired_work_location,
			desired_rank,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?
		)`,
		desiredWorkLocation.InitialQuestionnaireID,
		desiredWorkLocation.DesiredWorkLocation,
		desiredWorkLocation.DesiredRank,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredWorkLocation.ID = uint(lastID)
	return nil
}
