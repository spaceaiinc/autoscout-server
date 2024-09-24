-- 送客求職者の職歴を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_seeker_work_histories (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,                 -- 重複しないID
    sending_job_seeker_id	INT NOT NULL,	                        -- 送客求職者ID
    company_name	VARCHAR(255) NOT NULL,                  -- 会社名
    employee_number_single	INT,	                        -- 従業員数（単体）
    employee_number_group	INT,	                        -- 従業員数（連結）
    public_offering	INT,	                                -- 株式公開
    joining_year	CHAR(7) NOT NULL,	                    -- 入社年月
    employment_status	INT,	                            -- 雇用形態
    retire_reason_of_truth	TEXT NOT NULL,	                -- 退職理由（本音）
    retire_reason_of_public	TEXT NOT NULL,	                -- 退職理由（建前）
    retire_year	CHAR(7) NOT NULL,	                        -- 退職年月
    first_status INT,                                       -- 開始ステータス（入社, 入行, 入局など）
    last_status INT,                                        -- 終了ステータス（一身上の都合により退職, 派遣期間満了につき退職など）
    created_at	DATETIME,	                                -- 作成日時
    updated_at	DATETIME,	                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_job_seeker_work_histories_sending_job_seeker_id (sending_job_seeker_id)
);

ALTER TABLE sending_job_seeker_work_histories
    ADD CONSTRAINT fk_sending_job_seeker_work_histories_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_seeker_work_histories DROP FOREIGN KEY fk_sending_job_seeker_work_histories_sending_job_seeker_id;

DROP TABLE IF EXISTS sending_job_seeker_work_histories;