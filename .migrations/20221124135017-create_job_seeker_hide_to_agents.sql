-- 求職者の非公開エージェントを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_hide_to_agents (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                    -- 重複しないID
    job_seeker_id	INT NOT NULL,	                                -- 重複しないカラム毎の求職者ID
    agent_id INT NOT NULL,	                                    -- 重複しないカラム毎のエージェントID
    created_at	DATETIME,                                       -- 作成日時
    updated_at	DATETIME,                                       -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_hide_to_agents_job_seeker_id (job_seeker_id),
    INDEX idx_job_seeker_hide_to_agents_agent_id (agent_id)
);

ALTER TABLE job_seeker_hide_to_agents
    ADD CONSTRAINT fk_job_seeker_hide_to_agents_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE job_seeker_hide_to_agents
    ADD CONSTRAINT fk_job_seeker_hide_to_agents_agent_id
    FOREIGN KEY(agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_hide_to_agents DROP FOREIGN KEY fk_job_seeker_hide_to_agents_job_seeker_id;

DROP TABLE IF EXISTS job_seeker_hide_to_agents;
