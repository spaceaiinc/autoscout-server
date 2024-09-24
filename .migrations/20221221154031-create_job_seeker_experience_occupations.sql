-- 求職者の経験職種を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_seeker_experience_occupations (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	                  -- 重複しないID
    department_history_id INT NOT NULL,                       -- 部署履歴ID
    occupation INT,	                                          -- 職種
    created_at	DATETIME,                                     -- 作成日時
    updated_at	DATETIME,                                     -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_seeker_experience_occupations_department_history_id (department_history_id)
);

ALTER TABLE job_seeker_experience_occupations
    ADD CONSTRAINT fk_job_seeker_experience_occupations_department_history_id
    FOREIGN KEY(department_history_id)
    REFERENCES job_seeker_department_histories (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_seeker_experience_occupations DROP FOREIGN KEY fk_job_seeker_experience_occupations_department_history_id;

DROP TABLE IF EXISTS job_seeker_experience_occupations;
