-- Rollback: Remove Phase 2 reputation and match history tables
DROP INDEX IF EXISTS idx_scrim_requests_suspended;
DROP TABLE IF EXISTS screenshot_uploads;
DROP TABLE IF EXISTS reputation_events;
DROP TABLE IF EXISTS match_history;

ALTER TABLE scrim_matches
    DROP COLUMN IF EXISTS result_status,
    DROP COLUMN IF EXISTS winner_team_id,
    DROP COLUMN IF EXISTS score_team1,
    DROP COLUMN IF EXISTS score_team2,
    DROP COLUMN IF EXISTS ocr_confidence,
    DROP COLUMN IF EXISTS ocr_raw_text,
    DROP COLUMN IF EXISTS verified_at;

ALTER TABLE scrim_requests
    DROP COLUMN IF EXISTS reputation_score,
    DROP COLUMN IF EXISTS total_wins,
    DROP COLUMN IF EXISTS total_losses,
    DROP COLUMN IF EXISTS is_suspended,
    DROP COLUMN IF EXISTS suspended_until;

DROP TYPE IF EXISTS match_result_status;
