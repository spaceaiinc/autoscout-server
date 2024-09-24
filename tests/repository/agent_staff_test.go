package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/mock/mock_usecase"
	"go.uber.org/mock/gomock"
	"gopkg.in/guregu/null.v4"
)

var agentStaffEntity = &entity.AgentStaff{
	AgentID:               1,
	FirebaseID:            "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	Authority:             null.NewInt(0, true),
	StaffName:             "田中太郎",
	Furigana:              "タナカタロウ",
	Email:                 "test@example.com",
	StaffPhoneNumber:      "090-1234-5678",
	Department:            "営業部",
	Position:              "営業",
	Remarks:               "備考",
	UsageStatus:           null.NewInt(0, true),
	NotificationJobSeeker: true,
	NotificationUnwatched: true,
}

/****************************************************************************************/
// 作成
//
func Test_AgentStaff_Create(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgentStaff.EXPECT().Create(agentStaffEntity).Return(nil)

	err := mockAgentStaff.Create(agentStaffEntity)
	if err != nil {
		t.Fatal("Create error!", err)
	}
}

/****************************************************************************************/
// 更新
//
func Test_AgentStaff_Update(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgentStaff.EXPECT().Update(uint(1), agentStaffEntity).Return(nil)

	err := mockAgentStaff.Update(uint(1), agentStaffEntity)
	if err != nil {
		t.Fatal("Update error!", err)
	}
}

// スタッフのメールアドレスを更新する
func Test_AgentStaff_UpdateEmail(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgentStaff.EXPECT().UpdateEmail(uint(1), agentStaffEntity.Email).Return(nil)

	err := mockAgentStaff.UpdateEmail(uint(1), agentStaffEntity.Email)
	if err != nil {
		t.Fatal("UpdateEmail error!", err)
	}
}

// スタッフの利用状況を停止する
func Test_AgentStaff_UpdateUsageEnd(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgentStaff.EXPECT().UpdateUsageEnd(uint(1)).Return(nil)

	err := mockAgentStaff.UpdateUsageEnd(uint(1))
	if err != nil {
		t.Fatal("UpdateUsageEnd error!", err)
	}
}

// スタッフの利用状況を再開する
func Test_AgentStaff_UpdateUsageStart(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgentStaff.EXPECT().UpdateUsageStart(uint(1)).Return(nil)

	err := mockAgentStaff.UpdateUsageStart(uint(1))
	if err != nil {
		t.Fatal("UpdateUsageStart error!", err)
	}
}

// スタッフの求職者関連の通知を更新する
func Test_AgentStaff_UpdateNotificationJobSeeker(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgentStaff.EXPECT().UpdateNotificationJobSeeker(uint(1), true).Return(nil)

	err := mockAgentStaff.UpdateNotificationJobSeeker(uint(1), true)
	if err != nil {
		t.Fatal("UpdateNotificationJobSeeker error!", err)
	}
}

// スタッフの未読・未処理関連の通知を更新する
func Test_AgentStaff_UpdateNotificationUnwatched(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgentStaff.EXPECT().UpdateNotificationUnwatched(uint(1), true).Return(nil)

	err := mockAgentStaff.UpdateNotificationUnwatched(uint(1), true)
	if err != nil {
		t.Fatal("UpdateNotificationUnwatched error!", err)
	}
}

// スタッフの権限を更新する
func Test_AgentStaff_UpdateAuthority(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します
	mockAgentStaff.EXPECT().UpdateAuthority(uint(1), uint(1)).Return(nil)

	err := mockAgentStaff.UpdateAuthority(uint(1), uint(1))
	if err != nil {
		t.Fatal("UpdateAuthority error!", err)
	}
}

/****************************************************************************************/
// 単数取得
//
func Test_AgentStaff_FindByID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 作成したmockに対して期待する呼び出しと返り値を定義します
	// EXPECT()では呼び出されたかどうか
	// 「関数名」()ではそのメソッド名が指定した引数で呼び出されたかどうか
	// Return()では返り値を指定します

	mockAgentStaff.EXPECT().FindByID(gomock.Any()).DoAndReturn(
		func(id uint) (*entity.AgentStaff, error) {
			switch id {
			case 1:
				agentStaffEntity.ID = uint(1)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff1"
				return agentStaffEntity, nil
			case 2:
				agentStaffEntity.ID = uint(2)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff2"
				return agentStaffEntity, nil
			case 3:
				agentStaffEntity.ID = uint(3)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff3"
				return agentStaffEntity, nil
			}
			return nil, nil
		},
	)

	agentStaff, err := mockAgentStaff.FindByID(uint(1))
	if err != nil {
		t.Fatal("FindByID error!", err)
	}

	t.Log(agentStaff)
}

// FirebaseIDからスタッフを取得する
func Test_AgentStaff_FindByFirebaseID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)
	// 作成したmockに対して期待する呼び出しと返り値を定義します
	mockAgentStaff.EXPECT().FindByFirebaseID(gomock.Any()).DoAndReturn(
		func(firebaseID string) (*entity.AgentStaff, error) {
			switch firebaseID {
			case "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx":
				agentStaffEntity.ID = uint(1)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff1"
				return agentStaffEntity, nil
			case "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy":
				agentStaffEntity.ID = uint(2)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff2"
				return agentStaffEntity, nil
			case "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz":
				agentStaffEntity.ID = uint(3)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff3"
				return agentStaffEntity, nil
			}
			return nil, nil
		},
	)
	agentStaff, err := mockAgentStaff.FindByFirebaseID("xxxxxxxxxxxxxxxx")
	if err != nil {
		t.Fatal("FindByFirebaseID error!", err)
	}
	t.Log(agentStaff)
}

// 求職者IDからスタッフを取得する
func Test_AgentStaff_FindByJobSeekerID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)
	// 作成したmockに対して期待する呼び出しと返り値を定義します
	mockAgentStaff.EXPECT().FindByJobSeekerID(gomock.Any()).DoAndReturn(
		func(jobSeekerID uint) (*entity.AgentStaff, error) {
			switch jobSeekerID {
			case 1:
				agentStaffEntity.ID = uint(1)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff1"
				return agentStaffEntity, nil
			case 2:
				agentStaffEntity.ID = uint(2)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff2"
				return agentStaffEntity, nil
			case 3:
				agentStaffEntity.ID = uint(3)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff3"
				return agentStaffEntity, nil
			}
			return nil, nil
		},
	)
	agentStaff, err := mockAgentStaff.FindByJobSeekerID(1)
	if err != nil {
		t.Fatal("FindByJobSeekerID error!", err)
	}
	t.Log(agentStaff)
}

// スタッフ名から取得する　※送客求職者の担当者情報を更新する際に使用
func Test_AgentStaff_FindByName(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)
	// 作成したmockに対して期待する呼び出しと返り値を定義します
	mockAgentStaff.EXPECT().FindByName(gomock.Any()).DoAndReturn(
		func(staffName string) (*entity.AgentStaff, error) {
			switch staffName {
			case "Test AgentStaff1":
				agentStaffEntity.ID = uint(1)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff1"
				return agentStaffEntity, nil
			case "Test AgentStaff2":
				agentStaffEntity.ID = uint(2)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff2"
				return agentStaffEntity, nil
			case "Test AgentStaff3":
				agentStaffEntity.ID = uint(3)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff3"
				return agentStaffEntity, nil
			}
			return nil, nil
		},
	)
	agentStaff, err := mockAgentStaff.FindByName("Test AgentStaff1")
	if err != nil {
		t.Fatal("FindByName error!", err)
	}
	t.Log(agentStaff)
}

// IDからスタッフとエージェントのLINE情報を取得する
func Test_AgentStaff_FindStaffAndAgentLine(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)
	// 作成したmockに対して期待する呼び出しと返り値を定義します
	mockAgentStaff.EXPECT().FindStaffAndAgentLine(gomock.Any()).DoAndReturn(
		func(id uint) (*entity.AgentStaff, error) {
			switch id {
			case 1:
				agentStaffEntity.ID = uint(1)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff1"
				return agentStaffEntity, nil
			case 2:
				agentStaffEntity.ID = uint(2)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff2"
				return agentStaffEntity, nil
			case 3:
				agentStaffEntity.ID = uint(3)
				agentStaffEntity.UUID = uuid.New()
				agentStaffEntity.StaffName = "Test AgentStaff3"
				return agentStaffEntity, nil
			}
			return nil, nil
		},
	)
	agentStaff, err := mockAgentStaff.FindStaffAndAgentLine(1)
	if err != nil {
		t.Fatal("FindStaffAndAgentLine error!", err)
	}
	t.Log(agentStaff)
}

// 求職者IDからスタッフ名を取得する
func Test_AgentStaff_FindStaffNameByJobSeekerID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)
	// 作成したmockに対して期待する呼び出しと返り値を定義します
	mockAgentStaff.EXPECT().FindStaffNameByJobSeekerID(gomock.Any()).DoAndReturn(
		func(jobSeekerID uint) (string, error) {
			switch jobSeekerID {
			case 1:
				return "Test AgentStaff1", nil
			case 2:
				return "Test AgentStaff2", nil
			case 3:
				return "Test AgentStaff3", nil
			}
			return "", nil
		},
	)
	staffName, err := mockAgentStaff.FindStaffNameByJobSeekerID(1)
	if err != nil {
		t.Fatal("FindStaffNameByJobSeekerID error!", err)
	}
	t.Log(staffName)
}

/****************************************************************************************/
// 複数取得
//
// エージェントIDからスタッフ一覧を取得する
func Test_AgentStaff_GetByAgentID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 引数
	var agentID uint = 1

	// レスポンス
	var agentStaff1 = agentStaffEntity
	var agentStaff2 = agentStaffEntity
	var agentStaff3 = agentStaffEntity

	var agentStaffListEntity = []*entity.AgentStaff{agentStaff1, agentStaff2, agentStaff3}

	mockAgentStaff.EXPECT().GetByAgentID(agentID).Return(agentStaffListEntity, nil)

	agentList, err := mockAgentStaff.GetByAgentID(agentID)
	if err != nil {
		t.Fatal("GetByAgentID error!", err)
	}

	t.Log(agentList)
}

// エージェントIDが一致かつ売上のManagementIDが不一致のスタッフ一覧を取得する
func Test_AgentStaff_GetByAgentIDAndNotManagementID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 引数
	var agentID uint = 1
	var managementID uint = 1

	// レスポンス
	var agentStaff1 = agentStaffEntity
	var agentStaff2 = agentStaffEntity
	var agentStaff3 = agentStaffEntity

	var agentStaffListEntity = []*entity.AgentStaff{agentStaff1, agentStaff2, agentStaff3}

	mockAgentStaff.EXPECT().GetByAgentIDAndNotManagementID(agentID, managementID).Return(agentStaffListEntity, nil)

	agentList, err := mockAgentStaff.GetByAgentIDAndNotManagementID(agentID, managementID)
	if err != nil {
		t.Fatal("GetByAgentIDAndNotManagementID error!", err)
	}

	t.Log(agentList)
}

// スタッフIDが不一致かつエージェントIDとアライアンスエージェントIDが一致のスタッフ一覧を取得する
func Test_AgentStaff_GetByNotIDAndAgentIDAndAllianceAgentID(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 引数
	var id uint = 1
	var agentID uint = 1
	var allianceAgentID uint = 1

	// レスポンス
	var agentStaff1 = agentStaffEntity
	var agentStaff2 = agentStaffEntity
	var agentStaff3 = agentStaffEntity

	var agentStaffListEntity = []*entity.AgentStaff{agentStaff1, agentStaff2, agentStaff3}

	mockAgentStaff.EXPECT().GetByNotIDAndAgentIDAndAllianceAgentID(id, agentID, allianceAgentID).Return(agentStaffListEntity, nil)

	agentList, err := mockAgentStaff.GetByNotIDAndAgentIDAndAllianceAgentID(id, agentID, allianceAgentID)
	if err != nil {
		t.Fatal("GetByNotIDAndAgentIDAndAllianceAgentID error!", err)
	}

	t.Log(agentList)
}

// usage_status == 利用可能(0)
func Test_AgentStaff_GetByAgentIDAndUsageStatusAvailable(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 引数
	var agentID uint = 1

	// レスポンス
	var agentStaff1 = agentStaffEntity
	var agentStaff2 = agentStaffEntity
	var agentStaff3 = agentStaffEntity

	var agentStaffListEntity = []*entity.AgentStaff{agentStaff1, agentStaff2, agentStaff3}

	mockAgentStaff.EXPECT().GetByAgentIDAndUsageStatusAvailable(agentID).Return(agentStaffListEntity, nil)

	agentList, err := mockAgentStaff.GetByAgentIDAndUsageStatusAvailable(agentID)
	if err != nil {
		t.Fatal("GetByAgentIDAndUsageStatusAvailable error!", err)
	}

	t.Log(agentList)
}

// is_deleted == false
func Test_AgentStaff_GetByAgentIDAndIsDeletedFalse(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// 引数
	var agentID uint = 1

	// レスポンス
	var agentStaff1 = agentStaffEntity
	var agentStaff2 = agentStaffEntity
	var agentStaff3 = agentStaffEntity

	var agentStaffListEntity = []*entity.AgentStaff{agentStaff1, agentStaff2, agentStaff3}

	mockAgentStaff.EXPECT().GetByAgentIDAndIsDeletedFalse(agentID).Return(agentStaffListEntity, nil)

	agentList, err := mockAgentStaff.GetByAgentIDAndIsDeletedFalse(agentID)
	if err != nil {
		t.Fatal("GetByAgentIDAndIsDeletedFalse error!", err)
	}

	t.Log(agentList)
}

func Test_AgentStaff_All(t *testing.T) {
	// モックの呼び出しを管理するControllerを生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// モックの生成
	mockAgentStaff := mock_usecase.NewMockAgentStaffRepository(mockCtrl)

	// レスポンス
	var agentStaff1 = agentStaffEntity
	var agentStaff2 = agentStaffEntity
	var agentStaff3 = agentStaffEntity

	agentStaff1.ID = 1
	agentStaff2.ID = 2
	agentStaff3.ID = 3

	agentStaff1.AgentName = "ワン株式会社"
	agentStaff2.AgentName = "ツー株式会社"
	agentStaff3.AgentName = "スリー株式会社"

	var agentStaffListEntity = []*entity.AgentStaff{agentStaff1, agentStaff2, agentStaff3}

	mockAgentStaff.EXPECT().All().Return(agentStaffListEntity, nil)

	agentStaffList, err := mockAgentStaff.All()
	if err != nil {
		t.Fatal("All error!", err)
	}

	t.Log(agentStaffList)
}
