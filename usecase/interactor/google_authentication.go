package interactor

import (
	"context"
	"errors"
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

type GoogleAuthenticationInteractor interface {
	// Gest API
	GetGoogleAuthCodeURL(input GetGoogleAuthCodeURLInput) (GetGoogleAuthCodeURLOutput, error)
	UpdateGoogleOauthToken(input UpdateGoogleOauthTokenInput) (UpdateGoogleOauthTokenOutput, error)
}

type GoogleAuthenticationInteractorImpl struct {
	firebase                       usecase.Firebase
	sendgrid                       config.Sendgrid
	googleAPI                      config.GoogleAPI
	agentRepository                usecase.AgentRepository
	agentStaffRepository           usecase.AgentStaffRepository
	googleAuthenticationRepository usecase.GoogleAuthenticationRepository
}

func NewGoogleAuthenticationInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	ga config.GoogleAPI,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
	gaR usecase.GoogleAuthenticationRepository,
) GoogleAuthenticationInteractor {
	return &GoogleAuthenticationInteractorImpl{
		firebase:                       fb,
		sendgrid:                       sg,
		googleAPI:                      ga,
		agentRepository:                aR,
		agentStaffRepository:           asR,
		googleAuthenticationRepository: gaR,
	}
}

type GetGoogleAuthCodeURLInput struct {
	Token string
}

type GetGoogleAuthCodeURLOutput struct {
	AuthURL string
}

// 参照: https://qiita.com/Kentaro91011/items/e8f9c5764ff6792c1d98
// 参照: https://hongo.dev/blog/nodemailer-send-emails-using-alias-address-for-gmail/
func (i *GoogleAuthenticationInteractorImpl) GetGoogleAuthCodeURL(input GetGoogleAuthCodeURLInput) (GetGoogleAuthCodeURLOutput, error) {
	var (
		output = GetGoogleAuthCodeURLOutput{}
	)

	/************ ユーザー情報取得 **************/

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	agentStaff, err := i.agentStaffRepository.FindByFirebaseID(firebaseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if agentStaff.ID != 1 {
		wrapped := fmt.Errorf("エンジニアのアカウント以外では実行できません。")
		return output, wrapped
	}

	/************ 認証情報取得とURL生成 **************/

	// 認証情報からconfigを生成
	config, err := utility.NewGoogleAuthConf(i.googleAPI.JSONFilePath)
	if err != nil {
		return output, err
	}

	output.AuthURL = config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	return output, nil
}

type UpdateGoogleOauthTokenInput struct {
	Token string
	Code  string
}

type UpdateGoogleOauthTokenOutput struct {
	OK bool
}

func (i *GoogleAuthenticationInteractorImpl) UpdateGoogleOauthToken(input UpdateGoogleOauthTokenInput) (UpdateGoogleOauthTokenOutput, error) {
	var (
		output   UpdateGoogleOauthTokenOutput
		err      error
		newToken *oauth2.Token
	)

	/************ ユーザー情報を取得 **************/

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	agentStaff, err := i.agentStaffRepository.FindByFirebaseID(firebaseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if agentStaff.ID != 1 {
		wrapped := fmt.Errorf("エンジニアのアカウント以外では実行できません。")
		return output, wrapped
	}

	/************ 認証情報取得 **************/

	// 認証情報からconfigを生成
	config, err := utility.NewGoogleAuthConf(i.googleAPI.JSONFilePath)
	if err != nil {
		return output, err
	}

	/************ トークン生成 **************/

	tok, err := config.Exchange(context.Background(), input.Code)
	if err != nil {
		wrapped := fmt.Errorf("トークンを取得できません: %v", err)
		return output, wrapped
	}

	// 参考: https://kitagry.github.io/blog/programmings/2020/11/go-google-oauth2/
	// refresh_tokenは最初のログインの時しか返してくれない.(refresh_tokenは有効時間が長いので発行しすぎは注意しなければならないから)
	// もし、開発時にrefresh_tokenを発行しなおさないといけないってなったら、以下のリンクにアクセスして、permissionを消すとまた発行しなおされる.
	// https://myaccount.google.com/u/0/permissions?pli=1
	// トークンの再取得: https://qiita.com/yyoshiki41/items/6f5bc207ec3418d4f0c5

	/************ トークン情報をDBに登録 **************/

	ga := entity.NewGoogleAuthentication(tok.AccessToken, tok.RefreshToken, tok.Expiry)

	newToken = &oauth2.Token{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
	}

	latestGA, err := i.googleAuthenticationRepository.FindLatest()
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			err = i.googleAuthenticationRepository.Create(ga)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			fmt.Println(err)
			return output, err
		}
	} else {
		err = i.googleAuthenticationRepository.UpdateTokenAndExpiry(latestGA.ID, tok.AccessToken, tok.Expiry)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// すでにレコードがある場合は
		newToken.RefreshToken = latestGA.RefreshToken
	}

	/************ サブスクリプションで登録したURLを叩くようにする **************/

	srv, err := utility.GenerateService(config, newToken)
	if err != nil {
		return output, err
	}

	user := "me"

	// ここにwatchリクエストを追加
	var topicName = i.googleAPI.TopicName // 環境変数にする
	watchRes, err := srv.Users.Watch(user, &gmail.WatchRequest{TopicName: topicName}).Do()
	fmt.Println(watchRes)
	if err != nil {
		wrapped := fmt.Errorf("unable to retrieve watch action: %v", err)
		return output, wrapped
	}

	output.OK = true

	return output, nil
}
