-- 求職者の経験職種（LPからの登録時に使用）
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_experience_jobs (
  id INT AUTO_INCREMENT NOT NULL UNIQUE, -- 重複しないID
  job_seeker_id INT NOT NULL,            -- 部署履歴ID
  occupation INT,	                       -- 職種
  experience_year INT,                   -- 経験年数（0: 1年未満, 1: 1年以上, 2: 2年以上, ..., 10: 10年以上）
  created_at DATETIME,                   -- 作成日時
  updated_at DATETIME,                   -- 最終更新日時
  PRIMARY KEY(id),
  INDEX idx_job_seeker_experience_jobs_job_seeker_id (job_seeker_id)
);

-- +migrate Down
ALTER TABLE job_seeker_experience_jobs
  ADD CONSTRAINT fk_job_seeker_experience_jobs_job_seeker_id
  FOREIGN KEY(job_seeker_id)
  REFERENCES job_seekers (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_experience_jobs DROP FOREIGN KEY fk_job_seeker_experience_jobs_job_seeker_id;

DROP TABLE IF EXISTS job_seeker_experience_jobs;
