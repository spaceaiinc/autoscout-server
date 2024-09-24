package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type InitialQuestionnaireDesiredIndustryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewInitialQuestionnaireDesiredIndustryRepositoryImpl(ex interfaces.SQLExecuter) usecase.InitialQuestionnaireDesiredIndustryRepository {
	return &InitialQuestionnaireDesiredIndustryRepositoryImpl{
		Name:     "InitialQuestionnaireDesiredIndustryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
//求職者の面談前アンケートの作成
func (repo *InitialQuestionnaireDesiredIndustryRepositoryImpl) Create(desiredIndustry *entity.InitialQuestionnaireDesiredIndustry) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO initial_questionnaire_desired_industries (
			initial_questionnaire_id,
			desired_industry,
			desired_rank,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?
		)`,
		desiredIndustry.InitialQuestionnaireID,
		desiredIndustry.DesiredIndustry,
		desiredIndustry.DesiredRank,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredIndustry.ID = uint(lastID)
	return nil
}
