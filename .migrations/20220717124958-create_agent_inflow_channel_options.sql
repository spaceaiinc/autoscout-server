-- エージェントの流入経路マスタを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS agent_inflow_channel_options (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,          -- 重複しないカラム毎のid
    agent_id INT NOT NULL,                          -- agentsテーブルのid
    channel_name VARCHAR(255) NOT NULL,             -- 流入経路の名前
    is_open BOOLEAN NOT NULL,                       -- Open / Close
    created_at DATETIME,                            -- 作成日時
    updated_at DATETIME,                            -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_agent_inflow_channel_optionss (agent_id)
);

ALTER TABLE agent_inflow_channel_options
    ADD CONSTRAINT fk_agent_inflow_channel_options_agent_id
    FOREIGN KEY(agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE agent_inflow_channel_options DROP FOREIGN KEY fk_agent_inflow_channel_options_agent_id;

DROP TABLE IF EXISTS agent_inflow_channel_options;