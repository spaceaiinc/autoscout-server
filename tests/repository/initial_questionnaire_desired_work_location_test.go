package repository_test

import (
	"testing"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/mock/mock_usecase"
	"go.uber.org/mock/gomock"
	"gopkg.in/guregu/null.v4"
)

var (
	initialQuestionnaireDesiredWorkLocationEntity1 = entity.NewInitialQuestionnaireDesiredWorkLocation(
		1,
		null.NewInt(1, true),
		null.NewInt(1, true),
	)
)

/****************************************************************************************/
// 作成
//
func Test_InitialQuestionnaireDesiredWorkLocation_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockInitialQuestionnaireDesiredWorkLocation := mock_usecase.NewMockInitialQuestionnaireDesiredWorkLocationRepository(mockCtrl)

	mockInitialQuestionnaireDesiredWorkLocation.EXPECT().Create(initialQuestionnaireDesiredWorkLocationEntity1).Return(nil)

	err := mockInitialQuestionnaireDesiredWorkLocation.Create(initialQuestionnaireDesiredWorkLocationEntity1)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}
