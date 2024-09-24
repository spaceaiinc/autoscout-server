-- エージェントアライアンステーブルに、エージェント1とエージェント2の解除申請状況を表すカラムを追加する
-- +migrate Up
ALTER TABLE agent_alliances
    ADD COLUMN agent1_cancel_request BOOLEAN NOT NULL DEFAULT FALSE,          -- エージェント1の解除申請状況        
    ADD COLUMN agent2_cancel_request BOOLEAN NOT NULL DEFAULT FALSE;   -- エージェント2の解除申請状況        

-- +migrate Down