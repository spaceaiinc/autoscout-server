-- 求職者の経験業界を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_experience_industries (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,                 -- 重複しないID
    work_history_id	INT NOT NULL,	                        -- 職歴のID
    industry	INT,                                        -- 業界  
    created_at	DATETIME,                                   -- 作成日時
    updated_at	DATETIME,                                   -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_experience_industries_work_history_id (work_history_id)
);

ALTER TABLE job_seeker_experience_industries
    ADD CONSTRAINT fk_job_seeker_experience_industries_work_history_id
    FOREIGN KEY(work_history_id)
    REFERENCES job_seeker_work_histories (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_experience_industries DROP FOREIGN KEY fk_job_seeker_experience_industries_work_history_id;

DROP TABLE IF EXISTS job_seeker_experience_industries;