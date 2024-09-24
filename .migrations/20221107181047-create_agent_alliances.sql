-- エージェント間のアライアンス情報を管理するテーブル
-- +migrate Up
CREATE TABLE agent_alliances (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,	                -- 重複しないID
    agent1_id INT NOT NULL,	                                -- アライアンス申請元エージェントのID
    agent2_id INT NOT NULL,	                                -- アライアンス申請元エージェントのID
    agent1_request BOOLEAN NOT NULL,                        -- エージェント1の申請状況          
    agent2_request BOOLEAN NOT NULL,                        -- エージェント2の申請状況           
    created_at DATETIME,                                    -- 作成日時
    updated_at DATETIME,                                    -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_alliances_agent1_id (agent1_id),
    INDEX idx_alliances_agent2_id (agent2_id)
);

ALTER TABLE agent_alliances
    ADD CONSTRAINT fk_alliances_agent1_id 
    FOREIGN KEY(agent1_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE agent_alliances
    ADD CONSTRAINT fk_alliances_agent2_id 
    FOREIGN KEY(agent2_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE agent_alliances DROP FOREIGN KEY fk_alliances_agent1_id ;
ALTER TABLE agent_alliances DROP FOREIGN KEY fk_alliances_agent2_id ;

DROP TABLE IF EXISTS agent_alliances;

