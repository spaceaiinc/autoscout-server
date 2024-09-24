-- 送客求職者の自己PR情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_seeker_self_promotions (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,                 -- 重複しないID
    sending_job_seeker_id	INT NOT NULL,	                -- 送客求職者ID
    title	VARCHAR(255) NOT NULL,	                        -- 自己PRのタイトル
    contents	TEXT NOT NULL,	                            -- 自己PRの内容
    created_at	DATETIME,	                                -- 作成日時
    updated_at	DATETIME,	                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_job_seeker_self_promotions_sending_job_seeker_id (sending_job_seeker_id)
);

ALTER TABLE sending_job_seeker_self_promotions
    ADD CONSTRAINT fk_sending_job_seeker_self_promotions_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_seeker_self_promotions DROP FOREIGN KEY fk_sending_job_seeker_self_promotions_sending_job_seeker_id;

DROP TABLE IF EXISTS sending_job_seeker_self_promotions;