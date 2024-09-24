-- 送客先エージェントの求人企業の必要経験業界を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_information_required_experience_industries (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    experience_job_id	INT NOT NULL,	                    -- 必要職歴テーブルのID
    experience_industry	INT,	                          -- 経験業界
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_required_experience_industries_experience_job_id (experience_job_id)
);

ALTER TABLE sending_job_information_required_experience_industries
    ADD CONSTRAINT fk_sending_required_experience_industries_experience_job_id
    FOREIGN KEY(experience_job_id)
    REFERENCES sending_job_information_required_experience_jobs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_information_required_experience_industries DROP FOREIGN KEY fk_sending_required_experience_industries_experience_job_id;

DROP TABLE IF EXISTS sending_job_information_required_experience_industries;