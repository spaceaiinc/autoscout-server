-- 求人情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS job_informations (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,                  -- 重複しないID
    uuid CHAR(36) NOT NULL UNIQUE,                          -- 重複しないカラム毎のUUID
    billing_address_id	INT NOT NULL,	                    -- 請求先のID
    title	VARCHAR(255) NOT NULL,                          -- 求人タイトル
    recruitment_state	INT,                	            -- 募集状況
    expiration_date	CHAR(10) NOT NULL,	                    -- 募集期限（年月日）
    work_detail	TEXT NOT NULL,                              -- 仕事内容
    number_of_hires	INT,                                    -- 従業員数
    work_location	TEXT NOT NULL,                          -- 勤務地（雇入れ直後）
    transfer	INT,                                        -- 転勤有無
    transfer_detail	TEXT NOT NULL,                          -- 変更の範囲
    under_income	INT,                    	            -- 年収下限
    over_income	INT,                                        -- 年収上限
    salary	TEXT NOT NULL,                                  -- 給与詳細・昇給・賞与・諸手当
    insurance	TEXT NOT NULL,                              -- 諸手当・福利厚生
    work_time	TEXT NOT NULL,                              -- 勤務時間
    overtime_average	TEXT NOT NULL,                      -- 平均残業時間
    fixed_overtime_payment	INT,	                        -- 固定残業代超過分の支払い有無
    fixed_overtime_detail	TEXT NOT NULL,	                -- 固定残業代の詳細
    trial_period	INT,	                                -- 試用期間有無
    trial_period_detail	TEXT NOT NULL,	                    -- 試用期間の詳細
    employment_period	INT, 	                            -- 雇用期間の定め
    employment_period_detail	TEXT NOT NULL,	            -- 更新上限
    holiday_type	INT,	                                -- 休日タイプ
    holiday_detail	TEXT NOT NULL,                          -- 休日・休暇詳細
    passive_smoking	INT,                                    -- 受動喫煙対策タイプ
    selection_flow	TEXT NOT NULL,      	                -- 選考フロー
    gender	INT,                                            -- 募集性別
    nationality	INT,                                        -- 国籍
    final_education INT,                                    -- 最終学歴
    school_level	INT,                                    -- 大学レベル
    medical_history	INT,                                    -- 既往歴
    age_under	INT,                                        -- 年収想定年齢 下限
    age_over	INT,                                        -- 想定年齢 上限
    job_change	INT,                                        -- 転職回数限度
    short_resignation	INT,                                -- 短期離職条件(条件あり（備考参照） or 不可 or 不問)
    short_resignation_remarks	TEXT NOT NULL,              -- 短期離職の備考
    social_experience_year	INT,	                        -- 経験年数
    social_experience_month	INT,	                        -- 経験月数
    appearance	INT,                                        -- アピアランス
    communication	INT,                                    -- コミュニケーションスキル
    thinking	INT,                                        -- 論理的思考力
    target_detail	TEXT NOT NULL,                          -- 応募条件（エージェント向け情報）
    commission	INT,	                                    -- 成功報酬手数料（固定）
    commission_rate	INT,	                                -- 成功報酬手数料（料率）
    commission_detail TEXT NOT NULL,                        -- 手数料補足
    refund_policy	VARCHAR(255) NOT NULL,	                -- 返金規定
    required_experience_job_detail TEXT NOT NULL,           -- 必要業職種経験の詳細
    secret_memo TEXT NOT NULL,                              -- 社内限定メモ
    required_documents_detail TEXT NOT NULL,                -- 推薦時に必要な情報・書類の詳細
    employment_insurance BOOLEAN NOT NULL,	                -- 雇用保険の有無
    accident_insurance BOOLEAN NOT NULL,	                -- 労災保険の有無　
    health_insurance BOOLEAN NOT NULL,	                    -- 健康保険の有無　
    pension_insurance BOOLEAN NOT NULL,                     -- 厚生年金保険の有無
    register_phase INT,	                                    -- 求人の登録状況（0: 本登録, 1: 仮登録）
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,              -- 削除フラグ false: 有効, true: 削除済み
    study_category INT,	                                    -- 専攻学科の大分類(0:　理系尚可, 1: 理系のみ, 99: 不問)
    word_skill INT,                                         -- Wordのスキル
    excel_skill INT,                                        -- Excelのスキル
    power_point_skill INT,                                  -- PowerPointのスキル
    created_at	DATETIME,                                   -- 作成日時
    updated_at	DATETIME,                                   -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_job_informations_billing_address_id (billing_address_id),
    FULLTEXT idx_job_informations_title (title) WITH PARSER ngram
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE job_informations
    ADD CONSTRAINT fk_job_informations_billing_address_id
    FOREIGN KEY(billing_address_id)
    REFERENCES billing_addresses (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE job_informations DROP FOREIGN KEY fk_job_informations_billing_address_id;

DROP TABLE IF EXISTS job_informations;