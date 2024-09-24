-- 送客先エージェントの求人企業の職場の魅力情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_information_employment_statuses (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    sending_job_information_id	INT NOT NULL,	          -- 送客先エージェントの求人のID
    employment_status	INT,	                            -- 雇用形態
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_employment_statuses_sending_job_information_id (sending_job_information_id)
);

ALTER TABLE sending_job_information_employment_statuses
    ADD CONSTRAINT fk_sending_employment_statuses_sending_job_information_id
    FOREIGN KEY(sending_job_information_id)
    REFERENCES sending_job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_information_employment_statuses DROP FOREIGN KEY fk_sending_employment_statuses_sending_job_information_id;

DROP TABLE IF EXISTS sending_job_information_employment_statuses;