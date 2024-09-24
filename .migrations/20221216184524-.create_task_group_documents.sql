-- タスクグループの書類を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS task_group_documents (
  id	INT AUTO_INCREMENT NOT NULL UNIQUE,	              -- 重複しないID
  task_group_id	INT NOT NULL UNIQUE,	                  -- タスクグループID
  document1_url	VARCHAR(2000) NOT NULL,	                -- 添付ファイル1
  document2_url	VARCHAR(2000) NOT NULL,	                -- 添付ファイル2
  document3_url	VARCHAR(2000) NOT NULL,                 -- 添付ファイル3
  document4_url	VARCHAR(2000) NOT NULL,	                -- 添付ファイル4
  document5_url VARCHAR(2000) NOT NULL,	                -- 添付ファイル5
  created_at	DATETIME,	                                -- 作成日時
  updated_at	DATETIME,	                                -- 最終更新日時	    
  PRIMARY KEY(id),
  INDEX idx_task_group_documents_task_group_id (task_group_id)
);

ALTER TABLE task_group_documents
  ADD CONSTRAINT fk_task_group_documents_task_group_id
  FOREIGN KEY(task_group_id)
  REFERENCES task_groups (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE task_group_documents DROP FOREIGN KEY fk_task_group_documents_task_group_id;

DROP TABLE IF EXISTS task_group_documents;