-- 求人情報の募集対象を管理するテーブル（縦持ち）
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_targets (
    id INT AUTO_INCREMENT NOT NULL UNIQUE, -- 重複しないカラム毎のid
    job_information_id INT NOT NULL,       -- 求人のID
    target INT,                            -- 募集対象
    created_at DATETIME,                   -- 作成日時
    updated_at DATETIME,                   -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_targets_contents (job_information_id)
);

ALTER TABLE job_information_targets
    ADD CONSTRAINT fk_job_information_targets_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_targets DROP FOREIGN KEY fk_job_information_targets_job_information_id;

DROP TABLE IF EXISTS job_information_targets;
