-- 求人企業の必要経験職種を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_required_experience_occupations (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    experience_job_id	INT NOT NULL,	                    -- 必要職歴テーブルのID
    experience_occupation	INT,	                        -- 経験職種
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_required_experience_occupations (experience_job_id)
);

ALTER TABLE job_information_required_experience_occupations
    ADD CONSTRAINT fk_job_information_required_experience_occupations
    FOREIGN KEY(experience_job_id)
    REFERENCES job_information_required_experience_jobs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_required_experience_occupations DROP FOREIGN KEY fk_job_information_required_experience_occupations;

DROP TABLE IF EXISTS job_information_required_experience_occupations;