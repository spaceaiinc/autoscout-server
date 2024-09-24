-- 求人企業の請求情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS billing_addresses (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    uuid CHAR(36) NOT NULL UNIQUE,                      -- 重複しないカラム毎のUUID
    enterprise_id	INT NOT NULL,	                    -- 企業のID
    agent_staff_id	INT NOT NULL,	                    -- RA担当者のID
    contract_phase	INT,                                -- 契約フェーズ 0: リーガルチェック中, 1: リーガルチェック完了, 2: 契約締結済み
    contract_date	CHAR(10) NOT NULL,	                -- 基本契約締結日
    payment_policy	VARCHAR(255) NOT NULL,	            -- 支払い規定
    company_name	VARCHAR(255) NOT NULL,	            -- 会社名
    address VARCHAR(255) NOT NULL,                      -- 請求先住所 ※郵便番号も含む
    how_to_recommend	TEXT NOT NULL,                  -- 推薦方法
    title VARCHAR(255) NOT NULL,                        -- 請求先のタイトル
    created_at	DATETIME,                               -- 作成日時
    updated_at	DATETIME,                               -- 最終更新日時
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,          -- 削除フラグ false: 有効, true: 削除済み
    PRIMARY KEY(id),
    INDEX idx_billing_addresses_enterprise_id (enterprise_id),
    INDEX idx_billing_addresses_agent_staff_id (agent_staff_id)
);

ALTER TABLE billing_addresses
    ADD CONSTRAINT fk_billing_addresses_enterprise_id
    FOREIGN KEY(enterprise_id)
    REFERENCES enterprise_profiles (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE billing_addresses
    ADD CONSTRAINT fk_billing_addresses_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE billing_addresses DROP FOREIGN KEY fk_billing_addresses_enterprise_id;

DROP TABLE IF EXISTS billing_addresses;