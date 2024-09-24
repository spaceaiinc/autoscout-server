-- エージェントのスカウトサービスのテンプレートを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS scout_service_templates (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                -- 重複しないID
    scout_service_id INT NOT NULL,	                        -- スカウトサービスID
    start_hour INT,                                         -- 開始UTC時間(媒体共通/0:0時, 1:1時, 2:2時, ..., 23:23時)
    start_minute INT,                                       -- 開始UTC分(媒体共通/0:0分, 1:1分, 2:2分, ..., 59:59分)
    run_on_monday BOOLEAN NOT NULL DEFAULT FALSE,	        -- 月曜日に走らせるかどうか/false:走らせない true:走る(共通)
    run_on_tuesday BOOLEAN NOT NULL DEFAULT FALSE,	        -- 火曜日に走らせるかどうか/false:走らせない true:走る(共通)
    run_on_wednesday BOOLEAN NOT NULL DEFAULT FALSE,	    -- 水曜日に走らせるかどうか/false:走らせない true:走る(共通)
    run_on_thursday BOOLEAN NOT NULL DEFAULT FALSE,	        -- 木曜日に走らせるかどうか/false:走らせない true:走る(共通)
    run_on_friday BOOLEAN NOT NULL DEFAULT FALSE,	        -- 金曜日に走らせるかどうか/false:走らせない true:走る(共通)
    run_on_saturday BOOLEAN NOT NULL DEFAULT FALSE,	        -- 土曜日に走らせるかどうか/false:走らせない true:走る(共通)
    run_on_sunday BOOLEAN NOT NULL DEFAULT FALSE,	        -- 日曜日に走らせるかどうか/false:走らせない true:走る(共通)
    scout_count INT,	                                    -- スカウト件数(MYNAVI(50, 100, 300, 500), RAN(max100), AMBI(50, 200))
    search_title VARCHAR(255) NOT NULL,	                    -- スカウトに利用する検索条件のタイトル(媒体共通)
    message_title VARCHAR(255) NOT NULL,	                -- スカウトに利用するメッセージのタイトル(媒体共通)
    job_information_title VARCHAR(255) NOT NULL,	        -- スカウトに利用する求人情報のタイトル(RAN)
    job_information_id VARCHAR(255) NOT NULL,	            -- スカウトに利用する求人情報のID(RAN)
    age_limit INT,                                          -- 年齢(MYNAVI, RAN, AMBI)
    scout_type INT,	                                        -- スカウトタイプ(AMBI/0: 通常スカウト, 1: プレミアムスカウト)
    auto_remind BOOLEAN NOT NULL DEFAULT FALSE,	            -- 自動リマインド(AMBI/false:しない, true:する)
    reply_limit INT,	                                    -- 返信期限(AMBI, マイナビ. 例:2日後期限なら2)
    created_at DATETIME,                                    -- 作成日時
    updated_at DATETIME,                                    -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_scout_service_templates_scout_service_id (scout_service_id)
);

ALTER TABLE scout_service_templates
    ADD CONSTRAINT fk_scout_service_templates_scout_service_id
    FOREIGN KEY(scout_service_id)
    REFERENCES scout_services (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE scout_service_templates DROP FOREIGN KEY fk_scout_service_templates_scout_service_id;

DROP TABLE IF EXISTS scout_service_templates;