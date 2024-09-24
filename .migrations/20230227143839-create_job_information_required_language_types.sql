-- 求人企業の必要開発言語・OSのタイプを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_required_language_types (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    language_id	INT NOT NULL,	                -- 必要開発経験のID
    language_type  INT,	                        -- タイプ
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_required_language_types_language_id (language_id)
);

ALTER TABLE job_information_required_language_types
    ADD CONSTRAINT fk_job_information_required_language_types_language_id
    FOREIGN KEY(language_id)
    REFERENCES job_information_required_languages (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_required_language_types DROP FOREIGN KEY fk_job_information_required_language_types_language_id;

DROP TABLE IF EXISTS job_information_required_language_types;