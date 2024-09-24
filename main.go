package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/infrastructure/batch"
	"github.com/spaceaiinc/autoscout-server/infrastructure/database"
	"github.com/spaceaiinc/autoscout-server/infrastructure/driver"
	"github.com/spaceaiinc/autoscout-server/infrastructure/router"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

func init() {
	time.Local = utility.Tokyo
	if os.Getenv("APP_ENV") == "local" {
		if err := godotenv.Load(); err != nil {
			panic("Failed to load .env file")
		}
	}
}

// リリース時に更新
var version string

func main() {
	fmt.Printf("version: %s", version)

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	switch cfg.App.Service {
	case "api":
		// 一旦 apiコンテナを立ち上げる時にマイグレーションする
		db := database.NewDB(cfg.DB, true)
		err := db.MigrateUp(".migrations")
		if err != nil {
			fmt.Println(err)
		}
		// cache := driver.NewRedisCacheImpl(cfg.Redis)
		if cfg.App.Env == "local" {
			firebase := driver.NewFirebaseImpl(cfg.Firebase)
			getTestUserToken(firebase, uuid.New().String())
		}
		r := router.NewRouter(cfg)

		// エラーハンドラー（dev or prdのみSlack通知）
		if cfg.App.Env != "local" {
			r.Engine.HTTPErrorHandler = customHTTPErrorHandler
		}

		// ルーティング
		r.SetUp().Start()

	case "batch":
		db := database.NewDB(cfg.DB, false)
		err := db.MigrateUp(".migrations")
		if err != nil {
			fmt.Println(err)
		}
		firebase := driver.NewFirebaseImpl(cfg.Firebase)
		// バッチ処理開始
		b := batch.NewBatch(cfg, db, firebase)
		b.SetUp().Start()
	}
}

func getTestUserToken(fb usecase.Firebase, uuid string) {
	customToken, _ := fb.GetCustomToken(uuid)
	idToken, err := fb.GetIDToken(customToken)
	if err != nil {
		panic(err)
	}
	fmt.Println(idToken)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	var (
		cfg, _        = config.New()
		code, message = entity.ErrorInfo(err)
		statusCode    = strconv.Itoa(code)
		path          = c.Path()
		method        = c.Request().Method
		errText       = err.Error()
	)

	fmt.Println(err)

	if path != "/*" {
		te := "*autoscout開発環境 Error*\n" +
			">>>status: " + message + "（" + statusCode + "）" + "\n" +
			"method: " + method + "\n" +
			"uri: " + path + "\n" +
			"error: `" + errText + "` \n"

		// アクセストークンを使用してクライアントを生成する
		// https://api.slack.com/apps からトークン取得
		// 参考: https://risaki-masa.com/how-to-get-api-token-in-slack/
		tkn := cfg.Slack.AccessToken
		chanelID := cfg.Slack.ChanelID
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
			fmt.Println(err)
		}

		c.Logger().Error(err)
	}
}
