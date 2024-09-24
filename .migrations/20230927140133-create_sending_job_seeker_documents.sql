-- 送客求職者の応募書類を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_job_seeker_documents (
    id INT AUTO_INCREMENT NOT NULL UNIQUE,	  -- 重複しないID
    sending_job_seeker_id INT NOT NULL,	      -- 送客求職者のID
    resume_origin_url TEXT NOT NULL,	        -- 履歴書原本のURL（Word or Excel）
    resume_pdf_url TEXT NOT NULL,	            -- 履歴書のURL（PDF）
    cv_origin_url TEXT NOT NULL,	            -- 職務経歴書原本のURL（Word or Excel）
    cv_pdf_url TEXT NOT NULL,	                -- 職務経歴書のURL（PDF）
    recommendation_origin_url TEXT NOT NULL,	-- 推薦状のURL（Word or Excel）
    recommendation_pdf_url TEXT NOT NULL,	    -- 推薦状のURL（PDF）
    id_photo_url TEXT NOT NULL,	              -- 証明写真のURL（jpeg or png）
    other_document1_url TEXT NOT NULL,	      -- その他①のURL
    other_document2_url TEXT NOT NULL,	      -- その他②のURL
    other_document3_url TEXT NOT NULL,	      -- その他③のURL
    created_at	DATETIME,                     -- 作成日時
    updated_at	DATETIME,                     -- 最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sending_job_seeker_documents_sending_job_seeker_id (sending_job_seeker_id)
);

ALTER TABLE sending_job_seeker_documents
    ADD CONSTRAINT fk_sending_job_seeker_documents_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_job_seeker_documents DROP FOREIGN KEY fk_sending_job_seeker_documents_sending_job_seeker_id;

DROP TABLE IF EXISTS sending_job_seeker_documents;