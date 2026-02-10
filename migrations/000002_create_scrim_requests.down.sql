-- Rollback migration for scrim_requests system
-- Version: 000002

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_scrim_requests_updated_at ON scrim_requests;

-- Drop functions
DROP FUNCTION IF EXISTS update_scrim_updated_at();
DROP FUNCTION IF EXISTS expire_old_scrim_requests();
DROP FUNCTION IF EXISTS cancel_expired_scrim_matches();

-- Drop tables
DROP TABLE IF EXISTS scrim_matches CASCADE;
DROP TABLE IF EXISTS scrim_requests CASCADE;

-- Drop enum types
DROP TYPE IF EXISTS scrim_status;
DROP TYPE IF EXISTS scrim_category;
