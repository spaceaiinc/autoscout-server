-- 送客先エージェントを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_enterprises (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,      -- 重複しないカラム毎のid
    uuid CHAR(36) NOT NULL UNIQUE,              -- 重複しないカラム毎のUUID
    company_name VARCHAR(255) NOT NULL,         -- 企業名
    agent_staff_id INT NOT NULL,                -- 企業担当者ID
    corporate_site_url VARCHAR(255) NOT NULL,   -- 企業サイト
    representative VARCHAR(255) NOT NULL,       -- 代表者
    establishment VARCHAR(255) NOT NULL,        -- 設立
    post_code CHAR(8) NOT NULL,                 -- 郵便番号
    office_location VARCHAR(255) NOT NULL,      -- 住所
    employee_number_single INT,                 -- 従業員数(単体)
    public_offering INT,                        -- 株式公開
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,  -- 削除フラグ false: 有効, true: 削除済み
    created_at DATETIME,                        -- 作成日時
    updated_at DATETIME,                        -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_enterprises_agent_staff_id (agent_staff_id),
    INDEX idx_sending_enterprises_post_code (post_code),
    FULLTEXT idx_sending_enterprises_company_name (company_name) WITH PARSER ngram
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE sending_enterprises
    ADD CONSTRAINT fk_sending_enterprises_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_enterprises DROP FOREIGN KEY fk_sending_enterprises_agent_staff_id;

DROP TABLE IF EXISTS sending_enterprises;