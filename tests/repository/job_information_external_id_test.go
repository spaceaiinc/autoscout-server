package repository_test

import (
	"testing"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/mock/mock_usecase"
	"go.uber.org/mock/gomock"
	"gopkg.in/guregu/null.v4"
)

var (
	jobInformationExternalIDEntity1 = &entity.JobInformationExternalID{
		JobInformationID: 1,
		ExternalType:     null.NewInt(int64(entity.JobInformatinoExternalTypeAgentBank), true),
		ExternalID:       "1234567890",
	}
	jobInformationExternalIDEntity2 = &entity.JobInformationExternalID{
		JobInformationID: 2,
		ExternalType:     null.NewInt(int64(entity.JobInformatinoExternalTypeAgentBank), true),
		ExternalID:       "2345678901",
	}
	jobInformationExternalIDEntity3 = &entity.JobInformationExternalID{
		JobInformationID: 3,
		ExternalType:     null.NewInt(int64(entity.JobInformatinoExternalTypeAgentBank), true),
		ExternalID:       "3456789012",
	}
)

/****************************************************************************************/
// 作成
//
func Test_JobInformationExternalID_Create(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockJobInformationExternalID := mock_usecase.NewMockJobInformationExternalIDRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockJobInformationExternalID.EXPECT().Create(jobInformationExternalIDEntity1).Return(nil)

	err := mockJobInformationExternalID.Create(jobInformationExternalIDEntity1)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}

/****************************************************************************************/
// 削除
//
func Test_JobInformationExternalID_DeleteByJobInformationID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockJobInformationExternalID := mock_usecase.NewMockJobInformationExternalIDRepository(mockCtrl)

	var (
		jobInformationID uint = 1
	)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockJobInformationExternalID.EXPECT().DeleteByJobInformationID(jobInformationID).Return(nil)

	err := mockJobInformationExternalID.DeleteByJobInformationID(jobInformationID)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}

/****************************************************************************************/
// 複数取得
//
func Test_JobInformationExternalID_GetByJobInformationID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockJobInformationExternalID := mock_usecase.NewMockJobInformationExternalIDRepository(mockCtrl)

	var (
		jobInformationID uint = 1

		jobInformationExternalID1 = jobInformationExternalIDEntity1
		jobInformationExternalID2 = jobInformationExternalIDEntity2
		jobInformationExternalID3 = jobInformationExternalIDEntity3
	)

	jobInformationExternalID1.ID = 1
	jobInformationExternalID2.ID = 2
	jobInformationExternalID3.ID = 3

	var jobInformationExternalIDListEntity = []*entity.JobInformationExternalID{jobInformationExternalID1, jobInformationExternalID2, jobInformationExternalID3}

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockJobInformationExternalID.EXPECT().GetByJobInformationID(jobInformationID).Return(jobInformationExternalIDListEntity, nil)

	jobInformationExternalIDList, err := mockJobInformationExternalID.GetByJobInformationID(jobInformationID)
	if err != nil {
		t.Fatal("GetByJobInformationID error!", err)
	}

	t.Log(jobInformationExternalIDList)
}

func Test_JobInformationExternalID_GetByAgentIDAndExternalType(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockJobInformationExternalID := mock_usecase.NewMockJobInformationExternalIDRepository(mockCtrl)

	var (
		agentID      uint = 1
		externalType      = entity.JobInformatinoExternalTypeAgentBank

		jobInformationExternalID1 = jobInformationExternalIDEntity1
		jobInformationExternalID2 = jobInformationExternalIDEntity2
		jobInformationExternalID3 = jobInformationExternalIDEntity3
	)

	jobInformationExternalID1.ID = 1
	jobInformationExternalID2.ID = 2
	jobInformationExternalID3.ID = 3

	var jobInformationExternalIDListEntity = []*entity.JobInformationExternalID{jobInformationExternalID1, jobInformationExternalID2, jobInformationExternalID3}

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockJobInformationExternalID.EXPECT().GetByAgentIDAndExternalType(agentID, externalType).Return(jobInformationExternalIDListEntity, nil)

	jobInformationExternalIDList, err := mockJobInformationExternalID.GetByAgentIDAndExternalType(agentID, externalType)
	if err != nil {
		t.Fatal("GetByAgentIDAndExternalType error!", err)
	}

	t.Log(jobInformationExternalIDList)
}
