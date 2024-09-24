-- タスクグループを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS tasks (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                -- 重複しないID
    task_group_id INT NOT NULL,	                            -- タスクグループID
    phase_category INT NOT NULL,	                        -- フェーズ（大分類）
    phase_sub_category INT NOT NULL,   	                    -- フェーズ（中分類）
    staff_type INT,                                         -- タスクの担当者種別（0: CA, 1: RA, 2: CA上長, 3: RA上長）
    executed_staff_id INT NOT NULL,                         -- タスクを実行した担当者ID
    remarks	TEXT NOT NULL,	                                -- 依頼内容
    deadline_day CHAR(10) NOT NULL,                         -- 期限（日付）
    deadline_time INT,	                                    -- 期限（時間）
    talk_about_in_interview TEXT NOT NULL,                  -- オファー面談で話してほしいこと
    schedule_collection_condition TEXT NOT NULL,            -- 日程回収条件（曜日・時間・形式・候補日の数）
    exam_guide_content TEXT NOT NULL,                       -- 案内内容（適性検査や適性検査の受験案内・受験期日など）
    is_check_double_sided BOOLEAN NOT NULL,                 -- 両面タスクのチェックをつけたタスクかどうか？
    created_at DATETIME,	                                -- 作成日時
    updated_at DATETIME,	                                -- 最終更新日時	    
    PRIMARY KEY(id),
    INDEX idx_tasks_task_group_id (task_group_id),
    INDEX idx_tasks_executed_staff_id (executed_staff_id)
);

ALTER TABLE tasks
    ADD CONSTRAINT fk_tasks_task_group_id
    FOREIGN KEY(task_group_id)
    REFERENCES task_groups (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE tasks
    ADD CONSTRAINT fk_tasks_executed_staff_id
    FOREIGN KEY(executed_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE tasks DROP FOREIGN KEY fk_tasks_task_group_id;
ALTER TABLE tasks DROP FOREIGN KEY fk_tasks_executed_staff_id;

DROP TABLE IF EXISTS tasks;