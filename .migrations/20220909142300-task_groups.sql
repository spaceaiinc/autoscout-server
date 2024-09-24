-- タスクグループを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS task_groups (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    uuid CHAR(36) NOT NULL UNIQUE,	                    -- 重複しないカラム毎のUUID
    job_seeker_id	INT NOT NULL,	                    -- 求職者のID
    job_information_id	INT NOT NULL,	                -- 求人のID
    selection_flow_pattern_id INT,                      -- 選考フローパターンのID
    ra_last_request_at DATETIME,                        -- RAの最終タスク依頼時間
    ra_last_watched_at DATETIME,                        -- RAの最終タスク閲覧時間
    ca_last_request_at DATETIME,                        -- CAの最終タスク依頼時間
    ca_last_watched_at DATETIME,                        -- CAの最終タスク閲覧時間
    joining_date CHAR(10) NOT NULL,                     -- 入社日 形式：2020-12-31
    is_double_sided BOOLEAN NOT NULL,                   -- 両面タスクの有無
    created_at	DATETIME,	                            -- 作成日時
    updated_at	DATETIME,	                            -- 最終更新日時
    PRIMARY KEY(id),
    UNIQUE u_seeker_and_information (job_seeker_id, job_information_id),
    INDEX idx_task_groups_job_seeker_id (job_seeker_id),
    INDEX idx_task_groups_job_information_id (job_information_id),
    INDEX idx_task_groups_selection_flow_pattern_id (selection_flow_pattern_id)
);

ALTER TABLE task_groups
    ADD CONSTRAINT fk_task_groups_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE task_groups
    ADD CONSTRAINT fk_task_groups_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE task_groups
    ADD CONSTRAINT fk_task_groups_selection_flow_pattern_id
    FOREIGN KEY(selection_flow_pattern_id)
    REFERENCES job_information_selection_flow_patterns (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE task_groups DROP FOREIGN KEY fk_task_groups_job_seeker_id;
ALTER TABLE task_groups DROP FOREIGN KEY fk_task_groups_job_information_id;
ALTER TABLE task_groups DROP FOREIGN KEY fk_task_groups_selection_flow_pattern_id;

DROP TABLE IF EXISTS task_groups;