-- 求職者の面談前アンケートを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS initial_questionnaires (
  id	INT AUTO_INCREMENT NOT NULL UNIQUE,             -- 重複しないID
  uuid VARCHAR(255) NOT NULL UNIQUE,	                -- 重複しないUUID
  job_seeker_id	INT NOT NULL,	                        -- 求職者のID
  question	TEXT NOT NULL,	                          -- 質問や要望
  created_at	DATETIME,	                              -- 作成日時
  updated_at	DATETIME,	                              -- 最終更新日時
  PRIMARY KEY(id),
  INDEX idx_initial_questionnaires_job_seeker_id (job_seeker_id)
);

ALTER TABLE initial_questionnaires
    ADD CONSTRAINT fk_initial_questionnaires_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE initial_questionnaires DROP FOREIGN KEY fk_initial_questionnaires_job_seeker_id;

DROP TABLE IF EXISTS initial_questionnaires;
