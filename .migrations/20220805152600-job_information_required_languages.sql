-- 求人企業の必要言語スキルを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_required_languages (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    condition_id INT NOT NULL,                          -- 必要条件ID
    -- language_type	INT,	                            -- 言語の種類
    language_level    INT,	                            -- 語学レベル {0:日常会話, 1:ビジネス 99:不問}
    toeic	INT,	                                    -- TOEICの点数
    toefl_ibt	INT,	                                -- TOEFL iBT点数
    toefl_pbt	INT,	                                -- TOEFL PBT点数
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_required_languages_condition_id (condition_id)
);

ALTER TABLE job_information_required_languages
    ADD CONSTRAINT fk_job_information_required_languages_condition_id
    FOREIGN KEY(condition_id)
    REFERENCES job_information_required_conditions (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_required_languages DROP FOREIGN KEY fk_job_information_required_languages_condition_id;

DROP TABLE IF EXISTS job_information_required_languages;