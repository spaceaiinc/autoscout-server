package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type InitialQuestionnaireDesiredOccupationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewInitialQuestionnaireDesiredOccupationRepositoryImpl(ex interfaces.SQLExecuter) usecase.InitialQuestionnaireDesiredOccupationRepository {
	return &InitialQuestionnaireDesiredOccupationRepositoryImpl{
		Name:     "InitialQuestionnaireDesiredOccupationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// 求職者の面談前アンケートの希望職種を作成
func (repo *InitialQuestionnaireDesiredOccupationRepositoryImpl) Create(desiredOccupation *entity.InitialQuestionnaireDesiredOccupation) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO initial_questionnaire_desired_occupations (
			initial_questionnaire_id,
			desired_occupation,
			desired_rank,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?
		)`,
		desiredOccupation.InitialQuestionnaireID,
		desiredOccupation.DesiredOccupation,
		desiredOccupation.DesiredRank,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredOccupation.ID = uint(lastID)
	return nil
}
