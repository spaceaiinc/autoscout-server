package entity

import (
	"time"
)

type GoogleAuthentication struct {
	ID           uint      `db:"id" json:"id"`
	AccessToken  string    `db:"acess_token" json:"acess_token"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
	Expiry       time.Time `db:"expiry" json:"expiry"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func NewGoogleAuthentication(
	accessToken string,
	refreshToken string,
	expiry time.Time,
) *GoogleAuthentication {
	return &GoogleAuthentication{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}
}

type GoogleAuthParam struct {
	AuthCode string `json:"auth_code" validate:"required"`
}

type PubsubStruct struct {
	Subscription string `json:"subscription"`
	Message      MessageStruct
}

type MessageStruct struct {
	Data      string `json:"data"`
	MessageId string `json:"message_id"`
}

type MeessageData struct {
	EmailAddress string `json:"emailAddress"`
	HistoryID    int    `json:"historyId"`
}
