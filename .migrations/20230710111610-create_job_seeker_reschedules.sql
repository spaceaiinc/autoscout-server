-- 求職者のリスケ情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_reschedules (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	-- 重複しないID
    reschedule_id	INT NOT NULL,	            -- job_seeker_schedulesテーブルのID
    task_id	INT NOT NULL,	                  -- リスケを実行したタスクID
    created_at DATETIME,                    -- 作成日時
    updated_at DATETIME,                    -- 最終更新日時
    PRIMARY KEY(id),
    UNIQUE u_reschedule_id_and_task_id (reschedule_id, task_id),
    INDEX idx_job_seeker_reschedules_reschedule_id (reschedule_id),
    INDEX idx_job_seeker_reschedules_task_id (task_id)  
  );

ALTER TABLE job_seeker_reschedules
  ADD CONSTRAINT fk_job_seeker_reschedules_reschedule_id
  FOREIGN KEY(reschedule_id)
  REFERENCES job_seeker_schedules (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE job_seeker_reschedules
  ADD CONSTRAINT fk_job_seeker_reschedules_task_id
  FOREIGN KEY(task_id)
  REFERENCES tasks (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_reschedules DROP FOREIGN KEY fk_job_seeker_reschedules_reschedule_id
ALTER TABLE job_seeker_reschedules DROP FOREIGN KEY fk_job_seeker_reschedules_job_task_id

DROP TABLE IF EXISTS job_seeker_reschedules;


