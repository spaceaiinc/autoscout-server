-- 送客先エージェントの求人情報の募集職種を管理するテーブル（縦持ち）
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_information_occupations (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,         -- 重複しないカラム毎のid
    sending_job_information_id INT NOT NULL,       -- 送客先エージェントの求人のID
    occupation INT,                                -- 募集職種
    created_at DATETIME,                           -- 作成日時
    updated_at DATETIME,                           -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_occupations_sending_job_information_id (sending_job_information_id)
);

ALTER TABLE sending_job_information_occupations
    ADD CONSTRAINT fk_sending_occupations_sending_job_information_id
    FOREIGN KEY(sending_job_information_id)
    REFERENCES sending_job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_information_occupations DROP FOREIGN KEY fk_sending_occupations_sending_job_information_id;

DROP TABLE IF EXISTS sending_job_information_occupations;