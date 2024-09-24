-- スカウト媒体に面談設定メールテンプレートIDを管理するカラム
-- +migrate Up
ALTER TABLE scout_services
  ADD COLUMN interview_adjustment_template_id INT, -- 面談調整メールテンプレートID
  ADD COLUMN inflow_channel_id INT;             -- 流入経路ID

ALTER TABLE scout_services
    ADD INDEX idx_scout_services_interview_adjustment_template_id (interview_adjustment_template_id),
    ADD INDEX idx_scout_services_inflow_channel_id (inflow_channel_id);

ALTER TABLE scout_services
    ADD CONSTRAINT fk_scout_services_interview_adjustment_template_id
    FOREIGN KEY(interview_adjustment_template_id)
    REFERENCES interview_adjustment_templates (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE scout_services
    ADD CONSTRAINT fk_scout_services_inflow_channel_id
    FOREIGN KEY(inflow_channel_id)
    REFERENCES agent_inflow_channel_options (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE scout_services DROP FOREIGN KEY fk_scout_services_interview_adjustment_template_id;
ALTER TABLE scout_services DROP FOREIGN KEY fk_scout_services_inflow_channel_id;