ALTER TABLE users ADD COLUMN IF NOT EXISTS public_key TEXT DEFAULT NULL;
UPDATE users SET public_key = ssh_users.public_key FROM ssh_users WHERE users.uuid = ssh_users.user_uuid;
DROP TABLE IF EXISTS ssh_users;
