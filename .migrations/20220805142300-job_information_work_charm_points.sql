-- 求人企業の職場の魅力情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_information_work_charm_points (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    job_information_id	INT NOT NULL,	                -- 求人のID
    title	VARCHAR(255) NOT NULL,	                    -- タイトル
    contents	TEXT NOT NULL,	                        -- 魅力
    created_at DATETIME,                                -- 作成日時
    updated_at DATETIME,                                -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_information_work_charm_points_job_information_id (job_information_id)
);

ALTER TABLE job_information_work_charm_points
    ADD CONSTRAINT fk_job_information_work_charm_points_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_information_work_charm_points DROP FOREIGN KEY fk_job_information_work_charm_points_job_information_id;

DROP TABLE IF EXISTS job_information_work_charm_points;