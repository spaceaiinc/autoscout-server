package repository_test

import (
	"testing"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/mock/mock_usecase"
	"go.uber.org/mock/gomock"
	"gopkg.in/guregu/null.v4"
)

var (
	initialQuestionnaireDesiredIndustryEntity1 = entity.NewInitialQuestionnaireDesiredIndustry(
		1,
		null.NewInt(1, true),
		null.NewInt(1, true),
	)
)

/****************************************************************************************/
// 作成
//
func Test_InitialQuestionnaireDesiredIndustry_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockInitialQuestionnaireDesiredIndustry := mock_usecase.NewMockInitialQuestionnaireDesiredIndustryRepository(mockCtrl)

	mockInitialQuestionnaireDesiredIndustry.EXPECT().Create(initialQuestionnaireDesiredIndustryEntity1).Return(nil)

	err := mockInitialQuestionnaireDesiredIndustry.Create(initialQuestionnaireDesiredIndustryEntity1)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}
