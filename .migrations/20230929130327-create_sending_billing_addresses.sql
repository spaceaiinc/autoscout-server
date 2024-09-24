-- 送客先エージェントの請求情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_billing_addresses (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	          -- 重複しないID
    uuid CHAR(36) NOT NULL UNIQUE,                    -- 重複しないカラム毎のUUID
    sending_enterprise_id	INT NOT NULL,	              -- 送客先エージェントのID
    agent_staff_id	INT NOT NULL,	                  -- 請求先担当者のID デフォルト:企業担当を選択
    contract_phase	INT,                              -- 契約フェーズ 0: リーガルチェック中, 1: リーガルチェック完了, 2: 契約締結済み
    contract_date	CHAR(10) NOT NULL,	              -- 基本契約締結日
    payment_policy	VARCHAR(255) NOT NULL,	          -- 支払い規定
    agent_name	VARCHAR(255) NOT NULL,	              -- 送客先エージェント名
    address VARCHAR(255) NOT NULL,                    -- 請求先住所 ※郵便番号も含む
    title VARCHAR(255) NOT NULL,                      -- 請求先のタイトル デフォルト: 送客先エージェント名
    schedule_adjustment_url VARCHAR(255) NOT NULL,    -- 日程調整URL
    commission INT,                                   -- 送客単価
    created_at	DATETIME,                             -- 作成日時
    updated_at	DATETIME,                             -- 最終更新日時
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,        -- 削除フラグ false: 有効, true: 削除済み
    PRIMARY KEY(id),
    INDEX idx_sending_billing_addresses_sending_enterprise_id (sending_enterprise_id),
    INDEX idx_sending_billing_addresses_agent_staff_id (agent_staff_id)
);

ALTER TABLE sending_billing_addresses
    ADD CONSTRAINT fk_sending_billing_addresses_sending_enterprise_id
    FOREIGN KEY(sending_enterprise_id)
    REFERENCES sending_enterprises (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE sending_billing_addresses
    ADD CONSTRAINT fk_sending_billing_addresses_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_billing_addresses DROP FOREIGN KEY fk_sending_billing_addresses_sending_enterprise_id;

DROP TABLE IF EXISTS sending_billing_addresses;