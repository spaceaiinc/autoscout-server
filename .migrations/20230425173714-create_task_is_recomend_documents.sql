-- 送付をお勧めする書類
-- +migrate Up
CREATE TABLE IF NOT EXISTS task_is_recommend_documents (
  id	INT AUTO_INCREMENT NOT NULL UNIQUE,	              -- 重複しないID
  task_id INT NOT NULL UNIQUE,	                        -- タスクID
  is_recommend_upload_resume BOOLEAN,                   -- 送付依頼(手動アップロード履歴書)
  is_recommend_upload_cv BOOLEAN,                       -- 送付依頼(手動アップロード職務経歴書)
  is_recommend_upload_recommendation BOOLEAN,           -- 送付依頼(手動アップロード推薦状)
  is_recommend_generated_resume BOOLEAN,                -- 送付依頼(自動生成の履歴書)
  is_recommend_generated_cv BOOLEAN,                    -- 送付依頼(自動生成の職務経歴書)
  is_recommend_generated_recommendation BOOLEAN,        -- 送付依頼(自動生成の推薦状)
  is_recommend_generated_mask_resume BOOLEAN,           -- 送付依頼(自動生成のマスクレジュメ)
  created_at DATETIME,	                                -- 作成日時
  updated_at DATETIME,	                                -- 最終更新日時	    
  PRIMARY KEY(id),
  INDEX idx_tasks_task_id (task_id)
);

ALTER TABLE task_is_recommend_documents
    ADD CONSTRAINT fk_tasks_task_id
    FOREIGN KEY(task_id)
    REFERENCES tasks (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE task_is_recommend_documents DROP FOREIGN KEY fk_tasks_task_id;

DROP TABLE IF EXISTS task_is_recommend_documents;