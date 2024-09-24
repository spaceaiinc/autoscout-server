-- 企業プロフィールの業界を管理するテーブル（縦持ち）
-- +migrate Up
CREATE TABLE IF NOT EXISTS enterprise_industries (
    id INT AUTO_INCREMENT NOT NULL UNIQUE, -- 重複しないカラム毎のid
    enterprise_id INT NOT NULL,            -- 企業のID
    industry INT,                          -- 募集職種
    created_at DATETIME,                   -- 作成日時
    updated_at DATETIME,                   -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_enterprise_industries_contents (enterprise_id)
);

ALTER TABLE enterprise_industries
    ADD CONSTRAINT fk_enterprise_industries_enterprise_id
    FOREIGN KEY(enterprise_id)
    REFERENCES enterprise_profiles (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE enterprise_industries DROP FOREIGN KEY fk_enterprise_industries_enterprise_id;

DROP TABLE IF EXISTS enterprise_industries;
