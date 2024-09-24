-- 送客先エージェントの求人企業の必要条件を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_information_required_conditions (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    sending_job_information_id	INT NOT NULL,	        -- 送客先エージェントの求人のID
    is_common BOOLEAN NOT NULL,                         -- 共通条件orパターン条件 {true: 共通条件, false: パターン}
    required_management INT,                            -- マネジメント経験
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_required_conditions_sending_job_information_id (sending_job_information_id)
);

ALTER TABLE sending_job_information_required_conditions
    ADD CONSTRAINT fk_sending_required_conditions_sending_job_information_id
    FOREIGN KEY(sending_job_information_id)
    REFERENCES sending_job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_information_required_conditions DROP FOREIGN KEY fk_sending_required_conditions_sending_job_information_id;

DROP TABLE IF EXISTS sending_job_information_required_conditions;