-- 求人テーブルに普通自動車免許を追加
-- +migrate Up
ALTER TABLE job_informations
  ADD COLUMN driver_licence INT; -- 普通自動車免許（0: 必須, 1: 入社時までに取得必須, 99: 不要）
