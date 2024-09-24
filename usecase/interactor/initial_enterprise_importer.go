package interactor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"

	// ブラウザ操作
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type InitialEnterpriseImporterInteractor interface {
	// 汎用系 API
	CreateInitialEnterpriseImporter(input CreateInitialEnterpriseImporterInput) (CreateInitialEnterpriseImporterOutput, error)
	DeleteInitialEnterpriseImporter(input DeleteInitialEnterpriseImporterInput) (DeleteInitialEnterpriseImporterOutput, error)
	GetInitialEnterpriseImporterListByAgentID(input GetInitialEnterpriseImporterListByAgentIDInput) (GetInitialEnterpriseImporterListByAgentIDOutput, error)
	GetInitialEnterpriseImporterListByWeek() (GetInitialEnterpriseImporterListByWeekOutput, error) // 今後1週間の企業求人一括インポート情報一覧を取得 *予約日程調整時の重複チェック用

	// 他サービスから求人取得
	BatchInitialEnterpriseImporter(input BatchInitialEnterpriseImporterInput) (BatchInitialEnterpriseImporterOutput, error)
	ImportEnterpriseFromCircus(input ImportEnterpriseFromCircusInput) (ImportEnterpriseFromCircusOutput, error)
	ImportEnterpriseFromAgentBank(input ImportEnterpriseFromAgentBankInput) (ImportEnterpriseFromAgentBankOutput, error)
	BatchUpdateEnterpriseForAgentBank(input BatchUpdateEnterpriseForAgentBankInput) (BatchUpdateEnterpriseForAgentBankOutput, error)
	UpdateEnterpriseFromAgentBank(input UpdateEnterpriseFromAgentBankInput) (UpdateEnterpriseFromAgentBankOutput, error)
}

type InitialEnterpriseImporterInteractorImpl struct {
	firebase                                           usecase.Firebase
	sendgrid                                           config.Sendgrid
	oneSignal                                          config.OneSignal
	initialEnterpriseImporterRepository                usecase.InitialEnterpriseImporterRepository
	enterpriseProfileRepository                        usecase.EnterpriseProfileRepository
	enterpriseIndustryRepository                       usecase.EnterpriseIndustryRepository
	enterpriseReferenceMaterialRepository              usecase.EnterpriseReferenceMaterialRepository
	billingAddressRepository                           usecase.BillingAddressRepository
	billingAddressHRStaffRepository                    usecase.BillingAddressHRStaffRepository
	billingAddressRAStaffRepository                    usecase.BillingAddressRAStaffRepository
	jobInformationRepository                           usecase.JobInformationRepository
	jobInfoTargetRepository                            usecase.JobInformationTargetRepository
	jobInfoFeatureRepository                           usecase.JobInformationFeatureRepository
	jobInfoPrefectureRepository                        usecase.JobInformationPrefectureRepository
	jobInfoWorkCharmPointRepository                    usecase.JobInformationWorkCharmPointRepository
	jobInfoEmploymentStatusRepository                  usecase.JobInformationEmploymentStatusRepository
	jobInfoRequiredConditionRepository                 usecase.JobInformationRequiredConditionRepository
	jobInfoRequiredLicenseRepository                   usecase.JobInformationRequiredLicenseRepository
	jobInfoRequiredPCToolRepository                    usecase.JobInformationRequiredPCToolRepository
	jobInfoRequiredLanguageRepository                  usecase.JobInformationRequiredLanguageRepository
	jobInfoRequiredLanguageTypeRepository              usecase.JobInformationRequiredLanguageTypeRepository
	jobInfoRequiredExperienceDevelopmentRepository     usecase.JobInformationRequiredExperienceDevelopmentRepository
	jobInfoRequiredExperienceDevelopmentTypeRepository usecase.JobInformationRequiredExperienceDevelopmentTypeRepository
	jobInfoRequiredExperienceJobRepository             usecase.JobInformationRequiredExperienceJobRepository
	jobInfoRequiredExperienceIndustryRepository        usecase.JobInformationRequiredExperienceIndustryRepository
	jobInfoRequiredExperienceOccupationRepository      usecase.JobInformationRequiredExperienceOccupationRepository
	jobInfoRequiredSocialExperienceRepository          usecase.JobInformationRequiredSocialExperienceRepository
	jobInfoSelectionFlowPatternRepository              usecase.JobInformationSelectionFlowPatternRepository
	jobInfoSelectionInformationRepository              usecase.JobInformationSelectionInformationRepository
	jobInfoOccupationRepository                        usecase.JobInformationOccupationRepository
	agentRepository                                    usecase.AgentRepository
	agentStaffRepository                               usecase.AgentStaffRepository
}

// InitialEnterpriseImporterInteractorImpl is an implementation of InitialEnterpriseImporterInteractor
func NewInitialEnterpriseImporterInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	aeiR usecase.InitialEnterpriseImporterRepository,
	epR usecase.EnterpriseProfileRepository,
	eiR usecase.EnterpriseIndustryRepository,
	ermR usecase.EnterpriseReferenceMaterialRepository,
	baR usecase.BillingAddressRepository,
	bahsR usecase.BillingAddressHRStaffRepository,
	barsR usecase.BillingAddressRAStaffRepository,
	jR usecase.JobInformationRepository,
	jtR usecase.JobInformationTargetRepository,
	jfR usecase.JobInformationFeatureRepository,
	jpR usecase.JobInformationPrefectureRepository,
	jwcpR usecase.JobInformationWorkCharmPointRepository,
	jesR usecase.JobInformationEmploymentStatusRepository,
	jrcR usecase.JobInformationRequiredConditionRepository,
	jrlR usecase.JobInformationRequiredLicenseRepository,
	jrptR usecase.JobInformationRequiredPCToolRepository,
	jrlgR usecase.JobInformationRequiredLanguageRepository,
	jrlgtR usecase.JobInformationRequiredLanguageTypeRepository,
	jredR usecase.JobInformationRequiredExperienceDevelopmentRepository,
	jredtR usecase.JobInformationRequiredExperienceDevelopmentTypeRepository,
	jrejR usecase.JobInformationRequiredExperienceJobRepository,
	jreiR usecase.JobInformationRequiredExperienceIndustryRepository,
	jreoR usecase.JobInformationRequiredExperienceOccupationRepository,
	jrseR usecase.JobInformationRequiredSocialExperienceRepository,
	jsfpR usecase.JobInformationSelectionFlowPatternRepository,
	jsiR usecase.JobInformationSelectionInformationRepository,
	joR usecase.JobInformationOccupationRepository,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
) InitialEnterpriseImporterInteractor {
	return &InitialEnterpriseImporterInteractorImpl{
		firebase:                                           fb,
		sendgrid:                                           sg,
		oneSignal:                                          os,
		initialEnterpriseImporterRepository:                aeiR,
		enterpriseProfileRepository:                        epR,
		enterpriseIndustryRepository:                       eiR,
		enterpriseReferenceMaterialRepository:              ermR,
		billingAddressRepository:                           baR,
		billingAddressHRStaffRepository:                    bahsR,
		billingAddressRAStaffRepository:                    barsR,
		jobInformationRepository:                           jR,
		jobInfoTargetRepository:                            jtR,
		jobInfoFeatureRepository:                           jfR,
		jobInfoPrefectureRepository:                        jpR,
		jobInfoWorkCharmPointRepository:                    jwcpR,
		jobInfoEmploymentStatusRepository:                  jesR,
		jobInfoRequiredConditionRepository:                 jrcR,
		jobInfoRequiredLicenseRepository:                   jrlR,
		jobInfoRequiredPCToolRepository:                    jrptR,
		jobInfoRequiredLanguageRepository:                  jrlgR,
		jobInfoRequiredLanguageTypeRepository:              jrlgtR,
		jobInfoRequiredExperienceDevelopmentRepository:     jredR,
		jobInfoRequiredExperienceDevelopmentTypeRepository: jredtR,
		jobInfoRequiredExperienceJobRepository:             jrejR,
		jobInfoRequiredExperienceIndustryRepository:        jreiR,
		jobInfoRequiredExperienceOccupationRepository:      jreoR,
		jobInfoRequiredSocialExperienceRepository:          jrseR,
		jobInfoSelectionFlowPatternRepository:              jsfpR,
		jobInfoSelectionInformationRepository:              jsiR,
		jobInfoOccupationRepository:                        joR,
		agentRepository:                                    aR,
		agentStaffRepository:                               asR,
	}
}

/****************************************************************************************/
/// 汎用系
//
// 作成
type CreateInitialEnterpriseImporterInput struct {
	CreateParam entity.InitialEnterpriseImporter
}

type CreateInitialEnterpriseImporterOutput struct {
	InitialEnterpriseImporter *entity.InitialEnterpriseImporter
}

func (i *InitialEnterpriseImporterInteractorImpl) CreateInitialEnterpriseImporter(input CreateInitialEnterpriseImporterInput) (CreateInitialEnterpriseImporterOutput, error) {
	var (
		output CreateInitialEnterpriseImporterOutput
		err    error
	)

	// インポート情報を取得
	bookedImporterList, err := i.initialEnterpriseImporterRepository.GetByWeek()
	if err != nil {
		log.Println(err)
		return output, err
	}

	for _, importer := range bookedImporterList {
		if importer.StartDate == input.CreateParam.StartDate && importer.StartHour == input.CreateParam.StartHour {
			errMsg := "同じ日時でのインポート予約が既に存在します。"
			log.Println(errMsg)
			return output, errors.New(errMsg)
		}
	}

	// パスワードの暗号化
	input.CreateParam.Password, err = encrypt(input.CreateParam.Password)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 作成
	importer := entity.NewInitialEnterpriseImporter(
		input.CreateParam.AgentStaffID,
		input.CreateParam.ServiceType,
		input.CreateParam.LoginID,
		input.CreateParam.Password,
		input.CreateParam.StartDate,
		input.CreateParam.StartHour,
		input.CreateParam.MaxCount,
		input.CreateParam.Offset,
		input.CreateParam.SearchURL,
		false,
	)

	err = i.initialEnterpriseImporterRepository.Create(importer)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InitialEnterpriseImporter = importer

	return output, nil
}

// 削除
type DeleteInitialEnterpriseImporterInput struct {
	ID uint
}

type DeleteInitialEnterpriseImporterOutput struct {
	OK bool
}

func (i *InitialEnterpriseImporterInteractorImpl) DeleteInitialEnterpriseImporter(input DeleteInitialEnterpriseImporterInput) (DeleteInitialEnterpriseImporterOutput, error) {
	var (
		output DeleteInitialEnterpriseImporterOutput
		err    error
	)

	// 削除
	err = i.initialEnterpriseImporterRepository.Delete(input.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true
	return output, nil
}

// エージェントIDからエージェントロボットを取得する
type GetInitialEnterpriseImporterListByAgentIDInput struct {
	AgentID uint
}

type GetInitialEnterpriseImporterListByAgentIDOutput struct {
	InitialEnterpriseImporterList []*entity.InitialEnterpriseImporter
}

func (i *InitialEnterpriseImporterInteractorImpl) GetInitialEnterpriseImporterListByAgentID(input GetInitialEnterpriseImporterListByAgentIDInput) (GetInitialEnterpriseImporterListByAgentIDOutput, error) {
	var (
		output                        GetInitialEnterpriseImporterListByAgentIDOutput
		initialEnterpriseImporterList []*entity.InitialEnterpriseImporter
		err                           error
	)

	// インポート情報を取得
	initialEnterpriseImporterList, err = i.initialEnterpriseImporterRepository.GetByAgentID(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.InitialEnterpriseImporterList = initialEnterpriseImporterList

	return output, nil
}

// 今後1週間の企業求人一括インポート情報一覧を取得 *予約日程調整時の重複チェック用
type GetInitialEnterpriseImporterListByWeekOutput struct {
	InitialEnterpriseImporterList []*entity.InitialEnterpriseImporter
}

func (i *InitialEnterpriseImporterInteractorImpl) GetInitialEnterpriseImporterListByWeek() (GetInitialEnterpriseImporterListByWeekOutput, error) {
	var (
		output                        GetInitialEnterpriseImporterListByWeekOutput
		initialEnterpriseImporterList []*entity.InitialEnterpriseImporter
		err                           error
	)

	// インポート情報を取得
	initialEnterpriseImporterList, err = i.initialEnterpriseImporterRepository.GetByWeek()
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.InitialEnterpriseImporterList = initialEnterpriseImporterList

	return output, nil
}

/****************************************************************************************/
// Batch API
//

/*
他媒体から求人企業を一括インポートする バッチ
*/
type BatchInitialEnterpriseImporterInput struct {
	Now time.Time
}

type BatchInitialEnterpriseImporterOutput struct {
	OK bool
}

func (i *InitialEnterpriseImporterInteractorImpl) BatchInitialEnterpriseImporter(input BatchInitialEnterpriseImporterInput) (BatchInitialEnterpriseImporterOutput, error) {
	var (
		output           BatchInitialEnterpriseImporterOutput
		errMsg           string
		err              error
		selectedImporter entity.InitialEnterpriseImporter
		mailMessage      string
	)

	log.Println("interactor処理開始. input: ", input.Now)

	// 現在日時と実行時間が一致するインポート予約を取得
	importerList, err := i.initialEnterpriseImporterRepository.GetByStartDate(input.Now)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	log.Println("importerList", importerList)

	// タイムアウトを設定
	ctx, cancel := context.WithTimeout(context.Background(), 360*time.Minute)
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
			errMsg = fmt.Sprintf("求人取得処理でpanicエラーが発生しました。Recover: %v\nStack:\n%s", rec, errorLine)
			log.Println(errMsg)
			err = errors.New(errMsg)

			// エラーメール送信
			now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
			if strings.Contains(fmt.Sprint(rec), "deadline") {
				mailMessage = fmt.Sprintf(
					"RPA処理の中でタイムアウトが発生しました。\n発生時刻: %s\n対象ID: %d\n 媒体: %s\n 開始時間: %s\n取得開始位置: %d\n 最大取得件数: %d\n",
					now, selectedImporter.ID, "サーカスエージェント", fmt.Sprint(selectedImporter.StartDate, "-", selectedImporter.StartHour.Int64, ":00"), selectedImporter.Offset.Int64, selectedImporter.MaxCount.Int64,
				)
			} else {
				mailMessage = fmt.Sprintf(
					"RPA処理の中でパニックエラーが発生しました。\n発生時刻: %s\n対象ID: %d\n 媒体: %s\n 開始時間: %s\n取得開始位置: %d\n 最大取得件数: %d\n",
					now, selectedImporter.ID, "サーカスエージェント", fmt.Sprint(selectedImporter.StartDate, "-", selectedImporter.StartHour.Int64, ":00"), selectedImporter.Offset.Int64, selectedImporter.MaxCount.Int64,
				)
			}
			agentStaffs := []*entity.AgentStaff{}
			i.sendErrorMailToAgentStaffs(mailMessage, agentStaffs)
		}
	}()

ImporterLoop:
	for _, importer := range importerList {
		selectedImporter = *importer

		switch importer.ServiceType {

		// サーカス
		case null.NewInt(entity.EnterpriseImporterServiceTypeCircus, true):
			_, err = i.ImportEnterpriseFromCircus(ImportEnterpriseFromCircusInput{
				Param:   *importer,
				Context: ctx,
			})
			if err != nil {
				log.Println(err)
				mailMessage := fmt.Sprintf(
					"RPA処理の中でエラーが発生しました。\n対象ID: %d\n 媒体: サーカスエージェント\n 開始時間: %s\n取得開始位置: %d\n 最大取得件数: %d\n エラー詳細: %s\n",
					importer.ID, fmt.Sprint(importer.StartDate, "-", importer.StartHour.Int64, ":00"), importer.Offset.Int64, importer.MaxCount.Int64, err.Error(),
				)
				// エラーメール送信のためにエージェントスタッフを取得
				agentStaff, err := i.agentStaffRepository.FindByID(importer.AgentStaffID)
				if err != nil {
					log.Println(err)
					return output, err
				}
				agentStaffs, err := i.agentStaffRepository.GetByAgentID(agentStaff.AgentID)
				if err != nil {
					log.Println(err)
					return output, err
				}
				i.sendErrorMailToAgentStaffs(mailMessage, agentStaffs)
				return output, err
			}
			break ImporterLoop

			// エージェントバンク
		case null.NewInt(entity.EnterpriseImporterServiceTypeAgentBank, true):
			_, err = i.ImportEnterpriseFromAgentBank(ImportEnterpriseFromAgentBankInput{
				Param:   *importer,
				Context: ctx,
			})
			if err != nil {
				log.Println(err)
				mailMessage := fmt.Sprintf(
					"RPA処理の中でエラーが発生しました。\n対象ID: %d\n 媒体: エージェントバンク\n 開始時間: %s\n取得開始位置: %d\n 最大取得件数: %d\n エラー詳細: %s\n",
					importer.ID, fmt.Sprint(importer.StartDate, "-", importer.StartHour.Int64, ":00"), importer.Offset.Int64, importer.MaxCount.Int64, err.Error(),
				)
				// エラーメール送信のためにエージェントスタッフを取得
				agentStaff, err := i.agentStaffRepository.FindByID(importer.AgentStaffID)
				if err != nil {
					log.Println(err)
					return output, err
				}
				agentStaffs, err := i.agentStaffRepository.GetByAgentID(agentStaff.AgentID)
				if err != nil {
					log.Println(err)
					return output, err
				}
				i.sendErrorMailToAgentStaffs(mailMessage, agentStaffs)
				return output, err
			}
			break ImporterLoop
		}
	}
	log.Println("interactor処理終了. output: ", output)

	return output, nil
}

/*
サーカスから求人を取得する
*/
type ImportEnterpriseFromCircusInput struct {
	Param   entity.InitialEnterpriseImporter
	Context context.Context
}

type ImportEnterpriseFromCircusOutput struct {
	OK bool `json:"ok"`
}

func (i *InitialEnterpriseImporterInteractorImpl) ImportEnterpriseFromCircus(input ImportEnterpriseFromCircusInput) (ImportEnterpriseFromCircusOutput, error) {
	var (
		output         ImportEnterpriseFromCircusOutput
		browser        *rod.Browser
		page           *rod.Page
		errMsg         string
		err            error
		agentStaff     *entity.AgentStaff
		enterpriseList []*entity.EnterpriseAndJobInformation
		now            string = time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
		jobCnt         int    // ヒットした求人数
	)

	if input.Context == nil {
		errMsg = "Contextが指定されていません"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	agentStaff, err = i.agentStaffRepository.FindByID(input.Param.AgentStaffID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// ホストOSのchrome, chromium or microsoft edgeのパスを取得
	path, ok := launcher.LookPath()
	if !ok {
		errMsg = "chrome, chromium or microsoft edgeが見つかりません"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	// 取得したパスからchromeを起動
	l := launcher.New().
		Bin(path)

	// ローカル環境の場合は画面ありで起動
	if os.Getenv("APP_ENV") == "local" {
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

	// ログインページへ遷移
	page = browser.
		MustPage("https://circus-job.com/login?redirectTo=/login").
		MustWaitLoad()

	// ログインIDとパスワードを入力
	page.
		MustElement("input[name=email]").
		MustInput(input.Param.LoginID)

	// パスワードの複号
	input.Param.Password, err = decryption(input.Param.Password)
	if err != nil {
		log.Println(err)
		return output, err
	}

	page.
		MustElement("input[name=password]").
		MustInput(input.Param.Password)

	buttons := page.MustElements("button")
	isMatchedWithButton := false
	for _, button := range buttons {
		log.Println("button.MustText(): ", button.MustText())
		if button.MustText() == "ログインする" {
			button.MustClick()
			page.WaitLoad()
			isMatchedWithButton = true
		}
	}
	if !isMatchedWithButton {
		errMsg = "ログインボタンが見つかりません"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	page.WaitLoad()
	time.Sleep(4 * time.Second)

	// URLがログインページのままの場合はもう一度ログイン
	if strings.Contains(page.MustInfo().URL, "login") {
		button := page.MustElement("button.MuiButton-root.MuiButton-contained.MuiButton-containedConversion.MuiButton-sizeMedium.MuiButton-containedSizeMedium.MuiButtonBase-root.css-w0x5z")
		log.Println("button.MustText(): ", button.MustText())
		button.MustClick()
		page.WaitLoad()
		time.Sleep(10 * time.Second)
	}

	// URLがログインページのままの場合はログインに失敗していると判断
	if strings.Contains(page.MustInfo().URL, "login") {
		errMsg = "ログインに失敗しました。"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	// url(input.Param.SearchURL)で検索条件を指定
	// デフォルトURL: "https://circus-job.com/search?qJson=%5B%7B%22option%22%3A1%2C%22keyword%22%3A%22%22%2C%22logicType%22%3A%22and%22%7D%2C%7B%22option%22%3A1%2C%22keyword%22%3A%22%22%2C%22logicType%22%3A%22or%22%7D%2C%7B%22option%22%3A1%2C%22keyword%22%3A%22%22%2C%22logicType%22%3A%22excludeAnd%22%7D%2C%7B%22option%22%3A1%2C%22keyword%22%3A%22%22%2C%22logicType%22%3A%22excludeOr%22%7D%5D&resumePassRateIncludingDuringMeasurement=included&selectionDaysIncludingDuringMeasurement=included"
	// ソート順の選択
	// 掲載開始日が新しい順: &orderBy=publishStartedAt&order=desc
	// 掲載開始日が古い順: &orderBy=publishStartedAt&order=asc

	// 求人タイプの選択
	// 自社求人: &jobPostOwners=1
	// 事務局求人: &jobPostOwners=2
	// 企業掲載求人: &jobPostOwners=3
	// シェア求人: &jobPostOwners=4

	// 新規求人一覧: page.MustNavigate("https://circus-job.com/new-jobs").
	// クローズド求人一覧: page.MustNavigate("https://circus-job.com/closed-jobs").
	// RAに切り替え: https://circus-job.com/?requiredTheme=RA

	// 取得開始ページ
	// デフォルトは1ページ目
	startPageI := 1
	// 開始指定のある場合. 例:101件目から取得する場合(Offset:101)は5ページ目から取得
	if input.Param.Offset.Valid && input.Param.Offset.Int64 > 0 {
		startPageI += int(input.Param.Offset.Int64 / 25)
	}

pageLoop:
	for pageI := startPageI; pageI <= 50; pageI++ {
		// 指定ページへ遷移
		page.MustNavigate(input.Param.SearchURL + "&page=" + strconv.Itoa(pageI)).
			WaitLoad()
		time.Sleep(4 * time.Second)

		if !page.MustHas("div.JobSearchResultCard-root") {
			log.Println("該当する求人がありません. ページ: ", pageI)
			break
		}

		// JobPostHead-jobTitle
		jobDivs := page.MustElements("div.JobSearchResultCard-root")
		log.Println("len(jobDivs): ", len(jobDivs))
		jobDivsLen := len(jobDivs)
		for jobI := 0; jobI < jobDivsLen; jobI++ {
			log.Println("jobI: ", jobI)
			// 求人数の合計がmaxCountを超えたら終了
			jobCnt++
			if jobCnt > int(input.Param.MaxCount.Int64) {
				break pageLoop

				// ローカル環境の場合は5件のみ取得
			} else if os.Getenv("APP_ENV") == "local" {
				log.Println("jobI: ", jobI)
				if jobI > 7 {
					break pageLoop
				}
			}

			// 企業情報を初期化
			enterprise := entity.EnterpriseAndJobInformation{
				AgentStaffID:  input.Param.AgentStaffID,
				RegisterPhase: null.NewInt(1, true),
				SecretMemo:    "・取り込み媒体: サーカスエージェント\n\n・取り込み日時: " + now,
			}

			// 指定要素までスクロール *要素までスクロールしないと詳細がレンダリングされないため
			jobDiv := jobDivs[jobI]
			jobDiv.MustScrollIntoView()
			page.Mouse.MustScroll(0, 100)
			jobDiv = page.MustElements("div.JobSearchResultCard-root")[jobI]

			jobTitleDiv := jobDiv.MustElement("div.JobPostHead-jobTitle")
			enterprise.Title = jobTitleDiv.MustText()
			log.Println("jobTitle: ", enterprise.Title)

			jobTitleA := jobTitleDiv.MustElement("a")
			detailPageURLMemory, err := jobTitleA.Attribute("href")
			if err != nil {
				errMsg = "求人詳細URLの取得に失敗しました。求人名: " + enterprise.Title
				log.Println(errMsg)
				return output, errors.New(errMsg)
			}
			detailPageURL := *detailPageURLMemory
			log.Println("detailPageURL: ", detailPageURL)

			detailPage := browser.MustPage("https://circus-job.com" + detailPageURL)
			detailPage.WaitLoad()
			time.Sleep(4 * time.Second)

			// 株式会社カオナビ　= MuiTypography-root MuiTypography-inherit MuiLink-root MuiLink-underlineNone Link-textLink css-3vsmdf
			companyNameAs := detailPage.MustElements("a.MuiTypography-root.MuiTypography-inherit.MuiLink-root.MuiLink-underlineNone.Link-textLink.css-3vsmdf")
			for _, companyNameA := range companyNameAs {
				companyNameAText := companyNameA.MustText()
				log.Println("companyNameAText: ", companyNameAText)
				if companyNameAText == "" || (!strings.Contains(companyNameAText, "株式会社") && !strings.Contains(companyNameAText, "有限会社") && !strings.Contains(companyNameAText, "合同会社") && !strings.Contains(companyNameAText, "合名会社") && !strings.Contains(companyNameAText, "合資会社")) {
					continue
				}
				println("companyNameAText: ", companyNameAText)
				enterprise.CompanyName = companyNameAText
				if enterprise.CompanyName == "" {
					errMsg = "企業名の取得に失敗しました。求人名: " + enterprise.Title
					log.Println(errMsg)
					return output, errors.New(errMsg)
				}

				enterprise.BillingAddressTitle = enterprise.CompanyName + " 請求先1"
				corporateURLMemory, err := companyNameA.Attribute("href")
				if err != nil {
					errMsg = "企業URLの取得に失敗しました。求人名: " + enterprise.Title
					log.Println(errMsg)
					return output, errors.New(errMsg)
				}
				enterprise.CorporateSiteURL = *corporateURLMemory
				log.Println("corporateURL: ", enterprise.CorporateSiteURL)
				break
			}

			// ヘッダーチップ
			headerChips := detailPage.MustElements("div.JobPostHead-chipWrapper")
			enterprise.SecretMemo += "\n\n・付属情報:"
			for _, headerChip := range headerChips {
				headerChipText := headerChip.MustText()
				log.Println("headerChipText: ", headerChipText)
				enterprise.SecretMemo += headerChipText + "、"

				// 雇用形態
				outputNumber := entity.GetIntEmploymentStatusForCircus(headerChipText)
				if outputNumber.Valid && outputNumber.Int64 != 99 {
					enterprise.EmploymentStatuses = append(enterprise.EmploymentStatuses, entity.JobInformationEmploymentStatus{
						EmploymentStatus: outputNumber,
					})
				}

			}

			// （web・インターネット・ゲーム、人材ビジネス） -> web・インターネット・ゲーム, 人材ビジネス
			industriesPlainText := detailPage.MustElement("div.JobPostHead-industries").MustText()
			log.Println("industriesPlainText: ", industriesPlainText)

			enterprise.BusinessDetail += "\n\n業界: " + industriesPlainText

			replaced := strings.Replace(industriesPlainText, "（", "", -1)
			replaced = strings.Replace(replaced, "）", "", -1)
			splitedIndustries := strings.Split(replaced, "、")
			log.Println("splitedIndustries: ", splitedIndustries)
			for _, industryText := range splitedIndustries {
				industry1, industry2 := entity.GetIntIndustryForCircus(industryText)
				if industry1.Valid {
					enterprise.Industries = append(enterprise.Industries, industry1)
				}

				if industry2.Valid {
					enterprise.Industries = append(enterprise.Industries, industry2)
				}
			}

			jobSummaryDivs := detailPage.MustElement("div.JobSummaryTable-root").MustElements("div.MuiBox-root.css-0")
			for jobSummaryI, jobSummaryDiv := range jobSummaryDivs {
				if jobSummaryI == 1 {
					// 職種 = カスタマーサポート・ヘルプデスク、カスタマーサポート・ヘルプデスク、法人営業
					occupationPlainText := jobSummaryDiv.MustText()
					log.Println("occupationPlainText: ", occupationPlainText)
					enterprise.SecretMemo += "\n\n・職種: " + occupationPlainText
					splitedOccupations := strings.Split(occupationPlainText, "、")
					for _, occupationText := range splitedOccupations {
						occupation1, occupation2 := entity.GetIntOccupationForCircus(occupationText)
						if occupation1.Valid {
							enterprise.Occupations = append(enterprise.Occupations, entity.JobInformationOccupation{
								Occupation: occupation1,
							})
						}

						if occupation2.Valid {
							enterprise.Occupations = append(enterprise.Occupations, entity.JobInformationOccupation{
								Occupation: occupation2,
							})
						}
					}

				} else if jobSummaryI == 3 {
					// 職位 = リーダー、メンバー
					position := jobSummaryDiv.MustText()
					log.Println("position: ", position)
					enterprise.SecretMemo += "\n\n・ポジション:" + position

				} else if jobSummaryI == 5 {
					// 年収　= 400万円〜600万円
					annualIncome := jobSummaryDiv.MustText()
					log.Println("annualIncome: ", annualIncome)

					enterprise.Salary += "\n\n・年収:" + annualIncome

				} else if jobSummaryI == 7 {
					// 勤務地 = 東京
					workLocation := jobSummaryDiv.MustText()
					log.Println("workLocation: ", workLocation)

					enterprise.WorkLocation += "\n\n都道府県:" + workLocation
					// splite with "、"
					splitedWorkLocations := strings.Split(workLocation, "、")
					for _, splitedWorkLocation := range splitedWorkLocations {
						prefectureInt := entity.GetIntPrefectureForCircus(splitedWorkLocation)
						if prefectureInt.Valid {
							enterprise.Prefectures = append(enterprise.Prefectures, entity.JobInformationPrefecture{
								Prefecture: prefectureInt,
							})
						}
					}
				}
			}

			// 企業情報
			jobCompanyInformationTableTds := detailPage.MustElements("td.JobCompanyInformationTable-root.MuiTableCell-root.MuiTableCell-body.MuiTableCell-sizeSmall.css-wfyo7z")
			for jobCompanyInformationTableTdI, jobCompanyInformationTableTd := range jobCompanyInformationTableTds {
				if jobCompanyInformationTableTdI == 0 {
					// 住所 = 東京都渋谷区
					address := jobCompanyInformationTableTd.MustText()
					log.Println("address: ", address)
					enterprise.OfficeLocation += address
					enterprise.HowToRecommend += "\n\n・本拠所在地:" + address

				} else if jobCompanyInformationTableTdI == 1 {
					// 設立年月 = 2017年
					establishmentDate := jobCompanyInformationTableTd.MustText()
					log.Println("establishmentDate: ", establishmentDate)
					enterprise.BusinessDetail += "\n\n・設立年月:" + establishmentDate

				} else if jobCompanyInformationTableTdI == 2 {
					// 企業の規模 = ベンチャー、中小
					companyPhase := jobCompanyInformationTableTd.MustText()
					log.Println("companyPhase: ", companyPhase)
					enterprise.BusinessDetail += "\n\n・企業の規模:" + companyPhase
					enterprise.SecretMemo += "\n\n・企業の規模:" + companyPhase

				} else if jobCompanyInformationTableTdI == 3 {
					// 従業員数 = 1〜9人
					numberOfEmployees := jobCompanyInformationTableTd.MustText()
					log.Println("numberOfEmployees: ", numberOfEmployees)
					enterprise.BusinessDetail += "\n\n・従業員数:" + numberOfEmployees

				} else if jobCompanyInformationTableTdI == 4 {
					// 上場 = 未上場
					publicOffering := jobCompanyInformationTableTd.MustText()
					log.Println("publicOffering: ", publicOffering)
					enterprise.BusinessDetail += "\n\n・上場:" + publicOffering
					enterprise.PublicOffering = entity.GetIntPublicOfferingForCircus(publicOffering)

				} else if jobCompanyInformationTableTdI == 5 {
					// 平均年齢 = 20代
					averageAge := jobCompanyInformationTableTd.MustText()
					log.Println("averageAge: ", averageAge)

					enterprise.BusinessDetail += "\n\n・平均年齢:" + averageAge + "歳"
					enterprise.SecretMemo += "\n\n・平均年齢:" + averageAge + "歳"

				} else if jobCompanyInformationTableTdI == 6 {
					// 男女比 = 男性 100% 女性 0%
					parcentageOfGender := jobCompanyInformationTableTd.MustText()
					log.Println("parcentageOfGender: ", parcentageOfGender)

					enterprise.BusinessDetail += "\n\n・男女比:" + parcentageOfGender
					enterprise.SecretMemo += "\n\n・男女比:" + parcentageOfGender
				}
			}

			//JobRequiredTable-root MuiTableCell-root MuiTableCell-body MuiTableCell-sizeSmall css-1wyk6sn
			jobRequiredTableDivs := detailPage.MustElements("td.JobRequiredTable-root.MuiTableCell-root.MuiTableCell-body.MuiTableCell-sizeSmall.css-1wyk6sn")
			// 1 = 応募必須条件
			requiredChips := detailPage.MustElements("div.RequiredApplicationChips-chip")
			log.Println("requiredChips: ", len(requiredChips))
			for _, requiredChip := range requiredChips {
				requiredText := requiredChip.MustText()
				enterprise.RequiredExperienceJobDetail += "\n・" + requiredText
				log.Println("requiredText: ", requiredText)

				outputNumber := entity.GetIntFinalEducationForCircus(requiredText)
				if outputNumber.Valid {
					enterprise.FinalEducation = outputNumber
				}
				// 休日
				outputNumber = entity.GetIntHolidayForCircus(requiredText)
				if outputNumber.Valid {
					enterprise.HolidayType = outputNumber
				}
				// 国籍
				outputNumber = entity.GetIntNationalityForCircus(requiredText)
				if outputNumber.Valid {
					enterprise.Nationality = outputNumber
				}
				// 性別
				outputNumber = entity.GetIntGenderForCircus(requiredText)
				if outputNumber.Valid {
					enterprise.Gender = outputNumber
				}
			}
			// MuiTypography-root MuiTypography-body1 css-1ioozpd
			enterprise.RequiredExperienceJobDetail += "\n\n・その他条件:" + detailPage.MustElement("p.MuiTypography-root.MuiTypography-body1.css-1ioozpd").MustText()
			for jobRequiredTableDivI, jobRequiredTableDiv := range jobRequiredTableDivs {
				if jobRequiredTableDivI == 1 {
					continue

					//書類見送りの主な理由
				} else if jobRequiredTableDivI == 2 {
					enterprise.SecretMemo += "\n\n・書類見送りの主な理由:" + jobRequiredTableDiv.MustText()
					// 内定の可能性が高い人
				} else if jobRequiredTableDivI == 3 {
					enterprise.SecretMemo += "\n\n・内定の可能性が高い人:" + jobRequiredTableDiv.MustText()
					// 面接見送りの主な理由
				} else if jobRequiredTableDivI == 4 {
					enterprise.SecretMemo += "\n\n・面接見送りの主な理由:" + jobRequiredTableDiv.MustText()
				}
			}

			// MuiPaper-root MuiPaper-elevation MuiPaper-rounded MuiPaper-elevation3 JobSearchDetailsCard-root css-zf553n
			detailCards := detailPage.MustElements("div.JobSearchDetailsCard-root")
			log.Println("detailCards", len(detailCards))
			if len(detailCards) < 3 {
				errMsg = "detailCardsの長さが3未満です。"
				log.Println(errMsg)
				return output, errors.New("detailCards is less than 3")
			}

			// エージェント向け情報
			infoForAgentDiv := detailCards[1]
			infoForAgent := infoForAgentDiv.MustText()
			enterprise.SecretMemo += "\n\n・エージェント向け情報:" + infoForAgent

			// 求人概要
			// MuiPaper-root MuiPaper-elevation MuiPaper-rounded MuiPaper-elevation3 JobSearchDetailsCard-root css-zf553n
			jobSearchDetailsCardDiv := detailCards[2]
			jobSearchDetailsCard := jobSearchDetailsCardDiv.MustText()
			enterprise.SecretMemo += "\n\n・求人概要:" + jobSearchDetailsCard

			jobDetailSections := jobSearchDetailsCardDiv.MustElements("section.css-1ns80lu")
			for jobDetailSectionI, jobDetailSection := range jobDetailSections {
				if jobDetailSectionI == 0 {
					JobDetailTrs := jobDetailSection.MustElements("tr")
					if len(JobDetailTrs) == 0 {
						errMsg = "JobDetailTrs is less than 1"
						log.Println(errMsg)
						return output, errors.New(errMsg)
					}
					for _, JobDetailTr := range JobDetailTrs {
						jonDetailThTxt := JobDetailTr.MustElement("th").MustText()
						jobDetailTd := JobDetailTr.MustElement("td")
						jobDetailTdTxt := jobDetailTd.MustText()
						log.Println("jonDetailThTxt", jonDetailThTxt, "jobDetailTdTxt", jobDetailTdTxt)
						switch jonDetailThTxt {
						case "":
							continue

						case "事業内容と今後の事業展開":
							enterprise.BusinessDetail += "\n\n・事業内容と今後の事業展開:" + jobDetailTdTxt
							enterprise.WorkDetail += "\n\n・事業内容と今後の事業展開:" + jobDetailTdTxt

						// 増員募集 / 更なる組織強化
						case "募集背景":
							enterprise.SecretMemo += "\n\n・募集背景:" + jobDetailTdTxt

						// 募集予定人数・期間　: 1人 / 3ヶ月
						case "募集予定人数・期間":
							// 募集人数
							enterprise.SecretMemo += "\n\n・募集予定人数・期間:" + jobDetailTdTxt
							splitedRecruitNumAndPeriod := strings.Split(jobDetailTdTxt, " / ")
							if len(splitedRecruitNumAndPeriod) == 2 {
								recruitNamStr := splitedRecruitNumAndPeriod[0]
								// remove "人"
								recruitNamStr = strings.Replace(recruitNamStr, "人", "", -1)
								recruitNum, err := strconv.Atoi(recruitNamStr)
								if err == nil {
									enterprise.NumberOfHires = null.NewInt(int64(recruitNum), true)
								}

								// 募集期間
								// recruitPeriod, err := strconv.Atoi(splitedRecruitNumAndPeriod[1])
								// if err == nil {
								// 	enterprise.
								// }
							}

						// PRポイント
						case "PRポイント":
							enterprise.SecretMemo += "\n\n・PRポイント:" + jobDetailTdTxt

							// 仕事内容
						case "仕事内容":
							enterprise.WorkDetail += "\n\n・仕事内容:" + jobDetailTdTxt

						case "現在の組織構成":
							enterprise.WorkDetail += "\n\n・現在の組織構成:" + jobDetailTdTxt

							// 給与・年収例
						case "給与・年収例":
							enterprise.Salary = jobDetailTdTxt

							// 勤務地・勤務時間
						case "勤務地・勤務時間":
							if jobDetailTd.MustHas("div.MuiBox-root.css-11ze7cv") {
								separatedInfo := jobDetailTd.MustElements("div.MuiBox-root.css-11ze7cv")
								if len(separatedInfo) != 2 {
									continue
								}
								// 勤務地
								enterprise.WorkLocation += separatedInfo[0].MustText()
								// 勤務時間
								enterprise.WorkTime += separatedInfo[1].MustText()

								// 転勤の可能性: なし, あり, 可能性はあるが、最大限考慮
								if strings.Contains(enterprise.WorkLocation, "転勤の可能性") {
									if strings.Contains(enterprise.WorkLocation, "転勤の可能性：なし") {
										enterprise.Transfer = null.NewInt(0, true)
									} else if strings.Contains(enterprise.WorkLocation, "転勤の可能性：あり") {
										enterprise.Transfer = null.NewInt(1, true)
									} else if strings.Contains(enterprise.WorkLocation, "転勤の可能性：可能性はあるが、最大限考慮") {
										enterprise.Transfer = null.NewInt(2, true)
									}
								}

								// 月間平均残業時間: 指定しない, 残業なし, 10時間以下, 20時間以下, 30時間以下, 40時間以下, 50時間以下
								if strings.Contains(enterprise.WorkTime, "月間平均残業時間：") {
									if strings.Contains(enterprise.WorkTime, "月間平均残業時間：残業なし") {
										enterprise.OvertimeAverage = "月間平均残業時間：残業なし"
									} else if strings.Contains(enterprise.WorkTime, "月間平均残業時間：10時間以下") {
										enterprise.OvertimeAverage = "月間平均残業時間：10時間以下"
									} else if strings.Contains(enterprise.WorkTime, "月間平均残業時間：20時間以下") {
										enterprise.OvertimeAverage = "月間平均残業時間：20時間以下"
									} else if strings.Contains(enterprise.WorkTime, "月間平均残業時間：30時間以下") {
										enterprise.OvertimeAverage = "月間平均残業時間：30時間以下"
									} else if strings.Contains(enterprise.WorkTime, "月間平均残業時間：40時間以下") {
										enterprise.OvertimeAverage = "月間平均残業時間：40時間以下"
									} else if strings.Contains(enterprise.WorkTime, "月間平均残業時間：50時間以下") {
										enterprise.OvertimeAverage = "月間平均残業時間：50時間以下"
									}
								}
							}

							// 補足情報
							if jobDetailTd.MustHas("div.MuiBox-root.css-15bcdd2") {
								supInfo := jobDetailTd.MustElement("div.MuiBox-root.css-15bcdd2").MustText()
								enterprise.WorkDetail += supInfo
								enterprise.WorkTime += supInfo
							}

							// 休日休暇・福利厚生
						case "休日休暇・福利厚生":
							if jobDetailTd.MustHas("div.MuiBox-root.css-11ze7cv") {
								separatedInfo := jobDetailTd.MustElements("div.MuiBox-root.css-11ze7cv")
								if len(separatedInfo) != 2 {
									continue
								}
								// 休日休暇
								enterprise.HolidayDetail += separatedInfo[0].MustText()
								// 福利厚生
								enterprise.Insurance += separatedInfo[1].MustText()

								// 社会保険完備の記載がある場合は、各保険にチェックを入れる
								if strings.Contains(enterprise.Insurance, "社会保険完備") {
									enterprise.EmploymentInsurance = true
									enterprise.AccidentInsurance = true
									enterprise.HealthInsurance = true
									enterprise.PensionInsurance = true
								}

								if strings.Contains(enterprise.HolidayDetail, "休日：") {
									if strings.Contains(enterprise.HolidayDetail, "休日：土日祝休み") {
										enterprise.HolidayType = entity.GetIntHolidayForCircus("土日祝休み")
									} else if strings.Contains(enterprise.HolidayDetail, "休日：土日休み") {
										enterprise.HolidayType = entity.GetIntHolidayForCircus("土日休み")
									} else if strings.Contains(enterprise.HolidayDetail, "休日：週休2日（土日以外）") {
										enterprise.HolidayType = entity.GetIntHolidayForCircus("週休2日（土日以外）")
									} else if strings.Contains(enterprise.HolidayDetail, "休日：シフト制") {
										enterprise.HolidayType = entity.GetIntHolidayForCircus("シフト制")
									} else if strings.Contains(enterprise.HolidayDetail, "休日：その他") {
										enterprise.HolidayType = entity.GetIntHolidayForCircus("その他")
									}
								}

							}
							// 補足情報
							if jobDetailTd.MustHas("div.MuiBox-root.css-15bcdd2") {
								supInfo := jobDetailTd.MustElement("div.MuiBox-root.css-15bcdd2").MustText()
								// 休日休暇
								enterprise.HolidayDetail += supInfo
								// 福利厚生
								enterprise.Insurance += supInfo
							}

							// 受動喫煙対策
						case "受動喫煙対策":
							continue
							// smokingInfo := jobDetailTd.MustText()
							// if strings.Contains(smokingInfo, "受動喫煙対策の有無 : 対策あり") {
							// 	enterprise.PassiveSmoking = null.NewInt(0, true)
							// } else {
							// 	enterprise.PassiveSmoking = null.NewInt(1, true)
							// }

						default:
							continue
						}
					}

					// 選考情報
				} else if jobDetailSectionI == 1 {
					jobDetailSectionTxt := jobDetailSection.MustText()
					enterprise.SelectionFlow = jobDetailSectionTxt

					// 成約手数料
				} else if jobDetailSectionI == 2 {
					jobDetailSectionTxt := jobDetailSection.MustText()
					enterprise.CommissionDetail = jobDetailSectionTxt

					// commisionTrs := jobDetailSection.MustElements("tr")
					// for _, commisionTr := range commisionTrs {
					// 	commisionTh := commisionTr.MustElement("th")
					// 	commisionThTxt := commisionTh.MustText()
					// 	commisionTd := commisionTr.MustElement("td")
					// 	commisionTdTxt := commisionTd.MustText()

					// 	switch commisionThTxt {
					// 	case "成果報酬金額":
				}
			}

			log.Println("jobSearchDetailsCardDiv", jobSearchDetailsCard[0:20])

			detailPage.MustClose()
			time.Sleep(4 * time.Second)

			enterpriseList = append(enterpriseList, &enterprise)
		}
	}

	// 被り項目削除 *業界, 職種など
	for _, enterprise := range enterpriseList {
		log.Println("enterprise company name: ", enterprise.CompanyName)
		// 業界
		m := make(map[null.Int]bool)
		uniqueIndustries := []null.Int{}

		for _, ele := range enterprise.Industries {
			if !m[ele] {
				m[ele] = true
				uniqueIndustries = append(uniqueIndustries, ele)
			}
		}
		enterprise.Industries = uniqueIndustries

		// 職種
		m = make(map[null.Int]bool)
		uniqueOccupations := []entity.JobInformationOccupation{}
		for _, ele := range enterprise.Occupations {
			if !m[ele.Occupation] {
				m[ele.Occupation] = true
				uniqueOccupations = append(uniqueOccupations, ele)
			}
		}
		enterprise.Occupations = uniqueOccupations

	}

	// if input.Param.IsJSON {

	// 	// JSON ファイルとして出力
	// 	file, err := os.Create("circus-enterprise-" + now + ".json")
	// 	if err != nil {
	// 		log.Println(err)
	// 		return output, err
	// 	}
	// 	defer file.Close()

	// 	// JSON ファイルに書き込み
	// 	encoder := json.NewEncoder(file)
	// 	if err := encoder.Encode(output); err != nil {
	// 		log.Println(err)
	// 		return output, err
	// 	}
	// } else {

	// 企業情報を作成する
	err = i.createEnterpriseList(enterpriseList, agentStaff)
	if err != nil {
		log.Println("err: ", err)
		errMsg = "企業情報の作成に失敗しました"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}
	// }

	// isSuccessをtrueにする
	err = i.initialEnterpriseImporterRepository.UpdateIsSuccessTrue(input.Param.ID)
	if err != nil {
		log.Println("err: ", err)
		errMsg = "成功フラグ(isSuccess)の更新に失敗しました"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	log.Println("interactor処理終了. output: ", output)

	output.OK = true

	return output, nil
}

/*
エージェントバンクから求人を取得する
*/
type ImportEnterpriseFromAgentBankInput struct {
	Param   entity.InitialEnterpriseImporter
	Context context.Context
}

type ImportEnterpriseFromAgentBankOutput struct {
	OK bool `json:"ok"`
}

func (i *InitialEnterpriseImporterInteractorImpl) ImportEnterpriseFromAgentBank(input ImportEnterpriseFromAgentBankInput) (ImportEnterpriseFromAgentBankOutput, error) {
	var (
		output  ImportEnterpriseFromAgentBankOutput
		browser *rod.Browser
		page    *rod.Page
		errMsg  string
		err     error
		// agentStaff     *entity.AgentStaff
		enterpriseList []*entity.EnterpriseAndJobInformation
		now            string = time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
		jobCnt         int    // ヒットした求人数
	)

	/************ 1. ブラウザ起動 **************/

	if input.Context == nil {
		errMsg = "Contextが指定されていません"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	// ホストOSのchrome, chromium or microsoft edgeのパスを取得
	path, ok := launcher.LookPath()
	if !ok {
		errMsg = "chrome, chromium or microsoft edgeが見つかりません"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	// 取得したパスからchromeを起動
	l := launcher.New().
		Bin(path)

	// ローカル環境の場合は画面ありで起動
	if os.Getenv("APP_ENV") == "local" {
		l.Headless(true)
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

	/************ 2. agentBankへログイン **************/

	// ログインページへ遷移
	page = browser.
		MustPage("https://agent-bank.com/service/login").
		MustWaitLoad()

	// ログインIDとパスワードを入力
	page.
		MustElement("input[name=email]").
		MustInput(input.Param.LoginID)

	// パスワードの複号
	input.Param.Password, err = decryption(input.Param.Password)
	if err != nil {
		log.Println(err)
		return output, err
	}

	page.
		MustElement("input[name=password]").
		MustInput(input.Param.Password)

	buttons := page.MustElements("button")
	isMatchedWithButton := false
	for _, button := range buttons {
		log.Println("button.MustText(): ", button.MustText())
		if button.MustText() == "ログイン" {
			button.MustClick()
			page.WaitLoad()
			isMatchedWithButton = true
		}
	}
	if !isMatchedWithButton {
		errMsg = "ログインボタンが見つかりません"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	page.WaitLoad()
	time.Sleep(10 * time.Second)

	// URLがログインページのままの場合はもう一度ログイン
	if strings.Contains(page.MustInfo().URL, "login") {
		page.
			MustElement("input[name=password]").
			MustInput(input.Param.Password)

		buttons := page.MustElements("button")
		isMatchedWithButton := false
		for _, button := range buttons {
			log.Println("button.MustText(): ", button.MustText())
			if button.MustText() == "ログイン" {
				button.MustClick()
				page.WaitLoad()
				isMatchedWithButton = true
			}
		}
		if !isMatchedWithButton {
			errMsg = "ログインボタンが見つかりません"
			log.Println(errMsg)
			return output, errors.New(errMsg)
		}
		page.WaitLoad()
		time.Sleep(10 * time.Second)
	}

	// URLがログインページのままの場合はログインに失敗していると判断
	if strings.Contains(page.MustInfo().URL, "login") {
		errMsg = "ログインに失敗しました。"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	/************ 3. 開始ページと終了ページの定義 + 求人検索ページへ遷移 **************/

	// url(input.Param.SearchURL)で検索条件を指定
	var (
		startPageI   = 1
		lastPageI    = 1
		perPageCount = 15.0 // 1ページあたりの件数

		// 数字部分を抽出する正規表現
		regexpNumber = regexp.MustCompile(`\d+`)
	)
	// 開始指定のある場合. 例:24ページから取得する場合(Offset:24)から取得
	if input.Param.Offset.Valid && input.Param.Offset.Int64 > 0 {
		// 1~27ページ 405件
		// 28~54ページ 405件
		input.Param.SearchURL += fmt.Sprintf("page=%v", input.Param.Offset.Int64)
	}

	// agentBankは1ページあたり15件の求人が表示されている（取得は800件）
	if input.Param.MaxCount.Valid && float64(input.Param.MaxCount.Int64) > perPageCount {
		lastPageI = int(math.Ceil(float64((input.Param.MaxCount.Int64)) / perPageCount))
	}

	fmt.Printf("取得範囲: %v ~ %v\n", startPageI, lastPageI)

	// 指定ページへ遷移
	page.MustNavigate(input.Param.SearchURL)
	page.WaitLoad()
	time.Sleep(4 * time.Second)

	/************ 4. ページ遷移のループ **************/

pageLoop:
	for pageI := startPageI; pageI <= lastPageI; pageI++ {

		/************ 4-1. ページ遷移（1ページ目以外） **************/

		if pageI != startPageI {
			fmt.Println("ページ遷移")
			page.MustElement("button.btn-next").MustClick().WaitLoad()
			time.Sleep(2 * time.Second)
		}

		if !page.MustHas("div.s-job-item") {
			log.Println("該当する求人がありません. ページ: ", pageI)
			break
		}

		/************ 4-2. ページ内の求人コンテンツの要素を全て取得 **************/

		jobDivs := page.MustElements("div.s-job-item")
		jobDivsLen := len(jobDivs)

		/************ 4-3. 求人コンテンツのループ **************/

		for jobI := 0; jobI < jobDivsLen; jobI++ {
			// 現在の求人コンテンツを変数化
			var jobDiv = jobDivs[jobI]

			/************ 4-3-1. 取得求人のカウント（条件次第でループの終了） **************/

			// 求人数の合計がmaxCountを超えたら終了
			jobCnt++
			if jobCnt > int(input.Param.MaxCount.Int64) {
				break pageLoop

				// ローカル環境の場合は2Pのみ取得
			} else if os.Getenv("APP_ENV") == "local" {
				log.Println("jobI: ", jobI, "pageI: ", pageI)
				// if pageI == 2 {
				// 	break pageLoop
				// }
			}

			/************ 4-3-2. 企業情報を初期化 **************/

			enterprise := entity.EnterpriseAndJobInformation{
				AgentStaffID:                  input.Param.AgentStaffID,
				AgentStaffIDForBillingAddress: input.Param.AgentStaffID,
				SecretMemo:                    fmt.Sprintf("・取り込み媒体: エージェントバンク\n\n・取り込み日時: %s\n\n", now),
				RegisterPhase:                 null.NewInt(0, true),                                          // 登録状況: 本登録
				RecruitmentState:              null.NewInt(0, true),                                          // 募集状況: 募集中
				Targets:                       []entity.JobInformationTarget{{Target: null.NewInt(0, true)}}, // 募集対象: 中途
				EmploymentInsurance:           true,                                                          // 雇用保険
				AccidentInsurance:             true,                                                          // 労災保険
				HealthInsurance:               true,                                                          // 健康保険
				PensionInsurance:              true,                                                          // 厚生年金保険
				ContractPhase:                 null.NewInt(2, true),                                          // 契約フェーズ: 契約締結済み
				BillingAddressTitle:           "agent bank",                                                  // 請求先のタイトル
				ExternalType:                  null.NewInt(0, true),                                          // 他媒体タイプ: agentBank
				DriverLicence:                 null.NewInt(99, true),                                         // 普通自動車免許: 不要（99）← agentBankで指定がある場合は更新される
			}

			/************ 4-3-3. 求人タイトルの取得 + 求人詳細ページへの遷移 + 媒体URLの取得 **************/

			// 求人タイトルの取得
			titleAndLink, err := jobDiv.Element("a.title")
			if err != nil {
				errMsg = "求人タイトルの取得に失敗しました。"
				log.Println(errMsg)
				return output, errors.New(errMsg)
			}
			enterprise.Title = titleAndLink.MustText()

			// 求人詳細ページへ遷移
			detailPageURLMemory, err := titleAndLink.Attribute("href")
			if err != nil {
				errMsg = "求人詳細URLの取得に失敗しました。求人名: " + enterprise.Title
				log.Println(errMsg)
				return output, errors.New(errMsg)
			}

			detailPageURL := *detailPageURLMemory

			jobDetailURL := "https://agent-bank.com" + detailPageURL

			// 他媒体の求人IDを取得
			enterprise.ExternalID = regexpNumber.FindString(detailPageURL)

			// 媒体URLの取得
			enterprise.SecretMemo += "・媒体URL: " + jobDetailURL + "\n\n"
			fmt.Println("・媒体URL: " + jobDetailURL)

			detailPage := browser.MustPage(jobDetailURL)
			detailPage.WaitLoad()
			time.Sleep(4 * time.Second)

			/************ 4-3-4. 求人詳細  **************/

			divList, err := detailPage.Elements("div.list")
			if err != nil {
				errMsg = "求人詳細の取得に失敗しました。求人名: " + enterprise.Title
				log.Println(errMsg)
				return output, errors.New(errMsg)
			}

			for _, div := range divList {
				dt, err := div.Element("dt")
				if err != nil {
					errMsg = "dtの取得に失敗しました。求人名: " + enterprise.Title
					log.Println(errMsg)
					return output, errors.New(errMsg)
				}
				dtTxt := dt.MustText()
				dd, err := div.Element("dd")
				if err != nil {
					errMsg = "ddの取得に失敗しました。求人名: " + enterprise.Title
					log.Println(errMsg)
					return output, errors.New(errMsg)
				}
				ddTxt := dd.MustText()
				if ddTxt == "" || ddTxt == "未入力" {
					continue
				}

				switch dtTxt {
				/** 必須要件 */

				case "必須要件":
					// 応募資格に入れる
					enterprise.RequiredExperienceJobDetail += "■必須要件\n" + ddTxt + "\n\n"

				/** エージェント限定情報 */

				case "成功報酬":
					// 手数料（固定）に入れる
					if strings.Contains(ddTxt, "理論年収") {
						matches := regexpNumber.FindAllString(ddTxt, -1)
						// 抽出された部分を整数に変換
						if len(matches) > 0 {
							// 最初のマッチを取得して整数に変換
							rate, err := strconv.Atoi(matches[0])
							if err == nil {
								enterprise.CommissionRate = null.NewInt(int64(rate), true)
							}
						}
					}

					// 手数料（料率）に入れる
					if strings.Contains(ddTxt, "一律料金") {
						matches := regexpNumber.FindAllString(ddTxt, -1)
						// 抽出された部分を整数に変換
						if len(matches) > 0 {
							// 最初のマッチを取得して整数に変換
							commission, err := strconv.Atoi(matches[0])
							if err == nil {
								// 25万 -> 250000
								enterprise.Commission = null.NewInt(int64(commission*10000), true)
							}
						}
					}

					// 返金規定
					start := strings.Index(ddTxt, "【返金規定】")
					if start != -1 {
						// 返金規定の開始位置を調整
						startRefund := start + len("【返金規定】")
						enterprise.RefundPolicy = strings.TrimSpace(ddTxt[startRefund:])
					}

					// 手数料補足
					enterprise.CommissionDetail = strings.TrimSpace(ddTxt[:start])
				case "応募時に入力必須の項目":
					// 推薦時に必要な情報・書類に入れる
					enterprise.RequiredDocumentsDetail += "■応募時に入力必須の項目\n" + ddTxt + "\n\n"
				case "特記事項":
					// 推薦時に必要な情報・書類に入れる
					enterprise.RequiredDocumentsDetail += "・特記事項\n" + ddTxt + "\n\n" + "・媒体URL: " + jobDetailURL
				case "選考難易度":
					enterprise.TargetDetail += "■選考難易度\n" + ddTxt + "\n\n"
				case "直近の内定者情報":
					enterprise.TargetDetail += "■直近の内定者情報\n" + ddTxt + "\n\n"
				case "直近のNG情報（1面以降）":
					enterprise.TargetDetail += "■直近のNG情報（1面以降）\n" + ddTxt + "\n\n"
				case "実際の労働環境":
					enterprise.TargetDetail += "■実際の労働環境\n" + ddTxt + "\n\n"
				case "面接官情報&面接スタイル":
					enterprise.TargetDetail += "■面接官情報&面接スタイル\n" + ddTxt + "\n\n"
				case "ネガティブポイント":
					enterprise.TargetDetail += "■ネガティブポイント\n" + ddTxt + "\n\n"
				case "面接のポイント":
					enterprise.TargetDetail += "■面接のポイント\n" + ddTxt + "\n\n"
				case "応募承諾の取り方":
					enterprise.TargetDetail += "■応募承諾の取り方\n" + ddTxt + "\n\n"
				case "応募を後押しする情報":
					enterprise.TargetDetail += "■応募を後押しする情報\n" + ddTxt + "\n\n"
				case "通過者情報":
					enterprise.TargetDetail += "■通過者情報\n" + ddTxt + "\n\n"

				/** 仕事内容 */
				case "仕事内容":
					// 仕事内容と仕事内容（雇入直後）に入れる
					enterprise.WorkDetail = ddTxt
					enterprise.WorkDetailAfterHiring = ddTxt
				case "仕事の醍醐味":
					// 求人の魅力に入れる
					workCharmPoint := entity.JobInformationWorkCharmPoint{
						Title:    "仕事の醍醐味",
						Contents: ddTxt,
					}
					enterprise.WorkCharmPoints = append(enterprise.WorkCharmPoints, workCharmPoint)
				case "歓迎要件（活躍できる経験）":
					// 応募資格（経験・スキルなど）に入れる
					enterprise.RequiredExperienceJobDetail += "■歓迎要件（活躍できる経験）\n" + ddTxt

					/** 募集要項 */
				case "募集職種":
					// 募集職種に入れる
					occupations := entity.GetIntOccupationForAgentBank(ddTxt)
					for _, occupation := range occupations {
						enterprise.Occupations = append(enterprise.Occupations, entity.JobInformationOccupation{Occupation: occupation})
					}
				case "想定年収":
					// 年収下限と年収上限に入れる
					// 300 〜 450万円
					incomes := strings.Split(ddTxt, " 〜 ")
					if len(incomes) == 2 {
						overIncomeStr := strings.Replace(incomes[1], "万円", "", -1)
						overIncome, _ := strconv.Atoi(overIncomeStr)
						enterprise.OverIncome = null.NewInt(int64(overIncome), true)
					}
					underIncome, _ := strconv.Atoi(incomes[0])
					enterprise.UnderIncome = null.NewInt(int64(underIncome), true)
				case "予定募集人数":
					// 20人 or 20人~ or 20~30人のフォーマットで入力されている
					// 最低人数を募集人数として取得
					numberOfHiresStr := regexpNumber.FindString(ddTxt) // "32"
					numberOfHires, _ := strconv.Atoi(numberOfHiresStr)

					if numberOfHires > 0 {
						enterprise.NumberOfHires = null.NewInt(int64(numberOfHires), true)
					}
				case "雇用形態":
					// 雇用形態はautoscoutマスタを参照して入れる
					employmentStatusInt := entity.GetIntEmploymentStatusForAgentBank(ddTxt)
					if employmentStatusInt.Valid {
						enterprise.EmploymentStatuses = append(enterprise.EmploymentStatuses, entity.JobInformationEmploymentStatus{
							EmploymentStatus: employmentStatusInt,
						})

						// 雇用期間の定め
						if ddTxt == "正社員" || ddTxt == "無期雇用派遣" {
							enterprise.EmploymentPeriod = null.NewInt(1, true)
						} else {
							enterprise.EmploymentPeriod = null.NewInt(0, true)
						}
					}
				case "勤務時間":
					enterprise.WorkTime += ddTxt
				case "残業時間":
					// 平均残業時間に入れる
					enterprise.OvertimeAverage += ddTxt
				case "選考フロー":
					// 選考フローを読み込んで設定する
					// 筆記、webテスト or 面談 or 一次面接 or 二次面接 or 三次面接 or 四次面接 or. 五次面接 or 最終面接の数をカウントしてautoscoutで選考フローを設定する
					stepContents, _ := dd.Elements("div.step-content")
					steps := []string{}

					selectionFlowPattern := entity.JobInformationSelectionFlowPattern{
						// AgentID: uint,
						PublicStatus: null.NewInt(0, true),
						FlowTitle:    "通常選考",
						// FlowPattern:  選考数（len(steps)の値）に応じて決まる
						IsDeleted: false,
					}

					// 選考情報を選考フローにセット
					for _, stepContent := range stepContents {
						stepLabel := stepContent.MustText()

						if stepLabel == "筆記、webテスト" || stepLabel == "面談" || stepLabel == "一次面接" || stepLabel == "二次面接" || stepLabel == "三次面接" || stepLabel == "四次面接" || stepLabel == "五次面接" || stepLabel == "最終面接" {
							steps = append(steps, stepLabel)
						}
					}

					var (
						documentSelection = entity.JobInformationSelectionInformation{SelectionType: null.NewInt(1, true), IsQuestionnairy: true} // 書類選考
						firstSelection    = entity.JobInformationSelectionInformation{SelectionType: null.NewInt(2, true), IsQuestionnairy: true} // 一次選考
						secondSelection   = entity.JobInformationSelectionInformation{SelectionType: null.NewInt(3, true), IsQuestionnairy: true} // 二次選考
						thirdSelection    = entity.JobInformationSelectionInformation{SelectionType: null.NewInt(4, true), IsQuestionnairy: true} // 三次選考
						fourthSelection   = entity.JobInformationSelectionInformation{SelectionType: null.NewInt(5, true), IsQuestionnairy: true} // 四次選考
						fifthSelection    = entity.JobInformationSelectionInformation{SelectionType: null.NewInt(6, true), IsQuestionnairy: true} // 五次選考
						lastSelection     = entity.JobInformationSelectionInformation{SelectionType: null.NewInt(7, true), IsQuestionnairy: true} // 最終選考
						jobOffer          = entity.JobInformationSelectionInformation{SelectionType: null.NewInt(8, true), IsQuestionnairy: true} // 内定
						acceptJobOffer    = entity.JobInformationSelectionInformation{SelectionType: null.NewInt(9, true), IsQuestionnairy: true} // 内定承諾
					)

					// 書類選考をセット
					selectionFlowPattern.SelectionInformations = append(selectionFlowPattern.SelectionInformations, documentSelection)

					if len(steps) >= 2 {
						// 選考数2回以上
						selectionFlowPattern.SelectionInformations = append(selectionFlowPattern.SelectionInformations, firstSelection)
					}
					if len(steps) >= 3 {
						// 選考数3回以上
						selectionFlowPattern.SelectionInformations = append(selectionFlowPattern.SelectionInformations, secondSelection)
					}
					if len(steps) >= 4 {
						// 選考数4回以上
						selectionFlowPattern.SelectionInformations = append(selectionFlowPattern.SelectionInformations, thirdSelection)
					}
					if len(steps) >= 5 {
						// 選考数5回以上
						selectionFlowPattern.SelectionInformations = append(selectionFlowPattern.SelectionInformations, fourthSelection)
					}
					if len(steps) >= 6 {
						// 選考数5回以上
						selectionFlowPattern.SelectionInformations = append(selectionFlowPattern.SelectionInformations, fifthSelection)
					}

					// 最終選考、内定、内定承諾をセット
					selectionFlowPattern.SelectionInformations = append(selectionFlowPattern.SelectionInformations, lastSelection, jobOffer, acceptJobOffer)

					// 選考数に応じて選考フローパターンをセット
					flowPattern := len(stepContents) - 1
					selectionFlowPattern.FlowPattern = null.NewInt(int64(flowPattern), true)

					// 選考フローをセット
					enterprise.SelectionFlowPatterns = append(enterprise.SelectionFlowPatterns, selectionFlowPattern)
				case "選考詳細":
					// 選考フローに入れる
					enterprise.SelectionFlow += ddTxt
				case "勤務地":
					// コンテンツを取得
					locationContents, err := div.Elements("dd > div")
					if err == nil {
						// 重複チェックのmapを設定
						duplocateCheck := make(map[null.Int]bool)

						for _, locationContent := range locationContents {
							var (
								titleStr   string = ""
								addressStr string = ""
							)
							// フォーマットを整えて勤務地詳細に入れる
							titleContent, err := locationContent.Element("div.title")
							if err == nil {
								titleStr = titleContent.MustText()
							}
							addressContent, err := locationContent.Element("div.address-container")
							if err == nil {
								addressStr = addressContent.MustText()
							}

							if titleStr != "" && addressStr != "" {
								enterprise.WorkLocation += fmt.Sprintf("■%s\n%s\n\n", titleStr, addressStr)
							}

							prefectureContent, err := addressContent.Element("b")
							if err == nil {
								// 都道府県のテキストを勤務地マスタに変換
								prefectureStr := prefectureContent.MustText()
								prefectureInt := entity.GetIntPrefectureForCircus(prefectureStr)

								// 値がnullではなくかつ、重複がない場合はセット
								if prefectureInt.Valid && !duplocateCheck[prefectureInt] {
									// 重複チェックを更新
									duplocateCheck[prefectureInt] = true

									// 値をセット
									enterprise.Prefectures = append(enterprise.Prefectures, entity.JobInformationPrefecture{
										Prefecture: prefectureInt,
									})
								}
							}
						}
					}
				/** 仕事環境 */
				case "試用期間":
					if ddTxt == "あり" {
						enterprise.TrialPeriod = null.NewInt(0, true)
					} else if ddTxt == "なし" {
						enterprise.TrialPeriod = null.NewInt(1, true)
					}
				case "試用期間詳細":
					enterprise.TrialPeriodDetail += ddTxt
				case "給与・待遇":
					enterprise.Salary += ddTxt
				case "年間休日":
					enterprise.HolidayDetail += "■年間休日: " + ddTxt + "\n\n"
				case "休日・休暇":
					enterprise.HolidayDetail += "■休日・休暇\n" + ddTxt + "\n\n"

					// シフトという単語が含まれている場合は「休日タイプ」でシフトを選択
					var isShift = strings.Contains(ddTxt, "シフト")
					if isShift {
						enterprise.HolidayType = null.NewInt(3, true) // シフト制
					}
				case "福利厚生":
					enterprise.Insurance += ddTxt
				case "受動喫煙対策":
					if ddTxt == "禁煙" {
						// 「あり・屋内禁煙」で登録
						enterprise.PassiveSmoking = null.NewInt(0, true)
					} else if ddTxt == "喫煙スペースあり" {
						// 「あり・喫煙室設置・喫煙可能室設置」で登録
						enterprise.PassiveSmoking = null.NewInt(3, true)
					}
				case "受動喫煙対策（詳細）":
					enterprise.SecretMemo += "・受動喫煙対策\n" + ddTxt + "\n\n"

				default:
					continue
				}
			}

			/************ 4-3-5. 必要要件のテーブル  **************/

			if detailPage.MustHas("div.s-table__container.s-table__container_double") {
				tableContainers, err := detailPage.Elements("div.s-table__container.s-table__container_double")
				if err != nil {
					errMsg = "テーブルの取得に失敗しました。求人名: " + enterprise.Title
					log.Println(errMsg)
					return output, errors.New(errMsg)
				}

				for _, tableContainer := range tableContainers {
					dt, err := tableContainer.Element("dt")
					if err != nil {
						errMsg = "テーブルの取得に失敗しました。求人名: " + enterprise.Title
						log.Println(errMsg)
						return output, errors.New(errMsg)
					}
					dtTxt := dt.MustText()
					dd, err := tableContainer.Element("dd")
					if err != nil {
						errMsg = "テーブルの取得に失敗しました。求人名: " + enterprise.Title
						log.Println(errMsg)
						return output, errors.New(errMsg)
					}
					ddTxt := dd.MustText()
					if ddTxt == "" || ddTxt == "未入力" {
						continue
					}
					switch dtTxt {
					case "最終学歴":
						enterprise.FinalEducation = entity.GetIntFinalEducationForAgentBank(ddTxt)

					case "応募可能年齢":
						// 不問 or "18歳 〜 32歳" -> 18, 32
						if ddTxt != "不問" {
							ageStr := strings.Split(ddTxt, "歳 〜 ")

							if len(ageStr) == 2 {
								// 年齢の下限に入れる
								ageUnderStr := ageStr[0] // "18"
								ageUnder, _ := strconv.Atoi(ageUnderStr)
								enterprise.AgeUnder = null.NewInt(int64(ageUnder), true)

								// 年齢の上限に入れる
								ageOverStr := regexpNumber.FindString(ageStr[1]) // "32"
								ageOver, _ := strconv.Atoi(ageOverStr)
								if ageOver > ageUnder {
									enterprise.AgeOver = null.NewInt(int64(ageOver), true)
								}
							}
						}

					case "就業経験社数":
						enterprise.JobChange = entity.GetIntJobChangeForAgentBank(ddTxt)

					case "性別":
						// 男性 or 女性 or 不問
						if ddTxt == "男性" {
							enterprise.Gender = null.NewInt(0, true)
						} else if ddTxt == "女性" {
							enterprise.Gender = null.NewInt(1, true)
						} else if ddTxt == "不問" || ddTxt == "どちらでも" {
							enterprise.Gender = null.NewInt(99, true)
						}
					case "普通自動車免許":
						// 不要, 必須, 入社までに取得必要
						if ddTxt == "必須" {
							enterprise.DriverLicence = null.NewInt(0, true)
						} else if ddTxt == "入社時までに取得必須" {
							enterprise.DriverLicence = null.NewInt(1, true)
						} else if ddTxt == "不要" {
							enterprise.DriverLicence = null.NewInt(99, true)
						}
					case "未経験の可否":
						// ヘッダーのタグ部分で特徴は一括で入れるため、ここでは何もしない

						// if strings.Contains(ddTxt, "業界未経験可") {
						// 	enterprise.Features = append(enterprise.Features, entity.JobInformationFeature{
						// 		Feature: null.NewInt(0, true),
						// 	})
						// } else if strings.Contains(ddTxt, "職種未経験可") {
						// 	enterprise.Features = append(enterprise.Features, entity.JobInformationFeature{
						// 		Feature: null.NewInt(1, true),
						// 	})
						// } else if strings.Contains(ddTxt, "業界・職種未経験可") {
						// 	enterprise.Features = append(enterprise.Features, entity.JobInformationFeature{
						// 		Feature: null.NewInt(2, true),
						// 	})
						// }
					}
				}
			}

			/************ 4-3-6. 企業情報  **************/

			switchers, err := detailPage.Elements("div.switcher__content")
			if err != nil {
				errMsg = "企業情報の取得に失敗しました。求人名: " + enterprise.Title
				log.Println(errMsg)
				return output, errors.New(errMsg)
			}

			isMatchedWithSwitcher := false

			for _, switcher := range switchers {
				if strings.Contains(switcher.MustText(), "企業情報") {
					isMatchedWithSwitcher = true
					switcher.MustClick()
					time.Sleep(4 * time.Second)
					break
				}
			}

			if !isMatchedWithSwitcher {
				errMsg = "企業情報の取得に失敗しました。求人名: " + enterprise.Title
				log.Println(errMsg)
				return output, errors.New(errMsg)
			}

			companyDivList, err := detailPage.Elements("div.list")
			if err != nil {
				errMsg = "企業情報の取得に失敗しました。求人名: " + enterprise.Title
				log.Println(errMsg)
				return output, errors.New(errMsg)
			}

			// 会社概要を事業内容の下に入れるための変数
			var companyDetail = ""

			for _, companyDiv := range companyDivList {
				dt, err := companyDiv.Element("dt")
				if err != nil {
					errMsg = "企業情報の取得に失敗しました。求人名: " + enterprise.Title
					log.Println(errMsg)
					return output, errors.New(errMsg)
				}

				dtTxt := dt.MustText()

				dd, err := companyDiv.Element("dd")
				if err != nil {
					errMsg = "企業情報の取得に失敗しました。求人名: " + enterprise.Title
					log.Println(errMsg)
					return output, errors.New(errMsg)
				}

				ddTxt := dd.MustText()
				if ddTxt == "" || ddTxt == "未入力" {
					continue
				}

				switch dtTxt {
				case "会社名":
					enterprise.CompanyName = ddTxt
					enterprise.BillingAddressCompanyName = fmt.Sprintf("agent bank（%s）", ddTxt) // 請求先の会社名
				case "本社住所":
					enterprise.OfficeLocation = ddTxt
				case "業界":
					// 重複チェックのmapを設定
					duplocateCheck := make(map[null.Int]bool)

					// 複数の場合はカンマ区切りになっている
					industriesStrList := strings.Split(ddTxt, ",")
					for _, industryStr := range industriesStrList {
						industries := entity.GetIntIndustryForAgentBank(industryStr)

						for _, industryInt := range industries {
							if !duplocateCheck[industryInt] {
								// 重複チェックを更新
								duplocateCheck[industryInt] = true

								enterprise.Industries = append(enterprise.Industries, industryInt)
							}
						}
					}
				case "会社HP":
					enterprise.CorporateSiteURL = ddTxt
				case "従業員数":
					matches := regexpNumber.FindAllString(ddTxt, -1)
					// 抽出された部分を整数に変換
					if len(matches) > 0 {
						// 最初のマッチを取得して整数に変換
						employee, err := strconv.Atoi(matches[0])
						if err == nil {
							enterprise.EmployeeNumberSingle = null.NewInt(int64(employee), true)
						}
					}
				case "設立年月":
					enterprise.Establishment = ddTxt
				case "会社概要":
					companyDetail = ddTxt
				case "事業内容":
					enterprise.BusinessDetail = ddTxt
					if companyDetail != "" {
						enterprise.BusinessDetail += "\n\n" + companyDetail
					}
				default:
					continue
				}
			}

			/************ 4-3-7. 内定率、書類通過率、直近の応募数  **************/

			pointTagList, _ := detailPage.Elements("span.point-tag")

			if len(pointTagList) > 0 {
				var (
					offerRateTag                  = "内定率"
					documentPassingRateTag        = "書類通過率"
					recentNumberOfApplicationsTag = "直近の応募数"
				)

				for _, point := range pointTagList {
					pointStr := point.MustText()

					// 内定率
					if strings.Contains(pointStr, offerRateTag) {
						offerRate := entity.GetIntOfferRateForAgentBank(pointStr)
						enterprise.OfferRate = offerRate

						// 書類通過率
					} else if strings.Contains(pointStr, documentPassingRateTag) {
						documentPassingRate := entity.GetIntDocumentPassingRateForAgentBank(pointStr)
						enterprise.DocumentPassingRate = documentPassingRate

						// 直近の応募数
					} else if strings.Contains(pointStr, recentNumberOfApplicationsTag) {
						recentNumberOfApplications := entity.GetIntNumberOfRecentApplicationsForAgentBank(pointStr)
						enterprise.NumberOfRecentApplications = recentNumberOfApplications
					}
				}
			}

			/************ 4-3-8. 求人の特徴  **************/

			appealPointTagList, _ := detailPage.Elements("span.appeal-point-tag")

			if len(appealPointTagList) > 0 {
				var (
					noTransfer  = false // 転勤無しのチェック
					weekendsOff = false // 土日休みのチェック
				)

				for _, appealPoint := range appealPointTagList {
					featureStr := appealPoint.MustText()
					feature := entity.GetIntFeatureForAgentBank(featureStr)

					if feature.Valid {
						enterprise.Features = append(enterprise.Features, entity.JobInformationFeature{Feature: feature})
					}

					// featureがnullの可能性があるため、feature.Validのチェックの外に記述
					if !noTransfer && featureStr == "転勤なし" {
						noTransfer = true
						enterprise.Transfer = null.NewInt(1, true) // 転勤無し
					}

					if !weekendsOff && featureStr == "土日祝休み" {
						weekendsOff = true
						enterprise.HolidayType = null.NewInt(0, true) // 完全週休2日制（土日）
					}
				}
			}

			/************ 4-3-9. 企業リストの追加 + 求人詳細ページ閉じる  **************/

			enterpriseList = append(enterpriseList, &enterprise)

			detailPage.MustClose()
		}
	}

	/************ 5. Jsonファイルの作成（1ファイルだと大きくなりすぎるため、ファイル分割）  **************/

	var (
		// 分割数を定義
		numParts = 10

		// 1つの部分のサイズを計算
		partSize = (len(enterpriseList) + numParts - 1) / numParts
	)

	// if input.Param.IsJSON {
	type outputForJson struct {
		EnterpriseAndJobInformation []*entity.EnterpriseAndJobInformation `json:"enterprise_and_job_information_list"`
	}

	// 分割してファイルに書き込む
	for i := 0; i < numParts; i++ {
		// 各部分の開始インデックスと終了インデックスを計算
		startIndex := i * partSize
		endIndex := (i + 1) * partSize
		if endIndex > len(enterpriseList) {
			endIndex = len(enterpriseList)
		}
		if startIndex >= endIndex {
			break
		}

		// 部分リストを取得
		partList := enterpriseList[startIndex:endIndex]

		// ファイル名を生成
		fileName := "agentbank-enterprise-" + time.Now().Format("2006-01-02-150405") + "-" + strconv.Itoa(i+1) + ".json"

		// ファイルを作成
		file, err := os.Create(fileName)
		if err != nil {
			log.Println(err)
			return output, err
		}
		defer file.Close()

		// JSON ファイルに書き込み
		encoder := json.NewEncoder(file)
		if err := encoder.Encode(
			outputForJson{
				EnterpriseAndJobInformation: partList,
			},
		); err != nil {
			log.Println(err)
			return output, err
		}
	}

	// } else {

	// 企業情報を作成する
	// err = i.createEnterpriseList(enterpriseList, agentStaff)
	// if err != nil {
	// 	log.Println("err: ", err)
	// 	errMsg = "企業情報の作成に失敗しました"
	// 	log.Println(errMsg)
	// 	return output, errors.New(errMsg)
	// }

	/************ 6. 実行フラグの更新  **************/

	// isSuccessをtrueにする
	err = i.initialEnterpriseImporterRepository.UpdateIsSuccessTrue(input.Param.ID)
	if err != nil {
		log.Println("err: ", err)
		errMsg = "成功フラグ(isSuccess)の更新に失敗しました"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	log.Println("interactor処理終了. output: ", output)

	output.OK = true

	return output, nil
}

/**************************************************************************************/
// 求人企業のcsvファイルを読み込む (サーカス用)
//
func (i *InitialEnterpriseImporterInteractorImpl) createEnterpriseList(enterpriseList []*entity.EnterpriseAndJobInformation, agentStaff *entity.AgentStaff) error {
	var (
		err error
	)

	/*****被りを防ぐため既にDBに登録されている情報を取得****/
	// エージェントIDを指定して企業情報を取得
	enterpriseProfileListInDb, err := i.enterpriseProfileRepository.GetByAgentID(agentStaff.AgentID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// エージェントIDから請求先一覧を取得
	billingAddressListInDb, err := i.billingAddressRepository.GetByAgentID(agentStaff.AgentID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	/***************************/

	// 企業のかぶりチェック用
	type duplicatedEnterpriseChecker struct {
		ID                 uint
		CompanyName        string
		CircusEnterpriseID uint
	}

	var (
		duplicatedEnterpriseCheckerList     []duplicatedEnterpriseChecker
		isDuplicatedEnterpriseInCSV         bool
		isDuplicatedEnterpriseInDb          bool
		duplicatedBillingAddressCheckerList []entity.BillingAddress
		isDuplicatedBillingAddressInCSV     bool
		isDuplicatedBillingAddressInDb      bool
	)

	for index, enterprise := range enterpriseList {

		fmt.Println("csv企業情報", enterprise.AgentStaffID, enterprise.CompanyName)
		if enterprise.CompanyName == "" {
			fmt.Println("企業名が空欄のためスキップします。")
			continue
		}

		/***** 企業情報の登録 *****/
		isDuplicatedEnterpriseInCSV = false

		// CSV内で企業名とサーカス内の企業IDが一致する企業は保存しない
		// 一番初めは必ず保存
		if index > 0 {
			for _, duplicatedEnterpriseInCSV := range duplicatedEnterpriseCheckerList {

				if enterprise.CompanyName == duplicatedEnterpriseInCSV.CompanyName &&
					(enterprise.CircusEnterpriseID == duplicatedEnterpriseInCSV.CircusEnterpriseID || enterprise.CircusEnterpriseID == 0) {
					enterprise.EnterpriseID = duplicatedEnterpriseInCSV.ID
					fmt.Println("csv企業ID引き継ぎ:", duplicatedEnterpriseInCSV.ID)
					isDuplicatedEnterpriseInCSV = true
					break
				}

			}
		}

		if !isDuplicatedEnterpriseInCSV {

			// DBに既に入っている企業名と郵便番号が一致する企業は保存しない
			isDuplicatedEnterpriseInDb = false

			for _, enterpriseInDb := range enterpriseProfileListInDb {

				if enterprise.CompanyName == enterpriseInDb.CompanyName {

					enterprise.EnterpriseID = enterpriseInDb.ID
					isDuplicatedEnterpriseInDb = true
					break

				}
			}

			if !isDuplicatedEnterpriseInDb {

				enterpriseProfile := entity.NewEnterpriseProfile(
					enterprise.CompanyName,
					enterprise.AgentStaffID,
					enterprise.CorporateSiteURL,
					enterprise.Representative,
					enterprise.Establishment,
					enterprise.PostCode,
					enterprise.OfficeLocation,
					enterprise.EmployeeNumberSingle,
					enterprise.EmployeeNumberGroup,
					enterprise.Capital,
					enterprise.PublicOffering,
					enterprise.EarningsYear,
					enterprise.Earnings,
					enterprise.BusinessDetail,
				)

				err = i.enterpriseProfileRepository.Create(enterpriseProfile)
				if err != nil {
					fmt.Println(err)
					return err
				}

				// 空のレコードを作成
				referenceMaterial := entity.NewEnterpriseReferenceMaterial(
					enterpriseProfile.ID,
					"",
					"",
				)

				err = i.enterpriseReferenceMaterialRepository.Create(referenceMaterial)
				if err != nil {
					fmt.Println("企業登録時にエラーの求人の行:", enterprise.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", enterprise.EnterpriseID, "企業担当者ID:", enterprise.AgentStaffID)
					fmt.Println("請求先ID", enterprise.BillingAddressID, "請求先担当者ID", enterprise.AgentStaffIDForBillingAddress)
					fmt.Println(err)
					return err
				}

				for _, industry := range enterprise.Industries {
					industry := entity.NewEnterpriseIndustry(
						enterpriseProfile.ID,
						industry,
					)

					err = i.enterpriseIndustryRepository.Create(industry)
					if err != nil {
						fmt.Println(err)
						return err
					}
				}

				// 請求先で使用するために企業IDを保存する
				enterprise.EnterpriseID = enterpriseProfile.ID

			}

			// csvで被らなかった場合、IDを取得してから被りチェックリストにもアペンドする
			duplicatedEnterpriseCheckerList = append(
				duplicatedEnterpriseCheckerList,
				duplicatedEnterpriseChecker{
					ID:                 enterprise.EnterpriseID,
					CompanyName:        enterprise.CompanyName,
					CircusEnterpriseID: enterprise.CircusEnterpriseID,
				},
			)

		}

		/***** 請求先情報の登録 *****/
		// CSV内で請求先が一致する企業は保存しない
		isDuplicatedBillingAddressInCSV = false

		// 一番初めは必ず保存する
		if index > 0 {
			for _, duplicatedbBillingAddressInCSV := range duplicatedBillingAddressCheckerList {
				// すべての請求先担当者情報が一致した場合は飛ばす
				if enterprise.CompanyName == duplicatedbBillingAddressInCSV.CompanyName {
					enterprise.BillingAddressID = duplicatedbBillingAddressInCSV.ID
					isDuplicatedBillingAddressInCSV = true
					break
				}
			}
		}

		if !isDuplicatedBillingAddressInCSV {

			// DBに既に入っている企業名と郵便番号が一致する企業は保存しない
			isDuplicatedBillingAddressInDb = false

			for _, billingAddressInDb := range billingAddressListInDb {
				if enterprise.CompanyName == billingAddressInDb.CompanyName {
					enterprise.BillingAddressID = billingAddressInDb.ID
					isDuplicatedBillingAddressInDb = true
					break
				}
			}

			if !isDuplicatedBillingAddressInDb {

				billingAddress := entity.NewBillingAddress(
					enterprise.EnterpriseID,
					enterprise.AgentStaffID,
					enterprise.ContractPhase,
					enterprise.ContractDate,
					enterprise.PaymentPolicy,
					enterprise.CompanyName,
					enterprise.BillingAddressAddress,
					enterprise.HowToRecommend,
					enterprise.BillingAddressTitle,
				)

				err = i.billingAddressRepository.Create(billingAddress)
				if err != nil {
					fmt.Println(err)
					return err
				}

				// 求人登録時に使用するために請求先IDを保存する
				enterprise.BillingAddressID = billingAddress.ID

			}

			// csvで被らなかった場合、IDを取得してから被りチェックリストにもアペンドする
			duplicatedBillingAddressCheckerList = append(duplicatedBillingAddressCheckerList, entity.BillingAddress{
				ID:             enterprise.BillingAddressID,
				CompanyName:    enterprise.CompanyName,
				ContractDate:   enterprise.ContractDate,
				PaymentPolicy:  enterprise.PaymentPolicy,
				HowToRecommend: enterprise.HowToRecommend,
				AgentStaffID:   enterprise.AgentStaffIDForBillingAddress,
				HRStaffs:       enterprise.HRStaffs,
				RAStaffs:       enterprise.RAStaffs,
			})
		}

		/********* 求人情報 *********/
		jobInformation := entity.NewJobInformation(
			enterprise.BillingAddressID,
			enterprise.Title,
			enterprise.RecruitmentState,
			enterprise.ExpirationDate,
			enterprise.WorkDetail,
			enterprise.NumberOfHires,
			enterprise.WorkLocation,
			enterprise.Transfer,
			enterprise.TransferDetail,
			enterprise.UnderIncome,
			enterprise.OverIncome,
			enterprise.Salary,
			enterprise.Insurance,
			enterprise.WorkTime,
			enterprise.OvertimeAverage,
			enterprise.FixedOvertimePayment,
			enterprise.FixedOvertimeDetail,
			enterprise.TrialPeriod,
			enterprise.TrialPeriodDetail,
			enterprise.EmploymentPeriod,
			enterprise.EmploymentPeriodDetail,
			enterprise.HolidayType,
			enterprise.HolidayDetail,
			enterprise.PassiveSmoking,
			enterprise.SelectionFlow,
			enterprise.Gender,
			enterprise.Nationality,
			enterprise.FinalEducation,
			enterprise.SchoolLevel,
			enterprise.MedicalHistory,
			enterprise.AgeUnder,
			enterprise.AgeOver,
			enterprise.JobChange,
			enterprise.ShortResignation,
			enterprise.ShortResignationRemarks,
			enterprise.SocialExperienceYear,
			enterprise.SocialExperienceMonth,
			enterprise.Appearance,
			enterprise.Communication,
			enterprise.Thinking,
			enterprise.TargetDetail,
			enterprise.Commission,
			enterprise.CommissionRate,
			enterprise.CommissionDetail,
			enterprise.RefundPolicy,
			enterprise.RequiredExperienceJobDetail,
			enterprise.SecretMemo,
			enterprise.RequiredDocumentsDetail,
			enterprise.EmploymentInsurance,
			enterprise.AccidentInsurance,
			enterprise.HealthInsurance,
			enterprise.PensionInsurance,
			enterprise.RegisterPhase,
			enterprise.StudyCategory,
			enterprise.DriverLicence,
			enterprise.WordSkill,
			enterprise.ExcelSkill,
			enterprise.PowerPointSkill,
			false, // import時はautoscout求人のためfalse
			enterprise.WorkDetailAfterHiring,
			enterprise.WorkDetailScopeOfChange,
			enterprise.OfferRate,
			enterprise.DocumentPassingRate,
			enterprise.NumberOfRecentApplications,
			enterprise.IsGuaranteedInterview,
		)

		err = i.jobInformationRepository.Create(jobInformation)
		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, prefecture := range enterprise.Prefectures {
			p := entity.NewJobInformationPrefecture(
				jobInformation.ID,
				prefecture.Prefecture,
			)

			err = i.jobInfoPrefectureRepository.Create(p)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

		for _, employmentStatus := range enterprise.EmploymentStatuses {
			es := entity.NewJobInformationEmploymentStatus(
				jobInformation.ID,
				employmentStatus.EmploymentStatus,
			)

			err = i.jobInfoEmploymentStatusRepository.Create(es)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

		// 募集職種
		for _, occupation := range enterprise.Occupations {
			oc := entity.NewJobInformationOccupation(
				jobInformation.ID,
				occupation.Occupation,
			)

			err = i.jobInfoOccupationRepository.Create(oc)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

	}

	fmt.Println("企業情報の登録が完了しました。")
	fmt.Println("----------------------------------------")

	return err
}

// エージェントスタッフにメールを送信する
func (i *InitialEnterpriseImporterInteractorImpl) sendErrorMailToAgentStaffs(messageText string, agentStaffs []*entity.AgentStaff) error {
	var (
		err    error
		toList []entity.EmailUser
	)

	if os.Getenv("APP_ENV") != "prd" {
		return nil
	}

	if messageText == "" {
		return errors.New("メッセージが空です")
	} else if len(agentStaffs) == 0 {
		return errors.New("エージェントスタッフが空です")
	}

	// メール送信先リスト(管理者権限,利用可能の担当者)を作成 *通知の有無は考慮しない
	for _, staff := range agentStaffs {
		if staff.Authority != null.NewInt(int64(entity.AuthorityAdmin), true) ||
			staff.UsageStatus == null.NewInt(int64(entity.UsageStatusNotAvailable), true) {
			continue
		}
		toList = append(toList, entity.EmailUser{
			Name:  staff.StaffName,
			Email: staff.Email,
		})
	}

	// メール送信
	err = utility.SendMailList(
		i.sendgrid.APIKey,
		"RPA求人取り込み処理エラー通知",
		messageText,
		entity.EmailUser{
			Name:  "autoscout事務局",
			Email: "info@spaceai.jp",
		},
		toList,
	)

	if err != nil {
		log.Println("メール送信失敗")
		log.Println(err)
		return err
	} else {
		log.Println("メール送信成功")
	}

	return nil
}

type BatchUpdateEnterpriseForAgentBankInput struct {
	Now time.Time
}

type BatchUpdateEnterpriseForAgentBankOutput struct {
	OK bool
}

func (i *InitialEnterpriseImporterInteractorImpl) BatchUpdateEnterpriseForAgentBank(input BatchUpdateEnterpriseForAgentBankInput) (BatchUpdateEnterpriseForAgentBankOutput, error) {
	var (
		output           BatchUpdateEnterpriseForAgentBankOutput
		errMsg           string
		err              error
		selectedImporter entity.InitialEnterpriseImporter
		mailMessage      string
	)

	log.Println("interactor処理開始. input: ", input.Now)

	// 現在日時と実行時間が一致するインポート予約を取得
	importerList, err := i.initialEnterpriseImporterRepository.GetByStartDate(input.Now)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	log.Println("importerList", importerList)

	// タイムアウトを設定（2時間30分）
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Minute)
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
			errMsg = fmt.Sprintf("求人取得処理でpanicエラーが発生しました。Recover: %v\nStack:\n%s", rec, errorLine)
			log.Println(errMsg)
			err = errors.New(errMsg)

			// エラーメール送信
			now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
			if strings.Contains(fmt.Sprint(rec), "deadline") {
				mailMessage = fmt.Sprintf(
					"RPA処理の中でタイムアウトが発生しました。\n発生時刻: %s\n対象ID: %d\n 媒体: %s\n 開始時間: %s\n取得開始位置: %d\n 最大取得件数: %d\n",
					now, selectedImporter.ID, "サーカスエージェント", fmt.Sprint(selectedImporter.StartDate, "-", selectedImporter.StartHour.Int64, ":00"), selectedImporter.Offset.Int64, selectedImporter.MaxCount.Int64,
				)
			} else {
				mailMessage = fmt.Sprintf(
					"RPA処理の中でパニックエラーが発生しました。\n発生時刻: %s\n対象ID: %d\n 媒体: %s\n 開始時間: %s\n取得開始位置: %d\n 最大取得件数: %d\n",
					now, selectedImporter.ID, "サーカスエージェント", fmt.Sprint(selectedImporter.StartDate, "-", selectedImporter.StartHour.Int64, ":00"), selectedImporter.Offset.Int64, selectedImporter.MaxCount.Int64,
				)
			}
			agentStaffs := []*entity.AgentStaff{}
			i.sendErrorMailToAgentStaffs(mailMessage, agentStaffs)
		}
	}()

ImporterLoop:
	for _, importer := range importerList {
		selectedImporter = *importer

		switch importer.ServiceType {

		// サーカス
		case null.NewInt(entity.EnterpriseImporterServiceTypeCircus, true):
			log.Println("サーカスは実行しない")
			break ImporterLoop

			// エージェントバンク
		case null.NewInt(entity.EnterpriseImporterServiceTypeAgentBank, true):
			_, err = i.UpdateEnterpriseFromAgentBank(UpdateEnterpriseFromAgentBankInput{
				Param:   *importer,
				Context: ctx,
			})
			if err != nil {
				log.Println(err)
				mailMessage := fmt.Sprintf(
					"RPA処理の中でエラーが発生しました。\n対象ID: %d\n 媒体: エージェントバンク\n 開始時間: %s\n取得開始位置: %d\n 最大取得件数: %d\n エラー詳細: %s\n",
					importer.ID, fmt.Sprint(importer.StartDate, "-", importer.StartHour.Int64, ":00"), importer.Offset.Int64, importer.MaxCount.Int64, err.Error(),
				)
				// エラーメール送信のためにエージェントスタッフを取得
				agentStaff, err := i.agentStaffRepository.FindByID(importer.AgentStaffID)
				if err != nil {
					log.Println(err)
					return output, err
				}
				agentStaffs, err := i.agentStaffRepository.GetByAgentID(agentStaff.AgentID)
				if err != nil {
					log.Println(err)
					return output, err
				}
				i.sendErrorMailToAgentStaffs(mailMessage, agentStaffs)
				return output, err
			}
			break ImporterLoop
		}
	}
	log.Println("interactor処理終了. output: ", output)

	return output, nil
}

/*
エージェントバンクの求人を更新する
*/
type UpdateEnterpriseFromAgentBankInput struct {
	Param   entity.InitialEnterpriseImporter
	Context context.Context
}

type UpdateEnterpriseFromAgentBankOutput struct {
	OK bool `json:"ok"`
}

func (i *InitialEnterpriseImporterInteractorImpl) UpdateEnterpriseFromAgentBank(input UpdateEnterpriseFromAgentBankInput) (UpdateEnterpriseFromAgentBankOutput, error) {
	var (
		output  UpdateEnterpriseFromAgentBankOutput
		browser *rod.Browser
		page    *rod.Page
		errMsg  string
		err     error
		agentID uint = 1 // アンドイーズのエージェントID
	)

	/************ 1. エージェントバンクから取り込んだ求人取得を取得 **************/

	// 「job_information_external_idsテーブルにレコードを持ち、かつexternal_typeが0」の求人を取得
	// 求人ID, タイトル, 募集状況, 外部ID, 外部種別のみ取得する
	agentBankJobInformations, err := i.jobInformationRepository.GetByAgentIDAndExternalType(agentID, entity.JobInformatinoExternalTypeAgentBank)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	fmt.Println("length: ", len(agentBankJobInformations))

	/************ 2. ブラウザ起動 **************/

	if input.Context == nil {
		errMsg = "Contextが指定されていません"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	// ホストOSのchrome, chromium or microsoft edgeのパスを取得
	path, ok := launcher.LookPath()
	if !ok {
		errMsg = "chrome, chromium or microsoft edgeが見つかりません"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	// 取得したパスからchromeを起動
	l := launcher.New().
		Bin(path)

	// ローカル環境の場合は画面ありで起動
	if os.Getenv("APP_ENV") == "local" {
		l.Headless(false)
	}

	// ブラウザをexitし、UserDataDirを削除
	defer time.Sleep(10 * time.Second)
	defer l.Cleanup()

	// ブラウザを起動後、AMBIへ遷移
	mustLaunch := l.MustLaunch()
	browser = rod.New().
		ControlURL(mustLaunch).
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

	/************ 3. agentBankへログイン **************/

	// ログインページへ遷移
	page = browser.
		MustPage("https://agent-bank.com/service/login").
		MustWaitLoad()

	// ログインIDとパスワードを入力
	page.
		MustElement("input[name=email]").
		MustInput(input.Param.LoginID)

	// パスワードの複号
	input.Param.Password, err = decryption(input.Param.Password)
	if err != nil {
		log.Println(err)
		return output, err
	}

	page.
		MustElement("input[name=password]").
		MustInput(input.Param.Password)

	buttons := page.MustElements("button")
	isMatchedWithButton := false
	for _, button := range buttons {
		log.Println("button.MustText(): ", button.MustText())
		if button.MustText() == "ログイン" {
			button.MustClick()
			page.WaitLoad()
			isMatchedWithButton = true
		}
	}
	if !isMatchedWithButton {
		errMsg = "ログインボタンが見つかりません"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	page.WaitLoad()
	time.Sleep(10 * time.Second)

	// URLがログインページのままの場合はもう一度ログイン
	if strings.Contains(page.MustInfo().URL, "login") {
		page.
			MustElement("input[name=password]").
			MustInput(input.Param.Password)

		buttons := page.MustElements("button")
		isMatchedWithButton := false
		for _, button := range buttons {
			log.Println("button.MustText(): ", button.MustText())
			if button.MustText() == "ログイン" {
				button.MustClick()
				page.WaitLoad()
				isMatchedWithButton = true
			}
		}
		if !isMatchedWithButton {
			errMsg = "ログインボタンが見つかりません"
			log.Println(errMsg)
			return output, errors.New(errMsg)
		}
		page.WaitLoad()
		time.Sleep(10 * time.Second)
	}

	// URLがログインページのままの場合はログインに失敗していると判断
	if strings.Contains(page.MustInfo().URL, "login") {
		errMsg = "ログインに失敗しました。"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	/************ 4. ページ遷移のループ **************/

	for index, agentBankJob := range agentBankJobInformations {

		/************ 4-1. ページ遷移（求人詳細ページへ） **************/

		// 指定ページへ遷移
		jobInfoURL := fmt.Sprintf("https://agent-bank.com/service/job/%s", agentBankJob.ExternalID)
		loopLog := fmt.Sprintf("%v回目のループ: %v（%v）", index+1, agentBankJob.Title, agentBankJob.ExternalID)
		fmt.Println(loopLog)

		page.MustNavigate(jobInfoURL)
		page.WaitLoad()
		time.Sleep(4 * time.Second)

		var (
			recruitmentState null.Int // 募集状況: {0: 0pen, 1: Close}
		)

		if page.MustHas("div.s-table__container.s-table__container_double") {

			/************ 4-2. 募集停止ラベルの確認  **************/

			fmt.Println("求人ID: ", agentBankJob.ID)
			if page.MustHas("div.closedLabel > span.closedLabel__title") {
				closeLabelEL := page.MustElement("div.closedLabel > span.closedLabel__title")
				closeText := closeLabelEL.MustText()

				if closeText == "募集停止中" {
					recruitmentState = null.NewInt(1, true)
					fmt.Println("募集停止中...")
				}
			} else {
				recruitmentState = null.NewInt(0, true)
				fmt.Println("募集中!!")
			}

			/************ 4-3. 公開状況の更新  **************/

			// agentBank側で募集状況が変更されている場合は更新する
			if recruitmentState != agentBankJob.RecruitmentState {
				err = i.jobInformationRepository.UpdateRecruitmentState(agentBankJob.ID, recruitmentState)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}
	}

	/************ 5. 実行フラグの更新  **************/

	// isSuccessをtrueにする
	err = i.initialEnterpriseImporterRepository.UpdateIsSuccessTrue(input.Param.ID)
	if err != nil {
		log.Println("err: ", err)
		errMsg = "成功フラグ(isSuccess)の更新に失敗しました"
		log.Println(errMsg)
		return output, errors.New(errMsg)
	}

	output.OK = true

	return output, nil
}
