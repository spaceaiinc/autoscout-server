-- 送客先エージェントの求人企業の必要社会人経験を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_information_required_social_experiences (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	          -- 重複しないID
    sending_job_information_id	INT NOT NULL,	                -- 送客先エージェントの求人のID
    social_experience_type	INT,	                    -- 社会人経験タイプ（正社員, 契約社員, 業務委託 ...）
    created_at DATETIME,                              -- 作成日時
    updated_at DATETIME,                              -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_social_experiences_sending_job_information_id (sending_job_information_id)
);

ALTER TABLE sending_job_information_required_social_experiences
    ADD CONSTRAINT fk_sending_social_experiences_sending_job_information_id
    FOREIGN KEY(sending_job_information_id)
    REFERENCES sending_job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_information_required_social_experiences DROP FOREIGN KEY fk_sending_social_experiences_sending_job_information_id;

DROP TABLE IF EXISTS sending_job_information_required_social_experiences;