-- 送客先エージェントの「エージェントの特徴」を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_enterprise_specialities (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,	            -- 重複しないID
    sending_enterprise_id INT NOT NULL  UNIQUE,	        -- 送客先エージェントのID
    image_url TEXT NOT NULL,                            -- エージェント画像
    job_information_count INT,                          -- 保有求人数
    specialized_occupation TEXT NOT NULL,               -- 得意な職種
    specialized_industry TEXT NOT NULL,                 -- 得意な業種
    specialized_area TEXT NOT NULL,                     -- 得意エリア
    specialized_company_type TEXT NOT NULL,             -- 得意な企業タイプ
    specialized_job_seeker_type TEXT NOT NULL,          -- 得意な求職者タイプ
    consulting_strengths TEXT NOT NULL,                 -- コンサルティングの強み
    support_content TEXT NOT NULL,                      -- サポート内容
    pr_point TEXT NOT NULL,                             -- PRポイント
    created_at	DATETIME,                               -- 作成日時
    updated_at	DATETIME,                               -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_enterprise_specialities_sending_enterprise_id (sending_enterprise_id)
);

ALTER TABLE sending_enterprise_specialities
  ADD CONSTRAINT fk_sending_enterprise_specialities_sending_enterprise_id
  FOREIGN KEY(sending_enterprise_id)
  REFERENCES sending_enterprises (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

-- 初期値の設定（登録済みの送客企業分のレコードを作成）
INSERT INTO sending_enterprise_specialities (
    sending_enterprise_id,
    image_url,
    job_information_count,
    specialized_occupation,
    specialized_industry,
    specialized_area,
    specialized_company_type,
    specialized_job_seeker_type,
    consulting_strengths,
    support_content,
    pr_point,
    created_at, 
    updated_at
)
SELECT 
    id,
    "",
    0,
    "",
    "",
    "",
    "",
    "",
    "",
    "",
    "",
    created_at, 
    created_at -- 0は初期値として設定、必要に応じて変更
FROM sending_enterprises;

-- +migrate Down
ALTER TABLE sending_enterprise_specialities DROP FOREIGN KEY fk_sending_enterprise_specialities_sending_enterprise_id;

DROP TABLE IF EXISTS sending_enterprise_specialities;