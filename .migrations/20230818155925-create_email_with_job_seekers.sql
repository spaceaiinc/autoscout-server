-- 求職者とのメールを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS email_with_job_seekers (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	-- 重複しないID
    job_seeker_id	INT NOT NULL,	            -- 求職者ID
    subject VARCHAR(255) NOT NULL,          -- メールの件名
    content TEXT NOT NULL,                  -- メールの本文
    file_name TEXT NOT NULL,                -- 添付ファイル名(複数の場合は"|--|--T--|--|"区切りで保存)
    created_at DATETIME,                    -- 作成日時
    updated_at DATETIME,                    -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_email_with_job_seekers_job_seeker_id (job_seeker_id)
  );

ALTER TABLE email_with_job_seekers
  ADD CONSTRAINT fk_email_with_job_seekers_job_seeker_id
  FOREIGN KEY(job_seeker_id)
  REFERENCES job_seekers (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE email_with_job_seekers DROP FOREIGN KEY fk_email_with_job_seekers_job_seeker_id;

DROP TABLE IF EXISTS email_with_job_seekers;


