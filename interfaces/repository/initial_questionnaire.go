package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type InitialQuestionnaireRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewInitialQuestionnaireRepositoryImpl(ex interfaces.SQLExecuter) usecase.InitialQuestionnaireRepository {
	return &InitialQuestionnaireRepositoryImpl{
		Name:     "InitialQuestionnaireRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// 求職者の面談前アンケートの作成
func (repo *InitialQuestionnaireRepositoryImpl) Create(initialQuestion *entity.InitialQuestionnaire) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO initial_questionnaires (
			uuid,
			job_seeker_id,
			question,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?
		)`,
		utility.CreateUUID(),
		initialQuestion.JobSeekerID,
		initialQuestion.Question,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	initialQuestion.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 単数取得 API
//
//求職者IDから面談前アンケートの取得
func (repo *InitialQuestionnaireRepositoryImpl) FindByJobSeekerID(jobSeekerID uint) (*entity.InitialQuestionnaire, error) {
	var initialQuestion *entity.InitialQuestionnaire

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerID",
		initialQuestion,
		`
		SELECT * 
		FROM initial_questionnaires 
		WHERE job_seeker_id = ?`,
		jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return initialQuestion, nil
}
