package batch

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/slack-go/slack"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/infrastructure/database"
	"github.com/spaceaiinc/autoscout-server/infrastructure/di"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type Batch struct {
	cfg       config.Config
	Engine    *echo.Echo
	scheduler *gocron.Scheduler
	db        *database.DB
	firebase  usecase.Firebase
}

func NewBatch(
	cfg config.Config,
	db *database.DB,
	firebase usecase.Firebase,
) *Batch {
	return &Batch{
		cfg:       cfg,
		Engine:    echo.New(),
		scheduler: gocron.NewScheduler(time.FixedZone("Asia/Tokyo", 9*60*60)),
		db:        db,
		firebase:  firebase,
	}
}

func (b *Batch) SetUp() *Batch {
	// サーバーの設定
	b.Engine.HidePort = true
	b.Engine.HideBanner = true

	// パニック時のエラー処理
	b.Engine.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          0,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			stackStr := string(stack)
			c.Logger().Error(stackStr)

			// interactor or repositoryを含む文字列を取得
			var errorLine string
			for _, line := range strings.Split(stackStr, "\n") {
				if strings.Contains(line, "interactor/") || strings.Contains(line, "repository/") {
					errorLine += line + "\n"
				}
			}
			log.Println("errorLine:", errorLine)

			// slack通知
			if b.cfg.App.Env == "prd" || b.cfg.App.Env == "dev" {
				b.notifyError(fmt.Errorf("PANIC ERROR.\nerror line:\n%s", errorLine))
			}
			return nil
		},
	}))

	/****************************************************************************************/
	// バッチ処理
	// バッチの実行時間設定 *1時間毎に実行 *初回は現在時刻から次の00分00秒に実行
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	firstStartTime := now.Add(time.Duration(60-now.Minute())*time.Minute + time.Duration(60-now.Second())*time.Second)
	log.Println("現在時刻(JST):", now, "初回実行時間(JST):", firstStartTime)
	log.Println("APP_SERVICE:", b.cfg.App.Service, "AGENT_ROBOT_ID:", b.cfg.RPA.AgentRobotID)

	/*
		未読の通知
		slack:https://app.slack.com/client/T01EY1159HP/C04MGLZSH8Q/thread/C04MGLZSH8Q-1687426686.644339)
		未読チャットの通知 対象:エージェント向けチャット。@付きのメッセージ
		未処理のタスク 対象:未読or期限切れ
		08:00に通知
		※本番環境のみ
	*/
	if b.cfg.App.Env == "prd" {
		// batchNotifyUnwatchedJob, err := b.scheduler.
		// 	Every(1).
		// 	Monday().
		// 	Tuesday().
		// 	Wednesday().
		// 	Thursday().
		// 	Friday().
		// 	At("08:00").
		// 	Do(
		// 		func() {
		// 			now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
		// 			log.Println("batchNotifyUnwatched開始 現在時刻(JST):", now)
		// 			err := b.batchNotifyUnwatched(now)
		// 			// slack通知
		// 			if (b.cfg.App.Env == "prd" || b.cfg.App.Env == "dev") && err != nil {
		// 				b.notifyError(err)
		// 			}
		// 			for _, job := range b.scheduler.Jobs() {
		// 				log.Println("job tag:", job.Tags(), "\nlast run time:", job.LastRun(), "\nnext run time:", job.NextRun())
		// 			}
		// 			log.Println("batchNotifyUnwatched処理終了")
		// 		},
		// 	)
		// if err != nil {
		// 	log.Println("err:", err)
		// 	panic(err)
		// }

		// batchNotifyUnwatchedJob.Tag("batchNotifyUnwatched")

		/*
			求人一括インポート処理
			1時間おきに実行
		*/
		// initialEnterpriseImporterJob, err := b.scheduler.
		// 	Every(1).
		// 	Hour().
		// 	StartAt(firstStartTime).
		// 	Do(
		// 		func() {
		// 			now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
		// 			log.Println("BatchInitialEnterpriseImporter開始 現在時刻(JST):", now)
		// 			err := b.batchInitialEnterpriseImporter(now)
		// 			// slack通知
		// 			if (b.cfg.App.Env == "prd" || b.cfg.App.Env == "dev") && err != nil {
		// 				b.notifyError(err)
		// 			}
		// 			log.Println("BatchInitialEnterpriseImporter処理終了")
		// 		},
		// 	)
		// if err != nil {
		// 	log.Println("err:", err)
		// 	panic(err)
		// }

		// initialEnterpriseImporterJob.Tag("batchInitialEnterpriseImporter")

		// Pub/Subで取得したエントリーユーザーを媒体ごとに取り込む処理（15分に一度実行する）
		// AMBI&マイナビ転職スカウト
		batchEntryUser, err := b.scheduler.
			Every(1).
			Day().
			At("00:00;00:15;00:30;00:45;01:00;01:15;01:30;01:45;02:00;02:15;02:30;02:45;03:00;03:15;03:30;03:45;04:00;04:15;04:30;04:45;05:00;05:15;05:30;05:45;06:00;06:15;06:30;06:45;07:00;07:15;07:30;07:45;08:00;08:15;08:30;08:45;09:00;09:15;09:30;09:45;10:00;10:15;10:30;10:45;11:00;11:15;11:30;11:45;12:00;12:15;12:30;12:45;13:00;13:15;13:30;13:45;14:00;14:15;14:30;14:45;15:00;15:15;15:30;15:45;16:00;16:15;16:30;16:45;17:00;17:15;17:30;17:45;18:00;18:15;18:30;18:45;19:00;19:15;19:30;19:45;20:00;20:15;20:30;20:45;21:00;21:15;21:30;21:45;22:00;22:15;22:30;22:45;23:00;23:15;23:30;23:45").
			Do(
				func() {
					now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
					log.Println("batchEntryUser開始 現在時刻(JST):", now)
					err := b.batchEntry(now)
					// slack通知
					if (b.cfg.App.Env == "prd" || b.cfg.App.Env == "dev") && err != nil {
						b.notifyError(err)
					}
					for _, job := range b.scheduler.Jobs() {
						log.Println("job tag:", job.Tags(), "\nlast run time:", job.LastRun(), "\nnext run time:", job.NextRun())
					}
					log.Println("batchEntryUser処理終了")
				},
			)
		if err != nil {
			log.Println("err:", err)
			panic(err)
		}

		batchEntryUser.Tag("batchEntryUser")

		/*
			エージェントバンク求人の更新処理
			毎日2時に実行
		*/
		// initialUpdateEnterpriseForAgentBankJob, err := b.scheduler.
		// 	Every(1).
		// 	Day().
		// 	At("00:45").
		// 	Do(
		// 		func() {
		// 			now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
		// 			log.Println("BatchUpdateEnterpriseForAgentBank開始 現在時刻(JST):", now)
		// 			err := b.batchUpdateEnterpriseForAgentBank(now)
		// 			if err != nil {
		// 				log.Println("err:", err)
		// 			}
		// 			log.Println("BatchUpdateEnterpriseForAgentBank処理終了")
		// 		},
		// 	)
		// if err != nil {
		// 	log.Println("err:", err)
		// 	panic(err)
		// }

		// initialUpdateEnterpriseForAgentBankJob.Tag("batchInitialEnterpriseImporter")

	} else if b.cfg.App.Env == "local" {
		// // Pub/Subで取得したエントリーユーザーを媒体ごとに取り込む処理（15分に一度実行する）
		// // AMBI&マイナビ転職スカウト
		// batchEntryUser, err := b.scheduler.
		// 	Every(1).
		// 	Day().
		// 	At("00:00;00:15;00:30;00:45;01:00;01:15;01:30;01:45;02:00;02:15;02:30;02:45;03:00;03:15;03:30;03:45;04:00;04:15;04:30;04:45;05:00;05:15;05:30;05:45;06:00;06:15;06:30;06:45;07:00;07:15;07:30;07:45;08:00;08:15;08:30;08:45;09:00;09:15;09:30;09:45;10:00;10:15;10:30;10:45;11:00;11:15;11:30;11:45;12:00;12:15;12:30;12:45;13:00;13:15;13:30;13:45;14:00;14:15;14:30;14:45;15:00;15:15;15:30;15:45;16:00;16:15;16:30;16:45;17:00;17:15;17:30;17:45;18:00;18:15;18:30;18:45;19:00;19:15;19:30;19:45;20:00;20:15;20:30;20:45;21:00;21:15;21:30;21:45;22:00;22:15;22:30;22:45;23:00;23:15;23:30;23:45").
		// 	Do(
		// 		func() {
		// 			now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
		// 			log.Println("batchEntryUser開始 現在時刻(JST):", now)
		// 			err := b.batchEntry(now)
		// 			// slack通知
		// 			if (b.cfg.App.Env == "prd" || b.cfg.App.Env == "dev") && err != nil {
		// 				b.notifyError(err)
		// 			}
		// 			for _, job := range b.scheduler.Jobs() {
		// 				log.Println("job tag:", job.Tags(), "\nlast run time:", job.LastRun(), "\nnext run time:", job.NextRun())
		// 			}
		// 			log.Println("batchEntryUser処理終了")
		// 		},
		// 	)
		// if err != nil {
		// 	log.Println("err:", err)
		// 	panic(err)
		// }

		// batchEntryUser.Tag("batchEntryUser")

		/*
			スカウト処理
			1時間おきに実行
		*/
		batchScoutJob, err := b.scheduler.
			Every(1).
			Day().
			At("12:01;18:01").
			Do(
				func() {
					now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
					log.Println("BatchScout開始 現在時刻(JST):", now)
					err := b.batchScout(now)
					log.Println(err)
					for _, job := range b.scheduler.Jobs() {
						log.Println("job tag:", job.Tags(), "\nlast run time:", job.LastRun(), "\nnext run time:", job.NextRun())
					}
					log.Println("BatchScout処理終了")
				},
			)
		if err != nil {
			log.Println("err:", err)
			panic(err)
		}

		// 関数名をタグ付け
		batchScoutJob.Tag("batchScout")

		now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
		log.Println("BatchScout開始 現在時刻(JST):", now)
		err = b.batchScout(now)
		// slack通知
		if err != nil && b.cfg.App.BatchType == "scout" {
			b.notifyError(err)
		}
		log.Println("BatchScout処理終了")

		// /*
		// 	求人一括インポート処理 テスト用
		// 	立ち上げ時と15分おきに実行
		// */
		// initialEnterpriseImporterTestJob, err := b.scheduler.
		// 	Every(15).
		// 	Hour().
		// 	StartImmediately().
		// 	Do(
		// 		func() {
		// 			now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
		// 			log.Println("BatchInitialEnterpriseImporter開始(テスト用) 現在時刻(JST):", now)
		// 			err := b.batchInitialEnterpriseImporter(now)
		// 			if err != nil {
		// 				log.Println("err:", err)
		// 			}
		// 			log.Println("BatchInitialEnterpriseImporter処理(テスト用)終了")
		// 		},
		// 	)
		// if err != nil {
		// 	log.Println("err:", err)
		// 	panic(err)
		// }

		// initialEnterpriseImporterTestJob.Tag("batchInitialEnterpriseImporterTest")

		/*
			エージェントバンク求人の更新処理 テスト用
			毎日3時に実行
		*/
		// initialUpdateEnterpriseForAgentBankTestJob, err := b.scheduler.
		// 	StartImmediately().
		// 	Do(
		// 		func() {
		// 			now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
		// 			log.Println("BatchUpdateEnterpriseForAgentBank開始(テスト用) 現在時刻(JST):", now)
		// 			err := b.batchUpdateEnterpriseForAgentBank(now)
		// 			if err != nil {
		// 				log.Println("err:", err)
		// 			}
		// 			log.Println("BatchUpdateEnterpriseForAgentBank処理(テスト用)終了")
		// 		},
		// 	)
		// if err != nil {
		// 	log.Println("err:", err)
		// 	panic(err)
		// }

		// initialUpdateEnterpriseForAgentBankTestJob.Tag("batchInitialEnterpriseImporterTest")

	} else if b.cfg.App.Env == "dev" {
		// Pub/Subで取得したエントリーユーザーを媒体ごとに取り込む処理（15分に一度実行する）
		// AMBI&マイナビ転職スカウト
		batchEntryUser, err := b.scheduler.
			Every(1).
			Day().
			At("00:00;00:15;00:30;00:45;01:00;01:15;01:30;01:45;02:00;02:15;02:30;02:45;03:00;03:15;03:30;03:45;04:00;04:15;04:30;04:45;05:00;05:15;05:30;05:45;06:00;06:15;06:30;06:45;07:00;07:15;07:30;07:45;08:00;08:15;08:30;08:45;09:00;09:15;09:30;09:45;10:00;10:15;10:30;10:45;11:00;11:15;11:30;11:45;12:00;12:15;12:30;12:45;13:00;13:15;13:30;13:45;14:00;14:15;14:30;14:45;15:00;15:15;15:30;15:45;16:00;16:15;16:30;16:45;17:00;17:15;17:30;17:45;18:00;18:15;18:30;18:45;19:00;19:15;19:30;19:45;20:00;20:15;20:30;20:45;21:00;21:15;21:30;21:45;22:00;22:15;22:30;22:45;23:00;23:15;23:30;23:45").
			Do(
				func() {
					now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
					log.Println("batchEntryUser開始 現在時刻(JST):", now)
					err := b.batchEntry(now)
					// slack通知
					if (b.cfg.App.Env == "prd" || b.cfg.App.Env == "dev") && err != nil {
						b.notifyError(err)
					}
					for _, job := range b.scheduler.Jobs() {
						log.Println("job tag:", job.Tags(), "\nlast run time:", job.LastRun(), "\nnext run time:", job.NextRun())
					}
					log.Println("batchEntryUser処理終了")
				},
			)
		if err != nil {
			log.Println("err:", err)
			panic(err)
		}

		batchEntryUser.Tag("batchEntryUser")
	}

	// 非同期で実行。実行中の処理をブロックせずに処理を実行する
	b.scheduler.StartAsync()

	return b
}

func (b *Batch) Start() {
	log.Println("Batch Server Start. Port:", b.cfg.App.Port)
	b.Engine.Start(fmt.Sprintf(":%d", b.cfg.App.Port))
}

// Slack通知
func (b *Batch) notifyError(err error) {
	var (
		code, message = entity.ErrorInfo(err)
		statusCode    = strconv.Itoa(code)
		// path          = c.Path()
		// method        = c.Request().Method
		errText = err.Error()
	)

	te := "*autoscout開発環境 Error*\n" +
		">>>status: " + message + "（" + statusCode + "）" + "\n" +
		"service: " + "batch" + "\n" +
		"agent_robot_id: " + fmt.Sprint(b.cfg.RPA.AgentRobotID) + "\n" +
		"error: `" + errText + "` \n"

	// アクセストークンを使用してクライアントを生成する
	// https://api.slack.com/apps からトークン取得
	// 参考: https://risaki-masa.com/how-to-get-api-token-in-slack/
	tkn := b.cfg.Slack.AccessToken
	chanelID := b.cfg.Slack.ChanelID
	s := slack.New(tkn)

	// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
	_, _, err = s.PostMessage(chanelID, slack.MsgOptionBlocks(
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: te,
			},
		},
	))
	if err != nil {
		log.Println(err)
	}

	log.Println("slack通知完了")
}

// 15分おきにスカウト媒体から新規求職者を取得
func (b *Batch) batchEntry(now time.Time) error {
	tx, err := b.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	h := di.InitializeScoutServiceHandler(b.firebase, b.db, b.cfg.Sendgrid, b.cfg.OneSignal, b.cfg.App, b.cfg.GoogleAPI, b.cfg.Slack)
	_, err = h.BatchEntry(now)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// 各エージェントのスカウト媒体から新規求職者を取得
// func (b *Batch) batchEntry(now time.Time) error {
// 	tx, err := b.db.Begin()
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	h := di.InitializeScoutServiceHandler(b.firebase, b.db, b.cfg.Sendgrid, b.cfg.OneSignal, b.cfg.App, b.cfg.GoogleAPI)
// 	_, err = h.BatchEntry(now, uint(b.cfg.RPA.AgentRobotID))
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	tx.Commit()

// 	return nil
// }

// 各エージェントのスカウト媒体からスカウトを送信
func (b *Batch) batchScout(now time.Time) error {
	h := di.InitializeScoutServiceHandler(b.firebase, b.db, b.cfg.Sendgrid, b.cfg.OneSignal, b.cfg.App, b.cfg.GoogleAPI, b.cfg.Slack)
	_, err := h.BatchScout(now, uint(b.cfg.RPA.AgentRobotID))
	if err != nil {
		return err
	}

	return nil
}

// // 各エージェントの求人管理媒体から求人企業を一括インポート
// func (b *Batch) batchInitialEnterpriseImporter(now time.Time) error {
// 	tx, err := b.db.Begin()
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	h := di.InitializeInitialEnterpriseImporterHandler(b.firebase, tx, b.cfg.Sendgrid, b.cfg.OneSignal)
// 	_, err = h.BatchInitialEnterpriseImporter(now)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	tx.Commit()

// 	return nil
// }

// エージェントバンク求人の募集状況を一括更新
// func (b *Batch) batchUpdateEnterpriseForAgentBank(now time.Time) error {
// 	tx, err := b.db.Begin()
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	h := di.InitializeInitialEnterpriseImporterHandler(b.firebase, tx, b.cfg.Sendgrid, b.cfg.OneSignal)
// 	_, err = h.BatchUpdateEnterpriseForAgentBank(now)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	tx.Commit()

// 	return nil
// }

// 未読or未処理を通知(エージェントとのチャット, タスク)
// func (b *Batch) batchNotifyUnwatched(now time.Time) error {

// 	h := di.InitializeTaskHandler(b.firebase, b.db, b.cfg.Sendgrid, b.cfg.OneSignal)
// 	_, err := h.BatchNotifyUnwatched(now)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
