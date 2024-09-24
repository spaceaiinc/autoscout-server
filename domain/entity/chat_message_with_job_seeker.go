package entity

import (
	"mime/multipart"
	"time"

	"gopkg.in/guregu/null.v4"
)

type ChatMessageWithJobSeeker struct {
	ID                 uint      `db:"id" json:"id"`
	GroupID            uint      `db:"group_id" json:"group_id"`                         // チャットグループID
	UserType           null.Int  `db:"user_type" json:"user_type"`                       // エージェント（0）or 求職者（1）
	Message            string    `db:"message" json:"message"`                           // メッセージ
	PackageID          string    `db:"package_id" json:"package_id"`                     // スタンプの表示に使用するID
	StickerID          string    `db:"sticker_id" json:"sticker_id"`                     // スタンプの表示に使用するID
	OriginalContentURL string    `db:"original_content_url" json:"original_content_url"` // 画像ファイル or 動画ファイル or 音声ファイルのURL
	PreviewImageURL    string    `db:"preview_image_url" json:"preview_image_url"`       // 画像ファイル or 動画ファイルのプレビュー表示用のファイルURL
	Duration           null.Int  `db:"duration" json:"duration"`                         // 音声ファイルに使用する値
	LineMessageID      string    `db:"line_message_id" json:"line_message_id"`           // LINEのメッセージID
	LineMessageType    null.Int  `db:"line_message_type" json:"line_message_type"`       // LINEのメッセージタイプ （0: テキスト, 1: スタンプ, 2: 画像, 3: 動画, 4: 音声
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"-"`
}

func NewChatMessageWithJobSeeker(
	groupID uint,
	userType null.Int,
	message string,
	packageID string,
	stickerID string,
	originalContentURL string,
	previewImageURL string,
	duration null.Int,
	lineMessageID string,
	lineMessageType null.Int,
	sendAt time.Time, // 求職者がLINEを送信した時間を正確に取得するため
) *ChatMessageWithJobSeeker {
	return &ChatMessageWithJobSeeker{
		GroupID:            groupID,
		UserType:           userType,
		Message:            message,
		PackageID:          packageID,
		StickerID:          stickerID,
		OriginalContentURL: originalContentURL,
		PreviewImageURL:    previewImageURL,
		Duration:           duration,
		LineMessageID:      lineMessageID,
		LineMessageType:    lineMessageType,
		CreatedAt:          sendAt,
		UpdatedAt:          sendAt,
	}
}

type SendChatMessageWithJobSeekerLineParam struct {
	GroupID uint   `json:"group_id" validate:"required"` // チャットグループID
	Message string `json:"message" validate:"required"`  // メッセージ
	File    string `json:"file"`
}

type SendChatMessageWithJobSeekerLineImageParam struct {
	GroupID uint   // チャットグループID
	File    *multipart.FileHeader
}
