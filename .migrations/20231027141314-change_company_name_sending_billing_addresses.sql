-- 送客先エージェントの請求先企業名を修正
-- +migrate Up
ALTER TABLE sending_billing_addresses
    CHANGE agent_name company_name VARCHAR(255) NOT NULL;         -- 企業名

-- +migrate Down
