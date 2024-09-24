-- 送客応諾後の終了理由を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_phase_end_statuses (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,
    sending_phase_id INT NOT NULL UNIQUE,           -- 送客フェーズのID（送客フェーズに対して1レコード）
    end_reason TEXT NOT NULL,                       -- 終了理由
    end_status INT,                                 -- 終了ステータス
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY(id),
    INDEX idx_sending_phase_end_statuses_sending_phase_id (sending_phase_id)
);

ALTER TABLE sending_phase_end_statuses
  ADD CONSTRAINT fk_sending_phase_end_statuses_sending_phase_id
  FOREIGN KEY(sending_phase_id)
  REFERENCES sending_phases (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_phase_end_statuses DROP FOREIGN KEY fk_sending_phase_end_statuses_sending_phase_id;

DROP TABLE IF EXISTS sending_phase_end_statuses;