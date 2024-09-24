-- アンケート結果を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS selection_questionnaires (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    uuid VARCHAR(255) NOT NULL UNIQUE,	                -- 重複しないUUID
    job_seeker_id INT NOT NULL,	                        -- 求職者のID
    job_information_id INT NOT NULL,	                -- 求人のID
    selection_information_id INT NOT NULL,	            -- 選考情報のID
    my_ranking INT,	                                    -- 志望度
    my_ranking_reason TEXT NOT NULL,	                -- 志望度が上がった理由
    concern_point TEXT NOT NULL,	                    -- 心配な点・懸念点
    continue_selection INT,	                            -- 選考継続希望有無
    my_ranking_detail TEXT NOT NULL,	                -- 他社選考含めての志望順位
    selection_question TEXT NOT NULL,	                -- 選考で聞かれた質問
    remarks	TEXT NOT NULL,	                            -- その他質問・備考
    is_answer BOOLEAN NOT NULL,                         -- 回答状況
    is_self_introduction BOOLEAN NOT NULL,              -- 自己紹介（どんなこと聞かれた？）
    is_self_pr BOOLEAN NOT NULL,                        -- 自己PR（強み・長所）（どんなこと聞かれた？）
    is_retire_reason BOOLEAN NOT NULL,                  -- 退職理由（どんなこと聞かれた？）
    is_job_change_axis BOOLEAN NOT NULL,                -- 転職軸（どんなこと聞かれた？）
    is_applying_reason BOOLEAN NOT NULL,                -- 志望動機（志望理由）（どんなこと聞かれた？）
    is_career_vision BOOLEAN NOT NULL,                  -- キャリアビジョン（どんなこと聞かれた？）
    is_reverse_question BOOLEAN NOT NULL,               -- 逆質問（どんなこと聞かれた？）
    intention_to_job_offer INT,                         -- 内定が出た場合の意向
    intention_detail TEXT NOT NULL,                     -- 意向の詳細
    created_at DATETIME,	                            -- 作成日時
    updated_at DATETIME,	                            -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_selection_questionnaires_job_seeker_id (job_seeker_id),
    INDEX idx_selection_questionnaires_job_information_id (job_information_id),
    INDEX idx_selection_questionnaires_selection_information_id (selection_information_id)
);

ALTER TABLE selection_questionnaires
    ADD CONSTRAINT fk_selection_questionnaires_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE selection_questionnaires
    ADD CONSTRAINT fk_selection_questionnaires_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE selection_questionnaires
    ADD CONSTRAINT fk_selection_questionnaires_selection_information_id
    FOREIGN KEY(selection_information_id)
    REFERENCES job_information_selection_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE selection_questionnaires DROP FOREIGN KEY fk_selection_questionnaires_job_seeker_id;
ALTER TABLE selection_questionnaires DROP FOREIGN KEY fk_selection_questionnaires_job_information_id;
ALTER TABLE selection_questionnaires DROP FOREIGN KEY fk_selection_questionnaires_selection_information_id;

DROP TABLE IF EXISTS selection_questionnaires;