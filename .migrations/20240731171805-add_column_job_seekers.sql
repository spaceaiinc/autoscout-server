-- 求職者テーブルに「マッチング求人を閲覧可能かを管理する値」を追加する
-- +migrate Up
ALTER TABLE job_seekers
  ADD COLUMN can_view_matching_job BOOLEAN NOT NULL DEFAULT FALSE;