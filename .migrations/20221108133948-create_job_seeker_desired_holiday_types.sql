-- 求職者の希望休日タイプを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_desired_holiday_types (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                -- 重複しないID
    job_seeker_id	INT NOT NULL,	                        -- 求職者のID
    holiday_type INT,	                                    -- 休日タイプ
    created_at	DATETIME,                                   -- 作成日時
    updated_at	DATETIME,                                   -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_desired_holiday_types_job_seeker_id (job_seeker_id)
);

ALTER TABLE job_seeker_desired_holiday_types
    ADD CONSTRAINT fk_job_seeker_desired_holiday_types_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_desired_holiday_types DROP FOREIGN KEY fk_job_seeker_desired_holiday_types_job_seeker_id;

DROP TABLE IF EXISTS job_seeker_desired_holiday_types;
