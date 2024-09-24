package utility

import (
	"log"

	"github.com/slack-go/slack"
)

type Slack struct {
	Client *slack.Client
}

func NewSlack(accessToken string) *Slack {
	client := slack.New(accessToken)
	return &Slack{
		Client: client,
	}
}

// autoscoutからのお問い合わせを通知
func (i *Slack) SendContact(chanelID, content string) error {
	var (
		err error
	)

	// アクセストークンを使用してクライアントを生成する
	// https://api.slack.com/apps からトークン取得
	// 参考: https://risaki-masa.com/how-to-get-api-token-in-slack/
	// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
	_, _, err = i.Client.PostMessage(chanelID, slack.MsgOptionText(
		content,
		true,
		// &slack.SectionBlock{
		// 	Type: slack.MBTSection,
		// 	Text: &slack.TextBlockObject{
		// 		Type: "mrkdwn",
		// 		Text: te,
		// 	},
		// },
	))
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("slack通知完了")
	return err
}
