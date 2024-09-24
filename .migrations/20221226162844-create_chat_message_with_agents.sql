-- +migrate Up
CREATE TABLE IF NOT EXISTS chat_message_with_agents (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                    -- 重複しないカラム毎のid
    uuid CHAR(36) NOT NULL UNIQUE,	                            -- 重複しないカラム毎のUUID
    thread_id	    INT NOT NULL,	                            -- エージェントチャットのグループID
    agent_staff_id	INT NOT NULL,	                            -- エージェントスタッフのid
    message	        TEXT,	                                    -- メッセージの内容
    file_url        varchar(2000) NOT NULL,                     -- url
    created_at	    DATETIME,	                                -- 作成日時
    updated_at	    DATETIME,	                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_chat_message_with_agents_thread_id (thread_id),
    INDEX idx_chat_message_with_agents_agent_staff_id (agent_staff_id)
);

ALTER TABLE chat_message_with_agents
    ADD CONSTRAINT fk_chat_message_with_agents_thread_id
    FOREIGN KEY(thread_id)
    REFERENCES chat_thread_with_agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE chat_message_with_agents
    ADD CONSTRAINT fk_chat_message_with_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE chat_message_with_agents DROP FOREIGN KEY fk_chat_message_with_agents_thread_id;
ALTER TABLE chat_message_with_agents DROP FOREIGN KEY fk_chat_message_with_agents_agent_staff_id;

DROP TABLE IF EXISTS chat_message_with_agents;