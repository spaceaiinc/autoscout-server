-- 送客先エージェントの求人企業の必要語学のタイプを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_information_required_language_types (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    language_id	INT NOT NULL,	                        -- 必要語学経験のID
    language_type  INT,	                                -- タイプ
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_required_language_types_language_id (language_id)
);

ALTER TABLE sending_job_information_required_language_types
    ADD CONSTRAINT fk_sending_required_language_types_language_id
    FOREIGN KEY(language_id)
    REFERENCES sending_job_information_required_languages (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_information_required_language_types DROP FOREIGN KEY fk_sending_required_language_types_language_id;

DROP TABLE IF EXISTS sending_job_information_required_language_types;