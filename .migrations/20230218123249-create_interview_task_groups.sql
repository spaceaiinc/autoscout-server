-- 面談調整タスクグループのテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS interview_task_groups (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,                                -- 重複しないカラム毎のid
    uuid CHAR(36) NOT NULL UNIQUE,	                                      -- 重複しないカラム毎のUUID
    agent_id INT NOT NULL,	                                              -- エージェントのID
    job_seeker_id INT NOT NULL,	                                          -- 求職者のID
    interview_date DATETIME NOT NULL DEFAULT "0001-01-01 01:00:00",       -- 面談日時
    first_interview_date DATETIME NOT NULL DEFAULT "0001-01-01 01:00:00", -- 初回面談日時(KPIの計算の基準になる)
    last_request_at DATETIME,                                             -- 最終タスク依頼時間
    last_watched_at DATETIME,                                             -- 最終タスク閲覧時間
    created_at DATETIME,                                                  -- 作成日時
    updated_at DATETIME,                                                  -- 最終更新日時
    PRIMARY KEY(id),
    UNIQUE u_agent_and_seeker (agent_id, job_seeker_id),
    INDEX idx_interview_task_groups_agent_id (agent_id),
    INDEX idx_interview_task_groups_job_seeker_id (job_seeker_id),
    INDEX idx_interview_task_groups_interview_date (interview_date),
    INDEX idx_interview_task_groups_first_interview_date (first_interview_date)
);

ALTER TABLE interview_task_groups
    ADD CONSTRAINT fk_interview_task_groups_agent_id
    FOREIGN KEY(agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE interview_task_groups
    ADD CONSTRAINT fk_interview_task_groups_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE interview_task_groups DROP FOREIGN KEY fk_interview_task_groups_agent_id;
ALTER TABLE interview_task_groups DROP FOREIGN KEY fk_interview_task_groups_job_seeker_id;

DROP TABLE IF EXISTS interview_task_groups;