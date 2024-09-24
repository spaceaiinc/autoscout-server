-- スカウトサービスのエントリー取得時間を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS scout_service_get_entry_times (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                -- 重複しないID
    scout_service_id	INT NOT NULL,	                    -- スカウトサービスID
    start_hour INT,                                         -- 開始UTC時間(媒体共通/0:0時, 1:1時, 2:2時, ..., 23:23時)
    start_minute INT,                                       -- 開始UTC分(媒体共通/0:0分, 1:1分, 2:2分, ..., 59:59分)
    created_at DATETIME,                                    -- 作成日時
    updated_at DATETIME,                                    -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_scout_service_get_entry_times_scout_service_id (scout_service_id)
);

ALTER TABLE scout_service_get_entry_times
    ADD CONSTRAINT fk_scout_service_get_entry_times_scout_service_id
    FOREIGN KEY(scout_service_id)
    REFERENCES scout_services (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE scout_service_get_entry_times DROP FOREIGN KEY fk_scout_service_get_entry_times_scout_service_id;

DROP TABLE IF EXISTS scout_service_get_entry_times;