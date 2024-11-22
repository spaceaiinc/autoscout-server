package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"gopkg.in/guregu/null.v4"
)

/****************************************************************************************/
/// エージェント関連
//

type AgentRepository interface {
	/** 作成 */
	// エージェントを作成する
	Create(agent *entity.Agent) error

	/** 更新 */
	// エージェントを更新する
	Update(id uint, agent *entity.Agent) error

	// 管理側からCRMや送客機能の利用を更新する
	UpdateForAdmin(id uint, agent entity.AgentForAdminParam) error

	// LINE情報を更新する
	UpdateLineChannel(id uint, botID string, agentLineChannel entity.AgentLineChannelParam) error

	// 同意書URLを更新する
	UpdateAgreementFileURL(id uint, agentAgreement entity.AgentAgreementFileURLParam) error

	/** 単数取得 */
	// IDからエージェントを取得する
	FindByID(id uint) (*entity.Agent, error)

	// UUIDからエージェントを取得する
	FindByUUID(uuid uuid.UUID) (*entity.Agent, error)

	// スタッフIDからエージェントを取得する
	FindByAgentStaffID(agentStaffID uint) (*entity.Agent, error)

	// エージェントIDからLINE連携情報を取得する
	FindLineChannelByAgentID(agentID uint) (*entity.AgentLineChannelParam, error)

	// エージェントUUIDからLINE連携情報を取得する
	FindLineChannelByAgentUUID(agentUUID uuid.UUID) (*entity.AgentLineChannelParam, error)

	// エージェントIDからLINE連携情報を取得する
	FindLineChannelByBotID(botID string) (*entity.AgentLineChannelParam, error)

	/** 複数取得 */
	// 指定idのエージェントを全て取得
	GetByIDList(idList []uint) ([]*entity.Agent, error)

	// 指定id以外のエージェントを全て取得
	GetByNotIDList(idList []uint) ([]*entity.Agent, error)

	// 自社エージェントID以外のアライアンス未締結のエージェント一覧を取得
	GetNotAllianceByMyAgentID(myAgentID uint) ([]*entity.Agent, error)

	// 送客利用のエージェントをすべて取得する
	GetSendingActive() ([]*entity.Agent, error)

	// 全てのエージェント一覧を取得する
	All() ([]*entity.Agent, error)
}

type AgentStaffRepository interface {
	/** 作成 */
	// スタッフを作成する
	Create(agentStaff *entity.AgentStaff) error

	/** 更新 */
	// スタッフを更新する
	Update(id uint, agent *entity.AgentStaff) error

	// スタッフのメールアドレスを更新する
	UpdateEmail(id uint, email string) error

	// スタッフの利用状況を停止する
	UpdateUsageEnd(id uint) error

	// スタッフの利用状況を再開する
	UpdateUsageStart(id uint) error

	// スタッフの求職者関連の通知を更新する
	UpdateNotificationJobSeeker(id uint, notificationJobSeeker bool) error

	// スタッフの未読・未処理関連の通知を更新する
	UpdateNotificationUnwatched(id uint, notificationUnwatched bool) error

	// スタッフの権限を更新する
	UpdateAuthority(id uint, authority uint) error

	// 最終ログイン時間を更新
	UpdateLastLogin(id uint, lastLoginTime time.Time) error

	/** 削除 */
	// スタッフを削除する　※is_deletedをtrueにする
	Delete(id uint) error

	/** 単数取得 */
	// IDからスタッフを取得する
	FindByID(id uint) (*entity.AgentStaff, error)

	// FirebaseIDからスタッフを取得する
	FindByFirebaseID(firebaseID string) (*entity.AgentStaff, error)

	// 求職者IDからスタッフを取得する
	FindByJobSeekerID(jobSeekerID uint) (*entity.AgentStaff, error)

	// スタッフ名から取得する　※送客求職者の担当者情報を更新する際に使用
	FindByName(staffName string) (*entity.AgentStaff, error)

	// IDからスタッフとエージェントのLINE情報を取得する
	FindStaffAndAgentLine(id uint) (*entity.AgentStaff, error)

	// 求職者IDからスタッフ名を取得する
	FindStaffNameByJobSeekerID(jobSeekerID uint) (string, error)

	/** 複数取得 */
	// エージェントIDからスタッフ一覧を取得する
	GetByAgentID(agentID uint) ([]*entity.AgentStaff, error)

	// エージェントIDが一致かつ売上のManagementIDが不一致のスタッフ一覧を取得する
	GetByAgentIDAndNotManagementID(agentID, managementID uint) ([]*entity.AgentStaff, error)

	// スタッフIDが不一致かつエージェントIDとアライアンスエージェントIDが一致のスタッフ一覧を取得する
	GetByNotIDAndAgentIDAndAllianceAgentID(id, agentID, allianceAgentID uint) ([]*entity.AgentStaff, error)

	// usage_status == 利用可能(0)
	GetByAgentIDAndUsageStatusAvailable(agentID uint) ([]*entity.AgentStaff, error)

	// is_deleted == false
	GetByAgentIDAndIsDeletedFalse(agentID uint) ([]*entity.AgentStaff, error)

	// 全てのスタッフ一覧を取得する
	All() ([]*entity.AgentStaff, error)
}

/****************************************************************************************/
// エージェントのアライアンス関連
//
type AgentAllianceRepository interface {
	/** 作成 */
	// アライアンス情報を作成する
	Create(agentAlliance *entity.AgentAlliance) error

	/** 更新 */
	// アライアンス情報を更新する
	Update(id uint, agentAlliance *entity.AgentAlliance) error

	// アライアンスのリクエストを更新する
	UpdateAgentRequest(id uint, agentID uint) error

	/** 単数取得 */
	// IDからアライアンス情報を取得する
	FindByID(id uint) (*entity.AgentAlliance, error)

	// エージェントIDからアライアンス情報を取得する
	FindByAgentID(agent1ID uint, agent2ID uint) (*entity.AgentAlliance, error)

	/** 複数取得 */
	// 申請完了済みのアライアンス情報を取得する
	GetByAgentIDAndRequestDone(agentID uint) ([]*entity.AgentAlliance, error)

	// 自社IDと他社IDリストからアライアンス情報を取得する
	GetByMyAgentIDAndOtherIDList(myAgentID uint, otherAgentIDList []uint) ([]*entity.AgentAlliance, error)

	// 全てのアライアンス情報を取得する
	All() ([]*entity.AgentAlliance, error)
}

/****************************************************************************************/
// エージェントのRPAロボット scout_serviceを子テーブルとして持つ
//
type AgentRobotRepository interface {
	/** 作成 */
	// RPAロボットを作成する
	Create(robot *entity.AgentRobot) error

	/** 更新 */
	// RPAロボットを更新する
	Update(id uint, robot *entity.AgentRobot) error

	/** 削除 */
	// RPAロボットを削除する
	Delete(id uint) error

	/** 単数取得 */
	// IDからRPAロボットを取得する
	FindByID(id uint) (*entity.AgentRobot, error)

	/** 複数取得 */
	// エージェントIDからRPAロボット一覧を取得する
	GetByAgentID(agentID uint) ([]*entity.AgentRobot, error)

	// 全てのRPAロボット一覧を取得する
	All() ([]*entity.AgentRobot, error)
}

/****************************************************************************************/
// エージェントの流入経路マスタ
//
type AgentInflowChannelOptionRepository interface {
	/** 作成 */
	// 流入経路マスタを作成する
	Create(inflowChannelOption *entity.AgentInflowChannelOption) error

	/** 更新 */
	// 流入経路マスタを更新する
	Update(id uint, inflowChannelOption *entity.AgentInflowChannelOption) error

	/** 単数取得 */
	// IDから流入経路マスタを取得する
	FindByID(id uint) (*entity.AgentInflowChannelOption, error)

	/** 複数取得 */
	// エージェントIDから流入経路マスタ一覧を取得する
	GetByAgentID(agentID uint) ([]*entity.AgentInflowChannelOption, error)
}

/****************************************************************************************/

/****************************************************************************************/
/// 求人企業関連
//
type EnterpriseProfileRepository interface {
	/** 作成 */
	// 求人企業を作成する
	Create(enterprise *entity.EnterpriseProfile) error

	/** 更新 */
	// 求人企業を更新する
	Update(id uint, enterprise *entity.EnterpriseProfile) error

	// 担当者引き継ぎ用
	UpdateAgentStaffID(id, agentStaffID uint) error

	/** 削除 */
	// 求人企業を削除する
	Delete(id uint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.EnterpriseProfile, error)

	// billing_address_idから取得
	FindByJobInformationID(jobInformationID uint) (*entity.EnterpriseProfile, error)

	/** 複数取得 */
	// エージェントスタッフIDから求人企業一覧を取得する
	GetByAgentStaffID(agentStaffID uint) ([]*entity.EnterpriseProfile, error)

	// エージェントIDから求人企業一覧を取得する
	GetByAgentID(agentID uint) ([]*entity.EnterpriseProfile, error)

	// 絞り込み検索
	GetByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.EnterpriseProfile, error)

	// 全てのレコードを取得
	All() ([]*entity.EnterpriseProfile, error)

	/** その他 */
	// 郵便番号の下4桁チェック
	CheckPostCode(postCodeUnder string) (*entity.EnterpriseProfile, error)
}

type EnterpriseIndustryRepository interface {
	/** 作成 */
	// 求人企業の業界を作成する
	Create(industry *entity.EnterpriseIndustry) error

	/** 削除 */
	// 求人企業の業界を削除する
	DeleteByEnterpriseID(enterpriseID uint) error

	/** 複数取得 */
	// 企業IDから求人企業の業界を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.EnterpriseIndustry, error)

	// 請求先IDから取得
	GetByBillingAddressID(enterpriseID uint) ([]*entity.EnterpriseIndustry, error)

	// 求人IDから取得
	GetByJobInformationID(jobInformationID uint) ([]*entity.EnterpriseIndustry, error)

	// エージェントスタッフIDから業界一覧を取得する
	GetByAgentStaffID(agentStaffID uint) ([]*entity.EnterpriseIndustry, error)

	// エージェントIDから業界一覧を取得する
	GetByAgentID(agentID uint) ([]*entity.EnterpriseIndustry, error)

	// IDリストから業界一覧を取得
	GetByEnterpriseIDList(enterpriseIDList []uint) ([]*entity.EnterpriseIndustry, error)

	// 求人IDリストから業界一覧を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.EnterpriseIndustry, error)

	// 全てのレコードを取得
	All() ([]*entity.EnterpriseIndustry, error)
}

// 企業の参考資料
type EnterpriseReferenceMaterialRepository interface {
	/** 作成 */
	Create(referenceMaterial *entity.EnterpriseReferenceMaterial) error

	/** 更新 */
	// 企業IDから求人企業資料を更新
	UpdateByEnterpriseID(enterpriseID uint, referenceMaterial *entity.EnterpriseReferenceMaterial) error

	// 資料IDから求人企業資料を更新
	UpdateMaterialTypeByID(id, materialType uint) error

	/** 削除 */
	DeleteByEnterpriseID(enterpriseID uint) error

	/** 単数取得 */
	// 企業IDから求人企業資料を取得
	FindByEnterpriseID(enterpriseID uint) (*entity.EnterpriseReferenceMaterial, error)

	/** 複数取得 */
	// スタッフIDから求人企業資料を取得
	GetByAgentStaffID(agentStaffID uint) ([]*entity.EnterpriseReferenceMaterial, error)

	// エージェントIDから求人企業資料を取得
	GetByAgentID(agentID uint) ([]*entity.EnterpriseReferenceMaterial, error)

	// IDリストから一覧を取得
	GetByEnterpriseIDList(enterpriseIDList []uint) ([]*entity.EnterpriseReferenceMaterial, error)

	/** 全てのレコードを取得 */
	All() ([]*entity.EnterpriseReferenceMaterial, error)
}

// 企業の追加情報
type EnterpriseActivityRepository interface {
	/** 作成 */
	Create(activity *entity.EnterpriseActivity) error

	/** 複数取得 */
	// 企業IDから追加情報を取得
	GetByEnterpriseID(enterpriseID uint) ([]*entity.EnterpriseActivity, error)

	// 請求先IDから取得
	GetByAgentStaffID(agentStaffID uint) ([]*entity.EnterpriseActivity, error)

	// エージェントIDから取得
	GetByAgentID(agentID uint) ([]*entity.EnterpriseActivity, error)

	// IDリストから一覧を取得
	GetByEnterpriseIDList(enterpriseIDList []uint) ([]*entity.EnterpriseActivity, error)

	/** 全てのレコードを取得 */
	All() ([]*entity.EnterpriseActivity, error)
}

/****************************************************************************************/

/****************************************************************************************/
/// 請求先関連
//
type BillingAddressRepository interface {
	/** 作成 */
	// 請求先を作成する
	Create(billingAddress *entity.BillingAddress) error

	/** 更新 */
	// 請求先を更新する
	Update(id uint, billingAddress *entity.BillingAddress) error

	// agent_staff_idを更新する。引き継ぎ用
	UpdateAgentStaffID(id, agentStaffID uint) error

	// 指定の請求先idリストのagent_staff_idを更新
	UpdateAgentStaffIDByBillingAddressIDList(idList []uint, agentStaffID uint) error

	/** 削除 */
	// 請求先を削除する
	Delete(billingAddressID uint) error

	/** 単数取得 */
	// IDから請求先を取得する
	FindByID(billingAddressID uint) (*entity.BillingAddress, error)

	/** 複数取得 */
	// 求人企業IDから請求先一覧を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.BillingAddress, error)

	// 担当者IDから請求先一覧を取得する
	GetByAgentStaffID(agentStaffID uint) ([]*entity.BillingAddress, error)

	// エージェントIDから請求先一覧を取得する
	GetByAgentID(agentID uint) ([]*entity.BillingAddress, error)

	// 絞り込み検索
	GetByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.BillingAddress, error)

	// 全てのレコードを取得
	All() ([]*entity.BillingAddress, error)
}

type BillingAddressHRStaffRepository interface {
	/** 作成 */
	// 請求先の人事担当者を作成する
	Create(hrStaff *entity.BillingAddressHRStaff) error

	/** 削除 */
	// 請求先の人事担当者を削除する
	DeleteByBillingAddressID(billingAddressID uint) error

	/** 複数取得 */
	// 請求先IDから人事担当者一覧を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.BillingAddressHRStaff, error)

	// 企業IDから人事担当者一覧を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.BillingAddressHRStaff, error)

	// 全てのレコードを取得
	All() ([]*entity.BillingAddressHRStaff, error)
}

type BillingAddressRAStaffRepository interface {
	/** 作成 */
	// 請求先のRA担当者を作成する
	Create(raStaff *entity.BillingAddressRAStaff) error

	/** 削除 */
	// 請求先のRA担当者を削除する
	DeleteByBillingAddressID(billingAddressID uint) error

	/** 複数取得 */
	// 請求先IDからRA担当者一覧を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.BillingAddressRAStaff, error)

	// 企業IDからRA担当者一覧を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.BillingAddressRAStaff, error)

	// 全てのレコードを取得
	All() ([]*entity.BillingAddressRAStaff, error)
}

/****************************************************************************************/

/****************************************************************************************/
/// 求人関連
//
type JobInformationRepository interface {
	/** 作成 */
	// 求人を作成する
	Create(jobInformation *entity.JobInformation) error

	/** 更新 */
	// 求人を更新する
	Update(id uint, jobInformation *entity.JobInformation) error

	// 求人の募集状況のみを更新
	UpdateRecruitmentState(id uint, recruitmentState null.Int) error

	// 求人の面接確約フラグを更新する
	UpdateIsGuaranteedInterview(id uint, isGuaranteedInterview bool) error

	/** 削除 */
	// 求人を削除する
	Delete(id uint) error

	// 指定請求先IDの求人を削除
	DeleteByBillingAddressID(billingAddressID uint) error

	/** 単数取得 */
	// 指定IDの求人を取得する
	FindByID(jobInformationD uint) (*entity.JobInformation, error)

	// 指定UUIDの求人を取得する
	FindByUUID(uuid uuid.UUID) (*entity.JobInformation, error)

	// 指定タスクグループUUIDの求人を取得する
	FindByTaskGroupUUID(uuid uuid.UUID) (*entity.JobInformation, error)

	/** 複数取得 */
	// 指定請求先IDの求人を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformation, error)

	// 指定求人企業IDの求人を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformation, error)

	// IDリストの求人を取得
	GetByIDList(idList []uint) ([]*entity.JobInformation, error)

	// IDリストの求人を取得 タスク用
	GetByIDListForTask(idList []uint) ([]*entity.JobInformation, error)

	// 指定エージェントIDの求人を取得
	GetByAgentID(agentID uint) ([]*entity.JobInformation, error)

	// エージェントIDとアライアンスリストから求人を取得する
	GetByAgentIDAndAlliance(agentID uint, agentAllianceList []*entity.AgentAlliance) ([]*entity.JobInformation, error)

	// エージェントIDとフリーワードから求人を取得する
	GetByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.JobInformation, error)

	// エージェントIDとフリーワードとアライアンスリストから求人を取得する
	GetByAgentIDAndFreeWordAndAlliance(agentID uint, freeWord string, agentAllianceList []*entity.AgentAlliance) ([]*entity.JobInformation, error)

	// フリーワードからエージェントIDのアクティブな求人
	GetActiveByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.JobInformation, error)

	// エージェントIDと自社全て+他社の請求フェーズが契約済みのアクティブ求人
	GetActiveAllByAgentID(agentID uint) ([]*entity.JobInformation, error)

	// エージェントIDと自社全て+他社の請求フェーズが契約済みのアクティブ求人 外部求人を除く
	GetActiveAllByAgentIDWithoutExternal(agentID uint) ([]*entity.JobInformation, error)

	GetActiveAllByAgentIDAndDiagnosisParamWithoutExternal(agentID uint, diagnosisParam entity.DiagnosisParam) ([]*entity.JobInformation, error)

	// エージェントIDとフリーワードから自社全て+他社の請求フェーズが契約済みのアクティブ求人
	GetActiveAllByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.JobInformation, error)

	// エージェントIDとフリーワードとアライアンスから自社エージェント除く有効なアライアンス求人を取得
	GetActiveAllianceByAgentIDAndFreeWordAndAlliance(agentID uint, freeWord string, agentAllianceList []*entity.AgentAlliance) ([]*entity.JobInformation, error)

	// 求職者uuidを使って打診されている求人を取得
	GetAlreadySoundOutByJobSeekerUUID(jobSeekerUUID uuid.UUID) ([]*entity.JobInformation, error)

	// AgentIDと外部媒体タイプから求人を取得
	GetByAgentIDAndExternalType(agentID uint, externalType entity.JobInformatinoExternalType) ([]*entity.JobInformation, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformation, error)
}

// 募集対象
type JobInformationTargetRepository interface {
	/** 作成 */
	// 募集対象を作成する
	Create(target *entity.JobInformationTarget) error

	/** 削除 */
	// 募集対象を削除する
	DeleteByJobInformationID(jobInformationID uint) error

	/** 複数取得 */
	// 指定求人IDの募集対象を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationTarget, error)

	// 指定請求先IDの募集対象を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationTarget, error)

	// 指定企業IDの募集対象を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationTarget, error)

	// 指定エージェントIDの募集対象を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationTarget, error)

	// 指定求人IDリストの募集対象を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationTarget, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationTarget, error)
}

// 募集職種
type JobInformationOccupationRepository interface {
	/** 作成 */
	// 職種を作成する
	Create(occupation *entity.JobInformationOccupation) error

	/** 削除 */
	// 職種を削除する
	DeleteByJobInformationID(jobInformationID uint) error

	/** 複数取得 */
	// 指定求人IDの職種を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationOccupation, error)

	// 指定請求先IDの職種を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationOccupation, error)

	// 指定企業IDの職種を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationOccupation, error)

	// 指定エージェントIDの職種を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationOccupation, error)

	// 指定求人IDリストの職種を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationOccupation, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationOccupation, error)
}

// 求人の特徴
type JobInformationFeatureRepository interface {
	/** 作成 */
	// 特徴を作成する
	Create(feature *entity.JobInformationFeature) error

	/** 削除 */
	// 特徴を削除する
	DeleteByJobInformationID(jobInformationID uint) error

	/** 複数取得 */
	// 指定求人IDの特徴を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationFeature, error)

	// 指定請求先IDの特徴を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationFeature, error)

	// 指定企業IDの特徴を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationFeature, error)

	// 指定エージェントIDの特徴を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationFeature, error)

	// 指定求人IDリストの特徴を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationFeature, error)

	// 求職者ページ表示用項目のみ取得 0:業界未経験OK, 1:職種未経験OK, 2:業界・職種未経験OK, 6:転勤なし
	GetByJobInformationIDListForGuestJobSeeker(jobInformationIDList []uint) ([]*entity.JobInformationFeature, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationFeature, error)
}

// 都道府県
type JobInformationPrefectureRepository interface {
	/** 作成 */
	// 都道府県を作成する
	Create(prefecture *entity.JobInformationPrefecture) error

	/** 削除 */
	// 都道府県を削除する
	DeleteByJobInformationID(jobInformationID uint) error

	/** 複数取得 */
	// 指定求人IDの都道府県を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationPrefecture, error)

	// 指定請求先IDの都道府県を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationPrefecture, error)

	// 指定企業IDの都道府県を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationPrefecture, error)

	// 指定エージェントIDの都道府県を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationPrefecture, error)

	// 指定求人IDリストの都道府県を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationPrefecture, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationPrefecture, error)
}

// 企業の魅力
type JobInformationWorkCharmPointRepository interface {
	/** 作成 */
	// 企業の魅力を作成する
	Create(workCharmPoint *entity.JobInformationWorkCharmPoint) error

	/** 削除 */
	// 企業の魅力を削除する
	DeleteByJobInformationID(jobInformationID uint) error

	/** 複数取得 */
	// 指定求人IDの企業の魅力を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationWorkCharmPoint, error)

	// 指定請求先IDの企業の魅力を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationWorkCharmPoint, error)

	// 指定企業IDの企業の魅力を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationWorkCharmPoint, error)

	// 指定エージェントIDの企業の魅力を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationWorkCharmPoint, error)

	// 指定求人IDリストの企業の魅力を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationWorkCharmPoint, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationWorkCharmPoint, error)
}

// 雇用形態
type JobInformationEmploymentStatusRepository interface {
	/** 作成 */
	// 雇用形態を作成する
	Create(employmentStatus *entity.JobInformationEmploymentStatus) error

	/** 削除 */
	// 雇用形態を削除する
	DeleteByJobInformationID(jobInformationID uint) error

	/** 複数取得 */
	// 指定求人IDの雇用形態を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationEmploymentStatus, error)

	// 指定請求先IDの雇用形態を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationEmploymentStatus, error)

	// 指定企業IDの雇用形態を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationEmploymentStatus, error)

	// 指定エージェントIDの雇用形態を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationEmploymentStatus, error)

	// 指定求人IDリストの雇用形態を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationEmploymentStatus, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationEmploymentStatus, error)
}

// 必要条件　requiredCondition
type JobInformationRequiredConditionRepository interface {
	/** 作成 */
	// 必要条件を作成する
	Create(requiredCondition *entity.JobInformationRequiredCondition) error

	/** 削除 */
	// 必要条件を削除する
	DeleteByJobInformationID(jobInformationID uint) error

	/** 複数取得 */
	// 指定求人IDの必要条件を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredCondition, error)

	// 指定請求先IDの必要条件を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredCondition, error)

	// 指定企業IDの必要条件を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredCondition, error)

	// 指定エージェントIDの必要条件を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredCondition, error)

	// 指定求人IDリストの必要条件を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredCondition, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredCondition, error)
}

// 必要資格
type JobInformationRequiredLicenseRepository interface {
	/** 作成 */
	// 必要資格を作成する
	Create(requiredLicense *entity.JobInformationRequiredLicense) error

	/** 複数取得 */
	// 指定求人IDの必要資格を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredLicense, error)

	// 指定請求先IDの必要資格を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredLicense, error)

	// 指定企業IDの必要資格を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredLicense, error)

	// 指定エージェントIDの必要資格を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredLicense, error)

	// 指定求人IDリストの必要資格を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredLicense, error)

	// 求人リストと資格タイプから必要資格を取得
	GetByJobInformationIDListAndLicenceTypeList(jobInformationIDList []uint, licenceTypeList []null.Int) ([]*entity.JobInformationRequiredLicense, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredLicense, error)
}

// PCスキル
type JobInformationRequiredPCToolRepository interface {
	/** 作成 */
	// 必要PCスキルを作成する
	Create(requiredPCTool *entity.JobInformationRequiredPCTool) error

	/** 複数取得 */
	// 指定求人IDの必要PCスキルを取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredPCTool, error)

	// 指定請求先IDの必要PCスキルを取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredPCTool, error)

	// 指定企業IDの必要PCスキルを取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredPCTool, error)

	// 指定エージェントIDの必要PCスキルを取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredPCTool, error)

	// 指定求人IDリストの必要PCスキルを取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredPCTool, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredPCTool, error)
}

// 必要言語スキル
type JobInformationRequiredLanguageRepository interface {
	/** 作成 */
	// 必要言語スキルを作成する
	Create(requiredLanguage *entity.JobInformationRequiredLanguage) error

	/** 複数取得 */
	// 指定求人IDの必要言語スキルを取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredLanguage, error)

	// 指定請求先IDの必要言語スキルを取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredLanguage, error)

	// 指定企業IDの必要言語スキルを取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredLanguage, error)

	// 指定エージェントIDの必要言語スキルを取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredLanguage, error)

	// 指定求人IDリストの必要言語スキルを取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredLanguage, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredLanguage, error)
}

// 必要言語スキルの種類　requiredLanguageType
type JobInformationRequiredLanguageTypeRepository interface {
	/** 作成 */
	// 必要言語スキルの種類を作成する
	Create(requiredLanguageType *entity.JobInformationRequiredLanguageType) error

	/** 複数取得 */
	// 指定求人IDの必要言語スキルの種類を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredLanguageType, error)

	// 指定請求先IDの必要言語スキルの種類を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredLanguageType, error)

	// 指定企業IDの必要言語スキルの種類を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredLanguageType, error)

	// 指定エージェントIDの必要言語スキルの種類を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredLanguageType, error)

	// 指定求人IDリストの必要言語スキルの種類を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredLanguageType, error)

	// 指定求人IDリストと言語タイプの必要言語スキルの種類を取得
	GetByJobInformationIDListAndLanguageTypeList(jobInformationIDList []uint, languageTypeList []entity.Language) ([]*entity.JobInformationRequiredLanguageType, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredLanguageType, error)
}

// 必要開発経験
type JobInformationRequiredExperienceDevelopmentRepository interface {
	/** 作成 */
	// 必要開発経験を作成する
	Create(requiredExperienceDevelopment *entity.JobInformationRequiredExperienceDevelopment) error

	/** 複数取得 */
	// 指定求人IDの必要開発経験を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredExperienceDevelopment, error)

	// 指定請求先IDの必要開発経験を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredExperienceDevelopment, error)

	// 指定企業IDの必要開発経験を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredExperienceDevelopment, error)

	// 指定エージェントIDの必要開発経験を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredExperienceDevelopment, error)

	// 指定求人IDリストの必要開発経験を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredExperienceDevelopment, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredExperienceDevelopment, error)
}

// 必要開発経験の種類
type JobInformationRequiredExperienceDevelopmentTypeRepository interface {
	/** 作成 */
	// 必要開発経験の種類を作成する
	Create(requiredExperienceDevelopmentType *entity.JobInformationRequiredExperienceDevelopmentType) error

	/** 複数取得 */
	// 指定求人IDの必要開発経験の種類を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredExperienceDevelopmentType, error)

	// 指定請求先IDの必要開発経験の種類を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredExperienceDevelopmentType, error)

	// 指定企業IDの必要開発経験の種類を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredExperienceDevelopmentType, error)

	// 指定エージェントIDの必要開発経験の種類を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredExperienceDevelopmentType, error)

	// 指定求人IDリストの必要開発経験の種類を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredExperienceDevelopmentType, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredExperienceDevelopmentType, error)
}

// 必要業界・職種経験
type JobInformationRequiredExperienceJobRepository interface {
	/** 作成 */
	// 必要業職種経験を作成する
	Create(requiredExperienceJob *entity.JobInformationRequiredExperienceJob) error

	/** 複数取得 */
	// 指定求人IDの必要業職種経験を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredExperienceJob, error)

	// 指定請求先IDの必要業職種経験を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredExperienceJob, error)

	// 指定企業IDの必要業職種経験を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredExperienceJob, error)

	// 指定エージェントIDの必要業職種経験を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredExperienceJob, error)

	// 指定求人IDリストの必要業職種経験を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredExperienceJob, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredExperienceJob, error)
}

// 必要業界経験
type JobInformationRequiredExperienceIndustryRepository interface {
	/** 作成 */
	// 必要業界経験を作成する
	Create(requiredExperienceIndustry *entity.JobInformationRequiredExperienceIndustry) error

	/** 複数取得 */
	// 指定求人IDの必要業界経験を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredExperienceIndustry, error)

	// 指定請求先IDの必要業界経験を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredExperienceIndustry, error)

	// 指定企業IDの必要業界経験を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredExperienceIndustry, error)

	// 指定エージェントIDの必要業界経験を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredExperienceIndustry, error)

	// 指定求人IDリストの必要業界経験を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredExperienceIndustry, error)

	// 指定求人IDリストと業界の必要業界経験を取得
	GetByJobInformationIDListAndIndustryList(jobInformationIDList []uint, industryList []entity.ExperienceIndustry) ([]*entity.JobInformationRequiredExperienceIndustry, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredExperienceIndustry, error)
}

// 必要職種経験
type JobInformationRequiredExperienceOccupationRepository interface {
	/** 作成 */
	// 必要職種経験を作成する
	Create(requiredExperienceOccupation *entity.JobInformationRequiredExperienceOccupation) error

	/** 複数取得 */
	// 指定求人IDの必要職種経験を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredExperienceOccupation, error)

	// 指定請求先IDの必要職種経験を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredExperienceOccupation, error)

	// 指定企業IDの必要職種経験を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredExperienceOccupation, error)

	// 指定エージェントIDの必要職種経験を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredExperienceOccupation, error)

	// 指定求人IDリストの必要職種経験を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredExperienceOccupation, error)

	// 指定求人IDリストと職種の必要職種経験を取得
	GetByJobInformationIDListAndOccupationList(jobInformationIDList []uint, occupationList []entity.ExperienceOccupation) ([]*entity.JobInformationRequiredExperienceOccupation, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredExperienceOccupation, error)
}

// 必要社会人経験
type JobInformationRequiredSocialExperienceRepository interface {
	/** 作成 */
	// 必要社会人経験を作成する
	Create(requiredSocialExperience *entity.JobInformationRequiredSocialExperience) error

	/** 削除 */
	// 必要社会人経験を削除する
	DeleteByJobInformationID(jobInformationID uint) error

	/** 複数取得 */
	// 指定求人IDの必要社会人経験を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationRequiredSocialExperience, error)

	// 指定請求先IDの必要社会人経験を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationRequiredSocialExperience, error)

	// 指定企業IDの必要社会人経験を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationRequiredSocialExperience, error)

	// 指定エージェントIDの必要社会人経験を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationRequiredSocialExperience, error)

	// 指定求人IDリストの必要社会人経験を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationRequiredSocialExperience, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationRequiredSocialExperience, error)
}

// 選考フローのパターン
type JobInformationSelectionFlowPatternRepository interface {
	/** 作成 */
	// 選考フローのパターンを作成する
	Create(selectionFlowPattern *entity.JobInformationSelectionFlowPattern) error

	/** 更新 */
	// 選考フローのパターンを更新する
	Update(id uint, selectionFlowPattern *entity.JobInformationSelectionFlowPattern) error

	/** 削除 */
	// 選考フローのパターンを削除する
	Delete(id uint) error

	/** 単数取得 */
	// 指定IDの選考フローのパターンを取得する
	FindByID(id uint) (*entity.JobInformationSelectionFlowPattern, error)

	/** 複数取得 */
	// 指定求人IDの選考フローのパターンを取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationSelectionFlowPattern, error)

	// public_statusが0のレコードのみ取得
	GetOpenByJobInformationID(jobInformationID uint) ([]*entity.JobInformationSelectionFlowPattern, error)

	// 指定請求先IDの選考フローのパターンを取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationSelectionFlowPattern, error)

	// 指定企業IDの選考フローのパターンを取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationSelectionFlowPattern, error)

	// 指定エージェントIDの選考フローのパターンを取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationSelectionFlowPattern, error)

	// 指定IDリストの選考フローのパターンを取得
	GetByIDList(idList []uint) ([]*entity.JobInformationSelectionFlowPattern, error)

	// 指定求人IDリストの選考フローのパターンを取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationSelectionFlowPattern, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationSelectionFlowPattern, error)
}

// 選考情報
type JobInformationSelectionInformationRepository interface {
	/** 作成 */
	// 選考情報を作成する
	Create(selectionInformation *entity.JobInformationSelectionInformation) error

	/** 更新 */
	// 選考情報を更新する
	Update(id uint, selectionInformation *entity.JobInformationSelectionInformation) error

	/** 単数取得 */
	//
	FindBySelectionFlowIDAndSelectionType(selectionFlowID, selectionType uint) (*entity.JobInformationSelectionInformation, error)

	// タスクグループの求職者id、求人idと選考フェーズに合致するレコードを取得（アンケートに使用）
	FindByJobSeekerIDAndJobInformationIDAndSelectionType(jobSeekerID, jobInformationID uint, selectionType null.Int) (*entity.JobInformationSelectionInformation, error)

	/** 複数取得 */
	// 指定選考フローIDの選考情報を取得する
	GetBySelectionFlowID(selectionFlowID uint) ([]*entity.JobInformationSelectionInformation, error)

	// 指定求人IDの選考情報を取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationSelectionInformation, error)

	// 指定請求先IDの選考情報を取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationSelectionInformation, error)

	// 指定企業IDの選考情報を取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationSelectionInformation, error)

	// 指定エージェントIDの選考情報を取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationSelectionInformation, error)

	// 指定求人IDリストの選考情報を取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationSelectionInformation, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationSelectionInformation, error)
}

// 非公開エージェント
type JobInformationHideToAgentRepository interface {
	/** 作成 */
	// 非公開エージェントを作成する
	Create(hideToAgent *entity.JobInformationHideToAgent) error

	/** 削除 */
	// 非公開エージェントを削除する
	DeleteByJobInformationID(jobInformationID uint) error

	/** 複数取得 */
	// 指定求人IDの非公開エージェントを取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationHideToAgent, error)

	// 指定請求先IDの非公開エージェントを取得する
	GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationHideToAgent, error)

	// 指定企業IDの非公開エージェントを取得する
	GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationHideToAgent, error)

	// 指定エージェントIDの非公開エージェントを取得する
	GetByAgentID(agentID uint) ([]*entity.JobInformationHideToAgent, error)

	// エージェントIDから対象エージェント非公開設定の求人を取得
	GetHideByAgentID(agentID uint) ([]*entity.JobInformationHideToAgent, error)

	// 指定求人IDリストの非公開エージェントを取得
	GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationHideToAgent, error)

	// 指定エージェントIDリストの非公開エージェントを取得
	GetByAgentIDList(agentIDList []uint) ([]*entity.JobInformationHideToAgent, error)

	// 全てのレコードを取得
	All() ([]*entity.JobInformationHideToAgent, error)
}

// 求人の他媒体でのID
type JobInformationExternalIDRepository interface {
	/** 作成 */
	// 他媒体でのIDを作成する
	Create(externalID *entity.JobInformationExternalID) error

	/** 削除 */
	// 他媒体でのIDを削除する
	DeleteByJobInformationID(jobInformationID uint) error

	/** 複数取得 */
	// 指定求人IDの他媒体でのIDを取得する
	GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationExternalID, error)

	// 指定エージェントIDと媒体タイプの他媒体でのIDを取得する
	GetByAgentIDAndExternalType(agentID uint, externalType entity.JobInformatinoExternalType) ([]*entity.JobInformationExternalID, error)
}

/****************************************************************************************/
/****************************************************************************************/
/// 求職者プロフィール関連
//
type JobSeekerRepository interface {
	/** 作成 */
	Create(jobSeeker *entity.JobSeeker) error

	/** 更新 */
	Update(id uint, jobSeeker *entity.JobSeeker) error

	// テスト環境用に名前と電話番号を更新
	UpdateForDev(id uint) error

	// activity_memoカラムのみを更新する処理
	UpdateActivityMemo(id uint, activityMemo string) error

	// フェーズのみを更新する処理
	UpdatePhase(id uint, phase null.Int) error

	// 面談日を更新する処理
	UpdateInterviewDate(id uint, interviewDate time.Time) error

	// CAを更新する処理
	UpdateStaffID(id uint, staffID null.Int) error

	// created_atを更新する処理 *csvインポートでエントリー日が入力されている場合に使用
	UpdateCreatedAt(id uint, createdAt time.Time) error

	// 応募承諾のポイントを更新する処理 *求職者相談時に使用
	UpdateAcceptancePoints(id uint, acceptancePoints string) error

	// 個人情報同意の更新
	UpdateAgreement(id uint, agreement bool) error

	// LineIDを更新
	UpdateLineIDByUUID(uuid uuid.UUID, lineID string) error

	// ゲストログイン用のパスワードを更新
	UpdatePassword(uuid uuid.UUID, password string) error

	UpdateResetPasswordToken(uuid uuid.UUID, resetPasswordToken string) error

	// uuidを使って電話番を更新（LPでのみ使用）
	UpdatePhoneNumberByUUID(uuid uuid.UUID, phoneNumber string) error

	// uuidを使って電話番号と希望年収を更新（LPでのみ使用）
	UpdateDesiredIncomeByUUID(uuid uuid.UUID, desiredIncome null.Int) error

	// おすすめ求人の閲覧権限の更新
	UpdateCanViewMatchingJob(id uint, canView bool) error

	/** 削除 */
	Delete(id uint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.JobSeeker, error)

	FindByUUID(uuid uuid.UUID) (*entity.JobSeeker, error)

	FindByTaskGroupUUID(taskGroupUUID uuid.UUID) (*entity.JobSeeker, error)

	FindByAgentIDAndLineID(agentID uint, lineID string) (*entity.JobSeeker, error)

	FindByPhoneNumberForLP(phoneNumber string, agentID uint) (*entity.JobSeeker, error)

	FindByEmailForLP(email string, agentID uint) (*entity.JobSeeker, error)

	FindByResetPasswordTokenForLP(resetPasswordToken string, agentID uint) (*entity.JobSeeker, error)

	/** カウント */
	CountByEmail(email string) (float64, error)

	// 送客の重複登録判定
	FindByNameAndPhoneNumberBySystemAgent(firstName, lastName, firstFurigana, lastFurigana, phoneNumber string) (*entity.JobSeeker, error)

	/** 複数取得 */
	GetByAgentID(agentID uint) ([]*entity.JobSeeker, error)

	// 2日以内に登録された求職者を取得　*エントリー取得時の重複判定に使用
	GetByAgentIDWithinTwoDays(agentID uint) ([]*entity.JobSeeker, error)

	GetByIDList(idList []uint) ([]*entity.JobSeeker, error)

	// タスク用
	GetByIDListForTask(idList []uint) ([]*entity.JobSeeker, error)

	GetByAgentStaffIDAndNotRelease(agentStaffID uint) ([]*entity.JobSeeker, error)

	// 求人検索→求職者検索
	GetActiveAllAndFreeWord(agentID uint, freeWord string) ([]*entity.JobSeeker, error)

	// シェア求人検索→自社求職者検索
	GetActiveOwnAndFreeWord(agentID uint, freeWord string) ([]*entity.JobSeeker, error)

	// 絞り込み検索
	GetByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.JobSeeker, error)

	// アライアンスを含む
	GetByAgentIDAndFreeWordAndAlliance(agentID uint, freeWord string, allaicenAgentList []*entity.AgentAlliance) ([]*entity.JobSeeker, error)

	// アライアンス
	GetByOtherAgentIDAndFreeWord(agentID uint, freeWord string, agentAllianceList []*entity.AgentAlliance) ([]*entity.JobSeeker, error)

	// ダッシュボード
	GetReleaseByStaffID(staffID uint) ([]*entity.JobSeeker, error)

	GetReleaseByAgentID(agentID uint) ([]*entity.JobSeeker, error)

	// 引数の値で重複する求職者を取得する
	GetDuplicateByNameAndFuriganaAndEmailAndPhoneNumber(agentID uint, lastName, firstName, lastFurigana, firstFurigana, email, phoneNumber string) ([]*entity.JobSeeker, error)

	All() ([]*entity.JobSeeker, error)
}

type JobSeekerDocumentRepository interface {
	/** 作成 */
	Create(document *entity.JobSeekerDocument) error

	/** 更新 */
	UpdateByJobSeekerID(jobSeekerID uint, document *entity.JobSeekerDocument) error
	UpdateResumeOriginURLByJobSeekerID(jobSeekerID uint, resumeOriginURL string) error
	UpdateCVOriginURLByJobSeekerID(jobSeekerID uint, cvOriginURL string) error
	UpdateResumePDFURLByJobSeekerID(jobSeekerID uint, resumePDFURL string) error
	UpdateCVPDFURLByJobSeekerID(jobSeekerID uint, cvPDFURL string) error
	UpdateRecommendationOriginURLByJobSeekerID(jobSeekerID uint, recommendationOriginURL string) error
	UpdateRecommendationPDFURLByJobSeekerID(jobSeekerID uint, recommendationPDFURL string) error
	UpdateIDPhotoURLByJobSeekerID(jobSeekerID uint, idPhotoURL string) error
	UpdateOtherDocument1URLByJobSeekerID(jobSeekerID uint, otherDocument1URL string) error
	UpdateOtherDocument2URLByJobSeekerID(jobSeekerID uint, otherDocument2URL string) error
	UpdateOtherDocument3URLByJobSeekerID(jobSeekerID uint, otherDocument3URL string) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 単数取得 */
	FindByJobSeekerID(jobSeekerID uint) (*entity.JobSeekerDocument, error)
	FindByJobSeekerUUID(jobSeekerUUID uuid.UUID) (*entity.JobSeekerDocument, error)

	/** 複数取得 */
	GetByAgentID(agentID uint) ([]*entity.JobSeekerDocument, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerDocument, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerDocument, error)

	All() ([]*entity.JobSeekerDocument, error)
}

type JobSeekerStudentHistoryRepository interface {
	/** 作成 */
	Create(studentHistory *entity.JobSeekerStudentHistory) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerStudentHistory, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerStudentHistory, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerStudentHistory, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerStudentHistory, error)

	All() ([]*entity.JobSeekerStudentHistory, error)
}

type JobSeekerWorkHistoryRepository interface {
	/** 作成 */
	Create(workHistory *entity.JobSeekerWorkHistory) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerWorkHistory, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerWorkHistory, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerWorkHistory, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerWorkHistory, error)

	All() ([]*entity.JobSeekerWorkHistory, error)
}

type JobSeekerExperienceIndustryRepository interface {
	/** 作成 */
	Create(experienceIndustry *entity.JobSeekerExperienceIndustry) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerExperienceIndustry, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerExperienceIndustry, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerExperienceIndustry, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerExperienceIndustry, error)

	All() ([]*entity.JobSeekerExperienceIndustry, error)
}

type JobSeekerDepartmentHistoryRepository interface {
	/** 作成 */
	Create(departmentHistory *entity.JobSeekerDepartmentHistory) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDepartmentHistory, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerDepartmentHistory, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerDepartmentHistory, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerDepartmentHistory, error)

	All() ([]*entity.JobSeekerDepartmentHistory, error)
}

type JobSeekerExperienceOccupationRepository interface {
	/** 作成 */
	Create(experienceOccupation *entity.JobSeekerExperienceOccupation) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerExperienceOccupation, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerExperienceOccupation, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerExperienceOccupation, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerExperienceOccupation, error)

	All() ([]*entity.JobSeekerExperienceOccupation, error)
}

type JobSeekerLicenseRepository interface {
	/** 作成 */
	Create(license *entity.JobSeekerLicense) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerLicense, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerLicense, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerLicense, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerLicense, error)

	All() ([]*entity.JobSeekerLicense, error)
}

type JobSeekerSelfPromotionRepository interface {
	/** 作成 */
	Create(industry *entity.JobSeekerSelfPromotion) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerSelfPromotion, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerSelfPromotion, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerSelfPromotion, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerSelfPromotion, error)

	All() ([]*entity.JobSeekerSelfPromotion, error)
}

type JobSeekerDesiredIndustryRepository interface {
	/** 作成 */
	Create(desiredIndustry *entity.JobSeekerDesiredIndustry) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDesiredIndustry, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerDesiredIndustry, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerDesiredIndustry, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerDesiredIndustry, error)

	All() ([]*entity.JobSeekerDesiredIndustry, error)
}

type JobSeekerDesiredOccupationRepository interface {
	/** 作成 */
	Create(desiredOccupation *entity.JobSeekerDesiredOccupation) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDesiredOccupation, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerDesiredOccupation, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerDesiredOccupation, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerDesiredOccupation, error)

	All() ([]*entity.JobSeekerDesiredOccupation, error)
}

type JobSeekerDesiredWorkLocationRepository interface {
	/** 作成 */
	Create(desiredWorkLocation *entity.JobSeekerDesiredWorkLocation) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDesiredWorkLocation, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerDesiredWorkLocation, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerDesiredWorkLocation, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerDesiredWorkLocation, error)

	All() ([]*entity.JobSeekerDesiredWorkLocation, error)
}

type JobSeekerDesiredHolidayTypeRepository interface {
	/** 作成 */
	Create(desiredHolidayType *entity.JobSeekerDesiredHolidayType) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDesiredHolidayType, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerDesiredHolidayType, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerDesiredHolidayType, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerDesiredHolidayType, error)

	All() ([]*entity.JobSeekerDesiredHolidayType, error)
}

// 　求職者の開発スキル
type JobSeekerDevelopmentSkillRepository interface {
	/** 作成 */
	Create(desiredWorkLocation *entity.JobSeekerDevelopmentSkill) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDevelopmentSkill, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerDevelopmentSkill, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerDevelopmentSkill, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerDevelopmentSkill, error)

	All() ([]*entity.JobSeekerDevelopmentSkill, error)
}
type JobSeekerLanguageSkillRepository interface {
	/** 作成 */
	Create(desiredWorkLocation *entity.JobSeekerLanguageSkill) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerLanguageSkill, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerLanguageSkill, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerLanguageSkill, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerLanguageSkill, error)

	All() ([]*entity.JobSeekerLanguageSkill, error)
}

// pcTool
type JobSeekerPCToolRepository interface {
	/** 作成 */
	Create(desiredWorkLocation *entity.JobSeekerPCTool) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerPCTool, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerPCTool, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerPCTool, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerPCTool, error)

	All() ([]*entity.JobSeekerPCTool, error)
}

// 非公開エージェント
type JobSeekerHideToAgentRepository interface {
	/** 作成 */
	Create(hideToAgent *entity.JobSeekerHideToAgent) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerHideToAgent, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerHideToAgent, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerHideToAgent, error)
	// エージェントIDから対象エージェント非公開設定の求職者を取得
	GetHideByAgentID(agentID uint) ([]*entity.JobSeekerHideToAgent, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerHideToAgent, error)
	GetByAgentIDList(agentIDList []uint) ([]*entity.JobSeekerHideToAgent, error)

	All() ([]*entity.JobSeekerHideToAgent, error)
}

// desiredComapnyScale
type JobSeekerDesiredCompanyScaleRepository interface {
	/** 作成 */
	Create(desiredCompanyScale *entity.JobSeekerDesiredCompanyScale) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerDesiredCompanyScale, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerDesiredCompanyScale, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerDesiredCompanyScale, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerDesiredCompanyScale, error)

	All() ([]*entity.JobSeekerDesiredCompanyScale, error)
}

// 求職者の経験職種（LPからの登録時に使用）
type JobSeekerExperienceJobRepository interface {
	/** 作成 */
	Create(desiredCompanyScale *entity.JobSeekerExperienceJob) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerExperienceJob, error)
	GetByAgentID(agentID uint) ([]*entity.JobSeekerExperienceJob, error)
	GetByStaffID(staffID uint) ([]*entity.JobSeekerExperienceJob, error)
	GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerExperienceJob, error)

	All() ([]*entity.JobSeekerExperienceJob, error)
}

// 求職者の経験職種（LPからの登録時に使用）
type JobSeekerLPLoginTokenRepository interface {
	/** 作成 */
	Create(lpLoginToken *entity.JobSeekerLPLoginToken) error

	/** 更新 */
	UpdateByJobSeekerID(jobSeekerID uint, lpLoginToken *entity.JobSeekerLPLoginToken) error

	/** 削除 */
	DeleteByJobSeekerID(jobSeekerID uint) error

	/** 単数取得 */
	FindByJobSeekerID(jobSeekerID uint) (*entity.JobSeekerLPLoginToken, error)
}

type JobSeekerInterestedJobListingRepository interface {
	/** 作成 */
	Create(interestedJobListing *entity.JobSeekerInterestedJobListing) error

	/** 削除 */
	Delete(id uint) error

	/** 複数取得 */
	GetByJobSeekerUUIDAndInterestedType(jobSeekerUUID uuid.UUID, interestedType entity.InterestedType) ([]*entity.JobSeekerInterestedJobListing, error)
}

/****************************************************************************************/
/****************************************************************************************/
/// タスク関連
//
type TaskGroupRepository interface {
	/** 作成 */
	Create(taskGroup *entity.TaskGroup) error

	/** 更新 */
	UpdateRALastRequestAt(id uint) error
	UpdateRALastWatchedAt(id uint) error
	UpdateCALastRequestAt(id uint) error
	UpdateCALastWatchedAt(id uint) error
	UpdateLastRequestAt(id uint) error // RAとCAの依頼時間をどちらも更新
	UpdateSelectionFlowPatternID(id uint, selectionFlowID null.Int) error
	UpdateJoiningDate(id uint, joiningDate string) error
	UpdateIsDoubleSided(id uint, isDoubleSided bool) error
	UpdateListIsDoubleSided(idList []uint, isDoubleSided bool) error
	UpdateExternalJob(id uint, externalJob entity.ExternalJob) error

	/** 削除 */
	Delete(id uint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.TaskGroup, error)

	/** 複数取得 */
	GetByJobInformationIDList(idList []uint) ([]*entity.TaskGroup, error)
	GetByJobSeekerIDList(idList []uint) ([]*entity.TaskGroup, error)
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.TaskGroup, error)
	// is_double_sidedがfalseでra_staff_idとca_staff_idが一致しているtaskGroup
	GetNotDoubleSidedSameRAAndCAByStaffID(agentStaffID uint) ([]*entity.TaskGroup, error)

	// agent_idが一致するtaskGroup
	GetByAgentID(agentID uint) ([]*entity.TaskGroup, error)

	// GetByAgentIDAndSearchParam(agentID uint, searchParam entity.SearchTask) ([]*entity.TaskGroup, error)
}

type TaskGroupDocumentRepository interface {
	/** 作成 */
	Create(document *entity.TaskGroupDocument) error

	/** 更新 */
	Update(id uint, document *entity.TaskGroupDocument) error
	UpdateDocument1URL(id uint, document1URL string) error // 添付ファイル1の更新
	UpdateDocument2URL(id uint, document2URL string) error // 添付ファイル2の更新
	UpdateDocument3URL(id uint, document3URL string) error // 添付ファイル3の更新
	UpdateDocument4URL(id uint, document4URL string) error // 添付ファイル4の更新
	UpdateDocument5URL(id uint, document5URL string) error // 添付ファイル5の更新

	/** 単数取得 */
	FindByID(id uint) (*entity.TaskGroupDocument, error)
	FindByGroupID(groupID uint) (*entity.TaskGroupDocument, error)
}

type TaskRepository interface {
	/** 作成 */
	Create(task *entity.Task) error

	/** 更新 */
	Update(id uint, task *entity.Task) error

	/** 削除 */
	Delete(id uint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.Task, error)

	// 求職者UUIDと求人UUIDから最新のタスクを取得 *関連テーブル含まない
	FindByJobSeekerUUIDAndJobInformationUUID(jobseekerUUID, jobInformationUUID uuid.UUID) (*entity.Task, error)

	// 求職者UUIDと求人UUIDから最新のタスクを取得 *関連テーブル含む
	FindLatestWithRelatedByJobSeekerUUIDAndJobInformationUUID(jobInformationUUID, jobseekerUUID uuid.UUID) (*entity.Task, error)

	// 求職者IDと求人IDから最新のタスクを取得
	FindLatestByJobSeekerIDAndJobInformationID(jobSeekerID, jobInformationID uint) (*entity.Task, error)

	// タスクグループIDから最新のタスクを取得
	FindLatestByGroupID(taskGroupID uint) (*entity.Task, error)

	// 直近の「書類選考/結果回収中」or「○次選考/選考結果の回収依頼」を取得
	FindLatestByGroupIDAndCollectResultPhase(taskGroupID uint) (*entity.Task, error)

	//「phase_sub_category」が辞退処理する前（100未満）のタスクを取得 *選考継続時に使用
	FindLatestForContinue(taskGroupID uint, phase null.Int) (*entity.Task, error)

	/** 複数取得 */
	GetByTaskGroupID(taskGroupID uint) ([]*entity.Task, error)

	// タスクとタスク関連の子テーブルの情報をすべて取得
	GetWithRelatedByTaskGroupID(taskGroupID uint) ([]*entity.Task, error)

	// エージェントIDから取得
	GetByAgentID(agentID uint) ([]*entity.Task, error)

	// エージェントIDと最新を取得
	// GetLatestByAgentID(agentID uint) ([]*entity.Task, error)
	GetLatestByGroupIDList(groupIDList []uint) ([]*entity.Task, error)

	// 書類選考以降のタスク情報取得
	GetLatestAfterSelectPhaseByJobSeekerID(jobSeekerID uint, phase entity.TaskCategory) ([]*entity.Task, error)

	// エージェントIDとフリーワードから最新を取得
	GetLatestByAgentIDAndFreeWord(agentID uint, jobSeekerFreeWord, enterpriseFreeWord string) ([]*entity.Task, error)

	// 担当者IDから最新を取得
	GetLatestByStaffID(agentStaffID uint) ([]*entity.Task, error)

	// 「同一求職者」で「同一フェーズ」のタスクのリストを取得
	GetLatestSameByJobSeekerID(jobSeekerID, taskID uint, phase, phaseSub null.Int) ([]*entity.Task, error)

	// 各エージェントIDがCA/RAに該当するタスクを取得　*アライアンス解除時のタスク確認用
	GetByEachAgentID(agent1ID, agent2ID uint) ([]*entity.Task, error)

	// 引数に合致する最新タスクを取得する
	GetLatestByJobSeekerUUIDAndJobInformationIDList(jobSeekerUUID uuid.UUID, jobInformationIDList []uint) ([]*entity.Task, error)

	// 全ての最新フェーズタスクを取得 *未読通知で使用
	GetLatest() ([]*entity.Task, error)

	// カウント用
	GetActiveByBillingAddressID(billingAddressID uint) ([]*entity.Task, error)
	GetActiveByJobInformationID(jobInformationID uint) ([]*entity.Task, error)
	GetActiveBySelectionFlowPatternID(selectionID uint) ([]*entity.Task, error)

	// 売上管理関係のカウント（全体）
	GetSeekerOfferPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)
	GetSeekerFinalSelectionPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)
	GetSeekerSelectionPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)
	GetSeekerRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)
	GetSeekerJobIntroductionPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)
	GetInterviewOfferPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)
	GetInterviewFinalSelectionPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)
	GetInterviewSelectionPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)
	GetInterviewRecommendationCompletionPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)
	GetInterviewJobIntroductionPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)

	// 売上管理関係のカウント（個人）
	GetSeekerOfferPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)
	GetSeekerFinalSelectionPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)
	GetSeekerSelectionPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)
	GetSeekerRecommendationCompletionPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)
	GetSeekerJobIntroductionPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)
	GetInterviewOfferPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)
	GetInterviewFinalSelectionPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)
	GetInterviewSelectionPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)
	GetInterviewRecommendationCompletionPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)
	GetInterviewJobIntroductionPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)

	// ダッシュボードの売上管理関係のカウント（全体）
	GetInterviewOfferPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error)
	GetInterviewFinalSelectionPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error)
	GetInterviewSelectionPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error)
	GetInterviewRecommendationCompletionPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error)
	GetInterviewJobIntroductionPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error)

	// 売上管理関係のカウント（個人）
	GetInterviewOfferPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error)
	GetInterviewFinalSelectionPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error)
	GetInterviewSelectionPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error)
	GetInterviewRecommendationCompletionPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error)
	GetInterviewJobIntroductionPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error)
}

type TaskIsRecommendDocumentRepository interface {
	/** 作成 */
	Create(task *entity.TaskIsRecommendDocument) error

	/** 更新 */
	Update(id uint, task *entity.TaskIsRecommendDocument) error

	/** 単数取得 */
	FindByID(id uint) (*entity.TaskIsRecommendDocument, error)
	FindByTaskID(taskID uint) (*entity.TaskIsRecommendDocument, error)

	/** 複数取得 */
	// アライアンスエージェントIDのリストからエージェント一覧を取得
	GetByTaskIDList(taskIDList []uint) ([]*entity.TaskIsRecommendDocument, error)
}

type EvaluationPointRepository interface {
	/** 作成 */
	Create(evaluationPoint *entity.EvaluationPoint) error

	/** 更新 */
	Update(id uint, evaluationPoint *entity.EvaluationPoint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.EvaluationPoint, error)
	FindByTaskID(taskID uint) (*entity.EvaluationPoint, error)

	/** 複数取得 */
	GetByJobSeekerIDAndJobInformationID(jobSeekerID, jobInformationID uint) ([]*entity.EvaluationPoint, error)

	// 選考フローIDから評価ポイントを取得
	GetBySelectionFlowID(selectionFlowID uint) ([]*entity.EvaluationPoint, error)
}

/****************************************************************************************/
/****************************************************************************************/
/// メッセージテンプレート
//
type MessageTemplateRepository interface {
	/** 作成 */
	// メッセージテンプレートを作成
	Create(messageTemplate *entity.MessageTemplate) error

	/** 更新 */
	// メッセージテンプレートを更新
	Update(id uint, messageTemplate *entity.MessageTemplate) error

	/** 削除 */
	// メッセージテンプレートを削除
	Delete(id uint) error

	/** 単数取得 */
	// IDからメッセージテンプレートを取得
	FindByID(id uint) (*entity.MessageTemplate, error)

	/** 複数取得 */
	// 担当者IDからメッセージテンプレートを取得
	GetByAgentStaffID(agentStaffID uint) ([]*entity.MessageTemplate, error)

	// 担当者IDと送信シーンからメッセージテンプレートを取得
	GetByAgentStaffIDAndSendScene(agentStaffID, sendScene uint) ([]*entity.MessageTemplate, error)

	// エージェントIDからメッセージテンプレートを取得
	GetByAgentID(agentID uint) ([]*entity.MessageTemplate, error)
}

/****************************************************************************************/
/// 選考後アンケート
//
type SelectionQuestionnaireRepository interface {
	/** 作成 */
	// 選考後アンケートを作成
	Create(selectionQuestionnaire *entity.SelectionQuestionnaire) error

	// UUIDから選考後アンケートを作成
	CreateByUUID(selectionQuestionnaire *entity.SelectionQuestionnaire, questionnaireUUID uuid.UUID) error

	/** 更新 */
	// 選考後アンケートを更新
	Update(id uint, selectionQuestionnaire *entity.SelectionQuestionnaire) error

	/** 単数取得 */
	// IDから選考後アンケートを取得
	FindByID(id uint) (*entity.SelectionQuestionnaire, error)

	// UUIDから選考後アンケートを取得
	FindByUUID(uuid uuid.UUID) (*entity.SelectionQuestionnaire, error)

	// 求職者IDと選考フローIDから選考後アンケートを取得
	FindByJobSeekerIDAndJobInformationIDAndSelectionPhase(jobSeekerID, jobInformationID uint, selectionPhase uint) (*entity.SelectionQuestionnaire, error)

	// 求職者UUIDと選考フローUUIDから選考後アンケートを取得
	FindByJobSeekerUUIDAndJobInformationUUIDAndSelectionPhase(jobSeekerUUID, jobInformationUUID uuid.UUID, selectionPhase uint) (*entity.SelectionQuestionnaire, error)

	/** 複数取得 */
	// 求職者IDから選考後アンケートを取得
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.SelectionQuestionnaire, error)

	// 求職者IDと選考フローIDから選考後アンケートを取得
	GetBySelectionFlowID(selectionFlowID uint) ([]*entity.SelectionQuestionnaire, error)

	// 指定uuidに合致するis_answerdがfalseのレコードを取得
	GetUnanswerdByJobSeekerUUID(jobSeekerUUID uuid.UUID) ([]*entity.SelectionQuestionnaire, error)
}

type SelectionQuestionnaireMyRankingRepository interface {
	/** 作成 */
	// 選考後アンケートの希望順位を作成
	Create(selectionQuestionnaire *entity.SelectionQuestionnaireMyRanking) error

	/** 削除 */
	// 選考後アンケートの希望順位を削除
	DeleteByQuestionnaireID(questionnaireID uint) error

	/** 複数取得 */
	// 選考後アンケートIDから希望順位を取得
	GetByQuestionnaireID(questionnaireID uint) ([]*entity.SelectionQuestionnaireMyRanking, error)

	// 選考後アンケートIDリストから希望順位を取得
	GetByQuestionnaireIDList(idList []uint) ([]*entity.SelectionQuestionnaireMyRanking, error)
}

/****************************************************************************************/

/****************************************************************************************/
/****************************************************************************************/
/// エージェントと求職者のチャット
//
type ChatGroupWithJobSeekerRepository interface {
	/** 作成 */
	// 求職者とのチャットグループを作成
	Create(chatGroupWithJobSeeker *entity.ChatGroupWithJobSeeker) error

	/** 更新 */
	// エージェントの最終閲覧時間を更新
	UpdateAgentLastWatchedAt(id uint) error

	// エージェントの最終送信時間を更新
	UpdateAgentLastSendAt(id uint) error

	// 求職者の最終閲覧時間と最終送信時間を更新
	UpdateJobSeekerLastWatchedAtAndSendAt(id uint, sendAt time.Time) error

	// 求職者のLINEの有効/無効を更新
	UpdateJobSeekerLineActive(jobSeekerID uint, isActive bool) error

	/** 単数取得 */
	// IDからチャットグループを取得
	FindByID(id uint) (*entity.ChatGroupWithJobSeeker, error)

	// エージェントIDと求職者IDからチャットグループを取得
	FindByAgentIDAndJobSeekerID(agentID, jobSeekerID uint) (*entity.ChatGroupWithJobSeeker, error)

	/** 複数取得 */
	// エージェントIDからチャットグループを取得
	GetByAgentID(agentID uint) ([]*entity.ChatGroupWithJobSeeker, error)

	// エージェントIDとフリーワードからチャットグループを取得
	GetByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.ChatGroupWithJobSeeker, error)

	// エージェントIDからアクティブなチャットグループを取得 ※管理側一斉送信用
	GetActiveJobSeekerByAgentID(agentID uint) ([]*entity.ChatGroupWithJobSeeker, error)

	// 全てのチャットグループを取得
	All() ([]*entity.ChatGroupWithJobSeeker, error)

	/** その他 */
	// エージェントIDから通知数を取得
	CountNotificationByAgentID(agentID uint) (float64, error)
}

type ChatMessageWithJobSeekerRepository interface {
	/** 作成 */
	// 求職者とのチャットメッセージを作成
	Create(chatMessageJobSeeker *entity.ChatMessageWithJobSeeker) error

	/** 単数取得 */
	FindByLINEMessageID(messageID string) (*entity.ChatMessageWithJobSeeker, error)

	/** 複数取得 */
	// チャットグループIDからチャットメッセージを取得
	GetByGroupID(groupID uint) ([]*entity.ChatMessageWithJobSeeker, error)
}

// 求職者とのメールを管理する
type EmailWithJobSeekerRepository interface {
	/** 作成 */
	// 求職者とのメールを作成
	Create(emailWithJobSeeker *entity.EmailWithJobSeeker) error

	/** 複数取得 */
	// 求職者IDからメールを取得
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.EmailWithJobSeeker, error)
}

/****************************************************************************************/
/****************************************************************************************/
/// 求職者の面談前アンケート
//
type InitialQuestionnaireRepository interface {
	/** 作成 */
	// 求職者の面談前アンケートを作成
	Create(initialQuestionnaire *entity.InitialQuestionnaire) error

	/** 単数取得 */
	// 求職者IDから面談前アンケートを取得
	FindByJobSeekerID(jobSeekerID uint) (*entity.InitialQuestionnaire, error)
}

// 希望業界
type InitialQuestionnaireDesiredIndustryRepository interface {
	/** 作成 */
	// 求職者の希望業界を作成
	Create(desiredIndustry *entity.InitialQuestionnaireDesiredIndustry) error
}

// 希望職種
type InitialQuestionnaireDesiredOccupationRepository interface {
	/** 作成 */
	// 求職者の希望職種を作成
	Create(desiredOccupation *entity.InitialQuestionnaireDesiredOccupation) error
}

// 希望勤務地
type InitialQuestionnaireDesiredWorkLocationRepository interface {
	/** 作成 */
	// 求職者の希望勤務地を作成
	Create(desiredWorkLocation *entity.InitialQuestionnaireDesiredWorkLocation) error
}

/****************************************************************************************/
/****************************************************************************************/
/// エージェントとエージェントのチャット
//
type ChatGroupWithAgentRepository interface {
	/** 作成 */
	Create(chatGroupWithAgent *entity.ChatGroupWithAgent) error

	/** 更新 */
	UpdateAgentLastWatchedAt(id uint, agentID uint) error

	// メッセージの送受信ごとにLastSendAtを更新する
	UpdateLastSendAtByThreadID(threadID uint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.ChatGroupWithAgent, error)

	FindByIDAndAgentID(id, agentID uint) (*entity.ChatGroupWithAgent, error)

	FindByAgentID(agent1ID, agent2ID uint) (*entity.ChatGroupWithAgent, error)

	/** 複数取得 */
	GetByAgentID(agentID uint) ([]*entity.ChatGroupWithAgent, error)

	GetByMyAgentIDAndOtherIDList(myAgentID uint, otherAgentIDList []uint) ([]*entity.ChatGroupWithAgent, error)

	All() ([]*entity.ChatGroupWithAgent, error)
}

// エージェントチャットメッセージの返信先
type ChatThreadWithAgentRepository interface {
	/** 作成 */
	Create(chatThread *entity.ChatThreadWithAgent) error

	/** 単数取得 */
	FindByID(id uint) (*entity.ChatThreadWithAgent, error)

	/** 複数取得 */
	GetByGroupID(groupID uint) ([]*entity.ChatThreadWithAgent, error)

	All() ([]*entity.ChatThreadWithAgent, error)
}

type ChatMessageWithAgentRepository interface {
	/** 作成 */
	Create(chatMessageWithAgent *entity.ChatMessageWithAgent) error

	/** 複数取得 */
	GetByThreadID(threadID uint) ([]*entity.ChatMessageWithAgent, error)

	All() ([]*entity.ChatMessageWithAgent, error)
}

// エージェントチャットメッセージの宛先
type ChatMessageToUserWithAgentRepository interface {
	/** 作成 */
	Create(chatMessageToUser *entity.ChatMessageToUserWithAgent) error

	/** 更新 */
	UpdateWatchedAtByThreadID(threadID, agentStaffID uint) error

	/** 複数取得 */
	GetByThreadID(threadID uint) ([]*entity.ChatMessageToUserWithAgent, error)

	GetByGroupIDAndUnwatched(groupID, agentStaffID uint) ([]*entity.ChatMessageToUserWithAgent, error)

	GetByAgentStaffIDAndUnwatched(agentStaffID uint) ([]*entity.ChatMessageToUserWithAgent, error)

	GetByMessageIDList(idList []uint) ([]*entity.ChatMessageToUserWithAgent, error)

	All() ([]*entity.ChatMessageToUserWithAgent, error)
}

/****************************************************************************************/
// タスクの評価
//

// 面談調整タスクグループ
type InterviewTaskGroupRepository interface {
	/** 作成 */
	// 面談調整のタスクグループを作成
	Create(interviewTaskGroup *entity.InterviewTaskGroup) error

	/** 更新 */
	// 最終リクエスト日時を更新
	UpdateLastRequestAt(id uint) error
	// 最終閲覧日時を更新
	UpdateLastWatchedAt(id uint) error
	// 面談日時を更新
	UpdateInterviewDate(id uint, interviewDate time.Time) error
	// 初回面談日時と通常の面談日時を更新
	UpdateNormalAndFirstInterviewDate(id uint, normalInterviewDate, firstInterviewDate time.Time) error

	/** 単数取得 */
	// IDから面談調整のタスクグループを取得
	FindByID(id uint) (*entity.InterviewTaskGroup, error)

	// IDと求職者IDから面談調整のタスクグループを取得
	FindByAgentIDAndJobSeekerID(agentID, jobSeekerID uint) (*entity.InterviewTaskGroup, error)

	/** 複数取得 */
	// スケジュール
	GetByStaffIDAndPeriod(staffID uint, startDate, endDate time.Time) ([]*entity.InterviewTaskGroup, error)
	GetByStaffIDAndPeriodByStaffIDList(idList []uint, startDate, endDate time.Time) ([]*entity.InterviewTaskGroup, error)

	// 売上管理のカウント関連（全体）
	GetInterviewOfferAcceptancePerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)
	GetInterviewPerformanceCountByAgentIDAndSalesMonth(agentID uint, saleMonth string) (float64, error)

	// 売上管理のカウント関連（個人）
	GetInterviewOfferAcceptancePerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)
	GetInterviewPerformanceCountByStaffIDAndSalesMonth(staffID uint, saleMonth string) (float64, error)

	// ダッシュボード 売上管理のカウント関連（全体）
	GetInterviewOfferAcceptancePerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error)
	GetInterviewPerformanceCountByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) (float64, error)

	// ダッシュボード 売上管理のカウント関連（個人）
	GetInterviewOfferAcceptancePerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error)
	GetInterviewPerformanceCountByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) (float64, error)
}

// 面談調整タスク
type InterviewTaskRepository interface {
	/** 作成 */
	// 面談調整のタスクを作成
	Create(interviewTask *entity.InterviewTask) error

	/** 削除 */
	// 面談調整のタスクを削除
	Delete(id uint) error

	/** 単数取得 */
	// IDから面談調整のタスクを取得
	FindByID(id uint) (*entity.InterviewTask, error)

	// エージェントIDと求職者IDから最新の面談調整のタスクを取得
	FindLatestByAgentIDAndJobSeekerID(agentID, jobSeekerID uint) (*entity.InterviewTask, error)

	/** 複数取得 */
	// グループIDから面談調整のタスクを取得
	GetByGroupID(groupID uint) ([]*entity.InterviewTask, error)

	// エージェントIDから最新の面談調整のタスクを取得（phaseが「0 or 1」）
	GetLatestAdjustmentByAgentID(agentID uint) ([]*entity.InterviewTask, error)

	// エージェントIDから最新の参加確認のタスクを取得（phaseが「2 or 3」）
	GetLatestConfirmationByAgentID(agentID uint) ([]*entity.InterviewTask, error)
}

// 面談調整テンプレート
type InterviewAdjustmentTemplateRepository interface {
	/** 作成 */
	// 面談調整テンプレートを作成
	Create(template *entity.InterviewAdjustmentTemplate) error

	/** 更新 */
	// 面談調整テンプレートを更新
	Update(id uint, template *entity.InterviewAdjustmentTemplate) error

	/** 削除 */
	// 面談調整テンプレートを削除
	Delete(id uint) error

	/** 単数取得 */
	// IDから面談調整テンプレートを取得
	FindByID(id uint) (*entity.InterviewAdjustmentTemplate, error)

	/** 複数取得 */
	// 担当者IDと送信シーンからメッセージテンプレートを取得
	GetByAgentIDAndSendScene(agentID, sendScene uint) ([]*entity.InterviewAdjustmentTemplate, error)

	// エージェントIDから面談調整テンプレートを取得
	GetByAgentID(agentID uint) ([]*entity.InterviewAdjustmentTemplate, error)
}

/****************************************************************************************/
// スカウトサービス
//

// エージェントのスカウトサービス
type ScoutServiceRepository interface {
	/** 作成 */
	Create(scoutService *entity.ScoutService) error

	/** 更新 */
	Update(id uint, scoutService *entity.ScoutService) error

	UpdatePassword(param entity.UpdateScoutServicePasswordParam) error

	UpdateLastSendCount(id, lastSendCount uint) error

	/** 削除 */
	Delete(id uint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.ScoutService, error)

	/** 複数取得 */
	GetByAgentID(agentID uint) ([]*entity.ScoutService, error)

	GetByAgentRobotID(agentRobotID uint) ([]*entity.ScoutService, error)
}

// スカウトテンプレート
type ScoutServiceTemplateRepository interface {
	/** 作成 */
	Create(template *entity.ScoutServiceTemplate) error

	/** 更新 */
	UpdateLastSend(id, lastSendCount uint, lastSendAt time.Time) error

	/** 削除 */
	DeleteByScoutServiceID(scoutServiceID uint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.ScoutServiceTemplate, error)

	/** 複数取得 */
	GetByScoutServiceID(scoutServiceID uint) ([]*entity.ScoutServiceTemplate, error)

	GetByAgentID(agentID uint) ([]*entity.ScoutServiceTemplate, error)

	GetByAgentRobotID(agentRobotID uint) ([]*entity.ScoutServiceTemplate, error)

	GetByIDList(idList []uint) ([]*entity.ScoutServiceTemplate, error)
}

// スカウトサービスの求職者取得時刻
type ScoutServiceGetEntryTimeRepository interface {
	/** 作成 */
	Create(scoutServiceGetEntryTime *entity.ScoutServiceGetEntryTime) error

	/** 削除 */
	DeleteByScoutServiceID(scoutServiceID uint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.ScoutServiceGetEntryTime, error)

	/** 複数取得 */
	GetByScoutServiceID(scoutServiceID uint) ([]*entity.ScoutServiceGetEntryTime, error)

	GetByAgentID(agentID uint) ([]*entity.ScoutServiceGetEntryTime, error)

	GetByAgentRobotID(agentRobotID uint) ([]*entity.ScoutServiceGetEntryTime, error)
}

// 売上管理
type SaleRepository interface {
	/** 作成 */
	Create(sale *entity.Sale) error

	/** 更新 */
	Update(id uint, sale *entity.Sale) error

	/** 単数取得 */
	FindByID(id uint) (*entity.Sale, error)
	FindByJobSeekerID(seeker_id uint) (*entity.Sale, error)
	FindByJobSeekerIDAndJobInformationID(seekerID, jobInfoID uint) (*entity.Sale, error)

	/** 複数取得 */
	GetByAgentStaffID(agentStaffID uint) ([]*entity.Sale, error)
	GetByIDList(idList []uint) ([]*entity.Sale, error)
	GetSearchByAgent(searchParam entity.SearchAccuracy) ([]*entity.Sale, error)

	GetByAgentIDForMonthly(agentID uint, startMonth, endMonth string) ([]*entity.Sale, error)
	GetByStaffIDForMonthly(agentStaffID uint, startMonth, endMonth string) ([]*entity.Sale, error)
	GetContractSignedByStaffIDAndMonth(agentStaffID uint, thisMonth string) ([]*entity.Sale, error)
	GetBillingByStaffIDAndMonth(agentStaffID uint, thisMonth string) ([]*entity.Sale, error)
	GetContractSignedByAgentIDAndMonth(agentID uint, month string) ([]*entity.Sale, error)
	GetBillingByAgentIDAndMonth(agentID uint, month string) ([]*entity.Sale, error)

	// ダッシュボード
	GetContractSignedByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) ([]*entity.Sale, error)
	GetBillingByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) ([]*entity.Sale, error)
	GetContractSignedByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) ([]*entity.Sale, error)
	GetBillingByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) ([]*entity.Sale, error)

	// csv出力用
	GetByAgentIDForCSV(agentID uint) ([]*entity.Sale, error)
}

// Agent売上目標管理
type AgentMonthlySaleRepository interface {
	/** 作成 */
	// Agent売上目標管理を作成する
	Create(agentMonthlySale *entity.AgentMonthlySale) error

	/** 更新 */
	// Agent売上目標管理を更新する
	Update(agentMonthlySale *entity.AgentMonthlySale, agentMonthlySaleID uint) error

	/** 単数取得 */
	// 指定のManagementIDと月でAgent売上目標管理を取得する
	FindByManagementIDAndMonth(managementID uint, month string) (*entity.AgentMonthlySale, error)

	/** 複数取得 */
	// 指定のManagementIDでAgent売上目標管理をすべて取得する
	GetByManagementID(managementID uint) ([]*entity.AgentMonthlySale, error)

	// 指定のManagementIDとその期間でAgent売上目標管理をすべて取得する
	GetByManagementIDAndPeriod(managementID uint, startMonth, endMonth string) ([]*entity.AgentMonthlySale, error)
}

// AgentStaff売上目標管理
type AgentStaffMonthlySaleRepository interface {
	/** 作成 */
	// AgentStaff売上目標管理を作成する
	Create(agentStaffMonthlySale *entity.AgentStaffMonthlySale) error

	/** 更新 */
	// AgentStaff売上目標管理を更新する
	Update(agentStaffMonthlySale *entity.AgentStaffMonthlySale, agentMonthlySaleID uint) error

	/** 単数取得 */
	// 指定のManagementIDと月でAgentStaff売上目標管理を取得する
	FindByStaffManagementIDAndMonth(staffManagementID uint, thisMonth string) (*entity.AgentStaffMonthlySale, error)

	/** 複数取得 */
	// 指定のagent_staff_monthly_salesテーブルのStaffManagementIDでAgentStaff売上目標管理理をすべて取得する
	GetByStaffManagementID(staffManagementID uint) ([]*entity.AgentStaffMonthlySale, error)

	// 指定のagentSaleManagementIDとagent_staff_sale_managementsテーブルのmanagement_id合致する、AgentStaff売上目標管理をすべて取得する
	GetByAgentSaleManagementID(agentSaleManagementID uint) ([]*entity.AgentStaffMonthlySale, error)

	// 指定のagent_staff_monthly_salesテーブルのStaffManagementIDとその期間でAgentStaff売上目標管理理をすべて取得する
	GetByStaffManagementIDAndPeriod(staffManagementID uint, startMonth, endMonth string) ([]*entity.AgentStaffMonthlySale, error)
}

// エージェントの売上管理の決算月や公開状況の管理
type AgentSaleManagementRepository interface {
	/** 作成 */
	// エージェントの売上管理の決算月や公開状況の管理を作成する
	Create(agentMonthlySale *entity.AgentSaleManagement) error

	/** 更新 */
	//エージェントの売上管理の決算月や公開状況の管理を更新する
	Update(agentMonthlySale *entity.AgentSaleManagement, agentsaleManagementID uint) error

	//　指定のagentIDとsaleManagment以外のエージェントの売上管理の決算月や公開状況の管理を更新する
	UpdateIsOpenOtherThanID(agentID, saleManagementID uint) error // 指定IDのisOpen以外をCloseに更新

	/** 単数取得 */
	// IDからエージェントの売上管理の決算月や公開状況の管理を取得する
	FindByID(id uint) (*entity.AgentSaleManagement, error)

	// IDとAgentIDからエージェントの売上管理の決算月や公開状況の管理を取得する
	FindByIDAndAgentID(id, agentID uint) (*entity.AgentSaleManagement, error)

	// AgentIDとis_openがtrueのエージェントの売上管理の決算月や公開状況の管理を取得する
	FindByAgentIDAndIsOpen(agentID uint) (*entity.AgentSaleManagement, error)

	// agentIDから最新のエージェントの売上管理の決算月や公開状況の管理を取得する
	FindLatestByAgentID(agentID uint) (*entity.AgentSaleManagement, error)

	/** 複数取得 */
	// 指定のagentIDのすべてエージェントの売上管理の決算月や公開状況の管理を取得する
	GetByAgentID(agentID uint) ([]*entity.AgentSaleManagement, error)
}

// エージェントスタッフの売上管理の決算月や公開状況を管理する
type AgentStaffSaleManagementRepository interface {
	/** 作成 */
	// エージェントスタッフの売上管理の決算月や公開状況の管理を作成する
	Create(staffMonthlySale *entity.AgentStaffSaleManagement) error

	/** 単数取得 */
	// managementIDとagentStaffIDからエージェントスタッフの売上管理の決算月や公開状況の管理を取得する
	FindByManagementIDAndStaffID(managementID, staffID uint) (*entity.AgentStaffSaleManagement, error)

	/** 複数取得 */
	// agentStaffIDからすべてのエージェントスタッフの売上管理の決算月や公開状況の管理を取得する
	GetByAgentID(agentID uint) ([]*entity.AgentStaffSaleManagement, error)

	// managementIDからすべてのエージェントスタッフの売上管理の決算月や公開状況の管理とagent_staffsテーブルと結合して担当者名を取得を取得する
	GetStaffNameByManagementID(managementID uint) ([]*entity.AgentStaffSaleManagement, error)

	//  agentStaffIDから決算期間が当月を含むエージェントスタッフの売上管理の決算月や公開状況を管理を全て取得
	GetByStaffIDAndThidMonth(staffID uint, thisMonth string) ([]*entity.AgentStaffSaleManagement, error)
}

type JobSeekerScheduleRepository interface {
	/** 作成 */
	Create(schedule *entity.JobSeekerSchedule) error

	/** 更新 */
	Update(id uint, schedule *entity.JobSeekerSchedule) error
	UpdateTitle(id uint, title string) error
	UpdateIsSharedByIDList(idList []uint) error

	/** 削除 */
	Delete(id uint) error
	DeleteByTaskGroupIDAndScheduleType(taskGroupID, scheduleType uint) error
	DeleteScheduleInTaskGroupAboveTaskID(taskID, taskGroupID uint) error

	/** 複数取得 */
	GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerSchedule, error)
	GetByJobSeekerIDAndScheuldType(jobSeekerID, scheduleType uint) ([]*entity.JobSeekerSchedule, error)
	GetByTaskGroupID(taskGroupID uint) ([]*entity.JobSeekerSchedule, error)
	GetByTaskGroupIDList(taskGroupIDList []uint) ([]*entity.JobSeekerSchedule, error)
	GetByTaskGroupIDAndScheduleType(taskGroupID, scheduleType uint) ([]*entity.JobSeekerSchedule, error)
	GetByStaffIDAndPeriod(agentStaffID, scheduleType uint, startDate, endDate time.Time) ([]*entity.JobSeekerSchedule, error)
	GetByStaffIDAndPeriodByStaffIDList(idList []uint, scheduleType uint, startDate, endDate time.Time) ([]*entity.JobSeekerSchedule, error)
}

type JobSeekerRescheduleRepository interface {
	/** 作成 */
	Create(reschedule *entity.JobSeekerReschedule) error

	/** 削除 */
	Delete(id uint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.JobSeekerReschedule, error)
}

/****************************************************************************************/
/****************************************************************************************/
// 他媒体求人インポート関連
//

// 他媒体の求人企業を一括インポートするテーブル
type InitialEnterpriseImporterRepository interface {
	/** 作成 */
	Create(initialEnterpriseImporter *entity.InitialEnterpriseImporter) error

	// 実行済みに更新
	UpdateIsSuccessTrue(id uint) error

	/** 削除 */
	Delete(id uint) error

	/** 単数取得 */
	FindByID(id uint) (*entity.InitialEnterpriseImporter, error)

	/** 複数取得 */
	// エージェントIDから求人インポートの一覧を取得
	GetByAgentID(agentID uint) ([]*entity.InitialEnterpriseImporter, error)

	// 今後1週間の企業求人一括インポート情報一覧を取得 *予約日程調整時の重複チェック用
	GetByWeek() ([]*entity.InitialEnterpriseImporter, error)

	// 未実行の実行日時が現在日時を取得
	GetByStartDate(now time.Time) ([]*entity.InitialEnterpriseImporter, error)
}

/****************************************************************************************/
/****************************************************************************************/
// ユーザー向けお知らせ関連
//

// お知らせを管理するテーブル
type NotificationForUserRepository interface {
	/** 作成 */
	// お知らせを作成
	Create(notificationForUser *entity.NotificationForUser) error

	/** 単数取得 */
	// 指定IDのお知らせを取得
	FindByID(id uint) (*entity.NotificationForUser, error)

	// 最新のお知らせを取得
	FindLatest() (*entity.NotificationForUser, error)

	/** 複数取得 */
	// 指定ID以降、ターゲットのお知らせを取得
	GetAfterIDAndTargetList(id uint, targetList []entity.NotificationForUserTarget) ([]*entity.NotificationForUser, error)

	// 指定ターゲットのお知らせを取得
	GetByTargetList(targetList []entity.NotificationForUserTarget) ([]*entity.NotificationForUser, error)

	// 全てのお知らせを取得
	All() ([]*entity.NotificationForUser, error)
}

// お知らせを確認したユーザーをを管理するテーブル（ユーザーが閲覧したらレコードを追加）
type UserNotificationViewRepository interface {
	/** 作成 */
	// お知らせの既読を作成
	Create(userNotificationView *entity.UserNotificationView) error

	/** 単数取得 */
	// 指定IDのお知らせ既読レコードを取得
	FindByID(id uint) (*entity.UserNotificationView, error)

	// 最新のお知らせ既読レコードを取得
	FindLatestByAgentStaffID(agentStaffID uint) (*entity.UserNotificationView, error)

	/** 複数取得 */
	// 指定お知らせIDのお知らせ既読レコードを取得
	GetByNotificationID(notificationID uint) ([]*entity.UserNotificationView, error)

	// 指定お知らせIDリストのお知らせ既読レコードを取得
	GetByNotificationIDList(notificationIDList []uint) ([]*entity.UserNotificationView, error)

	// 全てのお知らせ既読レコードを取得
	All() ([]*entity.UserNotificationView, error)
}

/****************************************************************************************/

/****************************************************************************************/
// デプロイ情報関連
//

type DeploymentInformationRepository interface {
	/** 作成 */
	// デプロイ情報を作成
	Create(deploymentInfo *entity.DeploymentInformation) error

	/** 更新 */
	// デプロイ情報を更新
	Update(id uint, deploymentInfo *entity.DeploymentInformation) error

	/** 複数取得 */
	All() ([]*entity.DeploymentInformation, error)
}

// デプロイの反映状況
type DeploymentReflectionRepository interface {
	/** 作成 */
	// デプロイの反映状況を複数作成
	CreateMulti(deploymentID uint, agentStaffList []*entity.AgentStaff) error

	/** 更新 */
	// デプロイの反映状況フラグを更新
	UpdateIsReflectedByAgentStaffID(staffID uint, isReflected bool) error

	/** 単数取得 */
	// 指定スタッフIDのis_reflectedがtrueのリストを取得
	FindNotReflectedByAgentStaffID(staffID uint) (*entity.DeploymentReflection, error)
}

/****************************************************************************************/
// Google認証情報
//
type GoogleAuthenticationRepository interface {
	Create(googleAuth *entity.GoogleAuthentication) error
	Update(id uint, googleAuth *entity.GoogleAuthentication) error
	UpdateTokenAndExpiry(id uint, accessToken string, expiry time.Time) error
	FindLatest() (*entity.GoogleAuthentication, error)
}

/****************************************************************************************/
// エントリー処理実ユーザー
//
type UserEntryRepository interface {
	Create(userEntry *entity.UserEntry) error
	UpdateIsProcessedByUserID(userID string, isProcessed bool) error
	UpdateIsProcessedByUserIDList(userIDList []string, isProcessed bool) error
	GetUnprocessed() ([]*entity.UserEntry, error)
}
