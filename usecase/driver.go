package usecase

import (
	"io"
	"mime/multipart"
	"time"
)

type Firebase interface {
	VerifyIDToken(idToken string) (string, error)
	GetCustomToken(uid string) (string, error)
	GetCustomTokenWithTimeLimit(uid string, add60minuteTime time.Time) (string, error)
	GetIDToken(token string) (string, error)
	GetPhoneNumber(uid string) (string, error)
	Set(doc string, data map[string]interface{}) error
	CreateUser(email, password string) (string, error)
	UpdateEmail(email, uid string) error
	UpdatePassword(password, uid string) error
	DeleteUser(uid string) error
	UploadToStorageForJobSeekerLine(content io.ReadCloser, messageID string) (string, error)
	UploadToStorageForAgentLine(file *multipart.FileHeader, agentUUID string) (string, error)
	SignOut(uid string) error
}

type Cache interface {
	GetBytes(key string) ([]byte, error)
	GetString(key string) (string, error)
	Set(key string, obj interface{}, ttl int) (interface{}, error)
	Do(commandName string, args ...interface{}) (interface{}, error)
	Values(reply interface{}, err error) ([]interface{}, error)
}
