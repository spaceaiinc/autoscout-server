-- 求人の非公開エージェントを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_hide_to_agents (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                    -- 重複しないID
    job_information_id	INT NOT NULL,	                          -- 重複しないカラム毎の求人ID
    agent_id INT NOT NULL,	                                    -- 重複しないカラム毎のエージェントID
    created_at	DATETIME,                                       -- 作成日時
    updated_at	DATETIME,                                       -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_hide_to_agents_job_information_id (job_information_id),
    INDEX idx_job_information_hide_to_agents_agent_id (agent_id) 
);

ALTER TABLE job_information_hide_to_agents
    ADD CONSTRAINT fk_job_information_hide_to_agents_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE job_information_hide_to_agents
    ADD CONSTRAINT fk_job_information_hide_to_agents_agent_id
    FOREIGN KEY(agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_hide_to_agents DROP FOREIGN KEY fk_job_information_hide_to_agents_job_information_id;

DROP TABLE IF EXISTS job_information_hide_to_agents;