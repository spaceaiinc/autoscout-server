-- CA側のアライアンスエージェント情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS agents (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,                     -- 重複しないカラム毎のid
    uuid	CHAR(36) NOT NULL UNIQUE,	                       -- 重複しないカラム毎のUUID
    agent_name VARCHAR(255) NOT NULL,                          -- エージェント名
    office_location VARCHAR(255) NOT NULL,                     -- 住所
    representative    VARCHAR(255) NOT NULL,	               -- 代表者
    establish	CHAR(7) NOT NULL,	                           -- 設立（年月）
    corporate_site_url	VARCHAR(255) NOT NULL,	               -- ホームページ
    permission_code VARCHAR(255) NOT NULL,                     -- 職業紹介許可番号
    permission_year	CHAR(7) NOT NULL,	                       -- 紹介事業許可年月
    workers_count INT,	                                       -- 職業紹介の業務に従事する者の数
    phone_number VARCHAR(255) NOT NULL,	                       -- 電話番号
    interview_adjustment_email VARCHAR(255) NOT NULL,	       -- メールアドレス（面談調整用）
    agreement_file_url TEXT NOT NULL,                          -- 同意書ファイルのURL
    line_bot_id VARCHAR(255) NOT NULL,	                       -- LINE BotアカウントのID
    line_messaging_channel_secret VARCHAR(255) NOT NULL,	   -- LINE Messaging APIのチャネルシークレット
    line_messaging_channel_access_token VARCHAR(255) NOT NULL, -- LINE Messaging APIのチャネルアクセストークン
    line_loging_channel_id VARCHAR(255) NOT NULL,	           -- LINE LoginのチャネルID
    line_login_channel_secret VARCHAR(255) NOT NULL,           -- LINE Loginのチャネルシークレット
    created_at DATETIME,                                       -- 作成日時
    updated_at DATETIME,                                       -- 最終更新日時
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE IF EXISTS agents;