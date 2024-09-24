-- メッセージテンプレートを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS interview_adjustment_templates (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	    -- 重複しないID
    agent_id	INT NOT NULL,	                    -- 担当者のID
    send_scene	INT,	                          -- 送信シーン（マスクレジュメ or 書類推薦 or 請求）
    title	VARCHAR(255) NOT NULL,	              -- テンプレートのタイトル
    subject VARCHAR(255) NOT NULL,              -- 件名
    content	TEXT NOT NULL,	                    -- テンプレートの内容
    created_at	DATETIME,	                      -- 作成日時
    updated_at	DATETIME,	                      -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_interview_adjustment_templates_agent_id (agent_id)
);

ALTER TABLE interview_adjustment_templates
    ADD CONSTRAINT fk_interview_adjustment_templates_agent_id
    FOREIGN KEY(agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE interview_adjustment_templates DROP FOREIGN KEY fk_interview_adjustment_templates_agent_id;

DROP TABLE IF EXISTS interview_adjustment_templates;