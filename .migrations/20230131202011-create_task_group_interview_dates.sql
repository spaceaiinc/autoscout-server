-- +migrate Up
CREATE TABLE IF NOT EXISTS task_group_interview_dates (
  id	           INT AUTO_INCREMENT NOT NULL UNIQUE,	-- 重複しないカラム毎のid
  task_group_id	 INT NOT NULL,	                      -- タスクグループのid
  phase_category INT NOT NULL,	                                -- 選考フェーズ（2: 一次, 3: 二次, 4: 三次, 5: 四次, 6: 五次, 7: 最終）
  interview_date VARCHAR(255) NOT NULL,	              -- 確定日時（面談日時）
  created_at	DATETIME,	                              -- 作成日時
  updated_at	DATETIME,	                              -- 最終更新日時
  PRIMARY KEY(id),
  INDEX idx_task_group_interview_dates_task_group_id (task_group_id),
  UNIQUE u_task_group_id_and_phase_category (task_group_id, phase_category)
);

-- +migrate Down
ALTER TABLE task_group_interview_dates
  ADD CONSTRAINT fk_task_group_interview_dates_task_group_id
  FOREIGN KEY (task_group_id) 
  REFERENCES task_groups(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE task_group_interview_dates
  ADD CONSTRAINT UNIQUE u_task_group_id_and_phase_category (task_group_id, phase_category);

-- +migrate Down
ALTER TABLE task_group_interview_dates DROP FOREIGN KEY fk_task_group_interview_dates_task_group_id

DROP TABLE IF EXISTS task_group_interview_dates;