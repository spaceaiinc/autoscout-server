-- 送客先エージェントに「送客対象」と「送客時に必要な情報」と「備考」を追加
-- +migrate Up
ALTER TABLE sending_enterprises
    ADD COLUMN sending_target TEXT NOT NULL,	     -- 送客対象
    ADD COLUMN sending_required_info TEXT NOT NULL,	 -- 送客時に必要な情報
    ADD COLUMN remarks TEXT NOT NULL;                -- 備考
