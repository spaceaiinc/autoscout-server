-- 送客求職者の希望勤務地を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_seeker_desired_work_locations (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    sending_job_seeker_id	INT NOT NULL,	                -- 送客求職者のID
    desired_work_location INT,	                        -- 希望勤務地
    desired_rank    INT,	                              -- 希望順位
    created_at	DATETIME,                               -- 作成日時
    updated_at	DATETIME,                               -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_desired_work_locations_sending_job_seeker_id (sending_job_seeker_id)
);

ALTER TABLE sending_job_seeker_desired_work_locations
    ADD CONSTRAINT fk_sending_desired_work_locations_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_seeker_desired_work_locations DROP FOREIGN KEY fk_sending_desired_work_locations_sending_job_seeker_id;

DROP TABLE IF EXISTS sending_job_seeker_desired_work_locations;