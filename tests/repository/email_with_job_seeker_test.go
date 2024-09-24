package repository_test

import (
	"testing"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/mock/mock_usecase"
	"go.uber.org/mock/gomock"
)

var (
	emailWithJobSeekerEntity1Subject = "Text Subject 1"
	emailWithJobSeekerEntity1Message = "Text Message 1"
	emailWithJobSeekerEntity2Subject = "Text Subject 2"
	emailWithJobSeekerEntity2Message = "Text Message 2"

	emailWithJobSeekerEntity1 = entity.NewEmailWithJobSeeker(
		1,
		emailWithJobSeekerEntity1Subject,
		emailWithJobSeekerEntity1Message,
		"",
	)

	emailWithJobSeekerEntity2 = entity.NewEmailWithJobSeeker(
		2,
		emailWithJobSeekerEntity2Subject,
		emailWithJobSeekerEntity2Message,
		"",
	)
)

/****************************************************************************************/
// 作成
//
func Test_EmailWithJobSeeker_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockEmailWithJobSeeker := mock_usecase.NewMockEmailWithJobSeekerRepository(mockCtrl)

	mockEmailWithJobSeeker.EXPECT().Create(emailWithJobSeekerEntity1).Return(nil)

	err := mockEmailWithJobSeeker.Create(emailWithJobSeekerEntity1)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}

/****************************************************************************************/
// 複数取得
//
// チャットグループIDからチャットメッセージを取得
func Test_EmailWithJobSeeker_GetByJobSeekerID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockEmailWithJobSeeker := mock_usecase.NewMockEmailWithJobSeekerRepository(mockCtrl)

	mockEmailWithJobSeeker.EXPECT().GetByJobSeekerID(gomock.Any()).DoAndReturn(
		func(jobSeekerID uint) ([]*entity.EmailWithJobSeeker, error) {
			switch jobSeekerID {
			case 1:
				return []*entity.EmailWithJobSeeker{emailWithJobSeekerEntity1, emailWithJobSeekerEntity2}, nil
			case 2:
				return []*entity.EmailWithJobSeeker{emailWithJobSeekerEntity1}, nil
			case 3:
				return []*entity.EmailWithJobSeeker{}, nil
			}
			return nil, nil
		},
	)
	emailWithJobSeekers, err := mockEmailWithJobSeeker.GetByJobSeekerID(uint(1))
	if err != nil {
		t.Fatal("GetByJobSeekerID error!", err)
	}
	if len(emailWithJobSeekers) != 2 {
		t.Fatal("GetByJobSeekerID return value error! len(emailWithJobSeekers) != 2", emailWithJobSeekers)
	}
	if emailWithJobSeekers[0].JobSeekerID != 1 {
		t.Fatal("GetByJobSeekerID return value error! emailWithJobSeekers[0].ID != 2", emailWithJobSeekers)
	}
	if emailWithJobSeekers[1].Subject != emailWithJobSeekerEntity2Subject {
		t.Fatal("GetByJobSeekerID return value error! emailWithJobSeekers[1].Subject != Text Subject 2", emailWithJobSeekers)
	}
	t.Log(emailWithJobSeekers)
}
