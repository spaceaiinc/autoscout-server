package interactor

import (
	"fmt"
	"io"
	"os"

	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type ChatMessageWithJobSeekerInteractor interface {
	// 汎用系 API
	SendChatMessageWithJobSeekerLineMessage(input SendChatMessageWithJobSeekerLineMessageInput) (SendChatMessageWithJobSeekerLineMessageOutput, error)
	SendChatMessageWithJobSeekerLineImage(input SendChatMessageWithJobSeekerLineImageInput) (SendChatMessageWithJobSeekerLineImageOutput, error)
	GetChatMessageWithJobSeekerListByGroupID(input GetChatMessageWithJobSeekerListByGroupIDInput) (GetChatMessageWithJobSeekerListByGroupIDOutput, error)

	// LINE WebHook
	LineWebHook(input LineWebHookInput) (LineWebHookOutput, error)
}

type ChatMessageWithJobSeekerInteractorImpl struct {
	firebase                           usecase.Firebase
	sendgrid                           config.Sendgrid
	oneSignal                          config.OneSignal
	chatMessageWithJobSeekerRepository usecase.ChatMessageWithJobSeekerRepository
	chatGroupWithJobSeekerRepository   usecase.ChatGroupWithJobSeekerRepository
	jobSeekerRepository                usecase.JobSeekerRepository
	agentRepository                    usecase.AgentRepository
	agentStaffRepository               usecase.AgentStaffRepository
}

// ChatMessageWithJobSeekerInteractorImpl is an implementation of ChatMessageWithJobSeekerInteractor
func NewChatMessageWithJobSeekerInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	cmR usecase.ChatMessageWithJobSeekerRepository,
	cgR usecase.ChatGroupWithJobSeekerRepository,
	jsR usecase.JobSeekerRepository,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
) ChatMessageWithJobSeekerInteractor {
	return &ChatMessageWithJobSeekerInteractorImpl{
		firebase:                           fb,
		sendgrid:                           sg,
		oneSignal:                          os,
		chatMessageWithJobSeekerRepository: cmR,
		chatGroupWithJobSeekerRepository:   cgR,
		jobSeekerRepository:                jsR,
		agentRepository:                    aR,
		agentStaffRepository:               asR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
// エージェントが求職者にLINE送信
type SendChatMessageWithJobSeekerLineMessageInput struct {
	SendParam entity.SendChatMessageWithJobSeekerLineParam
}

type SendChatMessageWithJobSeekerLineMessageOutput struct {
	ChatMessageWithJobSeeker *entity.ChatMessageWithJobSeeker
}

func (i *ChatMessageWithJobSeekerInteractorImpl) SendChatMessageWithJobSeekerLineMessage(input SendChatMessageWithJobSeekerLineMessageInput) (SendChatMessageWithJobSeekerLineMessageOutput, error) {
	var (
		output SendChatMessageWithJobSeekerLineMessageOutput
		err    error
	)

	// GroupIDからGroupを取得
	chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByID(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送信元がエージェントかつ、テキストメッセージの場合、LINEに送信
	agent, err := i.agentRepository.FindByID(chatGroup.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// chatGroup.JobSeekerIDから求職者を取得
	jobSeeker, err := i.jobSeekerRepository.FindByID(chatGroup.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// シークレットとトークンを設定
	bot, err := linebot.New(agent.LineMessagingChannelSecret, agent.LineMessagingChannelAccessToken)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// LINE メッセージを送信する
	call, err := bot.PushMessage(
		jobSeeker.LineID,
		linebot.NewTextMessage(input.SendParam.Message),
	).Do()
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	fmt.Println("メッセージの送信に成功しました:", call)

	// エージェントから求職者へのチャットメッセージを作成
	chatMessage := entity.NewChatMessageWithJobSeeker(
		input.SendParam.GroupID,
		null.NewInt(0, true), // エージェント
		input.SendParam.Message,
		"",
		"",
		"",
		"",
		null.NewInt(0, false),
		"",
		null.NewInt(0, true),
		time.Now().In(time.UTC),
	)

	err = i.chatMessageWithJobSeekerRepository.Create(chatMessage)
	if err != nil {
		return output, err
	}

	// エージェントの最終送信時間を更新
	err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatMessageWithJobSeeker = chatMessage

	return output, nil
}

type SendChatMessageWithJobSeekerLineImageInput struct {
	SendParam entity.SendChatMessageWithJobSeekerLineImageParam
}

type SendChatMessageWithJobSeekerLineImageOutput struct {
	ChatMessageWithJobSeeker *entity.ChatMessageWithJobSeeker
}

func (i *ChatMessageWithJobSeekerInteractorImpl) SendChatMessageWithJobSeekerLineImage(input SendChatMessageWithJobSeekerLineImageInput) (SendChatMessageWithJobSeekerLineImageOutput, error) {
	var (
		output SendChatMessageWithJobSeekerLineImageOutput
		err    error
	)

	// GroupIDからGroupを取得
	chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByID(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agent, err := i.agentRepository.FindByID(chatGroup.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// chatGroup.JobSeekerIDから求職者を取得
	jobSeeker, err := i.jobSeekerRepository.FindByID(chatGroup.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	url, err := i.firebase.UploadToStorageForAgentLine(input.SendParam.File, agent.UUID.String())
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// シークレットとトークンを設定
	bot, err := linebot.New(agent.LineMessagingChannelSecret, agent.LineMessagingChannelAccessToken)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// LINE メッセージを送信する
	call, err := bot.PushMessage(
		jobSeeker.LineID,
		linebot.NewImageMessage(url, url),
	).Do()
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	fmt.Println("メッセージの送信に成功しました:", call)

	// エージェントから求職者へのチャットメッセージを作成
	chatMessage := entity.NewChatMessageWithJobSeeker(
		input.SendParam.GroupID,
		null.NewInt(0, true), // エージェント
		"",
		"",
		"",
		url,
		url,
		null.NewInt(0, false),
		"",
		null.NewInt(2, true),
		time.Now().In(time.UTC),
	)

	err = i.chatMessageWithJobSeekerRepository.Create(chatMessage)
	if err != nil {
		return output, err
	}

	// エージェントの最終送信時間を更新
	err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatMessageWithJobSeeker = chatMessage

	return output, nil
}

// チャットグループIDからチャットメッセージ情報を取得する
type GetChatMessageWithJobSeekerListByGroupIDInput struct {
	GroupID uint
}

type GetChatMessageWithJobSeekerListByGroupIDOutput struct {
	ChatMessageWithJobSeekerList []*entity.ChatMessageWithJobSeeker
}

func (i *ChatMessageWithJobSeekerInteractorImpl) GetChatMessageWithJobSeekerListByGroupID(input GetChatMessageWithJobSeekerListByGroupIDInput) (GetChatMessageWithJobSeekerListByGroupIDOutput, error) {
	var (
		output GetChatMessageWithJobSeekerListByGroupIDOutput
		err    error
	)

	chatMessageList, err := i.chatMessageWithJobSeekerRepository.GetByGroupID(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatMessageWithJobSeekerList = chatMessageList

	return output, nil
}

/****************************************************************************************/
/// LINE WebHook
//
type WebHookDestinationParam struct {
	Destination string `json:"destination"`
}

type LineWebHookInput struct {
	Request *http.Request
}

type LineWebHookOutput struct {
	OK bool
}

func (i *ChatMessageWithJobSeekerInteractorImpl) LineWebHook(input LineWebHookInput) (LineWebHookOutput, error) {
	var (
		request     = input.Request
		output      LineWebHookOutput
		destination WebHookDestinationParam
	)

	// リクエストの内容を取得
	bodyText, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(bodyText, &destination)
	if err != nil {
		fmt.Println(err)
	}

	// リクエストの再セット（これしないとEvent取れない）
	request.Body = io.NopCloser(bytes.NewBuffer(bodyText))

	agentLineChannel, err := i.agentRepository.FindLineChannelByBotID(destination.Destination)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 署名を検証する
	// 参照: https://developers.line.biz/ja/reference/messaging-api/#signature-validation
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("io.ReadAll(request.Body)のところでエラー")
		fmt.Println(err)
		return output, err
	}

	decoded, err := base64.StdEncoding.DecodeString(request.Header.Get("x-line-signature"))
	if err != nil {
		fmt.Println("base64.StdEncoding.DecodeString(request.Header.Get(x-line-signature))のところでエラー")
		fmt.Println(err)
		return output, err
	}

	hash := hmac.New(sha256.New, []byte(agentLineChannel.LineMessagingChannelSecret))
	_, err = hash.Write(body)
	if err != nil {
		fmt.Println("hash.Write(body)のところでエラー")
		fmt.Println(err)
		return output, err
	}

	// リクエストのx-line-signatureとシークレットをハッシュ化したものが一致するか検証
	if !hmac.Equal(decoded, hash.Sum(nil)) {
		fmt.Println("hmac.Equal(decoded, hash.Sum(nil))のところでエラー")
		fmt.Println(err)
		return output, err
	}

	// リクエストの再セット（これしないとEvent取れない）
	request.Body = io.NopCloser(bytes.NewBuffer(bodyText))

	bot, err := linebot.New(agentLineChannel.LineMessagingChannelSecret, agentLineChannel.LineMessagingChannelAccessToken)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// Webhookイベントオブジェクトの取得
	// 参考: https://developers.line.biz/ja/reference/messaging-api/#webhook-event-objects
	events, err := bot.ParseRequest(request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			fmt.Println("c.Response().WriteHeader(400)")
			return output, err
		} else {
			fmt.Println("c.Response().WriteHeader(500)")
			return output, err
		}
	}

	// イベントの重複を確認するための変数
	duplicateConfirmation := make(map[string]bool)

	for _, event := range events {
		fmt.Println("event", event)
		fmt.Println("送信者のID", event.Source.UserID)

		// 重複がある場合は処理しない
		if duplicateConfirmation[event.WebhookEventID] {
			fmt.Println("WebHookIDの重複をスキップ", event.WebhookEventID)
			continue
		}

		// LINEのIDがからの場合は
		if event.Source.UserID == "" {
			break
		}

		// LINEIDから求職者を取得
		jobSeeker, err := i.jobSeekerRepository.FindByAgentIDAndLineID(agentLineChannel.AgentID, event.Source.UserID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {

				/*********************** job_seekersに該当の求職者が見つからない場合 ***********************/
				output.OK = true
				return output, nil

				/*********************** job_seekersに該当の求職者が見つからない場合 ここまで ***********************/
			} else {
				// job_seekersテーブルの「FindByLineID」の実行結果が「Not Found」以外の場合はエラーを返す
				fmt.Println(err)
				return output, err
			}
		} else {
			// 通常求職者からのメッセージを処理
			err = MessageForJobSeeker(i, agentLineChannel, bot, event, jobSeeker)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 重複確認のために値を更新
		duplicateConfirmation[event.WebhookEventID] = true
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/

/****************************************************************************************/

// 通常求職者からのメッセージを処理する関数
func MessageForJobSeeker(i *ChatMessageWithJobSeekerInteractorImpl, agentLineChannel *entity.AgentLineChannelParam, bot *linebot.Client, event *linebot.Event, jobSeeker *entity.JobSeeker) error {
	var (
		err         error
		chatGroup   *entity.ChatGroupWithJobSeeker
		chatMessage *entity.ChatMessageWithJobSeeker
	)

	// 求職者IDとエージェントIDからグループを取得
	chatGroup, err = i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(agentLineChannel.AgentID, jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if event.Type == linebot.EventTypeMessage {
		// メッセージが送信された場合
		switch message := event.Message.(type) {
		// テキストメッセージの場合
		case *linebot.TextMessage:
			fmt.Println("message", message)
			fmt.Println("TextMessage")
			fmt.Println("----------------")

			// DB格納済みのメッセージと重複がないか確認
			_, err = i.chatMessageWithJobSeekerRepository.FindByLINEMessageID(message.ID)
			if err != nil {
				if errors.Is(err, entity.ErrNotFound) {
					// DBにメッセージを保存する処理
					chatMessage = entity.NewChatMessageWithJobSeeker(
						chatGroup.ID,         // chatGroup.IDを設定
						null.NewInt(1, true), // 求職者
						event.Message.(*linebot.TextMessage).Text,
						"",
						"",
						"",
						"",
						null.NewInt(0, false),
						message.ID,
						null.NewInt(0, true),
						event.Timestamp.In(time.UTC),
					)
				} else {
					fmt.Println(err)
					return err
				}
			}

		// スタンプの場合
		case *linebot.StickerMessage:
			fmt.Println("message", message)

			fmt.Println("StickerMessage")
			fmt.Println("----------------")

			// DB格納済みのメッセージと重複がないか確認
			_, err = i.chatMessageWithJobSeekerRepository.FindByLINEMessageID(message.ID)
			if err != nil {
				if errors.Is(err, entity.ErrNotFound) {
					// DBにメッセージを保存する処理
					chatMessage = entity.NewChatMessageWithJobSeeker(
						chatGroup.ID,         // chatGroup.IDを設定
						null.NewInt(1, true), // 求職者
						"スタンプが送信されましたがautoscoutでは表示できません。",
						message.PackageID,
						message.StickerID,
						"",
						"",
						null.NewInt(0, false),
						message.ID,
						null.NewInt(1, true),
						event.Timestamp.In(time.UTC),
					)
				} else {
					fmt.Println(err)
					return err
				}
			}

		// 画像データの場合
		case *linebot.ImageMessage:
			fmt.Println("message", message)
			fmt.Println("ImageMessage")
			fmt.Println("message.ContentProvider.Type", message.ContentProvider.Type)
			fmt.Println("----------------")

			// 画像のコンテンツを取得しにいく
			content, err := bot.GetMessageContent(message.ID).Do()
			if err != nil {
				fmt.Println("コンテンツの取得に失敗しました。")
			}
			defer content.Content.Close()

			url, err := i.firebase.UploadToStorageForJobSeekerLine(content.Content, message.ID)
			if err != nil {
				fmt.Println("URLの取得に失敗しました。")
			}

			// DB格納済みのメッセージと重複がないか確認
			_, err = i.chatMessageWithJobSeekerRepository.FindByLINEMessageID(message.ID)
			if err != nil {
				if errors.Is(err, entity.ErrNotFound) {
					// DBにメッセージを保存する処理
					chatMessage = entity.NewChatMessageWithJobSeeker(
						chatGroup.ID,         // chatGroup.IDを設定
						null.NewInt(1, true), // 求職者
						"",
						"",
						"",
						url,
						url,
						null.NewInt(0, false),
						message.ID,
						null.NewInt(2, true),
						event.Timestamp.In(time.UTC),
					)
				} else {
					fmt.Println(err)
					return err
				}
			}

		// 動画データの場合
		case *linebot.VideoMessage:
			fmt.Println("message", message)

			fmt.Println("VideoMessage")
			fmt.Println("----------------")

			// DB格納済みのメッセージと重複がないか確認
			_, err = i.chatMessageWithJobSeekerRepository.FindByLINEMessageID(message.ID)
			if err != nil {
				if errors.Is(err, entity.ErrNotFound) {
					// DBにメッセージを保存する処理
					chatMessage = entity.NewChatMessageWithJobSeeker(
						chatGroup.ID,         // chatGroup.IDを設定
						null.NewInt(1, true), // 求職者
						"動画データが送信されましたがautoscoutでは表示できません。",
						"",
						"",
						message.OriginalContentURL,
						message.PreviewImageURL,
						null.NewInt(0, false),
						message.ID,
						null.NewInt(3, true),
						event.Timestamp.In(time.UTC),
					)
				} else {
					fmt.Println(err)
					return err
				}
			}

		// 音声データの場合
		case *linebot.AudioMessage:
			fmt.Println("message", message)

			fmt.Println("AudioMessage")
			fmt.Println("----------------")

			// DB格納済みのメッセージと重複がないか確認
			_, err = i.chatMessageWithJobSeekerRepository.FindByLINEMessageID(message.ID)
			if err != nil {
				if errors.Is(err, entity.ErrNotFound) {
					// DBにメッセージを保存する処理
					chatMessage = entity.NewChatMessageWithJobSeeker(
						chatGroup.ID,         // chatGroup.IDを設定
						null.NewInt(1, true), // 求職者
						"音声データが送信されましたがautoscoutでは表示できません。",
						"",
						"",
						message.OriginalContentURL,
						"",
						null.NewInt(int64(message.Duration), true),
						message.ID,
						null.NewInt(4, true),
						event.Timestamp.In(time.UTC),
					)
				} else {
					fmt.Println(err)
					return err
				}
			}
		}

		fmt.Println("chatMessageの値", chatMessage)

		// chatMessageが有効な値の場合
		if chatMessage != nil {
			err = i.chatMessageWithJobSeekerRepository.Create(chatMessage)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// 求職者の送信時間と閲覧時間を更新する
			err = i.chatGroupWithJobSeekerRepository.UpdateJobSeekerLastWatchedAtAndSendAt(chatGroup.ID, event.Timestamp.In(time.UTC))
			if err != nil {
				fmt.Println(err)
				return err
			}

			// プッシュ通知
			redirectURL := os.Getenv("BASE_DOMAIN" + "/")
			topic := "Chat"

			if jobSeeker.AgentStaffID.Int64 != 0 {
				caStaff, err := i.agentStaffRepository.FindByID(uint(jobSeeker.AgentStaffID.Int64))
				if err != nil {
					fmt.Println(err)
					return err
				}

				err = utility.WebPush(
					i.oneSignal.AppID,
					i.oneSignal.APIKey,
					caStaff.FirebaseID,
					"チャット通知",
					"新着メッセージがあります",
					topic,
					redirectURL,
				)
				if err != nil {
					fmt.Println("WebPushの通知でエラー")
					fmt.Println(err)
				}
			}
		}

	} else if event.Type == linebot.EventTypeUnfollow {
		// ブロックされた場合 LINE Active
		err = i.chatGroupWithJobSeekerRepository.UpdateJobSeekerLineActive(jobSeeker.ID, false)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else if event.Type == linebot.EventTypeFollow {
		// ブロック解除された場合 LINE Activeをtrueにする
		// 友達新規追加の場合は上の「i.jobSeekerRepository.FindByLineID(event.Source.UserID)」の部分で処理してる
		err = i.chatGroupWithJobSeekerRepository.UpdateJobSeekerLineActive(jobSeeker.ID, true)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
