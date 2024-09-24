package repository_test

import (
	"testing"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/mock/mock_usecase"
	"go.uber.org/mock/gomock"
)

var (
	initialQuestionnaireEntity1Question = "Text Question 1"
	initialQuestionnaireEntity2Question = "Text Question 2"

	initialQuestionnaireEntity1 = entity.NewInitialQuestionnaire(
		1,
		initialQuestionnaireEntity1Question,
	)

	initialQuestionnaireEntity2 = entity.NewInitialQuestionnaire(
		2,
		initialQuestionnaireEntity2Question,
	)
)

/****************************************************************************************/
// 作成
//
func Test_InitialQuestionnaire_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockInitialQuestionnaire := mock_usecase.NewMockInitialQuestionnaireRepository(mockCtrl)

	mockInitialQuestionnaire.EXPECT().Create(initialQuestionnaireEntity1).Return(nil)

	err := mockInitialQuestionnaire.Create(initialQuestionnaireEntity1)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}

/****************************************************************************************/
// 複数取得
//
func Test_InitialQuestionnaire_FindByJobSeekerID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockInitialQuestionnaire := mock_usecase.NewMockInitialQuestionnaireRepository(mockCtrl)

	mockInitialQuestionnaire.EXPECT().FindByJobSeekerID(gomock.Any()).DoAndReturn(
		func(jobSeekerID uint) (*entity.InitialQuestionnaire, error) {
			switch jobSeekerID {
			case 1:
				initialQuestionnaireEntity1.ID = 1
				return initialQuestionnaireEntity1, nil
			case 2:
				initialQuestionnaireEntity2.ID = 2
				return initialQuestionnaireEntity2, nil
			case 3:
				return nil, nil
			}
			return nil, nil
		},
	)
	initialQuestionnaire, err := mockInitialQuestionnaire.FindByJobSeekerID(1)
	if err != nil {
		t.Fatal("FindByJobSeekerID error!", err)
	}
	if initialQuestionnaire.Question != initialQuestionnaireEntity1Question {
		t.Fatal("FindByJobSeekerID return value error! initialQuestionnaire.Question != initialQuestionnaireEntity1Question", initialQuestionnaire)
	}
	t.Log(initialQuestionnaire)
}
