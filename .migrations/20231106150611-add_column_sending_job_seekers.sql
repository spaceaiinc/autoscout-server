-- 送客求職者がアンケート回答時の質問・要望を管理するためのカラムを追加
-- +migrate Up
ALTER TABLE sending_job_seekers
  ADD COLUMN question TEXT NOT NULL; -- 質問や要望