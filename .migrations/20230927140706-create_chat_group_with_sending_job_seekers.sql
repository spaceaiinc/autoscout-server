-- エージェントと送客求職者のチャットグループ情報のデータを保存するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS chat_group_with_sending_job_seekers (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                -- 重複しないID
    sending_job_seeker_id	INT NOT NULL,	                        -- 重複しないカラム毎の送客求職者id
    agent_last_send_at	DATETIME,	                        -- エージェントの最終送信時間
    agent_last_watched_at	DATETIME,	                    -- エージェントの最終閲覧時間
    sending_job_seeker_last_send_at	DATETIME,	                    -- 求職者の最終送信時間
    sending_job_seeker_last_watched_at	DATETIME,	                -- 求職者の最終閲覧時間
    line_active BOOLEAN NOT NULL,                                   -- LINEのブロック有無
    created_at	DATETIME,	                                -- 作成日時
    updated_at	DATETIME,	                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_chat_group_with_sending_job_seekers_sending_job_seeker_id (sending_job_seeker_id)
);

ALTER TABLE chat_group_with_sending_job_seekers
    ADD CONSTRAINT fk_chat_group_with_sending_job_seekers_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;


-- +migrate Down
ALTER TABLE chat_group_with_sending_job_seekers DROP FOREIGN KEY fk_chat_group_with_sending_job_seekers_sending_job_seeker_id;

DROP TABLE IF EXISTS chat_group_with_sending_job_seekers;