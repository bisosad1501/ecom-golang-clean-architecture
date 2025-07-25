-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    old_values TEXT,
    new_values TEXT,
    ip_address INET,
    user_agent TEXT,
    session_id VARCHAR(255),
    metadata TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create security_events table
CREATE TABLE IF NOT EXISTS security_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    event_type VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')),
    description TEXT NOT NULL,
    ip_address INET,
    user_agent TEXT,
    location VARCHAR(255),
    metadata TEXT,
    resolved BOOLEAN DEFAULT FALSE,
    resolved_by UUID,
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create system_events table
CREATE TABLE IF NOT EXISTS system_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('SUCCESS', 'FAILED', 'IN_PROGRESS', 'CANCELLED')),
    duration BIGINT, -- Duration in milliseconds
    metadata TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for audit_logs
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_id ON audit_logs(resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_ip_address ON audit_logs(ip_address);
CREATE INDEX IF NOT EXISTS idx_audit_logs_session_id ON audit_logs(session_id);

-- Create indexes for security_events
CREATE INDEX IF NOT EXISTS idx_security_events_user_id ON security_events(user_id);
CREATE INDEX IF NOT EXISTS idx_security_events_event_type ON security_events(event_type);
CREATE INDEX IF NOT EXISTS idx_security_events_severity ON security_events(severity);
CREATE INDEX IF NOT EXISTS idx_security_events_resolved ON security_events(resolved);
CREATE INDEX IF NOT EXISTS idx_security_events_created_at ON security_events(created_at);
CREATE INDEX IF NOT EXISTS idx_security_events_ip_address ON security_events(ip_address);

-- Create indexes for system_events
CREATE INDEX IF NOT EXISTS idx_system_events_event_type ON system_events(event_type);
CREATE INDEX IF NOT EXISTS idx_system_events_category ON system_events(category);
CREATE INDEX IF NOT EXISTS idx_system_events_status ON system_events(status);
CREATE INDEX IF NOT EXISTS idx_system_events_created_at ON system_events(created_at);

-- Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_resource ON audit_logs(user_id, resource);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_action ON audit_logs(resource, action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_user ON audit_logs(created_at, user_id);

CREATE INDEX IF NOT EXISTS idx_security_events_severity_resolved ON security_events(severity, resolved);
CREATE INDEX IF NOT EXISTS idx_security_events_created_severity ON security_events(created_at, severity);

-- Create partial indexes for performance
CREATE INDEX IF NOT EXISTS idx_audit_logs_recent ON audit_logs(created_at) 
WHERE created_at > NOW() - INTERVAL '30 days';

CREATE INDEX IF NOT EXISTS idx_security_events_unresolved ON security_events(created_at, severity) 
WHERE resolved = FALSE;

-- Add foreign key constraints (if users table exists)
-- Note: Uncomment these if you have a users table
-- ALTER TABLE audit_logs ADD CONSTRAINT fk_audit_logs_user_id 
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
-- 
-- ALTER TABLE security_events ADD CONSTRAINT fk_security_events_user_id 
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
-- 
-- ALTER TABLE security_events ADD CONSTRAINT fk_security_events_resolved_by 
--     FOREIGN KEY (resolved_by) REFERENCES users(id) ON DELETE SET NULL;

-- Create function to automatically clean old audit logs (optional)
CREATE OR REPLACE FUNCTION cleanup_old_audit_logs()
RETURNS void AS $$
BEGIN
    -- Delete audit logs older than 2 years
    DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL '2 years';
    
    -- Delete resolved security events older than 1 year
    DELETE FROM security_events 
    WHERE resolved = TRUE AND created_at < NOW() - INTERVAL '1 year';
    
    -- Delete system events older than 6 months (except critical ones)
    DELETE FROM system_events 
    WHERE created_at < NOW() - INTERVAL '6 months' 
    AND event_type NOT IN ('BACKUP_CREATED', 'MIGRATION_RUN');
END;
$$ LANGUAGE plpgsql;

-- Create a scheduled job to run cleanup (requires pg_cron extension)
-- Note: Uncomment this if you have pg_cron extension installed
-- SELECT cron.schedule('cleanup-audit-logs', '0 2 * * 0', 'SELECT cleanup_old_audit_logs();');

-- Insert some initial system events
INSERT INTO system_events (event_type, category, description, status, metadata) VALUES
('MIGRATION_RUN', 'MIGRATION', 'Created audit tables', 'SUCCESS', '{"migration": "20240126000001_create_audit_tables"}'),
('CONFIG_RELOADED', 'MAINTENANCE', 'Audit system initialized', 'SUCCESS', '{"component": "audit_system"}');

-- Create view for recent audit activity
CREATE OR REPLACE VIEW recent_audit_activity AS
SELECT 
    al.id,
    al.user_id,
    al.action,
    al.resource,
    al.resource_id,
    al.ip_address,
    al.created_at,
    'audit_log' as event_source
FROM audit_logs al
WHERE al.created_at > NOW() - INTERVAL '7 days'

UNION ALL

SELECT 
    se.id,
    se.user_id,
    se.event_type as action,
    'security' as resource,
    se.severity as resource_id,
    se.ip_address,
    se.created_at,
    'security_event' as event_source
FROM security_events se
WHERE se.created_at > NOW() - INTERVAL '7 days'

ORDER BY created_at DESC;

-- Create view for audit statistics
CREATE OR REPLACE VIEW audit_statistics AS
SELECT 
    'audit_logs' as table_name,
    COUNT(*) as total_records,
    COUNT(DISTINCT user_id) as unique_users,
    MIN(created_at) as oldest_record,
    MAX(created_at) as newest_record
FROM audit_logs

UNION ALL

SELECT 
    'security_events' as table_name,
    COUNT(*) as total_records,
    COUNT(DISTINCT user_id) as unique_users,
    MIN(created_at) as oldest_record,
    MAX(created_at) as newest_record
FROM security_events

UNION ALL

SELECT 
    'system_events' as table_name,
    COUNT(*) as total_records,
    0 as unique_users,
    MIN(created_at) as oldest_record,
    MAX(created_at) as newest_record
FROM system_events;
