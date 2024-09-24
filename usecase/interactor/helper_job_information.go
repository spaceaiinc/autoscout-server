package interactor

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"gopkg.in/guregu/null.v4"
)

// 求人の最大ページ数を返す（本番実装までは1ページあたり20件）
func getJobInformationListMaxPage(jobInformationList []*entity.JobInformation) uint {
	var maxPage = len(jobInformationList) / 20
	if 0 < (len(jobInformationList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページの求人一覧を返す（本番実装までは1ページあたり20件）
func getJobInformationListWithPage(jobInformationList []*entity.JobInformation, page uint) []*entity.JobInformation {
	var (
		perPage uint = 20
		listLen uint = uint(len(jobInformationList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return jobInformationList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.JobInformation{}
	}

	if (listLen - first) <= perPage {
		return jobInformationList[first:]
	}
	return jobInformationList[first:last]
}

// 非公開エージェントに含まれているかどうかをチェックする
func checkJobInformationByHideToAgent(jobInformationList []*entity.JobInformation, hideToAgentList []*entity.JobInformationHideToAgent) []*entity.JobInformation {
	var (
		outputList         []*entity.JobInformation
		hideToAgentChecker = make(map[uint]bool)
	)

	// 非公開エージェントIDリストが空の場合は求人リストをそのまま返す
	if len(hideToAgentList) == 0 {
		return jobInformationList
	}

	// hideToAgentで弾かれる求人IDを取得
	for _, hta := range hideToAgentList {
		hideToAgentChecker[hta.JobInformationID] = true
	}

	// 検索元エージェントと一致する非公開エージェント情報を取得
	// 求人一覧から非公開エージェントに含まれているものを除外
	for _, jobInfo := range jobInformationList {
		if !hideToAgentChecker[jobInfo.ID] {
			outputList = append(outputList, jobInfo)
		}
	}

	return outputList
}

// 求人リストからIDリストを取得する
func getJobInformationIDList(jobInformationList []*entity.JobInformation) []uint {
	var idListUint []uint

	if len(jobInformationList) == 0 {
		return idListUint
	}

	for _, jobInformation := range jobInformationList {
		idListUint = append(idListUint, jobInformation.ID)
	}

	return idListUint
}

// 求人リストから企業IDリストを取得する
func getEnterpriseIDList(jobInformationList []*entity.JobInformation) []uint {
	var enterpriseIDList []uint

	if len(jobInformationList) == 0 {
		return enterpriseIDList
	}

	for _, jobInformation := range jobInformationList {
		enterpriseIDList = append(enterpriseIDList, jobInformation.EnterpriseID)
	}

	return enterpriseIDList
}

// 求人リストから求人IDリストと企業IDリストを取得する
func getJobInformationIDListAndEnterpriseIDList(jobInformationList []*entity.JobInformation) ([]uint, []uint) {
	var (
		lenJobInformationList = len(jobInformationList)
		jobInformationIDList  = make([]uint, 0, lenJobInformationList)
		enterpriseIDList      []uint
	)

	if lenJobInformationList == 0 {
		return jobInformationIDList, enterpriseIDList
	}

	for _, jobInformation := range jobInformationList {
		jobInformationIDList = append(jobInformationIDList, jobInformation.ID)
		enterpriseIDList = append(enterpriseIDList, jobInformation.EnterpriseID)
	}

	return jobInformationIDList, enterpriseIDList
}

// 会社名 + 郵便番号が一致する場合は除外する
func excludeDuplicateJobInformation(
	jobInformationListBeforeDuplicate []*entity.JobInformation,
	agentID uint,
) []*entity.JobInformation {
	var (
		jobInformationList                      []*entity.JobInformation
		companyNameAndPostCodeCheckerOwn        = make(map[string]uint) // 自社の会社名・郵便番号の重複
		companyNameAndPostCodeCheckerOther      = make(map[string]uint) // 他社の会社名・郵便番号の重複
		companyNameAndPostCodeAndAgentIDChecker = make(map[string]uint)
	)

	// 自社と他社エージェントで同一社名の求人がある場合は他社の同一社名の求人を除外する
	if len(jobInformationListBeforeDuplicate) > 0 {
		for _, jobInformation := range jobInformationListBeforeDuplicate {
			//この形式にする→「株式会社Motoyui_614-0014」（社名にスペースが入る事も考慮してスペース削除の処理も実行）
			cp := strings.Join(strings.Fields(jobInformation.CompanyName+"_"+jobInformation.PostCode), "")
			cpa := cp + "_" + strconv.Itoa(int(jobInformation.AgentID))

			if jobInformation.AgentID == agentID {
				companyNameAndPostCodeCheckerOwn[cp] = companyNameAndPostCodeCheckerOwn[cp] + 1
			} else {
				companyNameAndPostCodeCheckerOther[cp] = companyNameAndPostCodeCheckerOther[cp] + 1
			}
			companyNameAndPostCodeAndAgentIDChecker[cpa] = companyNameAndPostCodeAndAgentIDChecker[cpa] + 1
		}

		for _, jobInformation := range jobInformationListBeforeDuplicate {
			//この形式にする→「株式会社Motoyui_614-0014」（社名にスペースが入る事も考慮してスペース削除の処理も実行）
			cp := strings.Join(strings.Fields(jobInformation.CompanyName+"_"+jobInformation.PostCode), "")

			if jobInformation.AgentID == agentID {
				// 自社求人は取得
				jobInformationList = append(jobInformationList, jobInformation)
			} else if jobInformation.AgentID != agentID && (companyNameAndPostCodeCheckerOwn[cp] == 0 && companyNameAndPostCodeCheckerOther[cp] > 0) {
				// 自社求人との社名被りが無いシェア求人のみ取得
				jobInformationList = append(jobInformationList, jobInformation)
			}
		}
	}

	// 他社同士の重複チェック
	if len(jobInformationList) > 0 {
		for _, jobInformation := range jobInformationList {
			//この形式にする→「株式会社Motoyui_614-0014」（社名にスペースが入る事も考慮してスペース削除の処理も実行）
			cp := strings.Join(strings.Fields(jobInformation.CompanyName+"_"+jobInformation.PostCode), "")
			cpa := cp + "_" + strconv.Itoa(int(jobInformation.AgentID))

			if (companyNameAndPostCodeCheckerOwn[cp] + companyNameAndPostCodeCheckerOther[cp]) > companyNameAndPostCodeAndAgentIDChecker[cpa] {
				jobInformation.IsDuplicate = true
			}
		}
	}

	return jobInformationList
}

// シェア求人で「雇用形態」が「2: 派遣社員 or 3: 紹介予定派遣」のみの場合は除外する
// 仮に表示してこの求人が稼働した場合、二重派遣で法的にアウトなるための措置
func filterAllianceJobInformationExcludedTemporaryWorker(i *JobInformationInteractorImpl, allianceJobInformationList []*entity.JobInformation) ([]*entity.JobInformation, error) {

	var jobInformationListExcludedTemporaryWorker []*entity.JobInformation

	idList, _ := getJobInformationIDListAndEnterpriseIDList(allianceJobInformationList)

	// 該当の雇用形態の情報を取得
	employmentStatuses, err := i.jobInfoEmploymentStatusRepository.GetByJobInformationIDList(idList)
	if err != nil {
		fmt.Println(err)
		return allianceJobInformationList, err
	}

	for _, j := range allianceJobInformationList {
		var result = false

		for _, es := range employmentStatuses {
			if j.ID == es.JobInformationID {
				j.EmploymentStatuses = append(j.EmploymentStatuses, *es)
			}
		}

		if len(j.EmploymentStatuses) == 0 {
			// 雇用形態登録していない場合はresultをtrueにしてループを抜ける
			result = true
		} else {
			// 「2: 派遣社員 or 3: 紹介予定派遣」以外の雇用形態が登録されている場合はresultをtrueにしてループを抜ける
			for _, es := range j.EmploymentStatuses {
				if es.EmploymentStatus != null.NewInt(2, true) && es.EmploymentStatus != null.NewInt(3, true) {
					result = true
					break
				}
			}
		}

		// 「2: 派遣社員 or 3: 紹介予定派遣」以外の雇用形態が登録されている求人のみをセットし直す
		if result {
			jobInformationListExcludedTemporaryWorker = append(jobInformationListExcludedTemporaryWorker, j)
		}
	}

	allianceJobInformationList = jobInformationListExcludedTemporaryWorker

	return allianceJobInformationList, nil
}

// 特別仕様: 本番環境のみ「2: 株式会社テスト」と「3: 株式会社Motoyui（非公開求人管理用）」を除外して他社エージェントに非表示にするための関数
func excludeTestJobInformation(
	jobInformationList []*entity.JobInformation,
	agentID uint,
) []*entity.JobInformation {
	env := os.Getenv("APP_ENV")

	/**
	 * 条件
	 *
	 * ユーザーが「1: 株式会社Motoyui」と「2: 株式会社テスト」と「3: 株式会社Motoyui（非公開求人管理用）」以外で
	 * 求人の担当エージェントが「2: 株式会社テスト」と「3: 株式会社Motoyui（非公開求人管理用）」の場合
	 * スライスから除外する
	**/

	if env == "prd" {
		for i := 0; i < len(jobInformationList); i++ {

			// ユーザーが「1: 株式会社Motoyui」と「2: 株式会社テスト」と「3: 株式会社Motoyui（非公開求人管理用）」以外の場合
			if agentID != 1 && agentID != 2 && agentID != 3 {
				// 求人の担当エージェントが「2: 株式会社テスト」と「3: 株式会社Motoyui（非公開求人管理用）」の場合
				if jobInformationList[i].AgentID == 2 || jobInformationList[i].AgentID == 3 {
					// 除外する
					jobInformationList = append(jobInformationList[:i], jobInformationList[i+1:]...)
					i-- // スライスの要素が前にシフトされたため、現在のインデックスを調整する
				}
			}
		}
	}

	return jobInformationList
}
