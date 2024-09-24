package entity

import "mime/multipart"

// "mime/multipart"
// "net/textproto"

type EmailUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// type FileHeader struct {
// 	Filename string
// 	Header   textproto.MIMEHeader
// 	Size     int64
// 	content  []byte
// 	tmpfile  string
// }

type SendEmailParam struct {
	AgentStaffID uint        `json:"agent_staff_id" validate:"required"` // 重複しないカラム毎の担当者のid
	Tos          []EmailUser `json:"tos" validate:"required"`            // 送信先の情報
	Subject      string      `json:"subject" validate:"required"`        // メールの件名
	Content      string      `json:"content" validate:"required"`        // メールの本文
	Files        []string    `json:"files"`
}

type SendEmail struct {
	AgentStaffID uint
	Tos          []EmailUser
	Subject      string
	Content      string
	Files        []*multipart.FileHeader
}
