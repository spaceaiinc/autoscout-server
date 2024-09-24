-- タスクグループにRAエージェントIDとCAエージェントIDを追加
-- +migrate Up
ALTER TABLE task_groups
  ADD COLUMN ra_agent_id INT, -- RAエージェントID
  ADD COLUMN ca_agent_id INT; -- CAエージェントID

UPDATE task_groups
SET 
  ra_agent_id = (
    SELECT agent_id 
    FROM agent_staffs 
    WHERE agent_staffs.id = (
      SELECT agent_staff_id 
      FROM billing_addresses 
      WHERE billing_addresses.id = (
        SELECT billing_address_id 
        FROM job_informations 
        WHERE job_informations.id = task_groups.job_information_id
      )
    )
  ),
  ca_agent_id = (
    SELECT id 
    FROM agents
    WHERE agents.id = (
      SELECT agent_id 
      FROM job_seekers 
      WHERE job_seekers.id = task_groups.job_seeker_id
    )
  );
