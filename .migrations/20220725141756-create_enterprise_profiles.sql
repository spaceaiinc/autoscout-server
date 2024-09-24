-- 企業のプロフィール情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS enterprise_profiles (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,      -- 重複しないカラム毎のid
    uuid CHAR(36) NOT NULL UNIQUE,              -- 重複しないカラム毎のUUID
    company_name VARCHAR(255) NOT NULL,         -- 社名
    agent_staff_id INT NOT NULL,                -- 企業担当者ID
    corporate_site_url VARCHAR(255) NOT NULL,   -- 企業サイト
    representative VARCHAR(255) NOT NULL,       -- 代表者
    establishment VARCHAR(255) NOT NULL,        -- 設立
    post_code CHAR(8) NOT NULL,                 -- 郵便番号
    office_location VARCHAR(255) NOT NULL,      -- 住所
    employee_number_single INT,                 -- 従業員数(単体)
    employee_number_group INT,                  -- 従業員数(連結)
    capital VARCHAR(255) NOT NULL,              -- 資本金
    public_offering INT,                        -- 株式公開
    earnings_year INT,                          -- 売上高(年度)
    earnings VARCHAR(255) NOT NULL,             -- 売上高 
    business_detail TEXT NOT NULL,              -- 事業内容
    created_at DATETIME,                        -- 作成日時
    updated_at DATETIME,                        -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_enterprise_profiles_agent_staff_id (agent_staff_id),
    INDEX idx_enterprise_profiles_post_code (post_code),
    FULLTEXT idx_enterprise_profiles_company_name (company_name) WITH PARSER ngram
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE enterprise_profiles
    ADD CONSTRAINT fk_enteprise_profiles_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE enterprise_profiles DROP FOREIGN KEY fk_enterprise_profiles_agent_staff_id;

DROP TABLE IF EXISTS enterprise_profiles;