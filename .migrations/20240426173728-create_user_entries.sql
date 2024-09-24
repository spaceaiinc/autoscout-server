-- エントリー処理実行待ちユーザーを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS user_entries (
  id INT AUTO_INCREMENT NOT NULL UNIQUE,       -- 重複しないID
  user_id VARCHAR(255) NOT NULL,               -- メールに記載されているユーザーID
  service_type INT NOT NULL,                   -- 対象の媒体（0: RAN, 1: マイナビ転職スカウト, 2: AMBI, 3 マイナビエージェントスカウト）
  is_processed BOOLEAN NOT NULL DEFAULT FALSE, -- 取り込みを完了
  created_at	DATETIME,                        -- 作成日時
  updated_at	DATETIME,                        -- 最終更新日時
  PRIMARY KEY(id),
  INDEX idx_user_entries_user_id (user_id),
  INDEX idx_user_entries_is_processed (is_processed)
);

-- +migrate Down
DROP TABLE IF EXISTS user_entries;