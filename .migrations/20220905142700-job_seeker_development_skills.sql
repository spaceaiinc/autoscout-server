-- 求職者の希望勤務地を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_development_skills (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    job_seeker_id	INT NOT NULL,	                    -- 求職者のID
    development_category	INT,	                    -- カテゴリ(言語 or OS)
    development_type	INT,	                        -- タイプ（言語の種類 or OSの種類）
    experience_year	INT,	                            -- 経験年数
    experience_month	INT,	                        -- 経験月数
    created_at	DATETIME,                               -- 作成日時
    updated_at	DATETIME,                               -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_development_skills_job_seeker_id (job_seeker_id)
);

ALTER TABLE job_seeker_development_skills
    ADD CONSTRAINT fk_job_seeker_development_skills_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_development_skills DROP FOREIGN KEY fk_job_seeker_development_skills_job_seeker_id;

DROP TABLE IF EXISTS job_seeker_development_skills;