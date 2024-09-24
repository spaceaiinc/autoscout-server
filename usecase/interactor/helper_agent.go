package interactor

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"gopkg.in/guregu/null.v4"
)

// 自社のAgentとAllianceを区別して、AllianceのAgent型を取得する
func getAllianceAgentList(myAgentID uint, allianceList []*entity.AgentAlliance) []*entity.Agent {
	var allianceAgentList []*entity.Agent

	if len(allianceList) < 1 {
		return allianceAgentList
	}

	// アライアンス相手の情報を取得する
	for _, alliance := range allianceList {
		var agentID uint
		// myAgentID が Agent1ID と等しい場合、Agent2ID を使用
		if alliance.Agent1ID == myAgentID {
			agentID = alliance.Agent2ID
		}
		// myAgentID が Agent2ID と等しい場合、Agent1ID を使用
		if alliance.Agent2ID == myAgentID {
			agentID = alliance.Agent1ID
		}

		// Agent を作成し、リストに追加
		agent := entity.NewAgent(alliance.AgentName, alliance.OfficeLocation, alliance.Representative, alliance.Establish, "", "", "", null.NewInt(0, false), "", "", "", "", "", "", "", "", "", false, false, false, null.NewInt(0, false))
		agent.ID = agentID // ID を設定
		allianceAgentList = append(allianceAgentList, agent)
	}

	return allianceAgentList
}

// エージェントグループの最大ページ数を返す（本番実装までは1ページあたり20件）
func getAgentListMaxPage(agentList []*entity.Agent) uint {
	var maxPage = len(agentList) / 20

	if 0 < (len(agentList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページのエージェントグループ一覧を返す（本番実装までは1ページあたり20件）
func getAgentListWithPage(agentList []*entity.Agent, page uint) []*entity.Agent {
	var (
		perPage uint = 20
		listLen uint = uint(len(agentList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return agentList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.Agent{}
	}

	if (listLen - first) <= perPage {
		return agentList[first:]
	}
	return agentList[first:last]
}
