-- 選考フローのパターンを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_selection_flow_patterns (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    job_information_id INT NOT NULL,	                -- 求人のID
    public_status INT,	                                -- 公開設定(Open or Close)
    flow_title VARCHAR(255) NOT NULL,                   -- フローパターンのタイトル
    flow_pattern INT,	                                -- 選考フローのパターン
    is_deleted BOOLEAN NOT NULL,                        -- 削除状況
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_selection_flow_patterns_job_information_id (job_information_id)
);

ALTER TABLE job_information_selection_flow_patterns
    ADD CONSTRAINT fk_job_information_selection_flow_patterns_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_selection_flow_patterns DROP FOREIGN KEY fk_job_information_selection_flow_patterns_job_information_id;

DROP TABLE IF EXISTS job_information_selection_flow_patterns;
