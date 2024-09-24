-- タスクグループにRAエージェントIDとCAエージェントIDをINT NOT NULLに変更
-- +migrate Up
ALTER TABLE task_groups
  MODIFY COLUMN ra_agent_id INT NOT NULL, -- RAエージェントID
  MODIFY COLUMN ca_agent_id INT NOT NULL; -- CAエージェントID

ALTER TABLE task_groups
    ADD INDEX idx_task_groups_ra_agent_id (ra_agent_id),
    ADD INDEX idx_task_groups_ca_agent_id (ca_agent_id);

ALTER TABLE task_groups
    ADD CONSTRAINT fk_task_groups_ra_agent_id
    FOREIGN KEY(ra_agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE task_groups
    ADD CONSTRAINT fk_task_groups_ca_agent_id
    FOREIGN KEY(ca_agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE task_groups DROP FOREIGN KEY fk_task_groups_ra_agent_id;
ALTER TABLE task_groups DROP FOREIGN KEY fk_task_groups_ca_agent_id;