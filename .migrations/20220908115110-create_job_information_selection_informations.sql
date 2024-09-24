-- 選考ごとの情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_selection_informations (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	        -- 重複しないID
    selection_flow_id	INT NOT NULL,	            -- 選考フローパターンのID
    selection_type INT,	                            -- 選考タイプ（書類選考 or 1次選考 ...）
    selection_point	VARCHAR(255) NOT NULL,	        -- 選考ポイント
    passed_example VARCHAR(255) NOT NULL,	        -- 合格事例
    fail_example VARCHAR(255) NOT NULL,	            -- 不合格事例
    passing_rate INT,	                            -- 通過率
    is_questionnairy BOOLEAN NOT NULL,              -- アンケートの有無
    created_at DATETIME,                            -- 作成日時
    updated_at DATETIME,                            -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_selection_informations_selection_flow_id (selection_flow_id)
);

ALTER TABLE job_information_selection_informations
    ADD CONSTRAINT fk_job_information_selection_informations_selection_flow_id
    FOREIGN KEY(selection_flow_id)
    REFERENCES job_information_selection_flow_patterns (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_selection_informations DROP FOREIGN KEY fk_job_information_selection_informations_selection_flow_id;

DROP TABLE IF EXISTS job_information_selection_informations;
