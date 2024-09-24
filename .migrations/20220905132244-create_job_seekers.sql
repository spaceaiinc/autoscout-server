-- 求職者のプロフィール情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seekers (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,                          -- 重複しないID
    uuid CHAR(36) NOT NULL UNIQUE,                                  -- 重複しないカラム毎のUUID
    line_id VARCHAR(255) NOT NULL,                                  -- 重複しないLINEのID
    agent_id INT NOT NULL,                                          -- エージェントのID
    agent_staff_id INT,                                             -- CA担当者のID
    user_status INT,                                                -- ステータス（中途, 既卒, 新卒）
    last_name VARCHAR(255) NOT NULL,                                -- 名前（姓）
    first_name VARCHAR(255) NOT NULL,                               -- 名前（名）
    last_furigana VARCHAR(255) NOT NULL,                            -- フリガナ（セイ）
    first_furigana VARCHAR(255) NOT NULL,                           -- フリガナ（メイ）
    gender INT,                                                     -- 性別
    gender_remarks VARCHAR(255) NOT NULL,                           -- 性別の備考
    birthday CHAR(10) NOT NULL,                                     -- 生年月日
    spouse INT,                                                     -- 配偶者
    support_obligation INT,                                         -- 配偶者の扶養義務（有無）
    dependents INT,                                                 -- 扶養家族人数（配偶者を除く）
    phone_number CHAR(13) NOT NULL,                                 -- 電話番号
    email VARCHAR(255) NOT NULL,                                    -- メールアドレス
    emergency_phone_number CHAR(13) NOT NULL,                       -- 緊急連絡先（電話番号）
    post_code CHAR(8) NOT NULL,                                     -- 郵便番号
    prefecture INT,                                                 -- 都道府県
    address VARCHAR(255) NOT NULL,                                  -- 住所詳細（市町村 番地 建物名 部屋番号）
    address_furigana  VARCHAR(255) NOT NULL,                        -- 住所詳細（フリガナ）
    state_of_employment INT,                                        -- 就業状況
    job_summary TEXT NOT NULL,                                      -- 職務要約
    history_supplement TEXT NOT NULL,                               -- 経歴補足
    research_content TEXT NOT NULL,                                 -- 研究内容・学チカ
    join_company_period INT,                                        -- 入社可能時期
    job_change INT,                                                 -- 転職回数
    annual_income INT,                                              -- 直近の年収
    desired_annual_income INT,                                      -- 希望年収
    transfer INT,                                                   -- 転勤可否
    transfer_requirement TEXT NOT NULL,                             -- 転勤条件
    short_resignation INT,                                          -- 短期離職(有 or 無)
    short_resignation_remarks VARCHAR(255) NOT NULL,                -- 短期離職の備考
    medical_history INT,                                            -- 既往歴
    nationality INT,                                                -- 国籍
    appearance INT,                                                 -- アピアランス
    communication INT,                                              -- コミュニケーション能力
    thinking INT,                                                   -- 論理的思考力
    secret_memo TEXT NOT NULL,                                      -- 社内限定メモ
    job_hunting_state INT,                                          -- 転職・就活状況
    recommend_reason VARCHAR(255) NOT NULL,                         -- 推薦理由
    phase INT,                                                      -- 求職者のフェーズ（エントリー、面談調整中、面談予約、面談実施）
    interview_date DATETIME NOT NULL DEFAULT "0001-01-01 01:00:00", -- 面談日時
    agreement BOOLEAN NOT NULL,                                     -- 個人情報同意の有無
    register_phase INT,                                             -- 求職者の登録状況（0: 本登録, 1: 仮登録）
    activity_memo TEXT NOT NULL,                                    -- アクティビティメモ
    study_category INT,                                             -- 専攻学科の大分類(0:　理系 1: 文系)
    word_skill INT,                                                 -- Wordのスキル
    excel_skill INT,                                                -- Excelのスキル
    power_point_skill INT,                                          -- PowerPointのスキル
    inflow_channel_id INT,                                          -- 流入経路のid 
    nationality_remarks TEXT NOT NULL,                              -- 国籍 外国籍を選択→国籍備考（フリーテキスト）を表示
    medical_history_remarks TEXT NOT NULL,                          -- 既往歴 ありを選択→既往歴備考（フリーテキスト）を表示
    acceptance_points TEXT NOT NULL,                                -- 応募承諾のポイント
    created_at DATETIME,                                            -- 作成日時
    updated_at DATETIME,                                            -- 更新日時
    PRIMARY KEY (id),
    INDEX idx_job_seekers_agent_id (agent_id),
    INDEX idx_job_seekers_agent_staff_id (agent_staff_id),
    INDEX idx_job_seekers_last_name (last_name),
    INDEX idx_job_seekers_first_name (first_name),
    INDEX idx_job_seekers_last_furigana (last_furigana),
    INDEX idx_job_seekers_first_furigana (first_furigana),
    INDEX idx_job_seekers_inflow_channel_id (inflow_channel_id)                       
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE job_seekers
    ADD CONSTRAINT fk_job_seekers_agent_id
    FOREIGN KEY(agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE 
    ON UPDATE CASCADE;

ALTER TABLE job_seekers
    ADD CONSTRAINT fk_job_seekers_agent_staff_id
    FOREIGN KEY(agent_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE 
    ON UPDATE CASCADE;

ALTER TABLE job_seekers
    ADD CONSTRAINT fk_job_seekers_inflow_channel_id
    FOREIGN KEY(inflow_channel_id)
    REFERENCES agent_inflow_channel_options (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seekers DROP FOREIGN KEY fk_job_seekers_agent_id;
ALTER TABLE job_seekers DROP FOREIGN KEY fk_job_seekers_agent_staff_id;
ALTER TABLE job_seekers DROP FOREIGN KEY fk_job_seekers_inflow_channel_id;

-- add column（後から追加）
-- ALTER TABLE job_seekers ADD COLUMN history_supplement TEXT NOT NULL;

DROP TABLE IF EXISTS job_seekers;