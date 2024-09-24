-- 送客先エージェントの請求先の担当者情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_billing_address_staffs (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	          -- 重複しないID
    sending_billing_address_id INT NOT NULL,	        -- 請求先ID
    staff_name	VARCHAR(255) NOT NULL,	              -- 担当者名
    staff_email	VARCHAR(255) NOT NULL,	              -- 担当者メールアドレス
    staff_phone_number	CHAR(13) NOT NULL,	          -- 担当者電話番号
    created_at	DATETIME,                             -- 作成日時
    updated_at	DATETIME,                             -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_billing_address_staffs_sending_billing_address_id (sending_billing_address_id)
);

ALTER TABLE sending_billing_address_staffs
    ADD CONSTRAINT fk_sending_billing_address_staffs_sending_billing_address_id
    FOREIGN KEY(sending_billing_address_id)
    REFERENCES sending_billing_addresses (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_billing_address_staffs DROP FOREIGN KEY fk_sending_billing_address_staffs_sending_billing_address_id;

DROP TABLE IF EXISTS sending_billing_address_staffs;