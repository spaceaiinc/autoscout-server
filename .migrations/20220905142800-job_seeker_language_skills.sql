-- 求職者の希望勤務地を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_language_skills (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    job_seeker_id   INT NOT NULL,	                    -- 求職者のID
    language_type	INT,	                            -- 語学の種類
    language_level    INT,	                            -- 語学レベル {0:日常会話, 1:ビジネス 99:不問}
    toeic	INT,	                                    -- TOEICの点数
    toeic_examination_year CHAR(7) NOT NULL,	        -- TOEICの受験年月
    toefl_ibt	INT,	                                -- TOEFL iBT点数
    toefl_ibt_examination_year CHAR(7) NOT NULL,	    -- TOEFL iBTの受験年月
    toefl_pbt	INT,	                                -- TOEFL PBT点数
    toefl_pbt_examination_year CHAR(7) NOT NULL,	    -- TOEFL PBTの受験年月
    created_at	DATETIME,                               -- 作成日時
    updated_at	DATETIME,                               -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_language_skills_job_seeker_id (job_seeker_id)
);

ALTER TABLE job_seeker_language_skills
    ADD CONSTRAINT fk_job_seeker_language_skills_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_language_skills DROP FOREIGN KEY fk_job_seeker_language_skills_job_seeker_id;

DROP TABLE IF EXISTS job_seeker_language_skills;