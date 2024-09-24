-- 評価点を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS evaluation_points (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    task_id	INT NOT NULL UNIQUE,	                    -- 重複しないタスクのID
    selection_information_id INT,	                    -- 重複しない選考情報のID
    job_seeker_id INT NOT NULL,	                        -- 求職者のID
    job_information_id	INT NOT NULL,	                -- 求人のID
    good_point	VARCHAR(255) NOT NULL,	                -- 評価点
    ng_point	VARCHAR(255) NOT NULL,  	            -- 懸念点・NG理由
    is_passed BOOLEAN NOT NULL,                         -- 「true: 合格 or false: 不合格」
    is_re_interview BOOLEAN NOT NULL,                   -- 再面接の有無
    created_at	DATETIME,	                            -- 作成日時
    updated_at	DATETIME,	                            -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_evaluation_points_task_id (task_id),
    INDEX idx_evaluation_points_selection_information_id (selection_information_id),
    INDEX idx_evaluation_points_job_seeker_id (job_seeker_id),
    INDEX idx_evaluation_points_job_information_id (job_information_id)
);

ALTER TABLE evaluation_points
    ADD CONSTRAINT fk_evaluation_points_task_id
    FOREIGN KEY(task_id)
    REFERENCES tasks (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE evaluation_points
    ADD CONSTRAINT fk_evaluation_points_selection_information_id
    FOREIGN KEY(selection_information_id)
    REFERENCES job_information_selection_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE evaluation_points
    ADD CONSTRAINT fk_evaluation_points_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE evaluation_points
    ADD CONSTRAINT fk_evaluation_points_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE evaluation_points DROP FOREIGN KEY fk_evaluation_points_task_id;
ALTER TABLE evaluation_points DROP FOREIGN KEY fk_evaluation_points_selection_information_id;
ALTER TABLE evaluation_points DROP FOREIGN KEY fk_evaluation_points_job_seeker_id;
ALTER TABLE evaluation_points DROP FOREIGN KEY fk_evaluation_points_job_information_id;

DROP TABLE IF EXISTS evaluation_points;