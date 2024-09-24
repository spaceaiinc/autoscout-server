-- 企業プロフィールの参考資料を管理するテーブル（縦持ち）
-- +migrate Up
CREATE TABLE IF NOT EXISTS enterprise_reference_materials (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,              -- 重複しないカラム毎のid
    enterprise_id INT NOT NULL,                         -- 企業のID
    reference1_pdf_url	VARCHAR(2000) NOT NULL,	        -- 参考資料1（PDF）
    reference2_pdf_url	VARCHAR(2000) NOT NULL,	        -- 参考資料2（PDF）
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_enterprise_reference_materials_contents (enterprise_id)
);

ALTER TABLE enterprise_reference_materials
    ADD CONSTRAINT fk_enterprise_reference_materials_enterprise_id
    FOREIGN KEY(enterprise_id)
    REFERENCES enterprise_profiles (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE enterprise_reference_materials DROP FOREIGN KEY fk_enterprise_reference_materials_enterprise_id;

DROP TABLE IF EXISTS enterprise_reference_materials;
