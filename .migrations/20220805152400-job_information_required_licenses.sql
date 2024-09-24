-- 求人企業の必要資格を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_required_licenses (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    condition_id INT NOT NULL,                          -- 必要条件ID
    license	INT,                                        -- 資格
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_required_licenses_condition_id (condition_id)
);

ALTER TABLE job_information_required_licenses
    ADD CONSTRAINT fk_job_information_required_licenses_job_condition_id
    FOREIGN KEY(condition_id)
    REFERENCES job_information_required_conditions (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_required_licenses DROP FOREIGN KEY fk_job_information_required_licenses_job_condition_id;

DROP TABLE IF EXISTS job_information_required_licenses;