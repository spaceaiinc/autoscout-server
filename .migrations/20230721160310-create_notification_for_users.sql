-- Systemからのお知らせを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS notification_for_users (
    id INT AUTO_INCREMENT NOT NULL UNIQUE, -- 重複しないカラム毎のid
    title VARCHAR(255) NOT NULL,           -- お知らせのタイトル
    body TEXT NOT NULL,                    -- お知らせの本文
    created_at DATETIME,                   -- 作成日時
    updated_at DATETIME,                   -- 最終更新日時
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE IF EXISTS notification_for_users;