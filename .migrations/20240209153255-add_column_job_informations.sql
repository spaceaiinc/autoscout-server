-- 求人テーブルに外部求人かを管理するカラムを追加
-- +migrate Up
ALTER TABLE job_informations
  ADD COLUMN is_external BOOLEAN NOT NULL; -- 外部求人か
