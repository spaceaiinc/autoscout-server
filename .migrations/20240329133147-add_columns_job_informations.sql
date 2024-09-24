-- 求人テーブルに仕事内容（雇入れ直後）と仕事内容（変更の範囲）を追加
-- +migrate Up
ALTER TABLE job_informations
  ADD COLUMN work_detail_after_hiring TEXT NOT NULL,         -- 仕事内容（雇入れ直後）
  ADD COLUMN work_detail_scope_of_change TEXT NOT NULL; -- 仕事内容（変更の範囲）
