-- CA側のアライアンスユーザー情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS agent_staffs (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,                 -- 重複しないカラム毎のid
    uuid CHAR(36) NOT NULL UNIQUE,                         -- 重複しないカラム毎のUUID
    agent_id INT NOT NULL,                                 -- 重複しない所属アライアンスエージェントのid
    firebase_id VARCHAR(255) NOT NULL UNIQUE,              -- 重複しないfirebase_id
    authority INT NOT NULL,                                -- 管理者権限
    staff_name VARCHAR(255) NOT NULL,                      -- 担当者名
    furigana VARCHAR(255) NOT NULL,                        -- 担当者フリガナ（カナ）
    email VARCHAR(255) NOT NULL,                           -- メールアドレス
    staff_phone_number CHAR(13) NOT NULL,                  -- 電話番号 
    department VARCHAR(255) NOT NULL,                      -- 部署
    position VARCHAR(255) NOT NULL,                        -- 役職
    remarks VARCHAR(255) NOT NULL,                         -- 備考
    usage_status INT NOT NULL,                             -- 利用状況
    notification INT NOT NULL,                             -- 通知設定 *未使用
    notification_job_seeker BOOLEAN NOT NULL DEFAULT true, -- メール通知（求職者）
    notification_unwatched BOOLEAN NOT NULL DEFAULT false, -- メール通知（未処理・未読）
    last_login DATETIME,                                   -- 最終ログイン時間
    usage_start_date DATE NOT NULL,                        -- 利用開始日(2023-01-01形式)
    usage_end_date DATE NOT NULL,                          -- 利用終了日(2023-01-01形式)
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,      -- 削除フラグ false: 有効, true: 削除済み
    created_at DATETIME,                                   -- 作成日時
    updated_at DATETIME,                                   -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_agent_staffs (agent_id)
);

ALTER TABLE agent_staffs
    ADD CONSTRAINT fk_agent_staffs_agent_id
    FOREIGN KEY(agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE agent_staffs DROP FOREIGN KEY fk_agent_staffs_agent_id;

DROP TABLE IF EXISTS agent_staffs;
