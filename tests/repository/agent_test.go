package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/mock/mock_usecase"
	"go.uber.org/mock/gomock"
	"gopkg.in/guregu/null.v4"
)

var agentEntity = &entity.Agent{
	AgentName:                       "株式会社テストエージェント",
	OfficeLocation:                  "110-1111 東京都千代田区9-9-9-",
	Representative:                  "田中太郎",
	Establish:                       "2023-01-01",
	CorporateSiteURL:                "https://autoscout.spaceai.jp",
	PermissionCode:                  "",
	PermissionYear:                  "",
	WorkersCount:                    null.NewInt(10, true),
	PhoneNumber:                     "090-0000-0000",
	InterviewAdjustmentEmail:        "",
	AgreementFileURL:                "",
	LineBotID:                       "",
	LineMessagingChannelSecret:      "",
	LineMessagingChannelAccessToken: "",
	LineLoginChannelID:              "",
	LineLoginChannelSecret:          "",
	SendingAgreementFileURL:         "",
	IsCRMActive:                     true,
	IsAllianceActive:                true,
	IsSendingActive:                 true,
	SendingType:                     null.NewInt(0, true),
}

var agentLineChannelEntity = &entity.AgentLineChannelParam{
	AgentID:                         uint(1),
	AgentUUID:                       uuid.New(),
	LineMessagingChannelSecret:      "c414141414141414141414141414141",
	LineMessagingChannelAccessToken: "xVleZQFrprprprprprprprprprprprprprprprprprprprprprprprprprprprprprprprpFU=",
	LineLoginChannelID:              "169999999",
	LineLoginChannelSecret:          "74747474747474747474747474747474",
}

/****************************************************************************************/
// 作成
//
func Test_Agent_Create(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgent.EXPECT().Create(agentEntity).Return(nil)

	err := mockAgent.Create(agentEntity)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}

/****************************************************************************************/
// 更新
//
func Test_Agent_Update(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgent.EXPECT().Update(uint(1), agentEntity).Return(nil)

	err := mockAgent.Update(uint(1), agentEntity)
	if err != nil {
		t.Fatal("Update error!", err)
	}
}

func Test_Agent_UpdateForAdmin(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	var agentForAdminEntity = &entity.AgentForAdminParam{
		AgentID:          uint(1),
		IsCRMActive:      true,
		IsAllianceActive: true,
		IsSendingActive:  true,
		SendingType:      null.NewInt(0, true),
	}

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgent.EXPECT().UpdateForAdmin(uint(1), *agentForAdminEntity).Return(nil)

	err := mockAgent.UpdateForAdmin(uint(1), *agentForAdminEntity)
	if err != nil {
		t.Fatal("UpdateForAdmin error!", err)
	}
}

func Test_Agent_UpdateLineChannel(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)
	botID := "008293384928"

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgent.EXPECT().UpdateLineChannel(uint(1), botID, *agentLineChannelEntity).Return(nil)

	err := mockAgent.UpdateLineChannel(uint(1), botID, *agentLineChannelEntity)
	if err != nil {
		t.Fatal("UpdateLineChannel error!", err)
	}
}

func Test_Agent_UpdateAgreementFileURL(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	var agentAgreementFileURLEntity = &entity.AgentAgreementFileURLParam{
		AgentID:                 uint(1),
		AgreementFileURL:        "https://doda.jp/guide/rireki/template/file/rirekisyo_a4_mhlw.docx",
		SendingAgreementFileURL: "https://doda.jp/guide/rireki/template/file/rirekisyo_a4_mhlw.xlsx",
	}

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgent.EXPECT().UpdateAgreementFileURL(uint(1), *agentAgreementFileURLEntity).Return(nil)

	err := mockAgent.UpdateAgreementFileURL(uint(1), *agentAgreementFileURLEntity)
	if err != nil {
		t.Fatal("UpdateAgreementFileURL error!", err)
	}
}

/****************************************************************************************/
// 単数取得
//
func Test_Agent_FindByID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します

	// mockAgent.EXPECT().FindByID(gomock.Any()).Return(agentEntity, nil).Times(5)
	mockAgent.EXPECT().FindByID(gomock.Any()).DoAndReturn(
		func(id uint) (*entity.Agent, error) {
			switch id {
			case 1:
				agentEntity.ID = uint(1)
				agentEntity.UUID = uuid.New()
				agentEntity.AgentName = "Test Agent1"
				return agentEntity, nil
			case 2:
				agentEntity.ID = uint(2)
				agentEntity.UUID = uuid.New()
				agentEntity.AgentName = "Test Agent2"
				return agentEntity, nil
			case 3:
				agentEntity.ID = uint(3)
				agentEntity.UUID = uuid.New()
				agentEntity.AgentName = "Test Agent3"
				return agentEntity, nil
			}
			return nil, nil
		},
	)

	agent, err := mockAgent.FindByID(uint(2))
	if err != nil {
		t.Fatal("FindByID error!", err)
	}

	t.Log(agent)
}

func Test_Agent_FindByUUID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します

	var generateUUID = uuid.New()

	agentEntity.UUID = generateUUID
	mockAgent.EXPECT().FindByUUID(generateUUID).Return(agentEntity, nil)

	agent, err := mockAgent.FindByUUID(generateUUID)
	if err != nil {
		t.Fatal("FindByUUID error!", err)
	}

	t.Log(agent)
}

func Test_Agent_FindByAgentStaffID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します

	mockAgent.EXPECT().FindByAgentStaffID(uint(1)).Return(agentEntity, nil)

	agent, err := mockAgent.FindByAgentStaffID(uint(1))
	if err != nil {
		t.Fatal("FindByAgentStaffID error!", err)
	}

	t.Log(agent)
}

func Test_Agent_FindLineChannelByAgentID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します

	mockAgent.EXPECT().FindLineChannelByAgentID(uint(1)).Return(agentLineChannelEntity, nil)

	agentLineChannel, err := mockAgent.FindLineChannelByAgentID(uint(1))
	if err != nil {
		t.Fatal("FindLineChannelByAgentID error!", err)
	}

	t.Log(agentLineChannel)
}

func Test_Agent_FindLineChannelByAgentUUID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します

	var generateUUID = uuid.New()

	mockAgent.EXPECT().FindLineChannelByAgentUUID(generateUUID).Return(agentLineChannelEntity, nil)

	agentLineChannel, err := mockAgent.FindLineChannelByAgentUUID(generateUUID)
	if err != nil {
		t.Fatal("FindLineChannelByAgentUUID error!", err)
	}

	t.Log(agentLineChannel)
}

func Test_Agent_FindLineChannelByBotID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します

	botID := "008293384928"

	mockAgent.EXPECT().FindLineChannelByBotID(botID).Return(agentLineChannelEntity, nil)

	agentLineChannel, err := mockAgent.FindLineChannelByBotID(botID)
	if err != nil {
		t.Fatal("FindLineChannelByBotID error!", err)
	}

	t.Log(agentLineChannel)
}

/****************************************************************************************/
// 複数取得
//
func Test_Agent_GetByIDList(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 引数
	var idList = []uint{1, 2, 3, 4, 5}

	// レスポンス
	var agent1 = agentEntity
	var agent2 = agentEntity
	var agent3 = agentEntity

	agent1.ID = 1
	agent2.ID = 2
	agent3.ID = 3

	agent1.AgentName = "ワン株式会社"
	agent2.AgentName = "ツー株式会社"
	agent3.AgentName = "スリー株式会社"

	var agentListEntity = []*entity.Agent{agent1, agent2, agent3}

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgent.EXPECT().GetByIDList(idList).Return(agentListEntity, nil)

	agentList, err := mockAgent.GetByIDList(idList)
	if err != nil {
		t.Fatal("GetByIDList error!", err)
	}

	t.Log(agentList)
}

func Test_Agent_GetByNotIDList(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 引数
	var idList = []uint{1, 2, 3, 4, 5}

	// レスポンス
	var agent4 = agentEntity
	var agent5 = agentEntity

	agent4.ID = 4
	agent5.ID = 5

	agent4.AgentName = "フォー株式会社"
	agent5.AgentName = "ファイブ株式会社"

	var agentListEntity = []*entity.Agent{agent4, agent5}

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgent.EXPECT().GetByNotIDList(idList).Return(agentListEntity, nil)

	agentList, err := mockAgent.GetByNotIDList(idList)
	if err != nil {
		t.Fatal("GetByNotIDList error!", err)
	}

	t.Log(agentList)
}

func Test_Agent_GetNotAllianceByMyAgentID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// 引数
	var myAgentID = uint(1)

	// レスポンス
	var agent4 = agentEntity
	var agent5 = agentEntity

	agent4.ID = 4
	agent5.ID = 5

	agent4.AgentName = "フォー株式会社"
	agent5.AgentName = "ファイブ株式会社"

	var agentListEntity = []*entity.Agent{agent4, agent5}

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgent.EXPECT().GetNotAllianceByMyAgentID(myAgentID).Return(agentListEntity, nil)

	agentList, err := mockAgent.GetNotAllianceByMyAgentID(myAgentID)
	if err != nil {
		t.Fatal("GetNotAllianceByMyAgentID error!", err)
	}

	t.Log(agentList)
}

func Test_Agent_GetSendingActive(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// レスポンス
	var agent1 = agentEntity
	var agent2 = agentEntity
	var agent3 = agentEntity

	agent1.ID = 1
	agent2.ID = 2
	agent3.ID = 3

	agent1.AgentName = "ワン株式会社"
	agent2.AgentName = "ツー株式会社"
	agent3.AgentName = "スリー株式会社"

	var agentListEntity = []*entity.Agent{agent1, agent2, agent3}

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgent.EXPECT().GetSendingActive().Return(agentListEntity, nil)

	agentList, err := mockAgent.GetSendingActive()
	if err != nil {
		t.Fatal("GetSendingActive error!", err)
	}

	t.Log(agentList)
}

func Test_Agent_All(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgent := mock_usecase.NewMockAgentRepository(mockCtrl)

	// レスポンス
	var agent1 = agentEntity
	var agent2 = agentEntity
	var agent3 = agentEntity

	agent1.ID = 1
	agent2.ID = 2
	agent3.ID = 3

	agent1.AgentName = "ワン株式会社"
	agent2.AgentName = "ツー株式会社"
	agent3.AgentName = "スリー株式会社"

	var agentListEntity = []*entity.Agent{agent1, agent2, agent3}

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgent.EXPECT().All().Return(agentListEntity, nil)

	agentList, err := mockAgent.All()
	if err != nil {
		t.Fatal("All error!", err)
	}

	t.Log(agentList)
}
