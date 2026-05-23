ALTER TABLE users 
DROP COLUMN IF EXISTS is_transactions_enabled,
DROP COLUMN IF EXISTS is_reporting_enabled,
DROP COLUMN IF EXISTS is_notification_enabled;
