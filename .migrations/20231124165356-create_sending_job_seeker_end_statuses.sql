-- 送客終了理由を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_seeker_end_statuses (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,	        -- 重複しないID
    sending_job_seeker_id INT NOT NULL UNIQUE,	    -- 送客求職者のID（送客求職者に対して1レコード）
    end_reason TEXT NOT NULL,                       -- 終了理由
    end_status INT,                                 -- 終了ステータス
    created_at	DATETIME,                           -- 作成日時
    updated_at	DATETIME,                           -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_job_seeker_end_statuses_sending_job_seeker_id (sending_job_seeker_id)
);

ALTER TABLE sending_job_seeker_end_statuses
  ADD CONSTRAINT fk_sending_job_seeker_end_statuses_sending_job_seeker_id
  FOREIGN KEY(sending_job_seeker_id)
  REFERENCES sending_job_seekers (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_seeker_end_statuses DROP FOREIGN KEY fk_sending_job_seeker_end_statuses_sending_job_seeker_id;

DROP TABLE IF EXISTS sending_job_seeker_end_statuses;