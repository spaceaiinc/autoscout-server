-- 求人テーブルに面接確約フラグを追加
-- +migrate Up
ALTER TABLE job_informations
  ADD COLUMN is_guaranteed_interview BOOLEAN NOT NULL DEFAULT false; -- 面接確約フラグ