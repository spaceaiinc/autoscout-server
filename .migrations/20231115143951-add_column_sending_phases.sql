-- 送客後に参加確認のメールを送信したかを管理するカラム
-- +migrate Up
ALTER TABLE sending_phases
  ADD COLUMN is_attended BOOLEAN NOT NULL DEFAULT FALSE; -- 参加確認有無（メールを重複して送信しないために使用するカラム）