package entity

import (
	"time"
)

type DeploymentInformation struct {
	ID        uint      `db:"id" json:"id"`                 // 重複しないID
	BeVer     string    `db:"be_ver" json:"be_ver"`         // バックエンドバージョン
	FeVer     string    `db:"fe_ver" json:"fe_ver"`         // フロントエンドバージョン
	BeDetail  string    `db:"be_detail" json:"be_detail"`   // バックエンド詳細
	FeDetail  string    `db:"fe_detail" json:"fe_detail"`   // フロントエンド詳細
	CreatedAt time.Time `db:"created_at" json:"created_at"` // 作成日時
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"` // 更新日時
}

func NewDeploymentInformation(
	beVer string,
	feVer string,
	beDetail string,
	feDetail string,
) *DeploymentInformation {
	return &DeploymentInformation{
		BeVer:    beVer,
		FeVer:    feVer,
		BeDetail: beDetail,
		FeDetail: feDetail,
	}
}

// デプロイ情報登録時のパラム
type CreateOrUpdateDeploymentInformationParam struct {
	BeVer    string `json:"be_ver" validate:"required"` // バックエンドバージョン
	FeVer    string `json:"fe_ver" validate:"required"` // フロントエンドバージョン
	BeDetail string `json:"be_detail"`                  // バックエンド詳細
	FeDetail string `json:"fe_detail"`                  // フロントエンド詳細
}
