-- 求職者の所持資格を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_licenses (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,                 -- 重複しないID
    job_seeker_id	INT NOT NULL,	                        -- 求職者ID
    license_type	INT,	                                -- 資格種の種別
    acquisition_time    CHAR(7) NOT NULL,	                -- 資格取得時期
    created_at	DATETIME,	                                -- 作成日時
    updated_at	DATETIME,	                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_license_job_seeker_id (job_seeker_id)
);

ALTER TABLE job_seeker_licenses
    ADD CONSTRAINT fk_job_seeker_licenses_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_licenses DROP FOREIGN KEY fk_job_seeker_license_job_seeker_id;

DROP TABLE IF EXISTS job_seeker_licenses;