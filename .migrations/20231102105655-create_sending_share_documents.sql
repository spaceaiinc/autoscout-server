-- 送客先エージェントへの求職者の応募書類を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_share_documents (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    sending_job_seeker_id INT NOT NULL,	                -- 送客求職者のID
    sending_enterprise_id INT NOT NULL,	                -- 送客エージェントのID
    is_share_upload_resume BOOLEAN NOT NULL,	          -- アップロードした履歴書のシェア
    is_share_upload_cv BOOLEAN NOT NULL,	              -- アップロードした職務経歴書のシェア
    is_share_upload_recommendation BOOLEAN NOT NULL,	    -- アップロードした推薦状のシェア
    is_share_generated_resume BOOLEAN NOT NULL,	        -- 自動生成した履歴書のシェア
    is_share_generated_cv BOOLEAN NOT NULL,	            -- 自動生成した職務経歴書のシェア
    is_share_generated_recommendation BOOLEAN NOT NULL,	-- 自動生成した推薦状のシェア
    created_at	DATETIME,                               -- 作成日時
    updated_at	DATETIME,                               -- 最終更新日時
    PRIMARY KEY(id),
    UNIQUE u_sending_share_documents_seeker_and_enterprise (sending_job_seeker_id, sending_enterprise_id),
    INDEX idx_sending_share_documents_sending_job_seeker_id (sending_job_seeker_id),
    INDEX idx_sending_share_documents_sending_enterprise_id (sending_enterprise_id)
);

ALTER TABLE sending_share_documents
  ADD CONSTRAINT fk_sending_share_documents_sending_job_seeker_id
  FOREIGN KEY(sending_job_seeker_id)
  REFERENCES sending_job_seekers (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE sending_share_documents
  ADD CONSTRAINT fk_sending_share_documents_sending_enterprise_id
  FOREIGN KEY(sending_enterprise_id)
  REFERENCES sending_enterprises (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_share_documents DROP FOREIGN KEY fk_sending_share_documents_sending_job_seeker_id;
ALTER TABLE sending_share_documents DROP FOREIGN KEY fk_sending_share_documents_sending_enterprise_id;

DROP TABLE IF EXISTS sending_share_documents;