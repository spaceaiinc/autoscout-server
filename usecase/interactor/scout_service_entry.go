package interactor

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"google.golang.org/api/gmail/v1"
	"gopkg.in/guregu/null.v4"

	// ブラウザ操作
	"github.com/go-rod/rod"
	rodInput "github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"

	"golang.org/x/oauth2"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	// csvの文字化け変換
)

type BatchEntryInput struct {
	Now time.Time
}

type BatchEntryOutput struct {
	OK bool
}

func (i *ScoutServiceInteractorImpl) BatchEntry(input BatchEntryInput) (BatchEntryOutput, error) {
	var (
		output BatchEntryOutput
		err    error

		ambiUserIDList           = []string{} // AmbiのユーザーID
		mynaviScoutingUserIDList = []string{} // マイナビス転職スカウトのユーザーID
	)

	/********* 未実行のユーザー情報を取得 *********/

	unprocessedUserList, err := i.userEntryRepository.GetUnprocessed()
	if err != nil {
		return output, err
	}

	fmt.Println("unprocessedUserList", len(unprocessedUserList))

	/********* 未実行のユーザーのエントリー処理を実行（api実行中にメールを受信した場合の対策） *********/

	if len(unprocessedUserList) > 0 {
		/**
		 * エントリー取得処理
		 */
		var (
			errMessage           string
			selectedScoutService *entity.ScoutService
			now                  time.Time = time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
			agentRobotID         uint
		)

		agentRobotIDInt, err := strconv.Atoi(os.Getenv("AGENT_ROBOT_ID"))
		if err != nil {
			log.Println("エージェントロボットの変換に失敗しました", err)
			return output, err
		}
		agentRobotID = uint(agentRobotIDInt)

		agentRobot, err := i.agentRobotRepository.FindByID(agentRobotID)
		if err != nil {
			log.Println("エージェントロボットの取得に失敗しました", err)
			return output, err
		}

		scoutServices, err := i.scoutServiceRepository.GetByAgentRobotID(agentRobotID)
		if err != nil {
			log.Println("スカウトサービスの取得に失敗しました", err)
			return output, err
		}
		// AMBIの方がエラーが少ないため優先
		sort.Slice(scoutServices, func(i, j int) bool {
			return scoutServices[i].ServiceType.Int64 == entity.ScoutServiceTypeAmbi
		})

		// タイムアウトを設定
		ctx, cancel := context.WithTimeout(context.Background(), 14*time.Minute)
		defer cancel()

		// panicが発生した場合にエラーを返すためにdeferでrecoverを実行して、エラーとして返す
		defer func() {
			if rec := recover(); rec != nil {
				// interactor or repositoryを含む文字列を取得
				stackStr := string(debug.Stack())
				var errorLine string
				for _, line := range strings.Split(stackStr, "\n") {
					if strings.Contains(line, "interactor/") || strings.Contains(line, "repository/") {
						errorLine += line + "\n"
						log.Println("error line:", errorLine)
					}
				}
				errMessage = fmt.Sprintf("エントリー取得処理でpanicエラーが発生しました。Recover: %v\nStack:\n%s", rec, errorLine)
				log.Println(errMessage)
				err = errors.New(errMessage)

				// エラーメール送信
				now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
				errorIssue := "パニックエラー"
				if strings.Contains(fmt.Sprint(rec), "deadline") {
					errorIssue = "タイムアウト"
				}
				if selectedScoutService != nil {
					i.sendErrorMail(
						fmt.Sprintf(
							"%vの新規エントリー取得で%vが発生しました。\n発生時刻: %s\nロボット名: %s\nロボットID: %v\n・Recover: %v\n・Stack:\n%s",
							entity.ScoutServiceTypeLabel[selectedScoutService.ServiceType.Int64], errorIssue, now, agentRobot.Name, agentRobot.ID, rec, errorLine,
						),
					)
				}
			}
		}()

		// 未実行のユーザーIDをセット
		for _, unprocessedUser := range unprocessedUserList {
			if unprocessedUser.ServiceType == null.NewInt(entity.ScoutServiceTypeMynaviScouting, true) {
				mynaviScoutingUserIDList = append(mynaviScoutingUserIDList, unprocessedUser.UserID)
			}
			if unprocessedUser.ServiceType == null.NewInt(entity.ScoutServiceTypeAmbi, true) {
				ambiUserIDList = append(ambiUserIDList, unprocessedUser.UserID)
			}
		}

		for _, scoutService := range scoutServices {

			selectedScoutService = scoutService

			// マイナビスカウティング
			if scoutService.ServiceType == null.NewInt(entity.ScoutServiceTypeMynaviScouting, true) &&
				len(mynaviScoutingUserIDList) > 0 {
				log.Println("マイナビのエントリー取得を開始します\n現在:", now.Hour())

				// スカウト処理中の場合はログアウトさせないように終わるまで処理しない。
				serviceOutput, err := i.EntryOnMynaviScouting(EntryOnMynaviScoutingInput{
					AgentID:      agentRobot.AgentID,
					ScoutService: scoutService,
					Context:      ctx,

					UserIDList: mynaviScoutingUserIDList,
				})
				if err != nil {
					now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
					errMessage = fmt.Sprintf("マイナビスカウティングの新規エントリー求職者取得に失敗しました。\n発生時刻: %s\nロボット名: %s\nロボットID: %v\nエラー内容:%s", now, agentRobot.Name, agentRobot.ID, err.Error())
					log.Println(errMessage)
					i.sendErrorMail(errMessage)
					return output, errors.New(errMessage)
				}

				// 処理したマイナビエントリーユーザーのIDのフラグを更新
				err = i.userEntryRepository.UpdateIsProcessedByUserIDList(mynaviScoutingUserIDList, true)
				if err != nil {
					return output, err
				}

				log.Println("マイナビの取得に成功しました", serviceOutput.JobSeekerList)

				output.OK = true
				return output, err

				// AMBI
			} else if scoutService.ServiceType == null.NewInt(entity.ScoutServiceTypeAmbi, true) &&
				len(ambiUserIDList) > 0 {
				log.Println("AMBIのエントリー取得を開始します\n現在:", now.Hour())

				serviceOutput, err := i.EntryOnAmbi(EntryOnAmbiInput{
					AgentID:      agentRobot.AgentID,
					ScoutService: scoutService,
					Context:      ctx,

					UserIDList: ambiUserIDList,
				})
				if err != nil {
					now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
					errMessage = fmt.Sprintf("AMBIの新規エントリー求職者取得に失敗しました。\n発生時刻: %s\nロボット名: %s\nロボットID: %v\nエラー内容:%s", now, agentRobot.Name, agentRobot.ID, err.Error())
					log.Println(errMessage)
					i.sendErrorMail(errMessage)
					return output, errors.New(errMessage)
				}

				// 処理したAMBIエントリーユーザーのIDのフラグを更新
				err = i.userEntryRepository.UpdateIsProcessedByUserIDList(ambiUserIDList, true)
				if err != nil {
					return output, err
				}

				log.Println("AMBIの取得に成功しました", serviceOutput.JobSeekerList)

				output.OK = true
				return output, err
			}
		}
	} else {
		fmt.Println("---------------\nエントリーユーザーが存在しません。\n---------------\n")

		// 定期実行フラグ
	}

	// var (
	// 	amnbiIDListStr          string
	// 	mynaviScoutingIDListStr string
	// 	// mynaviAgentScoutIDListStr     string
	// )

	// // Ambiの求職者ID
	// if len(ambiUserIDList) > 0 {
	// 	idStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ambiUserIDList)), ", "), "[]")
	// 	amnbiIDListStr = fmt.Sprintf("【AMBI】\n%s", idStr)
	// 	fmt.Printf("---------------\n%s\n---------------\n", amnbiIDListStr)
	// }

	// // マイナビ転職の求職者ID
	// if len(mynaviScoutingUserIDList) > 0 {
	// 	idStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(mynaviScoutingUserIDList)), ", "), "[]")
	// 	mynaviScoutingIDListStr = fmt.Sprintf("【マイナビスカウティング】\n%s", idStr)
	// 	fmt.Printf("---------------\n%s\n---------------\n", mynaviScoutingIDListStr)
	// }

	output.OK = true
	return output, err
}

/****************************************************************************************/
// Gmail Hook処理 API
//
type GmailWebHookInput struct {
	PubSubStruct *entity.PubsubStruct
}

type GmailWebHookOutput struct {
	OK bool
}

func (i *ScoutServiceInteractorImpl) GmailWebHook(input GmailWebHookInput) (GmailWebHookOutput, error) {
	var (
		output       GmailWebHookOutput
		err          error
		messageData  entity.MeessageData
		mailDataList []string
	)

	/********* 認証情報 *********/

	// 認証情報からconfigを生成
	config, err := utility.NewGoogleAuthConf(i.googleAPI.JSONFilePath)
	if err != nil {
		return output, err
	}

	tok, err := i.googleAuthenticationRepository.FindLatest()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	token := &oauth2.Token{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
	}

	srv, err := utility.GenerateService(config, token)
	if err != nil {
		return output, err
	}

	// デコード前のData
	encodedData := input.PubSubStruct.Message.Data

	// Base64デコード
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		wrapped := fmt.Errorf("decode error: %v", err)
		return output, wrapped
	}

	// 形式: {"emailAddress":"info@spaceai.jp","historyId":526295}
	decodedData := string(decodedBytes)

	// JSON文字列を解析して構造体に格納
	if err := json.Unmarshal([]byte(decodedData), &messageData); err != nil {
		wrapped := fmt.Errorf("JSON Unmarshal error: %v", err)
		return output, wrapped
	}

	/********* 指定のメールを取得する *********/

	user := "me"

	var (
		ambiEmail             = "ambi-support@en-japan.com"    // AMBI事務局
		mynaviScoutingEmail   = "ags-admin@mynavi.jp"          // マイナビス転職スカウト
		mynaviAgentScoutEmail = "sk-mag-scout@mynavi-agent.jp" // マイナビAGENTスカウト
		// dodaXEmail           = "partner_support@doda-x.jp"    // doda X（dodax_search@persol.co.jpも追加で必要かも）

		ambiSubject             = "\"スカウトへのエントリー/AMBI\""                       // AMBIの件名
		mynaviScoutingSubject   = "\"スカウト応募（マイナビスカウティング）がありました【マイナビスカウティング】\"" // マイナビ転職の件名
		mynaviAgentScoutSubject = "\"【マイナビAGENTスカウト】求職者がエントリーしました\""           // マイナビAGENTの件名

		ambiUserIDList           = []string{} // AmbiのユーザーID
		mynaviScoutingUserIDList = []string{} // マイナビス転職スカウトのユーザーID
		// mynaviAgentScoutUserIDList     = []string{} // マイナビAGENTスカウトのユーザーID
		// dodaXUserIDList       = []string{} // doda XのユーザーID
	)

	// local用
	// if i.app.Env != "prd" {
	// 	ambiEmail = "tonbi0521@gmail.com"
	// 	mynaviScoutingEmail = "sibuhiro81@gmail.com"
	// }

	// Gmailの検索条件: https://support.google.com/mail/answer/7190?hl=ja
	var searchParam = fmt.Sprintf(
		"((from:%s subject:%s) OR (from:%s subject:%s) OR (from:%s subject:%s)) is:unread",
		ambiEmail, ambiSubject,
		mynaviScoutingEmail, mynaviScoutingSubject,
		mynaviAgentScoutEmail, mynaviAgentScoutSubject,
	)

	// Fromと未読の条件を指定してメールを取得
	mes, err := srv.Users.Messages.List(user).Q(searchParam).Do()
	if err != nil {
		wrapped := fmt.Errorf("unable to retrieve history: %s", err)
		return output, wrapped
	}
	if len(mes.Messages) == 0 {
		log.Println("新着メールはありません")
		return output, nil
	}

	fmt.Println("mes.Messagesの長さ", len(mes.Messages))

	for _, message := range mes.Messages {
		msg, err := srv.Users.Messages.Get(user, message.Id).Do()
		if err != nil {
			wrapped := fmt.Errorf("unable to retrieve message %v: %s", message.Id, err)
			return output, wrapped
		}
		var (
			from string
			// to      string
			subject string
			body    string
		)

		for _, header := range msg.Payload.Headers {
			switch header.Name {
			case "From":
				address, err := mail.ParseAddress(header.Value)
				if err != nil {
					wrapped := fmt.Errorf("error parsing address: %v", err)
					return output, wrapped
				}

				from = address.Address
			case "Subject":
				subject = header.Value
				// case "To":
				// 	to = header.Value
			}
		}

		// メール本文の表示
		if msg.Payload.Body.Size > 0 {
			// Base64デコード
			bodyBytes, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
			if err != nil {
				wrapped := fmt.Errorf("error decoding message body: %v", err)
				return output, wrapped
			} else {
				body = string(bodyBytes)
			}
		} else if len(msg.Payload.Parts) > 0 {
			// MIMEメッセージの場合、パーツを調べて本文を探す
			for _, part := range msg.Payload.Parts {
				if part.MimeType == "text/plain" {
					// Base64デコード
					bodyBytes, err := base64.URLEncoding.DecodeString(part.Body.Data)
					if err != nil {
						wrapped := fmt.Errorf("error decoding multipart message body: %v", err)
						return output, wrapped
					} else {
						body = string(bodyBytes)
						break
					}
				}
			}
		}

		var (
			userID      string
			serviceType null.Int
		)

		fmt.Printf("\n---------------\n%s\n---------------\n", body)

		/********* メール送信元の媒体に応じてユーザーIDを振り分け *********/

		switch from {

		// Ambiエントリーメールの処理
		case ambiEmail:
			matches := utility.RegexpForAmbiUserID.FindStringSubmatch(body)

			fmt.Println("Ambi matches: ", matches, len(matches))

			if len(matches) >= 2 {
				fmt.Println("会員番号:", matches[1])

				// ユーザーIDと媒体の種類をセット
				userID = matches[1]
				serviceType = null.NewInt(entity.ScoutServiceTypeAmbi, true)

				// ユーザーIDをリストに追加
				ambiUserIDList = append(ambiUserIDList, matches[1])
			}

			if len(ambiUserIDList) == 0 {
				fmt.Println("会員番号が見つかりませんでした（Ambi）")
				return output, errors.New("会員番号が見つかりませんでした（Ambi）")
			}

		// マイナビスカウティングのエントリーメールの処理
		case mynaviScoutingEmail:
			matches := utility.RegexpForMynaviScoutingUserID.FindStringSubmatch(body)

			fmt.Println("マイナビ転職 matches: ", matches, len(matches))

			if len(matches) >= 2 {
				fmt.Println("会員番号:", matches[1])

				// ユーザーIDと媒体の種類をセット
				userID = matches[1]
				serviceType = null.NewInt(entity.ScoutServiceTypeMynaviScouting, true)

				// ユーザーIDをリストに追加
				mynaviScoutingUserIDList = append(mynaviScoutingUserIDList, matches[1])
			}

			if len(mynaviScoutingUserIDList) == 0 {
				fmt.Println("会員番号が見つかりませんでした（マイナビスカウティング）")
				return output, errors.New("会員番号が見つかりませんでした（マイナビスカウティング）")
			}

		// マイナビエージェントスカウトのエントリーメールの処理
		case mynaviAgentScoutEmail:
			fmt.Println("マイナビエージェントスカウトは未対応です")
			// return output, nil
		// scoutServiceType = uint(entity.ScoutServiceTypemynaviAgentScoutScout)

		// matches := utility.RegexpFormynaviAgentScoutScoutUserID.FindStringSubmatch(body)

		// if len(matches) > 1 {
		// 	fmt.Println("会員番号:", matches[1], matches)
		// 	mynaviAgentScoutUserIDList = append(mynaviAgentScoutUserIDList, matches[1])
		// } else {
		// 	fmt.Println("会員番号が見つかりませんでした。")
		// }

		default:
			fmt.Println("対象スカウトサービスのメールではありません。from: ", from)
			// return output, errors.New("対象スカウトサービスのメールではありません")
		}

		// ログにメールの詳細をフォーマットして表示
		str := fmt.Sprintf("件名: %s\n", subject)
		mailDataList = append(mailDataList, str)
		fmt.Printf("From: %s\n", from)

		/********* 指定メールの既読更新 *********/

		// UNREADラベルを削除して既読にマーク
		// スコープの追加または削除から、Gmail APIの".../auth/gmail.modify"を追加する必要あり（https://zenn.dev/acn_jp_sdet/articles/f3f88523954626）
		if from == ambiEmail || from == mynaviScoutingEmail {
			fmt.Println("既読処理を実行！！")
			_, err = srv.Users.Messages.Modify(user, message.Id, &gmail.ModifyMessageRequest{
				RemoveLabelIds: []string{"UNREAD"},
			}).Do()
			if err != nil {
				errorMessage := fmt.Sprintf("メールを既読にすることができませんでした:  %s\n", err)
				fmt.Println(errorMessage)
				return output, errors.New(errorMessage)
			} else {
				fmt.Println("既読にしました")

				// エントリー処理を実行するユーザーを登録
				userEntry := entity.NewUserEntry(userID, serviceType)
				err = i.userEntryRepository.Create(userEntry)
				if err != nil {
					return output, err
				}
			}
		}
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/
// 新しくエントリーした応募者取得 API
//
/*
		Ranから新規応募者を取得する

	  開発環境でかかった時間 5分(5名登録):
		2023/04/14 15:01:00 RANのエントリー取得を開始します
		2023/04/14 15:05:03 RANの取得に成功しました [0xc0003ecc00 0xc0003ed200 0xc0003ed800 0xc0005d6000 0xc0005d6600]
*/
type EntryOnRanInput struct {
	AgentID      uint
	Context      context.Context
	ScoutService *entity.ScoutService
}

type EntryOnRanOutput struct {
	JobSeekerList []*entity.JobSeeker
}

func (i *ScoutServiceInteractorImpl) EntryOnRan(input EntryOnRanInput) (EntryOnRanOutput, error) {

	var (
		output        EntryOnRanOutput
		err           error
		errMessage    string
		jobSeekerList []*entity.JobSeeker
		browser       *rod.Browser
		page          *rod.Page
		detailPage    *rod.Page
		scoutService  *entity.ScoutService
	)

	scoutService = input.ScoutService

	if scoutService.LoginID == "" || scoutService.Password == "" {
		errMessage = "ScoutServiceの認証情報が未入力です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	} else if input.Context == nil {
		errMessage = "Contextが空です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	} else if input.AgentID == 0 {
		errMessage = "AgentIDが未入力です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// 日本時間00:00〜05:00の間は使用できない
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	if now.Hour() < 5 {
		errMessage = "日本時間00:00〜05:00の間は使用できません"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// 2日以内に登録されたの応募者を取得
	jobSeekerListInDb, err := i.jobSeekerRepository.GetByAgentIDWithinTwoDays(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// ホストOSのchrome, chromium or microsoft edgeのパスを取得
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
	// remove launcher.FlagUserDataDir
	defer l.Cleanup()

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

	page = browser.
		// ログインページへ遷移
		MustPage("https://ran.next.rikunabi.com/rnc/docs/ca_s01000.jsp").
		MustWaitLoad()

	time.Sleep(4 * time.Second)

	// ログイン情報を入力
	page.MustElement("input[name=login_nm]").
		MustInput(scoutService.LoginID)

	// パスワードの複号
	decryptedPassword, err := decryption(scoutService.Password)
	if err != nil {
		log.Println(err)
		return output, err
	}

	page.MustElement("input[name=pswd]").
		MustInput(decryptedPassword).
		MustType(rodInput.Enter)

	page.WaitLoad()
	time.Sleep(4 * time.Second)

	// URLがログインページのままの場合はログインに失敗していると判断
	if strings.Contains(page.MustInfo().URL, "ca_i01000.jsp") {
		errMessage = "ログインに失敗しました。"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// エントリー管理ページへ遷移
	page.MustNavigate("https://ran.next.rikunabi.com/rnc/docs/ca_f00010.jsp?tab_select_id=4")
	page.WaitLoad()
	time.Sleep(4 * time.Second)

	// 担当者ごとに選択 *ローカル環境は濱田さんを選択
	// if i.app.Env == "local" {
	// 	frame := page.MustElement("frame[name=lowerframe]").MustFrame()
	// 	frame.MustElement("select[name=charge_cnsltnt_user_no]").MustSelect("濱田 僚太")
	// 	frame.MustElement("a.btn01.btn_m.large.w300").MustClick().WaitLoad()
	// 	time.Sleep(2 * time.Second)
	// }

	// ページング処理 10ページまで
	// for pageI := 1; pageI <= 10; pageI++ {
	// ダミーから主要情報が載っているframeへ切り替え
	frame := page.MustElement("frame[name=lowerframe]").MustFrame()

	userTrs := frame.MustElement("table.cell_middle").MustElement("tbody").MustElements("tr")
	if len(userTrs) == 0 {
		log.Println("userTrsの長さが0です")
		return output, nil
	}

userLoop:
	for userI, userTr := range userTrs {
		if userI == 0 {
			continue
		}

		jobSeeker := entity.JobSeeker{}
		userTds := userTr.MustElements("td")
		if len(userTds) < 6 {
			log.Println("userTdsの長さが6未満")
			continue
		}
		entryDateStr := userTds[5].MustText()
		log.Println("entryDate:", entryDateStr)

		// 申請日が2日以上前の場合は、スキップする 2022/07/07 13:02
		entryDate, err := time.Parse("2006-01-02\n15:04", entryDateStr)
		if err != nil {
			log.Println("申請日のパースに失敗しました", entryDateStr)
			continue userLoop
		}

		// 2日以上前の場合は、スキップする
		now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
		if entryDate.Before(now.AddDate(0, 0, -2)) {
			log.Println("2日以上前のエントリーはスキップします", entryDate)
			break userLoop
		}

		nameWithLink := userTds[1].MustElement("a")

		nameWithSpace := nameWithLink.MustText()
		nameSplited := strings.Split(nameWithSpace, "\u00a0\u00a0")
		if len(nameSplited) < 2 {
			continue
		}
		jobSeeker.LastName = nameSplited[0]
		jobSeeker.FirstName = strings.Split(nameSplited[1], "\n")[0]

		log.Println("jobSeeker:", jobSeeker.LastName, jobSeeker.FirstName, jobSeeker.LastFurigana, jobSeeker.FirstFurigana)

		// 「名前」が既存の求職者と重複している場合はスキップ
		for _, jobSeekerInDb := range jobSeekerListInDb {
			if jobSeekerInDb.LastName == jobSeeker.LastName &&
				jobSeekerInDb.FirstName == jobSeeker.FirstName {
				continue userLoop
			}
		}

		// ユーザー詳細ページへ遷移
		nameWithLink.MustClick()
		time.Sleep(10 * time.Second)

		// 求職者の詳細ページを取得
		detailPage = nil
		pages := browser.MustPages()
		for _, page := range pages {
			log.Println("page", page.MustInfo().URL)
			if strings.Contains(page.MustInfo().URL, "ca_s02050") {
				detailPage = page
				break
			}
		}
		if detailPage == nil {
			log.Println("detailPage is nil")
			continue
		}

		// ユーザー情報を取得
		var (
			profileTrs    []*rod.Element
			educationTrs  []*rod.Element
			experienceTrs []*rod.Element
		)

		// 社内限定メモに記載する
		jobSeeker.SecretMemo = "・エントリー媒体\nRAN\n\n"

		// オファー履歴
		if detailPage.MustHas("p.mb10") {
			offerHistoryPs := detailPage.MustElements("p.mb10")
			for _, offerHistoryP := range offerHistoryPs {
				offerHistoryPText := offerHistoryP.MustText()
				log.Println("offerHistoryPText:", offerHistoryPText)
				if offerHistoryPText == "この会員に対して送信されたオファーの履歴です。" {
					offerHistoryTable := offerHistoryP.MustNext()
					offerHistoryTableText := offerHistoryTable.MustText()
					log.Println("offerHistoryTableText:", offerHistoryTableText)
					jobSeeker.SecretMemo += "・オファー履歴\n" + offerHistoryTableText + "\n\n"
				}
			}
		}

		tableTitlePs := detailPage.MustElements("p.bold.large.mb5")
		log.Println("tableTitlePs:", tableTitlePs)
		for _, tableTitleP := range tableTitlePs {
			tableTitle := tableTitleP.MustText()
			log.Println("tableTitle:", tableTitle)
			if tableTitle == "プロフィール" {
				profileTable := tableTitleP.MustNext()
				log.Println("profileTable:", profileTable)
				profileTrs = profileTable.MustElement("tbody").MustElements("tr")
			} else if tableTitle == "学歴・語学・資格" {
				educationTable := tableTitleP.MustNext()
				log.Println("educationTable:", educationTable)
				educationTrs = educationTable.MustElement("tbody").MustElements("tr")
			} else if tableTitle == "職務経歴" {
				experienceTable := tableTitleP.MustNext()
				log.Println("experienceTable:", experienceTable)
				experienceTrs = experienceTable.MustElement("tbody").MustElements("tr")
			}
		}

		for _, profileTr := range profileTrs {
			profileTh := profileTr.MustElement("th").MustText()
			log.Println("profileTh:", profileTh)

			// ヤマダ タロウ
			// 山田 太郎　（会員番号19113702）
			if profileTh == "氏名" {
				nameWithLineBreak := profileTr.MustElement("td").MustText()
				jobSeeker.SecretMemo += "氏名\n" + nameWithLineBreak + "\n\n"
				// 改行ごとにを分割
				nameWithSpace := utility.RegexpForLineBreak.Split(nameWithLineBreak, -1)

				// 空白ごとにを分割
				furiganaSplited := strings.Split(nameWithSpace[0], " ")
				if len(furiganaSplited) < 2 {
					log.Println("furiganaSplited is not 2")
					continue
				}
				jobSeeker.LastFurigana = furiganaSplited[0]
				jobSeeker.FirstFurigana = furiganaSplited[1]

				// 1987年8月15日生まれ（35歳）
				// ※年齢は、レジュメを表示するたびに生年月日から再計算しています。
			} else if profileTh == "生年月日" {
				birthdayWithLineBreak := profileTr.MustElement("td").MustText()
				jobSeeker.SecretMemo += "生年月日\n" + birthdayWithLineBreak + "\n\n"
				// 改行ごとにを分割
				birthdatWithSpace := utility.RegexpForLineBreak.Split(birthdayWithLineBreak, -1)
				if len(birthdatWithSpace) != 2 {
					continue
				}

				yearArr := strings.Split(birthdatWithSpace[0], "年")
				if len(yearArr) < 2 {
					log.Println("yearArrの長さが2未満です", yearArr)
					continue
				}
				monthArr := strings.Split(yearArr[1], "月")
				if len(monthArr) < 2 {
					log.Println("monthArrの長さが2未満です", monthArr)
					continue
				}
				monthInt, err := strconv.Atoi(monthArr[0])
				if err != nil {
					log.Println("monthIntへの変換に失敗しました", monthArr[0])
					continue
				}
				if monthInt < 10 {
					monthArr[0] = "0" + monthArr[0]
				}
				dayArr := strings.Split(monthArr[1], "日")
				if len(dayArr) < 2 {
					log.Println("dayArrの長さが2未満です", dayArr)
					continue
				}
				dayInt, err := strconv.Atoi(dayArr[0])
				if err != nil {
					log.Println("dayIntへの変換に失敗しました", dayArr[0])
					continue
				}
				if dayInt < 10 {
					dayArr[0] = "0" + dayArr[0]
				}
				birthday := yearArr[0] + "-" + monthArr[0] + "-" + dayArr[0]

				log.Println("birthday:", birthday)
				jobSeeker.Birthday = birthday

			} else if profileTh == "メールアドレス" {
				email := profileTr.MustElement("td").MustText()
				jobSeeker.SecretMemo += "メールアドレス\n" + email + "\n\n"
				log.Println("email:", email)
				jobSeeker.Email = email
				// テスト用メールアドレス
				// if i.app.Env != "prd" {
				// 	jobSeeker.Email = ""
				// }

			} else if profileTh == "電話番号" {
				phoneNumber := profileTr.MustElement("td").MustText()
				jobSeeker.SecretMemo += "電話番号\n" + phoneNumber + "\n\n"
				log.Println("phoneNumber:", phoneNumber)
				jobSeeker.PhoneNumber = phoneNumber

			} else if profileTh == "住所" {
				address := profileTr.MustElement("td").MustText()
				jobSeeker.SecretMemo += "住所\n" + address + "\n\n"
				log.Println("address:", address)
				jobSeeker.Address = address

				// 以下の都道府県のみを取得
				// 〒113-0022
				// 東京都文京区千駄木

				// 改行ごとにを分割
				addressSplited := utility.RegexpForLineBreak.Split(address, -1)
				if len(addressSplited) < 2 || len(addressSplited[1]) < 12 {
					log.Println("addressSplitedの長さが2未満です", addressSplited)
					continue
				}

				// 都道府県のみを取得
				inputPrefecture := addressSplited[1][:12]
				for i, prefectureName := range entity.Prefecture {
					if strings.Contains(inputPrefecture, prefectureName) {
						log.Println("prefectureName:", prefectureName)
						jobSeeker.Prefecture = null.NewInt(int64(i), true)
						break
					}
				}

			}
		}

		educationMemo := "<<学歴・語学・資格>>\n"
		for _, educationTr := range educationTrs {
			var (
				educationThText string
				educationTdText string
			)
			if educationTr.MustHas("th") {
				educationThText = educationTr.MustElement("th").
					MustText()
			} else {
				log.Println("educationTableにthがありません")
				continue
			}
			if educationTr.MustHas("td") {
				educationTdText = educationTr.MustElement("td").
					MustText()
				educationMemo += "・" + educationThText + "\n" + educationTdText + "\n\n"

				// 語学テーブル
				// if educationThText == "英語" {
				// 	var languageLevel null.Int
				// 	if strings.Contains(educationTdText, "ビジネス会話レベル") || strings.Contains(educationTdText, "ネイティブレベル") {
				// 		languageLevel = null.NewInt(0, true)
				// 	} else if strings.Contains(educationTdText, "日常会話レベル") {
				// 		languageLevel = null.NewInt(1, true)
				// 	}

				// 	languageSkill := entity.JobSeekerLanguageSkill{
				// 		LanguageType:  null.NewInt(0, true),
				// 		LanguageLevel: languageLevel,
				// 	}
				// 	jobSeeker.LanguageSkills = append(jobSeeker.LanguageSkills, languageSkill)
				// } else if educationThText == "北京語" {
				// 	var languageLevel null.Int
				// 	if strings.Contains(educationTdText, "ビジネス会話レベル") || strings.Contains(educationTdText, "ネイティブレベル") {
				// 		languageLevel = null.NewInt(0, true)
				// 	} else if strings.Contains(educationTdText, "日常会話レベル") {
				// 		languageLevel = null.NewInt(1, true)
				// 	}

				// 	languageSkill := entity.JobSeekerLanguageSkill{
				// 		LanguageType:  null.NewInt(6, true),
				// 		LanguageLevel: languageLevel,
				// 	}
				// 	jobSeeker.LanguageSkills = append(jobSeeker.LanguageSkills, languageSkill)
				// }

				// 資格
				// if educationThText == "資格" {
				// 	educationTdText = strings.ReplaceAll(educationTdText, "・", "")
				// 	licenseStrListSplitedByLineBreak := utility.RegexpForLineBreak.Split(educationTdText, -1)
				// 	for _, licenseStr := range licenseStrListSplitedByLineBreak {
				// 		licenseInt := entity.GetIntLicenseType(licenseStr)
				// 		if licenseInt.Valid {
				// 			license := entity.JobSeekerLicense{
				// 				LicenseType: licenseInt,
				// 			}
				// 			jobSeeker.Licenses = append(jobSeeker.Licenses, license)
				// 		}
				// 	}
				// }

			} else {
				educationMemo += "【" + educationThText + "】\n\n"
				continue
			}
		}
		jobSeeker.SecretMemo += educationMemo

		experienceMemo := "<<職務経歴>>\n"
		for _, experienceTr := range experienceTrs {
			var (
				experienceThText string
				experienceTdText string
			)
			if experienceTr.MustHas("th") {
				experienceThText = experienceTr.MustElement("th").
					MustText()
			} else {
				log.Println("experienceTableにthがありません")
				continue
			}
			if experienceTr.MustHas("td") {
				experienceTdText = experienceTr.MustElement("td").
					MustText()
				experienceMemo += "・" + experienceThText + "\n" + experienceTdText + "\n\n"

				if experienceThText == "現在または直前の年収" {
					annualIncomeStr := strings.Replace(experienceTdText, "万円", "", -1)
					annualIncomeStr = strings.Replace(annualIncomeStr, "約", "", -1)
					annualIncomeStr = strings.TrimSpace(annualIncomeStr)
					annualIncomeInt, err := strconv.Atoi(annualIncomeStr)
					if err != nil {
						log.Println("年収の変換に失敗しました", annualIncomeStr)
						continue
					}
					jobSeeker.AnnualIncome = null.NewInt(int64(annualIncomeInt), true)
				}

			} else {
				experienceMemo += "【" + experienceThText + "】\n\n"
				continue
			}

		}
		jobSeeker.SecretMemo += experienceMemo

		jobSeekerList = append(jobSeekerList, &jobSeeker)
		detailPage.MustClose()
	}

	// 求職者をDBに保存
	for _, jobSeeker := range jobSeekerList {
		jobSeeker.AgentID = input.AgentID
		jobSeeker.AgentStaffID = null.NewInt(int64(input.ScoutService.AgentStaffID), true)
		jobSeeker.Phase = null.NewInt(int64(entity.EntryInterview), true) // フェーズ： エントリー
		jobSeeker.RegisterPhase = null.NewInt(1, true)                    // 登録状況:下書き
		jobSeeker.InterviewDate = time.Now().UTC()
		jobSeeker.InflowChannelID = input.ScoutService.InflowChannelID

		err := i.jobSeekerRepository.Create(jobSeeker)
		if err != nil {
			log.Println(err)
			return output, err
		}

		// 開発環境の場合はユーザー名などを統一する
		// if os.Getenv("APP_ENV") != "prd" {
		// 	err = i.jobSeekerRepository.UpdateForDev(jobSeeker.ID)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return output, err
		// 	}
		// }

		jobSeekerDocument := entity.NewJobSeekerDocument(
			jobSeeker.ID,
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		)

		err = i.jobSeekerDocumentRepository.Create(jobSeekerDocument)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求職者作成時にメッセージグループ作成
		// エージェントと求職者のチャットグループを作成
		chatGroup := entity.NewChatGroupWithJobSeeker(
			jobSeeker.AgentID,
			jobSeeker.ID,
			false, // 初めはLINE連携してないから false
		)

		err = i.chatGroupWithJobSeekerRepository.Create(chatGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 面談実施待ちより前のフェーズの場合は「面談調整タスク」を作成
		var (
			phaseSub null.Int
			date     string
		)

		phaseSub = null.NewInt(0, true) //日程調整依頼

		// タスクの期限は当日で登録
		now := time.Now()
		date = now.Format("2006-01-02")

		// 面談調整タスクの作成
		interviewTaskGroup := entity.NewInterviewTaskGroup(
			jobSeeker.AgentID,
			jobSeeker.ID,
			jobSeeker.InterviewDate,
			utility.EarliestTime(), // 初期値,
		)

		err = i.interviewTaskGroupRepository.Create(interviewTaskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		interviewTask := entity.NewInterviewTask(
			interviewTaskGroup.ID,
			null.NewInt(0, false),
			null.NewInt(0, false),
			jobSeeker.Phase,
			phaseSub,
			"",
			date,
			null.NewInt(99, true),
			getStrPhaseForJobSeeker(jobSeeker.Phase),
		)

		err = i.interviewTaskRepository.Create(interviewTask)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	log.Println("処理完了。jobSeekerList:", jobSeekerList)
	output.JobSeekerList = jobSeekerList

	return output, nil
}

/*
		マイナビから新規応募者を取得する

	  開発環境でかかった時間 5分(5名登録):
		2023/04/14 16:01:00 マイナビのエントリー取得を開始します
		30分経過でタイムアウト
*/
type EntryOnMynaviScoutingInput struct {
	AgentID      uint
	ScoutService *entity.ScoutService
	Context      context.Context

	UserIDList []string
}

type EntryOnMynaviScoutingOutput struct {
	JobSeekerList []*entity.JobSeeker
}

func (i *ScoutServiceInteractorImpl) EntryOnMynaviScouting(input EntryOnMynaviScoutingInput) (EntryOnMynaviScoutingOutput, error) {
	var (
		output        EntryOnMynaviScoutingOutput
		err           error
		errMessage    string
		jobSeekerList []*entity.JobSeeker = []*entity.JobSeeker{}
		browser       *rod.Browser
		page          *rod.Page
	)

	if input.ScoutService.LoginID == "" || input.ScoutService.Password == "" {
		errMessage = "ScoutServiceの認証情報が未入力です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	} else if input.Context == nil {
		errMessage = "Contextが空です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	} else if input.AgentID == 0 {
		errMessage = "AgentIDが未入力です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// 2日以内に登録されたの応募者を取得
	jobSeekerListInDb, err := i.jobSeekerRepository.GetByAgentIDWithinTwoDays(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// ホストOSのchrome, chromium or microsoft edgeのパスを取得
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
	} else {
		l.Headless(true)
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

	// ブラウザ閉じる
	defer browser.MustClose()

	// browserにタイムアウトを設定
	browser, cancel := browser.
		Context(input.Context).
		WithCancel()
	defer cancel()

	// ログインページへ遷移
	fmt.Println("------------------\n\nブラウザOpen\n\n------------------")
	page = browser.
		MustPage("https://scouting.mynavi.jp/client/").
		MustWaitLoad()

	// ログイン情報を入力
	loginIDInput := page.MustElement("input[name=clientMailAddress]")
	loginIDInput.MustInput(input.ScoutService.LoginID)

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
	time.Sleep(6 * time.Second)

	// ログアウト
	defer page.MustNavigate("https://scouting.mynavi.jp/client/login/logout")

	// マイナビの場合は、ログイン失敗時、再度ログインする
	if page.MustHas("input[name=clientMailAddress]") {
		log.Println("ログイン失敗1回目、再度ログインします")
		page.
			MustElement("input[name=clientPassword]").
			MustType(rodInput.Enter)

		page.WaitLoad()
		time.Sleep(6 * time.Second)
	}

	// 再度ログイン時も失敗した場合は、エラーを返す
	if page.MustHas("input[name=clientMailAddress]") {
		errMessage = "ログイン失敗しました。IDとパスワードが間違っているか。すでに利用しているユーザーがいる可能性があります。"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	log.Println("ログイン成功")

	var index = 0
	for {
		// エントリー管理ページへ遷移
		page.MustNavigate("https://scouting.mynavi.jp/client/applicationTop/search")
		page.WaitLoad()
		time.Sleep(10 * time.Second)

		// エージェントサーチ経由とスカウト経由の応募者両方を見るために、2回ループする
		searchTargetsOnSearch := page.MustElements("input[name=searchTarget]")
		if len(searchTargetsOnSearch) != 2 {
			errMessage = "スカウト経由の応募者を見るための要素が見つかりません"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		// 検索対象で「マイナビスカウティング応募」を選択
		searchTargetsOnSearch[1].MustClick()

		// マイナビの条件検索はIDのOR検索ができるためIDをカンマ区切りにして検索
		var userIDListStr = strings.Join(input.UserIDList, ", ")
		page.MustElement("input[name=memberId]").MustInput(userIDListStr)

		searchSubmitOnSearch := page.MustElement("a.btn.sizeLM.blue")
		if searchSubmitOnSearch == nil || searchSubmitOnSearch.MustText() != "検索する" {
			errMessage = "マイナビスカウティング経由の応募者検索する送信要素が見つかりません"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}
		searchSubmitOnSearch.MustClick()
		page.WaitLoad()
		time.Sleep(10 * time.Second)

		usersOnSearch := page.MustElements("section.user-box")
		if len(usersOnSearch) == 0 {
			errMessage = "応募者が見つかりません"
			log.Println(errMessage)
			return output, errors.New(errMessage)
		}

		// 最終ループを終えたら
		if len(input.UserIDList) == index {
			break
		}

		userID := input.UserIDList[index]

		userSendingButtonSelecter := fmt.Sprintf("section#scroll-%v > div.user-box-inner > div.user-foot > div.floatR > ul.inline-box > li > a.btn.sizeML", userID)

		sendingButton, err := page.Element(userSendingButtonSelecter)
		if err != nil {
			log.Println("初回メッセージを送信するボタンが見つかりませんでした。")
			log.Println(err)
			return output, err
		}

		var buttonText = sendingButton.MustText()

		// 初回メッセージを送信ボタンが有効な場合はメッセージ送信ページに遷移、そうで無い場合は次のループへ
		if strings.Contains(buttonText, "初回メッセージを送信する") {
			messagePagePath, err := sendingButton.Attribute("href")
			if err != nil {
				log.Println("初回メッセージ送信ページのリンクの取得に失敗しました。")
				log.Println(err)
				return output, err
			}

			messagePageLink := fmt.Sprintf("https://scouting.mynavi.jp%v", *messagePagePath)
			page.MustNavigate(messagePageLink).WaitLoad()

			time.Sleep(5 * time.Second)
		} else if strings.Contains(buttonText, "初回メッセージ送信済") {
			log.Println("初回メッセージ送信済み")
			continue
		}

		fmt.Println("メッセージ送信ページに遷移")

		// ローカル環境ではメッセージ送信しない
		if i.app.BatchType == "entry" {
			// メッセージを作成する
			messagePageLinks, err := page.Elements("a.btn.blue")
			if err != nil {
				log.Println(err)
				return output, err
			}

			isMatchedWithMessagePageLink := false
			for _, messagePageLink := range messagePageLinks {
				if strings.Contains(messagePageLink.MustText(), "メッセージを作成する") {
					isMatchedWithMessagePageLink = true
					messagePageLink.MustClick()
					page.WaitLoad()
					time.Sleep(5 * time.Second)
					break
				}
			}
			if !isMatchedWithMessagePageLink {
				errMessage = "メッセージを作成するリンクが見つかりません"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			log.Println("メッセージを作成するページに遷移しました")

			// templateUser
			// templateUserSelect, err := page.Element("select[name=templateUser]")
			// if err != nil {
			// 	log.Println(err)
			// 	return output, err
			// }

			// templateUser := "濱田　僚太"
			// err = templateUserSelect.Select(
			// 	[]string{templateUser},
			// 	true,
			// 	"text",
			// )
			// if err != nil {
			// 	log.Println(err)
			// 	return output, err
			// }
			// time.Sleep(5 * time.Second)

			// select-scout-templete　【面談調整】テンプレ
			scoutTemplateSelect, err := page.Element("select.select-scout-templete")
			if err != nil {
				log.Println(err)
				return output, err
			}

			err = scoutTemplateSelect.Select(
				[]string{input.ScoutService.TemplateTitleForEmployed},
				true,
				"text",
			)
			if err != nil {
				log.Println(err)
				return output, err
			}

			setTemplateLinks, err := page.Elements("a.set-templete")
			if err != nil {
				log.Println(err)
				return output, err
			}

			isMatchedWithSetTemplateLink := false
			for _, setTemplateLink := range setTemplateLinks {
				if strings.Contains(setTemplateLink.MustText(), "テンプレートを読み込む") {
					isMatchedWithSetTemplateLink = true
					setTemplateLink.MustClick()
					break
				}
			}
			if !isMatchedWithSetTemplateLink {
				errMessage = "テンプレートを読み込むリンクが見つかりません"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			confirmLinks, err := page.Elements("a.btn.sizeMM")
			if err != nil {
				log.Println(err)
				return output, err
			}

			isMatchedWithConfirmLink := false
			for _, confirmLink := range confirmLinks {
				if strings.Contains(confirmLink.MustText(), "確認する") {
					isMatchedWithConfirmLink = true
					confirmLink.MustClick()
					time.Sleep(5 * time.Second)
					break
				}
			}
			if !isMatchedWithConfirmLink {
				errMessage = "確認するリンクが見つかりません"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			sendLinks, err := page.Elements("a.btn.blue")
			if err != nil {
				log.Println(err)
				return output, err
			}

			isMatchedWithSendLink := false
			for _, sendLink := range sendLinks {
				if strings.Contains(sendLink.MustText(), "送信する") {
					isMatchedWithSendLink = true
					if i.app.BatchType == "entry" {
						sendLink.MustClick()
						time.Sleep(10 * time.Second)
					}
					log.Println("メッセージを送信しました")
					break
				}
			}
			if !isMatchedWithSendLink {
				errMessage = "送信するリンクが見つかりません"
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}
		}

		index++
	}

	// csvダウンロード
	csvFile := browser.MustWaitDownload()

	downloadCSVBtn, err := page.Element("div.action-btn-box > div.floatR > a.btn.white")
	if err != nil {
		log.Println("CSVダウンロードボタンが見つかりません")
		log.Println(err)
		return output, err
	}

	log.Println("CSVダウンロードボタンを見つけました。")

	err = downloadCSVBtn.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		log.Println("CSVダウンロードボタンのクリックで失敗しました。")
		log.Println(err)
		return output, err
	}

	ouputFileName := fmt.Sprintf("getentry_mynavi_%v_%v.csv", input.ScoutService.AgentRobotID, time.Now().Format("20060102150405"))
	err = utils.OutputFile(ouputFileName, csvFile())
	if err != nil {
		log.Println(err)
		return output, err
	}

	// csvファイルを読み込み
	file, err := os.Open(ouputFileName)
	if err != nil {
		log.Println(err)
		return output, err
	}
	log.Println("csvファイルを開く file:", file.Name(), *file)

	// ShiftJISのCSVファイルを文字化けしない様に読み込む
	reader := csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))
	log.Println("csvファイルを読み込み reader:", reader)

	reader.TrimLeadingSpace = true // true の場合は、先頭の空白文字を無視する
	reader.ReuseRecord = true      // true の場合は、Read で戻ってくるスライスを次回再利用する。パフォーマンスが上がる

	// 項目説明のレコードを読み込み
	reader.Read()

recordLoop:
	for {
		record, err := reader.Read()
		// 空行があった場合は終了
		log.Println("record: ", record)
		if err == io.EOF || record[0] == "" {
			log.Println("行が空のため、ループを終了します")
			break
		}

		jobSeeker := entity.JobSeeker{}
		jobSeeker.SecretMemo = "・エントリー媒体\nマイナビスカウティング\n\n"
		isMatchedUserCount := 0

		for columnI, column := range record {
			log.Printf("column[%v]: %v", columnI, column)
			if column == "" {
				continue
			}

			switch columnI {
			// 0:最終応募経路_名
			case 0:
				jobSeeker.SecretMemo += "・最終応募経路\n" + column + "\n\n"
				// 1:最終応募経路_コード
				// 2:最終応募日時
				// 3:エントリー求人管理No.
				// 4:進捗状況_名
			case 4:
				jobSeeker.SecretMemo += "進捗状況: " + column + "\n\n"
			// 5:進捗状況_コード
			// 6:応募メッセージ
			case 7: // 会員ＮＯ
				jobSeeker.ExternalID = column
			// 8:姓
			case 8:
				jobSeeker.LastName = column
			// 9:名
			case 9:
				jobSeeker.FirstName = column
			// 10:姓カナ
			case 10:
				jobSeeker.LastFurigana = column
			//11:名カナ
			case 11:
				jobSeeker.FirstFurigana = column

				// 「名前」がDB登録済みの求職者と重複している場合はスキップ
				for _, jobSeekerInDb := range jobSeekerListInDb {
					if jobSeekerInDb.LastName == jobSeeker.LastName &&
						jobSeekerInDb.FirstName == jobSeeker.FirstName &&
						jobSeekerInDb.LastFurigana == jobSeeker.LastFurigana &&
						jobSeekerInDb.FirstFurigana == jobSeeker.FirstFurigana {
						fmt.Println("2日以内に登録した求職者と重複しています。", jobSeeker.LastName+jobSeeker.FirstName, jobSeeker.LastFurigana+jobSeeker.FirstFurigana)
						continue recordLoop
					}
				}

			// 12:郵便番号
			case 12:
				jobSeeker.PostCode = column
				jobSeeker.Address = column
			//13:都道府県
			case 13:
				jobSeeker.Prefecture = entity.GetIntPrefecture(column)
				jobSeeker.Address += column
			// 14:市区町村
			case 14:
				jobSeeker.Address += column
				//15:住所
			case 15:
				jobSeeker.Address += column
				//16:電話番号
			case 16:
				jobSeeker.SecretMemo += "・電話番号\n" + column + "\n\n"
				//17:携帯番号
			case 17:
				jobSeeker.PhoneNumber = column
				jobSeeker.SecretMemo += "・携帯番号\n" + column + "\n\n"
				//18:E-MAIL
			case 18:
				jobSeeker.Email = column

				// 多重で応募している場合、同一候補者が重複してcsvに出力されるため、名前、フリガナ、メールアドレスで重複をチェックする
				for _, jobSeekerFromList := range jobSeekerList {
					if jobSeekerFromList.FirstName == jobSeeker.FirstName &&
						jobSeekerFromList.LastName == jobSeeker.LastName &&
						jobSeekerFromList.FirstFurigana == jobSeeker.FirstFurigana &&
						jobSeekerFromList.LastFurigana == jobSeeker.LastFurigana &&
						jobSeekerFromList.Email == jobSeeker.Email {
						log.Println("リストに登録済み求職者と重複しているためスキップします。", jobSeeker.LastName, jobSeeker.FirstName, jobSeeker.FirstFurigana, jobSeeker.LastFurigana, jobSeeker.Email)
						continue recordLoop
					}
				}
				//19:生年月日
			case 19:
				// 2000/01/01 -> 2000-01-01
				birthday := strings.Replace(column, "/", "-", -1)
				jobSeeker.Birthday = birthday

				//20:年齢
				//21:性別
			case 21:
				if column == "m" {
					jobSeeker.Gender = null.NewInt(0, true)
					jobSeeker.SecretMemo += "・性別\n" + "男性" + "\n\n"
				} else if column == "w" {
					jobSeeker.Gender = null.NewInt(1, true)
					jobSeeker.SecretMemo += "・性別\n" + "女性" + "\n\n"
				}
				//22:経験職種大分類_名１
				//23:経験職種大分類_コード１
				//24:経験職種中分類_名１
				//25:経験職種中分類_コード１
			//26:経験職種小分類_名１
			case 26:
				jobSeeker.SecretMemo += "・経験職種:" + column
			//27:経験職種小分類_コード１
			//28:経験職種経験年数_名１
			case 28:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
				//29:経験職種経験年数_コード１
				//30:経験職種大分類_名２
				//31:経験職種大分類_コード２
				//32:経験職種中分類_名２
				//33:経験職種中分類_コード２
				//34:経験職種小分類_名２
			case 34:
				jobSeeker.SecretMemo += "・経験職種:" + column
				//35:経験職種小分類_コード２
				//36:経験職種経験年数_名２
			case 36:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
				//37:経験職種経験年数_コード２
				//38:経験職種大分類_名３
				//39:経験職種大分類_コード３
				//40:経験職種中分類_名３
				//41:経験職種中分類_コード３
				//42:経験職種小分類_名３
			case 42:
				jobSeeker.SecretMemo += "・経験職種:" + column
				//43:経験職種小分類_コード３
				//44:経験職種経験年数_名３
			case 44:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
				//45:経験職種経験年数_コード３
				//46:経験職種大分類_名４
				//47:経験職種大分類_コード４
				//48:経験職種中分類_名４
				//49:経験職種中分類_コード４
				//50:経験職種小分類_名４
			case 50:
				jobSeeker.SecretMemo += "・経験職種:" + column
				//51:経験職種小分類_コード４
				//52:経験職種経験年数_名４
			case 52:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
				//53:経験職種経験年数_コード４
				//54:経験職種大分類_名５
				//55:経験職種大分類_コード５
				//56:経験職種中分類_名５
				//57:経験職種中分類_コード５
				//58:経験職種小分類_名５
			case 58:
				jobSeeker.SecretMemo += "・経験職種:" + column
				//59:経験職種小分類_コード５
				//60:経験職種経験年数_名５
			case 60:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
				//61:経験職種経験年数_コード５
				//62:経験職種大分類_名６
				//63:経験職種大分類_コード６
				//64:経験職種中分類_名６
				//65:経験職種中分類_コード６
				//66:経験職種小分類_名６
			case 66:
				jobSeeker.SecretMemo += "・経験職種:" + column
				//67:経験職種小分類_コード６
				//68:経験職種経験年数_名６
			case 68:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
				//69:経験職種経験年数_コード６
				//70:経験職種大分類_名７
				//71:経験職種大分類_コード７
				//72:経験職種中分類_名７
				//73:経験職種中分類_コード７
				//74:経験職種小分類_名７
			case 74:
				jobSeeker.SecretMemo += "・経験職種:" + column
				//75:経験職種小分類_コード７
				//76:経験職種経験年数_名７
			case 76:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
				//77:経験職種経験年数_コード７
				//78:経験職種大分類_名８
				//79:経験職種大分類_コード８
				//80:経験職種中分類_名８
				//81:経験職種中分類_コード８
				//82:経験職種小分類_名８
			case 82:
				jobSeeker.SecretMemo += "・経験職種:" + column
				//83:経験職種小分類_コード８
				//84:経験職種経験年数_名８
			case 84:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
				//85:経験職種経験年数_コード８
				//86:経験職種大分類_名９
				//87:経験職種大分類_コード９
				//88:経験職種中分類_名９
				//89:経験職種中分類_コード９
				//90:経験職種小分類_名９
			case 90:
				jobSeeker.SecretMemo += "・経験職種:" + column
				//91:経験職種小分類_コード９
				//92:経験職種経験年数_名９
			case 92:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
				//93:経験職種経験年数_コード９
				//94:経験職種大分類_名１０
				//95:経験職種大分類_コード１０
				//96:経験職種中分類_名１０
				//97:経験職種中分類_コード１０
				//98:経験職種小分類_名１０
			case 98:
				jobSeeker.SecretMemo += "・経験職種:" + column
				//99:経験職種小分類_コード１０
				//100:経験職種経験年数_名１０
			case 100:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
			//101:経験職種経験年数_コード１０
			//102:最終学歴（最終）
			//103:学校名（最終）
			case 103:
				jobSeeker.SecretMemo += "・学校名(最終):" + column + "\n"
			//104:区分（最終）_名
			case 104:
				jobSeeker.SecretMemo += "・区分(最終):" + column + "\n"
			//105:区分（最終）_コード
			//106:学部学科名（最終）
			case 106:
				jobSeeker.SecretMemo += "・学部学科名(最終):" + column + "\n"
			//107:学部学科系統（最終）_名
			case 107:
				jobSeeker.SecretMemo += "・学部学科系統(最終):" + column + "\n"
			//108:学部学科系統（最終）_コード
			//109:在籍期間FROM（最終）_名
			case 109:
				jobSeeker.SecretMemo += "・在籍開始年(最終):" + column + "\n"
			//110:在籍期間TO（最終）_名
			case 110:
				jobSeeker.SecretMemo += "・在籍終了年(最終):" + column + "\n"
			//111:卒業状況（最終）_名
			case 111:
				jobSeeker.SecretMemo += "・卒業状況(最終):" + column + "\n\n"
			//112:卒業状況（最終）_コード
			//113:学校名１
			//114:区分１_名
			//115:区分１_コード
			//116:学部学科名１
			//117:学部学科系統１_名
			//118:学部学科系統１_コード
			//119:在籍期間FROM１_名
			//120:在籍期間TO１_名
			//121:卒業状況１_名
			//122:卒業状況１_コード
			//123:学校名２
			//124:区分２_名
			//125:区分２_コード
			//126:学部学科名２
			//127:学部学科系統２_名
			//128:学部学科系統２_コード
			//129:在籍期間FROM２_名
			//130:在籍期間TO２_名
			//131:卒業状況２_名
			//132:卒業状況２_コード
			//133:学校名３
			//134:区分３_名
			//135:区分３_コード
			//136:学部学科名３
			//137:学部学科系統３_名
			//138:学部学科系統３_コード
			//139:在籍期間FROM３_名
			//140:在籍期間TO３_名
			//141:卒業状況３_名
			//142:卒業状況３_コード
			//143:学校名４
			//144:区分４_名
			//145:区分４_コード
			//146:学部学科名４
			//147:学部学科系統４_名
			//148:学部学科系統４_コード
			//149:在籍期間FROM４_名
			//150:在籍期間TO４_名
			//151:卒業状況４_名
			//152:卒業状況４_コード
			//153:学校名５
			//154:区分５_名
			//155:区分５_コード
			//156:学部学科名５
			//157:学部学科系統５_名
			//158:学部学科系統５_コード
			//159:在籍期間FROM５_名
			//160:在籍期間TO５_名
			//161:卒業状況５_名
			//162:卒業状況５_コード
			//163:現在の就業状況_名
			case 163:
				jobSeeker.SecretMemo += "・現在の就業状況:" + column + "\n\n"
				if column == "離職中 " {
					jobSeeker.StateOfEmployment = null.NewInt(1, true)
				} else {
					jobSeeker.StateOfEmployment = null.NewInt(0, true)
				}
			//164:現在の就業状況_コード
			//165:年収実績_名
			case 165:
				jobSeeker.SecretMemo += "・年収実績:" + column + "\n\n"
				annualIncomeInt := 0
				if column == "199万円以下" {
					annualIncomeInt = 199
				} else {
					// 300〜349万円 → 300
					annualIncomeInt, err = strconv.Atoi(strings.Split(column, "～")[0])
					if err != nil {
						log.Println("直近の給与の変換に失敗しました", err)
						continue
					}
				}
				jobSeeker.AnnualIncome = null.NewInt(int64(annualIncomeInt), true)
			//166:年収実績_コード
			//167:配偶者_名
			case 167:
				jobSeeker.SecretMemo += "・配偶者:" + column + "\n\n"
				if column == "配偶者あり" {
					jobSeeker.Spouse = null.NewInt(0, true)
				} else {
					jobSeeker.Spouse = null.NewInt(1, true)
				}
			//168:配偶者_コード
			//169:希望職種大分類_名１
			//170:希望職種大分類_コード１
			//171:希望職種中分類_名１
			//172:希望職種中分類_コード１
			//173:希望職種小分類_名１
			case 173:
				jobSeeker.SecretMemo += "・希望職種:" + column + "\n\n"
				//174:希望職種小分類_コード１
				//175:希望職種大分類_名２
				//176:希望職種大分類_コード２
				//177:希望職種中分類_名２
				//178:希望職種中分類_コード２
				//179:希望職種小分類_名２
			case 179:
				jobSeeker.SecretMemo += "・希望職種:" + column + "\n"
				//180:希望職種小分類_コード２
				//181:希望職種大分類_名３
				//182:希望職種大分類_コード３
				//183:希望職種中分類_名３
				//184:希望職種中分類_コード３
				//185:希望職種小分類_名３
			case 185:
				jobSeeker.SecretMemo += "・希望職種:" + column + "\n"
				//186:希望職種小分類_コード３
				//187:希望職種大分類_名４
				//188:希望職種大分類_コード４
				//189:希望職種中分類_名４
				//190:希望職種中分類_コード４
				//191:希望職種小分類_名４
			case 191:
				jobSeeker.SecretMemo += "・希望職種:" + column + "\n"
				//192:希望職種小分類_コード４
				//193:希望職種大分類_名５
				//194:希望職種大分類_コード５
				//195:希望職種中分類_名５
				//196:希望職種中分類_コード５
				//197:希望職種小分類_名５
			case 197:
				jobSeeker.SecretMemo += "・希望職種:" + column + "\n"
				//198:希望職種小分類_コード５
				//199:希望業種大分類_名１
				//200:希望業種大分類_コード１
				//201:希望業種中分類_名１
				//202:希望業種中分類_コード１
				//203:希望業種小分類_名１
			case 203:
				jobSeeker.SecretMemo += "・希望職種:" + column + "\n"
				//204:希望業種小分類_コード１
				//205:希望業種大分類_名２
				//206:希望業種大分類_コード２
				//207:希望業種中分類_名２
				//208:希望業種中分類_コード２
				//209:希望業種小分類_名２
			case 209:
				jobSeeker.SecretMemo += "・希望職種:" + column
				//210:希望業種小分類_コード２
				//211:希望業種大分類_名３
				//212:希望業種大分類_コード３
				//213:希望業種中分類_名３
				//214:希望業種中分類_コード３
				//215:希望業種小分類_名３
			case 215:
				jobSeeker.SecretMemo += "・希望職種:" + column
				//216:希望業種小分類_コード３
				//217:希望業種大分類_名４
				//218:希望業種大分類_コード４
				//219:希望業種中分類_名４
				//220:希望業種中分類_コード４
				//221:希望業種小分類_名４
			case 221:
				jobSeeker.SecretMemo += "・希望職種:" + column
				//222:希望業種小分類_コード４
				//223:希望業種大分類_名５
				//224:希望業種大分類_コード５
				//225:希望業種中分類_名５
				//226:希望業種中分類_コード５
				//227:希望業種小分類_名５
			case 227:
				jobSeeker.SecretMemo += "・希望職種:" + column + "\n\n"
				//228:希望業種小分類_コード５
				//229:希望年収_名
			case 229:
				jobSeeker.SecretMemo += "\n\n・希望年収:" + column
				// 200万円以上 -> 200
				desiredAnnualIncomeInt, err := strconv.Atoi(strings.Split(column, "万")[0])
				if err != nil {
					log.Println(err)
					continue
				}
				jobSeeker.DesiredAnnualIncome = null.NewInt(int64(desiredAnnualIncomeInt), true)
				//230:希望年収_コード
				//231:希望勤務地_名１
			case 231:
				jobSeeker.SecretMemo += "\n\n・希望勤務地:" + column
				desiredWorkLocationInt := entity.GetIntPrefecture(column)
				if desiredWorkLocationInt.Valid {
					jobSeeker.DesiredWorkLocations = append(
						jobSeeker.DesiredWorkLocations,
						entity.JobSeekerDesiredWorkLocation{
							JobSeekerID:         jobSeeker.ID,
							DesiredWorkLocation: desiredWorkLocationInt,
							DesiredRank:         null.NewInt(1, true),
						},
					)
				}

			//232:希望勤務地_コード１
			//233:希望勤務地_名２
			case 233:
				jobSeeker.SecretMemo += "\n・希望勤務地:" + column
				desiredWorkLocationInt := entity.GetIntPrefecture(column)
				if desiredWorkLocationInt.Valid {
					jobSeeker.DesiredWorkLocations = append(
						jobSeeker.DesiredWorkLocations,
						entity.JobSeekerDesiredWorkLocation{
							JobSeekerID:         jobSeeker.ID,
							DesiredWorkLocation: desiredWorkLocationInt,
							DesiredRank:         null.NewInt(1, true),
						},
					)
				}
				//234:希望勤務地_コード２
				//235:希望勤務地_名３
			case 235:
				jobSeeker.SecretMemo += "\n・希望勤務地:" + column
				desiredWorkLocationInt := entity.GetIntPrefecture(column)
				if desiredWorkLocationInt.Valid {
					jobSeeker.DesiredWorkLocations = append(
						jobSeeker.DesiredWorkLocations,
						entity.JobSeekerDesiredWorkLocation{
							JobSeekerID:         jobSeeker.ID,
							DesiredWorkLocation: desiredWorkLocationInt,
							DesiredRank:         null.NewInt(1, true),
						},
					)
				}
				//236:希望勤務地_コード３
				//237:希望勤務地_名４
			case 237:
				jobSeeker.SecretMemo += "\n・希望勤務地:" + column
				desiredWorkLocationInt := entity.GetIntPrefecture(column)
				if desiredWorkLocationInt.Valid {
					jobSeeker.DesiredWorkLocations = append(
						jobSeeker.DesiredWorkLocations,
						entity.JobSeekerDesiredWorkLocation{
							JobSeekerID:         jobSeeker.ID,
							DesiredWorkLocation: desiredWorkLocationInt,
							DesiredRank:         null.NewInt(1, true),
						},
					)
				}
				//238:希望勤務地_コード４
				//239:希望勤務地_名５
			case 239:
				jobSeeker.SecretMemo += "\n・希望勤務地:" + column
				desiredWorkLocationInt := entity.GetIntPrefecture(column)
				if desiredWorkLocationInt.Valid {
					jobSeeker.DesiredWorkLocations = append(
						jobSeeker.DesiredWorkLocations,
						entity.JobSeekerDesiredWorkLocation{
							JobSeekerID:         jobSeeker.ID,
							DesiredWorkLocation: desiredWorkLocationInt,
							DesiredRank:         null.NewInt(1, true),
						},
					)
				}
				//240:希望勤務地_コード５
				//241:希望雇用形態_名
			case 241:
				jobSeeker.SecretMemo += "\n\n・希望雇用形態:" + column
				//242:希望雇用形態_コード
				//243:希望転職時期_名
			case 243:
				jobSeeker.SecretMemo += "\n\n・希望転職時期:" + column + "\n\n"
				// autoscout:即入社可能,1ヶ月以内,2ヶ月以内,3ヶ月以内,4ヶ月以内,5ヶ月以内,6ヶ月以内,7ヶ月〜
				// mynavi:すぐにでも,1カ月以内,3カ月以内,半年以内,1年以内
				if column == "すぐにでも" {
					jobSeeker.JoinCompanyPeriod = null.NewInt(0, true)
				} else if column == "1カ月以内" {
					jobSeeker.JoinCompanyPeriod = null.NewInt(1, true)
				} else if column == "3カ月以内" {
					jobSeeker.JoinCompanyPeriod = null.NewInt(3, true)
				} else if column == "半年以内" {
					jobSeeker.JoinCompanyPeriod = null.NewInt(6, true)
				} else if column == "1年以内" {
					jobSeeker.JoinCompanyPeriod = null.NewInt(7, true)
				}
				//244:希望転職時期_コード
				//245:保有資格_名（１つのセルに、カンマ区切りで格納）
			case 245:
				jobSeeker.SecretMemo += "・保有資格:" + column + "\n\n"
				//246:保有資格_コード（１つのセルに、カンマ区切りで格納）
				//247:英語資格・スキル（ヒアリング・スピーキング）_名
			case 247:
				jobSeeker.SecretMemo += "・英語資格・スキル（ヒアリング・スピーキング）:" + column + "\n\n"
				//248:英語資格・スキル（ヒアリング・スピーキング）_コード
				//249:英語資格・スキル（TOEIC）_名
			case 249:
				jobSeeker.SecretMemo += "・英語資格・スキル（TOEIC）:" + column + "\n\n"
				//250:英語資格・スキル（TOEIC）_コード
				//251:英語資格・スキル（英検）_名
			case 251:
				jobSeeker.SecretMemo += "・英語資格・スキル（英検）:" + column + "\n\n"
				//252:英語資格・スキル（英検）_コード
				//253:その他の英語資格・スキル_名（１つのセルに、カンマ区切りで格納）
			case 253:
				jobSeeker.SecretMemo += "・その他の英語資格・スキル:" + column + "\n\n"
				//254:その他の英語資格・スキル_コード（１つのセルに、カンマ区切りで格納）
				//255:その他の言語資格・スキル_名（１つのセルに、カンマ区切りで格納）
			case 255:
				jobSeeker.SecretMemo += "・その他の言語資格・スキル:" + column + "\n\n"
				//256:その他の言語資格・スキル_コード（１つのセルに、カンマ区切りで格納）
				//257:ITキャリアシート_名（１つのセルに、カンマ区切りで格納）
			case 257:
				jobSeeker.SecretMemo += "・ITキャリアシート:" + column + "\n\n"
				//258:ITキャリアシート_コード（１つのセルに、カンマ区切りで格納）
				//259:経験年数_名（１つのセルに、カンマ区切りで格納）
			case 259:
				jobSeeker.SecretMemo += "・経験年数:" + column + "\n\n"
				//260:経験年数_コード（１つのセルに、カンマ区切りで格納）
				//261:経験社数_名
			case 261:
				jobSeeker.SecretMemo += "・経験社数:" + column + "\n\n"
			//262:経験社数_コード
			//263:在籍期間FROM１_名
			case 263:
				jobSeeker.SecretMemo += "・在籍開始年:" + column + "\n"
			//264:在籍期間TO１_名
			case 264:
				jobSeeker.SecretMemo += "・在籍終了年:" + column + "\n"
			//265:雇用形態１_名
			case 265:
				jobSeeker.SecretMemo += "・雇用形態:" + column + "\n"
			//266:雇用形態１_コード
			//267:勤務先名１
			case 267:
				jobSeeker.SecretMemo += "・勤務先名:" + column + "\n"
			//268:勤務先業種大分類_名１
			//269:勤務先業種大分類_コード１
			//270:勤務先業種中分類_名１
			//271:勤務先業種中分類_コード１
			//272:勤務先業種小分類_名１
			case 272:
				jobSeeker.SecretMemo += "・勤務先業種:" + column + "\n"
			//273:勤務先業種小分類_コード１
			//274:勤務先規模（資本金）１_名
			case 274:
				jobSeeker.SecretMemo += "・勤務先規模（資本金）:" + column + "\n"
			//275:勤務先規模（資本金）１_コード
			//276:勤務先規模（従業員数）１_名
			case 276:
				jobSeeker.SecretMemo += "・勤務先規模（従業員数）:" + column + "\n"
			//277:勤務先規模（従業員数）１_コード
			//278:マネジメント１_名
			case 278:
				jobSeeker.SecretMemo += "・マネジメント経験:" + column + "\n"
			//279:マネジメント１_コード
			//280:業務上のポジション１_名
			case 280:
				jobSeeker.SecretMemo += "・業務上のポジション:" + column + "\n"
			//281:業務上のポジション１_コード
			//282:職務内容１
			case 282:
				jobSeeker.SecretMemo += "・職務内容:" + column + "\n\n"
			//283:在籍期間FROM２_名
			case 283:
				jobSeeker.SecretMemo += "・在籍開始年:" + column + "\n"
			//284:在籍期間TO２_名
			case 284:
				jobSeeker.SecretMemo += "・在籍終了年:" + column + "\n"
			//285:雇用形態２_名
			case 285:
				jobSeeker.SecretMemo += "・雇用形態:" + column + "\n"
			//286:雇用形態２_コード
			//287:勤務先名２
			case 287:
				jobSeeker.SecretMemo += "・勤務先名:" + column + "\n"
			//288:勤務先業種大分類_名２
			//289:勤務先業種大分類_コード２
			//290:勤務先業種中分類_名２
			//291:勤務先業種中分類_コード２
			//292:勤務先業種小分類_名２
			case 292:
				jobSeeker.SecretMemo += "・勤務先業種:" + column + "\n"
			//293:勤務先業種小分類_コード２
			//294:勤務先規模（資本金）２_名
			case 294:
				jobSeeker.SecretMemo += "・勤務先規模（資本金）:" + column + "\n"
			//295:勤務先規模（資本金）２_コード
			//296:勤務先規模（従業員数）２_名
			case 296:
				jobSeeker.SecretMemo += "・勤務先規模（従業員数）:" + column + "\n"
			//297:勤務先規模（従業員数）２_コード
			//298:マネジメント２_名
			case 298:
				jobSeeker.SecretMemo += "・マネジメント経験:" + column + "\n"
			//299:マネジメント２_コード
			//300:業務上のポジション２_名
			case 300:
				jobSeeker.SecretMemo += "・業務上のポジション:" + column + "\n"
			//301:業務上のポジション２_コード
			//302:職務内容２
			case 302:
				jobSeeker.SecretMemo += "・職務内容:" + column + "\n\n"
			//303:在籍期間FROM３_名
			case 303:
				jobSeeker.SecretMemo += "・在籍開始年:" + column + "\n"
			//304:在籍期間TO３_名
			case 304:
				jobSeeker.SecretMemo += "・在籍終了年:" + column + "\n"
			//305:雇用形態３_名
			case 305:
				jobSeeker.SecretMemo += "・雇用形態:" + column + "\n"
			//306:雇用形態３_コード
			//307:勤務先名３
			case 307:
				jobSeeker.SecretMemo += "・勤務先名:" + column + "\n"
			//308:勤務先業種大分類_名３
			//309:勤務先業種大分類_コード３
			//310:勤務先業種中分類_名３
			//311:勤務先業種中分類_コード３
			//312:勤務先業種小分類_名３
			case 312:
				jobSeeker.SecretMemo += "・勤務先業種:" + column + "\n"
			//313:勤務先業種小分類_コード３
			//314:勤務先規模（資本金）３_名
			case 314:
				jobSeeker.SecretMemo += "・勤務先規模（資本金）:" + column + "\n"
			//315:勤務先規模（資本金）３_コード
			//316:勤務先規模（従業員数）３_名
			case 316:
				jobSeeker.SecretMemo += "・勤務先規模（従業員数）:" + column + "\n"
			//317:勤務先規模（従業員数）３_コード
			//318:マネジメント３_名
			case 318:
				jobSeeker.SecretMemo += "・マネジメント経験:" + column + "\n"
			//319:マネジメント３_コード
			//320:業務上のポジション３_名
			case 320:
				jobSeeker.SecretMemo += "・業務上のポジション:" + column + "\n"
			//321:業務上のポジション３_コード
			//322:職務内容３
			case 322:
				jobSeeker.SecretMemo += "・職務内容:" + column + "\n\n"
			//323:在籍期間FROM４_名
			case 323:
				jobSeeker.SecretMemo += "・在籍開始年:" + column + "\n"
			//324:在籍期間TO４_名
			case 324:
				jobSeeker.SecretMemo += "・在籍終了年:" + column + "\n"
			//325:雇用形態４_名
			case 325:
				jobSeeker.SecretMemo += "・雇用形態:" + column + "\n"
			//326:雇用形態４_コード
			//327:勤務先名４
			case 327:
				jobSeeker.SecretMemo += "・勤務先名:" + column + "\n"
			//328:勤務先業種大分類_名４
			//329:勤務先業種大分類_コード４
			//330:勤務先業種中分類_名４
			//331:勤務先業種中分類_コード４
			//332:勤務先業種小分類_名４
			case 332:
				jobSeeker.SecretMemo += "・勤務先業種:" + column + "\n"
			//333:勤務先業種小分類_コード４
			//334:勤務先規模（資本金）４_名
			case 334:
				jobSeeker.SecretMemo += "・勤務先規模（資本金）:" + column + "\n"
			//335:勤務先規模（資本金）４_コード
			//336:勤務先規模（従業員数）４_名
			case 336:
				jobSeeker.SecretMemo += "・勤務先規模（従業員数）:" + column + "\n"
			//337:勤務先規模（従業員数）４_コード
			//338:マネジメント４_名
			case 338:
				jobSeeker.SecretMemo += "・マネジメント経験:" + column + "\n"
			//339:マネジメント４_コード
			//340:業務上のポジション４_名
			case 340:
				jobSeeker.SecretMemo += "・業務上のポジション:" + column + "\n"
			//341:業務上のポジション４_コード
			//342:職務内容４
			case 342:
				jobSeeker.SecretMemo += "・職務内容:" + column + "\n\n"
			//343:在籍期間FROM５_名
			case 343:
				jobSeeker.SecretMemo += "・在籍開始年:" + column + "\n"
			//344:在籍期間TO５_名
			//345:雇用形態５_名
			case 345:
				jobSeeker.SecretMemo += "・雇用形態:" + column + "\n"
			//346:雇用形態５_コード
			//347:勤務先名５
			case 347:
				jobSeeker.SecretMemo += "・勤務先名:" + column + "\n"
			//348:勤務先業種大分類_名５
			//349:勤務先業種大分類_コード５
			//350:勤務先業種中分類_名５
			//351:勤務先業種中分類_コード５
			//352:勤務先業種小分類_名５
			case 352:
				jobSeeker.SecretMemo += "・勤務先業種:" + column + "\n"
			//353:勤務先業種小分類_コード５
			//354:勤務先規模（資本金）５_名
			case 354:
				jobSeeker.SecretMemo += "・勤務先規模（資本金）:" + column + "\n"
			//355:勤務先規模（資本金）５_コード
			//356:勤務先規模（従業員数）５_名
			case 356:
				jobSeeker.SecretMemo += "・勤務先規模（従業員数）:" + column + "\n"
			//357:勤務先規模（従業員数）５_コード
			//358:マネジメント５_名
			case 358:
				jobSeeker.SecretMemo += "・マネジメント経験:" + column + "\n"
			//359:マネジメント５_コード
			//360:業務上のポジション５_名
			case 360:
				jobSeeker.SecretMemo += "・業務上のポジション:" + column + "\n"
			//361:業務上のポジション５_コード
			//362:職務内容５
			case 362:
				jobSeeker.SecretMemo += "・職務内容:" + column + "\n\n"
			//363:在籍期間FROM６_名
			case 363:
				jobSeeker.SecretMemo += "・在籍開始年:" + column + "\n"
				//364:在籍期間TO６_名
			case 364:
				jobSeeker.SecretMemo += "・在籍終了年:" + column + "\n"
			//365:雇用形態６_名
			case 365:
				jobSeeker.SecretMemo += "・雇用形態:" + column + "\n"
			//366:雇用形態６_コード
			//367:勤務先名６
			case 367:
				jobSeeker.SecretMemo += "・勤務先名:" + column + "\n"
			//368:勤務先業種大分類_名６
			//369:勤務先業種大分類_コード６
			//370:勤務先業種中分類_名６
			//371:勤務先業種中分類_コード６
			//372:勤務先業種小分類_名６
			case 372:
				jobSeeker.SecretMemo += "・勤務先業種:" + column + "\n"
			//373:勤務先業種小分類_コード６
			//374:勤務先規模（資本金）６_名
			case 374:
				jobSeeker.SecretMemo += "・勤務先規模（資本金）:" + column + "\n"
			//375:勤務先規模（資本金）６_コード
			//376:勤務先規模（従業員数）６_名
			case 376:
				jobSeeker.SecretMemo += "・勤務先規模（従業員数）:" + column + "\n"
			//377:勤務先規模（従業員数）６_コード
			//378:マネジメント６_名
			case 378:
				jobSeeker.SecretMemo += "・マネジメント経験:" + column + "\n"
			//379:マネジメント６_コード
			//380:業務上のポジション６_名
			case 380:
				jobSeeker.SecretMemo += "・業務上のポジション:" + column + "\n"
			//381:業務上のポジション６_コード
			//382:職務内容６
			case 382:
				jobSeeker.SecretMemo += "・職務内容:" + column + "\n\n"
			//383:在籍期間FROM７_名
			case 383:
				jobSeeker.SecretMemo += "・在籍開始年:" + column + "\n"
			//384:在籍期間TO７_名
			case 384:
				jobSeeker.SecretMemo += "・在籍終了年:" + column + "\n"
			//385:雇用形態７_名
			case 385:
				jobSeeker.SecretMemo += "・雇用形態:" + column + "\n"
			//386:雇用形態７_コード
			//387:勤務先名７
			case 387:
				jobSeeker.SecretMemo += "・勤務先名:" + column + "\n"
			//388:勤務先業種大分類_名７
			//389:勤務先業種大分類_コード７
			//390:勤務先業種中分類_名７
			//391:勤務先業種中分類_コード７
			//392:勤務先業種小分類_名７
			case 392:
				jobSeeker.SecretMemo += "・勤務先業種:" + column + "\n"
			//393:勤務先業種小分類_コード７
			//394:勤務先規模（資本金）７_名
			case 394:
				jobSeeker.SecretMemo += "・勤務先規模（資本金）:" + column + "\n"
			//395:勤務先規模（資本金）７_コード
			//396:勤務先規模（従業員数）７_名
			case 396:
				jobSeeker.SecretMemo += "・勤務先規模（従業員数）:" + column + "\n"
			//397:勤務先規模（従業員数）７_コード
			//398:マネジメント７_名
			case 398:
				jobSeeker.SecretMemo += "・マネジメント経験:" + column + "\n"
			//399:マネジメント７_コード
			//400:業務上のポジション７_名
			case 400:
				jobSeeker.SecretMemo += "・業務上のポジション:" + column + "\n"
			//401:業務上のポジション７_コード
			//402:職務内容７
			case 402:
				jobSeeker.SecretMemo += "・職務内容:" + column + "\n\n"
			//403:在籍期間FROM８_名
			case 403:
				jobSeeker.SecretMemo += "・在籍開始年:" + column + "\n"
			//404:在籍期間TO８_名
			case 404:
				jobSeeker.SecretMemo += "・在籍終了年:" + column + "\n"
			//405:雇用形態８_名
			case 405:
				jobSeeker.SecretMemo += "・雇用形態:" + column + "\n"
			//406:雇用形態８_コード
			//407:勤務先名８
			case 407:
				jobSeeker.SecretMemo += "・勤務先名:" + column + "\n"
			//408:勤務先業種大分類_名８
			//409:勤務先業種大分類_コード８
			//410:勤務先業種中分類_名８
			//411:勤務先業種中分類_コード８
			//412:勤務先業種小分類_名８
			case 412:
				jobSeeker.SecretMemo += "・勤務先業種:" + column + "\n"
			//413:勤務先業種小分類_コード８
			//414:勤務先規模（資本金）８_名
			case 414:
				jobSeeker.SecretMemo += "・勤務先規模（資本金）:" + column + "\n"
			//415:勤務先規模（資本金）８_コード
			//416:勤務先規模（従業員数）８_名
			case 416:
				jobSeeker.SecretMemo += "・勤務先規模（従業員数）:" + column + "\n"
			//417:勤務先規模（従業員数）８_コード
			//418:マネジメント８_名
			case 418:
				jobSeeker.SecretMemo += "・マネジメント経験:" + column + "\n"
			//419:マネジメント８_コード
			//420:業務上のポジション８_名
			case 420:
				jobSeeker.SecretMemo += "・業務上のポジション:" + column + "\n"
			//421:業務上のポジション８_コード
			//422:職務内容８
			case 422:
				jobSeeker.SecretMemo += "・職務内容:" + column + "\n\n"
			//423:在籍期間FROM９_名
			case 423:
				jobSeeker.SecretMemo += "・在籍開始年:" + column + "\n"
			//424:在籍期間TO９_名
			case 424:
				jobSeeker.SecretMemo += "・在籍終了年:" + column + "\n"
			//425:雇用形態９_名
			case 425:
				jobSeeker.SecretMemo += "・雇用形態:" + column + "\n"
			//426:雇用形態９_コード
			//427:勤務先名９
			case 427:
				jobSeeker.SecretMemo += "・勤務先名:" + column + "\n"
			//428:勤務先業種大分類_名９
			//429:勤務先業種大分類_コード９
			//430:勤務先業種中分類_名９
			//431:勤務先業種中分類_コード９
			//432:勤務先業種小分類_名９
			case 432:
				jobSeeker.SecretMemo += "・勤務先業種:" + column + "\n"
			//433:勤務先業種小分類_コード９
			//434:勤務先規模（資本金）９_名
			case 434:
				jobSeeker.SecretMemo += "・勤務先規模（資本金）:" + column + "\n"
			//435:勤務先規模（資本金）９_コード
			//436:勤務先規模（従業員数）９_名
			case 436:
				jobSeeker.SecretMemo += "・勤務先規模（従業員数）:" + column + "\n"
			//437:勤務先規模（従業員数）９_コード
			//438:マネジメント９_名
			case 438:
				jobSeeker.SecretMemo += "・マネジメント経験:" + column + "\n"
			//439:マネジメント９_コード
			//440:業務上のポジション９_名
			case 440:
				jobSeeker.SecretMemo += "・業務上のポジション:" + column + "\n"
			//441:業務上のポジション９_コード
			//442:職務内容９
			case 442:
				jobSeeker.SecretMemo += "・職務内容:" + column + "\n\n"
			//443:在籍期間FROM１０_名
			case 443:
				jobSeeker.SecretMemo += "・在籍開始年:" + column + "\n"
			//444:在籍期間TO１０_名
			//445:雇用形態１０_名
			case 445:
				jobSeeker.SecretMemo += "・雇用形態:" + column + "\n"
			//446:雇用形態１０_コード
			//447:勤務先名１０
			case 447:
				jobSeeker.SecretMemo += "・勤務先名:" + column + "\n"
			//448:勤務先業種大分類_名１０
			//449:勤務先業種大分類_コード１０
			//450:勤務先業種中分類_名１０
			//451:勤務先業種中分類_コード１０
			//452:勤務先業種小分類_名１０
			case 452:
				jobSeeker.SecretMemo += "・勤務先業種:" + column + "\n"
			//453:勤務先業種小分類_コード１０
			//454:勤務先規模（資本金）１０_名
			case 454:
				jobSeeker.SecretMemo += "・勤務先規模（資本金）:" + column + "\n"
			//455:勤務先規模（資本金）１０_コード
			//456:勤務先規模（従業員数）１０_名
			case 456:
				jobSeeker.SecretMemo += "・勤務先規模（従業員数）:" + column + "\n"
			//457:勤務先規模（従業員数）１０_コード
			//458:マネジメント１０_名
			case 458:
				jobSeeker.SecretMemo += "・マネジメント経験:" + column + "\n"
			//459:マネジメント１０_コード
			//460:業務上のポジション１０_名
			case 460:
				jobSeeker.SecretMemo += "・業務上のポジション:" + column + "\n"
			//461:業務上のポジション１０_コード
			//462:職務内容１０
			case 462:
				jobSeeker.SecretMemo += "・職務内容:" + column + "\n\n"
			//463:自己ＰＲ_名
			case 463:
				jobSeeker.SecretMemo += "・自己PR:\n" + column + "\n\n"
			//464:志望動機_名
			case 464:
				jobSeeker.SecretMemo += "・志望動機:\n" + column + "\n\n"
			//465:職務経歴（英文）_名
			case 465:
				jobSeeker.SecretMemo += "・職務経歴（英文）:\n" + column + "\n\n"
			//466:自己ＰＲ_名(英文)_名
			case 466:
				jobSeeker.SecretMemo += "・自己PR_名(英文):\n" + column + "\n\n"
			//467:担当者名
			case 467:
				jobSeeker.SecretMemo += "・担当者名:" + column + "\n"
			//468:担当者ID（メールアドレス）
			case 468:
				jobSeeker.SecretMemo += "・担当者ID（メールアドレス）:" + column + "\n"
				//469:
				//470:
			}
		}

		jobSeekerList = append(jobSeekerList, &jobSeeker)
		if isMatchedUserCount >= len(input.UserIDList) {
			log.Println("対象ユーザー数に達しました。 isMatchedUserCount:", isMatchedUserCount, " >= len(input.UserIDList):", len(input.UserIDList))
			break
		}
	}
	os.Remove(ouputFileName)
	file.Close()

	if len(jobSeekerList) == 0 {
		log.Println("対象ユーザーが見つかりませんでした")
		return output, nil
	}

	var (
		records [][]string
		record  []string
	)

	// csvの一行目を作成
	records = append(
		records,
		[]string{
			"candidate-externalId",
			"candidate-firstName",
			"candidate-Lastname",
			"candidate-email",
			"candidate-title",
			"candidate-middleName",
			"candidate-FirstNameKana",
			"candidate-LastNameKana",
			"candidate-workEmail",
			"candidate-employmentType",
			"candidate-jobTypes",
			"candidate-gender",
			"candidate-dob",
			"candidate-address",
			"candidate-citizenship",
			"candidate-Country",
			"candidate-city",
			"candidate-linkedln",
			"candidate-currentSalary",
			"candidate-desiredSalary",
			"candidate-company1",
			"candidate-company2",
			"candidate-company3",
			"candidate-contractInterval",
			"candidate-contractRate",
			"candidate-currency",
			"candidate-degreeName",
			"candidate-education",
			"candidate-educationLevel",
			"candidate-employer1",
			"candidate-employer2",
			"candidate-employer3",
			"candidate-gpa",
			"candidate-grade",
			"candidate-graduationDate",
			"candidate-homePhone",
			"candidate-phone",
			"candidate-mobile",
			"candidate-jobTitle1",
			"candidate-jobTitle2",
			"candidate-jobTitle3",
			"candidate-keyword",
			"candidate-note",
			"candidate-numberOfEmployers",
			"candidate-photo",
			"candidate-resume",
			"candidate-schoolName",
			"candidate-skills",
			"candidate-startDate1",
			"candidate-startDate2",
			"candidate-startDate3",
			"candidate-endDate1",
			"candidate-endDate2",
			"candidate-endDate3",
			"candidate-State",
			"candidate-workHistory",
			"candidate-workPhone",
			"candidate-zipCode",
			"candidate-owners",
			"candidate-comments",
		},
	)

	for _, jobSeeker := range jobSeekerList {
		// 求人ごとにデータを格納
		record =
			[]string{
				// 		"candidate-externalId",
				jobSeeker.ExternalID,
				// "candidate-firstName",
				jobSeeker.FirstName,
				// "candidate-Lastname",
				jobSeeker.LastName,
				// "candidate-email",
				jobSeeker.Email,
				// "candidate-title",
				"",
				// "candidate-middleName",
				"",
				// "candidate-FirstNameKana",
				jobSeeker.FirstFurigana,
				// "candidate-LastNameKana",
				jobSeeker.LastFurigana,
				// "candidate-workEmail",
				"",
				// "candidate-employmentType",
				"FULL_TIME",
				// "candidate-jobTypes",
				"PERMANENT",
				// "candidate-gender",
				entity.ConvertGenderToVincere(jobSeeker.Gender),
				// "candidate-dob",
				"",
				// "candidate-address",
				jobSeeker.Address,
				// "candidate-citizenship",
				"",
				// "candidate-Country",
				"",
				// "candidate-city",
				"",
				// "candidate-linkedln",
				"",
				// "candidate-currentSalary",
				fmt.Sprint(jobSeeker.AnnualIncome),
				// "candidate-desiredSalary",
				fmt.Sprint(jobSeeker.DesiredAnnualIncome),
				// "candidate-company1",
				"",
				// "candidate-company2",
				"",
				// "candidate-company3",
				"",
				// "candidate-contractInterval",
				"",
				// "candidate-contractRate",
				"",
				// "candidate-currency",
				"",
				// "candidate-degreeName",
				"",
				// "candidate-education",
				"",
				// "candidate-educationLevel",
				"",
				// "candidate-employer1",
				"",
				// "candidate-employer2",
				"",
				// "candidate-employer3",
				"",
				// "candidate-gpa",
				"",
				// "candidate-grade",
				"",
				// "candidate-graduationDate",
				"",
				// "candidate-homePhone",
				"",
				// "candidate-phone",
				jobSeeker.PhoneNumber,
				// "candidate-mobile",
				jobSeeker.PhoneNumber,
				// "candidate-jobTitle1",
				"",
				// "candidate-jobTitle2",
				"",
				// "candidate-jobTitle3",
				"",
				// "candidate-keyword",
				"",
				// "candidate-note",
				"",
				// "candidate-numberOfEmployers",
				"",
				// "candidate-photo",
				"",
				// "candidate-resume",
				"",
				// "candidate-schoolName",
				"",
				// "candidate-skills",
				"",
				// "candidate-startDate1",
				"",
				// "candidate-startDate2",
				"",
				// "candidate-startDate3",
				"",
				// "candidate-endDate1",
				"",
				// "candidate-endDate2",
				"",
				// "candidate-endDate3",
				"",
				// "candidate-State",
				"",
				// "candidate-workHistory",
				"",
				// "candidate-workPhone",
				"",
				// "candidate-zipCode",
				"",
				// "candidate-owners",
				"",
				// "candidate-comments",
				"",

				// jobSeeker.StaffName,
				// getStrPhaseForJobSeeker(jobSeeker.Phase),
				// jobSeeker.ChannelName,
				// jobSeeker.CreatedAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02 15:04:05"),
				// jobSeeker.InterviewDate.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02 15:04:05"),
				// getStrUserStatus(jobSeeker.UserStatus),
				// jobSeeker.LastName,
				// jobSeeker.FirstName,
				// jobSeeker.LastFurigana,
				// jobSeeker.FirstFurigana,
				// getStrGenderForJobSeeker(jobSeeker.Gender),
				// jobSeeker.GenderRemarks,
				// getStrNationalityForJobSeeker(jobSeeker.Nationality),
				// jobSeeker.NationalityRemarks,
				// getStrAvailable(jobSeeker.MedicalHistory),
				// jobSeeker.MedicalHistoryRemarks,
				// jobSeeker.Birthday,
				// getStrAvailable(jobSeeker.Spouse),
				// getStrAvailable(jobSeeker.SupportObligation),
				// fmt.Sprint(jobSeeker.Dependents.Int64),
				// jobSeeker.Email,
				// jobSeeker.PhoneNumber,
				// jobSeeker.EmergencyPhoneNumber,
				// getStrPrefecture(jobSeeker.Prefecture),
				// fmt.Sprint(jobSeeker.PostCode, "\n", getStrPrefecture(jobSeeker.Prefecture), jobSeeker.Address),
				// jobSeeker.AddressFurigana,
				// fmt.Sprint(jobSeeker.AnnualIncome.Int64) + "万円",

				// // 学歴
				// getStrStudyCategoryForJobSeeker(jobSeeker.StudyCategory),
				// getStrStudentHistoryList(jobSeeker.StudentHistories),

				// // 職歴
				// getStrStateOfEmployment(jobSeeker.StateOfEmployment),
				// getStrJobChangeForJobSeeker(jobSeeker.JobChange),
				// getStrAvailable(jobSeeker.ShortResignation),
				// jobSeeker.ShortResignationRemarks,
				// getStrWorkHistoryList(jobSeeker.WorkHistories),

				// // PCスキル
				// getStrExcelSkill(jobSeeker.ExcelSkill),
				// getStrWordSkill(jobSeeker.WordSkill),
				// getStrPowerPointSkill(jobSeeker.PowerPointSkill),

				// // 活かせるスキル・資格
				// getStrLicenseListForJobSeeker(jobSeeker.Licenses),
				// getStrPCToolListForJobSeeker(jobSeeker.PCTools),
				// getStrLanguageDevelopmentExperienceListForJobSeeker(jobSeeker.DevelopmentSkills),
				// getStrOSDevelopmentExperienceListForJobSeeker(jobSeeker.DevelopmentSkills),
				// getStrLanguageListForJobSeeker(jobSeeker.LanguageSkills),

				// // 希望条件
				// getStrDesiredHolidayTypeList(jobSeeker.DesiredHolidayTypes),
				// getStrDesiredCompanyScaleList(jobSeeker.DesiredCompanyScales),
				// getStrDesiredIndustryList(jobSeeker.DesiredIndustries),
				// getStrDesiredOccupationList(jobSeeker.DesiredOccupations),
				// getStrDesiredWorkLocationList(jobSeeker.DesiredWorkLocations),
				// getStrTransfer(jobSeeker.Transfer),
				// jobSeeker.TransferRequirement,
				// fmt.Sprint(jobSeeker.DesiredAnnualIncome.Int64) + "万円",
				// fmt.Sprint(getStrJoinCompanyPeriod(jobSeeker.JoinCompanyPeriod)),

				// getStrAppearanceForJobSeeker(jobSeeker.Appearance),
				// // jobSeeker.AppearanceDetailOfTruth,
				// // jobSeeker.AppearanceDetail,
				// getStrCommunicationForJobSeeker(jobSeeker.Communication),
				// // jobSeeker.CommunicationDetailOfTruth,
				// // jobSeeker.CommunicationDetail,
				// getStrThinkingForJobSeeker(jobSeeker.Thinking),
				// // jobSeeker.ThinkingDetailOfTruth,
				// // jobSeeker.ThinkingDetail,
				// getStrSelfPromotionList(jobSeeker.SelfPromotions),
				// jobSeeker.ResearchContent,
				// jobSeeker.AcceptancePoints,

				// // メモ
				// jobSeeker.SecretMemo,
				// // jobSeeker.PublicMemo,
			}

			// レコードを追加
		records = append(records, record)
	}

	//CSVファイルを作成
	filePath := ("./job-seeker-" + fmt.Sprint(utility.CreateUUID()) + ".csv")

	exportedFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	cw := csv.NewWriter(exportedFile)

	defer file.Close()
	defer cw.Flush()

	cw.WriteAll(records)

	if err := cw.Error(); err != nil {
		fmt.Println("error writing csv:", err)
		return output, err
	}

	// output.FilePath = entity.NewFilePath(filePath)

	// 求職者をDBに保存
	// for _, jobSeeker := range jobSeekerList {

	// 本番以外はメールアドレスに空白
	// if i.app.Env != "prd" {
	// 	jobSeeker.Email = ""
	// }

	// // 面談調整メールを送信
	// jobSeeker.Phase = null.NewInt(int64(entity.EntryInterview), true) // フェーズ： エントリー
	// jobSeeker.AgentID = input.AgentID
	// jobSeeker.AgentStaffID = null.NewInt(int64(input.ScoutService.AgentStaffID), true)
	// jobSeeker.RegisterPhase = null.NewInt(1, true) // 登録状況:下書き
	// jobSeeker.InterviewDate = time.Now().UTC()
	// jobSeeker.InflowChannelID = input.ScoutService.InflowChannelID

	// err := i.jobSeekerRepository.Create(jobSeeker)
	// if err != nil {
	// 	log.Println(err)
	// 	return output, err
	// }

	// // 開発環境の場合はユーザー名などを統一する
	// if os.Getenv("APP_ENV") != "prd" {
	// 	err = i.jobSeekerRepository.UpdateForDev(jobSeeker.ID)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return output, err
	// 	}
	// }

	// jobSeekerDocument := entity.NewJobSeekerDocument(
	// 	jobSeeker.ID,
	// 	"",
	// 	"",
	// 	"",
	// 	"",
	// 	"",
	// 	"",
	// 	"",
	// 	"",
	// 	"",
	// 	"",
	// )

	// err = i.jobSeekerDocumentRepository.Create(jobSeekerDocument)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// // 求職者作成時にメッセージグループ作成
	// // エージェントと求職者のチャットグループを作成
	// chatGroup := entity.NewChatGroupWithJobSeeker(
	// 	jobSeeker.AgentID,
	// 	jobSeeker.ID,
	// 	false, // 初めはLINE連携してないから false
	// )

	// err = i.chatGroupWithJobSeekerRepository.Create(chatGroup)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// // 面談実施待ちより前のフェーズの場合は「面談調整タスク」を作成
	// var (
	// 	phaseSub null.Int
	// 	date     string
	// )

	// phaseSub = null.NewInt(0, true) //日程調整依頼

	// // タスクの期限は当日で登録
	// now := time.Now()
	// date = now.Format("2006-01-02")

	// // 面談調整タスクの作成
	// interviewTaskGroup := entity.NewInterviewTaskGroup(
	// 	jobSeeker.AgentID,
	// 	jobSeeker.ID,
	// 	jobSeeker.InterviewDate,
	// 	utility.EarliestTime(), // 初期値,
	// )

	// err = i.interviewTaskGroupRepository.Create(interviewTaskGroup)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// interviewTask := entity.NewInterviewTask(
	// 	interviewTaskGroup.ID,
	// 	null.NewInt(0, false),
	// 	null.NewInt(0, false),
	// 	jobSeeker.Phase,
	// 	phaseSub,
	// 	"",
	// 	date,
	// 	null.NewInt(99, true),
	// 	getStrPhaseForJobSeeker(jobSeeker.Phase),
	// )

	// err = i.interviewTaskRepository.Create(interviewTask)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// // 希望勤務地
	// for _, desiredWorkLocation := range jobSeeker.DesiredWorkLocations {
	// 	desiredWorkLocation.JobSeekerID = jobSeeker.ID
	// 	err := i.jobSeekerDesiredWorkLocationRepository.Create(&desiredWorkLocation)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return output, err
	// 	}
	// }
	// }

	log.Println("処理成功. jobSeekerList:", jobSeekerList)

	output.JobSeekerList = jobSeekerList

	return output, nil
}

/*
AMBIから新規応募者を取得する

開発環境ログ 6名登録の場合8分かかる
2023/04/14 17:01:01 AMBIのエントリー取得を開始します
2023/04/14 17:09:19 AMBIの取得に成功しました [0xc0002eac00 0xc0002eb200 0xc0002eb800 0xc0005d4000 0xc0005d4600 0xc00067e000 0xc0005d4c00]
*/
type EntryOnAmbiInput struct {
	AgentID      uint
	ScoutService *entity.ScoutService
	Context      context.Context

	UserIDList []string
}

type EntryOnAmbiOutput struct {
	JobSeekerList []*entity.JobSeeker
}

func (i *ScoutServiceInteractorImpl) EntryOnAmbi(input EntryOnAmbiInput) (EntryOnAmbiOutput, error) {
	var (
		output        EntryOnAmbiOutput
		jobSeekerList []*entity.JobSeeker
		browser       *rod.Browser
		page          *rod.Page
		err           error
		errMessage    string
	)

	if input.ScoutService.LoginID == "" || input.ScoutService.Password == "" {
		errMessage = "ScoutServiceの認証情報が未入力です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	} else if input.Context == nil {
		errMessage = "Contextが空です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	} else if input.AgentID == 0 {
		errMessage = "AgentIDが未入力です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// 2日以内に登録されたの応募者を取得
	jobSeekerListInDb, err := i.jobSeekerRepository.GetByAgentIDWithinTwoDays(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// ホストOSのchrome, chromium or microsoft edgeのパスを取得
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

	// ブラウザをexitし、UserDataDirを削除
	defer time.Sleep(10 * time.Second)
	defer l.Cleanup()

	// ブラウザを起動後、AMBIへ遷移
	url := l.MustLaunch()
	browser = rod.New().
		ControlURL(url).
		Trace(true).                 // ログ出力
		SlowMotion(2 * time.Second). // 機械対策のページ用に遅延を入れる
		MustConnect()

	// 終了後ブラウザを閉じる
	defer browser.MustClose()

	// browserにタイムアウトを設定
	browser, cancel := browser.
		Context(input.Context).
		WithCancel()
	defer cancel()

	// AMBIのログインページへ遷移
	page = browser.
		MustPage("https://en-ambi.com/company_login/login/").
		MustWaitLoad()

	// ログインIDとパスワードを入力
	page.
		MustElement("input[name=accLoginID]").
		MustInput(input.ScoutService.LoginID)

	// パスワードの複号
	decryptedPassword, err := decryption(input.ScoutService.Password)
	if err != nil {
		log.Println(err)
		return output, err
	}

	page.
		MustElement("input[name=accLoginPW]").
		MustInput(decryptedPassword).
		MustType(rodInput.Enter).
		WaitLoad()

	time.Sleep(4 * time.Second)

	// URLがログインページのままの場合はログインに失敗していると判断
	if strings.Contains(page.MustInfo().URL, "login") {
		errMessage = "ログインに失敗しました。"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// 求職者のユーザーIDの回数分実行する
userIDLoop:
	for _, userID := range input.UserIDList {
		// エントリー一覧ページへ移動（ループに入れているのは2週目以降でID検索しても求職者が出てこないため）
		page.
			MustNavigate("https://en-ambi.com/company/entry/entry_list/?PK=6BF948").
			WaitLoad()

		// ユーザーIDで特定の求職者を検索
		page.
			MustElement("input[name=CName]").
			MustInput(userID).
			MustType(rodInput.Enter).
			WaitLoad()

		time.Sleep(4 * time.Second)

		// 該当の要素を全て取得
		users := page.
			MustElements("div.userSet")

		// 該当コンポーネントがない場合は次のユーザーIDを検索し直す
		if len(users) == 0 {
			continue
		}

		for userI, user := range users {
			var (
				jobSeeker     entity.JobSeeker = entity.JobSeeker{}
				experienceCnt int              = 0 // 何社目か
			)

			// ID検索のため一件しかヒットしない想定
			if userI >= 1 {
				break
			}

			// 退会済みの場合はスキップ
			quitCheckerEl := user.MustElement("div.data.btnSet")
			quitChecker := quitCheckerEl.MustText()
			if quitChecker == "退会済" || quitChecker == "利用停止中" {
				fmt.Println("quitChecker:", quitChecker)
				continue
			}

			userBasicDiv, err := user.
				Element("div.data.basic")
			if err != nil {
				errMessage = "ユーザーIDの取得に失敗しました。" + err.Error()
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			userIDOnPageDiv, err := userBasicDiv.
				Element("div.num")
			if err != nil {
				errMessage = "ユーザーIDの取得に失敗しました。" + err.Error()
				log.Println(errMessage)
				return output, errors.New(errMessage)
			}

			// No.912271 -> 912271
			userIDOnPage := strings.Replace(userIDOnPageDiv.MustText(), "No.", "", -1)
			log.Println("userIDOnPage:", userIDOnPage)

			// ユーザーIDと一致した場合のみ処理を行う
			if userID == userIDOnPage {
				// ユーザー詳細ページへ遷移
				user.
					MustElement("div.profileSet").
					MustClick()
				page.WaitLoad()
				time.Sleep(10 * time.Second)

				//職務経歴書更新日：21/09/07 最終ログイン日：23/03/14 山田 太郎（ヤマダ タロウ） / 26歳 エントリー済 スカウト済
				// -> 山田, 太郎, ヤマダ, タロウ, 26
				basicInfoWithLineBreak := page.MustElement("div.dataArea > div.name").MustText()
				// jobSeeker.SecretMemo = "・エントリー媒体\nAMBI\n\n・基本情報\n" + basicInfoWithLineBreak + "\n\n"
				// 改行を半角スペースに変換
				basicInfo := strings.Join(utility.RegexpForLineBreak.Split(basicInfoWithLineBreak, -1), " ")
				splitedBasicInfo := strings.Split(basicInfo, " ")

				log.Println("splitedBasicInfo:", splitedBasicInfo)

				if len(splitedBasicInfo) < 3 {
					log.Println("splitedBasicInfoの長さが3未満です", splitedBasicInfo)
					continue
				}

				// 名前を取得
				fullName := strings.Split(splitedBasicInfo[2], "（")[0]
				log.Println("fullName:", fullName)

				// 対応済みの場合は情報だけ取得
				if fullName == "＊＊＊＊＊" {
					// 名前がブラインドで対応済みの場合はNG対応のためスキップ
					if quitChecker == "対応済" {
						continue
					}

					// 面談設定テンプレートを送信 デフォルトで選択されている
					interviewLabel, err := page.Element("label[for=interview-01]")
					if err != nil {
						errMessage = "面談設定テンプレートを選択できませんでした"
						log.Println(errMessage)
						return output, errors.New(errMessage)
					}

					interviewLabel.MustClick()
					time.Sleep(5 * time.Second)

					templateSelectEl, err := page.
						Element("select#js_messageTemplate_select")
					if err != nil {
						errMessage = "面談設定テンプレートを選択できませんでした"
						log.Println(errMessage)
						return output, errors.New(errMessage)
					}

					err = templateSelectEl.
						Select(
							[]string{
								input.ScoutService.TemplateTitleForEmployed,
							},
							true,
							"text",
						)
					if err != nil {
						errMessage = "面談設定テンプレートを選択できませんでした"
						log.Println(errMessage)
						return output, errors.New(errMessage)
					}

					time.Sleep(10 * time.Second)

					// SMSも送信　input#smsCheck disabledの可能性もあるのでMustClick()ではなくClick()を使用
					if !page.MustHas("input.js_unavailable#smsCheck") {
						smsCheckBtn := page.MustElement("label[for=smsCheck]")
						smsCheckBtn.Click(proto.InputMouseButtonLeft, 1)
						time.Sleep(3 * time.Second)
					}

					// prd以外の場合は送信しないためここで終了
					if i.app.BatchType != "entry" {
						log.Println("テスト環境のため、メッセージ送信はスキップします")
						jobSeeker := entity.JobSeeker{
							SecretMemo: "・エントリー媒体：AMBI\n\n開発テスト用取り込み\n\n対象ID: " + fmt.Sprint(userID),
						}
						jobSeekerList = append(jobSeekerList, &jobSeeker)
						continue
					}

					// 送信ボタン
					// 送信ボタンが押せない場合はスキップ
					if page.MustHas("a.md_btn.md_btn--green.welcomeDone.md_btn.md_btn--disabled.js_disabled") {
						time.Sleep(10 * time.Second)
						if page.MustHas("a.md_btn.md_btn--green.welcomeDone.md_btn.md_btn--disabled.js_disabled") {
							continue userIDLoop
						}
					}
					page.MustElement("a.md_btn.md_btn--green.welcomeDone.js_tipOpen.js_sendMessage").MustClick()
					time.Sleep(5 * time.Second)

					// 確認ボタン　md_btn md_btn--min js_sendYes
					page.MustElement("a.md_btn.md_btn--min.js_sendYes").MustClick()
					page.WaitLoad()
					time.Sleep(15 * time.Second)

					// 送信後に情報取得
					//職務経歴書更新日：21/09/07 最終ログイン日：23/03/14 山田 太郎（ヤマダ タロウ） / 26歳 エントリー済 スカウト済
					// -> 山田, 太郎, ヤマダ, タロウ, 26
					basicInfoWithLineBreak = page.MustElement("div.dataArea > div.name").MustText()
					jobSeeker.SecretMemo = "・エントリー媒体\nAMBI\n\n・基本情報\n" + basicInfoWithLineBreak + "未対応\n\n"
					// 改行を半角スペースに変換
					basicInfo = strings.Join(utility.RegexpForLineBreak.Split(basicInfoWithLineBreak, -1), " ")
					splitedBasicInfo = strings.Split(basicInfo, " ")

					log.Println("splitedBasicInfo:", splitedBasicInfo)

					if len(splitedBasicInfo) < 5 {
						log.Println("splitedBasicInfoの長さが5未満です", splitedBasicInfo)
						continue
					}

					// 苗字を取得
					jobSeeker.LastName = splitedBasicInfo[2]
					log.Println("jobSeeker.LastName:", jobSeeker.LastName)

					// 名前を取得
					firstNameAndLastFurigana := strings.Split(splitedBasicInfo[3], "（")
					log.Println("firstNameAndLastFurigana:", firstNameAndLastFurigana)

					jobSeeker.FirstName = firstNameAndLastFurigana[0]
					log.Println("jobSeeker.FirstName:", jobSeeker.FirstName)

					// フリガナを取得
					if len(firstNameAndLastFurigana) > 1 {
						lastFurigana := strings.Replace(firstNameAndLastFurigana[1], "（", "", -1)
						log.Println("lastFurigana:", lastFurigana)
						jobSeeker.LastFurigana = lastFurigana
					}

					firstFurigana := strings.Replace(splitedBasicInfo[4], "）", "", -1)
					log.Println("firstFurigana:", firstFurigana)
					jobSeeker.FirstFurigana = firstFurigana

					// 名前がブラインドではない場合
				} else {
					// 名前を分割 &nbsp(\u00a0) で分割
					nameList := strings.Split(fullName, "\u00a0")

					for index, name := range nameList {
						if index == 0 {
							jobSeeker.LastName = name
						} else if index == 1 {
							jobSeeker.FirstName = name
						}
					}

					// フリガナを取得
					splitedFullFurigana := strings.Split(splitedBasicInfo[2], "（")
					if len(splitedFullFurigana) < 2 {
						log.Println("splitedFullFuriganaの長さが2未満です", splitedFullFurigana)
						jobSeeker.LastFurigana = splitedBasicInfo[2]
					} else {
						fullFurigana := strings.Split(splitedFullFurigana[1], "）")[0]
						log.Println("fullFurigana:", fullFurigana)

						// フリガナを分割
						furiganaList := strings.Split(fullFurigana, "\u00a0")
						log.Println("furiganaList:", furiganaList)

						for index, furigana := range furiganaList {
							if index == 0 {
								jobSeeker.LastFurigana = furigana
							} else if index == 1 {
								jobSeeker.FirstFurigana = furigana
							}
						}
					}
				}

				// 「名前」が既存の求職者と重複している場合はスキップ
				for _, jobSeekerInDb := range jobSeekerListInDb {
					if jobSeekerInDb.LastName == jobSeeker.LastName &&
						jobSeekerInDb.FirstName == jobSeeker.FirstName {
						// モーダルを閉じる
						fmt.Println("2日以内に登録した求職者と重複しています。", jobSeeker.LastName+jobSeeker.FirstName, jobSeeker.LastFurigana+jobSeeker.FirstFurigana)
						page.MustElement("a.closeBtn").MustClick()
						continue userIDLoop
					}
				}

				log.Println(
					"jobSeeker.LastName:", jobSeeker.LastName, "jobSeeker.FirstName:", jobSeeker.FirstName,
					"jobSeeker.LastFurigana:", jobSeeker.LastFurigana, "jobSeeker.FirstFurigana:", jobSeeker.FirstFurigana,
				)

				// 東京都 渋谷区 14-1 / 09012341234 / tarou.yamada@gmail.com
				// -> 東京都, 渋谷区 14-1, 09012341234, tarou.yamada@gmail.com
				addressAndContactWithSlash := page.MustElement("div.dataArea > div.data").MustText()

				splitedAddressAndContact := strings.Split(addressAndContactWithSlash, "/")
				for index, addressAndContact := range splitedAddressAndContact {
					// 都道府県と住所
					if index == 0 {
						prefectureAndAddress := strings.TrimSpace(addressAndContact)
						jobSeeker.Address = prefectureAndAddress

						splitedPrefectureAndAddress := strings.Split(prefectureAndAddress, "\u00a0")
						log.Println("splitedPrefectureAndAddress:", splitedPrefectureAndAddress)
						jobSeeker.Prefecture = entity.GetIntPrefecture(splitedPrefectureAndAddress[0])
					} else if index == 1 {
						// 電話番号
						phone := strings.TrimSpace(addressAndContact)
						jobSeeker.SecretMemo += "・電話番号\n" + phone + "\n\n"
						jobSeeker.PhoneNumber = phone
					} else if index == 2 {
						// メールアドレス
						email := strings.TrimSpace(addressAndContact)
						jobSeeker.SecretMemo += "・メールアドレス\n" + email + "\n\n"
						jobSeeker.Email = email
					}
				}

				// 求職者情報テーブルを取得
				userInfos := page.MustElements("table.md_tableForm > tbody > tr > td")
				userInfoLen := len(userInfos)
				log.Println("userInfoLen:", userInfos[0].MustText(), userInfoLen)

			userInfoLoop:
				for index, info := range userInfos {
					infoText := info.MustText()

					// 1990年（平成元年） 09月02日 -> 1990-09-03
					switch index {
					case 0:
						jobSeeker.SecretMemo += "・生年月日\n" + infoText + "\n\n"
						birthday := strings.Split(infoText, " ")
						if len(birthday) < 2 {
							log.Println("birthdayの長さが2未満です", birthday)
							continue userInfoLoop
						}

						yearArr := strings.Split(birthday[0], "年")
						year := strings.Split(yearArr[0], "（")
						monthAndDay := strings.Split(birthday[1], "月")
						month := monthAndDay[0]
						dayArr := strings.Split(monthAndDay[1], "日")
						day := dayArr[0]

						birthdayStr := year[0] + "-" + month + "-" + day
						log.Println("birthday:", birthdayStr)

						jobSeeker.Birthday = birthdayStr

					case 1:
						log.Println("最終学歴:", infoText)
						jobSeeker.SecretMemo += "・最終学歴\n" + infoText + "\n\n"

						// 大学卒 / テスト大学大学院テスト学部テスト学科 / 2019（令和元）年卒業
						// 大学卒
						// spaceを消す
						infoText = strings.Replace(infoText, "&nbsp;", "", -1)
						educationSplitedBySlash := strings.Split(infoText, "/")
						log.Println("educationSplitedBySlash:", educationSplitedBySlash)
						if len(educationSplitedBySlash) < 4 {
							log.Println("educationSplitedBySlashの長さが4未満です", educationSplitedBySlash)
							continue userInfoLoop
						}
						// 文系, 理系
						studyCategory := educationSplitedBySlash[2]
						if strings.Contains(studyCategory, "理系") {
							jobSeeker.StudyCategory = null.NewInt(0, true)
						} else if strings.Contains(studyCategory, "文系") {
							jobSeeker.StudyCategory = null.NewInt(1, true)
						}
						// 希望転職時期　1年以内
					case 2:
						jobSeeker.SecretMemo += "・希望転職時期\n" + infoText + "\n\n"
						// ambiマスタ:すぐに	3ヶ月以内	6ヶ月以内	1年以内	いいところがあれば	まだ未定	転職を考えていない
						if infoText == "すぐに" {
							jobSeeker.JoinCompanyPeriod = null.NewInt(0, true)
						} else if infoText == "3ヶ月以内" {
							jobSeeker.JoinCompanyPeriod = null.NewInt(3, true)
						} else if infoText == "6ヶ月以内" {
							jobSeeker.JoinCompanyPeriod = null.NewInt(6, true)
						} else if infoText == "1年以内" {
							jobSeeker.JoinCompanyPeriod = null.NewInt(7, true)
						}
						// 直近の年収 : 400万円以上
					case 3:
						jobSeeker.SecretMemo += "・直近の年収\n" + infoText + "\n\n"
						annualncomeStrTrimed := strings.Replace(infoText, "万円以上", "", -1)
						annualIncomeInt, err := strconv.Atoi(annualncomeStrTrimed)
						if err != nil {
							log.Println("年収の変換に失敗しました。")
							continue userInfoLoop
						}
						jobSeeker.AnnualIncome = null.NewInt(int64(annualIncomeInt), true)
						continue

						// 就業状況	就業している
					case 4:
						jobSeeker.SecretMemo += "・就業状況\n" + infoText + "\n\n"
						if infoText == "就業している" {
							jobSeeker.StateOfEmployment = null.NewInt(0, true)
						} else if infoText == "就業していない" {
							jobSeeker.StateOfEmployment = null.NewInt(1, true)
						}
						continue

						// 転職回数	1回
					case 5:
						jobSeeker.SecretMemo += "・転職回数\n" + infoText + "\n\n"
						if infoText == "転職経験なし" {
							jobSeeker.JobChange = null.NewInt(0, true)
						} else {
							jobChangeSplit := strings.Split(infoText, "回")
							jobChange, err := strconv.Atoi(jobChangeSplit[0])
							if err != nil {
								log.Println("転職回数の変換に失敗しました。")
								continue userInfoLoop
							}
							jobSeeker.JobChange = null.NewInt(int64(jobChange), true)
							if jobChange > 5 {
								jobSeeker.JobChange = null.NewInt(5, true)
							}
						}
						continue

						// 配偶者	なし
					case 6:
						jobSeeker.SecretMemo += "・配偶者\n" + infoText + "\n\n"
						if infoText == "なし" {
							jobSeeker.Spouse = null.NewInt(1, true)
						} else if infoText == "あり" {
							jobSeeker.Spouse = null.NewInt(0, true)
						}

						continue

						// 英語
						// 例:
						// 英語スキル
						// TOEIC：--- TOEFL：---
						// 会話：初級 読解：初級 作文：中級

						// その他の語学スキル（言語 ---語）
						// 会話：--- 読解：--- 作文：---
					case 7:
						jobSeeker.SecretMemo += "・語学スキル\n" + infoText + "\n\n"
						continue

						// 保有資格
						// 例:中学校・高校教諭専修免許状
					case 8:
						jobSeeker.SecretMemo += "・保有資格\n" + infoText + "\n\n"
						continue

						// 経験職種と年数
						// 例:
						// 店長・販売・店舗管理： 1年以上
						// 講師・教師・インストラクター： 1年以上
						// マーケティング・販促企画： 経験あり
					case 9:
						jobSeeker.SecretMemo += "・経験職種と年数\n" + infoText + "\n\n"
						continue

						// スキル	家電量販店での販売 / 教員 / 自営業
					case 10:
						jobSeeker.SecretMemo += "・スキル\n" + infoText + "\n\n"
						continue

						// 経験業界	教育・学校
					case 11:
						jobSeeker.SecretMemo += "・経験業界\n" + infoText + "\n\n"
						continue

						// マネジメント経験	あり（5人以下）
					case 12:
						jobSeeker.SecretMemo += "・マネジメント経験\n" + infoText + "\n\n"
						continue

						// キャリア要約
					case 13:
						jobSeeker.SecretMemo += "・キャリア要約\n" + infoText + "\n\n"
						continue

					// 希望条件
					// 希望職種 *職種経歴が可変のため、列の最後から取得する
					case userInfoLen - 4:
						desiredOccupationWithLineBreak := userInfos[userInfoLen-4].MustText()
						jobSeeker.SecretMemo += "・希望職種\n" + desiredOccupationWithLineBreak + "\n\n"

						desiredOccupationList := utility.RegexpForLineBreak.Split(desiredOccupationWithLineBreak, -1)

						// ambiの職種マスタをautoscoutに変換していく
						var duplicateCheck = make(map[int64]bool)

						for _, occupation := range desiredOccupationList {
							for _, ambiBigValue := range entity.AmbiOccupation {
								for ambiKey, ambiValue := range ambiBigValue {
									if occupation == ambiValue {
										autoscoutOccupation1, autoscoutOccupation2 := entity.ConvertAmbiOccupationInt(null.NewInt(int64(ambiKey), true))
										log.Println("ambi職種マスタ", ambiKey, ambiValue, "autoscoutマスタ", autoscoutOccupation1, autoscoutOccupation2)

										// 重複チェック
										if _, ok := duplicateCheck[autoscoutOccupation1]; !ok {
											duplicateCheck[autoscoutOccupation1] = true

											// 1つ目の職種
											if autoscoutOccupation1 != 0 {
												jobSeekerDesiredOccupation := entity.JobSeekerDesiredOccupation{
													JobSeekerID:       jobSeeker.ID,
													DesiredOccupation: null.NewInt(autoscoutOccupation1, true),
													DesiredRank:       null.NewInt(1, true),
												}
												jobSeeker.DesiredOccupations = append(jobSeeker.DesiredOccupations, jobSeekerDesiredOccupation)
											}
										}
										if _, ok := duplicateCheck[autoscoutOccupation2]; !ok {
											duplicateCheck[autoscoutOccupation2] = true
											// 2つ目の職種がある場合
											if autoscoutOccupation2 != 0 {
												jobSeekerDesiredOccupation := entity.JobSeekerDesiredOccupation{
													JobSeekerID:       jobSeeker.ID,
													DesiredOccupation: null.NewInt(autoscoutOccupation2, true),
													DesiredRank:       null.NewInt(1, true),
												}
												jobSeeker.DesiredOccupations = append(jobSeeker.DesiredOccupations, jobSeekerDesiredOccupation)
											}
										}
									}
								}
							}
						}
						continue userInfoLoop

					case userInfoLen - 3:

						// 希望都道府県  *職種経歴が可変のため、列の最後から取得する
						// 大阪府 / 東京都
						desiredPrefectureWithSlash := userInfos[userInfoLen-3].MustText()
						jobSeeker.SecretMemo += "・希望都道府県\n" + desiredPrefectureWithSlash + "\n\n"

						// 希望都道府県を分割
						desiredPrefectureList := strings.Split(desiredPrefectureWithSlash, "/")

						for _, desiredPrefecture := range desiredPrefectureList {
							desiredPrefecture = strings.TrimSpace(desiredPrefecture)
							prefecture := entity.GetIntPrefecture(desiredPrefecture)
							if prefecture.Valid {
								jobSeekerDesiredPrefecture := entity.JobSeekerDesiredWorkLocation{
									JobSeekerID:         jobSeeker.ID,
									DesiredWorkLocation: prefecture,
									DesiredRank:         null.NewInt(1, true),
								}
								jobSeeker.DesiredWorkLocations = append(jobSeeker.DesiredWorkLocations, jobSeekerDesiredPrefecture)
							}
						}
						continue userInfoLoop

					// 希望業界 *職種経歴が可変のため、列の最後から取得する
					case userInfoLen - 2:
						desiredIndustryWithLineBreak := userInfos[userInfoLen-2].MustText()
						jobSeeker.SecretMemo += "・希望業界\n" + desiredIndustryWithLineBreak + "\n\n"

						// ambiの業界マスタをautoscoutに変換していく
						duplicateCheck := make(map[int64]bool)
						desiredIndustryList := utility.RegexpForLineBreak.Split(desiredIndustryWithLineBreak, -1)
						for _, desiredIndustry := range desiredIndustryList {
							desiredIndustry = strings.TrimSpace(desiredIndustry)
							for _, ambiBigValue := range entity.AmbiIndustry {
								for ambiKey, ambiValue := range ambiBigValue {
									if desiredIndustry == ambiValue {
										autoscoutIndustry := entity.ConvertAmbiIndustryInt(null.NewInt(int64(ambiKey), true))
										// 重複チェック
										if _, ok := duplicateCheck[autoscoutIndustry]; !ok {
											duplicateCheck[autoscoutIndustry] = true

											log.Println("ambi業界マスタ", ambiKey, ambiValue, "autoscoutマスタ", autoscoutIndustry)
											// 1つ目の業界
											if autoscoutIndustry != 0 {
												jobSeekerDesiredIndustry := entity.JobSeekerDesiredIndustry{
													JobSeekerID:     jobSeeker.ID,
													DesiredIndustry: null.NewInt(autoscoutIndustry, true),
													DesiredRank:     null.NewInt(1, true),
												}
												jobSeeker.DesiredIndustries = append(jobSeeker.DesiredIndustries, jobSeekerDesiredIndustry)
											}
										}
									}
								}
							}
						}
						continue userInfoLoop
						// 希望年収 *職種経歴が可変のため、列の最後から取得する
						// 400万円以上
					case userInfoLen - 1:
						desiredIncomeStr := userInfos[userInfoLen-1].MustText()
						jobSeeker.SecretMemo += "・希望年収\n" + desiredIncomeStr + "\n\n"
						desiredIncomeStrTrimed := strings.Replace(desiredIncomeStr, "万円以上\n", "", -1)
						log.Println("希望年収", desiredIncomeStrTrimed)
						desiredIncomeInt, err := strconv.Atoi(desiredIncomeStrTrimed)
						if err != nil {
							log.Println(err, desiredIncomeStrTrimed)
							continue userInfoLoop
						}
						log.Println("希望年収", desiredIncomeInt)

						jobSeeker.DesiredAnnualIncome = null.NewInt(int64(desiredIncomeInt), true)
						log.Println("希望年収", jobSeeker.DesiredAnnualIncome)
						continue userInfoLoop

					default:
						// 職務経歴を取得
						experienceCnt++
						jobSeeker.SecretMemo += "・職務経歴" + fmt.Sprint(experienceCnt) + "\n" + infoText + "\n\n"
						continue userInfoLoop

					}
				}
				jobSeekerList = append(jobSeekerList, &jobSeeker)

				// モーダルを閉じる
				page.MustElement("a.closeBtn").MustClick()
				time.Sleep(2 * time.Second)
			} else {
				log.Println("ユーザーIDが一致しませんでしたのでスキップします")
				continue userIDLoop
			}
		}
	}

	// 求職者をDBに保存
	for _, jobSeeker := range jobSeekerList {

		// 本番以外はメールアドレスに空白
		// if i.app.Env != "prd" {
		// 	jobSeeker.Email = ""
		// }

		// 面談調整メールを送信
		jobSeeker.Phase = null.NewInt(int64(entity.EntryInterview), true) // フェーズ： エントリー
		jobSeeker.AgentID = input.AgentID
		jobSeeker.AgentStaffID = null.NewInt(int64(input.ScoutService.AgentStaffID), true)
		jobSeeker.RegisterPhase = null.NewInt(1, true) // 登録状況:下書き
		jobSeeker.InterviewDate = time.Now().UTC()
		jobSeeker.InflowChannelID = input.ScoutService.InflowChannelID

		err := i.jobSeekerRepository.Create(jobSeeker)
		if err != nil {
			log.Println(err)
			return output, err
		}

		// 開発環境の場合はユーザー名などを統一する
		// if os.Getenv("APP_ENV") != "prd" {
		// 	err = i.jobSeekerRepository.UpdateForDev(jobSeeker.ID)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return output, err
		// 	}
		// }

		jobSeekerDocument := entity.NewJobSeekerDocument(
			jobSeeker.ID,
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		)

		err = i.jobSeekerDocumentRepository.Create(jobSeekerDocument)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求職者作成時にメッセージグループ作成
		// エージェントと求職者のチャットグループを作成
		chatGroup := entity.NewChatGroupWithJobSeeker(
			jobSeeker.AgentID,
			jobSeeker.ID,
			false, // 初めはLINE連携してないから false
		)

		err = i.chatGroupWithJobSeekerRepository.Create(chatGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 面談実施待ちより前のフェーズの場合は「面談調整タスク」を作成
		var (
			phaseSub null.Int
			date     string
		)

		phaseSub = null.NewInt(0, true) //日程調整依頼

		// タスクの期限は当日で登録
		now := time.Now()
		date = now.Format("2006-01-02")

		// 面談調整タスクの作成
		interviewTaskGroup := entity.NewInterviewTaskGroup(
			jobSeeker.AgentID,
			jobSeeker.ID,
			jobSeeker.InterviewDate,
			utility.EarliestTime(), // 初期値,
		)

		err = i.interviewTaskGroupRepository.Create(interviewTaskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		interviewTask := entity.NewInterviewTask(
			interviewTaskGroup.ID,
			null.NewInt(0, false),
			null.NewInt(0, false),
			jobSeeker.Phase,
			phaseSub,
			"",
			date,
			null.NewInt(99, true),
			getStrPhaseForJobSeeker(jobSeeker.Phase),
		)

		err = i.interviewTaskRepository.Create(interviewTask)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 希望職種
		for _, desiredOccupation := range jobSeeker.DesiredOccupations {
			desiredOccupation.JobSeekerID = jobSeeker.ID
			err := i.jobSeekerDesiredOccupationRepository.Create(&desiredOccupation)
			if err != nil {
				log.Println(err)
				return output, err
			}
		}

		// 希望業界
		for _, desiredIndustry := range jobSeeker.DesiredIndustries {
			desiredIndustry.JobSeekerID = jobSeeker.ID
			err := i.jobSeekerDesiredIndustryRepository.Create(&desiredIndustry)
			if err != nil {
				log.Println(err)
				return output, err
			}
		}

		// 希望勤務地
		for _, desiredWorkLocation := range jobSeeker.DesiredWorkLocations {
			desiredWorkLocation.JobSeekerID = jobSeeker.ID
			err := i.jobSeekerDesiredWorkLocationRepository.Create(&desiredWorkLocation)
			if err != nil {
				log.Println(err)
				return output, err
			}
		}
	}

	log.Println("処理成功. jobSeekerList:", jobSeekerList)

	output.JobSeekerList = jobSeekerList

	return output, nil
}

/*
マイナビから新規応募者を取得する

開発環境でかかった時間 5分(5名登録):
2023/04/14 16:01:00 マイナビのエントリー取得を開始します
30分経過でタイムアウト
*/
type EntryOnMynaviAgentScoutInput struct {
	AgentID        uint
	Context        context.Context
	ScoutServiceID uint
}

type EntryOnMynaviAgentScoutOutput struct {
	JobSeekerList []*entity.JobSeeker
}

func (i *ScoutServiceInteractorImpl) EntryOnMynaviAgentScout(input EntryOnMynaviAgentScoutInput) (EntryOnMynaviAgentScoutOutput, error) {
	var (
		output        EntryOnMynaviAgentScoutOutput
		err           error
		errMessage    string
		jobSeekerList []*entity.JobSeeker = []*entity.JobSeeker{}
		browser       *rod.Browser
		page          *rod.Page
		scoutService  *entity.ScoutService
	)

	if input.ScoutServiceID == 0 {
		errMessage = "ScoutServiceIDが未入力です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	scoutService, err = i.scoutServiceRepository.FindByID(input.ScoutServiceID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	if scoutService.LoginID == "" || scoutService.Password == "" {
		errMessage = "ScoutServiceの認証情報が未入力です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	} else if input.Context == nil {
		errMessage = "Contextが空です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	} else if input.AgentID == 0 {
		errMessage = "AgentIDが未入力です"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	// 2日以内に登録されたの応募者を取得
	jobSeekerListInDb, err := i.jobSeekerRepository.GetByAgentIDWithinTwoDays(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}
	log.Println(jobSeekerListInDb)

	// ホストOSのchrome, chromium or microsoft edgeのパスを取得
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
	// remove launcher.FlagUserDataDir
	defer l.Cleanup()

	url := l.MustLaunch()

	// ブラウザを起動後、AMBIへ遷移
	browser = rod.New().
		ControlURL(url).
		Trace(true).                 // デバッグ用のためローカル以外で使わない
		SlowMotion(2 * time.Second). // 機械対策のページ用に遅延を入れる
		MustConnect()

	// ブラウザ閉じる
	defer browser.MustClose()

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

	companyIDInput.MustInput(scoutService.LoginID)

	// メールアドレス
	emailInput, err := page.Element("input#email")
	if err != nil {
		log.Println(err)
		return output, err
	}

	emailInput.MustInput("info@spaceai.jp")

	// パスワードの複号
	decryptedPassword, err := decryption(scoutService.Password)
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
	time.Sleep(10 * time.Second)

	// マイナビの場合は、ログイン失敗時、再度ログインする
	if page.MustHas("input#email") {
		log.Println("ログイン失敗1回目、再度ログインします")

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
		time.Sleep(10 * time.Second)
	}

	// 再度ログイン時も失敗した場合は、エラーを返す
	if page.MustHas("input#email") {
		errMessage = "ログイン失敗しました。IDとパスワードが間違っているか。すでに利用しているユーザーがいる可能性があります。"
		log.Println(errMessage)
		return output, errors.New(errMessage)
	}

	log.Println("ログイン成功")

	// エントリー管理ページへ遷移
	page.MustNavigate("https://scout.mynavi-agent.jp/progress/entried")
	page.WaitLoad()
	time.Sleep(10 * time.Second)

	// searchNavButton isDownload
	openModalBtn, err := page.Element("button.searchNavButton.isDownload")
	if err != nil {
		log.Println(err)
		return output, err
	}
	openModalBtn.MustClick()
	time.Sleep(10 * time.Second)

	modalContainer, err := page.Element("div.modalContainer")
	if err != nil {
		log.Println(err)
		return output, err
	}

	// csvダウンロード
	wait := browser.MustWaitDownload()
	downloadCSVBtn, err := modalContainer.Element("button.generalButton.blue")
	if err != nil {
		log.Println(err)
		return output, err
	}
	downloadCSVBtn.MustClick()
	ouputFileName := fmt.Sprintf("getentry_mynavi_agent_scout_%v_%v.csv", scoutService.AgentRobotID, time.Now().Format("20060102150405"))
	err = utils.OutputFile(ouputFileName, wait())
	if err != nil {
		log.Println(err)
		return output, err
	}
	defer os.Remove(ouputFileName)
	time.Sleep(20 * time.Second)

	// csvファイルを読み込み
	file, err := os.Open(ouputFileName)
	if err != nil {
		log.Println(err)
		return output, err
	}
	defer file.Close()
	log.Println("csvファイルを開く file:", file.Name(), *file)

	// ShiftJISのCSVファイルを文字化けしない様に読み込む
	reader := csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))
	log.Println("csvファイルを読み込み reader:", reader)

	reader.TrimLeadingSpace = true // true の場合は、先頭の空白文字を無視する
	reader.ReuseRecord = true      // true の場合は、Read で戻ってくるスライスを次回再利用する。パフォーマンスが上がる

	// 項目説明のレコードを読み込み
	reader.Read()

	// recordLoop:
	for {
		record, err := reader.Read()
		// 空行があった場合は終了
		log.Println("record: ", record)
		if err == io.EOF || record[0] == "" {
			log.Println("行が空のため、ループを終了します")
			break
		}

		jobSeeker := entity.JobSeeker{}
		jobSeeker.SecretMemo = "・エントリー媒体\nマイナビスカウティング\n\n"

		for columnI, column := range record {
			log.Printf("column[%v]: %v", columnI, column)
			if column == "" {
				continue
			}

			switch columnI {

			// 0	求職者ID
			// 1	求職者氏名（姓）
			case 1:
				jobSeeker.LastName = column
			// 2	求職者氏名（名）
			case 2:
				jobSeeker.FirstName = column
			// 3	フリガナ（姓）
			case 3:
				jobSeeker.LastFurigana = column
			// 4	フリガナ（名）
			case 4:
				jobSeeker.FirstFurigana = column
				// 5	性別
				// 6	生年月日
				// 7	年齢
				// 8	都道府県
				// 9	市区町村
				// 10	電話番号1
				// 11	メールアドレス1
				// 12	現在の状況（在職・離職・その他・不明）
				// 13	経験社数
				// 14	現年収（万円）
				// 15	希望業種1
				// 16	希望業種2
				// 17	希望業種3
				// 18	希望業種4
				// 19	希望業種5
				// 20	希望業種6
				// 21	希望職種1
				// 22	希望職種2
				// 23	希望職種3
				// 24	希望職種4
				// 25	希望職種5
				// 26	希望職種6
				// 27	希望年収（万円）
				// 28	最低希望年収（万円）
				// 29	希望する雇用形態
				// 30	希望する勤務地1
				// 31	希望する勤務地2
				// 32	希望する勤務地3
				// 33	経験業種1
				// 34	経験業種1 年数
				// 35	経験業種2
				// 36	経験業種2 年数
				// 37	経験業種3
				// 38	経験業種3 年数
				// 39	経験業種4
				// 40	経験業種4 年数
				// 41	経験業種5
				// 42	経験業種5 年数
				// 43	経験業種6
				// 44	経験業種6 年数
				// 45	経験職種1
				// 46	経験職種1 年数
				// 47	経験職種2
				// 48	経験職種2 年数
				// 49	経験職種3
				// 50	経験職種3 年数
				// 51	経験職種4
				// 52	経験職種4 年数
				// 53	経験職種5
				// 54	経験職種5 年数
				// 55	経験職種6
				// 56	経験職種6 年数
				// 57	企業名1
				// 58	部署名1
				// 59	業種1
				// 60	年収（万円）1
				// 61	雇用形態1
				// 62	役職1
				// 63	期間From（YY/MM）1
				// 64	期間From 区分（入社・出向など）1
				// 65	期間To（YY/MM）1
				// 66	期間To 区分（退職・出向終了など）1
				// 67	職務内容1
				// 68	転職理由1
				// 69	備考1
				// 70	企業名2
				// 71	部署名2
				// 72	業種2
				// 73	年収（万円）2
				// 74	雇用形態2
				// 75	役職2
				// 76	期間From（YY/MM）2
				// 77	期間From 区分（入社・出向など）2
				// 78	期間To（YY/MM）2
				// 79	期間To 区分（退職・出向終了など）2
				// 80	職務内容2
				// 81	転職理由2
				// 82	備考2
				// 83	企業名3
				// 84	部署名3
				// 85	業種3
				// 86	年収（万円）3
				// 87	雇用形態3
				// 88	役職3
				// 89	期間From（YY/MM）3
				// 90	期間From 区分（入社・出向など）3
				// 91	期間To（YY/MM）3
				// 92	期間To 区分（退職・出向終了など）3
				// 93	職務内容3
				// 94	転職理由3
				// 95	備考3
				// 96	企業名4
				// 97	部署名4
				// 98	業種4
				// 99	年収（万円）4
				// 100	雇用形態4
				// 101	役職4
				// 102	期間From（YY/MM）4
				// 103	期間From 区分（入社・出向など）4
				// 104	期間To（YY/MM）4
				// 105	期間To 区分（退職・出向終了など）4
				// 106	職務内容4
				// 107	転職理由4
				// 108	備考4
				// 109	企業名5
				// 110	部署名5
				// 111	業種5
				// 112	年収（万円）5
				// 113	雇用形態5
				// 114	役職5
				// 115	期間From（YY/MM）5
				// 116	期間From 区分（入社・出向など）5
				// 117	期間To（YY/MM）5
				// 118	期間To 区分（退職・出向終了など）5
				// 119	職務内容5
				// 120	転職理由5
				// 121	備考5
				// 122	企業名6
				// 123	部署名6
				// 124	業種6
				// 125	年収（万円）6
				// 126	雇用形態6
				// 127	役職6
				// 128	期間From（YY/MM）6
				// 129	期間From 区分（入社・出向など）6
				// 130	期間To（YY/MM）6
				// 131	期間To 区分（退職・出向終了など）6
				// 132	職務内容6
				// 133	転職理由6
				// 134	備考6
				// 135	企業名7
				// 136	部署名7
				// 137	業種7
				// 138	年収（万円）7
				// 139	雇用形態7
				// 140	役職7
				// 141	期間From（YY/MM）7
				// 142	期間From 区分（入社・出向など）7
				// 143	期間To（YY/MM）7
				// 144	期間To 区分（退職・出向終了など）7
				// 145	職務内容7
				// 146	転職理由7
				// 147	備考7
				// 148	企業名8
				// 149	部署名8
				// 150	業種8
				// 151	年収（万円）8
				// 152	雇用形態8
				// 153	役職8
				// 154	期間From（YY/MM）8
				// 155	期間From 区分（入社・出向など）8
				// 156	期間To（YY/MM）8
				// 157	期間To 区分（退職・出向終了など）8
				// 158	職務内容8
				// 159	転職理由8
				// 160	備考8
				// 161	企業名9
				// 162	部署名9
				// 163	業種9
				// 164	年収（万円）9
				// 165	雇用形態9
				// 166	役職9
				// 167	期間From（YY/MM）9
				// 168	期間From 区分（入社・出向など）9
				// 169	期間To（YY/MM）9
				// 170	期間To 区分（退職・出向終了など）9
				// 171	職務内容9
				// 172	転職理由9
				// 173	備考9
				// 174	企業名10
				// 175	部署名10
				// 176	業種10
				// 177	年収（万円）10
				// 178	雇用形態10
				// 179	役職10
				// 180	期間From（YY/MM）10
				// 181	期間From 区分（入社・出向など）10
				// 182	期間To（YY/MM）10
				// 183	期間To 区分（退職・出向終了など）10
				// 184	職務内容10
				// 185	転職理由10
				// 186	備考10
				// 187	最終学歴
				// 188	理系フラグ
				// 189	学校名1
				// 190	学部学科1
				// 191	期間To（YY/MM）1
				// 192	期間To 区分（卒業・中退・卒業見込）1
				// 193	学校名2
				// 194	学部学科2
				// 195	期間To（YY/MM）2
				// 196	期間To 区分（卒業・中退・卒業見込）2
				// 197	学校名3
				// 198	学部学科3
				// 199	期間To（YY/MM）3
				// 200	期間To 区分（卒業・中退・卒業見込）3
				// 201	学校名4
				// 202	学部学科4
				// 203	期間To（YY/MM）4
				// 204	期間To 区分（卒業・中退・卒業見込）4
				// 205	学校名5
				// 206	学部学科5
				// 207	期間To（YY/MM）5
				// 208	期間To 区分（卒業・中退・卒業見込）5
				// 209	資格名(取得年月)1
				// 210	資格名(取得年月)2
				// 211	資格名(取得年月)3
				// 212	資格名(取得年月)4
				// 213	資格名(取得年月)5
				// 214	資格名(取得年月)6
				// 215	資格名(取得年月)7
				// 216	資格名(取得年月)8
				// 217	資格名(取得年月)9
				// 218	資格名(取得年月)10
				// 219	その他資格
				// 220	TOEIC 点数
				// 221	TOEIC 取得年月
				// 222	TOEFL 点数
				// 223	TOEFL 取得年月
				// 224	英語レベル
				// 225	その他補足
				// 226	言語1
				// 227	言語1 レベル
				// 228	言語2
				// 229	言語2 レベル
				// 230	言語3
				// 231	言語3 レベル
				// 232	マネジメント経験有無
				// 233	マネジメント人数
				// 234	マネジメント年数
				// 235	転職理由
				// 236	転職で優先したいこと・叶えたいこと
				// 237	スカウトID
				// 238	スカウト応募日時
				// 239	エージェントアカウント名
				// 240	応募時メッセージ

			}
		}

		jobSeekerList = append(jobSeekerList, &jobSeeker)
	}

	// 求職者をDBに保存
	for _, jobSeeker := range jobSeekerList {
		jobSeeker.AgentID = input.AgentID
		jobSeeker.AgentStaffID = null.NewInt(int64(scoutService.AgentStaffID), true)
		jobSeeker.Phase = null.NewInt(0, true)         // フェーズ： エントリー
		jobSeeker.RegisterPhase = null.NewInt(1, true) // 登録状況:下書き

		// 本番以外はメールアドレスに空白
		// if i.app.Env != "prd" {
		// 	jobSeeker.Email = ""
		// }

		err := i.jobSeekerRepository.Create(jobSeeker)
		if err != nil {
			log.Println(err)
			return output, err
		}

		// 開発環境の場合はユーザー名などを統一する
		// if os.Getenv("APP_ENV") != "prd" {
		// 	err = i.jobSeekerRepository.UpdateForDev(jobSeeker.ID)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return output, err
		// 	}
		// }

		jobSeekerDocument := entity.NewJobSeekerDocument(
			jobSeeker.ID,
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		)

		err = i.jobSeekerDocumentRepository.Create(jobSeekerDocument)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求職者作成時にメッセージグループ作成
		// エージェントと求職者のチャットグループを作成
		chatGroup := entity.NewChatGroupWithJobSeeker(
			jobSeeker.AgentID,
			jobSeeker.ID,
			false, // 初めはLINE連携してないから false
		)

		err = i.chatGroupWithJobSeekerRepository.Create(chatGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 面談実施待ちより前のフェーズの場合は「面談調整タスク」を作成
		var (
			phaseSub null.Int
			date     string
		)

		phaseSub = null.NewInt(0, true) //日程調整依頼

		// タスクの期限は当日で登録
		now := time.Now()
		date = now.Format("2006-01-02")

		// 面談調整タスクの作成
		interviewTaskGroup := entity.NewInterviewTaskGroup(
			jobSeeker.AgentID,
			jobSeeker.ID,
			jobSeeker.InterviewDate,
			utility.EarliestTime(),
		)

		err = i.interviewTaskGroupRepository.Create(interviewTaskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		interviewTask := entity.NewInterviewTask(
			interviewTaskGroup.ID,
			null.NewInt(0, false),
			null.NewInt(0, false),
			jobSeeker.Phase,
			phaseSub,
			"",
			date,
			null.NewInt(99, true),
			getStrPhaseForJobSeeker(jobSeeker.Phase),
		)

		err = i.interviewTaskRepository.Create(interviewTask)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 希望勤務地
		for _, desiredWorkLocation := range jobSeeker.DesiredWorkLocations {
			desiredWorkLocation.JobSeekerID = jobSeeker.ID
			err := i.jobSeekerDesiredWorkLocationRepository.Create(&desiredWorkLocation)
			if err != nil {
				log.Println(err)
				return output, err
			}
		}
	}

	log.Println("処理成功. jobSeekerList:", jobSeekerList)

	output.JobSeekerList = jobSeekerList

	return output, nil
}
