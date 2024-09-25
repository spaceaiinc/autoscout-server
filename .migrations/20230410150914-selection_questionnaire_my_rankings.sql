-- スカウトサービスのエントリー取得時間を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS selection_questionnaire_my_rankings (
    id	INT AUTO_INCREMENT         NOT NULL UNIQUE,	               -- 重複しないID
    selection_questionnaire_id	   INT NOT NULL,	               -- 選考後アンケートのid
    ranking	                       INT, 	                       -- 志望順位
    company_name	               VARCHAR(255) NOT NULL,	       -- 企業名
    phase	                       INT, 	                       -- フェーズ
    selection_date	            VARCHAR(255) NOT NULL,	           -- 確定日時
    created_at	DATETIME,                               -- 作成日時
    updated_at	DATETIME,                               -- 最終更新日時
    PRIMARY KEY(id)
);

ALTER TABLE selection_questionnaire_my_rankings
    ADD CONSTRAINT fk_my_rankings_selection_questionnaire_id
    FOREIGN KEY(selection_questionnaire_id)
    REFERENCES selection_questionnaires (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE selection_questionnaire_my_rankings DROP FOREIGN KEY fk_my_rankings_selection_questionnaire_id;

DROP TABLE IF EXISTS selection_questionnaire_my_rankings;