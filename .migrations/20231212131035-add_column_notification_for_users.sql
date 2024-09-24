-- お知らせテーブルに送信対象を管理するカラムを追加（0: 全てのユーザー, 1: CRMのみ, 2: 送客のみ）
-- +migrate Up
ALTER TABLE notification_for_users
  ADD COLUMN target INT DEFAULT 0; -- 送信対象（0: 全てのユーザー, 1: CRMのみ, 2: 送客のみ）
