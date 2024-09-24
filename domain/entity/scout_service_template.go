package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type ScoutServiceTemplate struct {
	ID                  uint      `db:"id" json:"id"`
	ScoutServiceID      uint      `db:"scout_service_id" json:"scout_service_id"`           // スカウトサービスID
	StartHour           null.Int  `db:"start_hour" json:"start_hour"`                       // スカウト開始時間(媒体共通/0:0時, 1:1時, 2:2時, ..., 23:23時)
	StartMinute         null.Int  `db:"start_minute" json:"start_minute"`                   // スカウト開始分(媒体共通/0:0分, 1:1分, 2:2分, ..., 59:59分)
	RunOnMonday         bool      `db:"run_on_monday" json:"run_on_monday"`                 // 月曜日に走らせるかどうか/false:走らせない true:走る(共通)
	RunOnTuesday        bool      `db:"run_on_tuesday" json:"run_on_tuesday"`               // 火曜日に走らせるかどうか/false:走らせない true:走る(共通)
	RunOnWednesday      bool      `db:"run_on_wednesday" json:"run_on_wednesday"`           // 水曜日に走らせるかどうか/false:走らせない true:走る(共通)
	RunOnThursday       bool      `db:"run_on_thursday" json:"run_on_thursday"`             // 木曜日に走らせるかどうか/false:走らせない true:走る(共通)
	RunOnFriday         bool      `db:"run_on_friday" json:"run_on_friday"`                 // 金曜日に走らせるかどうか/false:走らせない true:走る(共通)
	RunOnSaturday       bool      `db:"run_on_saturday" json:"run_on_saturday"`             // 土曜日に走らせるかどうか/false:走らせない true:走る(共通)
	RunOnSunday         bool      `db:"run_on_sunday" json:"run_on_sunday"`                 // 日曜日に走らせるかどうか/false:走らせない true:走る(共通)
	ScoutCount          null.Int  `db:"scout_count" json:"scout_count"`                     // スカウト件数(媒体共通)
	SearchTitle         string    `db:"search_title" json:"search_title"`                   // スカウトに利用する検索条件のタイトル(媒体共通)
	MessageTitle        string    `db:"message_title" json:"message_title"`                 // スカウトに利用するメッセージのタイトル(媒体共通)
	JobInformationTitle string    `db:"job_information_title" json:"job_information_title"` // スカウトに利用する求人情報のタイトル
	JobInformationID    string    `db:"job_information_id" json:"job_information_id"`       // スカウトに利用する求人情報ID(RAN)
	AgeLimit            null.Int  `db:"age_limit" json:"age_limit"`                         // 年齢制限(Mynavi)
	ScoutType           null.Int  `db:"scout_type" json:"scout_type"`                       // スカウトタイプ(AMBI/0: 通常スカウト, 1: プレミアムスカウト)
	AutoRemind          bool      `db:"auto_remind" json:"auto_remind"`                     // 自動リマインド(AMBI/0:しない, 1:する)
	ReplyLimit          null.Int  `db:"reply_limit" json:"reply_limit"`                     // 返信制限(AMBI, マイナビ)
	LastSendCount       null.Int  `db:"last_send_count" json:"last_send_count"`             // 最終送信件数
	LastSendAt          time.Time `db:"last_send_at" json:"last_send_at"`                   // 最終送信日時
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`

	// 結合
	// ScoutServiceテーブル
	LoginID     string   `db:"login_id" json:"login_id"`
	Password    string   `db:"password" json:"password"`
	ServiceType null.Int `db:"service_type" json:"service_type"` // サービスタイプ(0: RAN, 1: マイナビ転職スカウト, 2:AMBI)
}

func NewScoutServiceTemplate(
	scounServiceID uint,
	startHour null.Int,
	startMinute null.Int,
	runOnMonday bool,
	runOnTuesday bool,
	runOnWednesday bool,
	runOnThursday bool,
	runOnFriday bool,
	runOnSaturday bool,
	runOnSunday bool,
	scoutCount null.Int,
	searchTitle string,
	messageTitle string,
	jobInformationTitle string,
	jobInformationID string,
	ageLimit null.Int,
	scoutType null.Int,
	autoRemind bool,
	replyLimit null.Int,
) *ScoutServiceTemplate {
	return &ScoutServiceTemplate{
		ScoutServiceID:      scounServiceID,
		StartHour:           startHour,
		StartMinute:         startMinute,
		RunOnMonday:         runOnMonday,
		RunOnTuesday:        runOnTuesday,
		RunOnWednesday:      runOnWednesday,
		RunOnThursday:       runOnThursday,
		RunOnFriday:         runOnFriday,
		RunOnSaturday:       runOnSaturday,
		RunOnSunday:         runOnSunday,
		ScoutCount:          scoutCount,
		SearchTitle:         searchTitle,
		MessageTitle:        messageTitle,
		JobInformationTitle: jobInformationTitle,
		JobInformationID:    jobInformationID,
		AgeLimit:            ageLimit,
		ScoutType:           scoutType,
		AutoRemind:          autoRemind,
		ReplyLimit:          replyLimit,
	}
}

// スカウトテンプレートのスカウトタイプ
const (
	// 通常スカウト
	AmbiScoutTypeNormal = iota
	// 通常スカウト(再スカウト)
	AmbiScoutTypeNormalAndAgain
	// プレミアムスカウト
	AmbiScoutTypePremium
	// プレミアムスカウト(再スカウト)
	AmbiScoutTypePremiumAndAgain
)

const (
	// 通常スカウト
	RanScoutTypeNormal = iota
	// 再スカウト
	RanScoutTypeAgain
	// 別求人を送る
	RanScoutTypeSendOther
)

/*
  テーブル構成
  エージェントが利用するロボットテーブル(agent_robots) *管理側のみ作成可能。オプション契約後管理側が作成。 *この数の分コンテナができる *agent_idに紐づく
  ↓
  スカウトサービステーブル([]scout_services) *サービスのID,パスワードなどの情報を持つ
  ↓
  スカウトテンプレートテーブル([]scout_service_templates) *スカウト送信の時間帯や、スカウトの内容を持つ
  エントリー取得時間テーブル([]scout_service_get_entry_times) *エントリー取得の開始時間を持つ

  スカウト仕様
  参考:大平さんのスプシ（https:ocs.google.com/spreadsheets/d/1DMbOKxAJkYOfy-DiLcbp5wuI7fIe1mGh1Wh264SK9bY/edit#gid=948666478）
  ・時間間隔:
  1時間ごと(0:0時,1:1時,2:2時,3:3時,4:4時,5:5時,6:6時,7:7時,8:8時,9:9時,10:10時,11:11時,12:12時,13:13時,14:14時,15:15時,16:16時,17:17時,18:18時,19:19時,20:20時,21:21時,22:22時,23:23時)

  ・最大送信件数マスタ(1回につき):
  RAN:100, 200, 300
  マイナビ:50, 100, 300, 500
  AMBI:50, 200

  ・最大送信年齢:
  RAN:なし
  マイナビ:あり
  AMBI:なし

  ・添付求人ID:
  RAN:1つのみ
  マイナビ:なし
  AMBI:なし

  ・送信タイプ
  RAN:あり(0:通常スカウト, 1:再スカウト, 2:別求人送付)
  マイナビ:なし
  AMBI:あり(0:通常スカウト, 1:通常スカウト(再スカウト), 2:プレミアムスカウト, 3:プレミアムスカウト(再スカウト))

  返信期限(例:2日後期限なら2)
  RAN:なし
  マイナビ:あり
  AMBI:あり
*/
