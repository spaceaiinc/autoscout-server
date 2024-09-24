-- ユーザーがデプロイの内容を反映したかを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS deployment_reflections (
  id INT AUTO_INCREMENT NOT NULL UNIQUE,      -- 重複しないID
  deployment_id INT NOT NULL,                 -- デプロイID
  agent_staff_id INT NOT NULL,                -- agentStaffID
  is_reflected BOOLEAN NOT NULL,              -- 反映の有無
  created_at DATETIME,                        -- 作成日時
  updated_at DATETIME,                        -- 最終更新日時
  PRIMARY KEY(id),
  INDEX idx_deployment_reflections_deployment_id (deployment_id),
  INDEX idx_deployment_reflections_agent_staff_id (agent_staff_id)
);

ALTER TABLE deployment_reflections
  ADD CONSTRAINT fk_deployment_reflections_deployment_id
  FOREIGN KEY(deployment_id)
  REFERENCES deployment_informations (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE deployment_reflections
  ADD CONSTRAINT fk_deployment_reflections_agent_staff_id
  FOREIGN KEY(agent_staff_id)
  REFERENCES agent_staffs (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE deployment_reflections DROP FOREIGN KEY fk_deployment_reflections_deployment_id;
ALTER TABLE deployment_reflections DROP FOREIGN KEY fk_deployment_reflections_agent_staff_id;

DROP TABLE IF EXISTS deployment_reflections;
