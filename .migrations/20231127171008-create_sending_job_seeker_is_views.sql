-- 送客求職者の進捗管理上での閲覧状況を管理するテーブル（進捗管理ページでのみ使用するため、sending_job_seekersテーブルとは分離）
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_seeker_is_views (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,                 -- 重複しないID
    sending_job_seeker_id INT NOT NULL UNIQUE,	            -- 送客求職者ID
    is_not_waiting_viewed BOOLEAN NOT NULL DEFAULT TRUE,    -- 面談実施待ちタブ上での閲覧有無
    is_not_unregister_viewed BOOLEAN NOT NULL DEFAULT TRUE, -- 未登録タブ上での閲覧有無
    created_at	DATETIME,	                                -- 作成日時
    updated_at	DATETIME,	                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_job_seeker_is_views_sending_job_seeker_id (sending_job_seeker_id)
);

ALTER TABLE sending_job_seeker_is_views
    ADD CONSTRAINT fk_sending_job_seeker_is_views_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- 初期値の設定（登録済みの送客求職者分のレコードを作成）
INSERT INTO sending_job_seeker_is_views (
    sending_job_seeker_id, 
    is_not_waiting_viewed, 
    is_not_unregister_viewed, 
    created_at, 
    updated_at
)
SELECT id, true, true, created_at, created_at -- 0は初期値として設定、必要に応じて変更
FROM sending_job_seekers;

-- +migrate Down
ALTER TABLE sending_job_seeker_is_views DROP FOREIGN KEY fk_sending_job_seeker_is_views_sending_job_seeker_id;

DROP TABLE IF EXISTS sending_job_seeker_is_views;
