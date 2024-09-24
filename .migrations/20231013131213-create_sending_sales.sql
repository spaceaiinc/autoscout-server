-- 送客売上を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sending_sales (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	    -- 重複しないID
    sending_job_seeker_id	INT NOT NULL,		-- 送客求職者のID
    sending_enterprise_id	INT NOT NULL,		-- 送客先エージェントのID
    motoyui_sales	INT,		                -- Motoyuiの売上
    kickback INT,		                        -- キックバック（送客元の売上）
    created_at DATETIME,	                    -- 作成日時
    updated_at DATETIME,	                    -- 最終更新日時
    PRIMARY KEY(id),
    UNIQUE u_sending_sales_seeker_and_enterprise (sending_job_seeker_id, sending_enterprise_id),
    INDEX idx_sending_sales_sending_job_seeker_id (sending_job_seeker_id),
    INDEX idx_sending_sales_sending_enterprise_id (sending_enterprise_id)
);

ALTER TABLE sending_sales
    ADD CONSTRAINT fk_sending_sales_sending_job_seeker_id
    FOREIGN KEY(sending_job_seeker_id)
    REFERENCES sending_job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE sending_sales
    ADD CONSTRAINT fk_sending_sales_sending_enterprise_id
    FOREIGN KEY(sending_enterprise_id)
    REFERENCES sending_enterprises (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sending_sales DROP FOREIGN KEY fk_sending_sales_sending_job_seeker_id;
ALTER TABLE sending_sales DROP FOREIGN KEY fk_sending_sales_sending_enterprise_id;

DROP TABLE IF EXISTS sending_sales;
