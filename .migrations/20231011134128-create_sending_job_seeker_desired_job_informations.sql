-- 送客求職者の希望求人を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_seeker_desired_job_informations (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    sending_job_seeker_id INT NOT NULL,	                -- 送客求職者ID
    sending_job_information_id INT NOT NULL,            -- 送客求人ID
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    UNIQUE u_sending_seeker_and_information (sending_job_seeker_id, sending_job_information_id),
    INDEX idx_sending_desired_job_informations_sending_job_seeker_id (sending_job_seeker_id),
    INDEX idx_sending_desired_job_informations_sending_job_information_id (sending_job_information_id)
);

ALTER TABLE sending_job_seeker_desired_job_informations
    ADD CONSTRAINT fk_sending_desired_job_informations_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE sending_job_seeker_desired_job_informations
    ADD CONSTRAINT fk_sending_desired_job_informations_sending_job_information_id
    FOREIGN KEY(sending_job_information_id)
    REFERENCES sending_job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_seeker_desired_job_informations DROP FOREIGN KEY fk_sending_desired_job_informations_sending_job_seeker_id
ALTER TABLE sending_job_seeker_desired_job_informations DROP FOREIGN KEY fk_sending_desired_job_informations_sending_job_information_id

DROP TABLE IF EXISTS sending_job_seeker_desired_job_informations;