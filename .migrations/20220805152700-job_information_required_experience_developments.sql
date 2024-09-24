-- 求人企業の必要開発経験を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_required_experience_developments (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    condition_id INT NOT NULL,                          -- 必要条件ID
    development_category	INT,	                    -- 開発経験カテゴリ
    -- development_type	INT,	                        -- タイプ
    experience_year	INT,	                            -- 経験年数
    experience_month	INT,	                        -- 経験月数
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_required_developments_condition_id (condition_id)
);

ALTER TABLE job_information_required_experience_developments
    ADD CONSTRAINT fk_job_information_required_experience_developments_condition_id
    FOREIGN KEY(condition_id)
    REFERENCES job_information_required_conditions (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_required_experience_developments DROP FOREIGN KEY fk_job_information_required_experience_developments_condition_id;

DROP TABLE IF EXISTS job_information_required_experience_developments;