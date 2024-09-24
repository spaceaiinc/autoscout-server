-- agent_monthly_salesテーブルの金額を管理するカラムを「INT」から「BIGINT」に変更
-- 金額の単位を「千円」から「円」に変更することで格納する数字が増えるため修正
-- +migrate Up
ALTER TABLE agent_monthly_sales
	MODIFY COLUMN order_sales_budget BIGINT,
	MODIFY COLUMN order_cost_budget BIGINT,
	MODIFY COLUMN order_gross_profit_budget BIGINT,
	MODIFY COLUMN order_assumed_unit_price BIGINT,
	MODIFY COLUMN claim_sales_revenue_budget BIGINT,
	MODIFY COLUMN claim_cost_budget BIGINT,
	MODIFY COLUMN claim_gross_margin_budget BIGINT,
	MODIFY COLUMN claim_assumed_unit_price BIGINT;
