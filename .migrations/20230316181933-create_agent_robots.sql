-- エージェントのロボットを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS agent_robots (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    uuid	CHAR(36) NOT NULL UNIQUE,	                -- 重複しないカラム毎のUUID
    agent_id	INT NOT NULL,	                        -- エージェントID
    name VARCHAR(255) NOT NULL,	                        -- ロボット名
    is_entry_active BOOLEAN NOT NULL DEFAULT FALSE,	    -- エントリーオプションが有効かどうか/false:走らせない true:走る(共通)
    is_scout_active BOOLEAN NOT NULL DEFAULT FALSE,	    -- スカウトオプションが有効かどうか/false:走らせない true:走る(共通)
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_agent_robots_agent_id (agent_id)
);

ALTER TABLE agent_robots
    ADD CONSTRAINT fk_agent_robots_agent_id
    FOREIGN KEY(agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE agent_robots DROP FOREIGN KEY fk_agent_robots_agent_id;

DROP TABLE IF EXISTS agent_robots;
