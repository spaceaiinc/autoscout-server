-- エージェントの売上管理の決算月や公開状況を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS agent_sale_managements (
  id INT AUTO_INCREMENT NOT NULL UNIQUE,      -- 重複しないID
  agent_id INT NOT NULL,                      -- agentID
  fiscal_year DATE,                           -- 決算年月
  is_open BOOLEAN,                            -- Open・Close
  created_at DATETIME,                        -- 作成日時
  updated_at DATETIME,                        -- 最終更新日時
  PRIMARY KEY(id),
  UNIQUE u_management_id_and_fiscal_year (id, fiscal_year),
  INDEX idx_agent_sale_managements_agent_id (agent_id),
  INDEX idx_agent_sale_managements_fiscal_year (fiscal_year) 
);

ALTER TABLE agent_sale_managements
  ADD CONSTRAINT fk_agent_sale_managements_agent_id
  FOREIGN KEY(agent_id)
  REFERENCES agents (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE agent_sale_managements DROP FOREIGN KEY fk_agent_sale_managements_agent_id;

DROP TABLE IF EXISTS agent_sale_managements;
