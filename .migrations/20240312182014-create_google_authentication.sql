-- google apiのToken管理
-- +migrate Up
CREATE TABLE IF NOT EXISTS google_authentication (
  id INT AUTO_INCREMENT NOT NULL UNIQUE,         -- 重複しないID
  acess_token TEXT NOT NULL,	                   -- アクセストークン
  refresh_token TEXT NOT NULL,                   -- リフレッシュトークン
  expiry DATETIME,                               -- 有効期限
  created_at	DATETIME,                          -- 作成日時
  updated_at	DATETIME,                          -- 最終更新日時
  PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE IF EXISTS google_authentication;