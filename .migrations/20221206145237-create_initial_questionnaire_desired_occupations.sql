-- 求職者の面談前アンケートの希望職種を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS initial_questionnaire_desired_occupations (
  id	INT AUTO_INCREMENT NOT NULL UNIQUE,             -- 重複しないID
  initial_questionnaire_id	INT NOT NULL,	            -- アンケートのID
  desired_occupation INT,	                                    -- 希望職種
  desired_rank   INT,	                                        -- 希望順位
  created_at	DATETIME,	                              -- 作成日時
  updated_at	DATETIME,	                              -- 最終更新日時
  PRIMARY KEY(id),
  INDEX idx_initial_questionnaire_desired_occupations_contents (initial_questionnaire_id)
);

ALTER TABLE initial_questionnaire_desired_occupations
    ADD CONSTRAINT fk_initial_questionnaire_desired_occupations_contents
    FOREIGN KEY(initial_questionnaire_id)
    REFERENCES initial_questionnaires (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE initial_questionnaire_desired_occupations DROP FOREIGN KEY fk_initial_questionnaire_desired_occupations_contents;

DROP TABLE IF EXISTS initial_questionnaire_desired_occupations;