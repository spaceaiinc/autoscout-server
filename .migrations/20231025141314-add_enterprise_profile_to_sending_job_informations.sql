-- 送客先エージェントの求人情報に企業情報を追加 ※求人票で使用
-- +migrate Up
ALTER TABLE sending_job_informations
    ADD COLUMN corporate_site_url VARCHAR(255) NOT NULL,   -- 企業サイト
    ADD COLUMN post_code CHAR(8) NOT NULL,                 -- 郵便番号
    ADD COLUMN office_location VARCHAR(255) NOT NULL,      -- 住所
    ADD COLUMN establishment VARCHAR(255) NOT NULL,        -- 設立
    ADD COLUMN employee_number_single INT,                 -- 従業員数(単体)
    ADD COLUMN employee_number_group INT,                  -- 従業員数(連結)
    ADD COLUMN public_offering INT,                        -- 株式公開
    ADD COLUMN earnings_year INT,                          -- 売上高(年度)
    ADD COLUMN earnings VARCHAR(255) NOT NULL,             -- 売上高 
    ADD COLUMN business_detail TEXT NOT NULL;              -- 事業内容
