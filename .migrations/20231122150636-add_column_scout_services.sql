-- スカウト媒体に面談設定メールのテンプレートを管理するカラム
-- +migrate Up
ALTER TABLE scout_services
  ADD COLUMN template_title_for_employed VARCHAR(255) NOT NULL,                    -- 面談調整メールのテンプレート ※就業中
  ADD COLUMN template_title_for_unemployed VARCHAR(255) NOT NULL; -- 面談調整メールのテンプレート ※離職中