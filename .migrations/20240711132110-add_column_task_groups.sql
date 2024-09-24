-- タスクグループテーブルに自己応募フラグを追加する
-- +migrate Up
ALTER TABLE task_groups
    ADD COLUMN is_self_application BOOLEAN NOT NULL DEFAULT FALSE; -- 自己応募フラグ（true: 自己応募, false: エージェント応募）