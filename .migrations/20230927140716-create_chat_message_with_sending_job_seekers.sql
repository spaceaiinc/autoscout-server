-- エージェントと送客求職者のチャットメッセージ情報のデータを保存するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS chat_message_with_sending_job_seekers (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    group_id	INT NOT NULL,                           -- チャットグループID	
    user_type	INT,	                                -- エージェント（0）or 送客求職者（1）
    message	TEXT NOT NULL,	                            -- メッセージ	
    package_id	VARCHAR(255) NOT NULL,	                -- スタンプの表示に使用するID
    sticker_id	VARCHAR(255) NOT NULL,                  -- スタンプの表示に使用するID
    original_content_url	VARCHAR(2000) NOT NULL,	    -- 画像ファイル or 動画ファイル or 音声ファイルのURL
    preview_image_url	VARCHAR(2000) NOT NULL,	        -- 画像ファイル or 動画ファイルのプレビュー表示用のファイルURL
    duration	INT,	                                -- 音声ファイルに使用する値
    line_message_id VARCHAR(255) NOT NULL,	            -- LINEのメッセージID
    line_message_type INT,                              -- LINEのメッセージタイプ （0: テキスト, 1: スタンプ, 2: 画像, 3: 動画, 4: 音声
    created_at	DATETIME,	                            -- 作成日時
    updated_at	DATETIME,	                            -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_chat_message_with_sending_job_seekers_group_id (group_id)
);

ALTER TABLE chat_message_with_sending_job_seekers
    ADD CONSTRAINT fk_chat_message_with_sending_job_seekers_group_id
    FOREIGN KEY(group_id)
    REFERENCES chat_group_with_sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE chat_message_with_sending_job_seekers DROP FOREIGN KEY fk_chat_message_with_sending_job_seekers_group_id

DROP TABLE IF EXISTS chat_message_with_sending_job_seekers;