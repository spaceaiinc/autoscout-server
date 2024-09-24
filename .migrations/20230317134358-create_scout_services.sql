-- エージェントのスカウトサービスを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS scout_services (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    uuid	CHAR(36) NOT NULL UNIQUE,	                -- 重複しないカラム毎のUUID
    agent_robot_id INT NOT NULL,                        -- エージェントロボットID
    agent_staff_id INT NOT NULL,                        -- エージェントスタッフID
    login_id VARCHAR(255) NOT NULL,	                    -- ログインID
    password VARCHAR(255) NOT NULL,	                    -- パスワード
    service_type INT,	                                -- サービスタイプ(0: RAN, 1: マイナビ転職スカウト, 2:AMBI)
    is_active BOOLEAN NOT NULL DEFAULT FALSE,	        -- アクティブかどうか/false:走らせない true:走る(共通)
    memo TEXT NOT NULL,	                                -- メモ
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_scout_services_agent_robot_id (agent_robot_id),
    INDEX idx_scout_services_agent_staff_id (agent_staff_id)
);

ALTER TABLE scout_services
    ADD CONSTRAINT fk_scout_services_agent_robot_id
    FOREIGN KEY(agent_robot_id)
    REFERENCES agent_robots (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE scout_services
    ADD CONSTRAINT fk_scout_services_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE scout_services DROP FOREIGN KEY fk_scout_services_agent_robot_id;
ALTER TABLE scout_services DROP FOREIGN KEY fk_scout_services_agent_staff_id;

DROP TABLE IF EXISTS scout_services;