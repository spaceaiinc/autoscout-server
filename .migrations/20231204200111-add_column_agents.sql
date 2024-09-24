-- agentsテーブルに送客のタイプを管理するカラムを追加（0: 通常, 1: 送客管理(アンドイーズ仕様)）
-- +migrate Up
ALTER TABLE agents
  ADD COLUMN sending_type INT DEFAULT 0; -- 送客のタイプ（0: 通常, 1: 送客管理(アンドイーズ仕様)）
