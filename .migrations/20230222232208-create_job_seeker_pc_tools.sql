-- 求職者の経験業務ツールを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_pc_tools (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    job_seeker_id	INT NOT NULL,	                        -- 求職者ID
    tool	INT,	                                        -- 業務ツール
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_pc_tools_job_seeker_id (job_seeker_id)
);

ALTER TABLE job_seeker_pc_tools
    ADD CONSTRAINT fk_job_seeker_pc_tools_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_pc_tools DROP FOREIGN KEY fk_job_seeker_pc_tools_job_seeker_id

DROP TABLE IF EXISTS job_seeker_pc_tools;