package repository_test

import (
	"testing"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/mock/mock_usecase"
	"go.uber.org/mock/gomock"
	"gopkg.in/guregu/null.v4"
)

var (
	initialQuestionnaireDesiredOccupationEntity1 = entity.NewInitialQuestionnaireDesiredOccupation(
		1,
		null.NewInt(1, true),
		null.NewInt(1, true),
	)
)

/****************************************************************************************/
// 作成
//
func Test_InitialQuestionnaireDesiredOccupation_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockInitialQuestionnaireDesiredOccupation := mock_usecase.NewMockInitialQuestionnaireDesiredOccupationRepository(mockCtrl)

	mockInitialQuestionnaireDesiredOccupation.EXPECT().Create(initialQuestionnaireDesiredOccupationEntity1).Return(nil)

	err := mockInitialQuestionnaireDesiredOccupation.Create(initialQuestionnaireDesiredOccupationEntity1)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}
