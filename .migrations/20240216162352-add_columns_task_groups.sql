-- タスクグループテーブルに外部求人関連の値を管理するカラムを追加
-- +migrate Up
ALTER TABLE task_groups
  ADD COLUMN external_job_information_title VARCHAR(255) NOT NULL, -- 外部求人のタイトル
  ADD COLUMN external_company_name VARCHAR(255) NOT NULL,          -- 外部会社名のタイトル
  ADD COLUMN external_job_listing_url TEXT NOT NULL;               -- 外部求人の求人票URL

ALTER TABLE task_groups ADD FULLTEXT idx_task_groups_external_job_information_title (external_job_information_title) WITH PARSER ngram;
ALTER TABLE task_groups ADD FULLTEXT idx_task_groups_external_company_name (external_company_name) WITH PARSER ngram;
ALTER TABLE task_groups ADD FULLTEXT idx_task_groups_external_job_listing_url (external_job_listing_url) WITH PARSER ngram;
