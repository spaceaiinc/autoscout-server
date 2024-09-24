-- エージェントテーブルに送客用同意書ファイルのURL、CRM機能の有無、アライアンス機能の有無を追加する
-- +migrate Up
ALTER TABLE agents
    ADD COLUMN sending_agreement_file_url TEXT NOT NULL,                  -- 送客用同意書ファイルのURL
    ADD COLUMN is_crm_active BOOLEAN NOT NULL DEFAULT TRUE,               -- CRM機能の有無
    ADD COLUMN is_alliance_active BOOLEAN NOT NULL DEFAULT FALSE;         -- アライアンス機能の有無

-- +migrate Down