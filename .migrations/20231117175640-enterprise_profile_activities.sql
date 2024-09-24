-- 企業の追加情報を管理するテーブル（縦持ち）
-- +migrate Up
CREATE TABLE IF NOT EXISTS enterprise_activities (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,              -- 重複しないカラム毎のid
    enterprise_id INT NOT NULL,                         -- 企業のID
    content TEXT NOT NULL,                              -- 内容
    added_at DATETIME,                                  -- 追加日時
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_enterprise_activities_contents (enterprise_id)
);

ALTER TABLE enterprise_activities
    ADD CONSTRAINT fk_enterprise_activities_enterprise_id
    FOREIGN KEY(enterprise_id)
    REFERENCES enterprise_profiles (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE enterprise_activities DROP FOREIGN KEY fk_enterprise_activities_enterprise_id;

DROP TABLE IF EXISTS enterprise_activities;
