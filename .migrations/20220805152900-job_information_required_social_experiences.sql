-- 求人企業の必要社会人経験を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_required_social_experiences (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    job_information_id	INT NOT NULL,	                -- 求人のID
    social_experience_type	INT,	                    -- 社会人経験タイプ（正社員, 契約社員, 業務委託 ...）
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_required_social_experiences (job_information_id)
);

ALTER TABLE job_information_required_social_experiences
    ADD CONSTRAINT fk_job_information_required_social_experiences
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_required_social_experiences DROP FOREIGN KEY fk_job_information_required_social_experiences;

DROP TABLE IF EXISTS job_information_required_social_experiences;