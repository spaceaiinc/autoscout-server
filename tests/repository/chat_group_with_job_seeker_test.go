package repository_test

import (
	"testing"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/mock/mock_usecase"
	"go.uber.org/mock/gomock"
)

var (
	chatGroupWithJobSeekerEntity1 = entity.NewChatGroupWithJobSeeker(
		1,
		1,
		true,
	)

	chatGroupWithJobSeekerEntity2 = entity.NewChatGroupWithJobSeeker(
		2,
		2,
		true,
	)

	chatGroupWithJobSeekerEntity3 = entity.NewChatGroupWithJobSeeker(
		3,
		3,
		false,
	)
)

/****************************************************************************************/
// 作成
//
func Test_ChatGroupWithJobSeeker_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	mockChatGroupWithJobSeeker.EXPECT().Create(chatGroupWithJobSeekerEntity1).Return(nil)

	err := mockChatGroupWithJobSeeker.Create(chatGroupWithJobSeekerEntity1)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}

/****************************************************************************************/
// 更新
//
// エージェントの最終閲覧時間を更新
func Test_ChatGroupWithJobSeeker_UpdateAgentLastWatchedAt(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	jobSeekerID := uint(1)

	mockChatGroupWithJobSeeker.EXPECT().UpdateAgentLastWatchedAt(jobSeekerID).Return(nil)

	err := mockChatGroupWithJobSeeker.UpdateAgentLastWatchedAt(jobSeekerID)
	if err != nil {
		t.Fatal("UpdateAgentLastWatchedAt error!", err)
	}
}

// エージェントの最終送信時間を更新
func Test_ChatGroupWithJobSeeker_UpdateAgentLastSendAt(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	jobSeekerID := uint(1)

	mockChatGroupWithJobSeeker.EXPECT().UpdateAgentLastSendAt(jobSeekerID).Return(nil)

	err := mockChatGroupWithJobSeeker.UpdateAgentLastSendAt(jobSeekerID)
	if err != nil {
		t.Fatal("UpdateAgentLastSendAt error!", err)
	}
}

// 求職者の最終閲覧時間と最終送信時間を更新
func Test_ChatGroupWithJobSeeker_UpdateJobSeekerLastWatchedAtAndSendAt(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	jobSeekerID := uint(1)
	sendAt := time.Now().UTC()

	mockChatGroupWithJobSeeker.EXPECT().UpdateJobSeekerLastWatchedAtAndSendAt(jobSeekerID, sendAt).Return(nil)

	err := mockChatGroupWithJobSeeker.UpdateJobSeekerLastWatchedAtAndSendAt(jobSeekerID, sendAt)
	if err != nil {
		t.Fatal("UpdateJobSeekerLastWatchedAtAndSendAt error!", err)
	}
}

// 求職者のLINEの有効/無効を更新
func Test_ChatGroupWithJobSeeker_UpdateJobSeekerLineActive(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	jobSeekerID := uint(1)
	isActive := true

	mockChatGroupWithJobSeeker.EXPECT().UpdateJobSeekerLineActive(jobSeekerID, isActive).Return(nil)

	err := mockChatGroupWithJobSeeker.UpdateJobSeekerLineActive(jobSeekerID, isActive)
	if err != nil {
		t.Fatal("UpdateJobSeekerLineActive error!", err)
	}
}

/****************************************************************************************/
// 単数取得
//
// IDからチャットグループを取得
func Test_ChatGroupWithJobSeeker_FindByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	mockChatGroupWithJobSeeker.EXPECT().FindByID(gomock.Any()).DoAndReturn(
		func(id uint) (*entity.ChatGroupWithJobSeeker, error) {
			switch id {
			case 1:
				chatGroupWithJobSeekerEntity1.ID = uint(1)
				return chatGroupWithJobSeekerEntity1, nil
			case 2:
				chatGroupWithJobSeekerEntity2.ID = uint(2)
				return chatGroupWithJobSeekerEntity2, nil
			case 3:
				chatGroupWithJobSeekerEntity3.ID = uint(3)
				return chatGroupWithJobSeekerEntity3, nil
			}
			return nil, nil
		},
	)
	chatGroupWithJobSeeker, err := mockChatGroupWithJobSeeker.FindByID(uint(1))
	if err != nil {
		t.Fatal("FindByID error!", err)
	}
	if chatGroupWithJobSeeker.ID != 1 {
		t.Fatal("FindByID return value error!", chatGroupWithJobSeeker)
	}
	t.Log(chatGroupWithJobSeeker)
}

// エージェントIDと求職者IDからチャットグループを取得
func Test_ChatGroupWithJobSeeker_FindByAgentIDAndJobSeekerID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	mockChatGroupWithJobSeeker.EXPECT().FindByAgentIDAndJobSeekerID(gomock.Any(), gomock.Any()).DoAndReturn(
		func(agentID, jobSeekerID uint) (*entity.ChatGroupWithJobSeeker, error) {
			if agentID == 1 && jobSeekerID == 1 {
				chatGroupWithJobSeekerEntity1.ID = uint(1)
				return chatGroupWithJobSeekerEntity1, nil
			} else if agentID == 2 && jobSeekerID == 2 {
				chatGroupWithJobSeekerEntity2.ID = uint(2)
				return chatGroupWithJobSeekerEntity2, nil
			} else if agentID == 3 && jobSeekerID == 3 {
				chatGroupWithJobSeekerEntity3.ID = uint(3)
				return chatGroupWithJobSeekerEntity3, nil
			}
			return nil, nil
		},
	)

	chatGroupWithJobSeeker, err := mockChatGroupWithJobSeeker.FindByAgentIDAndJobSeekerID(uint(1), uint(1))
	if err != nil {
		t.Fatal("FindByAgentIDAndJobSeekerID error!", err)
	}
	if chatGroupWithJobSeeker.ID != 1 {
		t.Fatal("FindByAgentIDAndJobSeekerID return value error!", chatGroupWithJobSeeker)
	}
	t.Log(chatGroupWithJobSeeker)
}

/****************************************************************************************/
// 複数取得
//
// エージェントIDからチャットグループを取得
func Test_ChatGroupWithJobSeeker_GetByAgentID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	mockChatGroupWithJobSeeker.EXPECT().GetByAgentID(gomock.Any()).DoAndReturn(
		func(agentID uint) ([]*entity.ChatGroupWithJobSeeker, error) {
			switch agentID {
			case 1:
				return []*entity.ChatGroupWithJobSeeker{chatGroupWithJobSeekerEntity1, chatGroupWithJobSeekerEntity2}, nil
			case 2:
				return []*entity.ChatGroupWithJobSeeker{chatGroupWithJobSeekerEntity3}, nil
			case 3:
				return []*entity.ChatGroupWithJobSeeker{}, nil
			}
			return nil, nil
		},
	)

	chatGroupWithJobSeekers, err := mockChatGroupWithJobSeeker.GetByAgentID(uint(1))
	if err != nil {
		t.Fatal("GetByAgentID error!", err)
	}
	if len(chatGroupWithJobSeekers) != 2 {
		t.Fatal("GetByAgentID return value error!", chatGroupWithJobSeekers)
	}
	t.Log(chatGroupWithJobSeekers)
}

// エージェントIDとフリーワードからチャットグループを取得
func Test_ChatGroupWithJobSeeker_GetByAgentIDAndFreeWord(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	mockChatGroupWithJobSeeker.EXPECT().GetByAgentIDAndFreeWord(gomock.Any(), gomock.Any()).DoAndReturn(
		func(agentID uint, freeWord string) ([]*entity.ChatGroupWithJobSeeker, error) {
			switch agentID {
			case 1:
				return []*entity.ChatGroupWithJobSeeker{chatGroupWithJobSeekerEntity1, chatGroupWithJobSeekerEntity2}, nil
			case 2:
				return []*entity.ChatGroupWithJobSeeker{chatGroupWithJobSeekerEntity3}, nil
			case 3:
				return []*entity.ChatGroupWithJobSeeker{}, nil
			}
			return nil, nil
		},
	)
	chatGroupWithJobSeekers, err := mockChatGroupWithJobSeeker.GetByAgentIDAndFreeWord(uint(1), "freeWord")
	if err != nil {
		t.Fatal("GetByAgentIDAndFreeWord error!", err)
	}
	if len(chatGroupWithJobSeekers) != 2 {
		t.Fatal("GetByAgentIDAndFreeWord return value error!", chatGroupWithJobSeekers)
	}
	t.Log(chatGroupWithJobSeekers)
}

// 全てのチャットグループを取得
func Test_ChatGroupWithJobSeeker_All(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	mockChatGroupWithJobSeeker.EXPECT().All().DoAndReturn(
		func() ([]*entity.ChatGroupWithJobSeeker, error) {
			return []*entity.ChatGroupWithJobSeeker{chatGroupWithJobSeekerEntity1, chatGroupWithJobSeekerEntity2, chatGroupWithJobSeekerEntity3}, nil
		},
	)
	chatGroupWithJobSeekers, err := mockChatGroupWithJobSeeker.All()
	if err != nil {
		t.Fatal("All error!", err)
	}
	if len(chatGroupWithJobSeekers) != 3 {
		t.Fatal("All return value error!", chatGroupWithJobSeekers)
	}
	t.Log(chatGroupWithJobSeekers)
}

// エージェントIDから通知数を取得
func Test_ChatGroupWithJobSeeker_CountNotificationByAgentID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockChatGroupWithJobSeeker := mock_usecase.NewMockChatGroupWithJobSeekerRepository(mockCtrl)

	mockChatGroupWithJobSeeker.EXPECT().CountNotificationByAgentID(gomock.Any()).DoAndReturn(
		func(agentID uint) (float64, error) {
			switch agentID {
			case 1:
				return float64(10), nil
			case 2:
				return float64(20), nil
			case 3:
				return float64(30), nil
			}
			return 0, nil
		},
	)
	count, err := mockChatGroupWithJobSeeker.CountNotificationByAgentID(uint(1))
	if err != nil {
		t.Fatal("CountNotificationByAgentID error!", err)
	}
	if count != 10 {
		t.Fatal("CountNotificationByAgentID return value error!", count)
	}
	t.Log(count)
}
