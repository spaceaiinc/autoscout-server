-- 求職者のスケジュール管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_schedules (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	-- 重複しないID
    job_seeker_id	INT NOT NULL,	            -- 求人申請ID
    task_id	INT,	                          -- タスクID（レコード作成時はnullになるため、nullを許容する）
    schedule_type INT,                      -- 種別（0: 求職者調整, 1: CA調整, 2: 企業調整(RA), 3: リスケ）
    title VARCHAR(255) NOT NULL,            -- スケジュールのタイトル
    start_time VARCHAR(255) NOT NULL,       -- スケジュールの開始時間（形式: 2023-04-01T12:00）
    end_time VARCHAR(255) NOT NULL,         -- スケジュールの終了時間（形式: 2023-04-01T12:00）
    seeker_description text NOT NULL,       -- 求職者の日程補足
    staff_description  text NOT NULL,       -- 担当者の日程補足
    is_share BOOLEAN,                       -- CA担当者に共有済みか？
    repetition_count INT,                   -- 選考の繰り返し回数（schedule_typeが2の時のみ保存される）
    created_at DATETIME,                    -- 作成日時
    updated_at DATETIME,                    -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_schedules_job_seeker_id (job_seeker_id),
    INDEX idx_job_seeker_schedules_task_id (task_id),
    INDEX idx_job_seeker_schedules_schedule_type (schedule_type)
  );

ALTER TABLE job_seeker_schedules
  ADD CONSTRAINT fk_job_seeker_schedules_job_seeker_id
  FOREIGN KEY(job_seeker_id)
  REFERENCES job_seekers (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE job_seeker_schedules
  ADD CONSTRAINT fk_job_seeker_schedules_task_id
  FOREIGN KEY(task_id)
  REFERENCES tasks (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_schedules DROP FOREIGN KEY fk_job_seeker_schedules_job_seeker_id
ALTER TABLE job_seeker_schedules DROP FOREIGN KEY fk_job_seeker_schedules_job_task_id

DROP TABLE IF EXISTS job_seeker_schedules;


