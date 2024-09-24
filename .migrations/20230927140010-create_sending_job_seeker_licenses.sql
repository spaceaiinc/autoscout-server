-- 送客求職者の所持資格を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_seeker_licenses (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,       -- 重複しないID
    sending_job_seeker_id	INT NOT NULL,	          -- 送客求職者ID
    license_type INT,	                            -- 資格種の種別
    acquisition_time CHAR(7) NOT NULL,	          -- 資格取得時期
    created_at DATETIME,	                        -- 作成日時
    updated_at DATETIME,	                        -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_job_seeker_license_sending_job_seeker_id (sending_job_seeker_id)
);

ALTER TABLE sending_job_seeker_licenses
    ADD CONSTRAINT fk_sending_job_seeker_licenses_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_seeker_licenses DROP FOREIGN KEY fk_sending_job_seeker_license_sending_job_seeker_id;

DROP TABLE IF EXISTS sending_job_seeker_licenses;