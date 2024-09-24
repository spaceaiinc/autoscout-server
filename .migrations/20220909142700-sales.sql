-- 売上を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS sales (
    id	INT AUTO_INCREMENT NOT NULL UNIQUE,	        -- 重複しないID
    job_seeker_id INT NOT NULL UNIQUE,	            -- 求職者のID
    job_information_id INT NOT NULL,	            -- 求人のID
    accuracy INT,	                                -- 確度（ヨミ）
    contract_signed_month CHAR(7) NOT NULL,         -- 受注月
    billing_month CHAR(7) NOT NULL,                 -- 請求月
    billing_amount INT,                             -- 請求金額
    cost INT,                                       -- 原価
    gross_profit INT,                               -- 粗利
    ra_staff_id INT NOT NULL,                       -- RA担当者のid
    ca_staff_id INT NOT NULL,                       -- CA担当者のid
    ra_sales_ratio INT,                             -- RAの売上比率
    ca_sales_ratio INT,                             -- CAの売上比率
    created_at DATETIME,	                        --  作成日時
    updated_at DATETIME,	                        --  最終更新日時
    PRIMARY KEY(id),
    INDEX idx_sales_job_seeker_id (job_seeker_id),
    INDEX idx_sales_job_information_id (job_information_id),
    INDEX idx_sales_ra_staff_id (ra_staff_id),
    INDEX idx_sales_ca_staff_id (ca_staff_id)
);

ALTER TABLE sales
    ADD CONSTRAINT fk_sales_job_seeker_id
    FOREIGN KEY(job_seeker_id)
    REFERENCES job_seekers (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE sales
    ADD CONSTRAINT fk_sales_job_information_id
    FOREIGN KEY(job_information_id)
    REFERENCES job_informations (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE sales
    ADD CONSTRAINT fk_sales_ra_staff_id
    FOREIGN KEY(ra_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE sales
    ADD CONSTRAINT fk_sales_ca_staff_id
    FOREIGN KEY(ca_staff_id)
    REFERENCES agent_staffs (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE sales DROP FOREIGN KEY fk_sales_job_seeker_id;
ALTER TABLE sales DROP FOREIGN KEY fk_sales_job_information_id;
ALTER TABLE sales DROP FOREIGN KEY fk_sales_ra_staff_id;
ALTER TABLE sales DROP FOREIGN KEY fk_sales_ca_staff_id;

DROP TABLE IF EXISTS sales;


-- ALTER TABLE sales ADD cost INT;
-- ALTER TABLE sales ADD gross_profit INT;
-- ALTER TABLE sales ADD ra_staff_id INT NOT NULL;
-- ALTER TABLE sales ADD ca_staff_id INT NOT NULL;
-- ALTER TABLE sales ADD ra_sales_ratio INT;
-- ALTER TABLE sales ADD ca_sales_ratio INT;

-- ALTER TABLE sales ADD INDEX idx_sales_ra_staff_id (ra_staff_id);
-- ALTER TABLE sales ADD INDEX idx_sales_ca_staff_id (ca_staff_id);


-- ALTER TABLE sales
--     ADD CONSTRAINT fk_sales_ra_staff_id
--     FOREIGN KEY(ra_staff_id)
--     REFERENCES agent_staffs (id)
--     ON DELETE CASCADE
--     ON UPDATE CASCADE;

-- ALTER TABLE sales
--     ADD CONSTRAINT fk_sales_ca_staff_id
--     FOREIGN KEY(ca_staff_id)
--     REFERENCES agent_staffs (id)
--     ON DELETE CASCADE
--     ON UPDATE CASCADE;