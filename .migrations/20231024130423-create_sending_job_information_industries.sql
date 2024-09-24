-- 送客先エージェントの求人情報の業界を管理するテーブル（縦持ち）
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_information_industries (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,         -- 重複しないカラム毎のid
    sending_job_information_id INT NOT NULL,       -- 送客先エージェントの求人のID
    industry INT,                                  -- 業界
    created_at DATETIME,                           -- 作成日時
    updated_at DATETIME,                           -- 最終更新日時
    PRIMARY KEY(id), 
    INDEX idx_sending_industries_sending_job_information_id (sending_job_information_id)
);

ALTER TABLE sending_job_information_industries
    ADD CONSTRAINT fk_sending_industries_sending_job_information_id
    FOREIGN KEY(sending_job_information_id)
    REFERENCES sending_job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_information_industries DROP FOREIGN KEY fk_sending_industries_sending_job_information_id;

DROP TABLE IF EXISTS sending_job_information_industries;