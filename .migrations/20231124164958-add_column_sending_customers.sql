-- 送客求職者(送客元表示)をテーブルに履歴書原本と職務経歴書原本を追加する
-- +migrate Up
ALTER TABLE sending_customers
    ADD COLUMN resume_origin_url TEXT NOT NULL,	        -- 履歴書原本のURL（Word or Excel）
    ADD COLUMN cv_origin_url TEXT NOT NULL;	            -- 職務経歴書原本のURL（Word or Excel）