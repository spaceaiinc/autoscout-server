-- 求人企業の必要条件を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_required_conditions (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    job_information_id	INT NOT NULL,	                -- 求人のID
    is_common BOOLEAN NOT NULL,                         -- 共通条件orパターン条件 {true: 共通条件, false: パターン}
    required_management INT,                            -- マネジメント経験
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_required_conditions_job_information_id (job_information_id)
);

ALTER TABLE job_information_required_conditions
    ADD CONSTRAINT fk_job_information_required_conditions_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_required_conditions DROP FOREIGN KEY fk_job_information_required_conditions_job_information_id;

DROP TABLE IF EXISTS job_information_required_conditions;