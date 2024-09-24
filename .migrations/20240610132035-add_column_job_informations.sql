-- 求人テーブルに書類通過率,内定率,直近の応募数を追加
-- +migrate Up
ALTER TABLE job_informations
  ADD COLUMN document_passing_rate INT,             -- 書類通過率
  ADD COLUMN offer_rate INT,                        -- 内定率
  ADD COLUMN number_of_recent_applications INT;     -- 直近の応募数