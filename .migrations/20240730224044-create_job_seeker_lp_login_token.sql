-- 求職者の経験職種（LPからの登録時に使用）
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_lp_login_token (
  id INT AUTO_INCREMENT NOT NULL UNIQUE, -- 重複しないID
  job_seeker_id INT NOT NULL UNIQUE,     -- 求職者ID
  login_token CHAR(36) NOT NULL UNIQUE,  -- 重複しないカラム毎のUUID
  created_at DATETIME,                   -- 作成日時
  updated_at DATETIME,                   -- 最終更新日時
  PRIMARY KEY(id),
  INDEX idx_job_seeker_lp_login_token_job_seeker_id (job_seeker_id)
);

-- +migrate Down
ALTER TABLE job_seeker_lp_login_token
  ADD CONSTRAINT fk_job_seeker_lp_login_token_job_seeker_id
  FOREIGN KEY(job_seeker_id)
  REFERENCES job_seekers (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_lp_login_token DROP FOREIGN KEY fk_job_seeker_lp_login_token_job_seeker_id;

DROP TABLE IF EXISTS job_seeker_lp_login_token;
