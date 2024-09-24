-- スカウトテンプレートに最終送信時間を追加する
-- +migrate Up
ALTER TABLE scout_service_templates
    ADD COLUMN last_send_count INT DEFAULT 0,
    ADD COLUMN last_send_at DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00';

-- +migrate Down