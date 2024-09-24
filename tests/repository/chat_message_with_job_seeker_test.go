package repository_test

import (
	"testing"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/mock/mock_usecase"
	"go.uber.org/mock/gomock"
	"gopkg.in/guregu/null.v4"
)

var (
	chatMessageWithJobSeekerEntity1Message = "Text Message 1"
	chatMessageWithJobSeekerEntity2Message = "Text Message 2"

	chatMessageWithJobSeekerEntity1 = entity.NewChatMessageWithJobSeeker(
		1,
		null.NewInt(0, true),
		chatMessageWithJobSeekerEntity1Message,
		"",
		"",
		"",
		"",
		null.NewInt(0, false),
		"",
		null.NewInt(0, true),
		time.Now().UTC(),
	)

	chatMessageWithJobSeekerEntity2 = entity.NewChatMessageWithJobSeeker(
		2,
		null.NewInt(0, true),
		chatMessageWithJobSeekerEntity2Message,
		"",
		"",
		"",
		"",
		null.NewInt(0, false),
		"",
		null.NewInt(0, true),
		time.Now().UTC(),
	)
)

/****************************************************************************************/
// 作成
//
func Test_ChatMessageWithJobSeeker_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatMessageWithJobSeeker := mock_usecase.NewMockChatMessageWithJobSeekerRepository(mockCtrl)

	mockChatMessageWithJobSeeker.EXPECT().Create(chatMessageWithJobSeekerEntity1).Return(nil)

	err := mockChatMessageWithJobSeeker.Create(chatMessageWithJobSeekerEntity1)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}

/****************************************************************************************/
// 複数取得
//
// チャットグループIDからチャットメッセージを取得
func Test_ChatMessageWithJobSeeker_GetByGroupID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatMessageWithJobSeeker := mock_usecase.NewMockChatMessageWithJobSeekerRepository(mockCtrl)

	mockChatMessageWithJobSeeker.EXPECT().GetByGroupID(gomock.Any()).DoAndReturn(
		func(chatMessageID uint) ([]*entity.ChatMessageWithJobSeeker, error) {
			switch chatMessageID {
			case 1:
				return []*entity.ChatMessageWithJobSeeker{chatMessageWithJobSeekerEntity1, chatMessageWithJobSeekerEntity2}, nil
			case 2:
				return []*entity.ChatMessageWithJobSeeker{chatMessageWithJobSeekerEntity1}, nil
			case 3:
				return []*entity.ChatMessageWithJobSeeker{}, nil
			}
			return nil, nil
		},
	)
	chatMessageWithJobSeekers, err := mockChatMessageWithJobSeeker.GetByGroupID(uint(1))
	if err != nil {
		t.Fatal("GetByGroupID error!", err)
	}
	if len(chatMessageWithJobSeekers) != 2 {
		t.Fatal("GetByGroupID return value error! len(chatMessageWithJobSeekers) != 2", chatMessageWithJobSeekers)
	}
	if chatMessageWithJobSeekers[0].GroupID != 1 {
		t.Fatal("GetByGroupID return value error! chatMessageWithJobSeekers[0].ID != 2", chatMessageWithJobSeekers)
	}
	if chatMessageWithJobSeekers[1].Message != chatMessageWithJobSeekerEntity2Message {
		t.Fatal("GetByGroupID return value error! chatMessageWithJobSeekers[1].Message != Text Message 2", chatMessageWithJobSeekers)
	}
	t.Log(chatMessageWithJobSeekers)
}
