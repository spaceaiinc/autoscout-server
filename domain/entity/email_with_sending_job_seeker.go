package entity

import (
	"mime/multipart"
	"time"
)

// 送客求職者とのメールを管理する
type EmailWithSendingJobSeeker struct {
	ID                 uint      `db:"id" json:"id"`
	SendingJobSeekerID uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	Subject            string    `db:"subject" json:"subject"`     // メール件名
	Content            string    `db:"content" json:"content"`     // メール本文
	FileName           string    `db:"file_name" json:"file_name"` // 添付ファイル名(複数の場合は"|--|--T--|--|"区切りで保存)
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

func NewEmailWithSendingJobSeeker(
	sendingJobSeekerID uint,
	subject string,
	content string,
	fileName string,
) *EmailWithSendingJobSeeker {
	return &EmailWithSendingJobSeeker{
		SendingJobSeekerID: sendingJobSeekerID,
		Subject:            subject,
		Content:            content,
		FileName:           fileName,
	}
}

// メール送信時&保存のパラメータ
type SendEmailWithSendingJobSeekerParam struct {
	SendingJobSeekerID uint   `json:"sending_job_seeker_id" validate:"required"` // 送客求職者ID
	Subject            string `json:"subject" validate:"required"`               // メールの件名
	Content            string `json:"content" validate:"required"`               // メールの本文
	FileName           string `json:"file_name"`                                 // 添付ファイル名

	AgentStaffID uint                    `json:"agent_staff_id" validate:"required"` // エージェントID
	GroupID      uint                    `json:"group_id" validate:"required"`       // チャットグループID
	To           EmailUser               `json:"to" validate:"required"`             // 送り先（求職者の名前とメールアドレス）
	Files        []*multipart.FileHeader `json:"files"`                              // 添付ファイル
}
