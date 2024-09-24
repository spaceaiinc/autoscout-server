-- 送客求職者の経験業務ツールを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_seeker_pc_tools (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    sending_job_seeker_id	INT NOT NULL,	                        -- 送客求職者ID
    tool	INT,	                                        -- 業務ツール
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_job_seeker_pc_tools_sending_job_seeker_id (sending_job_seeker_id)
);

ALTER TABLE sending_job_seeker_pc_tools
    ADD CONSTRAINT fk_sending_job_seeker_pc_tools_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_seeker_pc_tools DROP FOREIGN KEY fk_sending_job_seeker_pc_tools_sending_job_seeker_id;

DROP TABLE IF EXISTS sending_job_seeker_pc_tools;