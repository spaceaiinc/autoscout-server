-- 求職者テーブルにパスワードを追加する
-- +migrate Up
ALTER TABLE job_seekers
    ADD COLUMN password VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN reset_password_token VARCHAR(255) NOT NULL DEFAULT '';