-- 求職者の経験職種を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_department_histories (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	    -- 重複しないID
    work_history_id	INT NOT NULL,	            -- 職歴のID
    department VARCHAR(255) NOT NULL,	        -- マネジメント経験
    management_number INT,	                    -- マネジメント経験数
    management_detail VARCHAR(255) NOT NULL,    -- マネジメント経験詳細
    job_description TEXT NOT NULL,              -- 職務内容
    start_year	CHAR(7) NOT NULL,	            -- 開始時期
    end_year	CHAR(7) NOT NULL,	            -- 終了時期
    created_at	DATETIME,	                    -- 作成日時
    updated_at	DATETIME,	                    -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_department_histories_work_history_id (work_history_id)
);

ALTER TABLE job_seeker_department_histories
    ADD CONSTRAINT fk_job_seeker_department_histories_work_history_id
    FOREIGN KEY(work_history_id)
    REFERENCES job_seeker_work_histories (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_department_histories DROP FOREIGN KEY fk_job_seeker_department_histories_work_history_id;

DROP TABLE IF EXISTS job_seeker_department_histories;