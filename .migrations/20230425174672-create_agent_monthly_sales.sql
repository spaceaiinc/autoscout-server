-- エージェント全体の月間の売上情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS agent_monthly_sales (
  id	INT AUTO_INCREMENT NOT NULL UNIQUE,      -- 重複しないID
  management_id INT NOT NULL,	                 -- 親テーブルのID
  sales_month DATE,                             -- 売上月

  -- 受注関連
  order_sales_budget INT,                      -- 受注売上予算
  order_cost_budget INT,                       -- 原価予算
  order_gross_profit_budget INT,               -- 受注粗利予算
  order_assumed_unit_price INT,                -- 単価想定
  order_expected_offer_acceptance INT,         -- 内定承諾想定

  -- 請求関連
  claim_sales_revenue_budget INT,              -- 請求売上予算
  claim_cost_budget INT,                       -- 原価予算
  claim_gross_margin_budget INT,               -- 請求粗利予算
  claim_assumed_unit_price INT,                -- 単価想定
  claim_expected_new_employee_number INT,      -- 入社人数想定

  -- KPI（求職者ベース）
  seeker_offer_acceptance_target INT,          -- 内定承諾目標
  seeker_offer_target INT,                     -- 内定目標
  seeker_final_selection_target INT,           -- 最終選考目標
  seeker_selection_target INT,                 -- 選考目標
  seeker_recommendation_completion_target INT, -- 推薦完了目標
  seeker_job_introduction_target INT,          -- 求人紹介目標
  seeker_interview_target INT,                 -- 面談目標

  -- KPI (稼働件数ベース)
  active_offer_target INT,                     -- 内定目標 
  active_final_selection_target INT,           -- 最終選考目標 
  active_selection_target INT,                 -- 選考目標 
  active_recommendation_completion_target INT, -- 推薦完了目標 
  active_job_introduction_target INT,          -- 求人紹介目標 
    
  -- KPI（求職者面談実施月ベース）
  interview_offer_acceptance_target INT,          -- 内定承諾目標
  interview_offer_target INT,                     -- 内定目標
  interview_final_selection_target INT,           -- 最終選考目標
  interview_selection_target INT,                 -- 選考目標
  interview_recommendation_completion_target INT, -- 推薦完了目標
  interview_job_introduction_target INT,          -- 求人紹介目標
  interview_interview_target INT,                 -- 面談目標
  
  created_at	DATETIME,                        -- 作成日時
  updated_at	DATETIME,                        -- 最終更新日時
  PRIMARY KEY(id),
  INDEX idx_agent_monthly_sales_management_id (management_id)
);

ALTER TABLE agent_monthly_sales
  ADD CONSTRAINT fk_agent_monthly_sales_management_id
  FOREIGN KEY(management_id)
  REFERENCES agent_sale_managements (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE agent_monthly_sales DROP FOREIGN KEY fk_agent_monthly_sales_management_id;

DROP TABLE IF EXISTS agent_monthly_sales;