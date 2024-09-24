-- 求職者の興味あり求人情報
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_interested_job_listings (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	    -- 重複しないID
    uuid CHAR(36) NOT NULL UNIQUE,	            -- 重複しないカラム毎のUUID
    job_seeker_id INT NOT NULL,	                -- 求職者のID
    job_information_id INT NOT NULL,	        -- 求人のID
    interested_type INT,	                    -- 興味ありの種別(0: エントリー希望, 1: 興味あり)
    created_at DATETIME,	                    -- 作成日時
    updated_at DATETIME,	                    -- 最終更新日時
    PRIMARY KEY(id),
    UNIQUE u_jsijl_seeker_and_information (job_seeker_id, job_information_id),
    INDEX idx_jsijl_job_seeker_id (job_seeker_id),
    INDEX idx_jsijl_job_information_id (job_information_id),
    INDEX idx_jsijl_interested_type (interested_type)
);

ALTER TABLE job_seeker_interested_job_listings
    ADD CONSTRAINT fk_job_seeker_interested_job_listings_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE job_seeker_interested_job_listings
    ADD CONSTRAINT fk_job_seeker_interested_job_listings_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_interested_job_listings DROP FOREIGN KEY fk_job_seeker_interested_job_listings_job_seeker_id;
ALTER TABLE job_seeker_interested_job_listings DROP FOREIGN KEY fk_job_seeker_interested_job_listings_job_information_id;

DROP TABLE IF EXISTS job_seeker_interested_job_listings;