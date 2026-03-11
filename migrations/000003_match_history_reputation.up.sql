-- Migration: Match History, Reputation Points, and OCR Result Verification
-- Version: 000003
-- Description: Adds tables and columns to support automated result verification,
--              match history, and the Anti-Fraud reputation system.

-- ============================================================
-- 1. Match Result Status Enum
-- ============================================================
CREATE TYPE match_result_status AS ENUM (
    'pending_upload',   -- Waiting for screenshot upload
    'verifying',        -- OCR service is processing
    'verified',         -- AI confirmed the result
    'disputed',         -- AI flagged as suspicious / fraud
    'manual_review'     -- Escalated to human review
);

-- ============================================================
-- 2. Extend scrim_matches with result & verification columns
-- ============================================================
ALTER TABLE scrim_matches
    ADD COLUMN IF NOT EXISTS result_status   match_result_status DEFAULT 'pending_upload',
    ADD COLUMN IF NOT EXISTS winner_team_id  UUID REFERENCES scrim_requests(id),
    ADD COLUMN IF NOT EXISTS score_team1     VARCHAR(20),
    ADD COLUMN IF NOT EXISTS score_team2     VARCHAR(20),
    ADD COLUMN IF NOT EXISTS ocr_confidence  NUMERIC(5,2),       -- 0.00 - 100.00 %
    ADD COLUMN IF NOT EXISTS ocr_raw_text    TEXT,               -- raw OCR output
    ADD COLUMN IF NOT EXISTS verified_at     TIMESTAMP WITH TIME ZONE;

COMMENT ON COLUMN scrim_matches.result_status  IS 'Tracks AI verification pipeline state';
COMMENT ON COLUMN scrim_matches.ocr_confidence IS 'OCR confidence score 0-100 from Python service';

-- ============================================================
-- 3. Match History Table
-- ============================================================
CREATE TABLE match_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scrim_match_id  UUID NOT NULL REFERENCES scrim_matches(id) ON DELETE CASCADE,
    team_id         UUID NOT NULL REFERENCES scrim_requests(id),
    opponent_id     UUID NOT NULL REFERENCES scrim_requests(id),
    category        scrim_category NOT NULL,
    result          VARCHAR(10) NOT NULL CHECK (result IN ('win', 'loss', 'draw', 'cancelled')),
    rank_weight     INTEGER NOT NULL,
    rep_delta       INTEGER NOT NULL DEFAULT 0,  -- reputation points change (+/-)
    played_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_match_history_team_id    ON match_history(team_id);
CREATE INDEX idx_match_history_played_at  ON match_history(played_at DESC);

COMMENT ON TABLE  match_history    IS 'Per-team historical record of all completed scrims';
COMMENT ON COLUMN match_history.rep_delta IS 'Points added/deducted to reputation after this match';

-- ============================================================
-- 4. Reputation Events Table (audit trail)
-- ============================================================
CREATE TABLE reputation_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id         UUID NOT NULL,          -- references scrim_requests.id
    event_type      VARCHAR(50) NOT NULL,   -- 'win','loss','fraud_penalty','ghost_penalty','manual_adjust'
    points_delta    INTEGER NOT NULL,
    reason          TEXT,
    match_id        UUID,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rep_events_team_id    ON reputation_events(team_id);
CREATE INDEX idx_rep_events_created_at ON reputation_events(created_at DESC);

COMMENT ON TABLE reputation_events IS 'Audit log for every reputation point change';

-- ============================================================
-- 5. Reputation Score on scrim_requests
-- ============================================================
ALTER TABLE scrim_requests
    ADD COLUMN IF NOT EXISTS reputation_score INTEGER NOT NULL DEFAULT 100,
    ADD COLUMN IF NOT EXISTS total_wins        INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_losses      INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_suspended      BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS suspended_until   TIMESTAMP WITH TIME ZONE;

COMMENT ON COLUMN scrim_requests.reputation_score IS '0-200, default 100. Drops on fraud/ghosting, rises on wins';
COMMENT ON COLUMN scrim_requests.is_suspended     IS 'Anti-fraud suspension flag set by OCR service';

-- ============================================================
-- 6. OCR Screenshot Upload Log
-- ============================================================
CREATE TABLE screenshot_uploads (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scrim_match_id  UUID NOT NULL REFERENCES scrim_matches(id) ON DELETE CASCADE,
    uploader_team_id UUID NOT NULL,
    file_path       TEXT NOT NULL,           -- local/S3 path
    mime_type       VARCHAR(50),
    ocr_job_id      VARCHAR(100),            -- async task ID from FastAPI
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_screenshot_uploads_match_id ON screenshot_uploads(scrim_match_id);

COMMENT ON TABLE screenshot_uploads IS 'Tracks screenshots submitted for OCR anti-fraud verification';

-- ============================================================
-- 7. Index for suspension lookup in matchmaking
-- ============================================================
CREATE INDEX idx_scrim_requests_suspended ON scrim_requests(is_suspended)
    WHERE is_suspended = TRUE;
