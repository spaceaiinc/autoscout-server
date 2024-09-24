-- 求人企業請求先の人事担当者情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS billing_address_hr_staffs (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    billing_address_id INT NOT NULL,	                -- 請求先ID
    hr_staff_name	VARCHAR(255) NOT NULL,	            -- 人事担当者名
    hr_staff_email	VARCHAR(255) NOT NULL,	            -- 人事担当者メールアドレス
    hr_staff_phone_number	CHAR(13) NOT NULL,	        -- 人事担当者電話番号
    created_at	DATETIME,                               -- 作成日時
    updated_at	DATETIME,                               -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_billing_address_hr_staffs_billing_address_id (billing_address_id)
);

ALTER TABLE billing_address_hr_staffs
    ADD CONSTRAINT fk_billing_address_hr_staffs_billing_address_id
    FOREIGN KEY(billing_address_id)
    REFERENCES billing_addresses (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE billing_address_hr_staffs DROP FOREIGN KEY fk_billing_address_hr_staffs_billing_address_id;

DROP TABLE IF EXISTS billing_address_hr_staffs;