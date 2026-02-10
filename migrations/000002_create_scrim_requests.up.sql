-- Migration: Create scrim_requests table for new POKE/WARKOP matchmaking system
-- Version: 000002
-- Description: Implements category-based matchmaking with WhatsApp integration

-- Create enum type for category
CREATE TYPE scrim_category AS ENUM ('POKE', 'WARKOP');

-- Create enum type for status
CREATE TYPE scrim_status AS ENUM ('searching', 'matched', 'expired', 'cancelled');

-- Create scrim_requests table
CREATE TABLE scrim_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_name VARCHAR(100) NOT NULL,
    whatsapp_number VARCHAR(20) NOT NULL,
    category scrim_category NOT NULL,
    rank_weight INTEGER NOT NULL CHECK (rank_weight >= 1 AND rank_weight <= 10),
    status scrim_status NOT NULL DEFAULT 'searching',
    match_id UUID,
    ip_address VARCHAR(45),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP + INTERVAL '30 minutes'),
    matched_at TIMESTAMP WITH TIME ZONE
);

-- Create scrim_matches table
CREATE TABLE scrim_matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team1_id UUID NOT NULL REFERENCES scrim_requests(id),
    team2_id UUID NOT NULL REFERENCES scrim_requests(id),
    category scrim_category NOT NULL,
    rank_diff INTEGER,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP + INTERVAL '60 seconds')
);

-- Create indexes for performance
CREATE INDEX idx_scrim_requests_status ON scrim_requests(status);
CREATE INDEX idx_scrim_requests_category ON scrim_requests(category);
CREATE INDEX idx_scrim_requests_searching ON scrim_requests(category, rank_weight) 
    WHERE status = 'searching';
CREATE INDEX idx_scrim_requests_expires ON scrim_requests(expires_at) 
    WHERE status = 'searching';
CREATE INDEX idx_scrim_matches_status ON scrim_matches(status);
CREATE INDEX idx_scrim_matches_expires ON scrim_matches(expires_at) 
    WHERE status = 'pending';

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_scrim_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for auto-updating updated_at
CREATE TRIGGER trigger_scrim_requests_updated_at
    BEFORE UPDATE ON scrim_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_scrim_updated_at();

-- Create function to auto-expire old requests
CREATE OR REPLACE FUNCTION expire_old_scrim_requests()
RETURNS void AS $$
BEGIN
    UPDATE scrim_requests
    SET status = 'expired'
    WHERE status = 'searching'
      AND expires_at < CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;

-- Create function to auto-cancel expired matches
CREATE OR REPLACE FUNCTION cancel_expired_scrim_matches()
RETURNS void AS $$
BEGIN
    UPDATE scrim_matches
    SET status = 'cancelled'
    WHERE status = 'pending'
      AND expires_at < CURRENT_TIMESTAMP;
      
    -- Also update the scrim_requests
    UPDATE scrim_requests sr
    SET status = 'cancelled', match_id = NULL
    FROM scrim_matches sm
    WHERE (sr.id = sm.team1_id OR sr.id = sm.team2_id)
      AND sm.status = 'cancelled'
      AND sr.status = 'matched';
END;
$$ LANGUAGE plpgsql;

-- Add comments for documentation
COMMENT ON TABLE scrim_requests IS 'Scrim matchmaking requests with POKE/WARKOP categories';
COMMENT ON COLUMN scrim_requests.category IS 'POKE: rank_weight 1-8, WARKOP: rank_weight 9-10';
COMMENT ON COLUMN scrim_requests.rank_weight IS 'Player skill level: 1-10 (POKE uses ±2 tolerance, WARKOP ignores)';
COMMENT ON COLUMN scrim_requests.whatsapp_number IS 'Team captain WhatsApp number for match notification';
COMMENT ON COLUMN scrim_requests.ip_address IS 'Used for rate limiting (1 active request per IP)';
