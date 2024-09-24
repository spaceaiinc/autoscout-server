package utility

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

type Sendgrid struct {
	Client *sendgrid.Client
}

func NewSendGrid(apiKey string) *Sendgrid {
	client := sendgrid.NewSendClient(apiKey)
	return &Sendgrid{
		Client: client,
	}
}

func (s *Sendgrid) SendMail(from, to *mail.Email, subject, plainTextContent, htmlContent string) error {
	if to.Address == from.Address {
		err := fmt.Errorf("送信元と送信先が同じため、送信できません。")
		return err
	}

	htmlContent = ""
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	_, err := s.Client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

// 複数の宛先に送信
func SendMailToMultiple(
	apiKey string,
	subject string,
	content string,
	from entity.EmailUser,
	toList []entity.EmailUser,
	fileParams []*multipart.FileHeader,
) error {
	var (
		tos []*mail.Email
	)

	for _, to := range toList {
		if to.Email == from.Email {
			err := fmt.Errorf("送信元と送信先が同じため、送信できません。")
			return err
		}
	}

	// メッセージの構築
	message := mail.NewV3Mail()

	appEnv := os.Getenv("APP_ENV")

	// 本番以外は場合はinfo@spaceai.jpから送付
	// https://support.sendgrid.kke.co.jp/hc/ja/articles/203574969
	if appEnv != "prd" {
		if !strings.Contains(from.Email, "motoyui.com") {
			err := fmt.Errorf("送信元が「motoyui.com」ではありません。\nSendGridの「Domain Authentication」で設定したドメインをご確認ください。")
			return err
		}
	}

	// 送信元を設定
	f := mail.NewEmail(from.Name, from.Email)
	message.SetFrom(f)

	// 1つ目の宛先と、対応するSubstitutionタグを指定
	p := mail.NewPersonalization()

	// to
	for _, to := range toList {
		t := mail.NewEmail(to.Name, to.Email)
		tos = append(tos, t)
	}
	p.AddTos(tos...)

	// cc（送信者をCCに入れることで）
	cc := mail.NewEmail(from.Name, from.Email)
	p.AddCCs(cc)

	message.AddPersonalizations(p)

	// 件名を設定
	message.Subject = subject

	// 本文
	c := mail.NewContent("text/plain", content)
	message.AddContent(c)

	// カスタムヘッダを指定
	message.SetHeader("X-Sent-Using", "SendGrid-API")

	// ファイルがあれば添付の処理
	if fileParams != nil && len(fileParams) != 0 {
		for _, f := range fileParams {
			a := mail.NewAttachment()

			file, err := f.Open()
			if err != nil {
				log.Println("ファイルが開けません")
				log.Println(err)
				return err
			}

			defer file.Close()
			data, _ := io.ReadAll(file)

			data_enc := base64.StdEncoding.EncodeToString(data)
			a.SetContent(data_enc)

			a.SetFilename(f.Filename)

			a.SetDisposition("attachment")
			message.AddAttachment(a)
		}
	}

	// メール送信を行い、レスポンスを表示
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("StatusCode", response.StatusCode)
	fmt.Println("Headers", response.Headers)
	fmt.Println("Body", response.Body)

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrRequestError)
		return err
	} else if response.StatusCode >= 500 && response.StatusCode < 600 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrServerError)
		return err
	}

	return nil
}

func SendMailToSingle(
	apiKey string,
	subject string,
	content string,
	from entity.EmailUser,
	to entity.EmailUser,
	fileParams []*multipart.FileHeader, // 必要ない場合はnilが入ってる
) error {
	if to.Email == from.Email {
		err := fmt.Errorf("送信元と送信先が同じため、送信できません。")
		return err
	}

	// メッセージの構築
	message := mail.NewV3Mail()

	appEnv := os.Getenv("APP_ENV")

	// 本番以外は場合はinfo@spaceai.jpから送付
	// https://support.sendgrid.kke.co.jp/hc/ja/articles/203574969
	if appEnv != "prd" {
		if !strings.Contains(from.Email, "motoyui.com") {
			err := fmt.Errorf("送信元が「motoyui.com」ではありません。\nSendGridの「Domain Authentication」で設定したドメインをご確認ください。")
			return err
		}
	}

	// 送信元を設定
	f := mail.NewEmail(from.Name, from.Email)
	message.SetFrom(f)

	// 1つ目の宛先と、対応するSubstitutionタグを指定
	p := mail.NewPersonalization()

	// to
	t := mail.NewEmail(to.Name, to.Email)
	p.AddTos(t)

	// cc（送信者をCCに入れることで）
	cc := mail.NewEmail(from.Name, from.Email)
	p.AddCCs(cc)

	message.AddPersonalizations(p)

	// 件名を設定
	message.Subject = subject

	// 本文
	c := mail.NewContent("text/plain", content)
	message.AddContent(c)

	// カスタムヘッダを指定
	message.SetHeader("X-Sent-Using", "SendGrid-API")

	// ファイルがあれば添付の処理
	if fileParams != nil && len(fileParams) != 0 {
		for _, f := range fileParams {
			a := mail.NewAttachment()

			file, err := f.Open()
			if err != nil {
				log.Println("ファイルが開けません")
				log.Println(err)
				return err
			}

			defer file.Close()
			data, _ := io.ReadAll(file)

			data_enc := base64.StdEncoding.EncodeToString(data)
			a.SetContent(data_enc)

			a.SetFilename(f.Filename)

			a.SetDisposition("attachment")
			message.AddAttachment(a)
		}
	}

	// メール送信を行い、レスポンスを表示
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("StatusCode", response.StatusCode)
	fmt.Println("Headers", response.Headers)
	fmt.Println("Body", response.Body)

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrRequestError)
		return err
	} else if response.StatusCode >= 500 && response.StatusCode < 600 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrServerError)
		return err
	}

	return nil
}

// ccなし *未読通知で使用
func SendMailToSingleWithoutCC(
	apiKey string,
	subject string,
	content string,
	from entity.EmailUser,
	to entity.EmailUser,
	fileParams []*multipart.FileHeader, // 必要ない場合はnilが入ってる
) error {
	if to.Email == from.Email {
		err := fmt.Errorf("送信元と送信先が同じため、送信できません。")
		return err
	}

	// メッセージの構築
	message := mail.NewV3Mail()

	appEnv := os.Getenv("APP_ENV")

	// 本番以外は場合はinfo@spaceai.jpから送付
	// https://support.sendgrid.kke.co.jp/hc/ja/articles/203574969
	if appEnv != "prd" {
		if !strings.Contains(from.Email, "motoyui.com") {
			err := fmt.Errorf("送信元が「motoyui.com」ではありません。\nSendGridの「Domain Authentication」で設定したドメインをご確認ください。")
			return err
		}
	}

	// 送信元を設定
	f := mail.NewEmail(from.Name, from.Email)
	message.SetFrom(f)

	// 1つ目の宛先と、対応するSubstitutionタグを指定
	p := mail.NewPersonalization()

	// to
	t := mail.NewEmail(to.Name, to.Email)
	p.AddTos(t)

	message.AddPersonalizations(p)

	// 件名を設定
	message.Subject = subject

	// 本文
	c := mail.NewContent("text/plain", content)
	message.AddContent(c)

	// カスタムヘッダを指定
	message.SetHeader("X-Sent-Using", "SendGrid-API")

	// ファイルがあれば添付の処理
	if fileParams != nil && len(fileParams) != 0 {
		for _, f := range fileParams {
			a := mail.NewAttachment()

			file, err := f.Open()
			if err != nil {
				log.Println("ファイルが開けません")
				log.Println(err)
				return err
			}

			defer file.Close()
			data, _ := io.ReadAll(file)

			data_enc := base64.StdEncoding.EncodeToString(data)
			a.SetContent(data_enc)

			a.SetFilename(f.Filename)

			a.SetDisposition("attachment")
			message.AddAttachment(a)
		}
	}

	// メール送信を行い、レスポンスを表示
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("StatusCode", response.StatusCode)
	fmt.Println("Headers", response.Headers)
	fmt.Println("Body", response.Body)

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrRequestError)
		return err
	} else if response.StatusCode >= 500 && response.StatusCode < 600 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrServerError)
		return err
	}

	return nil
}

func SendMailToMultipleFileName(
	apiKey string,
	subject string,
	content string,
	from entity.EmailUser,
	toList []entity.EmailUser,
	data []byte, // 必要ない場合はnilが入ってる
	fileName string,
) error {
	var (
		tos []*mail.Email
	)

	for _, to := range toList {
		if to.Email == from.Email {
			err := fmt.Errorf("送信元と送信先が同じため、送信できません。")
			return err
		}
	}

	// メッセージの構築
	message := mail.NewV3Mail()

	appEnv := os.Getenv("APP_ENV")

	// 本番以外は場合はinfo@spaceai.jpから送付
	// https://support.sendgrid.kke.co.jp/hc/ja/articles/203574969
	if appEnv != "prd" {
		if !strings.Contains(from.Email, "motoyui.com") {
			err := fmt.Errorf("送信元が「motoyui.com」ではありません。\nSendGridの「Domain Authentication」で設定したドメインをご確認ください。")
			return err
		}
	}

	// 送信元を設定
	f := mail.NewEmail(from.Name, from.Email)
	message.SetFrom(f)

	// 1つ目の宛先と、対応するSubstitutionタグを指定
	p := mail.NewPersonalization()

	// to
	for _, to := range toList {
		t := mail.NewEmail(to.Name, to.Email)
		tos = append(tos, t)
	}
	p.AddTos(tos...)
	// t := mail.NewEmail(to.Name, to.Email)
	// p.AddTos(t)

	// cc（送信者をCCに入れることで）
	cc := mail.NewEmail(from.Name, from.Email)
	p.AddCCs(cc)

	message.AddPersonalizations(p)

	// 件名を設定
	message.Subject = subject

	// 本文
	c := mail.NewContent("text/plain", content)
	message.AddContent(c)

	// カスタムヘッダを指定
	message.SetHeader("X-Sent-Using", "SendGrid-API")

	// ファイルがあれば添付の処理
	if data != nil {
		a := mail.NewAttachment()

		data_enc := base64.StdEncoding.EncodeToString(data)
		a.SetContent(data_enc)

		a.SetFilename(fileName)

		a.SetDisposition("attachment")
		message.AddAttachment(a)
	}

	// メール送信を行い、レスポンスを表示
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("StatusCode", response.StatusCode)
	fmt.Println("Headers", response.Headers)
	fmt.Println("Body", response.Body)

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrRequestError)
		return err
	} else if response.StatusCode >= 500 && response.StatusCode < 600 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrServerError)
		return err
	}

	return nil
}

// 自分から自分へメール送信（求職者のアンケート回答通知用）
func SendMailToSingleByMyself(
	apiKey string,
	subject string,
	content string,
	fromAndTo entity.EmailUser,
) error {
	// メッセージの構築
	message := mail.NewV3Mail()

	// 送信元を設定
	ft := mail.NewEmail(fromAndTo.Name, fromAndTo.Email)
	message.SetFrom(ft)

	// 宛先を設定
	p := mail.NewPersonalization()
	p.AddTos(ft)

	message.AddPersonalizations(p)

	// 件名を設定
	message.Subject = subject

	// 本文
	c := mail.NewContent("text/plain", content)
	message.AddContent(c)

	// カスタムヘッダを指定
	message.SetHeader("X-Sent-Using", "SendGrid-API")

	// メール送信を行い、レスポンスを表示
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("StatusCode", response.StatusCode)
	fmt.Println("Headers", response.Headers)
	fmt.Println("Body", response.Body)

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrRequestError)
		return err
	} else if response.StatusCode >= 500 && response.StatusCode < 600 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrServerError)
		return err
	}

	return nil
}

// 複数の宛先にToでメール送信（RPAのエラー通知用）*ccとファイル添付なし
func SendMailList(
	apiKey string,
	subject string,
	content string,
	from entity.EmailUser,
	toList []entity.EmailUser,
) error {
	for _, to := range toList {
		if to.Email == from.Email {
			err := fmt.Errorf("送信元と送信先が同じため、送信できません。")
			return err
		}
	}

	// メッセージの構築
	message := mail.NewV3Mail()

	appEnv := os.Getenv("APP_ENV")

	// 本番以外は場合はinfo@spaceai.jpから送付
	// https://support.sendgrid.kke.co.jp/hc/ja/articles/203574969
	if appEnv != "prd" {
		if !strings.Contains(from.Email, "motoyui.com") {
			err := fmt.Errorf("送信元が「motoyui.com」ではありません。\nSendGridの「Domain Authentication」で設定したドメインをご確認ください。")
			return err
		}
	}

	// 送信元を設定
	f := mail.NewEmail(from.Name, from.Email)
	message.SetFrom(f)

	// 1つ目の宛先と、対応するSubstitutionタグを指定
	p := mail.NewPersonalization()

	// 宛先を追加
	for _, to := range toList {
		t := mail.NewEmail(to.Name, to.Email)
		p.AddTos(t)
	}

	message.AddPersonalizations(p)

	// 件名を設定
	message.Subject = subject

	// 本文
	c := mail.NewContent("text/plain", content)
	message.AddContent(c)

	// カスタムヘッダを指定
	message.SetHeader("X-Sent-Using", "SendGrid-API")

	// メール送信を行い、レスポンスを表示
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("StatusCode", response.StatusCode)
	fmt.Println("Headers", response.Headers)
	fmt.Println("Body", response.Body)

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrRequestError)
		return err
	} else if response.StatusCode >= 500 && response.StatusCode < 600 {
		err = fmt.Errorf("%s:%w", response.Body, entity.ErrServerError)
		return err
	}

	return nil
}
