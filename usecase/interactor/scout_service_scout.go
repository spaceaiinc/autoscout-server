package interactor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"gopkg.in/guregu/null.v4"

	// ブラウザ操作
	"github.com/go-rod/rod"
	rodInput "github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

/****************************************************************************************/
// Batch処理 API
//
/*
1時間ごとに該当するスカウト送信を実行する
*/
type BatchScoutInput struct {
	Now          time.Time
	AgentRobotID uint
}

type BatchScoutOutput struct {
	OK bool
}

func (i *ScoutServiceInteractorImpl) BatchScout(input BatchScoutInput) (BatchScoutOutput, error) {
	var (
		output                                      BatchScoutOutput
		err                                         error
		errMessage                                  string
		selectedScoutService                        *entity.ScoutService
		selectedScoutServiceTemplateList            []*entity.ScoutServiceTemplate
		scoutServiceTemplateListForRan              []*entity.ScoutServiceTemplate
		scoutServiceTemplateListForMynaviScouting   []*entity.ScoutServiceTemplate
		scoutServiceTemplateListForAmbi             []*entity.ScoutServiceTemplate
		scoutServiceTemplateListForMynaviAgentScout []*entity.ScoutServiceTemplate
	)

	agentRobot, err := i.agentRobotRepository.FindByID(input.AgentRobotID)
	if err != nil {
		log.Println("エージェントロボットの取得に失敗しました", err)
		return output, err
	}

	if !agentRobot.IsScoutActive {
		return output, nil
	}

	scoutServices, err := i.scoutServiceRepository.GetByAgentRobotID(input.AgentRobotID)
	if err != nil {
		log.Println("スカウトサービスの取得に失敗しました", err)
		return output, err
	}

	scoutServiceTemplates, err := i.scoutServiceTemplateRepository.GetByAgentRobotID(input.AgentRobotID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// タイムアウトを設定
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Minute)
	defer cancel()

	// panicが発生した場合にエラーを返すためにdeferでrecoverを実行して、エラーとして返す
	defer func() {
		if rec := recover(); rec != nil {
			// stack traceを取得
			stackStr := string(debug.Stack())

			// interactor or repositoryを含む文字列を取得
			errorLine := ""
			for _, line := range strings.Split(stackStr, "\n") {
				if strings.Contains(line, "interactor/") || strings.Contains(line, "repository/") {
					errorLine += line + "\n"
					log.Println("error line:", errorLine)
				}
			}
			errMessage = fmt.Sprintf("スカウトサービス処理でpanicエラーが発生しました。\nRecover: %v\nStack:\n%s", rec, errorLine)
			log.Println(errMessage)
			err = errors.New(errMessage)

			if selectedScoutService == nil {
				return
			}

			// エラーの原因を判定
			errorCause := "パニックエラー"
			if strings.Contains(fmt.Sprint(rec), "deadline") {
				errorCause = "タイムアウト"
			}

			// 更新後のスカウトサービステンプレートを取得
			scoutServiceTemplateIDList := make([]uint, 0)
			for _, scoutServiceTemplate := range selectedScoutServiceTemplateList {
				scoutServiceTemplateIDList = append(scoutServiceTemplateIDList, scoutServiceTemplate.ID)
			}

			updatedScoutServiceTemplateList, err := i.scoutServiceTemplateRepository.GetByIDList(scoutServiceTemplateIDList)
			if err != nil {
				log.Println(err)
			}

			// エラーメールの内容を作成
			searchAndMessageTitle := ""
			for _, updatedScoutServiceTemplate := range updatedScoutServiceTemplateList {
				isSuccessStr := "成功"
				lastSendAtStr := updatedScoutServiceTemplate.LastSendAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02 15:04:05")
				// 12時間以内の場合は成功とする
				if updatedScoutServiceTemplate.LastSendAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Before(input.Now.Add(-12 * time.Hour)) {
					isSuccessStr = "失敗"
					updatedScoutServiceTemplate.LastSendCount = null.NewInt(0, false)
					lastSendAtStr = ""
				}

				searchTitle := updatedScoutServiceTemplate.SearchTitle
				// マイナビスカウティングの場合はメッセージタイトルを表示する
				if selectedScoutService.ServiceType.Int64 == entity.ScoutServiceTypeMynaviScouting {
					searchTitle = updatedScoutServiceTemplate.MessageTitle
				}

				searchAndMessageTitle += fmt.Sprintf(
					"検索タイトル: %s, 送信時間: %v, 送信数(件): %v, 成功/失敗: %s\n\n",
					searchTitle,
					lastSendAtStr,
					updatedScoutServiceTemplate.LastSendCount.Int64,
					isSuccessStr,
				)
			}

			now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
			errMailMessage := fmt.Sprintf(
				"%sのスカウト送信に失敗しました。\n\n・送信進捗:\n%s\n\n・発生時刻: %s\n・ロボット名: %s\n・ロボットID: %v\n・エラー内容:%s\n・Recover: %v\n・Stack:\n%s",
				entity.ScoutServiceTypeLabel[selectedScoutService.ServiceType.Int64],
				searchAndMessageTitle, now, agentRobot.Name, agentRobot.ID, errorCause,
				rec, errorLine,
			)
			log.Println(errMailMessage)

			i.sendErrorMail(errMailMessage)
		}
	}()

	/********
	 スカウト送信
	********/

	// 曜日の判定
	var (
		weekdayToday = input.Now.Weekday()
		isMonday     bool
		isTuesday    bool
		isWednesday  bool
		isThursday   bool
		isFriday     bool
		isSaturday   bool
		isSunday     bool
	)
	switch weekdayToday {
	case time.Monday:
		isMonday = true
	case time.Tuesday:
		isTuesday = true
	case time.Wednesday:
		isWednesday = true
	case time.Thursday:
		isThursday = true
	case time.Friday:
		isFriday = true
	case time.Saturday:
		isSaturday = true
	case time.Sunday:
		isSunday = true
	}

	// スカウトサービステンプレートを媒体ごとに分ける
	for _, scoutService := range scoutServices {
		if !scoutService.IsActive {
			continue
		}
		selectedScoutService = scoutService

		for _, scoutServiceTemplate := range scoutServiceTemplates {
			if scoutService.ID == scoutServiceTemplate.ScoutServiceID &&
				scoutServiceTemplate.StartHour == null.NewInt(int64(input.Now.Hour()), true) &&
				((scoutServiceTemplate.RunOnMonday && isMonday) ||
					(scoutServiceTemplate.RunOnTuesday && isTuesday) ||
					(scoutServiceTemplate.RunOnWednesday && isWednesday) ||
					(scoutServiceTemplate.RunOnThursday && isThursday) ||
					(scoutServiceTemplate.RunOnFriday && isFriday) ||
					(scoutServiceTemplate.RunOnSaturday && isSaturday) ||
					(scoutServiceTemplate.RunOnSunday && isSunday)) {

				// 媒体ごとにスカウトサービステンプレートを分ける
				switch scoutService.ServiceType.Int64 {
				case entity.ScoutServiceTypeRan:
					scoutServiceTemplateListForRan = append(scoutServiceTemplateListForRan, scoutServiceTemplate)
				case entity.ScoutServiceTypeMynaviScouting:
					scoutServiceTemplateListForMynaviScouting = append(scoutServiceTemplateListForMynaviScouting, scoutServiceTemplate)
				case entity.ScoutServiceTypeAmbi:
					scoutServiceTemplateListForAmbi = append(scoutServiceTemplateListForAmbi, scoutServiceTemplate)
				case entity.ScoutServiceTypeMynaviAgentScout:
					scoutServiceTemplateListForMynaviAgentScout = append(scoutServiceTemplateListForMynaviAgentScout, scoutServiceTemplate)
				}
			}
		}

		// RAN
		if len(scoutServiceTemplateListForRan) > 0 {
			log.Println("RANのスカウト送信を開始します")
			selectedScoutServiceTemplateList = scoutServiceTemplateListForRan
			_, err = i.ScoutOnRan(ScoutOnRanInput{
				ScoutService:             scoutService,
				ScoutServiceTemplateList: scoutServiceTemplateListForRan,
				Context:                  ctx,
			})
			break

			// マイナビスカウティング
		} else if len(scoutServiceTemplateListForMynaviScouting) > 0 {
			log.Println("マイナビスカウティングのスカウト送信を開始します")
			selectedScoutServiceTemplateList = scoutServiceTemplateListForMynaviScouting
			_, err = i.ScoutOnMynaviScouting(ScoutOnMynaviScoutingInput{
				ScoutService:             scoutService,
				ScoutServiceTemplateList: scoutServiceTemplateListForMynaviScouting,
				Context:                  ctx,
			})
			break

			// AMBI
		} else if len(scoutServiceTemplateListForAmbi) > 0 {
			log.Println("AMBIのスカウト送信を開始します")
			selectedScoutServiceTemplateList = scoutServiceTemplateListForAmbi
			_, err = i.ScoutOnAmbi(ScoutOnAmbiInput{
				ScoutService:             scoutService,
				ScoutServiceTemplateList: scoutServiceTemplateListForAmbi,
				Context:                  ctx,
			})
			break

			// マイナビエージェントスカウト
		} else if len(scoutServiceTemplateListForMynaviAgentScout) > 0 {
			log.Println("マイナビエージェントスカウトのスカウト送信を開始します")
			selectedScoutServiceTemplateList = scoutServiceTemplateListForMynaviAgentScout
			_, err = i.ScoutOnMynaviAgentScout(ScoutOnMynaviAgentScoutInput{
				ScoutService:             scoutService,
				ScoutServiceTemplateList: scoutServiceTemplateListForMynaviAgentScout,
				Context:                  ctx,
			})
			break
		}

	}

	// エラーが発生した場合
	if err != nil {
		// 送信進捗をメールで送信
		scoutServiceTemplateIDList := make([]uint, 0)
		for _, scoutServiceTemplate := range selectedScoutServiceTemplateList {
			scoutServiceTemplateIDList = append(scoutServiceTemplateIDList, scoutServiceTemplate.ID)
		}

		updatedScoutServiceTemplateList, err := i.scoutServiceTemplateRepository.GetByIDList(scoutServiceTemplateIDList)
		if err != nil {
			log.Println(err)
			return output, err
		}

		searchAndMessageTitle := ""
		for _, updatedScoutServiceTemplate := range updatedScoutServiceTemplateList {
			isSuccessStr := "成功"
			lastSendAtStr := updatedScoutServiceTemplate.LastSendAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02 15:04:05")
			// 12時間以内の場合は成功とする
			if updatedScoutServiceTemplate.LastSendAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Before(input.Now.Add(-12 * time.Hour)) {
				isSuccessStr = "失敗"
				updatedScoutServiceTemplate.LastSendCount = null.NewInt(0, false)
				lastSendAtStr = ""
			}

			searchTitle := updatedScoutServiceTemplate.SearchTitle
			// マイナビスカウティングの場合はメッセージタイトルを表示する
			if selectedScoutService.ServiceType.Int64 == entity.ScoutServiceTypeMynaviScouting {
				searchTitle = updatedScoutServiceTemplate.MessageTitle
			}

			searchAndMessageTitle += fmt.Sprintf(
				"検索タイトル: %s, 送信時間: %v, 送信数(件): %v, 成功/失敗: %s\n\n",
				searchTitle,
				lastSendAtStr,
				updatedScoutServiceTemplate.LastSendCount.Int64,
				isSuccessStr,
			)

		}

		now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
		errMessage = fmt.Sprintf(
			"%vのスカウト送信に失敗しました。\n\n・送信進捗:\n%s\n\n・発生時刻: %s\n\n・ロボット名: %s\n\n・ロボットID: %v",
			entity.ScoutServiceTypeLabel[selectedScoutService.ServiceType.Int64], searchAndMessageTitle, now, agentRobot.Name, agentRobot.ID,
		)
		log.Println(errMessage)
		i.sendErrorMail(errMessage)
		return output, errors.New(errMessage)
	}

	output.OK = true
	return output, err
}

/****************************************************************************************/
// スカウト送信 API
//
/*
RANでスカウトを送信する

開発環境ログ 2つのテンプレートでヒットした人数分送信
2023/05/10 17:07:53 RANのスカウト送信を開始します

*/
type ScoutOnRanInput struct {
	Context                  context.Context
	ScoutService             *entity.ScoutService
	ScoutServiceTemplateList []*entity.ScoutServiceTemplate
}

type ScoutOnRanOutput struct {
	OK bool
}

func (i *ScoutServiceInteractorImpl) ScoutOnRan(input ScoutOnRanInput) (ScoutOnRanOutput, error) {
	var (
		output           ScoutOnRanOutput
		err              error
		errMessage       string
		browser          *rod.Browser
		page             *rod.Page
		searchResultPage *rod.Page
		searchResultCnt  int
	)

	// ホストOSのchromeパスを取得
	path, ok := launcher.LookPath()
	if !ok {
		errMessage = "chrome, chromium or microsoft edgeが見つかりません"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// 取得したパスからchromeを起動
	l := launcher.New().
		Bin(path)

	// ローカル環境の場合は画面ありで起動
	if i.app.Env == "local" {
		l.Headless(false)
	}

	defer time.Sleep(10 * time.Second)
	defer l.Cleanup() // remove launcher.FlagUserDataDir

	url := l.MustLaunch()

	// ブラウザを起動後、AMBIへ遷移
	browser = rod.New().
		ControlURL(url).
		Trace(true).                 // ログを出力
		SlowMotion(2 * time.Second). // 機械対策のページ用に遅延を入れる
		MustConnect()

	defer browser.MustClose()

	// browserにタイムアウトを設定
	browser, cancel := browser.
		Context(input.Context).
		WithCancel()
	defer cancel()

	// ログインページへ遷移
	page = browser.
		MustPage("https://ran.next.rikunabi.com/rnc/docs/ca_s01000.jsp")
	page.WaitLoad()
	time.Sleep(10 * time.Second)

	// ログイン情報を入力
	page.MustElement("input[name=login_nm]").
		MustInput(input.ScoutService.LoginID)

	// パスワードの複号
	decryptedPassword, err := decryption(input.ScoutService.Password)
	if err != nil {
		log.Println(err)
		return output, err
	}

	page.MustElement("input[name=pswd]").
		MustInput(decryptedPassword).
		MustType(rodInput.Enter)

	page.WaitLoad()
	time.Sleep(10 * time.Second)

	// URLがログインページのままの場合はログインに失敗していると判断
	if strings.Contains(page.MustInfo().URL, "ca_i01000.jsp") {
		errMessage = "ログインに失敗しました。"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	for _, scoutServiceTemplate := range input.ScoutServiceTemplateList {
		// ダミーから主要情報が載っているframeへ切り替え
		pageFrame := page.MustElement("frame[name=lowerframe]").MustFrame()

		// 保存済みの検索条件をクリック
		iconLinks := pageFrame.MustElements("a.icn_link02")
		if len(iconLinks) == 0 {
			errMessage = "保存済みの検索条件要素リストが取得できませんでした。"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}
		var isMatchedWithSearchLink bool
		for _, iconLink := range iconLinks {
			if iconLink.MustText() == "保存済み検索条件一覧" {
				iconLink.MustClick()
				time.Sleep(10 * time.Second)
				isMatchedWithSearchLink = true
				break
			}
		}
		log.Println("isMatchedWithSearchLink:", isMatchedWithSearchLink)
		if !isMatchedWithSearchLink {
			errMessage = "保存済みの検索条件要素が取得できませんでした。"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		/**
		通常スカウトの場合
		*/
		if scoutServiceTemplate.ScoutType == null.NewInt(entity.RanScoutTypeNormal, true) {
			// 保存済みの検索条件をクリック
			log.Println("保存済みの検索条件をクリック")
			var isMatchedWithSearch bool
			searchFrame := page.MustElement("frame[name=lowerframe]").MustFrame()
			searchTables := searchFrame.MustElements("table.cell_middle")
			for searchI, searchTable := range searchTables {
				if searchI == 1 {
					searchTrs := searchTable.MustElements("tr")
					for trI, searchTr := range searchTrs {
						if trI == 0 ||
							!searchTr.MustHas("p.large.bold") {
							log.Println("continue")
							continue
						}
						titleEl := searchTr.
							MustElement("p.large.bold")

						titleText := titleEl.MustText()
						log.Println("titleText", titleText)

						if titleText == scoutServiceTemplate.SearchTitle {
							searchTr.MustElement("a.btn01.btn_m.large").MustClick()
							searchFrame.WaitLoad()
							time.Sleep(10 * time.Second)
							isMatchedWithSearch = true
							break
						}
					}
				}
			}
			if !isMatchedWithSearch {
				errMessage = "保存済みの検索条件が一致しませんでした。入力した保存条件:" + scoutServiceTemplate.SearchTitle
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 検索結果数を取得
			pageFrame = page.MustElement("frame[name=lowerframe]").MustFrame()
			searchResultPage = pageFrame
			if !searchResultPage.MustHas("span.num_font.num_medium.red") {
				errMessage = "検索結果数の取得に失敗しました。対象者がいない可能性があります。"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			searchResultCntStr := searchResultPage.MustElement("span.num_font.num_medium.red").MustText()
			searchResultCnt, err = strconv.Atoi(searchResultCntStr)
			log.Println("searchResultCnt:", searchResultCnt)
			if err != nil {
				errMessage = "検索結果数の取得に失敗しました"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			if searchResultCnt == 0 {
				errMessage = "検索結果がありません"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 検索結果の全て選択(300件まで)
			for searchI := 0; searchI < searchResultCnt/100+1; searchI++ {
				// 検索結果数が指定の最大送信を超えている場合
				if int64(searchI) > scoutServiceTemplate.ScoutCount.Int64/100 {
					log.Println("検索結果数が最大送信を超えているため、選択を終了します")
					break
				}

				searchResultPage.MustElement("input.all_check.chkbox").MustClick()
				searchResultPage.WaitLoad()
				time.Sleep(10 * time.Second)

				// 一括送信候補者一覧へ移動させる
				sendBtnList := searchResultPage.MustElements("a.btn03.btn_s.w200")
				isMatchedWithSendBtn := false
				for _, sendBtn := range sendBtnList {
					sendBtnText, _ := sendBtn.Text()
					log.Println("sendBtnText:", sendBtnText)
					if sendBtnText == "一括送信候補者一覧へ移動" {
						sendBtn.MustClick()
						searchResultPage.WaitLoad()
						time.Sleep(10 * time.Second)
						isMatchedWithSendBtn = true
						break
					}
				}
				if !isMatchedWithSendBtn {
					errMessage = "一括送信候補者一覧へ移動ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// 3ページ目の場合は終了
				if searchI == 2 {
					break
				}

				// 次のページへ
				pagingLinksDiv := searchResultPage.MustElement("div.f_right.right.pt5")
				pagingLinks := pagingLinksDiv.MustElements("a")
				for _, pagingLink := range pagingLinks {
					pagingLinkText, err := pagingLink.Text()
					if err != nil {
						continue
					}
					if pagingLinkText == "次へ" {
						href := pagingLink.MustAttribute("href")
						log.Println("次のページへ移動します. href:", *href)
						searchResultPage = browser.MustPage("https://ran.next.rikunabi.com/rnc/docs/" + *href)
						searchResultPage.WaitLoad()
						time.Sleep(10 * time.Second)
						break
					}
				}
			}

			// 余分なページを閉じる
			openPages := browser.MustPages()
			for _, openPage := range openPages {
				log.Println("openPage:", openPage.MustInfo())
				if strings.Contains(openPage.MustInfo().URL, "ca_s02000") {
					openPage.MustClose()
				}
			}

			time.Sleep(10 * time.Second)
			linkTd := pageFrame.MustElement("td.bottom")

			var isMatchedWithSendAtOnceLink bool
			links := linkTd.MustElements("a")
			for _, link := range links {
				if link.MustText() == "一括送信候補者一覧へ" {
					link.MustClick()
					pageFrame.WaitLoad()
					time.Sleep(10 * time.Second)
					isMatchedWithSendAtOnceLink = true
					break
				}
			}
			if !isMatchedWithSendAtOnceLink {
				errMessage = "一括送信候補者一覧へ移動ボタンが見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 100件ずつ送信する(300件まで)
			for sendI := 0; sendI < searchResultCnt/100+1; sendI++ {
				if sendI >= int(scoutServiceTemplate.ScoutCount.Int64/100) {
					log.Println("指定された数の候補者を選択しました\nindex:", sendI, "\nscoutCount:", scoutServiceTemplate.ScoutCount.Int64)
					break
				} else if sendI > 2 {
					log.Println("最大送信300件までです\nindex:", sendI, "\nsearchResultCnt:", searchResultCnt)
					break
				}

				// すべての候補者を選択
				// if searchResultCnt/100 >= sendI+1 || searchResultCnt%100 == 0 {
				pageFrame.MustElement("input.all_check.chkbox").
					MustClick()
				pageFrame.WaitLoad()
				time.Sleep(10 * time.Second)
				// }

				// 指定された数の候補者を選択
				// } else {
				// 	selectSendUserInputsFromSavedList := pageFrame.MustElements("input[name=snd_trgt_mbr_chk]")
				// 	if len(selectSendUserInputsFromSavedList) == 0 || len(selectSendUserInputsFromSavedList) < int(scoutServiceTemplate.ScoutCount.Int64) {
				// 		errMessage = "保存済みリストの候補者が足りません"
				// 		log.Println(errMessage)
				// 		return output, errors.New(errMessage)
				// 	}
				// 	for selectSendUserI, selectSendUserInput := range selectSendUserInputsFromSavedList {
				// 		if selectSendUserI < searchResultCnt {
				// 			selectSendUserInput.MustClick()
				// 		} else {
				// 			break
				// 		}
				// 	}
				// }

				// 求人IDを選択ページへ遷移するボタンをクリック
				selectJobInfoLinks := pageFrame.MustElements("a.btn03.w100.btn_s.bold")
				var isMatchedWithSelectJobInfoLink bool
				for _, selectJobInfoLink := range selectJobInfoLinks {
					if selectJobInfoLink.MustText() == "求人ID選択" {
						selectJobInfoLink.MustClick()
						pageFrame.WaitLoad()
						time.Sleep(10 * time.Second)
						isMatchedWithSelectJobInfoLink = true
						break
					}
				}
				if !isMatchedWithSelectJobInfoLink {
					errMessage = "求人IDを選択ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// 求人選択ページへ遷移
				var selectJobInfoPage *rod.Page
				pages := browser.MustPages()
				for _, page := range pages {
					if strings.Contains(page.MustInfo().URL, "ca_s02070.jsp") {
						selectJobInfoPage = page
						break
					}
				}
				if selectJobInfoPage == nil {
					errMessage = "求人IDを選択ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				selectJobInfoPageBody := selectJobInfoPage.MustElement("body")
				if !selectJobInfoPageBody.MustHas("div#pp_contents") {
					errMessage = "オファー送信対象でない候補者が含まれています。該当候補者を対象から外してください。"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// メインのbodyのdiv要素を取得
				selectJobInfoPageDiv := selectJobInfoPageBody.
					MustElement("div#pp_contents")

				// 求人IDを選択
				if !selectJobInfoPageDiv.MustHas("table#tmpltList") {
					errMessage = "求人IDを選択ページのテーブルが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}
				selectJobInfoTrs := selectJobInfoPageDiv.
					MustElement("table#tmpltList").
					MustElements("tr")
				var isMatchedWithJobInfo bool
			selectJobInfoLoop:
				for selectJobInfoI, selectJobInfoTr := range selectJobInfoTrs {
					if selectJobInfoI == 0 {
						continue
					}
					if selectJobInfoTr.MustHas("td") {
						selectJobInfoTds := selectJobInfoTr.MustElements("td")
						for selectJobInfoTdI, selectJobInfoTd := range selectJobInfoTds {
							if selectJobInfoTdI == 1 {
								savedJobInfoIDStr := selectJobInfoTd.MustText()
								log.Println("savedJobInfoID", savedJobInfoIDStr)

								if savedJobInfoIDStr == scoutServiceTemplate.JobInformationID {
									selectJobInfoTds[0].MustElement("a.btn01.btn_s.w80").MustClick()
									time.Sleep(2 * time.Second)
									isMatchedWithJobInfo = true
									break selectJobInfoLoop
								}
							}
						}
					}
				}

				if !isMatchedWithJobInfo {
					errMessage = "求人IDが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// 求人IDを反映して閉じる
				selectJobInfoPageDiv.MustElement("a.btn01.btn_m.large.w200").MustClick()
				time.Sleep(10 * time.Second)

				// 一括オファーを作成ボタンをクリック
				createOfferBtn := pageFrame.MustElement("a.btn01.btn_m.bold.large.w200.middle")
				if createOfferBtn == nil || createOfferBtn.MustText() != "一括オファー作成" {
					errMessage = "一括オファーを作成ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				createOfferBtn.MustClick()
				pageFrame.WaitLoad()
				time.Sleep(10 * time.Second)

				// メッセージテンプレートを選択
				messageTemplateSelect := pageFrame.MustElement("select#tmplt_id")
				err = messageTemplateSelect.Select(
					[]string{
						scoutServiceTemplate.MessageTitle,
					},
					true,
					rod.SelectorTypeText,
				)
				if err != nil {
					errMessage = "メッセージテンプレートが見つかりませんでした。DBから取得したタイトル" + scoutServiceTemplate.SearchTitle
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// テンプレートを反映
				applyTemplateBtn := pageFrame.MustElement("a.btn03.btn_s.w150")
				if applyTemplateBtn == nil || applyTemplateBtn.MustText() != "本文テンプレートを反映" {
					errMessage = "テンプレートを反映ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}
				// クリック後のwindow.confirmをOKにする
				waitConfirm, handleConfirm := pageFrame.MustHandleDialog()
				go applyTemplateBtn.MustClick()
				waitConfirm()
				handleConfirm(true, "")

				time.Sleep(2 * time.Second)

				// 送信内容確認ボタンをクリック
				confirmBtn := pageFrame.MustElement("a.btn01.btn_m.large.w200")
				if confirmBtn == nil || confirmBtn.MustText() != "送信内容確認" {
					errMessage = "送信内容確認ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				confirmBtn.MustClick()
				pageFrame.WaitLoad()
				time.Sleep(10 * time.Second)

				// 送信内容確認画面へ移動
				confirmPage := page.MustElement("frame[name=lowerframe]").MustFrame()

				// 送信ボタンをクリック
				sendBtn := confirmPage.MustElement("a.btn01.btn_m.large.w300")
				log.Println("送信ボタン。sendBtn: ", sendBtn.MustText())
				if sendBtn == nil || sendBtn.MustText() != "送信" {
					errMessage = "送信ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				if i.app.Env == "prd" {
					sendBtn.MustClick()
					confirmPage.WaitLoad()
					time.Sleep(10 * time.Second)

					// 一括送信候補者一覧へ戻る
					returnBtnList := confirmPage.MustElements("a")
					for _, returnBtn := range returnBtnList {
						if returnBtn.MustText() == "一括送信候補者一覧へ戻る" {
							returnBtn.MustClick()
							confirmPage.WaitLoad()
							time.Sleep(10 * time.Second)
							break
						}
					}
				}

				// 次のテンプレートのためにスカウトのトップページへ戻る
				page.MustNavigate("https://ran.next.rikunabi.com/rnc/docs/ca_f00010.jsp?tab_select_id=1&screen_id=ca_s02120&__u=16849106369894720112352085858491")
				page.WaitLoad()
				time.Sleep(10 * time.Second)
				pageFrame = page.MustElement("frame[name=lowerframe]").MustFrame()
			}

			log.Println("スカウトサービスを実行しました")

			// 次のテンプレートのためにスカウトのトップページへ戻る
			page.MustNavigate("https://ran.next.rikunabi.com/rnc/docs/ca_f00010.jsp?__u=16806743487066738118737323900111")
			page.WaitLoad()
			time.Sleep(10 * time.Second)

			/**
			再度スカウトの場合
			*/
		} else if scoutServiceTemplate.ScoutType == null.NewInt(entity.RanScoutTypeAgain, true) {
			// スカウト済み候補者一覧
			pageFrame := page.MustElement("frame[name=lowerframe]").MustFrame()
			btns := pageFrame.MustElements("a")
			if len(btns) == 0 {
				errMessage = "スカウト済み候補者一覧要素リストが取得できませんでした。"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			var isMatchedWithScoutedCandidateList bool
			for _, btn := range btns {
				if btn.MustText() == "スカウト済み候補者一覧" {
					btn.MustClick()
					time.Sleep(10 * time.Second)
					isMatchedWithScoutedCandidateList = true
					break
				}
			}
			log.Println("isMatchedWithScoutedCandidateList:", isMatchedWithScoutedCandidateList)
			if !isMatchedWithScoutedCandidateList {
				errMessage = "スカウト済み候補者一覧へ遷移できませんでした。"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 保存済みの検索条件をクリック
			log.Println("保存済みの検索条件をクリック")
			var isMatchedWithSearch bool
			searchFrame := page.MustElement("frame[name=lowerframe]").MustFrame()
			searchTables := searchFrame.MustElements("table.cell_middle")
			for _, searchTable := range searchTables {
				searchTrs := searchTable.MustElements("tr")
				for trI, searchTr := range searchTrs {
					if trI == 0 ||
						!searchTr.MustHas("p.large.bold") {
						log.Println("continue")
						continue
					}
					titleEl := searchTr.
						MustElement("p.large.bold")

					titleText := titleEl.MustText()
					log.Println("titleText", titleText)

					if titleText == scoutServiceTemplate.SearchTitle {
						searchTr.MustElement("a.btn01.btn_m.large").MustClick()
						searchFrame.WaitLoad()
						time.Sleep(10 * time.Second)
						isMatchedWithSearch = true
						break
					}
				}
			}

			if !isMatchedWithSearch {
				errMessage = "保存済みの検索条件が一致しませんでした。入力した保存条件:" + scoutServiceTemplate.SearchTitle
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 検索結果数を取得
			pageFrame = page.MustElement("frame[name=lowerframe]").MustFrame()
			searchResultPage = pageFrame
			if !searchResultPage.MustHas("span.num_font.num_medium.red") {
				errMessage = "検索結果数の取得に失敗しました。対象者がいない可能性があります。"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			searchResultCntStr := searchResultPage.MustElement("span.num_font.num_medium.red").MustText()
			searchResultCnt, err = strconv.Atoi(searchResultCntStr)
			log.Println("searchResultCnt:", searchResultCnt)
			if err != nil {
				errMessage = "検索結果数の取得に失敗しました"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			if searchResultCnt == 0 {
				errMessage = "検索結果がありません"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 100件ずつ送信する
			for searchI := 0; searchI > searchResultCnt/100; searchI++ {
				if searchI >= int(scoutServiceTemplate.ScoutCount.Int64/100) {
					log.Println("指定された最大数を超えました\nindex:", searchI, "\nscoutCount:", scoutServiceTemplate.ScoutCount.Int64)
					break
				} else if searchI > 2 {
					log.Println("3ページ分送信しました\nindex:", searchI, "\nsearchResultCnt:", searchResultCnt)
					break
				}

				// すべての候補者を選択
				pageFrame.MustElement("input.all_check.chkbox").
					MustClick()
				pageFrame.
					WaitLoad()
				time.Sleep(10 * time.Second)

				// 再送ボタンをクリック
				createOfferBtns := pageFrame.MustElements("a")
				var isMatchedWithReSendBtn bool
				for _, createOfferBtn := range createOfferBtns {
					if createOfferBtn.MustText() == "再送" {
						createOfferBtn.MustClick()
						pageFrame.WaitLoad()
						time.Sleep(10 * time.Second)
						isMatchedWithReSendBtn = true
						break
					}
				}
				if !isMatchedWithReSendBtn {
					errMessage = "再送ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// メッセージテンプレートを選択
				messageTemplateSelect := pageFrame.MustElement("select#tmplt_id")
				err = messageTemplateSelect.Select(
					[]string{
						scoutServiceTemplate.MessageTitle,
					},
					true,
					rod.SelectorTypeText,
				)
				if err != nil {
					errMessage = "メッセージテンプレートが見つかりませんでした。DBから取得したタイトル" + scoutServiceTemplate.SearchTitle
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// テンプレートを反映
				applyTemplateBtn := pageFrame.MustElement("a.btn03.btn_s.w150")
				if applyTemplateBtn == nil || applyTemplateBtn.MustText() != "本文テンプレートを反映" {
					errMessage = "テンプレートを反映ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}
				// クリック後のwindow.confirmをOKにする
				waitConfirm, handleConfirm := pageFrame.MustHandleDialog()
				go applyTemplateBtn.MustClick()
				waitConfirm()
				handleConfirm(true, "")

				time.Sleep(2 * time.Second)

				// 送信内容確認ボタンをクリック
				confirmBtn := pageFrame.MustElement("a.btn01.btn_m.large.w200")
				if confirmBtn == nil || confirmBtn.MustText() != "送信内容確認" {
					errMessage = "送信内容確認ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				confirmBtn.MustClick()
				pageFrame.WaitLoad()
				time.Sleep(10 * time.Second)

				// 送信内容確認画面へ移動
				confirmPage := page.MustElement("frame[name=lowerframe]").MustFrame()

				// 送信ボタンをクリック
				sendBtn := confirmPage.MustElement("a.btn01.btn_m.large.w300")
				log.Println("送信ボタン。sendBtn: ", sendBtn.MustText())
				if sendBtn == nil || sendBtn.MustText() != "送信" {
					errMessage = "送信ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				if i.app.Env == "prd" {
					sendBtn.MustClick()
					confirmPage.WaitLoad()
					time.Sleep(10 * time.Second)

					// 一括送信候補者一覧へ戻る
					returnBtnList := confirmPage.MustElements("a")
					for _, returnBtn := range returnBtnList {
						if returnBtn.MustText() == "一括送信候補者一覧へ戻る" {
							returnBtn.MustClick()
							confirmPage.WaitLoad()
							time.Sleep(10 * time.Second)
							break
						}
					}
					// 次のテンプレートのためにスカウトのトップページへ戻る
					// page.MustNavigate("https://ran.next.rikunabi.com/rnc/docs/ca_f00010.jsp?tab_select_id=1&screen_id=ca_s02010&src_cnd_tab_id=A&__u=16861253765566400756907650246970")
					page.WaitLoad()
					time.Sleep(10 * time.Second)
					pageFrame = page.MustElement("frame[name=lowerframe]").MustFrame()
				} else {
					break
				}
			}

			log.Println("スカウトサービスを実行しました")

			// 次のテンプレートのためにスカウトのトップページへ戻る
			page.MustNavigate("https://ran.next.rikunabi.com/rnc/docs/ca_f00010.jsp?__u=16806743487066738118737323900111")
			page.WaitLoad()
			time.Sleep(10 * time.Second)

			/**
			別求人送付の場合
			*/
		} else if scoutServiceTemplate.ScoutType == null.NewInt(entity.RanScoutTypeSendOther, true) {
			// スカウト済み候補者一覧
			candidateFrame := page.MustElement("frame[name=lowerframe]").MustFrame()
			btns := candidateFrame.MustElements("a")
			if len(btns) == 0 {
				errMessage = "スカウト済み候補者一覧要素リストが取得できませんでした。"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			var isMatchedWithScoutedCandidateList bool
			for _, btn := range btns {
				if btn.MustText() == "スカウト済み候補者一覧" {
					btn.MustClick()
					candidateFrame.WaitLoad()
					time.Sleep(10 * time.Second)
					isMatchedWithScoutedCandidateList = true
					break
				}
			}
			log.Println("isMatchedWithScoutedCandidateList:", isMatchedWithScoutedCandidateList)
			if !isMatchedWithScoutedCandidateList {
				errMessage = "スカウト済み候補者一覧へ遷移できませんでした。"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 保存済みの検索条件をクリック
			log.Println("保存済みの検索条件をクリック")
			var isMatchedWithSearch bool
			searchFrame := page.MustElement("frame[name=lowerframe]").MustFrame()
			log.Println("searchFrame:", searchFrame)
			searchTables := searchFrame.MustElements("table.cell_middle")
			log.Println("searchTables:", searchTables)
			for _, searchTable := range searchTables {
				searchTrs := searchTable.MustElements("tr")
				for trI, searchTr := range searchTrs {
					if trI == 0 ||
						!searchTr.MustHas("p.large.bold") {
						log.Println("continue")
						continue
					}
					titleEl := searchTr.
						MustElement("p.large.bold")

					titleText := titleEl.MustText()
					log.Println("titleText", titleText)

					if titleText == scoutServiceTemplate.SearchTitle {
						searchTr.MustElement("a.btn01.btn_m.large").MustClick()
						searchFrame.WaitLoad()
						time.Sleep(10 * time.Second)
						isMatchedWithSearch = true
						break
					}
				}
			}

			if !isMatchedWithSearch {
				errMessage = "保存済みの検索条件が一致しませんでした。入力した保存条件:" + scoutServiceTemplate.SearchTitle
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 検索結果数を取得
			pageFrame = page.MustElement("frame[name=lowerframe]").MustFrame()
			searchResultPage = pageFrame
			if !searchResultPage.MustHas("span.num_font.num_medium.red") {
				errMessage = "検索結果数の取得に失敗しました。対象者がいない可能性があります。"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			searchResultCntStr := searchResultPage.MustElement("span.num_font.num_medium.red").MustText()
			searchResultCnt, err = strconv.Atoi(searchResultCntStr)
			log.Println("searchResultCnt:", searchResultCnt)
			if err != nil {
				errMessage = "検索結果数の取得に失敗しました"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			if searchResultCnt == 0 {
				errMessage = "検索結果がありません"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 検索結果の全て選択(300件まで)
			for searchI := 0; searchI < searchResultCnt/100+1; searchI++ {
				// 検索結果数が指定の最大送信を超えている場合は終了
				if int64(searchI) > scoutServiceTemplate.ScoutCount.Int64/100 {
					log.Println("検索結果数が最大送信を超えているため、選択を終了します")
					break
				}

				searchResultPage.MustElement("input.all_check.chkbox").MustClick()
				searchResultPage.WaitLoad()
				time.Sleep(10 * time.Second)

				// 一括送信させる
				sendBtnList := searchResultPage.MustElements("a")
				isMatchedWithSendBtn := false
				for _, sendBtn := range sendBtnList {
					sendBtnText, _ := sendBtn.Text()
					log.Println("sendBtnText:", sendBtnText)
					if sendBtnText == "一括送信" {
						sendBtn.MustClick()
						searchResultPage.WaitLoad()
						time.Sleep(10 * time.Second)
						isMatchedWithSendBtn = true
						break
					}
				}
				if !isMatchedWithSendBtn {
					errMessage = "一括送信ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// 前のページに戻る
				searchResultPage.NavigateBack()
				searchResultPage.WaitLoad()
				time.Sleep(10 * time.Second)

				// 3ページ目の場合は終了(最大送信が300件までのため)
				if searchI == 2 {
					break
				}

				// 次のページへ
				pagingLinksDiv := searchResultPage.MustElement("div.f_right.right.pt5")
				pagingLinks := pagingLinksDiv.MustElements("a")
				for _, pagingLink := range pagingLinks {
					pagingLinkText, err := pagingLink.Text()
					if err != nil {
						continue
					}
					if pagingLinkText == "次へ" {
						href := pagingLink.MustAttribute("href")
						log.Println("次のページへ移動します. href:", *href)
						searchResultPage = browser.MustPage("https://ran.next.rikunabi.com/rnc/docs/" + *href)
						searchResultPage.WaitLoad()
						time.Sleep(10 * time.Second)
						break
					}
				}
			}

			// 余分なページを閉じる
			openPages := browser.MustPages()
			for _, openPage := range openPages {
				log.Println("openPage:", openPage.MustInfo())
				if strings.Contains(openPage.MustInfo().URL, "ca_s02200") {
					openPage.MustClose()
				}
			}

			time.Sleep(10 * time.Second)
			linkTd := pageFrame.MustElement("td.bottom")

			var isMatchedWithSendAtOnceLink bool
			links := linkTd.MustElements("a")
			for _, link := range links {
				if link.MustText() == "一括送信候補者一覧へ" {
					link.MustClick()
					pageFrame.WaitLoad()
					time.Sleep(10 * time.Second)
					isMatchedWithSendAtOnceLink = true
					break
				}
			}
			if !isMatchedWithSendAtOnceLink {
				errMessage = "一括送信候補者一覧へ移動ボタンが見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 100件ずつ送信する
			for sendI := 0; sendI < searchResultCnt/100+1; sendI++ {
				if sendI >= int(scoutServiceTemplate.ScoutCount.Int64/100) {
					log.Println("指定された数の候補者を選択しました\nindex:", sendI, "\nscoutCount:", scoutServiceTemplate.ScoutCount.Int64)
					break
				} else if sendI >= 3 {
					log.Println("最大送信300件までです\nindex:", sendI, "\nsearchResultCnt:", searchResultCnt)
					break
				}

				// すべての候補者を選択
				pageFrame.MustElement("input.all_check.chkbox").
					MustClick()
				pageFrame.WaitLoad()
				time.Sleep(10 * time.Second)

				// 求人IDを選択ページへ遷移するボタンをクリック
				selectJobInfoLinks := pageFrame.MustElements("a.btn03.w100.btn_s.bold")
				var isMatchedWithSelectJobInfoLink bool
				for _, selectJobInfoLink := range selectJobInfoLinks {
					if selectJobInfoLink.MustText() == "求人ID選択" {
						selectJobInfoLink.MustClick()
						pageFrame.WaitLoad()
						time.Sleep(10 * time.Second)
						isMatchedWithSelectJobInfoLink = true
						break
					}
				}
				if !isMatchedWithSelectJobInfoLink {
					errMessage = "求人IDを選択ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// 求人選択ページへ遷移
				var selectJobInfoPage *rod.Page
				pages := browser.MustPages()
				for _, page := range pages {
					if strings.Contains(page.MustInfo().URL, "ca_s02070.jsp") {
						selectJobInfoPage = page
						break
					}
				}
				if selectJobInfoPage == nil {
					errMessage = "求人IDを選択ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				selectJobInfoPageBody := selectJobInfoPage.MustElement("body")
				if !selectJobInfoPageBody.MustHas("div#pp_contents") {
					errMessage = "オファー送信対象でない候補者が含まれています。該当候補者を対象から外してください。"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// メインのbodyのdiv要素を取得
				selectJobInfoPageDiv := selectJobInfoPageBody.
					MustElement("div#pp_contents")

				// 求人IDを選択
				if !selectJobInfoPageDiv.MustHas("table#tmpltList") {
					errMessage = "求人IDを選択ページのテーブルが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}
				selectJobInfoTrs := selectJobInfoPageDiv.
					MustElement("table#tmpltList").
					MustElements("tr")
				var isMatchedWithJobInfo bool
			selectJobInfoLoopForSendOther:
				for selectJobInfoI, selectJobInfoTr := range selectJobInfoTrs {
					if selectJobInfoI == 0 {
						continue
					}
					if selectJobInfoTr.MustHas("td") {
						selectJobInfoTds := selectJobInfoTr.MustElements("td")
						for selectJobInfoTdI, selectJobInfoTd := range selectJobInfoTds {
							if selectJobInfoTdI == 1 {
								savedJobInfoIDStr := selectJobInfoTd.MustText()
								log.Println("savedJobInfoID", savedJobInfoIDStr)

								if savedJobInfoIDStr == scoutServiceTemplate.JobInformationID {
									selectJobInfoTds[0].MustElement("a.btn01.btn_s.w80").MustClick()
									time.Sleep(2 * time.Second)
									isMatchedWithJobInfo = true
									break selectJobInfoLoopForSendOther
								}
							}
						}
					}
				}

				if !isMatchedWithJobInfo {
					errMessage = "求人IDが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// 求人IDを反映して閉じる
				selectJobInfoPageDiv.MustElement("a.btn01.btn_m.large.w200").MustClick()
				time.Sleep(10 * time.Second)

				// 一括オファーを作成ボタンをクリック
				createOfferBtn := pageFrame.MustElement("a.btn01.btn_m.bold.large.w200.middle")
				if createOfferBtn == nil || createOfferBtn.MustText() != "一括オファー作成" {
					errMessage = "一括オファーを作成ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				createOfferBtn.MustClick()
				pageFrame.WaitLoad()
				time.Sleep(10 * time.Second)

				// メッセージテンプレートを選択
				messageTemplateSelect := pageFrame.MustElement("select#tmplt_id")
				err = messageTemplateSelect.Select(
					[]string{
						scoutServiceTemplate.MessageTitle,
					},
					true,
					rod.SelectorTypeText,
				)
				if err != nil {
					errMessage = "メッセージテンプレートが見つかりませんでした。DBから取得したタイトル" + scoutServiceTemplate.SearchTitle
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				// テンプレートを反映
				applyTemplateBtn := pageFrame.MustElement("a.btn03.btn_s.w150")
				if applyTemplateBtn == nil || applyTemplateBtn.MustText() != "本文テンプレートを反映" {
					errMessage = "テンプレートを反映ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}
				// クリック後のwindow.confirmをOKにする
				waitConfirm, handleConfirm := pageFrame.MustHandleDialog()
				go applyTemplateBtn.MustClick()
				waitConfirm()
				handleConfirm(true, "")

				time.Sleep(2 * time.Second)

				// 送信内容確認ボタンをクリック
				confirmBtn := pageFrame.MustElement("a.btn01.btn_m.large.w200")
				if confirmBtn == nil || confirmBtn.MustText() != "送信内容確認" {
					errMessage = "送信内容確認ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				confirmBtn.MustClick()
				pageFrame.WaitLoad()
				time.Sleep(10 * time.Second)

				// 送信内容確認画面へ移動
				confirmPage := page.MustElement("frame[name=lowerframe]").MustFrame()

				// 送信ボタンをクリック
				sendBtn := confirmPage.MustElement("a.btn01.btn_m.large.w300")
				log.Println("送信ボタン。sendBtn: ", sendBtn.MustText())
				if sendBtn == nil || sendBtn.MustText() != "送信" {
					errMessage = "送信ボタンが見つかりませんでした"
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				if i.app.Env == "prd" {
					sendBtn.MustClick()
					confirmPage.WaitLoad()
					time.Sleep(10 * time.Second)

				}

				// 次のテンプレートのためにスカウトのトップページへ戻る
				page.MustNavigate("https://ran.next.rikunabi.com/rnc/docs/ca_f00010.jsp?tab_select_id=1&screen_id=ca_s02120&__u=16849106369894720112352085858491")
				page.WaitLoad()
				time.Sleep(10 * time.Second)
				pageFrame = page.MustElement("frame[name=lowerframe]").MustFrame()
			}

			log.Println("スカウトサービスを実行しました")

			// 次のテンプレートのためにスカウトのトップページへ戻る
			page.MustNavigate("https://ran.next.rikunabi.com/rnc/docs/ca_f00010.jsp?__u=16806743487066738118737323900111")
			page.WaitLoad()
			time.Sleep(10 * time.Second)
		} else {
			errMessage = "スカウトタイプが不正です。"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}
	}

	// ログアウト
	// logoutBtn := page.MustElement("a.btn03.btn_s.w80.ml15")
	// if logoutBtn == nil || logoutBtn.MustText() != "ログアウト" {
	// errMessage = "ログアウトボタンが見つかりませんでした"
	// 	log.Println(errMessage)
	// 			return ok, errors.New(errMessage)
	// }

	output.OK = true

	return output, nil
}

/*
マイナビエージェントでスカウトを送信する

開発環境ログ 2つのテンプレートで各々1人ずつ選択
*/
type ScoutOnMynaviScoutingInput struct {
	Context                  context.Context
	ScoutService             *entity.ScoutService
	ScoutServiceTemplateList []*entity.ScoutServiceTemplate
}

type ScoutOnMynaviScoutingOutput struct {
	OK bool
}

func (i *ScoutServiceInteractorImpl) ScoutOnMynaviScouting(input ScoutOnMynaviScoutingInput) (ScoutOnMynaviScoutingOutput, error) {
	var (
		output     ScoutOnMynaviScoutingOutput
		err        error
		errMessage string
		browser    *rod.Browser
		page       *rod.Page
		scoutPage  *rod.Page
		now        = time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	)

	// ホストOSのchromeパスを取得
	path, ok := launcher.LookPath()
	if !ok {
		errMessage = "chrome, chromium or microsoft edgeが見つかりません"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// 取得したパスからchromeを起動
	l := launcher.New().
		Bin(path)

	// ローカル環境の場合は画面ありで起動
	if i.app.Env == "local" {
		l.Headless(false)
	}

	defer time.Sleep(10 * time.Second)
	defer l.Cleanup()

	url := l.MustLaunch()

	// ブラウザを起動後、AMBIへ遷移
	browser = rod.New().
		ControlURL(url).
		Trace(true).                 // デバッグ用のためローカル以外で使わない
		SlowMotion(2 * time.Second). // 機械対策のページ用に遅延を入れる
		MustConnect()

	defer browser.MustClose()

	// browserにタイムアウトを設定
	browser, cancel := browser.
		Context(input.Context).
		WithCancel()
	defer cancel()

	// ブラウザを起動後、AMBIへ遷移
	page = browser.
		// ログインページへ遷移
		MustPage("https://scouting.mynavi.jp/client/").
		MustWaitLoad()
	time.Sleep(4 * time.Second)

	// ログイン情報を入力
	page.
		MustElement("input[name=clientMailAddress]").
		MustInput(input.ScoutService.LoginID)

	// パスワードの複号
	decryptedPassword, err := decryption(input.ScoutService.Password)
	if err != nil {
		log.Println(err)
		return output, err
	}

	page.
		MustElement("input[name=clientPassword]").
		MustInput(decryptedPassword).
		MustType(rodInput.Enter)

	page.WaitLoad()
	time.Sleep(10 * time.Second)

	// マイナビの場合は、ログイン失敗時、再度ログインする
	if page.MustHas("input[name=clientMailAddress]") {
		log.Println("ログイン失敗1回目、再度ログインします")
		page.
			MustElement("input[name=clientPassword]").
			MustType(rodInput.Enter)

		page.WaitLoad()
		time.Sleep(10 * time.Second)
	}

	// マイナビの場合は、ログイン失敗時、再度ログインする
	if page.MustHas("input[name=clientMailAddress]") {
		log.Println("ログイン失敗1回目、再度ログインします")
		page.
			MustElement("input[name=clientPassword]").
			MustType(rodInput.Enter)

		page.WaitLoad()
		time.Sleep(10 * time.Second)
	}

	// 再度ログイン時も失敗した場合は、エラーを返す
	if page.MustHas("input[name=clientMailAddress]") {
		errMessage = "ログイン失敗しました。IDとパスワードが間違っているか。すでに利用しているユーザーがいる可能性があります。"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	log.Println("ログイン成功")

	// ログアウト
	defer page.MustNavigate("https://scouting.mynavi.jp/client/login/logout")

	maxCount := 330
	currentCount := 0
templateLoop:
	for _, scoutServiceTemplate := range input.ScoutServiceTemplateList {
		currentCountPerTemplate := 0
		if currentCount >= maxCount {
			log.Println("指定された最大数を超えたためスカウトを終了します。\ncurrentCount:", currentCount, "\nmaxCount:", maxCount)
			break templateLoop
		}

		// 最終送信が12時間以内の場合はスカウト送信をスキップ
		if scoutServiceTemplate.LastSendAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).After(now.Add(-12 * time.Hour)) {
			log.Println("スカウト送信をスキップします。Title", scoutServiceTemplate.JobInformationTitle, "LastSendAt:", scoutServiceTemplate.LastSendAt)
			continue
		}

		// 保存条件の検索結果を表示する
		page.MustNavigate(
			fmt.Sprintf(
				"https://scouting.mynavi.jp/client/mynaviScoutSearchResult/?useMemSelectId=%s&transition=1",
				scoutServiceTemplate.SearchTitle,
			),
		)
		page.WaitLoad()
		time.Sleep(20 * time.Second)

		// 広告が表示される場合は、閉じる
		// if page.MustHas("button.karte-close") {
		// 	page.MustElement("button.karte-close").MustClick()
		// }

		if !page.MustHas("select[name=leadScout]") {
			log.Println("検索結果がありません。20秒待機")
			time.Sleep(20 * time.Second)
		}
		if !page.MustHas("select[name=leadScout]") {
			log.Println("検索結果がありません。さらに20秒待機")
			time.Sleep(20 * time.Second)
		}

		if !page.MustHas("select[name=leadScout]") {
			log.Println("検索結果がありません")
			err = i.scoutServiceTemplateRepository.UpdateLastSend(
				scoutServiceTemplate.ID,
				0,
				time.Now().UTC(),
			)
			if err != nil {
				log.Println(err)
				return output, err
			}
			continue templateLoop
		}

		// 送信可能通数を計算
		// if scoutServiceTemplateI == 0 {
		// 	// dl.point-list
		// 	pointListDl, err := page.
		// 		Element("dl.point-list")
		// 	if err != nil {
		// 		log.Println(err)
		// 		return output, err
		// 	}

		// 	// 送信数
		// 	pointListDdList, err := pointListDl.
		// 		Elements("dd")
		// 	if err != nil {
		// 		log.Println(err)
		// 		return output, err
		// 	}

		// 	// 送信数を取得
		// 	totalCountThisMonth := 0
		// 	maxCountThisMonth := 0
		// 	for pointListDdI, pointListDd := range pointListDdList {
		// 		if pointListDdI == 0 {
		// 			totalCountThisMonthStr, err := pointListDd.Text()
		// 			if err != nil {
		// 				log.Println(err)
		// 				return output, err
		// 			}
		// 			totalCountThisMonthStr = strings.Replace(totalCountThisMonthStr, "通", "", -1)
		// 			totalCountThisMonth, err = strconv.Atoi(totalCountThisMonthStr)
		// 			if err != nil {
		// 				log.Println(err)
		// 				return output, err
		// 			}
		// 		} else if pointListDdI == 1 {
		// 			maxCountThisMonthStr, err := pointListDd.Text()
		// 			if err != nil {
		// 				log.Println(err)
		// 				return output, err
		// 			}
		// 			maxCountThisMonthStr = strings.Replace(maxCountThisMonthStr, "通", "", -1)
		// 			maxCountThisMonth, err = strconv.Atoi(maxCountThisMonthStr)
		// 			if err != nil {
		// 				log.Println(err)
		// 				return output, err
		// 			}
		// 		}
		// 	}
		// 	totalAvailableCount := maxCountThisMonth - totalCountThisMonth
		// 	if totalAvailableCount < maxCount {
		// 		maxCount = totalAvailableCount
		// 	}
		// 	log.Println("totalAvailableCount:", totalAvailableCount, "\nmaxCount:", maxCount)
		// }

		// 現在の送信数を最大送信数から引いた数に一番近い数を選択する
		if currentCount+50 > maxCount {
			log.Println("50通選択. currentCount:", currentCount, "maxCount:", maxCount)
			scoutServiceTemplate.ScoutCount = null.NewInt(50, true)
		} else if currentCount+100 > maxCount {
			log.Println("100通選択. currentCount:", currentCount, "maxCount:", maxCount)
			scoutServiceTemplate.ScoutCount = null.NewInt(100, true)
		} else if currentCount+300 > maxCount {
			log.Println("300通選択. currentCount:", currentCount, "maxCount:", maxCount)
			scoutServiceTemplate.ScoutCount = null.NewInt(300, true)
		} else if currentCount+500 > maxCount {
			log.Println("500通選択. currentCount:", currentCount, "maxCount:", maxCount)
			scoutServiceTemplate.ScoutCount = null.NewInt(500, true)
		}

		// 通数を選択する
		leadScoutSelect, err := page.
			Element("select[name=leadScout]")
		if err != nil {
			log.Println(err)
			return output, err
		}
		if scoutServiceTemplate.ScoutCount == null.NewInt(50, true) ||
			scoutServiceTemplate.ScoutCount == null.NewInt(100, true) ||
			scoutServiceTemplate.ScoutCount == null.NewInt(300, true) ||
			scoutServiceTemplate.ScoutCount == null.NewInt(500, true) {
			err = leadScoutSelect.
				Select([]string{strconv.Itoa(int(scoutServiceTemplate.ScoutCount.Int64)) + "人"}, true, rod.SelectorTypeText)
			if err != nil {
				log.Println(err)
				return output, err
			}
		} else {
			errMessage = fmt.Sprint("スカウト人数が不正です。DBから取得した数:", scoutServiceTemplate.ScoutCount)
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		time.Sleep(2 * time.Second)

		// 実行ボタンをクリックする
		exeBtn, err := page.
			Element("a.btn.sizeMS.white.btn-blukregistration-action")
		if err != nil {
			log.Println(err)
			return output, err
		}
		exeBtn.MustClick()
		page.WaitLoad()
		time.Sleep(10 * time.Second)

		// ブラウザで開いているページを取得する
		pages := browser.MustPages()

		// スカウト送信ページへ移動
		for _, page := range pages {
			log.Println(page.MustInfo().URL)
			if page.MustInfo().URL == "https://scouting.mynavi.jp/client/mynaviScoutMail/leadSend" {
				scoutPage = page
				break
			}
		}
		if scoutPage == nil {
			errMessage = "スカウト送信ページが見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		if !scoutPage.MustHas("select.select-scout-templete") {
			time.Sleep(5 * time.Second)
		}

		if !scoutPage.MustHas("select.select-scout-templete") {
			log.Println("送信対象者がいません")
			scoutPage.MustClose()
			err = i.scoutServiceTemplateRepository.UpdateLastSend(
				scoutServiceTemplate.ID,
				0,
				time.Now().UTC(),
			)
			if err != nil {
				log.Println(err)
				return output, err
			}
			continue templateLoop
		}

		// 保存したテンプレートを選択する
		savedTemplateSelect, err := scoutPage.
			Element("select.select-scout-templete")
		if err != nil {
			log.Println(err)
			return output, err
		}

		err = savedTemplateSelect.Select(
			[]string{scoutServiceTemplate.MessageTitle}, true, rod.SelectorTypeText)
		if err != nil {
			errMessage = "テンプレートが見つかりませんでした" + scoutServiceTemplate.MessageTitle
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		// テンプレートを適用する
		setTemplateBtn, err := scoutPage.
			Element("a.set-templete")
		if err != nil {
			log.Println(err)
			return output, err
		}

		setTemplateBtn.MustClick()
		scoutPage.WaitLoad()
		time.Sleep(10 * time.Second)

		// 年齢チェック
		jobSeekerInfos, err := scoutPage.
			Elements("p.acdp")
		if err != nil {
			log.Println(err)
			return output, err
		}

		if len(jobSeekerInfos) == 0 {
			// p.first-child
			jobSeekerInfo, err := scoutPage.
				Element("p.first-child")
			if err != nil {
				log.Println(err)
				return output, err
			}

			jobSeekerInfoText, err := jobSeekerInfo.Text()
			if err != nil {
				log.Println(err)
				return output, err
			}
			log.Println(jobSeekerInfoText)

			// （ 23歳 / 女性 / 愛知県春日井市藤山台 ）から年齢を取得
			ageStr := strings.Split(jobSeekerInfoText, " ")[1]
			ageStr = strings.Replace(ageStr, "歳", "", -1)
			age, err := strconv.Atoi(ageStr)
			if err != nil {
				log.Println(err)
				return output, err
			}

			// 年齢制限を超えている場合は、スカウトを送信しない
			if int64(age) >= scoutServiceTemplate.AgeLimit.Int64 {
				log.Println("年齢制限を超えているため、スカウトを送信しません", age)
				scoutPage.MustClose()
				continue templateLoop
			}
		}

		for jobSeekerInfoI, jobSeekerInfo := range jobSeekerInfos {
			// 年齢が超えている場合 or 最大送信数を超えた場合は、チェックボックスを外す
			isExclude := false
			if jobSeekerInfoI == 50 {
				// 50人以上の場合は、全て表示ボタンを押す
				allBtn, err := scoutPage.
					Element("div.acd-all-btn")
				if err != nil {
					log.Println(err)
					return output, err
				}

				allBtn.MustClick()
				time.Sleep(15 * time.Second)
			}
			if currentCount > maxCount {
				log.Println("最大送信数を超えました", currentCount)
				isExclude = true
			} else {
				infoSpans, err := jobSeekerInfo.
					Elements("span")
				if err != nil {
					log.Println(err)
					return output, err
				}

				if len(infoSpans) < 2 {
					log.Println("infoSpansが2未満です")
					isExclude = true
				} else {

					ageStr := infoSpans[1].MustText()
					age, err := strconv.Atoi(ageStr)
					if err != nil {
						log.Println("年齢情報が見つかりませんでした。ageStr:", ageStr)
						isExclude = true
					}

					// 年齢制限を超えている場合は、チェックボックスを外す
					if int64(age) >= scoutServiceTemplate.AgeLimit.Int64 {
						isExclude = true
					}
				}
			}

			if isExclude {
				log.Println("年齢制限を超えているため、スカウトを送信しません。")
				sendMemberIdInput, err := jobSeekerInfo.
					Element("input[name=sendMemberId]")
				if err != nil {
					log.Println(err)
					return output, err
				}

				sendMemberIdInput.MustClick()
				continue
			}

			// 現在の送信数をカウント
			currentCount++
			currentCountPerTemplate++
		}
		log.Println("currentCountPerTemplate:", currentCountPerTemplate, "currentCount:", currentCount)

		if currentCountPerTemplate == 0 {
			if scoutPage.MustHas("div.user-head.cf") {
				currentCountPerTemplate = 1
			} else {
				log.Println("送信対象者がいません")
				scoutPage.MustClose()
				err = i.scoutServiceTemplateRepository.UpdateLastSend(
					scoutServiceTemplate.ID,
					0,
					time.Now().UTC(),
				)
				if err != nil {
					log.Println(err)
					return output, err
				}
				continue templateLoop
			}
		}

		time.Sleep(2 * time.Second)

		// 確認
		commonSubmitBtn, err := scoutPage.
			Element("a.js_commonSubmit")
		if err != nil {
			log.Println(err)
			return output, err
		}

		commonSubmitBtn.MustClick()
		scoutPage.WaitLoad()
		time.Sleep(10 * time.Second)

		log.Println("スカウト確認ボタンをクリックしました")

		// 送信
		if scoutPage.MustHas("a.js-combined.disabled") {
			log.Println("送信ボタンが無効です")
			break templateLoop
		}

		if scoutPage.MustHas("a.js-combined") {
			// btn sizeMM blue inline js-btn js-indivisual
			sendBtn, err := scoutPage.
				Element("a.js-combined")
			if err != nil {
				log.Println(err)
				return output, err
			}
			log.Println("送信ボタン。sendBtn: ", sendBtn.MustText())

			// if i.app.BatchType == "scout" {
			sendBtn.MustClick()
			time.Sleep(40 * time.Second)
			// }
		} else {

			time.Sleep(5 * time.Second)

			if !scoutPage.MustHas("a.js-indivisual") {
				continue templateLoop
			}

			sendBtn, err := scoutPage.
				Element("a.js-indivisual")
			if err != nil {
				log.Println(err)
				return output, err
			}
			log.Println("送信ボタン。sendBtn: ", sendBtn.MustText())

			// if i.app.BatchType == "scout" {
			sendBtn.MustClick()

			time.Sleep((time.Duration(currentCountPerTemplate/5) + 5) * time.Second)
			// }
		}

		// ページを閉じる
		scoutPage.MustClose()
		log.Println("スカウト送信ボタンをクリックしました")

		// スカウト送信完了後、スカウト情報を更新
		err = i.scoutServiceTemplateRepository.UpdateLastSend(
			scoutServiceTemplate.ID,
			uint(currentCountPerTemplate),
			time.Now().UTC(),
		)
		if err != nil {
			log.Println(err)
			return output, err
		}
		continue templateLoop
	}

	// 成功メールを送信
	err = i.sendScoutSuccessMail(input.ScoutService.AgentStaffID, input.ScoutService, input.ScoutServiceTemplateList)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.OK = true
	return output, err
}

/*
AMBIでスカウトを送信する

開発環境ログ 2つのテンプレートで各々1人ずつ選択
[rod] 2023/05/10 17:01:01 [wait] load <page:B2A9206C>
2023/05/10 17:07:53 AMBIのスカウト送信が完了しました true
*/
type ScoutOnAmbiInput struct {
	Context                  context.Context
	ScoutService             *entity.ScoutService
	ScoutServiceTemplateList []*entity.ScoutServiceTemplate
}

type ScoutOnAmbiOutput struct {
	OK bool
}

func (i *ScoutServiceInteractorImpl) ScoutOnAmbi(input ScoutOnAmbiInput) (ScoutOnAmbiOutput, error) {
	var (
		output     ScoutOnAmbiOutput
		err        error
		errMessage string
		browser    *rod.Browser
		page       *rod.Page
		now        = time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	)

	// 全ての最終送信日時が4日以前の場合は、スカウトの合計送信数をリセット
	runnedCount := 0
	for _, scoutServiceTemplate := range input.ScoutServiceTemplateList {
		if strings.Contains(scoutServiceTemplate.SearchTitle, "数合わせ") {
			continue
		}
		if scoutServiceTemplate.LastSendAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Before(now.Add(-4 * 24 * time.Hour)) {
			runnedCount++
		}
	}

	if runnedCount == len(input.ScoutServiceTemplateList) {
		input.ScoutService.LastSendCount = null.NewInt(0, true)
		err = i.scoutServiceRepository.UpdateLastSendCount(
			input.ScoutService.ID,
			0,
		)
		if err != nil {
			log.Println(err)
			return output, err
		}
	}

	// ホストOSのchromeパスを取得
	path, ok := launcher.LookPath()
	if !ok {
		errMessage = "chrome, chromium or microsoft edgeが見つかりません"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// 取得したパスからchromeを起動
	l := launcher.New().
		Bin(path)

	// ローカル環境の場合は画面ありで起動
	if i.app.Env == "local" {
		l.Headless(false)
	}

	url := l.MustLaunch()

	// ブラウザを起動後、AMBIへ遷移
	browser = rod.New().
		ControlURL(url).
		Trace(true).                 // ログ出力
		SlowMotion(2 * time.Second). // 機械対策のページ用に遅延を入れる
		MustConnect()

	// browserにタイムアウトを設定
	browser, cancel := browser.
		Context(input.Context).
		WithCancel()
	defer cancel()

	page = browser.
		// ログインページへ遷移
		MustPage("https://en-ambi.com/company_login/login/").
		MustWaitLoad()

	time.Sleep(10 * time.Second)

	// ログイン情報を入力
	loginIDEl, err := page.
		Element("input[name=accLoginID]")
	if err != nil {
		errMessage = "ログインID入力要素が見つかりませんでした"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	loginIDEl.
		MustInput(input.ScoutService.LoginID)

	// パスワードの複号
	decryptedPassword, err := decryption(input.ScoutService.Password)
	if err != nil {
		log.Println(err)
		return output, err
	}

	loginPWEl, err := page.
		Element("input[name=accLoginPW]")
	if err != nil {
		errMessage = "パスワード入力要素が見つかりませんでした"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	loginPWEl.
		MustInput(decryptedPassword).
		MustType(rodInput.Enter)

	page.WaitLoad()
	time.Sleep(10 * time.Second)

	// AMBIでスカウト送信
templateLoop:
	for _, scoutServiceTemplate := range input.ScoutServiceTemplateList {

		// 最終送信が12時間以内の場合はスカウト送信をスキップ
		if scoutServiceTemplate.LastSendAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).After(now.Add(-12 * time.Hour)) {
			log.Println("スカウト送信をスキップします。Title", scoutServiceTemplate.SearchTitle, "LastSendAt:", scoutServiceTemplate.LastSendAt)
			continue
		}

		// 数合わせ用の場合は、1000 - scoutService.LastSendCount分のスカウトを送信する
		if strings.Contains(scoutServiceTemplate.SearchTitle, "数合わせ") {
			if 1000 <= input.ScoutService.LastSendCount.Int64 {
				log.Println("スカウト送信数が1000件を超えたため終了します。")
				break
			}

			balanceCount := 1000 - input.ScoutService.LastSendCount.Int64
			// 200を超える場合は、200にする
			if balanceCount > 200 {
				balanceCount = 200
			}
			scoutServiceTemplate.ScoutCount = null.NewInt(balanceCount, true)
			log.Println("数合わせ用のスカウト送信数を設定しました. balanceCount:", balanceCount, "input.ScoutService.LastSendCount", input.ScoutService.LastSendCount)
			if scoutServiceTemplate.ScoutCount.Int64 == 0 {
				log.Println("スカウト送信数が0のため、スカウト送信をスキップします")
				continue
			}
		}

		// 保存条件ページへ遷移
		page.MustNavigate("https://en-ambi.com/company/scout/condition_list/?PK=9ABBAF")
		page.WaitLoad()
		time.Sleep(10 * time.Second)

		// URLがログインページのままの場合はログインに失敗していると判断
		if strings.Contains(page.MustInfo().URL, "login") {
			errMessage = "ログインに失敗しました。"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		/*
			保存条件の検索結果を表示する
		*/
		searchList, err := page.
			Elements("table.md_tableList.md_tableHoverList > tbody > tr")
		if err != nil {
			log.Println(err)
			return output, err
		}

		isMatchedWithSearchTitle := false
		for _, search := range searchList {
			searchTitle, err := search.
				Element("td.data.name")
			if err != nil {
				errMessage = "保存条件が見つかりませんでした" + scoutServiceTemplate.SearchTitle
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			if searchTitle.MustText() == scoutServiceTemplate.SearchTitle {
				savedFolder, err := search.
					Element("a.md_btn.md_btn--green")
				if err != nil {
					errMessage = "保存条件が見つかりませんでした" + scoutServiceTemplate.SearchTitle
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}
				savedFolder.MustClick()

				page.WaitLoad()
				time.Sleep(10 * time.Second)
				isMatchedWithSearchTitle = true
				break
			}
		}
		if !isMatchedWithSearchTitle {
			errMessage = "保存条件が見つかりませんでした" + scoutServiceTemplate.SearchTitle
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		// スカウト件数を選択
		if !page.MustHas("select[name=pageLimit]") {
			time.Sleep(15 * time.Second)
			if !page.MustHas("select[name=pageLimit]") {
				log.Println("条件検索の結果が0件でした")
				// スカウト送信完了後、スカウト情報を更新
				err = i.scoutServiceTemplateRepository.UpdateLastSend(
					scoutServiceTemplate.ID,
					0,
					time.Now().UTC(),
				)
				if err != nil {
					log.Println(err)
					return output, err
				}

				continue templateLoop
			}
		}

		pageLimitSelectForSearch, err := page.
			Element("select[name=pageLimit]")
		if err != nil {
			errMessage = "条件検索のスカウト件数選択要素が見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}
		if scoutServiceTemplate.ScoutCount.Int64 <= 50 {
			log.Println("50件表示")
			pageLimitSelectForSearch.MustSelect("50件表示")
			page.WaitLoad()
			time.Sleep(10 * time.Second)
		} else if scoutServiceTemplate.ScoutCount.Int64 > 50 {
			log.Println("200件表示")
			pageLimitSelectForSearch.MustSelect("200件表示")
			page.WaitLoad()
			time.Sleep(10 * time.Second)
		} else {
			errMessage = "AMBIのスカウト件数が不正です。50件または200件で入力してください" + strconv.FormatInt(scoutServiceTemplate.ScoutCount.Int64, 10)
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		/*
			検索条件の該当者を検討人材リストへ追加する
		*/
		// 検討人材リストに追加した人数をカウント
		var savedUserCnt int64
	pageLoop:
		for pageI := 1; pageI <= 5; pageI++ {

			if savedUserCnt >= scoutServiceTemplate.ScoutCount.Int64 {
				log.Println("検討人材リストに追加した人数が上限に達したため、追加を終了します")
				time.Sleep(5 * time.Second)
				break pageLoop
			}

			time.Sleep(5 * time.Second)

			users, err := page.
				Elements("div.userSet")
			if err != nil {
				log.Println(err)
				return output, err
			}

			for _, user := range users {
				log.Println("user:", user)

				// 36歳以上は省く
				profDiv, err := user.
					Element("div.prof")
				if err != nil {
					log.Println(err)
					return output, err
				}

				profText := profDiv.MustText()

				/*
					年齢チェック
				*/
				agenAndPrefecture := strings.Split(profText, "性\n")[1]
				ageStr := strings.Split(agenAndPrefecture, "歳")[0]
				age, err := strconv.Atoi(ageStr)
				if err != nil {
					log.Println(err)
					return output, err
				}
				if int64(age) >= scoutServiceTemplate.AgeLimit.Int64 {
					continue
				}

				// 動作確認のため、2人まで
				// if i.app.BatchType != "scout" && savedUserCnt >= 2 {
				if savedUserCnt >= 2 {
					break pageLoop
				}

				// スカウトボタン md_btn md_btn--green js_slideOpen js_abled_duplication_scout
				has, _, err := user.
					Has("a.md_btn.md_btn--green")
				if err != nil {
					log.Println(err)
					return output, err
				}
				if !has {
					continue
				}

				scoutBtn, err := user.
					Element("a.md_btn--green")
				if err != nil {
					log.Println(err)
					return output, err
				}
				if scoutBtn == nil {
					continue
				}
				scoutBtnText := scoutBtn.MustText()
				// 再スカウトかどうかで処理を分ける
				if scoutBtnText == "再スカウト" &&
					(scoutServiceTemplate.ScoutType == null.NewInt(entity.AmbiScoutTypeNormalAndAgain, true) ||
						scoutServiceTemplate.ScoutType == null.NewInt(entity.AmbiScoutTypePremiumAndAgain, true)) {
					// 検討人材リストに追加
					if user.MustHas("a.md_btn.md_btn--white2.js_consider.js_onlist") {
						considerBtn, err := user.
							Element("a.md_btn.md_btn--white2.js_consider.js_onlist")
						if err != nil {
							errMessage = "検討ボタンが見つかりませんでした"
							return output, errors.New(errMessage)
						}

						considerBtn.MustClick()
						time.Sleep(1 * time.Second)

						savedUserCnt++
						scoutServiceTemplate.LastSendCount = null.NewInt(savedUserCnt, true)
					}
				} else if scoutBtnText == "スカウト" &&
					(scoutServiceTemplate.ScoutType == null.NewInt(entity.AmbiScoutTypeNormal, true) ||
						scoutServiceTemplate.ScoutType == null.NewInt(entity.AmbiScoutTypePremium, true)) {
					// 検討人材リストに追加
					if user.MustHas("a.md_btn.md_btn--white2.js_consider.js_onlist") {
						considerBtn, err := user.
							Element("a.md_btn.md_btn--white2.js_consider.js_onlist")
						if err != nil {
							log.Println(err)
							return output, err
						}

						considerBtn.MustClick()
						time.Sleep(1 * time.Second)

						savedUserCnt++
						scoutServiceTemplate.LastSendCount = null.NewInt(savedUserCnt, true)
					}
				}

				// 次のユーザーへ
				continue

			}
			// 次のページへ
			if page.MustHas("li.list.next") {
				page.MustElement("li.list.next").MustElement("a").MustClick()
				page.WaitLoad()
				time.Sleep(10 * time.Second)
			} else {
				break pageLoop
			}
		}

		/*
			保存条件リストページへ遷移
		*/
		page.MustNavigate("https://en-ambi.com/company/scout/condition_list/?PK=6BF948")
		page.WaitLoad()
		time.Sleep(10 * time.Second)

		searchList, err = page.
			Elements("table.md_tableList.md_tableHoverList > tbody > tr")
		if err != nil {
			log.Println(err)
			return output, err
		}

		isMatchedWithSearchTitleForScout := false
		for _, search := range searchList {
			searchTitle, err := search.
				Element("td.data.name")
			if err != nil {
				errMessage = "保存条件が見つかりませんでした" + scoutServiceTemplate.SearchTitle
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			if searchTitle.MustText() == scoutServiceTemplate.SearchTitle {
				savedFolder, err := search.
					Element("td.data.folder")
				if err != nil {
					errMessage = "保存条件が見つかりませんでした" + scoutServiceTemplate.SearchTitle
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}
				savedFolderA, err := savedFolder.
					Element("a")
				if err != nil {
					errMessage = "保存条件が見つかりませんでした" + scoutServiceTemplate.SearchTitle
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}

				savedFolderA.
					MustClick()

				page.WaitLoad()
				time.Sleep(10 * time.Second)
				isMatchedWithSearchTitleForScout = true
				break
			}
		}
		if !isMatchedWithSearchTitleForScout {
			errMessage = "保存条件が見つかりませんでした" + scoutServiceTemplate.SearchTitle
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		/*
			検討人材リストを表示する
		*/
		// 再度スカウトの場合スカウト送信済みで検索
		if scoutServiceTemplate.ScoutType == null.NewInt(entity.AmbiScoutTypeNormalAndAgain, true) ||
			scoutServiceTemplate.ScoutType == null.NewInt(entity.AmbiScoutTypePremiumAndAgain, true) {
			labelEl, err := page.
				Element("label[for=status_003]")
			if err != nil {
				errMessage = "スカウト送信済みのラベルが見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			labelEl.MustClick()
			time.Sleep(4 * time.Second)
			searchSubmitBtn, err := page.
				Element("a.md_btn.md_btn--min.search.js_submit")
			if err != nil {
				errMessage = "検索ボタンが見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			searchSubmitBtn.MustClick()
			page.WaitLoad()
			time.Sleep(10 * time.Second)
		}

		// スカウト件数を選択
		if !page.MustHas("select[name=pageLimit]") {
			time.Sleep(15 * time.Second)
			if !page.MustHas("select[name=pageLimit]") {
				log.Println("保存済みの検討人材リストが0件でした")
				// スカウト送信完了後、スカウト情報を更新
				err = i.scoutServiceTemplateRepository.UpdateLastSend(
					scoutServiceTemplate.ID,
					0,
					time.Now().UTC(),
				)
				if err != nil {
					log.Println(err)
					return output, err
				}

				continue templateLoop
			}
		}
		pageLimitSelect, err := page.
			Element("select[name=pageLimit]")
		if err != nil {
			errMessage = "スカウト件数選択要素が見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		if scoutServiceTemplate.ScoutCount.Int64 <= 50 {
			log.Println("50件表示")
			pageLimitSelect.MustSelect("50件表示")
			page.WaitLoad()
			time.Sleep(10 * time.Second)
		} else if scoutServiceTemplate.ScoutCount.Int64 > 50 {
			log.Println("200件表示")
			pageLimitSelect.MustSelect("200件表示")
			page.WaitLoad()
			time.Sleep(10 * time.Second)
		} else {
			errMessage = "AMBIのスカウト件数が不正です。50件または200件で入力してください" + strconv.FormatInt(scoutServiceTemplate.ScoutCount.Int64, 10)
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		/*
			全て選択
		*/
		amuoutCheckAllLabel, err := page.
			Element("label[for=amountCheck_all]")
		if err != nil {
			errMessage = "全て選択のラベルが見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		amuoutCheckAllLabel.MustClick()
		time.Sleep(5 * time.Second)

		if !page.MustHas("a.md_btn.md_btn--min.js_slideOpen.md_btn--green") {
			log.Println("検討人材リストが空のため、スカウトを送信できません。")

			// スカウト送信完了後、スカウト情報を更新
			err = i.scoutServiceTemplateRepository.UpdateLastSend(
				scoutServiceTemplate.ID,
				0,
				time.Now().UTC(),
			)
			if err != nil {
				log.Println(err)
				return output, err
			}
			continue templateLoop
		}

		// md_btn--gray

		slideOpenBtn, err := page.
			Element("a.md_btn.md_btn--min.js_slideOpen.md_btn--green")
		if err != nil {
			errMessage = "スカウト送信ボタンが見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		slideOpenBtn.MustClick()
		time.Sleep(5 * time.Second)

		/*
			スカウト件数を取得
		*/
		copyNumDiv, err := page.
			Element("div.copyNum")
		if err != nil {
			errMessage = "スカウト件数が見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		numEms, err := copyNumDiv.Elements("em.num")
		if err != nil {
			errMessage = "スカウト件数が見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}
		if len(numEms) < 2 {
			errMessage = "スカウト件数が見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		// スカウト件数を取得
		scoutCountText, err := numEms[1].Text()
		if err != nil {
			errMessage = "スカウト件数が見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		// 10通　のような文字列から数字のみを取得
		scoutCountText = strings.Replace(scoutCountText, "通", "", -1)
		scoutCount, err := strconv.ParseInt(scoutCountText, 10, 64)
		if err != nil {
			errMessage = "スカウト件数が見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}
		log.Println("スカウト件数:", scoutCount)
		scoutServiceTemplate.LastSendCount = null.NewInt(scoutCount, true)
		input.ScoutService.LastSendCount.Int64 += scoutCount

		/*
			スカウトタイプの選択
		*/
		if scoutServiceTemplate.ScoutType == null.NewInt(0, true) {
			scoutTypePrivateBtn, err := page.
				Element("label[for=scoutType-01-private]")
			if err != nil {
				errMessage = "スカウトタイプが見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			// scoutTypePrivateBtn, err := page.
			// 	Element("label[for=scoutType-01-platinum]")
			// if err != nil {
			// 	errMessage = "スカウトタイプが見つかりませんでした"
			// 	log.Println(errMessage)
			// 	return output, errors.New(errMessage)
			// }

			scoutTypePrivateBtn.MustClick()
			page.WaitLoad()
			time.Sleep(5 * time.Second)
		} else if scoutServiceTemplate.ScoutType == null.NewInt(1, true) {
			scoutTypePlatinumBtn, err := page.
				Element("label[for=scoutType-01-platinum]")
			if err != nil {
				errMessage = "スカウトタイプが見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			scoutTypePlatinumBtn.MustClick()
			page.WaitLoad()
			time.Sleep(5 * time.Second)
		} else {
			errMessage = "スカウトタイプが不正です。入力されたスカウトタイプ:" + strconv.FormatInt(scoutServiceTemplate.ScoutType.Int64, 10)
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		/*
			テンプレートを選択
		*/
		if scoutServiceTemplate.MessageTitle != "" {
			templateSelect, err := page.Element("select[name=TemplateID]")
			if err != nil {
				errMessage = "テンプレート選択要素が見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			err = templateSelect.Select(
				[]string{scoutServiceTemplate.MessageTitle}, true, rod.SelectorTypeText)
			if err != nil {
				errMessage = "テンプレートが見つかりませんでした" + scoutServiceTemplate.MessageTitle
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
		}

		/*
			返信期限は最短日を選択する
		*/
		replyDeadline, err := page.
			Element("input[name=ReplyDeadline]")
		if err != nil {
			errMessage = "返信期限が見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}
		replyDeadline.MustInput("")

		// ui-datepicker-div
		datePickerDiv, err := page.
			Element("div#ui-datepicker-div")
		if err != nil {
			errMessage = "返信期限のカレンダーが見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		// 前月選択ボタンがdisabledでない場合はクリック
		if !datePickerDiv.MustHas("a.ui-datepicker-prev.ui-corner-all.ui-state-disabled") {
			prevBtn, err := datePickerDiv.Element("a.ui-datepicker-prev")
			if err != nil {
				errMessage = "返信期限のカレンダーの前月ボタンが見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			prevBtn.MustClick()
			datePickerDiv, err = page.
				Element("div#ui-datepicker-div")
			if err != nil {
				errMessage = "返信期限のカレンダーの日付選択が見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// 最短を取得 td.ui-datepicker-days-cell-over
			daysCellOver, err := datePickerDiv.
				Element("td.ui-datepicker-days-cell-over")

			if err != nil {
				errMessage = "返信期限のカレンダーの最短日が見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			selectDay, err := daysCellOver.
				Element("a.ui-state-default")
			if err != nil {
				errMessage = "返信期限のカレンダーの最短日が見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
			selectDay.MustClick()

		} else {
			selectDay, err := datePickerDiv.
				Element("a.ui-state-default")
			if err != nil {
				errMessage = "返信期限のカレンダーの最短日が見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			selectDay.MustClick()
		}

		time.Sleep(5 * time.Second)

		/*
			スカウト送信
		*/
		if page.MustHas("a.md_btn.md_btn--green.welcomeDone.js_tipOpen.js_submitBtn.js_disabled.md_btn--disabled") {
			errMessage = "スカウト送信ボタンが無効です。通数が上限に達しているまたはメッセージテンプレートが設定されていない可能性があります"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		sendBtn, err := page.
			Element("a.md_btn--green.welcomeDone.js_tipOpen.js_submitBtn")
		if err != nil {
			errMessage = "スカウト送信ボタンが見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}
		log.Println("送信ボタン。sendBtn: ", sendBtn.MustText())

		// if i.app.BatchType == "scout" {
		sendBtn.MustClick()
		time.Sleep(5 * time.Second)

		// confirmTip
		confirmTip, err := page.
			Element("div.confirmTip")
		if err != nil {
			errMessage = "スカウト送信確認が見つかりませんでした"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		confirmTip.MustElement("a.md_btn.md_btn--min.js_cellChange.js_loadingIconFlg").MustClick()
		page.WaitLoad()
		// 送信数/2+20秒待つ
		time.Sleep(time.Duration(scoutServiceTemplate.LastSendCount.Int64/2+20) * time.Second)
		// }

		/*
			スカウト送信完了後、スカウト情報を更新
		*/
		err = i.scoutServiceTemplateRepository.UpdateLastSend(
			scoutServiceTemplate.ID,
			uint(scoutServiceTemplate.LastSendCount.Int64),
			time.Now().UTC(),
		)
		if err != nil {
			log.Println(err)
			return output, err
		}

		err = i.scoutServiceRepository.UpdateLastSendCount(
			input.ScoutService.ID,
			uint(input.ScoutService.LastSendCount.Int64),
		)
		if err != nil {
			log.Println(err)
			return output, err
		}

		// 次のテンプレートへ
	}

	// 成功メールを送信
	err = i.sendScoutSuccessMail(input.ScoutService.AgentStaffID, input.ScoutService, input.ScoutServiceTemplateList)
	if err != nil {
		log.Println(err)
		return output, err
	}

	browser.MustClose()
	l.Cleanup()

	// ログアウト
	log.Println("スカウト完了")
	output.OK = true
	return output, nil
}

/*
マイナビエージェントスカウトでスカウトを送信する

開発環境ログ 2つのテンプレートで各々1人ずつ選択
*/
type ScoutOnMynaviAgentScoutInput struct {
	Context                  context.Context
	ScoutService             *entity.ScoutService
	ScoutServiceTemplateList []*entity.ScoutServiceTemplate
}

type ScoutOnMynaviAgentScoutOutput struct {
	OK bool
}

func (i *ScoutServiceInteractorImpl) ScoutOnMynaviAgentScout(input ScoutOnMynaviAgentScoutInput) (ScoutOnMynaviAgentScoutOutput, error) {
	var (
		output     ScoutOnMynaviAgentScoutOutput
		err        error
		errMessage string
		browser    *rod.Browser
		page       *rod.Page
		now        = time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	)

	// スカウトの合計送信数をリセット
	for _, scoutServiceTemplate := range input.ScoutServiceTemplateList {
		if scoutServiceTemplate.LastSendAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Before(now.Add(-12 * time.Hour)) {
			scoutServiceTemplate.LastSendCount = null.NewInt(0, true)

			err = i.scoutServiceTemplateRepository.UpdateLastSend(
				scoutServiceTemplate.ID,
				0,
				scoutServiceTemplate.LastSendAt,
			)
			if err != nil {
				log.Println(err)
				return output, err
			}
		}
	}

	// ホストOSのchromeパスを取得
	path, ok := launcher.LookPath()
	if !ok {
		errMessage = "chrome, chromium or microsoft edgeが見つかりません"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// 取得したパスからchromeを起動
	l := launcher.New().
		Bin(path)

	// ローカル環境の場合は画面ありで起動
	if i.app.Env == "local" {
		l.Headless(false)
	}

	defer time.Sleep(10 * time.Second)

	url := l.MustLaunch()

	// ブラウザを起動後、AMBIへ遷移
	browser = rod.New().
		ControlURL(url).
		Trace(true).                 // デバッグ用のためローカル以外で使わない
		SlowMotion(2 * time.Second). // 機械対策のページ用に遅延を入れる
		MustConnect()

	// browserにタイムアウトを設定
	browser, cancel := browser.
		Context(input.Context).
		WithCancel()
	defer cancel()

	// ブラウザを起動後、AMBIへ遷移
	// 認証情報を設定
	go browser.MustHandleAuth("mynavi-ag-scout", "aJJJ83AsdYULPQTe")()

	page = browser.
		// ログインページへ遷移
		MustPage("https://scout.mynavi-agent.jp/login")

	time.Sleep(4 * time.Second)

	// ログイン情報を入力
	// 企業ID
	companyIDInput, err := page.Element("input#companyid")
	if err != nil {
		log.Println(err)
		return output, err
	}

	companyIDInput.MustInput(input.ScoutService.LoginID)

	// メールアドレス
	emailInput, err := page.Element("input#email")
	if err != nil {
		log.Println(err)
		return output, err
	}

	emailInput.MustInput("info@spaceai.jp")

	// パスワードの複号
	decryptedPassword, err := decryption(input.ScoutService.Password)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// パスワード
	passwordInput, err := page.Element("input#password")
	if err != nil {
		log.Println(err)
		return output, err
	}

	passwordInput.MustInput(decryptedPassword)

	// ログインボタンをクリック
	loginBtn, err := page.Element("button")
	if err != nil {
		log.Println(err)
		return output, err
	}
	if loginBtn.MustText() != "ログイン" {
		errMessage = "ログインボタンが見つかりませんでした"
		return output, errors.New(errMessage)
	}
	loginBtn.MustClick()

	page.WaitLoad()
	time.Sleep(40 * time.Second)

	// 再度ログイン時も失敗した場合は、エラーを返す
	if page.MustHas("input#email") {
		time.Sleep(30 * time.Second)
	}
	if page.MustHas("input#email") {
		errMessage = "ログイン失敗しました。IDとパスワードが間違っているか。すでに利用しているユーザーがいる可能性があります。"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	log.Println("ログイン成功")

	/*
		ログアウト
	*/
	defer func() {
		page.
			MustNavigate("https://scout.mynavi-agent.jp/master/scout/search/conditions")
		page.WaitLoad()
		time.Sleep(10 * time.Second)

		// 設定を開く
		settingNavDiv, err := page.
			Element("div.nav.isSetting")
		if err != nil {
			log.Println(err)
			return
		}

		openSettingNavBtn, err := settingNavDiv.
			Element("a")
		if err != nil {
			log.Println(err)
			return
		}

		openSettingNavBtn.MustClick()
		time.Sleep(2 * time.Second)

		// subNavButton subNavFunctionButton
		logoutBtns, err := page.
			Elements("button.subNavButton.subNavFunctionButton")
		if err != nil {
			log.Println(err)
			return
		}

		isMatchedWithLogoutBtn := false
		for _, logoutBtn := range logoutBtns {
			if logoutBtn.MustText() == "ログアウト" {
				isMatchedWithLogoutBtn = true
				logoutBtn.MustClick()
				time.Sleep(10 * time.Second)
			}
		}
		if !isMatchedWithLogoutBtn {
			errMessage = "ログアウトボタンが見つかりませんでした"
			log.Println(errMessage)
			return
		}

		// generalButton small blue fz-medium
		confirmLogoutBtns, err := page.
			Elements("button.generalButton")
		if err != nil {
			log.Println(err)
			return
		}

		isMatchedWithConfirmLogoutBtn := false
		for _, confirmLogoutBtn := range confirmLogoutBtns {
			if confirmLogoutBtn.MustText() == "ログアウトする" {
				isMatchedWithConfirmLogoutBtn = true
				confirmLogoutBtn.MustClick()
				time.Sleep(10 * time.Second)
			}
		}
		if !isMatchedWithConfirmLogoutBtn {
			errMessage = "ログアウト確認ボタンが見つかりませんでした"
			log.Println(errMessage)
			return
		}

		browser.MustClose()
		l.Cleanup()
	}()

	// 別テンプレートに移動
	for _, scoutServiceTemplate := range input.ScoutServiceTemplateList {
		// 同じテンプレートを500件ずつ複数回送信するため、ページを再度開く
	templateLoop:
		for templateI := 1; templateI <= 30; templateI++ {

			// 送信数を超えた場合は終了
			if scoutServiceTemplate.ScoutCount.Int64 <= scoutServiceTemplate.LastSendCount.Int64 {
				break templateLoop
			}

			// スカウトページへ遷移
			page.
				MustNavigate("https://scout.mynavi-agent.jp/master/scout/search/conditions")
			page.WaitLoad()

			time.Sleep(20 * time.Second)

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			if !page.MustHas("div.filterCard") {
				time.Sleep(30 * time.Second)
				// errorページが表示されているか確認
				if strings.Contains(page.MustInfo().URL, "error") {
					log.Println("エラーページが表示されました")
					continue templateLoop
				}
			}

			// 保存した検索条件を選択する
			filterCardList, err := page.
				Elements("div.filterCard")
			if err != nil {
				log.Println(err)
				return output, err
			}

			if len(filterCardList) == 0 {
				time.Sleep(30 * time.Second)
				// errorページが表示されているか確認
				if strings.Contains(page.MustInfo().URL, "error") {
					log.Println("エラーページが表示されました")
					continue templateLoop
				}

				filterCardList, err = page.
					Elements("div.filterCard")
				if err != nil {
					log.Println(err)
					return output, err
				}
				if len(filterCardList) == 0 {
					errMessage = "保存した検索条件が見つかりませんでした.scoutServiceTemplate.SearchTitle:" + scoutServiceTemplate.SearchTitle
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}
			}

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			isMatchedWithSearchTitle := false
			for _, filterCard := range filterCardList {
				filterCardTitleDiv, err := filterCard.
					Element("dl.filterCardTitle")
				if err != nil {
					log.Println(err)
					return output, err
				}

				// errorページが表示されているか確認
				if strings.Contains(page.MustInfo().URL, "error") {
					log.Println("エラーページが表示されました")
					continue templateLoop
				}

				filterCardTitle, err := filterCardTitleDiv.
					Element("a")
				if err != nil {
					log.Println(err)
					return output, err
				}

				// errorページが表示されているか確認
				if strings.Contains(page.MustInfo().URL, "error") {
					log.Println("エラーページが表示されました")
					continue templateLoop
				}

				if filterCardTitle.MustText() == scoutServiceTemplate.SearchTitle {
					isMatchedWithSearchTitle = true
					filterCardBtnDiv, err := filterCard.
						Element("div.filterCardButtonMain")
					if err != nil {
						log.Println(err)
						return output, err
					}

					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					filterCardBtn, err := filterCardBtnDiv.
						Element("button")
					if err != nil {
						log.Println(err)
						return output, err
					}

					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					filterCardBtn.MustClick()
					page.WaitLoad()
					time.Sleep(5 * time.Second)
					break
				}
			}
			if !isMatchedWithSearchTitle {
				errMessage = "保存した検索条件が見つかりませんでした.scoutServiceTemplate.SearchTitle:" + scoutServiceTemplate.SearchTitle

				time.Sleep(20 * time.Second)
				// 保存した検索条件を選択する
				filterCardList, err = page.
					Elements("div.filterCard")
				if err != nil {
					log.Println(err)
					return output, err
				}

				// errorページが表示されているか確認
				if strings.Contains(page.MustInfo().URL, "error") {
					log.Println("エラーページが表示されました")
					continue templateLoop
				}

				isMatchedWithSearchTitle := false
				for _, filterCard := range filterCardList {
					filterCardTitleDiv, err := filterCard.
						Element("dl.filterCardTitle")
					if err != nil {
						log.Println(err)
						return output, err
					}

					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					filterCardTitle, err := filterCardTitleDiv.
						Element("a")
					if err != nil {
						log.Println(err)
						return output, err
					}

					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					if filterCardTitle.MustText() == scoutServiceTemplate.SearchTitle {
						isMatchedWithSearchTitle = true
						filterCardBtnDiv, err := filterCard.
							Element("div.filterCardButtonMain")
						if err != nil {
							log.Println(err)
							return output, err
						}

						// errorページが表示されているか確認
						if strings.Contains(page.MustInfo().URL, "error") {
							log.Println("エラーページが表示されました")
							continue templateLoop
						}

						filterCardBtn, err := filterCardBtnDiv.
							Element("button")
						if err != nil {
							log.Println(err)
							return output, err
						}

						// errorページが表示されているか確認
						if strings.Contains(page.MustInfo().URL, "error") {
							log.Println("エラーページが表示されました")
							continue templateLoop
						}

						filterCardBtn.MustClick()
						page.WaitLoad()
						time.Sleep(5 * time.Second)
						break
					}
				}
				if !isMatchedWithSearchTitle {
					errMessage = "保存した検索条件が見つかりませんでした.scoutServiceTemplate.SearchTitle:" + scoutServiceTemplate.SearchTitle
					log.Println(errMessage)
					return output, errors.New(errMessage)
				}
			}

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			// 全て選択
			allSelectBtn, err := page.
				Element("button.searchNavButton.isCheck")
			if err != nil {
				log.Println(err)
				return output, err
			}
			allSelectBtn.MustClick()

			time.Sleep(5 * time.Second)

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			if page.MustHas("a.searchStateNavMailButton.isNoSelect") {
				log.Println("page.MustHas(a.searchStateNavMailButton.isNoSelect)")
				time.Sleep(10 * time.Second)
				if page.MustHas("a.searchStateNavMailButton.isNoSelect") {
					log.Println("page.MustHas(a.searchStateNavMailButton.isNoSelect)2")
					time.Sleep(10 * time.Second)
					if page.MustHas("a.searchStateNavMailButton.isNoSelect") {
						log.Println("page.MustHas(a.searchStateNavMailButton.isNoSelect)2")
						time.Sleep(10 * time.Second)
					}
				}
			}

			// メール文を作成する searchStateNavMailButton
			mailBtns, err := page.
				Elements("a.searchStateNavMailButton")
			if err != nil {
				log.Println(err)
				return output, err
			}

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			isMatchedWithMailBtn := false
			for _, mailBtn := range mailBtns {
				if mailBtn.MustText() == "メール文を作成する" {
					isMatchedWithMailBtn = true
					mailBtn.MustClick()
					page.WaitLoad()
					time.Sleep(10 * time.Second)
					break
				}
			}
			if !isMatchedWithMailBtn {
				errMessage = "メール文を作成するボタンが見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			// mailModalInner
			mailModalInner, err := page.
				Element("div.mailModalInner")
			if err != nil {
				log.Println(err)
				return output, err
			}

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			// テンプレートを利用する
			templateSelectBtns, err := mailModalInner.
				Elements("button.generalButton")
			if err != nil {
				log.Println(err)
				return output, err
			}

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			isMatchedWithTemplateSelectBtn := false
			for _, templateSelectBtn := range templateSelectBtns {
				if templateSelectBtn.MustText() == "テンプレートを利用する" {
					isMatchedWithTemplateSelectBtn = true
					templateSelectBtn.MustClick()
					page.WaitLoad()
					time.Sleep(20 * time.Second)
					break
				}
			}
			if !isMatchedWithTemplateSelectBtn {
				errMessage = "テンプレートを利用するボタンが見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			isMatchedWithMessageTitle := false
		messageTitlePageLoop:
			for pageI := 1; pageI <= 5; pageI++ {
				// mailModalInner
				modalInner, err := page.
					Element("div.modalInner")
				if err != nil {
					log.Println(err)
					return output, err
				}

				// errorページが表示されているか確認
				if strings.Contains(page.MustInfo().URL, "error") {
					log.Println("エラーページが表示されました")
					continue templateLoop
				}

				// filterCard
				if !modalInner.MustHas("div.filterCard") {
					time.Sleep(20 * time.Second)
				}

				templateCardList, err := modalInner.
					Elements("div.filterCard")
				if err != nil {
					log.Println(err)
					return output, err
				}

				log.Println("templateCardList:", len(templateCardList))

				if len(templateCardList) == 0 {
					time.Sleep(30 * time.Second)
					// filterCard
					templateCardList, err = modalInner.
						Elements("div.filterCard")
					if err != nil {
						log.Println(err)
						return output, err
					}
					if len(templateCardList) == 0 {
						err = errors.New("テンプレートが見つかりませんでした")
						return output, err
					}
				}

				for _, templateCard := range templateCardList {
					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					templateCardTemplateDiv, err := templateCard.
						Element("dl.filterCardTemplate")
					if err != nil {
						log.Println(err)
						return output, err
					}

					log.Println("templateCardTemplateDiv.MustText():", templateCardTemplateDiv.MustText())

					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					templateCardA, err := templateCardTemplateDiv.
						Element("a")
					if err != nil {
						log.Println(err)
						return output, err
					}

					log.Println("templateCardA.MustText():", templateCardA.MustText())

					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					if templateCardA.MustText() == scoutServiceTemplate.MessageTitle {
						isMatchedWithMessageTitle = true
						templateCardBtn, err := templateCard.
							Element("button.generalButton")
						if err != nil {
							log.Println(err)
							return output, err
						}
						// errorページが表示されているか確認
						if strings.Contains(page.MustInfo().URL, "error") {
							log.Println("エラーページが表示されました")
							continue templateLoop
						}
						templateCardBtn.MustClick()
						page.WaitLoad()
						time.Sleep(2 * time.Second)
						break messageTitlePageLoop
					}
				}
			}
			if !isMatchedWithMessageTitle {
				errMessage = "テンプレートが見つかりませんでした。:" + scoutServiceTemplate.MessageTitle
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			// スカウトを送信 mailModalFooter
			mailModalFooter, err := page.
				Element("div.mailModalFooter")
			if err != nil {
				log.Println(err)
				return output, err
			}

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			sendBtns, err := mailModalFooter.
				Elements("button.generalButton")
			if err != nil {
				log.Println(err)
				return output, err
			}

			// errorページが表示されているか確認
			if strings.Contains(page.MustInfo().URL, "error") {
				log.Println("エラーページが表示されました")
				continue templateLoop
			}

			isMatchedWithSendBtn := false
			for _, sendBtn := range sendBtns {
				if sendBtn.MustText() == "スカウトを送信" {
					isMatchedWithSendBtn = true
					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					err := sendBtn.Click(proto.InputMouseButtonLeft, 1)
					if err != nil {
						log.Println(err)
						continue templateLoop
					}
					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					page.WaitLoad()
					time.Sleep(2 * time.Second)

					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					// 送信可能な時間は、8:00～22:00です。

					if !page.MustHas("div.modalContainer") {
						errMessage = "スカウト送信確認が見つかりませんでした"
						log.Println(errMessage)
						return output, errors.New(errMessage)
					}

					modalContainer, err := page.
						Element("div.modalContainer")
					if err != nil {
						log.Println(err)
						return output, err
					}

					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

					confirmSendBtns, err := modalContainer.
						Elements("button.generalButton")
					if err != nil {
						log.Println(err)
						return output, err
					}

					isMatchedWithConfirmSendBtn := false
					for _, confirmSendBtn := range confirmSendBtns {
						// errorページが表示されているか確認
						if strings.Contains(page.MustInfo().URL, "error") {
							log.Println("エラーページが表示されました")
							continue templateLoop
						}
						log.Println("confirmSendBtn.MustText():", confirmSendBtn.MustText())

						if confirmSendBtn.MustText() == "送信する" {
							isMatchedWithConfirmSendBtn = true
							if strings.Contains(page.MustInfo().URL, "error") {
								log.Println("エラーページが表示されました")
								continue templateLoop
							}
							// if i.app.BatchType == "scout" {
							confirmSendBtn.MustClick()
							page.WaitLoad()
							// }
							time.Sleep(2 * time.Second)

							// errorページが表示されているか確認
							if strings.Contains(page.MustInfo().URL, "error") {
								log.Println("エラーページが表示されました")
								continue templateLoop
							}

							finalConfirmModalContainer, err := page.
								Element("div.modalContainer")
							if err != nil {
								log.Println(err)
								return output, err
							}

							// errorページが表示されているか確認
							if strings.Contains(page.MustInfo().URL, "error") {
								log.Println("エラーページが表示されました")
								continue templateLoop
							}

							finalConfirmSendBtns, err := finalConfirmModalContainer.
								Elements("button.generalButton")
							if err != nil {
								log.Println(err)
								return output, err
							}

							isMatchedWithfinalConfirmSendBtn := false
							for _, finalConfirmSendBtn := range finalConfirmSendBtns {
								// errorページが表示されているか確認
								if strings.Contains(page.MustInfo().URL, "error") {
									log.Println("エラーページが表示されました")
									continue templateLoop
								}
								log.Println("finalConfirmSendBtn.MustText():", finalConfirmSendBtn.MustText())

								if finalConfirmSendBtn.MustText() == "送信する" {
									isMatchedWithfinalConfirmSendBtn = true
									// errorページが表示されているか確認
									if strings.Contains(page.MustInfo().URL, "error") {
										log.Println("エラーページが表示されました")
										continue templateLoop
									}
									// if i.app.BatchType == "scout" {
									finalConfirmSendBtn.MustClick()
									page.WaitLoad()
									// }
									time.Sleep(40 * time.Second)

									// 送信されたか確認
									page.MustNavigate("https://scout.mynavi-agent.jp/analyse/member/")
									page.WaitLoad()
									time.Sleep(20 * time.Second)

									// errorページが表示されているか確認
									if strings.Contains(page.MustInfo().URL, "error") {
										page.MustNavigate("https://scout.mynavi-agent.jp/analyse/member/")
										page.WaitLoad()
										time.Sleep(30 * time.Second)
									}

									tableItemTexts, err := page.
										Elements("span.tableItemText")
									if err != nil {
										log.Println(err)
										return output, err
									}
									if len(tableItemTexts) < 2 {
										time.Sleep(30 * time.Second)
										tableItemTexts, err = page.
											Elements("span.tableItemText")
										if err != nil {
											log.Println(err)
											return output, err
										}
										if len(tableItemTexts) < 2 {
											errMessage = "スカウト送信確認ボタンが見つかりませんでした"
											log.Println(errMessage)
											return output, errors.New(errMessage)
										}
									}

									totalScoutCountStr := tableItemTexts[1].MustText()
									totalScoutCount, err := strconv.Atoi(totalScoutCountStr)
									if err != nil {
										log.Println(err)
										return output, err
									}
									log.Println("totalScoutCount:", totalScoutCount)

									log.Println("scoutServiceTemplate.LastSendCount.Int64:", scoutServiceTemplate.LastSendCount.Int64, "totalScoutCount:", totalScoutCount)
									if int(input.ScoutService.LastSendCount.Int64) == totalScoutCount {
										log.Println("スカウトが送信されていません。")
										continue templateLoop
									}

									scoutServiceTemplate.LastSendCount.Int64 += 500
									input.ScoutService.LastSendCount.Int64 = int64(totalScoutCount)
									log.Println("スカウト送信完了しました")

									// スカウト送信完了後、スカウト情報を更新
									err = i.scoutServiceTemplateRepository.UpdateLastSend(
										scoutServiceTemplate.ID,
										uint(scoutServiceTemplate.LastSendCount.Int64),
										time.Now().UTC(),
									)
									if err != nil {
										log.Println(err)
										return output, err
									}

									err = i.scoutServiceRepository.UpdateLastSendCount(
										input.ScoutService.ID,
										uint(input.ScoutService.LastSendCount.Int64),
									)
									if err != nil {
										log.Println(err)
										return output, err
									}
									continue templateLoop
								}
							}
							if !isMatchedWithfinalConfirmSendBtn {
								errMessage = "スカウト送信確認ボタンが見つかりませんでした"
								log.Println(errMessage)
								return output, errors.New(errMessage)
							}
						}
					}
					if !isMatchedWithConfirmSendBtn {
						errMessage = "スカウト送信確認ボタンが見つかりませんでした"
						log.Println(errMessage)
						return output, errors.New(errMessage)
					}

					// errorページが表示されているか確認
					if strings.Contains(page.MustInfo().URL, "error") {
						log.Println("エラーページが表示されました")
						continue templateLoop
					}

				}
			}
			if !isMatchedWithSendBtn {
				errMessage = "スカウト送信ボタンが見つかりませんでした"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

		}
	}

	// 成功メールを送信
	err = i.sendScoutSuccessMail(input.ScoutService.AgentStaffID, input.ScoutService, input.ScoutServiceTemplateList)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.OK = true
	return output, err
}
