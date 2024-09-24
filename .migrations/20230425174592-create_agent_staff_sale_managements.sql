-- エージェントスタッフの売上管理の決算月や公開状況を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS agent_staff_sale_managements (
  id INT AUTO_INCREMENT NOT NULL UNIQUE,      -- 重複しないID
  management_id INT NOT NULL,                 -- 親テーブルのID
  agent_staff_id INT NOT NULL,                -- agentID
  created_at DATETIME,                        -- 作成日時
  updated_at DATETIME,                        -- 最終更新日時
  PRIMARY KEY(id),
  UNIQUE u_management_id_and_agent_staff_id (management_id, agent_staff_id),
  INDEX idx_agent_staff_sale_managements_management_id (management_id),
  INDEX idx_agent_staff_sale_managements_agent_staff_id (agent_staff_id)
);

ALTER TABLE agent_staff_sale_managements
  ADD CONSTRAINT fk_agent_staff_sale_managements_management_id
  FOREIGN KEY(management_id)
  REFERENCES agent_sale_managements (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE agent_staff_sale_managements
  ADD CONSTRAINT fk_agent_staff_sale_managements_agent_staff_id
  FOREIGN KEY(agent_staff_id)
  REFERENCES agent_staffs (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE agent_staff_sale_managements DROP FOREIGN KEY fk_agent_staff_sale_managements_management_id;
ALTER TABLE agent_staff_sale_managements DROP FOREIGN KEY fk_agent_staff_sale_managements_agent_staff_id;

DROP TABLE IF EXISTS agent_staff_sale_managements;
