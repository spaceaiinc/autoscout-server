-- 送客求職者(送客元表示)を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_customers (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,   -- 重複しないカラム毎のid
    agent_id INT NOT NULL,                   -- エージェントのID
    last_name VARCHAR(255) NOT NULL,         -- 名前（姓）
    first_name VARCHAR(255) NOT NULL,        -- 名前（名）
    last_furigana VARCHAR(255) NOT NULL,     -- フリガナ（セイ）
    first_furigana VARCHAR(255) NOT NULL,    -- フリガナ（メイ）
    phone_number CHAR(13) NOT NULL,          -- 電話番号
    email VARCHAR(255) NOT NULL,             -- メールアドレス
    resume_pdf_url TEXT NOT NULL,	         -- 履歴書のURL（PDF）
    cv_pdf_url TEXT NOT NULL,	             -- 職務経歴書のURL（PDF）
    interview_date DATETIME,                 -- 面談日時
    interview_information TEXT NOT NULL,     -- 面談情報
    remarks TEXT NOT NULL,                   -- 備考
    gender INT,                              -- 性別
    nationality INT,                         -- 国籍
    nationality_remarks TEXT NOT NULL,       -- 国籍 外国籍を選択→国籍備考（フリーテキスト）を表示
    birthday CHAR(10) NOT NULL,              -- 生年月日
    post_code CHAR(8) NOT NULL,              -- 郵便番号
    prefecture INT,                          -- 都道府県
    address VARCHAR(255) NOT NULL,           -- 住所詳細（市町村 番地 建物名 部屋番号）
    address_furigana  VARCHAR(255) NOT NULL, -- 住所詳細（フリガナ）
    school_name	VARCHAR(255) NOT NULL,	     -- 学校名（最終学歴）
    subject	VARCHAR(255) NOT NULL,	         -- 学部・学科・コース（最終学歴）
    entrance_year	CHAR(7) NOT NULL,	     -- 入学年月（最終学歴）
    graduation_year	CHAR(7) NOT NULL,        -- 卒業年月（最終学歴）
    state_of_employment INT,                 -- 就業状況
    job_change INT,                          -- 転職回数
    job_summary TEXT NOT NULL,               -- 職務要約
    company_name VARCHAR(255) NOT NULL,      -- 会社名（直近の就業先）
    joining_year	CHAR(7) NOT NULL,	     -- 入社年月（直近の就業先）
    retire_year	CHAR(7) NOT NULL,	         -- 退職年月（直近の就業先）
    first_status INT,                        -- 開始ステータス（入社, 入行, 入局など）（直近の就業先）
    last_status INT,                         -- 終了ステータス（一身上の都合により退職, 派遣期間満了につき退職など）（直近の就業先）
    job_description TEXT NOT NULL,           -- 職務内容（直近の就業先）
    history_supplement TEXT NOT NULL,        -- 経歴補足
    phase INT,                               -- 相談状況（0: 面談実施待ち, 1: 送客応諾, 2: 送客完了, 3: 送客なし/終了）
    sending_at DATETIME,                     -- 送客の実行時間
    created_at DATETIME,                     -- 作成日時
    updated_at DATETIME,                     -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_customers_agent_id (agent_id)
);

ALTER TABLE sending_customers
    ADD CONSTRAINT fk_sending_customers_agent_id
    FOREIGN KEY(agent_id)
    REFERENCES agents (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_customers DROP FOREIGN KEY fk_sending_customers_agent_id

DROP TABLE IF EXISTS sending_customers;