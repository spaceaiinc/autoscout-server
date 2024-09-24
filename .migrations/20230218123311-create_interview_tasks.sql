-- 面談調整タスクのテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS interview_tasks (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,          -- 重複しないカラム毎のid
    interview_task_group_id	INT NOT NULL,	        -- 面談調整タスクグループのID
    agent_staff_id INT,	                            -- タスク依頼者のID
    ca_staff_id INT,                                -- タスク実行時の求職者の担当CAのID
    phase_category INT NOT NULL,	                -- フェーズ（大分類）
    phase_sub_category INT NOT NULL,   	            -- フェーズ（中分類）
    remarks	TEXT NOT NULL,	                        -- 依頼内容
    deadline_day CHAR(10) NOT NULL,                 -- 期限（日付）
    deadline_time INT,	                            -- 期限（時間）
    select_action_label VARCHAR(255) NOT NULL,      -- 面談調整タスクで選択したアクションのラベル
    created_at DATETIME,                            -- 作成日時
    updated_at DATETIME,                            -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_interview_tasks_interview_task_group_id (interview_task_group_id),
    INDEX idx_interview_tasks_agent_staff_id (agent_staff_id),
    INDEX idx_interview_tasks_ca_staff_id (ca_staff_id)
);

ALTER TABLE interview_tasks 
    ADD CONSTRAINT fk_interview_tasks_interview_task_group_id
    FOREIGN KEY(interview_task_group_id)
    REFERENCES interview_task_groups (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE interview_tasks
    ADD CONSTRAINT fk_interview_tasks_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE interview_tasks
    ADD CONSTRAINT fk_interview_tasks_ca_staff_id
    FOREIGN KEY(ca_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE interview_tasks DROP FOREIGN KEY fk_interview_tasks_interview_task_group_id;
ALTER TABLE interview_tasks DROP FOREIGN KEY fk_interview_tasks_agent_staff_id;
ALTER TABLE interview_tasks DROP FOREIGN KEY fk_interview_tasks_ca_staff_id;

DROP TABLE IF EXISTS interview_tasks;
