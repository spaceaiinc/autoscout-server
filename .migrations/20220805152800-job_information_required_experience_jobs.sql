-- 求人企業の必要経験業界・職種を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_required_experience_jobs (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    condition_id	INT NOT NULL,	                -- 求人のID
    experience_year	INT,	                            -- 経験年数
    experience_month INT,	                            -- 経験月数
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_required_experience_jobs_condition_id (condition_id)
);

ALTER TABLE job_information_required_experience_jobs
    ADD CONSTRAINT fk_job_information_required_experience_jobs_condition_id
    FOREIGN KEY(condition_id)
    REFERENCES job_information_required_conditions (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_required_experience_jobs DROP FOREIGN KEY fk_job_information_required_experience_jobs_condition_id;

DROP TABLE IF EXISTS job_information_required_experience_jobs;