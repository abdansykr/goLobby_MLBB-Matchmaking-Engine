-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY,
    captain_id UUID NOT NULL,
    captain_name VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) NOT NULL,
    average_rank INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'WAITING',
    match_id UUID,
    reputation_score INT NOT NULL DEFAULT 100,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Indexes for performance
    CONSTRAINT chk_rank CHECK (average_rank >= 0 AND average_rank <= 100),
    CONSTRAINT chk_reputation CHECK (reputation_score >= 0 AND reputation_score <= 200)
);

-- Create index on captain_id for fast lookups
CREATE INDEX idx_teams_captain_id ON teams(captain_id);

-- Create index on status for matchmaking queries
CREATE INDEX idx_teams_status ON teams(status);

-- Create index on average_rank for rank-based matching
CREATE INDEX idx_teams_average_rank ON teams(average_rank) WHERE status = 'WAITING';

-- Create matches table
CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY,
    team1_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    team2_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    rank_diff INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    
    -- Ensure teams are different
    CONSTRAINT chk_different_teams CHECK (team1_id != team2_id)
);

-- Create indexes for matches
CREATE INDEX idx_matches_team1_id ON matches(team1_id);
CREATE INDEX idx_matches_team2_id ON matches(team2_id);
CREATE INDEX idx_matches_status ON matches(status);
CREATE INDEX idx_matches_expires_at ON matches(expires_at) WHERE status = 'PENDING';

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for auto-updating updated_at
CREATE TRIGGER update_teams_updated_at BEFORE UPDATE ON teams
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_matches_updated_at BEFORE UPDATE ON matches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
