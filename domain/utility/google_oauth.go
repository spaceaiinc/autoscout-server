package utility

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GoogleOAuthのリダイレクトURL作成用Configを生成
func NewGoogleAuthConf(filePath string) (*oauth2.Config, error) {

	// OAuth2.0クライアント設定を読み込む
	b, err := os.ReadFile(filePath)
	if err != nil {
		wrapped := fmt.Errorf("認証情報ファイルの読み込みに失敗しました: %v", err)
		return nil, wrapped
	}

	// 認証情報からconfigを生成（gmail.GmailReadonlyScope || gmail.GmailModifyScope）
	config, err := google.ConfigFromJSON(b, gmail.GmailModifyScope)
	if err != nil {
		wrapped := fmt.Errorf("認証情報からのconfig生成に失敗しました: %v", err)
		return nil, wrapped
	}

	return config, nil
}

// Gmailサービスを生成
func GenerateService(config *oauth2.Config, token *oauth2.Token) (*gmail.Service, error) {
	ctx := context.Background()

	// 新しいHTTPクライアントを生成
	// client := config.Client(ctx, token)

	// persistent storage から, oauth2.Token を取得
	tokenSource := config.TokenSource(context.Background(), token)

	// アクセストークンの有効期限が切れていた場合、refresh_token から新たなアクセストークンを発行する
	client := oauth2.NewClient(context.Background(), tokenSource)

	// newTok, err := tokenSource.Token()
	// if err != nil {
	// 	wrapped := fmt.Errorf("tokenSource.Token(): %v", err)
	// 	return nil, wrapped
	// }

	// fmt.Println("newTok.AccessToken", newTok.AccessToken)
	// fmt.Println("newTok.RefreshToken", newTok.RefreshToken)
	// fmt.Println("newTok.Expiry", newTok.Expiry)

	// Gmailサービスを生成
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		wrapped := fmt.Errorf("gmailサービスの生成に失敗しました: %v", err)
		return nil, wrapped
	}

	return srv, nil
}
