-- +migrate Up
CREATE TABLE IF NOT EXISTS chat_group_with_agents (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                -- 重複しないID
    uuid CHAR(36) NOT NULL UNIQUE,	                        -- 重複しないカラム毎のUUID
    agent1_id	INT NOT NULL,	                            -- エージェントid
    agent2_id	INT NOT NULL,	                            -- エージェントid
    last_send_at DATETIME,	                                -- 最終送信日時
    created_at	DATETIME,	                                -- 作成日時
    updated_at	DATETIME,	                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_chat_group_with_agents_agent1_id (agent1_id),
    INDEX idx_chat_group_with_agents_agent2_id (agent2_id),
    UNIQUE KEY u_chat_group_with_agents_agent1_id_agent2_id (agent1_id, agent2_id)
);

ALTER TABLE chat_group_with_agents
    ADD CONSTRAINT fk_chat_group_with_agents_agent1_id
    FOREIGN KEY(agent1_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE chat_group_with_agents
    ADD CONSTRAINT fk_chat_group_with_agents_agent2_id
    FOREIGN KEY(agent2_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- ALTER TABLE chat_group_with_agents
--   ADD CONSTRAINT UNIQUE u_chat_group_with_agents_agent1_id_agent2_id (agent1_id, agent2_id);

-- +migrate Down
ALTER TABLE chat_group_with_agents DROP FOREIGN KEY fk_chat_group_with_agents_agent1_id;
ALTER TABLE chat_group_with_agents DROP FOREIGN KEY fk_chat_group_with_agents_agent2_id;

DROP TABLE IF EXISTS chat_group_with_agents;