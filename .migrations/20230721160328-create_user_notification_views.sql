-- お知らせを確認したユーザーをを管理するテーブル（ユーザーが閲覧したらレコードを追加）
-- +migrate Up
CREATE TABLE IF NOT EXISTS user_notification_views (
    id INT AUTO_INCREMENT NOT NULL UNIQUE, -- 重複しないカラム毎のid
    notification_id INT NOT NULL,          -- notification_for_usersテーブルのid
    agent_staff_id INT NOT NULL,           -- agent_staffsテーブルのid
    created_at DATETIME,                   -- 作成日時
    updated_at DATETIME,                   -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_user_notification_views_notification_id (notification_id),
    INDEX idx_user_notification_views_agent_staff_id (agent_staff_id),
    UNIQUE u_notification_id_and_agent_staff_id (notification_id, agent_staff_id)
);

ALTER TABLE user_notification_views
    ADD CONSTRAINT fk_user_notification_views_notification_id
    FOREIGN KEY(notification_id)
    REFERENCES notification_for_users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE user_notification_views
    ADD CONSTRAINT fk_user_notification_views_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE user_notification_views DROP FOREIGN KEY fk_user_notification_views_notification_id;
ALTER TABLE user_notification_views DROP FOREIGN KEY fk_user_notification_views_agent_staff_id;

DROP TABLE IF EXISTS user_notification_views;