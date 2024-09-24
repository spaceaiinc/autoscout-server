-- スカウトテンプレートに最終送信時間を追加する
-- +migrate Up
ALTER TABLE scout_services
    ADD COLUMN last_send_count INT;

-- +migrate Down