-- 送客先エージェントが求職者の応募書類にアクセスするため追加
-- +migrate Up
ALTER TABLE sending_enterprises
  ADD COLUMN password VARCHAR(255) NOT NULL; -- 応募書類にアクセスするためのパスワード