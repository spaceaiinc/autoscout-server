-- 送客先エージェントの求人企業の必要経験業界・職種を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_information_required_experience_jobs (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	          -- 重複しないID
    condition_id	INT NOT NULL,	                      -- 必要条件ID
    experience_year	INT,	                            -- 経験年数
    experience_month INT,	                            -- 経験月数
    created_at DATETIME,                              -- 作成日時
    updated_at DATETIME,                              -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_required_experience_jobs_condition_id (condition_id)
);

ALTER TABLE sending_job_information_required_experience_jobs
    ADD CONSTRAINT fk_sending_required_experience_jobs_condition_id
    FOREIGN KEY(condition_id)
    REFERENCES sending_job_information_required_conditions (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_information_required_experience_jobs DROP FOREIGN KEY fk_sending_required_experience_jobs_condition_id;

DROP TABLE IF EXISTS sending_job_information_required_experience_jobs;