-- 求職者の希望企業規模タイプを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_desired_company_scales (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                          -- 重複しないID
    job_seeker_id	INT NOT NULL,	                                  -- 求職者のID
    desired_company_scale INT,	                                      -- 企業規模タイプ
    -- desired_rank INT,                                                 -- 企業規模の順位
    created_at	DATETIME,                                             -- 作成日時
    updated_at	DATETIME,                                             -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_desired_company_scales_job_seeker_id (job_seeker_id)
);

ALTER TABLE job_seeker_desired_company_scales
    ADD CONSTRAINT fk_job_seeker_desired_company_scales_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_desired_company_scales DROP FOREIGN KEY fk_job_seeker_desired_company_scales_job_seeker_id;

DROP TABLE IF EXISTS job_seeker_desired_company_scales;
