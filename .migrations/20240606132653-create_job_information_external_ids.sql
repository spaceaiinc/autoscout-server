-- 求人の他媒体でのIDを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_external_ids (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,	                    -- 重複しないID
    job_information_id INT NOT NULL,	                        -- 重複しないカラム毎の求人ID
    external_type INT,	                                        -- 他媒体のタイプ(0: エージェントバンク)
    external_id VARCHAR(255) NOT NULL,	                        -- 他媒体でのID
    created_at DATETIME,                                        -- 作成日時
    updated_at DATETIME,                                        -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_external_ids_job_information_id (job_information_id)
);

ALTER TABLE job_information_external_ids
    ADD CONSTRAINT fk_job_information_external_ids_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_external_ids DROP FOREIGN KEY fk_job_information_external_ids_job_information_id;

DROP TABLE IF EXISTS job_information_external_ids;