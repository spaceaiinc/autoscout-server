-- 送客先エージェントの参考資料を管理するテーブル（縦持ち）
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_enterprise_reference_materials (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,              -- 重複しないカラム毎のid
    sending_enterprise_id INT NOT NULL,                 -- 企業のID
    reference1_pdf_url	VARCHAR(2000) NOT NULL,	        -- 参考資料1（PDF）
    reference2_pdf_url	VARCHAR(2000) NOT NULL,	        -- 参考資料2（PDF）
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_enterprise_reference_materials_contents (sending_enterprise_id)
);

ALTER TABLE sending_enterprise_reference_materials
    ADD CONSTRAINT fk_sending_enterprise_reference_materials_sending_enterprise_id
    FOREIGN KEY(sending_enterprise_id)
    REFERENCES sending_enterprises (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_enterprise_reference_materials DROP FOREIGN KEY fk_sending_enterprise_reference_materials_sending_enterprise_id;

DROP TABLE IF EXISTS sending_enterprise_reference_materials;
