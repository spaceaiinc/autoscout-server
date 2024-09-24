-- 送客求職者の学歴情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_seeker_student_histories (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,                 -- 重複しないID
    sending_job_seeker_id	INT NOT NULL,	                        -- 求職者ID
    school_category	INT,	                                -- 学校カテゴリ
    school_name	VARCHAR(255) NOT NULL,	                    -- 学校名
    school_level	INT,	                                -- 学校レベル
    subject	VARCHAR(255) NOT NULL,	                        -- 学部・学科・コース
    entrance_year	CHAR(7) NOT NULL,	                    -- 入学年月
    first_status INT,                                       -- 開始ステータス（入学, 編集学, 転入学）
    graduation_year	CHAR(7) NOT NULL,	                    -- 卒業年月
    last_status	INT,	                                    -- 終了ステータス（卒業, 中退、卒業見込み）
    created_at	DATETIME,	                                -- 作成日時
    updated_at	DATETIME,	                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_job_seeker_student_histories_sending_job_seeker_id (sending_job_seeker_id)
);

ALTER TABLE sending_job_seeker_student_histories
    ADD CONSTRAINT fk_sending_job_seeker_student_histories_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_seeker_student_histories DROP FOREIGN KEY fk_sending_job_seeker_student_histories_sending_job_seeker_id;

DROP TABLE IF EXISTS sending_job_seeker_student_histories;