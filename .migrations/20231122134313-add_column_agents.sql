-- エージェントテーブルに送客機能の有無を追加する
-- +migrate Up
ALTER TABLE agents
    ADD COLUMN is_sending_active BOOLEAN NOT NULL DEFAULT TRUE; -- 送客機能の有無

-- +migrate Down