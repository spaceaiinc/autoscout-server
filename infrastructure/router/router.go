package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/infrastructure/database"
	"github.com/spaceaiinc/autoscout-server/infrastructure/di"
	"github.com/spaceaiinc/autoscout-server/infrastructure/driver"
	"github.com/spaceaiinc/autoscout-server/infrastructure/router/routes"
)

type Router struct {
	cfg    config.Config
	Engine *echo.Echo
}

func NewRouter(cfg config.Config) *Router {
	return &Router{
		cfg:    cfg,
		Engine: echo.New(),
	}
}

// レスポンスのJSONをログで出力する関数
func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("Request Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
}

func (r *Router) SetUp() *Router {
	var (
		db       = database.NewDB(r.cfg.DB, true)
		firebase = driver.NewFirebaseImpl(r.cfg.Firebase)

		// AdminAPI用で使っていたためコメントアウト
		// basicAuth = utility.NewBasicAuth(r.cfg)
	)

	r.Engine.HidePort = true
	r.Engine.HideBanner = true

	// パニック時のエラー処理
	// r.Engine.Use(middleware.Recover())
	r.Engine.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          0,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			stackStr := string(stack)
			c.Logger().Error(stackStr)
			c.JSON(http.StatusInternalServerError, "PANIC ERROR")

			// interactor or repositoryを含む文字列を取得
			var errorLine string
			for _, line := range strings.Split(stackStr, "\n") {
				if strings.Contains(line, "interactor/") || strings.Contains(line, "repository/") {
					errorLine = line
					log.Println("error line:", errorLine)
					break
				}
			}

			return fmt.Errorf("PANIC ERROR.\nerror line:\n%s", errorLine)
		},
	}))

	fmt.Println("------------")
	fmt.Println(r.cfg.App.Env)
	fmt.Println("------------")

	// NOTE: 環境変数を使用するのはLPのプレビュー環境ドメインが可変のため（ワイルドカードの指定不可）
	var allowOrigins = r.cfg.App.CorsDomains
	fmt.Println("allowOrigins: ", len(allowOrigins), allowOrigins)

	r.Engine.Use(middleware.CORSWithConfig((middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowHeaders: []string{
			echo.HeaderAuthorization,
			echo.HeaderAccessControlAllowHeaders,
			echo.HeaderContentType,
			echo.HeaderOrigin,
			echo.HeaderAccessControlAllowOrigin,
			"FirebaseAuthorization",
			"LineAuthorization",
		},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	})))

	r.Engine.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "healthz") {
				return true
			} else {
				return false
			}
		},
	}))

	// middlewareセット
	// r.Engine.Use(middleware.BodyDump(bodyDumpHandler))

	// 連続リクエストの制限 *Gmail Webhookは1/15minのリクエストを許可(https://echo.labstack.com/docs/middleware/rate-limiter)
	// var (
	// 	// 間隔
	// 	duration = 5 * time.Minute
	// 	// 最大リクエスト数
	// 	burst = 1
	// 	// burst/interval秒間に許可されるリクエスト数
	// 	limit = rate.NewLimiter(rate.Every(duration), burst).Limit()
	// )
	// rateLimiterConfig := middleware.RateLimiterConfig{
	// 	Skipper: middleware.DefaultSkipper,
	// 	Store: middleware.NewRateLimiterMemoryStoreWithConfig(
	// 		middleware.RateLimiterMemoryStoreConfig{Rate: limit, Burst: 1, ExpiresIn: 1 * time.Minute},
	// 	),
	// 	IdentifierExtractor: func(ctx echo.Context) (string, error) {
	// 		id := ctx.RealIP()
	// 		return id, nil
	// 	},
	// 	ErrorHandler: func(context echo.Context, err error) error {
	// 		return context.JSON(http.StatusForbidden, nil)
	// 	},
	// 	DenyHandler: func(context echo.Context, identifier string, err error) error {
	// 		// slackに通知を飛ばさないようにする
	// 		return context.JSON(http.StatusAccepted, nil)
	// 		// return context.JSON(http.StatusTooManyRequests, nil)
	// 	},
	// }

	api := r.Engine.Group("")
	{
		api.GET("/healthz", func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})

		api.GET("/*", func(c echo.Context) error {
			return c.NoContent(http.StatusNotFound)
		})

		api.POST("/*", func(c echo.Context) error {
			return c.NoContent(http.StatusNotFound)
		})

		api.PUT("/*", func(c echo.Context) error {
			return c.NoContent(http.StatusNotFound)
		})
	}

	/****************************************************************************************/
	/// No Auth API
	//
	noAuthAPI := api.Group("api")
	{
		noAuthAPI.GET("/healthz", func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})

		// ユーザーのログイン処理（ログイン状況やログイン時間の更新）
		noAuthAPI.PUT("/signin", routes.SignIn(db, firebase, r.cfg.Sendgrid))

		// ユーザーのログアウト処理（ログイン状況やログアウト時間の更新）
		noAuthAPI.PUT("/signout", routes.SignOut(db, firebase, r.cfg.Sendgrid))

		// ユーザーのログイン情報を取得
		noAuthAPI.GET("/signin/user", routes.GetSignInUser(db, firebase, r.cfg.Sendgrid))

		// LINEのメッセージ受信
		noAuthAPI.POST("/line", routes.LineWebHook(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// google認証のURLを発行
		noAuthAPI.GET("/auth_code_url", routes.GetGoogleAuthCodeURL(db, firebase, r.cfg.Sendgrid, r.cfg.GoogleAPI))

		// google認証のURLを発行
		noAuthAPI.PUT("/google_certification", routes.UpdateGoogleOauthToken(db, firebase, r.cfg.Sendgrid, r.cfg.GoogleAPI))
	}
	/****************************************************************************************/
	/// AgentStaff API
	//
	guestAPI := noAuthAPI.Group("/guest")
	{
		// ゲストユーザーのログイン
		{
			// ゲスト企業
			guestAPI.PUT("/signin/for/enterprise/:job_information_uuid", routes.SignInForGuestEnterprise(db, firebase, r.cfg.Sendgrid))

			// タスクグループのuuidでゲスト企業のログイン
			guestAPI.PUT("/signin/for/enterprise/task_group/:task_group_uuid", routes.SignInForGuestEnterpriseByTaskGroupUUID(db, firebase, r.cfg.Sendgrid))

			// ゲスト求職者
			guestAPI.PUT("/signin/for/job_seeker/:job_seeker_uuid", routes.SignInForGuestJobSeeker(db, firebase, r.cfg.Sendgrid))

			// マイページログイン（LPからログイントークンを使ってログイン）
			guestAPI.PUT("/signin/from_lp/for/job_seeker/:job_seeker_uuid", routes.SignInForGuestJobSeekerFromLP(db, firebase, r.cfg.Sendgrid))

			// ゲスト送客求職者
			guestAPI.PUT("/signin/for/sending_job_seeker/:sending_job_seeker_uuid", routes.SignInForGuestSendingJobSeeker(db, firebase, r.cfg.Sendgrid))
		}
	}
	/****************************************************************************************/
	/// AgentStaff API
	//
	userAPI := noAuthAPI.Group("/agent_staff")
	{
		/************************************** POSTメソッド **************************************/

		/************************************** PUTメソッド **************************************/

		// 担当者の基本情報を更新
		userAPI.PUT("/update/:agent_staff_id", routes.UpdateAgentStaff(db, firebase, r.cfg.Sendgrid))

		// 担当者のログイン情報を更新（メールアドレス）
		userAPI.PUT("/update/email/:agent_staff_id", routes.UpdateAgentStaffEmail(db, firebase, r.cfg.Sendgrid))

		// 担当者のログイン情報を更新（パスワード）
		userAPI.PUT("/update/password/:agent_staff_id", routes.UpdateAgentStaffPassword(db, firebase, r.cfg.Sendgrid))

		// 利用停止
		userAPI.PUT("/update/usage_end", routes.UpdateAgentStaffUsageEnd(db, firebase, r.cfg.Sendgrid))

		// 利用再開
		userAPI.PUT("/update/usage_restart/:agent_staff_id", routes.UpdateAgentStaffUsageReStart(db, firebase, r.cfg.Sendgrid))

		// メール通知（求職者）の更新 body: {agent_staff_id, notification_job_seeker}
		userAPI.PUT("/update/notification_job_seeker", routes.UpdateAgentStaffNotificationJobSeeker(db, firebase, r.cfg.Sendgrid))

		// メール通知（未処理・未読）の更新 body: {agent_staff_id, notification_unwatched}
		userAPI.PUT("/update/notification_unwatched", routes.UpdateAgentStaffNotificationUnwatched(db, firebase, r.cfg.Sendgrid))

		// 管理権限の更新 body: {agent_staff_id, authority}
		userAPI.PUT("/update/authority", routes.UpdateAgentStaffAuthority(db, firebase, r.cfg.Sendgrid))

		// 担当者削除　*firebaseのアイパスを削除&DBのis_deletedをtrueにする
		userAPI.PUT("/delete", routes.DeleteAgentStaff(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/

		// 自身の担当者情報を取得
		userAPI.GET("/me", routes.GetAgentStaffMe(db, firebase, r.cfg.Sendgrid))

		// 指定IDのエージェントの担当者一覧の取得
		userAPI.GET("/:agent_id/staff_list", routes.GetAgentStaffListByAgentID(db, firebase, r.cfg.Sendgrid))

		// 指定IDのエージェントの担当者一覧の取得 *ログイン画面から遷移した場合、Firebase.auth.currentUserがnullになるため、agentStaffIDを直接渡す
		userAPI.GET("/:agent_id/staff_list/order_by_id/:agent_staff_id", routes.GetAgentStaffListByAgentIDOrderByID(db, firebase, r.cfg.Sendgrid))

		// 自社エージェントのIDとアライアンスエージェントのIDで担当者一覧を取得（自分は除く）
		userAPI.GET("/agent/alliance/other_staff_list", routes.GetOhterAgentStaffListByAgentIDAndAllianceAgentID(db, firebase, r.cfg.Sendgrid))

		// 個人売上が未作成の担当者一覧を取得
		userAPI.GET("/list/agent/:agent_id/not_created_sale/:management_id", routes.GetAgentStaffListWithSaleNotCreated(db, firebase, r.cfg.Sendgrid))

		// すべてのエージェント担当者を取得
		userAPI.GET("/all", routes.GetAllAgentStaffList(db, firebase, r.cfg.Sendgrid))

		// 指定AgentIDの利用可能の担当者一覧を取得
		userAPI.GET("/list/agent/usage_status_available/:agent_id", routes.GetAgentStaffListByAgentIDAndUsageStatusAvailable(db, firebase, r.cfg.Sendgrid))

		// 指定AgentIDの未削除の担当者一覧を取得
		userAPI.GET("/list/agent/is_deleted_false/:agent_id", routes.GetAgentStaffListByAgentIDAndIsDeletedFalse(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	/// Agent API
	//
	agentAPI := noAuthAPI.Group("/agent")
	{
		/************************************** POSTメソッド **************************************/
		// エージェントの担当者アカウントの作成
		agentAPI.POST("/:agent_id/staff/signup", routes.AgentStaffSignUp(db, firebase, r.cfg.Sendgrid))

		// エージェントのアカウント作成
		agentAPI.POST("/agent_and_agent_staff", routes.AgentAndAgentStaffSignUp(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/
		// エージェント情報を更新
		agentAPI.PUT("/update/:agent_id", routes.UpdateAgent(db, firebase, r.cfg.Sendgrid))

		// LINE連携情報の更新
		agentAPI.PUT("/line_channel/update", routes.UpdateAgentLineChannel(db, firebase, r.cfg.Sendgrid))

		// 同意書ファイルのURLを更新
		agentAPI.PUT("/agreement_file/update", routes.UpdateAgentAgreementFileURL(db, firebase, r.cfg.Sendgrid))

		// CRM機能を更新
		agentAPI.PUT("/admin/update", routes.UpdateAgentForAdmin(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/
		// エージェント情報の取得
		agentAPI.GET("/:agent_id", routes.GetAgentByAgentID(db, firebase, r.cfg.Sendgrid))

		// 指定UUIDのエージェントの担当者一覧の取得
		agentAPI.GET("/uuid/:agent_uuid", routes.GetAgentByAgentUUID(db, firebase, r.cfg.Sendgrid))

		// LINE情報を取得
		agentAPI.GET("/line_channel/:agent_id", routes.GetAgentLineChannelByAgentID(db, firebase, r.cfg.Sendgrid))

		// agent uuidでLINEChannelIDを取得
		agentAPI.GET("/line_login_channel_id/uuid/:agent_uuid", routes.GetAgentLineLoginChannelIDByAgentUUID(db, firebase, r.cfg.Sendgrid))

		// すべてのエージェントを取得
		agentAPI.GET("/all", routes.GetAllAgentList(db, firebase, r.cfg.Sendgrid))

		// 自社とアライアンスを組んでいるエージェントの一覧を取得
		agentAPI.GET("/alliance/list/:agent_id", routes.GetAllianceAgentListByAgentID(db, firebase, r.cfg.Sendgrid))

		// 自社とアライアンスを組んでいるエージェントの一覧を取得（セレクトボックス用）
		agentAPI.GET("/alliance/list/select/:agent_id", routes.GetAllianceAgentListByAgentIDForSelect(db, firebase, r.cfg.Sendgrid))

		// エージェントの同意書ファイルのURLを取得する
		agentAPI.GET("/agreement_file_url", routes.GetAgreementFileURL(db, firebase, r.cfg.Sendgrid))

		// 請求書リストを取得する
		agentAPI.GET("/claim/list/page", routes.GetAgentClaimList(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	/// AgentAlliance API
	//
	agentAllianceAPI := noAuthAPI.Group("/agent_alliance")
	{
		/************************************** POSTメソッド **************************************/
		// エージェントアライアンス情報の作成
		agentAllianceAPI.POST("/request", routes.RequestAlliance(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/
		// エージェントアライアンス情報の更新
		agentAllianceAPI.PUT("/update_or_create/request/list", routes.UpdateOrCreateAgentAllianceList(db, firebase, r.cfg.Sendgrid))

		// エージェントアライアンス情報の解除申請
		agentAllianceAPI.PUT("/update/cancel_request", routes.UpdateAgentAllianceCancelRequest(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/
		// IDからエージェントアライアンス情報の取得
		agentAllianceAPI.GET("/:agent_alliance_id", routes.GetAgentAllianceByID(db, firebase, r.cfg.Sendgrid))

		// エージェントIDからエージェントアライアンス情報の取得
		agentAllianceAPI.GET("/list/agent/:agent_id", routes.GetAgentAllianceListByAgentID(db, firebase, r.cfg.Sendgrid))

		// アライアンス未申請がひとつでも存在することを確認する
		agentAllianceAPI.GET("/check/without_application", routes.CheckAnyAllianceWithoutApplication(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	// JobSeekerSchedule API
	//
	jobSeekerScheduleAPI := noAuthAPI.Group("/job_seeker_schedule")
	{
		/************************************** POSTメソッド **************************************/

		// 求職者の日程情報の作成
		jobSeekerScheduleAPI.POST("/create", routes.CreateJobSeekerSchedule(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/

		// 求職者の日程情報の更新
		jobSeekerScheduleAPI.PUT("/update/:schedule_id", routes.UpdateJobSeekerSchedule(db, firebase, r.cfg.Sendgrid))

		// 求職者の日程情報の更新
		jobSeekerScheduleAPI.PUT("/share/ca_staff/:job_seeker_id", routes.ShareScheduleToCAStaff(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/

		// 求職者の日程情報の取得
		jobSeekerScheduleAPI.GET("/list/:job_seeker_id", routes.GetJobSeekerScheduleListByJobSeekerID(db, firebase, r.cfg.Sendgrid))

		// 求職者の日程情報の取得
		jobSeekerScheduleAPI.GET("/type_list/:job_seeker_id", routes.GetJobSeekerScheduleTypeListByJobSeekerID(db, firebase, r.cfg.Sendgrid))

		/************************************** DELETEメソッド **************************************/

		// 求職者の日程情報の更新
		jobSeekerScheduleAPI.DELETE("/delete/:schedule_id", routes.DeleteJobSeekerSchedule(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	/// AgentRobot API
	//
	agentRobotAPI := noAuthAPI.Group("/agent_robot")
	{
		// RPA用エージェントロボット情報の作成
		agentRobotAPI.POST("/create", routes.CreateAgentRobot(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// RPA用エージェントロボット情報の更新
		agentRobotAPI.PUT("/update/:agent_robot_id", routes.UpdateAgentRobot(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// IDからRPA用エージェントロボット情報の取得
		agentRobotAPI.GET("/:agent_robot_id", routes.GetAgentRobotByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェントIDからRPA用エージェントロボット情報の取得
		agentRobotAPI.GET("/list/agent/:agent_id", routes.GetAgentRobotListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// エージェントの流入経路マスタ API
	//
	agentInflowChannelOptionAPI := noAuthAPI.Group("/agent_inflow_channel_option")
	{
		// エージェントの流入経路マスタの作成
		agentInflowChannelOptionAPI.POST("/create", routes.CreateAgentInflowChannelOption(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェントの流入経路マスタの更新
		agentInflowChannelOptionAPI.PUT("/update/:agent_inflow_channel_option_id", routes.UpdateAgentInflowChannelOption(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// IDからエージェントの流入経路マスタの取得
		agentInflowChannelOptionAPI.GET("/:agent_inflow_channel_option_id", routes.GetAgentInflowChannelOptionByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェントIDからエージェントの流入経路マスタの取得
		agentInflowChannelOptionAPI.GET("/list/agent/:agent_id", routes.GetAgentInflowChannelOptionListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// 求職者 API
	//
	jobSeekerAPI := noAuthAPI.Group("/job_seeker")
	{
		/************************************** POSTメソッド **************************************/

		// 求職者情報の登録
		jobSeekerAPI.POST("/create/:agent_staff_id", routes.CreateJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者資料の登録
		jobSeekerAPI.POST("/document", routes.CreateJobSeekerDocument(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 面談前アンケートの登録 (求職者の同意項目、業界、職種、勤務地、質問要望、ファイル（履歴書（原本）、職務経歴書（原本）））
		jobSeekerAPI.POST("/create/initial_questionnaire", routes.CreateInitialQuestionnaire(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		//csvファイルを使って求職者情報の登録
		jobSeekerAPI.POST("/import_csv/:agent_id", routes.ImportJobSeekerCSV(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		//csvファイルを使って求職者情報の登録 *プレビュー用
		jobSeekerAPI.POST("/preview_csv/:agent_id", routes.PreviewJobSeekerCSV(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		//csvファイルを使って求職者情報の出力
		jobSeekerAPI.POST("/export_csv/:agent_id", routes.ExportJobSeekerCSV(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のUUIDと名前から一致するか確認　*面談設定フォームの署名時に使用
		jobSeekerAPI.POST("/check/uuid/name", routes.CheckJobSeekerByUUIDAndName(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者UUIDとメールアドレスが一致するか確認&送信 *ゲスト求職者のパスワード変更時に使用
		jobSeekerAPI.POST("/send/reset_password_email", routes.SendJobSeekerResetPasswordEmail(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のお問い合わせメール送信
		jobSeekerAPI.POST("/send/contact", routes.SendJobSeekerContact(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		/************************************** PUTメソッド **************************************/

		// 求職者情報の更新
		jobSeekerAPI.PUT("/update/:job_seeker_id/:agent_staff_id", routes.UpdateJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者資料の更新
		jobSeekerAPI.PUT("/document/:job_seeker_id", routes.UpdateJobSeekerDocument(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者資料の更新 + 現在のタスクが「書類選考/応募書類準備」の場合は「書類選考/エントリー依頼」に変える
		jobSeekerAPI.PUT("/document/job_seeker/:job_seeker_uuid/job_information/:job_information_uuid", routes.UpdateJobSeekerDocumentForTask(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のLINEIDを更新 body:{ job_seeker_uuid, agent_uuid, code }
		jobSeekerAPI.PUT("/update/line_id", routes.UpdateJobSeekerLineID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// アクティビティメモを更新
		jobSeekerAPI.PUT("/update/activity_memo/:job_seeker_id", routes.UpdateActivityMemoByJobSeekerID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のマッチング求人を閲覧可能かを管理する値を更新
		jobSeekerAPI.PUT("/update/can_view_matching_job", routes.UpdateCanViewMatchingJob(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のログイン情報を更新
		jobSeekerAPI.PUT("/update/password", routes.UpdateJobSeekerPassword(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 面談日時の更新
		jobSeekerAPI.PUT("/update/interview_date", routes.UpdateInterviewDateByJobSeekerID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		/************************************** GETメソッド **************************************/

		// 求職者IDから求職者情報の取得
		jobSeekerAPI.GET("/:job_seeker_id", routes.GetJobSeekerByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のuuidから求職者情報の取得
		jobSeekerAPI.GET("/uuid/:job_seeker_uuid", routes.GetJobSeekerByUUID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のuuidから求職者の応募書類情報の取得
		jobSeekerAPI.GET("/document/uuid/:job_seeker_uuid", routes.GetJobSeekerDocumentByUUID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のuuidから求職者情報の取得
		jobSeekerAPI.GET("/task_group/:task_group_uuid", routes.GetJobSeekerByTaskGroupUUID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者IDから求職者資料の取得
		jobSeekerAPI.GET("/document/:job_seeker_id", routes.GetJobSeekerDocumentByJobSeekerID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者の絞り込み検索（全ての自社求職者）
		jobSeekerAPI.GET("/list/search/:agent_id", routes.GetSearchJobSeekerListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者の絞り込み検索（自社求職者（転職活動中&面談実施））
		jobSeekerAPI.GET("/active/list/search/:agent_id", routes.GetSearchActiveJobSeekerListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		//アライアンス求職者の絞り込み検索
		jobSeekerAPI.GET("/list/search/alliance/:agent_id", routes.GetSearchAllianceJobSeekerListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// job_seeker_idのリストで求職者の一覧を取得
		jobSeekerAPI.GET("/list/by/id_list", routes.GetJobSeekerListByIDList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// クエリパラム（last_name, first_name, last_furigana, first_furigana, email, phone_number）に合致する求職者情報を取得
		jobSeekerAPI.GET("/duplicate/search", routes.GetDuplicateJobSeekerList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者編集ページで使用されるセレクトボックスに必要なデータを取得（自社の担当者一覧、アライアンスエージェント一覧、流入経路）
		jobSeekerAPI.GET("/list_used_for_create_or_update/agent/:agent_id", routes.GetSelectListForCreateOrUpdateJobSeekerByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求人検索→求職者検索(絞り込み)
		jobSeekerAPI.GET("/search_list/page/agent/:agent_id/type", routes.GetSearchJobSeekerListByAgentIDAndType(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// シェア求人検索→自社求職者検索(絞り込み)
		jobSeekerAPI.GET("/public/search_list/page/agent/:agent_id", routes.GetSearchPublicJobSeekerListByAgentIDAndPage(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のuuidから求職者情報（初回面談~LINE連携）の取得
		jobSeekerAPI.GET("/initial_step/uuid/:job_seeker_uuid", routes.GetJobSeekerForInitialStepByUUID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のuuidから求職者情報（ゲスト用）の取得
		jobSeekerAPI.GET("/guest/uuid/:job_seeker_uuid", routes.GetGuestJobSeekerForByUUID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のuuidから求職者情報（ゲスト用）の取得
		jobSeekerAPI.GET("/guest/desired/uuid/:job_seeker_uuid", routes.GetJobSeekerDesiredForGuestByUUID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 求職者のUUIDからエージェントIDを取得 ＊ゲストログイン用
		jobSeekerAPI.GET("/uuid/agent_id/:job_seeker_uuid", routes.GetJobSeekerAgentIDByUUID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		/************************************** DELETEメソッド **************************************/

		// 求職者情報の削除
		jobSeekerAPI.DELETE("/delete", routes.DeleteJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 履歴書PDFの削除（resume_pdf_urlカラムを空文字で更新）
		jobSeekerAPI.DELETE("/:job_seeker_id/document/resume_pdf_url/delete", routes.DeleteJobSeekerResumePDFURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 履歴書原本の削除（resume_origin_urlカラムを空文字で更新）
		jobSeekerAPI.DELETE("/:job_seeker_id/document/resume_origin_url/delete", routes.DeleteJobSeekerResumeOriginURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 職務経歴書PDFの削除（cv_pdf_urlカラムを空文字で更新）
		jobSeekerAPI.DELETE("/:job_seeker_id/document/cv_pdf_url/delete", routes.DeleteJobSeekerCVPDFURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 職務経歴書原本の削除（cv_origin_urlカラムを空文字で更新）
		jobSeekerAPI.DELETE("/:job_seeker_id/document/cv_origin_url/delete", routes.DeleteJobSeekerCVOriginURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 推薦状PDFの削除（recommendation_pdf_urlカラムを空文字で更新）
		jobSeekerAPI.DELETE("/:job_seeker_id/document/recommendation_pdf_url/delete", routes.DeleteJobSeekerRecommendationPDFURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 推薦状原本の削除（recommendation_origin_urlカラムを空文字で更新）
		jobSeekerAPI.DELETE("/:job_seeker_id/document/recommendation_origin_url/delete", routes.DeleteJobSeekerRecommendationOriginURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// 証明写真の削除（id_photo_urlカラムを空文字で更新）
		jobSeekerAPI.DELETE("/:job_seeker_id/document/id_photo_url/delete", routes.DeleteJobSeekerIDPhotoURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// その他①の削除（other_document1_urlカラムを空文字で更新）
		jobSeekerAPI.DELETE("/:job_seeker_id/document/other_document1_url/delete", routes.DeleteJobSeekerOtherDocument1URL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// その他②の削除（other_document2_urlカラムを空文字で更新）
		jobSeekerAPI.DELETE("/:job_seeker_id/document/other_document2_url/delete", routes.DeleteJobSeekerOtherDocument2URL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// その他③の削除（other_document3_urlカラムを空文字で更新）
		jobSeekerAPI.DELETE("/:job_seeker_id/document/other_document3_url/delete", routes.DeleteJobSeekerOtherDocument3URL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))
	}

	/****************************************************************************************/
	/// 企業 API
	//
	enterpriseAPI := noAuthAPI.Group("/enterprise")
	{
		/************************************** POSTメソッド **************************************/

		// 企業情報の登録
		enterpriseAPI.POST("/create", routes.CreateEnterpriseProfile(db, firebase, r.cfg.Sendgrid))

		// 企業資料の登録
		enterpriseAPI.POST("/reference_material", routes.CreateEnterpriseReferenceMaterial(db, firebase, r.cfg.Sendgrid))

		// 企業の追加情報を登録
		enterpriseAPI.POST("/activity", routes.CreateEnterpriseActivity(db, firebase, r.cfg.Sendgrid))

		//csvファイルを使って企業情報を登録  企業・請求先
		enterpriseAPI.POST("/import_csv/:agent_id", routes.ImportEnterpriseCSV(db, firebase, r.cfg.Sendgrid))

		//csvファイルを使って企業情報を登録  企業・請求先　プレビュー用
		enterpriseAPI.POST("/preview_csv/:agent_id", routes.PreviewEnterpriseCSV(db, firebase, r.cfg.Sendgrid))

		//csvファイルを使って企業情報を登録  求人
		enterpriseAPI.POST("/import_csv/job_info/:agent_id", routes.ImportJobInformationCSV(db, firebase, r.cfg.Sendgrid))

		//csvファイルを使って企業情報を登録  企業・請求先・求人 (サーカスエージェント)
		enterpriseAPI.POST("/import_csv/circus/:agent_id/:agent_staff_id", routes.ImportEnterpriseCSVForCircus(db, firebase, r.cfg.Sendgrid))

		//csvファイルを使って企業情報を登録  企業・請求先・求人 (エージェントバンク)
		enterpriseAPI.POST("/import_csv/agent_bank/:agent_id/:agent_staff_id", routes.ImportEnterpriseCSVForAgentBank(db, firebase, r.cfg.Sendgrid))

		//csvファイルを使って企業情報を登録  求人　プレビュー用
		enterpriseAPI.POST("/preview_csv/job_info/:agent_id", routes.PreviewJobInformationCSV(db, firebase, r.cfg.Sendgrid))

		// 企業と求人リストから作成
		enterpriseAPI.POST("/json/:agent_staff_id", routes.ImportEnterpriseJSON(db, firebase, r.cfg.Sendgrid))

		//企業情報をcsvファイルに出力
		enterpriseAPI.POST("/export_csv/:agent_id", routes.ExportEnterpriseCSV(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/

		// 企業情報の更新
		enterpriseAPI.PUT("/update/:enterprise_id", routes.UpdateEnterpriseProfile(db, firebase, r.cfg.Sendgrid))

		// 企業資料の更新
		enterpriseAPI.PUT("/reference_material/:enterprise_id", routes.UpdateEnterpriseReferenceMaterial(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/

		// 企業IDを使って企業情報取得する関数
		enterpriseAPI.GET("/:enterprise_id", routes.GetEnterpriseByID(db, firebase, r.cfg.Sendgrid))

		// 企業IDを使って企業資料取得する関数
		enterpriseAPI.GET("/reference_material/:enterprise_id", routes.GetEnterpriseReferenceMaterialByEnterpriseID(db, firebase, r.cfg.Sendgrid))

		// agentStaffのIDを使って企業情報一覧を取得
		enterpriseAPI.GET("/list/agent_staff/:agent_staff_id", routes.GetEnterpriseListByAgentStaffID(db, firebase, r.cfg.Sendgrid))

		// agentIDを使って企業情報一覧を取得
		enterpriseAPI.GET("/list/agent/:agent_id", routes.GetEnterpriseListByAgentID(db, firebase, r.cfg.Sendgrid))

		//エージェントIDとページ番号で50件取得
		enterpriseAPI.GET("/list/agent/page/:agent_id", routes.GetEnterpriseListByAgentIDAndPage(db, firebase, r.cfg.Sendgrid))

		//企業の絞り込み検索
		enterpriseAPI.GET("/list/search/:agent_id", routes.GetSearchEnterpriseListByAgentID(db, firebase, r.cfg.Sendgrid))

		/************************************** DELETEメソッド **************************************/

		// 企業情報の削除
		enterpriseAPI.DELETE("/delete", routes.DeleteEnterpriseProfile(db, firebase, r.cfg.Sendgrid))

		// 企業の参考資料の削除(material_type = 1 or 2)
		enterpriseAPI.DELETE("/:reference_material_id/material/:material_type/delete", routes.DeleteEnterpriseReferenceMaterial(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/

	/// 請求先 API
	//
	billingAddressAPI := noAuthAPI.Group("/billing_address")
	{
		/************************************** POSTメソッド **************************************/

		// 請求先情報の登録
		billingAddressAPI.POST("/create/:enterprise_id", routes.CreateBillingAddress(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/

		// 請求先情報の更新
		billingAddressAPI.PUT("/update/:billing_address_id", routes.UpdateBillingAddress(db, firebase, r.cfg.Sendgrid))

		// 複数の請求先のagent_staff_idカラムを一括更新
		billingAddressAPI.PUT("/update/list/agent_staff", routes.UpdateBillingAddressStaffIDByIDListAtOnce(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/

		// IDを使って請求先情報取得する関数
		billingAddressAPI.GET("/:billing_address_id", routes.GetBillingAddressByID(db, firebase, r.cfg.Sendgrid))

		// 企業IDを使って請求先情報一覧取得する関数
		billingAddressAPI.GET("/list/:enterprise_id", routes.GetBillingAddressListByEnterpriseID(db, firebase, r.cfg.Sendgrid))

		// AgentIDを使って請求先と企業名情報一覧取得する関数
		billingAddressAPI.GET("/list/page/agent/:agent_id", routes.GetBillingAddressListByPageAndAgentID(db, firebase, r.cfg.Sendgrid))

		//企業の絞り込み検索
		billingAddressAPI.GET("/list/search/agent/:agent_id", routes.GetSearchBillingAddressListByPageAndAgentID(db, firebase, r.cfg.Sendgrid))

		/************************************** DELETEメソッド **************************************/

		// 請求先情報の削除
		billingAddressAPI.DELETE("/delete/:billing_address_id", routes.DeleteBillingAddress(db, firebase, r.cfg.Sendgrid))

	}

	/****************************************************************************************/
	/// 求人 API
	//
	jobInformationAPI := noAuthAPI.Group("/job_information")
	{
		/************************************** POSTメソッド **************************************/

		// 求人情報の登録
		jobInformationAPI.POST("/create/:billing_address_id", routes.CreateJobInformation(db, firebase, r.cfg.Sendgrid))

		// 選考フローの登録
		jobInformationAPI.POST("/selection_flow/create", routes.CreateSelectionFlowPattern(db, firebase, r.cfg.Sendgrid))

		// 求職者のUUIDから求人リスト、合致求人数と年収期待値
		jobInformationAPI.POST("/list/search/job_seeker_uuid", routes.GetSearchJobListingListByJobSeekerUUID(db, firebase, r.cfg.Sendgrid))

		// 求職者のエントリー希望と興味あり求人の取得
		jobInformationAPI.POST("/list/interested_type/job_seeker_uuid", routes.GetJobListingListByJobSeekerUUIDAndInterestedType(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/

		// 求人情報の更新
		jobInformationAPI.PUT("/update/:job_information_id", routes.UpdateJobInformation(db, firebase, r.cfg.Sendgrid))

		// 選考フローの更新
		jobInformationAPI.PUT("/selection_flow/update/:selection_flow_id", routes.UpdateSelectionFlowPattern(db, firebase, r.cfg.Sendgrid))

		// 選考フローの削除
		jobInformationAPI.PUT("/delete/selection_flow/:selection_flow_id", routes.DeltedSelectionFlowPattern(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/

		// 求人IDを使って求人情報取得する関数
		jobInformationAPI.GET("/:job_information_id", routes.GetJobInformationByID(db, firebase, r.cfg.Sendgrid))

		// 求人のuuidを使って求人情報取得する関数
		jobInformationAPI.GET("/uuid/:job_information_uuid", routes.GetJobInformationByUUID(db, firebase, r.cfg.Sendgrid))

		// 求人のuuidを使って求人票の情報を取得する関数
		jobInformationAPI.GET("/job_listing/job_information_uuid/:job_information_uuid", routes.GetJobListingByJobInformationUUID(db, firebase, r.cfg.Sendgrid))

		// 求職者が確認する求人票情報を取得（求人票 + タスクに紐づいた選考情報）
		jobInformationAPI.GET("/job_listing/for_job_seeker", routes.GetJobListingForJobSeeker(db, firebase, r.cfg.Sendgrid))

		// 請求先IDを使って求人情報一覧取得する関数
		jobInformationAPI.GET("/list/billing_address/:billing_address_id", routes.GetJobInformationListByBillingAddressID(db, firebase, r.cfg.Sendgrid))

		// 請求先IDを使って求人情報一覧取得する関数
		jobInformationAPI.GET("/list/enterprise/:enterprise_id", routes.GetJobInformationListByEnterpriseID(db, firebase, r.cfg.Sendgrid))

		// agentIDを使って求人情報一覧を取得する関数
		jobInformationAPI.GET("/list/agent/:agent_id", routes.GetJobInformationListByAgentID(db, firebase, r.cfg.Sendgrid))

		// 選考IDを使って「選考情報 + アンケート情報 + 評価点（タスク）」を取得を取得する
		jobInformationAPI.GET("/selection_flow/:selection_flow_id", routes.GetSelectionFlowPatternByID(db, firebase, r.cfg.Sendgrid))

		//求人の絞り込み検索
		jobInformationAPI.GET("/active/list/search/:agent_id", routes.GetSearchActiveJobInformationListByAgentID(db, firebase, r.cfg.Sendgrid))

		//アライアンス求人の絞り込み検索
		jobInformationAPI.GET("/list/search/alliance/:agent_id", routes.GetSearchJobInformationListByOtherAgentID(db, firebase, r.cfg.Sendgrid))

		// job_information_idのリストで求職者の一覧を取得
		jobInformationAPI.GET("/list/by/id_list", routes.GetJobInformationListByIDList(db, firebase, r.cfg.Sendgrid))

		//求人IDを使って選考フローパターンの一覧を取得する
		jobInformationAPI.GET("/selection_flow/list/:job_information_id", routes.GetSelectionFlowPatternListByJobInformationID(db, firebase, r.cfg.Sendgrid))

		// 求人IDを使ってOpenになっている選考フローパターンの一覧を取得する
		jobInformationAPI.GET("/open/selection_flow/list/:job_information_id", routes.GetOpenSelectionFlowPatternListByJobInformationID(db, firebase, r.cfg.Sendgrid))

		// 求職者uuidで打診されている求人一覧を取得する
		jobInformationAPI.GET("/job_listing/for/job_seeker/:job_seeker_uuid", routes.GetJobListingListByJobSeekerUUID(db, firebase, r.cfg.Sendgrid))

		// LPから診断に使用する求人リストと求職者の希望条件を取得
		jobInformationAPI.GET("/list/job_seeker/:job_seeker_uuid/for_diagnosis", routes.GetJobListingListAndJobSeekerDesiredForDiagnosis(db, firebase, r.cfg.Sendgrid))

		/****** 求職者検索→求人検索 API ******/

		// 求職者検索→求人検索
		jobInformationAPI.GET("/list/page/agent/:agent_id/type", routes.GetJobInformationListByAgentIDAndType(db, firebase, r.cfg.Sendgrid))

		// 求職者検索→求人検索(絞り込み)
		jobInformationAPI.GET("/search_list/page/agent/:agent_id/type", routes.GetSearchJobInformationListByAgentIDAndType(db, firebase, r.cfg.Sendgrid))

		/****** シェア求職者検索→自社求人検索 API ******/

		// シェア求職者検索→自社求人検索(絞り込み)
		jobInformationAPI.GET("/public/search_list/page/agent/:agent_id", routes.GetSearchPublicJobInformationListByAgentIDAndPage(db, firebase, r.cfg.Sendgrid))

		/************************************** DELETEメソッド **************************************/

		// 求人情報の削除
		jobInformationAPI.DELETE("/delete/:job_information_id", routes.DeleteJobInformation(db, firebase, r.cfg.Sendgrid))

	}

	/****************************************************************************************/
	/// task API
	//
	taskAPI := noAuthAPI.Group("/task")
	{
		// タスクの取得
		taskAPI.GET("/:task_id", routes.GetTaskByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// タスクグループの取得（自分が関わっているタスク）
		taskAPI.GET("/group/:task_group_id", routes.GetTaskGroupByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// タスクグループの一覧取得（エージェントが関わっているタスク）
		taskAPI.GET("/list/agent/page/:agent_id", routes.GetTaskListByAgentIDAndPage(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// タスクグループの一覧を絞り込みで取得（エージェントが関わっているタスク）
		taskAPI.GET("/list/search/:agent_id", routes.GetSearchTaskListByAgentIDAndPage(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// アクティブなタスクの数を取得（請求先のID）
		taskAPI.GET("/active_task/count/billing_address/:billing_address_id", routes.GetActiveTaskCountByBillingAddressID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// アクティブなタスクの数を取得（求人のID）
		taskAPI.GET("/active_task/count/job_information/:job_information_id", routes.GetActiveTaskCountByJobInformationID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// アクティブなタスクの数を取得（選考フローのID）
		taskAPI.GET("/active_task/count/selection/:selection_id", routes.GetActiveTaskCountBySelectionID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 最新のフェーズが書類選考以降のタスク一覧
		taskAPI.GET("/after_entry/list/job_seeker/:job_seeker_id", routes.GetTaskListAfterEntryByJobSeekerID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 求人打診可能なグループと不可能なグループを取得
		taskAPI.GET("/sound_out_group/list", routes.GetSoundOutGroupList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// トップページのタスク一覧を取得するapi
		taskAPI.GET("/job_seeker_list/agent_staff/:agent_staff_id", routes.GetJobSeekerTaskListByAgentStaffID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** POSTメソッド **************************************/

		/************* タスク開始 API *************/
		// Memo: 旧版の為、新しいapiの動作確認後に削除（が対応）

		// 求人打診
		taskAPI.POST("/sound_out/for/job_information", routes.SoundOutForJobInformation(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// マスクレジュメの打診依頼
		taskAPI.POST("/sound_out/for/mask_resume", routes.SoundOutForMaskResume(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 求職者に求人票を直接送信
		taskAPI.POST("/sound_out/for/send_job_listing", routes.SoundOutForSendJobListing(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// シェア依頼処理（求人検索ページから）
		taskAPI.POST("/sound_out/for/request_share/for/job_information_search", routes.SoundOutForRequestShareJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// シェア依頼処理（求職者検索ページから）
		taskAPI.POST("/sound_out/for/request_share/for/job_seeker_search", routes.SoundOutForRequestShareJobInformation(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************* 新タスク開始 API *************/

		// 求人打診
		taskAPI.POST("/sound_out_group/for/job_information", routes.SoundOutGroupForJobInformation(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 求職者に求人票を直接送信
		taskAPI.POST("/sound_out_group/for/send_job_listing", routes.SoundOutGroupForSendJobListing(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/***************************************/

		// 次のタスク作成（現在のタスクがEntryの場合）
		taskAPI.POST("/create/batch_processing", routes.CreateTaskInBatchProcessing(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 次のタスク作成（現在のタスクがEntryの場合）
		taskAPI.POST("/entry/create/:agent_staff_id", routes.CreateNextTaskAfterEntryPhase(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 次のタスク作成（現在のタスクが書類選考の場合）
		taskAPI.POST("/document_selection/create/:agent_staff_id", routes.CreateNextTaskAfterDocumentSelectionPhase(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 次のタスク作成（現在のタスクが選考（1次-最終）の場合）
		taskAPI.POST("/selection/create/:agent_staff_id", routes.CreateNextTaskAfterSelectionPhase(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 次のタスク作成（現在のタスクが辞退の場合）
		taskAPI.POST("/decline/create/:agent_staff_id", routes.CreateNextTaskAfterDeclinePhase(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 次のタスク作成（現在のタスクが内定保留の場合）
		taskAPI.POST("/hold_job_offer/create/:agent_staff_id", routes.CreateNextTaskAfterHoldJobOfferPhase(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 次のタスク作成（現在のタスクが内定承諾の場合）
		taskAPI.POST("/accept_job_offer/create/:agent_staff_id", routes.CreateNextTaskAfterAcceptJobOfferPhase(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 次のタスク作成（同一タスクをまとめて処理する場合）
		taskAPI.POST("/same_list/create/:agent_staff_id", routes.CreateNextSameTaskList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// マイページのマッチ求人からエントリー {job_seeker_uuid, job_information_uuid, is_entry}
		taskAPI.POST("/entry/create/matchin_job", routes.CreateEntryTaskFromMatchingJob(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** PUTメソッド **************************************/

		// 求職者資料の登録
		taskAPI.PUT("/document", routes.UpdateTaskGroupDocument(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// RAの最終閲覧時間を更新
		taskAPI.PUT("/update/ra/last_watched/:group_id", routes.UpdateRALastWatched(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// RAの最終依頼時間を更新
		taskAPI.PUT("/update/ra/last_request/:group_id", routes.UpdateRALastRequest(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// CAの最終閲覧時間を更新
		taskAPI.PUT("/update/ca/last_watched/:group_id", routes.UpdateCALastWatched(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// CAの最終依頼時間を更新
		taskAPI.PUT("/update/ca/last_request/:group_id", routes.UpdateCALastRequest(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 外部求人情報を更新する
		taskAPI.PUT("/external_job/update/task_group/:group_id", routes.UpdateExternalJob(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/

		// タスクの取得
		taskAPI.GET("/latest_list/agent/:agent_id/agent_staff/:agent_staff_id", routes.GetLatestTaskListByAgentStaffID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// タスクの取得
		taskAPI.GET("/latest/job_seeker/:job_seeker_id/job_information/:job_information_id", routes.GetLatestTaskByJobSeekerIDAndJobInformationID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 「同一求職者」で「同一フェーズ」のタスクのリストを取得
		taskAPI.GET("/latest_list/same_phase/job_seeker/:job_seeker_id", routes.GetLatestSameTaskListByJobSeekerID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** DELETEメソッド **************************************/

		// タスクの削除
		taskAPI.DELETE("/delete", routes.DeleteTask(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// 面談調整タスク API
	//
	interviewTaskAPI := noAuthAPI.Group("/interview_task")
	{
		/************************************** PUTメソッド **************************************/
		// 最終閲覧時間を更新
		interviewTaskAPI.POST("/create", routes.CreateNextInterviewTask(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** DELETEメソッド **************************************/

		//面談タスクの削除
		interviewTaskAPI.DELETE("/delete/latest", routes.DeleteLatestInterviewTask(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** PUTメソッド **************************************/
		// 最終閲覧時間を更新
		interviewTaskAPI.PUT("/update/last_watched/:group_id", routes.UpdateInterviewTaskGroupLastWatched(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 求職者の担当CAを更新
		interviewTaskAPI.PUT("/update/task/ca_staff_id", routes.UpdateInterviewTaskCAStaffID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 求職者の担当CAを更新
		interviewTaskAPI.PUT("/update/task/interview_date", routes.UpdateInterviewTaskInterviewDate(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/
		// 面談調整タスク（phase_category = (0 or 1)）の取得
		interviewTaskAPI.GET("/adjustment/latest_list/:agent_id", routes.GetLatestAdjustmentTaskListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 参加確認（phase_category = (2 or 3)）タスクの取得
		interviewTaskAPI.GET("/confirmation/latest_list/:agent_id", routes.GetLatestConfirmationTaskListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// GroupIDから面談調整タスク一覧を取得 *アクティビティ表示に使用
		interviewTaskAPI.GET("/list/group/:group_id", routes.GetInterviewTaskListByGroupID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	// 面談調整用のテンプレート API
	//
	interviewTemplateAPI := noAuthAPI.Group("/interview_template")
	{
		// 面談調整メッセージテンプレートの登録
		interviewTemplateAPI.POST("/create", routes.CreateInterviewAdjustmentTemplate(db, firebase, r.cfg.Sendgrid))

		// 面談調整メッセージテンプレートの更新
		interviewTemplateAPI.PUT("/update/:template_id", routes.UpdateInterviewAdjustmentTemplate(db, firebase, r.cfg.Sendgrid))

		// 面談調整メッセージテンプレートの削除
		interviewTemplateAPI.DELETE("/delete/:template_id", routes.DeleteInterviewAdjustmentTemplate(db, firebase, r.cfg.Sendgrid))

		// IDから面談調整メッセージテンプレートを取得
		interviewTemplateAPI.GET("/:template_id", routes.GetInterviewAdjustmentTemplateByID(db, firebase, r.cfg.Sendgrid))

		//エージェントIDと送信シーンから面談調整メッセージテンプレートを取得
		interviewTemplateAPI.GET("/list/agent/:agent_id/send_scene/:send_scene", routes.GetInterviewAdjustmentTemplateListByAgentIDAndSendScene(db, firebase, r.cfg.Sendgrid))

		// エージェントIDから面談調整メッセージテンプレートを取得
		interviewTemplateAPI.GET("/list/agent/:agent_id", routes.GetInterviewAdjustmentTemplateListByAgentID(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	/// 選考後アンケート API
	//
	selectionQuestionnaireAPI := noAuthAPI.Group("/selection_questionnaire")
	{
		// 選考後アンケートを更新
		selectionQuestionnaireAPI.POST("/create/:questionnaire_uuid", routes.CreateSelectionQuestionnaire(db, firebase, r.cfg.Sendgrid))

		// 選考後アンケートを更新
		selectionQuestionnaireAPI.PUT("/update/:questionnaire_id", routes.UpdateSelectionQuestionnaire(db, firebase, r.cfg.Sendgrid))

		// 選考後アンケートをuuidで取得
		selectionQuestionnaireAPI.GET("/or_null/uuid/:questionnaire_uuid", routes.GetSelectionQuestionnaireOrNullByUUID(db, firebase, r.cfg.Sendgrid))

		// 選考後アンケートをuuidで取得
		selectionQuestionnaireAPI.GET("/generate/uuid", routes.GenerateSelectionQuestionnaireByUUID(db, firebase, r.cfg.Sendgrid))

		// 求職者のuuidで未回答のアンケートを全て取得
		selectionQuestionnaireAPI.GET("/unanswer/for_job_seeker/:job_seeker_uuid", routes.GetUnansweredQuestionnaireListByJobSeekerUUID(db, firebase, r.cfg.Sendgrid))

		// 求職者のアンケート情報を取得する（queryParam: job_seeker_uuid, job_information_uuid, selection_information_id）
		selectionQuestionnaireAPI.GET("/for_job_seeker", routes.GetQuestionnaireForJobSeeker(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	/// メッセージテンプレート API
	//
	messageTemplateAPI := noAuthAPI.Group("/template")
	{
		// メッセージテンプレートの登録
		messageTemplateAPI.POST("/create", routes.CreateMessageTemplate(db, firebase, r.cfg.Sendgrid))

		// メッセージテンプレートの更新
		messageTemplateAPI.PUT("/update/:template_id", routes.UpdateMessageTemplate(db, firebase, r.cfg.Sendgrid))

		// メッセージテンプレートの削除
		messageTemplateAPI.DELETE("/delete/:template_id", routes.DeleteMessageTemplate(db, firebase, r.cfg.Sendgrid))

		// IDからメッセージテンプレートを取得
		messageTemplateAPI.GET("/:template_id", routes.GetMessageTemplateByID(db, firebase, r.cfg.Sendgrid))

		// 担当者IDからメッセージテンプレートを取得
		messageTemplateAPI.GET("/list/staff/:agent_staff_id", routes.GetMessageTemplateListByAgentStaffID(db, firebase, r.cfg.Sendgrid))

		// 担当者IDと送信シーンからメッセージテンプレートを取得
		messageTemplateAPI.GET("/list/staff/:agent_staff_id/send_scene/:send_scene", routes.GetMessageTemplateListByAgentStaffIDAndSendScene(db, firebase, r.cfg.Sendgrid))

		// エージェントIDからメッセージテンプレートを取得
		messageTemplateAPI.GET("/list/agent/:agent_id", routes.GetMessageTemplateListByAgentID(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	/// エージェントと求職者のチャットグループ API
	//
	chatGroupWithJobSeekerAPI := noAuthAPI.Group("/chat_group_with_job_seeker")
	{
		/************************************** POST	メソッド **************************************/

		// エージェントと求職者のチャットグループの登録
		chatGroupWithJobSeekerAPI.POST("/create", routes.CreateChatGroupWithJobSeeker(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/

		// エージェントの最終閲覧時間を更新
		chatGroupWithJobSeekerAPI.PUT("/update/agent/last_watched/:group_id", routes.UpdateAgentLastWatched(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/

		// IDからエージェントと求職者のチャットグループを取得
		chatGroupWithJobSeekerAPI.GET("/:id", routes.GetChatGroupWithJobSeekerByID(db, firebase, r.cfg.Sendgrid))

		// エージェントIDからエージェントと求職者のチャットグループを取得
		chatGroupWithJobSeekerAPI.GET("/list/agent/:agent_id", routes.GetChatGroupWithJobSeekerListByAgentID(db, firebase, r.cfg.Sendgrid))

		// エージェントIDからエージェントと求職者のチャットグループを取得（ページングあり）
		chatGroupWithJobSeekerAPI.GET("/search/list/page/agent/:agent_id", routes.GetChatGroupWithJobSeekerSearchListAndPageByAgentID(db, firebase, r.cfg.Sendgrid))

		// チャットの未読件数を返す
		chatGroupWithJobSeekerAPI.GET("/notification/agent/:agent_id", routes.GetJobSeekerChatNotificationByAgentID(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	/// エージェントと求職者のチャットメッセージAPI
	//
	chatMessageWithJobSeekerAPI := noAuthAPI.Group("/chat_message_with_job_seeker")
	{
		/************************************** POSTメソッド **************************************/

		// エージェントが求職者にLINEでメッセージ送信
		chatMessageWithJobSeekerAPI.POST("/line/message/create", routes.SendChatMessageWithJobSeekerLineMessage(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェントが求職者にLINEでメッセージ送信
		chatMessageWithJobSeekerAPI.POST("/line/image/create", routes.SendChatMessageWithJobSeekerLineImage(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/

		// グループIDからエージェントと求職者のチャットメッセージを取得
		chatMessageWithJobSeekerAPI.GET("/list/group/:group_id", routes.GetChatMessageWithJobSeekerListByGroupID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// エージェントと求職者のメールAPI
	//
	emailWithJobSeekerAPI := noAuthAPI.Group("/email_with_job_seeker")
	{
		/************************************** POSTメソッド **************************************/

		// エージェントが求職者にメール送信
		emailWithJobSeekerAPI.POST("/create", routes.SendEmailWithJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/

		// 求職者IDからエージェントと求職者のメール一覧を取得
		emailWithJobSeekerAPI.GET("/list/job_seeker/:job_seeker_id", routes.GetEmailWithJobSeekerListByJobSeekerID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// エージェント同士のチャットグループ API
	//
	chatGroupWithAgentAPI := noAuthAPI.Group("/chat_group_with_agent")
	{
		/************************************** POST	メソッド **************************************/

		// エージェントとエージェントのチャットグループの登録
		chatGroupWithAgentAPI.POST("/create", routes.CreateChatGroupWithAgent(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェントとエージェントのチャットグループの登録
		chatGroupWithAgentAPI.POST("/create/list", routes.CreateChatGroupWithAgentList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** PUTメソッド **************************************/

		// エージェントの最終閲覧時間を更新
		chatGroupWithAgentAPI.PUT("/update/agent/last_watched/:group_id/:agent_id", routes.UpdateAgentLastWatchedForAgentChat(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/

		// IDからエージェントとエージェントのチャットグループを取得
		chatGroupWithAgentAPI.GET("/:id/agent/:agent_id", routes.GetChatGroupWithAgentByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェントIDと担当者からエージェントとエージェントのチャットグループを取得　未読件数含む
		chatGroupWithAgentAPI.GET("/agent/:agent_id/:agent_staff_id", routes.GetChatGroupWithAgentListByAgentIDAndAgentStaffID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 自社エージェントIDと他社エージェントIDでチャットグループを取得
		chatGroupWithAgentAPI.GET("/my/:agent_id/other/:other_agent_id", routes.GetChatGroupWithAgentByAgentIDAndOtherAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェントIDからエージェントとエージェントのチャットグループを取得
		chatGroupWithAgentAPI.GET("/list/agent/:agent_id", routes.GetChatGroupWithAgentListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// グループIDからエージェントとエージェントのチャットスレッドを取得
		chatGroupWithAgentAPI.GET("/thread/:group_id/:agent_staff_id", routes.GetChatGroupAndThreadWithAgentByChatGroupID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// チャットグループ未作成が一つでも存在することを確認する。
		chatGroupWithAgentAPI.GET("/check/without_group", routes.CheckAnyChatGroupWithoutGroup(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}
	/****************************************************************************************/
	/// エージェント同士のチャットスレッド API
	//
	chatThreadWithAgentAPI := noAuthAPI.Group("/chat_thread_with_agent")
	{
		// /************************************** POST	メソッド **************************************/

		// // エージェントとエージェントのチャットスレッドの登録
		chatThreadWithAgentAPI.POST("/create", routes.CreateChatThreadWithAgent(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// /************************************** PUTメソッド **************************************/

		// // エージェントの最終閲覧時間を更新
		// chatGroupWithAgentAPI.PUT("/update/agent/last_watched/:group_id/:agent_id", routes.UpdateAgentLastWatchedForAgentChat(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// /************************************** GETメソッド **************************************/

		// スレッドIDからスレッドとメッセージ一覧を取得 未読も判定
		chatThreadWithAgentAPI.GET("/:thread_id", routes.GetChatThreadAndMessageWithAgentByChatThreadID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// エージェントチャットのメッセージグループ API
	//
	chatMessageWithAgentAPI := noAuthAPI.Group("/chat_message_with_agent")
	{
		/************************************** POSTメソッド **************************************/
		// メッセージ送信
		chatMessageWithAgentAPI.POST("/create", routes.SendChatMessageWithAgent(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** PUTメソッド **************************************/
		// メッセージ送信
		chatMessageWithAgentAPI.PUT("/update/watched_at", routes.UpdateChatMessageWithAgentWatchedAtByThreadID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/
		// メッセージ一覧を取得
		chatMessageWithAgentAPI.GET("/list/thread/:thread_id/:agent_staff_id", routes.GetChatMessageWithAgentListByThreadID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// Sale API
	//
	saleAPI := noAuthAPI.Group("/sales")
	{
		/************************************** POSTメソッド **************************************/

		// ヨミの登録
		saleAPI.POST("/create", routes.CreateSale(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** PUTメソッド **************************************/

		// ヨミの更新
		saleAPI.PUT("/update/:sale_id", routes.UpdateSale(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/

		// ヨミ情報の取得
		saleAPI.GET("/:sale_id", routes.GetSaleByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 求職者IDでヨミ情報の取得
		saleAPI.GET("/job_seeker/:job_seeker_id", routes.GetSaleByJobSeekerID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// idのリストでヨミ情報リストを取得
		saleAPI.GET("/list/by/id_list", routes.GetSaleListByIDList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 自社が絡むヨミ情報を全て取得
		saleAPI.GET("/accuracy/search/list", routes.GetAccuracySearchList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// カレンダーのスケジュール API
	//
	scheduleAPI := noAuthAPI.Group("/schedule")
	{
		/************************************** GETメソッド **************************************/

		// カレンダに表示するスケジュール一覧を取得
		scheduleAPI.GET("/list/staff/:agent_staff_id", routes.GetScheduleListWithInPeriod(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// カレンダに表示するスケジュール一覧を取得（）
		scheduleAPI.GET("/list/select_staffs", routes.GetScheduleListWithInPeriodByStaffIDList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	// ダッシュボード API
	//
	dashboardAPI := noAuthAPI.Group("/dashboard")
	{
		/************************************** GETメソッド **************************************/

		// ダッシュボードのデフォルト表示に必要な情報を取得（当月が含まれる売上期間の受注ベースの個人売上情報を取得する）
		dashboardAPI.GET("/default/:agent_staff_id", routes.GetDashboardForDefaultPreview(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// デフォルト表示のヨミ情報（ページングあり）
		dashboardAPI.GET("/default/sale_list/:agent_staff_id", routes.GetSaleListForDefaultPreview(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// デフォルト表示のリリース求職者情報（ページングあり）
		dashboardAPI.GET("/default/release_list/:agent_staff_id", routes.GetReleaseJobSeekerListForDefaultPreview(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// ダッシュボードの検索時に必要な情報
		dashboardAPI.GET("/search", routes.GetDashboardForSearch(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// ダッシュボードの検索時のヨミ情報（ページングあり）
		dashboardAPI.GET("/search/sale_list", routes.GetSaleListForSearch(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// ダッシュボードの検索時のリリース求職者情報（ページングあり）
		dashboardAPI.GET("/search/release_list", routes.GetReleaseJobSeekerListForSearch(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// CSVファイルを出力 (確定・ヨミ)
		dashboardAPI.GET("/export_csv/accuracy/:agent_id", routes.ExportAccuracyCSV(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// スカウトサービス API
	//
	scoutServiceAPI := noAuthAPI.Group("/scout_service")

	{
		scoutServiceHandler := di.InitializeScoutServiceHandler(firebase, db, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.App, r.cfg.GoogleAPI, r.cfg.Slack)
		/************************************** POSTメソッド **************************************/
		// スカウトサービス作成
		scoutServiceAPI.POST("/create", routes.CreateScoutService(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.App, r.cfg.GoogleAPI, r.cfg.Slack))

		/************************************** PUTメソッド **************************************/
		// スカウトサービス更新
		scoutServiceAPI.PUT("/update/:scout_service_id", routes.UpdateScoutService(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.App, r.cfg.GoogleAPI, r.cfg.Slack))

		// スカウトサービスのパスワード更新
		scoutServiceAPI.PUT("/update/password", routes.UpdateScoutServicePassword(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.App, r.cfg.GoogleAPI, r.cfg.Slack))

		/************************************** GETメソッド **************************************/
		// IDからスカウトサービスを取得
		scoutServiceAPI.GET("/:scout_service_id", routes.GetByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.App, r.cfg.GoogleAPI, r.cfg.Slack))

		// エージェントIDからスカウトサービスを取得
		scoutServiceAPI.GET("/list/agent/:agent_id", scoutServiceHandler.GetListByAgentID())
	}

	/****************************************************************************************/
	/// Agent API
	//
	saleManagementAPI := noAuthAPI.Group("/sale_management")
	{
		/************************************** GETメソッド **************************************/
		// 売上期間の一覧(決算月単位で取得)
		saleManagementAPI.GET("/list/agent/:agent_id", routes.GetSaleManagementListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 売上期間の一覧＋担当者情報(決算月単位で取得)
		saleManagementAPI.GET("/list/staff_management/agent/:agent_id", routes.GetSaleManagementListAndStaffManagementByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 売上情報のみを取得
		saleManagementAPI.GET("/:management_id", routes.GetSaleManagementByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 売上情報(＋年間全体売上)を取得 年間売上の編集時のデフォルトバリューに使用
		saleManagementAPI.GET("/agent_monthly_list/:management_id", routes.GetSaleManagementAndAgentMonthlyByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 売上情報(＋年間個人売上)を取得 年間売上の編集時のデフォルトバリューに使用
		saleManagementAPI.GET("/staff_monthly_list/:management_id/staff/:agent_staff_id", routes.GetStaffSaleManagementAndAgentMonthlyByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 個人売上の合算値を全体売り上げとして取得する
		saleManagementAPI.GET("/sum_of/staff_monthly/:management_id", routes.GetSumOfStaffMonthlyByManagementID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	// applicationAPI := noAuthAPI.Group("/application")
	// {
	// 	// サービス申し込みのデータを作成
	// 	applicationAPI.POST("/create", routes.CreateServiceApplication(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

	// 	// サービス利用停止に関するフィードバックのデータを作成
	// 	applicationAPI.POST("/cancellation/create", routes.CreateCancellationFeedback(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	// }

	/****************************************************************************************/
	/// Agent API
	//
	agentStaffMonthlySaleAPI := noAuthAPI.Group("/staff_monthly_sale")
	{
		/************************************** POSTメソッド **************************************/
		// エージェントの担当者アカウントの作成
		agentStaffMonthlySaleAPI.POST("/create", routes.CreateAgentStaffMonthlySale(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** PUTメソッド **************************************/
		// エージェント情報を更新
		agentStaffMonthlySaleAPI.PUT("/update", routes.UpdateAgentStaffMonthlySale(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/
		// 担当者IDとマネジメントIDに合致する情報を取得
		agentStaffMonthlySaleAPI.GET("/list/staff/:agent_staff_id/management/:management_id", routes.GetStaffMonthlySaleList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// Agent API
	//
	agentMonthlySaleAPI := noAuthAPI.Group("/agent_monthly_sale")
	{
		/************************************** POSTメソッド **************************************/
		// エージェントの担当者アカウントの作成
		agentMonthlySaleAPI.POST("/create", routes.CreateAgentMonthlySale(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** PUTメソッド **************************************/
		// エージェント情報を更新
		agentMonthlySaleAPI.PUT("/update", routes.UpdateAgentMonthlySale(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/
		// エージェント売り上げ目標取得
		agentMonthlySaleAPI.GET("/list/agent/:agent_id/management/:management_id", routes.GetAgentMonthlySaleList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェント売り上げ目標取得（デフォルト表示）
		agentMonthlySaleAPI.GET("/default/list/agent/:agent_id", routes.GetDefaultPreviewAgentMonthlySaleListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// 企業求人一括インポート情報 API
	//
	initialEnterpriseImporterAPI := noAuthAPI.Group("/initial_enterprise_importer")
	{
		// 企業求人一括インポート情報の登録
		initialEnterpriseImporterAPI.POST("/create", routes.CreateInitialEnterpriseImporter(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 企業求人一括インポート情報の削除
		initialEnterpriseImporterAPI.DELETE("/delete/:initial_enterprise_importer_id", routes.DeleteInitialEnterpriseImporter(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェントIDから企業求人一括インポート情報一覧を取得
		initialEnterpriseImporterAPI.GET("/list/agent/:agent_id", routes.GetInitialEnterpriseImporterListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 今後1週間の企業求人一括インポート情報一覧を取得 *予約日程調整時の重複チェック用
		initialEnterpriseImporterAPI.GET("/list/week", routes.GetInitialEnterpriseImporterListByWeek(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	// / 送客先エージェント API
	sendingEnterpriseAPI := noAuthAPI.Group("/sending_enterprise")
	{
		/************************************** POSTメソッド **************************************/

		// 送客先エージェント情報の登録
		sendingEnterpriseAPI.POST("/create", routes.CreateSendingEnterprise(db, firebase, r.cfg.Sendgrid))

		// 送客先エージェント資料の登録
		sendingEnterpriseAPI.POST("/reference_material", routes.CreateSendingEnterpriseReferenceMaterial(db, firebase, r.cfg.Sendgrid))

		// 送客メールの送信
		sendingEnterpriseAPI.POST("/send_mail", routes.SendSendingMail(db, firebase, r.cfg.Sendgrid))

		// 参加確認のメール送信
		sendingEnterpriseAPI.POST("/send_mail/rsvp", routes.SendMailForRSVP(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/

		// 送客先エージェント情報の更新
		sendingEnterpriseAPI.PUT("/update/:sending_enterprise_id", routes.UpdateSendingEnterprise(db, firebase, r.cfg.Sendgrid))

		// 送客先エージェント資料の更新
		sendingEnterpriseAPI.PUT("/reference_material/:sending_enterprise_id", routes.UpdateSendingEnterpriseReferenceMaterial(db, firebase, r.cfg.Sendgrid))

		// 送客先のパスワードを変更
		sendingEnterpriseAPI.PUT("/update/password/:sending_enterprise_id", routes.UpdateSendingEnterprisePassword(db, firebase, r.cfg.Sendgrid))

		// 送客先のログイン
		sendingEnterpriseAPI.PUT("/signin", routes.SigninSendingEnterprise(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/

		// 送客先エージェントIDを使って送客先エージェント情報取得する関数
		sendingEnterpriseAPI.GET("/:sending_enterprise_id", routes.GetSendingEnterpriseByID(db, firebase, r.cfg.Sendgrid))

		// 送客先エージェントIDを使って送客先エージェント + 請求先情報取得する関数
		sendingEnterpriseAPI.GET("/billing_address/:sending_enterprise_id", routes.GetSendingEnterpriseAndBillingAddressByID(db, firebase, r.cfg.Sendgrid))

		// 送客先エージェントUUIDを使って送客先エージェント情報取得する関数
		sendingEnterpriseAPI.GET("/uuid/:sending_enterprise_uuid", routes.GetSendingEnterpriseByUUID(db, firebase, r.cfg.Sendgrid))

		// 送客先エージェントIDを使って送客先エージェント資料取得する関数
		sendingEnterpriseAPI.GET("/reference_material/:sending_enterprise_id", routes.GetSendingEnterpriseReferenceMaterialBySendingEnterpriseID(db, firebase, r.cfg.Sendgrid))

		// 指定ページの送客先エージェント情報を取得 query: page_number
		sendingEnterpriseAPI.GET("/all/page/free_word", routes.GetAllSendingEnterpriseByPageAndFreeWord(db, firebase, r.cfg.Sendgrid))

		// 送客メール送信時に必要な情報をまとめて取得するapi
		sendingEnterpriseAPI.GET("/for_sending_mail", routes.GetSendingInformationForSendingMail(db, firebase, r.cfg.Sendgrid))

		// 参加確認(RSVP)に必要な企業情報 + 求職者情報を取得するapi
		sendingEnterpriseAPI.GET("/rsvp/agent/:agent_id", routes.GetSendingEnterpriseAndAcceptJobSeekerListByAgentID(db, firebase, r.cfg.Sendgrid))

		/************************************** DELETEメソッド **************************************/

		// 送客先エージェント情報の削除
		sendingEnterpriseAPI.DELETE("/delete", routes.DeleteSendingEnterprise(db, firebase, r.cfg.Sendgrid))

		// 送客先エージェントの参考資料の削除 sending_enterprise_id, file_type: '送客先エージェント画像' | '参考資料1' | '参考資料2'
		sendingEnterpriseAPI.DELETE("/:sending_enterprise_id/material/:file_type/delete", routes.DeleteSendingEnterpriseReferenceMaterial(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/

	// / 送客先エージェントの請求先 API
	sendingSendingBillingAddressAPI := noAuthAPI.Group("/sending_billing_address")
	{
		/************************************** PUTメソッド **************************************/

		// 送客先エージェントの請求先情報の更新
		sendingSendingBillingAddressAPI.PUT("/update/:sending_billing_address_id", routes.UpdateSendingBillingAddress(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/

		// IDを使って送客先エージェントの請求先情報取得する関数
		sendingSendingBillingAddressAPI.GET("/:sending_billing_address_id", routes.GetSendingBillingAddressByID(db, firebase, r.cfg.Sendgrid))

		// 送客先エージェントIDを使って送客先エージェントの請求先情報一覧取得する関数
		sendingSendingBillingAddressAPI.GET("/enterprise/:sending_enterprise_id", routes.GetSendingBillingAddressBySendingEnterpriseID(db, firebase, r.cfg.Sendgrid))

		// 送客先エージェントIDを使って送客先エージェントの請求先情報一覧取得する関数
		sendingSendingBillingAddressAPI.GET("/list/enterprise/id_list", routes.GetSendingBillingAddressListBySendingEnterpriseIDList(db, firebase, r.cfg.Sendgrid))

		// 全ての送客先エージェントの請求先情報を取得
		sendingSendingBillingAddressAPI.GET("/all", routes.GetAllSendingBillingAddress(db, firebase, r.cfg.Sendgrid))

	}

	/****************************************************************************************/
	/// 求人 API
	//
	sendingSendingJobInformationAPI := noAuthAPI.Group("/sending_job_information")
	{
		/************************************** POSTメソッド **************************************/

		// 送客先エージェントの求人情報の登録
		sendingSendingJobInformationAPI.POST("/create", routes.CreateSendingJobInformation(db, firebase, r.cfg.Sendgrid))

		//csvファイルを使って送客先エージェントの求人情報の登録
		sendingSendingJobInformationAPI.POST("/import_csv/:sending_billing_address_id", routes.ImportSendingJobInformationCSV(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/

		// 送客先エージェントの求人情報の更新
		sendingSendingJobInformationAPI.PUT("/update/:sending_job_information_id", routes.UpdateSendingJobInformation(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/

		// 求人IDを使って送客先エージェントの求人情報取得する関数
		sendingSendingJobInformationAPI.GET("/:sending_job_information_id", routes.GetSendingJobInformationByID(db, firebase, r.cfg.Sendgrid))

		// 求人のuuidを使って送客先エージェントの求人情報取得する関数
		sendingSendingJobInformationAPI.GET("/uuid/:sending_job_information_uuid", routes.GetSendingJobInformationByUUID(db, firebase, r.cfg.Sendgrid))

		// 求人のuuidを使って送客先エージェントの求人情報取得する関数
		sendingSendingJobInformationAPI.GET("/job_listing/uuid/:sending_job_information_uuid", routes.GetJobListingBySendingJobInformationUUID(db, firebase, r.cfg.Sendgrid))

		// 送客先エージェントの請求先IDを使って送客先エージェントの求人情報一覧取得する関数
		sendingSendingJobInformationAPI.GET("/list/sending_enterprise/:sending_enterprise_id", routes.GetSendingJobInformationListBySendingEnterpriseID(db, firebase, r.cfg.Sendgrid))

		// まだ送客していない送客先のリストと指定IDの送客先が保有する求人の一覧と最大ページ数を取得するapi
		sendingSendingJobInformationAPI.GET("/list/page/have_not/sent/yet", routes.GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYet(db, firebase, r.cfg.Sendgrid))

		// まだ送客していない送客先のリストと指定IDの送客先が保有する求人の一覧と最大ページ数を取得するapi
		sendingSendingJobInformationAPI.GET("/list/page/search/have_not/sent/yet", routes.GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYet(db, firebase, r.cfg.Sendgrid))

		/************************************** DELETEメソッド **************************************/

		// 送客先エージェントの求人情報の削除
		sendingSendingJobInformationAPI.DELETE("/delete/:sending_job_information_id", routes.DeleteSendingJobInformation(db, firebase, r.cfg.Sendgrid))

	}

	/****************************************************************************************/
	// 送客求職者(送客元表示)のapi
	//
	sendingCustomerAPI := noAuthAPI.Group("/sending_customer")
	{
		/************************************** GETメソッド **************************************/
		// 指定IDの送客相談情報を取得する
		sendingCustomerAPI.GET("/:sending_customer_id", routes.GetSendingCustomerByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェントIDの送客求職者リストを取得
		sendingCustomerAPI.GET("/list/agent/:agent_id", routes.GetSendingCustomerListByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 指定ページ, エージェントID, フェーズの送客求職者リストを取得 query: page_number, agent_id, panel_number
		sendingCustomerAPI.GET("/search/list/page/tab", routes.GetSearchSendingCustomerListByPageAndTabAndAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	// 送客求職者(送客先表示)のapi
	//
	sendingJobSeekerAPI := noAuthAPI.Group("/sending_job_seeker")
	{
		/************************************** POSTメソッド **************************************/
		// 最初の作成 *作成時は名前と連絡先のみで、それ以外の項目は作成後最初の更新時に作成する. sending_customerも同時に作成される
		sendingJobSeekerAPI.POST("/create", routes.CreateSendingJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// CRM求職者から送客求職者を作成
		sendingJobSeekerAPI.POST("/create_from_job_seeker", routes.CreateSendingJobSeekerFromJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 面談前アンケートの登録 (求職者の同意項目、業界、職種、勤務地、質問要望、ファイル（履歴書（原本）、職務経歴書（原本）））
		sendingJobSeekerAPI.POST("/create/initial_questionnaire", routes.CreateSendingInitialQuestionnaire(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 送客終了理由
		sendingJobSeekerAPI.POST("/create/end_status", routes.CreateSendingJobSeekerEndStatus(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** PUTメソッド **************************************/
		// 作成後最初の更新 *作成時は名前と連絡先のみで、それ以外の項目は作成後最初の更新時に作成する. sending_customerも同時に更新される
		sendingJobSeekerAPI.PUT("/first_update/:sending_customer_id", routes.FirstUpdateSendingJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 更新
		sendingJobSeekerAPI.PUT("/update/:sending_job_seeker_id", routes.UpdateSendingJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 書類の更新
		sendingJobSeekerAPI.PUT("/update/document", routes.UpdateSendingJobSeekerDocument(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// フェーズの更新
		sendingJobSeekerAPI.PUT("/update/phase", routes.UpdateSendingJobSeekerPhase(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 面談日時の更新
		sendingJobSeekerAPI.PUT("/update/interview_date", routes.UpdateSendingInterviewDateBySendingJobSeekerID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// アクティビティメモの更新
		sendingJobSeekerAPI.PUT("/update/activity_memo", routes.UpdateSendingJobSeekerActivityMemo(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 面談実施待ちの未閲覧の更新（送客進捗管理上で未閲覧クリックした時に実行）
		sendingJobSeekerAPI.PUT("/update/is_view/for_waiting/:sending_job_seeker_id", routes.UpdateIsVewForWating(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 未登録の未閲覧の更新（送客進捗管理上で未閲覧クリックした時に実行）
		sendingJobSeekerAPI.PUT("/update/is_view/for_unregister/:sending_job_seeker_id", routes.UpdateIsVewForUnregister(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/
		// 指定IDの送客相談情報を取得する
		sendingJobSeekerAPI.GET("/:sending_job_seeker_id", routes.GetSendingJobSeekerByID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 指定IDの送客相談情報を取得する
		sendingJobSeekerAPI.GET("/uuid/:sending_job_seeker_uuid", routes.GetSendingJobSeekerByUUID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 指定IDの書類を取得する
		sendingJobSeekerAPI.GET("/document/:sending_job_seeker_id", routes.GetSendingJobSeekerDocumentBySendingJobSeekerID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 未閲覧の送客求職者の数を取得するapi
		sendingJobSeekerAPI.GET("/is_not_view/count/:agent_staff_id", routes.GetIsNotViewSendingJobSeekerCountByAgentStaffID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 送客側の進捗一覧画面の絞り込み検索に使用する担当者, 送客元, 送客先のリストを取得する
		sendingJobSeekerAPI.GET("/list_used_for_search/agent/:agent_id", routes.GetSearchListForSendingJobSeekerManagementByAgentID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** DELETEメソッド **************************************/

		// 送客求職者情報の削除
		sendingJobSeekerAPI.DELETE("/delete/:sending_job_seeker_id", routes.DeleteSendingJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 履歴書PDFの削除（resume_pdf_urlカラムを空文字で更新）
		sendingJobSeekerAPI.DELETE("/:sending_job_seeker_id/document/resume_pdf_url/delete", routes.DeleteSendingJobSeekerResumePDFURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 履歴書原本の削除（resume_origin_urlカラムを空文字で更新）
		sendingJobSeekerAPI.DELETE("/:sending_job_seeker_id/document/resume_origin_url/delete", routes.DeleteSendingJobSeekerResumeOriginURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 職務経歴書PDFの削除（cv_pdf_urlカラムを空文字で更新）
		sendingJobSeekerAPI.DELETE("/:sending_job_seeker_id/document/cv_pdf_url/delete", routes.DeleteSendingJobSeekerCVPDFURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 職務経歴書原本の削除（cv_origin_urlカラムを空文字で更新）
		sendingJobSeekerAPI.DELETE("/:sending_job_seeker_id/document/cv_origin_url/delete", routes.DeleteSendingJobSeekerCVOriginURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 推薦状PDFの削除（recommendation_pdf_urlカラムを空文字で更新）
		sendingJobSeekerAPI.DELETE("/:sending_job_seeker_id/document/recommendation_pdf_url/delete", routes.DeleteSendingJobSeekerRecommendationPDFURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 推薦状原本の削除（recommendation_origin_urlカラムを空文字で更新）
		sendingJobSeekerAPI.DELETE("/:sending_job_seeker_id/document/recommendation_origin_url/delete", routes.DeleteSendingJobSeekerRecommendationOriginURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 証明写真の削除（id_photo_urlカラムを空文字で更新）
		sendingJobSeekerAPI.DELETE("/:sending_job_seeker_id/document/id_photo_url/delete", routes.DeleteSendingJobSeekerIDPhotoURL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// その他①の削除（other_document1_urlカラムを空文字で更新）
		sendingJobSeekerAPI.DELETE("/:sending_job_seeker_id/document/other_document1_url/delete", routes.DeleteSendingJobSeekerOtherDocument1URL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// その他②の削除（other_document2_urlカラムを空文字で更新）
		sendingJobSeekerAPI.DELETE("/:sending_job_seeker_id/document/other_document2_url/delete", routes.DeleteSendingJobSeekerOtherDocument2URL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// その他③の削除（other_document3_urlカラムを空文字で更新）
		sendingJobSeekerAPI.DELETE("/:sending_job_seeker_id/document/other_document3_url/delete", routes.DeleteSendingJobSeekerOtherDocument3URL(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// エージェントと送客求職者のチャットグループ API
	//
	chatGroupWithSendingJobSeekerAPI := noAuthAPI.Group("/chat_group_with_sending_job_seeker")
	{
		/************************************** POST	メソッド **************************************/

		// エージェントと送客求職者のチャットグループの登録
		chatGroupWithSendingJobSeekerAPI.POST("/create", routes.CreateChatGroupWithSendingJobSeeker(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/

		// エージェントの最終閲覧時間を更新
		chatGroupWithSendingJobSeekerAPI.PUT("/update/agent/last_watched/:group_id", routes.UpdateSendingJobSeekerAgentLastWatched(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/

		// IDからエージェントと送客求職者のチャットグループを取得
		chatGroupWithSendingJobSeekerAPI.GET("/:id", routes.GetChatGroupWithSendingJobSeekerByID(db, firebase, r.cfg.Sendgrid))

		// Freewordで絞り込んで送客求職者のチャットグループを取得
		chatGroupWithSendingJobSeekerAPI.GET("/search/list/page/agent/:agent_id", routes.GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentID(db, firebase, r.cfg.Sendgrid))

		// チャットの通知件数を返す
		chatGroupWithSendingJobSeekerAPI.GET("/notification/agent/:agent_id", routes.GetSendingJobSeekerChatNotificationByAgentID(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	/// エージェントと送客求職者のチャットメッセージAPI
	//
	chatMessageWithSendingJobSeekerAPI := noAuthAPI.Group("/chat_message_with_sending_job_seeker")
	{
		/************************************** POSTメソッド **************************************/

		// エージェントが送客求職者にLINEでメッセージ送信
		chatMessageWithSendingJobSeekerAPI.POST("/line/message/create", routes.SendChatMessageWithSendingJobSeekerLineMessage(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// エージェントが送客求職者にLINEでメッセージ送信
		chatMessageWithSendingJobSeekerAPI.POST("/line/image/create", routes.SendChatMessageWithSendingJobSeekerLineImage(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/

		// グループIDからエージェントと送客求職者のチャットメッセージを取得
		chatMessageWithSendingJobSeekerAPI.GET("/list/group/:group_id", routes.GetChatMessageWithSendingJobSeekerListByGroupID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// エージェントと送客求職者のメールAPI
	//
	emailWithSendingJobSeekerAPI := noAuthAPI.Group("/email_with_sending_job_seeker")
	{
		/************************************** POSTメソッド **************************************/

		// エージェントが送客求職者にメール送信
		emailWithSendingJobSeekerAPI.POST("/create", routes.SendEmailWithSendingJobSeeker(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/

		// 送客求職者IDからエージェントと送客求職者のメール一覧を取得
		emailWithSendingJobSeekerAPI.GET("/list/sending_job_seeker/:sending_job_seeker_id", routes.GetEmailWithSendingJobSeekerListBySendingJobSeekerID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	// 送客進捗のapi
	//
	sendingPhaseAPI := noAuthAPI.Group("/sending_phase")
	{
		/************************************** POSTメソッド **************************************/
		// 送客応諾後の送客終了理由
		sendingPhaseAPI.POST("/create/end_status", routes.CreateSendingPhaseEndStatus(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/

		// 送客進捗の送客予定日更新
		sendingPhaseAPI.PUT("/update/sending_date", routes.UpdateSendingPhaseSendingDate(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/
		// 指定IDの送客相談情報を取得する
		sendingPhaseAPI.GET("/:sending_phase_id", routes.GetSendingPhaseByID(db, firebase, r.cfg.Sendgrid))

		// 送客の進捗一覧ページに表示するリストを取得するapi（query: page_number, tab_number, staff_id_list, sender_id_list, send_agent_id_list）
		sendingPhaseAPI.GET("/search/list/page/tab/agent/:agent_id", routes.GetSearchSendingPhaseListByPageAndTabAndAgentID(db, firebase, r.cfg.Sendgrid))

		// 求職者IDの送客進捗リストを取得
		sendingPhaseAPI.GET("/list/job_seeker/:sending_job_seeker_id", routes.GetSendingPhaseListBySendingJobSeekerID(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	// 送客求職者の興味あり求人のapi
	//
	sendingJobSeekerDesiredJobInformationAPI := noAuthAPI.Group("/sending_job_seeker_desired_job_information")
	{
		/************************************** POSTメソッド **************************************/
		// 興味あり求人を複数作成
		sendingJobSeekerDesiredJobInformationAPI.POST("/create/multi", routes.CreateMultiSendingJobSeekerDesiredJobInformation(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/
		// 指定求職者IDの興味あり求人を取得する
		sendingJobSeekerDesiredJobInformationAPI.GET("/list/sending_job_seeker/:sending_job_seeker_id", routes.GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerID(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	// 送客売上のapi
	//
	sendingSaleAPI := noAuthAPI.Group("/sending_sale")
	{
		/************************************** POSTメソッド **************************************/
		// 送客売上の作成
		sendingSaleAPI.POST("/create", routes.CreateSendingSale(db, firebase, r.cfg.Sendgrid))

		/************************************** PUTメソッド **************************************/
		// 送客売上の更新
		sendingSaleAPI.PUT("/update/:sending_sale_id", routes.UpdateSendingSale(db, firebase, r.cfg.Sendgrid))

		/************************************** GETメソッド **************************************/
		// 指定IDの送客売上を取得する
		sendingSaleAPI.GET("/:sending_sale_id", routes.GetSendingSaleByID(db, firebase, r.cfg.Sendgrid))

		// 指定求職者と送客先の送客売上を取得する
		sendingSaleAPI.GET("/job_seeker_and_enterprise/:job_seeker_id/:enterprise_id", routes.GetSendingSaleByJobSeekerIDAndEnterpriseID(db, firebase, r.cfg.Sendgrid))

		// 送客元のエージェントの指定した月範囲とエージェントの送客売上を取得 query: start_month, end_month
		sendingSaleAPI.GET("/list/sender_agent/monthly/:sender_agent_id", routes.GetSendingSaleListBySenderAgentIDForMonthly(db, firebase, r.cfg.Sendgrid))

		// 送客管理のエージェントの指定した月範囲とエージェントの送客売上を取得 query: start_month, end_month
		sendingSaleAPI.GET("/list/agent/monthly/:agent_id", routes.GetSendingSaleListByAgentIDForMonthly(db, firebase, r.cfg.Sendgrid))

		// 指定した月範囲のすべての送客売上を取得 query: start_month, end_month
		sendingSaleAPI.GET("/all/monthly", routes.GetAllSendingSaleForMonthly(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	// お知らせ関連API
	//
	notificationForUserAPI := noAuthAPI.Group("/notification_for_user")
	{
		/************************************** POSTメソッド **************************************/
		// お知らせを作成
		notificationForUserAPI.POST("/create", routes.CreateNotificationForUser(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// お知らせを確認したユーザーを作成
		notificationForUserAPI.POST("/create/view/:agent_staff_id", routes.CreateUserNotificationView(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** PUTメソッド **************************************/

		/************************************** GETメソッド **************************************/
		// ページと送信対象からお知らせを取得 @param: page_number, target_list
		notificationForUserAPI.GET("/list/page/target", routes.GetNotificationForUserListByPageAndTargetList(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// 担当者IDからお知らせの未読件数を取得
		notificationForUserAPI.GET("/unwatched/count/:agent_staff_id", routes.GetUnwatchedNotificationForUserCountByAgentStaffID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	// デプロイ反映関連API
	//
	deploymentReflectionAPI := noAuthAPI.Group("/deployment_reflection")
	{
		/************************************** POSTメソッド **************************************/
		// デプロイ反映の情報を一括作成
		deploymentReflectionAPI.POST("/create", routes.CreateDeploymentInformation(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** PUTメソッド **************************************/
		// デプロイ反映のステータスを更新
		deploymentReflectionAPI.PUT("/update/:deployment_id", routes.UpdateDeployMenConfirmStatusByStaffID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// デプロイ反映のステータスを更新
		deploymentReflectionAPI.PUT("/update/information/:deployment_id", routes.UpdateDeploymentInformation(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// デプロイ反映のステータスを更新
		deploymentReflectionAPI.PUT("/update/is_reflected/:agent_staff_id", routes.UpdateDeployMenConfirmStatusByStaffID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		/************************************** GETメソッド **************************************/
		// 全てのデプロイ情報
		deploymentReflectionAPI.GET("/all", routes.GetAllDeploymentInformations(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))

		// デプロイ情報が反映済みかを確認 (全て反映ずみの場合はTrueを返す)
		deploymentReflectionAPI.GET("/reflected/:agent_staff_id", routes.CheckDeploymentReflection(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal))
	}

	/****************************************************************************************/
	/// Admin API
	//
	adminAPI := r.Engine.Group("/admin")
	adminNoAuthAPI := adminAPI.Group("")
	{
		// 管理者サイトのログイン
		adminNoAuthAPI.PUT("/authorize", routes.AdminAuthorize(db, r.cfg.App))
	}

	/****************************************************************************************/

	/****************************************************************************************/
	/// LP API
	//
	lpAPI := noAuthAPI.Group("/lp")
	{
		/************************************** POSTメソッド **************************************/
		// LPのログインフォームのログイン処理
		lpAPI.POST("/job_seeker/login", routes.LoginGuestJobSeekerForLP(db, firebase, r.cfg.Sendgrid))

		// パスワード再設定メールを送信する（LP）
		lpAPI.POST("/job_seeker/send/reset_password_email", routes.SendJobSeekerResetPasswordEmailForLP(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// LPから求職者登録
		lpAPI.POST("/job_seeker/create", routes.CreateJobSeekerFromLP(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// LPの診断ページから合致求人数と年収期待値を取得
		lpAPI.POST("/job_information/search/diagnosis", routes.GetSearchJobInformationCountByLPDiagnosis(db, firebase, r.cfg.Sendgrid))

		// LPからのお問い合わせ
		lpAPI.POST("/contact", routes.SendLPContact(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		/************************************** PUTメソッド **************************************/
		// LPから求職者の電話番号を更新
		lpAPI.PUT("/job_seeker/update/phone_number", routes.UpdateJobSeekerPhoneFromLP(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// LPから求職者の希望条件を更新
		lpAPI.PUT("/job_seeker/update/desired", routes.UpdateJobSeekerDesiredFromLP(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// LPから求職者の電話番号と希望条件を更新
		lpAPI.PUT("/job_seeker/reset/password", routes.ResetPasswordForLP(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		/************************************** GETメソッド **************************************/
		// LPの登録状況を取得する
		lpAPI.GET("/job_seeker/register_status/uuid/:job_seeker_uuid", routes.GetJobSeekerLPRegisterStatusByUUID(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// LPからパスワードトークンの有効性を確認
		lpAPI.GET("/job_seeker/check/password_token/:password_token", routes.CheckResetPasswordToken(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.Slack))

		// LPから診断に使用する求人を取得
		lpAPI.GET("/job_information/list/for_diagnosis", routes.GetJobInformationListForDiagnosis(db, firebase, r.cfg.Sendgrid))
	}

	/****************************************************************************************/
	// Admin API
	//
	rpaAPI := noAuthAPI.Group("/rpa")
	{
		// Gmail Webhook
		rpaAPI.POST(
			"/gmail_hook",
			routes.GmailWebHook(db, firebase, r.cfg.Sendgrid, r.cfg.OneSignal, r.cfg.App, r.cfg.GoogleAPI, r.cfg.Slack),
			// リクエストを制限(https://echo.labstack.com/docs/middleware/rate-limiter)
			// middleware.RateLimiterWithConfig(rateLimiterConfig),
		)
	}

	/****************************************************************************************/

	return r
}

func (r *Router) Start() {
	r.Engine.Start(fmt.Sprintf(":%d", 8080))
}
