-- 採用管理サービスの求人企業を一括インポートするテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS initial_enterprise_importers (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,          -- 重複しないID
    uuid CHAR(36) NOT NULL UNIQUE,	                -- 重複しないカラム毎のUUID
    agent_staff_id INT NOT NULL,	                -- エージェント担当者ID
    service_type INT NOT NULL,	                    -- サービスタイプ 0: サーカス
    login_id VARCHAR(255) NOT NULL,	                -- ログインID
    password VARCHAR(255) NOT NULL,	                -- パスワード
    start_date VARCHAR(255) NOT NULL,	            -- 取得開始日 fmt: YYYY-MM-DD
    start_hour INT,	                                -- 取得開始時間(時)
    max_count INT,	                                -- 取得件数 サーカスの場合(MAX100)
    offset INT,	                                    -- 取得開始位置 サーカスの場合(100ずつ)
    search_url TEXT,	                            -- 検索URL サーカスの場合
    is_success BOOLEAN NOT NULL DEFAULT false,	    -- 成功したかどうか true: 成功, false: 失敗or未実行
    created_at DATETIME,                            -- 作成日時
    updated_at DATETIME,                            -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_initial_enterprise_importers_agent_staff_id (agent_staff_id)
);

ALTER TABLE initial_enterprise_importers
    ADD CONSTRAINT fk_initial_enterprise_importers_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;


-- +migrate Down
ALTER TABLE initial_enterprise_importers DROP FOREIGN KEY fk_initial_enterprise_importers_agent_staff_id;

DROP TABLE IF EXISTS initial_enterprise_importers;
