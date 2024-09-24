-- ユーザーがデプロイの内容を反映したかを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS deployment_informations (
  id INT AUTO_INCREMENT NOT NULL UNIQUE,      -- 重複しないID
  be_ver VARCHAR(255) NOT NULL,               -- バックエンドのバージョン
  fe_ver VARCHAR(255) NOT NULL,               -- フロントエンドのバージョン
  be_detail TEXT NOT NULL,                    -- バックエンドの詳細
  fe_detail TEXT NOT NULL,                    -- フロントエンドの詳細
  created_at DATETIME,                        -- 作成日時
  updated_at DATETIME,                        -- 最終更新日時
  PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE IF EXISTS deployment_informations;
