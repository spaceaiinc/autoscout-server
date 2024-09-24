package handler

import (
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type InitialEnterpriseImporterHandler interface {
	// 汎用系 API
	CreateInitialEnterpriseImporter(param entity.InitialEnterpriseImporter) (presenter.Presenter, error)
	DeleteInitialEnterpriseImporter(id uint) (presenter.Presenter, error)
	GetInitialEnterpriseImporterListByAgentID(agentID uint) (presenter.Presenter, error)
	GetInitialEnterpriseImporterListByWeek() (presenter.Presenter, error) // 今後1週間の企業求人一括インポート情報一覧を取得 *予約日程調整時の重複チェック用

	// 他サービスから求人取得・更新
	BatchInitialEnterpriseImporter(now time.Time) (presenter.Presenter, error) // 他サービスから求人取得用 API(テスト用)
	BatchUpdateEnterpriseForAgentBank(now time.Time) (presenter.Presenter, error)
}

type InitialEnterpriseImporterHandlerImpl struct {
	initialEnterpriseImporterInteractor interactor.InitialEnterpriseImporterInteractor
}

func NewInitialEnterpriseImporterHandlerImpl(epI interactor.InitialEnterpriseImporterInteractor) InitialEnterpriseImporterHandler {
	return &InitialEnterpriseImporterHandlerImpl{
		initialEnterpriseImporterInteractor: epI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *InitialEnterpriseImporterHandlerImpl) CreateInitialEnterpriseImporter(param entity.InitialEnterpriseImporter) (presenter.Presenter, error) {
	output, err := h.initialEnterpriseImporterInteractor.CreateInitialEnterpriseImporter(interactor.CreateInitialEnterpriseImporterInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInitialEnterpriseImporterJSONPresenter(responses.NewInitialEnterpriseImporter(output.InitialEnterpriseImporter)), nil
}

// 求人インポートの削除
func (h *InitialEnterpriseImporterHandlerImpl) DeleteInitialEnterpriseImporter(id uint) (presenter.Presenter, error) {
	output, err := h.initialEnterpriseImporterInteractor.DeleteInitialEnterpriseImporter(interactor.DeleteInitialEnterpriseImporterInput{
		ID: id,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// エージェントIDから求人インポートの一覧を取得する
func (h *InitialEnterpriseImporterHandlerImpl) GetInitialEnterpriseImporterListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.initialEnterpriseImporterInteractor.GetInitialEnterpriseImporterListByAgentID(interactor.GetInitialEnterpriseImporterListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInitialEnterpriseImporterListJSONPresenter(responses.NewInitialEnterpriseImporterList(output.InitialEnterpriseImporterList)), nil
}

// 今後1週間の企業求人一括インポート情報一覧を取得 *予約日程調整時の重複チェック用
func (h *InitialEnterpriseImporterHandlerImpl) GetInitialEnterpriseImporterListByWeek() (presenter.Presenter, error) {
	output, err := h.initialEnterpriseImporterInteractor.GetInitialEnterpriseImporterListByWeek()

	if err != nil {
		return nil, err
	}

	return presenter.NewInitialEnterpriseImporterListJSONPresenter(responses.NewInitialEnterpriseImporterList(output.InitialEnterpriseImporterList)), nil
}

/****************************************************************************************/
/****************************************************************************************/
// 他サービスから求人取得・更新
//

// 他媒体の求人企業を一括インポートするテーブル
func (h *InitialEnterpriseImporterHandlerImpl) BatchInitialEnterpriseImporter(now time.Time) (presenter.Presenter, error) {
	output, err := h.initialEnterpriseImporterInteractor.BatchInitialEnterpriseImporter(interactor.BatchInitialEnterpriseImporterInput{
		Now: now,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 他媒体の求人企業を一括インポートするテーブル
func (h *InitialEnterpriseImporterHandlerImpl) BatchUpdateEnterpriseForAgentBank(now time.Time) (presenter.Presenter, error) {
	output, err := h.initialEnterpriseImporterInteractor.BatchUpdateEnterpriseForAgentBank(interactor.BatchUpdateEnterpriseForAgentBankInput{
		Now: now,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
