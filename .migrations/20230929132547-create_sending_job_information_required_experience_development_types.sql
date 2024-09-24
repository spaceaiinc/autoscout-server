-- 送客先エージェントの求人企業の必要開発言語・OSのタイプを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_information_required_experience_development_types (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    experience_development_id	INT NOT NULL,	            -- 必要開発経験のID
    development_type  INT,	                            -- タイプ
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_development_types_experience_development_id (experience_development_id)
);

ALTER TABLE sending_job_information_required_experience_development_types
    ADD CONSTRAINT fk_sending_development_types_experience_development_id
    FOREIGN KEY(experience_development_id)
    REFERENCES sending_job_information_required_experience_developments (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_information_required_experience_development_types DROP FOREIGN KEY fk_sending_development_types_experience_development_id;

DROP TABLE IF EXISTS sending_job_information_required_experience_development_types;