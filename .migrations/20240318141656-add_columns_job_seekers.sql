-- 求職者テーブルに人物像カラムを追加
-- +migrate Up
ALTER TABLE job_seekers
  ADD COLUMN recommendation_profile TEXT NOT NULL, -- 人物像（推薦状用）
  ADD COLUMN candid_profile TEXT NOT NULL; -- 人物像（本音）
