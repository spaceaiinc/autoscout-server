-- +migrate Up
CREATE TABLE IF NOT EXISTS chat_thread_with_agents (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                  -- 重複しないカラム毎のid
    uuid CHAR(36) NOT NULL UNIQUE,	                          -- 重複しないカラム毎のUUID
    group_id	    INT NOT NULL,	                              -- エージェントチャットのグループID
    agent_staff_id	INT NOT NULL,	                            -- エージェントスタッフのid
    title	VARCHAR(255) NOT NULL,	                            -- メッセージの内容
    created_at	    DATETIME,	                                -- 作成日時
    updated_at	    DATETIME,	                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_chat_thread_with_agents_group_id (group_id),
    INDEX idx_chat_thread_with_agents_agent_staff_id (agent_staff_id)
);

ALTER TABLE chat_thread_with_agents
    ADD CONSTRAINT fk_chat_thread_with_agents_group_id
    FOREIGN KEY(group_id)
    REFERENCES chat_group_with_agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE chat_thread_with_agents
    ADD CONSTRAINT fk_chat_thread_with_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE chat_thread_with_agents DROP FOREIGN KEY fk_chat_thread_with_agents_group_id;
ALTER TABLE chat_thread_with_agents DROP FOREIGN KEY fk_chat_thread_with_agents_agent_staff_id;

DROP TABLE IF EXISTS chat_thread_with_agents;